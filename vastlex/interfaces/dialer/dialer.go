package dialer

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/translators/blocks"
	"github.com/VastleLLC/VastleX/vastlex/translators/entities"
	"sync"
)

// ...
type Dialer struct {
	minecraft.Dialer
	entityId int64
	uniqueId int64
	player   interfaces.InternalPlayer
	mutex    sync.Mutex
	ready    chan struct{}
	blocks   *blocks.Store
	entities sync.Map
	leaving  bool
}

func (dialer *Dialer) EntityId() int64 {
	return dialer.entityId
}

func (dialer *Dialer) UniqueId() int64 {
	return dialer.uniqueId
}

func (dialer *Dialer) Player() interfaces.InternalPlayer {
	return dialer.player
}

func (dialer *Dialer) listenPackets() {
	for {
		pk, err := dialer.ReadPacket()
		if err != nil {
			break
		}
		if Handles[pk.ID()] != nil {
			Handles[pk.ID()](dialer, pk)
		}
		if config.Config.Debug.BlockTranslating {
			blocks.TranslatePacket(pk, dialer.player.Blocks(), dialer.blocks)
		}
		if entities.Pool[pk.ID()] != nil {
			entities.Pool[pk.ID()]().Translate(pk, int64(dialer.entityId), 1, int64(dialer.uniqueId), 1)
		}
		_ = dialer.player.WritePacket(pk)
	}
	_ = dialer.Close()
}

func (dialer *Dialer) Blocks() *blocks.Store {
	return dialer.blocks
}
