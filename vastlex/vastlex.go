package vastlex

import (
	"errors"
	"fmt"
	"github.com/VastleLLC/VastleX/config"
	"github.com/VastleLLC/VastleX/log"
	"github.com/VastleLLC/VastleX/vastlex/server"
	"github.com/VastleLLC/VastleX/vastlex/session"
	"github.com/sandertv/gophertunnel/minecraft"
)

// VastleX is the main structure for the proxy.
var VastleX = &Structure{
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
	players  map[string]Player
}

// Start starts the proxy.
func Start() (err error) {
	err = VastleX.listener.Listen("raknet", fmt.Sprintf("%v:%v", VastleX.info.Host, VastleX.info.Port))
	if err != nil {
		return err
	}
	log.Info().Str("host", VastleX.info.Host).Int("port", VastleX.info.Port).Msg("VastleX is listening for players")

	for {
		conn, err := VastleX.listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			VastleX.players[conn.(*minecraft.Conn).IdentityData().DisplayName] = session.New(conn.(*minecraft.Conn))
			log.Info().Str("username", conn.(*minecraft.Conn).IdentityData().DisplayName).Msg("Player connected")
			if config.Config.Lobby.Enabled {
				err = VastleX.players[conn.(*minecraft.Conn).IdentityData().DisplayName].Send(server.Info{})
				if err != nil {
					log.Err().Str("username", conn.(*minecraft.Conn).IdentityData().DisplayName).Err(err).Msg("Player failed to connect to lobby")
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
