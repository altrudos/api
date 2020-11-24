package main

import (
	"fmt"
	"net/http"

	grtest "github.com/Vindexus/go-router-test"
	"github.com/jmoiron/sqlx"

	. "github.com/altrudos/api"
	altrudos "github.com/altrudos/api"
)

var testServerRunning = false
var testServer *http.Server
var testConfig *altrudos.Config

// This sets up our HTTP server to run on a test port
// We will test it by making real HTTP requests to localhost:{TEST_PORT}
func MustSetupTestServer() *Config {
	conf := MustGetTestConfig()
	if testServerRunning {
		return conf
	}
	testServerRunning = true

	server, err := NewServer(conf)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	return conf
}

func MustSetupTestServerDB() (*Config, *sqlx.DB) {
	conf := MustSetupTestServer()
	service, err := conf.Connect()
	if err != nil {
		panic(err)
	}

	return conf, service.DB
}

func runTests(tests []*grtest.RouteTest) error {
	c := MustSetupTestServer()

	for i, v := range tests {
		tests[i].URL = fmt.Sprintf("http://localhost:%d%s", c.Port, v.Path)
	}

	return grtest.RunTests(tests)
}

func runTest(test *grtest.RouteTest) error {
	c := MustSetupTestServer()
	if test.URL == "" {
		test.URL = fmt.Sprintf("http://localhost:%d%s", c.Port, test.Path)
	}
	return test.Run()
}
