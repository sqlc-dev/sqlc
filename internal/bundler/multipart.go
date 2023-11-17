package bundler

import (
	"os"
	"path/filepath"

	"github.com/sqlc-dev/sqlc/internal/plugin"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func readFiles(dir string, paths []string) ([]*plugin.File, error) {
	files, err := sqlpath.Glob(paths)
	if err != nil {
		return nil, err
	}
	var out []*plugin.File
	for _, file := range files {
		f, err := readFile(dir, file)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func readFile(dir string, path string) (*plugin.File, error) {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return nil, err
	}
	blob, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &plugin.File{
		Name:     rel,
		Contents: blob,
	}, nil
}
