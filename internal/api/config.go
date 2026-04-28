package api

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sqlc-dev/sqlc/internal/config"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 or 2 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The supported version can only be "1" or "2".
`

const errMessageNoPackages = `No packages are configured`

func readConfig(stderr io.Writer, dir, filename string) (string, *config.Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		var yamlMissing, jsonMissing, ymlMissing bool
		yamlPath := filepath.Join(dir, "sqlc.yaml")
		ymlPath := filepath.Join(dir, "sqlc.yml")
		jsonPath := filepath.Join(dir, "sqlc.json")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			yamlMissing = true
		}
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			jsonMissing = true
		}

		if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
			ymlMissing = true
		}

		if yamlMissing && ymlMissing && jsonMissing {
			fmt.Fprintln(stderr, "error parsing configuration files. sqlc.(yaml|yml) or sqlc.json: file does not exist")
			return "", nil, errors.New("config file missing")
		}

		if (!yamlMissing || !ymlMissing) && !jsonMissing {
			fmt.Fprintln(stderr, "error: both sqlc.json and sqlc.(yaml|yml) files present")
			return "", nil, errors.New("sqlc.json and sqlc.(yaml|yml) present")
		}

		if jsonMissing {
			if yamlMissing {
				configPath = ymlPath
			} else {
				configPath = yamlPath
			}
		} else {
			configPath = jsonPath
		}
	}

	base := filepath.Base(configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}
	defer file.Close()

	conf, err := config.ParseConfig(file)
	if err != nil {
		switch err {
		case config.ErrMissingVersion:
			fmt.Fprint(stderr, errMessageNoVersion)
		case config.ErrUnknownVersion:
			fmt.Fprint(stderr, errMessageUnknownVersion)
		case config.ErrNoPackages:
			fmt.Fprint(stderr, errMessageNoPackages)
		}
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return configPath, &conf, nil
}
