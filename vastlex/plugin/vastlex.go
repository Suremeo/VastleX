package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
	log "github.com/VastleLLC/VastleX/vastlex/logging"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
	"sync"
)

type pluginVastleX struct {
	motd         string
	players      *sync.Map
	handleEvents bool
}

func (vastlex *pluginVastleX) Motd() string {
	return vastlex.motd
}

func (vastlex *pluginVastleX) SetMotd(s string) {
	vastlex.motd = s
	log.DefaultLogger.Error(WriteAction(&actions.SetMotd{SetMotdAction: &protobuf.SetMotdAction{Motd: s}}))
}

func (vastlex *pluginVastleX) Players() *sync.Map {
	return vastlex.players
}

func (vastlex *pluginVastleX) GetPlayer(xuid string) (interfaces.Player, bool) {
	p, ok := vastlex.players.Load(xuid)
	if !ok {
		return nil, ok
	}
	pl, ok := p.(interfaces.Player)
	return pl, ok
}
