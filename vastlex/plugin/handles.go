package plugin

import "github.com/VastleLLC/VastleX/vastlex/plugin/actions"

var handles = map[int16]func(a actions.Action){
	actions.IDInfo:  handleInfo,
	actions.IDEvent: handleEvent,
}
