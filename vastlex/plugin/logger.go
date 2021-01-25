package plugin

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"os"
)

var _ log.Logger = &logger{}

type logger struct{}

func (l *logger) Info(message string, sources ...string) {
	_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Info, message, append([]string{Name}, sources...)))
}

func (l *logger) Debug(message string, sources ...string) {
	_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Debug, message, append([]string{Name}, sources...)))
}

func (l *logger) SetDebug(enabled bool) {
	_ = WriteAction(&actions.SetDebug{SetDebugAction: &protobuf.SetDebugAction{Debug: enabled}})
}

func (l *logger) Warn(message string, sources ...string) {
	_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Warn, message, append([]string{Name}, sources...)))
}

func (l *logger) Error(err error, sources ...string) {
	if err != nil {
		_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Error, err.Error(), append([]string{Name}, sources...)))
	}
}

func (l *logger) Fatal(err error, sources ...string) {
	if err != nil {
		_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Fatal, err.Error(), append([]string{Name}, sources...)))
		os.Exit(0)
	}
}

func (l *logger) Title(message string) {
	_ = WriteAction(actions.NewLogAction(protobuf.LogAction_Title, message, []string{}))
}
