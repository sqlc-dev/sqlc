package main

import (
	"context"
	"github.com/sqlc-dev/sqlc/internal/cmd"
	"os"
	"path/filepath"
)

func main() {
	stderr := os.Stderr
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd = filepath.Join(wd, "examples/stand_alone_name")
	output, err := cmd.Generate(context.TODO(), cmd.Env{}, wd, "sqlc.yaml", stderr)
	if err != nil {
		panic(err)
	}
	for filename, source := range output {
		os.MkdirAll(filepath.Dir(filename), 0755)
		if err := os.WriteFile(filename, []byte(source), 0644); err != nil {
			panic(err)
		}
	}
}
