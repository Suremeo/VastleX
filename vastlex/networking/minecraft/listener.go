package minecraft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/config"
	log2 "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/ddos"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.uber.org/atomic"
	"log"
	"net"
	"sync"
)

// Listener represents an instance of the proxy listening for connecting players.
type Listener interface {
	SetMotd(string)
	Motd() string
	Accept() Player
	SetPlayerCount(int)
	PlayerCount() int
}

var _ Listener = &listener{}

// listener is an internal structure for the default listener.
type listener struct {
	net      *raknet.Listener
	mutex    sync.Mutex
	count    *atomic.Int32
	incoming chan *Connection
	key      *ecdsa.PrivateKey
}

// Listen listens for new connections on the address specified in the configuration.
func Listen() (_ Listener, err error) {
	listener := &listener{
		count:    atomic.NewInt32(0),
		incoming: make(chan *Connection),
	}
	l := &raknet.ListenConfig{ErrorLog: log.New(&ddos.W{}, "", 0)}
	listener.key, _ = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	listener.net, err = l.Listen(fmt.Sprintf("%v:%v", config.Config.Listener.Host, config.Config.Listener.Port))
	if err == nil {
		listener.updatePongData()
		go func() {
			for {
				c, err := listener.net.Accept()
				if err != nil {
					panic(err) // We panic because listener should never be closed and if it is something is really messed up.
				}
				if listener.count.Load() > int32(config.Config.Minecraft.MaxPlayers) && config.Config.Minecraft.MaxPlayers != 0 {
					// The server is full.
					conn := &Connection{net: c, encoder: packet.NewEncoder(c), closed: make(chan struct{}), writeBuffer: bytes.NewBuffer(make([]byte, 0, 4096))}
					_ = conn.WritePacket(&packet.PlayStatus{Status: packet.PlayStatusLoginFailedServerFull})
					_ = conn.Flush()
					conn = nil
					_ = c.Close()
				} else {
					go func() {
						listener.count.Add(1)
						listener.updatePongData()
						listener.handleConnection(initializeConnection(c, listener.key, connectionTypePlayer))
					}()
				}
			}
		}()
	}
	return listener, err
}

// updatePongData updates the pong data of the listener using the current only players, maximum players and motd for the proxy.
func (l *listener) updatePongData() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	current := l.count.Load()
	max := config.Config.Minecraft.MaxPlayers
	if max == 0 {
		max = int(current + 1)
	}

	l.net.PongData([]byte(fmt.Sprintf("MCPE;%v;%v;%v;%v;%v;%v;Minecraft Server;%v;%v;%v;%v;",
		l.Motd(), protocol.CurrentProtocol, protocol.CurrentVersion, current, max, l.net.ID(),
		"Creative", 1, l.net.Addr().(*net.UDPAddr).Port, l.net.Addr().(*net.UDPAddr).Port,
	)))
}

// SetMotd updates the MOTD for the listener.
func (l *listener) SetMotd(s string) {
	config.Config.Minecraft.Motd = s
}

// Motd returns the MOTD for the listener.
func (l *listener) Motd() string {
	return text.Colourf(config.Config.Minecraft.Motd)
}

// Accept accepts a new connection from the listener.
func (l *listener) Accept() Player {
	return <-l.incoming
}

// SetPlayerCount manually sets the player count of the proxy.
func (l *listener) SetPlayerCount(count int) {
	l.count.Store(int32(count))
}

// PlayerCount returns the player count of the proxy.
func (l *listener) PlayerCount() int {
	return int(l.count.Load())
}

// handleConnection handles an incoming connection of the Listener. It will first attempt to get the connection to
// log in, after which it will expose packets received to the user.
func (l *listener) handleConnection(conn *Connection) {
	defer func() {
		_ = conn.Close()
		l.count.Add(-1)
		l.updatePongData()
	}()
	conn.expect(packet.IDLogin)
	for {
		// We finally arrived at the packet decoding loop. We constantly decode packets that arrive
		// and push them to the Conn so that they may be processed.
		packets, err := conn.decoder.Decode()
		if err != nil {
			return
		}
		for _, data := range packets {
			loggedInBefore := conn.loggedIn.Load()
			if err := conn.receive(data); err != nil {
				log2.DefaultLogger.Error(err)
				return
			}
			if !loggedInBefore && conn.loggedIn.Load() {
				l.incoming <- conn
			}
		}
	}
}

// handleLogin handles a login packet from a player.
func (connection *Connection) handleLogin(pk *packet.Login) error {
	identityData, clientData, authResult, err := login.Parse(pk.ConnectionRequest)
	if err != nil {
		return fmt.Errorf("parse login request: %w", err)
	}
	connection.identityData = &identityData
	connection.clientData = &clientData
	if !authResult.XBOXLiveAuthenticated && config.Config.Minecraft.Auth {
		_ = connection.WritePacket(&packet.Disconnect{Message: text.Colourf("<red>You must be logged in with XBOX Live to join.</red>")})
		_ = connection.Close()
		return fmt.Errorf("connection %v was not authenticated to XBOX Live", connection.RemoteAddr())
	}
	if pk.ClientProtocol != protocol.CurrentProtocol {
		status := packet.PlayStatusLoginFailedClient
		if pk.ClientProtocol > protocol.CurrentProtocol {
			status = packet.PlayStatusLoginFailedServer
		}
		_ = connection.WritePacket(&packet.PlayStatus{Status: status})
		return fmt.Errorf("%v connected with an incompatible protocol: expected protocol = %v, client protocol = %v", connection.identityData.DisplayName, protocol.CurrentProtocol, pk.ClientProtocol)
	}
	connection.authResult.Store(true)
	if err := connection.enableEncryption(authResult.PublicKey); err != nil {
		return fmt.Errorf("error enabling encryption: %v", err)
	}
	return nil
}
