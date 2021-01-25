package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type StructureBlockUpdate struct{}

func (StructureBlockUpdate) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	pk.(*packet.StructureBlockUpdate).Settings.LastEditingPlayerUniqueID = translateUid(pk.(*packet.StructureBlockUpdate).Settings.LastEditingPlayerUniqueID, int64(uid1), int64(uid2))
}
