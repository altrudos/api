package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/cyc-ttn/gorouter"

	. "github.com/altrudos/api"
)

var ErrNoConfig = errors.New("-config flag missing")

type Server struct {
	Port   int
	S      *Services
	Config *Config
	R      *gorouter.RouterNode
}

func (s *Server) ParseFlags() error {
	var confFile string
	flag.IntVar(&s.Port, "port", 8080, "Server port")
	flag.StringVar(&confFile, "config", "", "Configuration File")
	flag.Parse()

	if confFile == "" {
		return ErrNoConfig
	}

	config, err := ParseConfig(confFile)
	if err != nil {
		return err
	}
	services, err := config.Connect()
	if err != nil {
		return err
	}
	s.S = services
	s.Config = config
	return nil
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
	log.Printf("Starting server on port %d", s.Port)
	return http.ListenAndServe(":"+strconv.Itoa(s.Port), s)
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
