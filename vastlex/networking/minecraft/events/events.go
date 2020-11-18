package events

// Event is an event.
type Event interface {
	// ID returns the id of the event.
	ID() int

	// Handle handles the event using the function and arguments provided.
	Handle(function interface{}, args ...interface{})
}

const (
	IDClose = iota
	IDLogin
)
