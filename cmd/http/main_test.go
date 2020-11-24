package main

/*
import (
	"net/http"

	altrudos "github.com/altrudos/api"
	"github.com/jmoiron/sqlx"
)

var testServerRunning = false
var testServer *http.Server

// This sets up our HTTP server to run on a test port
// We will test it by making real HTTP requests to localhost:{TEST_PORT}
func MustSetupTestServer() *altrudos.Config {
	conf := MustGetTestConfig()
	if testServerRunning {
		return conf
	}
	testServerRunning = true
	server := NewServer(conf)
	testServer = &http.Server{
		Addr:    ":" + conf.Port,
		Handler: server,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	return conf
}

func MustSetupTestServerDB() (*Config, *sqlx.DB) {
	conf := MustSetupTestServer()
	db := conf.MustGetDB()
	return conf, db
}
*/
