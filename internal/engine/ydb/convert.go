package ydb

import (
	"log"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/debug"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	parser "github.com/ydb-platform/yql-parsers/go"
)

type cc struct {
	paramCount int
	content    string
}

func (c *cc) pos(token antlr.Token) int {
	if token == nil {
		return 0
	}
	runeIdx := token.GetStart()
	return byteOffsetFromRuneIndex(c.content, runeIdx)
}

type node interface {
	GetParser() antlr.Parser
}

func todo(funcname string, n node) *ast.TODO {
	if debug.Active {
		log.Printf("ydb.%s: Unknown node type %T\n", funcname, n)
	}
	return &ast.TODO{}
}

func identifier(id string) string {
	if len(id) >= 2 && id[0] == '"' && id[len(id)-1] == '"' {
		unquoted, _ := strconv.Unquote(id)
		return unquoted
	}
	return strings.ToLower(id)
}

func stripQuotes(s string) string {
	if len(s) >= 2 && (s[0] == '\'' || s[0] == '"') && s[0] == s[len(s)-1] {
		return s[1 : len(s)-1]
	}
	return s
}

func NewIdentifier(t string) *ast.String {
	return &ast.String{Str: identifier(t)}
}

func (c *cc) convertDrop_role_stmtCOntext(n *parser.Drop_role_stmtContext) ast.Node {
	if n.DROP() == nil || (n.USER() == nil && n.GROUP() == nil) || len(n.AllRole_name()) == 0 {
		return todo("Drop_role_stmtContext", n)
	}

	stmt := &ast.DropRoleStmt{
		MissingOk: n.IF() != nil && n.EXISTS() != nil,
		Roles:     &ast.List{},
	}

	for _, role := range n.AllRole_name() {
		member, isParam, _ := c.extractRoleSpec(role, ast.RoleSpecType(1))
		if member == nil {
			return todo("Drop_role_stmtContext", n)
		}

		if debug.Active && isParam {
			log.Printf("YDB does not currently support parameters in the DROP ROLE statement")
		}

		stmt.Roles.Items = append(stmt.Roles.Items, member)
	}

	return stmt
}

func (c *cc) convertAlter_group_stmtContext(n *parser.Alter_group_stmtContext) ast.Node {
	if n.ALTER() == nil || n.GROUP() == nil || len(n.AllRole_name()) == 0 {
		return todo("convertAlter_group_stmtContext", n)
	}
	role, paramFlag, _ := c.extractRoleSpec(n.Role_name(0), ast.RoleSpecType(1))
	if role == nil {
		return todo("convertAlter_group_stmtContext", n)
	}

	if debug.Active && paramFlag {
		log.Printf("YDB does not currently support parameters in the ALTER GROUP statement")
	}

	stmt := &ast.AlterRoleStmt{
		Role:    role,
		Action:  1,
		Options: &ast.List{},
	}

	switch {
	case n.RENAME() != nil && n.TO() != nil && len(n.AllRole_name()) > 1:
		newName := c.convert(n.Role_name(1))
		action := "rename"

		defElem := &ast.DefElem{
			Defname:   &action,
			Defaction: ast.DefElemAction(1),
			Location:  c.pos(n.Role_name(1).GetStart()),
		}

		bindFlag := true
		switch v := newName.(type) {
		case *ast.A_Const:
			switch val := v.Val.(type) {
			case *ast.String:
				bindFlag = false
				defElem.Arg = val
			case *ast.Boolean:
				defElem.Arg = val
			default:
				return todo("convertAlter_group_stmtContext", n)
			}
		case *ast.ParamRef, *ast.A_Expr:
			defElem.Arg = newName
		default:
			return todo("convertAlter_group_stmtContext", n)
		}

		if debug.Active && !paramFlag && bindFlag {
			log.Printf("YDB does not currently support parameters in the ALTER GROUP statement")
		}

		stmt.Options.Items = append(stmt.Options.Items, defElem)

	case (n.ADD() != nil || n.DROP() != nil) && len(n.AllRole_name()) > 1:
		defname := "rolemembers"
		optionList := &ast.List{}
		for _, role := range n.AllRole_name()[1:] {
			member, isParam, _ := c.extractRoleSpec(role, ast.RoleSpecType(1))
			if member == nil {
				return todo("convertAlter_group_stmtContext", n)
			}

			if debug.Active && isParam && !paramFlag {
				log.Printf("YDB does not currently support parameters in the ALTER GROUP statement")
			}

			optionList.Items = append(optionList.Items, member)
		}

		var action ast.DefElemAction
		if n.ADD() != nil {
			action = 3
		} else {
			action = 4
		}

		stmt.Options.Items = append(stmt.Options.Items, &ast.DefElem{
			Defname:   &defname,
			Arg:       optionList,
			Defaction: action,
			Location:  c.pos(n.GetStart()),
		})
	}

	return stmt
}

func (c *cc) convertAlter_user_stmtContext(n *parser.Alter_user_stmtContext) ast.Node {
	if n.ALTER() == nil || n.USER() == nil || len(n.AllRole_name()) == 0 {
		return todo("Alter_user_stmtContext", n)
	}

	role, paramFlag, _ := c.extractRoleSpec(n.Role_name(0), ast.RoleSpecType(1))
	if role == nil {
		return todo("convertAlter_group_stmtContext", n)
	}

	if debug.Active && paramFlag {
		log.Printf("YDB does not currently support parameters in the ALTER USER statement")
	}

	stmt := &ast.AlterRoleStmt{
		Role:    role,
		Action:  1,
		Options: &ast.List{},
	}

	switch {
	case n.RENAME() != nil && n.TO() != nil && len(n.AllRole_name()) > 1:
		newName := c.convert(n.Role_name(1))
		action := "rename"

		defElem := &ast.DefElem{
			Defname:   &action,
			Defaction: ast.DefElemAction(1),
			Location:  c.pos(n.Role_name(1).GetStart()),
		}

		bindFlag := true
		switch v := newName.(type) {
		case *ast.A_Const:
			switch val := v.Val.(type) {
			case *ast.String:
				bindFlag = false
				defElem.Arg = val
			case *ast.Boolean:
				defElem.Arg = val
			default:
				return todo("Alter_user_stmtContext", n)
			}
		case *ast.ParamRef, *ast.A_Expr:
			defElem.Arg = newName
		default:
			return todo("Alter_user_stmtContext", n)
		}

		if debug.Active && !paramFlag && bindFlag {
			log.Printf("YDB does not currently support parameters in the ALTER USER statement")
		}

		stmt.Options.Items = append(stmt.Options.Items, defElem)

	case len(n.AllUser_option()) > 0:
		for _, opt := range n.AllUser_option() {
			if node := c.convert(opt); node != nil {
				stmt.Options.Items = append(stmt.Options.Items, node)
			}
		}
	}

	return stmt
}

func (c *cc) convertCreate_group_stmtContext(n *parser.Create_group_stmtContext) ast.Node {
	if n.CREATE() == nil || n.GROUP() == nil || len(n.AllRole_name()) == 0 {
		return todo("Create_group_stmtContext", n)
	}
	groupName := c.convert(n.Role_name(0))

	stmt := &ast.CreateRoleStmt{
		StmtType: ast.RoleStmtType(3),
		Options:  &ast.List{},
	}

	paramFlag := true
	switch v := groupName.(type) {
	case *ast.A_Const:
		switch val := v.Val.(type) {
		case *ast.String:
			paramFlag = false
			stmt.Role = &val.Str
		case *ast.Boolean:
			stmt.BindRole = groupName
		default:
			return todo("convertCreate_group_stmtContext", n)
		}
	case *ast.ParamRef, *ast.A_Expr:
		stmt.BindRole = groupName
	default:
		return todo("convertCreate_group_stmtContext", n)
	}

	if debug.Active && paramFlag {
		log.Printf("YDB does not currently support parameters in the CREATE GROUP statement")
	}

	if n.WITH() != nil && n.USER() != nil && len(n.AllRole_name()) > 1 {
		defname := "rolemembers"
		optionList := &ast.List{}
		for _, role := range n.AllRole_name()[1:] {
			member, isParam, _ := c.extractRoleSpec(role, ast.RoleSpecType(1))
			if member == nil {
				return todo("convertCreate_group_stmtContext", n)
			}

			if debug.Active && isParam && !paramFlag {
				log.Printf("YDB does not currently support parameters in the CREATE GROUP statement")
			}

			optionList.Items = append(optionList.Items, member)
		}

		stmt.Options.Items = append(stmt.Options.Items, &ast.DefElem{
			Defname:  &defname,
			Arg:      optionList,
			Location: c.pos(n.GetStart()),
		})
	}

	return stmt
}

func (c *cc) convertUse_stmtContext(n *parser.Use_stmtContext) ast.Node {
	if n.USE() != nil && n.Cluster_expr() != nil {
		clusterExpr := c.convert(n.Cluster_expr())
		stmt := &ast.UseStmt{
			Xpr:      clusterExpr,
			Location: c.pos(n.GetStart()),
		}
		return stmt
	}
	return todo("convertUse_stmtContext", n)
}

