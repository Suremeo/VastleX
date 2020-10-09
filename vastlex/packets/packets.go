package packets

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

const (
	IDVastleXTransfer = iota + 1000
)

// init registers all packets
func init() {
	packet.Register(IDVastleXTransfer, func() packet.Packet { return &VastleXTransfer{} })
}
