package sqlanalyze

import (
	"testing"

	"github.com/sqlc-dev/sqlc/internal/analysis/scope"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// buildTestCatalog creates a catalog with users and orders tables.
func buildTestCatalog(defaultSchema string) *catalog.Catalog {
	c := catalog.New(defaultSchema)

	intType := ast.TypeName{Name: "integer"}
	textType := ast.TypeName{Name: "text"}
	numericType := ast.TypeName{Name: "numeric"}
	boolType := ast.TypeName{Name: "boolean"}

	c.Update(ast.Statement{
		Raw: &ast.RawStmt{
			Stmt: &ast.CreateTableStmt{
				Name: &ast.TableName{Name: "users"},
				Cols: []*ast.ColumnDef{
					{Colname: "id", TypeName: &intType, IsNotNull: true},
					{Colname: "name", TypeName: &textType, IsNotNull: true},
					{Colname: "email", TypeName: &textType, IsNotNull: false},
					{Colname: "age", TypeName: &intType, IsNotNull: false},
					{Colname: "active", TypeName: &boolType, IsNotNull: true},
				},
			},
		},
	}, nil)

	c.Update(ast.Statement{
		Raw: &ast.RawStmt{
			Stmt: &ast.CreateTableStmt{
				Name: &ast.TableName{Name: "orders"},
				Cols: []*ast.ColumnDef{
					{Colname: "id", TypeName: &intType, IsNotNull: true},
					{Colname: "user_id", TypeName: &intType, IsNotNull: true},
					{Colname: "total", TypeName: &numericType, IsNotNull: true},
					{Colname: "status", TypeName: &textType, IsNotNull: false},
				},
			},
		},
	}, nil)

	return c
}

func TestAnalyzeSelectSimple(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT name, email FROM users WHERE id = $1
	nameStr := "name"
	emailStr := "email"
	relname := "users"
	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "name"}}},
						},
						Name: &nameStr,
					},
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "email"}}},
						},
						Name: &emailStr,
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
			WhereClause: &ast.A_Expr{
				Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
				Lexpr: &ast.ColumnRef{
					Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
				},
				Rexpr: &ast.ParamRef{Number: 1, Location: 40},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// Check output columns
	if len(result.OutputColumns) != 2 {
		t.Fatalf("expected 2 output columns, got %d", len(result.OutputColumns))
	}
	if result.OutputColumns[0].Name != "name" {
		t.Errorf("col 0 name: got %q, want 'name'", result.OutputColumns[0].Name)
	}
	if result.OutputColumns[0].Type.Name != "text" {
		t.Errorf("col 0 type: got %q, want 'text'", result.OutputColumns[0].Type.Name)
	}
	if result.OutputColumns[1].Name != "email" {
		t.Errorf("col 1 name: got %q, want 'email'", result.OutputColumns[1].Name)
	}

	// Check parameter type inference
	if len(result.ParamTypes) != 1 {
		t.Fatalf("expected 1 param, got %d", len(result.ParamTypes))
	}
	p := result.ParamTypes[1]
	if p == nil {
		t.Fatal("param $1 not found")
	}
	if p.Type.Name != "integer" {
		t.Errorf("param $1 type: got %q, want 'integer'", p.Type.Name)
	}
}

func TestAnalyzeSelectWithAlias(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT u.name FROM users AS u WHERE u.id = $1
	nameStr := "name"
	relname := "users"
	aliasname := "u"
	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{
								&ast.String{Str: "u"},
								&ast.String{Str: "name"},
							}},
						},
						Name: &nameStr,
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{
						Relname: &relname,
						Alias:   &ast.Alias{Aliasname: &aliasname},
					},
				},
			},
			WhereClause: &ast.A_Expr{
				Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
				Lexpr: &ast.ColumnRef{
					Fields: &ast.List{Items: []ast.Node{
						&ast.String{Str: "u"},
						&ast.String{Str: "id"},
					}},
				},
				Rexpr: &ast.ParamRef{Number: 1, Location: 45},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if len(result.OutputColumns) != 1 {
		t.Fatalf("expected 1 output column, got %d", len(result.OutputColumns))
	}
	if result.OutputColumns[0].Type.Name != "text" {
		t.Errorf("col type: got %q, want 'text'", result.OutputColumns[0].Type.Name)
	}

	// $1 should be inferred as integer from u.id
	if p := result.ParamTypes[1]; p == nil {
		t.Error("param $1 not found")
	} else if p.Type.Name != "integer" {
		t.Errorf("param $1: got %q, want 'integer'", p.Type.Name)
	}
}

