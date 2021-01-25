package minecraftevents

// Event is an event.
type Event interface {
	// ID returns the id of the event.
	ID() EventId

	// Handle handles the event using the function and arguments provided.
	Handle(function interface{}, args ...interface{})
}

type EventId int16

const (
	IDClose = EventId(iota)
	IDLogin
)
