package minecraft

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login/jwt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
)

// Player represents a Player connected to a Listener.
type Player interface {
	Conn
	Identity() login.IdentityData
	ClientData() login.ClientData
	RemoteAddr() net.Addr
}

// Identity returns the login.IdentityData of the Player.
func (connection *Connection) Identity() login.IdentityData {
	return *connection.identityData
}

// ClientData returns the login.ClientData of the Player.
func (connection *Connection) ClientData() login.ClientData {
	return *connection.clientData
}

// enableEncryption enables encryption on the server side over the connection. It sends an unencrypted
// handshake packet to the client and enables encryption after that.
func (connection *Connection) enableEncryption(clientPublicKey *ecdsa.PublicKey) error {
	connection.expect(packet.IDClientToServerHandshake)
	pubKey := jwt.MarshalPublicKey(&connection.privateKey.PublicKey)
	header := jwt.Header{
		Algorithm: "ES384",
		X5U:       pubKey,
	}
	payload := map[string]interface{}{
		"salt": base64.StdEncoding.EncodeToString(connection.salt),
	}

	// We produce an encoded JWT using the header and payload above, then we send the JWT in a ServerToClient-
	// Handshake packet so that the client can initialise encryption.
	serverJWT, err := jwt.New(header, payload, connection.privateKey)
	if err != nil {
		return fmt.Errorf("error creating encoded JWT: %v", err)
	}
	if err := connection.WritePacket(&packet.ServerToClientHandshake{JWT: serverJWT}); err != nil {
		return fmt.Errorf("error sending ServerToClientHandshake packet: %v", err)
	}
	// Flush immediately as we'll enable encryption after this.
	_ = connection.Flush()

	// We first compute the shared secret.
	x, _ := clientPublicKey.Curve.ScalarMult(clientPublicKey.X, clientPublicKey.Y, connection.privateKey.D.Bytes())
	sharedSecret := x.Bytes()
	keyBytes := sha256.Sum256(append(connection.salt, sharedSecret...))

	// Finally we enable encryption for the encoder and decoder using the secret key bytes we produced.
	connection.encoder.EnableEncryption(keyBytes)
	connection.decoder.EnableEncryption(keyBytes)

	return nil
}