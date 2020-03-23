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

	if d.Id == "" {
		t.Error("Drive id should probably be set")
	}

	tx.Rollback()
}

func TestDriveSelect(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	source := "https://www.reddit.com/r/pathofexile/comments/c7wdss/for_fellow_ssf_bow_users_the_lion_card_farming/eshxtna/"

	//Do some cleanup
	b := QueryBuilder.Delete(TABLE_DRIVES).Where("source_url=?", source).RunWith(db)
	_, err = b.Exec()
	if err != nil {
		t.Fatal(err)
	}

	d := Drive{
		SourceUrl: source,
	}

	err = d.Insert(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	drive, err := GetDriveBySourceUrl(db, source)
	if err != nil {
		t.Fatal(err)
	}

	if drive.SourceUrl != source {
		t.Error("Wrong source URL returned")
	}
}
