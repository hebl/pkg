package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hebl/pkg/httputils"
	"github.com/sirupsen/logrus"
)

//Mux create mux
// func (rt *Router) Mux() *mux.Router {
// 	m := mux.NewRouter()
// 	var s *mux.Router
// 	if rt.PathPrefix == "" {
// 		s = m
// 	} else {
// 		s = m.PathPrefix(rt.PathPrefix).Subrouter()
// 	}
// 	for _, r := range rt.Routes {
// 		f := makeHTTPHandler(r.Handler)
// 		s.Path(r.Path).Methods(r.Method).Handler(f)
// 	}
// 	m.NotFoundHandler = notFoundHandler{}
// 	m.MethodNotAllowedHandler = methodNotAllowedHandler{}
// 	return m
// }

func MakeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		handlerFunc := handler
		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(ctx, w, r, vars); err != nil {
			statusCode := httputils.GetHTTPErrorStatusCode(err)
			if statusCode >= 500 {
				logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			}
			httputils.MakeErrorHandler(err)(w, r)
		}
	}
}
