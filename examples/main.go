package main

import (
	"../src"
	"net/http"
	"fmt"
)

func test(rw *http.ResponseWriter, req *http.Request, w gorouter.Wildcards) {

}

func main(){
	a := gorouter.NewRouter()
	b := a.Fork("user")
	c := b.Wildcard("username")
	d := c.Fork("status")
	d.RegisterHandler(gorouter.GET, test)
	println()
	a.Iterate(func (route *gorouter.Route) {
		path := route.String()
		var methods []string
		if route.HasMethod(gorouter.GET) {
			methods = append(methods, "GET")
		}
		if route.HasMethod(gorouter.POST) {
			methods = append(methods, "POST")
		}
		fmt.Printf("%v %s\n", methods, path)
	})
	println()
}
