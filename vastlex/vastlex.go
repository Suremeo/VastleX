package vastlex

import (
	"errors"
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/interfaces/player"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/events"
	"github.com/nakabonne/gosivy/agent"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"net/http"
)

// VastleX is the main structure for the proxy.
var VastleX interfaces.VastleX = vastlex

// vastlex is the non interface version of the VastleX variable.
var vastlex = &Structure{
	players: map[string]interfaces.Player{},
}

// Structure is the structure of VastleX.
type Structure struct {
	listener minecraft.Listener
	players  map[string]interfaces.Player // Will be put to use once events are added.
}

// Start starts the proxy.
func Start() (err error) {
	log.UpdatePlayerCount(0)
	if config.Config.Debug.Profiling.PPROF.Enabled {
		go func() {
			log.DefaultLogger.Info("The PPROF (https://github.com/google/pprof) profiler is running on " + config.Config.Debug.Profiling.PPROF.Address)
			log.DefaultLogger.Fatal(http.ListenAndServe(config.Config.Debug.Profiling.PPROF.Address, nil))
		}()
	}
	if config.Config.Debug.Profiling.GOSIVY.Enabled {
		log.DefaultLogger.Info("The GOSIVY (https://github.com/nakabonne/gosivy) profiler is running on " + config.Config.Debug.Profiling.GOSIVY.Address)
		if err := agent.Listen(agent.Options{
			Addr: config.Config.Debug.Profiling.GOSIVY.Address,
		}); err != nil {
			log.DefaultLogger.Fatal(err)
		}
		defer agent.Close()
	}
	l, err := minecraft.Listen()
	if err != nil {
		return err
	}
	vastlex.listener = l
	log.DefaultLogger.Info("VastleX is running on " + fmt.Sprintf("%v:%v",  config.Config.Listener.Host,  config.Config.Listener.Port))
	for {
		conn := vastlex.listener.Accept()
		go func() {
			p := player.New(conn)
			log.DefaultLogger.Info(conn.Identity().DisplayName + " logged in.")
			vastlex.players[conn.Identity().XUID] = p
			p.HandleEvent(&events.Close{}, func() {
				delete(vastlex.players, conn.Identity().XUID)
				log.UpdatePlayerCount(len(vastlex.players))
				log.DefaultLogger.Warn(conn.Identity().DisplayName + " logged out.")
			})
			log.UpdatePlayerCount(len(vastlex.players))
			if config.Config.Lobby.Enabled {
				err = p.Send(config.Config.Lobby.Host, config.Config.Lobby.Port)
				if err != nil {
					log.DefaultLogger.Warn(conn.Identity().DisplayName + " failed to connect to a lobby.")
					p.KickOrFallback(text.Colourf("<red>We had an error connecting you to a lobby</red>"))
				}
			}
		}()
	}
}

// ...
func (vastlex *Structure) Config() config.Structure {
	return config.Config
}

// ...
func (vastlex *Structure) Motd() string {
	return vastlex.listener.Motd()
}

// ...
func (vastlex *Structure) SetMotd(motd string) {
	vastlex.listener.SetMotd(motd)
}

// ...
func (vastlex *Structure) Players() map[string]interfaces.Player {
	return vastlex.players
}

// ...
func (vastlex *Structure) GetPlayer(username string) (interfaces.Player, error) {
	if vastlex.players[username] != nil {
		return vastlex.players[username], nil
	} else {
		return nil, errors.New("player not found")
	}
}