package canopy

const (
	GET = 0
	HEAD = 1
	POST = 2
	PUT = 3
	DELETE = 4
	CONNECT = 5
	OPTIONS = 6
	TRACE = 7
	PATCH = 8
	METHOD_COUNT = PATCH + 1
)

func methodCode(method string) int {
	switch (method) {
		case "GET": return GET
		case "HEAD" : return HEAD
		case "POST": return POST
		case "PUT": return PUT
		case "DELETE": return DELETE
		case "CONNECT": return CONNECT
		case "OPTIONS": return OPTIONS
		case "TRACE": return TRACE
		case "PATCH": return PATCH
		default: return -1
	}
}

func NewRouter() {
}