func (c *cc) convertCluster_exprContext(n *parser.Cluster_exprContext) ast.Node {
	var node ast.Node

	switch {
	case n.Pure_column_or_named() != nil:
		pureCtx := n.Pure_column_or_named()
		if anID := pureCtx.An_id(); anID != nil {
			name := parseAnId(anID)
			node = &ast.ColumnRef{
				Fields:   &ast.List{Items: []ast.Node{NewIdentifier(name)}},
				Location: c.pos(anID.GetStart()),
			}
		} else if bp := pureCtx.Bind_parameter(); bp != nil {
			node = c.convert(bp)
		}
	case n.ASTERISK() != nil:
		node = &ast.A_Star{}
	default:
		return todo("convertCluster_exprContext", n)
	}

	if n.An_id() != nil && n.COLON() != nil {
		name := parseAnId(n.An_id())
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: ":"}}},
			Lexpr:    &ast.String{Str: name},
			Rexpr:    node,
			Location: c.pos(n.GetStart()),
		}
	}

	return node
}

func (c *cc) convertCreate_user_stmtContext(n *parser.Create_user_stmtContext) ast.Node {
	if n.CREATE() == nil || n.USER() == nil || n.Role_name() == nil {
		return todo("convertCreate_user_stmtContext", n)
	}
	roleNode := c.convert(n.Role_name())

	stmt := &ast.CreateRoleStmt{
		StmtType: ast.RoleStmtType(2),
		Options:  &ast.List{},
	}

	paramFlag := true
	switch v := roleNode.(type) {
	case *ast.A_Const:
		switch val := v.Val.(type) {
		case *ast.String:
			paramFlag = false
			stmt.Role = &val.Str
		case *ast.Boolean:
			stmt.BindRole = roleNode
		default:
			return todo("convertCreate_user_stmtContext", n)
		}
	case *ast.ParamRef, *ast.A_Expr:
		stmt.BindRole = roleNode
	default:
		return todo("convertCreate_user_stmtContext", n)
	}

	if debug.Active && paramFlag {
		log.Printf("YDB does not currently support parameters in the CREATE USER statement")
	}

	if len(n.AllUser_option()) > 0 {
		options := []ast.Node{}
		for _, opt := range n.AllUser_option() {
			if node := c.convert(opt); node != nil {
				options = append(options, node)
			}
		}
		if len(options) > 0 {
			stmt.Options = &ast.List{Items: options}
		}
	}
	return stmt
}

func (c *cc) convertUser_optionContext(n *parser.User_optionContext) ast.Node {
	switch {
	case n.Authentication_option() != nil:
		aOpt := n.Authentication_option()
		if pOpt := aOpt.Password_option(); pOpt != nil {
			if pOpt.PASSWORD() != nil {
				name := "password"
				pValue := pOpt.Password_value()
				var password ast.Node
				if pValue.STRING_VALUE() != nil {
					password = &ast.String{Str: stripQuotes(pValue.STRING_VALUE().GetText())}
				} else {
					password = &ast.Null{}
				}
				return &ast.DefElem{
					Defname:  &name,
					Arg:      password,
					Location: c.pos(pOpt.GetStart()),
				}
			}
		} else if hOpt := aOpt.Hash_option(); hOpt != nil {
			if debug.Active {
				log.Printf("YDB does not currently support HASH in CREATE USER statement")
			}
			var pass string
			if hOpt.HASH() != nil && hOpt.STRING_VALUE() != nil {
				pass = stripQuotes(hOpt.STRING_VALUE().GetText())
			}
			name := "hash"
			return &ast.DefElem{
				Defname:  &name,
				Arg:      &ast.String{Str: pass},
				Location: c.pos(hOpt.GetStart()),
			}
		}

	case n.Login_option() != nil:
		lOpt := n.Login_option()
		var name string
		if lOpt.LOGIN() != nil {
			name = "login"
		} else if lOpt.NOLOGIN() != nil {
			name = "nologin"
		}
		return &ast.DefElem{
			Defname:  &name,
			Arg:      &ast.Boolean{Boolval: lOpt.LOGIN() != nil},
			Location: c.pos(lOpt.GetStart()),
		}
	default:
		return todo("convertUser_optionContext", n)
	}
	return nil
}

func (c *cc) convertRole_nameContext(n *parser.Role_nameContext) ast.Node {
	switch {
	case n.An_id_or_type() != nil:
		name := parseAnIdOrType(n.An_id_or_type())
		return &ast.A_Const{Val: NewIdentifier(name), Location: c.pos(n.GetStart())}
	case n.Bind_parameter() != nil:
		bindPar := c.convert(n.Bind_parameter())
		return bindPar
	}
	return todo("convertRole_nameContext", n)
}

func (c *cc) convertCommit_stmtContext(n *parser.Commit_stmtContext) ast.Node {
	if n.COMMIT() != nil {
		return &ast.TransactionStmt{Kind: ast.TransactionStmtKind(3)}
	}
	return todo("convertCommit_stmtContext", n)
}

func (c *cc) convertRollback_stmtContext(n *parser.Rollback_stmtContext) ast.Node {
	if n.ROLLBACK() != nil {
		return &ast.TransactionStmt{Kind: ast.TransactionStmtKind(4)}
	}
	return todo("convertRollback_stmtContext", n)
}

func (c *cc) convertDrop_table_stmtContext(n *parser.Drop_table_stmtContext) ast.Node {
	if n.DROP() != nil && (n.TABLESTORE() != nil || (n.EXTERNAL() != nil && n.TABLE() != nil) || n.TABLE() != nil) {
		name := parseTableName(n.Simple_table_ref().Simple_table_ref_core())
		stmt := &ast.DropTableStmt{
			IfExists: n.IF() != nil && n.EXISTS() != nil,
			Tables:   []*ast.TableName{name},
		}
		return stmt
	}
	return todo("convertDrop_Table_stmtContxt", n)
}

func (c *cc) convertDelete_stmtContext(n *parser.Delete_stmtContext) ast.Node {
	batch := n.BATCH() != nil

	tableName := identifier(n.Simple_table_ref().Simple_table_ref_core().GetText())
	rel := &ast.RangeVar{Relname: &tableName}

	var where ast.Node
	if n.WHERE() != nil && n.Expr() != nil {
		where = c.convert(n.Expr())
	}
	var cols *ast.List
	var source ast.Node
	if n.ON() != nil && n.Into_values_source() != nil {
		nVal := n.Into_values_source()
		// todo: handle default values when implemented
		if pureCols := nVal.Pure_column_list(); pureCols != nil {
			cols = &ast.List{}
			for _, anID := range pureCols.AllAn_id() {
				name := identifier(parseAnId(anID))
				cols.Items = append(cols.Items, &ast.ResTarget{
					Name:     &name,
					Location: c.pos(anID.GetStart()),
				})
			}
		}

		valSource := nVal.Values_source()
		if valSource != nil {
			switch {
			case valSource.Values_stmt() != nil:
				source = &ast.SelectStmt{
					ValuesLists: c.convert(valSource.Values_stmt()).(*ast.List),
					FromClause:  &ast.List{},
					TargetList:  &ast.List{},
				}

			case valSource.Select_stmt() != nil:
				source = c.convert(valSource.Select_stmt())
			}
		}
	}

	returning := &ast.List{}
	if ret := n.Returning_columns_list(); ret != nil {
		returning = c.convert(ret).(*ast.List)
	}

	stmts := &ast.DeleteStmt{
		Relations:     &ast.List{Items: []ast.Node{rel}},
		WhereClause:   where,
		ReturningList: returning,
		Batch:         batch,
		OnCols:        cols,
		OnSelectStmt:  source,
	}

	return stmts
}

func (c *cc) convertPragma_stmtContext(n *parser.Pragma_stmtContext) ast.Node {
	if n.PRAGMA() != nil && n.An_id() != nil {
		prefix := ""
		if p := n.Opt_id_prefix_or_type(); p != nil {
			prefix = parseAnIdOrType(p.An_id_or_type())
		}
		items := []ast.Node{}
		if prefix != "" {
			items = append(items, &ast.A_Const{Val: NewIdentifier(prefix)})
		}

		name := parseAnId(n.An_id())
		items = append(items, &ast.A_Const{Val: NewIdentifier(name)})

		stmt := &ast.Pragma_stmt{
			Name:     &ast.List{Items: items},
			Location: c.pos(n.An_id().GetStart()),
		}

		if n.EQUALS() != nil {
			stmt.Equals = true
			if val := n.Pragma_value(0); val != nil {
				stmt.Values = &ast.List{Items: []ast.Node{c.convert(val)}}
			}
		} else if lp := n.LPAREN(); lp != nil {
			values := []ast.Node{}
			for _, v := range n.AllPragma_value() {
				values = append(values, c.convert(v))
			}
			stmt.Values = &ast.List{Items: values}
		}

		return stmt
	}
	return todo("convertPragma_stmtContext", n)
}

