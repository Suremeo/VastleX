package plugin

import (
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type player struct {
	identity login.IdentityData
}

func (p *player) Identity() login.IdentityData {
	return p.identity
}

func (p *player) Send(ip string, port int) error {
	err := WriteAction(&actions.Send{SendAction: &protobuf.SendAction{
		Xuid: p.identity.XUID,
		Ip:   ip,
		Port: int32(port),
	}})
	log.DefaultLogger.Error(err)
	return err
}

func (p *player) Message(message string) error {
	err := p.WritePacket(&packet.Text{Message: message})
	log.DefaultLogger.Error(err)
	return err
}

func (p *player) Kick(reason ...string) {
	err := WriteAction(&actions.Kick{KickAction: &protobuf.KickAction{
		Xuid:   p.identity.XUID,
		Reason: reason,
	}})
	log.DefaultLogger.Error(err)
}

func (p *player) WritePacket(packet packet.Packet) error {
	err := WriteAction(actions.WritePacket{}.New(p.identity.XUID).Encode(packet))
	log.DefaultLogger.Error(err)
	return err
}
