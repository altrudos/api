package main

import (
	"net/http"

	. "github.com/charityhonor/ch-api"
	"github.com/cyc-ttn/gorouter"
)

var DonationRoutes = []*gorouter.Route{
	NewGET("/donations/recent", getDonationsRecent),
}

func getDonationsRecent(c *RouteContext) {
	donations, err := GetDonationsRecent(c.DB, &DonationOperators{
		BaseOperator: &BaseOperator{
			Limit: 10,
		},
	})
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, M{
		"Donations": donations,
	})
}
