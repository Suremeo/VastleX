package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// handleRemoveActor handles the RemoveActor packet sent by the remote minecraft.Connection.
func handleRemoveActor(dialer *Dialer, pak packet.Packet) {
	s := pak.(*packet.RemoveActor).EntityUniqueID
	if s == 1 {
		s = dialer.uniqueId
	}
	dialer.entities.Delete(s)
}
