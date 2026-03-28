package scope

import (
	"testing"
)

func TestResolveUnqualified(t *testing.T) {
	// Build scope graph:
	//   [FROM scope] has table "users" with columns {id, name, email}
	//   [SELECT scope] → PARENT → [FROM scope]

	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", Type{Name: "integer", NotNull: true}, 0)
	usersScope.DeclareColumn("name", Type{Name: "text", NotNull: true}, 0)
	usersScope.DeclareColumn("email", Type{Name: "text", NotNull: false}, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.DeclareTable("users", usersScope, 0)

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	// Resolve "name" from SELECT scope
	resolved, err := selectScope.Resolve("name")
	if err != nil {
		t.Fatalf("expected to resolve 'name', got error: %v", err)
	}
	if resolved.Declaration.Name != "name" {
		t.Errorf("expected declaration name 'name', got %q", resolved.Declaration.Name)
	}
	if resolved.Declaration.Type.Name != "text" {
		t.Errorf("expected type 'text', got %q", resolved.Declaration.Type.Name)
	}
	if !resolved.Declaration.Type.NotNull {
		t.Error("expected 'name' to be NOT NULL")
	}
}

func TestResolveQualified(t *testing.T) {
	// SELECT u.name FROM users AS u

	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", Type{Name: "integer", NotNull: true}, 0)
	usersScope.DeclareColumn("name", Type{Name: "text", NotNull: true}, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.AddAlias("u", usersScope)
	fromScope.Declare(&Declaration{
		Name:  "u",
		Kind:  DeclAlias,
		Type:  TypeUnknown,
		Scope: usersScope,
	})

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	// Resolve "u.name"
	resolved, err := selectScope.ResolveQualified("u", "name")
	if err != nil {
		t.Fatalf("expected to resolve 'u.name', got error: %v", err)
	}
	if resolved.Declaration.Name != "name" {
		t.Errorf("expected 'name', got %q", resolved.Declaration.Name)
	}
}

func TestResolveAmbiguous(t *testing.T) {
	// SELECT id FROM users JOIN orders ON ...
	// Both tables have an 'id' column → should be ambiguous

	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", TypeInt, 0)
	usersScope.DeclareColumn("name", TypeText, 0)

	ordersScope := NewScope(ScopeFrom)
	ordersScope.DeclareColumn("id", TypeInt, 0)
	ordersScope.DeclareColumn("total", TypeNumeric, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.DeclareTable("users", usersScope, 0)
	fromScope.DeclareTable("orders", ordersScope, 0)

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	_, err := selectScope.Resolve("id")
	if err == nil {
		t.Fatal("expected ambiguity error for 'id', got nil")
	}
	resErr, ok := err.(*ResolutionError)
	if !ok {
		t.Fatalf("expected *ResolutionError, got %T", err)
	}
	if resErr.Kind != ErrAmbiguous {
		t.Errorf("expected ErrAmbiguous, got %v", resErr.Kind)
	}
}

func TestResolveNotFound(t *testing.T) {
	fromScope := NewScope(ScopeFrom)
	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	_, err := selectScope.Resolve("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent column")
	}
	resErr, ok := err.(*ResolutionError)
	if !ok {
		t.Fatalf("expected *ResolutionError, got %T", err)
	}
	if resErr.Kind != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", resErr.Kind)
	}
}

func TestResolveQualifierNotFound(t *testing.T) {
	fromScope := NewScope(ScopeFrom)
	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	_, err := selectScope.ResolveQualified("nonexistent", "col")
	if err == nil {
		t.Fatal("expected error for nonexistent qualifier")
	}
	resErr, ok := err.(*ResolutionError)
	if !ok {
		t.Fatalf("expected *ResolutionError, got %T", err)
	}
	if resErr.Kind != ErrQualifierNotFound {
		t.Errorf("expected ErrQualifierNotFound, got %v", resErr.Kind)
	}
}

func TestResolveColumnRef(t *testing.T) {
	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", TypeInt, 0)
	usersScope.DeclareColumn("name", TypeText, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.AddAlias("u", usersScope)
	fromScope.Declare(&Declaration{Name: "u", Kind: DeclAlias, Scope: usersScope})

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	tests := []struct {
		parts    []string
		wantName string
		wantErr  bool
	}{
		{[]string{"name"}, "name", false},
		{[]string{"u", "name"}, "name", false},
		{[]string{"public", "u", "name"}, "name", false},
		{[]string{"nonexistent"}, "", true},
		{[]string{"u", "nonexistent"}, "", true},
	}

	for _, tt := range tests {
		resolved, err := selectScope.ResolveColumnRef(tt.parts)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ResolveColumnRef(%v): expected error, got nil", tt.parts)
			}
			continue
		}
		if err != nil {
			t.Errorf("ResolveColumnRef(%v): unexpected error: %v", tt.parts, err)
			continue
		}
		if resolved.Declaration.Name != tt.wantName {
			t.Errorf("ResolveColumnRef(%v): got name %q, want %q", tt.parts, resolved.Declaration.Name, tt.wantName)
		}
	}
}

