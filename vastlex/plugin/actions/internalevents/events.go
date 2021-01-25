package internalevents

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

type Event interface {
	ID() protobuf.EventActionId
}

var Pool = map[protobuf.EventActionId]func() Event{
	protobuf.EventAction_AddPlayer:    func() Event { return &AddPlayer{} },
	protobuf.EventAction_RemovePlayer: func() Event { return &RemovePlayer{} },
}
