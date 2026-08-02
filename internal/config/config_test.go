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

func TestQueryCommentsConfig(t *testing.T) {
	const raw = `{
  "version": "2",
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      "query_comments": {
        "enabled": true,
        "format": "marginalia",
        "tags": ["name", "cmd", "filename"]
      },
      "gen": {
        "go": {
          "out": "db"
        }
      }
    }
  ]
}`
	conf, err := ParseConfig(strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	got := conf.SQL[0].QueryComments
	want := QueryComments{
		Enabled: true,
		Format:  "marginalia",
		Tags:    []string{"name", "cmd", "filename"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("query comments differed (-want +got):\n%s", diff)
	}
}

func TestInvalidQueryCommentsConfig(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
		err  string
	}{
		{
			name: "invalid format",
			body: `"query_comments": {"enabled": true, "format": "unknown"}`,
			err:  "invalid query_comments format: unknown",
		},
		{
			name: "invalid tag",
			body: `"query_comments": {"enabled": true, "tags": ["name", "database"]}`,
			err:  "invalid query_comments tag: database",
		},
	} {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			raw := `{
  "version": "2",
  "sql": [
    {
      "schema": "schema.sql",
      "queries": "query.sql",
      "engine": "postgresql",
      ` + tt.body + `,
      "gen": {
        "go": {
          "out": "db"
        }
      }
    }
  ]
}`
			_, err := ParseConfig(strings.NewReader(raw))
			if err == nil {
				t.Fatalf("expected err; got nil")
			}
			if diff := cmp.Diff(tt.err, err.Error()); diff != "" {
				t.Errorf("error differed (-want +got):\n%s", diff)
			}
		})
	}
}