func TestAnalyzeSelectStar(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT * FROM users
	relname := "users"
	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{&ast.A_Star{}}},
						},
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// users has 5 columns: id, name, email, age, active
	if len(result.OutputColumns) != 5 {
		t.Fatalf("expected 5 output columns, got %d", len(result.OutputColumns))
	}

	// Check that types are resolved
	expectedCols := []struct {
		name     string
		typeName string
	}{
		{"id", "integer"},
		{"name", "text"},
		{"email", "text"},
		{"age", "integer"},
		{"active", "boolean"},
	}
	for i, expected := range expectedCols {
		if i >= len(result.OutputColumns) {
			break
		}
		col := result.OutputColumns[i]
		if col.Name != expected.name {
			t.Errorf("col %d name: got %q, want %q", i, col.Name, expected.name)
		}
		if col.Type.Name != expected.typeName {
			t.Errorf("col %d (%s) type: got %q, want %q", i, expected.name, col.Type.Name, expected.typeName)
		}
	}
}

func TestAnalyzeJoin(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT u.name, o.total FROM users AS u JOIN orders AS o ON u.id = o.user_id
	nameStr := "name"
	totalStr := "total"
	usersStr := "users"
	ordersStr := "orders"
	uAlias := "u"
	oAlias := "o"

	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{
								&ast.String{Str: "u"}, &ast.String{Str: "name"},
							}},
						},
						Name: &nameStr,
					},
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{
								&ast.String{Str: "o"}, &ast.String{Str: "total"},
							}},
						},
						Name: &totalStr,
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.JoinExpr{
						Jointype: ast.JoinTypeInner,
						Larg: &ast.RangeVar{
							Relname: &usersStr,
							Alias:   &ast.Alias{Aliasname: &uAlias},
						},
						Rarg: &ast.RangeVar{
							Relname: &ordersStr,
							Alias:   &ast.Alias{Aliasname: &oAlias},
						},
					},
				},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if len(result.OutputColumns) != 2 {
		t.Fatalf("expected 2 output columns, got %d", len(result.OutputColumns))
	}
	if result.OutputColumns[0].Type.Name != "text" {
		t.Errorf("u.name type: got %q, want 'text'", result.OutputColumns[0].Type.Name)
	}
	if result.OutputColumns[1].Type.Name != "numeric" {
		t.Errorf("o.total type: got %q, want 'numeric'", result.OutputColumns[1].Type.Name)
	}
}

func TestAnalyzeInsertParamTypes(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// INSERT INTO users (name, email) VALUES ($1, $2)
	relname := "users"
	nameCol := "name"
	emailCol := "email"

	stmt := &ast.RawStmt{
		Stmt: &ast.InsertStmt{
			Relation: &ast.RangeVar{Relname: &relname},
			Cols: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{Name: &nameCol},
					&ast.ResTarget{Name: &emailCol},
				},
			},
			SelectStmt: &ast.SelectStmt{
				ValuesLists: &ast.List{
					Items: []ast.Node{
						&ast.List{
							Items: []ast.Node{
								&ast.ParamRef{Number: 1, Location: 40},
								&ast.ParamRef{Number: 2, Location: 44},
							},
						},
					},
				},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// $1 should be text (name column), $2 should be text (email column)
	if len(result.ParamTypes) != 2 {
		t.Fatalf("expected 2 params, got %d", len(result.ParamTypes))
	}
	if p := result.ParamTypes[1]; p == nil {
		t.Error("$1 not found")
	} else if p.Type.Name != "text" {
		t.Errorf("$1: got %q, want 'text'", p.Type.Name)
	}
	if p := result.ParamTypes[2]; p == nil {
		t.Error("$2 not found")
	} else if p.Type.Name != "text" {
		t.Errorf("$2: got %q, want 'text'", p.Type.Name)
	}
}

func TestAnalyzeUpdate(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// UPDATE users SET name = $1 WHERE id = $2
	relname := "users"
	nameCol := "name"

	stmt := &ast.RawStmt{
		Stmt: &ast.UpdateStmt{
			Relations: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Name: &nameCol,
						Val:  &ast.ParamRef{Number: 1, Location: 25},
					},
				},
			},
			WhereClause: &ast.A_Expr{
				Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
				Lexpr: &ast.ColumnRef{
					Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
				},
				Rexpr: &ast.ParamRef{Number: 2, Location: 45},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// $1 should be text (name column type), $2 should be integer (id column type)
	if p := result.ParamTypes[1]; p == nil {
		t.Error("$1 not found")
	} else if p.Type.Name != "text" {
		t.Errorf("$1: got %q, want 'text'", p.Type.Name)
	}
	if p := result.ParamTypes[2]; p == nil {
		t.Error("$2 not found")
	} else if p.Type.Name != "integer" {
		t.Errorf("$2: got %q, want 'integer'", p.Type.Name)
	}
}

func TestAnalyzeDelete(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// DELETE FROM users WHERE id = $1
	relname := "users"
	stmt := &ast.RawStmt{
		Stmt: &ast.DeleteStmt{
			Relations: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
			WhereClause: &ast.A_Expr{
				Name: &ast.List{Items: []ast.Node{&ast.String{Str: "="}}},
				Lexpr: &ast.ColumnRef{
					Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "id"}}},
				},
				Rexpr: &ast.ParamRef{Number: 1, Location: 30},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	if p := result.ParamTypes[1]; p == nil {
		t.Error("$1 not found")
	} else if p.Type.Name != "integer" {
		t.Errorf("$1: got %q, want 'integer'", p.Type.Name)
	}
}

func TestAnalyzeLimitOffset(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT name FROM users LIMIT $1 OFFSET $2
	nameStr := "name"
	relname := "users"
	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "name"}}},
						},
						Name: &nameStr,
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
			LimitCount:  &ast.ParamRef{Number: 1, Location: 30},
			LimitOffset: &ast.ParamRef{Number: 2, Location: 40},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// Both LIMIT and OFFSET params should be integer
	if p := result.ParamTypes[1]; p == nil {
		t.Error("$1 (LIMIT) not found")
	} else if p.Type.Name != "integer" {
		t.Errorf("$1 (LIMIT): got %q, want 'integer'", p.Type.Name)
	}
	if p := result.ParamTypes[2]; p == nil {
		t.Error("$2 (OFFSET) not found")
	} else if p.Type.Name != "integer" {
		t.Errorf("$2 (OFFSET): got %q, want 'integer'", p.Type.Name)
	}
}

