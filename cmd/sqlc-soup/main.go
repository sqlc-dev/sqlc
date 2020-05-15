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

const soupTmpl = `package astutils

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

type Visitor interface {
	Visit(ast.Node) Visitor
}

type VisitorFunc func(ast.Node)

func (vf VisitorFunc) Visit(node ast.Node) Visitor {
	vf(node)
	return vf
}

func walkn(f Visitor, node ast.Node) {
	if node != nil {
		Walk(f, node)
	}
}

func Walk(f Visitor, node ast.Node) {
	if f = f.Visit(node); f == nil {
		return
	}
	switch n := node.(type) {
	{{range .}}
	case *{{.Pkg}}.{{.Name}}:
		{{- range .Nodes}}
		walkn(f, n.{{.Name}})
		{{- else}}
		// pass
		{{- end}}
	{{end}}
	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))
	}

	f.Visit(nil)
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

	cf, err := os.Create("../../internal/sql/astutils/walk.go")
	if err != nil {
		log.Fatal(err)
	}
	if tmpl.Execute(cf, ctx.Structs); err != nil {
		log.Fatal(err)
	}
	cf.Close()

}
