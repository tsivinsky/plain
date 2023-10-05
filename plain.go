package plain

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

// struct that stores server config and runs it
type Server struct {
	Host       string
	Port       int
	WorkingDir string

	pagesPath string
	routes    []route
}

// Run runs server on Host:Port address
func (s *Server) Run() error {
	s.pagesPath = path.Join(s.WorkingDir, PagesDir)
	if _, err := os.Stat(s.pagesPath); os.IsNotExist(err) {
		return err
	}

	var err error

	s.routes, err = getRoutes(s.pagesPath, s.WorkingDir)
	if err != nil {
		return err
	}

	portStr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err = http.ListenAndServe(portStr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		route := matchRoute(r, s.routes)
		if route == nil {
			fp := path.Join(s.WorkingDir, StaticDir, r.URL.Path)
			http.ServeFile(w, r, fp)
			return
		}

		err := renderHTMLFile(w, route.filepath)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
		}
	}))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(s.pagesPath)
	if err != nil {
		return err
	}

	fmt.Println("[plain]: started watching files")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				s.routes, err = getRoutes(s.pagesPath, s.WorkingDir)
				if err != nil {
					fmt.Printf("Error happened while updating routes list on file change: %s\n", err.Error())
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			fmt.Printf("Error happened while watching files for changing: %s\n", err.Error())
		}
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
