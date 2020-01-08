package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/mysql"
)

func TestCodeGeneration(t *testing.T) {
	// Change to the top-level directory of the project
	os.Chdir(filepath.Join("..", ".."))

	rd, err := os.Open("sqlc.json")
	if err != nil {
		t.Fatal(err)
	}

	conf, err := dinosql.ParseConfig(rd)
	if err != nil {
		t.Fatal(err)
	}

	for n, s := range conf.PackageMap {
		pkg := s
		t.Run(n, func(t *testing.T) {
			var result dinosql.Generateable
			switch pkg.Engine {
			case dinosql.EngineMySQL:
				q, err := mysql.GeneratePkg(pkg.Name, pkg.Schema, pkg.Queries, conf)
				if err != nil {
					t.Fatal(err)
				}
				result = q
			case dinosql.EnginePostgreSQL:
				c, err := dinosql.ParseCatalog(pkg.Schema)
				if err != nil {
					fmt.Printf("%#v\n", err)
					t.Fatal(err)
				}
				q, err := dinosql.ParseQueries(c, pkg)
				if err != nil {
					t.Fatal(err)
				}
				result = q
			}
			output, err := dinosql.Generate(result, conf)
			if err != nil {
				t.Fatal(err)
			}
			cmpDirectory(t, pkg.Path, output)
		})
	}

}

func cmpDirectory(t *testing.T, dir string, actual map[string]string) {
	t.Helper()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("error reading dir %s: %s", dir, err)
	}

	expected := map[string]string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(file.Name(), "_test.go") {
			continue
		}
		blob, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		expected[file.Name()] = string(blob)
	}

	if !cmp.Equal(expected, actual) {
		t.Errorf("%s contents differ", dir)
		for name, contents := range expected {
			if actual[name] == "" {
				t.Errorf("%s is empty", name)
				continue
			}
			if diff := cmp.Diff(contents, actual[name]); diff != "" {
				t.Errorf("%s differed (-want +got):\n%s", name, diff)
			}
		}
	}
}
