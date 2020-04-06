package main

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/tmathews/commander"
)

var (
	blue    = color.FgBlue.Render
	lblue   = color.FgLightBlue.Render
	green   = color.FgGreen.Render
	lgreen  = color.LightGreen.Render
	lyellow = color.FgLightYellow.Render
	gray    = color.Gray.Render
	italic  = color.OpItalic.Render
	red     = color.FgRed.Render
	lred    = color.FgLightRed.Render
)

func spl(args ...interface{}) {
	fmt.Println(fmt.Sprintf(args[0].(string), args[1:]...))
}

func maybeEmpty(msg string, color func(...interface{}) string) string {
	if msg == "" {
		return gray(italic("Empty"))
	}

	return color(msg)
}

func main() {
	var args []string
	if len(os.Args) < 2 {
		args = []string{}
	} else {
		args = os.Args[1:]
	}
	err := cmd.Exec(args, cmd.DefaultHelper, cmd.M{
		"create-donation":    createDonation,
		"check-donations":    checkDonations,
		"populate-charities": populateCharities,
		"drives":             cmdDrives,
		"show-drive":         showDrive,
	})
	if err != nil {
		switch v := err.(type) {
		case cmd.Error:
			fmt.Print("ERROR:", v.Help())
			os.Exit(2)
		default:
			fmt.Println(color.FgRed.Render("ERROR:" + err.Error()))
			os.Exit(1)
		}
	}
}

func Pls(msg string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, args...))
}
