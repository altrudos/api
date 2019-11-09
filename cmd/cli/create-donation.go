package main

import (
	"errors"
	"flag"
	"fmt"

	. "github.com/charityhonor/ch-api"
)

/**
 * This command line tool is mostly used for testing.
 * It replicates the logic that would happen if a user were to
 * submit a donation from a form on a website.
 * The difference is that they use a CLI and instead of being
 * redirected to a URL that URL is just outputted.
 */
func createDonation(name string, args []string) error {
	db := MustGetDefaultDb()
	jg := MustGetDefaultJustGiving()
	fmt.Println("jg", jg)
	var amount float64
	var currency string
	var message string
	var charityid int
	var sourceUrl string

	set := flag.NewFlagSet("", flag.ExitOnError)
	set.Float64Var(&amount, "amount", 0.00, "Donation amount")
	set.IntVar(&charityid, "charityid", 0, "The charity ID locally")
	set.StringVar(&currency, "currency", "USD", "The currency code to use.")
	set.StringVar(&sourceUrl, "url", "", "The URL of the content to honor.")
	set.StringVar(&message, "message", "", "Message.")

	if err := set.Parse(args); err != nil {
		return err
	}

	if charityid == 0 {
		return errors.New("--charityid is required. It is the postgres ID in our database.")
	}

	if sourceUrl == "" {
		return errors.New("--url is required. It is the source URL to honor.")
	}

	charity, err := GetCharityById(db, charityid)
	if err != nil {
		return err
	}

	drive, err := GetOrCreateDriveBySourceUrl(db, sourceUrl)
	if err != nil {
		return err
	}

	spl("Mode:     %s", lyellow(jg.Mode))
	spl("Charity:  %s", green(charity.Name))
	spl("Drive:    %s", blue(drive.Name))
	spl("Message:  %s", maybeEmpty(message, lyellow))
	spl("Amount:   %s", lgreen(AmountToString(amount)))
	spl("Currency: %s", lyellow(currency))

	donation := drive.GenerateDonation()
	donation.CharityId = charityid
	donation.Amount = amount
	donation.Message = message

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = donation.Create(tx)
	if err != nil {
		return err
	}
	spl("")
	link, err := donation.GetDonationLink(jg)
	if err != nil {
		return err
	}
	spl("Donation Link: %s", lblue(link))

	return nil
}
