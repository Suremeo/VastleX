package vastlex

import (
	"errors"
	"fmt"
	"github.com/VastleLLC/VastleX/config"
	"github.com/VastleLLC/VastleX/log"
	"github.com/VastleLLC/VastleX/vastlex/server"
	"github.com/VastleLLC/VastleX/vastlex/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"net/http"
	_ "net/http/pprof"
)

// VastleX is the main structure for the proxy.
var VastleX Proxy = vastlex

// vastlex is the non interface version of the VastleX variable.
var vastlex = &Structure{
	listener: &minecraft.Listener{
		ErrorLog:               nil,
		AuthenticationDisabled: !config.Config.Minecraft.Auth,
		ServerName:             config.Config.Minecraft.Motd,
		ShowVersion:            config.Config.Minecraft.ShowVersion,
	},
	info: server.Info{
		Host: config.Config.Listener.Host,
		Port: config.Config.Listener.Port,
	},
	players: map[string]Player{},
}

// Structure is the structure of VastleX.
type Structure struct {
	listener *minecraft.Listener
	info     server.Info
	players  map[string]Player // Will be put to use once events are added.
}

// Start starts the proxy.
func Start() (err error) {
	if config.Config.Minecraft.MaxPlayers == 0 {
		log.Title(fmt.Sprintf("%v", log.TotalPlayers))
	} else {
		log.Title(fmt.Sprintf("%v/%v", log.TotalPlayers, config.Config.Minecraft.MaxPlayers))
	}
	if config.Config.Debug.Profiling {
		go func() {
			log.Debug().Str("host", "localhost").Int("port", 6060).Msg("The profiling server is running")
			log.FatalError("Error occured with the profiling server", http.ListenAndServe("localhost:6060", nil))
		}()
	}
	err = vastlex.listener.Listen("raknet", fmt.Sprintf("%v:%v", vastlex.info.Host, vastlex.info.Port))
	if err != nil {
		return err
	}
	log.Info().Str("host", vastlex.info.Host).Int("port", vastlex.info.Port).Str("checksum", log.Checksum).Msg("VastleX is listening for players")

	for {
		conn, err := vastlex.listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			p := session.New(conn.(*minecraft.Conn))
			log.Info().Str("username", p.Conn().IdentityData().DisplayName).Msg("Player connected")
			if config.Config.Lobby.Enabled {
				err = p.Send(server.Info{
					Host: config.Config.Lobby.Host,
					Port: config.Config.Lobby.Port,
				})
				if err != nil {
					p.Kick(text.Red()("We had an error connecting you to a lobby"))
					log.Err().Str("username", p.Identity().DisplayName).Err(err).Msg("Player failed to connect to lobby")
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
	return vastlex.listener.ServerName
}

// ...
func (vastlex *Structure) SetMotd(motd string) {
	vastlex.listener.ServerName = motd
}

// ...
func (vastlex *Structure) Players() map[string]Player {
	return vastlex.players
}

// ...
func (vastlex *Structure) GetPlayer(username string) (Player, error) {
	if vastlex.players[username] != nil {
		return vastlex.players[username], nil
	} else {
		return nil, errors.New("player not found")
	}
}
