package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type BlockActorData struct {}

func (BlockActorData) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
}