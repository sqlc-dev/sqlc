package cmd

import (
	"context"
	"fmt"
	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func InitTmpl(ctx context.Context, e Env, dir, filename string, stderr io.Writer) (map[string]string, error) {
	_, conf, err := readConfig(stderr, dir, filename)
	if err != nil {
		return nil, err
	}

	for _, sqlConfig := range conf.SQL {
		templatesPath := sqlConfig.Gen.Go.TemplatePath
		if templatesPath != "" {
			templateFS := golang.GetTemplates()

			templateDir, err := fs.ReadDir(templateFS, ".")
			if err != nil {
				return nil, err
			}
			templateName := templateDir[0].Name()

			sub, err := fs.Sub(templateFS, templateName)
			if err != nil {
				return nil, err
			}

			err = copyDir(sub, ".", templatesPath)
			if err != nil {
				return nil, err
			}

		}
	}
	return nil, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func copyFile(srcFile fs.File, dstFile string) error {
	if fileExists(dstFile) {
		fmt.Println("File already exists: ", dstFile, " skipping")
		return nil
	}

	fmt.Println("Copying file: ", dstFile)

	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			panic(err)
		}
	}(out)

	_, err = io.Copy(out, srcFile)
	return err
}

func copyDir(fsys fs.FS, srcDir string, dstDir string) error {
	err := fs.WalkDir(fsys, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, path)
		if d.IsDir() {
			return os.MkdirAll(dstPath, os.ModePerm)
		}

		if d.Type().IsRegular() {
			srcFile, err := fsys.Open(path)
			if err != nil {
				return err
			}
			defer func(srcFile fs.File) {
				err := srcFile.Close()
				if err != nil {
					panic(err)
				}
			}(srcFile)

			return copyFile(srcFile, dstPath)
		}

		return nil
	})

	return err
}
