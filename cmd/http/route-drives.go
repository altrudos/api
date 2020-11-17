package main

import (
	"net/http"

	. "github.com/altrudos/api"
	"github.com/cyc-ttn/gorouter"
)

var DriveRoutes = []*gorouter.Route{
	NewGET("/drives", getDrives),
	NewGET("/drives/top/:range", getTopDrives),
	NewGET("/drive/:uri", getDrive),
	NewPOST("/drive", createDrive),
	NewPOST("/drive/:id/donate", createDriveDonation),
}

var DriveColMap = map[string]string{
	"total": "usd_amount_total",
}

func getTopDrives(c *RouteContext) {
	drives, err := GetTopDrives(c.DB)
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, M{
		"Drives": drives,
	})
}

func getDrives(c *RouteContext) {
	cond := GetDefaultCondFromQuery(c.Query)
	cond.OrderBys = GetSortFromQueryWithDefault(c.Query, DriveColMap, []string{"-total"})

	var xs []*Drive
	defaultGetAll(c, "Drives", ViewDrives, &xs, cond)
}

func getDrive(c *RouteContext) {
	uri := c.Params["uri"]
	if c.HandledMissingParam(uri) {
		return
	}
	drive, err := GetDriveByUri(c.DB, uri)
	if c.HandledError(err) {
		return
	}
	if drive == nil {
		c.HandleError(ErrNotFound)
		return
	}
	topDonations, err := GetDriveTopDonations(c.DB, drive.Id, 10)
	if c.HandledError(err) {
		return
	}
	recentDonations, err := GetDriveRecentDonations(c.DB, drive.Id, 10)
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, M{
		"Drive":           drive,
		"TopDonations":    topDonations,
		"RecentDonations": recentDonations,
	})
}

func createDrive(c *RouteContext) {
	var nd NewDrive
	if err := c.ShouldBindJSON(&nd); c.HandledError(err) {
		return
	}

	tx, err := c.Services.DB.Beginx()
	if c.HandledError(err) {
		return
	}

	if err := nd.Process(tx); c.HandledError(err) {
		tx.Rollback()
		return
	}

	link, err := nd.Donation.GetDonationLink(c.Services.JG, c.Config.BaseUrl)
	if c.HandledError(err) {
		tx.Rollback()
		return
	}

	if err := tx.Commit(); c.HandledError(err) {
		return
	}

	c.JSON(http.StatusOK, M{
		"DonateLink": link,
		"Drive":      nd.Drive,
	})
}

func createDriveDonation(c *RouteContext) {
	var submitted struct {
		SubmittedDonation SubmittedDonation
	}
	if err := c.ShouldBindJSON(&submitted); c.HandledError(err) {
		return
	}
	var err error

	tx, err := c.Services.DB.Beginx()
	if c.HandledError(err) {
		return
	}

	var donation *Donation
	if donation, err = CreateDonation(tx, c.Params["id"], &submitted.SubmittedDonation); err != nil {
		c.HandleError(err)
		return
	}

	link, err := donation.GetDonationLink(c.Services.JG, c.Config.BaseUrl)
	if c.HandledError(err) {
		tx.Rollback()
		return
	}

	if err := tx.Commit(); c.HandledError(err) {
		return
	}

	c.JSON(http.StatusOK, M{
		"DonateLink": link,
	})
}
