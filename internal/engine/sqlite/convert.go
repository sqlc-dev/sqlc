package sqlite

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	meyer "github.com/sqlc-dev/meyer/ast"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// cc converts a meyer syntax tree into sqlc's engine-independent AST. One
// converter is used per statement; paramCount numbers the anonymous "?"
// markers within it.
type cc struct {
	paramCount int
}

func todo(funcname string, n meyer.Node) *ast.TODO {
	if debug.Active {
		log.Printf("sqlite.%s: Unknown node type %T\n", funcname, n)
	}
	return &ast.TODO{}
}

func (c *cc) convert(node meyer.Node) ast.Node {
	switch n := node.(type) {

	// Statements

	case *meyer.AlterTableStmt:
		return c.convertAlterTableStmt(n)

	case *meyer.AttachStmt:
		return c.convertAttachStmt(n)

	case *meyer.CreateTableStmt:
		return c.convertCreateTableStmt(n)

	case *meyer.CreateViewStmt:
		return c.convertCreateViewStmt(n)

	case *meyer.CreateVirtualTableStmt:
		return c.convertCreateVirtualTableStmt(n)

	case *meyer.DeleteStmt:
		return c.convertDeleteStmt(n)

	case *meyer.DropTableStmt:
		return c.convertDropStmt(n.Name, n.IfExists)

	case *meyer.DropViewStmt:
		return c.convertDropStmt(n.Name, n.IfExists)

	case *meyer.ExplainStmt:
		return c.convert(n.Stmt)

	case *meyer.InsertStmt:
		return c.convertInsertStmt(n)

	case *meyer.SelectStmt:
		return c.convertSelectStmt(n)

	case *meyer.UpdateStmt:
		return c.convertUpdateStmt(n)

	// Expressions

	case *meyer.BetweenExpr:
		return c.convertBetweenExpr(n)

	case *meyer.BinaryExpr:
		return c.convertBinaryExpr(n)

	case *meyer.BindParam:
		return c.convertBindParam(n)

	case *meyer.CaseExpr:
		return c.convertCaseExpr(n)

	case *meyer.CastExpr:
		return c.convertCastExpr(n)

	case *meyer.CollateExpr:
		return c.convertCollateExpr(n)

	case *meyer.ExistsExpr:
		return &ast.SubLink{
			SubLinkType: ast.EXISTS_SUBLINK,
			Subselect:   c.convert(n.Select),
			Location:    n.Pos(),
		}

	case *meyer.FuncCall:
		return c.convertFuncCall(n)

	case *meyer.Ident:
		return c.columnRef(n.Pos(), n)

	case *meyer.InExpr:
		return c.convertInExpr(n)

	case *meyer.IsExpr:
		return c.convertIsExpr(n)

	case *meyer.LikeExpr:
		return c.convertLikeExpr(n)

	case *meyer.Literal:
		return c.convertLiteral(n)

	case *meyer.NullCheckExpr:
		return c.convertNullCheckExpr(n)

	case *meyer.ParenExpr:
		return c.convert(n.X)

	case *meyer.QualifiedRef:
		return c.columnRef(n.Pos(), n.Parts...)

	case *meyer.Star:
		return c.convertStar(n)

	case *meyer.SubqueryExpr:
		return &ast.SubLink{
			SubLinkType: ast.EXPR_SUBLINK,
			Subselect:   c.convert(n.Select),
			Location:    n.Pos(),
		}

	case *meyer.UnaryExpr:
		return c.convertUnaryExpr(n)

	case *meyer.VectorExpr:
		return c.convertExprList(n.List)

	default:
		return todo("convert(case=default)", node)
	}
}

// convertExpr converts an optional expression, mapping an absent one to a
// nil node rather than to a TODO.
func (c *cc) convertExpr(e meyer.Expr) ast.Node {
	if e == nil {
		return nil
	}
	return c.convert(e)
}

func (c *cc) convertExprList(exprs []meyer.Expr) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, e := range exprs {
		list.Items = append(list.Items, c.convert(e))
	}
	return list
}

