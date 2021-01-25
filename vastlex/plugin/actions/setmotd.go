package actions

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type SetMotd struct {
	*protobuf.SetMotdAction
}

func (i *SetMotd) ID() int16 {
	return IDSetMotd
}
