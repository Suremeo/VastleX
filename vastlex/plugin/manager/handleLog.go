package manager

import (
	"errors"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

func handleLog(plugin *Plugin, a actions.Action) {
	action := a.(*actions.Log)
	switch action.Type {
	case protobuf.LogAction_Info:
		log.DefaultLogger.Info(action.Message, action.Sources...)
	case protobuf.LogAction_Debug:
		log.DefaultLogger.Debug(action.Message, action.Sources...)
	case protobuf.LogAction_Warn:
		log.DefaultLogger.Warn(action.Message, action.Sources...)
	case protobuf.LogAction_Error:
		log.DefaultLogger.Error(errors.New(action.Message), action.Sources...)
	case protobuf.LogAction_Fatal:
		log.DefaultLogger.Fatal(errors.New(action.Message), action.Sources...)
	case protobuf.LogAction_Title:
		log.DefaultLogger.Title(action.Message)
	}
}