// columnRef builds the dotted reference sqlc uses for a column: a one, two
// or three part name of the form [[schema.]table.]column.
func (c *cc) columnRef(loc int, parts ...*meyer.Ident) *ast.ColumnRef {
	items := make([]ast.Node, 0, len(parts))
	for _, p := range parts {
		items = append(items, NewIdentifier(p))
	}
	return &ast.ColumnRef{
		Fields:   &ast.List{Items: items},
		Location: loc,
	}
}

func (c *cc) convertStar(n *meyer.Star) *ast.ColumnRef {
	items := []ast.Node{}
	if n.Table != nil {
		items = append(items, NewIdentifier(n.Table))
	}
	items = append(items, &ast.A_Star{})
	return &ast.ColumnRef{
		Fields:   &ast.List{Items: items},
		Location: n.Pos(),
	}
}

//
// Statements
//

func (c *cc) convertAlterTableStmt(n *meyer.AlterTableStmt) ast.Node {
	switch n.Action {
	case meyer.AlterRenameTable:
		name := identifier(n.NewName)
		return &ast.RenameTableStmt{
			Table:   parseTableName(n.Table),
			NewName: &name,
		}

	case meyer.AlterRenameColumn:
		name := identifier(n.NewName)
		return &ast.RenameColumnStmt{
			Table: parseTableName(n.Table),
			Col: &ast.ColumnRef{
				Name: identifier(n.Column),
			},
			NewName: &name,
		}

	case meyer.AlterAddColumn:
		def := n.ColumnDef
		if def == nil {
			break
		}
		name := identifier(def.Name)
		return &ast.AlterTableStmt{
			Table: parseTableName(n.Table),
			Cmds: &ast.List{Items: []ast.Node{
				&ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def: &ast.ColumnDef{
						Colname:   name,
						TypeName:  &ast.TypeName{Name: columnTypeName(def.Type)},
						IsNotNull: hasNotNullConstraint(def.Constraints),
					},
				},
			}},
		}

	case meyer.AlterDropColumn:
		name := identifier(n.Column)
		return &ast.AlterTableStmt{
			Table: parseTableName(n.Table),
			Cmds: &ast.List{Items: []ast.Node{
				&ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_DropColumn,
				},
			}},
		}
	}
	return todo("convertAlterTableStmt", n)
}

// convertAttachStmt records an attached database as a schema. Only the name
// matters to sqlc; the file expression is a runtime concern.
func (c *cc) convertAttachStmt(n *meyer.AttachStmt) ast.Node {
	name := schemaName(n.Name)
	return &ast.CreateSchemaStmt{
		Name: &name,
	}
}

// schemaName reads the schema out of the "AS <name>" operand of ATTACH,
// which the grammar accepts as an arbitrary expression.
func schemaName(e meyer.Expr) string {
	switch n := e.(type) {
	case *meyer.Ident:
		return identifier(n)
	case *meyer.Literal:
		return n.Value
	default:
		return ""
	}
}

// columnTypeName returns the declared type of a column, defaulting to
// SQLite's untyped "any" when the declaration omits one.
func columnTypeName(t *meyer.TypeName) string {
	if t == nil {
		return "any"
	}
	return typeName(t)
}

func (c *cc) convertCreateTableStmt(n *meyer.CreateTableStmt) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n.Name),
		IfNotExists: n.IfNotExists,
	}
	for _, def := range n.Columns {
		stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
			Colname:   identifier(def.Name),
			IsNotNull: hasNotNullConstraint(def.Constraints),
			TypeName:  &ast.TypeName{Name: columnTypeName(def.Type)},
		})
	}
	return stmt
}

func (c *cc) convertCreateVirtualTableStmt(n *meyer.CreateVirtualTableStmt) ast.Node {
	switch moduleName := identifier(n.Module); moduleName {
	case "fts5":
		// https://www.sqlite.org/fts5.html
		return c.convertCreateVirtualTableFTS5(n)
	default:
		return todo(
			fmt.Sprintf("create_virtual_table. unsupported module name: %q", moduleName),
			n,
		)
	}
}

