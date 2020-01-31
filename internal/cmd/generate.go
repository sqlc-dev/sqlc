package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/mysql"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoPackages = `No packages are configured`

func printFileErr(stderr io.Writer, dir string, fileErr dinosql.FileErr) {
	filename := strings.TrimPrefix(fileErr.Filename, dir+"/")
	fmt.Fprintf(stderr, "%s:%d:%d: %s\n", filename, fileErr.Line, fileErr.Column, fileErr.Err)
}

func Generate(dir string, stderr io.Writer) (map[string]string, error) {
	blob, err := ioutil.ReadFile(filepath.Join(dir, "sqlc.json"))
	if err != nil {
		fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
		return nil, err
	}

	conf, err := config.ParseConfig(bytes.NewReader(blob))
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprintf(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprintf(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprintf(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing sqlc.json: %s\n", err)
		return nil, err
	}

	output := map[string]string{}
	errored := false

	for _, sql := range conf.SQL {
		combo := config.Combine(conf, sql)
		name := combo.Go.Package
		var result dinosql.Generateable

		// TODO: This feels like a hack that will bite us later
		sql.Schema = filepath.Join(dir, sql.Schema)
		sql.Queries = filepath.Join(dir, sql.Queries)

		switch sql.Engine {
		case config.EngineMySQL:
			// Experimental MySQL support
			q, err := mysql.GeneratePkg(name, sql.Schema, sql.Queries, combo)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						printFileErr(stderr, dir, fileErr)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}
			result = q

		case config.EnginePostgreSQL:
			c, err := dinosql.ParseCatalog(sql.Schema)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						printFileErr(stderr, dir, fileErr)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}

			q, err := dinosql.ParseQueries(c, sql)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						printFileErr(stderr, dir, fileErr)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing queries: %s\n", err)
				}
				errored = true
				continue
			}
			result = q

		}

		files, err := dinosql.Generate(result, combo)
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error generating code: %s\n", err)
			errored = true
			continue
		}

		for n, source := range files {
			filename := filepath.Join(dir, combo.Go.Out, n)
			output[filename] = source
		}
	}

	if errored {
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}
