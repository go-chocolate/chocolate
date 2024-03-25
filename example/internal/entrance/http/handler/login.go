package handler

import (
	"net/http"

	"github.com/go-chocolate/chocolate/pkg/chocolate/chocohttp/chocomux"
	"github.com/julienschmidt/httprouter"
)

func Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	login(chocomux.WithStd(w, r))
}

func login(ctx chocomux.Context) {

}
