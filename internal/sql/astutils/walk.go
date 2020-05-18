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

func walkn(f Visitor, node ast.Node) {
	if node != nil {
		Walk(f, node)
	}
}

func Walk(f Visitor, node ast.Node) {
	if f = f.Visit(node); f == nil {
		return
	}
	switch n := node.(type) {

	case *ast.AlterTableCmd:
		walkn(f, n.Def)

	case *ast.AlterTableSetSchemaStmt:
		walkn(f, n.Table)

	case *ast.AlterTableStmt:
		walkn(f, n.Table)
		walkn(f, n.Cmds)

	case *ast.AlterTypeAddValueStmt:
		walkn(f, n.Type)

	case *ast.AlterTypeRenameValueStmt:
		walkn(f, n.Type)

	case *ast.ColumnDef:
		walkn(f, n.TypeName)

	case *ast.ColumnRef:
		// pass

	case *ast.CommentOnColumnStmt:
		walkn(f, n.Table)
		walkn(f, n.Col)

	case *ast.CommentOnSchemaStmt:
		walkn(f, n.Schema)

	case *ast.CommentOnTableStmt:
		walkn(f, n.Table)

	case *ast.CommentOnTypeStmt:
		walkn(f, n.Type)

	case *ast.CreateEnumStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Vals)

	case *ast.CreateFunctionStmt:
		walkn(f, n.ReturnType)
		walkn(f, n.Func)

	case *ast.CreateSchemaStmt:
		// pass

	case *ast.CreateTableStmt:
		walkn(f, n.Name)

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
		walkn(f, n.Type)
		walkn(f, n.DefExpr)

	case *ast.FuncSpec:
		walkn(f, n.Name)

	case *ast.List:
		// pass

	case *ast.RawStmt:
		walkn(f, n.Stmt)

	case *ast.RenameColumnStmt:
		walkn(f, n.Table)
		walkn(f, n.Col)

	case *ast.RenameTableStmt:
		walkn(f, n.Table)

	case *ast.ResTarget:
		walkn(f, n.Val)

	case *ast.SelectStmt:
		walkn(f, n.Fields)
		walkn(f, n.From)

	case *ast.Statement:
		walkn(f, n.Raw)

	case *ast.String:
		// pass

	case *ast.TODO:
		// pass

	case *ast.TableName:
		// pass

	case *ast.TypeName:
		// pass

	case *pg.A_ArrayExpr:
		walkn(f, n.Elements)

	case *pg.A_Const:
		walkn(f, n.Val)

	case *pg.A_Expr:
		walkn(f, n.Name)
		walkn(f, n.Lexpr)
		walkn(f, n.Rexpr)

	case *pg.A_Indices:
		walkn(f, n.Lidx)
		walkn(f, n.Uidx)

	case *pg.A_Indirection:
		walkn(f, n.Arg)
		walkn(f, n.Indirection)

	case *pg.A_Star:
		// pass

	case *pg.AccessPriv:
		walkn(f, n.Cols)

	case *pg.Aggref:
		walkn(f, n.Xpr)
		walkn(f, n.Aggargtypes)
		walkn(f, n.Aggdirectargs)
		walkn(f, n.Args)
		walkn(f, n.Aggorder)
		walkn(f, n.Aggdistinct)
		walkn(f, n.Aggfilter)

	case *pg.Alias:
		walkn(f, n.Colnames)

	case *pg.AlterCollationStmt:
		walkn(f, n.Collname)

	case *pg.AlterDatabaseSetStmt:
		walkn(f, n.Setstmt)

	case *pg.AlterDatabaseStmt:
		walkn(f, n.Options)

	case *pg.AlterDefaultPrivilegesStmt:
		walkn(f, n.Options)
		walkn(f, n.Action)

	case *pg.AlterDomainStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Def)

	case *pg.AlterEnumStmt:
		walkn(f, n.TypeName)

	case *pg.AlterEventTrigStmt:
		// pass

	case *pg.AlterExtensionContentsStmt:
		walkn(f, n.Object)

	case *pg.AlterExtensionStmt:
		walkn(f, n.Options)

	case *pg.AlterFdwStmt:
		walkn(f, n.FuncOptions)
		walkn(f, n.Options)

	case *pg.AlterForeignServerStmt:
		walkn(f, n.Options)

	case *pg.AlterFunctionStmt:
		walkn(f, n.Func)
		walkn(f, n.Actions)

	case *pg.AlterObjectDependsStmt:
		walkn(f, n.Relation)
		walkn(f, n.Object)
		walkn(f, n.Extname)

	case *pg.AlterObjectSchemaStmt:
		walkn(f, n.Relation)
		walkn(f, n.Object)

	case *pg.AlterOpFamilyStmt:
		walkn(f, n.Opfamilyname)
		walkn(f, n.Items)

	case *pg.AlterOperatorStmt:
		walkn(f, n.Opername)
		walkn(f, n.Options)

	case *pg.AlterOwnerStmt:
		walkn(f, n.Relation)
		walkn(f, n.Object)
		walkn(f, n.Newowner)

	case *pg.AlterPolicyStmt:
		walkn(f, n.Table)
		walkn(f, n.Roles)
		walkn(f, n.Qual)
		walkn(f, n.WithCheck)

	case *pg.AlterPublicationStmt:
		walkn(f, n.Options)
		walkn(f, n.Tables)

	case *pg.AlterRoleSetStmt:
		walkn(f, n.Role)
		walkn(f, n.Setstmt)

	case *pg.AlterRoleStmt:
		walkn(f, n.Role)
		walkn(f, n.Options)

	case *pg.AlterSeqStmt:
		walkn(f, n.Sequence)
		walkn(f, n.Options)

	case *pg.AlterSubscriptionStmt:
		walkn(f, n.Publication)
		walkn(f, n.Options)

	case *pg.AlterSystemStmt:
		walkn(f, n.Setstmt)

	case *pg.AlterTSConfigurationStmt:
		walkn(f, n.Cfgname)
		walkn(f, n.Tokentype)
		walkn(f, n.Dicts)

	case *pg.AlterTSDictionaryStmt:
		walkn(f, n.Dictname)
		walkn(f, n.Options)

	case *pg.AlterTableCmd:
		walkn(f, n.Newowner)
		walkn(f, n.Def)

	case *pg.AlterTableMoveAllStmt:
		walkn(f, n.Roles)

	case *pg.AlterTableSpaceOptionsStmt:
		walkn(f, n.Options)

	case *pg.AlterTableStmt:
		walkn(f, n.Relation)
		walkn(f, n.Cmds)

	case *pg.AlterUserMappingStmt:
		walkn(f, n.User)
		walkn(f, n.Options)

	case *pg.AlternativeSubPlan:
		walkn(f, n.Xpr)
		walkn(f, n.Subplans)

	case *pg.ArrayCoerceExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.ArrayExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Elements)

	case *pg.ArrayRef:
		walkn(f, n.Xpr)
		walkn(f, n.Refupperindexpr)
		walkn(f, n.Reflowerindexpr)
		walkn(f, n.Refexpr)
		walkn(f, n.Refassgnexpr)

	case *pg.BitString:
		// pass

	case *pg.BlockIdData:
		// pass

	case *pg.BoolExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.BooleanTest:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.CaseExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)
		walkn(f, n.Args)
		walkn(f, n.Defresult)

	case *pg.CaseTestExpr:
		walkn(f, n.Xpr)

	case *pg.CaseWhen:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)
		walkn(f, n.Result)

	case *pg.CheckPointStmt:
		// pass

	case *pg.ClosePortalStmt:
		// pass

	case *pg.ClusterStmt:
		walkn(f, n.Relation)

	case *pg.CoalesceExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.CoerceToDomain:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.CoerceToDomainValue:
		walkn(f, n.Xpr)

	case *pg.CoerceViaIO:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.CollateClause:
		walkn(f, n.Arg)
		walkn(f, n.Collname)

	case *pg.CollateExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.ColumnDef:
		walkn(f, n.TypeName)
		walkn(f, n.RawDefault)
		walkn(f, n.CookedDefault)
		walkn(f, n.CollClause)
		walkn(f, n.Constraints)
		walkn(f, n.Fdwoptions)

	case *pg.ColumnRef:
		walkn(f, n.Fields)

	case *pg.CommentStmt:
		walkn(f, n.Object)

	case *pg.CommonTableExpr:
		walkn(f, n.Aliascolnames)
		walkn(f, n.Ctequery)
		walkn(f, n.Ctecolnames)
		walkn(f, n.Ctecoltypes)
		walkn(f, n.Ctecoltypmods)
		walkn(f, n.Ctecolcollations)

	case *pg.CompositeTypeStmt:
		walkn(f, n.Typevar)
		walkn(f, n.Coldeflist)

	case *pg.Const:
		walkn(f, n.Xpr)

	case *pg.Constraint:
		walkn(f, n.RawExpr)
		walkn(f, n.Keys)
		walkn(f, n.Exclusions)
		walkn(f, n.Options)
		walkn(f, n.WhereClause)
		walkn(f, n.Pktable)
		walkn(f, n.FkAttrs)
		walkn(f, n.PkAttrs)
		walkn(f, n.OldConpfeqop)

	case *pg.ConstraintsSetStmt:
		walkn(f, n.Constraints)

	case *pg.ConvertRowtypeExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.CopyStmt:
		walkn(f, n.Relation)
		walkn(f, n.Query)
		walkn(f, n.Attlist)
		walkn(f, n.Options)

	case *pg.CreateAmStmt:
		walkn(f, n.HandlerName)

	case *pg.CreateCastStmt:
		walkn(f, n.Sourcetype)
		walkn(f, n.Targettype)
		walkn(f, n.Func)

	case *pg.CreateConversionStmt:
		walkn(f, n.ConversionName)
		walkn(f, n.FuncName)

	case *pg.CreateDomainStmt:
		walkn(f, n.Domainname)
		walkn(f, n.TypeName)
		walkn(f, n.CollClause)
		walkn(f, n.Constraints)

	case *pg.CreateEnumStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Vals)

	case *pg.CreateEventTrigStmt:
		walkn(f, n.Whenclause)
		walkn(f, n.Funcname)

	case *pg.CreateExtensionStmt:
		walkn(f, n.Options)

	case *pg.CreateFdwStmt:
		walkn(f, n.FuncOptions)
		walkn(f, n.Options)

	case *pg.CreateForeignServerStmt:
		walkn(f, n.Options)

	case *pg.CreateForeignTableStmt:
		walkn(f, n.Base)
		walkn(f, n.Options)

	case *pg.CreateFunctionStmt:
		walkn(f, n.Funcname)
		walkn(f, n.Parameters)
		walkn(f, n.ReturnType)
		walkn(f, n.Options)
		walkn(f, n.WithClause)

	case *pg.CreateOpClassItem:
		walkn(f, n.Name)
		walkn(f, n.OrderFamily)
		walkn(f, n.ClassArgs)
		walkn(f, n.Storedtype)

	case *pg.CreateOpClassStmt:
		walkn(f, n.Opclassname)
		walkn(f, n.Opfamilyname)
		walkn(f, n.Datatype)
		walkn(f, n.Items)

	case *pg.CreateOpFamilyStmt:
		walkn(f, n.Opfamilyname)

	case *pg.CreatePLangStmt:
		walkn(f, n.Plhandler)
		walkn(f, n.Plinline)
		walkn(f, n.Plvalidator)

	case *pg.CreatePolicyStmt:
		walkn(f, n.Table)
		walkn(f, n.Roles)
		walkn(f, n.Qual)
		walkn(f, n.WithCheck)

	case *pg.CreatePublicationStmt:
		walkn(f, n.Options)
		walkn(f, n.Tables)

	case *pg.CreateRangeStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Params)

	case *pg.CreateRoleStmt:
		walkn(f, n.Options)

	case *pg.CreateSchemaStmt:
		walkn(f, n.Authrole)
		walkn(f, n.SchemaElts)

	case *pg.CreateSeqStmt:
		walkn(f, n.Sequence)
		walkn(f, n.Options)

	case *pg.CreateStatsStmt:
		walkn(f, n.Defnames)
		walkn(f, n.StatTypes)
		walkn(f, n.Exprs)
		walkn(f, n.Relations)

	case *pg.CreateStmt:
		walkn(f, n.Relation)
		walkn(f, n.TableElts)
		walkn(f, n.InhRelations)
		walkn(f, n.Partbound)
		walkn(f, n.Partspec)
		walkn(f, n.OfTypename)
		walkn(f, n.Constraints)
		walkn(f, n.Options)

	case *pg.CreateSubscriptionStmt:
		walkn(f, n.Publication)
		walkn(f, n.Options)

	case *pg.CreateTableAsStmt:
		walkn(f, n.Query)
		walkn(f, n.Into)

	case *pg.CreateTableSpaceStmt:
		walkn(f, n.Owner)
		walkn(f, n.Options)

	case *pg.CreateTransformStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Fromsql)
		walkn(f, n.Tosql)

	case *pg.CreateTrigStmt:
		walkn(f, n.Relation)
		walkn(f, n.Funcname)
		walkn(f, n.Args)
		walkn(f, n.Columns)
		walkn(f, n.WhenClause)
		walkn(f, n.TransitionRels)
		walkn(f, n.Constrrel)

	case *pg.CreateUserMappingStmt:
		walkn(f, n.User)
		walkn(f, n.Options)

	case *pg.CreatedbStmt:
		walkn(f, n.Options)

	case *pg.CurrentOfExpr:
		walkn(f, n.Xpr)

	case *pg.DeallocateStmt:
		// pass

	case *pg.DeclareCursorStmt:
		walkn(f, n.Query)

	case *pg.DefElem:
		walkn(f, n.Arg)

	case *pg.DefineStmt:
		walkn(f, n.Defnames)
		walkn(f, n.Args)
		walkn(f, n.Definition)

	case *pg.DeleteStmt:
		walkn(f, n.Relation)
		walkn(f, n.UsingClause)
		walkn(f, n.WhereClause)
		walkn(f, n.ReturningList)
		walkn(f, n.WithClause)

	case *pg.DiscardStmt:
		// pass

	case *pg.DoStmt:
		walkn(f, n.Args)

	case *pg.DropOwnedStmt:
		walkn(f, n.Roles)

	case *pg.DropRoleStmt:
		walkn(f, n.Roles)

	case *pg.DropStmt:
		walkn(f, n.Objects)

	case *pg.DropSubscriptionStmt:
		// pass

	case *pg.DropTableSpaceStmt:
		// pass

	case *pg.DropUserMappingStmt:
		walkn(f, n.User)

	case *pg.DropdbStmt:
		// pass

	case *pg.ExecuteStmt:
		walkn(f, n.Params)

	case *pg.ExplainStmt:
		walkn(f, n.Query)
		walkn(f, n.Options)

	case *pg.Expr:
		// pass

	case *pg.FetchStmt:
		// pass

	case *pg.FieldSelect:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.FieldStore:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)
		walkn(f, n.Newvals)
		walkn(f, n.Fieldnums)

	case *pg.Float:
		// pass

	case *pg.FromExpr:
		walkn(f, n.Fromlist)
		walkn(f, n.Quals)

	case *ast.FuncCall:
		walkn(f, n.Func)
		walkn(f, n.Funcname)
		walkn(f, n.Args)
		walkn(f, n.AggOrder)
		walkn(f, n.AggFilter)
		walkn(f, n.Over)

	case *pg.FuncExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.FunctionParameter:
		walkn(f, n.ArgType)
		walkn(f, n.Defexpr)

	case *pg.GrantRoleStmt:
		walkn(f, n.GrantedRoles)
		walkn(f, n.GranteeRoles)
		walkn(f, n.Grantor)

	case *pg.GrantStmt:
		walkn(f, n.Objects)
		walkn(f, n.Privileges)
		walkn(f, n.Grantees)

	case *pg.GroupingFunc:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Refs)
		walkn(f, n.Cols)

	case *pg.GroupingSet:
		walkn(f, n.Content)

	case *pg.ImportForeignSchemaStmt:
		walkn(f, n.TableList)
		walkn(f, n.Options)

	case *pg.IndexElem:
		walkn(f, n.Expr)
		walkn(f, n.Collation)
		walkn(f, n.Opclass)

	case *pg.IndexStmt:
		walkn(f, n.Relation)
		walkn(f, n.IndexParams)
		walkn(f, n.Options)
		walkn(f, n.WhereClause)
		walkn(f, n.ExcludeOpNames)

	case *pg.InferClause:
		walkn(f, n.IndexElems)
		walkn(f, n.WhereClause)

	case *pg.InferenceElem:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)

	case *pg.InlineCodeBlock:
		// pass

	case *pg.InsertStmt:
		walkn(f, n.Relation)
		walkn(f, n.Cols)
		walkn(f, n.SelectStmt)
		walkn(f, n.OnConflictClause)
		walkn(f, n.ReturningList)
		walkn(f, n.WithClause)

	case *pg.Integer:
		// pass

	case *pg.IntoClause:
		walkn(f, n.Rel)
		walkn(f, n.ColNames)
		walkn(f, n.Options)
		walkn(f, n.ViewQuery)

	case *pg.JoinExpr:
		walkn(f, n.Larg)
		walkn(f, n.Rarg)
		walkn(f, n.UsingClause)
		walkn(f, n.Quals)
		walkn(f, n.Alias)

	case *pg.ListenStmt:
		// pass

	case *pg.LoadStmt:
		// pass

	case *pg.LockStmt:
		walkn(f, n.Relations)

	case *pg.LockingClause:
		walkn(f, n.LockedRels)

	case *pg.MinMaxExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.MultiAssignRef:
		walkn(f, n.Source)

	case *pg.NamedArgExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.NextValueExpr:
		walkn(f, n.Xpr)

	case *pg.NotifyStmt:
		// pass

	case *pg.Null:
		// pass

	case *pg.NullTest:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.ObjectWithArgs:
		walkn(f, n.Objname)
		walkn(f, n.Objargs)

	case *pg.OnConflictClause:
		walkn(f, n.Infer)
		walkn(f, n.TargetList)
		walkn(f, n.WhereClause)

	case *pg.OnConflictExpr:
		walkn(f, n.ArbiterElems)
		walkn(f, n.ArbiterWhere)
		walkn(f, n.OnConflictSet)
		walkn(f, n.OnConflictWhere)
		walkn(f, n.ExclRelTlist)

	case *pg.OpExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.Param:
		walkn(f, n.Xpr)

	case *pg.ParamExecData:
		// pass

	case *pg.ParamExternData:
		// pass

	case *pg.ParamListInfoData:
		// pass

	case *pg.ParamRef:
		// pass

	case *pg.PartitionBoundSpec:
		walkn(f, n.Listdatums)
		walkn(f, n.Lowerdatums)
		walkn(f, n.Upperdatums)

	case *pg.PartitionCmd:
		walkn(f, n.Name)
		walkn(f, n.Bound)

	case *pg.PartitionElem:
		walkn(f, n.Expr)
		walkn(f, n.Collation)
		walkn(f, n.Opclass)

	case *pg.PartitionRangeDatum:
		walkn(f, n.Value)

	case *pg.PartitionSpec:
		walkn(f, n.PartParams)

	case *pg.PrepareStmt:
		walkn(f, n.Argtypes)
		walkn(f, n.Query)

	case *pg.Query:
		walkn(f, n.UtilityStmt)
		walkn(f, n.CteList)
		walkn(f, n.Rtable)
		walkn(f, n.Jointree)
		walkn(f, n.TargetList)
		walkn(f, n.OnConflict)
		walkn(f, n.ReturningList)
		walkn(f, n.GroupClause)
		walkn(f, n.GroupingSets)
		walkn(f, n.HavingQual)
		walkn(f, n.WindowClause)
		walkn(f, n.DistinctClause)
		walkn(f, n.SortClause)
		walkn(f, n.LimitOffset)
		walkn(f, n.LimitCount)
		walkn(f, n.RowMarks)
		walkn(f, n.SetOperations)
		walkn(f, n.ConstraintDeps)
		walkn(f, n.WithCheckOptions)

	case *pg.RangeFunction:
		walkn(f, n.Functions)
		walkn(f, n.Alias)
		walkn(f, n.Coldeflist)

	case *pg.RangeSubselect:
		walkn(f, n.Subquery)
		walkn(f, n.Alias)

	case *pg.RangeTableFunc:
		walkn(f, n.Docexpr)
		walkn(f, n.Rowexpr)
		walkn(f, n.Namespaces)
		walkn(f, n.Columns)
		walkn(f, n.Alias)

	case *pg.RangeTableFuncCol:
		walkn(f, n.TypeName)
		walkn(f, n.Colexpr)
		walkn(f, n.Coldefexpr)

	case *pg.RangeTableSample:
		walkn(f, n.Relation)
		walkn(f, n.Method)
		walkn(f, n.Args)
		walkn(f, n.Repeatable)

	case *pg.RangeTblEntry:
		walkn(f, n.Tablesample)
		walkn(f, n.Subquery)
		walkn(f, n.Joinaliasvars)
		walkn(f, n.Functions)
		walkn(f, n.Tablefunc)
		walkn(f, n.ValuesLists)
		walkn(f, n.Coltypes)
		walkn(f, n.Coltypmods)
		walkn(f, n.Colcollations)
		walkn(f, n.Alias)
		walkn(f, n.Eref)
		walkn(f, n.SecurityQuals)

	case *pg.RangeTblFunction:
		walkn(f, n.Funcexpr)
		walkn(f, n.Funccolnames)
		walkn(f, n.Funccoltypes)
		walkn(f, n.Funccoltypmods)
		walkn(f, n.Funccolcollations)

	case *pg.RangeTblRef:
		// pass

	case *pg.RangeVar:
		walkn(f, n.Alias)

	case *pg.RawStmt:
		walkn(f, n.Stmt)

	case *pg.ReassignOwnedStmt:
		walkn(f, n.Roles)
		walkn(f, n.Newrole)

	case *pg.RefreshMatViewStmt:
		walkn(f, n.Relation)

	case *pg.ReindexStmt:
		walkn(f, n.Relation)

	case *pg.RelabelType:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case *pg.RenameStmt:
		walkn(f, n.Relation)
		walkn(f, n.Object)

	case *pg.ReplicaIdentityStmt:
		// pass

	case *pg.ResTarget:
		walkn(f, n.Indirection)
		walkn(f, n.Val)

	case *pg.RoleSpec:
		// pass

	case *pg.RowCompareExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Opnos)
		walkn(f, n.Opfamilies)
		walkn(f, n.Inputcollids)
		walkn(f, n.Largs)
		walkn(f, n.Rargs)

	case *pg.RowExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Colnames)

	case *pg.RowMarkClause:
		// pass

	case *pg.RuleStmt:
		walkn(f, n.Relation)
		walkn(f, n.WhereClause)
		walkn(f, n.Actions)

	case *pg.SQLValueFunction:
		walkn(f, n.Xpr)

	case *pg.ScalarArrayOpExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case *pg.SecLabelStmt:
		walkn(f, n.Object)

	case *pg.SelectStmt:
		walkn(f, n.DistinctClause)
		walkn(f, n.IntoClause)
		walkn(f, n.TargetList)
		walkn(f, n.FromClause)
		walkn(f, n.WhereClause)
		walkn(f, n.GroupClause)
		walkn(f, n.HavingClause)
		walkn(f, n.WindowClause)
		walkn(f, n.SortClause)
		walkn(f, n.LimitOffset)
		walkn(f, n.LimitCount)
		walkn(f, n.LockingClause)
		walkn(f, n.WithClause)
		walkn(f, n.Larg)
		walkn(f, n.Rarg)

	case *pg.SetOperationStmt:
		walkn(f, n.Larg)
		walkn(f, n.Rarg)
		walkn(f, n.ColTypes)
		walkn(f, n.ColTypmods)
		walkn(f, n.ColCollations)
		walkn(f, n.GroupClauses)

	case *pg.SetToDefault:
		walkn(f, n.Xpr)

	case *pg.SortBy:
		walkn(f, n.Node)
		walkn(f, n.UseOp)

	case *pg.SortGroupClause:
		// pass

	case *pg.String:
		// pass

	case *pg.SubLink:
		walkn(f, n.Xpr)
		walkn(f, n.Testexpr)
		walkn(f, n.OperName)
		walkn(f, n.Subselect)

	case *pg.SubPlan:
		walkn(f, n.Xpr)
		walkn(f, n.Testexpr)
		walkn(f, n.ParamIds)
		walkn(f, n.SetParam)
		walkn(f, n.ParParam)
		walkn(f, n.Args)

	case *pg.TableFunc:
		walkn(f, n.NsUris)
		walkn(f, n.NsNames)
		walkn(f, n.Docexpr)
		walkn(f, n.Rowexpr)
		walkn(f, n.Colnames)
		walkn(f, n.Coltypes)
		walkn(f, n.Coltypmods)
		walkn(f, n.Colcollations)
		walkn(f, n.Colexprs)
		walkn(f, n.Coldefexprs)

	case *pg.TableLikeClause:
		walkn(f, n.Relation)

	case *pg.TableSampleClause:
		walkn(f, n.Args)
		walkn(f, n.Repeatable)

	case *pg.TargetEntry:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)

	case *pg.TransactionStmt:
		walkn(f, n.Options)

	case *pg.TriggerTransition:
		// pass

	case *pg.TruncateStmt:
		walkn(f, n.Relations)

	case *pg.TypeCast:
		walkn(f, n.Arg)
		walkn(f, n.TypeName)

	case *pg.TypeName:
		walkn(f, n.Names)
		walkn(f, n.Typmods)
		walkn(f, n.ArrayBounds)

	case *pg.UnlistenStmt:
		// pass

	case *pg.UpdateStmt:
		walkn(f, n.Relation)
		walkn(f, n.TargetList)
		walkn(f, n.WhereClause)
		walkn(f, n.FromClause)
		walkn(f, n.ReturningList)
		walkn(f, n.WithClause)

	case *pg.VacuumStmt:
		walkn(f, n.Relation)
		walkn(f, n.VaCols)

	case *pg.Var:
		walkn(f, n.Xpr)

	case *pg.VariableSetStmt:
		walkn(f, n.Args)

	case *pg.VariableShowStmt:
		// pass

	case *pg.ViewStmt:
		walkn(f, n.View)
		walkn(f, n.Aliases)
		walkn(f, n.Query)
		walkn(f, n.Options)

	case *pg.WindowClause:
		walkn(f, n.PartitionClause)
		walkn(f, n.OrderClause)
		walkn(f, n.StartOffset)
		walkn(f, n.EndOffset)

	case *ast.WindowDef:
		walkn(f, n.PartitionClause)
		walkn(f, n.OrderClause)
		walkn(f, n.StartOffset)
		walkn(f, n.EndOffset)

	case *pg.WindowFunc:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Aggfilter)

	case *pg.WithCheckOption:
		walkn(f, n.Qual)

	case *pg.WithClause:
		walkn(f, n.Ctes)

	case *pg.XmlExpr:
		walkn(f, n.Xpr)
		walkn(f, n.NamedArgs)
		walkn(f, n.ArgNames)
		walkn(f, n.Args)

	case *pg.XmlSerialize:
		walkn(f, n.Expr)
		walkn(f, n.TypeName)

	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))
	}

	f.Visit(nil)
}
