package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// handleAddActor handles the addactor packet sent by the remote minecraft.Connection.
func handleAddActor(dialer *Dialer, pak packet.Packet) {
	s := pak.(*packet.AddActor).EntityUniqueID
	if s == 1 {
		s = dialer.uniqueId
	}
	dialer.entities.Store(s, nil)
}