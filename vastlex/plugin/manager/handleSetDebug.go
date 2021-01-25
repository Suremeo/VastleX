package manager

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

func handleSetDebug(plugin *Plugin, a actions.Action) {
	action := a.(*actions.SetDebug)
	config.Config.Debug.Logging = action.Debug
	log.DefaultLogger.SetDebug(action.Debug)
	info := &actions.Info{
		InfoAction: &protobuf.InfoAction{},
	}
	WriteAll(info.EncodeConfig(config.Config))
}
