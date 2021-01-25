package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type PlayerList struct{}

func (PlayerList) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	ppk := pk.(*packet.PlayerList)
	for index, item := range ppk.Entries {
		if item.EntityUniqueID == int64(uid1) {
			item.EntityUniqueID = int64(uid2)
		} else if item.EntityUniqueID == int64(uid2) {
			item.EntityUniqueID = int64(uid1)
		}
		ppk.Entries[index] = item
	}
}
