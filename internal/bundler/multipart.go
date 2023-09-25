package bundler

import (
	"os"
	"path/filepath"

	"github.com/sqlc-dev/sqlc/internal/config"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

func readInputs(file string, conf *config.Config) ([]*pb.File, error) {
	refs := map[string]struct{}{}
	refs[filepath.Base(file)] = struct{}{}

	for _, pkg := range conf.SQL {
		for _, paths := range []config.Paths{pkg.Schema, pkg.Queries} {
			files, err := sqlpath.Glob(paths)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				refs[file] = struct{}{}
			}
		}
	}

	var files []*pb.File
	for file, _ := range refs {
		contents, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		files = append(files, &pb.File{
			Name:      file,
			MediaType: "application/octet-stream",
			Contents:  contents,
		})
	}
	return files, nil
}

func readOutputs(dir string, output map[string]string) ([]*pb.File, error) {
	var files []*pb.File
	for filename, contents := range output {
		rel, err := filepath.Rel(dir, filename)
		if err != nil {
			return nil, err
		}
		files = append(files, &pb.File{
			Name:      rel,
			MediaType: "application/octet-stream",
			Contents:  []byte(contents),
		})
	}
	return files, nil
}
