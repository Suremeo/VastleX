package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions"
	"github.com/VastleLLC/VastleX/vastlex/plugin/actions/internalevents"
)

func handleEvent(a actions.Action) {
	action := a.(*actions.Event)
	evt, exists := action.Decode()
	if exists {
		switch evt := evt.(type) {
		case *internalevents.AddPlayer:
			p := &player{identity: evt.IdentityData}
			pluginVastleXinit.players.Store(evt.IdentityData.XUID, p)
			if pluginVastleXinit.handleEvents && HandleAddPlayer != nil {
				HandleAddPlayer(p)
			}
			break
		case *internalevents.RemovePlayer:
			pluginVastleXinit.players.Delete(evt.XUID)
			if pluginVastleXinit.handleEvents && HandleRemovePlayer != nil {
				HandleRemovePlayer(evt.XUID)
			}
			break
		}
	}
}
