package charityhonor

import (
	"fmt"
	"testing"
)

func TestDriveInsert(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()

	if err != nil {
		t.Fatal(err)
	}

	d := Drive{
		SourceUrl: "https://reddit.com/r/gaming",
	}

	err = d.Insert(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if d.Uri == "" {
		t.Error("Failed to generate URI")
	} else {
		fmt.Println("URI:", d.Uri)
	}

	if d.Id < 1 {
		t.Error("Drive id should probably be 1")
	}

	tx.Rollback()
}

func TestDonationGenerate(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()

	if err != nil {
		t.Fatal(err)
	}

	d := Drive{
		SourceUrl: "https://reddit.com/r/gaming",
	}

	err = d.Insert(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	donation, err := d.CreateDonation(tx)
	if err != nil {
		t.Fatal(err)
	}

	donation.CharityId = 1
	donation.Amount = 2000

	if err := donation.Insert(tx); err != nil {
		t.Error(err)
	}

	if donation.ReferenceCode == "" {
		t.Errorf("No reference code was created")
	}

	if donation.DriveId != d.Id {
		t.Error("Drive Id doesn't match")
	}

	tx.Commit()

	tx, err = db.Beginx()

	copy, err := GetDonationByReferenceCode(tx, donation.ReferenceCode)

	if err != nil {
		t.Error(err)
	}

	if copy == nil || copy.ReferenceCode != donation.ReferenceCode {
		t.Error("Ref code not made")
	}

	tx.Rollback()
}
