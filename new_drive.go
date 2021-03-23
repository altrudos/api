package altrudos

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	SourceUrl         string
	Name              string
	SubmittedDonation *SubmittedDonation

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
	drive, err := GetDriveBySourceUrl(ext, nd.SourceUrl)
	fmt.Println("drive", drive)
	fmt.Println("err", err)
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
	dono, err := CreateDonation(ext, nd.Drive.Id, nd.SubmittedDonation)
	if err != nil {
		return err
	}
	nd.Donation = dono
	return nil
}
