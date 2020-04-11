package charityhonor

import (
	"fmt"
	"testing"

	"github.com/monstercat/pgnull"

	"github.com/monstercat/golib/expectm"

	"github.com/charityhonor/ch-api/pkg/fixtures"
)

func TestDriveSource(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	drive, err := CreatedDriveBySourceUrl(tx, "https://np.reddit.com/r/pathofexile/comments/c6oy9e/to_everyone_that_feels_bored_by_the_game_or/esai27c/?context=3")
	if err != nil {
		t.Fatal(err)
	}

	d, err := GetDriveByUri(tx, drive.Uri)
	if err != nil {
		t.Error(err)
	}

	meta := d.SourceMeta
	if v, ok := meta["subreddit"]; !ok {
		t.Error("no subreddit in meta")
	} else if v != "pathofexile" {
		t.Errorf("Expected subreddit pathofexile but found %s", v)
	}

	tx.Rollback()
}

func TestDriveInsert(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	d := Drive{
		SourceUrl:  "https://reddit.com/r/gaming",
		SourceType: STURL,
		SourceKey:  "rgaming",
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
		SourceUrl:  source,
		SourceKey:  "eshxtna",
		SourceType: STRedditComment,
		Uri:        uri,
	}
	if err := d.Insert(db); err != nil {
		t.Fatal(err)
		return
	}

	drive, err := GetDriveBySourceUrl(db, source)
	if err != nil {
		t.Fatal(err)
	}

	if drive.SourceType != STRedditComment {
		t.Error("Wrong source type")
	}

	// A similar source URL that creates the same source (reddit comment / eshxtna)
	// should find the same drive
	sourceSimilar := "https://www.reddit.com/r/WRONGSUBREDDIT/comments/c7wdss/for_fellow_ssf_bow_users_the_lion_card_farming/eshxtna/?context=3"
	drive2, err := GetDriveBySourceUrl(db, sourceSimilar)
	if err != nil {
		t.Fatal(err)
	}

	if drive2.Uri != drive.Uri {
		t.Errorf("Wrong drive found. Expected %v but found %v", drive, drive2)
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

func TestGetTopDrives(t *testing.T) {
	db := GetTestDb()
	drives, err := GetTopDrives(db)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 780 + 31001 + 1332
	if len(drives) != 1 {
		t.Fatalf("Expected 1 drives found %d", len(drives))
	}

	if drives[0].TopAmount != expectedTotal {
		t.Errorf("Top amount should be %d not %d", expectedTotal, drives[0].TopAmount)
	}

	// Create a new drive
	drive := &Drive{
		Uri:        "MadeUpNewDrive",
		SourceType: STURL,
	}
	if err := drive.Create(db); err != nil {
		t.Fatal(err)
	}

	drives, err = GetTopDrives(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(drives) == 2 {
		t.Fatal("Top drives shouldn't have increased because total of new one is still zero")
	}

	if drives[0].Uri == drive.Uri {
		t.Error("New drive shouldn't be top")
	}

	donation := Donation{
		CharityId:     fixtures.CharityId1,
		Status:        DonationAccepted,
		FinalAmount:   32,
		USDAmount:     50,
		DonorCurrency: "GBP",
		DriveId:       drive.Id,
		FinalCurrency: pgnull.NullString{"EUR", true},
	}
	if err := donation.Create(db); err != nil {
		t.Fatal(err)
	}

	drives, err = GetTopDrives(db)
	if err != nil {
		t.Fatal(err)
	}
	if drives[0].Uri == drive.Uri {
		t.Error("New drive still shouldn't be top")
	}

	donation2 := Donation{
		CharityId:     fixtures.CharityId1,
		Status:        DonationAccepted,
		FinalAmount:   1000000,
		USDAmount:     1000000,
		DonorCurrency: "CAD",
		DriveId:       drive.Id,
		FinalCurrency: pgnull.NullString{"USD", true},
	}
	if err := donation2.Create(db); err != nil {
		t.Fatal(err)
	}

	drives, err = GetTopDrives(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(drives) < 2 {
		t.Fatal("Should have two drives")
	}
	if drives[0].Uri != drive.Uri {
		t.Error("New drive should be top now")
	}
	if drives[0].TopAmount != 1000050 {
		for _, v := range drives {
			fmt.Println("top", v.TopAmount)
		}
		t.Errorf("Top amount should be 100050 not %d", drives[0].TopAmount)
	}

	// Cleanup

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

	recent, err := drive.GetRecentDonations(db, 5)
	if err != nil {
		t.Error(err)
	}

	if len(top) != 3 {
		t.Errorf("Expected 3 donations found %d", len(top))
	}

	if err := expectm.CheckJSON(recent, &expectm.ExpectedM{
		"0.FinalAmount": 1332,
		"1.DonorName":   "Big Spender",
		"1.FinalAmount": 31001,
		"2.FinalAmount": 780,
	}); err != nil {
		t.Error(err)
	}
}
