package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type AddPlayer struct {}

func (AddPlayer) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.AddPlayer).EntityUniqueID == int64(uid1) {
		pk.(*packet.AddPlayer).EntityUniqueID = int64(uid2)
	} else if pk.(*packet.AddPlayer).EntityUniqueID == int64(uid2) {
		pk.(*packet.AddPlayer).EntityUniqueID = int64(uid1)
	}
	if pk.(*packet.AddPlayer).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.AddPlayer).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.AddPlayer).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.AddPlayer).EntityRuntimeID = uint64(eid1)
	}
}