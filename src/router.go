package canopy

import (
	"net/http"
	"strings"
)

type HttpHandler func (rw http.ResponseWriter, req *http.Request)

type Router struct {
	Root *Route
	errHandler *[600]HttpHandler
}

func NewRouter() *Router {
	m := new([600]HttpHandler)
	return &Router{
		Root: newRoute(),
		errHandler: m,
	}
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
		r.Error(404, rw, req)
	} else {
		method := methodCode(req.Method)
		handler := params.Route.handlers[method]
		if handler != nil {
			handler(&rw, req, params)
		} else {
			r.Error(405, rw, req)
		}
	}
}

func (r *Router) OnError(code int, handler HttpHandler) {
	(*r.errHandler)[code] = handler
}

func (r *Router) Error(code int, rw http.ResponseWriter, req *http.Request) {
	handler := (*r.errHandler)[code]
	if handler == nil {
		rw.WriteHeader(code)
	} else {
		handler(rw, req)
	}
}
