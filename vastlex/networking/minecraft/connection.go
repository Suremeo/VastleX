package minecraft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

// A lot of this code came from https://github.com/Sandertv/gophertunnel/blob/master/minecraft/conn.go

// Conn is an interface version of minecraft.Connection.
type Conn interface {
	WritePacket(packet.Packet) error
	ReadPacket() (packet.Packet, error)
	Close() error

	// OnEvent only executes once and returns when the event happens.
	OnEvent(events.Event, interface{}) chan struct{}

	// HandleEvent is called whenever an event happens.
	HandleEvent(events.Event, interface{})
}

// Connection represents all types of raknet connections.
type Connection struct {
	encoder        *packet.Encoder
	decoder        *packet.Decoder
	connectionType byte
	net            net.Conn
	salt           []byte
	privateKey     *ecdsa.PrivateKey
	close          sync.Once
	closed         chan struct{}
	writeBuffer    *bytes.Buffer
	bufferedWrites [][]byte
	writeMutex     sync.Mutex

	identityData *login.IdentityData
	clientData   *login.ClientData

	eventMutex sync.Mutex
	events     map[int]map[int]map[int]interface{}
	packets    chan *packetData

	deferredPackets []*packetData
	deferredPacketsMutex sync.Mutex
	expected atomic.Value

	loggedIn atomic.Bool
}

const (
	connectionTypePlayer = byte(iota)
	connectionTypeDialer
)

// initializeConnection initializes a net.Conn into a minecraft.Connection
func initializeConnection(conn net.Conn, key *ecdsa.PrivateKey, connectionType byte) *Connection {
	connection := &Connection{
		encoder:        packet.NewEncoder(conn),
		decoder:        packet.NewDecoder(conn),
		connectionType: connectionType,
		net:            conn,
		salt:           make([]byte, 16),
		packets:        make(chan *packetData, 256),
		privateKey:     key,
		writeBuffer:    bytes.NewBuffer(make([]byte, 0, 4096)),
		closed:         make(chan struct{}),
		events:         make(map[int]map[int]map[int]interface{}),
	}
	connection.expected.Store([]uint32{})
	_, _ = rand.Read(connection.salt)
	go func() {
		ticker := time.NewTicker(time.Second / 20)
		defer ticker.Stop()
		for range ticker.C {
			if err := connection.Flush(); err != nil {
				_ = connection.Close()
				return
			}
		}
	}()
	return connection
}

// Close closes the Connection.
func (connection *Connection) Close() error {
	var err error
	connection.close.Do(func() {
		close(connection.closed)
		err = connection.net.Close()
		connection.executeEvent(&events.Close{})
	})
	return err
}

// WritePacket writes a packet to the Connection.
func (connection *Connection) WritePacket(pk packet.Packet) error {
	select {
	case <-connection.closed:
		return fmt.Errorf("connection closed")
	default:
	}
	connection.writeMutex.Lock()
	defer connection.writeMutex.Unlock()
	header := &packet.Header{PacketID: pk.ID()}
	_ = header.Write(connection.writeBuffer)

	pk.Marshal(protocol.NewWriter(connection.writeBuffer))
	connection.bufferedWrites = append(connection.bufferedWrites, append([]byte(nil), connection.writeBuffer.Bytes()...))
	connection.writeBuffer.Reset()
	return nil
}

// ReadPacket reads a packet from the Connection.
func (connection *Connection) ReadPacket() (packet.Packet, error) {
	if data, ok := connection.takePushedBackPacket(); ok {
		pk, err := data.decode()
		if err != nil {
			return connection.ReadPacket()
		}
		return pk, nil
	}
	select {
	case <-connection.closed:
		return nil, fmt.Errorf("error reading packet: connection closed")
	case data := <-connection.packets:
		pk, err := data.decode()
		if err != nil {
			return connection.ReadPacket()
		}
		return pk, nil
	}
}

// Flush flushes the packets currently buffered by the connections to the underlying net.Conn, so that they are directly sent.
func (connection *Connection) Flush() error {
	select {
	case <-connection.closed:
		return fmt.Errorf("connection closed")
	default:
	}
	connection.writeMutex.Lock()
	defer connection.writeMutex.Unlock()

	if len(connection.bufferedWrites) > 0 {
		if err := connection.encoder.Encode(connection.bufferedWrites); err != nil {
			return fmt.Errorf("error encoding packet batch: %v", err)
		}
		// Reset the send slice so that we don't accidentally send the same packets.
		connection.bufferedWrites = connection.bufferedWrites[:0]
	}
	return nil
}

