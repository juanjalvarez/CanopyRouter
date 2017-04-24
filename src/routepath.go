package gorouter

type RoutePath struct {
	isRoot bool
	name string
	parent *RoutePath
	children []*RoutePath
	endpoints *EndPointMap
	isWildcard bool
	wildcard *RoutePath
}

func NewRouter() *RoutePath {
	rp := new(RoutePath)
	rp.isRoot = true
	rp.name = "_root_"
	rp.parent = nil
	rp.children = *new([]*RoutePath)
	rp.endpoints = new(EndPointMap)
	rp.isWildcard = false
	rp.wildcard = nil
	return rp
}

func (rp RoutePath) Fork(name string) *RoutePath {
	child := NewRouter()
	child.isRoot = false
	child.name = name
	child.parent = &rp
	rp.children = append(rp.children, child)
	child.isWildcard = false
	child.wildcard = nil
	return child
}

func (rp RoutePath) Wildcard() *RoutePath {
	fork := rp.Fork("*")
	fork.isWildcard = true
	rp.wildcard = fork
	return fork
}

func (rp RoutePath) String() string {
	if rp.isRoot {
		return "/"
	}else {
		return (*(rp.parent)).String() + rp.name + "/"
	}
}