func (c *cc) convertCreateVirtualTableFTS5(n *meyer.CreateVirtualTableStmt) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n.Name),
		IfNotExists: n.IfNotExists,
	}

	// The module arguments of a virtual table are an arbitrary token
	// sequence, so meyer hands them back as source text. For fts5 an
	// argument is either a column ("b", or "c UNINDEXED") or one of the
	// module's own options, which are always written as "key = value".
	for _, arg := range n.Args {
		name := fts5ColumnName(arg)
		if name == "" {
			continue
		}
		stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
			Colname: name,
			// fts5 accepts no column constraints, so we supply them here
			IsNotNull: true,
			TypeName:  &ast.TypeName{Name: "text"},
		})
	}

	return stmt
}

// fts5ColumnName returns the column an fts5 module argument declares, or the
// empty string when the argument configures the module instead.
func fts5ColumnName(arg string) string {
	if strings.ContainsRune(arg, '=') {
		return ""
	}
	name, _, _ := strings.Cut(strings.TrimSpace(arg), " ")
	return unquoteIdent(name)
}

// unquoteIdent folds a name that has not been through the lexer, as happens
// for the raw module arguments of a virtual table.
func unquoteIdent(s string) string {
	if len(s) >= 2 {
		switch {
		case s[0] == '"' && s[len(s)-1] == '"':
			return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
		case s[0] == '`' && s[len(s)-1] == '`':
			return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
		case s[0] == '[' && s[len(s)-1] == ']':
			return s[1 : len(s)-1]
		case s[0] == '\'' && s[len(s)-1] == '\'':
			return strings.ReplaceAll(s[1:len(s)-1], `''`, `'`)
		}
	}
	return strings.ToLower(s)
}

func (c *cc) convertCreateViewStmt(n *meyer.CreateViewStmt) ast.Node {
	return &ast.ViewStmt{
		View:            parseRangeVar(n.Name, nil),
		Aliases:         &ast.List{},
		Query:           c.convert(n.Select),
		Replace:         false,
		Options:         &ast.List{},
		WithCheckOption: ast.ViewCheckOption(0),
	}
}

func (c *cc) convertDropStmt(name *meyer.QualifiedName, ifExists bool) ast.Node {
	return &ast.DropTableStmt{
		IfExists: ifExists,
		Tables:   []*ast.TableName{parseTableName(name)},
	}
}

func (c *cc) convertDeleteStmt(n *meyer.DeleteStmt) ast.Node {
	return &ast.DeleteStmt{
		Relations:     &ast.List{Items: []ast.Node{parseRangeVar(n.Table, n.Alias)}},
		WhereClause:   c.convertExpr(n.Where),
		LimitCount:    c.limitCount(n.Limit),
		ReturningList: c.convertReturning(n.Returning),
		WithClause:    c.convertWith(n.With),
	}
}

// limitCount converts the LIMIT of an UPDATE or DELETE. sqlc has no place to
// put the OFFSET or the ORDER BY that can accompany it, so they are dropped.
func (c *cc) limitCount(n *meyer.Limit) ast.Node {
	if n == nil {
		return nil
	}
	return c.convertExpr(n.Count)
}

func (c *cc) convertInsertStmt(n *meyer.InsertStmt) ast.Node {
	insert := &ast.InsertStmt{
		Relation:      parseRangeVar(n.Table, n.Alias),
		Cols:          c.convertColumnNames(n.Columns),
		ReturningList: c.convertReturning(n.Returning),
		WithClause:    c.convertWith(n.With),
		DefaultValues: n.DefaultValues,
	}
	if n.Select != nil {
		if sel, ok := c.convertSelectStmt(n.Select).(*ast.SelectStmt); ok {
			insert.SelectStmt = sel
		}
	}
	if len(n.Upserts) > 0 {
		insert.OnConflictClause = c.convertUpsert(n.Upserts[0])
	}
	return insert
}

func (c *cc) convertUpsert(n *meyer.Upsert) *ast.OnConflictClause {
	clause := &ast.OnConflictClause{
		Action:   ast.OnConflictActionNothing,
		Location: n.Pos(),
	}
	if len(n.Target) > 0 || n.TargetWhere != nil {
		clause.Infer = &ast.InferClause{
			IndexElems:  &ast.List{Items: c.convertOrderBy(n.Target)},
			WhereClause: c.convertExpr(n.TargetWhere),
		}
	}
	if n.DoNothing {
		return clause
	}
	clause.Action = ast.OnConflictActionUpdate
	clause.TargetList = c.convertSetPairs(n.Set)
	clause.WhereClause = c.convertExpr(n.Where)
	return clause
}