// OnEvent handles an event once and returns a channel of when the event is called.
func (connection *Connection) OnEvent(event events.Event, function interface{}) chan struct{} {
	connection.eventMutex.Lock()
	defer connection.eventMutex.Unlock()
	if connection.events[event.ID()] == nil {
		connection.events[event.ID()] = make(map[int]map[int]interface{})
	}
	this := len(connection.events[event.ID()])
	connection.events[event.ID()][this] = map[int]interface{}{
		1: &struct {
			Function interface{}
			Channel chan struct{}
		}{
			Function: function,
			Channel: make(chan struct{}),
		},
	}
	return connection.events[event.ID()][this][1].(*struct {
		Function interface{}
		Channel chan struct{}
	}).Channel
}

// HandleEvent registers a handler for an event.
func (connection *Connection) HandleEvent(event events.Event, function interface{}) {
	connection.eventMutex.Lock()
	defer connection.eventMutex.Unlock()
	if connection.events[event.ID()] == nil {
		connection.events[event.ID()] = make(map[int]map[int]interface{})
	}
	connection.events[event.ID()][len(connection.events[event.ID()])] = map[int]interface{}{
		0: function,
	}
}

// executeEvent executes an event so that it is handled by all the event handlers.
func (connection *Connection) executeEvent(event events.Event, args ...interface{}) {
	connection.eventMutex.Lock()
	defer connection.eventMutex.Unlock()
	if connection.events[event.ID()] != nil {
		for index, handler := range connection.events[event.ID()] {
			if handler[1] != nil { // 0 is an OnEvent handler.
				handle := handler[1].(*struct {
					Function interface{}
					Channel chan struct{}
				})
				close(handle.Channel)
				if handle.Function != nil {
					go event.Handle(handle.Function, args)
				}
				delete(connection.events[event.ID()], index)
			} else if handler[0] != nil { // 1 is a HandleEvent handler.
				go event.Handle(handler[0], args)
			}
		}
	}
}


// takePushedBackPacket locks the pushed back packets lock and takes the next packet from the list of
// pushed back packets. If none was found, it returns false, and if one was found, the data and true is
// returned.
func (connection *Connection) takePushedBackPacket() (*packetData, bool) {
	connection.deferredPacketsMutex.Lock()
	defer connection.deferredPacketsMutex.Unlock()

	if len(connection.deferredPackets) == 0 {
		return nil, false
	}
	data := connection.deferredPackets[0]
	connection.deferredPackets = connection.deferredPackets[1:]
	return data, true
}

// receive receives an incoming serialised packet from the underlying connection. If the connection is not yet
// logged in, the packet is immediately handled.
func (connection *Connection) receive(data []byte) error {
	pkData, err := parseData(data)
	if err != nil {
		return err
	}
	if pkData.h.PacketID == packet.IDDisconnect {
		// We always handle disconnect packets and close the connection if one comes in.
		_ = connection.Close()
		return nil
	}
	if connection.loggedIn.Load() {
		select {
		case <-connection.closed:
		case connection.packets <- pkData:
		}
		return nil
	}
	return connection.handle(pkData)
}

// handle tries to handle the incoming packetData.
func (connection *Connection) handle(pkData *packetData) error {
	for _, id := range connection.expected.Load().([]uint32) {
		if id == pkData.h.PacketID {
			// If the packet was expected, so we handle it right now.
			pk, err := pkData.decode()
			if err != nil {
				return err
			}
			return connection.handlePacket(pk)
		}
	}
	// This is not the packet we expected next in the login sequence. We push it back so that it may
	// be handled by the user.
	connection.deferredPackets = append(connection.deferredPackets, pkData)
	return nil
}

// handlePacket handles packets required for the login process.
func (connection *Connection) handlePacket(pk packet.Packet) error {
	defer func() {
		_ = connection.Flush()
	}()
	switch pk := pk.(type) {
	case *packet.Login:
		return connection.handleLogin(pk)
	case *packet.ClientToServerHandshake:
		return connection.handleClientToServerHandshake()
	case *packet.ResourcePackClientResponse:
		return connection.handleResourcePackClientResponse(pk)

	// Internal packets destined for the client.
	case *packet.ServerToClientHandshake:
		return connection.handleServerToClientHandshake(pk)
	case *packet.PlayStatus:
		return connection.handlePlayStatus(pk)
	case *packet.ResourcePacksInfo:
		return connection.handleResourcePacksInfo(pk)
	case *packet.ResourcePackStack:
		return connection.handleResourcePackStack(pk)
	}
	return nil
}

