package chocomux

import (
	"net/http"
)

type Handler func(ctx Context)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(WithStd(w, r))
}
