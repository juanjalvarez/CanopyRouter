package main

import (
	"../src"
	"net/http"
	"fmt"
)

func test(rw *http.ResponseWriter, req *http.Request, w canopy.Wildcards) {

}

func main(){
	a := canopy.NewRouter()
	b := a.Fork("user")
	c := b.Wildcard("username")
	d := c.Fork("status")
	d.RegisterHandler(canopy.GET, test)
	println()
	a.Iterate(func (route *canopy.Route) {
		path := route.String()
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
	http.HandleFunc("/", a.ToHandler())
	http.ListenAndServe(":8080", nil)
}
