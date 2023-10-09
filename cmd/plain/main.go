package main

import (
	"flag"
	"log"
	"os"

	"github.com/tsivinsky/plain"
)

var (
	port  = flag.Int("p", 5000, "Port to run application on")
	host  = flag.String("H", "localhost", "Host to run application on")
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

	s, err := plain.New(plain.Options{
		Host:       *host,
		Port:       *port,
		WorkingDir: wd,
	})
	if err != nil {
		log.Fatal(err)
	}

	if *watch {
		go s.Watch()
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
