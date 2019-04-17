package router

import (
	"net/http"

	"github.com/hebl/pkg/httputils"
)

//404
type notFoundHandler struct {
}

func (n notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := httputils.ErrorResponse{
		Code:    http.StatusNotFound,
		Message: r.URL.Path,
	}
	httputils.WriteJSON(w, http.StatusNotFound, v)
}

//405
type methodNotAllowedHandler struct{}

func (m methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := httputils.ErrorResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: r.URL.Path,
	}
	httputils.WriteJSON(w, http.StatusMethodNotAllowed, v)
}
