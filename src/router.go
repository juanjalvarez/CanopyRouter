package canopy

import (
	"net/http"
	"strings"
)

const (
	GET = 0
	HEAD = 1
	POST = 2
	PUT = 3
	DELETE = 4
	CONNECT = 5
	OPTIONS = 6
	TRACE = 7
	PATCH = 8
	METHOD_COUNT = PATCH + 1
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

func NewRouter() *Route {
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
	child := NewRouter()
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
		handler := r.handlers[method]
		if handler != nil {
			handler(&rw, req, wildcards)
		} else {
			rw.WriteHeader(405)
		}
	}
}

func (r *Route) resolveRoute(stack []string, idx int, wildcards *Wildcards) *Route {
	cur := stack[idx]
	child := r.children[cur]
	if child == nil {
		if r.wildcard != nil {
			(*wildcards)[r.wildcard.name]=stack[idx]
			if idx == len(stack) - 1 {
				return r.wildcard
			} else {
				return r.wildcard.resolveRoute(stack, idx + 1, wildcards)
			}
		} else {
			return nil
		}
	} else {
		if idx == len(stack) - 1 {
			return r
		} else {
			return child.resolveRoute(stack, idx + 1, wildcards)
		}
	}
}

func methodCode(method string) int {
	switch (method) {
		case "GET": return GET
		case "HEAD" : return HEAD
		case "POST": return POST
		case "PUT": return PUT
		case "DELETE": return DELETE
		case "CONNECT": return CONNECT
		case "OPTIONS": return OPTIONS
		case "TRACE": return TRACE
		case "PATCH": return PATCH
		default: return -1
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
