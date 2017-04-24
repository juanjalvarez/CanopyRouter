package main

import (
	"../src"
)

func main(){
	a := gorouter.NewRouter()
	b := a.Fork("user")
	c := b.Wildcard("username")
	d := c.Fork("status")
	println(a.String())
	println(b.String())
	println(c.String())
	println(d.String())
}
