package sqlpath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyleconroy/sqlc/internal/migrations"
)

// Return a list of SQL files in the listed paths. Only includes files ending
// in .sql. Omits hidden files, directories, and migrations.
func Glob(paths []string) ([]string, error) {
	var files []string
	for _, path := range paths {
		f, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("path %s does not exist", path)
		}
		if f.IsDir() {
			listing, err := ioutil.ReadDir(path)
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
