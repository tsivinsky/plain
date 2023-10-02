package main

import (
	"flag"
	"log"

	"github.com/tsivinsky/plain"
)

var (
	port  = flag.Int("p", 5000, "Port to run application on")
	watch = flag.Bool("w", false, "Watch html files for changes")
)

func main() {
	flag.Parse()

	s := &plain.Server{
		Port:  *port,
		Watch: *watch,
	}

	err := s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