func (c *cc) convertPragma_valueContext(n *parser.Pragma_valueContext) ast.Node {
	switch {
	case n.Signed_number() != nil:
		if n.Signed_number().Integer() != nil {
			text := n.Signed_number().GetText()
			val, err := parseIntegerValue(text)
			if err != nil {
				if debug.Active {
					log.Printf("Failed to parse integer value '%s': %v", text, err)
				}
				return &ast.TODO{}
			}
			return &ast.A_Const{Val: &ast.Integer{Ival: val}, Location: c.pos(n.GetStart())}
		}
		if n.Signed_number().Real_() != nil {
			text := n.Signed_number().GetText()
			return &ast.A_Const{Val: &ast.Float{Str: text}, Location: c.pos(n.GetStart())}
		}
	case n.STRING_VALUE() != nil:
		val := n.STRING_VALUE().GetText()
		if len(val) >= 2 {
			val = val[1 : len(val)-1]
		}
		return &ast.A_Const{Val: &ast.String{Str: val}, Location: c.pos(n.GetStart())}
	case n.Bool_value() != nil:
		var i bool
		if n.Bool_value().TRUE() != nil {
			i = true
		}
		return &ast.A_Const{Val: &ast.Boolean{Boolval: i}, Location: c.pos(n.GetStart())}
	case n.Bind_parameter() != nil:
		bindPar := c.convert(n.Bind_parameter())
		return bindPar
	}

	return todo("convertPragma_valueContext", n)
}

func (c *cc) convertUpdate_stmtContext(n *parser.Update_stmtContext) ast.Node {
	if n.UPDATE() == nil {
		return nil
	}
	batch := n.BATCH() != nil

	tableName := identifier(n.Simple_table_ref().Simple_table_ref_core().GetText())
	rel := &ast.RangeVar{Relname: &tableName}

	var where ast.Node
	var setList *ast.List
	var cols *ast.List
	var source ast.Node

	if n.SET() != nil && n.Set_clause_choice() != nil {
		nSet := n.Set_clause_choice()
		setList = &ast.List{Items: []ast.Node{}}

		switch {
		case nSet.Set_clause_list() != nil:
			for _, clause := range nSet.Set_clause_list().AllSet_clause() {
				targetCtx := clause.Set_target()
				columnName := identifier(targetCtx.Column_name().GetText())
				expr := c.convert(clause.Expr())
				resTarget := &ast.ResTarget{
					Name:     &columnName,
					Val:      expr,
					Location: c.pos(clause.Expr().GetStart()),
				}
				setList.Items = append(setList.Items, resTarget)
			}

		case nSet.Multiple_column_assignment() != nil:
			multiAssign := nSet.Multiple_column_assignment()
			targetsCtx := multiAssign.Set_target_list()
			valuesCtx := multiAssign.Simple_values_source()

			var colNames []string
			for _, target := range targetsCtx.AllSet_target() {
				targetCtx := target.(*parser.Set_targetContext)
				colNames = append(colNames, targetCtx.Column_name().GetText())
			}

			var rowExpr *ast.RowExpr
			if exprList := valuesCtx.Expr_list(); exprList != nil {
				rowExpr = &ast.RowExpr{
					Args: &ast.List{},
				}
				for _, expr := range exprList.AllExpr() {
					rowExpr.Args.Items = append(rowExpr.Args.Items, c.convert(expr))
				}
			}

			for i, colName := range colNames {
				name := identifier(colName)
				setList.Items = append(setList.Items, &ast.ResTarget{
					Name: &name,
					Val: &ast.MultiAssignRef{
						Source:   rowExpr,
						Colno:    i + 1,
						Ncolumns: len(colNames),
					},
					Location: c.pos(targetsCtx.Set_target(i).GetStart()),
				})
			}
		}

		if n.WHERE() != nil && n.Expr() != nil {
			where = c.convert(n.Expr())
		}
	} else if n.ON() != nil && n.Into_values_source() != nil {

		// todo: handle default values when implemented

		nVal := n.Into_values_source()

		if pureCols := nVal.Pure_column_list(); pureCols != nil {
			cols = &ast.List{}
			for _, anID := range pureCols.AllAn_id() {
				name := identifier(parseAnId(anID))
				cols.Items = append(cols.Items, &ast.ResTarget{
					Name:     &name,
					Location: c.pos(anID.GetStart()),
				})
			}
		}

		valSource := nVal.Values_source()
		if valSource != nil {
			switch {
			case valSource.Values_stmt() != nil:
				source = &ast.SelectStmt{
					ValuesLists: c.convert(valSource.Values_stmt()).(*ast.List),
					FromClause:  &ast.List{},
					TargetList:  &ast.List{},
				}

			case valSource.Select_stmt() != nil:
				source = c.convert(valSource.Select_stmt())
			}
		}
	}

	returning := &ast.List{}
	if ret := n.Returning_columns_list(); ret != nil {
		returning = c.convert(ret).(*ast.List)
	}

	stmts := &ast.UpdateStmt{
		Relations:     &ast.List{Items: []ast.Node{rel}},
		TargetList:    setList,
		WhereClause:   where,
		ReturningList: returning,
		FromClause:    &ast.List{},
		WithClause:    nil,
		Batch:         batch,
		OnCols:        cols,
		OnSelectStmt:  source,
	}

	return stmts
}

func (c *cc) convertInto_table_stmtContext(n *parser.Into_table_stmtContext) ast.Node {
	tableName := identifier(n.Into_simple_table_ref().Simple_table_ref().Simple_table_ref_core().GetText())
	rel := &ast.RangeVar{
		Relname:  &tableName,
		Location: c.pos(n.Into_simple_table_ref().GetStart()),
	}

	onConflict := &ast.OnConflictClause{}
	switch {
	case n.INSERT() != nil && n.OR() != nil && n.ABORT() != nil:
		onConflict.Action = ast.OnConflictAction_INSERT_OR_ABORT
	case n.INSERT() != nil && n.OR() != nil && n.REVERT() != nil:
		onConflict.Action = ast.OnConflictAction_INSERT_OR_REVERT
	case n.INSERT() != nil && n.OR() != nil && n.IGNORE() != nil:
		onConflict.Action = ast.OnConflictAction_INSERT_OR_IGNORE
	case n.UPSERT() != nil:
		onConflict.Action = ast.OnConflictAction_UPSERT
	case n.REPLACE() != nil:
		onConflict.Action = ast.OnConflictAction_REPLACE
	}

	var cols *ast.List
	var source ast.Node
	if nVal := n.Into_values_source(); nVal != nil {
		// todo: handle default values when implemented

		if pureCols := nVal.Pure_column_list(); pureCols != nil {
			cols = &ast.List{}
			for _, anID := range pureCols.AllAn_id() {
				name := identifier(parseAnId(anID))
				cols.Items = append(cols.Items, &ast.ResTarget{
					Name:     &name,
					Location: c.pos(anID.GetStart()),
				})
			}
		}

		valSource := nVal.Values_source()
		if valSource != nil {
			switch {
			case valSource.Values_stmt() != nil:
				source = &ast.SelectStmt{
					ValuesLists: c.convert(valSource.Values_stmt()).(*ast.List),
					FromClause:  &ast.List{},
					TargetList:  &ast.List{},
				}

			case valSource.Select_stmt() != nil:
				source = c.convert(valSource.Select_stmt())
			}
		}
	}

	returning := &ast.List{}
	if ret := n.Returning_columns_list(); ret != nil {
		returning = c.convert(ret).(*ast.List)
	}

	stmts := &ast.InsertStmt{
		Relation:         rel,
		Cols:             cols,
		SelectStmt:       source,
		OnConflictClause: onConflict,
		ReturningList:    returning,
	}

	return stmts
}

func (c *cc) convertValues_stmtContext(n *parser.Values_stmtContext) ast.Node {
	mainList := &ast.List{}

	for _, rowCtx := range n.Values_source_row_list().AllValues_source_row() {
		rowList := &ast.List{}
		exprListCtx := rowCtx.Expr_list().(*parser.Expr_listContext)

		for _, exprCtx := range exprListCtx.AllExpr() {
			if converted := c.convert(exprCtx); converted != nil {
				rowList.Items = append(rowList.Items, converted)
			}
		}

		mainList.Items = append(mainList.Items, rowList)

	}

	return mainList
}

func (c *cc) convertReturning_columns_listContext(n *parser.Returning_columns_listContext) ast.Node {
	list := &ast.List{Items: []ast.Node{}}

	if n.ASTERISK() != nil {
		target := &ast.ResTarget{
			Indirection: &ast.List{},
			Val: &ast.ColumnRef{
				Fields:   &ast.List{Items: []ast.Node{&ast.A_Star{}}},
				Location: c.pos(n.ASTERISK().GetSymbol()),
			},
			Location: c.pos(n.ASTERISK().GetSymbol()),
		}
		list.Items = append(list.Items, target)
		return list
	}

	for _, idCtx := range n.AllAn_id() {
		target := &ast.ResTarget{
			Indirection: &ast.List{},
			Val: &ast.ColumnRef{
				Fields: &ast.List{
					Items: []ast.Node{NewIdentifier(parseAnId(idCtx))},
				},
				Location: c.pos(idCtx.GetStart()),
			},
			Location: c.pos(idCtx.GetStart()),
		}
		list.Items = append(list.Items, target)
	}

	return list
}

