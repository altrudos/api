package main

func NewGET(path string, h HandlerFunc) *Route {
	return NewRoute("GET", path, h)
}

func NewPOST(path string, h HandlerFunc) *Route {
	return NewRoute("POST", path, h)
}

func NewAuthedPOST(path string, h HandlerFunc) *Route {
	return NewPOST(path, func(c *RouteContext) {
		if c.HandledError( c.Authenticate() ) {
			return
		}
		h(c)
	})
}

func NewRoute(method, path string, h HandlerFunc) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		HandlerFunc: h,
	}
}