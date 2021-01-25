package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type SetScore struct{}

func (SetScore) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	ppk := pk.(*packet.SetScore)
	for i := range ppk.Entries {
		ppk.Entries[i].EntityUniqueID = translateUid(ppk.Entries[i].EntityUniqueID, uid1, uid2)
	}
}