func (c *cc) convertColumnNames(cols []*meyer.Ident) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, col := range cols {
		name := identifier(col)
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
		})
	}
	return list
}

// convertSetPairs converts the assignment list of an UPDATE or of an upsert's
// DO UPDATE. SQLite allows "(a, b) = expr"; sqlc has no node for a
// multi-column assignment, so each name is given the same value expression.
func (c *cc) convertSetPairs(pairs []*meyer.SetPair) *ast.List {
	list := &ast.List{}
	for _, pair := range pairs {
		val := c.convertExpr(pair.Value)
		for _, col := range pair.Columns {
			name := identifier(col)
			list.Items = append(list.Items, &ast.ResTarget{
				Name:     &name,
				Val:      val,
				Location: col.Pos(),
			})
		}
	}
	return list
}

func (c *cc) convertUpdateStmt(n *meyer.UpdateStmt) ast.Node {
	return &ast.UpdateStmt{
		Relations:     &ast.List{Items: []ast.Node{parseRangeVar(n.Table, n.Alias)}},
		TargetList:    c.convertSetPairs(n.Set),
		FromClause:    &ast.List{Items: c.convertFrom(n.From)},
		WhereClause:   c.convertExpr(n.Where),
		LimitCount:    c.limitCount(n.Limit),
		ReturningList: c.convertReturning(n.Returning),
		WithClause:    c.convertWith(n.With),
	}
}

func (c *cc) convertReturning(cols []*meyer.ResultColumn) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, col := range cols {
		list.Items = append(list.Items, &ast.ResTarget{
			Indirection: &ast.List{},
			Val:         c.convertExpr(col.Expr),
			Location:    col.Pos(),
		})
	}
	return list
}

//
// SELECT
//

func (c *cc) convertSelectStmt(n *meyer.SelectStmt) ast.Node {
	var stmt *ast.SelectStmt
	for i, core := range n.Cores {
		sel := c.convertSelectCore(core)
		if stmt == nil {
			stmt = sel
			continue
		}
		op, all := ast.None, false
		if i-1 < len(n.Ops) {
			switch n.Ops[i-1].Op {
			case "UNION":
				op = ast.Union
				all = n.Ops[i-1].All
			case "INTERSECT":
				op = ast.Intersect
			case "EXCEPT":
				op = ast.Except
			}
		}
		stmt = &ast.SelectStmt{
			TargetList: &ast.List{},
			FromClause: &ast.List{},
			Op:         op,
			All:        all,
			Larg:       stmt,
			Rarg:       sel,
		}
		// ORDER BY and LIMIT are written on the last core of a compound
		// select but apply to the compound as a whole.
		hoistTail(stmt, stmt.Rarg)
	}
	if stmt == nil {
		return todo("convertSelectStmt", n)
	}
	stmt.WithClause = c.convertWith(n.With)
	return stmt
}

// hoistTail moves the sort and limit clauses of a compound select's final
// core up onto the compound itself.
func hoistTail(outer, inner *ast.SelectStmt) {
	outer.SortClause, inner.SortClause = inner.SortClause, nil
	outer.LimitCount, inner.LimitCount = inner.LimitCount, nil
	outer.LimitOffset, inner.LimitOffset = inner.LimitOffset, nil
}

