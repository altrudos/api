package main

import (
	"errors"
	"flag"

	"github.com/monstercat/pgnull"
	errs "github.com/pkg/errors"

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
	var confFile string
	flag.StringVar(&confFile,"config", "./config.toml", "Configuration file")
	flag.Parse()
	services := MustGetConfigServices(confFile)
	db := services.DB
	jg := services.JG

	var amount float64
	var currency string
	var message string
	var charityid string
	var charityname string
	var donorname string
	var sourceUrl string

	set := flag.NewFlagSet("", flag.ExitOnError)
	set.Float64Var(&amount, "amount", 0.00, "Donation amount")
	set.StringVar(&charityname, "charityname", "", "Exact charity name in our db")
	set.StringVar(&charityid, "charityid", "", "The charity ID in our db")
	set.StringVar(&currency, "currency", "USD", "The currency code to use.")
	set.StringVar(&sourceUrl, "url", "", "The URL of the content to honor.")
	set.StringVar(&message, "message", "", "Message.")
	set.StringVar(&donorname, "donorname", "", "Donor's name.")

	if err := set.Parse(args); err != nil {
		return err
	}

	if charityid == "" && charityname == "" {
		return errors.New("Either --charityid or --charityname is required. It is the value in our database.")
	}

	if sourceUrl == "" {
		return errors.New("--url is required. It is the source URL to honor.")
	}

	if amount <= 0 {
		return errors.New("--amount must be positive")
	}

	var charity *Charity
	var err error
	if charityid != "" {
		charity, err = GetCharityById(db, charityid)
		if err != nil {
			return errs.Wrap(err, "could not find charity by id")
		}
	} else {
		charity, err = GetCharityByName(db, charityname)
		if err != nil {
			return errs.Wrap(err, "could not find charity by name "+charityname)
		}
	}

	drive, err := GetOrCreateDriveBySourceUrl(db, sourceUrl)
	if err != nil {
		return errs.Wrap(err, "could not get or create drive from source")
	}

	donation := drive.GenerateDonation()
	donation.CharityId = charity.Id
	donation.DriveId = drive.Id
	donation.Message = pgnull.NullString{message, message != ""}
	donation.DonorAmount = amount
	donation.DonorCurrencyCode = currency
	donation.DonorName = pgnull.NullString{donorname, donorname != ""}

	spl("Mode:           %s", lyellow(jg.Mode))
	spl("Charity:        %s", green(charity.Name))
	spl("Drive:          %s", blue(drive.Name))
	spl("Message:        %s", maybeEmpty(donation.Message.String, lyellow))
	spl("Donor Amount:   %s", lgreen(AmountToString(donation.DonorAmount)))
	spl("Donor Currency: %s", lyellow(donation.DonorCurrencyCode))

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = donation.Create(tx)
	if err != nil {
		Pls("Error creating donation")
		return err
	}
	spl("")
	link, err := donation.GetDonationLink(jg)
	if err != nil {
		return err
	}
	spl("Donation Link: %s", lblue(link))

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
