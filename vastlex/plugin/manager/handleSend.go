package manager

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
)

func handleSend(plugin *Plugin, a actions.Action) {
	action := a.(*actions.Send)
	p, ok := vastlex.GetPlayer(action.Xuid)
	if ok {
		log.DefaultLogger.Error(p.Send(action.Ip, int(action.Port)), plugin.Name)
	}
}
