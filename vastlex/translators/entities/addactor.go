package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type AddActor struct{}

func (AddActor) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.AddActor).EntityUniqueID == int64(uid1) {
		pk.(*packet.AddActor).EntityUniqueID = int64(uid2)
	} else if pk.(*packet.AddActor).EntityUniqueID == int64(uid2) {
		pk.(*packet.AddActor).EntityUniqueID = int64(uid1)
	}
	if pk.(*packet.AddActor).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.AddActor).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.AddActor).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.AddActor).EntityRuntimeID = uint64(eid1)
	}
}
