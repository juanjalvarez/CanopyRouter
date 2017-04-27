package canopy

import (
	"net/http"
	"strings"
)

type Router struct {
	Root *Route
}

func NewRouter() *Router {
	r := new(Router)
	r.Root = newRoute()
	return r
}

func (r *Router) Handler() HttpHandler {
	return HttpHandler(http.HandlerFunc( func (rw http.ResponseWriter, req *http.Request) {
		r.solve(rw, req)
	}))
}

func (r *Router) solve(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	path := strings.Split(reqPath, "/")
	lo, hi := 0, len(path) - 1
	if len(path[lo]) == 0 {
		lo++
	}
	if len(path[hi]) == 0 {
		hi--
	}
	path = path[lo:hi + 1]
	params := r.Root.parse(path, 0)
	if params == nil {
		rw.WriteHeader(404)
	} else {
		method := methodCode(req.Method)
		handler := params.Route.handlers[method]
		if handler != nil {
			handler(&rw, req, params)
		} else {
			rw.WriteHeader(405)
		}
	}
}
