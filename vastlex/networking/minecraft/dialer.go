package minecraft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Dialer is a connection to a remote minecraft server.
type Dialer interface {
	Conn
}

// Dial dials a minecraft server and returns a Dialer
func Dial(identity login.IdentityData, client login.ClientData, address string) (Dialer, error) {
	key, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	var netConn net.Conn

	dialer := raknet.Dialer{ErrorLog: log.New(ioutil.Discard, "", 0)}
	var pong []byte
	pong, err := dialer.Ping(address)
	if err != nil {
		err = fmt.Errorf("raknet ping: %w", err)
		return nil, err
	}
	netConn, err = dialer.Dial(addressWithPongPort(pong, address))
	if err != nil {
		err = fmt.Errorf("raknet: %w", err)
		return nil, err
	}
	conn := initializeConnection(netConn, key, connectionTypeDialer)
	conn.identityData = &identity
	conn.clientData = &client

	// Disable the batch packet limit so that the server can send packets as often as it wants to.
	conn.decoder.DisableBatchPacketLimit()

	request := login.EncodeOffline(*conn.identityData, *conn.clientData, key)
	go conn.handlePackets()

	conn.expect(packet.IDServerToClientHandshake, packet.IDPlayStatus)
	if err := conn.WritePacket(&packet.Login{ConnectionRequest: request, ClientProtocol: protocol.CurrentProtocol}); err != nil {
		return nil, err
	}
	select {
	case <-conn.closed:
		// The connection was closed before we even were fully 'connected', so we return an error.
		return nil, fmt.Errorf("connection timeout")
	case <-conn.OnEvent(&events.Login{}, func() {}):
		// We've connected successfully. We return the connection and no error.
		return conn, nil
	}
}

// handlePackets handles incoming packets from the Connection.
func (connection *Connection) handlePackets() {
	defer func() {
		_ = connection.Close()
	}()
	for {
		// We finally arrived at the packet decoding loop. We constantly decode packets that arrive
		// and push them to the Conn so that they may be processed.
		packets, err := connection.decoder.Decode()
		if err != nil {
			if !raknet.ErrConnectionClosed(err) {

			} else {
				return
			}
		}
		for _, data := range packets {
			loggedInBefore := connection.loggedIn.Load()
			if err := connection.receive(data); err != nil {
				println(err.Error())
				return
			}
			if !loggedInBefore && connection.loggedIn.Load() {
				connection.executeEvent(&events.Login{})
			}
		}
	}
}

// handleServerToClientHandshake handles the handleServerToClientHandshake packet.
func (connection *Connection) handleServerToClientHandshake(pk *packet.ServerToClientHandshake) error {
	tok, err := jwt.ParseSigned(string(pk.JWT))
	if err != nil {
		return fmt.Errorf("parse server token: %w", err)
	}
	//lint:ignore S1005 Double assignment is done explicitly to prevent panics.
	raw, _ := tok.Headers[0].ExtraHeaders["x5u"]
	kStr, _ := raw.(string)

	pub := new(ecdsa.PublicKey)
	if err := login.ParsePublicKey(kStr, pub); err != nil {
		return fmt.Errorf("parse server public key: %w", err)
	}

	var c saltClaims
	if err := tok.Claims(pub, &c); err != nil {
		return fmt.Errorf("verify claims: %w", err)
	}
	c.Salt = strings.TrimRight(c.Salt, "=")
	salt, err := base64.RawStdEncoding.DecodeString(c.Salt)
	if err != nil {
		return fmt.Errorf("error base64 decoding ServerToClientHandshake salt: %v", err)
	}

	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, connection.privateKey.D.Bytes())
	// Make sure to pad the shared secret up to 96 bytes.
	sharedSecret := append(bytes.Repeat([]byte{0}, 48-len(x.Bytes())), x.Bytes()...)

	keyBytes := sha256.Sum256(append(salt, sharedSecret...))

	// Finally we enable encryption for the enc and dec using the secret pubKey bytes we produced.
	connection.encoder.EnableEncryption(keyBytes)
	connection.decoder.EnableEncryption(keyBytes)

	// We write a ClientToServerHandshake packet (which has no payload) as a response.
	return connection.WritePacket(&packet.ClientToServerHandshake{})
}

// saltClaims holds the claims for the salt sent by the server in the ServerToClientHandshake packet.
type saltClaims struct {
	Salt string `json:"salt"`
}

var regex = regexp.MustCompile(`[^\\];`)

// addressWithPongPort parses the redirect IPv4 port from the pong and returns the address passed with the port
// found if present, or the original address if not.
func addressWithPongPort(pong []byte, address string) string {
	indices := regex.FindAllStringIndex(string(pong), -1)
	frag := make([]string, len(indices)+1)

	first := 0
	for i, index := range indices {
		frag[i] = string(pong[first : index[1]-1])
		first = index[1]
	}
	if len(frag) > 10 {
		portStr := frag[10]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return address
		}
		// Remove the port from the address.
		addressParts := strings.Split(address, ":")
		address = strings.Join(strings.Split(address, ":")[:len(addressParts)-1], ":")
		return address + ":" + strconv.Itoa(port)
	}
	return address
}
