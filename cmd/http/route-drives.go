package main

import (
	"net/http"

	. "github.com/charityhonor/ch-api"
	"github.com/cyc-ttn/gorouter"
)

var DriveRoutes = []*gorouter.Route{
	NewGET("/drives", getDrives),
	NewGET("/drives/top/:range", getTopDrives),
	NewGET("/drive/:uri", getDrive),
	NewPOST("/drive", createDrive),
}

var DriveColMap = map[string]string{
	"total": "final_amount_total",
	"max":   "final_amount_max",
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

	link, err := nd.Donation.GetDonationLink(c.Services.JG)
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
