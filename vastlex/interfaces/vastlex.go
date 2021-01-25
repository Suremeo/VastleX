package interfaces

import "sync"

// VastleX is as the name implies, the VastleX instance.
type VastleX interface {
	Motd() string
	SetMotd(string)
	Players() *sync.Map
	GetPlayer(string) (Player, bool)
	Count() int
}
