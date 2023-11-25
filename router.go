package plain

import (
	"os"
	"path"
	"strings"
)

const (
	pagesDir          = "pages"
	staticDir         = "public"
	indexPageFileName = "index"
)

// function for reading files when building routes' list.
// also, it will skip files which returned `nil` for `data` when building routes
type ReadPageFileFunc = func(filepath string) (data []byte, err error)

type route struct {
	filepath string
	urlpath  string
	data     []byte
}

func getRoutes(p string, wd string, readPageFile ReadPageFileFunc) ([]route, error) {
	var routes []route

	dir, err := os.ReadDir(p)
	if err != nil {
		return routes, err
	}

	for _, file := range dir {
		filepath := path.Join(p, file.Name())

		if file.IsDir() {
			nestedRoutes, err := getRoutes(filepath, wd, readPageFile)
			if err != nil {
				return routes, err
			}

			routes = append(routes, nestedRoutes...)
			continue
		}

		fileExt := path.Ext(filepath)
		route := route{
			filepath: filepath,
			urlpath:  filePathToUrl(wd, filepath, fileExt),
		}

		data, err := readPageFile(filepath)
		if err != nil {
			return nil, err
		}

		if data == nil {
			continue
		}

		route.data = data

		routes = append(routes, route)
	}

	return routes, nil
}

// removes 'index' part from filename, so need to check if file is directory before calling this function
func filePathToUrl(wd, filepath, ext string) string {
	fileUrlPath := strings.ReplaceAll(filepath, ext, "")

	wdWithPagesDir := path.Join(wd, pagesDir)
	urlpath := strings.ReplaceAll(fileUrlPath, wdWithPagesDir, "")

	if strings.HasSuffix(urlpath, "index") {
		urlpath = strings.ReplaceAll(urlpath, "index", "")
	}

	return urlpath
}
