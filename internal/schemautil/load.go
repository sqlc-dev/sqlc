package schemautil

import (
	"fmt"
	"os"

	"github.com/sqlc-dev/sqlc/internal/migrations"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
)

// LoadSchemasForApply expands globs, preprocesses each schema in order, and
// reports any warnings through warn. The returned DDL is suitable for callers
// that will apply schema text to a live database.
func LoadSchemasForApply(globs []string, engine string, warn func(string)) ([]string, error) {
	files, err := sqlpath.Glob(globs)
	if err != nil {
		return nil, err
	}

	ddl := make([]string, 0, len(files))
	for _, schema := range files {
		contents, err := os.ReadFile(schema)
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}
		ddlText, warnings, err := migrations.PreprocessSchemaForApply(string(contents), engine)
		if err != nil {
			return nil, err
		}
		for _, warning := range warnings {
			if warn != nil {
				warn(warning)
			}
		}
		ddl = append(ddl, ddlText)
	}

	return ddl, nil
}
