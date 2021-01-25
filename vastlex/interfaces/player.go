package interfaces

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// InternalPlayer represents the internal version of a player connected to the proxy.
type InternalPlayer interface {
	minecraft.Player
	Send(ip string, port int) error
	Config() *config.Player
	SetConfig(config *config.Player)
	Dialer() Dialer
	SetDialer(dialer Dialer)
	State() State
	SetState(state State)
	Kick(...string)
	Message(string) error
	KickOrFallback(string)
	ChunkRadius() int32
}

// State represents the state of a player.
type State int

const (
	StateWaitingForFirstServer = State(iota)
	StateWaitingForNewServer
	StateConnected
	StateDisconnected
)

// Player is a player connected to the proxy.
type Player interface {
	// Identity returns the IdentityData of the player.
	Identity() login.IdentityData

	// Client returns the Client of the player
	Client() login.ClientData

	// Send transfers the player to a server.
	Send(ip string, port int) error

	// Message sends a chat message to a player.
	Message(message string) error

	// Kick kicks a player from the proxy, if no reason is provided a default reason is used.
	Kick(reason ...string)

	// WritePacket writes a packet directly to the player.
	WritePacket(packet packet.Packet) error

	// State returns the state of the player.
	State() State
}
