package dinosql

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	pg "github.com/lfittl/pg_query_go"
)

// Make sure that each example has one SQL code block and one Go code block.
// Ensure that both are valid.
func TestExamples(t *testing.T) {
	matches, err := filepath.Glob("examples/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, match := range matches {
		m := match
		t.Run(m, func(t *testing.T) {
			blob, err := ioutil.ReadFile(m)
			if err != nil {
				t.Fatal(err)
			}
			var sql, goc string
			var captureSQL, captureGo bool
			for _, line := range strings.Split(string(blob), "\n") {
				if strings.HasPrefix(line, "```sql") {
					captureSQL = true
					continue
				}
				if strings.HasPrefix(line, "```go") {
					captureGo = true
					continue
				}
				if strings.HasPrefix(line, "```") {
					captureSQL = false
					captureGo = false
					continue
				}
				if captureSQL {
					sql += line + "\n"
				}
				if captureGo {
					goc += line + "\n"
				}
			}
			if _, err := pg.Parse(sql); err != nil {
				t.Errorf("could not parse SQL: %s", err)
			}
			if _, err := parser.ParseFile(token.NewFileSet(), "", goc, parser.AllErrors); err != nil {
				t.Errorf("could not parse Go: %s", err)
			}
		})
	}
}
