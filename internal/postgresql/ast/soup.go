package ast

import (
	"fmt"

	nodes "github.com/lfittl/pg_query_go/nodes"
)

type Visitor interface {
	Visit(nodes.Node) Visitor
}

type VisitorFunc func(nodes.Node)

func (vf VisitorFunc) Visit(node nodes.Node) Visitor {
	vf(node)
	return vf
}

func walkn(f Visitor, node nodes.Node) {
	if node != nil {
		Walk(f, node)
	}
}

func Walk(f Visitor, node nodes.Node) {
	if f = f.Visit(node); f == nil {
		return
	}

	switch n := node.(type) {

	case nodes.A_ArrayExpr:
		walkn(f, n.Elements)

	case nodes.A_Const:
		walkn(f, n.Val)

	case nodes.A_Expr:
		walkn(f, n.Name)
		walkn(f, n.Lexpr)
		walkn(f, n.Rexpr)

	case nodes.A_Indices:
		walkn(f, n.Lidx)
		walkn(f, n.Uidx)

	case nodes.A_Indirection:
		walkn(f, n.Arg)
		walkn(f, n.Indirection)

	case nodes.A_Star:
		// pass

	case nodes.AccessPriv:
		walkn(f, n.Cols)

	case nodes.Aggref:
		walkn(f, n.Xpr)
		walkn(f, n.Aggargtypes)
		walkn(f, n.Aggdirectargs)
		walkn(f, n.Args)
		walkn(f, n.Aggorder)
		walkn(f, n.Aggdistinct)
		walkn(f, n.Aggfilter)

	case nodes.Alias:
		walkn(f, n.Colnames)

	case nodes.AlterCollationStmt:
		walkn(f, n.Collname)

	case nodes.AlterDatabaseSetStmt:
		if n.Setstmt != nil {
			walkn(f, *n.Setstmt)
		}

	case nodes.AlterDatabaseStmt:
		walkn(f, n.Options)

	case nodes.AlterDefaultPrivilegesStmt:
		if n.Action != nil {
			walkn(f, *n.Action)
		}
		walkn(f, n.Options)

	case nodes.AlterDomainStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Def)

	case nodes.AlterEnumStmt:
		walkn(f, n.TypeName)

	case nodes.AlterEventTrigStmt:
		// pass

	case nodes.AlterExtensionContentsStmt:
		walkn(f, n.Object)

	case nodes.AlterExtensionStmt:
		walkn(f, n.Options)

	case nodes.AlterFdwStmt:
		walkn(f, n.FuncOptions)
		walkn(f, n.Options)

	case nodes.AlterForeignServerStmt:
		walkn(f, n.Options)

	case nodes.AlterFunctionStmt:
		if n.Func != nil {
			walkn(f, n.Func)
		}
		walkn(f, n.Actions)

	case nodes.AlterObjectDependsStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Object)
		walkn(f, n.Extname)

	case nodes.AlterObjectSchemaStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Object)

	case nodes.AlterOpFamilyStmt:
		walkn(f, n.Opfamilyname)
		walkn(f, n.Items)

	case nodes.AlterOperatorStmt:
		if n.Opername != nil {
			walkn(f, *n.Opername)
		}
		walkn(f, n.Options)

	case nodes.AlterOwnerStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Object)
		if n.Newowner != nil {
			walkn(f, *n.Newowner)
		}

	case nodes.AlterPolicyStmt:
		if n.Table != nil {
			walkn(f, *n.Table)
		}
		walkn(f, n.Roles)
		walkn(f, n.Qual)
		walkn(f, n.WithCheck)

	case nodes.AlterPublicationStmt:
		walkn(f, n.Options)
		walkn(f, n.Tables)

	case nodes.AlterRoleSetStmt:
		if n.Role != nil {
			walkn(f, *n.Role)
		}
		walkn(f, n.Setstmt)

	case nodes.AlterRoleStmt:
		if n.Role != nil {
			walkn(f, *n.Role)
		}
		walkn(f, n.Options)

	case nodes.AlterSeqStmt:
		if n.Sequence != nil {
			walkn(f, *n.Sequence)
		}
		walkn(f, n.Options)

	case nodes.AlterSubscriptionStmt:
		walkn(f, n.Publication)
		walkn(f, n.Options)

	case nodes.AlterSystemStmt:
		walkn(f, n.Setstmt)

	case nodes.AlterTSConfigurationStmt:
		walkn(f, n.Cfgname)
		walkn(f, n.Tokentype)
		walkn(f, n.Dicts)

	case nodes.AlterTSDictionaryStmt:
		walkn(f, n.Dictname)
		walkn(f, n.Options)

	case nodes.AlterTableCmd:
		if n.Newowner != nil {
			walkn(f, *n.Newowner)
		}
		walkn(f, n.Def)

	case nodes.AlterTableMoveAllStmt:
		walkn(f, n.Roles)

	case nodes.AlterTableSpaceOptionsStmt:
		walkn(f, n.Options)

	case nodes.AlterTableStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Cmds)

	case nodes.AlterUserMappingStmt:
		if n.User != nil {
			walkn(f, *n.User)
		}
		walkn(f, n.Options)

	case nodes.AlternativeSubPlan:
		walkn(f, n.Xpr)
		walkn(f, n.Subplans)

	case nodes.ArrayCoerceExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.ArrayExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Elements)

	case nodes.ArrayRef:
		walkn(f, n.Xpr)
		walkn(f, n.Refupperindexpr)
		walkn(f, n.Reflowerindexpr)
		walkn(f, n.Refexpr)
		walkn(f, n.Refassgnexpr)

	case nodes.BitString:
		// pass

	case nodes.BlockIdData:
		// pass

	case nodes.BoolExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.BooleanTest:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.CaseExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)
		walkn(f, n.Args)
		walkn(f, n.Defresult)

	case nodes.CaseTestExpr:
		walkn(f, n.Xpr)

	case nodes.CaseWhen:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)
		walkn(f, n.Result)

	case nodes.CheckPointStmt:
		// pass

	case nodes.ClosePortalStmt:
		// pass

	case nodes.ClusterStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}

	case nodes.CoalesceExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.CoerceToDomain:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.CoerceToDomainValue:
		walkn(f, n.Xpr)

	case nodes.CoerceViaIO:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.CollateClause:
		walkn(f, n.Arg)
		walkn(f, n.Collname)

	case nodes.CollateExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.ColumnDef:
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}
		walkn(f, n.RawDefault)
		walkn(f, n.CookedDefault)
		walkn(f, n.Constraints)
		walkn(f, n.Fdwoptions)

	case nodes.ColumnRef:
		walkn(f, n.Fields)

	case nodes.CommentStmt:
		walkn(f, n.Object)

	case nodes.CommonTableExpr:
		walkn(f, n.Aliascolnames)
		walkn(f, n.Ctequery)
		walkn(f, n.Ctecolnames)
		walkn(f, n.Ctecolcollations)

	case nodes.CompositeTypeStmt:
		if n.Typevar != nil {
			walkn(f, *n.Typevar)
		}
		walkn(f, n.Coldeflist)

	case nodes.Const:
		walkn(f, n.Xpr)

	case nodes.Constraint:
		walkn(f, n.RawExpr)
		walkn(f, n.Keys)
		walkn(f, n.Exclusions)
		walkn(f, n.Options)
		walkn(f, n.WhereClause)
		if n.Pktable != nil {
			walkn(f, *n.Pktable)
		}
		walkn(f, n.FkAttrs)
		walkn(f, n.PkAttrs)
		walkn(f, n.OldConpfeqop)

	case nodes.ConstraintsSetStmt:
		walkn(f, n.Constraints)

	case nodes.ConvertRowtypeExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.CopyStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Query)
		walkn(f, n.Attlist)
		walkn(f, n.Options)

	case nodes.CreateAmStmt:
		walkn(f, n.HandlerName)

	case nodes.CreateCastStmt:
		if n.Sourcetype != nil {
			walkn(f, *n.Sourcetype)
		}
		if n.Targettype != nil {
			walkn(f, *n.Targettype)
		}
		walkn(f, n.Func)

	case nodes.CreateConversionStmt:
		walkn(f, n.ConversionName)
		walkn(f, n.FuncName)

	case nodes.CreateDomainStmt:
		walkn(f, n.Domainname)
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}
		if n.CollClause != nil {
			walkn(f, *n.CollClause)
		}
		walkn(f, n.Constraints)

	case nodes.CreateEnumStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Vals)

	case nodes.CreateEventTrigStmt:
		walkn(f, n.Whenclause)
		walkn(f, n.Funcname)

	case nodes.CreateExtensionStmt:
		walkn(f, n.Options)

	case nodes.CreateFdwStmt:
		walkn(f, n.FuncOptions)
		walkn(f, n.Options)

	case nodes.CreateForeignServerStmt:
		walkn(f, n.Options)

	case nodes.CreateForeignTableStmt:
		walkn(f, n.Base)
		walkn(f, n.Options)

	case nodes.CreateFunctionStmt:
		walkn(f, n.Funcname)
		walkn(f, n.Parameters)
		if n.ReturnType != nil {
			walkn(f, *n.ReturnType)
		}
		walkn(f, n.Options)
		walkn(f, n.WithClause)

	case nodes.CreateOpClassItem:
		walkn(f, n.Name)
		walkn(f, n.OrderFamily)
		walkn(f, n.ClassArgs)
		if n.Storedtype != nil {
			walkn(f, *n.Storedtype)
		}

	case nodes.CreateOpClassStmt:
		walkn(f, n.Opclassname)
		walkn(f, n.Opfamilyname)
		if n.Datatype != nil {
			walkn(f, *n.Datatype)
		}
		walkn(f, n.Items)

	case nodes.CreateOpFamilyStmt:
		walkn(f, n.Opfamilyname)

	case nodes.CreatePLangStmt:
		walkn(f, n.Plhandler)
		walkn(f, n.Plinline)
		walkn(f, n.Plvalidator)

	case nodes.CreatePolicyStmt:
		if n.Table != nil {
			walkn(f, *n.Table)
		}
		walkn(f, n.Roles)
		walkn(f, n.Qual)
		walkn(f, n.WithCheck)

	case nodes.CreatePublicationStmt:
		walkn(f, n.Options)
		walkn(f, n.Tables)

	case nodes.CreateRangeStmt:
		walkn(f, n.TypeName)
		walkn(f, n.Params)

	case nodes.CreateRoleStmt:
		walkn(f, n.Options)

	case nodes.CreateSchemaStmt:
		if n.Authrole != nil {
			walkn(f, *n.Authrole)
		}
		walkn(f, n.SchemaElts)

	case nodes.CreateSeqStmt:
		if n.Sequence != nil {
			walkn(f, *n.Sequence)
		}
		walkn(f, n.Options)

	case nodes.CreateStatsStmt:
		walkn(f, n.Defnames)
		walkn(f, n.StatTypes)
		walkn(f, n.Exprs)
		walkn(f, n.Relations)

	case nodes.CreateStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.TableElts)
		walkn(f, n.InhRelations)
		if n.Partbound != nil {
			walkn(f, *n.Partbound)
		}
		if n.Partspec != nil {
			walkn(f, *n.Partspec)
		}
		walkn(f, n.Constraints)
		walkn(f, n.Options)
		if n.OfTypename != nil {
			walkn(f, *n.OfTypename)
		}

	case nodes.CreateSubscriptionStmt:
		walkn(f, n.Publication)
		walkn(f, n.Options)

	case nodes.CreateTableAsStmt:
		walkn(f, n.Query)
		walkn(f, n.Into)

	case nodes.CreateTableSpaceStmt:
		if n.Owner != nil {
			walkn(f, *n.Owner)
		}
		walkn(f, n.Options)

	case nodes.CreateTransformStmt:
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}
		if n.Fromsql != nil {
			walkn(f, *n.Fromsql)
		}
		if n.Tosql != nil {
			walkn(f, *n.Tosql)
		}

	case nodes.CreateTrigStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Funcname)
		walkn(f, n.Args)
		walkn(f, n.Columns)
		walkn(f, n.WhenClause)
		walkn(f, n.TransitionRels)
		if n.Constrrel != nil {
			walkn(f, *n.Constrrel)
		}

	case nodes.CreateUserMappingStmt:
		if n.User != nil {
			walkn(f, *n.User)
		}
		walkn(f, n.Options)

	case nodes.CreatedbStmt:
		walkn(f, n.Options)

	case nodes.CurrentOfExpr:
		walkn(f, n.Xpr)

	case nodes.DeallocateStmt:
		// pass

	case nodes.DeclareCursorStmt:
		walkn(f, n.Query)

	case nodes.DefElem:
		walkn(f, n.Arg)

	case nodes.DefineStmt:
		walkn(f, n.Defnames)
		walkn(f, n.Args)
		walkn(f, n.Definition)

	case nodes.DeleteStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.UsingClause)
		walkn(f, n.WhereClause)
		walkn(f, n.ReturningList)
		if n.WithClause != nil {
			walkn(f, *n.WithClause)
		}

	case nodes.DiscardStmt:
		// pass

	case nodes.DoStmt:
		walkn(f, n.Args)

	case nodes.DropOwnedStmt:
		walkn(f, n.Roles)

	case nodes.DropRoleStmt:
		walkn(f, n.Roles)

	case nodes.DropStmt:
		walkn(f, n.Objects)

	case nodes.DropSubscriptionStmt:
		// pass

	case nodes.DropTableSpaceStmt:
		// pass

	case nodes.DropUserMappingStmt:
		if n.User != nil {
			walkn(f, *n.User)
		}

	case nodes.DropdbStmt:
		// pass

	case nodes.ExecuteStmt:
		walkn(f, n.Params)

	case nodes.ExplainStmt:
		walkn(f, n.Query)
		walkn(f, n.Options)

	case nodes.Expr:
		// pass

	case nodes.FetchStmt:
		// pass

	case nodes.FieldSelect:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.FieldStore:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)
		walkn(f, n.Newvals)
		walkn(f, n.Fieldnums)

	case nodes.Float:
		// pass

	case nodes.FromExpr:
		walkn(f, n.Fromlist)
		walkn(f, n.Quals)

	case nodes.FuncCall:
		walkn(f, n.Funcname)
		walkn(f, n.Args)
		walkn(f, n.AggOrder)
		walkn(f, n.AggFilter)
		if n.Over != nil {
			walkn(f, *n.Over)
		}

	case nodes.FuncExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.FunctionParameter:
		if n.ArgType != nil {
			walkn(f, *n.ArgType)
		}
		walkn(f, n.Defexpr)

	case nodes.GrantRoleStmt:
		walkn(f, n.GrantedRoles)
		walkn(f, n.GranteeRoles)
		if n.Grantor != nil {
			walkn(f, *n.Grantor)
		}

	case nodes.GrantStmt:
		walkn(f, n.Objects)
		walkn(f, n.Privileges)
		walkn(f, n.Grantees)

	case nodes.GroupingFunc:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Refs)
		walkn(f, n.Cols)

	case nodes.GroupingSet:
		walkn(f, n.Content)

	case nodes.ImportForeignSchemaStmt:
		walkn(f, n.TableList)
		walkn(f, n.Options)

	case nodes.IndexElem:
		walkn(f, n.Expr)
		walkn(f, n.Collation)
		walkn(f, n.Opclass)

	case nodes.IndexStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.IndexParams)
		walkn(f, n.Options)
		walkn(f, n.WhereClause)
		walkn(f, n.ExcludeOpNames)

	case nodes.InferClause:
		walkn(f, n.IndexElems)
		walkn(f, n.WhereClause)

	case nodes.InferenceElem:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)

	case nodes.InlineCodeBlock:
		// pass

	case nodes.InsertStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Cols)
		walkn(f, n.SelectStmt)
		if n.OnConflictClause != nil {
			walkn(f, *n.OnConflictClause)
		}
		walkn(f, n.ReturningList)
		if n.WithClause != nil {
			walkn(f, *n.WithClause)
		}

	case nodes.Integer:
		// pass

	case nodes.IntoClause:
		if n.Rel != nil {
			walkn(f, *n.Rel)
		}
		walkn(f, n.ColNames)
		walkn(f, n.Options)
		walkn(f, n.ViewQuery)

	case nodes.JoinExpr:
		walkn(f, n.Larg)
		walkn(f, n.Rarg)
		walkn(f, n.UsingClause)
		walkn(f, n.Quals)
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}

	case nodes.List:
		for _, item := range n.Items {
			walkn(f, item)
		}

	case nodes.ListenStmt:
		// pass

	case nodes.LoadStmt:
		// pass

	case nodes.LockStmt:
		walkn(f, n.Relations)

	case nodes.LockingClause:
		walkn(f, n.LockedRels)

	case nodes.MinMaxExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.MultiAssignRef:
		walkn(f, n.Source)

	case nodes.NamedArgExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.NextValueExpr:
		walkn(f, n.Xpr)

	case nodes.NotifyStmt:
		// pass

	case nodes.Null:
		// pass

	case nodes.NullTest:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.ObjectWithArgs:
		walkn(f, n.Objname)
		walkn(f, n.Objargs)

	case nodes.OnConflictClause:
		if n.Infer != nil {
			walkn(f, *n.Infer)
		}
		walkn(f, n.TargetList)
		walkn(f, n.WhereClause)

	case nodes.OnConflictExpr:
		walkn(f, n.ArbiterElems)
		walkn(f, n.ArbiterWhere)
		walkn(f, n.OnConflictSet)
		walkn(f, n.OnConflictWhere)
		walkn(f, n.ExclRelTlist)

	case nodes.OpExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.Param:
		walkn(f, n.Xpr)

	case nodes.ParamExecData:
		// pass

	case nodes.ParamExternData:
		// pass

	case nodes.ParamListInfoData:
		// pass

	case nodes.ParamRef:
		// pass

	case nodes.PartitionBoundSpec:
		walkn(f, n.Listdatums)
		walkn(f, n.Lowerdatums)
		walkn(f, n.Upperdatums)

	case nodes.PartitionCmd:
		if n.Name != nil {
			walkn(f, *n.Name)
		}
		if n.Bound != nil {
			walkn(f, *n.Bound)
		}

	case nodes.PartitionElem:
		walkn(f, n.Expr)
		walkn(f, n.Collation)
		walkn(f, n.Opclass)

	case nodes.PartitionRangeDatum:
		walkn(f, n.Value)

	case nodes.PartitionSpec:
		walkn(f, n.PartParams)

	case nodes.PrepareStmt:
		walkn(f, n.Argtypes)
		walkn(f, n.Query)

	case nodes.Query:
		walkn(f, n.UtilityStmt)
		walkn(f, n.CteList)
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
		walkn(f, n.LimitCount)
		walkn(f, n.RowMarks)
		walkn(f, n.SetOperations)
		walkn(f, n.ConstraintDeps)
		walkn(f, n.WithCheckOptions)

	case nodes.RangeFunction:
		walkn(f, n.Functions)
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}
		walkn(f, n.Coldeflist)

	case nodes.RangeSubselect:
		walkn(f, n.Subquery)
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}

	case nodes.RangeTableFunc:
		walkn(f, n.Docexpr)
		walkn(f, n.Rowexpr)
		walkn(f, n.Namespaces)
		walkn(f, n.Columns)
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}

	case nodes.RangeTableFuncCol:
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}
		walkn(f, n.Colexpr)
		walkn(f, n.Coldefexpr)

	case nodes.RangeTableSample:
		walkn(f, n.Relation)
		walkn(f, n.Method)
		walkn(f, n.Args)

	case nodes.RangeTblEntry:
		walkn(f, n.Tablesample)
		walkn(f, n.Subquery)
		walkn(f, n.Joinaliasvars)
		walkn(f, n.Functions)
		walkn(f, n.Tablefunc)
		walkn(f, n.ValuesLists)
		walkn(f, n.Coltypes)
		walkn(f, n.Colcollations)
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}
		walkn(f, n.Eref)
		walkn(f, n.SecurityQuals)

	case nodes.RangeTblFunction:
		walkn(f, n.Funcexpr)
		walkn(f, n.Funccolnames)
		walkn(f, n.Funccoltypes)
		walkn(f, n.Funccoltypmods)
		walkn(f, n.Funccolcollations)

	case nodes.RangeTblRef:
		// pass

	case nodes.RangeVar:
		if n.Alias != nil {
			walkn(f, *n.Alias)
		}

	case nodes.RawStmt:
		walkn(f, n.Stmt)

	case nodes.ReassignOwnedStmt:
		walkn(f, n.Roles)
		if n.Newrole != nil {
			walkn(f, *n.Newrole)
		}

	case nodes.RefreshMatViewStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}

	case nodes.ReindexStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}

	case nodes.RelabelType:
		walkn(f, n.Xpr)
		walkn(f, n.Arg)

	case nodes.RenameStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.Object)

	case nodes.ReplicaIdentityStmt:
		// pass

	case nodes.ResTarget:
		walkn(f, n.Indirection)
		walkn(f, n.Val)

	case nodes.RoleSpec:
		// pass

	case nodes.RowCompareExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Opnos)
		walkn(f, n.Opfamilies)
		walkn(f, n.Inputcollids)
		walkn(f, n.Largs)
		walkn(f, n.Rargs)

	case nodes.RowExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Colnames)

	case nodes.RowMarkClause:
		// pass

	case nodes.RuleStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.WhereClause)
		walkn(f, n.Actions)

	case nodes.SQLValueFunction:
		walkn(f, n.Xpr)

	case nodes.ScalarArrayOpExpr:
		walkn(f, n.Xpr)
		walkn(f, n.Args)

	case nodes.SecLabelStmt:
		walkn(f, n.Object)

	case nodes.SelectStmt:
		walkn(f, n.DistinctClause)
		if n.IntoClause != nil {
			walkn(f, *n.IntoClause)
		}
		walkn(f, n.TargetList)
		walkn(f, n.FromClause)
		walkn(f, n.WhereClause)
		walkn(f, n.GroupClause)
		walkn(f, n.HavingClause)
		walkn(f, n.WindowClause)
		for _, vs := range n.ValuesLists {
			for _, v := range vs {
				walkn(f, v)
			}
		}
		walkn(f, n.SortClause)
		walkn(f, n.LimitOffset)
		walkn(f, n.LimitCount)
		walkn(f, n.LockingClause)
		if n.WithClause != nil {
			walkn(f, *n.WithClause)
		}
		if n.Larg != nil {
			walkn(f, *n.Larg)
		}
		if n.Rarg != nil {
			walkn(f, *n.Rarg)
		}

	case nodes.SetOperationStmt:
		walkn(f, n.Larg)
		walkn(f, n.Rarg)
		walkn(f, n.ColTypes)
		walkn(f, n.ColTypmods)
		walkn(f, n.ColCollations)
		walkn(f, n.GroupClauses)

	case nodes.SetToDefault:
		walkn(f, n.Xpr)

	case nodes.SortBy:
		walkn(f, n.Node)
		walkn(f, n.UseOp)

	case nodes.SortGroupClause:
		// pass

	case nodes.String:
		// pass

	case nodes.SubLink:
		walkn(f, n.Xpr)
		walkn(f, n.Testexpr)
		walkn(f, n.OperName)
		walkn(f, n.Subselect)

	case nodes.SubPlan:
		walkn(f, n.Xpr)
		walkn(f, n.Testexpr)
		walkn(f, n.ParamIds)
		walkn(f, n.SetParam)
		walkn(f, n.ParParam)
		walkn(f, n.Args)

	case nodes.TableFunc:
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

	case nodes.TableLikeClause:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}

	case nodes.TableSampleClause:
		walkn(f, n.Args)
		walkn(f, n.Repeatable)

	case nodes.TargetEntry:
		walkn(f, n.Xpr)
		walkn(f, n.Expr)

	case nodes.TransactionStmt:
		walkn(f, n.Options)

	case nodes.TriggerTransition:
		// pass

	case nodes.TruncateStmt:
		walkn(f, n.Relations)

	case nodes.TypeCast:
		walkn(f, n.Arg)
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}

	case nodes.TypeName:
		walkn(f, n.Names)
		walkn(f, n.Typmods)
		walkn(f, n.ArrayBounds)

	case nodes.UnlistenStmt:
		// pass

	case nodes.UpdateStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.TargetList)
		walkn(f, n.WhereClause)
		walkn(f, n.FromClause)
		walkn(f, n.ReturningList)
		if n.WithClause != nil {
			walkn(f, *n.WithClause)
		}

	case nodes.VacuumStmt:
		if n.Relation != nil {
			walkn(f, *n.Relation)
		}
		walkn(f, n.VaCols)

	case nodes.Var:
		walkn(f, n.Xpr)

	case nodes.VariableSetStmt:
		walkn(f, n.Args)

	case nodes.VariableShowStmt:
		// pass

	case nodes.ViewStmt:
		if n.View != nil {
			walkn(f, *n.View)
		}
		walkn(f, n.Aliases)
		walkn(f, n.Query)
		walkn(f, n.Options)

	case nodes.WindowClause:
		walkn(f, n.PartitionClause)
		walkn(f, n.OrderClause)
		walkn(f, n.StartOffset)
		walkn(f, n.EndOffset)

	case nodes.WindowDef:
		walkn(f, n.PartitionClause)
		walkn(f, n.OrderClause)
		walkn(f, n.StartOffset)
		walkn(f, n.EndOffset)

	case nodes.WindowFunc:
		walkn(f, n.Xpr)
		walkn(f, n.Args)
		walkn(f, n.Aggfilter)

	case nodes.WithCheckOption:
		walkn(f, n.Qual)

	case nodes.WithClause:
		walkn(f, n.Ctes)

	case nodes.XmlExpr:
		walkn(f, n.Xpr)
		walkn(f, n.NamedArgs)
		walkn(f, n.ArgNames)
		walkn(f, n.Args)

	case nodes.XmlSerialize:
		walkn(f, n.Expr)
		if n.TypeName != nil {
			walkn(f, *n.TypeName)
		}

	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))

	}

	f.Visit(nil)
}
