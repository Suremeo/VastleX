package packets

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type VastleXTransfer struct {
	Host        string
	Port        int32
	Message     string
	HideMessage bool
}

// ID ...
func (*VastleXTransfer) ID() uint32 {
	return IDVastleXTransfer
}

// Marshal ...
func (pk *VastleXTransfer) Marshal(w *protocol.Writer) {
	w.String(&pk.Host)
	w.Varint32(&pk.Port)
	w.String(&pk.Message)
	w.Bool(&pk.HideMessage)
}

// Unmarshal ...
func (pk *VastleXTransfer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Host)
	r.Varint32(&pk.Port)
	r.String(&pk.Message)
	r.Bool(&pk.HideMessage)
}
