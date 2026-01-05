// Package poet provides Go code generation with custom AST nodes
// that properly support comment placement.
package poet

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
