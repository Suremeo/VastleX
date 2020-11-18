package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Respawn struct{}

func (Respawn) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.Respawn).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.Respawn).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.Respawn).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.Respawn).EntityRuntimeID = uint64(eid1)
	}
}
