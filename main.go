package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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

	routes, err := getRoutes(pagesPath, wd)
	if err != nil {
		log.Fatal(err)
	}

	portStr := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(portStr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := matchRoute(r, routes)
		if route == nil {
			fp := path.Join(wd, StaticDir, r.URL.Path)
			staticFile, err := getStaticFile(fp)
			if err != nil {
				w.WriteHeader(404)
				return
			}

			mime := http.DetectContentType(staticFile)
			if strings.Contains(mime, "text/plain") {
				mime = getFileTypeByName(fp)
			}

			w.Header().Set("Content-Type", mime)

			w.Write(staticFile)
			return
		}

		err := renderHTMLFile(w, route.filepath)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
		}
	}))
	if err != nil {
		log.Fatal(err)
	}
}

func matchRoute(r *http.Request, routes []route) *route {
	uri := r.URL

	for _, route := range routes {
		if route.urlpath == uri.Path || route.urlpath == uri.Path+"/" || route.urlpath+"/" == uri.Path {
			return &route
		}
	}

	return nil
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

func getFileTypeByName(filename string) string {
	ext := path.Ext(filename)

	switch ext {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	}

	return "text/plain"
}
