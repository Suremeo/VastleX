package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type SetActorLink struct{}

func (SetActorLink) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.SetActorLink).EntityLink.RiddenEntityUniqueID == int64(eid1) {
		pk.(*packet.SetActorLink).EntityLink.RiddenEntityUniqueID = int64(eid2)
	} else if pk.(*packet.SetActorLink).EntityLink.RiddenEntityUniqueID == int64(eid2) {
		pk.(*packet.SetActorLink).EntityLink.RiddenEntityUniqueID = int64(eid1)
	}
	if pk.(*packet.SetActorLink).EntityLink.RiderEntityUniqueID == int64(eid1) {
		pk.(*packet.SetActorLink).EntityLink.RiderEntityUniqueID = int64(eid2)
	} else if pk.(*packet.SetActorLink).EntityLink.RiderEntityUniqueID == int64(eid2) {
		pk.(*packet.SetActorLink).EntityLink.RiderEntityUniqueID = int64(eid1)
	}
}
