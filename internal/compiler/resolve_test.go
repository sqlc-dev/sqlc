package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/named"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func TestParamTypeString(t *testing.T) {
	t.Parallel()

	t.Run("postgresql type aliases", func(t *testing.T) {
		t.Parallel()
		comp := &Compiler{parser: postgresql.NewParser()}

		got := comp.paramTypeString(&Column{DataType: "pg_catalog.int4", ArrayDims: 2})
		if got != "integer[][]" {
			t.Fatalf("expected integer[][], got %q", got)
		}
	})

	t.Run("structured type metadata is preferred", func(t *testing.T) {
		t.Parallel()
		comp := &Compiler{parser: postgresql.NewParser()}

		got := comp.paramTypeString(&Column{
			DataType: "catalog.pg_catalog.int4",
			Type:     &ast.TypeName{Schema: "pg_catalog", Name: "bpchar"},
		})
		if got != "character" {
			t.Fatalf("expected character, got %q", got)
		}
	})

	t.Run("sqlite keeps names unchanged", func(t *testing.T) {
		t.Parallel()
		comp := &Compiler{parser: sqlite.NewParser()}

		got := comp.paramTypeString(&Column{DataType: "custom_type", ArrayDims: 1})
		if got != "custom_type[]" {
			t.Fatalf("expected custom_type[], got %q", got)
		}
	})
}

func TestIncompatibleParamRefErrorFormatsTypeNames(t *testing.T) {
	t.Parallel()

	comp := &Compiler{parser: postgresql.NewParser()}
	err := comp.incompatibleParamRefError(paramRef{ref: &ast.ParamRef{Number: 1}}, Parameter{
		Number: 1,
		Column: &Column{DataType: "text"},
	}, Parameter{
		Number: 1,
		Column: &Column{DataType: "pg_catalog.int4"},
	})

	sqlErr, ok := err.(*sqlerr.Error)
	if !ok {
		t.Fatalf("expected *sqlerr.Error, got %T", err)
	}
	if sqlErr.Message != "parameter $1 has incompatible types: text, integer" {
		t.Fatalf("unexpected message: %q", sqlErr.Message)
	}
}

func TestMergeResolvedParamKeepsFirstNameForCompatibleTypes(t *testing.T) {
	t.Parallel()

	merged := mergeResolvedParam(
		Parameter{Number: 1, Column: &Column{Name: "user", DataType: "text"}},
		Parameter{Number: 1, Column: &Column{Name: "student_user", DataType: "text"}},
	)

	if merged.Column == nil {
		t.Fatal("expected merged column")
	}
	if merged.Column.Name != "user" {
		t.Fatalf("expected first inferred name to win, got %q", merged.Column.Name)
	}
}

func TestResolvedFuncCallArgType(t *testing.T) {
	t.Parallel()

	fun := &catalog.Function{Args: []*catalog.Argument{
		{Name: "lhs", Type: &ast.TypeName{Name: "int8"}},
		{Name: "rhs", Type: &ast.TypeName{Name: "text"}},
	}}

	if got := resolvedFuncCallArgType(fun, 0, ""); got == nil || got.Name != "int8" {
		t.Fatalf("expected positional arg type int8, got %#v", got)
	}
	if got := resolvedFuncCallArgType(fun, 0, "rhs"); got == nil || got.Name != "text" {
		t.Fatalf("expected named arg type text, got %#v", got)
	}
	if got := resolvedFuncCallArgType(fun, 2, ""); got != nil {
		t.Fatalf("expected nil for out-of-range positional arg, got %#v", got)
	}
}

func TestResolveCatalogRefsInsertTargetStringInfersColumnName(t *testing.T) {
	t.Parallel()

	comp := &Compiler{parser: postgresql.NewParser(), catalog: postgresql.NewCatalog()}

	var schema *catalog.Schema
	for _, s := range comp.catalog.Schemas {
		if s.Name == comp.catalog.DefaultSchema {
			schema = s
			break
		}
	}
	if schema == nil {
		t.Fatal("default schema not found")
	}

	tableName := "solar_commcard_mapping"
	schema.Tables = append(schema.Tables, &catalog.Table{
		Rel: &ast.TableName{Schema: schema.Name, Name: tableName},
		Columns: []*catalog.Column{&catalog.Column{
			Name:      "deviceId",
			Type:      ast.TypeName{Schema: "pg_catalog", Name: "int8"},
			IsNotNull: true,
		}},
	})

	rv := &ast.RangeVar{Relname: &tableName}
	params, err := comp.resolveCatalogRefs(nil, []*ast.RangeVar{rv}, []paramRef{{
		parent: &ast.String{Str: "deviceId"},
		rv:     rv,
		ref:    &ast.ParamRef{Number: 1},
	}}, named.NewParamSet(nil, true), nil)
	if err != nil {
		t.Fatalf("resolveCatalogRefs returned error: %v", err)
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	if params[0].Column == nil {
		t.Fatal("expected resolved column metadata")
	}
	if params[0].Column.Name != "deviceId" {
		t.Fatalf("expected inferred name deviceId, got %q", params[0].Column.Name)
	}
	if params[0].Column.OriginalName != "deviceId" {
		t.Fatalf("expected original name deviceId, got %q", params[0].Column.OriginalName)
	}
	if params[0].Column.DataType != "pg_catalog.int8" {
		t.Fatalf("expected data type pg_catalog.int8, got %q", params[0].Column.DataType)
	}
	if params[0].Column.Table == nil || params[0].Column.Table.Name != tableName {
		t.Fatalf("expected table %q, got %#v", tableName, params[0].Column.Table)
	}
}

func TestResolveCatalogRefsAExprUsesScopedRangeVar(t *testing.T) {
	t.Parallel()

	comp := &Compiler{parser: postgresql.NewParser(), catalog: postgresql.NewCatalog()}

	var schema *catalog.Schema
	for _, s := range comp.catalog.Schemas {
		if s.Name == comp.catalog.DefaultSchema {
			schema = s
			break
		}
	}
	if schema == nil {
		t.Fatal("default schema not found")
	}

	tableName := "solar_commcard_mapping"
	schema.Tables = append(schema.Tables, &catalog.Table{
		Rel: &ast.TableName{Schema: schema.Name, Name: tableName},
		Columns: []*catalog.Column{{
			Name:      "deviceId",
			Type:      ast.TypeName{Schema: "pg_catalog", Name: "int8"},
			IsNotNull: true,
		}},
	})

	rv := &ast.RangeVar{Relname: &tableName}
	params, err := comp.resolveCatalogRefs(nil, []*ast.RangeVar{rv, rv}, []paramRef{{
		parent: &ast.A_Expr{
			Lexpr: &ast.ColumnRef{Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "deviceId"}}}},
			Rexpr: &ast.ParamRef{Number: 1},
		},
		rv:  rv,
		ref: &ast.ParamRef{Number: 1},
	}}, named.NewParamSet(nil, true), nil)
	if err != nil {
		t.Fatalf("resolveCatalogRefs returned error: %v", err)
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	if params[0].Column == nil {
		t.Fatal("expected resolved column metadata")
	}
	if params[0].Column.Name != "deviceId" {
		t.Fatalf("expected inferred name deviceId, got %q", params[0].Column.Name)
	}
}
