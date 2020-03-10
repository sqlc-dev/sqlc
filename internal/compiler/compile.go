// +build exp

package compiler

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/dinosql"
	"github.com/kyleconroy/sqlc/internal/dolphin"
	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgresql"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
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

func Run(conf config.SQL, combo config.CombinedSettings) (*Result, error) {
	var p Parser

	switch conf.Engine {
	case config.EngineXLemon:
		p = sqlite.NewParser()
	case config.EngineXDolphin:
		p = dolphin.NewParser()
	case config.EngineXElephant:
		p = postgresql.NewParser()
	default:
		return nil, fmt.Errorf("unknown engine: %s", conf.Engine)
	}

	rd, err := os.Open(conf.Schema)
	if err != nil {
		return nil, err
	}

	stmts, err := p.Parse(rd)
	if err != nil {
		return nil, err
	}

	c, err := catalog.Build(stmts)
	if err != nil {
		return nil, err
	}

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
				// TODO: This name should be public, not main
				if schema.Name == "main" {
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
