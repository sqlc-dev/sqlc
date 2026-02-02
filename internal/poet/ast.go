// Package poet provides Go code generation with custom AST nodes
// that properly support comment placement.
package poet

import "strings"

// File represents a Go source file.
type File struct {
	BuildTags    string
	Comments     []string // File-level comments
	Package      string
	ImportGroups [][]Import // Groups separated by blank lines
	Decls        []Decl
}

// Import represents an import statement.
type Import struct {
	Alias string // Optional alias
	Path  string
}

// Decl represents a declaration.
type Decl interface {
	isDecl()
}

// Raw is raw Go code (escape hatch).
type Raw struct {
	Code string
}

func (Raw) isDecl() {}

// Const represents a const declaration.
type Const struct {
	Comment string
	Name    string
	Type    string
	Value   string
}

func (Const) isDecl() {}

// ConstBlock represents a const block.
type ConstBlock struct {
	Consts []Const
}

func (ConstBlock) isDecl() {}

// Var represents a var declaration.
type Var struct {
	Comment string
	Name    string
	Type    string
	Value   string
}

func (Var) isDecl() {}

// VarBlock represents a var block.
type VarBlock struct {
	Vars []Var
}

func (VarBlock) isDecl() {}

// TypeDef represents a type declaration.
type TypeDef struct {
	Comment string
	Name    string
	Type    TypeExpr
}

func (TypeDef) isDecl() {}

// Func represents a function declaration.
type Func struct {
	Comment string
	Recv    *Param // nil for non-methods
	Name    string
	Params  []Param
	Results []Param
	Stmts   []Stmt
}

func (Func) isDecl() {}

// Param represents a function parameter or result.
type Param struct {
	Name    string
	Type    string
	Pointer bool // If true, type is rendered as *Type
}

// TypeExpr represents a type expression.
type TypeExpr interface {
	isTypeExpr()
}

// Struct represents a struct type.
type Struct struct {
	Fields []Field
}

func (Struct) isTypeExpr() {}

// Field represents a struct field.
type Field struct {
	Comment         string // Leading comment (above the field)
	Name            string
	Type            string
	Tag             string
	TrailingComment string // Trailing comment (on same line)
}

// Interface represents an interface type.
type Interface struct {
	Methods []Method
}

func (Interface) isTypeExpr() {}

// Method represents an interface method.
type Method struct {
	Comment string
	Name    string
	Params  []Param
	Results []Param
}

// TypeName represents a type alias or named type.
type TypeName struct {
	Name string
}

func (TypeName) isTypeExpr() {}

// Stmt represents a statement in a function body.
type Stmt interface {
	isStmt()
}

// RawStmt is raw Go code as a statement.
type RawStmt struct {
	Code string
}

func (RawStmt) isStmt() {}

// Return represents a return statement.
type Return struct {
	Values []string // Expressions to return
}

func (Return) isStmt() {}

// For represents a for loop.
type For struct {
	Init  string // e.g., "i := 0"
	Cond  string // e.g., "i < 10"
	Post  string // e.g., "i++"
	Range string // If set, renders as "for Range {" (e.g., "_, v := range items")
	Body  []Stmt
}

func (For) isStmt() {}

// If represents an if statement.
type If struct {
	Init string // Optional init statement (e.g., "err := foo()")
	Cond string // Condition expression
	Body []Stmt
	Else []Stmt // Optional else body
}

func (If) isStmt() {}

// Switch represents a switch statement.
type Switch struct {
	Init  string // Optional init statement
	Expr  string // Expression to switch on (empty for type switch or bool switch)
	Cases []Case
}

func (Switch) isStmt() {}

// Case represents a case clause in a switch statement.
type Case struct {
	Values []string // Case values (empty for default case)
	Body   []Stmt
}

// Defer represents a defer statement.
type Defer struct {
	Call string // The function call to defer
}

func (Defer) isStmt() {}

// Assign represents an assignment statement.
type Assign struct {
	Left  []string // Left-hand side (variable names)
	Op    string   // Assignment operator: "=", ":=", "+=", etc.
	Right []string // Right-hand side (expressions)
}

func (Assign) isStmt() {}

