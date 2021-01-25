package actions

import "github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"

// Log is sent from the plugin to log things to the console
type Log struct {
	*protobuf.LogAction
}

func (i *Log) ID() int16 {
	return IDLog
}

func NewLogAction(t protobuf.LogAction_Type, message string, sources []string) *Log {
	return &Log{&protobuf.LogAction{
		Type:    t,
		Message: message,
		Sources: sources,
	}}
}
