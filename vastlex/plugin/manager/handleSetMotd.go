package manager

import (
	"github.com/VastleLLC/VastleX/vastlex/config"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/protobuf"
)

func handleSetMotd(plugin *Plugin, a actions.Action) {
	action := a.(*actions.SetMotd)
	config.Config.Minecraft.Motd = action.Motd
	info := &actions.Info{
		InfoAction: &protobuf.InfoAction{},
	}
	WriteAll(info.EncodeConfig(config.Config))
}
