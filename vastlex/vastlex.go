package vastlex

import (
	"fmt"
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/interfaces/player"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/minecraftevents"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/internalevents"
	"github.com/VastleLLC/VastleX/vastlex/plugin/manager"
	"github.com/nakabonne/gosivy/agent"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"net/http"
	"sync"
)

// VastleX is the main structure for the proxy.
var VastleX interfaces.VastleX = vastlex

// vastlex is the non interface version of the VastleX variable.
var vastlex = &Structure{
	players: &sync.Map{},
}

// Structure is the structure of VastleX.
type Structure struct {
	listener minecraft.Listener
	players  *sync.Map // Will be put to use once minecraftevents are added.
}

// Start starts the proxy.
func Start() (err error) {
	config.LoadConfig()
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
	log.DefaultLogger.Info("VastleX is listening on " + fmt.Sprintf("%v:%v", config.Config.Listener.Host, config.Config.Listener.Port))
	for {
		conn := vastlex.listener.Accept()
		go func() {
			p := player.New(conn)
			vastlex.players.Store(conn.Identity().XUID, p)
			if manager.HandleAddPlayer != nil {
				go manager.HandleAddPlayer(p)
			}
			p.HandleEvent(&minecraftevents.Close{}, func() {
				vastlex.players.Delete(conn.Identity().XUID)
				if manager.HandleRemovePlayer != nil {
					go manager.HandleRemovePlayer(conn.Identity().XUID)
				}
				manager.WriteAll(actions.Event{}.New().Encode(&internalevents.RemovePlayer{XUID: p.Identity().XUID}))
				log.UpdatePlayerCount(int64(vastlex.listener.PlayerCount()))
				//log.DefaultLogger.Warn(conn.Identity().DisplayName + " logged out.")
			})
			log.UpdatePlayerCount(int64(vastlex.listener.PlayerCount()))
			//log.DefaultLogger.Info(conn.Identity().DisplayName + " logged in.")
			manager.WriteAll(actions.Event{}.New().Encode(&internalevents.AddPlayer{IdentityData: p.Identity()}))
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

func (vastlex *Structure) Count() int {
	return vastlex.listener.PlayerCount()
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
func (vastlex *Structure) Players() *sync.Map {
	return vastlex.players
}

// ...
func (vastlex *Structure) GetPlayer(xuid string) (interfaces.Player, bool) {
	p, ok := vastlex.players.Load(xuid)
	if !ok {
		return nil, ok
	}
	pl, ok := p.(interfaces.Player)
	return pl, ok
}
