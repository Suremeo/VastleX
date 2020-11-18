package events

// Login represents the login event for when the login process is complete for a player.
type Login struct {}

// ...
func (*Login) ID() int {
	return IDLogin
}

// ...
func (*Login) Handle(function interface{}, args ...interface{}) {
	function.(func())()
}