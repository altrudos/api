package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/monstercat/pgnull"

	. "github.com/charityhonor/ch-api"
)

func checkDonations(name string, args []string) error {
	var confFile string
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.StringVar(&confFile, "config", "./config.toml", "Configuration file")
	if err := set.Parse(args); err != nil {
		return err
	}
	services := MustGetConfigServices(confFile)
	db := services.DB
	jg := services.JG

	Pls("Querying for donations to check...")

	donos, err := GetDonationsToCheck(db)
	if err != nil {
		return err
	}

	Pls("Found %d donations to check.", len(donos))

	for i, dono := range donos {
		if i > 3 {
			break
		}

		Pls("-=-=-=-=-=-=-=-=-=-=-=-=-=")
		Pls("Reference:  %s", dono.ReferenceCode)
		Pls("Charity ID: %s", green(dono.CharityId))
		Pls("Drive ID    %s", blue(dono.DriveId))
		Pls("Message:    %s", maybeEmpty(dono.Message.String, lyellow))
		Pls("Amount:     %s", lgreen(AmountToString(dono.DonorAmount)))
		Pls("Currency:   %s", lyellow(dono.DonorCurrencyCode))
		jdon, err := jg.GetDonationByReference(dono.ReferenceCode)
		if err != nil {
			fmt.Println("Error finding dono on JG", err)
			dono.LastChecked = pgnull.NullTime{time.Now(), true}
			if err := dono.Save(db); err != nil {
				return err
			}
			continue
		}

		amount, err := strconv.ParseFloat(jdon.Amount, 64)
		if err != nil {
			return err
		}

		dono.FinalAmount = int(amount * 100)
		dono.Status = DonationAccepted
		if err := dono.Save(db); err != nil {
			return err
		}
	}

	Pls("Done.")

	return nil
}
