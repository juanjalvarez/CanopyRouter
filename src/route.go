package canopy

import (
	"net/http"
)

/**
 * Wildcards is an alias for a map with a key type of string and value type of string.
 */
type Wildcards map[string]string

/**
 * RouteHandler is an alias for the standard way of denoting route handlers with the canopy router. Defined as: func(rw *http.ResponseWriter, req *http.Request, rp *RouteParameters)
 */
type RouteHandler func(rw *http.ResponseWriter, req *http.Request, rp *RouteParameters)

/**
 * RouteHandlers is an array of RouteHandlers of size METHOD_COUNT, which is the amount of http request methods. This array contains a RouteHandler for each request method.
 */
type RouteHandlers [METHOD_COUNT]RouteHandler

/**
 * Route defines a URL endpoint, every route is a child of another route except for the root route. Routes can have multiple children, they can be directories/dynamic paths and they can also have a single wildcard child. All routes have a list of RouteHandlers, one for each request method.
 * @type {struct}
 */
type Route struct {
	name        string
	parent      *Route
	children    map[string]*Route
	wildcard    *Route
	isDirectory bool
	handlers    RouteHandlers
}

/**
 * RouteContext contains information relevant to a particular route being processed.
 * @type {struct}
 */
type RouteContext struct {
	stack    *[]string
	idx      int
	endSlash bool
}

/**
 * RouteParameters contains information parsed from a requested path, like wildcard values and requested paths for directories/dynamic paths.
 * @type {struct}
 */
type RouteParameters struct {
	Route         *Route
	Wildcards     Wildcards
	RequestedPath string
}

func (rc *RouteContext) next() *RouteContext {
	return &RouteContext{
		stack:    rc.stack,
		idx:      rc.idx + 1,
		endSlash: rc.endSlash,
	}
}

func newRoute() *Route {
	return &Route{
		name:        "",
		children:    make(map[string]*Route),
		wildcard:    nil,
		isDirectory: false,
		handlers:    *(new(RouteHandlers)),
	}
}

/**
 * Fork creates a child from the given route.
 * @param  {*Route} r Parent route for the new child route.
 * @param  {string} name The name of the new child route.
 * @return {*Route}   Newly created child route.
 */
func (r *Route) Fork(name string) *Route {
	child := newRoute()
	child.name = name
	child.parent = r
	r.children[name] = child
	return child
}

/**
 * Wildcard gives the given parent route its only wildcard child, this route is special as instead of being part of URLs statically, it can be a string of any web-friendly characters excluding the slash and backslash.
 * @param  {*Route} r Parent route for the new wildcard route.
 * @param  {string} name The key for the wildcard value.
 * @return {*Route}   Newly created wildcard route.
 */
func (r *Route) Wildcard(name string) *Route {
	fork := r.Fork(":" + name)
	r.wildcard = fork
	return fork
}

/**
 * Directory can enable or disable isDirectory status on a route. Being a directory means the route handler will receive a RequestedPath value in it's context, which is a directory style path that was extracted from the requested url.
 * @param  {*Route} r Route to start/stop being a directory route.
 * @param  {bool} b Boolean status for whether or not the route is a directory.
 */
func (r *Route) Directory(b bool) {
	r.isDirectory = b
}

func (r *Route) parse(router *Router, context *RouteContext) *RouteParameters {
	stack := *context.stack
	if len(stack) == context.idx || r.isDirectory {
		params := RouteParameters{
			Route:         r,
			Wildcards:     make(Wildcards),
			RequestedPath: "/",
		}
		if len(stack) == context.idx {
			if context.endSlash && router.Config.SensitiveSlashes {
				child := r.children[""]
				if child == nil {
					return nil
				} else {
					params.Route = child
					return &params
				}
			}
			return &params
		}
		if r.isDirectory {
			for index, val := range stack[context.idx:] {
				if index != 0 {
					params.RequestedPath += "/"
				}
				params.RequestedPath += val
			}
			if context.endSlash {
				params.RequestedPath += "/"
			}
			return &params
		}
	}
	cur := stack[context.idx]
	child := r.children[cur]
	if child == nil {
		if r.wildcard != nil {
			p := r.wildcard.parse(router, context.next())
			if p == nil {
				return nil
			}
			p.Wildcards[r.wildcard.name] = stack[context.idx]
			return p
		} else {
			return nil
		}
	} else {
		return child.parse(router, context.next())
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
	} else {
		name := r.name
		if r.parent.parent != nil {
			name = "/" + name
		}
		return (*(r.parent)).Path() + name
	}
}

func (r *Route) Iterate(handler func(route *Route)) {
	handler(r)
	for _, child := range r.children {
		(*child).Iterate(handler)
	}
}
