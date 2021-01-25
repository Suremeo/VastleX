package actions

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

// Init is sent from the plugin to server in order to initialize the plugin.
type Init struct {
	*protobuf.InitAction
}

func (i *Init) ID() int16 {
	return IDInit
}
