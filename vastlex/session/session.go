package session

import (
	"errors"
	"github.com/VastleLLC/VastleX/vastlex/blocks"
	"github.com/VastleLLC/VastleX/vastlex/entity"
	"github.com/VastleLLC/VastleX/vastlex/server"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
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
	server, err := server.Connect(info, p, conf)
	if err != nil {
		return err
	}
	if p.remote == nil {
		// first server (let client handle startgame etc)
		gameData := server.Conn.GameData()
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
		case <-server.HandleStartGame:
			// connected
		}
	}
	p.remote = server
	return err
}

func (p *Player) handlePackets() {
	go func() {
		for {
			pk, err := p.conn.ReadPacket()
			if err != nil {
				if err.Error() == "error reading packet: connection closed" {
					// Connection closed rip
					if p.remote != nil {
						if p.remote.Conn != nil {
							_ = p.remote.Conn.Close()
						}
					}
					println("Connection closed")
					break
				} else {
					println(err.Error())
				}
				continue
			}
			if p.remote != nil {
				if p.remote.Conn != nil {
					// server is connected
					switch pk := pk.(type) {
					case *packet.CommandRequest:
						strarr := strings.Split(strings.TrimPrefix(pk.CommandLine, "/"), " ")
						cmd, args := strarr[0], strarr[1:]
						switch cmd {
						case "send":
							if len(args) > 0 {
								err := p.Send(p.remote.Info(), server.ConnectConfig{HideMessage: true})
								if err != nil {
									_ = p.conn.WritePacket(&packet.Text{Message: text.Red()("Error transfering you: " + err.Error())})
								} else {
									_ = p.conn.WritePacket(&packet.Text{Message: text.Green()("Successfully transfered you")})
								}
							} else {
								_ = p.conn.WritePacket(&packet.Text{Message: text.Red()("You didn't provide a server address to be transfered to")})
							}
							continue // we shouldn't let pmmp know it was sent
						}
						break
					}
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

func (p *Player) Remote() *server.Remote {
	return p.remote
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
		p.Kick("Unknown error occured in your connection", err.Error())
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
