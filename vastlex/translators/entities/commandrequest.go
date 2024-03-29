package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type CommandRequest struct{}

func (CommandRequest) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	pk.(*packet.CommandRequest).CommandOrigin.PlayerUniqueID = translateUid(pk.(*packet.CommandRequest).CommandOrigin.PlayerUniqueID, uid1, uid2)
}
