package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	. "github.com/charityhonor/ch-api"
)

var (
	GetCharitiesRoute         = NewGET("/charities", getCharities)
	GetFeaturedCharitiesRoute = NewGET("/charities/featured", getFeaturedCharities)
	GetCharityRoute           = NewGET("/charity/:id", getById("id", "Charity", getCharity))
)

var CharityRoutes = []*gorouter.Route{
	GetCharitiesRoute,
	GetFeaturedCharitiesRoute,
	GetCharityRoute,
}

var CharityColMap = map[string]string{
	"total": "final_amount_total",
	"max":   "final_amount_max",
}

func getCharities(c *RouteContext) {
	search := c.Query.Get("search")
	if search == "" {
		getFeaturedCharities(c)
		return
	}
	var found bool
	cacheItem, err := GetSearchItem(c.DB, search)
	if err == nil && cacheItem != nil {
		if !cacheItem.Expired() {
			found = true
		}else{
			_ = cacheItem.Delete(c.DB)
		}
	}
	if !found {
		resp, err := c.JG.SearchCharitiesWithLimit(search, 20)
		if c.HandledError(err) {
			return
		}
		charities := make([]*Charity, 0, resp.Count)
		ids := make(pq.Int64Array, 0, resp.Count)
		for _, d := range resp.Results {
			charities = append(charities, &Charity{
				Name:                d.Name,
				Description:         d.Description,
				JustGivingCharityId: d.Id,
				WebsiteUrl:          d.WebsiteUrl,
			})
			ids = append(ids, int64(d.Id))
		}
		if err := CreateSearchCache(c.DB, search, ids); err != nil  {
			log.Print("Could not create search cache")
		}
		go func() {
			for _, charity := range charities {
				if err := charity.Insert(c.DB); err != nil && err != ErrDuplicateJGCharityId {
					log.Print("Could not create charity for " + strconv.Itoa(charity.JustGivingCharityId))
				}
			}
		}()
		c.JSON(http.StatusOK, charities)
		return
	}

	charities, err := GetCharitiesByJGId(c.DB, cacheItem.Ids)
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, charities)
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
