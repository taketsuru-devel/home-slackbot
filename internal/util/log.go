package util

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
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

func DebugLog(msg string) {
	log.Debug().Msg(msg)
}

func InfoLog(msg string) {
	log.Info().Msg(msg)
}

func ErrorLog(msg string) {
	log.Error().Msg(msg)
}
