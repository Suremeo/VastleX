package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
)

func handleInfo(a actions.Action) {
	action := a.(*actions.Info)
	config.Config = action.DecodeConfig()
	if pluginVastleXinit.handleEvents && HandleAddPlayer != nil {
		HandleInfoUpdate()
	}
}
