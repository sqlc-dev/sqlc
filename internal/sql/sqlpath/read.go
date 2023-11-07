package sqlpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/migrations"
)

// Return a list of SQL files in the listed paths. Only includes files ending
// in .sql. Omits hidden files, directories, and migrations.
func Glob(paths []string) ([]string, error) {
	paths, err := expandGlobs(paths)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, path := range paths {
		f, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("path %s does not exist", path)
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
			files = append(files, path)
		}
	}
	var sqlFiles []string //nolint:prealloc // can be empty
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

func expandGlobs(paths []string) ([]string, error) {
	expandedPatterns := make([]string, 0, len(paths))
	for _, pattern := range paths {
		expansion, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to expand pattern %q: %w", pattern, err)
		}
		if len(expansion) == 0 {
			fi, err := os.Lstat(pattern)
			if err != nil {
				return nil, fmt.Errorf("failed to stat path %q: %w", pattern, err)
			}
			if fi == nil {
				return nil, fmt.Errorf("failed to stat path %q: %w", pattern, os.ErrNotExist)
			}
			var isFilepath bool
			for _, mask := range []os.FileMode{os.ModeDir, os.ModeSymlink, os.FileMode(0x400)} {
				if fi.Mode()&mask == 0 {
					isFilepath = true
					break
				}
			}
			if !isFilepath {
				continue
			}
		}
		expandedPatterns = append(expandedPatterns, expansion...)
	}
	return expandedPatterns, nil
}
