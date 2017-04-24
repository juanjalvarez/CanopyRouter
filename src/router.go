package gorouter

import (
	"net/http"
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
	children []*Route
	isWildcard bool
	handlers RouteHandlers
}

func (r *Route) ToHandler() (func (http.ResponseWriter, *http.Request)) {
	return func (rw http.ResponseWriter, req *http.Request) {
		r.ResolveRoute(req.URL.Path, &rw, req)
	}
}

func NewRouter() *Route {
	r := new(Route)
	r.isRoot = true
	r.name = "_root_"
	r.parent = nil
	r.children = *(new([]*Route))
	r.isWildcard = false
	r.handlers = *(new(RouteHandlers))
	return r
}

func (r *Route) ResolveRoute(path string, rw *http.ResponseWriter, req *http.Request) {
	println("REQUEST MADE")
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

func (r *Route) Fork(name string) *Route {
	child := NewRouter()
	child.isRoot = false
	child.name = name
	child.parent = r
	r.children = append(r.children, child)
	child.isWildcard = false
	return child
}

func (r *Route) Wildcard(name string) *Route {
	fork := r.Fork(":" + name)
	fork.isWildcard = true
	return fork
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
