package dolphin

import (
	"log"
	"strings"

	pcast "github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/pingcap/tidb/pkg/parser/opcode"
	driver "github.com/pingcap/tidb/pkg/parser/test_driver"
	"github.com/pingcap/tidb/pkg/parser/types"

	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type cc struct {
	paramCount int
}

func todo(n pcast.Node) *ast.TODO {
	if debug.Active {
		log.Printf("dolphin.convert: Unknown node type %T\n", n)
	}
	return &ast.TODO{}
}

func identifier(id string) string {
	return strings.ToLower(id)
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func (c *cc) convertAlterTableStmt(n *pcast.AlterTableStmt) ast.Node {
	alt := &ast.AlterTableStmt{
		Table: parseTableName(n.Table),
		Cmds:  &ast.List{},
	}
	for _, spec := range n.Specs {
		switch spec.Tp {
		case pcast.AlterTableAddColumns:
			for _, def := range spec.NewColumns {
				name := def.Name.String()
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def:     convertColumnDef(def),
				})
			}

		case pcast.AlterTableDropColumn:
			name := spec.OldColumnName.String()
			alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
				Name:      &name,
				Subtype:   ast.AT_DropColumn,
				MissingOk: spec.IfExists,
			})

		case pcast.AlterTableChangeColumn:
			oldName := spec.OldColumnName.String()
			alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
				Name:    &oldName,
				Subtype: ast.AT_DropColumn,
			})

			for _, def := range spec.NewColumns {
				name := def.Name.String()
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def:     convertColumnDef(def),
				})
			}

		case pcast.AlterTableModifyColumn:
			for _, def := range spec.NewColumns {
				name := def.Name.String()
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_DropColumn,
				})
				alt.Cmds.Items = append(alt.Cmds.Items, &ast.AlterTableCmd{
					Name:    &name,
					Subtype: ast.AT_AddColumn,
					Def:     convertColumnDef(def),
				})
			}

		case pcast.AlterTableAlterColumn:
			// 	spew.Dump("alter column", spec)

		case pcast.AlterTableAddConstraint:
			// 	spew.Dump("add const", spec)

		case pcast.AlterTableRenameColumn:
			// TODO: Returning here may be incorrect if there are multiple specs
			oldName := spec.OldColumnName.String()
			newName := spec.NewColumnName.String()
			return &ast.RenameColumnStmt{
				Table:   parseTableName(n.Table),
				Col:     &ast.ColumnRef{Name: oldName},
				NewName: &newName,
			}

		case pcast.AlterTableRenameTable:
			// TODO: Returning here may be incorrect if there are multiple specs
			return &ast.RenameTableStmt{
				Table:   parseTableName(n.Table),
				NewName: &parseTableName(spec.NewTable).Name,
			}

		default:
			if debug.Active {
				log.Printf("dolphin.convert: Unknown alter table cmd %v\n", spec.Tp)
			}
			continue
		}
	}
	return alt
}

func (c *cc) convertAssignment(n *pcast.Assignment) *ast.ResTarget {
	name := identifier(n.Column.Name.String())
	return &ast.ResTarget{
		Name: &name,
		Val:  c.convert(n.Expr),
	}
}

// TODO: These codes should be defined in the sql/lang package
func opToName(o opcode.Op) string {
	switch o {
	// case opcode.And:
	// case opcode.BitNeg:
	// case opcode.Case:
	// case opcode.Div:
	case opcode.EQ:
		return "="
	case opcode.GE:
		return ">="
	case opcode.GT:
		return ">"
		// case opcode.In:
	case opcode.IntDiv:
		return "/"
	// case opcode.IsFalsity:
	// case opcode.IsNull:
	// case opcode.IsTruth:
	case opcode.LE:
		return "<="
	case opcode.LT:
		return "<"
	case opcode.LeftShift:
		return "<<"
		// case opcode.Like:
	case opcode.LogicAnd:
		return "&"
	case opcode.LogicOr:
		return "|"
	// case opcode.LogicXor:
	case opcode.Minus:
		return "-"
	case opcode.Mod:
		return "%"
	case opcode.Mul:
		return "*"
	case opcode.NE:
		return "!="
	case opcode.Not:
		return "!"
	// case opcode.NullEQ:
	// case opcode.Or:
	case opcode.Plus:
		return "+"
	case opcode.Regexp:
		return "~"
	case opcode.RightShift:
		return ">>"
	case opcode.Xor:
		return "#"
	default:
		return o.String()
	}
}

func (c *cc) convertBinaryOperationExpr(n *pcast.BinaryOperationExpr) ast.Node {
	if n.Op == opcode.LogicAnd || n.Op == opcode.LogicOr {
		return &ast.BoolExpr{
			// TODO: Set op
			Args: &ast.List{
				Items: []ast.Node{
					c.convert(n.L),
					c.convert(n.R),
				},
			},
		}
	} else {
		return &ast.A_Expr{
			// TODO: Set kind
			Name: &ast.List{
				Items: []ast.Node{
					&ast.String{Str: opToName(n.Op)},
				},
			},
			Lexpr: c.convert(n.L),
			Rexpr: c.convert(n.R),
		}
	}
}

func (c *cc) convertCreateTableStmt(n *pcast.CreateTableStmt) ast.Node {
	create := &ast.CreateTableStmt{
		Name:        parseTableName(n.Table),
		IfNotExists: n.IfNotExists,
	}
	if n.ReferTable != nil {
		create.ReferTable = parseTableName(n.ReferTable)
	}
	for _, def := range n.Cols {
		create.Cols = append(create.Cols, convertColumnDef(def))
	}
	for _, opt := range n.Options {
		switch opt.Tp {
		case pcast.TableOptionComment:
			create.Comment = opt.StrValue
		}
	}
	return create
}

