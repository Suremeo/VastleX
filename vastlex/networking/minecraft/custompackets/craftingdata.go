package custompackets

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type CraftingData struct {
	Data []byte
}

// ID ...
func (*CraftingData) ID() uint32 {
	return packet.IDCraftingData
}

// Marshal ...
func (pk *CraftingData) Marshal(w *protocol.Writer) {
	w.Bytes(&pk.Data)
}

// Unmarshal ...
func (pk *CraftingData) Unmarshal(r *protocol.Reader) {
	r.Bytes(&pk.Data)
}
