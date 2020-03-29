package main

import (
	"github.com/jmoiron/sqlx"

	. "github.com/charityhonor/ch-api"
)

var (
	GetCharitiesRoute = NewGET("/charities", getCharities)
	GetCharityRoute   = NewGET("/charity/:id", getById("id", "Charity", getCharity))
)

var CharityRoutes = []*Route{
	GetCharitiesRoute,
	GetCharityRoute,
}

func getCharities(c *RouteContext) {
	cond := GetDefaultCondFromQuery(c.Query)
	//TODO: add search params here into cond

	var xs []*Charity
	defaultGetAll(c, "Charities", ViewCharities, &xs, cond)
}

func getCharity(db sqlx.Queryer, id string) (interface{}, error) {
	return GetCharityById(db, id)
}
