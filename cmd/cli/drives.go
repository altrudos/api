package main

import (
	"flag"

	. "github.com/charityhonor/ch-api"

	cmd "github.com/tmathews/commander"
)

func cmdDrives(name string, args []string) error {
	return cmd.Exec(args, cmd.DefaultHelper, cmd.M{
		"list": listDrives,
	})
}

func listDrives(name string, args []string) error {
	var confFile string
	set := flag.NewFlagSet(name, flag.ExitOnError)
	set.StringVar(&confFile, "config", "./config.toml", "Configuration file")
	if err := set.Parse(args); err != nil {
		return err
	}
	s := MustGetConfigServices(confFile)

	drives, err := GetDrives(s.DB, nil)
	if err != nil {
		return err
	}

	for _, d := range drives {
		Pls("-=-=-=-=-=-=-=-=-=-=-=-=-=")
		Pls("Drive Id:     %s", d.Id)
		Pls("Drive URI:    %s", d.Uri)
	}
	Pls("-=-=-=-=-=-=-=-=-=-=-=-=-=")

	return nil
}
