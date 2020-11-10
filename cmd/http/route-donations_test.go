package main

import (
	"net/http"
	"testing"

	"github.com/altrudos/api/pkg/fixtures"

	. "github.com/altrudos/api"

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
		"Donations.0.DonorName":   "FordonGreeman",
		"Donations.0.CharityName": "The Demo Charity",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCheckDonation(t *testing.T) {
	ts, _ := MustGetTestServer(DonationRoutes...)
	services := MustGetTestServices()
	db := services.DB
	defer db.Close()

	donation, err := GetDonationById(db, fixtures.DonationId1)
	if err != nil {
		t.Fatal(err)
	}
	donation.Status = DonationPending
	saveDonation := func() {
		if err := donation.Save(db); err != nil {
			t.Fatal(err)
		}
	}

	res, err := CallJson(ts, http.MethodGet, "/donations/check/"+donation.ReferenceCode, nil)
	if err != nil {
		t.Fatal(err)
	}
	/*
		if resp.StatusCode != http.StatusTemporaryRedirect {
			t.Error(respBody(resp.Body))
			t.Error("Should be status ok got", resp.StatusCode)
		}
	*/
	donation, err = GetDonationById(db, fixtures.DonationId1)
	if err != nil {
		t.Fatal(err)
	}
	if donation.Status != DonationAccepted {
		logBody(res.Body, t)
		t.Errorf("Expected status Accepted but got %s", donation.Status)
	}

	//Cleanup
	donation.Status = DonationPending
	saveDonation()
}
