package plain

import (
	"os"
	"path"
)

type Options struct {
	Host         string
	Port         int
	WorkingDir   string
	ReadPageFile ReadPageFileFunc
}

// New initializes server with provided options
func New(o Options) (*Server, error) {
	pagesPath := path.Join(o.WorkingDir, pagesDir)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		return nil, err
	}

	s := &Server{
		Host:         o.Host,
		Port:         o.Port,
		WorkingDir:   o.WorkingDir,
		pagesPath:    pagesPath,
		readPageFile: o.ReadPageFile,
	}

	var err error

	s.routes, err = getRoutes(s.pagesPath, s.WorkingDir, s.readPageFile)
	if err != nil {
		return nil, err
	}

	return s, nil
}
