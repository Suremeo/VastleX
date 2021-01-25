package plugin

import (
	"github.com/VastleLLC/VastleX/vastlex/interfaces"
)

var Version = int32(0)
var Name = ""

// Events

// HandleAddPlayer is called when a player enters the proxy.
var HandleAddPlayer func(player interfaces.Player)

// HandleRemovePlayer is called when a player leaves the proxy.
var HandleRemovePlayer func(xuid string)

// HandleInfoUpdate is called when the proxy information is updated.
var HandleInfoUpdate func()
