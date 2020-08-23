package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func regenerate(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "sqlc.json") || strings.HasSuffix(path, "sqlc.yaml") {
			cwd := filepath.Dir(path)
			cmd := exec.Command("sqlc-dev", "generate")
			cmd.Dir = cwd
			failed := cmd.Run()
			if _, err := os.Stat(filepath.Join(cwd, "stderr.txt")); os.IsNotExist(err) && failed != nil {
				return fmt.Errorf("%s: sqlc-dev generate failed", cwd)
			}
		}
		return nil
	})
}

func main() {
	dirs := []string{
		filepath.Join("internal", "endtoend", "testdata"),
		filepath.Join("examples"),
	}
	for _, d := range dirs {
		if err := regenerate(d); err != nil {
			log.Fatal(err)
		}
	}
}