func (c *cc) convertSelectStmtContext(n *parser.Select_stmtContext) ast.Node {
	skp := n.Select_kind_parenthesis(0)
	if skp == nil {
		return nil
	}
	partial := skp.Select_kind_partial()
	if partial == nil {
		return nil
	}
	sk := partial.Select_kind()
	if sk == nil {
		return nil
	}
	selectStmt := &ast.SelectStmt{}

	switch {
	case sk.Process_core() != nil:
		cnode := c.convert(sk.Process_core())
		stmt, ok := cnode.(*ast.SelectStmt)
		if !ok {
			return nil
		}
		selectStmt = stmt
	case sk.Select_core() != nil:
		cnode := c.convert(sk.Select_core())
		stmt, ok := cnode.(*ast.SelectStmt)
		if !ok {
			return nil
		}
		selectStmt = stmt
	case sk.Reduce_core() != nil:
		cnode := c.convert(sk.Reduce_core())
		stmt, ok := cnode.(*ast.SelectStmt)
		if !ok {
			return nil
		}
		selectStmt = stmt
	}

	// todo: cover process and reduce core,
	// todo: cover LIMIT and OFFSET

	return selectStmt
}

func (c *cc) convertSelectCoreContext(n *parser.Select_coreContext) ast.Node {
	stmt := &ast.SelectStmt{
		TargetList: &ast.List{},
		FromClause: &ast.List{},
	}
	if n.Opt_set_quantifier() != nil {
		oq := n.Opt_set_quantifier()
		if oq.DISTINCT() != nil {
			// todo: add distinct support
			stmt.DistinctClause = &ast.List{}
		}
	}
	resultCols := n.AllResult_column()
	if len(resultCols) > 0 {
		var items []ast.Node
		for _, rc := range resultCols {
			resCol, ok := rc.(*parser.Result_columnContext)
			if !ok {
				continue
			}
			convNode := c.convertResultColumn(resCol)
			if convNode != nil {
				items = append(items, convNode)
			}
		}
		stmt.TargetList = &ast.List{
			Items: items,
		}
	}
	jsList := n.AllJoin_source()
	if len(n.AllFROM()) > 0 && len(jsList) > 0 {
		var fromItems []ast.Node
		for _, js := range jsList {
			jsCon, ok := js.(*parser.Join_sourceContext)
			if !ok {
				continue
			}

			joinNode := c.convertJoinSource(jsCon)
			if joinNode != nil {
				fromItems = append(fromItems, joinNode)
			}
		}
		stmt.FromClause = &ast.List{
			Items: fromItems,
		}
	}
	if n.WHERE() != nil {
		whereCtx := n.Expr(0)
		if whereCtx != nil {
			stmt.WhereClause = c.convert(whereCtx)
		}
	}
	return stmt
}

func (c *cc) convertResultColumn(n *parser.Result_columnContext) ast.Node {
	// todo: support opt_id_prefix
	target := &ast.ResTarget{
		Location: c.pos(n.GetStart()),
	}
	var val ast.Node
	iexpr := n.Expr()
	switch {
	case n.ASTERISK() != nil:
		val = c.convertWildCardField(n)
	case iexpr != nil:
		val = c.convert(iexpr)
	}

	if val == nil {
		return nil
	}
	switch {
	case n.AS() != nil && n.An_id_or_type() != nil:
		name := parseAnIdOrType(n.An_id_or_type())
		target.Name = &name
	case n.An_id_as_compat() != nil: //nolint
		// todo: parse as_compat
	}
	target.Val = val
	return target
}

func (c *cc) convertJoinSource(n *parser.Join_sourceContext) ast.Node {
	if n == nil {
		return nil
	}
	fsList := n.AllFlatten_source()
	if len(fsList) == 0 {
		return nil
	}
	joinOps := n.AllJoin_op()
	joinConstraints := n.AllJoin_constraint()

	// todo: add ANY support

	leftNode := c.convertFlattenSource(fsList[0])
	if leftNode == nil {
		return nil
	}
	for i, jopCtx := range joinOps {
		if i+1 >= len(fsList) {
			break
		}
		rightNode := c.convertFlattenSource(fsList[i+1])
		if rightNode == nil {
			return leftNode
		}
		jexpr := &ast.JoinExpr{
			Larg: leftNode,
			Rarg: rightNode,
		}
		if jopCtx.NATURAL() != nil {
			jexpr.IsNatural = true
		}
		// todo: cover semi/only/exclusion/
		switch {
		case jopCtx.LEFT() != nil:
			jexpr.Jointype = ast.JoinTypeLeft
		case jopCtx.RIGHT() != nil:
			jexpr.Jointype = ast.JoinTypeRight
		case jopCtx.FULL() != nil:
			jexpr.Jointype = ast.JoinTypeFull
		case jopCtx.INNER() != nil:
			jexpr.Jointype = ast.JoinTypeInner
		case jopCtx.COMMA() != nil:
			jexpr.Jointype = ast.JoinTypeInner
		default:
			jexpr.Jointype = ast.JoinTypeInner
		}
		if i < len(joinConstraints) {
			if jc := joinConstraints[i]; jc != nil {
				switch {
				case jc.ON() != nil:
					if exprCtx := jc.Expr(); exprCtx != nil {
						jexpr.Quals = c.convert(exprCtx)
					}
				case jc.USING() != nil:
					if pureListCtx := jc.Pure_column_or_named_list(); pureListCtx != nil {
						var using ast.List
						pureItems := pureListCtx.AllPure_column_or_named()
						for _, pureCtx := range pureItems {
							if anID := pureCtx.An_id(); anID != nil {
								using.Items = append(using.Items, NewIdentifier(parseAnId(anID)))
							} else if bp := pureCtx.Bind_parameter(); bp != nil {
								bindPar := c.convert(bp)
								using.Items = append(using.Items, bindPar)
							}
						}
						jexpr.UsingClause = &using
					}
				}
			}
		}
		leftNode = jexpr
	}
	return leftNode
}

func (c *cc) convertFlattenSource(n parser.IFlatten_sourceContext) ast.Node {
	if n == nil {
		return nil
	}
	nss := n.Named_single_source()
	if nss == nil {
		return nil
	}
	namedSingleSource, ok := nss.(*parser.Named_single_sourceContext)
	if !ok {
		return nil
	}
	return c.convertNamedSingleSource(namedSingleSource)
}

func (c *cc) convertNamedSingleSource(n *parser.Named_single_sourceContext) ast.Node {
	ss := n.Single_source()
	if ss == nil {
		return nil
	}
	SingleSource, ok := ss.(*parser.Single_sourceContext)
	if !ok {
		return nil
	}
	base := c.convertSingleSource(SingleSource)

	if n.AS() != nil && n.An_id() != nil {
		aliasText := parseAnId(n.An_id())
		switch source := base.(type) {
		case *ast.RangeVar:
			source.Alias = &ast.Alias{Aliasname: &aliasText}
		case *ast.RangeSubselect:
			source.Alias = &ast.Alias{Aliasname: &aliasText}
		}
	} else if n.An_id_as_compat() != nil { //nolint
		// todo: parse as_compat
	}
	return base
}

func (c *cc) convertSingleSource(n *parser.Single_sourceContext) ast.Node {
	if n.Table_ref() != nil {
		tableName := n.Table_ref().GetText() // !! debug !!
		return &ast.RangeVar{
			Relname:  &tableName,
			Location: c.pos(n.GetStart()),
		}
	}

	if n.Select_stmt() != nil {
		subquery := c.convert(n.Select_stmt())
		return &ast.RangeSubselect{
			Subquery: subquery,
		}

	}
	// todo: Values stmt

	return nil
}

func (c *cc) convertBindParameter(n *parser.Bind_parameterContext) ast.Node {
	// !!debug later!!
	if n.DOLLAR() != nil {
		if n.TRUE() != nil {
			return &ast.A_Const{Val: &ast.Boolean{Boolval: true}, Location: c.pos(n.GetStart())}
		}
		if n.FALSE() != nil {
			return &ast.A_Const{Val: &ast.Boolean{Boolval: false}, Location: c.pos(n.GetStart())}
		}

		if an := n.An_id_or_type(); an != nil {
			idText := parseAnIdOrType(an)
			return &ast.A_Expr{
				Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "@"}}},
				Rexpr:    &ast.String{Str: idText},
				Location: c.pos(n.GetStart()),
			}
		}
		c.paramCount++
		return &ast.ParamRef{
			Number:   c.paramCount,
			Location: c.pos(n.GetStart()),
			Dollar:   true,
		}
	}
	return &ast.TODO{}
}

func (c *cc) convertWildCardField(n *parser.Result_columnContext) *ast.ColumnRef {
	prefixCtx := n.Opt_id_prefix()
	prefix := c.convertOptIdPrefix(prefixCtx)

	items := []ast.Node{}
	if prefix != "" {
		items = append(items, NewIdentifier(prefix))
	}

	items = append(items, &ast.A_Star{})
	return &ast.ColumnRef{
		Fields:   &ast.List{Items: items},
		Location: c.pos(n.GetStart()),
	}
}

func (c *cc) convertOptIdPrefix(ctx parser.IOpt_id_prefixContext) string {
	if ctx == nil {
		return ""
	}
	if ctx.An_id() != nil {
		return ctx.An_id().GetText()
	}
	return ""
}

