package server

import (
	"errors"
	"fmt"
	"github.com/VastleLLC/VastleX/config"
	"github.com/VastleLLC/VastleX/log"
	"github.com/VastleLLC/VastleX/vastlex/blocks"
	"github.com/VastleLLC/VastleX/vastlex/entity"
	"github.com/VastleLLC/VastleX/vastlex/packets"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.uber.org/atomic"
)

// A player connected to a server.
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
	Kick(...string)
}

type Server interface {
	Info() Info
}

type Remote struct {
	Player          Player
	Conn            *minecraft.Conn
	Entities        *entity.Store
	UniqueEntities  *entity.Store
	Blocks          *blocks.Store
	HandleStartGame chan bool
	serverInfo      Info
}

func (remote *Remote) Info() Info {
	return remote.serverInfo
}

type ConnectConfig struct {
	Message         string
	HideMessage     bool
	HandleStartgame bool
}

type Info struct {
	Host string
	Port int
}

func Connect(info Info, player Player, connectConfig ...ConnectConfig) (remote *Remote, err error) {
	clientData := player.Conn().ClientData()
	clientData.ThirdPartyName = config.Config.Proxy.Secret          // ThirdPartyName is used as a placeholder for the connection secret
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

	remote.handlePackets()
	go func() {
		err = conn.DoSpawn()
		if err != nil {
			if player.Server() == remote {
				player.RemoteDisconnect(err)
			}
			_ = conn.Close()
		} else {
			log.Debug().Str("host", info.Host).Int("port", info.Port).Str("username", player.Conn().IdentityData().DisplayName).Msg("Player spawned in")
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
	return
}

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
			case *packet.AddActor:
				newId := remote.Player.CurrentId().Inc()
				remote.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
				remote.Entities.Set(int64(pk.EntityRuntimeID), newId)
				remote.UniqueEntities.Set(pk.EntityUniqueID, int64(pk.EntityRuntimeID))
				remote.Player.UniqueEntities().Set(int64(pk.EntityRuntimeID), pk.EntityUniqueID)
				break
			case *packet.AddItemActor:
				newId := remote.Player.CurrentId().Inc()
				remote.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
				remote.Entities.Set(int64(pk.EntityRuntimeID), newId)
				remote.UniqueEntities.Set(pk.EntityUniqueID, int64(pk.EntityRuntimeID))
				remote.Player.UniqueEntities().Set(int64(pk.EntityRuntimeID), pk.EntityUniqueID)
				break
			case *packet.AddPlayer:
				if pk.EntityRuntimeID != 1 {
					newId := remote.Player.CurrentId().Inc()
					remote.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
					remote.Entities.Set(int64(pk.EntityRuntimeID), newId)
					remote.UniqueEntities.Set(pk.EntityUniqueID, int64(pk.EntityRuntimeID))
					remote.Player.UniqueEntities().Set(int64(pk.EntityRuntimeID), pk.EntityUniqueID)
				}
				break
			case *packet.RemoveActor:
				_ = remote.Player.Conn().WritePacket(&packet.RemoveActor{EntityUniqueID: remote.Player.UniqueEntities().Get(remote.UniqueEntities.Get(pk.EntityUniqueID))})
				rid := remote.UniqueEntities.Get(pk.EntityUniqueID)
				eid := remote.Entities.Get(rid)
				remote.Player.Entities().Delete(eid)
				remote.Entities.Delete(rid)
				remote.UniqueEntities.Delete(pk.EntityUniqueID)
				remote.Player.UniqueEntities().Delete(rid)
				continue
			case *packet.Disconnect:
				remote.Player.RemoteDisconnect(errors.New(pk.Message))
				break
			case *packet.ChangeDimension:
				remote.Player.Dimension().Store(pk.Dimension)
				break
			case *packet.ResourcePacksInfo:
				if remote.HandleStartGame != nil {
					_ = remote.Conn.WritePacket(&packet.ResourcePackClientResponse{
						Response: packet.PackResponseAllPacksDownloaded,
					})
					continue
				}
				break
			case *packet.ResourcePackStack:
				if remote.HandleStartGame != nil {
					_ = remote.Conn.WritePacket(&packet.ResourcePackClientResponse{
						Response: packet.PackResponseCompleted,
					})
					continue
				}
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
			blocks.TranslatePacket(pk, remote.Player.Blocks(), remote.Blocks)
			if entity.TranslatePacket(pk, remote.Entities) {
				_ = remote.Player.Conn().WritePacket(pk)
			}
		}
	}()
}
