package main

import (
	"../src"
	"net/http"
	"fmt"
)

func main(){

	router := canopy.NewRouter()

	// http://localhost:8080/
	root := router.Root

	// GET http://localhost:8080/
	root.GET(func (rw *http.ResponseWriter, req *http.Request, rp *canopy.RouteParameters) {
		(*rw).Write([]byte("Hello, world!"))
	})

	// Basic route example
	// http://localhost:8080/user
	user := root.Fork("user")

	// Wildcard example
	// http://localhost:8080/user/:username
	username := user.Wildcard("username")

	// Another basic route example
	// http://localhost:8080/user/:username/status
	status := username.Fork("status")

	// "Directory / dynamic path / catch-all" example
	//http://localhost:8080/user/:username/dir/%
	dir := username.Fork("dir")
	dir.Directory(true)
	dir.GET(func (rw *http.ResponseWriter, req *http.Request, rp *canopy.RouteParameters) {
		(*rw).Write([]byte("Requested Path: " + rp.RequestedPath))
	})

	// GET method for the status route
	// GET http://localhost:8080/user/:username/status
	status.GET(func (rw *http.ResponseWriter, req *http.Request, rp *canopy.RouteParameters) {
		(*rw).Write([]byte("Hello, " + rp.Wildcards[":username"] +"!"))
	})

	// Error 404 handler
	router.OnError(404, func (rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(req.URL.Path + " does not exist, 404."))
	})

	// Error 405 handler
	router.OnError(405, func (rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(req.URL.Path + " does not contain a controller for the " + req.Method + " method."))
	})

	// Prints all registered routes and their methods
	println()
	root.Iterate(func (route *canopy.Route) {
		path := route.Path()
		var methods []string
		if route.HasMethod(canopy.GET) {
			methods = append(methods, "GET")
		}
		if route.HasMethod(canopy.POST) {
			methods = append(methods, "POST")
		}
		fmt.Printf("%v %s\n", methods, path)
	})
	println()

	// Router.Handler() yields a single HandlerFunc compliant with the http package to be hooked into an http server.
	http.HandleFunc("/", router.Handler())
	http.ListenAndServe(":8080", nil)

}