func (c *cc) convertCreate_table_stmtContext(n *parser.Create_table_stmtContext) ast.Node {
	stmt := &ast.CreateTableStmt{
		Name:        parseTableName(n.Simple_table_ref().Simple_table_ref_core()),
		IfNotExists: n.EXISTS() != nil,
	}
	for _, idef := range n.AllCreate_table_entry() {
		if def, ok := idef.(*parser.Create_table_entryContext); ok {
			switch {
			case def.Column_schema() != nil:
				if colCtx, ok := def.Column_schema().(*parser.Column_schemaContext); ok {
					colDef := c.convertColumnSchema(colCtx)
					if colDef != nil {
						stmt.Cols = append(stmt.Cols, colDef)
					}
				}
			case def.Table_constraint() != nil:
				if conCtx, ok := def.Table_constraint().(*parser.Table_constraintContext); ok {
					switch {
					case conCtx.PRIMARY() != nil && conCtx.KEY() != nil:
						for _, cname := range conCtx.AllAn_id() {
							for _, col := range stmt.Cols {
								if col.Colname == parseAnId(cname) {
									col.IsNotNull = true
								}
							}
						}
					case conCtx.PARTITION() != nil && conCtx.BY() != nil:
						_ = conCtx
						// todo: partition by constraint
					case conCtx.ORDER() != nil && conCtx.BY() != nil:
						_ = conCtx
						// todo: order by constraint
					}
				}

			case def.Table_index() != nil:
				if indCtx, ok := def.Table_index().(*parser.Table_indexContext); ok {
					_ = indCtx
					// todo
				}
			case def.Family_entry() != nil:
				if famCtx, ok := def.Family_entry().(*parser.Family_entryContext); ok {
					_ = famCtx
					// todo
				}
			case def.Changefeed() != nil: // таблица ориентированная
				if cgfCtx, ok := def.Changefeed().(*parser.ChangefeedContext); ok {
					_ = cgfCtx
					// todo
				}
			}
		}
	}
	return stmt
}

func (c *cc) convertColumnSchema(n *parser.Column_schemaContext) *ast.ColumnDef {

	col := &ast.ColumnDef{}

	if anId := n.An_id_schema(); anId != nil {
		col.Colname = identifier(parseAnIdSchema(anId))
	}
	if tnb := n.Type_name_or_bind(); tnb != nil {
		col.TypeName = c.convertTypeNameOrBind(tnb)
	}
	if colCons := n.Opt_column_constraints(); colCons != nil {
		col.IsNotNull = colCons.NOT() != nil && colCons.NULL() != nil
		//todo: cover exprs if needed
	}
	// todo: family

	return col
}

func (c *cc) convertTypeNameOrBind(n parser.IType_name_or_bindContext) *ast.TypeName {
	if t := n.Type_name(); t != nil {
		return c.convertTypeName(t)
	} else if b := n.Bind_parameter(); b != nil {
		return &ast.TypeName{Name: "BIND:" + identifier(parseAnIdOrType(b.An_id_or_type()))}
	}
	return nil
}

func (c *cc) convertTypeName(n parser.IType_nameContext) *ast.TypeName {
	if n == nil {
		return nil
	}

	if composite := n.Type_name_composite(); composite != nil {
		if node := c.convertTypeNameComposite(composite); node != nil {
			if typeName, ok := node.(*ast.TypeName); ok {
				return typeName
			}
		}
	}

	if decimal := n.Type_name_decimal(); decimal != nil {
		if integerOrBinds := decimal.AllInteger_or_bind(); len(integerOrBinds) >= 2 {
			return &ast.TypeName{
				Name:    "Decimal",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{
						c.convertIntegerOrBind(integerOrBinds[0]),
						c.convertIntegerOrBind(integerOrBinds[1]),
					},
				},
			}
		}
	}

	// Handle simple types
	if simple := n.Type_name_simple(); simple != nil {
		return &ast.TypeName{
			Name:    simple.GetText(),
			TypeOid: 0,
		}
	}

	return nil
}

func (c *cc) convertIntegerOrBind(n parser.IInteger_or_bindContext) ast.Node {
	if n == nil {
		return nil
	}

	if integer := n.Integer(); integer != nil {
		val, err := parseIntegerValue(integer.GetText())
		if err != nil {
			return &ast.TODO{}
		}
		return &ast.Integer{Ival: val}
	}

	if bind := n.Bind_parameter(); bind != nil {
		return c.convertBindParameter(bind.(*parser.Bind_parameterContext))
	}

	return nil
}

func (c *cc) convertTypeNameComposite(n parser.IType_name_compositeContext) ast.Node {
	if n == nil {
		return nil
	}

	if opt := n.Type_name_optional(); opt != nil {
		if typeName := opt.Type_name_or_bind(); typeName != nil {
			return &ast.TypeName{
				Name:    "Optional",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{c.convertTypeNameOrBind(typeName)},
				},
			}
		}
	}

	if tuple := n.Type_name_tuple(); tuple != nil {
		if typeNames := tuple.AllType_name_or_bind(); len(typeNames) > 0 {
			var items []ast.Node
			for _, tn := range typeNames {
				items = append(items, c.convertTypeNameOrBind(tn))
			}
			return &ast.TypeName{
				Name:    "Tuple",
				TypeOid: 0,
				Names:   &ast.List{Items: items},
			}
		}
	}

	if struct_ := n.Type_name_struct(); struct_ != nil {
		if structArgs := struct_.AllStruct_arg(); len(structArgs) > 0 {
			var items []ast.Node
			for range structArgs {
				// TODO: Handle struct field names and types
				items = append(items, &ast.TODO{})
			}
			return &ast.TypeName{
				Name:    "Struct",
				TypeOid: 0,
				Names:   &ast.List{Items: items},
			}
		}
	}

	if variant := n.Type_name_variant(); variant != nil {
		if variantArgs := variant.AllVariant_arg(); len(variantArgs) > 0 {
			var items []ast.Node
			for range variantArgs {
				// TODO: Handle variant arguments
				items = append(items, &ast.TODO{})
			}
			return &ast.TypeName{
				Name:    "Variant",
				TypeOid: 0,
				Names:   &ast.List{Items: items},
			}
		}
	}

	if list := n.Type_name_list(); list != nil {
		if typeName := list.Type_name_or_bind(); typeName != nil {
			return &ast.TypeName{
				Name:    "List",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{c.convertTypeNameOrBind(typeName)},
				},
			}
		}
	}

	if stream := n.Type_name_stream(); stream != nil {
		if typeName := stream.Type_name_or_bind(); typeName != nil {
			return &ast.TypeName{
				Name:    "Stream",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{c.convertTypeNameOrBind(typeName)},
				},
			}
		}
	}

	if flow := n.Type_name_flow(); flow != nil {
		if typeName := flow.Type_name_or_bind(); typeName != nil {
			return &ast.TypeName{
				Name:    "Flow",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{c.convertTypeNameOrBind(typeName)},
				},
			}
		}
	}

	if dict := n.Type_name_dict(); dict != nil {
		if typeNames := dict.AllType_name_or_bind(); len(typeNames) >= 2 {
			return &ast.TypeName{
				Name:    "Dict",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{
						c.convertTypeNameOrBind(typeNames[0]),
						c.convertTypeNameOrBind(typeNames[1]),
					},
				},
			}
		}
	}

	if set := n.Type_name_set(); set != nil {
		if typeName := set.Type_name_or_bind(); typeName != nil {
			return &ast.TypeName{
				Name:    "Set",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{c.convertTypeNameOrBind(typeName)},
				},
			}
		}
	}

	if enum := n.Type_name_enum(); enum != nil {
		if typeTags := enum.AllType_name_tag(); len(typeTags) > 0 {
			var items []ast.Node
			for range typeTags { // todo: Handle enum tags
				items = append(items, &ast.TODO{})
			}
			return &ast.TypeName{
				Name:    "Enum",
				TypeOid: 0,
				Names:   &ast.List{Items: items},
			}
		}
	}

	if resource := n.Type_name_resource(); resource != nil {
		if typeTag := resource.Type_name_tag(); typeTag != nil {
			// TODO: Handle resource tag
			return &ast.TypeName{
				Name:    "Resource",
				TypeOid: 0,
				Names: &ast.List{
					Items: []ast.Node{&ast.TODO{}},
				},
			}
		}
	}

	if tagged := n.Type_name_tagged(); tagged != nil {
		if typeName := tagged.Type_name_or_bind(); typeName != nil {
			if typeTag := tagged.Type_name_tag(); typeTag != nil {
				// TODO: Handle tagged type and tag
				return &ast.TypeName{
					Name:    "Tagged",
					TypeOid: 0,
					Names: &ast.List{
						Items: []ast.Node{
							c.convertTypeNameOrBind(typeName),
							&ast.TODO{},
						},
					},
				}
			}
		}
	}

	if callable := n.Type_name_callable(); callable != nil {
		// TODO: Handle callable argument list and return type
		return &ast.TypeName{
			Name:    "Callable",
			TypeOid: 0,
			Names: &ast.List{
				Items: []ast.Node{&ast.TODO{}},
			},
		}
	}

	return nil
}

