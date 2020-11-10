package altrudos

import (
	"testing"

	"github.com/altrudos/api/pkg/fixtures"

	"github.com/altrudos/api/pkg/justgiving"
)

func TestGetCharity(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	charity, err := GetCharityById(tx, fixtures.CharityId1)
	if err != nil {
		t.Fatal(err)
	}

	if charity.Name != "The Demo Charity" {
		t.Error("Name is wrong!")
	}

	if charity.JustGivingCharityId != justgiving.Fixtures.CharityId {
		t.Errorf("Expected JG Charity ID %v but got %v", justgiving.Fixtures.CharityId, charity.JustGivingCharityId)
	}

	tx.Rollback()
}

func TestInsertCharity(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	charity := &Charity{
		JustGivingCharityId: 835260, // American Red Cross
		Description:         "The American Red Cross prevents and alleviates human suffering in the face of emergencies by mobilizing the power of volunteers and the generosity of donors.",
		Name:                "American Red Cross",
	}
	if err := charity.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if charity.Id == "" {
		t.Fatal("ID should be filled in")
	}

	charity2 := &Charity{
		JustGivingCharityId: 835260, // American Red Cross
		Description:         "This is a duplicate",
		Name:                "American Red Cross Again",
	}
	if err := charity2.Insert(tx); err != ErrDuplicateJGCharityId {
		t.Error("Expected duplicate JGChairtyId error")
	}

	tx.Rollback()
}
