package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Animate struct{}

func (Animate) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.Animate).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.Animate).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.Animate).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.Animate).EntityRuntimeID = uint64(eid1)
	}
}