func (c *cc) convertSqlStmtCore(n parser.ISql_stmt_coreContext) ast.Node {
	if n == nil {
		return nil
	}

	if stmt := n.Pragma_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Select_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Named_nodes_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Use_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Into_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Commit_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Update_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Delete_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Rollback_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Declare_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Import_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Export_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_external_table_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Do_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Define_action_or_subquery_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.If_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.For_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Values_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_user_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_user_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_group_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_group_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_role_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_object_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_object_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_object_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_external_data_source_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_external_data_source_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_external_data_source_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_replication_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_replication_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_topic_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_topic_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_topic_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Grant_permissions_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Revoke_permissions_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_table_store_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Upsert_object_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_view_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_view_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_replication_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_resource_pool_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_resource_pool_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_resource_pool_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_backup_collection_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_backup_collection_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_backup_collection_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Analyze_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Create_resource_pool_classifier_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_resource_pool_classifier_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Drop_resource_pool_classifier_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Backup_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Restore_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	if stmt := n.Alter_sequence_stmt(); stmt != nil {
		return c.convert(stmt)
	}
	return nil
}

func (c *cc) convertExpr(n *parser.ExprContext) ast.Node {
	if n == nil {
		return nil
	}

	if tn := n.Type_name_composite(); tn != nil {
		return c.convertTypeNameComposite(tn)
	}

	orSubs := n.AllOr_subexpr()
	if len(orSubs) == 0 {
		return nil
	}

	orSub, ok := orSubs[0].(*parser.Or_subexprContext)
	if !ok {
		return nil
	}

	left := c.convertOrSubExpr(orSub)
	for i := 1; i < len(orSubs); i++ {
		orSub, ok = orSubs[i].(*parser.Or_subexprContext)
		if !ok {
			return nil
		}
		right := c.convertOrSubExpr(orSub)
		left = &ast.BoolExpr{
			Boolop:   ast.BoolExprTypeOr,
			Args:     &ast.List{Items: []ast.Node{left, right}},
			Location: c.pos(n.GetStart()),
		}
	}
	return left
}

func (c *cc) convertOrSubExpr(n *parser.Or_subexprContext) ast.Node {
	if n == nil {
		return nil
	}
	andSubs := n.AllAnd_subexpr()
	if len(andSubs) == 0 {
		return nil
	}
	andSub, ok := andSubs[0].(*parser.And_subexprContext)
	if !ok {
		return nil
	}

	left := c.convertAndSubexpr(andSub)
	for i := 1; i < len(andSubs); i++ {
		andSub, ok = andSubs[i].(*parser.And_subexprContext)
		if !ok {
			return nil
		}
		right := c.convertAndSubexpr(andSub)
		left = &ast.BoolExpr{
			Boolop:   ast.BoolExprTypeAnd,
			Args:     &ast.List{Items: []ast.Node{left, right}},
			Location: c.pos(n.GetStart()),
		}
	}
	return left
}

func (c *cc) convertAndSubexpr(n *parser.And_subexprContext) ast.Node {
	if n == nil {
		return nil
	}

	xors := n.AllXor_subexpr()
	if len(xors) == 0 {
		return nil
	}

	xor, ok := xors[0].(*parser.Xor_subexprContext)
	if !ok {
		return nil
	}

	left := c.convertXorSubexpr(xor)
	for i := 1; i < len(xors); i++ {
		xor, ok = xors[i].(*parser.Xor_subexprContext)
		if !ok {
			return nil
		}
		right := c.convertXorSubexpr(xor)
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "XOR"}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}
	return left
}

func (c *cc) convertXorSubexpr(n *parser.Xor_subexprContext) ast.Node {
	if n == nil {
		return nil
	}
	es := n.Eq_subexpr()
	if es == nil {
		return nil
	}
	subExpr, ok := es.(*parser.Eq_subexprContext)
	if !ok {
		return nil
	}
	base := c.convertEqSubexpr(subExpr)
	if cond := n.Cond_expr(); cond != nil {
		condCtx, ok := cond.(*parser.Cond_exprContext)
		if !ok {
			return base
		}

		switch {
		case condCtx.IN() != nil:
			if inExpr := condCtx.In_expr(); inExpr != nil {
				return &ast.A_Expr{
					Name:  &ast.List{Items: []ast.Node{&ast.String{Str: "IN"}}},
					Lexpr: base,
					Rexpr: c.convert(inExpr),
				}
			}
		case condCtx.BETWEEN() != nil:
			if eqSubs := condCtx.AllEq_subexpr(); len(eqSubs) >= 2 {
				return &ast.BetweenExpr{
					Expr:     base,
					Left:     c.convert(eqSubs[0]),
					Right:    c.convert(eqSubs[1]),
					Not:      condCtx.NOT() != nil,
					Location: c.pos(n.GetStart()),
				}
			}
		case condCtx.ISNULL() != nil:
			return &ast.NullTest{
				Arg:          base,
				Nulltesttype: 1, // IS NULL
				Location:     c.pos(n.GetStart()),
			}
		case condCtx.NOTNULL() != nil:
			return &ast.NullTest{
				Arg:          base,
				Nulltesttype: 2, // IS NOT NULL
				Location:     c.pos(n.GetStart()),
			}
		case condCtx.IS() != nil && condCtx.NULL() != nil:
			return &ast.NullTest{
				Arg:          base,
				Nulltesttype: 1, // IS NULL
				Location:     c.pos(n.GetStart()),
			}
		case condCtx.IS() != nil && condCtx.NOT() != nil && condCtx.NULL() != nil:
			return &ast.NullTest{
				Arg:          base,
				Nulltesttype: 2, // IS NOT NULL
				Location:     c.pos(n.GetStart()),
			}
		case condCtx.Match_op() != nil:
			// debug!!!
			matchOp := condCtx.Match_op().GetText()
			if eqSubs := condCtx.AllEq_subexpr(); len(eqSubs) >= 1 {
				expr := &ast.A_Expr{
					Name:  &ast.List{Items: []ast.Node{&ast.String{Str: matchOp}}},
					Lexpr: base,
					Rexpr: c.convert(eqSubs[0]),
				}
				if condCtx.ESCAPE() != nil && len(eqSubs) >= 2 { //nolint
					// todo: Add ESCAPE support
				}
				return expr
			}
		case len(condCtx.AllEQUALS()) > 0 || len(condCtx.AllEQUALS2()) > 0 ||
			len(condCtx.AllNOT_EQUALS()) > 0 || len(condCtx.AllNOT_EQUALS2()) > 0:
			// debug!!!
			var op string
			switch {
			case len(condCtx.AllEQUALS()) > 0:
				op = "="
			case len(condCtx.AllEQUALS2()) > 0:
				op = "=="
			case len(condCtx.AllNOT_EQUALS()) > 0:
				op = "!="
			case len(condCtx.AllNOT_EQUALS2()) > 0:
				op = "<>"
			}
			if eqSubs := condCtx.AllEq_subexpr(); len(eqSubs) >= 1 {
				return &ast.A_Expr{
					Name:  &ast.List{Items: []ast.Node{&ast.String{Str: op}}},
					Lexpr: base,
					Rexpr: c.convert(eqSubs[0]),
				}
			}
		case len(condCtx.AllDistinct_from_op()) > 0:
			// debug!!!
			distinctOps := condCtx.AllDistinct_from_op()
			for _, distinctOp := range distinctOps {
				if eqSubs := condCtx.AllEq_subexpr(); len(eqSubs) >= 1 {
					not := distinctOp.NOT() != nil
					op := "IS DISTINCT FROM"
					if not {
						op = "IS NOT DISTINCT FROM"
					}
					return &ast.A_Expr{
						Name:  &ast.List{Items: []ast.Node{&ast.String{Str: op}}},
						Lexpr: base,
						Rexpr: c.convert(eqSubs[0]),
					}
				}
			}
		}
	}
	return base
}

func (c *cc) convertEqSubexpr(n *parser.Eq_subexprContext) ast.Node {
	if n == nil {
		return nil
	}
	neqList := n.AllNeq_subexpr()
	if len(neqList) == 0 {
		return nil
	}
	neq, ok := neqList[0].(*parser.Neq_subexprContext)
	if !ok {
		return nil
	}
	left := c.convertNeqSubexpr(neq)
	ops := c.collectComparisonOps(n)
	for i := 1; i < len(neqList); i++ {
		neq, ok = neqList[i].(*parser.Neq_subexprContext)
		if !ok {
			return nil
		}
		right := c.convertNeqSubexpr(neq)
		opText := ops[i-1].GetText()
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: opText}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}
	return left
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

