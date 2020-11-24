package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cyc-ttn/gorouter"

	. "github.com/altrudos/api"
)

var (
	ErrBlankConfigPath = errors.New("config file path is blank")
	ErrNoPort          = errors.New("port in config is empty")
)

type Server struct {
	S      *Services
	Config *Config
	R      *gorouter.RouterNode
}

func NewServer(config *Config) (*Server, error) {
	services, err := config.Connect()
	if err != nil {
		return nil, err
	}
	s := &Server{}
	s.S = services
	s.Config = config

	s.AddRoutes(
		CharityRoutes,
		DonationRoutes,
		DriveRoutes,
	)

	return s, nil
}

func NewServerFromConfigFile(confFile string) (*Server, error) {
	if confFile == "" {
		return nil, ErrBlankConfigPath
	}

	config, err := ParseConfig(confFile)
	if err != nil {
		return nil, err
	}
	return NewServer(config)
}

func (s *Server) AddRoutes(rs ...[]gorouter.Route) error {
	if s.R == nil {
		s.R = gorouter.NewRouter()
	}
	for _, rr := range rs {
		for _, r := range rr {
			if err := s.R.AddRoute(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Server) Run() error {
	if s.Config.Port == 0 {
		return ErrNoPort
	}
	log.Printf("Starting server on port %d", s.Config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Config.Port), s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//CORS
	w.Header().Set("Access-Control-Allow-Origin", s.Config.WebsiteUrl)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	ctx := &gorouter.RouteContext{
		W:      w,
		R:      r,
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.Query(),
	}
	route, err := s.R.Match(r.Method, r.URL.Path, ctx)
	if err == gorouter.ErrPathNotFound || route == nil {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	route.GetHandler()(&RouteContext{
		Services:     s.S,
		RouteContext: *ctx,
		Config:       s.Config,
	})
}
