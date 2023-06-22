package main

import (
	"log"
	"os"
	"path"
	"strings"
)

const (
	PagesDir          = "pages"
	IndexPageFileName = "index"
)

type route struct {
	filepath string
	urlpath  string
}

func getRoutes(p string, wd string) []route {
	var routes []route

	dir, err := os.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		filepath := path.Join(p, file.Name())

		if file.IsDir() {
			nestedRoutes := getRoutes(filepath, wd)
			routes = append(routes, nestedRoutes...)
			continue
		}

		fileExt := path.Ext(filepath)
		if fileExt != ".html" {
			continue
		}

		fileUrlPath := strings.ReplaceAll(filepath, fileExt, "")

		wdWithPagesDir := path.Join(wd, PagesDir)
		urlpath := strings.ReplaceAll(fileUrlPath, wdWithPagesDir, "")

		// It's safe because we handle case if it's directory above
		if strings.HasSuffix(urlpath, "index") {
			urlpath = strings.ReplaceAll(urlpath, "index", "")
		}

		route := route{
			filepath: filepath,
			urlpath:  urlpath,
		}

		routes = append(routes, route)
	}

	return routes
}
