package blocks

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// TranslatePacket translates the block runtime ids for the packet.
func TranslatePacket(pk packet.Packet, store1, store2 *Store) {
	switch pk := pk.(type) {
	case *packet.InventoryTransaction:
		switch trans := pk.TransactionData.(type) {
		case *protocol.UseItemTransactionData:
			trans.BlockRuntimeID = uint32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(trans.BlockRuntimeID))))
			break
		}
		break
	case *packet.UpdateBlock:
		pk.NewBlockRuntimeID = uint32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.NewBlockRuntimeID))))
		break
	case *packet.UpdateBlockSynced:
		pk.NewBlockRuntimeID = uint32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.NewBlockRuntimeID))))
		break
	case *packet.LevelSoundEvent:
		switch pk.SoundType {
		case packet.SoundEventPlace:
			pk.ExtraData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.ExtraData))))
			break
		case packet.SoundEventBreakBlock:
			pk.ExtraData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.ExtraData))))
			break
		}
		break
	case *packet.LevelEvent:
		switch pk.EventType {
		case packet.EventParticleDestroy:
			pk.EventData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.EventData))))
			break
		case packet.EventParticlePunchBlock:
			pk.EventData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.EventData&0xffffff)))) | (pk.EventData>>24)<<24
			break
		case packet.EventBlockStartBreak:
			pk.EventData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.EventData))))
			break
		case packet.EventBlockStopBreak:
			pk.EventData = int32(store1.RuntimeFromHash(store2.HashFromRuntime(int64(pk.EventData))))
			break
		}
		break
	}
}