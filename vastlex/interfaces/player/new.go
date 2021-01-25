package player

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/minecraftevents"
	"time"
)

// New initializes a interfaces.InternalPlayer for the provided minecraft.Player.
func New(conn minecraft.Player) interfaces.InternalPlayer {
	player := &Player{
		Player:      conn,
		config:      &config.Player{},
		dialer:      nil,
		state:       interfaces.StateWaitingForFirstServer,
		chunkradius: 16,
	}
	go player.listenPackets()
	player.HandleEvent(&minecraftevents.Close{}, func() {
		if player.dialer != nil {
			_ = player.dialer.Close()
			player.dialer = nil
		}
		player.state = interfaces.StateDisconnected
		go func() {
			time.Sleep(5*time.Second)
			if player.dialer != nil {
				_ = player.dialer.Close()
				player.dialer = nil
			}
		}()
	})
	return player
}
