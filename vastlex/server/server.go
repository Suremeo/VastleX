package server

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/vastlellc/vastlex/config"
	"github.com/vastlellc/vastlex/log"
	"github.com/vastlellc/vastlex/vastlex/blocks"
	"github.com/vastlellc/vastlex/vastlex/entity"
	"github.com/vastlellc/vastlex/vastlex/packets"
	"go.uber.org/atomic"
)

// Player represents a player who is connected to the server.
type Player interface {
	Conn() *minecraft.Conn
	CurrentId() *atomic.Int64
	Entities() *entity.Store
	UniqueEntities() *entity.Store
	Blocks() *blocks.Store
	RemoteDisconnect(error)
	Server() Server
	Dimension() *atomic.Int32
	Send(info Info, config ...ConnectConfig) error
}

// Server represents a Minecraft server.
type Server interface {
	Info() Info
}

// Remote is a connection to a server.
type Remote struct {
	Player          Player
	Conn            *minecraft.Conn
	Entities        *entity.Store
	UniqueEntities  *entity.Store
	Blocks          *blocks.Store
	HandleStartGame chan bool
	serverInfo      Info
	connected       bool
}

// Info returns info about the server.
func (remote *Remote) Info() Info {
	return remote.serverInfo
}

// ConnectConfig is the configuration for sending a player.
type ConnectConfig struct {
	Message         string
	HideMessage     bool
	HandleStartgame bool
}

// Info is the info about a Minecraft server.
type Info struct {
	Host string
	Port int
}

// Connect opens a connection to a server to transfer the player.
func Connect(info Info, player Player, connectConfig ...ConnectConfig) (remote *Remote, err error) {
	clientData := player.Conn().ClientData()
	clientData.ThirdPartyName = config.Config.Proxy.Secret // ThirdPartyName is used as a placeholder for the connection secret
	clientData.PlatformOfflineID = player.Conn().RemoteAddr().String()
	clientData.PlatformOnlineID = player.Conn().IdentityData().XUID // Pmmp has an issue getting the XUID with auth disabled so the PlatformOnlineID is set to the players XUID to solve the issue.
	conn, err := minecraft.Dialer{
		ClientData:   clientData,
		IdentityData: player.Conn().IdentityData(),
	}.Dial("raknet", fmt.Sprintf("%v:%v", info.Host, info.Port))
	if err != nil {
		return
	}
	remote = &Remote{
		Player:         player,
		Entities:       &entity.Store{},
		Blocks:         &blocks.Store{},
		UniqueEntities: &entity.Store{},
		Conn:           conn,
		serverInfo:     info,
	}
	if len(connectConfig) > 0 {
		if connectConfig[0].HandleStartgame {
			remote.HandleStartGame = make(chan bool, 1)
		} else {
			player.Blocks().Import(conn.GameData().Blocks)
		}
	}
	if !connectConfig[0].HideMessage {
		msg := text.Bold()(text.Green()("Teleporting..."))
		if connectConfig[0].Message != "" {
			msg = connectConfig[0].Message
		}
		_ = player.Conn().WritePacket(&packet.SetTitle{
			ActionType: packet.TitleActionSetTitle,
			Text:       msg,
		})
	}
	remote.Blocks.Import(conn.GameData().Blocks)

	// Get the session up to date with the new things from GameData

	if remote.Player.Dimension().Load() != remote.Conn.GameData().Dimension {
		_ = remote.Player.Conn().WritePacket(&packet.ChangeDimension{
			Dimension: remote.Conn.GameData().Dimension,
			Position:  mgl32.Vec3{float32(remote.Conn.GameData().WorldSpawn.X()), float32(remote.Conn.GameData().WorldSpawn.Y()), float32(remote.Conn.GameData().WorldSpawn.Z())},
			Respawn:   true,
		})
	}
	_ = remote.Player.Conn().WritePacket(&packet.SetPlayerGameType{GameType: remote.Conn.GameData().PlayerGameMode})
	_ = remote.Player.Conn().WritePacket(&packet.GameRulesChanged{GameRules: remote.Conn.GameData().GameRules})
	_ = remote.Player.Conn().WritePacket(&packet.MovePlayer{
		EntityRuntimeID: 1,
		Position:        mgl32.Vec3{float32(remote.Conn.GameData().WorldSpawn.X()), float32(remote.Conn.GameData().WorldSpawn.Y()), float32(remote.Conn.GameData().WorldSpawn.Z())},
		Pitch:           remote.Conn.GameData().Pitch,
		Yaw:             remote.Conn.GameData().Yaw,
		HeadYaw:         remote.Conn.GameData().Yaw,
	})

	remote.clearEntities()

	player.Entities().Set(1, int64(conn.GameData().EntityRuntimeID))
	remote.Entities.Set(int64(conn.GameData().EntityRuntimeID), 1)

	player.UniqueEntities().Set(1, int64(conn.GameData().EntityUniqueID))
	remote.UniqueEntities.Set(conn.GameData().EntityUniqueID, 1)

	go func() {
		err = conn.DoSpawn()
		if err != nil {
			if player.Server() == remote {
				player.RemoteDisconnect(err)
			}
			_ = conn.Close()
		} else {
			log.Debug().Str("host", info.Host).Int("port", info.Port).Str("username", player.Conn().IdentityData().DisplayName).Msg("Player spawned in")
			remote.handlePackets()
			// Spawn is done so we can clear the previous title
			_ = player.Conn().WritePacket(&packet.SetTitle{
				ActionType: packet.TitleActionSetTitle,
				Text:       " ",
			})
			if remote.HandleStartGame != nil {
				remote.HandleStartGame <- true
			}
		}
	}()
	remote.connected = true
	return
}

// handlePackets handles all packets coming from the server, translates them and sends to the client.
func (remote *Remote) handlePackets() {
	go func() {
		for {
			pk, err := remote.Conn.ReadPacket()
			if err != nil {
				if err.Error() == "error reading packet: connection closed" {
					if remote.Player.Server() == remote {
						remote.Player.RemoteDisconnect(err)
					}
					break
				} else {
					println("Error reading packet from the server: " + err.Error())
				}
				continue
			}
			switch pk := pk.(type) {
			case *packet.Disconnect:
				println("Disconnected: " + pk.Message)
				break
			case *packet.ChangeDimension:
				remote.Player.Dimension().Store(pk.Dimension)
				break
			case *packets.VastleXTransfer:
				err = remote.Player.Send(Info{
					Host: pk.Host,
					Port: int(pk.Port),
				}, ConnectConfig{
					Message:     pk.Message,
					HideMessage: pk.HideMessage,
				})

				break
			}
			if !remote.connected {
				continue
			}
			remote.handleAddRemoveEntities(pk)
			blocks.TranslatePacket(pk, remote.Player.Blocks(), remote.Blocks)
			if entity.TranslatePacket(pk, remote.Entities) {
				_ = remote.Player.Conn().WritePacket(pk)
			}
		}
	}()
}
