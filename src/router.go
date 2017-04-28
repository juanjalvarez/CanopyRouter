package canopy

import (
	"net/http"
	"strings"
)

/**
 * HTTPHandler is an alias for the net/http standard func for handling http requests.
 * @type {struct}
 */
type HTTPHandler func(rw http.ResponseWriter, req *http.Request)

/**
 * Router is composed of the root route, its configuration and error handlers.
 * @type {struct}
 */
type Router struct {
	Root       *Route
	Config     *RouterConfig
	errHandler *[600]HTTPHandler
}

/**
 * Collection of values that determine the behaviour of a router.
 * @type {struct}
 */
type RouterConfig struct {
	SensitiveSlashes bool
}

/**
 * NewRouter creates a new router with the default root route.
 * @type {*RouterConfig} config Settings to create the router with.
 * @return {*Router}   A pointer to the newly created router.
 */
func NewRouter(config *RouterConfig) *Router {
	m := new([600]HTTPHandler)
	c := config
	if c == nil {
		c = &RouterConfig{
			SensitiveSlashes: true,
		}
	}
	return &Router{
		Root:       newRoute(),
		Config:     c,
		errHandler: m,
	}
}

/**
 * Handler returns an HTTPHandler func to be attached to an http server, this func will handle all of the routing and handling.
 * @param  {*Router} r The router to create an HTTPHandler func from.
 * @return {HTTPHandler}   A HandlerFunc according to the net/http package.
 */
func (r *Router) Handler() HTTPHandler {
	return HTTPHandler(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		r.solve(rw, req)
	}))
}

func (r *Router) solve(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	path := strings.Split(reqPath, "/")
	lo, hi := 0, len(path)-1
	if len(path[lo]) == 0 {
		lo++
	}
	context := &RouteContext{
		stack:    &path,
		idx:      0,
		endSlash: false,
	}
	if len(path[hi]) == 0 {
		hi--
		context.endSlash = true
	}
	path = path[lo : hi+1]
	params := r.Root.parse(r, context)
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

/**
 * OnError enables the router to use the given handler func in the case that the router receives the given error code.
 * @param  {*Router} r Router to be manipulated.
 * @param  {int} code Error code to listen for.
 * @param  {HTTPHandler} handler Handler func to be callbed when the given error code ocurred.
 */
func (r *Router) OnError(code int, handler HTTPHandler) {
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
