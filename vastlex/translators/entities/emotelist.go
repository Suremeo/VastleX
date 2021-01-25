package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type EmoteList struct{}

func (EmoteList) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.EmoteList).PlayerRuntimeID == uint64(uid1) {
		pk.(*packet.EmoteList).PlayerRuntimeID = uint64(uid2)
	} else if pk.(*packet.EmoteList).PlayerRuntimeID == uint64(uid2) {
		pk.(*packet.EmoteList).PlayerRuntimeID = uint64(uid1)
	}
}
