package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type MobArmourEquipment struct {}

func (MobArmourEquipment) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.MobArmourEquipment).EntityRuntimeID == uint64(eid1) {
		pk.(*packet.MobArmourEquipment).EntityRuntimeID = uint64(eid2)
	} else if pk.(*packet.MobArmourEquipment).EntityRuntimeID == uint64(eid2) {
		pk.(*packet.MobArmourEquipment).EntityRuntimeID = uint64(eid1)
	}
}