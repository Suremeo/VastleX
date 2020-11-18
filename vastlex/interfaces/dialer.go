package interfaces

import (
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/translators/blocks"
)

// Dialer is a connection to a remote server.
type Dialer interface {
	minecraft.Dialer
	EntityId() int64
	UniqueId() int64
	Player() InternalPlayer
	Blocks() *blocks.Store
}