package main

import (
	"net/http"

	. "github.com/charityhonor/ch-api"
	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"
)

var DriveRoutes = []*gorouter.Route{
	NewGET("/drives", getDrives),
	NewGET("/drives/top/:range", getTopDrives),
	NewGET("/drive/:uri", getById("uri", "Drive", getDrive)),
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

func getDrive(db sqlx.Queryer, uri string) (interface{}, error) {
	drive, err := GetDriveByUri(db, uri)
	if err != nil {
		return nil, err
	}
	drive.Top10Donations, err = GetDriveTop10Donations(db, drive.Id)
	if err != nil {
		return nil, err
	}
	return drive, nil
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