func TestScopeGraphStructure(t *testing.T) {
	cat := buildTestCatalog("public")
	a := New(cat, config.EnginePostgreSQL)

	// SELECT name FROM users
	nameStr := "name"
	relname := "users"
	stmt := &ast.RawStmt{
		Stmt: &ast.SelectStmt{
			TargetList: &ast.List{
				Items: []ast.Node{
					&ast.ResTarget{
						Val: &ast.ColumnRef{
							Fields: &ast.List{Items: []ast.Node{&ast.String{Str: "name"}}},
						},
						Name: &nameStr,
					},
				},
			},
			FromClause: &ast.List{
				Items: []ast.Node{
					&ast.RangeVar{Relname: &relname},
				},
			},
		},
	}

	result, err := a.AnalyzeQuery(stmt)
	if err != nil {
		t.Fatalf("AnalyzeQuery failed: %v", err)
	}

	// Verify scope graph structure
	rootScope := result.RootScope
	if rootScope.Kind != scope.ScopeSelect {
		t.Errorf("root scope kind: got %v, want scope.ScopeSelect", rootScope.Kind)
	}

	// Should have a parent edge to the FROM scope
	if len(rootScope.Edges) == 0 {
		t.Fatal("root scope has no edges")
	}

	parentEdge := rootScope.Edges[0]
	if parentEdge.Kind != scope.EdgeParent {
		t.Errorf("edge kind: got %v, want scope.EdgeParent", parentEdge.Kind)
	}

	fromScope := parentEdge.Target
	if fromScope.Kind != scope.ScopeFrom {
		t.Errorf("from scope kind: got %v, want scope.ScopeFrom", fromScope.Kind)
	}

	// FROM scope should have a table declaration for "users"
	if len(fromScope.Declarations) == 0 {
		t.Fatal("FROM scope has no declarations")
	}

	tableDecl := fromScope.Declarations[0]
	if tableDecl.Name != "users" {
		t.Errorf("table declaration: got %q, want 'users'", tableDecl.Name)
	}
	if tableDecl.Kind != scope.DeclTable {
		t.Errorf("declaration kind: got %v, want scope.DeclTable", tableDecl.Kind)
	}

	// The table's scope should have column declarations
	if tableDecl.Scope == nil {
		t.Fatal("table declaration has no scope")
	}
	if len(tableDecl.Scope.Declarations) != 5 {
		t.Errorf("table scope has %d declarations, want 5", len(tableDecl.Scope.Declarations))
	}
}

