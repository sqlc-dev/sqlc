package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error generating env: %s", err)
		os.Exit(2)
	}
}

func run() error {
	env := os.Environ()
	blob, err := json.Marshal(env)
	if err != nil {
		return err
	}
	resp := &plugin.CodeGenResponse{
		Files: []*plugin.File{
			{
				Name:     "env.json",
				Contents: append(blob, '\n'),
			},
		},
	}
	respBlob, err := resp.MarshalVT()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(os.Stdout)
	if _, err := w.Write(respBlob); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return nil
}