func TestAllColumns(t *testing.T) {
	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", TypeInt, 0)
	usersScope.DeclareColumn("name", TypeText, 0)

	ordersScope := NewScope(ScopeFrom)
	ordersScope.DeclareColumn("id", TypeInt, 0)
	ordersScope.DeclareColumn("total", TypeNumeric, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.AddAlias("u", usersScope)
	fromScope.Declare(&Declaration{Name: "u", Kind: DeclAlias, Scope: usersScope})
	fromScope.AddAlias("o", ordersScope)
	fromScope.Declare(&Declaration{Name: "o", Kind: DeclAlias, Scope: ordersScope})

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	// All columns (SELECT *)
	all := selectScope.AllColumns("")
	if len(all) != 4 {
		t.Errorf("AllColumns(''): got %d columns, want 4", len(all))
	}

	// Qualified (SELECT u.*)
	uCols := selectScope.AllColumns("u")
	if len(uCols) != 2 {
		t.Errorf("AllColumns('u'): got %d columns, want 2", len(uCols))
	}
	for _, c := range uCols {
		if c.Name != "id" && c.Name != "name" {
			t.Errorf("AllColumns('u'): unexpected column %q", c.Name)
		}
	}
}

func TestCTEScope(t *testing.T) {
	// WITH active_users AS (SELECT ...) SELECT * FROM active_users

	cteScope := NewScope(ScopeCTE)
	cteScope.DeclareColumn("id", TypeInt, 0)
	cteScope.DeclareColumn("name", TypeText, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.Declare(&Declaration{
		Name:  "active_users",
		Kind:  DeclCTE,
		Type:  TypeUnknown,
		Scope: cteScope,
	})

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	// Resolve "name" from the CTE
	resolved, err := selectScope.Resolve("name")
	if err != nil {
		t.Fatalf("expected to resolve 'name' from CTE, got error: %v", err)
	}
	if resolved.Declaration.Type.Name != "text" {
		t.Errorf("expected type 'text', got %q", resolved.Declaration.Type.Name)
	}

	// Resolve qualified "active_users.id"
	resolved, err = selectScope.ResolveQualified("active_users", "id")
	if err != nil {
		t.Fatalf("expected to resolve 'active_users.id', got error: %v", err)
	}
	if resolved.Declaration.Type.Name != "integer" {
		t.Errorf("expected type 'integer', got %q", resolved.Declaration.Type.Name)
	}
}

func TestJoinScope(t *testing.T) {
	// SELECT u.name, o.total
	// FROM users AS u
	// JOIN orders AS o ON u.id = o.user_id

	usersScope := NewScope(ScopeFrom)
	usersScope.DeclareColumn("id", TypeInt, 0)
	usersScope.DeclareColumn("name", TypeText, 0)

	ordersScope := NewScope(ScopeFrom)
	ordersScope.DeclareColumn("id", TypeInt, 0)
	ordersScope.DeclareColumn("user_id", TypeInt, 0)
	ordersScope.DeclareColumn("total", TypeNumeric, 0)

	fromScope := NewScope(ScopeFrom)
	fromScope.AddAlias("u", usersScope)
	fromScope.Declare(&Declaration{Name: "u", Kind: DeclAlias, Scope: usersScope})
	fromScope.AddAlias("o", ordersScope)
	fromScope.Declare(&Declaration{Name: "o", Kind: DeclAlias, Scope: ordersScope})

	selectScope := NewScope(ScopeSelect)
	selectScope.AddParent(fromScope)

	// u.name should resolve
	resolved, err := selectScope.ResolveQualified("u", "name")
	if err != nil {
		t.Fatalf("u.name: %v", err)
	}
	if resolved.Declaration.Type.Name != "text" {
		t.Errorf("u.name type: got %q, want 'text'", resolved.Declaration.Type.Name)
	}

	// o.total should resolve
	resolved, err = selectScope.ResolveQualified("o", "total")
	if err != nil {
		t.Fatalf("o.total: %v", err)
	}
	if resolved.Declaration.Type.Name != "numeric" {
		t.Errorf("o.total type: got %q, want 'numeric'", resolved.Declaration.Type.Name)
	}

	// Unqualified "total" should resolve (only in orders)
	resolved, err = selectScope.Resolve("total")
	if err != nil {
		t.Fatalf("total: %v", err)
	}
	if resolved.Declaration.Type.Name != "numeric" {
		t.Errorf("total type: got %q, want 'numeric'", resolved.Declaration.Type.Name)
	}

	// Unqualified "id" should be ambiguous
	_, err = selectScope.Resolve("id")
	if err == nil {
		t.Fatal("expected ambiguity for 'id'")
	}
}
