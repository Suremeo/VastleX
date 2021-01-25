package manager

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
)

func handleKick(plugin *Plugin, a actions.Action) {
	action := a.(*actions.Kick)
	p, ok := vastlex.GetPlayer(action.Xuid)
	if ok {
		p.Kick(action.Reason...)
	}
}