func (c *cc) convertNeqSubexpr(n *parser.Neq_subexprContext) ast.Node {
	if n == nil {
		return nil
	}
	bitList := n.AllBit_subexpr()
	if len(bitList) == 0 {
		return nil
	}

	bl, ok := bitList[0].(*parser.Bit_subexprContext)
	if !ok {
		return nil
	}
	left := c.convertBitSubexpr(bl)
	ops := c.collectBitwiseOps(n)
	for i := 1; i < len(bitList); i++ {
		bl, ok = bitList[i].(*parser.Bit_subexprContext)
		if !ok {
			return nil
		}
		right := c.convertBitSubexpr(bl)
		opText := ops[i-1].GetText()
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: opText}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}

	if n.Double_question() != nil {
		nextCtx := n.Neq_subexpr()
		if nextCtx != nil {
			neq, ok2 := nextCtx.(*parser.Neq_subexprContext)
			if !ok2 {
				return nil
			}
			right := c.convertNeqSubexpr(neq)
			left = &ast.A_Expr{
				Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "??"}}},
				Lexpr:    left,
				Rexpr:    right,
				Location: c.pos(n.GetStart()),
			}
		}
	} else {
		// !! debug !!
		qCount := len(n.AllQUESTION())
		if qCount > 0 {
			questionOp := "?"
			if qCount > 1 {
				questionOp = strings.Repeat("?", qCount)
			}
			left = &ast.A_Expr{
				Name:     &ast.List{Items: []ast.Node{&ast.String{Str: questionOp}}},
				Lexpr:    left,
				Location: c.pos(n.GetStart()),
			}
		}
	}

	return left
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

func (c *cc) convertBitSubexpr(n *parser.Bit_subexprContext) ast.Node {
	addList := n.AllAdd_subexpr()
	left := c.convertAddSubexpr(addList[0].(*parser.Add_subexprContext))

	ops := c.collectBitOps(n)
	for i := 1; i < len(addList); i++ {
		right := c.convertAddSubexpr(addList[i].(*parser.Add_subexprContext))
		opText := ops[i-1].GetText()
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: opText}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}
	return left
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

func (c *cc) convertAddSubexpr(n *parser.Add_subexprContext) ast.Node {
	mulList := n.AllMul_subexpr()
	left := c.convertMulSubexpr(mulList[0].(*parser.Mul_subexprContext))

	ops := c.collectAddOps(n)
	for i := 1; i < len(mulList); i++ {
		right := c.convertMulSubexpr(mulList[i].(*parser.Mul_subexprContext))
		opText := ops[i-1].GetText()
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: opText}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}
	return left
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

func (c *cc) convertMulSubexpr(n *parser.Mul_subexprContext) ast.Node {
	conList := n.AllCon_subexpr()
	left := c.convertConSubexpr(conList[0].(*parser.Con_subexprContext))

	for i := 1; i < len(conList); i++ {
		right := c.convertConSubexpr(conList[i].(*parser.Con_subexprContext))
		left = &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: "||"}}},
			Lexpr:    left,
			Rexpr:    right,
			Location: c.pos(n.GetStart()),
		}
	}
	return left
}

func (c *cc) convertConSubexpr(n *parser.Con_subexprContext) ast.Node {
	if opCtx := n.Unary_op(); opCtx != nil {
		op := opCtx.GetText()
		operand := c.convertUnarySubexpr(n.Unary_subexpr().(*parser.Unary_subexprContext))
		return &ast.A_Expr{
			Name:     &ast.List{Items: []ast.Node{&ast.String{Str: op}}},
			Rexpr:    operand,
			Location: c.pos(n.GetStart()),
		}
	}
	return c.convertUnarySubexpr(n.Unary_subexpr().(*parser.Unary_subexprContext))
}

func (c *cc) convertUnarySubexpr(n *parser.Unary_subexprContext) ast.Node {
	if casual := n.Unary_casual_subexpr(); casual != nil {
		return c.convertUnaryCasualSubexpr(casual.(*parser.Unary_casual_subexprContext))
	}
	if jsonExpr := n.Json_api_expr(); jsonExpr != nil {
		return c.convertJsonApiExpr(jsonExpr.(*parser.Json_api_exprContext))
	}
	return nil
}

func (c *cc) convertJsonApiExpr(n *parser.Json_api_exprContext) ast.Node {
	return todo("Json_api_exprContext", n)
}

func (c *cc) convertUnaryCasualSubexpr(n *parser.Unary_casual_subexprContext) ast.Node {
	var current ast.Node
	switch {
	case n.Id_expr() != nil:
		current = c.convertIdExpr(n.Id_expr().(*parser.Id_exprContext))
	case n.Atom_expr() != nil:
		current = c.convertAtomExpr(n.Atom_expr().(*parser.Atom_exprContext))
	default:
		return todo("Unary_casual_subexprContext", n)
	}

	if suffix := n.Unary_subexpr_suffix(); suffix != nil {
		current = c.processSuffixChain(current, suffix.(*parser.Unary_subexpr_suffixContext))
	}

	return current
}

func (c *cc) processSuffixChain(base ast.Node, suffix *parser.Unary_subexpr_suffixContext) ast.Node {
	current := base
	for i := 0; i < suffix.GetChildCount(); i++ {
		child := suffix.GetChild(i)
		switch elem := child.(type) {
		case *parser.Key_exprContext:
			current = c.handleKeySuffix(current, elem)
		case *parser.Invoke_exprContext:
			current = c.handleInvokeSuffix(current, elem, i)
		case antlr.TerminalNode:
			if elem.GetText() == "." {
				current = c.handleDotSuffix(current, suffix, &i)
			}
		}
	}
	return current
}

func (c *cc) handleKeySuffix(base ast.Node, keyCtx *parser.Key_exprContext) ast.Node {
	keyNode := c.convertKey_exprContext(keyCtx)
	ind, ok := keyNode.(*ast.A_Indirection)
	if !ok {
		return todo("Key_exprContext", keyCtx)
	}

	if indirection, ok := base.(*ast.A_Indirection); ok {
		indirection.Indirection.Items = append(indirection.Indirection.Items, ind.Indirection.Items...)
		return indirection
	}

	return &ast.A_Indirection{
		Arg: base,
		Indirection: &ast.List{
			Items: []ast.Node{keyNode},
		},
	}
}

func (c *cc) handleInvokeSuffix(base ast.Node, invokeCtx *parser.Invoke_exprContext, idx int) ast.Node {
	funcCall, ok := c.convertInvoke_exprContext(invokeCtx).(*ast.FuncCall)
	if !ok {
		return todo("Invoke_exprContext", invokeCtx)
	}

	if idx == 0 {
		switch baseNode := base.(type) {
		case *ast.ColumnRef:
			if len(baseNode.Fields.Items) > 0 {
				var nameParts []string
				for _, item := range baseNode.Fields.Items {
					if s, ok := item.(*ast.String); ok {
						nameParts = append(nameParts, s.Str)
					}
				}
				funcName := strings.Join(nameParts, ".")

				if funcName == "coalesce" {
					return &ast.CoalesceExpr{
						Args:     funcCall.Args,
						Location: baseNode.Location,
					}
				}

				funcCall.Func = &ast.FuncName{Name: funcName}
				funcCall.Funcname.Items = append(funcCall.Funcname.Items, &ast.String{Str: funcName})

				return funcCall
			}
		default:
			return todo("Invoke_exprContext", invokeCtx)
		}
	}

	stmt := &ast.RecursiveFuncCall{
		Func:        base,
		Funcname:    funcCall.Funcname,
		AggStar:     funcCall.AggStar,
		Location:    funcCall.Location,
		Args:        funcCall.Args,
		AggDistinct: funcCall.AggDistinct,
	}
	stmt.Funcname.Items = append(stmt.Funcname.Items, base)
	return stmt
}

func (c *cc) handleDotSuffix(base ast.Node, suffix *parser.Unary_subexpr_suffixContext, idx *int) ast.Node {
	if *idx+1 >= suffix.GetChildCount() {
		return base
	}

	next := suffix.GetChild(*idx + 1)
	*idx++

	var field ast.Node
	switch v := next.(type) {
	case *parser.Bind_parameterContext:
		field = c.convertBindParameter(v)
	case *parser.An_id_or_typeContext:
		field = &ast.String{Str: parseAnIdOrType(v)}
	case antlr.TerminalNode:
		if val, err := parseIntegerValue(v.GetText()); err == nil {
			field = &ast.A_Const{Val: &ast.Integer{Ival: val}}
		} else {
			return &ast.TODO{}
		}
	}

	if field == nil {
		return base
	}

	if cr, ok := base.(*ast.ColumnRef); ok {
		cr.Fields.Items = append(cr.Fields.Items, field)
		return cr
	}
	return &ast.ColumnRef{
		Fields: &ast.List{Items: []ast.Node{base, field}},
	}
}

func (c *cc) convertKey_exprContext(n *parser.Key_exprContext) ast.Node {
	if n.LBRACE_SQUARE() == nil || n.RBRACE_SQUARE() == nil || n.Expr() == nil {
		return todo("Key_exprContext", n)
	}

	stmt := &ast.A_Indirection{
		Indirection: &ast.List{},
	}

	expr := c.convert(n.Expr())

	stmt.Indirection.Items = append(stmt.Indirection.Items, &ast.A_Indices{
		Uidx: expr,
	})

	return stmt
}

