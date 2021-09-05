package util

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
)

func InitLog(debug bool, pretty bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func DebugLog(msg string, addStack int) {
	log.Debug().Msg(decolateLog(msg, addStack))
}

func InfoLog(msg string, addStack int) {
	log.Info().Msg(decolateLog(msg, addStack))
}

func ErrorLog(msg string, addStack int) {
	log.Error().Msg(decolateLog(msg, addStack))
}

func decolateLog(msg string, addStack int) string {
	_, file, line, ok := runtime.Caller(2 + addStack) //this, **Log, caller
	if !ok {
		return msg
	} else {
		return fmt.Sprintf("file: %s, line: %d, msg: %s", file, line, msg)
	}
}
