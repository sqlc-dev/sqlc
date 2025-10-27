package ydb

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	parser "github.com/ydb-platform/yql-parsers/go"
)

type objectRefProvider interface {
	antlr.ParserRuleContext
	Object_ref() parser.IObject_refContext
}

func parseTableName(ctx objectRefProvider) *ast.TableName {
	return parseObjectRef(ctx.Object_ref())
}

func parseObjectRef(r parser.IObject_refContext) *ast.TableName {
	if r == nil {
		return nil
	}
	ref := r.(*parser.Object_refContext)

	parts := []string{}

	if cl := ref.Cluster_expr(); cl != nil {
		parts = append(parts, parseClusterExpr(cl))
	}

	if idOrAt := ref.Id_or_at(); idOrAt != nil {
		parts = append(parts, parseIdOrAt(idOrAt))
	}

	objectName := strings.Join(parts, ".")

	return &ast.TableName{
		Schema: "",
		Name:   identifier(objectName),
	}
}

func parseClusterExpr(ctx parser.ICluster_exprContext) string {
	if ctx == nil {
		return ""
	}
	return identifier(ctx.GetText())
}

func parseIdOrAt(ctx parser.IId_or_atContext) string {
	if ctx == nil {
		return ""
	}
	idOrAt := ctx.(*parser.Id_or_atContext)

	if ao := idOrAt.An_id_or_type(); ao != nil {
		return identifier(parseAnIdOrType(ao))
	}
	return ""
}

func parseAnIdOrType(ctx parser.IAn_id_or_typeContext) string {
	if ctx == nil {
		return ""
	}
	anId := ctx.(*parser.An_id_or_typeContext)

	if anId.Id_or_type() != nil {
		return identifier(parseIdOrType(anId.Id_or_type()))
	}

	if anId.STRING_VALUE() != nil {
		return identifier(anId.STRING_VALUE().GetText())
	}

	return ""
}

func parseIdOrType(ctx parser.IId_or_typeContext) string {
	if ctx == nil {
		return ""
	}
	Id := ctx.(*parser.Id_or_typeContext)
	if Id.Id() != nil {
		return identifier(parseId(Id.Id()))
	}

	return ""
}

func parseAnId(ctx parser.IAn_idContext) string {
	if id := ctx.Id(); id != nil {
		return id.GetText()
	} else if str := ctx.STRING_VALUE(); str != nil {
		return str.GetText()
	}
	return ""
}

func parseAnIdSchema(ctx parser.IAn_id_schemaContext) string {
	if ctx == nil {
		return ""
	}
	if id := ctx.Id_schema(); id != nil {
		return id.GetText()
	} else if str := ctx.STRING_VALUE(); str != nil {
		return str.GetText()
	}
	return ""
}

func parseId(ctx parser.IIdContext) string {
	if ctx == nil {
		return ""
	}
	return ctx.GetText()
}

func parseAnIdTable(ctx parser.IAn_id_tableContext) string {
	if ctx == nil {
		return ""
	}
	if id := ctx.Id_table(); id != nil {
		return id.GetText()
	} else if str := ctx.STRING_VALUE(); str != nil {
		return str.GetText()
	}
	return ""
}

func parseIntegerValue(text string) (int64, error) {
	text = strings.ToLower(text)
	base := 10

	switch {
	case strings.HasPrefix(text, "0x"):
		base = 16
		text = strings.TrimPrefix(text, "0x")

	case strings.HasPrefix(text, "0o"):
		base = 8
		text = strings.TrimPrefix(text, "0o")

	case strings.HasPrefix(text, "0b"):
		base = 2
		text = strings.TrimPrefix(text, "0b")
	}

	// debug!!!
	text = strings.TrimRight(text, "pulstibn")

	return strconv.ParseInt(text, base, 64)
}

func (c *cc) extractRoleSpec(n parser.IRole_nameContext, roletype ast.RoleSpecType) (*ast.RoleSpec, bool, ast.Node) {
	if n == nil {
		return nil, false, nil
	}
	roleNode, ok := n.Accept(c).(ast.Node)
	if !ok {
		return nil, false, nil
	}

	roleSpec := &ast.RoleSpec{
		Roletype: roletype,
		Location: n.GetStart().GetStart(),
	}

	isParam := true
	switch v := roleNode.(type) {
	case *ast.A_Const:
		switch val := v.Val.(type) {
		case *ast.String:
			roleSpec.Rolename = &val.Str
			isParam = false
		case *ast.Boolean:
			roleSpec.BindRolename = roleNode
		default:
			return nil, false, nil
		}
	case *ast.ParamRef, *ast.A_Expr:
		roleSpec.BindRolename = roleNode
	default:
		return nil, false, nil
	}

	return roleSpec, isParam, roleNode
}