func (c *cc) convertInvoke_exprContext(n *parser.Invoke_exprContext) ast.Node {
	if n.LPAREN() == nil || n.RPAREN() == nil {
		return todo("Invoke_exprContext", n)
	}

	distinct := false
	if n.Opt_set_quantifier() != nil {
		distinct = n.Opt_set_quantifier().DISTINCT() != nil
	}

	stmt := &ast.FuncCall{
		AggDistinct: distinct,
		Funcname:    &ast.List{},
		AggOrder:    &ast.List{},
		Args:        &ast.List{},
		Location:    c.pos(n.GetStart()),
	}

	if nList := n.Named_expr_list(); nList != nil {
		for _, namedExpr := range nList.AllNamed_expr() {
			name := parseAnIdOrType(namedExpr.An_id_or_type())
			expr := c.convert(namedExpr.Expr())

			var res ast.Node
			if rt, ok := expr.(*ast.ResTarget); ok {
				if name != "" {
					rt.Name = &name
				}
				res = rt
			} else if name != "" {
				res = &ast.ResTarget{
					Name:     &name,
					Val:      expr,
					Location: c.pos(namedExpr.Expr().GetStart()),
				}
			} else {
				res = expr
			}

			stmt.Args.Items = append(stmt.Args.Items, res)
		}
	} else if n.ASTERISK() != nil {
		stmt.AggStar = true
	}

	return stmt
}

func (c *cc) convertIdExpr(n *parser.Id_exprContext) ast.Node {
	if id := n.Identifier(); id != nil {
		return &ast.ColumnRef{
			Fields: &ast.List{
				Items: []ast.Node{
					NewIdentifier(id.GetText()),
				},
			},
			Location: c.pos(id.GetStart()),
		}
	}
	return &ast.TODO{}
}

func (c *cc) convertAtomExpr(n *parser.Atom_exprContext) ast.Node {
	switch {
	case n.An_id_or_type() != nil && n.NAMESPACE() != nil:
		return NewIdentifier(parseAnIdOrType(n.An_id_or_type()) + "::" + parseIdOrType(n.Id_or_type()))
	case n.An_id_or_type() != nil:
		return NewIdentifier(parseAnIdOrType(n.An_id_or_type()))
	case n.Literal_value() != nil:
		return c.convertLiteralValue(n.Literal_value().(*parser.Literal_valueContext))
	case n.Bind_parameter() != nil:
		return c.convertBindParameter(n.Bind_parameter().(*parser.Bind_parameterContext))
	default:
		return &ast.TODO{}
	}
}

func (c *cc) convertLiteralValue(n *parser.Literal_valueContext) ast.Node {
	switch {
	case n.Integer() != nil:
		text := n.Integer().GetText()
		val, err := parseIntegerValue(text)
		if err != nil {
			if debug.Active {
				log.Printf("Failed to parse integer value '%s': %v", text, err)
			}
			return &ast.TODO{}
		}
		return &ast.A_Const{Val: &ast.Integer{Ival: val}, Location: c.pos(n.GetStart())}

	case n.Real_() != nil:
		text := n.Real_().GetText()
		return &ast.A_Const{Val: &ast.Float{Str: text}, Location: c.pos(n.GetStart())}

	case n.STRING_VALUE() != nil: // !!! debug !!! (problem with quoted strings)
		val := n.STRING_VALUE().GetText()
		if len(val) >= 2 {
			val = val[1 : len(val)-1]
		}
		return &ast.A_Const{Val: &ast.String{Str: val}, Location: c.pos(n.GetStart())}

	case n.Bool_value() != nil:
		var i bool
		if n.Bool_value().TRUE() != nil {
			i = true
		}
		return &ast.A_Const{Val: &ast.Boolean{Boolval: i}, Location: c.pos(n.GetStart())}

	case n.NULL() != nil:
		return &ast.Null{}

	case n.CURRENT_TIME() != nil:
		if debug.Active {
			log.Printf("TODO: Implement CURRENT_TIME")
		}
		return &ast.TODO{}

	case n.CURRENT_DATE() != nil:
		if debug.Active {
			log.Printf("TODO: Implement CURRENT_DATE")
		}
		return &ast.TODO{}

	case n.CURRENT_TIMESTAMP() != nil:
		if debug.Active {
			log.Printf("TODO: Implement CURRENT_TIMESTAMP")
		}
		return &ast.TODO{}

	case n.BLOB() != nil:
		blobText := n.BLOB().GetText()
		return &ast.A_Const{Val: &ast.String{Str: blobText}, Location: c.pos(n.GetStart())}

	case n.EMPTY_ACTION() != nil:
		if debug.Active {
			log.Printf("TODO: Implement EMPTY_ACTION")
		}
		return &ast.TODO{}

	default:
		if debug.Active {
			log.Printf("Unknown literal value type: %T", n)
		}
		return &ast.TODO{}
	}
}

func (c *cc) convertSqlStmt(n *parser.Sql_stmtContext) ast.Node {
	if n == nil {
		return nil
	}
	// todo: handle explain
	if core := n.Sql_stmt_core(); core != nil {
		return c.convert(core)
	}

	return nil
}

func (c *cc) convert(node node) ast.Node {
	switch n := node.(type) {
	case *parser.Sql_stmtContext:
		return c.convertSqlStmt(n)

	case *parser.Sql_stmt_coreContext:
		return c.convertSqlStmtCore(n)

	case *parser.Create_table_stmtContext:
		return c.convertCreate_table_stmtContext(n)

	case *parser.Select_stmtContext:
		return c.convertSelectStmtContext(n)

	case *parser.Select_coreContext:
		return c.convertSelectCoreContext(n)

	case *parser.Result_columnContext:
		return c.convertResultColumn(n)

	case *parser.Join_sourceContext:
		return c.convertJoinSource(n)

	case *parser.Flatten_sourceContext:
		return c.convertFlattenSource(n)

	case *parser.Named_single_sourceContext:
		return c.convertNamedSingleSource(n)

	case *parser.Single_sourceContext:
		return c.convertSingleSource(n)

	case *parser.Bind_parameterContext:
		return c.convertBindParameter(n)

	case *parser.ExprContext:
		return c.convertExpr(n)

	case *parser.Or_subexprContext:
		return c.convertOrSubExpr(n)

	case *parser.And_subexprContext:
		return c.convertAndSubexpr(n)

	case *parser.Xor_subexprContext:
		return c.convertXorSubexpr(n)

	case *parser.Eq_subexprContext:
		return c.convertEqSubexpr(n)

	case *parser.Neq_subexprContext:
		return c.convertNeqSubexpr(n)

	case *parser.Bit_subexprContext:
		return c.convertBitSubexpr(n)

	case *parser.Add_subexprContext:
		return c.convertAddSubexpr(n)

	case *parser.Mul_subexprContext:
		return c.convertMulSubexpr(n)

	case *parser.Con_subexprContext:
		return c.convertConSubexpr(n)

	case *parser.Unary_subexprContext:
		return c.convertUnarySubexpr(n)

	case *parser.Unary_casual_subexprContext:
		return c.convertUnaryCasualSubexpr(n)

	case *parser.Id_exprContext:
		return c.convertIdExpr(n)

	case *parser.Atom_exprContext:
		return c.convertAtomExpr(n)

	case *parser.Literal_valueContext:
		return c.convertLiteralValue(n)

	case *parser.Json_api_exprContext:
		return c.convertJsonApiExpr(n)

	case *parser.Type_name_compositeContext:
		return c.convertTypeNameComposite(n)

	case *parser.Type_nameContext:
		return c.convertTypeName(n)

	case *parser.Integer_or_bindContext:
		return c.convertIntegerOrBind(n)

	case *parser.Type_name_or_bindContext:
		return c.convertTypeNameOrBind(n)

	case *parser.Into_table_stmtContext:
		return c.convertInto_table_stmtContext(n)

	case *parser.Values_stmtContext:
		return c.convertValues_stmtContext(n)

	case *parser.Returning_columns_listContext:
		return c.convertReturning_columns_listContext(n)

	case *parser.Delete_stmtContext:
		return c.convertDelete_stmtContext(n)

	case *parser.Update_stmtContext:
		return c.convertUpdate_stmtContext(n)

	case *parser.Drop_table_stmtContext:
		return c.convertDrop_table_stmtContext(n)

	case *parser.Commit_stmtContext:
		return c.convertCommit_stmtContext(n)

	case *parser.Rollback_stmtContext:
		return c.convertRollback_stmtContext(n)

	case *parser.Pragma_valueContext:
		return c.convertPragma_valueContext(n)

	case *parser.Pragma_stmtContext:
		return c.convertPragma_stmtContext(n)

	case *parser.Use_stmtContext:
		return c.convertUse_stmtContext(n)

	case *parser.Cluster_exprContext:
		return c.convertCluster_exprContext(n)

	case *parser.Create_user_stmtContext:
		return c.convertCreate_user_stmtContext(n)

	case *parser.Role_nameContext:
		return c.convertRole_nameContext(n)

	case *parser.User_optionContext:
		return c.convertUser_optionContext(n)

	case *parser.Create_group_stmtContext:
		return c.convertCreate_group_stmtContext(n)

	case *parser.Alter_user_stmtContext:
		return c.convertAlter_user_stmtContext(n)

	case *parser.Alter_group_stmtContext:
		return c.convertAlter_group_stmtContext(n)

	case *parser.Drop_role_stmtContext:
		return c.convertDrop_role_stmtCOntext(n)

	default:
		return todo("convert(case=default)", n)
	}
}
