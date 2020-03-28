package main

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrPathNotFound   = errors.New("path not found")
	ErrInvalidMatcher = errors.New("invalid matcher")
)

type RouterNode struct {
	Children []*RouterNode
	Matcher  RouteMatcher
	Route    map[string]*Route
}

func NewRouter() *RouterNode {
	return &RouterNode{
		Children: make([]*RouterNode, 0, 1),
		Matcher:  &RouteMatcherRoot{},
	}
}

func (r *RouterNode) AddRoute(route *Route) error {
	return r.add(route, route.Path)
}

// Adds a route to RouterNode
func (r *RouterNode) add(route *Route, path string) error {
	if r.Matcher == nil {
		return ErrInvalidMatcher
	}

	// Make sure that we can get the token
	// for the current matcher. "GetToken" should error
	// when the matcher doesn't fit.
	rem, ok := r.Matcher.TokenMatch(path, route)
	if !ok {
		return ErrInvalidMatcher
	}

	for _, c := range r.Children {
		// If a child node is able to add then the
		// route has been added. Otherwise, we need to add the route!
		if err := c.add(route, rem); err == nil {
			return nil
		}
	}

	// Create matchers until there are no matchers left to create!
	node := r
	for rem != "" {
		var matcher RouteMatcher
		var err error

		// Add the route to myself by splitting into tokens!
		matcher, rem, err = MatchPathToMatcher(rem, route)
		if err != nil {
			return err
		}

		newNode := &RouterNode{
			Children: make([]*RouterNode, 0, 1),
			Matcher:  matcher,
		}
		node.Children = append(node.Children, newNode)
		node = newNode
	}
	node.addLeaf(route)

	return nil
}

func (r *RouterNode) addLeaf(route *Route) {
	if r.Route == nil {
		r.Route = make(map[string]*Route)
	}
	r.Route[route.Method] = route
}

func (r *RouterNode) getLeaf(method string) *Route {
	if r.Route == nil {
		return nil
	}
	return r.Route[method]
}

type RouteParams []string

func (r *RouteParams) Add(param string) {
	*r = append(*r, param)
}

func (r *RouterNode) Match(method, path string, ctx *RouteContext) (*Route, error) {

	// If there is a # in the path, completely ignore it.
	hashIdx := strings.Index(path, "#")
	if hashIdx > -1 {
		path = path[:hashIdx]
	}

	// If there is a ? in the path, parse separately.
	queryIdx := strings.Index(path, "?")
	if queryIdx > -1 {
		//Ignore erroneous query strings.
		ctx.Query, _ = url.ParseQuery(path[queryIdx+1:])
		path = path[:queryIdx]
	}

	params := &RouteParams{}
	route, err := r.match(method, path, params)
	if err != nil {
		return nil, err
	}

	ctx.Params = make(map[string]string)
	for i, p := range *params {
		name := route.ParamNames[i]
		ctx.Params[name] = p
	}

	return route, nil
}

func (r *RouterNode) match(method, path string, params *RouteParams) (*Route, error) {
	if r.Matcher == nil {
		return nil, ErrInvalidMatcher
	}
	rem, ok := r.Matcher.Match(method, path, params)
	if !ok {
		return nil, ErrPathNotFound
	}

	if rem == "" {
		route := r.getLeaf(method)
		if route == nil {
			return nil, ErrPathNotFound
		}
		return route, nil
	}

	for _, c := range r.Children {
		r, err := c.match(method, rem, params)
		if err == ErrInvalidMatcher {
			return nil, ErrInvalidMatcher
		}
		if err == nil {
			return r, nil
		}
	}
	return nil, ErrPathNotFound
}