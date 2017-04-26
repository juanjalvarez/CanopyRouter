package canopy

import (
	"net/http"
)

type Router struct {
	Root *Route
}

func NewRouter() *Router {
	r := new(Router)
	r.Root = newRoute()
	return r
}

func (r *Router) Handler() (func (http.ResponseWriter, *http.Request)) {
	return http.HandlerFunc( func (rw http.ResponseWriter, req *http.Request) {
		r.Root.parseRoute(rw, req)
	})
}
