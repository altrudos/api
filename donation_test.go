package charityhonor

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
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

func TestDonationInsert(t *testing.T) {
	drive, tx, db := getDriveForTesting()
	donation := Donation{
		DriveId:      drive.Id,
		CharityId:    1,
		Amount:       1234,
		CurrencyCode: "USD",
		DonorName:    "Vindexus",
		Message:      `I'm just trying this <strong>OUT!</strong>`,
	}

	if err := donation.Insert(tx); err != nil {
		t.Error(err)
	}

	if donation.ReferenceCode == "" {
		t.Errorf("No reference code was created")
	}

	if donation.DriveId != drive.Id {
		t.Error("Drive Id doesn't match")
	}

	tx.Commit()

	tx, err := db.Beginx()
	copy, err := GetDonationByReferenceCode(tx, donation.ReferenceCode)
	if err != nil {
		t.Error(err)
	}

	if copy == nil || copy.ReferenceCode != donation.ReferenceCode {
		t.Error("Ref code not made")
	}

	url := donation.GetDonationLink()
	fmt.Println("url", url)

	tx.Rollback()
}
