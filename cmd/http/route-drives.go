package main

import (
	"github.com/jmoiron/sqlx"

	. "github.com/charityhonor/ch-api"
)

var DriveRoutes = []*Route{
	NewGET("/drives", getDrives),
	NewGET("/drive/:id", getById("id", "Drive", getDrive)),
//	NewAuthedPOST("/drive", createDrive),
//	NewAuthedPOST("/drive/:id", updateDrive),
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