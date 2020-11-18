package dialer

import (
	"errors"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/VastleLLC/VastleX/vastlex/translators/blocks"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"sync"
)

// Connect connects a player to a remote server.
func Connect(player interfaces.InternalPlayer, address string, port int) error {
	client := player.ClientData()
	client.ThirdPartyName = config.Config.Security.Secret // ThirdPartyName is used as a placeholder for the connection secret.
	client.ServerAddress = player.RemoteAddr().String() // ServerAddress is used as a placeholder for the players ip.
	client.PlatformOnlineID = player.Identity().XUID // Pmmp has an issue getting the XUID with auth disabled so the PlatformOnlineID is set to the players XUID to solve the issue.
	conn, err := minecraft.Dial(player.Identity(), client, fmt.Sprintf("%s:%v", address, port))
	if err != nil {
		return errors.New("error connecting: dialer returned error: " + err.Error())
	}
	dialer := &Dialer{
		Dialer:   conn,
		player:   player,
		mutex:    sync.Mutex{},
		ready: make(chan struct{}),
	}
	if config.Config.Debug.BlockTranslating {
		dialer.blocks = &blocks.Store{}
	}
	dialer.HandleEvent(&events.Close{}, func() {
		dialer.entities.Range(func(key, value interface{}) bool {
			_ = player.WritePacket(&packet.RemoveActor{EntityUniqueID: key.(int64)})
			return true
		})
		dialer.entities = sync.Map{}
	})
	if player.Dialer() != nil {
		_ = player.Dialer().Close()
	}
	go dialer.listenPackets()
	<-dialer.ready
	player.SetDialer(dialer)
	player.SetState(interfaces.StateConnected)
	dialer.HandleEvent(&events.Close{}, func() {
		log.DefaultLogger.Debug("Remote connection for " + player.Identity().DisplayName + " on " + fmt.Sprintf("%v:%v", address, port) + " was closed.")
		if !dialer.leaving {
			dialer.player.KickOrFallback(text.Colourf("<red>We had an error connecting you to a lobby</red>"))
		}
	})
	log.DefaultLogger.Debug(player.Identity().DisplayName + " connected to " + fmt.Sprintf("%v:%v", address, port))
	return nil
}