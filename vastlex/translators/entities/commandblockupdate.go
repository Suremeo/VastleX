package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type CommandBlockUpdate struct{}

func (CommandBlockUpdate) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.CommandBlockUpdate).MinecartEntityRuntimeID == uint64(uid1) {
		pk.(*packet.CommandBlockUpdate).MinecartEntityRuntimeID = uint64(uid2)
	} else if pk.(*packet.CommandBlockUpdate).MinecartEntityRuntimeID == uint64(uid2) {
		pk.(*packet.CommandBlockUpdate).MinecartEntityRuntimeID = uint64(uid1)
	}
}
