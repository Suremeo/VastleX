package log

// Logger represents a logger for all messages.
type Logger interface {
	Info(message string)
	Debug(message string)
	SetDebug(enabled bool)
	Warn(message string)
	Error(err error)
	Fatal(err error)
	Title(message string)
}