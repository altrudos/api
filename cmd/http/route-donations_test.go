package main

import (
	"net/http"
	"testing"

	"github.com/monstercat/golib/expectm"
)

func TestGetDonationsRecent(t *testing.T) {
	ts, _ := MustGetTestServer(DonationRoutes...)

	resp, err := CallJson(ts, http.MethodGet, "/donations/recent", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error(respBody(resp.Body))
		t.Error("Should be status ok got", resp.StatusCode)
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Donations.#":             3,
		"Donations.0.Drive.Uri":   "PrettyPinkMoon",
		"Donations.0.FinalAmount": 1332,
	}); err != nil {
		t.Fatal(err)
	}
}
