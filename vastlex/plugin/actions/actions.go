package actions

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"github.com/golang/protobuf/proto"
)

type Action interface {
	proto.Message
	ID() int16
}

var Pool = map[int16]func() Action{
	IDInit:        func() Action { return &Init{InitAction: &protobuf.InitAction{}} },
	IDInfo:        func() Action { return &Info{InfoAction: &protobuf.InfoAction{}} },
	IDLog:         func() Action { return &Log{LogAction: &protobuf.LogAction{}} },
	IDSetDebug:    func() Action { return &SetDebug{SetDebugAction: &protobuf.SetDebugAction{}} },
	IDSetMotd:     func() Action { return &SetMotd{SetMotdAction: &protobuf.SetMotdAction{}} },
	IDEvent:       func() Action { return &Event{EventAction: &protobuf.EventAction{}} },
	IDSend:        func() Action { return &Send{SendAction: &protobuf.SendAction{}} },
	IDWritePacket: func() Action { return &WritePacket{WritePacketAction: &protobuf.WritePacketAction{}} },
	IDKick:        func() Action { return &Kick{KickAction: &protobuf.KickAction{}} },
}

const (
	IDInit = int16(iota)
	IDInfo
	IDLog
	IDSetDebug
	IDSetMotd
	IDEvent
	IDSend
	IDWritePacket
	IDKick
)
