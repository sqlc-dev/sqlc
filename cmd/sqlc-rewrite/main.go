package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string
}

func isBuiltIn(typ string) bool {
	typ = strings.TrimPrefix(typ, "[]")
	typ = strings.TrimPrefix(typ, "*")

	switch typ {
	case "bool":
		return true
	case "int":
		return true
	case "int32":
		return true
	case "uint32":
		return true
	case "string":
		return true
	case "byte":
		return true
	case "int64":
		return true
	case "uint16":
		return true
	case "int16":
		return true
	case "float64":
		return true
	}
	return false
}

type Struct struct {
	Pkg         string
	Name        string
	HasLocation bool
	Fields      []Field
}

func (s *Struct) Nodes() []Field {
	nodes := []Field{}
	for _, field := range s.Fields {
		if isBuiltIn(field.Type) {
			continue
		}
		if field.Type == "ast.Node" {
			nodes = append(nodes, field)
			continue
		}
		if strings.HasPrefix(field.Type, "*") {
			nodes = append(nodes, field)
			continue
		}

		// log.Println(s.Name, field.Name, field.Type)
	}
	return nodes
}

type File struct {
	ImportsAST bool
	Structs    []Struct
}

const soupTmpl = `// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package astutils

import (
	"fmt"
	"reflect"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

// An ApplyFunc is invoked by Apply for each node n, even if n is nil,
// before and/or after the node's children, using a Cursor describing
// the current node and providing operations on it.
//
// The return value of ApplyFunc controls the syntax tree traversal.
// See Apply for details.
type ApplyFunc func(*Cursor) bool

// Apply traverses a syntax tree recursively, starting with root,
// and calling pre and post for each node as described below.
// Apply returns the syntax tree, possibly modified.
//
// If pre is not nil, it is called for each node before the node's
// children are traversed (pre-order). If pre returns false, no
// children are traversed, and post is not called for that node.
//
// If post is not nil, and a prior call of pre didn't return false,
// post is called for each node after its children are traversed
// (post-order). If post returns false, traversal is terminated and
// Apply returns immediately.
//
// Only fields that refer to AST nodes are considered children;
// i.e., token.Pos, Scopes, Objects, and fields of basic types
// (strings, etc.) are ignored.
//
// Children are traversed in the order in which they appear in the
// respective node's struct definition. A package's files are
// traversed in the filenames' alphabetical order.
//
func Apply(root ast.Node, pre, post ApplyFunc) (result ast.Node) {
	parent := &struct{ ast.Node }{root}
	defer func() {
		if r := recover(); r != nil && r != abort {
			panic(r)
		}
		result = parent.Node
	}()
	a := &application{pre: pre, post: post}
	a.apply(parent, "Node", nil, root)
	return
}

var abort = new(int) // singleton, to signal termination of Apply

// A Cursor describes a node encountered during Apply.
// Information about the node and its parent is available
// from the Node, Parent, Name, and Index methods.
//
// If p is a variable of type and value of the current parent node
// c.Parent(), and f is the field identifier with name c.Name(),
// the following invariants hold:
//
//   p.f            == c.Node()  if c.Index() <  0
//   p.f[c.Index()] == c.Node()  if c.Index() >= 0
//
// The methods Replace, Delete, InsertBefore, and InsertAfter
// can be used to change the AST without disrupting Apply.
type Cursor struct {
	parent ast.Node
	name   string
	iter   *iterator // valid if non-nil
	node   ast.Node
}

// Node returns the current Node.
func (c *Cursor) Node() ast.Node { return c.node }

// Parent returns the parent of the current Node.
func (c *Cursor) Parent() ast.Node { return c.parent }

// Name returns the name of the parent Node field that contains the current Node.
// If the parent is a *ast.Package and the current Node is a *ast.File, Name returns
// the filename for the current Node.
func (c *Cursor) Name() string { return c.name }

// Index reports the index >= 0 of the current Node in the slice of Nodes that
// contains it, or a value < 0 if the current Node is not part of a slice.
// The index of the current node changes if InsertBefore is called while
// processing the current node.
func (c *Cursor) Index() int {
	if c.iter != nil {
		return c.iter.index
	}
	return -1
}

// field returns the current node's parent field value.
func (c *Cursor) field() reflect.Value {
	return reflect.Indirect(reflect.ValueOf(c.parent)).FieldByName(c.name)
}

// Replace replaces the current Node with n.
// The replacement node is not walked by Apply.
func (c *Cursor) Replace(n ast.Node) {
	v := c.field()
	if i := c.Index(); i >= 0 {
		v = v.Index(i)
	}
	v.Set(reflect.ValueOf(n))
}

// D// application carries all the shared data so we can pass it around cheaply.
type application struct {
	pre, post ApplyFunc
	cursor    Cursor
	iter      iterator
}

func (a *application) apply(parent ast.Node, name string, iter *iterator, n ast.Node) {
	// convert typed nil into untyped nil
	if v := reflect.ValueOf(n); v.Kind() == reflect.Ptr && v.IsNil() {
		n = nil
	}

	// avoid heap-allocating a new cursor for each apply call; reuse a.cursor instead
	saved := a.cursor
	a.cursor.parent = parent
	a.cursor.name = name
	a.cursor.iter = iter
	a.cursor.node = n

	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}

	// walk children
	// (the order of the cases matches the order of the corresponding node types in go/ast)
	switch n := n.(type) {
	case nil:
		// nothing to do

	{{range .}}
	case *{{.Pkg}}.{{.Name}}:
		{{- range .Nodes}}
		a.apply(n, "{{.Name}}", nil, n.{{.Name}})
		{{- else}}
		// pass
		{{- end}}
	{{end}}

	// Comments and fields
	default:
		panic(fmt.Sprintf("Apply: unexpected node type %T", n))
	}

	if a.post != nil && !a.post(&a.cursor) {
		panic(abort)
	}

	a.cursor = saved
}

// An iterator controls iteration over a slice of nodes.
type iterator struct {
	index, step int
}

func (a *application) applyList(parent ast.Node, name string) {
	// avoid heap-allocating a new iterator for each applyList call; reuse a.iter instead
	saved := a.iter
	a.iter.index = 0
	for {
		// must reload parent.name each time, since cursor modifications might change it
		v := reflect.Indirect(reflect.ValueOf(parent)).FieldByName(name)
		if a.iter.index >= v.Len() {
			break
		}

		// element x may be nil in a bad AST - be cautious
		var x ast.Node
		if e := v.Index(a.iter.index); e.IsValid() {
			x = e.Interface().(ast.Node)
		}

		a.iter.step = 1
		a.apply(parent, name, &a.iter, x)
		a.iter.index += a.iter.step
	}
	a.iter = saved
}
`

