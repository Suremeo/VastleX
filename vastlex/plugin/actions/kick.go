package actions

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type Kick struct {
	*protobuf.KickAction
}

func (i *Kick) ID() int16 {
	return IDKick
}
