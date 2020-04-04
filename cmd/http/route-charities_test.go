package main

import (
	"net/http"
	"testing"

	"github.com/monstercat/golib/expectm"
)

const (
	CharityId = "9d0b23cd-657b-4cc4-8258-a8cabb1f6847"
)

func TestGetCharities(t *testing.T) {
	ts, _ := MustGetTestServer(GetFeaturedCharitiesRoute)

	resp, err := CallJson(ts, http.MethodGet, "/charities", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Charities.Data.#": 2,
		"Charities.Total":     2,
		"Charities.Limit":     50,
		"Charities.Offset":    0,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGetCharity(t *testing.T) {
	ts, _ := MustGetTestServer(GetCharityRoute)
	resp, err := CallJson(ts, http.MethodGet, "/charity/"+CharityId, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Charity.Id":  CharityId,
	}); err != nil {
		t.Fatal(err)
	}

}
