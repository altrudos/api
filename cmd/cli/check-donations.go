package main

import (
	"fmt"

	. "github.com/charityhonor/ch-api"
)

func checkDonations(name string, args []string) error {
	db := MustGetDefaultDb()
	jg := MustGetDefaultJustGiving()

	donos, err := GetDonationsToCheck(db)
	if err != nil {
		return err
	}

	for i, dono := range donos {
		if i > 3 {
			break
		}

		fmt.Println("Check this", dono.ReferenceCode)
		jdon, err := jg.GetDonationByReference(dono.ReferenceCode)
		if err != nil {
			fmt.Println("err", err)
			continue
		}

		fmt.Println("jdon amount", jdon.Amount)
	}

	return nil
}