func (c *cc) convertSelectCore(core meyer.SelectCore) *ast.SelectStmt {
	switch n := core.(type) {
	case *meyer.SelectQuery:
		stmt := &ast.SelectStmt{
			TargetList:   &ast.List{Items: c.convertResultColumns(n.Columns)},
			FromClause:   &ast.List{Items: c.convertFrom(n.From)},
			WhereClause:  c.convertExpr(n.Where),
			GroupClause:  &ast.List{},
			HavingClause: c.convertExpr(n.Having),
			WindowClause: &ast.List{},
			ValuesLists:  &ast.List{},
		}
		for _, g := range n.GroupBy {
			stmt.GroupClause.Items = append(stmt.GroupClause.Items, c.convert(g))
		}
		if n.Distinct == meyer.DistinctDistinct {
			// SQLite has no DISTINCT ON, so the clause carries no columns.
			stmt.DistinctClause = &ast.List{Items: []ast.Node{&ast.TODO{}}}
		}
		for _, w := range n.Windows {
			stmt.WindowClause.Items = append(stmt.WindowClause.Items, c.convertWindowDef(w))
		}
		if len(n.OrderBy) > 0 {
			stmt.SortClause = &ast.List{Items: c.convertOrderBy(n.OrderBy)}
		}
		if n.Limit != nil {
			stmt.LimitCount = c.convertExpr(n.Limit.Count)
			stmt.LimitOffset = c.convertExpr(n.Limit.Offset)
		}
		return stmt

	case *meyer.ValuesClause:
		values := &ast.List{}
		for _, row := range n.Rows {
			values.Items = append(values.Items, c.convertExprList(row))
		}
		return &ast.SelectStmt{
			TargetList:  &ast.List{},
			FromClause:  &ast.List{},
			ValuesLists: values,
		}

	default:
		return &ast.SelectStmt{
			TargetList:  &ast.List{},
			FromClause:  &ast.List{},
			ValuesLists: &ast.List{},
		}
	}
}

func (c *cc) convertResultColumns(cols []*meyer.ResultColumn) []ast.Node {
	var out []ast.Node
	for _, col := range cols {
		val := c.convertExpr(col.Expr)
		if val == nil {
			continue
		}
		target := &ast.ResTarget{
			Val:      val,
			Location: col.Pos(),
		}
		if col.Alias != nil {
			name := identifier(col.Alias)
			target.Name = &name
		}
		out = append(out, target)
	}
	return out
}

func (c *cc) convertOrderBy(terms []*meyer.OrderingTerm) []ast.Node {
	var out []ast.Node
	for _, term := range terms {
		sortBy := &ast.SortBy{
			Node:     c.convertExpr(term.Expr),
			UseOp:    &ast.List{},
			Location: term.Pos(),
		}
		switch term.Order {
		case meyer.SortAsc:
			sortBy.SortbyDir = ast.SortByDirAsc
		case meyer.SortDesc:
			sortBy.SortbyDir = ast.SortByDirDesc
		default:
			sortBy.SortbyDir = ast.SortByDirDefault
		}
		switch term.Nulls {
		case meyer.NullsFirst:
			sortBy.SortbyNulls = ast.SortByNullsFirst
		case meyer.NullsLast:
			sortBy.SortbyNulls = ast.SortByNullsLast
		default:
			sortBy.SortbyNulls = ast.SortByNullsDefault
		}
		out = append(out, sortBy)
	}
	return out
}

func (c *cc) convertWindowDef(n *meyer.WindowDef) ast.Node {
	def := &ast.WindowDef{
		PartitionClause: c.convertExprList(n.Partition),
		OrderClause:     &ast.List{Items: c.convertOrderBy(n.OrderBy)},
		Location:        n.Pos(),
	}
	if n.Name != nil {
		name := identifier(n.Name)
		def.Name = &name
	}
	if n.Base != nil {
		base := identifier(n.Base)
		def.Refname = &base
	}
	return def
}

func (c *cc) convertWith(n *meyer.With) *ast.WithClause {
	if n == nil || len(n.CTEs) == 0 {
		return nil
	}
	ctes := &ast.List{}
	for _, cte := range n.CTEs {
		name := identifier(cte.Name)
		cols := &ast.List{}
		for _, col := range cte.Columns {
			cols.Items = append(cols.Items, NewIdentifier(col))
		}
		ctes.Items = append(ctes.Items, &ast.CommonTableExpr{
			Ctename:      &name,
			Ctequery:     c.convert(cte.Select),
			Location:     cte.Pos(),
			Cterecursive: n.Recursive,
			Ctecolnames:  cols,
		})
	}
	return &ast.WithClause{
		Ctes:      ctes,
		Recursive: n.Recursive,
		Location:  n.Pos(),
	}
}

