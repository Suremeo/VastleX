package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type LevelEvent struct{}

func (LevelEvent) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
}
