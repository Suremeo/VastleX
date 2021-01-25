package minecraft

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"net"
)

// Player represents a Player connected to a Listener.
type Player interface {
	Conn
	Identity() login.IdentityData
	Client() login.ClientData
	RemoteAddr() net.Addr
}

// Identity returns the login.IdentityData of the Player.
func (connection *Connection) Identity() login.IdentityData {
	return *connection.identityData
}

// Client returns the login.ClientData of the Player.
func (connection *Connection) Client() login.ClientData {
	return *connection.clientData
}

// enableEncryption enables encryption on the server side over the connection. It sends an unencrypted
// handshake packet to the client and enables encryption after that.
func (connection *Connection) enableEncryption(clientPublicKey *ecdsa.PublicKey) error {
	connection.expect(packet.IDClientToServerHandshake)
	signer, _ := jose.NewSigner(jose.SigningKey{Key: connection.privateKey, Algorithm: jose.ES384}, &jose.SignerOptions{
		ExtraHeaders: map[jose.HeaderKey]interface{}{"x5u": login.MarshalPublicKey(&connection.privateKey.PublicKey)},
	})
	// We produce an encoded JWT using the header and payload above, then we send the JWT in a ServerToClient-
	// Handshake packet so that the client can initialise encryption.
	serverJWT, err := jwt.Signed(signer).Claims(saltClaims{Salt: base64.RawStdEncoding.EncodeToString(connection.salt)}).CompactSerialize()
	if err != nil {
		return fmt.Errorf("compact serialise server JWT: %w", err)
	}
	if err := connection.WritePacket(&packet.ServerToClientHandshake{JWT: []byte(serverJWT)}); err != nil {
		return fmt.Errorf("error sending ServerToClientHandshake packet: %v", err)
	}
	// Flush immediately as we'll enable encryption after this.
	_ = connection.Flush()

	// We first compute the shared secret.
	x, _ := clientPublicKey.Curve.ScalarMult(clientPublicKey.X, clientPublicKey.Y, connection.privateKey.D.Bytes())

	sharedSecret := append(bytes.Repeat([]byte{0}, 48-len(x.Bytes())), x.Bytes()...)

	keyBytes := sha256.Sum256(append(connection.salt, sharedSecret...))

	// Finally we enable encryption for the encoder and decoder using the secret key bytes we produced.
	connection.encoder.EnableEncryption(keyBytes)
	connection.decoder.EnableEncryption(keyBytes)

	return nil
}
