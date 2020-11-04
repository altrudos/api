package main

import (
	"net/url"
	"strconv"
	"strings"

	. "github.com/altrudos/api"
)

const (
	DefaultLimit = 50
)

type ColMap map[string]string

func (m ColMap) Get(s string) string {
	if m == nil {
		return s
	}
	c := m[s]
	if c == "" {
		return s
	}
	return c
}
func (m ColMap) Decode(s string) string {
	isDesc := s[0] == '-'
	if isDesc {
		s = s[1:]
	}
	g := m.Get(s)
	if isDesc {
		return g + " DESC"
	}
	return g + " ASC"
}

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

func GetSortFromQueryWithDefault(query url.Values, m ColMap, defaults []string) []string {
	sort := GetSortFromQuery(query, m)
	if len(sort) == 0 {
		return GetSortFromStringArray(defaults, m)
	}
	return sort
}

func GetSortFromQuery(query url.Values, m ColMap) []string {
	sort := query.Get("sort")
	if sort == "" {
		return nil
	}
	return GetSortFromStringArray(strings.Split(sort, ","), m)
}

func GetSortFromStringArray(sortCols []string, m ColMap) []string {
	sorts := make([]string, 0, len(sortCols))
	for _, c := range sortCols {
		sorts = append(sorts, m.Decode(c))
	}
	return sorts
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
