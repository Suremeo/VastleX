package player

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Handles contains handlers for certain packets.
var Handles = map[uint32]func(player *Player, pk packet.Packet){
	packet.IDRequestChunkRadius: func(player *Player, pk packet.Packet) {
		player.chunkradius = pk.(*packet.RequestChunkRadius).ChunkRadius
	},
}