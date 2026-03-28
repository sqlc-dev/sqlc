// Package scope implements scope graphs for SQL name resolution.
//
// A scope graph models the visibility and accessibility of names (columns,
// tables, aliases, functions) in a SQL query. Each scope is a node in the
// graph, containing declarations and connected to other scopes via labeled
// edges. Name resolution is path-finding in this graph.
//
// This approach is inspired by the Statix/scope graph framework from
// programming language theory, adapted for SQL's particular scoping rules.
package scope

import "fmt"

// EdgeKind labels the relationship between two scopes.
type EdgeKind int

const (
	// EdgeParent links a child scope to its parent (e.g., WHERE -> FROM).
	EdgeParent EdgeKind = iota
	// EdgeAlias links an alias name to the scope it refers to (e.g., "u" -> users table scope).
	EdgeAlias
	// EdgeLateral links a LATERAL subquery to preceding FROM items.
	EdgeLateral
	// EdgeOuter links a correlated subquery to its outer query scope.
	EdgeOuter
)

func (k EdgeKind) String() string {
	switch k {
	case EdgeParent:
		return "PARENT"
	case EdgeAlias:
		return "ALIAS"
	case EdgeLateral:
		return "LATERAL"
	case EdgeOuter:
		return "OUTER"
	default:
		return fmt.Sprintf("EdgeKind(%d)", int(k))
	}
}

// DeclKind describes what kind of entity a declaration represents.
type DeclKind int

const (
	DeclColumn DeclKind = iota
	DeclTable
	DeclCTE
	DeclFunction
	DeclAlias // A table alias (e.g., "u" in "FROM users AS u")
)

func (k DeclKind) String() string {
	switch k {
	case DeclColumn:
		return "column"
	case DeclTable:
		return "table"
	case DeclCTE:
		return "CTE"
	case DeclFunction:
		return "function"
	case DeclAlias:
		return "alias"
	default:
		return fmt.Sprintf("DeclKind(%d)", int(k))
	}
}

// Type represents a SQL type within the scope system. It's kept simple
// and engine-agnostic — detailed type information lives in the catalog.
type Type struct {
	Name     string // e.g., "integer", "text", "boolean"
	Schema   string // e.g., "pg_catalog", "" for default
	NotNull  bool
	IsArray  bool
	ArrayDims int
	Unsigned bool
	Length   *int
}

// Equals checks structural type equality (ignoring nullability).
func (t Type) Equals(other Type) bool {
	return t.Name == other.Name && t.Schema == other.Schema && t.IsArray == other.IsArray
}

// IsUnknown returns true if this type hasn't been determined yet.
func (t Type) IsUnknown() bool {
	return t.Name == "" || t.Name == "any"
}

var (
	TypeUnknown = Type{Name: "any"}
	TypeInt     = Type{Name: "integer", NotNull: true}
	TypeText    = Type{Name: "text", NotNull: true}
	TypeBool    = Type{Name: "boolean", NotNull: true}
	TypeFloat   = Type{Name: "float", NotNull: true}
	TypeNumeric = Type{Name: "numeric", NotNull: true}
)

// Declaration is a named entity visible within a scope.
type Declaration struct {
	Name     string
	Kind     DeclKind
	Type     Type
	Scope    *Scope // For table/CTE declarations, the scope containing their columns
	Location int    // Source position for error reporting
}

// Edge connects one scope to another with a labeled relationship.
type Edge struct {
	Kind   EdgeKind
	Label  string // For EdgeAlias, the alias name
	Target *Scope
}

// ScopeKind describes the syntactic context that created this scope.
type ScopeKind int

const (
	ScopeRoot ScopeKind = iota
	ScopeFrom
	ScopeJoin
	ScopeWhere
	ScopeSelect
	ScopeHaving
	ScopeOrderBy
	ScopeSubquery
	ScopeCTE
	ScopeInsert
	ScopeUpdate
	ScopeDelete
	ScopeValues
	ScopeReturning
	ScopeFunction
)

func (k ScopeKind) String() string {
	names := [...]string{
		"ROOT", "FROM", "JOIN", "WHERE", "SELECT", "HAVING",
		"ORDER_BY", "SUBQUERY", "CTE", "INSERT", "UPDATE",
		"DELETE", "VALUES", "RETURNING", "FUNCTION",
	}
	if int(k) < len(names) {
		return names[k]
	}
	return fmt.Sprintf("ScopeKind(%d)", int(k))
}

// Scope is a node in the scope graph. It contains declarations and
// edges to other scopes.
type Scope struct {
	Kind         ScopeKind
	Declarations []*Declaration
	Edges        []*Edge
	Location     int // Source position of the construct that created this scope
}

// NewScope creates a new empty scope of the given kind.
func NewScope(kind ScopeKind) *Scope {
	return &Scope{
		Kind: kind,
	}
}

// Declare adds a declaration to this scope.
func (s *Scope) Declare(d *Declaration) {
	s.Declarations = append(s.Declarations, d)
}

// DeclareColumn is a convenience method for declaring a column.
func (s *Scope) DeclareColumn(name string, typ Type, location int) *Declaration {
	d := &Declaration{
		Name:     name,
		Kind:     DeclColumn,
		Type:     typ,
		Location: location,
	}
	s.Declare(d)
	return d
}

// DeclareTable adds a table declaration with its own column scope.
func (s *Scope) DeclareTable(name string, columnScope *Scope, location int) *Declaration {
	d := &Declaration{
		Name:     name,
		Kind:     DeclTable,
		Type:     TypeUnknown,
		Scope:    columnScope,
		Location: location,
	}
	s.Declare(d)
	return d
}

// AddEdge connects this scope to another scope via a labeled edge.
func (s *Scope) AddEdge(kind EdgeKind, label string, target *Scope) {
	s.Edges = append(s.Edges, &Edge{
		Kind:   kind,
		Label:  label,
		Target: target,
	})
}

// AddParent connects this scope to a parent scope.
func (s *Scope) AddParent(parent *Scope) {
	s.AddEdge(EdgeParent, "", parent)
}

// AddAlias connects an alias name to a target scope.
func (s *Scope) AddAlias(alias string, target *Scope) {
	s.AddEdge(EdgeAlias, alias, target)
}
