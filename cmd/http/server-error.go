package main

import (
	"database/sql"
	"net/http"

	. "github.com/altrudos/api"
)

type HttpError interface {
	HttpCode() int
	HumanReadable() ErrorMessage
}

type ErrorMessage struct {
	Message  string
	RawError string
}

type HttpErrorMessage struct {
	Code     int
	Message  string
	RawError string
}

func (m HttpErrorMessage) Error() string {
	return m.Message
}

func (m HttpErrorMessage) HttpCode() int {
	if m.Code == 0 {
		return http.StatusInternalServerError
	}
	return m.Code
}

func (m HttpErrorMessage) HumanReadable() ErrorMessage {
	return ErrorMessage{
		Message:  m.Message,
		RawError: m.RawError,
	}
}

type RouteError struct {
	Method   string
	Path     string
	Code     int
	Message  string
	RawError error
}

func (e RouteError) Error() string {
	return e.RawError.Error()
}

func (e RouteError) HttpCode() int {
	return e.Code
}

func (e RouteError) HumanReadable() ErrorMessage {
	return ErrorMessage{
		Message:  e.Message,
		RawError: e.RawError.Error(),
	}
}

var ErrorMap = map[error]RouteError{
	sql.ErrNoRows: {
		Code:    http.StatusNotFound,
		Message: "Not found",
	},
	ErrNotFound: {
		Code:    http.StatusNotFound,
		Message: "Not found",
	},
	ErrSourceInvalidURL: {
		Code:    http.StatusBadRequest,
		Message: "Invalid source URL provided.",
	},
	ErrNoCharity: {
		Code:    http.StatusBadRequest,
		Message: "No charity provided.",
	},
	ErrInvalidCurrency: {
		Code:    http.StatusBadRequest,
		Message: "Invalid currency.",
	},
	ErrInvalidAmount: {
		Code:    http.StatusBadRequest,
		Message: "Invalid donation amount.",
	},
	ErrNegativeAmount: {
		Code:    http.StatusBadRequest,
		Message: "Donation amount can't be negative.",
	},
	ErrNilDonation: {
		Code:    http.StatusBadRequest,
		Message: "No donation submitted.",
	},
}

func (c *RouteContext) HandledError(err error) bool {
	if err == nil {
		return false
	}
	c.HandleError(err)
	return true
}

func (c *RouteContext) HandleError(err error) {
	conv := c.ConvertError(err)
	c.JSON(conv.HttpCode(), conv.HumanReadable())
}

func (c *RouteContext) HandledMissingParam(param string) bool {
	if param == "" {
		c.HandleError(ErrNotFound)
		return true
	}
	return false
}

func (c *RouteContext) ConvertError(err error) RouteError {
	if v, ok := ErrorMap[err]; ok {
		v.Method = c.Method
		v.Path = c.Path
		v.RawError = err
		return v
	}
	if v, ok := err.(HttpError); ok {
		return RouteError{
			Method:   c.Method,
			Path:     c.Path,
			Code:     v.HttpCode(),
			Message:  v.HumanReadable().Message,
			RawError: err,
		}
	}
	return RouteError{
		Method:   c.Method,
		Path:     c.Path,
		Code:     http.StatusInternalServerError,
		Message:  err.Error(),
		RawError: err,
	}
}
