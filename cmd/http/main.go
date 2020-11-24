package main

import "flag"

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "", "Configuration File")
	flag.Parse()
	s, err := NewServerFromConfigFile(confFile)
	if err != nil {
		panic(err)
	}
	if err := s.Run(); err != nil {
		panic(err)
	}
}
