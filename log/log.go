package log

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

// Checksum is the checksum of the proxy file.
var Checksum = func() string {
	f, err := os.Open(os.Args[0])
	FatalError("Error generating checksum", err)
	hash := sha1.New()
	if _, err := io.Copy(hash, f); err != nil {
		FatalError("Error generating checksum", err)
	}
	err = f.Close()
	FatalError("Error generating checksum", err)
	return hex.EncodeToString(hash.Sum(nil))
}()

// Diode is a non-blocking writer that ensures it doesn't cause lag with too much logging.
var Diode diode.Writer

// Debugging represents whether or not the proxy will log debug messages (Its auto updated to whatever is set in the configuration).
var Debugging = false

// TotalPlayers is temporary until events are added.
var TotalPlayers = 0

// init initializes the logger & diode writer so that the logging doesn't cause lag.
func init() {
	wr := diode.NewWriter(os.Stdout, 1000, 5*time.Millisecond, func(missed int) {
		fmt.Printf("Logger Dropped %d messages", missed)
	})
	Diode = wr
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wr})
}

// Success logs a success message.
func Success(msg string) {
	log.Info().Msg(msg)
}

// Info returns a logger which has a prefix that indicates a warn.
func Info() *zerolog.Event {
	return log.Info()
}

// Warn returns a logger which has a prefix that indicates a warn.
func Warn(msg string) {
	log.Warn().Msg(msg)
}

// Debug returns a debug logger if debug logging is enabled in the configuration.
func Debug() *zerolog.Event {
	if Debugging {
		return log.Debug()
	} else {
		return log.Debug().Discard()
	}
}

// Error logs an error to the console.
func Error(msg string, err error) {
	log.Debug().Err(err).Msg(msg)
}

// Err returns a error logger.
func Err() *zerolog.Event {
	return log.Error()
}

// Title updates the title of the terminal.
func Title(msg string) {
	_, _ = Diode.Write([]byte(fmt.Sprintf("\033]0;VastleX | (%v) | Players: %v\007", Checksum, msg)))
}

// FatalError logs an error to console exits the program.
func FatalError(str string, err error) {
	if err != nil {
		log.Err(err).Msg(str)
		log.Error().Msg("A fatal error has occured so the program has been exited.")
		os.Exit(0)
	}
}
