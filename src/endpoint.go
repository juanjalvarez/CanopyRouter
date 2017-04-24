package gorouter

import (
	"net/http"
)

type EndPoint struct {
	name string
	handler func (w http.ResponseWriter, r *http.Request)
}

type EndPointMap map[string][]*EndPoint
