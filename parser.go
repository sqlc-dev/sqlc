package strongdb

import (
	"bufio"
	"bytes"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/strongdb/postgres"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func ParseSchmea(dir string) (*postgres.Schema, error) {
	// Keep the import around
	if false {
		spew.Dump(dir)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := postgres.Schema{}
	for _, f := range files {
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		tree, err := pg.Parse(string(blob))
		if err != nil {
			return nil, err
		}
		parse(&s, tree)
	}
	return &s, nil
}

func parse(s *postgres.Schema, tree pg.ParsetreeList) {
	for _, stmt := range tree.Statements {
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch n := raw.Stmt.(type) {
		case nodes.CreateStmt:
			table := postgres.Table{
				Name:   *n.Relation.Relname,
				GoName: structName(*n.Relation.Relname),
			}
			for _, elt := range n.TableElts.Items {
				switch n := elt.(type) {
				case nodes.ColumnDef:
					// spew.Dump(n)
					// log.Printf("not null: %t", n.IsNotNull)
					table.Columns = append(table.Columns, postgres.Column{
						Name:    *n.Colname,
						GoName:  structName(*n.Colname),
						NotNull: isNotNull(n),
					})
				}
			}
			s.Tables = append(s.Tables, table)
		default:
			// spew.Dump(n)
		}
	}
}

func isNotNull(n nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints.Items {
		switch n := c.(type) {
		case nodes.Constraint:
			if n.Contype == nodes.CONSTR_NOTNULL {
				return true
			}
			if n.Contype == nodes.CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func ParseQueries(dir string) (*postgres.Schema, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := postgres.Schema{}
	for _, f := range files {
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		tree, err := pg.Parse(string(blob))
		if err != nil {
			return nil, err
		}
		parseFuncs(&s, tree)
	}
	return &s, nil
}

func parseFuncs(s *postgres.Schema, tree pg.ParsetreeList) {
	for _, stmt := range tree.Statements {
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch n := raw.Stmt.(type) {
		case nodes.SelectStmt:
			t := tableName(n)
			spew.Dump(t)
			spew.Dump(n)
			// log.Printf("%T\n", n)
		default:
			log.Printf("%T\n", n)
		}
	}
}

func tableName(n nodes.SelectStmt) string {
	for _, item := range n.FromClause.Items {
		switch i := item.(type) {
		case nodes.RangeVar:
			return *i.Relname
		}
	}
	return ""
}

var hh = `package equinox
{{range .Tables}}
type {{.GoName}} struct { {{- range .Columns}}
  {{.GoName}} {{if .NotNull }}string{{else}}sql.NullString{{end}}
  {{- end}}
}
{{end}}
`

func structName(name string) string {
	if strings.HasSuffix(name, "s") {
		name = name[:len(name)-1]
	}
	out := ""
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func generate(s *postgres.Schema) string {
	funcMap := template.FuncMap{}

	fileTmpl := template.Must(template.New("table").Funcs(funcMap).Parse(hh))
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	fileTmpl.Execute(w, s)
	w.Flush()
	code, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	return string(code)
}
