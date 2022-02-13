package bundler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/debug"
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
	b := bytes.NewBuffer([]byte{})
	gw := gzip.NewWriter(b)
	defer gw.Close()
	w := tar.NewWriter(gw)
	defer w.Close()

	for file, _ := range refs {
		if err := addFile(w, file); err != nil {
			return nil, err
		}
	}
	if err := addMetadata(w); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	if err := gw.Flush(); err != nil {
		return nil, err
	}
	w.Close()
	gw.Close()
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
	debug.Dump(header)
	header.Name = file
	if err := w.WriteHeader(header); err != nil {
		return err
	}
	// copy the file data to the tarball
	if _, err := io.Copy(w, h); err != nil {
		return err
	}
	return nil
}
