package actions

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type Send struct {
	*protobuf.SendAction
}

func (i *Send) ID() int16 {
	return IDSend
}
