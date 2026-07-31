package core

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// TypeNameString is the catalog's name for the type an AST node names: the type
// as it was written, lowercased, with "[]" appended for an array. Engines
// report a type either as a plain name or as a list of qualifying parts.
func TypeNameString(tn *ast.TypeName) string {
	if tn == nil {
		return ""
	}
	name := strings.TrimSpace(tn.Name)
	if name == "" && tn.Names != nil {
		parts := make([]string, 0, len(tn.Names.Items))
		for _, item := range tn.Names.Items {
			s, ok := item.(*ast.String)
			if !ok || s.Str == "pg_catalog" {
				continue
			}
			parts = append(parts, s.Str)
		}
		name = strings.Join(parts, ".")
	}
	name = strings.ToLower(name)
	if name != "" && tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0 {
		name += ArraySuffix
	}
	return name
}

// ResolveType returns the type an AST node names, registering it when the
// dialect's seed did not: a schema is free to declare types of its own.
func (c *Catalog) ResolveType(tn *ast.TypeName) (int64, error) {
	name := TypeNameString(tn)
	if name == "" {
		return 0, fmt.Errorf("missing type name")
	}
	return c.ResolveTypeName(name)
}

// ResolveTypeName is ResolveType for a type already reduced to its name.
func (c *Catalog) ResolveTypeName(name string) (int64, error) {
	if oid, err := c.TypeOID(name); err == nil {
		return oid, nil
	}
	if element, ok := strings.CutSuffix(name, ArraySuffix); ok {
		elementOID, err := c.ResolveTypeName(element)
		if err != nil {
			return 0, err
		}
		return c.CreateArrayType(name, elementOID)
	}
	return c.CreateUserType(name, "U")
}
