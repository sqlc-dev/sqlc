package config

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const missingVersion = `{
}`

const missingPackages = `{
  "version": "1"
}`

const unknownVersion = `{
  "version": "foo"
}`

const unknownFields = `{
  "version": "1",
  "foo": "bar"
}`

func TestBadConfigs(t *testing.T) {
	for _, test := range []struct {
		name string
		err  string
		json string
	}{
		{
			"missing version",
			"no version number",
			missingVersion,
		},
		{
			"missing packages",
			"no packages",
			missingPackages,
		},
		{
			"unknown version",
			"invalid version number",
			unknownVersion,
		},
		{
			"unknown fields",
			`yaml: unmarshal errors:
  line 3: field foo not found in type config.V1GenerateSettings`,
			unknownFields,
		},
	} {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseConfig(strings.NewReader(tt.json))
			if err == nil {
				t.Fatalf("expected err; got nil")
			}
			if diff := cmp.Diff(err.Error(), tt.err); diff != "" {
				t.Errorf("differed (-want +got):\n%s", diff)
			}
		})
	}
}

const validConfigOne = `{
  "version": "1"
  "packages": []
}`

func FuzzConfig(f *testing.F) {
	f.Add(validConfigOne)
	f.Fuzz(func(t *testing.T, orig string) {
		ParseConfig(strings.NewReader(orig))
	})
}

func TestInvalidConfig(t *testing.T) {
	err := Validate(&Config{
		SQL: []SQL{{
			Database: &Database{
				URI:     "",
				Managed: false,
			},
		}},
	})
	if err == nil {
		t.Errorf("expected err; got nil")
	}
}
