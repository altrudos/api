package main

type Route struct {
	Method      string
	Path        string
	HandlerFunc HandlerFunc

	ParamNames []string
}

func (r *Route) AddParamName(name string) {
	r.ParamNames = append(r.ParamNames, name)
}

type HandlerFunc func(ctx *RouteContext)
