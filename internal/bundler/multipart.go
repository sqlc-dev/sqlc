package bundler

import (
	"os"
	"path/filepath"

	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func readFiles(dir string, paths []string) ([]*pb.File, error) {
	files, err := sqlpath.Glob(paths)
	if err != nil {
		return nil, err
	}
	var out []*pb.File
	for _, file := range files {
		f, err := readFile(dir, file)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func readFile(dir string, path string) (*pb.File, error) {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return nil, err
	}
	blob, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &pb.File{
		Name:     rel,
		Contents: blob,
	}, nil
}
