package actions

import (
	"encoding/json"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type WritePacket struct {
	*protobuf.WritePacketAction
}

func (i *WritePacket) ID() int16 {
	return IDWritePacket
}

func (i WritePacket) New(xuid string) *WritePacket {
	i.WritePacketAction = &protobuf.WritePacketAction{Xuid: xuid}
	return &i
}

func (i *WritePacket) Encode(packet packet.Packet) *WritePacket {
	i.Id = int32(packet.ID())
	i.Data, _ = json.Marshal(packet)
	return i
}
