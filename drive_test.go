package charityhonor

import (
	"fmt"
	"testing"

	"github.com/monstercat/golib/expectm"

	"github.com/charityhonor/ch-api/pkg/fixtures"
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

	err = d.Create(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if d.Uri == "" {
		t.Error("Failed to generate URI")
	} else {
		fmt.Println("URI:", d.Uri)
	}

	if d.Id == "" {
		t.Error("Drive id should probably be set")
	}

	tx.Rollback()
}

func TestDriveSelect(t *testing.T) {
	db := GetTestDb()
	source := "https://www.reddit.com/r/pathofexile/comments/c7wdss/for_fellow_ssf_bow_users_the_lion_card_farming/eshxtna/"
	uri := GenerateUri()
	d := Drive{
		SourceUrl: source,
		Uri:       uri,
	}
	if err := d.Insert(db); err != nil {
		t.Fatal(err)
		return
	}

	drive, err := GetDriveBySourceUrl(db, source)
	if err != nil {
		t.Fatal(err)
	}

	if drive.SourceUrl != source {
		t.Error("Wrong source URL returned")
	}

	drive, err = GetDriveById(db, fixtures.DriveId)
	if err != nil {
		t.Error(err)
	}

	drive, err = GetDriveByUri(db, fixtures.DriveUri)
	if err != nil {
		t.Error(err)
	}

	// This needs to be tested last because it causes the Tx to have an error
	d2 := Drive{
		SourceUrl: source,
		Uri:       uri,
	}
	err = d2.Insert(db)
	if err == nil {
		t.Error("Should have error for duplicate uri")
	}

	// Cleanup
	_, err = db.Exec(`DELETE FROM drives WHERE uri = $1`, uri)
}

func TestGetDrives(t *testing.T) {
	db := GetTestDb()

	drives, err := GetDrives(db, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(drives) == 0 {
		t.Error("Expecting some drives")
	}
}

func TestGetDriveDonations(t *testing.T) {
	db := GetTestDb()
	drive, err := GetDriveById(db, fixtures.DriveId)
	if err != nil {
		t.Fatal(err)
	}

	top, err := drive.GetTopDonations(db, 5)
	if err != nil {
		t.Error(err)
	}

	if len(top) != 3 {
		t.Errorf("Expected 3 donations found %d", len(top))
	}

	if err := expectm.CheckJSON(top, &expectm.ExpectedM{
		"0.FinalAmount": 31001,
		"0.DonorName":   "Big Spender",
		"1.FinalAmount": 1332,
		"2.FinalAmount": 780,
	}); err != nil {
		t.Error(err)
	}
}
