package vastlex

import (
	"github.com/VastleLLC/VastleX/vastlex/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Player interface {
	Identity() login.IdentityData

	Send(info server.Info, config ...server.ConnectConfig) error
	Message(message string) error
	Kick(reason ...string)

	Server() server.Server

	WritePacket(packet packet.Packet) error
}
