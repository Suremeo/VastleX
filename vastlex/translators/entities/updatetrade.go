package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type UpdateTrade struct{}

func (UpdateTrade) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.UpdateTrade).VillagerUniqueID == int64(uid1) {
		pk.(*packet.UpdateTrade).VillagerUniqueID = int64(uid2)
	} else if pk.(*packet.UpdateTrade).VillagerUniqueID == int64(uid2) {
		pk.(*packet.UpdateTrade).VillagerUniqueID = int64(uid1)
	}
	if pk.(*packet.UpdateTrade).EntityUniqueID == int64(uid1) {
		pk.(*packet.UpdateTrade).EntityUniqueID = int64(uid2)
	} else if pk.(*packet.UpdateTrade).EntityUniqueID == int64(uid2) {
		pk.(*packet.UpdateTrade).EntityUniqueID = int64(uid1)
	}
}
