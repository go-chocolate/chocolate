package handler

import "net/http"

const (
	HealthPath = "/__health__"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("success"))
}

type HealthHandler struct {
	Next http.Handler
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == HealthPath {
		Health(w, r)
		return
	}
	h.Next.ServeHTTP(w, r)
}
