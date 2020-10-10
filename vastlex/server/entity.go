package server

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

// clearEntities clears all entities that have been sent to the client.
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

// handleAddRemoveEntities handles packets related to adding and removing entities.
func (remote *Remote) handleAddRemoveEntities(pk packet.Packet) bool {
	switch pk := pk.(type) {
	case *packet.AddActor:
		newId := remote.Player.CurrentId().Inc()
		remote.Player.Entities().Set(newId, int64(pk.EntityRuntimeID))
		remote.Entities.Set(int64(pk.EntityRuntimeID), newId)
		remote.UniqueEntities.Set(pk.EntityUniqueID, int64(pk.EntityRuntimeID))
		remote.Player.UniqueEntities().Set(int64(pk.EntityRuntimeID), pk.EntityUniqueID)
		break
	case *packet.AddItemActor:
		if pk.Item.NetworkID == 446 && len(pk.Item.NBTData) != 0 { // Banners with NBT data seem to crash the client when routed through proxy and don't crash when sent directly to the player. (Might be PMMP).
			return false
		}
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
		old := pk.EntityUniqueID
		rid := remote.UniqueEntities.Get(pk.EntityUniqueID)
		eid := remote.Entities.Get(rid)
		remote.Player.Entities().Delete(eid)
		remote.Entities.Delete(rid)
		pk.EntityUniqueID = remote.Player.UniqueEntities().Get(remote.UniqueEntities.Get(pk.EntityUniqueID))
		remote.UniqueEntities.Delete(old)
		remote.Player.UniqueEntities().Delete(rid)
	}
	return true
}