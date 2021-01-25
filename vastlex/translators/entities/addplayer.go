package entities

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddPlayer struct{}

func (AddPlayer) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.AddPlayer).EntityUniqueID == int64(uid1) {
		pk.(*packet.AddPlayer).EntityUniqueID = int64(uid2)
	} else if pk.(*packet.AddPlayer).EntityUniqueID == int64(uid2) {
		pk.(*packet.AddPlayer).EntityUniqueID = int64(uid1)
	}
	if pk.(*packet.AddPlayer).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.AddPlayer).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.AddPlayer).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.AddPlayer).EntityRuntimeID = uint64(eid1)
	}
	if pk.(*packet.AddPlayer).PlayerUniqueID == int64(uid1) {
		pk.(*packet.AddPlayer).PlayerUniqueID = int64(uid2)
	} else if pk.(*packet.AddPlayer).PlayerUniqueID == int64(uid2) {
		pk.(*packet.AddPlayer).PlayerUniqueID = int64(uid1)
	}
	ppk := pk.(*packet.AddPlayer)
	for index, link := range pk.(*packet.AddPlayer).EntityLinks {
		if link.RiderEntityUniqueID == int64(uid1) {
			link.RiderEntityUniqueID = int64(uid2)
		} else if link.RiderEntityUniqueID == int64(uid2) {
			link.RiderEntityUniqueID = int64(uid1)
		}
		if link.RiddenEntityUniqueID == int64(uid1) {
			link.RiddenEntityUniqueID = int64(uid2)
		} else if link.RiddenEntityUniqueID == int64(uid2) {
			link.RiddenEntityUniqueID = int64(uid1)
		}
		ppk.EntityLinks[index] = link
	}
}
