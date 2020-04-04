package main

import (
	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"

	. "github.com/charityhonor/ch-api"
)

var (
	GetCharitiesRoute         = NewGET("/charities", getCharities)
	GetFeaturedCharitiesRoute = NewGET("/charities/featured", getFeaturedCharities)
	GetCharityRoute           = NewGET("/charity/:id", getById("id", "Charity", getCharity))
)

var CharityRoutes = []*gorouter.Route{
	GetFeaturedCharitiesRoute,
	GetCharityRoute,
}

var CharityColMap = map[string]string{
	"total": "final_amount_total",
	"max":   "final_amount_max",
}

func getCharities(c *RouteContext) {
}

func getFeaturedCharities(c *RouteContext) {
	cond := GetDefaultCondFromQuery(c.Query)
	cond.OrderBys = GetSortFromQueryWithDefault(c.Query, CharityColMap, []string{"-total"})
	cond.OrderBys = append(cond.OrderBys, "feature_score DESC")

	var xs []*Charity
	defaultGetAll(c, "Charities", ViewFeaturedCharities, &xs, cond)
}

func getCharity(db sqlx.Queryer, id string) (interface{}, error) {
	charity, err := GetCharityById(db, id)
	if err != nil {
		return nil, err
	}
	charity.Top10Donations, err = GetCharityTop10Donations(db, id)
	if err != nil {
		return nil, err
	}
	return charity, nil
}
