package actions

import "github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"

type SetDebug struct {
	*protobuf.SetDebugAction
}

func (i *SetDebug) ID() int16 {
	return IDSetDebug
}
