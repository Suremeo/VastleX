package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type BossEvent struct{}

func (BossEvent) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.BossEvent).BossEntityUniqueID == int64(uid1) {
		pk.(*packet.BossEvent).BossEntityUniqueID = int64(uid2)
	} else if pk.(*packet.BossEvent).BossEntityUniqueID == int64(uid2) {
		pk.(*packet.BossEvent).BossEntityUniqueID = int64(uid1)
	}
	if pk.(*packet.BossEvent).PlayerUniqueID == int64(uid1) {
		pk.(*packet.BossEvent).PlayerUniqueID = int64(uid2)
	} else if pk.(*packet.BossEvent).PlayerUniqueID == int64(uid2) {
		pk.(*packet.BossEvent).PlayerUniqueID = int64(uid1)
	}
}
