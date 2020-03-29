package main

import (
	"net/url"
	"strconv"

	. "github.com/charityhonor/ch-api"
)

const (
	DefaultLimit = 50
)

// Returns a Cond to be passed into a query. This function only retrieves the values
// from the request. Any other details such as default sorts and search parameters need
// to be handled by the route.
//
// In particular, this function will retrieve the limit and offset.
func GetDefaultCondFromQuery(query url.Values) *Cond {
	c := GetCondFromQuery(query)
	c.DefaultLimit(DefaultLimit)
	return c
}

func GetCondFromQuery(query url.Values) *Cond {
	return &Cond{
		Limit:  GetLimitFromQuery(query),
		Offset: GetOffsetFromQuery(query),
	}
}

func GetLimitFromQuery(query url.Values) int {
	return GetIntFromQuery(query, "limit")
}

func GetOffsetFromQuery(query url.Values) int {
	return GetIntFromQuery(query, "offset")
}

func GetIntFromQuery(query url.Values, key string) int {
	str := query.Get(key)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}