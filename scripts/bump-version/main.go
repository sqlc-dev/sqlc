package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func main() {
	curr := flag.String("c", "", "current version")
	next := flag.String("n", "", "next version")
	write := flag.Bool("w", false, "write out changes")
	flag.Parse()
	if err := run(*curr, *next, *write); err != nil {
		log.Fatal(err)
	}
}

func run(current, next string, realmode bool) error {
	write := func(path, old, new string) error {
		if realmode {
			if err := os.WriteFile(path, []byte(new), 0644); err != nil {
				return fmt.Errorf("write error: %s: %w", path, err)
			}
		} else {
			if diff := cmp.Diff(old, new); diff != "" {
				log.Printf("%s: %s\n", path, diff)
			}
		}
		return nil
	}

	{
		path := filepath.Join(".github", "ISSUE_TEMPLATE", "BUG_REPORT.yml")
		c, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		old := string(c)
		if !strings.Contains(old, "- "+next) {
			item := "- " + current
			new := strings.ReplaceAll(old, item, "- "+next+"\n        "+item)
			if err := write(path, old, new); err != nil {
				return err
			}
		}
	}

	{
		path := filepath.Join("docs", "overview", "install.md")
		c, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		old := string(c)
		new := strings.ReplaceAll(old, "v"+current, "v"+next)
		new = strings.ReplaceAll(new, "sqlc_"+current, "sqlc_"+next)
		if err := write(path, old, new); err != nil {
			return err
		}
	}

	{
		path := filepath.Join("internal", "info", "facts.go")
		c, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		old := string(c)
		new := strings.ReplaceAll(old, "v"+current, "v"+next)
		if err := write(path, old, new); err != nil {
			return err
		}
	}

	{
		path := filepath.Join("docs", "conf.py")
		c, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		old := string(c)
		new := strings.ReplaceAll(old, "release = '"+current, "release = '"+next)
		if err := write(path, old, new); err != nil {
			return err
		}
	}

	walker := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		switch filepath.Ext(path) {
		case ".go", ".kt", ".py", ".json", ".md":
			c, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			old := string(c)
			new := strings.ReplaceAll(old,
				`"sqlc_version": "v`+current,
				`"sqlc_version": "v`+next)
			new = strings.ReplaceAll(new,
				`sqlc-version: "`+current,
				`sqlc-version: "`+next)
			new = strings.ReplaceAll(new,
				`sqlc-version: '`+current,
				`sqlc-version: '`+next)
			new = strings.ReplaceAll(new, "sqlc v"+current, "sqlc v"+next)
			new = strings.ReplaceAll(new, "SQLC_VERSION=v"+current, "SQLC_VERSION=v"+next)
			if err := write(path, old, new); err != nil {
				return err
			}
		default:
		}
		return nil
	}

	{
		p := filepath.Join("internal", "endtoend", "testdata")
		if err := filepath.Walk(p, walker); err != nil {
			return err
		}
	}

	{
		p := filepath.Join("examples")
		if err := filepath.Walk(p, walker); err != nil {
			return err
		}
	}

	{
		p := filepath.Join("docs")
		if err := filepath.Walk(p, walker); err != nil {
			return err
		}
	}

	return nil
}
