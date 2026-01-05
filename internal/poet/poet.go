// Package poet provides helpers for generating Go source code using the go/ast package.
// It offers a fluent API for building Go AST nodes that can be formatted into source code.
package poet

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strconv"
	"strings"
)

// File represents a Go source file being built.
type File struct {
	name       string
	pkg        string
	buildTags  string
	comments   []string // File-level comments (before package)
	imports    []ImportSpec
	decls      []ast.Decl
	fset       *token.FileSet
	nextPos    token.Pos
	commentMap ast.CommentMap
}

// ImportSpec represents an import declaration.
type ImportSpec struct {
	Name string // Optional alias (empty for default)
	Path string // Import path
}

// NewFile creates a new file builder with the given package name.
func NewFile(pkg string) *File {
	return &File{
		pkg:        pkg,
		fset:       token.NewFileSet(),
		nextPos:    1,
		commentMap: make(ast.CommentMap),
	}
}

// SetBuildTags sets the build tags for the file.
func (f *File) SetBuildTags(tags string) *File {
	f.buildTags = tags
	return f
}

// AddComment adds a file-level comment (appears before package declaration).
func (f *File) AddComment(comment string) *File {
	f.comments = append(f.comments, comment)
	return f
}

// AddImport adds an import to the file.
func (f *File) AddImport(path string) *File {
	f.imports = append(f.imports, ImportSpec{Path: path})
	return f
}

// AddImportWithAlias adds an import with an alias to the file.
func (f *File) AddImportWithAlias(alias, path string) *File {
	f.imports = append(f.imports, ImportSpec{Name: alias, Path: path})
	return f
}

// AddImports adds multiple imports to the file, organized by groups.
func (f *File) AddImports(groups [][]ImportSpec) *File {
	for _, group := range groups {
		f.imports = append(f.imports, group...)
	}
	return f
}

// AddDecl adds a declaration to the file.
func (f *File) AddDecl(decl ast.Decl) *File {
	f.decls = append(f.decls, decl)
	return f
}

// allocPos allocates a new position for AST nodes.
func (f *File) allocPos() token.Pos {
	pos := f.nextPos
	f.nextPos++
	return pos
}

// Render generates the Go source code for the file.
func (f *File) Render() ([]byte, error) {
	var buf bytes.Buffer

	// Build tags
	if f.buildTags != "" {
		buf.WriteString("//go:build ")
		buf.WriteString(f.buildTags)
		buf.WriteString("\n\n")
	}

	// File-level comments
	for _, comment := range f.comments {
		buf.WriteString(comment)
		buf.WriteString("\n")
	}

	// Package declaration
	buf.WriteString("package ")
	buf.WriteString(f.pkg)
	buf.WriteString("\n")

	// Imports
	if len(f.imports) > 0 {
		buf.WriteString("\nimport (\n")
		prevWasStd := true
		for i, imp := range f.imports {
			// Add blank line between std and external packages
			isStd := !strings.Contains(imp.Path, ".")
			if i > 0 && prevWasStd && !isStd {
				buf.WriteString("\n")
			}
			prevWasStd = isStd

			buf.WriteString("\t")
			if imp.Name != "" {
				buf.WriteString(imp.Name)
				buf.WriteString(" ")
			}
			buf.WriteString(strconv.Quote(imp.Path))
			buf.WriteString("\n")
		}
		buf.WriteString(")\n")
	}

	// Declarations
	for _, decl := range f.decls {
		buf.WriteString("\n")
		declBuf, err := f.renderDecl(decl)
		if err != nil {
			return nil, err
		}
		buf.Write(declBuf)
		buf.WriteString("\n")
	}

	// Format the generated code
	return format.Source(buf.Bytes())
}

func (f *File) renderDecl(decl ast.Decl) ([]byte, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	// Create a minimal file to format the declaration
	file := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{decl},
	}

	if err := format.Node(&buf, fset, file); err != nil {
		return nil, err
	}

	// Extract just the declaration part (skip "package main\n")
	result := buf.Bytes()
	idx := bytes.Index(result, []byte("\n"))
	if idx >= 0 {
		result = result[idx+1:]
	}
	return result, nil
}
