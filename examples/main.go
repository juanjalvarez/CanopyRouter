package main

import (
	"../src"
)

func main(){
	a := gorouter.NewRouter()
	b := a.Fork("test")
	c := b.Wildcard()
	d := c.Fork("verify")
	println(a.String())
	println(b.String())
	println(c.String())
	println(d.String())
}
