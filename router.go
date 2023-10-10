package plain

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const (
	PagesDir          = "pages"
	StaticDir         = "public"
	IndexPageFileName = "index"
)

type route struct {
	filepath string
	urlpath  string
	data     []byte
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
		if fileExt != ".html" && fileExt != ".md" {
			continue
		}

		fileUrlPath := strings.ReplaceAll(filepath, fileExt, "")

		wdWithPagesDir := path.Join(wd, PagesDir)
		urlpath := strings.ReplaceAll(fileUrlPath, wdWithPagesDir, "")

		// It's safe because we handle case if it's directory above
		if strings.HasSuffix(urlpath, "index") {
			urlpath = strings.ReplaceAll(urlpath, "index", "")
		}

		f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}

		route := route{
			filepath: filepath,
			urlpath:  urlpath,
		}

		if fileExt == ".md" {
			route.data, err = parseMarkdown(filepath)
			if err != nil {
				return nil, err
			}
		} else {
			route.data, err = io.ReadAll(f)
			if err != nil {
				return nil, err
			}
		}

		f.Close()

		routes = append(routes, route)
	}

	return routes, nil
}

func parseMarkdown(fp string) ([]byte, error) {
	f, err := os.OpenFile(fp, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock)
	doc := p.Parse(data)

	r := html.NewRenderer(html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	})

	return markdown.Render(doc, r), nil
}
