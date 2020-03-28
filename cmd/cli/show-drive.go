package main

import (
	"flag"

	"github.com/pkg/errors"

	. "github.com/charityhonor/ch-api"
)

func showDrive(name string, args []string) error {
	var confFile string
	var id string
	var uri string
	var sourceurl string
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.StringVar(&confFile, "config", "../../config.toml", "Configuration file")
	set.StringVar(&id, "id", "", "")
	set.StringVar(&uri, "uri", "", "")
	set.StringVar(&sourceurl, "sourceurl", "", "")
	if err := set.Parse(args); err != nil {
		return err
	}
	s := MustGetConfigServices(confFile)
	db := s.DB

	var err error
	var drive *Drive
	if id != "" {
		drive, err = GetDriveById(db, id)
	} else if uri != "" {
		drive, err = GetDriveByUri(db, uri)
	} else if sourceurl != "" {
		drive, err = GetDriveBySourceUrl(db, sourceurl)
	}
	if err != nil {
		return errors.Wrap(err, "could not find drive")
	}

	Pls("URL:     %s", green(drive.SourceUrl))
	Pls("Raised:  %s", green(AmountToString(drive.Amount)))
	Pls("")

	donos, err := drive.GetTopDonations(db, 5)
	if err != nil {
		return errors.Wrap(err, "error getting top donations")
	}

	Pls("Top %d Donations", len(donos))
	for _, v := range donos {
		var name = v.GetDonorName()
		if v.IsAnonymous() {
			name = gray(name)
		}
		Pls("$%s from %s", v.AmountString(), name)
	}

	return nil
}