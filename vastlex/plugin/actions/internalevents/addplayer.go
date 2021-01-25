package internalevents

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
)

type AddPlayer struct {
	IdentityData login.IdentityData
}

func (i *AddPlayer) ID() protobuf.EventActionId {
	return protobuf.EventAction_AddPlayer
}
