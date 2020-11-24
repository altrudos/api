package main

import (
	"net/http"

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