// convertFrom converts a SQLite source list. The list is flat, with each
// item recording how it attaches to the one before it, so a chain of joins
// is rebuilt here as a left-leaning tree. A list joined only by commas has
// no join conditions to carry and stays flat.
func (c *cc) convertFrom(refs []*meyer.TableRef) []ast.Node {
	if len(refs) == 0 {
		return nil
	}
	if commaOnly(refs) {
		out := make([]ast.Node, 0, len(refs))
		for _, ref := range refs {
			out = append(out, c.convertTableRef(ref))
		}
		return out
	}
	table := c.convertTableRef(refs[0])
	for _, ref := range refs[1:] {
		join := &ast.JoinExpr{
			Larg: table,
			Rarg: c.convertTableRef(ref),
		}
		if op := ref.Join; op != nil {
			join.IsNatural = op.Type&meyer.JoinNatural != 0
			left := op.Type&meyer.JoinLeft != 0
			right := op.Type&meyer.JoinRight != 0
			switch {
			case left && right:
				join.Jointype = ast.JoinTypeFull
			case left:
				join.Jointype = ast.JoinTypeLeft
			case right:
				join.Jointype = ast.JoinTypeRight
			case op.Type&meyer.JoinComma == 0:
				join.Jointype = ast.JoinTypeInner
			}
		}
		if ref.On != nil {
			join.Quals = c.convert(ref.On)
		}
		if len(ref.Using) > 0 {
			using := &ast.List{}
			for _, col := range ref.Using {
				using.Items = append(using.Items, NewIdentifier(col))
			}
			join.UsingClause = using
		}
		table = join
	}
	return []ast.Node{table}
}

func commaOnly(refs []*meyer.TableRef) bool {
	for _, ref := range refs[1:] {
		if ref.Join == nil || ref.Join.Type&meyer.JoinComma == 0 {
			return false
		}
	}
	return true
}

func (c *cc) convertTableRef(n *meyer.TableRef) ast.Node {
	switch {
	case n.Name != nil && n.HasArgs:
		// A table-valued function, such as generate_series(1, 10).
		fn := &ast.FuncCall{
			Func:     &ast.FuncName{Name: identifier(n.Name.Name)},
			Funcname: &ast.List{Items: []ast.Node{NewIdentifier(n.Name.Name)}},
			Args:     c.convertExprList(n.Args),
			Location: n.Pos(),
		}
		if n.Name.Schema != nil {
			fn.Func.Schema = identifier(n.Name.Schema)
		}
		rf := &ast.RangeFunction{
			Functions: &ast.List{Items: []ast.Node{fn}},
		}
		if n.Alias != nil {
			alias := identifier(n.Alias)
			rf.Alias = &ast.Alias{Aliasname: &alias}
		}
		return rf

	case n.Name != nil:
		rv := parseRangeVar(n.Name, n.Alias)
		rv.Location = n.Pos()
		return rv

	case n.Select != nil:
		rs := &ast.RangeSubselect{
			Subquery: c.convert(n.Select),
		}
		if n.Alias != nil {
			alias := identifier(n.Alias)
			rs.Alias = &ast.Alias{Aliasname: &alias}
		}
		return rs

	case len(n.List) > 0:
		// A parenthesized source list, as in "FROM (a JOIN b ON ...)".
		nested := c.convertFrom(n.List)
		if len(nested) == 1 {
			return nested[0]
		}
		return &ast.List{Items: nested}
	}
	return todo("convertTableRef", n)
}

//
// Expressions
//

func (c *cc) convertLiteral(n *meyer.Literal) ast.Node {
	switch n.Kind {
	case meyer.LitInteger:
		if i, ok := parseInteger(n.Value); ok {
			return &ast.A_Const{
				Val:      &ast.Integer{Ival: i},
				Location: n.Pos(),
			}
		}
		// Too wide for an int64, or written with the digit separators
		// SQLite accepts. Either way the value is still numeric.
		return &ast.A_Const{
			Val:      &ast.Float{Str: n.Value},
			Location: n.Pos(),
		}

	case meyer.LitFloat:
		return &ast.A_Const{
			Val:      &ast.Float{Str: n.Value},
			Location: n.Pos(),
		}

	case meyer.LitString, meyer.LitBlob:
		return &ast.A_Const{
			Val:      &ast.String{Str: n.Value},
			Location: n.Pos(),
		}

	case meyer.LitNull:
		return &ast.A_Const{
			Val:      &ast.Null{},
			Location: n.Pos(),
		}
	}
	// CURRENT_DATE, CURRENT_TIME and CURRENT_TIMESTAMP are keywords rather
	// than function calls, but behave like niladic functions.
	name := strings.ToLower(n.Raw)
	return &ast.FuncCall{
		Func:     &ast.FuncName{Name: name},
		Funcname: &ast.List{Items: []ast.Node{&ast.String{Str: name}}},
		Args:     &ast.List{},
		AggOrder: &ast.List{},
		Location: n.Pos(),
	}
}

