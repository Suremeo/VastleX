package log

// Logger represents a logger for all messages.
type Logger interface {
	Info(message string, source ...string)
	Debug(message string, source ...string)
	SetDebug(enabled bool)
	Warn(message string, source ...string)
	Error(err error, source ...string)
	Fatal(err error, source ...string)
	Title(message string)
}
