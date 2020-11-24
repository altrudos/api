package main

import "github.com/cyc-ttn/gorouter"

type HandlerFunc func(ctx *RouteContext)

func NewGET(path string, h HandlerFunc) gorouter.Route {
	return NewRoute("GET", path, h)
}

func NewPOST(path string, h HandlerFunc) gorouter.Route {
	return NewRoute("POST", path, h)
}

func NewAuthedPOST(path string, h HandlerFunc) gorouter.Route {
	return NewPOST(path, func(c *RouteContext) {
		if c.HandledError(c.Authenticate()) {
			return
		}
		h(c)
	})
}

func NewRoute(method string, path string, handler HandlerFunc) *gorouter.DefaultRoute {
	return &gorouter.DefaultRoute{
		Method: method,
		Path:   path,
		HandlerFunc: func(ctx interface{}) {
			handler(ctx.(*RouteContext))
		},
	}
}
