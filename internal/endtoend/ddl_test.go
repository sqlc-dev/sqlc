package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestValidSchema(t *testing.T) {
	for _, replay := range FindTests(t, "testdata", "base") {

		if replay.Exec != nil {
			if replay.Exec.Meta.InvalidSchema {
				continue
			}
		}

		file := filepath.Join(replay.Path, replay.ConfigName)
		rd, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		conf, err := config.ParseConfig(rd)
		if err != nil {
			t.Fatal(err)
		}

		for j, pkg := range conf.SQL {
			switch pkg.Engine {
			case config.EnginePostgreSQL:
				// pass
			case config.EngineMySQL:
				// pass
			default:
				continue
			}
			t.Run(fmt.Sprintf("endtoend-%s-%d", file, j), func(t *testing.T) {
				t.Parallel()

				var schema []string
				for _, path := range pkg.Schema {
					schema = append(schema, filepath.Join(filepath.Dir(file), path))
				}

				switch pkg.Engine {
				case config.EnginePostgreSQL:
					local.PostgreSQL(t, schema)
				case config.EngineMySQL:
					local.MySQL(t, schema)
				}
			})
		}
	}
}
