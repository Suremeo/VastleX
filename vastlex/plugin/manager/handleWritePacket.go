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
		pkf, ok := p.(*player.Player).Player.(*minecraft.Connection).Pool[uint32(action.Id)]
		var pk packet.Packet
		if !ok {
			pk = &packet.Unknown{PacketID: uint32(action.Id)}
		} else {
			pk = pkf()
		}
		err := json.Unmarshal(action.Data, pk)
		if err == nil {
			log.DefaultLogger.Error(p.WritePacket(pk), plugin.Name)
		}
		log.DefaultLogger.Error(err, plugin.Name)
	}
}
