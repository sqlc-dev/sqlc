package astutils

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

type Visitor interface {
	Visit(ast.Node) Visitor
}

type VisitorFunc func(ast.Node)

func (vf VisitorFunc) Visit(node ast.Node) Visitor {
	vf(node)
	return vf
}

func Walk(f Visitor, node ast.Node) {
	if f = f.Visit(node); f == nil {
		return
	}
	switch n := node.(type) {

	case *ast.AlterTableCmd:
		if n.Def != nil {
			Walk(f, n.Def)
		}

	case *ast.AlterTableSetSchemaStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}

	case *ast.AlterTableStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}
		if n.Cmds != nil {
			Walk(f, n.Cmds)
		}

	case *ast.AlterTypeAddValueStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

	case *ast.AlterTypeRenameValueStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

	case *ast.ColumnDef:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *ast.ColumnRef:
		// pass

	case *ast.CommentOnColumnStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}
		if n.Col != nil {
			Walk(f, n.Col)
		}

	case *ast.CommentOnSchemaStmt:
		if n.Schema != nil {
			Walk(f, n.Schema)
		}

	case *ast.CommentOnTableStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}

	case *ast.CommentOnTypeStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

	case *ast.CreateEnumStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Vals != nil {
			Walk(f, n.Vals)
		}

	case *ast.CreateFunctionStmt:
		if n.ReturnType != nil {
			Walk(f, n.ReturnType)
		}
		if n.Func != nil {
			Walk(f, n.Func)
		}

	case *ast.CreateSchemaStmt:
		// pass

	case *ast.CreateTableStmt:
		if n.Name != nil {
			Walk(f, n.Name)
		}

	case *ast.DropFunctionStmt:
		// pass

	case *ast.DropSchemaStmt:
		// pass

	case *ast.DropTableStmt:
		// pass

	case *ast.DropTypeStmt:
		// pass

	case *ast.FuncName:
		// pass

	case *ast.FuncParam:
		if n.Type != nil {
			Walk(f, n.Type)
		}
		if n.DefExpr != nil {
			Walk(f, n.DefExpr)
		}

	case *ast.FuncSpec:
		if n.Name != nil {
			Walk(f, n.Name)
		}

	case *ast.List:
		// pass

	case *ast.RawStmt:
		if n.Stmt != nil {
			Walk(f, n.Stmt)
		}

	case *ast.RenameColumnStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}
		if n.Col != nil {
			Walk(f, n.Col)
		}

	case *ast.RenameTableStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}

	case *ast.ResTarget:
		if n.Val != nil {
			Walk(f, n.Val)
		}

	case *ast.SelectStmt:
		if n.Fields != nil {
			Walk(f, n.Fields)
		}
		if n.From != nil {
			Walk(f, n.From)
		}

	case *ast.Statement:
		if n.Raw != nil {
			Walk(f, n.Raw)
		}

	case *ast.String:
		// pass

	case *ast.TODO:
		// pass

	case *ast.TableName:
		// pass

	case *ast.TypeName:
		// pass

	case *pg.A_ArrayExpr:
		if n.Elements != nil {
			Walk(f, n.Elements)
		}

	case *pg.A_Const:
		if n.Val != nil {
			Walk(f, n.Val)
		}

	case *pg.A_Expr:
		if n.Name != nil {
			Walk(f, n.Name)
		}
		if n.Lexpr != nil {
			Walk(f, n.Lexpr)
		}
		if n.Rexpr != nil {
			Walk(f, n.Rexpr)
		}

	case *pg.A_Indices:
		if n.Lidx != nil {
			Walk(f, n.Lidx)
		}
		if n.Uidx != nil {
			Walk(f, n.Uidx)
		}

	case *pg.A_Indirection:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Indirection != nil {
			Walk(f, n.Indirection)
		}

	case *pg.A_Star:
		// pass

	case *pg.AccessPriv:
		if n.Cols != nil {
			Walk(f, n.Cols)
		}

	case *pg.Aggref:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Aggargtypes != nil {
			Walk(f, n.Aggargtypes)
		}
		if n.Aggdirectargs != nil {
			Walk(f, n.Aggdirectargs)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Aggorder != nil {
			Walk(f, n.Aggorder)
		}
		if n.Aggdistinct != nil {
			Walk(f, n.Aggdistinct)
		}
		if n.Aggfilter != nil {
			Walk(f, n.Aggfilter)
		}

	case *pg.Alias:
		if n.Colnames != nil {
			Walk(f, n.Colnames)
		}

	case *pg.AlterCollationStmt:
		if n.Collname != nil {
			Walk(f, n.Collname)
		}

	case *pg.AlterDatabaseSetStmt:
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *pg.AlterDatabaseStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterDefaultPrivilegesStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Action != nil {
			Walk(f, n.Action)
		}

	case *pg.AlterDomainStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Def != nil {
			Walk(f, n.Def)
		}

	case *pg.AlterEnumStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *pg.AlterEventTrigStmt:
		// pass

	case *pg.AlterExtensionContentsStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *pg.AlterExtensionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterFdwStmt:
		if n.FuncOptions != nil {
			Walk(f, n.FuncOptions)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterForeignServerStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterFunctionStmt:
		if n.Func != nil {
			Walk(f, n.Func)
		}
		if n.Actions != nil {
			Walk(f, n.Actions)
		}

	case *pg.AlterObjectDependsStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}
		if n.Extname != nil {
			Walk(f, n.Extname)
		}

	case *pg.AlterObjectSchemaStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *pg.AlterOpFamilyStmt:
		if n.Opfamilyname != nil {
			Walk(f, n.Opfamilyname)
		}
		if n.Items != nil {
			Walk(f, n.Items)
		}

	case *pg.AlterOperatorStmt:
		if n.Opername != nil {
			Walk(f, n.Opername)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterOwnerStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}
		if n.Newowner != nil {
			Walk(f, n.Newowner)
		}

	case *pg.AlterPolicyStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}
		if n.Roles != nil {
			Walk(f, n.Roles)
		}
		if n.Qual != nil {
			Walk(f, n.Qual)
		}
		if n.WithCheck != nil {
			Walk(f, n.WithCheck)
		}

	case *pg.AlterPublicationStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Tables != nil {
			Walk(f, n.Tables)
		}

	case *pg.AlterRoleSetStmt:
		if n.Role != nil {
			Walk(f, n.Role)
		}
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *pg.AlterRoleStmt:
		if n.Role != nil {
			Walk(f, n.Role)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterSeqStmt:
		if n.Sequence != nil {
			Walk(f, n.Sequence)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterSubscriptionStmt:
		if n.Publication != nil {
			Walk(f, n.Publication)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterSystemStmt:
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *pg.AlterTSConfigurationStmt:
		if n.Cfgname != nil {
			Walk(f, n.Cfgname)
		}
		if n.Tokentype != nil {
			Walk(f, n.Tokentype)
		}
		if n.Dicts != nil {
			Walk(f, n.Dicts)
		}

	case *pg.AlterTSDictionaryStmt:
		if n.Dictname != nil {
			Walk(f, n.Dictname)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterTableCmd:
		if n.Newowner != nil {
			Walk(f, n.Newowner)
		}
		if n.Def != nil {
			Walk(f, n.Def)
		}

	case *pg.AlterTableMoveAllStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *pg.AlterTableSpaceOptionsStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlterTableStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Cmds != nil {
			Walk(f, n.Cmds)
		}

	case *pg.AlterUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.AlternativeSubPlan:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Subplans != nil {
			Walk(f, n.Subplans)
		}

	case *pg.ArrayCoerceExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.ArrayExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Elements != nil {
			Walk(f, n.Elements)
		}

	case *pg.ArrayRef:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Refupperindexpr != nil {
			Walk(f, n.Refupperindexpr)
		}
		if n.Reflowerindexpr != nil {
			Walk(f, n.Reflowerindexpr)
		}
		if n.Refexpr != nil {
			Walk(f, n.Refexpr)
		}
		if n.Refassgnexpr != nil {
			Walk(f, n.Refassgnexpr)
		}

	case *pg.BitString:
		// pass

	case *pg.BlockIdData:
		// pass

	case *pg.BoolExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.BooleanTest:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.CaseExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Defresult != nil {
			Walk(f, n.Defresult)
		}

	case *pg.CaseTestExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.CaseWhen:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Result != nil {
			Walk(f, n.Result)
		}

	case *pg.CheckPointStmt:
		// pass

	case *pg.ClosePortalStmt:
		// pass

	case *pg.ClusterStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *pg.CoalesceExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.CoerceToDomain:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.CoerceToDomainValue:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.CoerceViaIO:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.CollateClause:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Collname != nil {
			Walk(f, n.Collname)
		}

	case *pg.CollateExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.ColumnDef:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.RawDefault != nil {
			Walk(f, n.RawDefault)
		}
		if n.CookedDefault != nil {
			Walk(f, n.CookedDefault)
		}
		if n.CollClause != nil {
			Walk(f, n.CollClause)
		}
		if n.Constraints != nil {
			Walk(f, n.Constraints)
		}
		if n.Fdwoptions != nil {
			Walk(f, n.Fdwoptions)
		}

	case *pg.ColumnRef:
		if n.Fields != nil {
			Walk(f, n.Fields)
		}

	case *pg.CommentStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *pg.CommonTableExpr:
		if n.Aliascolnames != nil {
			Walk(f, n.Aliascolnames)
		}
		if n.Ctequery != nil {
			Walk(f, n.Ctequery)
		}
		if n.Ctecolnames != nil {
			Walk(f, n.Ctecolnames)
		}
		if n.Ctecoltypes != nil {
			Walk(f, n.Ctecoltypes)
		}
		if n.Ctecoltypmods != nil {
			Walk(f, n.Ctecoltypmods)
		}
		if n.Ctecolcollations != nil {
			Walk(f, n.Ctecolcollations)
		}

	case *pg.CompositeTypeStmt:
		if n.Typevar != nil {
			Walk(f, n.Typevar)
		}
		if n.Coldeflist != nil {
			Walk(f, n.Coldeflist)
		}

	case *pg.Const:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.Constraint:
		if n.RawExpr != nil {
			Walk(f, n.RawExpr)
		}
		if n.Keys != nil {
			Walk(f, n.Keys)
		}
		if n.Exclusions != nil {
			Walk(f, n.Exclusions)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.Pktable != nil {
			Walk(f, n.Pktable)
		}
		if n.FkAttrs != nil {
			Walk(f, n.FkAttrs)
		}
		if n.PkAttrs != nil {
			Walk(f, n.PkAttrs)
		}
		if n.OldConpfeqop != nil {
			Walk(f, n.OldConpfeqop)
		}

	case *pg.ConstraintsSetStmt:
		if n.Constraints != nil {
			Walk(f, n.Constraints)
		}

	case *pg.ConvertRowtypeExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.CopyStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Attlist != nil {
			Walk(f, n.Attlist)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateAmStmt:
		if n.HandlerName != nil {
			Walk(f, n.HandlerName)
		}

	case *pg.CreateCastStmt:
		if n.Sourcetype != nil {
			Walk(f, n.Sourcetype)
		}
		if n.Targettype != nil {
			Walk(f, n.Targettype)
		}
		if n.Func != nil {
			Walk(f, n.Func)
		}

	case *pg.CreateConversionStmt:
		if n.ConversionName != nil {
			Walk(f, n.ConversionName)
		}
		if n.FuncName != nil {
			Walk(f, n.FuncName)
		}

	case *pg.CreateDomainStmt:
		if n.Domainname != nil {
			Walk(f, n.Domainname)
		}
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.CollClause != nil {
			Walk(f, n.CollClause)
		}
		if n.Constraints != nil {
			Walk(f, n.Constraints)
		}

	case *pg.CreateEnumStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Vals != nil {
			Walk(f, n.Vals)
		}

	case *pg.CreateEventTrigStmt:
		if n.Whenclause != nil {
			Walk(f, n.Whenclause)
		}
		if n.Funcname != nil {
			Walk(f, n.Funcname)
		}

	case *pg.CreateExtensionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateFdwStmt:
		if n.FuncOptions != nil {
			Walk(f, n.FuncOptions)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateForeignServerStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateForeignTableStmt:
		if n.Base != nil {
			Walk(f, n.Base)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateFunctionStmt:
		if n.Funcname != nil {
			Walk(f, n.Funcname)
		}
		if n.Parameters != nil {
			Walk(f, n.Parameters)
		}
		if n.ReturnType != nil {
			Walk(f, n.ReturnType)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *pg.CreateOpClassItem:
		if n.Name != nil {
			Walk(f, n.Name)
		}
		if n.OrderFamily != nil {
			Walk(f, n.OrderFamily)
		}
		if n.ClassArgs != nil {
			Walk(f, n.ClassArgs)
		}
		if n.Storedtype != nil {
			Walk(f, n.Storedtype)
		}

	case *pg.CreateOpClassStmt:
		if n.Opclassname != nil {
			Walk(f, n.Opclassname)
		}
		if n.Opfamilyname != nil {
			Walk(f, n.Opfamilyname)
		}
		if n.Datatype != nil {
			Walk(f, n.Datatype)
		}
		if n.Items != nil {
			Walk(f, n.Items)
		}

	case *pg.CreateOpFamilyStmt:
		if n.Opfamilyname != nil {
			Walk(f, n.Opfamilyname)
		}

	case *pg.CreatePLangStmt:
		if n.Plhandler != nil {
			Walk(f, n.Plhandler)
		}
		if n.Plinline != nil {
			Walk(f, n.Plinline)
		}
		if n.Plvalidator != nil {
			Walk(f, n.Plvalidator)
		}

	case *pg.CreatePolicyStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}
		if n.Roles != nil {
			Walk(f, n.Roles)
		}
		if n.Qual != nil {
			Walk(f, n.Qual)
		}
		if n.WithCheck != nil {
			Walk(f, n.WithCheck)
		}

	case *pg.CreatePublicationStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Tables != nil {
			Walk(f, n.Tables)
		}

	case *pg.CreateRangeStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Params != nil {
			Walk(f, n.Params)
		}

	case *pg.CreateRoleStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateSchemaStmt:
		if n.Authrole != nil {
			Walk(f, n.Authrole)
		}
		if n.SchemaElts != nil {
			Walk(f, n.SchemaElts)
		}

	case *pg.CreateSeqStmt:
		if n.Sequence != nil {
			Walk(f, n.Sequence)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateStatsStmt:
		if n.Defnames != nil {
			Walk(f, n.Defnames)
		}
		if n.StatTypes != nil {
			Walk(f, n.StatTypes)
		}
		if n.Exprs != nil {
			Walk(f, n.Exprs)
		}
		if n.Relations != nil {
			Walk(f, n.Relations)
		}

	case *pg.CreateStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.TableElts != nil {
			Walk(f, n.TableElts)
		}
		if n.InhRelations != nil {
			Walk(f, n.InhRelations)
		}
		if n.Partbound != nil {
			Walk(f, n.Partbound)
		}
		if n.Partspec != nil {
			Walk(f, n.Partspec)
		}
		if n.OfTypename != nil {
			Walk(f, n.OfTypename)
		}
		if n.Constraints != nil {
			Walk(f, n.Constraints)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateSubscriptionStmt:
		if n.Publication != nil {
			Walk(f, n.Publication)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateTableAsStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Into != nil {
			Walk(f, n.Into)
		}

	case *pg.CreateTableSpaceStmt:
		if n.Owner != nil {
			Walk(f, n.Owner)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreateTransformStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Fromsql != nil {
			Walk(f, n.Fromsql)
		}
		if n.Tosql != nil {
			Walk(f, n.Tosql)
		}

	case *pg.CreateTrigStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Funcname != nil {
			Walk(f, n.Funcname)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Columns != nil {
			Walk(f, n.Columns)
		}
		if n.WhenClause != nil {
			Walk(f, n.WhenClause)
		}
		if n.TransitionRels != nil {
			Walk(f, n.TransitionRels)
		}
		if n.Constrrel != nil {
			Walk(f, n.Constrrel)
		}

	case *pg.CreateUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CreatedbStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.CurrentOfExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.DeallocateStmt:
		// pass

	case *pg.DeclareCursorStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}

	case *pg.DefElem:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.DefineStmt:
		if n.Defnames != nil {
			Walk(f, n.Defnames)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Definition != nil {
			Walk(f, n.Definition)
		}

	case *pg.DeleteStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.UsingClause != nil {
			Walk(f, n.UsingClause)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *pg.DiscardStmt:
		// pass

	case *pg.DoStmt:
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.DropOwnedStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *pg.DropRoleStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *pg.DropStmt:
		if n.Objects != nil {
			Walk(f, n.Objects)
		}

	case *pg.DropSubscriptionStmt:
		// pass

	case *pg.DropTableSpaceStmt:
		// pass

	case *pg.DropUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}

	case *pg.DropdbStmt:
		// pass

	case *pg.ExecuteStmt:
		if n.Params != nil {
			Walk(f, n.Params)
		}

	case *pg.ExplainStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.Expr:
		// pass

	case *pg.FetchStmt:
		// pass

	case *pg.FieldSelect:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.FieldStore:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Newvals != nil {
			Walk(f, n.Newvals)
		}
		if n.Fieldnums != nil {
			Walk(f, n.Fieldnums)
		}

	case *pg.Float:
		// pass

	case *pg.FromExpr:
		if n.Fromlist != nil {
			Walk(f, n.Fromlist)
		}
		if n.Quals != nil {
			Walk(f, n.Quals)
		}

	case *ast.FuncCall:
		if n.Func != nil {
			Walk(f, n.Func)
		}
		if n.Funcname != nil {
			Walk(f, n.Funcname)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.AggOrder != nil {
			Walk(f, n.AggOrder)
		}
		if n.AggFilter != nil {
			Walk(f, n.AggFilter)
		}
		if n.Over != nil {
			Walk(f, n.Over)
		}

	case *pg.FuncExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.FunctionParameter:
		if n.ArgType != nil {
			Walk(f, n.ArgType)
		}
		if n.Defexpr != nil {
			Walk(f, n.Defexpr)
		}

	case *pg.GrantRoleStmt:
		if n.GrantedRoles != nil {
			Walk(f, n.GrantedRoles)
		}
		if n.GranteeRoles != nil {
			Walk(f, n.GranteeRoles)
		}
		if n.Grantor != nil {
			Walk(f, n.Grantor)
		}

	case *pg.GrantStmt:
		if n.Objects != nil {
			Walk(f, n.Objects)
		}
		if n.Privileges != nil {
			Walk(f, n.Privileges)
		}
		if n.Grantees != nil {
			Walk(f, n.Grantees)
		}

	case *pg.GroupingFunc:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Refs != nil {
			Walk(f, n.Refs)
		}
		if n.Cols != nil {
			Walk(f, n.Cols)
		}

	case *pg.GroupingSet:
		if n.Content != nil {
			Walk(f, n.Content)
		}

	case *pg.ImportForeignSchemaStmt:
		if n.TableList != nil {
			Walk(f, n.TableList)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.IndexElem:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Collation != nil {
			Walk(f, n.Collation)
		}
		if n.Opclass != nil {
			Walk(f, n.Opclass)
		}

	case *pg.IndexStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.IndexParams != nil {
			Walk(f, n.IndexParams)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.ExcludeOpNames != nil {
			Walk(f, n.ExcludeOpNames)
		}

	case *pg.InferClause:
		if n.IndexElems != nil {
			Walk(f, n.IndexElems)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}

	case *pg.InferenceElem:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}

	case *pg.InlineCodeBlock:
		// pass

	case *pg.InsertStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Cols != nil {
			Walk(f, n.Cols)
		}
		if n.SelectStmt != nil {
			Walk(f, n.SelectStmt)
		}
		if n.OnConflictClause != nil {
			Walk(f, n.OnConflictClause)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *pg.Integer:
		// pass

	case *pg.IntoClause:
		if n.Rel != nil {
			Walk(f, n.Rel)
		}
		if n.ColNames != nil {
			Walk(f, n.ColNames)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.ViewQuery != nil {
			Walk(f, n.ViewQuery)
		}

	case *pg.JoinExpr:
		if n.Larg != nil {
			Walk(f, n.Larg)
		}
		if n.Rarg != nil {
			Walk(f, n.Rarg)
		}
		if n.UsingClause != nil {
			Walk(f, n.UsingClause)
		}
		if n.Quals != nil {
			Walk(f, n.Quals)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *pg.ListenStmt:
		// pass

	case *pg.LoadStmt:
		// pass

	case *pg.LockStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
		}

	case *pg.LockingClause:
		if n.LockedRels != nil {
			Walk(f, n.LockedRels)
		}

	case *pg.MinMaxExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.MultiAssignRef:
		if n.Source != nil {
			Walk(f, n.Source)
		}

	case *pg.NamedArgExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.NextValueExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.NotifyStmt:
		// pass

	case *pg.Null:
		// pass

	case *pg.NullTest:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.ObjectWithArgs:
		if n.Objname != nil {
			Walk(f, n.Objname)
		}
		if n.Objargs != nil {
			Walk(f, n.Objargs)
		}

	case *pg.OnConflictClause:
		if n.Infer != nil {
			Walk(f, n.Infer)
		}
		if n.TargetList != nil {
			Walk(f, n.TargetList)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}

	case *pg.OnConflictExpr:
		if n.ArbiterElems != nil {
			Walk(f, n.ArbiterElems)
		}
		if n.ArbiterWhere != nil {
			Walk(f, n.ArbiterWhere)
		}
		if n.OnConflictSet != nil {
			Walk(f, n.OnConflictSet)
		}
		if n.OnConflictWhere != nil {
			Walk(f, n.OnConflictWhere)
		}
		if n.ExclRelTlist != nil {
			Walk(f, n.ExclRelTlist)
		}

	case *pg.OpExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.Param:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.ParamExecData:
		// pass

	case *pg.ParamExternData:
		// pass

	case *pg.ParamListInfoData:
		// pass

	case *pg.ParamRef:
		// pass

	case *pg.PartitionBoundSpec:
		if n.Listdatums != nil {
			Walk(f, n.Listdatums)
		}
		if n.Lowerdatums != nil {
			Walk(f, n.Lowerdatums)
		}
		if n.Upperdatums != nil {
			Walk(f, n.Upperdatums)
		}

	case *pg.PartitionCmd:
		if n.Name != nil {
			Walk(f, n.Name)
		}
		if n.Bound != nil {
			Walk(f, n.Bound)
		}

	case *pg.PartitionElem:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Collation != nil {
			Walk(f, n.Collation)
		}
		if n.Opclass != nil {
			Walk(f, n.Opclass)
		}

	case *pg.PartitionRangeDatum:
		if n.Value != nil {
			Walk(f, n.Value)
		}

	case *pg.PartitionSpec:
		if n.PartParams != nil {
			Walk(f, n.PartParams)
		}

	case *pg.PrepareStmt:
		if n.Argtypes != nil {
			Walk(f, n.Argtypes)
		}
		if n.Query != nil {
			Walk(f, n.Query)
		}

	case *pg.Query:
		if n.UtilityStmt != nil {
			Walk(f, n.UtilityStmt)
		}
		if n.CteList != nil {
			Walk(f, n.CteList)
		}
		if n.Rtable != nil {
			Walk(f, n.Rtable)
		}
		if n.Jointree != nil {
			Walk(f, n.Jointree)
		}
		if n.TargetList != nil {
			Walk(f, n.TargetList)
		}
		if n.OnConflict != nil {
			Walk(f, n.OnConflict)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.GroupClause != nil {
			Walk(f, n.GroupClause)
		}
		if n.GroupingSets != nil {
			Walk(f, n.GroupingSets)
		}
		if n.HavingQual != nil {
			Walk(f, n.HavingQual)
		}
		if n.WindowClause != nil {
			Walk(f, n.WindowClause)
		}
		if n.DistinctClause != nil {
			Walk(f, n.DistinctClause)
		}
		if n.SortClause != nil {
			Walk(f, n.SortClause)
		}
		if n.LimitOffset != nil {
			Walk(f, n.LimitOffset)
		}
		if n.LimitCount != nil {
			Walk(f, n.LimitCount)
		}
		if n.RowMarks != nil {
			Walk(f, n.RowMarks)
		}
		if n.SetOperations != nil {
			Walk(f, n.SetOperations)
		}
		if n.ConstraintDeps != nil {
			Walk(f, n.ConstraintDeps)
		}
		if n.WithCheckOptions != nil {
			Walk(f, n.WithCheckOptions)
		}

	case *pg.RangeFunction:
		if n.Functions != nil {
			Walk(f, n.Functions)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}
		if n.Coldeflist != nil {
			Walk(f, n.Coldeflist)
		}

	case *pg.RangeSubselect:
		if n.Subquery != nil {
			Walk(f, n.Subquery)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *pg.RangeTableFunc:
		if n.Docexpr != nil {
			Walk(f, n.Docexpr)
		}
		if n.Rowexpr != nil {
			Walk(f, n.Rowexpr)
		}
		if n.Namespaces != nil {
			Walk(f, n.Namespaces)
		}
		if n.Columns != nil {
			Walk(f, n.Columns)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *pg.RangeTableFuncCol:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Colexpr != nil {
			Walk(f, n.Colexpr)
		}
		if n.Coldefexpr != nil {
			Walk(f, n.Coldefexpr)
		}

	case *pg.RangeTableSample:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Method != nil {
			Walk(f, n.Method)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Repeatable != nil {
			Walk(f, n.Repeatable)
		}

	case *pg.RangeTblEntry:
		if n.Tablesample != nil {
			Walk(f, n.Tablesample)
		}
		if n.Subquery != nil {
			Walk(f, n.Subquery)
		}
		if n.Joinaliasvars != nil {
			Walk(f, n.Joinaliasvars)
		}
		if n.Functions != nil {
			Walk(f, n.Functions)
		}
		if n.Tablefunc != nil {
			Walk(f, n.Tablefunc)
		}
		if n.ValuesLists != nil {
			Walk(f, n.ValuesLists)
		}
		if n.Coltypes != nil {
			Walk(f, n.Coltypes)
		}
		if n.Coltypmods != nil {
			Walk(f, n.Coltypmods)
		}
		if n.Colcollations != nil {
			Walk(f, n.Colcollations)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}
		if n.Eref != nil {
			Walk(f, n.Eref)
		}
		if n.SecurityQuals != nil {
			Walk(f, n.SecurityQuals)
		}

	case *pg.RangeTblFunction:
		if n.Funcexpr != nil {
			Walk(f, n.Funcexpr)
		}
		if n.Funccolnames != nil {
			Walk(f, n.Funccolnames)
		}
		if n.Funccoltypes != nil {
			Walk(f, n.Funccoltypes)
		}
		if n.Funccoltypmods != nil {
			Walk(f, n.Funccoltypmods)
		}
		if n.Funccolcollations != nil {
			Walk(f, n.Funccolcollations)
		}

	case *pg.RangeTblRef:
		// pass

	case *pg.RangeVar:
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *pg.RawStmt:
		if n.Stmt != nil {
			Walk(f, n.Stmt)
		}

	case *pg.ReassignOwnedStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}
		if n.Newrole != nil {
			Walk(f, n.Newrole)
		}

	case *pg.RefreshMatViewStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *pg.ReindexStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *pg.RelabelType:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *pg.RenameStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *pg.ReplicaIdentityStmt:
		// pass

	case *pg.ResTarget:
		if n.Indirection != nil {
			Walk(f, n.Indirection)
		}
		if n.Val != nil {
			Walk(f, n.Val)
		}

	case *pg.RoleSpec:
		// pass

	case *pg.RowCompareExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Opnos != nil {
			Walk(f, n.Opnos)
		}
		if n.Opfamilies != nil {
			Walk(f, n.Opfamilies)
		}
		if n.Inputcollids != nil {
			Walk(f, n.Inputcollids)
		}
		if n.Largs != nil {
			Walk(f, n.Largs)
		}
		if n.Rargs != nil {
			Walk(f, n.Rargs)
		}

	case *pg.RowExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Colnames != nil {
			Walk(f, n.Colnames)
		}

	case *pg.RowMarkClause:
		// pass

	case *pg.RuleStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.Actions != nil {
			Walk(f, n.Actions)
		}

	case *pg.SQLValueFunction:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.ScalarArrayOpExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.SecLabelStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *pg.SelectStmt:
		if n.DistinctClause != nil {
			Walk(f, n.DistinctClause)
		}
		if n.IntoClause != nil {
			Walk(f, n.IntoClause)
		}
		if n.TargetList != nil {
			Walk(f, n.TargetList)
		}
		if n.FromClause != nil {
			Walk(f, n.FromClause)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.GroupClause != nil {
			Walk(f, n.GroupClause)
		}
		if n.HavingClause != nil {
			Walk(f, n.HavingClause)
		}
		if n.WindowClause != nil {
			Walk(f, n.WindowClause)
		}
		if n.SortClause != nil {
			Walk(f, n.SortClause)
		}
		if n.LimitOffset != nil {
			Walk(f, n.LimitOffset)
		}
		if n.LimitCount != nil {
			Walk(f, n.LimitCount)
		}
		if n.LockingClause != nil {
			Walk(f, n.LockingClause)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}
		if n.Larg != nil {
			Walk(f, n.Larg)
		}
		if n.Rarg != nil {
			Walk(f, n.Rarg)
		}

	case *pg.SetOperationStmt:
		if n.Larg != nil {
			Walk(f, n.Larg)
		}
		if n.Rarg != nil {
			Walk(f, n.Rarg)
		}
		if n.ColTypes != nil {
			Walk(f, n.ColTypes)
		}
		if n.ColTypmods != nil {
			Walk(f, n.ColTypmods)
		}
		if n.ColCollations != nil {
			Walk(f, n.ColCollations)
		}
		if n.GroupClauses != nil {
			Walk(f, n.GroupClauses)
		}

	case *pg.SetToDefault:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.SortBy:
		if n.Node != nil {
			Walk(f, n.Node)
		}
		if n.UseOp != nil {
			Walk(f, n.UseOp)
		}

	case *pg.SortGroupClause:
		// pass

	case *pg.String:
		// pass

	case *pg.SubLink:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Testexpr != nil {
			Walk(f, n.Testexpr)
		}
		if n.OperName != nil {
			Walk(f, n.OperName)
		}
		if n.Subselect != nil {
			Walk(f, n.Subselect)
		}

	case *pg.SubPlan:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Testexpr != nil {
			Walk(f, n.Testexpr)
		}
		if n.ParamIds != nil {
			Walk(f, n.ParamIds)
		}
		if n.SetParam != nil {
			Walk(f, n.SetParam)
		}
		if n.ParParam != nil {
			Walk(f, n.ParParam)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.TableFunc:
		if n.NsUris != nil {
			Walk(f, n.NsUris)
		}
		if n.NsNames != nil {
			Walk(f, n.NsNames)
		}
		if n.Docexpr != nil {
			Walk(f, n.Docexpr)
		}
		if n.Rowexpr != nil {
			Walk(f, n.Rowexpr)
		}
		if n.Colnames != nil {
			Walk(f, n.Colnames)
		}
		if n.Coltypes != nil {
			Walk(f, n.Coltypes)
		}
		if n.Coltypmods != nil {
			Walk(f, n.Coltypmods)
		}
		if n.Colcollations != nil {
			Walk(f, n.Colcollations)
		}
		if n.Colexprs != nil {
			Walk(f, n.Colexprs)
		}
		if n.Coldefexprs != nil {
			Walk(f, n.Coldefexprs)
		}

	case *pg.TableLikeClause:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *pg.TableSampleClause:
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Repeatable != nil {
			Walk(f, n.Repeatable)
		}

	case *pg.TargetEntry:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}

	case *pg.TransactionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.TriggerTransition:
		// pass

	case *pg.TruncateStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
		}

	case *pg.TypeCast:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *pg.TypeName:
		if n.Names != nil {
			Walk(f, n.Names)
		}
		if n.Typmods != nil {
			Walk(f, n.Typmods)
		}
		if n.ArrayBounds != nil {
			Walk(f, n.ArrayBounds)
		}

	case *pg.UnlistenStmt:
		// pass

	case *pg.UpdateStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.TargetList != nil {
			Walk(f, n.TargetList)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.FromClause != nil {
			Walk(f, n.FromClause)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *pg.VacuumStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.VaCols != nil {
			Walk(f, n.VaCols)
		}

	case *pg.Var:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *pg.VariableSetStmt:
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.VariableShowStmt:
		// pass

	case *pg.ViewStmt:
		if n.View != nil {
			Walk(f, n.View)
		}
		if n.Aliases != nil {
			Walk(f, n.Aliases)
		}
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *pg.WindowClause:
		if n.PartitionClause != nil {
			Walk(f, n.PartitionClause)
		}
		if n.OrderClause != nil {
			Walk(f, n.OrderClause)
		}
		if n.StartOffset != nil {
			Walk(f, n.StartOffset)
		}
		if n.EndOffset != nil {
			Walk(f, n.EndOffset)
		}

	case *ast.WindowDef:
		if n.PartitionClause != nil {
			Walk(f, n.PartitionClause)
		}
		if n.OrderClause != nil {
			Walk(f, n.OrderClause)
		}
		if n.StartOffset != nil {
			Walk(f, n.StartOffset)
		}
		if n.EndOffset != nil {
			Walk(f, n.EndOffset)
		}

	case *pg.WindowFunc:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Aggfilter != nil {
			Walk(f, n.Aggfilter)
		}

	case *pg.WithCheckOption:
		if n.Qual != nil {
			Walk(f, n.Qual)
		}

	case *pg.WithClause:
		if n.Ctes != nil {
			Walk(f, n.Ctes)
		}

	case *pg.XmlExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.NamedArgs != nil {
			Walk(f, n.NamedArgs)
		}
		if n.ArgNames != nil {
			Walk(f, n.ArgNames)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *pg.XmlSerialize:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))
	}

	f.Visit(nil)
}
