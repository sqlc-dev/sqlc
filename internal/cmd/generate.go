package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

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

func Generate(dir string, stderr io.Writer) (map[string]string, error) {
	blob, err := ioutil.ReadFile(filepath.Join(dir, "sqlc.json"))
	if err != nil {
		fmt.Fprintln(stderr, "error parsing sqlc.json: file does not exist")
		return nil, err
	}

	settings, err := dinosql.ParseConfig(bytes.NewReader(blob))
	if err != nil {
		switch err {
		case dinosql.ErrMissingVersion:
			fmt.Fprintf(stderr, errMessageNoVersion)
		case dinosql.ErrUnknownVersion:
			fmt.Fprintf(stderr, errMessageUnknownVersion)
		case dinosql.ErrNoPackages:
			fmt.Fprintf(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing sqlc.json: %s\n", err)
		return nil, err
	}

	output := map[string]string{}
	errored := false

	for _, pkg := range settings.Packages {
		name := pkg.Name
		combo := dinosql.Combine(settings, pkg)
		var result dinosql.Generateable

		// TODO: This feels like a hack that will bite us later
		pkg.Schema = filepath.Join(dir, pkg.Schema)
		pkg.Queries = filepath.Join(dir, pkg.Queries)

		switch pkg.Engine {

		case dinosql.EngineMySQL:
			// Experimental MySQL support
			q, err := mysql.GeneratePkg(name, pkg.Schema, pkg.Queries, combo)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}
			result = q

		case dinosql.EnginePostgreSQL:
			c, err := dinosql.ParseCatalog(pkg.Schema)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing schema: %s\n", err)
				}
				errored = true
				continue
			}

			q, err := dinosql.ParseQueries(c, pkg)
			if err != nil {
				fmt.Fprintf(stderr, "# package %s\n", name)
				if parserErr, ok := err.(*dinosql.ParserErr); ok {
					for _, fileErr := range parserErr.Errs {
						fmt.Fprintf(stderr, "%s:%d:%d: %s\n", fileErr.Filename, fileErr.Line, fileErr.Column, fileErr.Err)
					}
				} else {
					fmt.Fprintf(stderr, "error parsing queries: %s\n", err)
				}
				errored = true
				continue
			}
			result = q

		}

		var files map[string]string
		switch pkg.Language {
		case dinosql.LanguageGo:
			files, err = dinosql.Generate(result, combo)
		case dinosql.LanguageKotlin:
			ktRes, ok := result.(dinosql.KtGenerateable)
			if !ok {
				err = fmt.Errorf("kotlin not supported for engine %s", pkg.Engine)
				break
			}
			files, err = dinosql.KtGenerate(ktRes, combo)
		}
		if err != nil {
			fmt.Fprintf(stderr, "# package %s\n", name)
			fmt.Fprintf(stderr, "error generating code: %s\n", err)
			errored = true
			continue
		}

		for n, source := range files {
			filename := filepath.Join(dir, pkg.Path, n)
			output[filename] = source
		}
	}

	if errored {
		return nil, fmt.Errorf("errored")
	}
	return output, nil
}
