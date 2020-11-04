package main

import (
	"net/http"

	. "github.com/altrudos/api"
	"github.com/cyc-ttn/gorouter"
)

var DonationRoutes = []*gorouter.Route{
	NewGET("/donations/recent", getDonationsRecent),
	NewGET("/donations/check/:reference", checkDonation),
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

func checkDonation(c *RouteContext) {
	ref := c.Params["reference"]
	if c.HandledMissingParam("reference") {
		return
	}
	donation, err := GetDonationByReferenceCode(c.DB, ref)
	if c.HandledError(err) {
		return
	}
	if err := donation.CheckStatus(c.DB, c.JG); err != nil {
		c.HandleError(err)
		return
	}
	if donation.Status == DonationPending {
		c.String(http.StatusOK, "Could not verify donation on JustGiving. Make sure you finished donating on JustGiving.")
		return
	}
	if donation.Status == DonationRejected {
		c.String(http.StatusOK, "Donation took too long to verify. Try donating again.")
		return
	}
	drive, err := GetDriveById(c.DB, donation.DriveId)
	if c.HandledError(err) {
		return
	}
	destination := c.Config.WebsiteUrl + "/d/" + drive.Uri + "?donation=" + donation.ReferenceCode
	c.Redirect(http.StatusTemporaryRedirect, destination)
}