func byteOffset(s string, runeIndex int) int {
	count := 0
	for i := range s {
		if count == runeIndex {
			return i
		}
		count++
	}
	return len(s)
}

func byteOffsetFromRuneIndex(s string, runeIndex int) int {
	if runeIndex <= 0 {
		return 0
	}
	bytePos := 0
	for i := 0; i < runeIndex && bytePos < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[bytePos:])
		bytePos += size
	}
	return bytePos
}

func emptySelectStmt() *ast.SelectStmt {
	return &ast.SelectStmt{
		DistinctClause: &ast.List{},
		TargetList:     &ast.List{},
		FromClause:     &ast.List{},
		GroupClause:    &ast.List{},
		WindowClause:   &ast.List{},
		ValuesLists:    &ast.List{},
		SortClause:     &ast.List{},
		LockingClause:  &ast.List{},
	}
}

func (c *cc) collectComparisonOps(n parser.IEq_subexprContext) []antlr.TerminalNode {
	var ops []antlr.TerminalNode
	for _, child := range n.GetChildren() {
		if tn, ok := child.(antlr.TerminalNode); ok {
			switch tn.GetText() {
			case "<", "<=", ">", ">=":
				ops = append(ops, tn)
			}
		}
	}
	return ops
}

func (c *cc) collectBitwiseOps(ctx parser.INeq_subexprContext) []antlr.TerminalNode {
	var ops []antlr.TerminalNode
	children := ctx.GetChildren()
	for _, child := range children {
		if tn, ok := child.(antlr.TerminalNode); ok {
			txt := tn.GetText()
			switch txt {
			case "<<", ">>", "<<|", ">>|", "&", "|", "^":
				ops = append(ops, tn)
			}
		}
	}
	return ops
}

func (c *cc) collectBitOps(ctx parser.IBit_subexprContext) []antlr.TerminalNode {
	var ops []antlr.TerminalNode
	children := ctx.GetChildren()
	for _, child := range children {
		if tn, ok := child.(antlr.TerminalNode); ok {
			txt := tn.GetText()
			switch txt {
			case "+", "-":
				ops = append(ops, tn)
			}
		}
	}
	return ops
}

func (c *cc) collectAddOps(ctx parser.IAdd_subexprContext) []antlr.TerminalNode {
	var ops []antlr.TerminalNode
	for _, child := range ctx.GetChildren() {
		if tn, ok := child.(antlr.TerminalNode); ok {
			switch tn.GetText() {
			case "*", "/", "%":
				ops = append(ops, tn)
			}
		}
	}
	return ops
}

func (c *cc) collectEqualityOps(ctx parser.ICond_exprContext) []antlr.TerminalNode {
	var ops []antlr.TerminalNode
	children := ctx.GetChildren()
	for _, child := range children {
		if tn, ok := child.(antlr.TerminalNode); ok {
			switch tn.GetText() {
			case "=", "==", "!=", "<>":
				ops = append(ops, tn)
			}
		}
	}
	return ops
}

// parseStringLiteral parses a string literal from a YQL query and returns the value and whether it has a suffix.
// If a valid suffix is found, it is stripped and the content is returned.
// FIXME: rewrite this logic to correctly handle the type based on the suffix.
func parseStringLiteral(s string) (value string, hasSuffix bool) {
	if len(s) < 2 {
		return s, false
	}

	quote := s[0]
	if quote != '\'' && quote != '"' {
		return s, false
	}

	quotePos := -1
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == quote {
			quotePos = i
			break
		}
	}

	if quotePos == -1 || quotePos == 0 {
		return s, false
	}

	content := s[1:quotePos]

	if quotePos < len(s)-1 {
		suffix := s[quotePos+1:]
		if isValidYDBStringSuffix(suffix) {
			return content, true
		}
	}

	return content, false
}

func isValidYDBStringSuffix(suffix string) bool {
	switch suffix {
	case "s", "S", "u", "U", "y", "Y", "j", "J",
		"pt", "PT", "pb", "PB", "pv", "PV":
		return true
	default:
		return false
	}
}
