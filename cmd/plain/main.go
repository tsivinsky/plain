package main

import (
	"flag"
	"log"
	"os"

	"github.com/tsivinsky/plain"
)

var (
	port  = flag.Int("p", 5000, "Port to run application on")
	watch = flag.Bool("w", false, "Watch html files for changes")
)

func main() {
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	providedWd := flag.Arg(0)
	if providedWd != "" {
		wd = providedWd
	}

	s := &plain.Server{
		Port:       *port,
		Watch:      *watch,
		WorkingDir: wd,
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
