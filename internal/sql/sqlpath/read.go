package sqlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/migrations"
)

// Return a list of SQL files in the listed paths.
//
// Only includes files ending in .sql. Omits hidden files, directories, and
// down migrations.

// If a path contains *, ?, [, or ], treat the path as a pattern and expand it
// filepath.Glob.
func Glob(patterns []string) ([]string, error) {
	var files, paths []string
	for _, pattern := range patterns {
		if strings.ContainsAny(pattern, "*?[]") {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return nil, err
			}
			// if len(matches) == 0 {
			// 	slog.Warn("zero files matched", "pattern", pattern)
			// }
			paths = append(paths, matches...)
		} else {
			paths = append(paths, pattern)
		}
	}
	for _, path := range paths {
		f, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("path error: %w", err)
		}
		if f.IsDir() {
			listing, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			for _, f := range listing {
				files = append(files, filepath.Join(path, f.Name()))
			}
		} else {
			files = append(files, filepath.Clean(path))
		}
	}
	var sqlFiles []string
	for _, file := range files {
		if !strings.HasSuffix(file, ".sql") {
			continue
		}
		if strings.HasPrefix(filepath.Base(file), ".") {
			continue
		}
		if migrations.IsDown(filepath.Base(file)) {
			continue
		}
		sqlFiles = append(sqlFiles, file)
	}
	return sqlFiles, nil
}
