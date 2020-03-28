package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/charityhonor/ch-api/pkg/justgiving"

	. "github.com/charityhonor/ch-api"
)

/**
 * This command line tool is mostly used for testing.
 * It replicates the logic that would happen if a user were to
 * submit a donation from a form on a website.
 * The difference is that they use a CLI and instead of being
 * redirected to a URL that URL is just outputted.
 */
func populateCharities(name string, args []string) error {
	db := MustGetDefaultDb()
	jg := MustGetDefaultJustGiving()
	var search string
	var charityIds string

	set := flag.NewFlagSet("", flag.ExitOnError)
	set.StringVar(&search, "search", "", "Search for charities with this name.")
	set.StringVar(&charityIds, "charityids", "", "ID on JustGiving of the charity to add.")

	if err := set.Parse(args); err != nil {
		return err
	}

	if search == "" && charityIds == "" {
		return errors.New("Either --search or --charityid is required")
	}

	if search != "" && charityIds != "" {
		return errors.New("Either --search or --charityid is required")
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	var charities []*justgiving.Charity
	if search != "" {
		Pls("Adding charities that match search '%s'", search)
		result, err := jg.SearchCharities(search)
		if err != nil {
			return err
		}

		charities = result.Results
	} else {
		Pls("Adding charities with JustGiving IDs %s", charityIds)
		jgids := strings.Split(charityIds, ",")
		for _, jgid := range jgids {
			i, err := strconv.Atoi(jgid)
			if err != nil {
				return err
			}
			charity, err := jg.GetCharityById(i)
			if err != nil {
				return err
			}
			charities = append(charities, charity)
		}
	}

	Pls("Attempting to add %d charities to db", len(charities))
	for _, v := range charities {
		fmt.Println("JG", v)
		charity := ConvertJGCharity(v)
		fmt.Println("converted", charity)
		if err := charity.Insert(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
