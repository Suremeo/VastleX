package entities

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryTransaction struct{}

func (InventoryTransaction) Translate(pk packet.Packet, eid1, eid2 int64, uid1, uid2 int64) {
	switch trans := pk.(*packet.InventoryTransaction).TransactionData.(type) {
	case *protocol.UseItemOnEntityTransactionData:
		if trans.TargetEntityRuntimeID == uint64(eid1) {
			trans.TargetEntityRuntimeID = uint64(eid2)
		} else if trans.TargetEntityRuntimeID == uint64(eid2) {
			trans.TargetEntityRuntimeID = uint64(eid1)
		}
	}
}
