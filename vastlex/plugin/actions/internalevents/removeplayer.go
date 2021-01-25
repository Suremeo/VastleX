package internalevents

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type RemovePlayer struct {
	XUID string
}

func (i *RemovePlayer) ID() protobuf.EventActionId {
	return protobuf.EventAction_RemovePlayer
}
