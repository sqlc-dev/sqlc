// +build exp

package compiler

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/dolphin"
	"github.com/kyleconroy/sqlc/internal/migrations"
	"github.com/kyleconroy/sqlc/internal/multierr"
	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlpath"
	"github.com/kyleconroy/sqlc/internal/sqlite"
)

type Parser interface {
	Parse(io.Reader) ([]ast.Statement, error)
}

// copied over from gen.go
func structName(name string) string {
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

var identPattern = regexp.MustCompile("[^a-zA-Z0-9_]+")

func enumValueName(value string) string {
	name := ""
	id := strings.Replace(value, "-", "_", -1)
	id = strings.Replace(id, ":", "_", -1)
	id = strings.Replace(id, "/", "_", -1)
	id = identPattern.ReplaceAllString(id, "")
	for _, part := range strings.Split(id, "_") {
		name += strings.Title(part)
	}
	return name
}

// end copypasta
func parseCatalog(p Parser, c *catalog.Catalog, schema []string) error {
	files, err := sqlpath.Glob(schema)
	if err != nil {
		return err
	}
	merr := multierr.New()
	for _, filename := range files {
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		contents := migrations.RemoveRollbackStatements(string(blob))
		stmts, err := p.Parse(strings.NewReader(contents))
		if err != nil {
			merr.Add(filename, contents, 0, err)
			continue
		}
		for i := range stmts {
			if err := c.Update(stmts[i]); err != nil {
				merr.Add(filename, contents, stmts[i].Pos(), err)
				continue
			}
		}
	}
	if len(merr.Errs()) > 0 {
		return merr
	}
	return nil
}

func buildResult(c *catalog.Catalog) (*Result, error) {
	var structs []dinosql.GoStruct
	var enums []dinosql.GoEnum
	for _, schema := range c.Schemas {
		for _, table := range schema.Tables {
			s := dinosql.GoStruct{
				Table:   pg.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    strings.Title(table.Rel.Name),
				Comment: table.Comment,
			}
			for _, col := range table.Columns {
				s.Fields = append(s.Fields, dinosql.GoField{
					Name:    structName(col.Name),
					Type:    "string",
					Tags:    map[string]string{"json:": col.Name},
					Comment: col.Comment,
				})
			}
			structs = append(structs, s)
		}
		for _, typ := range schema.Types {
			switch t := typ.(type) {
			case *catalog.Enum:
				var name string
				if schema.Name == c.DefaultSchema {
					name = t.Name
				} else {
					name = schema.Name + "_" + t.Name
				}
				e := dinosql.GoEnum{
					Name:    structName(name),
					Comment: t.Comment,
				}
				for _, v := range t.Vals {
					e.Constants = append(e.Constants, dinosql.GoConstant{
						Name:  e.Name + enumValueName(v),
						Value: v,
						Type:  e.Name,
					})
				}
				enums = append(enums, e)
			}
		}
	}

	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return &Result{structs: structs, enums: enums}, nil
}

func Run(conf config.SQL, combo config.CombinedSettings) (*Result, error) {
	var c *catalog.Catalog
	var p Parser

	switch conf.Engine {
	case config.EngineXLemon:
		p = sqlite.NewParser()
		c = catalog.New("main")
	case config.EngineXDolphin:
		p = dolphin.NewParser()
		c = catalog.New("public") // TODO: What is the default database for MySQL?
	case config.EngineXElephant:
		p = postgresql.NewParser()
		c = postgresql.NewCatalog()
	default:
		return nil, fmt.Errorf("unknown engine: %s", conf.Engine)
	}

	if err := parseCatalog(p, c, conf.Schema); err != nil {
		return nil, err
	}

	return buildResult(c)
}
