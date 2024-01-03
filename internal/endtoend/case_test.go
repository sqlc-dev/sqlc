package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type Testcase struct {
	Name       string
	Path       string
	ConfigName string
	Stderr     []byte
	Exec       *Exec
}

type ExecMeta struct {
	InvalidSchema bool `json:"invalid_schema"`
}

type Exec struct {
	Command  string            `json:"command"`
	Contexts []string          `json:"contexts"`
	Process  string            `json:"process"`
	OS       []string          `json:"os"`
	Env      map[string]string `json:"env"`
	Meta     ExecMeta          `json:"meta"`
}

func parseStderr(t *testing.T, dir, testctx string) []byte {
	t.Helper()
	paths := []string{
		filepath.Join(dir, "stderr", fmt.Sprintf("%s.txt", testctx)),
		filepath.Join(dir, "stderr.txt"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			blob, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			return blob
		}
	}
	return nil
}

func parseExec(t *testing.T, dir string) *Exec {
	t.Helper()
	path := filepath.Join(dir, "exec.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	var e Exec
	blob, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s: %s", path, err)
	}
	if err := json.Unmarshal(blob, &e); err != nil {
		t.Fatalf("%s: %s", path, err)
	}
	if e.Command == "" {
		e.Command = "generate"
	}
	return &e
}

func FindTests(t *testing.T, root, testctx string) []*Testcase {
	var tcs []*Testcase
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.json" || info.Name() == "sqlc.yaml" || info.Name() == "sqlc.yml" {
			dir := filepath.Dir(path)
			tcs = append(tcs, &Testcase{
				Path:       dir,
				Name:       strings.TrimPrefix(dir, root+string(filepath.Separator)),
				ConfigName: info.Name(),
				Stderr:     parseStderr(t, dir, testctx),
				Exec:       parseExec(t, dir),
			})
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return tcs
}
