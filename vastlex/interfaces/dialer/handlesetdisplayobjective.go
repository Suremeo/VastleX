package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func handleSetDisplayObjective(dialer *Dialer, pak packet.Packet) {
	pk := pak.(*packet.SetDisplayObjective)
	dialer.scoreboards.Store(pk.ObjectiveName, struct{}{})
}
