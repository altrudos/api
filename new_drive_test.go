package altrudos

import (
	"strings"
	"testing"

	"github.com/altrudos/api/pkg/fixtures"
)

func TestNewDrive(t *testing.T) {
	config := MustGetTestConfig()
	services := GetTestServices()
	db := services.DB
	tx, err := db.Beginx()
	defer db.Close()
	defer tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}
	nd := NewDrive{
		SourceUrl: "https://www.reddit.com/r/ANormalDayInRussia/comments/fucm89/goodbye_anatoly/fmcx3tt/",
		SubmittedDonation: &SubmittedDonation{
			Amount:    "25.75",
			CharityId: fixtures.CharityId1,
			Currency:  "usd",
		},
	}

	if err := nd.Process(tx); err != nil {
		t.Fatal(err)
	}

	if nd.Drive.Uri == "" {
		t.Error("Should have URI")
	}

	if nd.Donation == nil {
		t.Fatal("Donation should not be nil")
	}

	if url, err := nd.Donation.GetDonationLink(services.JG, config.BaseUrl); err != nil {
		t.Error(err)
	} else if url == "" {
		t.Error("URL to donate shouldn't be blank")
		if !strings.Contains(url, "exitUrl") {
			t.Error("donation link should have exitUrl")
		}
	}

	// Create another with the same source should return the same drive
	// but a new donation
	nd2 := NewDrive{
		SourceUrl: "https://np.reddit.com/r/ANormalDayInRussia/comments/fucm89/goodbye_anatoly/fmcx3tt/?context=3",
		SubmittedDonation: &SubmittedDonation{
			Amount:    "13.75",
			CharityId: fixtures.CharityId1,
			Currency:  "cad",
		},
	}

	if err := nd2.Process(tx); err != nil {
		t.Fatal(err)
	}

	if nd.Drive.Uri != nd2.Drive.Uri {
		t.Error("Drive should be same")
	}

	if nd2.Donation.DonorAmount != 1375 {
		t.Error("Donation amount wrong")
	}

	if nd2.Donation.DonorCurrency != "CAD" {
		t.Errorf("Code should be CAD not %s", nd2.Donation.DonorCurrency)
	}
}

func TestNewDriveValidation(t *testing.T) {
	services := GetTestServices()
	nd := NewDrive{
		SourceUrl: "https://www.reddit.com/r/ANormalDayInRussia/comments/fucm89/goodbye_anatoly/fmcx3tt/",
		SubmittedDonation: &SubmittedDonation{
			Amount:    "-25.75",
			CharityId: fixtures.CharityId1,
			Currency:  "usd",
		},
	}
	if err := nd.Process(services.DB); err != ErrNegativeAmount {
		t.Errorf("Wrong error expected %s but found %s", ErrNegativeAmount, err)
	}

	nd.SubmittedDonation.Amount = "twenty bucks"
	if err := nd.Process(services.DB); err != ErrInvalidAmount {
		t.Errorf("Wrong error expected %s but found %s", ErrInvalidAmount, err)
	}
	nd.SubmittedDonation.Amount = "25.75"

	nd.SubmittedDonation.CharityId = ""
	if err := nd.Process(services.DB); err != ErrNoCharity {
		t.Errorf("Wrong error expected %s but found %s", ErrCharityNotFound, err)
	}

	nd.SubmittedDonation.CharityId = fixtures.DonationId1
	if err := nd.Process(services.DB); err != ErrCharityNotFound {
		t.Errorf("Wrong error expected %s but found %s", ErrCharityNotFound, err)
	}

	nd.SubmittedDonation.CharityId = fixtures.CharityId1
	nd.SubmittedDonation.Currency = "american"
	if err := nd.Process(services.DB); err != ErrInvalidCurrency {
		t.Errorf("Wrong error expected %s but found %s", ErrInvalidCurrency, err)
	}
}
