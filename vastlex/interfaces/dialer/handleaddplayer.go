package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func handleAddPlayer(dialer *Dialer, pak packet.Packet) {
	s := pak.(*packet.AddPlayer).EntityUniqueID
	if s == 1 {
		s = dialer.uniqueId
	}
	dialer.entities.Store(s, nil)
}
