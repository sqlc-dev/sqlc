package compiler

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
)

// This is mainly copy-pasted from internal/postgresql/parse.go
func stringSlice(list *ast.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
			continue
		}
		if n, ok := item.(*ast.String); ok {
			items = append(items, n.Str)
			continue
		}
	}
	return items
}

type Relation struct {
	Catalog string
	Schema  string
	Name    string
}

func parseRelation(node ast.Node) (*Relation, error) {
	switch n := node.(type) {

	case *ast.List:
		parts := stringSlice(n)
		switch len(parts) {
		case 1:
			return &Relation{
				Name: parts[0],
			}, nil
		case 2:
			return &Relation{
				Schema: parts[0],
				Name:   parts[1],
			}, nil
		case 3:
			return &Relation{
				Catalog: parts[0],
				Schema:  parts[1],
				Name:    parts[2],
			}, nil
		default:
			return nil, fmt.Errorf("invalid name: %s", astutils.Join(n, "."))
		}

	case *ast.RangeVar:
		name := Relation{}
		if n.Catalogname != nil {
			name.Catalog = *n.Catalogname
		}
		if n.Schemaname != nil {
			name.Schema = *n.Schemaname
		}
		if n.Relname != nil {
			name.Name = *n.Relname
		}
		return &name, nil

	case *ast.TypeName:
		return parseRelation(n.Names)

	default:
		return nil, fmt.Errorf("unexpected node type: %T", n)
	}
}

func ParseTableName(node ast.Node) (*ast.TableName, error) {
	rel, err := parseRelation(node)
	if err != nil {
		return nil, fmt.Errorf("parse table name: %w", err)
	}
	return &ast.TableName{
		Catalog: rel.Catalog,
		Schema:  rel.Schema,
		Name:    rel.Name,
	}, nil
}

func ParseTypeName(node ast.Node) (*ast.TypeName, error) {
	rel, err := parseRelation(node)
	if err != nil {
		return nil, fmt.Errorf("parse table name: %w", err)
	}
	return &ast.TypeName{
		Catalog: rel.Catalog,
		Schema:  rel.Schema,
		Name:    rel.Name,
	}, nil
}

func ParseRelationString(name string) (*Relation, error) {
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 1:
		return &Relation{
			Name: parts[0],
		}, nil
	case 2:
		return &Relation{
			Schema: parts[0],
			Name:   parts[1],
		}, nil
	case 3:
		return &Relation{
			Catalog: parts[0],
			Schema:  parts[1],
			Name:    parts[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid name: %s", name)
	}
}
