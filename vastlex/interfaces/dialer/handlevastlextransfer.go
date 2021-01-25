package dialer

import (
	"fmt"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/custompackets"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// handleVastleXTransfer handles the VastleXTransfer custom packet sent by the remote minecraft.Connection.
func handleVastleXTransfer(dialer *Dialer, pak packet.Packet) {
	pk := pak.(*custompackets.VastleXTransfer)
	log.DefaultLogger.Info(fmt.Sprintf("VastleX transfer was called | IP: %v | PORT: %v", pk.Host, pk.Port))
	err := dialer.player.Send(pk.Host, int(pk.Port))
	if err != nil {
		log.DefaultLogger.Error(err)
	}
}
