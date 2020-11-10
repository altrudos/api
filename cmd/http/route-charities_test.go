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
	ts, _ := MustGetTestServer(CharityRoutes...)

	resp, err := CallJson(ts, http.MethodGet, "/charities", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error(respBody(resp.Body))
		t.Error("Should be status ok got", resp.StatusCode)
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Charities.Data.#": 2,
		"Charities.Total":  2,
		"Charities.Limit":  50,
		"Charities.Offset": 0,
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
		t.Log(resp.StatusCode)
		logBody(resp.Body, t)
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Charity.Id": CharityId,
	}); err != nil {
		t.Fatal(err)
	}

}
