package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type ShowCredits struct{}

func (ShowCredits) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.ShowCredits).PlayerRuntimeID == uint64(uid1) {
		pk.(*packet.ShowCredits).PlayerRuntimeID = uint64(uid2)
	} else if pk.(*packet.ShowCredits).PlayerRuntimeID == uint64(uid2) {
		pk.(*packet.ShowCredits).PlayerRuntimeID = uint64(uid1)
	}
}
