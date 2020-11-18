package interfaces

// VastleX is as the name implies, the VastleX instance.
type VastleX interface {
	Motd() string
	SetMotd(string)
	Players() map[string]Player
	GetPlayer(string) (Player, error)
}
