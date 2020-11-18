package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type AdventureSettings struct{}

func (AdventureSettings) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.AdventureSettings).PlayerUniqueID == int64(uid1) {
		pk.(*packet.AdventureSettings).PlayerUniqueID = int64(uid2)
	} else if pk.(*packet.AdventureSettings).PlayerUniqueID == int64(uid2) {
		pk.(*packet.AdventureSettings).PlayerUniqueID = int64(uid1)
	}
}
