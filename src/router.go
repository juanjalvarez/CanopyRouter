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
)

type RouteHandler func(rw http.ResponseWriter, req *http.Request)

type RouteHandlers [9]*RouteHandler

type Route struct {
	isRoot bool
	name string
	parent *Route
	children []*Route
	isWildcard bool
}

func NewRouter() *Route {
	rp := new(Route)
	rp.isRoot = true
	rp.name = "_root_"
	rp.parent = nil
	rp.children = *(new([]*Route))
	rp.isWildcard = false
	return rp
}

func (rp *Route) Fork(name string) *Route {
	child := NewRouter()
	child.isRoot = false
	child.name = name
	child.parent = rp
	rp.children = append(rp.children, child)
	child.isWildcard = false
	return child
}

func (rp *Route) Wildcard() *Route {
	fork := rp.Fork("*")
	fork.isWildcard = true
	return fork
}

func (rp *Route) String() string {
	if rp.isRoot {
		return "/"
	}else {
		name := rp.name
		if len(rp.children) > 0 {
			name = name + "/"
		}
		return (*(rp.parent)).String() + name
	}
}
