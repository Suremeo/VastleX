package vastlex

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/vastlellc/vastlex/vastlex/server"
)

// Player is a player connected to the proxy.
type Player interface {
	// Identity returns the IdentityData of the player.
	Identity() login.IdentityData

	// Send transfers the player to a server.
	Send(info server.Info, config ...server.ConnectConfig) error

	// Message sends a chat message to a player.
	Message(message string) error

	// Kick kicks a player from the proxy, if no reason is provided a default reason is used.
	Kick(reason ...string)

	// Server returns the current server that a player is in, if they aren't in a server it will be nil.
	Server() server.Server

	// WritePacket writes a packet directly to the player.
	WritePacket(packet packet.Packet) error
}
