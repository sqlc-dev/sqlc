package bundler

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
)

func writeContents(w *multipart.Writer, file string, conf *config.Config) error {
	refs := map[string]struct{}{}
	refs[filepath.Base(file)] = struct{}{}

	for _, pkg := range conf.SQL {
		for _, paths := range []config.Paths{pkg.Schema, pkg.Queries} {
			files, err := sqlpath.Glob(paths)
			if err != nil {
				return err
			}
			for _, file := range files {
				refs[file] = struct{}{}
			}
		}
	}

	for file, _ := range refs {
		if err := addPart(w, file); err != nil {
			return err
		}
	}

	params, err := projectMetadata()
	if err != nil {
		return err
	}
	for _, val := range params {
		if err = w.WriteField(val[0], val[1]); err != nil {
			return err
		}
	}
	return nil
}

func addPart(w *multipart.Writer, file string) error {
	h, err := os.Open(file)
	if err != nil {
		return err
	}
	defer h.Close()
	part, err := w.CreateFormFile("contents", file)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, h)
	if err != nil {
		return err
	}
	return nil
}
