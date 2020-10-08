package player

import (
	"errors"
	"github.com/SaiCoDev/gophertunnel/minecraft/text"
	"github.com/VastleLLC/VastleX/vastlex/blocks"
	"github.com/VastleLLC/VastleX/vastlex/entity"
	"github.com/VastleLLC/VastleX/vastlex/remote"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
	"strings"
	"time"
)

var _ remote.Player = &Player{} // check if it implements remote.Player

type Player struct {
	currentId      *atomic.Int64
	conn           *minecraft.Conn
	remote         *remote.Remote
	entities       *entity.Store
	uniqueEntities *entity.Store
	blocks         *blocks.Store
	dimension      *atomic.Int32
	sending        bool
}

func NewPlayer(conn *minecraft.Conn) *Player {
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

func (p *Player) Send(address string, config ...remote.ConnectConfig) error {
	p.sending = true
	defer func() {
		p.sending = false
	}()

	conf := remote.ConnectConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	if p.remote != nil {
		conf.HandleStartgame = true
	}
	server, err := remote.Connect(address, p, conf)
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
					// remote is connected
					switch pk := pk.(type) {
					case *packet.CommandRequest:
						strarr := strings.Split(strings.TrimPrefix(pk.CommandLine, "/"), " ")
						cmd, args := strarr[0], strarr[1:]
						switch cmd {
						case "send":
							if len(args) > 0 {
								err := p.Send(args[0], remote.ConnectConfig{HideMessage: true})
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

func (p *Player) Remote() *remote.Remote {
	return p.remote
}

func (p *Player) Dimension() *atomic.Int32 {
	return p.dimension
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
