package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	vinscraper "github.com/Vindexus/go-scraper"

	altrudos "github.com/altrudos/api"
	"github.com/altrudos/api/pkg/fixtures"

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
		t.Error(respBody(resp.Body))
		t.Error("Should be status ok but got", resp.StatusCode)
	}

	if err := CheckResponseBody(resp.Body, &expectm.ExpectedM{
		"Drives.#":                 1,
		"Drives.0.Uri":             "PrettyPinkMoon",
		"Drives.0.TopAmount":       32333,
		"Drives.0.TopNumDonations": 2, // Only 2 are accepted in last 7 days
		"Drives.0.NumDonations":    3, // 3 total
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
		"Drive.Id":                   DriveId,
		"Drive.Uri":                  DriveUri,
		"Drive.NumDonations":         3,
		"RecentDonations.#":          3,
		"TopDonations.#":             3,
		"TopDonations.0.CharityName": "The Demo Charity",
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
				"RawError": vinscraper.ErrSourceInvalidURL.Error(),
			},
		},
		{
			Payload: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedM: &expectm.ExpectedM{
				"RawError": altrudos.ErrNilDonation.Error(),
			},
		},
		{
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
				"SourceUrl": "https://www.reddit.com/r/DunderMifflin/comments/fv3vz0/why_waste_time_say_lot_word_when_few_word_do_trick/fmgtyqq/",
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.50",
					"Currency":  "eur",
					"DonorName": "Elder",
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

	db := altrudos.GetTestDb()

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

	_, err = db.Exec("DELETE FROM "+altrudos.TableDrives+" WHERE source_key = $1", "fmgtyqq")
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
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
			Payload: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.87",
					"Currency":  "eur",
				},
			},
			ExpectedStatus: http.StatusOK,
		},
		{
			Payload: altrudos.FlatMap{
				"SubmittedDonation": altrudos.M{
					"CharityId": fixtures.CharityId1,
					"Amount":    "100.87",
					"Currency":  "eur",
					"DonorName": "Shaper",
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

	db := altrudos.GetTestDb()
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
