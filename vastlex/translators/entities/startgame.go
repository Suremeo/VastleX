package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type StartGame struct {}

func (StartGame) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.StartGame).EntityUniqueID == int64(uid1) {
		pk.(*packet.StartGame).EntityUniqueID = int64(uid2)
	} else if pk.(*packet.StartGame).EntityUniqueID == int64(uid2) {
		pk.(*packet.StartGame).EntityUniqueID = int64(uid1)
	}
	if pk.(*packet.StartGame).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.StartGame).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.StartGame).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.StartGame).EntityRuntimeID = uint64(eid1)
	}
}