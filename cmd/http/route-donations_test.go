package main

import (
	"net/http"
	"testing"

	grtest "github.com/Vindexus/go-router-test"

	"github.com/altrudos/api/pkg/fixtures"

	. "github.com/altrudos/api"

	"github.com/monstercat/golib/expectm"
)

func TestGetDonationsRecent(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/donations/recent",
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Donations.#":             4,
			"Donations.0.Drive.Uri":   "PrettyPinkMoon",
			"Donations.0.FinalAmount": 1332,
			"Donations.0.DonorName":   "FordonGreeman",
			"Donations.0.CharityName": "The Demo Charity",
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}

func TestCheckDonation(t *testing.T) {
	_, db := MustSetupTestServerDB()

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

	test := &grtest.RouteTest{
		Path:           "/donations/check/" + donation.ReferenceCode,
		ExpectedStatus: http.StatusTemporaryRedirect,
		NilResponse:    true,
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}

	donation, err = GetDonationById(db, fixtures.DonationId1)
	if err != nil {
		t.Fatal(err)
	}

	//Cleanup
	donation.Status = DonationPending
	saveDonation()
}

func TestGetDonation(t *testing.T) {
	_, db := MustSetupTestServerDB()

	donation, err := GetDonationById(db, fixtures.DonationId1)
	if err != nil {
		t.Fatal(err)
	}

	test := &grtest.RouteTest{
		Path:           "/donations/byref/" + donation.ReferenceCode,
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Donation.Id": fixtures.DonationId1,
			"Drive.Uri":   fixtures.DriveUri,
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}
