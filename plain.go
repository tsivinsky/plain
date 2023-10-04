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
	Watch      bool
}

// Run runs server on Host:Port address
func (s *Server) Run() error {
	pagesPath := path.Join(s.WorkingDir, PagesDir)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		return err
	}

	routes, err := getRoutes(pagesPath, s.WorkingDir)
	if err != nil {
		return err
	}

	if s.Watch {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		fmt.Println("[plain]: started watching files")

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
						routes, err = getRoutes(pagesPath, s.WorkingDir)
						if err != nil {
							fmt.Printf("Error happened while updating routes list on file change: %s\n", err.Error())
						}
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					fmt.Printf("Error happened while watching files for changing: %s\n", err.Error())
				}
			}
		}()

		err = watcher.Add(pagesPath)
		if err != nil {
			return err
		}
	}

	portStr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err = http.ListenAndServe(portStr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		route := matchRoute(r, routes)
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
