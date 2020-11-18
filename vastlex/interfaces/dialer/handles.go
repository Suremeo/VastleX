package dialer

import (
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft/custompackets"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Handles contains handlers for certain packets.
var Handles = map[uint32]func(dialer *Dialer, pak packet.Packet){
	packet.IDStartGame: handleStartGame,
	packet.IDAddActor: handleAddActor,
	packet.IDAddPlayer: handleAddPlayer,
	packet.IDAddItemActor: handleAddItemActor,
	packet.IDRemoveActor: handleRemoveActor,
	custompackets.IDVastleXTransfer: handleVastleXTransfer,
}