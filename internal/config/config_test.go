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
	cases := []struct {
		name    string
		db      *Database
		wantErr bool
	}{
		{"no uri no managed no testcontainers", &Database{}, true},
		{"testcontainers_image empty", &Database{TestcontainersImage: ""}, true},
		{"testcontainers_image set", &Database{TestcontainersImage: "postgres:18-alpine"}, false},
		{"managed true", &Database{Managed: true}, false},
		{"uri set", &Database{URI: "postgres://localhost/db"}, false},
		{"uri and testcontainers_image both set", &Database{URI: "postgres://localhost/db", TestcontainersImage: "postgres:18-alpine"}, true},
		{"managed and testcontainers_image both set", &Database{Managed: true, TestcontainersImage: "postgres:18-alpine"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(&Config{SQL: []SQL{{Database: tc.db}}})
			if tc.wantErr && err == nil {
				t.Errorf("expected err; got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected err: %v", err)
			}
		})
	}
}
