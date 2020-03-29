package main

import (
	"github.com/jmoiron/sqlx"

	. "github.com/charityhonor/ch-api"
)

var CharityRoutes = []*Route{
	NewGET("/charities", getCharities),
	NewGET("/charity/:id", getById("id", "Charity", getCharity)),
}

func getCharities(c *RouteContext) {
	cond := GetDefaultCondFromQuery(c.Query)
	//TODO: add search params here into cond

	var xs []*Drive
	defaultGetAll(c, "Charities", ViewCharities, &xs, cond)
}

func getCharity(db sqlx.Queryer, id string) (interface{}, error) {
	return GetCharityById(db, id)
}