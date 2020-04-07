package main

import (
	"net/http"
	"testing"

	"github.com/charityhonor/ch-api"
	"github.com/charityhonor/ch-api/pkg/fixtures"

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
		"Drives.Data.#":     1,
		"Drives.Total":      1,
		"Drives.Limit":      50,
		"Drives.Offset":     0,
		"Drives.Data.0.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGetTopDrives(t *testing.T) {
	ts, _ := MustGetTestServer(
		DriveRoutes...,
	)

	resp, err := CallJson(ts, http.MethodGet, "/drives/top/week", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drives.#":     1,
		"Drives.0.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestGetDrive(t *testing.T) {
	ts, _ := MustGetTestServer(
		NewGET("/drive/:id", getById("id", "Drive", getDrive)),
	)
	resp, err := CallJson(ts, http.MethodGet, "/drive/"+DriveId, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Should be status ok")
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drive.Id":  DriveId,
		"Drive.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDrive(t *testing.T) {
	ts, _ := MustGetTestServer(
		DriveRoutes...,
	)

	type test struct {
		Payload        interface{}
		ExpectedM      *expectm.ExpectedM
		ExpectedStatus int
	}

	tests := []test{
		{
			Payload:        nil,
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrSourceInvalidURL.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
			},
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"Amount":    "-100.50",
			},
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNegativeAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"Amount":    "100.50",
			},
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNoCharity.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"CharityId": fixtures.CharityId1,
				"Amount":    "100.50",
			},
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"CharityId": fixtures.CharityId1,
				"Amount":    "100.50",
				"Currency":  "fjdksalfjdsla",
			},
			ExpectedStatus: 500,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"CharityId": fixtures.CharityId1,
				"Amount":    "100.50",
				"Currency":  "eur",
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for i, test := range tests {
		resp, err := CallJson(ts, http.MethodPost, "/drive", test.Payload)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("[%d] Status should be %d but got %d", i, test.ExpectedStatus, resp.StatusCode)
		}

		if test.ExpectedM != nil {
			if err := CheckResponseBody(resp.Body, test.ExpectedM); err != nil {
				t.Errorf("[%d] %s", i, err)
			}
		}
	}

	// Cleanup
	db := charityhonor.GetTestDb()
	_, err := db.Exec("DELETE FROM " + charityhonor.TableDonations + " WHERE donor_amount = 10050 AND donor_currency_code = 'EUR'")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM "+charityhonor.TableDrives+" WHERE source_key = $1", "fmgtyqq")
	if err != nil {
		t.Fatal(err)
	}
}
