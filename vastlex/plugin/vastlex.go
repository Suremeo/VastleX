package plugin

import "github.com/VastleLLC/VastleX/vastlex/interfaces"

type pluginVastleX struct {
	players map[string]interfaces.Player
}

func (vastlex *pluginVastleX) Motd() string {
	panic("implement me")
}

func (vastlex *pluginVastleX) SetMotd(s string) {
	panic("implement me")
}

func (vastlex *pluginVastleX) Players() map[string]interfaces.Player {
	panic("implement me")
}

func (vastlex *pluginVastleX) GetPlayer(s string) (interfaces.Player, error) {
	panic("implement me")
}
