package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type httpWriter struct {
	Endpoint string
	Debug    bool
}

func (h *httpWriter) Write(p []byte) (n int, err error) {
	if h.Debug {
		fmt.Println(p)
	}
	request, err := http.NewRequest("POST", h.Endpoint, bytes.NewBuffer(p))
	if err != nil {
		if h.Debug {
			fmt.Println(err)
		}
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		if h.Debug {
			fmt.Println(err)
		}
		return 0, err
	}
	if h.Debug {
		body, _ := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		fmt.Println(body)
	}
	return len(p), nil
}

func InitLog(debug bool) {

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Logger = log.Output(&httpWriter{Endpoint: os.Getenv("LOGGER_ENDPOINT_URL"), Debug: debug})
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
