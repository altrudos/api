package charityhonor

import (
	"github.com/charityhonor/ch-api/pkg/fixtures"
	"testing"
)

func TestNewDrive(t *testing.T) {
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
		Amount: "25.75",
		CharityId: fixtures.CharityId1,
		Currency: "usd",
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

	if url, err := nd.Donation.GetDonationLink(services.JG); err != nil{
		t.Error(err)
	} else if url == "" {
		t.Error("URL to donate shouldn't be blank")
	}

	// Create another with the same source should return the same drive
	// but a new donation
	nd2 := NewDrive{
		SourceUrl: "https://np.reddit.com/r/ANormalDayInRussia/comments/fucm89/goodbye_anatoly/fmcx3tt/?context=3",
		Amount: "13.75",
		CharityId: fixtures.CharityId1,
		Currency: "cad",
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

	if nd2.Donation.DonorCurrencyCode != "CAD" {
		t.Errorf("Code should be CAD not %s", nd2.Donation.DonorCurrencyCode)
	}
}


func TestNewDriveValidation(t *testing.T) {
	services := GetTestServices()
	nd := NewDrive{
		SourceUrl: "https://www.reddit.com/r/ANormalDayInRussia/comments/fucm89/goodbye_anatoly/fmcx3tt/",
		Amount: "-25.75",
		CharityId: fixtures.CharityId1,
		Currency: "usd",
	}
	if err := nd.Process(services.DB); err != ErrNegativeAmount {
		t.Errorf("Wrong error expected %s but found %s", ErrNegativeAmount, err)
	}

	nd.Amount = "twenty bucks"
	if err := nd.Process(services.DB); err != ErrInvalidAmount {
		t.Errorf("Wrong error expected %s but found %s", ErrInvalidAmount, err)
	}
	nd.Amount = "25.75"

	nd.CharityId = ""
	if err := nd.Process(services.DB); err != ErrNoCharity {
		t.Errorf("Wrong error expected %s but found %s", ErrCharityNotFound, err)
	}


	nd.CharityId = fixtures.DonationId1
	if err := nd.Process(services.DB); err != ErrCharityNotFound {
		t.Errorf("Wrong error expected %s but found %s", ErrCharityNotFound, err)
	}

	nd.CharityId = fixtures.CharityId1
	nd.Currency = "american"
	if err := nd.Process(services.DB); err != ErrInvalidCurrency {
		t.Errorf("Wrong error expected %s but found %s", ErrInvalidCurrency, err)
	}
}