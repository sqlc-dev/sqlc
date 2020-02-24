// +build exp

package compiler

import (
	"fmt"
	"io"
	"os"
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
	for _, schema := range c.Schemas {
		for _, table := range schema.Tables {
			s := dinosql.GoStruct{
				Table: pg.FQN{Schema: table.Rel.Schema, Rel: table.Rel.Name},
				Name:  strings.Title(table.Rel.Name),
			}
			for _, col := range table.Columns {
				s.Fields = append(s.Fields, dinosql.GoField{
					Name: strings.Title(col.Name),
					Type: "string",
					Tags: map[string]string{"json:": col.Name},
				})
			}
			structs = append(structs, s)
		}
	}
	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}

	return &Result{structs: structs}, nil
}
