package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"
)

type Field struct {
	Name string
	Type string

	IsNode bool
	IsList bool
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

func (f Field) Convert() string {
	switch {
	case f.Name == "Base":
		return "convertCreateStmt(&n.Base)"
	case f.Name == "ValuesLists" && strings.HasPrefix(f.Type, "[][]"):
		return fmt.Sprintf("convertValuesList(n.%s)", f.Name)
	case f.Type == "interface{}":
		return "&ast.TODO{}"
	case isBuiltIn(f.Type):
		return "n." + f.Name
	case strings.HasSuffix(f.Type, "ast.List"):
		return fmt.Sprintf("convertList(n.%s)", f.Name)
	case strings.HasSuffix(f.Type, "ast.Node"):
		return fmt.Sprintf("convertNode(n.%s)", f.Name)
	case strings.HasPrefix(f.Type, "*"):
		return fmt.Sprintf("convert%s(n.%s)", strings.TrimPrefix(f.Type, "*"), f.Name)
	}
	return fmt.Sprintf("pg.%s(n.%s)", strings.TrimPrefix(f.Type, "*"), f.Name)
}

type Struct struct {
	Alias       string
	Name        string
	HasLocation bool
	Fields      []Field
}

type File struct {
	ImportsAST bool
	Structs    []Struct
}

const structTmpl = `package pg
{{- if .ImportsAST}}

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
){{end}}
{{range .Structs}}
{{if .Alias}}
type {{.Name}} {{.Alias}}
{{else}}
type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{.Type}}
{{- end}}
}
{{end}}

func (n *{{.Name}}) Pos() int {
	{{if .HasLocation}}return n.Location{{else}}return 0{{end}}
}
{{end}}
`

const convertTmpl = `package postgresql

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func convertList(l nodes.List) *ast.List {
	out := &ast.List{}
	for _, item := range l.Items {
		out.Items = append(out.Items, convertNode(item))
	}
	return out
}

func convertValuesList(l [][]nodes.Node) [][]ast.Node {
	out := [][]ast.Node{}
	for _, outer := range l {
		o := []ast.Node{}
		for _, inner := range outer {
			o = append(o, convertNode(inner))
		}
		out = append(out, o)
	}
	return out
}

func convert(node nodes.Node) (ast.Node, error) {
	return convertNode(node), nil
}

{{range .}}
func convert{{.Name}}(n *nodes.{{.Name}}) *pg.{{.Name}} {
	if n == nil {
		return nil
	}
	return &pg.{{.Name}}{
		{{- range .Fields}}
		  {{.Name}}: {{.Convert}},
		{{- end}}
	}
}
{{end}}


func convertNode(node nodes.Node) ast.Node {
	switch n := node.(type) {
	{{range .}}
	case nodes.{{.Name}}:
		return convert{{.Name}}(&n)
	{{end}}
	default:
		return &ast.TODO{}
	}
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

	case *ast.InterfaceType:
		return "interface{}"

	default:
		fmt.Printf("%#v\n", n)
		return "string"
	}
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, "testdata/pg_query_go/nodes", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmpl := template.Must(template.New("").Parse(structTmpl))
	cTmpl := template.Must(template.New("").Parse(convertTmpl))
	cOut := []Struct{}

	for _, pkg := range f {
		for name, file := range pkg.Files {
			if filepath.Base(name) == "node.go" {
				continue
			}
			if filepath.Base(name) == "list.go" {
				continue
			}
			ctx := File{}
			for _, decl := range file.Decls {
				if gen, ok := decl.(*ast.GenDecl); ok {
					if gen.Tok != token.TYPE {
						continue
					}
					if len(gen.Specs) != 1 {
						panic("expected one spec")
					}
					spec := gen.Specs[0].(*ast.TypeSpec)
					// if spec.Name.String() != "Oid" {
					// 	continue
					// }
					switch typ := spec.Type.(type) {

					case *ast.Ident:
						ctx.Structs = append(ctx.Structs, Struct{
							Name:  spec.Name.String(),
							Alias: typ.String(),
						})

					case *ast.StructType:
						out := Struct{
							Name: spec.Name.String(),
						}
						for _, field := range typ.Fields.List {
							fieldType := typName(field.Type)

							if strings.HasSuffix(fieldType, "ast.Node") {
								ctx.ImportsAST = true
							}
							if strings.HasSuffix(fieldType, "ast.List") {
								ctx.ImportsAST = true
							}
							if field.Names[0].String() == "Location" && fieldType == "int" {
								out.HasLocation = true
							}
							out.Fields = append(out.Fields, Field{
								Name: field.Names[0].String(),
								Type: fieldType,
							})
						}

						ctx.Structs = append(ctx.Structs, out)
						if unicode.IsUpper(rune(out.Name[0])) {
							cOut = append(cOut, out)
						}

					default:
						log.Printf("%T\n", typ)
						log.Printf("%#v\n", typ)
					}
				}
			}
			if len(ctx.Structs) == 0 {
				continue
			}

			f, err := os.Create(filepath.Join("../../internal/sql/ast/pg", filepath.Base(name)))
			if err != nil {
				log.Fatal(err)
			}
			if tmpl.Execute(f, ctx); err != nil {
				log.Fatal(err)
			}
			f.Close()
		}
	}

	// TODO: Sort cOut
	sort.Slice(cOut, func(i, j int) bool { return cOut[i].Name < cOut[j].Name })

	cf, err := os.Create("../../internal/postgresql/convert.go")
	if err != nil {
		log.Fatal(err)
	}
	if cTmpl.Execute(cf, cOut); err != nil {
		log.Fatal(err)
	}
	cf.Close()
}
