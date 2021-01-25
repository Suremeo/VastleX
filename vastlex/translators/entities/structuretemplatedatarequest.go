package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type StructureTemplateDataRequest struct{}

func (StructureTemplateDataRequest) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	pk.(*packet.StructureTemplateDataRequest).Settings.LastEditingPlayerUniqueID = translateUid(pk.(*packet.StructureTemplateDataRequest).Settings.LastEditingPlayerUniqueID, uid1, uid2)
}
