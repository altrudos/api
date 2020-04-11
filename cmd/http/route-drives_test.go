package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/charityhonor/ch-api"
	"github.com/charityhonor/ch-api/pkg/fixtures"

	"github.com/monstercat/golib/expectm"
)

var (
	DriveId  = "3656cf1d-8826-404c-8f85-77f3e1f50464"
	DriveUri = "PrettyPinkMoon"
)

func respBody(body io.Reader) string {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return err.Error()
	}
	return string(byt)
}

func TestGetDrives(t *testing.T) {
	ts, _ := MustGetTestServer(
		NewGET("/drives", getDrives),
	)

	resp, err := CallJson(ts, http.MethodGet, "/drives", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error(respBody(resp.Body))
		t.Error("Should be status ok")
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
		//t.Error(respBody(resp.Body))
		t.Error("Should be status ok but got", resp.StatusCode)
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drives.#":     1,
		"Drives.0.Uri": "PrettyPinkMoon",
	}); err != nil {
		t.Error(err)
	}
}

func TestGetDrive(t *testing.T) {
	ts, _ := MustGetTestServer(
		DriveRoutes...,
	)
	resp, err := CallJson(ts, http.MethodGet, "/drive/"+DriveUri, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(respBody(resp.Body))
		t.Error("Should be status ok got", resp.StatusCode)
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drive.Id":           DriveId,
		"Drive.Uri":          DriveUri,
		"Drive.NumDonations": 3,
		"RecentDonations.#":  3,
		"TopDonations.#":     3,
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
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrSourceInvalidURL.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNilDonation.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"Amount": "twenty bucks",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"Amount": "-100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNegativeAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"Amount": "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNoCharity.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "fjdksalfjdsla",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "eur",
				},
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
			t.Error(respBody(resp.Body))
		}

		if test.ExpectedM != nil {
			if err := CheckResponseBody(resp.Body, test.ExpectedM); err != nil {
				t.Errorf("[%d] %s", i, err)
			}
		}
	}

	// Cleanup
	db := charityhonor.GetTestDb()
	_, err := db.Exec("DELETE FROM " + charityhonor.TableDonations + " WHERE donor_amount = 10050 AND donor_currency = 'EUR'")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM "+charityhonor.TableDrives+" WHERE source_key = $1", "fmgtyqq")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateDriveDonation(t *testing.T) {
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
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"Amount": "twenty bucks",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"Amount": "-100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNegativeAmount.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"Amount": "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrNoCharity.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "fjdksalfjdsla",
				},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": charityhonor.ErrInvalidCurrency.Error(),
			},
		},
		{
			Payload: charityhonor.FlatMap{
				"SubmittedDonation": charityhonor.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.87",
					"Currency":  "eur",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
	}

	for i, test := range tests {
		resp, err := CallJson(ts, http.MethodPost, "/drive/"+DriveId+"/donate", test.Payload)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf("[%d] Status should be %d but got %d", i, test.ExpectedStatus, resp.StatusCode)
			t.Error(respBody(resp.Body))
		}

		if test.ExpectedM != nil {
			if err := CheckResponseBody(resp.Body, test.ExpectedM); err != nil {
				t.Errorf("[%d] %s", i, err)
			}
		}
	}

	// Cleanup
	db := charityhonor.GetTestDb()
	_, err := db.Exec("DELETE FROM " + charityhonor.TableDonations + " WHERE donor_amount = 10087 AND donor_currency = 'EUR'")
	if err != nil {
		t.Fatal(err)
	}

}
