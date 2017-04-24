package canopy

import (
	"net/http"
	"strings"
	"fmt"
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
	isWildcard bool
	handlers RouteHandlers
}

func NewRouter() *Route {
	r := new(Route)
	r.isRoot = true
	r.name = "_root_"
	r.parent = nil
	r.children = make(map[string]*Route)
	r.isWildcard = false
	r.handlers = *(new(RouteHandlers))
	return r
}

func (r *Route) Fork(name string) *Route {
	child := NewRouter()
	child.isRoot = false
	child.name = name
	child.parent = r
	r.children[name] = child
	child.isWildcard = false
	return child
}

func (r *Route) Wildcard(name string) *Route {
	fork := r.Fork(":" + name)
	fork.isWildcard = true
	return fork
}

func (r *Route) ToHandler() (func (http.ResponseWriter, *http.Request)) {
	return func (rw http.ResponseWriter, req *http.Request) {
		r.resolveRoute(&rw, req)
	}
}

func (r *Route) resolveRoute(rw *http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	path := strings.Split(reqPath, "/")
	fmt.Printf("%v\n", path)
	lo, hi := 0, len(path) - 1
	if len(path[lo]) == 0 {
		lo++
	}
	if len(path[hi]) == 0 {
		hi--
	}
	path = path[lo:hi + 1]
	for _, val := range(path) {
		fmt.Printf("'%s'\n", val)
	}
}

func (r *Route) findRoute(stack []string, idx int) {

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
