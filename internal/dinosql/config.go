package dinosql

import (
	"encoding/json"
	"errors"
	"io"
)

type PackageSettings struct {
	Name                string `json:"name"`
	Path                string `json:"path"`
	Schema              string `json:"schema"`
	Queries             string `json:"queries"`
	EmitPreparedQueries bool   `json:"emit_prepared_queries"`
	EmitJSONTags        bool   `json:"emit_json_tags"`
}

type GenerateSettings struct {
	Version   string            `json:"version"`
	Packages  []PackageSettings `json:"packages"`
	Overrides []TypeOverride    `json:"overrides,omitempty"`
	Rename    map[string]string `json:"rename,omitempty"`
}

var ErrMissingVersion = errors.New("no version number")
var ErrUnknownVersion = errors.New("invalid version number")
var ErrNoPackages = errors.New("no packages")

func ParseConfig(rd io.Reader) (GenerateSettings, error) {
	dec := json.NewDecoder(rd)
	dec.DisallowUnknownFields()
	var config GenerateSettings
	if err := dec.Decode(&config); err != nil {
		return config, err
	}
	if config.Version == "" {
		return config, ErrMissingVersion
	}
	if config.Version != "1" {
		return config, ErrUnknownVersion
	}
	if len(config.Packages) == 0 {
		return config, ErrNoPackages
	}
	return config, nil
}
