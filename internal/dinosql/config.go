package dinosql

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PackageSettings struct {
	Name                string           `json:"name"`
	Path                string           `json:"path"`
	Schema              string           `json:"schema"`
	Queries             string           `json:"queries"`
	EmitPreparedQueries bool             `json:"emit_prepared_queries"`
	EmitJSONTags        bool             `json:"emit_json_tags"`
	Overrides           []ColumnOverride `json:"overrides"`
}

type ColumnOverride struct {
	// fully qualified name of the column, e.g. `accounts.id`
	Column string `json:"column"`

	// fully qualified name of the Go type, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType string `json:"go_type"`

	columnName    string
	tableName     string
	goTypeName    string
	goPackageName string
}

func (o *ColumnOverride) Parse() error {
	colParts := strings.Split(o.Column, ".")
	if len(colParts) != 2 {
		return fmt.Errorf("Package override `column` specifier '%s' is not the proper format, expected 'colname.tablename'", o.Column)
	}

	lastDot := strings.LastIndex(o.GoType, ".")
	if lastDot == -1 {
		return fmt.Errorf("Package override `go_type` specificier '%s' is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
	}
	lastSlash := strings.LastIndex(o.GoType, "/")
	if lastSlash == -1 {
		return fmt.Errorf("Package override `go_type` specificier '%s' is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
	}

	o.columnName = colParts[1]
	o.tableName = colParts[0]
	o.goTypeName = o.GoType[lastSlash+1:]
	o.goPackageName = o.GoType[:lastDot]
	isPointer := o.GoType[0] == '*'
	if isPointer {
		o.goPackageName = o.goPackageName[1:]
		o.goTypeName = "*" + o.goTypeName
	}
	return nil
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
	for i := range config.Packages {
		for j := range config.Packages[i].Overrides {
			if err := config.Packages[i].Overrides[j].Parse(); err != nil {
				return config, err
			}
		}
	}
	return config, nil
}
