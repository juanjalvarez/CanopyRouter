package canopy

import (
	"net/http"
)

type Wildcards map[string]string

type RouteHandler func (rw *http.ResponseWriter, req *http.Request, rp *RouteParameters)

type RouteHandlers [METHOD_COUNT]RouteHandler

type Route struct {
	name string
	parent *Route
	children map[string]*Route
	wildcard *Route
	isDirectory bool
	handlers RouteHandlers
}

type RouteParameters struct {
	Route *Route
	Wildcards Wildcards
	RequestedPath string
}

func newRoute() *Route {
	return &Route{
		name: "_root_",
		children: make(map[string]*Route),
		wildcard: nil,
		isDirectory: false,
		handlers: *(new(RouteHandlers)),
	}
}

func (r *Route) Fork(name string) *Route {
	child := newRoute()
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

func (r *Route) Directory(b bool) {
	r.isDirectory = b
}

func (r *Route) parse(stack []string, idx int) *RouteParameters {
	params := RouteParameters{
		Route: r,
		Wildcards: make(Wildcards),
		RequestedPath: "",
	}
	if len(stack) == idx {
		return &params
	}
	if r.isDirectory {
		for _, val := range(stack[idx:]){
			params.RequestedPath = params.RequestedPath + "/" + val
		}
		return &params
	}
	cur := stack[idx]
	child := r.children[cur]
	if child == nil {
		if r.wildcard != nil {
			p := r.wildcard.parse(stack, idx + 1)
			p.Wildcards[r.wildcard.name]=stack[idx]
			return p
		} else {
			return nil
		}
	} else {
		return child.parse(stack, idx + 1)
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

func (r *Route) GET(handler RouteHandler) {
	r.handlers[GET] = handler
}

func (r *Route) POST(handler RouteHandler) {
	r.handlers[POST] = handler
}

func (r *Route) Path() string {
	if r.parent == nil {
		return "/"
	}else {
		name := r.name
		if len(r.children) > 0 {
			name = name + "/"
		}
		return (*(r.parent)).Path() + name
	}
}

func (r *Route) Iterate(handler func (route *Route)) {
	handler(r)
	for _, child := range(r.children) {
		(*child).Iterate(handler)
	}
}
