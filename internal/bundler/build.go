package bundler

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
)

func Build(file string, conf *config.Config) ([]byte, error) {
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

	// TODO: Checksum
	// TODO: Gzip this
	b := bytes.NewBuffer([]byte{})
	// bz := gzip.NewWriter(b)
	w := tar.NewWriter(b)
	defer w.Close()

	for file, _ := range refs {
		if err := addFile(w, file); err != nil {
			return nil, err
		}
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	// if err := bz.Flush(); err != nil {
	// 	return nil, err
	// }

	return b.Bytes(), nil
}

func addFile(w *tar.Writer, file string) error {
	h, err := os.Open(file)
	if err != nil {
		return err
	}
	defer h.Close()
	info, err := h.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = file
	if err := w.WriteHeader(header); err != nil {
		return err
	}
	body, err := io.ReadAll(h)
	if err != nil {
		return err
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}