func convertColumnDef(def *pcast.ColumnDef) *ast.ColumnDef {
	var vals *ast.List
	if len(def.Tp.GetElems()) > 0 {
		vals = &ast.List{}
		for i := range def.Tp.GetElems() {
			vals.Items = append(vals.Items, &ast.String{
				Str: def.Tp.GetElems()[i],
			})
		}
	}
	comment := ""
	for _, opt := range def.Options {
		switch opt.Tp {
		case pcast.ColumnOptionComment:
			if value, ok := opt.Expr.(*driver.ValueExpr); ok {
				comment = value.GetString()
			}
		}
	}
	columnDef := ast.ColumnDef{
		Colname:    def.Name.String(),
		TypeName:   &ast.TypeName{Name: types.TypeToStr(def.Tp.GetType(), def.Tp.GetCharset())},
		IsNotNull:  isNotNull(def),
		IsUnsigned: isUnsigned(def),
		Comment:    comment,
		Vals:       vals,
	}
	if def.Tp.GetFlen() >= 0 {
		length := def.Tp.GetFlen()
		columnDef.Length = &length
	}

	return &columnDef
}

func (c *cc) convertColumnNameExpr(n *pcast.ColumnNameExpr) *ast.ColumnRef {
	var items []ast.Node
	if schema := n.Name.Schema.String(); schema != "" {
		items = append(items, NewIdentifier(schema))
	}
	if table := n.Name.Table.String(); table != "" {
		items = append(items, NewIdentifier(table))
	}
	items = append(items, NewIdentifier(n.Name.Name.String()))
	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
		Location: n.OriginTextPosition(),
	}
}

func (c *cc) convertColumnNames(cols []*pcast.ColumnName) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for i := range cols {
		name := identifier(cols[i].Name.String())
		list.Items = append(list.Items, &ast.ResTarget{
			Name: &name,
		})
	}
	return list
}

func (c *cc) convertDeleteStmt(n *pcast.DeleteStmt) *ast.DeleteStmt {
	rels := c.convertTableRefsClause(n.TableRefs)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}
	relations := &ast.List{}
	convertToRangeVarList(rels, relations)

	stmt := &ast.DeleteStmt{
		Relations:     relations,
		WhereClause:   c.convert(n.Where),
		ReturningList: &ast.List{},
		WithClause:    c.convertWithClause(n.With),
	}
	if n.Limit != nil {
		stmt.LimitCount = c.convert(n.Limit.Count)
	}
	return stmt
}

func (c *cc) convertDropTableStmt(n *pcast.DropTableStmt) ast.Node {
	drop := &ast.DropTableStmt{IfExists: n.IfExists}
	for _, name := range n.Tables {
		drop.Tables = append(drop.Tables, parseTableName(name))
	}
	return drop
}

func (c *cc) convertRenameTableStmt(n *pcast.RenameTableStmt) ast.Node {
	list := &ast.List{Items: []ast.Node{}}
	for _, table := range n.TableToTables {
		list.Items = append(list.Items, &ast.RenameTableStmt{
			Table:   parseTableName(table.OldTable),
			NewName: &parseTableName(table.NewTable).Name,
		})
	}
	return list
}

func (c *cc) convertExistsSubqueryExpr(n *pcast.ExistsSubqueryExpr) *ast.SubLink {
	sublink := &ast.SubLink{}
	if ss, ok := c.convert(n.Sel).(*ast.SelectStmt); ok {
		sublink.Subselect = ss
	}
	return sublink
}

func (c *cc) convertFieldList(n *pcast.FieldList) *ast.List {
	fields := make([]ast.Node, len(n.Fields))
	for i := range n.Fields {
		fields[i] = c.convertSelectField(n.Fields[i])
	}
	return &ast.List{Items: fields}
}

func (c *cc) convertFuncCallExpr(n *pcast.FuncCallExpr) ast.Node {
	schema := n.Schema.String()
	name := strings.ToLower(n.FnName.String())

	// TODO: Deprecate the usage of Funcname
	items := []ast.Node{}
	if schema != "" {
		items = append(items, NewIdentifier(schema))
	}
	items = append(items, NewIdentifier(name))

	args := &ast.List{}
	for _, arg := range n.Args {
		args.Items = append(args.Items, c.convert(arg))
	}

	if schema == "" && name == "coalesce" {
		return &ast.CoalesceExpr{
			Args: args,
		}
	} else {
		return &ast.FuncCall{
			Args: args,
			Func: &ast.FuncName{
				Schema: schema,
				Name:   name,
			},
			Funcname: &ast.List{
				Items: items,
			},
			Location: n.OriginTextPosition(),
		}
	}
}

func (c *cc) convertInsertStmt(n *pcast.InsertStmt) *ast.InsertStmt {
	rels := c.convertTableRefsClause(n.Table)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}
	rel := rels.Items[0]
	rangeVar, ok := rel.(*ast.RangeVar)
	if !ok {
		panic("expected range var")
	}

	insert := &ast.InsertStmt{
		Relation:      rangeVar,
		Cols:          c.convertColumnNames(n.Columns),
		ReturningList: &ast.List{},
	}
	if ss, ok := c.convert(n.Select).(*ast.SelectStmt); ok {
		ss.ValuesLists = c.convertLists(n.Lists)
		insert.SelectStmt = ss
	} else {
		insert.SelectStmt = &ast.SelectStmt{
			FromClause:  &ast.List{},
			TargetList:  &ast.List{},
			ValuesLists: c.convertLists(n.Lists),
		}
	}

	if n.OnDuplicate != nil {
		targetList := &ast.List{}
		for _, a := range n.OnDuplicate {
			targetList.Items = append(targetList.Items, c.convertAssignment(a))
		}
		insert.OnConflictClause = &ast.OnConflictClause{
			TargetList: targetList,
			Location:   n.OriginTextPosition(),
		}
	}

	return insert
}

func (c *cc) convertLists(lists [][]pcast.ExprNode) *ast.List {
	list := &ast.List{Items: []ast.Node{}}
	for _, exprs := range lists {
		inner := &ast.List{Items: []ast.Node{}}
		for _, expr := range exprs {
			inner.Items = append(inner.Items, c.convert(expr))
		}
		list.Items = append(list.Items, inner)
	}
	return list
}

func (c *cc) convertParamMarkerExpr(n *driver.ParamMarkerExpr) *ast.ParamRef {
	// Parameter numbers start at one
	c.paramCount += 1
	return &ast.ParamRef{
		Number:   c.paramCount,
		Location: n.Offset,
	}
}

