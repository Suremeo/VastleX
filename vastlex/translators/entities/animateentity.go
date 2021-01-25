package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type AnimateEntity struct{}

func (AnimateEntity) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	ppk := pk.(*packet.AnimateEntity)
	for index, runtime := range ppk.EntityRuntimeIDs {
		if runtime == uint64(eid1) {
			runtime = uint64(eid2)
		} else if runtime == uint64(eid2) {
			runtime = uint64(eid1)
		}
		ppk.EntityRuntimeIDs[index] = translateEid(runtime, eid1, eid2)
	}
}
