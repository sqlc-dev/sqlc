package poet

import (
	"go/ast"
	"go/token"
)

// Assign creates a simple assignment statement (lhs = rhs).
func Assign(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{rhs},
	}
}

// AssignMulti creates a multi-value assignment statement (lhs1, lhs2 = rhs1, rhs2).
func AssignMulti(lhs []ast.Expr, rhs []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.ASSIGN,
		Rhs: rhs,
	}
}

// Define creates a short variable declaration (lhs := rhs).
func Define(lhs, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}

// DefineMulti creates a multi-value short variable declaration.
func DefineMulti(lhs []ast.Expr, rhs []ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}

// DefineNames creates a short variable declaration with named variables.
func DefineNames(names []string, rhs ast.Expr) *ast.AssignStmt {
	var lhs []ast.Expr
	for _, name := range names {
		lhs = append(lhs, Ident(name))
	}
	return &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{rhs},
	}
}

// DeclStmt creates a declaration statement.
func DeclStmt(decl ast.Decl) *ast.DeclStmt {
	return &ast.DeclStmt{Decl: decl}
}

// ExprStmt creates an expression statement.
func ExprStmt(expr ast.Expr) *ast.ExprStmt {
	return &ast.ExprStmt{X: expr}
}

// Return creates a return statement.
func Return(results ...ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{Results: results}
}

// If creates an if statement.
func If(cond ast.Expr, body ...ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
	}
}

// IfInit creates an if statement with an init clause.
func IfInit(init ast.Stmt, cond ast.Expr, body ...ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{
		Init: init,
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
	}
}

// IfElse creates an if-else statement.
func IfElse(cond ast.Expr, body []ast.Stmt, elseBody []ast.Stmt) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
		Else: &ast.BlockStmt{List: elseBody},
	}
}

// IfElseIf creates an if-else if chain.
func IfElseIf(cond ast.Expr, body []ast.Stmt, elseStmt *ast.IfStmt) *ast.IfStmt {
	return &ast.IfStmt{
		Cond: cond,
		Body: &ast.BlockStmt{List: body},
		Else: elseStmt,
	}
}

// For creates a for loop.
func For(init ast.Stmt, cond ast.Expr, post ast.Stmt, body ...ast.Stmt) *ast.ForStmt {
	return &ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: &ast.BlockStmt{List: body},
	}
}

// ForRange creates a for-range loop.
func ForRange(key, value, x ast.Expr, body ...ast.Stmt) *ast.RangeStmt {
	return &ast.RangeStmt{
		Key:   key,
		Value: value,
		Tok:   token.DEFINE,
		X:     x,
		Body:  &ast.BlockStmt{List: body},
	}
}

// ForRangeAssign creates a for-range loop with assignment (=).
func ForRangeAssign(key, value, x ast.Expr, body ...ast.Stmt) *ast.RangeStmt {
	return &ast.RangeStmt{
		Key:   key,
		Value: value,
		Tok:   token.ASSIGN,
		X:     x,
		Body:  &ast.BlockStmt{List: body},
	}
}

// Switch creates a switch statement.
func Switch(tag ast.Expr, body ...ast.Stmt) *ast.SwitchStmt {
	return &ast.SwitchStmt{
		Tag:  tag,
		Body: &ast.BlockStmt{List: body},
	}
}

// SwitchInit creates a switch statement with an init clause.
func SwitchInit(init ast.Stmt, tag ast.Expr, body ...ast.Stmt) *ast.SwitchStmt {
	return &ast.SwitchStmt{
		Init: init,
		Tag:  tag,
		Body: &ast.BlockStmt{List: body},
	}
}

// TypeSwitch creates a type switch statement.
func TypeSwitch(assign ast.Stmt, body ...ast.Stmt) *ast.TypeSwitchStmt {
	return &ast.TypeSwitchStmt{
		Assign: assign,
		Body:   &ast.BlockStmt{List: body},
	}
}

// Case creates a case clause for switch statements.
func Case(list []ast.Expr, body ...ast.Stmt) *ast.CaseClause {
	return &ast.CaseClause{
		List: list,
		Body: body,
	}
}

// Default creates a default case clause.
func Default(body ...ast.Stmt) *ast.CaseClause {
	return &ast.CaseClause{
		List: nil,
		Body: body,
	}
}

// Block creates a block statement.
func Block(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}

// Defer creates a defer statement.
func Defer(call *ast.CallExpr) *ast.DeferStmt {
	return &ast.DeferStmt{Call: call}
}

// Go creates a go statement.
func Go(call *ast.CallExpr) *ast.GoStmt {
	return &ast.GoStmt{Call: call}
}

// Send creates a channel send statement.
func Send(ch, value ast.Expr) *ast.SendStmt {
	return &ast.SendStmt{Chan: ch, Value: value}
}

// Inc creates an increment statement (x++).
func Inc(x ast.Expr) *ast.IncDecStmt {
	return &ast.IncDecStmt{X: x, Tok: token.INC}
}

// Dec creates a decrement statement (x--).
func Dec(x ast.Expr) *ast.IncDecStmt {
	return &ast.IncDecStmt{X: x, Tok: token.DEC}
}

// Break creates a break statement.
func Break() *ast.BranchStmt {
	return &ast.BranchStmt{Tok: token.BREAK}
}

// BreakLabel creates a break statement with a label.
func BreakLabel(label string) *ast.BranchStmt {
	return &ast.BranchStmt{Tok: token.BREAK, Label: ast.NewIdent(label)}
}

// Continue creates a continue statement.
func Continue() *ast.BranchStmt {
	return &ast.BranchStmt{Tok: token.CONTINUE}
}

// ContinueLabel creates a continue statement with a label.
func ContinueLabel(label string) *ast.BranchStmt {
	return &ast.BranchStmt{Tok: token.CONTINUE, Label: ast.NewIdent(label)}
}

// Goto creates a goto statement.
func Goto(label string) *ast.BranchStmt {
	return &ast.BranchStmt{Tok: token.GOTO, Label: ast.NewIdent(label)}
}

// Label creates a labeled statement.
func Label(name string, stmt ast.Stmt) *ast.LabeledStmt {
	return &ast.LabeledStmt{Label: ast.NewIdent(name), Stmt: stmt}
}

// Empty creates an empty statement.
func Empty() *ast.EmptyStmt {
	return &ast.EmptyStmt{}
}

// Select creates a select statement.
func Select(body ...ast.Stmt) *ast.SelectStmt {
	return &ast.SelectStmt{Body: &ast.BlockStmt{List: body}}
}

// CommClause creates a communication clause for select statements.
func CommClause(comm ast.Stmt, body ...ast.Stmt) *ast.CommClause {
	return &ast.CommClause{Comm: comm, Body: body}
}

// CommDefault creates a default communication clause.
func CommDefault(body ...ast.Stmt) *ast.CommClause {
	return &ast.CommClause{Comm: nil, Body: body}
}
