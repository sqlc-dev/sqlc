package compiler

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/kyleconroy/sqlc/internal/codegen/golang"
	"github.com/kyleconroy/sqlc/internal/config"
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
func parseCatalog(p Parser, c *catalog.Catalog, schemas []string) error {
	files, err := sqlpath.Glob(schemas)
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

func parseQueries(p Parser, c *catalog.Catalog, queries []string) (*Result, error) {
	var q []*Query
	merr := multierr.New()
	set := map[string]struct{}{}
	files, err := sqlpath.Glob(queries)
	if err != nil {
		return nil, err
	}
	for _, filename := range files {
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		src := string(blob)
		stmts, err := p.Parse(strings.NewReader(src))
		if err != nil {
			merr.Add(filename, src, 0, err)
			continue
		}
		for _, stmt := range stmts {
			query, err := parseQuery(p, c, stmt.Raw, src, false)
			if err == ErrUnsupportedStatementType {
				continue
			}
			if err != nil {
				merr.Add(filename, src, stmt.Raw.Pos(), err)
				continue
			}
			if query.Name != "" {
				if _, exists := set[query.Name]; exists {
					merr.Add(filename, src, 0, fmt.Errorf("duplicate query name: %s", query.Name))
					continue
				}
				set[query.Name] = struct{}{}
			}
			query.Filename = filepath.Base(filename)
			if query != nil {
				q = append(q, query)
			}
		}
	}
	if len(merr.Errs()) > 0 {
		return nil, merr
	}
	if len(q) == 0 {
		return nil, fmt.Errorf("no queries contained in paths %s", strings.Join(queries, ","))
	}
	return &Result{
		Catalog: c,
		Queries: q,
	}, nil
}

// Deprecated.
func buildResult(c *catalog.Catalog) (*BuildResult, error) {
	var structs []golang.Struct
	var enums []golang.Enum
	for _, schema := range c.Schemas {
		for _, table := range schema.Tables {
			s := golang.Struct{
				Table:   pg.FQN{Schema: schema.Name, Rel: table.Rel.Name},
				Name:    strings.Title(table.Rel.Name),
				Comment: table.Comment,
			}
			for _, col := range table.Columns {
				s.Fields = append(s.Fields, golang.Field{
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
				e := golang.Enum{
					Name:    structName(name),
					Comment: t.Comment,
				}
				for _, v := range t.Vals {
					e.Constants = append(e.Constants, golang.Constant{
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
	return &BuildResult{structs: structs, enums: enums}, nil
}

func Run(conf config.SQL, combo config.CombinedSettings) (*BuildResult, error) {
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
