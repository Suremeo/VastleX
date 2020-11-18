package dialer

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// handleStartGame handles the StartGame packet sent by the remote minecraft.Connection.
func handleStartGame(dialer *Dialer, pak packet.Packet) {
	pk := pak.(*packet.StartGame)
	dialer.entityId = int64(pk.EntityRuntimeID)
	dialer.uniqueId = pk.EntityUniqueID
	close(dialer.ready)
	if dialer.player.State() != interfaces.StateWaitingForFirstServer {
		_ = dialer.player.WritePacket(&packet.SetPlayerGameType{GameType: pk.PlayerGameMode})
		_ = dialer.player.WritePacket(&packet.GameRulesChanged{GameRules: pk.GameRules})
		_ = dialer.player.WritePacket(&packet.MovePlayer{
			EntityRuntimeID:          1,
			Position:                 pk.PlayerPosition,
			Pitch:                    pk.Pitch,
			Yaw:                      pk.Yaw,
			HeadYaw:                  pk.Yaw,
		})
		_ = dialer.WritePacket(&packet.SetLocalPlayerAsInitialised{EntityRuntimeID: pk.EntityRuntimeID})
		_ = dialer.WritePacket(&packet.RequestChunkRadius{ChunkRadius: dialer.player.ChunkRadius()})
	} else {
		if config.Config.Debug.BlockTranslating {
			dialer.player.Blocks().Import(pk.Blocks)
		}
	}
	if config.Config.Debug.BlockTranslating {
		dialer.Blocks().Import(pk.Blocks)
	}
}