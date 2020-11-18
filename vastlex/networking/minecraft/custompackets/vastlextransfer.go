package custompackets

import (
	"fmt"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// VastleXTransfer is the custom packet used for transferring a player.
type VastleXTransfer struct {
	Host        string
	Port        int32
}

// ID ...
func (*VastleXTransfer) ID() uint32 {
	return IDVastleXTransfer
}

// Marshal ...
func (pk *VastleXTransfer) Marshal(w *protocol.Writer) {
	w.String(&pk.Host)
	w.Varint32(&pk.Port)
}

// Unmarshal ...
func (pk *VastleXTransfer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Host)
	r.Varint32(&pk.Port)
	log.DefaultLogger.Info(fmt.Sprintf("Unmarshalling the packet %v:%v", pk.Host, pk.Port))
}