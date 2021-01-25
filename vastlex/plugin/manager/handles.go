package manager

import "github.com/VastleLLC/VastleX/vastlex/plugin/actions"

var handles = map[int16]func(plugin *Plugin, action actions.Action){
	actions.IDInit:        handleInit,
	actions.IDLog:         handleLog,
	actions.IDSetDebug:    handleSetDebug,
	actions.IDSetMotd:     handleSetMotd,
	actions.IDWritePacket: handleWritePacket,
	actions.IDSend:        handleSend,
	actions.IDKick:        handleKick,
}
