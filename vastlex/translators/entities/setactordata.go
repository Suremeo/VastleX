package entities

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SetActorData struct{}

func (SetActorData) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.SetActorData).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.SetActorData).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.SetActorData).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.SetActorData).EntityRuntimeID = uint64(eid1)
	}
}
