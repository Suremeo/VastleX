package remote

import (
	"encoding/json"
	"github.com/VastleLLC/VastleX/vastlex/blocks"
	"github.com/VastleLLC/VastleX/vastlex/entity"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.uber.org/atomic"
)

type Player interface {
	Conn() *minecraft.Conn
	CurrentId() *atomic.Int64
	Entities() *entity.Store
	Blocks() *blocks.Store
	RemoteDisconnect(error)
	Remote() *Remote
	Dimension() *atomic.Int32
	Send(address string, config ...ConnectConfig) error
}

type Remote struct {
	Player          Player
	Conn            *minecraft.Conn
	Entities        *entity.Store
	Blocks          *blocks.Store
	HandleStartGame chan bool
	connected       bool
}

type ConnectConfig struct {
	Message         string
	HideMessage     bool
	HandleStartgame bool
}

func Connect(address string, player Player, config ...ConnectConfig) (remote *Remote, err error) {
	clientData := player.Conn().ClientData()
	clientData.ThirdPartyName = "VaStLeXiScOoL"                      // ThirdPartyName can be used as a sort of shared secret (Its not implemented yet which is why it can't be set from outside this file)
	clientData.PlatformOfflineID = player.Conn().IdentityData().XUID // Pmmp has an issue getting the XUID with auth disabled so XUID is here to solve the issue.
	conn, err := minecraft.Dialer{
		ClientData:   clientData,
		IdentityData: player.Conn().IdentityData(),
	}.Dial("raknet", address)
	if err != nil {
		return
	}
	remote = &Remote{
		Player:   player,
		Entities: &entity.Store{},
		Blocks:   &blocks.Store{},
		Conn:     conn,
	}
	if len(config) > 0 {
		if config[0].HandleStartgame {
			remote.HandleStartGame = make(chan bool, 1)
		}
	}
	if !config[0].HideMessage {
		msg := text.Bold()(text.Green()("Teleporting..."))
		if config[0].Message != "" {
			msg = config[0].Message
		}
		_ = player.Conn().WritePacket(&packet.SetTitle{
			ActionType: packet.TitleActionSetTitle,
			Text:       msg,
		})
	}
	remote.Blocks.Import(conn.GameData().Blocks)
	if player.Remote() == nil {
		player.Blocks().Import(conn.GameData().Blocks)
	}

	// Get the player up to date with the new things from GameData

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
	remote.handlePackets()
	go func() {
		err = conn.DoSpawn()
		if err != nil {
			player.RemoteDisconnect(err)
			_ = conn.Close()
		} else {
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
	player.Entities().Set(1, int64(conn.GameData().EntityRuntimeID))
	remote.Entities.Set(int64(conn.GameData().EntityRuntimeID), 1)
	remote.connected = true
	return
}

func (server *Remote) handlePackets() {
	go func() {
		for {
			pk, err := server.Conn.ReadPacket()
			if err != nil {
				if err.Error() == "error reading packet: connection closed" {
					server.Player.RemoteDisconnect(err)
					break
				} else {
					println("Error reading packet from the server: " + err.Error())
				}
				continue
			}
			switch pk := pk.(type) {
			case *packet.AddActor:
				newId := server.Player.CurrentId().Inc()
				server.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
				server.Entities.Set(int64(pk.EntityRuntimeID), newId)
				break
			case *packet.AddItemActor:
				newId := server.Player.CurrentId().Inc()
				server.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
				server.Entities.Set(int64(pk.EntityRuntimeID), newId)
				break
			case *packet.AddEntity:
				newId := server.Player.CurrentId().Inc()
				server.Player.Entities().Set(newId, int64(pk.EntityNetworkID))
				server.Entities.Set(int64(pk.EntityNetworkID), newId)
				break
			case *packet.AddPlayer:
				if pk.EntityRuntimeID != 1 {
					newId := server.Player.CurrentId().Inc()
					server.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
					server.Entities.Set(int64(pk.EntityRuntimeID), newId)
				}
				break
			case *packet.RemoveActor:
				server.Player.Entities().Delete(int64(pk.EntityUniqueID))
				server.Entities.Delete(int64(pk.EntityUniqueID))
				break
			case *packet.RemoveEntity:
				server.Player.Entities().Delete(int64(pk.EntityNetworkID))
				server.Entities.Delete(int64(pk.EntityNetworkID))
				break
			case *packet.ChangeDimension:
				server.Player.Dimension().Store(pk.Dimension)
				break

				// StartGame packets
			//case *packet.ResourcePacksInfo:
			//	if server.HandleStartGame != nil {
			//		_ = server.Conn.WritePacket(&packet.ResourcePackClientResponse{
			//			Response: packet.PackResponseAllPacksDownloaded,
			//		})
			//		continue
			//	}
			//	break
			//case *packet.ResourcePackStack:
			//	if server.HandleStartGame != nil {
			//		_ = server.Conn.WritePacket(&packet.ResourcePackClientResponse{
			//			Response: packet.PackResponseCompleted,
			//		})
			//		continue
			//	}
			//	break
			case *packet.ScriptCustomEvent: // Custom packet for transferring
				switch pk.EventName {
				case "vastlex:transfer":
					var send struct {
						Address     string
						Message     string
						HideMessage bool
					}
					err := json.Unmarshal(pk.EventData, &send)
					if err != nil {
						_ = server.Player.Conn().WritePacket(&packet.ScriptCustomEvent{
							EventName: "vastlex:error",
							EventData: []byte(err.Error()),
						})
					} else {
						err := server.Player.Send(send.Address, ConnectConfig{
							Message:     send.Message,
							HideMessage: send.HideMessage,
						})
						if err != nil {
							if server.Conn != nil {
								_ = server.Conn.WritePacket(&packet.ScriptCustomEvent{
									EventName: "vastlex:error",
									EventData: []byte(err.Error()),
								})
							}
						}
					}
					break
				}
				break
			}
			if !server.connected {
				continue
			}
			blocks.TranslatePacket(pk, server.Player.Blocks(), server.Blocks)
			if entity.TranslatePacket(pk, server.Entities) {
				_ = server.Player.Conn().WritePacket(pk)
			}
		}
	}()
}
