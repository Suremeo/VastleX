package server

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func (remote *Remote) clearEntities() {
	remote.Player.UniqueEntities().Range(func(key, value interface{}) bool {
		uid := value.(int64)
		if uid != 1 {
			_ = remote.Player.Conn().WritePacket(&packet.RemoveActor{EntityUniqueID: uid})
		}
		return true
	})
	remote.Entities.Clear()
	remote.UniqueEntities.Clear()
	remote.Player.UniqueEntities().Clear()
	remote.Player.Entities().Clear()
}
