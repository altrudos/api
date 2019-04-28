package main

import (
	"fmt"

	ch "github.com/Vindexus/ch-api"
)

func main() {
	d := ch.Drive{
		SourceUrl: "Whatever",
	}
	d.GenerateUri()
	fmt.Println("Welcome to CLI tool")
	fmt.Println("You made", d.Uri)
	DoThing()
}
func DoThing() {
	fmt.Println("Do stuff")
}
