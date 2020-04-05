package charityhonor

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/monstercat/pgnull"
)

var (
	ErrDonationNotCreated = errors.New("donation on new drive not created")
	ErrDriveNotCreated    = errors.New("drive on new drive not created")
)

// A new drive is the result of a user submitting the New Drive form
// on the home page
// This form is for both creating a new drive and creating a donation for that
// drive all in one
type NewDrive struct {
	SourceUrl string
	Amount    string
	Currency  string
	CharityId string
	Name      string

	Drive    *Drive
	Donation *Donation
}

// Creates or finds the drive
// Creates the donation
func (nd *NewDrive) Process(ext sqlx.Ext) error {
	if _, err := nd.FetchOrCreateDrive(ext); err != nil {
		return err
	}
	if err := nd.CreateDonation(ext); err != nil {
		return err
	}
	return nil
}

// Based on what the user has submitted, this will either find an existing drive
// for that same source or create a new drive if none is found
func (nd *NewDrive) FetchOrCreateDrive(ext sqlx.Ext) (*Drive, error) {
	source, err := ParseSourceURL(nd.SourceUrl)
	if err != nil {
		return nil, err
	}

	drive, err := GetDriveBySource(ext, source)
	if err != nil {
		return nil, err
	} else if drive != nil {
		nd.Drive = drive
		return drive, nil
	}

	drive, err = CreatedDriveBySourceUrl(ext, nd.SourceUrl)
	if err == nil {
		nd.Drive = drive
	}
	return drive, err
}

func (nd *NewDrive) CreateDonation(ext sqlx.Ext) error {
	amt, err := AmountFromString(nd.Amount)
	if err != nil {
		return err
	}
	if nd.Drive == nil {
		return ErrDriveNotCreated
	}
	donation := &Donation{
		DonorAmount:       amt,
		DonorCurrencyCode: nd.Currency,
		CharityId:         nd.CharityId,
		DriveId:           nd.Drive.Id,
		Message:           pgnull.NullString{"", false},
	}
	if err := donation.Create(ext); err != nil {
		return err
	}
	nd.Donation = donation
	return nil
}
