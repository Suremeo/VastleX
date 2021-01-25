package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func handleRemoveObjective(dialer *Dialer, pak packet.Packet) {
	pk := pak.(*packet.RemoveObjective)
	dialer.scoreboards.Delete(pk.ObjectiveName)
}
