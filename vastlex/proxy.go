package vastlex

// Proxy is as the name implies, a proxy.
type Proxy interface {
	Motd() string
	SetMotd(string)
	Players() map[string]Player
	GetPlayer(string) (Player, error)
}
