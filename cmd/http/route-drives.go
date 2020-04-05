package main

import (
	"net/http"

	. "github.com/charityhonor/ch-api"
	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"
)

var DriveRoutes = []*gorouter.Route{
	NewGET("/drives", getDrives),
	NewGET("/drive/:id", getById("id", "Drive", getDrive)),
	NewPOST("/drive", createDrive),
}

func getDrives(c *RouteContext) {
	cond := GetDefaultCondFromQuery(c.Query)
	//TODO: add search params here into cond

	var xs []*Drive
	defaultGetAll(c, "Drives", ViewDrives, &xs, cond)
}

func getDrive(db sqlx.Queryer, id string) (interface{}, error) {
	return GetDriveById(db, id)
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