func (c *cc) convertSelectField(n *pcast.SelectField) *ast.ResTarget {
	var val ast.Node
	if n.WildCard != nil {
		val = c.convertWildCardField(n.WildCard)
	} else {
		val = c.convert(n.Expr)
	}
	var name *string
	if n.AsName.O != "" {
		asname := identifier(n.AsName.O)
		name = &asname
	}
	return &ast.ResTarget{
		// TODO: Populate Indirection field
		Name:     name,
		Val:      val,
		Location: n.Offset,
	}
}

func (c *cc) convertSelectStmt(n *pcast.SelectStmt) *ast.SelectStmt {
	windowClause := &ast.List{Items: make([]ast.Node, 0)}
	orderByClause := c.convertOrderByClause(n.OrderBy)
	if orderByClause != nil {
		windowClause.Items = append(windowClause.Items, orderByClause)
	}

	op, all := c.convertSetOprType(n.AfterSetOperator)
	stmt := &ast.SelectStmt{
		TargetList:   c.convertFieldList(n.Fields),
		FromClause:   c.convertTableRefsClause(n.From),
		GroupClause:  c.convertGroupByClause(n.GroupBy),
		HavingClause: c.convertHavingClause(n.Having),
		WhereClause:  c.convert(n.Where),
		WithClause:   c.convertWithClause(n.With),
		WindowClause: windowClause,
		Op:           op,
		All:          all,
	}
	if n.Limit != nil {
		stmt.LimitCount = c.convert(n.Limit.Count)
		stmt.LimitOffset = c.convert(n.Limit.Offset)
	}
	return stmt
}

func (c *cc) convertSubqueryExpr(n *pcast.SubqueryExpr) ast.Node {
	return c.convert(n.Query)
}

func (c *cc) convertTableRefsClause(n *pcast.TableRefsClause) *ast.List {
	if n == nil {
		return &ast.List{}
	}
	return c.convertJoin(n.TableRefs)
}

func (c *cc) convertCommonTableExpression(n *pcast.CommonTableExpression) *ast.CommonTableExpr {
	if n == nil {
		return nil
	}

	name := n.Name.String()

	columns := &ast.List{}
	for _, col := range n.ColNameList {
		columns.Items = append(columns.Items, NewIdentifier(col.String()))
	}

	return &ast.CommonTableExpr{
		Ctename:     &name,
		Ctequery:    c.convert(n.Query),
		Ctecolnames: columns,
	}
}

func (c *cc) convertWithClause(n *pcast.WithClause) *ast.WithClause {
	if n == nil {
		return nil
	}
	list := &ast.List{}
	for _, n := range n.CTEs {
		list.Items = append(list.Items, c.convertCommonTableExpression(n))
	}

	return &ast.WithClause{
		Ctes:      list,
		Recursive: n.IsRecursive,
		Location:  n.OriginTextPosition(),
	}
}

func (c *cc) convertUpdateStmt(n *pcast.UpdateStmt) *ast.UpdateStmt {
	rels := c.convertTableRefsClause(n.TableRefs)
	if len(rels.Items) != 1 {
		panic("expected one range var")
	}

	relations := &ast.List{}
	convertToRangeVarList(rels, relations)

	// TargetList
	list := &ast.List{}
	for _, a := range n.List {
		list.Items = append(list.Items, c.convertAssignment(a))
	}
	stmt := &ast.UpdateStmt{
		Relations:     relations,
		TargetList:    list,
		WhereClause:   c.convert(n.Where),
		FromClause:    &ast.List{},
		ReturningList: &ast.List{},
		WithClause:    c.convertWithClause(n.With),
	}
	if n.Limit != nil {
		stmt.LimitCount = c.convert(n.Limit.Count)
	}
	return stmt
}

func (c *cc) convertValueExpr(n *driver.ValueExpr) *ast.A_Const {
	switch n.TexprNode.Type.GetType() {
	case mysql.TypeBit:
	case mysql.TypeDate:
	case mysql.TypeDatetime:
	case mysql.TypeGeometry:
	case mysql.TypeJSON:
	case mysql.TypeNull:
	case mysql.TypeSet:
	case mysql.TypeShort:
	case mysql.TypeDuration:
	case mysql.TypeTimestamp:
		// TODO: Create an AST type for these?

	case mysql.TypeTiny,
		mysql.TypeInt24,
		mysql.TypeYear,
		mysql.TypeLong,
		mysql.TypeLonglong:
		return &ast.A_Const{
			Val: &ast.Integer{
				Ival: n.Datum.GetInt64(),
			},
			Location: n.OriginTextPosition(),
		}

	case mysql.TypeDouble,
		mysql.TypeFloat,
		mysql.TypeNewDecimal:
		return &ast.A_Const{
			Val: &ast.Float{
				// TODO: Extract the value from n.TexprNode
			},
			Location: n.OriginTextPosition(),
		}

	case mysql.TypeBlob, mysql.TypeString, mysql.TypeVarchar, mysql.TypeVarString, mysql.TypeLongBlob, mysql.TypeMediumBlob, mysql.TypeTinyBlob, mysql.TypeEnum:
	}
	return &ast.A_Const{
		Val: &ast.String{
			Str: n.Datum.GetString(),
		},
		Location: n.OriginTextPosition(),
	}
}

func (c *cc) convertWildCardField(n *pcast.WildCardField) *ast.ColumnRef {
	items := []ast.Node{}
	if t := n.Table.String(); t != "" {
		items = append(items, NewIdentifier(t))
	}
	items = append(items, &ast.A_Star{})

	return &ast.ColumnRef{
		Fields: &ast.List{
			Items: items,
		},
	}
}

