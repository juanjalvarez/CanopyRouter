package main

import (
	"../src"
)

func main(){
	a := gorouter.NewRouter()
	println(a.String())
	b := a.Fork("test")
	println(b.String())
	c := b.Wildcard()
	println(c.String())
}
