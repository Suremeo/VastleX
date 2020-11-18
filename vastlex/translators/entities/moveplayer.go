package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type MovePlayer struct{}

func (MovePlayer) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.MovePlayer).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.MovePlayer).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.MovePlayer).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.MovePlayer).EntityRuntimeID = uint64(eid1)
	}
	if pk.(*packet.MovePlayer).RiddenEntityRuntimeID == uint64(eid1) {
		pk.(*packet.MovePlayer).RiddenEntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.MovePlayer).RiddenEntityRuntimeID == uint64(eid2) {
		pk.(*packet.MovePlayer).RiddenEntityRuntimeID = uint64(eid1)
	}
}
