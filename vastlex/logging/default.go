package log

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"github.com/rs/zerolog/diode"
	"io"
	"os"
	"strings"
	"time"
)

// Checksum is the checksum of the proxy file.
var Checksum = func() string {
	f, err := os.Open(os.Args[0])
	DefaultLogger.Fatal(err)
	hash := sha1.New()
	if _, err := io.Copy(hash, f); err != nil {
		DefaultLogger.Fatal(err)
	}
	err = f.Close()
	DefaultLogger.Fatal(err)
	return hex.EncodeToString(hash.Sum(nil))
}()

// DefaultLogger is the default logger for the session.
var DefaultLogger Logger = &logger{
	diode: diode.NewWriter(os.Stdout, 1000, 5*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	}),
}

// logger is an internal type to represent the stderr logger.
type logger struct {
	debug bool
	diode diode.Writer
}

// players is the current player count of the proxy.
var players = 0

// message is a custom message that is placed in the title.
var message = ""

// MaxPlayerSuffix is the suffix for the player count as displayed in the console title.
var MaxPlayerSuffix = ""

// UpdatePlayerCount updates the player count in the console title.
func UpdatePlayerCount(count int) {
	players = count
	updateTitle()
}

// updateTitle updates the title with the variables defines above.
func updateTitle() {
	_, _ = DefaultLogger.(*logger).diode.Write([]byte(fmt.Sprintf("\033]0;VastleX (%v) | Players: %v%v | %v\007", Checksum, players, MaxPlayerSuffix, message)))
}

// Title sets the message for the window title.
func (l *logger) Title(m string) {
	message = m
	updateTitle()
}

// Info logs an INFO message to the console.
func (l *logger) Info(message string) {
	_, _ = l.diode.Write([]byte(addPrefixToNewLine(getTime() + " " + color.HiGreenString("INFO") + color.HiBlackString(" > "), message)))
}

// Debug logs a DEBUG message to the console only if debug logging is enabled in the proxy configuration.
func (l *logger) Debug(message string) {
	if l.debug {
		_, _ = l.diode.Write([]byte(addPrefixToNewLine(getTime() + " " + color.YellowString("DEBUG") + color.HiBlackString(" > "), message)))
	}
}

// SetDebug sets whether or not debug messages should be logged.
func (l *logger) SetDebug(enabled bool) {
	l.debug = enabled
}

// Warn logs a WARN message to the console.
func (l *logger) Warn(message string) {
	_, _ = l.diode.Write([]byte(addPrefixToNewLine(getTime() + " " + color.RedString("WARN") + color.HiBlackString(" > "), message)))
}

// Error logs a ERROR message to the console if the error provided is valid.
func (l *logger) Error(err error) {
	if err != nil {
		_, _ = l.diode.Write([]byte(addPrefixToNewLine(getTime() + " " + color.HiRedString("ERROR") + color.HiBlackString(" > "), err.Error())))
	}
}

// Fatal logs a FATAL message to the console and exits the program if the error provided is valid.
func (l *logger) Fatal(err error) {
	if err != nil {
		_, _ = l.diode.Write([]byte(addPrefixToNewLine(getTime() + " " + color.HiRedString("FATAL") + color.HiBlackString(" > "), err.Error())))
		time.Sleep(100*time.Millisecond)
		os.Exit(0)
	}
}

// addPrefixToNewLine adds the prefix provided to every single line of the provided body.
func addPrefixToNewLine(prefix, body string) string {
	var done string
	for _, line := range strings.Split(body, "\n") {
		done = done + prefix + color.WhiteString(line) + "\n"
	}
	return done
}

// getTime returns a formatted verison of the current time.
func getTime() string {
	return color.HiBlackString(time.Now().Format("15:04:05"))
}