func typName(node ast.Node) string {
	switch n := node.(type) {
	case *ast.ArrayType:
		return "[]" + typName(n.Elt)

	case *ast.Ident:
		tn := n.String()
		if tn == "Node" {
			tn = "ast.Node"
		}
		if tn == "List" {
			tn = "*ast.List"
		}
		return tn

	case *ast.StarExpr:
		return "*" + typName(n.X)

	case *ast.SelectorExpr:
		return typName(n.X) + "." + n.Sel.String()

	case *ast.InterfaceType:
		return "interface{}"

	default:
		fmt.Printf("%#v\n", n)
		return "string"
	}
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, "../../internal/sql/ast", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	fset = token.NewFileSet()
	ff, err := parser.ParseDir(fset, "../../internal/sql/ast/pg", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}
	for name, pkg := range ff {
		f[name] = pkg
	}

	tmpl := template.Must(template.New("").Parse(soupTmpl))
	// cTmpl := template.Must(template.New("").Parse(convertTmpl))
	ctx := File{}

	for _, pkg := range f {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if gen, ok := decl.(*ast.GenDecl); ok {
					if gen.Tok != token.TYPE {
						continue
					}
					if len(gen.Specs) != 1 {
						panic("expected one spec")
					}
					spec := gen.Specs[0].(*ast.TypeSpec)

					switch typ := spec.Type.(type) {

					case *ast.StructType:
						if spec.Name.String() == "varatt_external" {
							continue
						}

						out := Struct{
							Name: spec.Name.String(),
							Pkg:  pkg.Name,
						}
						for _, field := range typ.Fields.List {
							fieldType := typName(field.Type)
							out.Fields = append(out.Fields, Field{
								Name: field.Names[0].String(),
								Type: fieldType,
							})
						}

						// sort.Slice(out.Fields, func(i, j int) bool { return out.Fields[i].Name < out.Fields[j].Name })
						ctx.Structs = append(ctx.Structs, out)
					}
				}
			}
		}
	}

	// TODO: Sort cOut
	sort.Slice(ctx.Structs, func(i, j int) bool {
		if ctx.Structs[i].Pkg == ctx.Structs[j].Pkg {
			return ctx.Structs[i].Name < ctx.Structs[j].Name
		} else {
			return ctx.Structs[i].Pkg < ctx.Structs[j].Pkg
		}
	})

	cf, err := os.Create("../../internal/sql/astutils/rewrite.go")
	if err != nil {
		log.Fatal(err)
	}
	if tmpl.Execute(cf, ctx.Structs); err != nil {
		log.Fatal(err)
	}
	cf.Close()

}
