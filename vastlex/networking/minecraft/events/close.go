package events

// Close represents a close event which is sent when the connection is closed.
type Close struct {}

// ...
func (*Close) ID() int {
	return IDClose
}

// ...
func (*Close) Handle(function interface{}, args ...interface{}) {
	function.(func())()
}