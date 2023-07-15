package main

import (
	"io"
	"os"
	"path"
	"strings"
)

const (
	PagesDir          = "pages"
	StaticDir         = "public"
	IndexPageFileName = "index"
)

type route struct {
	filepath string
	urlpath  string
}

func getRoutes(p string, wd string) ([]route, error) {
	var routes []route

	dir, err := os.ReadDir(p)
	if err != nil {
		return routes, err
	}

	for _, file := range dir {
		filepath := path.Join(p, file.Name())

		if file.IsDir() {
			nestedRoutes, err := getRoutes(filepath, wd)
			if err != nil {
				return routes, err
			}

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

	return routes, nil
}

func getStaticFile(p string) ([]byte, error) {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil, err
	}

	f, err := os.OpenFile(p, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
