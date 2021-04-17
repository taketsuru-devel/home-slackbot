package serverwrap

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type serverWrap struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewServer(addr string) *serverWrap {
	sw := serverWrap{}
	sw.mux = http.NewServeMux()
	sw.server = &http.Server{
		Addr:         addr,
		Handler:      sw.mux,
		ReadTimeout:  10 * time.Second, //request読み込み
		WriteTimeout: 20 * time.Second, //response出力
		IdleTimeout:  10 * time.Second, //再利用(headerにkeep-aliveが含まれる場合)
	}

	return &sw
}

func (sw *serverWrap) AddHandle(path string, handler http.Handler) {
	sw.mux.Handle(path, handler)
}

func (sw *serverWrap) Start() {
	go func() {
		fmt.Println("Server listening")
		fmt.Println(sw.server.ListenAndServe())
	}()
}

func (sw *serverWrap) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := sw.server.Shutdown(ctx); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Server Done")
	}
}
