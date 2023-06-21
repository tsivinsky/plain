package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	port := flag.Int("p", 5000, "Port to run application on")

	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pagesPath := path.Join(wd, PagesDir)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		log.Fatal(err)
	}

	routes := getRoutes(pagesPath, wd)

	portStr := fmt.Sprintf(":%d", *port)
	http.ListenAndServe(portStr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL

		for _, route := range routes {
			if route.urlpath == uri.Path || route.urlpath == uri.Path+"/" || route.urlpath+"/" == uri.Path {
				err := renderHTMLFile(w, route.filepath)
				if err != nil {
					fmt.Fprintf(w, "error: %v", err)
				}
				return
			}
		}
	}))
}

func renderHTMLFile(w http.ResponseWriter, filepath string) error {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html")
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	return nil
}
