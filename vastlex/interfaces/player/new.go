package player

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/VastleLLC/VastleX/vastlex/translators/blocks"
)

// New initializes a interfaces.InternalPlayer for the provided minecraft.Player.
func New(conn minecraft.Player) interfaces.InternalPlayer {
	player := &Player{
		Player: conn,
		config: &config.Player{},
		dialer: nil,
		state:  interfaces.StateWaitingForFirstServer,
		chunkradius: 16,
	}
	if config.Config.Debug.BlockTranslating {
		player.blocks =  &blocks.Store{}
	}
	go player.listenPackets()
	player.HandleEvent(&events.Close{}, func() {
		if player.dialer != nil {
			_ = player.dialer.Close()
		}
	})
	return player
}