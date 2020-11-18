package player

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/interfaces/dialer"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/translators/blocks"
	"github.com/VastleLLC/VastleX/vastlex/translators/entities"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	"time"
)

// ...
type Player struct {
	minecraft.Player
	config *config.Player
	dialer interfaces.Dialer
	state interfaces.State
	blocks *blocks.Store
	chunkradius int32
	onFallback bool
}

// Blocks returns the block store for the player (Used for translating block runtime ids).
func (player *Player) Blocks() *blocks.Store {
	return player.blocks
}

// Send connects the player to the provided server.
func (player *Player) Send(ip string, port int) error {
	player.onFallback = false
	if player.state != interfaces.StateWaitingForFirstServer {
		player.state = interfaces.StateWaitingForNewServer
	}
	return dialer.Connect(player, ip, port)
}

// Config returns the configuration for the player.
func (player *Player) Config() *config.Player {
	return player.config
}

// SetConfig sets the configuration for the player.
func (player *Player) SetConfig(config *config.Player) {
	player.config = config
}

// Dialer returns the interfaces.Dialer the player is currently connected to.
func (player *Player) Dialer() interfaces.Dialer {
	return player.dialer
}

// SetDialer sets the interfaces.Dialer for the player.
func (player *Player) SetDialer(dialer interfaces.Dialer) {
	player.dialer = dialer
}

// State returns the current interfaces.State of the player.
func (player *Player) State() interfaces.State {
	return player.state
}

// SetState updates the interfaces.State of the player.
func (player *Player) SetState(state interfaces.State) {
	player.state = state
}

// listenPackets listens and forwards packets for the player.
func (player *Player) listenPackets() {
	for {
		pk, err := player.ReadPacket()
		if err != nil {
			break
		}
		if player.state == interfaces.StateConnected {
			if Handles[pk.ID()] != nil {
				Handles[pk.ID()](player, pk)
			}
			if config.Config.Debug.BlockTranslating {
				blocks.TranslatePacket(pk, player.dialer.Blocks(), player.Blocks())
			}
			entities.Pool[pk.ID()]().Translate(pk, 1, int64(player.dialer.EntityId()), 1, int64(player.dialer.UniqueId()))
			_ = player.dialer.WritePacket(pk)
		}
	}
	_ = player.Close()
}

// Kick kicks the player from the server using the provided message.
func (player *Player) Kick(reason ...string) {
	if len(reason) == 0 {
		reason = []string{text.Red()("No reason provided")}
	}
	_ = player.WritePacket(&packet.Disconnect{Message: strings.Join(reason, "\n")})
	time.Sleep(time.Second / 20 * 2)
	_ = player.Close()
}

// Message sends a chat message to the player.
func (player *Player) Message(msg string) error {
	return player.WritePacket(&packet.Text{Message: msg})
}

// KickOrFallback attempts to connect the player to the fallback server (If enabled), if it is unable to connect it kicks the player using the provided message.
func (player *Player) KickOrFallback(msg string) {
	if player.onFallback {
		player.Kick(text.Red()("The fallback server went down."))
		return
	}
	if config.Config.Fallback.Enabled {
		err := player.Send(config.Config.Fallback.Host, config.Config.Fallback.Port)
		if err != nil {
			player.Kick(text.Red()("Unable to connect you to a fallback server."))
		} else {
			player.onFallback = true
			_ = player.Message(text.Red()("Oof! The server you were on went down, so you were connected to a fallback server."))
		}
	} else {
		player.Kick(msg)
	}
}

// ChunkRadius returns the current chunk radius for the player.
func (player *Player) ChunkRadius() int32 {
	return player.chunkradius
}