// CallStmt represents a function call as a statement.
type CallStmt struct {
	Call string // The function call expression
}

func (CallStmt) isStmt() {}

// VarDecl represents a variable declaration statement.
type VarDecl struct {
	Name  string // Variable name
	Type  string // Type (optional if Value is set)
	Value string // Initial value (optional)
}

func (VarDecl) isStmt() {}

// GoStmt represents a go statement (goroutine).
type GoStmt struct {
	Call string // The function call to run as a goroutine
}

func (GoStmt) isStmt() {}

// Continue represents a continue statement.
type Continue struct {
	Label string // Optional label
}

func (Continue) isStmt() {}

// Break represents a break statement.
type Break struct {
	Label string // Optional label
}

func (Break) isStmt() {}

// Expr is an interface for expression types that can be rendered to strings.
// These can be used in Return.Values, Assign.Right, etc.
type Expr interface {
	Render() string
}

// CallExpr represents a function or method call expression.
type CallExpr struct {
	Func string   // Function name or receiver.method
	Args []string // Arguments
}

func (c CallExpr) Render() string {
	var b strings.Builder
	b.WriteString(c.Func)
	b.WriteString("(")
	for i, arg := range c.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(arg)
	}
	b.WriteString(")")
	return b.String()
}

// StructLit represents a struct literal expression.
type StructLit struct {
	Type      string      // Type name (e.g., "Queries")
	Pointer   bool        // If true, prefix with &
	Multiline bool        // If true, always use multi-line format
	Fields    [][2]string // Field name-value pairs (use slice to preserve order)
}

func (s StructLit) Render() string {
	var b strings.Builder
	if s.Pointer {
		b.WriteString("&")
	}
	b.WriteString(s.Type)
	b.WriteString("{")
	if len(s.Fields) <= 2 && !s.Multiline {
		// Compact format for small struct literals
		for i, f := range s.Fields {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(f[0])
			b.WriteString(": ")
			b.WriteString(f[1])
		}
	} else if len(s.Fields) > 0 {
		// Multi-line format for larger struct literals or when explicitly requested
		b.WriteString("\n")
		for _, f := range s.Fields {
			b.WriteString("\t\t")
			b.WriteString(f[0])
			b.WriteString(": ")
			b.WriteString(f[1])
			b.WriteString(",\n")
		}
		b.WriteString("\t")
	}
	b.WriteString("}")
	return b.String()
}

// SliceLit represents a slice literal expression.
type SliceLit struct {
	Type      string   // Element type (e.g., "interface{}")
	Multiline bool     // If true, always use multi-line format
	Values    []string // Elements
}

func (s SliceLit) Render() string {
	var b strings.Builder
	b.WriteString("[]")
	b.WriteString(s.Type)
	b.WriteString("{")
	if len(s.Values) <= 3 && !s.Multiline {
		// Compact format for small slice literals
		for i, v := range s.Values {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(v)
		}
	} else if len(s.Values) > 0 {
		// Multi-line format for larger slice literals or when explicitly requested
		b.WriteString("\n")
		for _, v := range s.Values {
			b.WriteString("\t\t")
			b.WriteString(v)
			b.WriteString(",\n")
		}
		b.WriteString("\t")
	}
	b.WriteString("}")
	return b.String()
}

// TypeCast represents a type conversion expression.
type TypeCast struct {
	Type  string // Target type
	Value string // Value to convert
}

func (t TypeCast) Render() string {
	return t.Type + "(" + t.Value + ")"
}

// FuncLit represents an anonymous function literal.
type FuncLit struct {
	Params  []Param
	Results []Param
	Body    []Stmt
	Indent  string // Base indentation for body statements (default: "\t")
}

// Note: FuncLit.Render() is implemented in render.go since it needs renderStmts

// Selector represents a field or method selector expression (a.b.c).
type Selector struct {
	Parts []string // e.g., ["r", "rows", "0", "Field"] for r.rows[0].Field
}

func (s Selector) Render() string {
	return strings.Join(s.Parts, ".")
}

// Index represents an index or slice expression.
type Index struct {
	Expr  string // Base expression
	Index string // Index value (or "start:end" for slice)
}

func (i Index) Render() string {
	return i.Expr + "[" + i.Index + "]"
}
