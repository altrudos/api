package main

func main() {
	s := &Server{
		Port: 8080,
	}
	if err := s.ParseFlags(); err != nil {
		panic(err)
	}

	s.AddRoutes(
		DriveRoutes,
	)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
