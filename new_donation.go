package altrudos

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/monstercat/pgnull"
)

var (
	ErrNilDonation = errors.New("submitted donation is nil")
)

type SubmittedDonation struct {
	Amount    string
	CharityId string
	Currency  string
	DonorName string
}

func CreateDonation(ext sqlx.Ext, driveId string, dono *SubmittedDonation) (*Donation, error) {
	if dono == nil {
		return nil, ErrNilDonation
	}
	amt, err := AmountFromString(dono.Amount)
	if err != nil {
		return nil, err
	}
	if driveId == "" {
		return nil, ErrDriveNotCreated
	}
	donation := &Donation{
		DonorAmount:   amt,
		DonorCurrency: dono.Currency,
		DonorName:     pgnull.NewNullString(dono.DonorName),
		CharityId:     dono.CharityId,
		DriveId:       driveId,
		Message:       pgnull.NullString{"", false},
	}
	if err := donation.Create(ext); err != nil {
		return nil, err
	}
	return donation, nil
}
