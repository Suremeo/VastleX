package vastlex

type Proxy interface {
	Motd() string
	SetMotd(string)
	Players() map[string]Player
	GetPlayer(string) (Player, error)
}