// parseInteger reads an integer literal in either of the two forms SQLite
// writes: decimal, where a leading zero is not an octal prefix, and
// hexadecimal.
func parseInteger(s string) (int64, bool) {
	if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		i, err := strconv.ParseInt(s[2:], 16, 64)
		return i, err == nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	return i, err == nil
}

func (c *cc) convertBindParam(n *meyer.BindParam) ast.Node {
	switch n.Kind {
	case meyer.ParamAnon, meyer.ParamNumber:
		// Parameter numbers start at one.
		c.paramCount += 1
		number := c.paramCount
		if n.Kind == meyer.ParamNumber {
			number = n.Number
		}
		return &ast.ParamRef{
			Number:   number,
			Location: n.Pos(),
			Dollar:   n.Kind == meyer.ParamNumber,
		}
	default:
		// A named parameter. sqlc carries the name through as the operand of
		// a pseudo-operator and resolves it later.
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "@"}}},
			Rexpr:    &ast.String{Str: n.Name},
			Location: n.Pos(),
		}
	}
}

func (c *cc) convertBinaryExpr(n *meyer.BinaryExpr) ast.Node {
	switch n.Op {
	case meyer.OpAnd, meyer.OpOr:
		op := ast.BoolExprTypeAnd
		if n.Op == meyer.OpOr {
			op = ast.BoolExprTypeOr
		}
		return &ast.BoolExpr{
			Boolop: op,
			Args: &ast.List{Items: []ast.Node{
				c.convert(n.X),
				c.convert(n.Y),
			}},
			Location: n.Pos(),
		}
	}
	return &ast.A_Expr{
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: n.Op.String()}}},
		Lexpr:    c.convert(n.X),
		Rexpr:    c.convert(n.Y),
		Location: n.Pos(),
	}
}

func (c *cc) convertUnaryExpr(n *meyer.UnaryExpr) ast.Node {
	expr := c.convert(n.X)
	switch n.Op {
	case meyer.OpNot:
		return &ast.BoolExpr{
			Boolop:   ast.BoolExprTypeNot,
			Args:     &ast.List{Items: []ast.Node{expr}},
			Location: n.Pos(),
		}
	case meyer.OpPlus:
		// A unary plus is a no-op on every operand SQLite accepts.
		return expr
	default:
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: n.Op.String()}}},
			Rexpr:    expr,
			Location: n.Pos(),
		}
	}
}

func (c *cc) convertIsExpr(n *meyer.IsExpr) ast.Node {
	expr := &ast.A_Expr{
		Lexpr:    c.convert(n.X),
		Rexpr:    c.convert(n.Y),
		Location: n.Pos(),
	}
	var name string
	switch {
	case n.Distinct && n.Not:
		expr.Kind = ast.A_Expr_Kind_NOT_DISTINCT
		name = "IS NOT DISTINCT FROM"
	case n.Distinct:
		expr.Kind = ast.A_Expr_Kind_DISTINCT
		name = "IS DISTINCT FROM"
	case n.Not:
		name = "IS NOT"
	default:
		name = "IS"
	}
	expr.Name = &ast.List{Items: []ast.Node{&ast.String{Str: name}}}
	return expr
}

func (c *cc) convertNullCheckExpr(n *meyer.NullCheckExpr) ast.Node {
	name := "NOT NULL"
	switch n.Test {
	case meyer.TestIsNull:
		name = "ISNULL"
	case meyer.TestNotNull:
		name = "NOTNULL"
	}
	return &ast.A_Expr{
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: name}}},
		Lexpr:    c.convert(n.X),
		Location: n.Pos(),
	}
}

