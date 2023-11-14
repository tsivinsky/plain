package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/tsivinsky/plain"
)

var (
	port  = flag.Int("p", 5000, "Port to run application on")
	host  = flag.String("H", "localhost", "Host to run application on")
	watch = flag.Bool("w", false, "Watch html files for changes")
)

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

func readPageFile(fp string) ([]byte, error) {
	ext := path.Ext(fp)

	if ext == ".md" {
		return parseMarkdown(fp)
	}

	if ext == ".html" {

		f, err := os.OpenFile(fp, os.O_RDONLY, 0644)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return io.ReadAll(f)
	}

	return nil, nil
}

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
		Host:         *host,
		Port:         *port,
		WorkingDir:   wd,
		ReadPageFile: readPageFile,
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
