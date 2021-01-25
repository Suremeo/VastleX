package manager

import (
	"encoding/json"
	"github.com/VastleLLC/VastleX/vastlex/interfaces/player"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/networking/minecraft"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func handleWritePacket(plugin *Plugin, a actions.Action) {
	action := a.(*actions.WritePacket)
	p, ok := vastlex.GetPlayer(action.Xuid)
	if ok {
		pk, ok := p.(*player.Player).Player.(*minecraft.Connection).Pool[uint32(action.Id)]
		if !ok {
			// No packet with the ID. This may be a custom packet of some sorts.
			pk = &packet.Unknown{PacketID: uint32(action.Id)}
		}
		err := json.Unmarshal(action.Data, pk)
		if err == nil {
			log.DefaultLogger.Error(p.WritePacket(pk), plugin.Name)
		}
		log.DefaultLogger.Error(err, plugin.Name)
	}
}
