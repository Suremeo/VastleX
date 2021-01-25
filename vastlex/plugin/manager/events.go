package manager

import "github.com/VastleLLC/VastleX/vastlex/interfaces"

// HandleAddPlayer is called when a player enters the proxy.
var HandleAddPlayer func(player interfaces.Player)

// HandleRemovePlayer is called when a player leaves the proxy.
var HandleRemovePlayer func(xuid string)
