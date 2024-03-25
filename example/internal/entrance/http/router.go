package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Router() http.Handler {
	router := httprouter.New()
	router.GET("/",http.)
}
