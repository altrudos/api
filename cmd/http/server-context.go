package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	. "github.com/charityhonor/ch-api"
)

type RouteContext struct {
	*Services
	W http.ResponseWriter
	R *http.Request

	UserId string // LoggedOn User
	Params  map[string]string
	Query   url.Values
}

func (c *RouteContext) HandledError(err error) bool {
	if err == nil {
		return false
	}
	c.HandleError(err)
	return true
}

func (c *RouteContext) HandleError(err error) {
	//TODO: implement
}


func (c *RouteContext) Status(status int) {
	c.W.WriteHeader(status)
}

func (c *RouteContext) JSON(status int, v interface{}) error {
	c.Status(status)
	c.W.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(c.W).Encode(v)
}

func (c *RouteContext) String(status int, format string, data ...interface{}) {
	c.Status(status)
	c.W.Header().Add("Content-Type", "text/plain; charset=utf-8")
	if _, err := io.WriteString(c.W, fmt.Sprintf(format, data...)); err != nil {
		log.Println(err)
	}
}

func (c *RouteContext) ShouldBindJSON(v interface{}) error {
	return json.NewDecoder(c.R.Body).Decode(v)
}