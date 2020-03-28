package main

import (
	"flag"
	"net/http"
	"strconv"

	. "github.com/charityhonor/ch-api"
)

type Server struct {
	Port int
	S    *Services
	R    *RouterNode
}

func (s *Server) ParseFlags() error {
	var confFile string
	flag.StringVar(&confFile, "config", "", "Configuration File")
	flag.Parse()

	services, err := GetConfigServices(confFile)
	if err != nil {
		return err
	}
	s.S = services
	return nil
}

func (s *Server) AddRoutes(rs ...[]Route) error {
	if s.R == nil {
		s.R = NewRouter()
	}
	for _, rr := range rs {
		for _, r := range rr {
			if err := s.R.AddRoute(&r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) Run() error {
	return http.ListenAndServe(":" + strconv.Itoa(s.Port), s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r * http.Request) {
	ctx := &RouteContext{
		Services: s.S,
		W: w,
		R: r,
	}
	route, err := s.R.Match(r.Method, r.URL.Path, ctx)
	if err == ErrPathNotFound || route == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	route.HandlerFunc(ctx)
}