package charityhonor

import (
	"testing"

	"github.com/charityhonor/ch-api/pkg/justgiving"
)

func TestGetCharity(t *testing.T) {
	db := GetTestDb()
	tx, err := db.Beginx()

	charity, err := GetCharityById(tx, 1)
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