func (c *cc) convertAdminStmt(n *pcast.AdminStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertAggregateFuncExpr(n *pcast.AggregateFuncExpr) *ast.FuncCall {
	name := strings.ToLower(n.F)
	fn := &ast.FuncCall{
		Func: &ast.FuncName{
			Name: name,
		},
		Funcname: &ast.List{
			Items: []ast.Node{
				NewIdentifier(name),
			},
		},
		Args:     &ast.List{},
		AggOrder: &ast.List{},
	}
	for _, a := range n.Args {
		if value, ok := a.(*driver.ValueExpr); ok {
			if value.GetInt64() == int64(1) {
				fn.AggStar = true
				continue
			}
		}
		fn.Args.Items = append(fn.Args.Items, c.convert(a))
	}
	if n.Distinct {
		fn.AggDistinct = true
	}
	return fn
}

func (c *cc) convertAlterDatabaseStmt(n *pcast.AlterDatabaseStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertAlterInstanceStmt(n *pcast.AlterInstanceStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertAlterTableSpec(n *pcast.AlterTableSpec) ast.Node {
	return todo(n)
}

func (c *cc) convertAlterUserStmt(n *pcast.AlterUserStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertAnalyzeTableStmt(n *pcast.AnalyzeTableStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertBRIEStmt(n *pcast.BRIEStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertBeginStmt(n *pcast.BeginStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertBetweenExpr(n *pcast.BetweenExpr) ast.Node {
	return &ast.BetweenExpr{
		Expr:     c.convert(n.Expr),
		Left:     c.convert(n.Left),
		Right:    c.convert(n.Right),
		Location: n.OriginTextPosition(),
		Not:      n.Not,
	}
}

func (c *cc) convertBinlogStmt(n *pcast.BinlogStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertByItem(n *pcast.ByItem) ast.Node {
	switch n.Expr.(type) {
	case *pcast.PositionExpr:
		return c.convertPositionExpr(n.Expr.(*pcast.PositionExpr))
	case *pcast.ColumnNameExpr:
		return c.convertColumnNameExpr(n.Expr.(*pcast.ColumnNameExpr))
	default:
		return todo(n)
	}
}

func (c *cc) convertCaseExpr(n *pcast.CaseExpr) ast.Node {
	if n == nil {
		return nil
	}
	list := &ast.List{Items: []ast.Node{}}
	for _, n := range n.WhenClauses {
		list.Items = append(list.Items, c.convertWhenClause(n))
	}
	return &ast.CaseExpr{
		Arg:       c.convert(n.Value),
		Args:      list,
		Defresult: c.convert(n.ElseClause),
		Location:  n.OriginTextPosition(),
	}
}

func (c *cc) convertCleanupTableLockStmt(n *pcast.CleanupTableLockStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertColumnDef(n *pcast.ColumnDef) ast.Node {
	return todo(n)
}

func (c *cc) convertColumnName(n *pcast.ColumnName) ast.Node {
	return todo(n)
}

func (c *cc) convertColumnPosition(n *pcast.ColumnPosition) ast.Node {
	return todo(n)
}

func (c *cc) convertCommitStmt(n *pcast.CommitStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCompareSubqueryExpr(n *pcast.CompareSubqueryExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertConstraint(n *pcast.Constraint) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateBindingStmt(n *pcast.CreateBindingStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateDatabaseStmt(n *pcast.CreateDatabaseStmt) ast.Node {
	return &ast.CreateSchemaStmt{
		Name:        &n.Name.O,
		IfNotExists: n.IfNotExists,
	}
}

func (c *cc) convertCreateIndexStmt(n *pcast.CreateIndexStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateSequenceStmt(n *pcast.CreateSequenceStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateStatisticsStmt(n *pcast.CreateStatisticsStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateUserStmt(n *pcast.CreateUserStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertCreateViewStmt(n *pcast.CreateViewStmt) ast.Node {
	return &ast.ViewStmt{
		View:            c.convertTableName(n.ViewName),
		Aliases:         &ast.List{},
		Query:           c.convert(n.Select),
		Replace:         n.OrReplace,
		Options:         &ast.List{},
		WithCheckOption: ast.ViewCheckOption(n.CheckOption),
	}
}

func (c *cc) convertDeallocateStmt(n *pcast.DeallocateStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDefaultExpr(n *pcast.DefaultExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertDeleteTableList(n *pcast.DeleteTableList) ast.Node {
	return todo(n)
}

func (c *cc) convertDoStmt(n *pcast.DoStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropBindingStmt(n *pcast.DropBindingStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropDatabaseStmt(n *pcast.DropDatabaseStmt) ast.Node {
	return &ast.DropSchemaStmt{
		MissingOk: !n.IfExists,
		Schemas: []*ast.String{
			NewIdentifier(n.Name.O),
		},
	}
}

func (c *cc) convertDropIndexStmt(n *pcast.DropIndexStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropSequenceStmt(n *pcast.DropSequenceStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropStatisticsStmt(n *pcast.DropStatisticsStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropStatsStmt(n *pcast.DropStatsStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertDropUserStmt(n *pcast.DropUserStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertExecuteStmt(n *pcast.ExecuteStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertExplainForStmt(n *pcast.ExplainForStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertExplainStmt(n *pcast.ExplainStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertFlashBackTableStmt(n *pcast.FlashBackTableStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertFlushStmt(n *pcast.FlushStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertFrameBound(n *pcast.FrameBound) ast.Node {
	return todo(n)
}

func (c *cc) convertFrameClause(n *pcast.FrameClause) ast.Node {
	return todo(n)
}

func (c *cc) convertFuncCastExpr(n *pcast.FuncCastExpr) ast.Node {
	return &ast.TypeCast{
		Arg:      c.convert(n.Expr),
		TypeName: &ast.TypeName{Name: types.TypeStr(n.Tp.GetType())},
	}
}

func (c *cc) convertGetFormatSelectorExpr(n *pcast.GetFormatSelectorExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertGrantRoleStmt(n *pcast.GrantRoleStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertGrantStmt(n *pcast.GrantStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertGroupByClause(n *pcast.GroupByClause) *ast.List {
	if n == nil {
		return &ast.List{}
	}

	var items []ast.Node
	for _, item := range n.Items {
		items = append(items, c.convertByItem(item))
	}

	return &ast.List{
		Items: items,
	}
}

func (c *cc) convertHavingClause(n *pcast.HavingClause) ast.Node {
	if n == nil {
		return nil
	}
	return c.convert(n.Expr)
}

func (c *cc) convertIndexLockAndAlgorithm(n *pcast.IndexLockAndAlgorithm) ast.Node {
	return todo(n)
}

func (c *cc) convertIndexPartSpecification(n *pcast.IndexPartSpecification) ast.Node {
	return todo(n)
}

func (c *cc) convertIsNullExpr(n *pcast.IsNullExpr) ast.Node {
	op := ast.BoolExprTypeIsNull
	if n.Not {
		op = ast.BoolExprTypeIsNotNull
	}
	return &ast.BoolExpr{
		Boolop: op,
		Args: &ast.List{
			Items: []ast.Node{
				c.convert(n.Expr),
			},
		},
	}
}

func (c *cc) convertIsTruthExpr(n *pcast.IsTruthExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertJoin(n *pcast.Join) *ast.List {
	if n == nil {
		return &ast.List{}
	}
	if n.Right != nil && n.Left != nil {
		// MySQL doesn't have a FULL join type
		joinType := ast.JoinType(n.Tp)
		if joinType >= ast.JoinTypeFull {
			joinType++
		}

		return &ast.List{
			Items: []ast.Node{&ast.JoinExpr{
				Jointype: joinType,
				Larg:     c.convert(n.Left),
				Rarg:     c.convert(n.Right),
				Quals:    c.convert(n.On),
			}},
		}
	}
	var tables []ast.Node
	if n.Right != nil {
		tables = append(tables, c.convert(n.Right))
	}
	if n.Left != nil {
		tables = append(tables, c.convert(n.Left))
	}
	return &ast.List{Items: tables}
}

func (c *cc) convertKillStmt(n *pcast.KillStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertLimit(n *pcast.Limit) ast.Node {
	return todo(n)
}

func (c *cc) convertLoadDataStmt(n *pcast.LoadDataStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertLoadStatsStmt(n *pcast.LoadStatsStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertLockTablesStmt(n *pcast.LockTablesStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertMatchAgainst(n *pcast.MatchAgainst) ast.Node {
	searchTerm := c.convert(n.Against)

	stringSearchTerm := &ast.TypeCast{
		Arg: searchTerm,
		TypeName: &ast.TypeName{
			Name: "text", // Use 'text' type which maps to string in Go
		},
		Location: n.OriginTextPosition(),
	}

	matchOperation := &ast.A_Const{
		Val: &ast.String{Str: "MATCH_AGAINST"},
	}

	return &ast.A_Expr{
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "AGAINST"},
			},
		},
		Lexpr:    matchOperation,
		Rexpr:    stringSearchTerm,
		Location: n.OriginTextPosition(),
	}
}

func (c *cc) convertMaxValueExpr(n *pcast.MaxValueExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertOnCondition(n *pcast.OnCondition) ast.Node {
	if n == nil {
		return nil
	}
	return c.convert(n.Expr)
}

func (c *cc) convertOnDeleteOpt(n *pcast.OnDeleteOpt) ast.Node {
	return todo(n)
}

func (c *cc) convertOnUpdateOpt(n *pcast.OnUpdateOpt) ast.Node {
	return todo(n)
}

func (c *cc) convertOrderByClause(n *pcast.OrderByClause) ast.Node {
	if n == nil {
		return nil
	}
	list := &ast.List{Items: []ast.Node{}}
	for _, item := range n.Items {
		list.Items = append(list.Items, c.convert(item.Expr))
	}
	return list
}

func (c *cc) convertParenthesesExpr(n *pcast.ParenthesesExpr) ast.Node {
	if n == nil {
		return nil
	}
	return c.convert(n.Expr)
}

func (c *cc) convertPartitionByClause(n *pcast.PartitionByClause) ast.Node {
	return todo(n)
}

func (c *cc) convertPatternInExpr(n *pcast.PatternInExpr) ast.Node {
	var list []ast.Node
	var val ast.Node

	expr := c.convert(n.Expr)

	for _, v := range n.List {
		val = c.convert(v)
		if val != nil {
			list = append(list, val)
		}
	}

	sel := c.convert(n.Sel)

	in := &ast.In{
		Expr:     expr,
		List:     list,
		Not:      n.Not,
		Sel:      sel,
		Location: n.OriginTextPosition(),
	}

	return in
}

func (c *cc) convertPatternLikeExpr(n *pcast.PatternLikeOrIlikeExpr) ast.Node {
	return &ast.A_Expr{
		Kind: ast.A_Expr_Kind(9),
		Name: &ast.List{
			Items: []ast.Node{
				&ast.String{Str: "~~"},
			},
		},
		Lexpr: c.convert(n.Expr),
		Rexpr: c.convert(n.Pattern),
	}
}

func (c *cc) convertPatternRegexpExpr(n *pcast.PatternRegexpExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertPositionExpr(n *pcast.PositionExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertPrepareStmt(n *pcast.PrepareStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertPrivElem(n *pcast.PrivElem) ast.Node {
	return todo(n)
}

func (c *cc) convertRecoverTableStmt(n *pcast.RecoverTableStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertReferenceDef(n *pcast.ReferenceDef) ast.Node {
	return todo(n)
}

func (c *cc) convertRepairTableStmt(n *pcast.RepairTableStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertRevokeRoleStmt(n *pcast.RevokeRoleStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertRevokeStmt(n *pcast.RevokeStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertRollbackStmt(n *pcast.RollbackStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertRowExpr(n *pcast.RowExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertSetCollationExpr(n *pcast.SetCollationExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertSetConfigStmt(n *pcast.SetConfigStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertSetDefaultRoleStmt(n *pcast.SetDefaultRoleStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertSetOprType(n *pcast.SetOprType) (op ast.SetOperation, all bool) {
	if n == nil {
		return
	}

	switch *n {
	case pcast.Union:
		op = ast.Union
	case pcast.UnionAll:
		op = ast.Union
		all = true
	case pcast.Intersect:
		op = ast.Intersect
	case pcast.IntersectAll:
		op = ast.Intersect
		all = true
	case pcast.Except:
		op = ast.Except
	case pcast.ExceptAll:
		op = ast.Except
		all = true
	}
	return
}

// convertSetOprSelectList converts a list of SELECT from the Pingcap parser
// into a tree. It is called for UNION, INTERSECT or EXCLUDE operation.
//
// Given an union with the following nodes:
//
//	[Select{1}, Select{2}, Select{3}, Select{4}]
//
// The function will return:
//
//	Select{
//	    Larg: Select{
//	        Larg: Select{
//	            Larg: Select{1},
//	            Rarg: Select{2},
//	            Op: Union
//	        },
//	        Rarg: Select{3},
//	        Op: Union,
//	    },
//	    Rarg: Select{4},
//	    Op: Union,
//	}
func (c *cc) convertSetOprSelectList(n *pcast.SetOprSelectList) ast.Node {
	selectStmts := make([]*ast.SelectStmt, len(n.Selects))
	for i, node := range n.Selects {
		switch node := node.(type) {
		case *pcast.SelectStmt:
			selectStmts[i] = c.convertSelectStmt(node)
		case *pcast.SetOprSelectList:
			selectStmts[i] = c.convertSetOprSelectList(node).(*ast.SelectStmt)
		}
	}

	op, all := c.convertSetOprType(n.AfterSetOperator)
	tree := &ast.SelectStmt{
		TargetList:  &ast.List{},
		FromClause:  &ast.List{},
		WhereClause: nil,
		Op:          op,
		All:         all,
		WithClause:  c.convertWithClause(n.With),
	}
	for _, stmt := range selectStmts {
		// We move Op and All from the child to the parent.
		op, all := stmt.Op, stmt.All
		stmt.Op, stmt.All = ast.None, false

		switch {
		case tree.Larg == nil:
			tree.Larg = stmt
		case tree.Rarg == nil:
			tree.Rarg = stmt
			tree.Op = op
			tree.All = all
		default:
			tree = &ast.SelectStmt{
				TargetList:  &ast.List{},
				FromClause:  &ast.List{},
				WhereClause: nil,
				Larg:        tree,
				Rarg:        stmt,
				Op:          op,
				All:         all,
				WithClause:  c.convertWithClause(n.With),
			}
		}
	}

	return tree
}

func (c *cc) convertSetOprStmt(n *pcast.SetOprStmt) ast.Node {
	if n.SelectList != nil {
		sn := c.convertSetOprSelectList(n.SelectList)
		if ss, ok := sn.(*ast.SelectStmt); ok && n.Limit != nil {
			ss.LimitOffset = c.convert(n.Limit.Offset)
			ss.LimitCount = c.convert(n.Limit.Count)
		}
		return sn
	}
	return todo(n)
}

func (c *cc) convertSetPwdStmt(n *pcast.SetPwdStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertSetRoleStmt(n *pcast.SetRoleStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertSetStmt(n *pcast.SetStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertShowStmt(n *pcast.ShowStmt) ast.Node {
	if n.Tp != pcast.ShowWarnings {
		return todo(n)
	}
	level := "level"
	code := "code"
	message := "message"
	stmt := &ast.SelectStmt{
		FromClause: &ast.List{},
		TargetList: &ast.List{
			Items: []ast.Node{
				&ast.ResTarget{
					Name: &level,
					Val:  &ast.A_Const{Val: &ast.String{}},
				},
				&ast.ResTarget{
					Name: &code,
					Val:  &ast.A_Const{Val: &ast.Integer{}},
				},
				&ast.ResTarget{
					Name: &message,
					Val:  &ast.A_Const{Val: &ast.String{}},
				},
			},
		},
	}
	return stmt
}

func (c *cc) convertShutdownStmt(n *pcast.ShutdownStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertSplitRegionStmt(n *pcast.SplitRegionStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertTableName(n *pcast.TableName) *ast.RangeVar {
	schema := identifier(n.Schema.String())
	rel := identifier(n.Name.String())
	return &ast.RangeVar{
		Schemaname: &schema,
		Relname:    &rel,
	}
}

func (c *cc) convertTableNameExpr(n *pcast.TableNameExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertTableOptimizerHint(n *pcast.TableOptimizerHint) ast.Node {
	return todo(n)
}

func (c *cc) convertTableSource(node *pcast.TableSource) ast.Node {
	if node == nil {
		return nil
	}
	alias := node.AsName.String()
	switch n := node.Source.(type) {

	case *pcast.SelectStmt, *pcast.SetOprStmt:
		rs := &ast.RangeSubselect{
			Subquery: c.convert(n),
		}
		if alias != "" {
			rs.Alias = &ast.Alias{Aliasname: &alias}
		}
		return rs

	case *pcast.TableName:
		rv := c.convertTableName(n)
		if alias != "" {
			rv.Alias = &ast.Alias{Aliasname: &alias}
		}
		return rv

	default:
		return todo(n)
	}
}

func (c *cc) convertTableToTable(n *pcast.TableToTable) ast.Node {
	return todo(n)
}

func (c *cc) convertTimeUnitExpr(n *pcast.TimeUnitExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertTraceStmt(n *pcast.TraceStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertTrimDirectionExpr(n *pcast.TrimDirectionExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertTruncateTableStmt(n *pcast.TruncateTableStmt) *ast.TruncateStmt {
	return &ast.TruncateStmt{
		Relations: toList(n.Table),
	}
}

func (c *cc) convertUnaryOperationExpr(n *pcast.UnaryOperationExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertUnlockTablesStmt(n *pcast.UnlockTablesStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertUseStmt(n *pcast.UseStmt) ast.Node {
	return todo(n)
}

func (c *cc) convertValuesExpr(n *pcast.ValuesExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertVariableAssignment(n *pcast.VariableAssignment) ast.Node {
	return todo(n)
}

func (c *cc) convertVariableExpr(n *pcast.VariableExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertWhenClause(n *pcast.WhenClause) ast.Node {
	if n == nil {
		return nil
	}
	return &ast.CaseWhen{
		Expr:     c.convert(n.Expr),
		Result:   c.convert(n.Result),
		Location: n.OriginTextPosition(),
	}
}

func (c *cc) convertWindowFuncExpr(n *pcast.WindowFuncExpr) ast.Node {
	return todo(n)
}

func (c *cc) convertWindowSpec(n *pcast.WindowSpec) ast.Node {
	return todo(n)
}

func (c *cc) convertCallStmt(n *pcast.CallStmt) ast.Node {
	var funcname ast.List
	for _, s := range []string{n.Procedure.Schema.L, n.Procedure.FnName.L} {
		if s != "" {
			funcname.Items = append(funcname.Items, NewIdentifier(s))
		}
	}
	var args ast.List
	for _, a := range n.Procedure.Args {
		args.Items = append(args.Items, c.convert(a))
	}
	return &ast.CallStmt{
		FuncCall: &ast.FuncCall{
			Func: &ast.FuncName{
				Schema: n.Procedure.Schema.L,
				Name:   n.Procedure.FnName.L,
			},
			Funcname: &funcname,
			Args:     &args,
			Location: n.OriginTextPosition(),
		},
	}
}

func (c *cc) convertProcedureInfo(n *pcast.ProcedureInfo) ast.Node {
	var params ast.List
	for _, sp := range n.ProcedureParam {
		paramName := sp.ParamName
		params.Items = append(params.Items, &ast.FuncParam{
			Name: &paramName,
			Type: &ast.TypeName{Name: types.TypeToStr(sp.ParamType.GetType(), sp.ParamType.GetCharset())},
		})
	}
	return &ast.CreateFunctionStmt{
		Params: &params,
		Func: &ast.FuncName{
			Schema: n.ProcedureName.Schema.L,
			Name:   n.ProcedureName.Name.L,
		},
	}
}

func (c *cc) convert(node pcast.Node) ast.Node {
	switch n := node.(type) {

	case *driver.ParamMarkerExpr:
		return c.convertParamMarkerExpr(n)

	case *driver.ValueExpr:
		return c.convertValueExpr(n)

	case *pcast.AdminStmt:
		return c.convertAdminStmt(n)

	case *pcast.AggregateFuncExpr:
		return c.convertAggregateFuncExpr(n)

	case *pcast.AlterDatabaseStmt:
		return c.convertAlterDatabaseStmt(n)

	case *pcast.AlterInstanceStmt:
		return c.convertAlterInstanceStmt(n)

	case *pcast.AlterTableSpec:
		return c.convertAlterTableSpec(n)

	case *pcast.AlterTableStmt:
		return c.convertAlterTableStmt(n)

	case *pcast.AlterUserStmt:
		return c.convertAlterUserStmt(n)

	case *pcast.AnalyzeTableStmt:
		return c.convertAnalyzeTableStmt(n)

	case *pcast.Assignment:
		return c.convertAssignment(n)

	case *pcast.BRIEStmt:
		return c.convertBRIEStmt(n)

	case *pcast.BeginStmt:
		return c.convertBeginStmt(n)

	case *pcast.BetweenExpr:
		return c.convertBetweenExpr(n)

	case *pcast.BinaryOperationExpr:
		return c.convertBinaryOperationExpr(n)

	case *pcast.BinlogStmt:
		return c.convertBinlogStmt(n)

	case *pcast.ByItem:
		return c.convertByItem(n)

	case *pcast.CallStmt:
		return c.convertCallStmt(n)

	case *pcast.CaseExpr:
		return c.convertCaseExpr(n)

	case *pcast.CleanupTableLockStmt:
		return c.convertCleanupTableLockStmt(n)

	case *pcast.ColumnDef:
		return c.convertColumnDef(n)

	case *pcast.ColumnName:
		return c.convertColumnName(n)

	case *pcast.ColumnNameExpr:
		return c.convertColumnNameExpr(n)

	case *pcast.ColumnPosition:
		return c.convertColumnPosition(n)

	case *pcast.CommitStmt:
		return c.convertCommitStmt(n)

	case *pcast.CompareSubqueryExpr:
		return c.convertCompareSubqueryExpr(n)

	case *pcast.Constraint:
		return c.convertConstraint(n)

	case *pcast.CreateBindingStmt:
		return c.convertCreateBindingStmt(n)

	case *pcast.CreateDatabaseStmt:
		return c.convertCreateDatabaseStmt(n)

	case *pcast.CreateIndexStmt:
		return c.convertCreateIndexStmt(n)

	case *pcast.CreateSequenceStmt:
		return c.convertCreateSequenceStmt(n)

	case *pcast.CreateStatisticsStmt:
		return c.convertCreateStatisticsStmt(n)

	case *pcast.CreateTableStmt:
		return c.convertCreateTableStmt(n)

	case *pcast.CreateUserStmt:
		return c.convertCreateUserStmt(n)

	case *pcast.CreateViewStmt:
		return c.convertCreateViewStmt(n)

	case *pcast.DeallocateStmt:
		return c.convertDeallocateStmt(n)

	case *pcast.DefaultExpr:
		return c.convertDefaultExpr(n)

	case *pcast.DeleteStmt:
		return c.convertDeleteStmt(n)

	case *pcast.DeleteTableList:
		return c.convertDeleteTableList(n)

	case *pcast.DoStmt:
		return c.convertDoStmt(n)

	case *pcast.DropBindingStmt:
		return c.convertDropBindingStmt(n)

	case *pcast.DropDatabaseStmt:
		return c.convertDropDatabaseStmt(n)

	case *pcast.DropIndexStmt:
		return c.convertDropIndexStmt(n)

	case *pcast.DropSequenceStmt:
		return c.convertDropSequenceStmt(n)

	case *pcast.DropStatisticsStmt:
		return c.convertDropStatisticsStmt(n)

	case *pcast.DropStatsStmt:
		return c.convertDropStatsStmt(n)

	case *pcast.DropTableStmt:
		return c.convertDropTableStmt(n)

	case *pcast.DropUserStmt:
		return c.convertDropUserStmt(n)

	case *pcast.ExecuteStmt:
		return c.convertExecuteStmt(n)

	case *pcast.ExistsSubqueryExpr:
		return c.convertExistsSubqueryExpr(n)

	case *pcast.ExplainForStmt:
		return c.convertExplainForStmt(n)

	case *pcast.ExplainStmt:
		return c.convertExplainStmt(n)

	case *pcast.FieldList:
		return c.convertFieldList(n)

	case *pcast.FlashBackTableStmt:
		return c.convertFlashBackTableStmt(n)

	case *pcast.FlushStmt:
		return c.convertFlushStmt(n)

	case *pcast.FrameBound:
		return c.convertFrameBound(n)

	case *pcast.FrameClause:
		return c.convertFrameClause(n)

	case *pcast.FuncCallExpr:
		return c.convertFuncCallExpr(n)

	case *pcast.FuncCastExpr:
		return c.convertFuncCastExpr(n)

	case *pcast.GetFormatSelectorExpr:
		return c.convertGetFormatSelectorExpr(n)

	case *pcast.GrantRoleStmt:
		return c.convertGrantRoleStmt(n)

	case *pcast.GrantStmt:
		return c.convertGrantStmt(n)

	case *pcast.GroupByClause:
		return c.convertGroupByClause(n)

	case *pcast.HavingClause:
		return c.convertHavingClause(n)

	case *pcast.IndexLockAndAlgorithm:
		return c.convertIndexLockAndAlgorithm(n)

	case *pcast.IndexPartSpecification:
		return c.convertIndexPartSpecification(n)

	case *pcast.InsertStmt:
		return c.convertInsertStmt(n)

	case *pcast.IsNullExpr:
		return c.convertIsNullExpr(n)

	case *pcast.IsTruthExpr:
		return c.convertIsTruthExpr(n)

	case *pcast.Join:
		return c.convertJoin(n)

	case *pcast.KillStmt:
		return c.convertKillStmt(n)

	case *pcast.Limit:
		return c.convertLimit(n)

	case *pcast.LoadDataStmt:
		return c.convertLoadDataStmt(n)

	case *pcast.LoadStatsStmt:
		return c.convertLoadStatsStmt(n)

	case *pcast.LockTablesStmt:
		return c.convertLockTablesStmt(n)

	case *pcast.MatchAgainst:
		return c.convertMatchAgainst(n)

	case *pcast.MaxValueExpr:
		return c.convertMaxValueExpr(n)

	case *pcast.OnCondition:
		return c.convertOnCondition(n)

	case *pcast.OnDeleteOpt:
		return c.convertOnDeleteOpt(n)

	case *pcast.OnUpdateOpt:
		return c.convertOnUpdateOpt(n)

	case *pcast.OrderByClause:
		return c.convertOrderByClause(n)

	case *pcast.ParenthesesExpr:
		return c.convertParenthesesExpr(n)

	case *pcast.PartitionByClause:
		return c.convertPartitionByClause(n)

	case *pcast.PatternInExpr:
		return c.convertPatternInExpr(n)

	case *pcast.PatternLikeOrIlikeExpr:
		return c.convertPatternLikeExpr(n)

	case *pcast.PatternRegexpExpr:
		return c.convertPatternRegexpExpr(n)

	case *pcast.PositionExpr:
		return c.convertPositionExpr(n)

	case *pcast.PrepareStmt:
		return c.convertPrepareStmt(n)

	case *pcast.PrivElem:
		return c.convertPrivElem(n)

	case *pcast.ProcedureInfo:
		return c.convertProcedureInfo(n)

	case *pcast.RecoverTableStmt:
		return c.convertRecoverTableStmt(n)

	case *pcast.ReferenceDef:
		return c.convertReferenceDef(n)

	case *pcast.RenameTableStmt:
		return c.convertRenameTableStmt(n)

	case *pcast.RepairTableStmt:
		return c.convertRepairTableStmt(n)

	case *pcast.RevokeRoleStmt:
		return c.convertRevokeRoleStmt(n)

	case *pcast.RevokeStmt:
		return c.convertRevokeStmt(n)

	case *pcast.RollbackStmt:
		return c.convertRollbackStmt(n)

	case *pcast.RowExpr:
		return c.convertRowExpr(n)

	case *pcast.SelectField:
		return c.convertSelectField(n)

	case *pcast.SelectStmt:
		return c.convertSelectStmt(n)

	case *pcast.SetCollationExpr:
		return c.convertSetCollationExpr(n)

	case *pcast.SetConfigStmt:
		return c.convertSetConfigStmt(n)

	case *pcast.SetDefaultRoleStmt:
		return c.convertSetDefaultRoleStmt(n)

	case *pcast.SetOprSelectList:
		return c.convertSetOprSelectList(n)

	case *pcast.SetOprStmt:
		return c.convertSetOprStmt(n)

	case *pcast.SetPwdStmt:
		return c.convertSetPwdStmt(n)

	case *pcast.SetRoleStmt:
		return c.convertSetRoleStmt(n)

	case *pcast.SetStmt:
		return c.convertSetStmt(n)

	case *pcast.ShowStmt:
		return c.convertShowStmt(n)

	case *pcast.ShutdownStmt:
		return c.convertShutdownStmt(n)

	case *pcast.SplitRegionStmt:
		return c.convertSplitRegionStmt(n)

	case *pcast.SubqueryExpr:
		return c.convertSubqueryExpr(n)

	case *pcast.TableName:
		return c.convertTableName(n)

	case *pcast.TableNameExpr:
		return c.convertTableNameExpr(n)

	case *pcast.TableOptimizerHint:
		return c.convertTableOptimizerHint(n)

	case *pcast.TableRefsClause:
		return c.convertTableRefsClause(n)

	case *pcast.TableSource:
		return c.convertTableSource(n)

	case *pcast.TableToTable:
		return c.convertTableToTable(n)

	case *pcast.TimeUnitExpr:
		return c.convertTimeUnitExpr(n)

	case *pcast.TraceStmt:
		return c.convertTraceStmt(n)

	case *pcast.TrimDirectionExpr:
		return c.convertTrimDirectionExpr(n)

	case *pcast.TruncateTableStmt:
		return c.convertTruncateTableStmt(n)

	case *pcast.UnaryOperationExpr:
		return c.convertUnaryOperationExpr(n)

	case *pcast.UnlockTablesStmt:
		return c.convertUnlockTablesStmt(n)

	case *pcast.UpdateStmt:
		return c.convertUpdateStmt(n)

	case *pcast.UseStmt:
		return c.convertUseStmt(n)

	case *pcast.ValuesExpr:
		return c.convertValuesExpr(n)

	case *pcast.VariableAssignment:
		return c.convertVariableAssignment(n)

	case *pcast.VariableExpr:
		return c.convertVariableExpr(n)

	case *pcast.WhenClause:
		return c.convertWhenClause(n)

	case *pcast.WildCardField:
		return c.convertWildCardField(n)

	case *pcast.WindowFuncExpr:
		return c.convertWindowFuncExpr(n)

	case *pcast.WindowSpec:
		return c.convertWindowSpec(n)

	case nil:
		return nil

	default:
		return todo(n)
	}
}
