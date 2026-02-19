package validate

import (
	"strings"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func makeTestCatalog(t *testing.T) (*catalog.Catalog, *ast.TableName) {
	t.Helper()

	p := postgresql.NewParser()
	stmts, err := p.Parse(strings.NewReader(`
		CREATE TABLE cart_items (
			owner_id       VARCHAR(255) NOT NULL,
			product_id     UUID         NOT NULL,
			price_amount   DECIMAL      NOT NULL,
			price_currency VARCHAR(3)   NOT NULL,
			PRIMARY KEY (owner_id, product_id)
		);
	`))
	if err != nil {
		t.Fatalf("parse schema: %v", err)
	}

	cat := catalog.New("public")
	for _, stmt := range stmts {
		if err := cat.Update(stmt, nil); err != nil {
			t.Fatalf("update catalog: %v", err)
		}
	}

	tableName := &ast.TableName{Schema: "public", Name: "cart_items"}
	return cat, tableName
}

func makeStmt(action ast.OnConflictAction, setItems []struct{ col, val string }) *ast.InsertStmt {
	stmt := &ast.InsertStmt{
		Relation: &ast.RangeVar{
			Schemaname: strPtr("public"),
			Relname:    strPtr("cart_items"),
		},
	}

	if action == ast.OnConflictActionNone {
		return stmt
	}

	items := make([]ast.Node, 0, len(setItems))
	for _, si := range setItems {
		colName := si.col
		items = append(items, &ast.ResTarget{
			Name: &colName,
			Val: &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{
						&ast.String{Str: "excluded"},
						&ast.String{Str: si.val},
					},
				},
			},
		})
	}

	stmt.OnConflictClause = &ast.OnConflictClause{
		Action:     action,
		TargetList: &ast.List{Items: items},
	}
	return stmt
}

func strPtr(s string) *string { return &s }

func TestOnConflictClause(t *testing.T) {
	cat, tableName := makeTestCatalog(t)

	tests := []struct {
		name    string
		stmt    *ast.InsertStmt
		wantErr bool
	}{
		{
			name: "valid columns in SET and EXCLUDED",
			stmt: makeStmt(ast.OnConflictActionUpdate, []struct{ col, val string }{
				{"price_amount", "price_amount"},
				{"price_currency", "price_currency"},
			}),
			wantErr: false,
		},
		{
			name: "invalid column on left side of SET",
			stmt: makeStmt(ast.OnConflictActionUpdate, []struct{ col, val string }{
				{"price_amount1", "price_amount"},
			}),
			wantErr: true,
		},
		{
			name: "invalid EXCLUDED reference on right side",
			stmt: makeStmt(ast.OnConflictActionUpdate, []struct{ col, val string }{
				{"price_amount", "price_amount1"},
			}),
			wantErr: true,
		},
		{
			name:    "DO NOTHING skips column validation",
			stmt:    makeStmt(ast.OnConflictActionNothing, nil),
			wantErr: false,
		},
		{
			name:    "no OnConflictClause passes without error",
			stmt:    makeStmt(ast.OnConflictActionNone, nil),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OnConflictClause(cat, tt.stmt, tableName)
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
