package compiler

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
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
