package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type TakeItemActor struct {}

func (TakeItemActor) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.TakeItemActor).TakerEntityRuntimeID == uint64(eid1) {
		pk.(*packet.TakeItemActor).TakerEntityRuntimeID = uint64(uid2)
	} else if pk.(*packet.TakeItemActor).TakerEntityRuntimeID == uint64(uid2) {
		pk.(*packet.TakeItemActor).TakerEntityRuntimeID = uint64(uid1)
	}
	if pk.(*packet.TakeItemActor).ItemEntityRuntimeID == uint64(uid1) {
		pk.(*packet.TakeItemActor).ItemEntityRuntimeID = uint64(uid2)
	} else if pk.(*packet.TakeItemActor).ItemEntityRuntimeID == uint64(uid2) {
		pk.(*packet.TakeItemActor).ItemEntityRuntimeID = uint64(uid1)
	}
}