package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/cyc-ttn/gorouter"

	"github.com/monstercat/golib/expectm"

	. "github.com/altrudos/api"
)

var (
	TestConfigPath = os.Getenv("TESTCONFIG")
)

func MustGetTestServer(routes ...*gorouter.Route) (*httptest.Server, *Services) {
	s := &Server{
		S:      MustGetTestServices(),
		R:      gorouter.NewRouter(),
		Config: MustGetTestConfig(),
	}
	if err := s.AddRoutes(routes); err != nil {
		panic(err)
	}
	ts := httptest.NewServer(s)
	return ts, s.S
}

func MustGetTestServices() *Services {
	if TestConfigPath == "" {
		TestConfigPath = "./config.example.toml"
	}
	s, err := GetConfigServices(TestConfigPath)
	if err != nil {
		panic(err)
	}
	return s
}

func CallJson(ts *httptest.Server, method, path string, data interface{}) (*http.Response, error) {
	byt, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return ts.Client().Do(req)
}

func CheckResponseBody(r io.Reader, m *expectm.ExpectedM) error {
	byt, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return expectm.CheckJSONBytes(byt, m)
}
