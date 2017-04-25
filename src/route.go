package canopy

import (
	"net/http"
	"strings"
)

type Wildcards map[string]string

type RouteHandler func(rw *http.ResponseWriter, req *http.Request, w Wildcards)

type RouteHandlers [METHOD_COUNT]RouteHandler

type Route struct {
	isRoot bool
	name string
	parent *Route
	children map[string]*Route
	wildcard *Route
	handlers RouteHandlers
}

func NewRoute() *Route {
	r := new(Route)
	r.isRoot = true
	r.name = "_root_"
	r.parent = nil
	r.children = make(map[string]*Route)
	r.wildcard = nil
	r.handlers = *(new(RouteHandlers))
	return r
}

func (r *Route) Fork(name string) *Route {
	child := NewRoute()
	child.isRoot = false
	child.name = name
	child.parent = r
	r.children[name] = child
	return child
}

func (r *Route) Wildcard(name string) *Route {
	fork := r.Fork(":" + name)
	r.wildcard = fork
	return fork
}

func (r *Route) ToHandler() (func (http.ResponseWriter, *http.Request)) {
	return func (rw http.ResponseWriter, req *http.Request) {
		r.parseRoute(rw, req)
	}
}

func (r *Route) parseRoute(rw http.ResponseWriter, req *http.Request) {
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
	wildcards := (make(Wildcards))
	route := r.resolveRoute(path, 0, &wildcards)
	if route == nil {
		rw.WriteHeader(404)
	} else {
		method := methodCode(req.Method)
		handler := route.handlers[method]
		if handler != nil {
			handler(&rw, req, wildcards)
		} else {
			rw.WriteHeader(405)
		}
	}
}

func (r *Route) resolveRoute(stack []string, idx int, wildcards *Wildcards) *Route {
	cur := stack[idx]
	isLast := idx == len(stack) - 1
	child := r.children[cur]
	if child == nil {
		if r.wildcard != nil {
			(*wildcards)[r.wildcard.name]=stack[idx]
			if isLast {
				return r.wildcard
			} else {
				return r.wildcard.resolveRoute(stack, idx + 1, wildcards)
			}
		} else {
			return nil
		}
	} else {
		if isLast {
			return child
		} else {
			return child.resolveRoute(stack, idx + 1, wildcards)
		}
	}
}

func (r *Route) HasMethod(method int) bool {
	if r.handlers[method] != nil {
		return true
	} else {
		return false
	}
}

func (r *Route) RegisterHandler(method int, handler RouteHandler) {
	r.handlers[method] = handler
}

func (r *Route) String() string {
	if r.isRoot {
		return "/"
	}else {
		name := r.name
		if len(r.children) > 0 {
			name = name + "/"
		}
		return (*(r.parent)).String() + name
	}
}

func (r *Route) Iterate(handler func (route *Route)) {
	handler(r)
	for _, child := range(r.children) {
		(*child).Iterate(handler)
	}
}
