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
}
