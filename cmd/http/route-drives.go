package main

import (
	"github.com/cyc-ttn/gorouter"
	"github.com/jmoiron/sqlx"
	"net/http"

	. "github.com/charityhonor/ch-api"
)

var DriveRoutes = []*gorouter.Route{
	NewGET("/drives", getDrives),
	NewGET("/drive/:id", getById("id", "Drive", getDrive)),
//	NewAuthedPOST("/drive", createDrive),
//	NewAuthedPOST("/drive/:id", updateDrive),
	NewPOST("/drive", func(c *RouteContext) {

		var payload NewDrive
		if err := c.ShouldBindJSON(&payload); c.HandledError(err) {
			return
		}

		c.JSON(http.StatusOK, M{
			"DonateLink": "https://www.justgiving.com",
			"Drive": Drive{
				Uri: GenerateUri(),
				SourceUrl: "http://www.reddit.com",
			},
		})
	}),
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

//func createDrive(c *RouteContext) {
//	var payload struct {
//		Drive Drive
//	}
//	if c.HandledError(c.ShouldBindJSON(&payload)) {
//		return
//	}
//	if c.HandledError(payload.Drive.Insert(c.DB)) {
//		return
//	}
//	c.Status(http.StatusNoContent)
//}
//
//func updateDrive(c *RouteContext) {
//	id := c.Params["id"]
//	if c.HandledMissingParam(id) {
//		return
//	}
//	drive, err := GetDriveById(c.DB, id)
//	if c.HandledError(err) {
//		return
//	}
//	var payload struct{
//		Drive map[string]interface{}
//	}
//	if c.HandledError(c.ShouldBindJSON(&payload)) {
//		return
//	}
//	// TODO: some type of apply from map
//	//  for partial updates.
//	log.Print(drive)
//
//	c.Status(http.StatusNoContent)
//}