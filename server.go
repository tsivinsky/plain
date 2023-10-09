package plain

import (
	"fmt"
	"net/http"
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

func (s *Server) getStaticFile(fp string) string {
	return path.Join(s.WorkingDir, StaticDir, fp)
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	err := http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		route := matchRoute(r, s.routes)
		if route == nil {
			http.ServeFile(w, r, s.getStaticFile(r.URL.Path))
		} else {
			http.ServeFile(w, r, route.filepath)
		}
	}))

	return err
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
