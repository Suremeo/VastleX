package minecraftevents

// Close represents a close event which is sent when the connection is closed.
type Close struct{}

// ...
func (*Close) ID() EventId {
	return IDClose
}

// ...
func (*Close) Handle(function interface{}, args ...interface{}) {
	function.(func())()
}
