package charityhonor

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/charityhonor/ch-api/pkg/fixtures"

	"github.com/charityhonor/ch-api/pkg/justgiving"
	"github.com/jmoiron/sqlx"
	. "github.com/monstercat/pgnull"
)

var numTestDrives = 0

func getDriveForTesting() (*Drive, *sqlx.Tx, *sqlx.DB) {
	numTestDrives++
	db := GetTestDb()
	drive := &Drive{
		SourceUrl: "http://www.reddit.com/r/number" + strconv.Itoa(numTestDrives),
	}

	tx, err := db.Beginx()
	if err != nil {
		panic(err)
	}
	if err = drive.Insert(tx); err != nil {
		panic(err)
	}

	return drive, tx, db
}

func TestDonationCRUD(t *testing.T) {
	drive, tx, _ := getDriveForTesting()

	donation := Donation{
		DriveId:           drive.Id,
		CharityId:         fixtures.CharityId1,
		DonorAmount:       float64(12.34),
		DonorCurrencyCode: "USD",
		DonorName:         NullString{"Vindexus", true},
		Message:           NullString{`I'm just trying this <strong>OUT!</strong>`, true},
	}

	if err := donation.Create(tx); err != nil {
		t.Error(err)
	}

	if donation.ReferenceCode == "" {
		t.Fatal("No reference code was created")
	}

	if donation.DriveId != drive.Id {
		t.Error("Drive Id doesn't match")
	}

	dono2, err := GetDonationByReferenceCode(tx, donation.ReferenceCode)
	if err != nil {
		t.Error(err)
	}

	if dono2.DonorAmount != donation.DonorAmount {
		t.Errorf("Expected Amount '%v' but got '%v'", donation.DonorAmount, dono2.DonorAmount)
	}

	if dono2.DonorCurrencyCode != donation.DonorCurrencyCode {
		t.Errorf("Expected CurrencyCode '%v' but got '%v'", donation.DonorCurrencyCode, dono2.DonorCurrencyCode)
	}

	if dono2.Message != donation.Message {
		t.Errorf("Expected Message '%v' but got '%v'", donation.Message, dono2.Message)
	}

	if dono2.Charity == nil {
		t.Fatal("Donation's Charity property was nil")
	}

	if dono2 == nil || dono2.ReferenceCode != donation.ReferenceCode {
		t.Error("Ref code not made")
	}

	jg := justgiving.GetTestJG()

	donError := Donation{}

	_, err = donError.GetDonationLink(jg)
	if err == nil {
		t.Error("Get donation link should fail if missing data")
	}

	donError.DonorAmount = 30
	_, err = donError.GetDonationLink(jg)
	if err == nil {
		t.Error("Get donation link should fail if missing amount")
	}

	donError.DonorCurrencyCode = "USD"
	_, err = donError.GetDonationLink(jg)
	if err == nil {
		t.Error("Get donation link should fail if missing currency code")
	}

	donError.CharityId = fixtures.CharityId1
	_, err = donError.GetDonationLink(jg)
	if err == nil {
		t.Error("Get donation link should fail if missing charity")
	}

	url, err := dono2.GetDonationLink(jg)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(url, strconv.Itoa(justgiving.Fixtures.CharityId)) {
		t.Errorf("Url should contain %v, got %s", justgiving.Fixtures.CharityId, url)
	}

	newAmount := float64(1337)
	newName := "Colin 9430843290"
	dono2.DonorAmount = newAmount
	dono2.DonorName = NullString{newName, true}
	err = dono2.Save(tx)
	if err != nil {
		t.Fatal(err)
	}

	dono3, err := GetDonationByReferenceCode(tx, dono2.ReferenceCode)
	if err != nil {
		t.Fatal(err)
	}
	if dono3.DonorName.String != newName {
		t.Errorf("Expected name %v got %v", newName, dono3.DonorName)
	}

	if dono3.DonorAmount != newAmount {
		t.Errorf("Expected amount %v got %v", newAmount, dono3.DonorAmount)
	}

	// Get multiple donations
	donos, err := GetDonations(tx, &DonationOperators{})
	if err != nil {
		t.Fatal(err)
	}

	if len(donos) == 0 {
		t.Error("Found 0 donations")
	}

	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
}

func TestGetDonationsToCheck(t *testing.T) {
	db := GetTestDb()

	// Get multiple donations
	donos, err := GetDonationsToCheck(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(donos) != 2 {
		t.Errorf("Expected 2 donations, found %d", len(donos))
	}

	for i, dono := range donos {
		if dono.Status != DonationPending {
			t.Errorf("[%d] Expected pending, found %s", i, dono.Status)
		}
	}
}

func TestDonationChecking(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}
	dono, err := GetDonationById(tx, fixtures.DonationId1)
	if err != nil {
		t.Fatal(err)
	}

	jg := justgiving.GetTestJG()

	//Change the data from whatever's in the db to this.
	dono.ReferenceCode = justgiving.Fixtures.DonationReferenceCode
	dono.Status = DonationPending

	err = dono.CheckStatus(tx, jg)
	if err != nil {
		t.Fatal(err)
	}

	if dono.Status != DonationAccepted {
		t.Errorf("Expectd donation %s but was %v", DonationAccepted, dono.Status)
	}

	if dono.GetLastChecked().IsZero() {
		t.Error("last checked is zero")
	}

	if dono.GetLastChecked().Before(time.Now().Add(time.Second * -1)) {
		t.Error("Last checked should be older than 1s ago")
	}

	if dono.GetLastChecked().After(time.Now().Add(time.Second)) {
		t.Error("Last checked shouldn't be in the future")
	}

	dono.ReferenceCode = "nonexistantcode"
	dono.Status = DonationPending

	err = dono.CheckStatus(tx, jg)
	if err != nil {
		t.Fatal(err)
	}

	if dono.Status != DonationRejected {
		t.Error("Donation should be rejected if we can't find it in JG by reference code")
	}

	if err := tx.Rollback(); err != nil {
		t.Error(err)
	}
}
