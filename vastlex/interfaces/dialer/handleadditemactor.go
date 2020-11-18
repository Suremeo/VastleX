package dialer

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// handleAddItemActor handles the AddItemActor packet sent by the remote minecraft.Connection.
func handleAddItemActor(dialer *Dialer, pak packet.Packet) {
	s := pak.(*packet.AddItemActor).EntityUniqueID
	if s == 1 {
		s = dialer.uniqueId
	}
	dialer.entities.Store(s, nil)
}