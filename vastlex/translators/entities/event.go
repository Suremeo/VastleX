package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Event struct {}

func (Event) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.Event).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.Event).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.Event).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.Event).EntityRuntimeID = uint64(eid1)
	}
}