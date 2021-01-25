package manager

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

func handleInit(plugin *Plugin, a actions.Action) {
	action := a.(*actions.Init)
	plugin.Name = action.Name
	plugin.Version = action.Version
	info := &actions.Info{
		InfoAction: &protobuf.InfoAction{},
	}
	err := plugin.WriteAction(info.EncodeConfig(config.Config))
	if err != nil {
		log.DefaultLogger.Warn("Error writing init packet, exiting plugin.")
		plugin.Close()
	} else {
		log.DefaultLogger.Debug("Loaded plugin: " + action.Name)
	}
}
