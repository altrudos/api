package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	. "github.com/altrudos/api"
	"github.com/altrudos/api/pkg/justgiving"
)

var (
	GetCharitiesRoute         = NewGET("/charities", getCharities)
	GetFeaturedCharitiesRoute = NewGET("/charities/featured", getFeaturedCharities)
	GetCharityRoute           = NewGET("/charity/:id", getById("id", "Charity", getCharity))
)

var CharityRoutes = []gorouter.Route{
	GetCharitiesRoute,
	GetFeaturedCharitiesRoute,
	GetCharityRoute,
}

var CharityColMap = map[string]string{
	"total": "usd_amount_total",
}

func getCharities(c *RouteContext) {
	search := c.Query.Get("search")
	if search == "" {
		getFeaturedCharities(c)
		return
	}
	limit := 20
	if c.Query.Get("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query.Get("limit"))
		if err != nil {
			limit = 20
		}
	}
	var found bool
	cacheItem, err := GetSearchItem(c.DB, search)
	if err == nil && cacheItem != nil {
		if !cacheItem.Expired() {
			found = true
		} else {
			_ = cacheItem.Delete(c.DB)
		}
	}
	if !found {
		resp, err := c.JG.SearchCharitiesWithLimit(search, limit)
		if err != nil {
			if err == justgiving.ErrGroupedResultsNot1 {
				c.JSON(http.StatusOK, M{
					"Charities": M{
						"Data": []string{},
					},
				})
				return
			}
			c.HandleError(err)
			return
		}
		charities := make([]*Charity, 0, resp.Count)
		ids := make(pq.Int64Array, 0, resp.Count)
		// TODO: Batch these into one insert call
		for _, d := range resp.Results {
			newCharity := &Charity{
				CountryCode:         d.CountryCode,
				Description:         d.Description,
				JustGivingCharityId: d.Id,
				LogoUrl:             d.LogoUrl,
				Subtext:             d.Subtext,
				Name:                d.Name,
				WebsiteUrl:          d.WebsiteUrl,
			}

			if err := newCharity.Insert(c.DB); err != nil && err != ErrDuplicateJGCharityId {
				c.HandleError(err)
				return
			} else if err == ErrDuplicateJGCharityId {
				newCharity, _ = GetCharityByJGId(c.DB, newCharity.JustGivingCharityId)
			}

			charities = append(charities, newCharity)
			ids = append(ids, int64(d.Id))
		}
		// TODO: Sort these results
		if err := CreateSearchCache(c.DB, search, ids); err != nil {
			log.Print("Could not create search cache")
		}
		SortCharities(charities)
		c.JSON(http.StatusOK, M{
			"Charities": M{
				"Data": charities,
			},
		})
		return
	}

	charities, err := GetCharitiesByJGId(c.DB, cacheItem.Ids)
	if c.HandledError(err) {
		return
	}
	SortCharities(charities)
	c.JSON(http.StatusOK, M{
		"Charities": M{
			"Data": charities,
		},
	})
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
