package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	// Assume it exists
	loc := flag.Arg(0)

	dir := filepath.Join("internal", "codegen", "golang")
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newdir := filepath.Join(loc, "internal")
		newpath := strings.Replace(path, dir, newdir, 1)

		os.MkdirAll(filepath.Dir(newpath), 0755)

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"`),
			[]byte(`"github.com/sqlc-dev/sqlc-gen-go/internal/opts"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/plugin"`),
			[]byte(`"github.com/sqlc-dev/plugin-sdk-go/plugin"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/codegen/sdk"`),
			[]byte(`"github.com/sqlc-dev/plugin-sdk-go/sdk"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/metadata"`),
			[]byte(`"github.com/sqlc-dev/plugin-sdk-go/metadata"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/pattern"`),
			[]byte(`"github.com/sqlc-dev/plugin-sdk-go/pattern"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/debug"`),
			[]byte(`"github.com/sqlc-dev/sqlc-gen-go/internal/debug"`))

		contents = bytes.ReplaceAll(contents,
			[]byte(`"github.com/sqlc-dev/sqlc/internal/inflection"`),
			[]byte(`"github.com/sqlc-dev/sqlc-gen-go/internal/inflection"`))

		if err := os.WriteFile(newpath, contents, 0644); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path: %v\n", err)
		return
	}

	{
		path := filepath.Join("internal", "inflection", "singular.go")
		contents, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}
		newpath := filepath.Join(loc, "internal", "inflection", "singular.go")
		if err := os.WriteFile(newpath, contents, 0644); err != nil {
			log.Fatal(err)
		}
	}

}
