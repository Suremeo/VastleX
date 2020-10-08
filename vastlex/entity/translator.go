package entity

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"reflect"
)

func TranslatePacket(pk packet.Packet, store *Store) bool {
	switch pk := pk.(type) {
	//case *packet.BossEvent:
	//	pk.PlayerUniqueID = store.Get(pk.PlayerUniqueID)
	//	pk.BossEntityUniqueID= store.Get(pk.BossEntityUniqueID)
	//	break
	//case *packet.Camera:
	//	pk.CameraEntityUniqueID = store.Get(pk.CameraEntityUniqueID)
	//	pk.TargetPlayerUniqueID = store.Get(pk.TargetPlayerUniqueID)
	//	break
	//case *packet.SetActorLink:
	//	if store.Get(int64(pk.EntityLink.RiddenEntityUniqueID)) != 0 {
	//		pk.EntityLink.RiddenEntityUniqueID = store.Get(pk.EntityLink.RiddenEntityUniqueID)
	//	}
	//	if store.Get(int64(pk.EntityLink.RiderEntityUniqueID)) != 0 {
	//		pk.EntityLink.RiderEntityUniqueID = store.Get(pk.EntityLink.RiderEntityUniqueID)
	//	}
	//	break
	case *packet.CommandBlockUpdate:
		if store.Get(int64(pk.MinecartEntityRuntimeID)) != 0 {
			pk.MinecartEntityRuntimeID = uint64(store.Get(int64(pk.MinecartEntityRuntimeID)))
		}
		break
	case *packet.MovePlayer:
		if store.Get(int64(pk.RiddenEntityRuntimeID)) != 0 {
			pk.RiddenEntityRuntimeID = uint64(store.Get(int64(pk.RiddenEntityRuntimeID)))
		}
		break
	case *packet.TakeItemActor:
		if store.Get(int64(pk.ItemEntityRuntimeID)) != 0 {
			pk.ItemEntityRuntimeID = uint64(store.Get(int64(pk.ItemEntityRuntimeID)))
		}
		if store.Get(int64(pk.TakerEntityRuntimeID)) != 0 {
			pk.TakerEntityRuntimeID = uint64(store.Get(int64(pk.TakerEntityRuntimeID)))
		}
		break
	case *packet.InventoryTransaction:
		switch trans := pk.TransactionData.(type) {
		case *protocol.UseItemOnEntityTransactionData:
			trans.TargetEntityRuntimeID = uint64(store.Get(int64(trans.TargetEntityRuntimeID)))
			break
		}
		break
	case *packet.Interact:
		if store.Get(int64(pk.TargetEntityRuntimeID)) != 0 {
			pk.TargetEntityRuntimeID = uint64(store.Get(int64(pk.TargetEntityRuntimeID)))
		}
		break
	}
	if reflect.ValueOf(pk).Elem().FieldByName("EntityRuntimeID").IsValid() {
		if store.Get(int64(reflect.ValueOf(pk).Elem().FieldByName("EntityRuntimeID").Uint())) != 0 {
			reflect.ValueOf(pk).Elem().FieldByName("EntityRuntimeID").SetUint(uint64(store.Get(int64(reflect.ValueOf(pk).Elem().FieldByName("EntityRuntimeID").Uint()))))
		} else {
			return false
		}
	}
	return true
}
