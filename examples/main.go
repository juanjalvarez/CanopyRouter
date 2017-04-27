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
	root.GET(func (rw *http.ResponseWriter, req *http.Request, w canopy.Wildcards) {
		(*rw).Write([]byte("Hello, world!"))
	})

	// http://localhost:8080/user
	user := root.Fork("user")

	// http://localhost:8080/user/:username
	username := user.Wildcard("username")

	// http://localhost:8080/user/:username/status
	status := username.Fork("status")

	// GET http://localhost:8080/user/:username/status
	status.GET(func (rw *http.ResponseWriter, req *http.Request, w canopy.Wildcards) {
		(*rw).Write([]byte("Hello, " + w[":username"] +"!"))
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
	http.HandleFunc("/", router.Handler())
	http.ListenAndServe(":8080", nil)

}
