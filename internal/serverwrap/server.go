package serverwrap

import (
	"context"
	"errors"
	"github.com/followedwind/slackbot/internal/util"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type serverWrap struct {
	server *http.Server
	mux    *mux.Router
}

func NewServer(addr string) *serverWrap {
	sw := serverWrap{}
	sw.mux = mux.NewRouter()
	sw.server = &http.Server{
		Addr:         addr,
		Handler:      sw.mux,
		ReadTimeout:  10 * time.Second, //request読み込み
		WriteTimeout: 20 * time.Second, //response出力
		IdleTimeout:  10 * time.Second, //再利用(headerにkeep-aliveが含まれる場合)
	}

	return &sw
}

func (sw *serverWrap) AddHandle(path string, handler http.Handler) *mux.Route {
	return sw.mux.Handle(path, handler)
}

func (sw *serverWrap) Start() {
	go func() {
		util.InfoLog("Server listening")
		if serverErr := sw.server.ListenAndServe(); !errors.Is(serverErr, http.ErrServerClosed) {
			util.ErrorLog(serverErr.Error())
		}
	}()
}

func (sw *serverWrap) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := sw.server.Shutdown(ctx); err != nil {
		util.ErrorLog(err.Error())
	} else {
		util.InfoLog("Server Done")
	}
}
