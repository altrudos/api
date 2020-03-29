package main

import (
	"net/http"

	. "github.com/charityhonor/ch-api"
)

var DriveRoutes = []Route{
	NewGET("/drives", getDrives),
}

func getDrives(c *RouteContext) {
	drives, err := GetDrives(c.Services.DB, nil)
	if c.HandledError(err) {
		return
	}
	c.JSON(http.StatusOK, M{
		"Drives": drives,
	})
}
