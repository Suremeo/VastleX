package dialer

import (
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/translators/entities"
	"sync"
)

var _ interfaces.Dialer = &Dialer{}

// ...
type Dialer struct {
	minecraft.Dialer
	entityId    int64
	uniqueId    int64
	player      interfaces.InternalPlayer
	mutex       sync.Mutex
	ready       chan struct{}
	entities    *sync.Map
	scoreboards *sync.Map
	leaving     bool
}

func (dialer *Dialer) SetLeaving(b bool) {
	dialer.leaving = b
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
		if dialer.player.State() == interfaces.StateDisconnected {
			break
		}
		pk, err := dialer.ReadPacket()
		if err != nil {
			break
		}
		if Handles[pk.ID()] != nil {
			Handles[pk.ID()](dialer, pk)
		}
		if entities.Pool[pk.ID()] != nil {
			entities.Pool[pk.ID()]().Translate(pk, int64(dialer.entityId), 1, int64(dialer.uniqueId), 1)
		}
		err = dialer.player.WritePacket(pk)
		if err != nil {
			dialer.player.Kick("We recieved something wierd from your client, please rejoin")
			break
		}
	}
	_ = dialer.Close()
}
