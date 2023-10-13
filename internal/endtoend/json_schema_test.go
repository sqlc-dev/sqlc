package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

type conf struct {
	Version string `json:"version"`
}

func loadSchema(t *testing.T, path string) *gojsonschema.Schema {
	t.Helper()

	schemaBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	loader := gojsonschema.NewStringLoader(string(schemaBytes))
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		t.Fatalf("invalid schema: %s", err)
	}
	return schema
}

func TestJsonSchema(t *testing.T) {
	t.Parallel()

	schemaOne := loadSchema(t, filepath.Join("..", "config", "v_one.json"))
	schemaTwo := loadSchema(t, filepath.Join("..", "config", "v_two.json"))

	srcs := []string{
		filepath.Join("..", "..", "examples"),
		filepath.Join("testdata"),
	}

	for _, dir := range srcs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Base(path) != "sqlc.json" {
				return nil
			}
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				contents, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var c conf
				if err := json.Unmarshal(contents, &c); err != nil {
					t.Fatal(err)
				}
				l := gojsonschema.NewStringLoader(string(contents))
				switch c.Version {
				case "1":
					if _, err := schemaOne.Validate(l); err != nil {
						t.Fatal(err)
					}
				case "2":
					if _, err := schemaTwo.Validate(l); err != nil {
						t.Fatal(err)
					}
				default:
					t.Fatalf("unknown schema version: %s", c.Version)
				}
			})
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	}
}
