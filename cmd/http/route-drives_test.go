package main

import (
	"fmt"
	"net/http"
	"testing"

	grtest "github.com/Vindexus/go-router-test"

	vinscraper "github.com/Vindexus/go-scraper"

	altrudos "github.com/altrudos/api"
	"github.com/altrudos/api/pkg/fixtures"

	"github.com/monstercat/golib/expectm"
)

var (
	DriveId  = "3656cf1d-8826-404c-8f85-77f3e1f50464"
	DriveUri = "PrettyPinkMoon"
)

func TestGetDrives(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/drives",
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Drives.Data.#":                3,
			"Drives.Total":                 3,
			"Drives.Limit":                 50,
			"Drives.Offset":                0,
			"Drives.Data.0.Uri":            "PrettyRedMoon",
			"Drives.Data.0.USDAmountTotal": 31001,
		},
	}
	if err := runTest(test); err != nil {
		fmt.Println("resp", test.Response)
		t.Error(err)
	}
}

func TestGetTopDrives(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/drives/top/week",
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Drives.#":                 3,
			"Drives.0.Uri":             "PrettyRedMoon",
			"Drives.0.TopAmount":       31001,
			"Drives.2.TopNumDonations": 1,
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}

func TestGetDrive(t *testing.T) {
	test := &grtest.RouteTest{
		Path:           "/drive/" + DriveUri,
		ExpectedStatus: http.StatusOK,
		ExpectedM: &expectm.ExpectedM{
			"Drive.Id":                   DriveId,
			"Drive.Uri":                  DriveUri,
			"Drive.NumDonations":         2,
			"RecentDonations.#":          2,
			"TopDonations.#":             2,
			"TopDonations.0.CharityName": "The Demo Charity",
		},
	}
	if err := runTest(test); err != nil {
		t.Error(err)
	}
}

func TestCreateDrive(t *testing.T) {
	_, db := MustSetupTestServerDB()

	test := &grtest.RouteTest{
		Path:   "/drive",
		Method: http.MethodPost,
	}

	tests := test.Apply([]*grtest.RouteTest{
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": vinscraper.ErrSourceInvalidURL.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNilDonation.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"Amount": "twenty bucks",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidAmount.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"Amount": "-100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNegativeAmount.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"Amount": "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNoCharity.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidCurrency.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "fjdksalfjdsla",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidCurrency.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "eur",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "eur",
					"DonorName": "Elder",
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedM: &expectm.ExpectedM{
				"Drive.SourceMeta.Title": "Comment by undercooktheonionz",
			},
		},
		{
			Body: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "eur",
					"DonorName": "Elder",
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedM: &expectm.ExpectedM{
				"Drive.SourceMeta.Title": "Why waste time say lot word when few word do trick",
			},
		},
	})

	if err := runTests(tests); err != nil {
		t.Fatal(err)
	}

	dono, err := altrudos.GetDonationByField(db, "donor_name", "Elder")
	if err != nil {
		t.Error(err)
	}
	if dono == nil {
		t.Error("No find Elder's donation")
	}

	// Cleanup
	_, err = db.Exec("DELETE FROM " + altrudos.TableDonations + " WHERE donor_amount = 10050 AND donor_currency = 'EUR'")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM "+altrudos.TableDrives+" WHERE source_key = $1 OR source_key = $2", "fmgtyqq", "fv3vz0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateDriveDonation(t *testing.T) {
	_, db := MustSetupTestServerDB()

	test := &grtest.RouteTest{
		Path:   "/drive/" + DriveId + "/donate",
		Method: http.MethodPost,
	}

	tests := test.Apply([]*grtest.RouteTest{

		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"Amount": "twenty bucks",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidAmount.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"Amount": "-100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNegativeAmount.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"Amount": "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNoCharity.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidCurrency.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "fjdksalfjdsla",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrInvalidCurrency.Error(),
			},
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.87",
					"Currency":  "eur",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Body: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.87",
					"Currency":  "eur",
					"DonorName": "Shaper",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
	})

	if err := runTests(tests); err != nil {
		t.Error(err)
	}

	// Should find one with name
	dono, err := altrudos.GetDonationByField(db, "donor_name", "Shaper")
	if err != nil {
		t.Error(err)
	}
	if dono == nil {
		t.Error("No find Shaper's donation")
	}

	// Cleanup
	_, err = db.Exec("DELETE FROM " + altrudos.TableDonations + " WHERE donor_amount = 10087 AND donor_currency = 'EUR'")
	if err != nil {
		t.Fatal(err)
	}
}
