package main

import (
	"flag"
	"fmt"

	. "github.com/altrudos/api"
)

func checkDonations(name string, args []string) error {
	var confFile string
	var limit int
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.StringVar(&confFile, "config", "../../config.toml", "Configuration file")
	set.IntVar(&limit, "limit", 100, "Max number of donations to check")
	if err := set.Parse(args); err != nil {
		return err
	}
	services := MustGetConfigServices(confFile)
	db := services.DB
	jg := services.JG
	jg.Debug = false

	Pls("Querying for donations to check...")

	donos, err := GetDonationsToCheck(db, limit)
	if err != nil {
		return err
	}

	Pls("Found %d donations to check.", len(donos))

	for _, dono := range donos {
		Pls("-=-=-=-=-=-=-=-=-=-=-=-=-=")
		Pls("Reference:      %s", dono.ReferenceCode)
		Pls("Charity ID:     %s", lblue(dono.CharityId))
		Pls("Drive ID        %s", blue(dono.DriveId))
		Pls("Message:        %s", maybeEmpty(dono.Message.String, lyellow))
		Pls("Donor Amount:   %s", lgreen(AmountToString(dono.DonorAmount)))
		Pls("Donor Currency: %s", lyellow(dono.DonorCurrency))
		if err := dono.CheckStatus(db, jg); err != nil {
			fmt.Println("err", err)
		}
		var status string
		if dono.Status == DonationAccepted {
			status = green(dono.Status)
		} else if dono.Status == DonationRejected {
			status = red(dono.Status)
		} else {
			status = string(dono.Status)
		}
		Pls("Status:         %s", status)
		Pls("Final Amount:   %s", lgreen(AmountToString(dono.FinalAmount)))
		Pls("Final Currency: %s", lyellow(dono.FinalCurrency.String))
	}

	Pls("Done.")

	return nil
}
