package router

import (
	"context"
	"net/http"

	"github.com/hebl/pkg/httputils"
)

//Route route
type Router interface {
	Routes() []Route
}

type Route interface {
	// Handler returns the raw function to create the http handler.
	Handler() httputils.APIFunc
	// Method returns the http method that the route responds to.
	Method() string
	// Path returns the subpath where the route responds to.
	Path() string
}

type localRoute struct {
	method  string
	path    string
	handler httputils.APIFunc
}

// Handler returns the APIFunc to let the server wrap it in middlewares.
func (l localRoute) Handler() httputils.APIFunc {
	return l.handler
}

// Method returns the http method that the route responds to.
func (l localRoute) Method() string {
	return l.method
}

// Path returns the subpath where the route responds to.
func (l localRoute) Path() string {
	return l.path
}

// NewRoute initializes a new local route for the router.
func NewRoute(method, path string, handler httputils.APIFunc) Route {
	return localRoute{method, path, handler}
}

//Get GET method
func Get(path string, handler httputils.APIFunc) Route {
	return NewRoute("GET", path, handler)
}

//Post POST method
func Post(path string, handler httputils.APIFunc) Route {
	return NewRoute("POST", path, handler)
}

//Put PUT method
func Put(path string, handler httputils.APIFunc) Route {
	return NewRoute("PUT", path, handler)
}

//Delete DELETE method
func Delete(path string, handler httputils.APIFunc) Route {
	return NewRoute("DELETE", path, handler)
}

//Option OPTION method
func Option(path string, handler httputils.APIFunc) Route {
	return NewRoute("OPTION", path, handler)
}

//Head get method
func Head(path string, handler httputils.APIFunc) Route {
	return NewRoute("HEAD", path, handler)
}

func cancellableHandler(h httputils.APIFunc) httputils.APIFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		if notifier, ok := w.(http.CloseNotifier); ok {
			notify := notifier.CloseNotify()
			notifyCtx, cancel := context.WithCancel(ctx)
			finished := make(chan struct{})
			defer close(finished)
			ctx = notifyCtx
			go func() {
				select {
				case <-notify:
					cancel()
				case <-finished:
				}
			}()
		}
		return h(ctx, w, r, vars)
	}
}

// WithCancel makes new route which embeds http.CloseNotifier feature to
// context.Context of handler.
func WithCancel(r Route) Route {
	return NewRoute(r.Method(), r.Path(), cancellableHandler(r.Handler()))
}