// expect sets the packet IDs that are next expected to arrive.
func (connection *Connection) expect(packetIDs ...uint32) {
	connection.expected.Store(packetIDs)
}

// handleClientToServerHandshake handles an incoming ClientToServerHandshake packet.
func (connection *Connection) handleClientToServerHandshake() error {
	if err := connection.WritePacket(&packet.NetworkSettings{CompressionThreshold: 512}); err != nil {
		return fmt.Errorf("error sending network settings: %v", err)
	}
	if err := connection.WritePacket(&packet.PlayStatus{Status: packet.PlayStatusLoginSuccess}); err != nil {
		return fmt.Errorf("error sending play status login success: %v", err)
	}
	if err := connection.WritePacket(&packet.ResourcePacksInfo{}); err != nil {
		return fmt.Errorf("error sending resource pack packet: %v", err)
	}
	connection.expect(packet.IDResourcePackClientResponse)
	return nil
}

// handleResourcePacksInfo handles a ResourcePacksInfo packet sent by the server. The client responds by
// sending the packs it needs downloaded.
func (connection *Connection) handleResourcePacksInfo(pk *packet.ResourcePacksInfo) error {
	return connection.WritePacket(&packet.ResourcePackClientResponse{Response: packet.PackResponseAllPacksDownloaded})
}

// handleResourcePackStack handles a ResourcePackStack packet sent by the server. The stack defines the order
// that resource packs are applied in.
func (connection *Connection) handleResourcePackStack(pk *packet.ResourcePackStack) error {
	connection.loggedIn.Store(true)
	connection.executeEvent(&events.Login{})
	return connection.WritePacket(&packet.ResourcePackClientResponse{Response: packet.PackResponseCompleted})
}

// handleResourcePackClientResponse handles the clients response to the ResourcePackStack and ResourcePacksInfo packets.
func (connection *Connection) handleResourcePackClientResponse(pk *packet.ResourcePackClientResponse) error {
	switch pk.Response {
	case packet.PackResponseRefused:
		return connection.Close()
	case packet.PackResponseAllPacksDownloaded:
		pk := &packet.ResourcePackStack{
			TexturePackRequired: false,
			BaseGameVersion:     protocol.CurrentVersion,
		}
		if err := connection.WritePacket(pk); err != nil {
			return fmt.Errorf("error writing resource pack stack packet: %v", err)
		}
		connection.expect(packet.IDResourcePackClientResponse)
	case packet.PackResponseCompleted:
		connection.loggedIn.Store(true)
	default:
		return fmt.Errorf("unknown resource pack client response: %v", pk.Response)
	}
	return nil
}

// handlePlayStatus handles an incoming PlayStatus packet. It reacts differently depending on the status
// found in the packet.
func (connection *Connection) handlePlayStatus(pk *packet.PlayStatus) error {
	switch pk.Status {
	case packet.PlayStatusLoginSuccess:
		connection.expect(packet.IDResourcePacksInfo, packet.IDResourcePackStack)
		if err := connection.WritePacket(&packet.ClientCacheStatus{Enabled: false}); err != nil {
			return fmt.Errorf("error sending client cache status: %v", err)
		}
		return connection.Flush()
	case packet.PlayStatusLoginFailedClient:
		_ = connection.Close()
		return fmt.Errorf("client outdated")
	case packet.PlayStatusLoginFailedServer:
		_ = connection.Close()
		return fmt.Errorf("server outdated")
	case packet.PlayStatusPlayerSpawn:
	case packet.PlayStatusLoginFailedInvalidTenant:
		_ = connection.Close()
		return fmt.Errorf("invalid edu edition game owner")
	case packet.PlayStatusLoginFailedVanillaEdu:
		_ = connection.Close()
		return fmt.Errorf("cannot join an edu edition game on vanilla")
	case packet.PlayStatusLoginFailedEduVanilla:
		_ = connection.Close()
		return fmt.Errorf("cannot join a vanilla game on edu edition")
	case packet.PlayStatusLoginFailedServerFull:
		_ = connection.Close()
		return fmt.Errorf("server full")
	default:
		return fmt.Errorf("unknown play status in PlayStatus packet %v", pk.Status)
	}
	return nil
}

// RemoteAddr returns the remote address of the underlying connection.
func (connection *Connection) RemoteAddr() net.Addr {
	return connection.net.RemoteAddr()
}