func (c *cc) convertLikeExpr(n *meyer.LikeExpr) ast.Node {
	name := strings.ToUpper(n.Op.Name)
	if n.Not {
		name = "NOT " + name
	}
	return &ast.A_Expr{
		Name:     &ast.List{Items: []ast.Node{&ast.String{Str: name}}},
		Lexpr:    c.convert(n.X),
		Rexpr:    c.convert(n.Y),
		Location: n.Pos(),
	}
}

func (c *cc) convertInExpr(n *meyer.InExpr) ast.Node {
	if n.Select != nil {
		sublink := &ast.SubLink{
			SubLinkType: ast.ANY_SUBLINK,
			Testexpr:    c.convert(n.X),
			Subselect:   c.convertSelectStmt(n.Select),
			Location:    n.Pos(),
		}
		if !n.Not {
			return sublink
		}
		return &ast.BoolExpr{
			Boolop:   ast.BoolExprTypeNot,
			Args:     &ast.List{Items: []ast.Node{sublink}},
			Location: n.Pos(),
		}
	}

	in := &ast.In{
		Expr:     c.convert(n.X),
		Not:      n.Not,
		Location: n.Pos(),
	}
	if n.Table != nil {
		// "expr IN table" tests membership in the table's single column.
		in.Sel = parseRangeVar(n.Table, nil)
		return in
	}
	for _, e := range n.List {
		in.List = append(in.List, c.convert(e))
	}
	return in
}

func (c *cc) convertBetweenExpr(n *meyer.BetweenExpr) ast.Node {
	return &ast.BetweenExpr{
		Expr:     c.convert(n.X),
		Left:     c.convert(n.Lo),
		Right:    c.convert(n.Hi),
		Not:      n.Not,
		Location: n.Pos(),
	}
}

func (c *cc) convertCastExpr(n *meyer.CastExpr) ast.Node {
	name := typeName(n.Type)
	return &ast.TypeCast{
		Arg: c.convert(n.X),
		TypeName: &ast.TypeName{
			Name:        name,
			Names:       &ast.List{Items: []ast.Node{&ast.String{Str: strings.ToLower(name)}}},
			ArrayBounds: &ast.List{},
		},
		Location: n.Pos(),
	}
}

func (c *cc) convertCollateExpr(n *meyer.CollateExpr) ast.Node {
	return &ast.CollateExpr{
		Xpr:      c.convert(n.X),
		Arg:      NewIdentifier(n.Name),
		Location: n.Pos(),
	}
}

func (c *cc) convertCaseExpr(n *meyer.CaseExpr) ast.Node {
	expr := &ast.CaseExpr{
		Arg:       c.convertExpr(n.Operand),
		Args:      &ast.List{},
		Defresult: c.convertExpr(n.Else),
		Location:  n.Pos(),
	}
	for _, when := range n.Whens {
		expr.Args.Items = append(expr.Args.Items, &ast.CaseWhen{
			Expr:     c.convert(when.When),
			Result:   c.convert(when.Then),
			Location: when.Pos(),
		})
	}
	return expr
}

func (c *cc) convertFuncCall(n *meyer.FuncCall) ast.Node {
	funcName := identifier(n.Name)
	args := c.convertExprList(n.Args)

	if funcName == "coalesce" {
		return &ast.CoalesceExpr{
			Args:     args,
			Location: n.Pos(),
		}
	}

	call := &ast.FuncCall{
		Func: &ast.FuncName{
			Name: funcName,
		},
		Funcname: &ast.List{Items: []ast.Node{
			&ast.String{Str: funcName},
		}},
		AggStar:     n.Star,
		Args:        args,
		AggOrder:    &ast.List{Items: c.convertOrderBy(n.OrderBy)},
		AggDistinct: n.Distinct,
		AggFilter:   c.convertExpr(n.Filter),
		Location:    n.Pos(),
	}
	if n.Over != nil {
		if over, ok := c.convertWindowDef(n.Over).(*ast.WindowDef); ok {
			call.Over = over
		}
	}
	return call
}
