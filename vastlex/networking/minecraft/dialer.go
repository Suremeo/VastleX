package minecraft

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/sandertv/go-raknet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login/jwt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
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

			}
			return
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
	headerData, err := jwt.HeaderFrom(pk.JWT)
	if err != nil {
		return fmt.Errorf("error reading ServerToClientHandshake JWT header: %v", err)
	}
	header := &jwt.Header{}
	if err := json.Unmarshal(headerData, header); err != nil {
		return fmt.Errorf("error parsing ServerToClientHandshake JWT header JSON: %v", err)
	}
	if !jwt.AllowedAlg(header.Algorithm) {
		return fmt.Errorf("ServerToClientHandshake JWT header had unexpected alg: expected %v, got %v", "ES384", header.Algorithm)
	}
	// First parse the public pubKey, so that we can use it to verify the entire JWT afterwards. The JWT is self-
	// signed by the server.
	pubKey := &ecdsa.PublicKey{}
	if err := jwt.ParsePublicKey(header.X5U, pubKey); err != nil {
		return fmt.Errorf("error parsing ServerToClientHandshake header x5u public pubKey: %v", err)
	}
	if _, err := jwt.Verify(pk.JWT, pubKey, false); err != nil {
		return fmt.Errorf("error verifying ServerToClientHandshake JWT: %v", err)
	}
	// We already know the JWT is valid as we verified it, so no need to error check.
	body, _ := jwt.Payload(pk.JWT)
	m := make(map[string]string)
	if err := json.Unmarshal(body, &m); err != nil {
		return fmt.Errorf("error parsing ServerToClientHandshake JWT payload JSON: %v", err)
	}
	b64Salt, ok := m["salt"]
	if !ok {
		return fmt.Errorf("ServerToClientHandshake JWT payload contained no 'salt'")
	}
	// Some (faulty) JWT implementations use padded base64, whereas it should be raw. We trim this off.
	b64Salt = strings.TrimRight(b64Salt, "=")
	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return fmt.Errorf("error base64 decoding ServerToClientHandshake salt: %v", err)
	}

	x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, connection.privateKey.D.Bytes())
	sharedSecret := x.Bytes()
	keyBytes := sha256.Sum256(append(salt, sharedSecret...))

	// Finally we enable encryption for the enc and dec using the secret pubKey bytes we produced.
	connection.encoder.EnableEncryption(keyBytes)
	connection.decoder.EnableEncryption(keyBytes)

	// We write a ClientToServerHandshake packet (which has no payload) as a response.
	return connection.WritePacket(&packet.ClientToServerHandshake{})
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
