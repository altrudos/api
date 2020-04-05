package main

import (
	"net/http"
	"testing"

	"github.com/monstercat/golib/expectm"
)

var (
	DriveId = "3656cf1d-8826-404c-8f85-77f3e1f50464"
)

func TestGetDrives(t *testing.T) {
	ts, _ := MustGetTestServer(
		NewGET("/drives", getDrives),
	)

	resp, err := CallJson(ts, http.MethodGet, "/drives", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drives.Data.#": 1,
		"Drives.Total":  1,
		"Drives.Limit":  50,
		"Drives.Offset": 0,
		"Drives.Data.0.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGetDrive(t *testing.T) {
	ts, _ := MustGetTestServer(
		NewGET("/drive/:id", getById("id", "Drive", getDrive)),
	)
	resp, err := CallJson(ts, http.MethodGet, "/drive/" + DriveId, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drive.Id": DriveId,
		"Drive.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDrive(t *testing.T) {
	ts, _ := MustGetTestServer(
		DriveRoutes...
	)

	type test struct {
		Payload interface{}
		ExpectedM *expectm.ExpectedM
		ExpecttedStatus int
	}

	tests := []test{
		{
			Payload: nil,
			ExpectedM: &expectm.ExpectedM{
				"RawError": ErrInvalidSourceUrl,
			},
		},
	}

	resp, err := CallJson(ts, http.MethodPost, "/drive", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drive.Id": DriveId,
		"Drive.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}