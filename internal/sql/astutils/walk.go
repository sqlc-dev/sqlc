package astutils

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
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

	case *ast.AlterTableSetSchemaStmt:
		if n.Table != nil {
			Walk(f, n.Table)
		}

	case *ast.AlterTableStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
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

	case *ast.AlterTypeSetSchemaStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

	case *ast.AlterTypeRenameValueStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

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

	case *ast.CommentOnViewStmt:
		if n.View != nil {
			Walk(f, n.View)
		}

	case *ast.CompositeTypeStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

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
		for _, item := range n.Items {
			Walk(f, item)
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

	case *ast.RenameTypeStmt:
		if n.Type != nil {
			Walk(f, n.Type)
		}

	case *ast.Statement:
		if n.Raw != nil {
			Walk(f, n.Raw)
		}

	case *ast.TODO:
		// pass

	case *ast.TableName:
		// pass

	case *ast.A_ArrayExpr:
		if n.Elements != nil {
			Walk(f, n.Elements)
		}

	case *ast.A_Const:
		if n.Val != nil {
			Walk(f, n.Val)
		}

	case *ast.A_Expr:
		if n.Name != nil {
			Walk(f, n.Name)
		}
		if n.Lexpr != nil {
			Walk(f, n.Lexpr)
		}
		if n.Rexpr != nil {
			Walk(f, n.Rexpr)
		}

	case *ast.A_Indices:
		if n.Lidx != nil {
			Walk(f, n.Lidx)
		}
		if n.Uidx != nil {
			Walk(f, n.Uidx)
		}

	case *ast.A_Indirection:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Indirection != nil {
			Walk(f, n.Indirection)
		}

	case *ast.A_Star:
		// pass

	case *ast.AccessPriv:
		if n.Cols != nil {
			Walk(f, n.Cols)
		}

	case *ast.Aggref:
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

	case *ast.Alias:
		if n.Colnames != nil {
			Walk(f, n.Colnames)
		}

	case *ast.AlterCollationStmt:
		if n.Collname != nil {
			Walk(f, n.Collname)
		}

	case *ast.AlterDatabaseSetStmt:
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *ast.AlterDatabaseStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterDefaultPrivilegesStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Action != nil {
			Walk(f, n.Action)
		}

	case *ast.AlterDomainStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Def != nil {
			Walk(f, n.Def)
		}

	case *ast.AlterEnumStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *ast.AlterEventTrigStmt:
		// pass

	case *ast.AlterExtensionContentsStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *ast.AlterExtensionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterFdwStmt:
		if n.FuncOptions != nil {
			Walk(f, n.FuncOptions)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterForeignServerStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterFunctionStmt:
		if n.Func != nil {
			Walk(f, n.Func)
		}
		if n.Actions != nil {
			Walk(f, n.Actions)
		}

	case *ast.AlterObjectDependsStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}
		if n.Extname != nil {
			Walk(f, n.Extname)
		}

	case *ast.AlterObjectSchemaStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *ast.AlterOpFamilyStmt:
		if n.Opfamilyname != nil {
			Walk(f, n.Opfamilyname)
		}
		if n.Items != nil {
			Walk(f, n.Items)
		}

	case *ast.AlterOperatorStmt:
		if n.Opername != nil {
			Walk(f, n.Opername)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterOwnerStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}
		if n.Newowner != nil {
			Walk(f, n.Newowner)
		}

	case *ast.AlterPolicyStmt:
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

	case *ast.AlterPublicationStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Tables != nil {
			Walk(f, n.Tables)
		}

	case *ast.AlterRoleSetStmt:
		if n.Role != nil {
			Walk(f, n.Role)
		}
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *ast.AlterRoleStmt:
		if n.Role != nil {
			Walk(f, n.Role)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterSeqStmt:
		if n.Sequence != nil {
			Walk(f, n.Sequence)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterSubscriptionStmt:
		if n.Publication != nil {
			Walk(f, n.Publication)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterSystemStmt:
		if n.Setstmt != nil {
			Walk(f, n.Setstmt)
		}

	case *ast.AlterTSConfigurationStmt:
		if n.Cfgname != nil {
			Walk(f, n.Cfgname)
		}
		if n.Tokentype != nil {
			Walk(f, n.Tokentype)
		}
		if n.Dicts != nil {
			Walk(f, n.Dicts)
		}

	case *ast.AlterTSDictionaryStmt:
		if n.Dictname != nil {
			Walk(f, n.Dictname)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterTableCmd:
		if n.Newowner != nil {
			Walk(f, n.Newowner)
		}
		if n.Def != nil {
			Walk(f, n.Def)
		}

	case *ast.AlterTableMoveAllStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *ast.AlterTableSpaceOptionsStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlterUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.AlternativeSubPlan:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Subplans != nil {
			Walk(f, n.Subplans)
		}

	case *ast.ArrayCoerceExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.ArrayExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Elements != nil {
			Walk(f, n.Elements)
		}

	case *ast.ArrayRef:
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

	case *ast.BetweenExpr:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Left != nil {
			Walk(f, n.Left)
		}
		if n.Right != nil {
			Walk(f, n.Right)
		}

	case *ast.BitString:
		// pass

	case *ast.BlockIdData:
		// pass

	case *ast.BoolExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.Boolean:
		// pass

	case *ast.BooleanTest:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.CallStmt:
		if n.FuncCall != nil {
			Walk(f, n.FuncCall)
		}

	case *ast.CaseExpr:
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

	case *ast.CaseTestExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.CaseWhen:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Result != nil {
			Walk(f, n.Result)
		}

	case *ast.CheckPointStmt:
		// pass

	case *ast.ClosePortalStmt:
		// pass

	case *ast.ClusterStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *ast.CoalesceExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.CoerceToDomain:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.CoerceToDomainValue:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.CoerceViaIO:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.CollateClause:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.Collname != nil {
			Walk(f, n.Collname)
		}

	case *ast.CollateExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.ColumnDef:
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

	case *ast.ColumnRef:
		if n.Fields != nil {
			Walk(f, n.Fields)
		}

	case *ast.CommentStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *ast.CommonTableExpr:
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

	case *ast.Const:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.Constraint:
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

	case *ast.ConstraintsSetStmt:
		if n.Constraints != nil {
			Walk(f, n.Constraints)
		}

	case *ast.ConvertRowtypeExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.CopyStmt:
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

	case *ast.CreateAmStmt:
		if n.HandlerName != nil {
			Walk(f, n.HandlerName)
		}

	case *ast.CreateCastStmt:
		if n.Sourcetype != nil {
			Walk(f, n.Sourcetype)
		}
		if n.Targettype != nil {
			Walk(f, n.Targettype)
		}
		if n.Func != nil {
			Walk(f, n.Func)
		}

	case *ast.CreateConversionStmt:
		if n.ConversionName != nil {
			Walk(f, n.ConversionName)
		}
		if n.FuncName != nil {
			Walk(f, n.FuncName)
		}

	case *ast.CreateDomainStmt:
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

	case *ast.CreateEnumStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Vals != nil {
			Walk(f, n.Vals)
		}

	case *ast.CreateEventTrigStmt:
		if n.Whenclause != nil {
			Walk(f, n.Whenclause)
		}
		if n.Funcname != nil {
			Walk(f, n.Funcname)
		}

	case *ast.CreateExtensionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateFdwStmt:
		if n.FuncOptions != nil {
			Walk(f, n.FuncOptions)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateForeignServerStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateForeignTableStmt:
		if n.Base != nil {
			Walk(f, n.Base)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateFunctionStmt:
		if n.Func != nil {
			Walk(f, n.Func)
		}
		if n.Params != nil {
			Walk(f, n.Params)
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

	case *ast.CreateOpClassItem:
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

	case *ast.CreateOpClassStmt:
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

	case *ast.CreateOpFamilyStmt:
		if n.Opfamilyname != nil {
			Walk(f, n.Opfamilyname)
		}

	case *ast.CreatePLangStmt:
		if n.Plhandler != nil {
			Walk(f, n.Plhandler)
		}
		if n.Plinline != nil {
			Walk(f, n.Plinline)
		}
		if n.Plvalidator != nil {
			Walk(f, n.Plvalidator)
		}

	case *ast.CreatePolicyStmt:
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

	case *ast.CreatePublicationStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}
		if n.Tables != nil {
			Walk(f, n.Tables)
		}

	case *ast.CreateRangeStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Params != nil {
			Walk(f, n.Params)
		}

	case *ast.CreateRoleStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateSchemaStmt:
		if n.Authrole != nil {
			Walk(f, n.Authrole)
		}
		if n.SchemaElts != nil {
			Walk(f, n.SchemaElts)
		}

	case *ast.CreateSeqStmt:
		if n.Sequence != nil {
			Walk(f, n.Sequence)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateStatsStmt:
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

	case *ast.CreateStmt:
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

	case *ast.CreateSubscriptionStmt:
		if n.Publication != nil {
			Walk(f, n.Publication)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateTableAsStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Into != nil {
			Walk(f, n.Into)
		}

	case *ast.CreateTableSpaceStmt:
		if n.Owner != nil {
			Walk(f, n.Owner)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreateTransformStmt:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Fromsql != nil {
			Walk(f, n.Fromsql)
		}
		if n.Tosql != nil {
			Walk(f, n.Tosql)
		}

	case *ast.CreateTrigStmt:
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

	case *ast.CreateUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CreatedbStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.CurrentOfExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.DeallocateStmt:
		// pass

	case *ast.DeclareCursorStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}

	case *ast.DefElem:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.DefineStmt:
		if n.Defnames != nil {
			Walk(f, n.Defnames)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Definition != nil {
			Walk(f, n.Definition)
		}

	case *ast.DeleteStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
		}
		if n.UsingClause != nil {
			Walk(f, n.UsingClause)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.LimitCount != nil {
			Walk(f, n.LimitCount)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *ast.DiscardStmt:
		// pass

	case *ast.DoStmt:
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.DropOwnedStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *ast.DropRoleStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}

	case *ast.DropStmt:
		if n.Objects != nil {
			Walk(f, n.Objects)
		}

	case *ast.DropSubscriptionStmt:
		// pass

	case *ast.DropTableSpaceStmt:
		// pass

	case *ast.DropUserMappingStmt:
		if n.User != nil {
			Walk(f, n.User)
		}

	case *ast.DropdbStmt:
		// pass

	case *ast.ExecuteStmt:
		if n.Params != nil {
			Walk(f, n.Params)
		}

	case *ast.ExplainStmt:
		if n.Query != nil {
			Walk(f, n.Query)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.Expr:
		// pass

	case *ast.FetchStmt:
		// pass

	case *ast.FieldSelect:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.FieldStore:
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

	case *ast.Float:
		// pass

	case *ast.FromExpr:
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

	case *ast.FuncExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.FunctionParameter:
		if n.ArgType != nil {
			Walk(f, n.ArgType)
		}
		if n.Defexpr != nil {
			Walk(f, n.Defexpr)
		}

	case *ast.GrantRoleStmt:
		if n.GrantedRoles != nil {
			Walk(f, n.GrantedRoles)
		}
		if n.GranteeRoles != nil {
			Walk(f, n.GranteeRoles)
		}
		if n.Grantor != nil {
			Walk(f, n.Grantor)
		}

	case *ast.GrantStmt:
		if n.Objects != nil {
			Walk(f, n.Objects)
		}
		if n.Privileges != nil {
			Walk(f, n.Privileges)
		}
		if n.Grantees != nil {
			Walk(f, n.Grantees)
		}

	case *ast.GroupingFunc:
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

	case *ast.GroupingSet:
		if n.Content != nil {
			Walk(f, n.Content)
		}

	case *ast.ImportForeignSchemaStmt:
		if n.TableList != nil {
			Walk(f, n.TableList)
		}
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.IndexElem:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Collation != nil {
			Walk(f, n.Collation)
		}
		if n.Opclass != nil {
			Walk(f, n.Opclass)
		}

	case *ast.IndexStmt:
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

	case *ast.InferClause:
		if n.IndexElems != nil {
			Walk(f, n.IndexElems)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}

	case *ast.InferenceElem:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}

	case *ast.InlineCodeBlock:
		// pass

	case *ast.InsertStmt:
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

	case *ast.Integer:
		// pass

	case *ast.IntoClause:
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

	case *ast.JoinExpr:
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

	case *ast.ListenStmt:
		// pass

	case *ast.LoadStmt:
		// pass

	case *ast.LockStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
		}

	case *ast.LockingClause:
		if n.LockedRels != nil {
			Walk(f, n.LockedRels)
		}

	case *ast.MinMaxExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.MultiAssignRef:
		if n.Source != nil {
			Walk(f, n.Source)
		}

	case *ast.NamedArgExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.NextValueExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.NotifyStmt:
		// pass

	case *ast.Null:
		// pass

	case *ast.NullTest:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.ObjectWithArgs:
		if n.Objname != nil {
			Walk(f, n.Objname)
		}
		if n.Objargs != nil {
			Walk(f, n.Objargs)
		}

	case *ast.OnConflictClause:
		if n.Infer != nil {
			Walk(f, n.Infer)
		}
		if n.TargetList != nil {
			Walk(f, n.TargetList)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}

	case *ast.OnConflictExpr:
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

	case *ast.OpExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.Param:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.ParamExecData:
		// pass

	case *ast.ParamExternData:
		// pass

	case *ast.ParamListInfoData:
		// pass

	case *ast.ParamRef:
		// pass

	case *ast.PartitionBoundSpec:
		if n.Listdatums != nil {
			Walk(f, n.Listdatums)
		}
		if n.Lowerdatums != nil {
			Walk(f, n.Lowerdatums)
		}
		if n.Upperdatums != nil {
			Walk(f, n.Upperdatums)
		}

	case *ast.PartitionCmd:
		if n.Name != nil {
			Walk(f, n.Name)
		}
		if n.Bound != nil {
			Walk(f, n.Bound)
		}

	case *ast.PartitionElem:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.Collation != nil {
			Walk(f, n.Collation)
		}
		if n.Opclass != nil {
			Walk(f, n.Opclass)
		}

	case *ast.PartitionRangeDatum:
		if n.Value != nil {
			Walk(f, n.Value)
		}

	case *ast.PartitionSpec:
		if n.PartParams != nil {
			Walk(f, n.PartParams)
		}

	case *ast.PrepareStmt:
		if n.Argtypes != nil {
			Walk(f, n.Argtypes)
		}
		if n.Query != nil {
			Walk(f, n.Query)
		}

	case *ast.Query:
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

	case *ast.RangeFunction:
		if n.Functions != nil {
			Walk(f, n.Functions)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}
		if n.Coldeflist != nil {
			Walk(f, n.Coldeflist)
		}

	case *ast.RangeSubselect:
		if n.Subquery != nil {
			Walk(f, n.Subquery)
		}
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *ast.RangeTableFunc:
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

	case *ast.RangeTableFuncCol:
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}
		if n.Colexpr != nil {
			Walk(f, n.Colexpr)
		}
		if n.Coldefexpr != nil {
			Walk(f, n.Coldefexpr)
		}

	case *ast.RangeTableSample:
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

	case *ast.RangeTblEntry:
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

	case *ast.RangeTblFunction:
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

	case *ast.RangeTblRef:
		// pass

	case *ast.RangeVar:
		if n.Alias != nil {
			Walk(f, n.Alias)
		}

	case *ast.RawStmt:
		if n.Stmt != nil {
			Walk(f, n.Stmt)
		}

	case *ast.ReassignOwnedStmt:
		if n.Roles != nil {
			Walk(f, n.Roles)
		}
		if n.Newrole != nil {
			Walk(f, n.Newrole)
		}

	case *ast.RefreshMatViewStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *ast.ReindexStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *ast.RelabelType:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Arg != nil {
			Walk(f, n.Arg)
		}

	case *ast.RenameStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *ast.ReplicaIdentityStmt:
		// pass

	case *ast.ResTarget:
		if n.Indirection != nil {
			Walk(f, n.Indirection)
		}
		if n.Val != nil {
			Walk(f, n.Val)
		}

	case *ast.RoleSpec:
		// pass

	case *ast.RowCompareExpr:
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

	case *ast.RowExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Colnames != nil {
			Walk(f, n.Colnames)
		}

	case *ast.RowMarkClause:
		// pass

	case *ast.RuleStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.WhereClause != nil {
			Walk(f, n.WhereClause)
		}
		if n.Actions != nil {
			Walk(f, n.Actions)
		}

	case *ast.SQLValueFunction:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.ScalarArrayOpExpr:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.SecLabelStmt:
		if n.Object != nil {
			Walk(f, n.Object)
		}

	case *ast.SelectStmt:
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
		if n.ValuesLists != nil {
			Walk(f, n.ValuesLists)
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

	case *ast.SetOperationStmt:
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

	case *ast.SetToDefault:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.SortBy:
		if n.Node != nil {
			Walk(f, n.Node)
		}
		if n.UseOp != nil {
			Walk(f, n.UseOp)
		}

	case *ast.SortGroupClause:
		// pass

	case *ast.String:
		// pass

	case *ast.SubLink:
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

	case *ast.SubPlan:
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

	case *ast.TableFunc:
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

	case *ast.TableLikeClause:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}

	case *ast.TableSampleClause:
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Repeatable != nil {
			Walk(f, n.Repeatable)
		}

	case *ast.TargetEntry:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Expr != nil {
			Walk(f, n.Expr)
		}

	case *ast.TransactionStmt:
		if n.Options != nil {
			Walk(f, n.Options)
		}

	case *ast.TriggerTransition:
		// pass

	case *ast.TruncateStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
		}

	case *ast.TypeCast:
		if n.Arg != nil {
			Walk(f, n.Arg)
		}
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *ast.TypeName:
		if n.Names != nil {
			Walk(f, n.Names)
		}
		if n.Typmods != nil {
			Walk(f, n.Typmods)
		}
		if n.ArrayBounds != nil {
			Walk(f, n.ArrayBounds)
		}

	case *ast.UnlistenStmt:
		// pass

	case *ast.UpdateStmt:
		if n.Relations != nil {
			Walk(f, n.Relations)
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
		if n.LimitCount != nil {
			Walk(f, n.LimitCount)
		}
		if n.ReturningList != nil {
			Walk(f, n.ReturningList)
		}
		if n.WithClause != nil {
			Walk(f, n.WithClause)
		}

	case *ast.VacuumStmt:
		if n.Relation != nil {
			Walk(f, n.Relation)
		}
		if n.VaCols != nil {
			Walk(f, n.VaCols)
		}

	case *ast.Var:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}

	case *ast.VariableSetStmt:
		if n.Args != nil {
			Walk(f, n.Args)
		}

	case *ast.VariableShowStmt:
		// pass

	case *ast.ViewStmt:
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

	case *ast.WindowClause:
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

	case *ast.WindowFunc:
		if n.Xpr != nil {
			Walk(f, n.Xpr)
		}
		if n.Args != nil {
			Walk(f, n.Args)
		}
		if n.Aggfilter != nil {
			Walk(f, n.Aggfilter)
		}

	case *ast.WithCheckOption:
		if n.Qual != nil {
			Walk(f, n.Qual)
		}

	case *ast.WithClause:
		if n.Ctes != nil {
			Walk(f, n.Ctes)
		}

	case *ast.XmlExpr:
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

	case *ast.XmlSerialize:
		if n.Expr != nil {
			Walk(f, n.Expr)
		}
		if n.TypeName != nil {
			Walk(f, n.TypeName)
		}

	case *ast.In:
		for _, l := range n.List {
			Walk(f, l)
		}
		if n.Sel != nil {
			Walk(f, n.Sel)
		}

	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))
	}

	f.Visit(nil)
}
