package justgiving

import (
	"testing"
	"time"
)

func TestGetDonationByReferenceCode(t *testing.T) {
	code := Fixtures.DonationReferenceCode //This is a code I usd to manually make a dono a while ago
	jg := GetTestJG()
	dono, err := jg.GetDonationByReference(code)
	if err != nil {
		t.Fatal(err)
	}

	expectedAmount := float64(9.7412)
	if dono.GetAmount() != expectedAmount {
		t.Errorf("Expected amount %v but got %v", expectedAmount, dono.Amount)
	}

	if dono.GetDate().IsZero() {
		t.Error("Date is zero")
	}

	if dono.GetDate().After(time.Now()) {
		t.Error("Get Date returns date in future")
	}

	if dono.ThirdPartyReference != code {
		t.Errorf("Expected ThirdPartyReference %v but got %v", code, dono.ThirdPartyReference)
	}

	expectedLocalAmount := float64(12.34)
	if dono.GetLocalAmount() != expectedLocalAmount {
		t.Errorf("Expected GetLocalAmount %v but got %v", expectedLocalAmount, dono.GetLocalAmount())
	}

	expectedLocalCurrency := "USD"
	if dono.LocalCurrencyCode != expectedLocalCurrency {
		t.Errorf("Expected LocalCurrencyCode %v but got %v", expectedLocalCurrency, dono.LocalCurrencyCode)
	}
}

func TestGetDonationById(t *testing.T) {
	jg := GetTestJG()
	_, err := jg.GetDonationById(483905)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCharity(t *testing.T) {
	jg := GetTestJG()
	id := Fixtures.CharityId //this is their demo ID in their docs
	charity, err := jg.GetCharityById(id)
	if err != nil {
		t.Fatal(err)
	}

	expectedName := "The Demo Charity"
	if charity.Name != expectedName {
		t.Errorf("Expected charity with name '%s' but got '%s'", expectedName, charity.Name)
	}

	//This UUID is actually what is in the description
	expDescription := "29c50192-e194-4fd8-9ae5-333d54e9c357"
	if charity.Description != expDescription {
		t.Errorf("Expected charity with description '%s' but got '%s'", expDescription, charity.Description)
	}
}
