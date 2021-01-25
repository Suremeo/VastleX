package entities

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type Camera struct{}

func (Camera) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	if pk.(*packet.Camera).CameraEntityUniqueID == int64(uid1) {
		pk.(*packet.Camera).CameraEntityUniqueID = int64(uid2)
	} else if pk.(*packet.Camera).CameraEntityUniqueID == int64(uid2) {
		pk.(*packet.Camera).CameraEntityUniqueID = int64(uid1)
	}
	if pk.(*packet.Camera).TargetPlayerUniqueID == int64(uid1) {
		pk.(*packet.Camera).TargetPlayerUniqueID = int64(uid2)
	} else if pk.(*packet.Camera).TargetPlayerUniqueID == int64(uid2) {
		pk.(*packet.Camera).TargetPlayerUniqueID = int64(uid1)
	}
}
