package main

func NewGET(path string, h HandlerFunc) Route {
	return Route{
		Method:      "GET",
		Path:        path,
		HandlerFunc: h,
	}
}
