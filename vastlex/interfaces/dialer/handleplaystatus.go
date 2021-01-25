package dialer

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func handlePlayStatus(dialer *Dialer, pak packet.Packet) {
	pk := pak.(*packet.PlayStatus)
	if pk.Status == packet.PlayStatusPlayerSpawn {
		_ = dialer.WritePacket(&packet.SetLocalPlayerAsInitialised{EntityRuntimeID: uint64(dialer.entityId)})
	}
}
