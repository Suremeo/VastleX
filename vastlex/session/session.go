package session

import (
	"errors"
	"github.com/VastleLLC/VastleX/config"
	"github.com/VastleLLC/VastleX/log"
	"github.com/VastleLLC/VastleX/vastlex/blocks"
	"github.com/VastleLLC/VastleX/vastlex/entity"
	"github.com/VastleLLC/VastleX/vastlex/server"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"strings"
	"time"
)

var _ server.Player = &Player{} // check if it implements server.Player

type Player struct {
	currentId      *atomic.Int64
	conn           *minecraft.Conn
	remote         *server.Remote
	entities       *entity.Store
	uniqueEntities *entity.Store
	blocks         *blocks.Store
	dimension      *atomic.Int32
	sending        bool
}

func New(conn *minecraft.Conn) *Player {
	p := &Player{
		currentId:      atomic.NewInt64(2),
		conn:           conn,
		entities:       &entity.Store{},
		uniqueEntities: &entity.Store{},
		blocks:         &blocks.Store{},
		dimension:      atomic.NewInt32(0),
	}
	p.handlePackets()
	return p
}

func (p *Player) Send(info server.Info, config ...server.ConnectConfig) error {
	p.sending = true
	defer func() {
		p.sending = false
	}()

	conf := server.ConnectConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	if p.remote != nil {
		conf.HandleStartgame = true
	}
	remote, err := server.Connect(info, p, conf)
	if err != nil {
		return err
	}
	if p.remote == nil {
		p.remote = remote
		// first remote (let client handle startgame etc)
		gameData := remote.Conn.GameData()
		gameData.EntityUniqueID = 1
		gameData.EntityRuntimeID = 1
		err = p.conn.StartGame(gameData)
		if err != nil {
			p.RemoteDisconnect(err)
		}
		_ = p.conn.WritePacket(&packet.AdventureSettings{
			Flags:             0,
			PermissionLevel:   packet.PermissionLevelMember,
			PlayerUniqueID:    1,
			ActionPermissions: uint32(packet.ActionPermissionBuildAndMine | packet.ActionPermissionDoorsAndSwitched | packet.ActionPermissionOpenContainers | packet.ActionPermissionAttackPlayers | packet.ActionPermissionAttackMobs),
		})
	} else {
		p.remote = remote
		_ = p.remote.Conn.Close()
		select {
		case <-time.After(10 * time.Second):
			return errors.New("startgame timed out")
		case <-remote.HandleStartGame:
			// connected
		}
	}
	return err
}

func (p *Player) handlePackets() {
	go func() {
		for {
			pk, err := p.conn.ReadPacket()
			if err != nil {
				if err.Error() == "error reading packet: connection closed" {
					if p.remote != nil {
						if p.remote.Conn != nil {
							_ = p.remote.Conn.Close()
						}
					}
					p.Kick()
					break
				} else {
					println(err.Error())
				}
				continue
			}
			if p.remote != nil {
				if p.remote.Conn != nil {
					// server is connected
					blocks.TranslatePacket(pk, p.blocks, p.remote.Blocks)
					if entity.TranslatePacket(pk, p.entities) {
						_ = p.remote.Conn.WritePacket(pk)
					}
				}
			}
		}
	}()
}

func (p *Player) CurrentId() *atomic.Int64 {
	return p.currentId
}

func (p *Player) Conn() *minecraft.Conn {
	return p.conn
}

func (p *Player) Entities() *entity.Store {
	return p.entities
}

func (p *Player) UniqueEntities() *entity.Store {
	return p.uniqueEntities
}

func (p *Player) Blocks() *blocks.Store {
	return p.blocks
}

func (p *Player) Dimension() *atomic.Int32 {
	return p.dimension
}

func (p *Player) Server() server.Server {
	return p.remote
}

func (p *Player) Identity() login.IdentityData {
	return p.conn.IdentityData()
}

func (p *Player) Message(message string) error {
	return p.conn.WritePacket(&packet.Text{Message: message})
}

func (p *Player) WritePacket(packet packet.Packet) error {
	return p.conn.WritePacket(packet)
}

func (p *Player) RemoteDisconnect(err error) {
	if !p.sending {
		log.Debug().Str("username", p.Identity().DisplayName).Err(err).Msg("Player disconnected from server")
		if config.Config.Fallback.Enabled {
			if p.remote.Info().Host == config.Config.Fallback.Host && p.remote.Info().Port == config.Config.Fallback.Port {
				// They got disconnected from the fallback server so we shouldn't send them to it again.
				p.remote = nil
				p.Kick("Unknown error occured in your connection", err.Error())
			}
			err = p.Send(server.Info{
				Host: config.Config.Fallback.Host,
				Port: config.Config.Fallback.Port,
			})
			if err != nil {
				p.remote = nil
				p.Kick("Unknown error occured in your connection", err.Error())
			}
		} else {
			p.remote = nil
			p.Kick("Unknown error occured in your connection", err.Error())
		}
	}
}

func (p *Player) Kick(msg ...string) {
	m := "No reason provided"
	if len(msg) > 0 {
		m = strings.Join(msg, "\n")
	}
	_ = p.conn.WritePacket(&packet.Disconnect{
		Message: m,
	})
	if p.remote != nil {
		if p.remote.Conn != nil {
			_ = p.remote.Conn.Close()
		}
	}
	_ = p.conn.Close()
}
