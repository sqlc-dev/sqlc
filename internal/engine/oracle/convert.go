package oracle

import (
	"log"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/engine/oracle/parser"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// cc is the Oracle parse-tree -> shared-AST converter. It mirrors the SQLite
// engine's converter: one convertXxx method per grammar construct, dispatched
// from convert(). Anything not yet handled returns *ast.TODO, which Parse skips.
type cc struct {
	paramCount int
}

type node interface {
	GetParser() antlr.Parser
}

func todo(funcname string, n node) *ast.TODO {
	if debug.Active {
		log.Printf("oracle.%s: Unknown node type %T\n", funcname, n)
	}
	return &ast.TODO{}
}

// newString is a small constructor for ast.String used throughout the converter.
func newString(s string) *ast.String {
	return &ast.String{Str: s}
}

// identifier normalizes an Oracle identifier. Oracle folds unquoted identifiers
// to UPPER case, but sqlc's catalog comparisons are lower-case throughout, so we
// fold to lower case for unquoted identifiers and strip quotes for quoted ones.
func identifier(id string) string {
	if len(id) >= 2 && id[0] == '"' && id[len(id)-1] == '"' {
		return id[1 : len(id)-1]
	}
	return strings.ToLower(id)
}

// convert is the top-level dispatch for a unit_statement.
func (c *cc) convert(n *parser.Unit_statementContext) ast.Node {
	if n == nil {
		return &ast.TODO{}
	}
	if ct := n.Create_table(); ct != nil {
		if ctx, ok := ct.(*parser.Create_tableContext); ok {
			return c.convertCreateTable(ctx)
		}
	}
	if dml := n.Data_manipulation_language_statements(); dml != nil {
		if ctx, ok := dml.(*parser.Data_manipulation_language_statementsContext); ok {
			return c.convertDML(ctx)
		}
	}
	return todo("convert", n)
}

func (c *cc) convertDML(n *parser.Data_manipulation_language_statementsContext) ast.Node {
	if s := n.Select_statement(); s != nil {
		if ctx, ok := s.(*parser.Select_statementContext); ok {
			return c.convertSelectStatement(ctx)
		}
	}
	if s := n.Insert_statement(); s != nil {
		if ctx, ok := s.(*parser.Insert_statementContext); ok {
			return c.convertInsertStatement(ctx)
		}
	}
	if s := n.Update_statement(); s != nil {
		if ctx, ok := s.(*parser.Update_statementContext); ok {
			return c.convertUpdateStatement(ctx)
		}
	}
	if s := n.Delete_statement(); s != nil {
		if ctx, ok := s.(*parser.Delete_statementContext); ok {
			return c.convertDeleteStatement(ctx)
		}
	}
	return todo("convertDML", n)
}

// ---------------------------------------------------------------------------
// CREATE TABLE
// ---------------------------------------------------------------------------

func (c *cc) convertCreateTable(n *parser.Create_tableContext) ast.Node {
	tn := &ast.TableName{}
	if sc := n.Schema_name(); sc != nil {
		if ctx, ok := sc.(*parser.Schema_nameContext); ok {
			tn.Schema = identifier(ctx.GetText())
		}
	}
	if t := n.Table_name(); t != nil {
		tn.Name = identifier(t.GetText())
	}

	stmt := &ast.CreateTableStmt{
		Name:        tn,
		IfNotExists: n.EXISTS() != nil,
	}

	rt := n.Relational_table()
	if rt == nil {
		return stmt
	}
	rtc, ok := rt.(*parser.Relational_tableContext)
	if !ok {
		return stmt
	}
	for _, iprop := range rtc.AllRelational_property() {
		prop, ok := iprop.(*parser.Relational_propertyContext)
		if !ok {
			continue
		}
		coldef := prop.Column_definition()
		if coldef == nil {
			continue
		}
		cd, ok := coldef.(*parser.Column_definitionContext)
		if !ok {
			continue
		}
		col := c.convertColumnDefinition(cd)
		if col != nil {
			stmt.Cols = append(stmt.Cols, col)
		}
	}
	return stmt
}

func (c *cc) convertColumnDefinition(n *parser.Column_definitionContext) *ast.ColumnDef {
	name := ""
	if cn := n.Column_name(); cn != nil {
		name = identifier(cn.GetText())
	}
	typeName := "any"
	if dt := n.Datatype(); dt != nil {
		if ctx, ok := dt.(*parser.DatatypeContext); ok {
			typeName = normalizeDatatype(ctx)
		}
	} else if tn := n.Type_name(); tn != nil {
		typeName = identifier(tn.GetText())
	}
	return &ast.ColumnDef{
		Colname:   name,
		TypeName:  &ast.TypeName{Name: typeName},
		IsNotNull: hasNotNull(n),
	}
}

// normalizeDatatype extracts the base Oracle native type name (lower-cased),
// dropping precision/scale so the catalog can key on the bare type.
func normalizeDatatype(n *parser.DatatypeContext) string {
	if nde := n.Native_datatype_element(); nde != nil {
		return strings.ToLower(nde.GetText())
	}
	if n.INTERVAL() != nil {
		return "interval"
	}
	return strings.ToLower(n.GetText())
}

// hasNotNull reports whether any inline constraint on the column is NOT NULL.
func hasNotNull(n *parser.Column_definitionContext) bool {
	for _, ic := range n.AllInline_constraint() {
		if strings.Contains(strings.ToUpper(ic.GetText()), "NOTNULL") ||
			strings.Contains(strings.ToUpper(ic.GetText()), "NOT NULL") {
			return true
		}
		// Text has whitespace stripped by GetText(); "NOT NULL" becomes "NOTNULL".
	}
	return false
}

// ---------------------------------------------------------------------------
// SELECT
// ---------------------------------------------------------------------------

func (c *cc) convertSelectStatement(n *parser.Select_statementContext) ast.Node {
	only := n.Select_only_statement()
	if only == nil {
		return todo("convertSelectStatement", n)
	}
	oc, ok := only.(*parser.Select_only_statementContext)
	if !ok {
		return todo("convertSelectStatement", n)
	}
	sub := oc.Subquery()
	if sub == nil {
		return todo("convertSelectStatement", n)
	}
	sc, ok := sub.(*parser.SubqueryContext)
	if !ok {
		return todo("convertSelectStatement", n)
	}
	basic := sc.Subquery_basic_elements()
	if basic == nil {
		return todo("convertSelectStatement", n)
	}
	be, ok := basic.(*parser.Subquery_basic_elementsContext)
	if !ok {
		return todo("convertSelectStatement", n)
	}
	qb := be.Query_block()
	if qb == nil {
		return todo("convertSelectStatement", n)
	}
	qbc, ok := qb.(*parser.Query_blockContext)
	if !ok {
		return todo("convertSelectStatement", n)
	}
	return c.convertQueryBlock(qbc)
}

func (c *cc) convertQueryBlock(n *parser.Query_blockContext) ast.Node {
	stmt := &ast.SelectStmt{
		TargetList:     &ast.List{},
		FromClause:     &ast.List{},
		DistinctClause: &ast.List{},
	}

	if sl := n.Selected_list(); sl != nil {
		if slc, ok := sl.(*parser.Selected_listContext); ok {
			stmt.TargetList = c.convertSelectedList(slc)
		}
	}

	if fc := n.From_clause(); fc != nil {
		if fcc, ok := fc.(*parser.From_clauseContext); ok {
			stmt.FromClause = c.convertFromClause(fcc)
		}
	}

	if wc := n.Where_clause(); wc != nil {
		if wcc, ok := wc.(*parser.Where_clauseContext); ok {
			stmt.WhereClause = c.convertWhereClause(wcc)
		}
	}

	return stmt
}

func (c *cc) convertSelectedList(n *parser.Selected_listContext) *ast.List {
	list := &ast.List{}
	// SELECT *
	if n.ASTERISK() != nil {
		list.Items = append(list.Items, &ast.ResTarget{
			Val: &ast.ColumnRef{
				Fields: &ast.List{Items: []ast.Node{&ast.A_Star{}}},
			},
		})
		return list
	}
	for _, iel := range n.AllSelect_list_elements() {
		el, ok := iel.(*parser.Select_list_elementsContext)
		if !ok {
			continue
		}
		list.Items = append(list.Items, c.convertSelectListElement(el))
	}
	return list
}

func (c *cc) convertSelectListElement(n *parser.Select_list_elementsContext) ast.Node {
	// tableview_name.*
	if n.ASTERISK() != nil {
		cols := &ast.List{}
		if tv := n.Tableview_name(); tv != nil {
			cols.Items = append(cols.Items, newString(identifier(tv.GetText())))
		}
		cols.Items = append(cols.Items, &ast.A_Star{})
		return &ast.ResTarget{
			Val: &ast.ColumnRef{Fields: cols},
		}
	}

	rt := &ast.ResTarget{}
	if ex := n.Expression(); ex != nil {
		if exc, ok := ex.(*parser.ExpressionContext); ok {
			rt.Val = c.convertExpression(exc)
		}
	}
	if ca := n.Column_alias(); ca != nil {
		if cac, ok := ca.(*parser.Column_aliasContext); ok {
			name := columnAlias(cac)
			if name != "" {
				rt.Name = &name
			}
		}
	}
	return rt
}

func columnAlias(n *parser.Column_aliasContext) string {
	if id := n.Identifier(); id != nil {
		return identifier(id.GetText())
	}
	if qs := n.Quoted_string(); qs != nil {
		return strings.Trim(qs.GetText(), "'")
	}
	return ""
}

func (c *cc) convertFromClause(n *parser.From_clauseContext) *ast.List {
	list := &ast.List{}
	trl := n.Table_ref_list()
	if trl == nil {
		return list
	}
	trlc, ok := trl.(*parser.Table_ref_listContext)
	if !ok {
		return list
	}
	for _, itr := range trlc.AllTable_ref() {
		tr, ok := itr.(*parser.Table_refContext)
		if !ok {
			continue
		}
		if rv := c.convertTableRef(tr); rv != nil {
			list.Items = append(list.Items, rv)
		}
	}
	return list
}

// convertTableRef handles the simple `table [alias]` case; joins and subqueries
// fall back to a best-effort table name extraction.
func (c *cc) convertTableRef(n *parser.Table_refContext) ast.Node {
	aux := n.Table_ref_aux()
	if aux == nil {
		return nil
	}
	name, alias := tableRefAuxName(aux)
	if name == "" {
		return nil
	}
	rv := &ast.RangeVar{
		Relname: &name,
	}
	if alias != "" {
		rv.Alias = &ast.Alias{Aliasname: &alias}
	}
	return rv
}

func tableRefAuxName(aux parser.ITable_ref_auxContext) (string, string) {
	auxc, ok := aux.(*parser.Table_ref_auxContext)
	if !ok {
		return "", ""
	}
	alias := ""
	if ta := auxc.Table_alias(); ta != nil {
		alias = identifier(ta.GetText())
	}
	internal := auxc.Table_ref_aux_internal()
	if internal == nil {
		return "", alias
	}
	// The one/two/three alternatives all expose a Dml_table_expression_clause
	// through the labeled context; use text-based extraction as a fallback since
	// the labeled subtypes differ.
	name := extractTableViewName(internal)
	return name, alias
}

func extractTableViewName(n antlr.ParseTree) string {
	// Depth-first search for a Tableview_nameContext.
	switch v := n.(type) {
	case *parser.Tableview_nameContext:
		return identifier(v.GetText())
	}
	for i := 0; i < n.GetChildCount(); i++ {
		if child, ok := n.GetChild(i).(antlr.ParseTree); ok {
			if name := extractTableViewName(child); name != "" {
				return name
			}
		}
	}
	return ""
}

func (c *cc) convertWhereClause(n *parser.Where_clauseContext) ast.Node {
	cond := n.Condition()
	if cond == nil {
		return nil
	}
	condc, ok := cond.(*parser.ConditionContext)
	if !ok {
		return nil
	}
	if ex := condc.Expression(); ex != nil {
		if exc, ok := ex.(*parser.ExpressionContext); ok {
			return c.convertExpression(exc)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// INSERT
// ---------------------------------------------------------------------------

func (c *cc) convertInsertStatement(n *parser.Insert_statementContext) ast.Node {
	single := n.Single_table_insert()
	if single == nil {
		return todo("convertInsertStatement", n)
	}
	sc, ok := single.(*parser.Single_table_insertContext)
	if !ok {
		return todo("convertInsertStatement", n)
	}

	stmt := &ast.InsertStmt{
		Cols:          &ast.List{},
		ReturningList: &ast.List{},
	}

	if into := sc.Insert_into_clause(); into != nil {
		if intoc, ok := into.(*parser.Insert_into_clauseContext); ok {
			stmt.Relation = c.insertRelation(intoc)
			stmt.Cols = c.insertColumns(intoc)
		}
	}

	if vals := sc.Values_clause(); vals != nil {
		if valsc, ok := vals.(*parser.Values_clauseContext); ok {
			stmt.SelectStmt = c.convertValuesClause(valsc)
		}
	} else if sel := sc.Select_statement(); sel != nil {
		if selc, ok := sel.(*parser.Select_statementContext); ok {
			stmt.SelectStmt = c.convertSelectStatement(selc)
		}
	}

	return stmt
}

func (c *cc) insertRelation(n *parser.Insert_into_clauseContext) *ast.RangeVar {
	gtr := n.General_table_ref()
	if gtr == nil {
		return nil
	}
	name := extractTableViewName(gtr)
	if name == "" {
		return nil
	}
	return &ast.RangeVar{Relname: &name}
}

func (c *cc) insertColumns(n *parser.Insert_into_clauseContext) *ast.List {
	list := &ast.List{}
	pcl := n.Paren_column_list()
	if pcl == nil {
		return list
	}
	pclc, ok := pcl.(*parser.Paren_column_listContext)
	if !ok {
		return list
	}
	cl := pclc.Column_list()
	if cl == nil {
		return list
	}
	clc, ok := cl.(*parser.Column_listContext)
	if !ok {
		return list
	}
	for _, icn := range clc.AllColumn_name() {
		name := identifier(icn.GetText())
		list.Items = append(list.Items, &ast.ResTarget{Name: &name})
	}
	return list
}

func (c *cc) convertValuesClause(n *parser.Values_clauseContext) ast.Node {
	sel := &ast.SelectStmt{
		// TargetList/FromClause are initialized (empty) because downstream
		// analysis (compiler.paramSearch) dereferences SelectStmt.TargetList.Items
		// unconditionally for INSERT ... VALUES.
		TargetList:  &ast.List{},
		FromClause:  &ast.List{},
		ValuesLists: &ast.List{},
	}
	exprs := n.Expressions_()
	if exprs == nil {
		return sel
	}
	exc, ok := exprs.(*parser.Expressions_Context)
	if !ok {
		return sel
	}
	row := &ast.List{}
	for _, iex := range exc.AllExpression() {
		if e, ok := iex.(*parser.ExpressionContext); ok {
			row.Items = append(row.Items, c.convertExpression(e))
		}
	}
	sel.ValuesLists.Items = append(sel.ValuesLists.Items, row)
	return sel
}

// ---------------------------------------------------------------------------
// UPDATE
// ---------------------------------------------------------------------------

func (c *cc) convertUpdateStatement(n *parser.Update_statementContext) ast.Node {
	stmt := &ast.UpdateStmt{
		Relations:     &ast.List{},
		TargetList:    &ast.List{},
		ReturningList: &ast.List{},
		FromClause:    &ast.List{},
	}

	if gtr := n.General_table_ref(); gtr != nil {
		name := extractTableViewName(gtr)
		if name != "" {
			stmt.Relations.Items = append(stmt.Relations.Items, &ast.RangeVar{Relname: &name})
		}
	}

	if usc := n.Update_set_clause(); usc != nil {
		if uscc, ok := usc.(*parser.Update_set_clauseContext); ok {
			stmt.TargetList = c.convertUpdateSet(uscc)
		}
	}

	if wc := n.Where_clause(); wc != nil {
		if wcc, ok := wc.(*parser.Where_clauseContext); ok {
			stmt.WhereClause = c.convertWhereClause(wcc)
		}
	}

	return stmt
}

func (c *cc) convertUpdateSet(n *parser.Update_set_clauseContext) *ast.List {
	list := &ast.List{}
	for _, icb := range n.AllColumn_based_update_set_clause() {
		cb, ok := icb.(*parser.Column_based_update_set_clauseContext)
		if !ok {
			continue
		}
		cn := cb.Column_name()
		if cn == nil {
			continue
		}
		name := identifier(cn.GetText())
		rt := &ast.ResTarget{Name: &name}
		if ex := cb.Expression(); ex != nil {
			if exc, ok := ex.(*parser.ExpressionContext); ok {
				rt.Val = c.convertExpression(exc)
			}
		}
		list.Items = append(list.Items, rt)
	}
	return list
}

// ---------------------------------------------------------------------------
// DELETE
// ---------------------------------------------------------------------------

func (c *cc) convertDeleteStatement(n *parser.Delete_statementContext) ast.Node {
	stmt := &ast.DeleteStmt{
		Relations:     &ast.List{},
		ReturningList: &ast.List{},
	}
	if gtr := n.General_table_ref(); gtr != nil {
		name := extractTableViewName(gtr)
		if name != "" {
			stmt.Relations.Items = append(stmt.Relations.Items, &ast.RangeVar{Relname: &name})
		}
	}
	if wc := n.Where_clause(); wc != nil {
		if wcc, ok := wc.(*parser.Where_clauseContext); ok {
			stmt.WhereClause = c.convertWhereClause(wcc)
		}
	}
	return stmt
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

func (c *cc) convertExpression(n *parser.ExpressionContext) ast.Node {
	if le := n.Logical_expression(); le != nil {
		if lec, ok := le.(*parser.Logical_expressionContext); ok {
			return c.convertLogicalExpression(lec)
		}
	}
	return todo("convertExpression", n)
}

func (c *cc) convertLogicalExpression(n *parser.Logical_expressionContext) ast.Node {
	// AND / OR
	subs := n.AllLogical_expression()
	if len(subs) == 2 {
		var kind ast.BoolExprType
		if n.AND() != nil {
			kind = ast.BoolExprTypeAnd
		} else {
			kind = ast.BoolExprTypeOr
		}
		left := c.convertLogicalExpression(subs[0].(*parser.Logical_expressionContext))
		right := c.convertLogicalExpression(subs[1].(*parser.Logical_expressionContext))
		return &ast.BoolExpr{
			Boolop: kind,
			Args:   &ast.List{Items: []ast.Node{left, right}},
		}
	}

	if ule := n.Unary_logical_expression(); ule != nil {
		if ulec, ok := ule.(*parser.Unary_logical_expressionContext); ok {
			return c.convertUnaryLogicalExpression(ulec)
		}
	}
	return todo("convertLogicalExpression", n)
}

func (c *cc) convertUnaryLogicalExpression(n *parser.Unary_logical_expressionContext) ast.Node {
	me := n.Multiset_expression()
	if me == nil {
		return todo("convertUnaryLogicalExpression", n)
	}
	mec, ok := me.(*parser.Multiset_expressionContext)
	if !ok {
		return todo("convertUnaryLogicalExpression", n)
	}
	re := mec.Relational_expression()
	if re == nil {
		return todo("convertUnaryLogicalExpression", n)
	}
	rec, ok := re.(*parser.Relational_expressionContext)
	if !ok {
		return todo("convertUnaryLogicalExpression", n)
	}
	return c.convertRelationalExpression(rec)
}

func (c *cc) convertRelationalExpression(n *parser.Relational_expressionContext) ast.Node {
	// Binary comparison: relational_expression relational_operator relational_expression
	subs := n.AllRelational_expression()
	if len(subs) == 2 {
		op := "="
		if ro := n.Relational_operator(); ro != nil {
			op = relationalOperator(ro.(*parser.Relational_operatorContext))
		}
		left := c.convertRelationalExpression(subs[0].(*parser.Relational_expressionContext))
		right := c.convertRelationalExpression(subs[1].(*parser.Relational_expressionContext))
		return &ast.A_Expr{
			Kind:  ast.A_Expr_Kind_OP,
			Name:  &ast.List{Items: []ast.Node{newString(op)}},
			Lexpr: left,
			Rexpr: right,
		}
	}

	if ce := n.Compound_expression(); ce != nil {
		if cec, ok := ce.(*parser.Compound_expressionContext); ok {
			return c.convertCompoundExpression(cec)
		}
	}
	return todo("convertRelationalExpression", n)
}

func relationalOperator(n *parser.Relational_operatorContext) string {
	switch {
	case n.EQUALS_OP() != nil:
		return "="
	case n.NOT_EQUAL_OP() != nil:
		return "<>"
	case n.LESS_THAN_OP() != nil:
		return "<"
	case n.GREATER_THAN_OP() != nil:
		return ">"
	default:
		return strings.TrimSpace(n.GetText())
	}
}

func (c *cc) convertCompoundExpression(n *parser.Compound_expressionContext) ast.Node {
	// Only handle the plain concatenation path for now.
	concs := n.AllConcatenation()
	if len(concs) >= 1 {
		if cc0, ok := concs[0].(*parser.ConcatenationContext); ok {
			return c.convertConcatenation(cc0)
		}
	}
	return todo("convertCompoundExpression", n)
}

func (c *cc) convertConcatenation(n *parser.ConcatenationContext) ast.Node {
	if me := n.Model_expression(); me != nil {
		if mec, ok := me.(*parser.Model_expressionContext); ok {
			return c.convertModelExpression(mec)
		}
	}
	return todo("convertConcatenation", n)
}

func (c *cc) convertModelExpression(n *parser.Model_expressionContext) ast.Node {
	if ue := n.Unary_expression(); ue != nil {
		if uec, ok := ue.(*parser.Unary_expressionContext); ok {
			return c.convertUnaryExpression(uec)
		}
	}
	return todo("convertModelExpression", n)
}

func (c *cc) convertUnaryExpression(n *parser.Unary_expressionContext) ast.Node {
	if core := n.Unary_expression_core(); core != nil {
		if corec, ok := core.(*parser.Unary_expression_coreContext); ok {
			return c.convertUnaryExpressionCore(corec)
		}
	}
	return todo("convertUnaryExpression", n)
}

func (c *cc) convertUnaryExpressionCore(n *parser.Unary_expression_coreContext) ast.Node {
	if a := n.Atom(); a != nil {
		if ac, ok := a.(*parser.AtomContext); ok {
			return c.convertAtom(ac)
		}
	}
	return todo("convertUnaryExpressionCore", n)
}

func (c *cc) convertAtom(n *parser.AtomContext) ast.Node {
	if bv := n.Bind_variable(); bv != nil {
		if bvc, ok := bv.(*parser.Bind_variableContext); ok {
			return c.convertBindVariable(bvc)
		}
	}
	if cst := n.Constant(); cst != nil {
		if cstc, ok := cst.(*parser.ConstantContext); ok {
			return c.convertConstant(cstc)
		}
	}
	if ge := n.General_element(); ge != nil {
		if gec, ok := ge.(*parser.General_elementContext); ok {
			return c.convertGeneralElement(gec)
		}
	}
	return todo("convertAtom", n)
}

// convertBindVariable maps Oracle bind variables to sqlc ParamRefs.
//
//	:1, :2 ...   -> positional, ParamRef.Number = N
//	:name        -> named; represented like SQLite's @name via A_Expr so the
//	                downstream named-parameter machinery can pick it up.
func (c *cc) convertBindVariable(n *parser.Bind_variableContext) ast.Node {
	loc := n.GetStart().GetStart()

	// Numbered bind spelled as ':' UNSIGNED_INTEGER (two tokens).
	if ui := n.UNSIGNED_INTEGER(0); ui != nil {
		number, _ := strconv.Atoi(ui.GetText())
		if number == 0 {
			c.paramCount++
			number = c.paramCount
		}
		return &ast.ParamRef{
			Number:   number,
			Location: loc,
			Dollar:   true,
		}
	}

	// BINDVAR is a single token, either ":name" or ":1". The Oracle lexer emits
	// ":1" as one BINDVAR token, so distinguish numbered from named here.
	if bindvar := n.BINDVAR(0); bindvar != nil {
		text := bindvar.GetText()
		name := strings.TrimPrefix(text, ":")
		if number, err := strconv.Atoi(name); err == nil {
			return &ast.ParamRef{
				Number:   number,
				Location: loc,
				Dollar:   true,
			}
		}
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{newString("@")}},
			Rexpr:    newString(name),
			Location: loc,
		}
	}

	c.paramCount++
	return &ast.ParamRef{Number: c.paramCount, Location: loc, Dollar: true}
}

func (c *cc) convertConstant(n *parser.ConstantContext) ast.Node {
	loc := n.GetStart().GetStart()

	if num := n.Numeric(); num != nil {
		text := num.GetText()
		if iv, err := strconv.ParseInt(text, 10, 64); err == nil {
			return &ast.A_Const{Val: &ast.Integer{Ival: iv}, Location: loc}
		}
		return &ast.A_Const{Val: &ast.Float{Str: text}, Location: loc}
	}
	if qs := n.Quoted_string(0); qs != nil {
		return &ast.A_Const{Val: newString(strings.Trim(qs.GetText(), "'")), Location: loc}
	}
	if n.NULL_() != nil {
		return &ast.A_Const{Val: &ast.Null{}, Location: loc}
	}
	return &ast.A_Const{Val: newString(n.GetText()), Location: loc}
}

// convertGeneralElement maps `a`, `a.b`, `t.col` into a ColumnRef.
func (c *cc) convertGeneralElement(n *parser.General_elementContext) ast.Node {
	fields := &ast.List{}
	collectGeneralElementParts(n, fields)
	if len(fields.Items) == 0 {
		return todo("convertGeneralElement", n)
	}
	return &ast.ColumnRef{
		Fields:   fields,
		Location: n.GetStart().GetStart(),
	}
}

func collectGeneralElementParts(n *parser.General_elementContext, out *ast.List) {
	if inner := n.General_element(); inner != nil {
		if innerc, ok := inner.(*parser.General_elementContext); ok {
			collectGeneralElementParts(innerc, out)
		}
	}
	for _, ipart := range n.AllGeneral_element_part() {
		pc, ok := ipart.(*parser.General_element_partContext)
		if !ok {
			continue
		}
		if id := pc.Id_expression(); id != nil {
			out.Items = append(out.Items, newString(identifier(id.GetText())))
		}
	}
}
