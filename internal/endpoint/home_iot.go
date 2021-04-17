package endpoint

import (
	"net/http"
)

type HomeIotEndpoint struct{}

func (h *HomeIotEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text")
	w.Write([]byte("success from lambda"))
}
