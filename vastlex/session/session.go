package session

import (
	"errors"
	"fmt"
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

// Runtime check to see if the Player structure implements the server.Player interface.
var _ server.Player = &Player{}

// Player represents a player connected to the proxy.
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

// New initializes a player using the supplied *minecraft.Conn.
func New(conn *minecraft.Conn) *Player {
	log.TotalPlayers++
	if config.Config.Minecraft.MaxPlayers == 0 {
		log.Title(fmt.Sprintf("%v", log.TotalPlayers))
	} else {
		log.Title(fmt.Sprintf("%v/%v", log.TotalPlayers, config.Config.Minecraft.MaxPlayers))
	}
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

// Send transfers a player to a different server.
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
		_ = p.remote.Conn.Close()
		select {
		case <-time.After(10 * time.Second):
			return errors.New("startgame timed out")
		case <-remote.HandleStartGame:
			// connected
		}
	}
	p.remote = remote
	return err
}

// handlePacket handles the packets sent by the client and directs them to the server they are currently connected to.
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

// CurrentId returns the entity id counter for the player.
func (p *Player) CurrentId() *atomic.Int64 {
	return p.currentId
}

// Conn returns the *minecraft.Conn for the player.
func (p *Player) Conn() *minecraft.Conn {
	return p.conn
}

// Entities returns the entity store for the player.
func (p *Player) Entities() *entity.Store {
	return p.entities
}

// UniqueEntities returns the unique entity store for the player.
func (p *Player) UniqueEntities() *entity.Store {
	return p.uniqueEntities
}

// Blocks returns the block store for the player.
func (p *Player) Blocks() *blocks.Store {
	return p.blocks
}

// Dimension returns the player's current dimension.
func (p *Player) Dimension() *atomic.Int32 {
	return p.dimension
}

// Server returns the current server that the player is connected to.
func (p *Player) Server() server.Server {
	return p.remote
}

// Identity returns the players login.IdentityData.
func (p *Player) Identity() login.IdentityData {
	return p.conn.IdentityData()
}

// Message sends a chat message to the player.
func (p *Player) Message(message string) error {
	return p.conn.WritePacket(&packet.Text{Message: message})
}

// WritePacket writes a packet directly to the player.
func (p *Player) WritePacket(packet packet.Packet) error {
	return p.conn.WritePacket(packet)
}

// RemoteDisconnect is called when a server disconnects the player.
func (p *Player) RemoteDisconnect(err error) {
	if !p.sending {
		if config.Config.Fallback.Enabled && config.Config.Fallback.Host != p.remote.Info().Host && config.Config.Fallback.Port != p.remote.Info().Port {
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

// Kick kicks the player from the proxy.
func (p *Player) Kick(msg ...string) {
	log.TotalPlayers--
	if config.Config.Minecraft.MaxPlayers == 0 {
		log.Title(fmt.Sprintf("%v", log.TotalPlayers))
	} else {
		log.Title(fmt.Sprintf("%v/%v", log.TotalPlayers, config.Config.Minecraft.MaxPlayers))
	}
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
