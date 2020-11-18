package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex"
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
)

func Init(name string, version int) {
	Name = name
	Version = version
	vastlex.VastleX = &pluginVastleX{
		players: map[string]interfaces.Player{},
	}
}
