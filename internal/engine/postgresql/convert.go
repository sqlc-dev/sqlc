package postgresql

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
)

func convertList(l nodes.List) *ast.List {
	out := &ast.List{}
	for _, item := range l.Items {
		out.Items = append(out.Items, convertNode(item))
	}
	return out
}

func convertValuesList(l [][]nodes.Node) *ast.List {
	out := &ast.List{}
	for _, outer := range l {
		o := &ast.List{}
		for _, inner := range outer {
			o.Items = append(o.Items, convertNode(inner))
		}
		out.Items = append(out.Items, o)
	}
	return out
}

func convert(node nodes.Node) (ast.Node, error) {
	return convertNode(node), nil
}

func convertA_ArrayExpr(n *nodes.A_ArrayExpr) *pg.A_ArrayExpr {
	if n == nil {
		return nil
	}
	return &pg.A_ArrayExpr{
		Elements: convertList(n.Elements),
		Location: n.Location,
	}
}

func convertA_Const(n *nodes.A_Const) *pg.A_Const {
	if n == nil {
		return nil
	}
	return &pg.A_Const{
		Val:      convertNode(n.Val),
		Location: n.Location,
	}
}

func convertA_Expr(n *nodes.A_Expr) *pg.A_Expr {
	if n == nil {
		return nil
	}
	return &pg.A_Expr{
		Kind:     pg.A_Expr_Kind(n.Kind),
		Name:     convertList(n.Name),
		Lexpr:    convertNode(n.Lexpr),
		Rexpr:    convertNode(n.Rexpr),
		Location: n.Location,
	}
}

func convertA_Indices(n *nodes.A_Indices) *pg.A_Indices {
	if n == nil {
		return nil
	}
	return &pg.A_Indices{
		IsSlice: n.IsSlice,
		Lidx:    convertNode(n.Lidx),
		Uidx:    convertNode(n.Uidx),
	}
}

func convertA_Indirection(n *nodes.A_Indirection) *pg.A_Indirection {
	if n == nil {
		return nil
	}
	return &pg.A_Indirection{
		Arg:         convertNode(n.Arg),
		Indirection: convertList(n.Indirection),
	}
}

func convertA_Star(n *nodes.A_Star) *pg.A_Star {
	if n == nil {
		return nil
	}
	return &pg.A_Star{}
}

func convertAccessPriv(n *nodes.AccessPriv) *pg.AccessPriv {
	if n == nil {
		return nil
	}
	return &pg.AccessPriv{
		PrivName: n.PrivName,
		Cols:     convertList(n.Cols),
	}
}

func convertAggref(n *nodes.Aggref) *pg.Aggref {
	if n == nil {
		return nil
	}
	return &pg.Aggref{
		Xpr:           convertNode(n.Xpr),
		Aggfnoid:      pg.Oid(n.Aggfnoid),
		Aggtype:       pg.Oid(n.Aggtype),
		Aggcollid:     pg.Oid(n.Aggcollid),
		Inputcollid:   pg.Oid(n.Inputcollid),
		Aggtranstype:  pg.Oid(n.Aggtranstype),
		Aggargtypes:   convertList(n.Aggargtypes),
		Aggdirectargs: convertList(n.Aggdirectargs),
		Args:          convertList(n.Args),
		Aggorder:      convertList(n.Aggorder),
		Aggdistinct:   convertList(n.Aggdistinct),
		Aggfilter:     convertNode(n.Aggfilter),
		Aggstar:       n.Aggstar,
		Aggvariadic:   n.Aggvariadic,
		Aggkind:       n.Aggkind,
		Agglevelsup:   pg.Index(n.Agglevelsup),
		Aggsplit:      pg.AggSplit(n.Aggsplit),
		Location:      n.Location,
	}
}

func convertAlias(n *nodes.Alias) *pg.Alias {
	if n == nil {
		return nil
	}
	return &pg.Alias{
		Aliasname: n.Aliasname,
		Colnames:  convertList(n.Colnames),
	}
}

func convertAlterCollationStmt(n *nodes.AlterCollationStmt) *pg.AlterCollationStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterCollationStmt{
		Collname: convertList(n.Collname),
	}
}

func convertAlterDatabaseSetStmt(n *nodes.AlterDatabaseSetStmt) *pg.AlterDatabaseSetStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterDatabaseSetStmt{
		Dbname:  n.Dbname,
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterDatabaseStmt(n *nodes.AlterDatabaseStmt) *pg.AlterDatabaseStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterDatabaseStmt{
		Dbname:  n.Dbname,
		Options: convertList(n.Options),
	}
}

func convertAlterDefaultPrivilegesStmt(n *nodes.AlterDefaultPrivilegesStmt) *pg.AlterDefaultPrivilegesStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterDefaultPrivilegesStmt{
		Options: convertList(n.Options),
		Action:  convertGrantStmt(n.Action),
	}
}

func convertAlterDomainStmt(n *nodes.AlterDomainStmt) *pg.AlterDomainStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterDomainStmt{
		Subtype:   n.Subtype,
		TypeName:  convertList(n.TypeName),
		Name:      n.Name,
		Def:       convertNode(n.Def),
		Behavior:  pg.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterEnumStmt(n *nodes.AlterEnumStmt) *pg.AlterEnumStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterEnumStmt{
		TypeName:           convertList(n.TypeName),
		OldVal:             n.OldVal,
		NewVal:             n.NewVal,
		NewValNeighbor:     n.NewValNeighbor,
		NewValIsAfter:      n.NewValIsAfter,
		SkipIfNewValExists: n.SkipIfNewValExists,
	}
}

func convertAlterEventTrigStmt(n *nodes.AlterEventTrigStmt) *pg.AlterEventTrigStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterEventTrigStmt{
		Trigname:  n.Trigname,
		Tgenabled: n.Tgenabled,
	}
}

func convertAlterExtensionContentsStmt(n *nodes.AlterExtensionContentsStmt) *pg.AlterExtensionContentsStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterExtensionContentsStmt{
		Extname: n.Extname,
		Action:  n.Action,
		Objtype: pg.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
	}
}

func convertAlterExtensionStmt(n *nodes.AlterExtensionStmt) *pg.AlterExtensionStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterExtensionStmt{
		Extname: n.Extname,
		Options: convertList(n.Options),
	}
}

func convertAlterFdwStmt(n *nodes.AlterFdwStmt) *pg.AlterFdwStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterFdwStmt{
		Fdwname:     n.Fdwname,
		FuncOptions: convertList(n.FuncOptions),
		Options:     convertList(n.Options),
	}
}

func convertAlterForeignServerStmt(n *nodes.AlterForeignServerStmt) *pg.AlterForeignServerStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterForeignServerStmt{
		Servername: n.Servername,
		Version:    n.Version,
		Options:    convertList(n.Options),
		HasVersion: n.HasVersion,
	}
}

func convertAlterFunctionStmt(n *nodes.AlterFunctionStmt) *pg.AlterFunctionStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterFunctionStmt{
		Func:    convertObjectWithArgs(n.Func),
		Actions: convertList(n.Actions),
	}
}

func convertAlterObjectDependsStmt(n *nodes.AlterObjectDependsStmt) *pg.AlterObjectDependsStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterObjectDependsStmt{
		ObjectType: pg.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Extname:    convertNode(n.Extname),
	}
}

func convertAlterObjectSchemaStmt(n *nodes.AlterObjectSchemaStmt) *pg.AlterObjectSchemaStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterObjectSchemaStmt{
		ObjectType: pg.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Newschema:  n.Newschema,
		MissingOk:  n.MissingOk,
	}
}

func convertAlterOpFamilyStmt(n *nodes.AlterOpFamilyStmt) *pg.AlterOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterOpFamilyStmt{
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
		IsDrop:       n.IsDrop,
		Items:        convertList(n.Items),
	}
}

func convertAlterOperatorStmt(n *nodes.AlterOperatorStmt) *pg.AlterOperatorStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterOperatorStmt{
		Opername: convertObjectWithArgs(n.Opername),
		Options:  convertList(n.Options),
	}
}

func convertAlterOwnerStmt(n *nodes.AlterOwnerStmt) *pg.AlterOwnerStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterOwnerStmt{
		ObjectType: pg.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Newowner:   convertRoleSpec(n.Newowner),
	}
}

func convertAlterPolicyStmt(n *nodes.AlterPolicyStmt) *pg.AlterPolicyStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterPolicyStmt{
		PolicyName: n.PolicyName,
		Table:      convertRangeVar(n.Table),
		Roles:      convertList(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertAlterPublicationStmt(n *nodes.AlterPublicationStmt) *pg.AlterPublicationStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterPublicationStmt{
		Pubname:      n.Pubname,
		Options:      convertList(n.Options),
		Tables:       convertList(n.Tables),
		ForAllTables: n.ForAllTables,
		TableAction:  pg.DefElemAction(n.TableAction),
	}
}

func convertAlterRoleSetStmt(n *nodes.AlterRoleSetStmt) *pg.AlterRoleSetStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterRoleSetStmt{
		Role:     convertRoleSpec(n.Role),
		Database: n.Database,
		Setstmt:  convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterRoleStmt(n *nodes.AlterRoleStmt) *pg.AlterRoleStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterRoleStmt{
		Role:    convertRoleSpec(n.Role),
		Options: convertList(n.Options),
		Action:  n.Action,
	}
}

func convertAlterSeqStmt(n *nodes.AlterSeqStmt) *pg.AlterSeqStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertList(n.Options),
		ForIdentity: n.ForIdentity,
		MissingOk:   n.MissingOk,
	}
}

func convertAlterSubscriptionStmt(n *nodes.AlterSubscriptionStmt) *pg.AlterSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterSubscriptionStmt{
		Kind:        pg.AlterSubscriptionType(n.Kind),
		Subname:     n.Subname,
		Conninfo:    n.Conninfo,
		Publication: convertList(n.Publication),
		Options:     convertList(n.Options),
	}
}

func convertAlterSystemStmt(n *nodes.AlterSystemStmt) *pg.AlterSystemStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterSystemStmt{
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterTSConfigurationStmt(n *nodes.AlterTSConfigurationStmt) *pg.AlterTSConfigurationStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterTSConfigurationStmt{
		Kind:      pg.AlterTSConfigType(n.Kind),
		Cfgname:   convertList(n.Cfgname),
		Tokentype: convertList(n.Tokentype),
		Dicts:     convertList(n.Dicts),
		Override:  n.Override,
		Replace:   n.Replace,
		MissingOk: n.MissingOk,
	}
}

func convertAlterTSDictionaryStmt(n *nodes.AlterTSDictionaryStmt) *pg.AlterTSDictionaryStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterTSDictionaryStmt{
		Dictname: convertList(n.Dictname),
		Options:  convertList(n.Options),
	}
}

func convertAlterTableCmd(n *nodes.AlterTableCmd) *pg.AlterTableCmd {
	if n == nil {
		return nil
	}
	return &pg.AlterTableCmd{
		Subtype:   pg.AlterTableType(n.Subtype),
		Name:      n.Name,
		Newowner:  convertRoleSpec(n.Newowner),
		Def:       convertNode(n.Def),
		Behavior:  pg.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterTableMoveAllStmt(n *nodes.AlterTableMoveAllStmt) *pg.AlterTableMoveAllStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterTableMoveAllStmt{
		OrigTablespacename: n.OrigTablespacename,
		Objtype:            pg.ObjectType(n.Objtype),
		Roles:              convertList(n.Roles),
		NewTablespacename:  n.NewTablespacename,
		Nowait:             n.Nowait,
	}
}

func convertAlterTableSpaceOptionsStmt(n *nodes.AlterTableSpaceOptionsStmt) *pg.AlterTableSpaceOptionsStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterTableSpaceOptionsStmt{
		Tablespacename: n.Tablespacename,
		Options:        convertList(n.Options),
		IsReset:        n.IsReset,
	}
}

func convertAlterTableStmt(n *nodes.AlterTableStmt) *pg.AlterTableStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterTableStmt{
		Relation:  convertRangeVar(n.Relation),
		Cmds:      convertList(n.Cmds),
		Relkind:   pg.ObjectType(n.Relkind),
		MissingOk: n.MissingOk,
	}
}

func convertAlterUserMappingStmt(n *nodes.AlterUserMappingStmt) *pg.AlterUserMappingStmt {
	if n == nil {
		return nil
	}
	return &pg.AlterUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: n.Servername,
		Options:    convertList(n.Options),
	}
}

func convertAlternativeSubPlan(n *nodes.AlternativeSubPlan) *pg.AlternativeSubPlan {
	if n == nil {
		return nil
	}
	return &pg.AlternativeSubPlan{
		Xpr:      convertNode(n.Xpr),
		Subplans: convertList(n.Subplans),
	}
}

func convertArrayCoerceExpr(n *nodes.ArrayCoerceExpr) *pg.ArrayCoerceExpr {
	if n == nil {
		return nil
	}
	return &pg.ArrayCoerceExpr{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Elemfuncid:   pg.Oid(n.Elemfuncid),
		Resulttype:   pg.Oid(n.Resulttype),
		Resulttypmod: n.Resulttypmod,
		Resultcollid: pg.Oid(n.Resultcollid),
		IsExplicit:   n.IsExplicit,
		Coerceformat: pg.CoercionForm(n.Coerceformat),
		Location:     n.Location,
	}
}

func convertArrayExpr(n *nodes.ArrayExpr) *pg.ArrayExpr {
	if n == nil {
		return nil
	}
	return &pg.ArrayExpr{
		Xpr:           convertNode(n.Xpr),
		ArrayTypeid:   pg.Oid(n.ArrayTypeid),
		ArrayCollid:   pg.Oid(n.ArrayCollid),
		ElementTypeid: pg.Oid(n.ElementTypeid),
		Elements:      convertList(n.Elements),
		Multidims:     n.Multidims,
		Location:      n.Location,
	}
}

func convertArrayRef(n *nodes.ArrayRef) *pg.ArrayRef {
	if n == nil {
		return nil
	}
	return &pg.ArrayRef{
		Xpr:             convertNode(n.Xpr),
		Refarraytype:    pg.Oid(n.Refarraytype),
		Refelemtype:     pg.Oid(n.Refelemtype),
		Reftypmod:       n.Reftypmod,
		Refcollid:       pg.Oid(n.Refcollid),
		Refupperindexpr: convertList(n.Refupperindexpr),
		Reflowerindexpr: convertList(n.Reflowerindexpr),
		Refexpr:         convertNode(n.Refexpr),
		Refassgnexpr:    convertNode(n.Refassgnexpr),
	}
}

func convertBitString(n *nodes.BitString) *pg.BitString {
	if n == nil {
		return nil
	}
	return &pg.BitString{
		Str: n.Str,
	}
}

func convertBlockIdData(n *nodes.BlockIdData) *pg.BlockIdData {
	if n == nil {
		return nil
	}
	return &pg.BlockIdData{
		BiHi: n.BiHi,
		BiLo: n.BiLo,
	}
}

func convertBoolExpr(n *nodes.BoolExpr) *pg.BoolExpr {
	if n == nil {
		return nil
	}
	return &pg.BoolExpr{
		Xpr:      convertNode(n.Xpr),
		Boolop:   pg.BoolExprType(n.Boolop),
		Args:     convertList(n.Args),
		Location: n.Location,
	}
}

func convertBooleanTest(n *nodes.BooleanTest) *pg.BooleanTest {
	if n == nil {
		return nil
	}
	return &pg.BooleanTest{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Booltesttype: pg.BoolTestType(n.Booltesttype),
		Location:     n.Location,
	}
}

func convertCaseExpr(n *nodes.CaseExpr) *pg.CaseExpr {
	if n == nil {
		return nil
	}
	return &pg.CaseExpr{
		Xpr:        convertNode(n.Xpr),
		Casetype:   pg.Oid(n.Casetype),
		Casecollid: pg.Oid(n.Casecollid),
		Arg:        convertNode(n.Arg),
		Args:       convertList(n.Args),
		Defresult:  convertNode(n.Defresult),
		Location:   n.Location,
	}
}

func convertCaseTestExpr(n *nodes.CaseTestExpr) *pg.CaseTestExpr {
	if n == nil {
		return nil
	}
	return &pg.CaseTestExpr{
		Xpr:       convertNode(n.Xpr),
		TypeId:    pg.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: pg.Oid(n.Collation),
	}
}

func convertCaseWhen(n *nodes.CaseWhen) *pg.CaseWhen {
	if n == nil {
		return nil
	}
	return &pg.CaseWhen{
		Xpr:      convertNode(n.Xpr),
		Expr:     convertNode(n.Expr),
		Result:   convertNode(n.Result),
		Location: n.Location,
	}
}

func convertCheckPointStmt(n *nodes.CheckPointStmt) *pg.CheckPointStmt {
	if n == nil {
		return nil
	}
	return &pg.CheckPointStmt{}
}

func convertClosePortalStmt(n *nodes.ClosePortalStmt) *pg.ClosePortalStmt {
	if n == nil {
		return nil
	}
	return &pg.ClosePortalStmt{
		Portalname: n.Portalname,
	}
}

func convertClusterStmt(n *nodes.ClusterStmt) *pg.ClusterStmt {
	if n == nil {
		return nil
	}
	return &pg.ClusterStmt{
		Relation:  convertRangeVar(n.Relation),
		Indexname: n.Indexname,
		Verbose:   n.Verbose,
	}
}

func convertCoalesceExpr(n *nodes.CoalesceExpr) *pg.CoalesceExpr {
	if n == nil {
		return nil
	}
	return &pg.CoalesceExpr{
		Xpr:            convertNode(n.Xpr),
		Coalescetype:   pg.Oid(n.Coalescetype),
		Coalescecollid: pg.Oid(n.Coalescecollid),
		Args:           convertList(n.Args),
		Location:       n.Location,
	}
}

func convertCoerceToDomain(n *nodes.CoerceToDomain) *pg.CoerceToDomain {
	if n == nil {
		return nil
	}
	return &pg.CoerceToDomain{
		Xpr:            convertNode(n.Xpr),
		Arg:            convertNode(n.Arg),
		Resulttype:     pg.Oid(n.Resulttype),
		Resulttypmod:   n.Resulttypmod,
		Resultcollid:   pg.Oid(n.Resultcollid),
		Coercionformat: pg.CoercionForm(n.Coercionformat),
		Location:       n.Location,
	}
}

func convertCoerceToDomainValue(n *nodes.CoerceToDomainValue) *pg.CoerceToDomainValue {
	if n == nil {
		return nil
	}
	return &pg.CoerceToDomainValue{
		Xpr:       convertNode(n.Xpr),
		TypeId:    pg.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: pg.Oid(n.Collation),
		Location:  n.Location,
	}
}

func convertCoerceViaIO(n *nodes.CoerceViaIO) *pg.CoerceViaIO {
	if n == nil {
		return nil
	}
	return &pg.CoerceViaIO{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Resulttype:   pg.Oid(n.Resulttype),
		Resultcollid: pg.Oid(n.Resultcollid),
		Coerceformat: pg.CoercionForm(n.Coerceformat),
		Location:     n.Location,
	}
}

func convertCollateClause(n *nodes.CollateClause) *pg.CollateClause {
	if n == nil {
		return nil
	}
	return &pg.CollateClause{
		Arg:      convertNode(n.Arg),
		Collname: convertList(n.Collname),
		Location: n.Location,
	}
}

func convertCollateExpr(n *nodes.CollateExpr) *pg.CollateExpr {
	if n == nil {
		return nil
	}
	return &pg.CollateExpr{
		Xpr:      convertNode(n.Xpr),
		Arg:      convertNode(n.Arg),
		CollOid:  pg.Oid(n.CollOid),
		Location: n.Location,
	}
}

func convertColumnDef(n *nodes.ColumnDef) *pg.ColumnDef {
	if n == nil {
		return nil
	}
	return &pg.ColumnDef{
		Colname:       n.Colname,
		TypeName:      convertTypeName(n.TypeName),
		Inhcount:      n.Inhcount,
		IsLocal:       n.IsLocal,
		IsNotNull:     n.IsNotNull,
		IsFromType:    n.IsFromType,
		IsFromParent:  n.IsFromParent,
		Storage:       n.Storage,
		RawDefault:    convertNode(n.RawDefault),
		CookedDefault: convertNode(n.CookedDefault),
		Identity:      n.Identity,
		CollClause:    convertCollateClause(n.CollClause),
		CollOid:       pg.Oid(n.CollOid),
		Constraints:   convertList(n.Constraints),
		Fdwoptions:    convertList(n.Fdwoptions),
		Location:      n.Location,
	}
}

func convertColumnRef(n *nodes.ColumnRef) *pg.ColumnRef {
	if n == nil {
		return nil
	}
	return &pg.ColumnRef{
		Fields:   convertList(n.Fields),
		Location: n.Location,
	}
}

func convertCommentStmt(n *nodes.CommentStmt) *pg.CommentStmt {
	if n == nil {
		return nil
	}
	return &pg.CommentStmt{
		Objtype: pg.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
		Comment: n.Comment,
	}
}

func convertCommonTableExpr(n *nodes.CommonTableExpr) *pg.CommonTableExpr {
	if n == nil {
		return nil
	}
	return &pg.CommonTableExpr{
		Ctename:          n.Ctename,
		Aliascolnames:    convertList(n.Aliascolnames),
		Ctequery:         convertNode(n.Ctequery),
		Location:         n.Location,
		Cterecursive:     n.Cterecursive,
		Cterefcount:      n.Cterefcount,
		Ctecolnames:      convertList(n.Ctecolnames),
		Ctecoltypes:      convertList(n.Ctecoltypes),
		Ctecoltypmods:    convertList(n.Ctecoltypmods),
		Ctecolcollations: convertList(n.Ctecolcollations),
	}
}

func convertCompositeTypeStmt(n *nodes.CompositeTypeStmt) *pg.CompositeTypeStmt {
	if n == nil {
		return nil
	}
	return &pg.CompositeTypeStmt{
		Typevar:    convertRangeVar(n.Typevar),
		Coldeflist: convertList(n.Coldeflist),
	}
}

func convertConst(n *nodes.Const) *pg.Const {
	if n == nil {
		return nil
	}
	return &pg.Const{
		Xpr:         convertNode(n.Xpr),
		Consttype:   pg.Oid(n.Consttype),
		Consttypmod: n.Consttypmod,
		Constcollid: pg.Oid(n.Constcollid),
		Constlen:    n.Constlen,
		Constvalue:  pg.Datum(n.Constvalue),
		Constisnull: n.Constisnull,
		Constbyval:  n.Constbyval,
		Location:    n.Location,
	}
}

func convertConstraint(n *nodes.Constraint) *pg.Constraint {
	if n == nil {
		return nil
	}
	return &pg.Constraint{
		Contype:        pg.ConstrType(n.Contype),
		Conname:        n.Conname,
		Deferrable:     n.Deferrable,
		Initdeferred:   n.Initdeferred,
		Location:       n.Location,
		IsNoInherit:    n.IsNoInherit,
		RawExpr:        convertNode(n.RawExpr),
		CookedExpr:     n.CookedExpr,
		GeneratedWhen:  n.GeneratedWhen,
		Keys:           convertList(n.Keys),
		Exclusions:     convertList(n.Exclusions),
		Options:        convertList(n.Options),
		Indexname:      n.Indexname,
		Indexspace:     n.Indexspace,
		AccessMethod:   n.AccessMethod,
		WhereClause:    convertNode(n.WhereClause),
		Pktable:        convertRangeVar(n.Pktable),
		FkAttrs:        convertList(n.FkAttrs),
		PkAttrs:        convertList(n.PkAttrs),
		FkMatchtype:    n.FkMatchtype,
		FkUpdAction:    n.FkUpdAction,
		FkDelAction:    n.FkDelAction,
		OldConpfeqop:   convertList(n.OldConpfeqop),
		OldPktableOid:  pg.Oid(n.OldPktableOid),
		SkipValidation: n.SkipValidation,
		InitiallyValid: n.InitiallyValid,
	}
}

func convertConstraintsSetStmt(n *nodes.ConstraintsSetStmt) *pg.ConstraintsSetStmt {
	if n == nil {
		return nil
	}
	return &pg.ConstraintsSetStmt{
		Constraints: convertList(n.Constraints),
		Deferred:    n.Deferred,
	}
}

func convertConvertRowtypeExpr(n *nodes.ConvertRowtypeExpr) *pg.ConvertRowtypeExpr {
	if n == nil {
		return nil
	}
	return &pg.ConvertRowtypeExpr{
		Xpr:           convertNode(n.Xpr),
		Arg:           convertNode(n.Arg),
		Resulttype:    pg.Oid(n.Resulttype),
		Convertformat: pg.CoercionForm(n.Convertformat),
		Location:      n.Location,
	}
}

func convertCopyStmt(n *nodes.CopyStmt) *pg.CopyStmt {
	if n == nil {
		return nil
	}
	return &pg.CopyStmt{
		Relation:  convertRangeVar(n.Relation),
		Query:     convertNode(n.Query),
		Attlist:   convertList(n.Attlist),
		IsFrom:    n.IsFrom,
		IsProgram: n.IsProgram,
		Filename:  n.Filename,
		Options:   convertList(n.Options),
	}
}

func convertCreateAmStmt(n *nodes.CreateAmStmt) *pg.CreateAmStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateAmStmt{
		Amname:      n.Amname,
		HandlerName: convertList(n.HandlerName),
		Amtype:      n.Amtype,
	}
}

func convertCreateCastStmt(n *nodes.CreateCastStmt) *pg.CreateCastStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateCastStmt{
		Sourcetype: convertTypeName(n.Sourcetype),
		Targettype: convertTypeName(n.Targettype),
		Func:       convertObjectWithArgs(n.Func),
		Context:    pg.CoercionContext(n.Context),
		Inout:      n.Inout,
	}
}

func convertCreateConversionStmt(n *nodes.CreateConversionStmt) *pg.CreateConversionStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateConversionStmt{
		ConversionName:  convertList(n.ConversionName),
		ForEncodingName: n.ForEncodingName,
		ToEncodingName:  n.ToEncodingName,
		FuncName:        convertList(n.FuncName),
		Def:             n.Def,
	}
}

func convertCreateDomainStmt(n *nodes.CreateDomainStmt) *pg.CreateDomainStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateDomainStmt{
		Domainname:  convertList(n.Domainname),
		TypeName:    convertTypeName(n.TypeName),
		CollClause:  convertCollateClause(n.CollClause),
		Constraints: convertList(n.Constraints),
	}
}

func convertCreateEnumStmt(n *nodes.CreateEnumStmt) *pg.CreateEnumStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateEnumStmt{
		TypeName: convertList(n.TypeName),
		Vals:     convertList(n.Vals),
	}
}

func convertCreateEventTrigStmt(n *nodes.CreateEventTrigStmt) *pg.CreateEventTrigStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateEventTrigStmt{
		Trigname:   n.Trigname,
		Eventname:  n.Eventname,
		Whenclause: convertList(n.Whenclause),
		Funcname:   convertList(n.Funcname),
	}
}

func convertCreateExtensionStmt(n *nodes.CreateExtensionStmt) *pg.CreateExtensionStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateExtensionStmt{
		Extname:     n.Extname,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreateFdwStmt(n *nodes.CreateFdwStmt) *pg.CreateFdwStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateFdwStmt{
		Fdwname:     n.Fdwname,
		FuncOptions: convertList(n.FuncOptions),
		Options:     convertList(n.Options),
	}
}

func convertCreateForeignServerStmt(n *nodes.CreateForeignServerStmt) *pg.CreateForeignServerStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateForeignServerStmt{
		Servername:  n.Servername,
		Servertype:  n.Servertype,
		Version:     n.Version,
		Fdwname:     n.Fdwname,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreateForeignTableStmt(n *nodes.CreateForeignTableStmt) *pg.CreateForeignTableStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateForeignTableStmt{
		Base:       convertCreateStmt(&n.Base),
		Servername: n.Servername,
		Options:    convertList(n.Options),
	}
}

func convertCreateFunctionStmt(n *nodes.CreateFunctionStmt) *pg.CreateFunctionStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateFunctionStmt{
		Replace:    n.Replace,
		Funcname:   convertList(n.Funcname),
		Parameters: convertList(n.Parameters),
		ReturnType: convertTypeName(n.ReturnType),
		Options:    convertList(n.Options),
		WithClause: convertList(n.WithClause),
	}
}

func convertCreateOpClassItem(n *nodes.CreateOpClassItem) *pg.CreateOpClassItem {
	if n == nil {
		return nil
	}
	return &pg.CreateOpClassItem{
		Itemtype:    n.Itemtype,
		Name:        convertObjectWithArgs(n.Name),
		Number:      n.Number,
		OrderFamily: convertList(n.OrderFamily),
		ClassArgs:   convertList(n.ClassArgs),
		Storedtype:  convertTypeName(n.Storedtype),
	}
}

func convertCreateOpClassStmt(n *nodes.CreateOpClassStmt) *pg.CreateOpClassStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateOpClassStmt{
		Opclassname:  convertList(n.Opclassname),
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
		Datatype:     convertTypeName(n.Datatype),
		Items:        convertList(n.Items),
		IsDefault:    n.IsDefault,
	}
}

func convertCreateOpFamilyStmt(n *nodes.CreateOpFamilyStmt) *pg.CreateOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateOpFamilyStmt{
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
	}
}

func convertCreatePLangStmt(n *nodes.CreatePLangStmt) *pg.CreatePLangStmt {
	if n == nil {
		return nil
	}
	return &pg.CreatePLangStmt{
		Replace:     n.Replace,
		Plname:      n.Plname,
		Plhandler:   convertList(n.Plhandler),
		Plinline:    convertList(n.Plinline),
		Plvalidator: convertList(n.Plvalidator),
		Pltrusted:   n.Pltrusted,
	}
}

func convertCreatePolicyStmt(n *nodes.CreatePolicyStmt) *pg.CreatePolicyStmt {
	if n == nil {
		return nil
	}
	return &pg.CreatePolicyStmt{
		PolicyName: n.PolicyName,
		Table:      convertRangeVar(n.Table),
		CmdName:    n.CmdName,
		Permissive: n.Permissive,
		Roles:      convertList(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertCreatePublicationStmt(n *nodes.CreatePublicationStmt) *pg.CreatePublicationStmt {
	if n == nil {
		return nil
	}
	return &pg.CreatePublicationStmt{
		Pubname:      n.Pubname,
		Options:      convertList(n.Options),
		Tables:       convertList(n.Tables),
		ForAllTables: n.ForAllTables,
	}
}

func convertCreateRangeStmt(n *nodes.CreateRangeStmt) *pg.CreateRangeStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateRangeStmt{
		TypeName: convertList(n.TypeName),
		Params:   convertList(n.Params),
	}
}

func convertCreateRoleStmt(n *nodes.CreateRoleStmt) *pg.CreateRoleStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateRoleStmt{
		StmtType: pg.RoleStmtType(n.StmtType),
		Role:     n.Role,
		Options:  convertList(n.Options),
	}
}

func convertCreateSchemaStmt(n *nodes.CreateSchemaStmt) *pg.CreateSchemaStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateSchemaStmt{
		Schemaname:  n.Schemaname,
		Authrole:    convertRoleSpec(n.Authrole),
		SchemaElts:  convertList(n.SchemaElts),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateSeqStmt(n *nodes.CreateSeqStmt) *pg.CreateSeqStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertList(n.Options),
		OwnerId:     pg.Oid(n.OwnerId),
		ForIdentity: n.ForIdentity,
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStatsStmt(n *nodes.CreateStatsStmt) *pg.CreateStatsStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateStatsStmt{
		Defnames:    convertList(n.Defnames),
		StatTypes:   convertList(n.StatTypes),
		Exprs:       convertList(n.Exprs),
		Relations:   convertList(n.Relations),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStmt(n *nodes.CreateStmt) *pg.CreateStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateStmt{
		Relation:       convertRangeVar(n.Relation),
		TableElts:      convertList(n.TableElts),
		InhRelations:   convertList(n.InhRelations),
		Partbound:      convertPartitionBoundSpec(n.Partbound),
		Partspec:       convertPartitionSpec(n.Partspec),
		OfTypename:     convertTypeName(n.OfTypename),
		Constraints:    convertList(n.Constraints),
		Options:        convertList(n.Options),
		Oncommit:       pg.OnCommitAction(n.Oncommit),
		Tablespacename: n.Tablespacename,
		IfNotExists:    n.IfNotExists,
	}
}

func convertCreateSubscriptionStmt(n *nodes.CreateSubscriptionStmt) *pg.CreateSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateSubscriptionStmt{
		Subname:     n.Subname,
		Conninfo:    n.Conninfo,
		Publication: convertList(n.Publication),
		Options:     convertList(n.Options),
	}
}

func convertCreateTableAsStmt(n *nodes.CreateTableAsStmt) *pg.CreateTableAsStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateTableAsStmt{
		Query:        convertNode(n.Query),
		Into:         convertIntoClause(n.Into),
		Relkind:      pg.ObjectType(n.Relkind),
		IsSelectInto: n.IsSelectInto,
		IfNotExists:  n.IfNotExists,
	}
}

func convertCreateTableSpaceStmt(n *nodes.CreateTableSpaceStmt) *pg.CreateTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		Owner:          convertRoleSpec(n.Owner),
		Location:       n.Location,
		Options:        convertList(n.Options),
	}
}

func convertCreateTransformStmt(n *nodes.CreateTransformStmt) *pg.CreateTransformStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateTransformStmt{
		Replace:  n.Replace,
		TypeName: convertTypeName(n.TypeName),
		Lang:     n.Lang,
		Fromsql:  convertObjectWithArgs(n.Fromsql),
		Tosql:    convertObjectWithArgs(n.Tosql),
	}
}

func convertCreateTrigStmt(n *nodes.CreateTrigStmt) *pg.CreateTrigStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateTrigStmt{
		Trigname:       n.Trigname,
		Relation:       convertRangeVar(n.Relation),
		Funcname:       convertList(n.Funcname),
		Args:           convertList(n.Args),
		Row:            n.Row,
		Timing:         n.Timing,
		Events:         n.Events,
		Columns:        convertList(n.Columns),
		WhenClause:     convertNode(n.WhenClause),
		Isconstraint:   n.Isconstraint,
		TransitionRels: convertList(n.TransitionRels),
		Deferrable:     n.Deferrable,
		Initdeferred:   n.Initdeferred,
		Constrrel:      convertRangeVar(n.Constrrel),
	}
}

func convertCreateUserMappingStmt(n *nodes.CreateUserMappingStmt) *pg.CreateUserMappingStmt {
	if n == nil {
		return nil
	}
	return &pg.CreateUserMappingStmt{
		User:        convertRoleSpec(n.User),
		Servername:  n.Servername,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreatedbStmt(n *nodes.CreatedbStmt) *pg.CreatedbStmt {
	if n == nil {
		return nil
	}
	return &pg.CreatedbStmt{
		Dbname:  n.Dbname,
		Options: convertList(n.Options),
	}
}

func convertCurrentOfExpr(n *nodes.CurrentOfExpr) *pg.CurrentOfExpr {
	if n == nil {
		return nil
	}
	return &pg.CurrentOfExpr{
		Xpr:         convertNode(n.Xpr),
		Cvarno:      pg.Index(n.Cvarno),
		CursorName:  n.CursorName,
		CursorParam: n.CursorParam,
	}
}

func convertDeallocateStmt(n *nodes.DeallocateStmt) *pg.DeallocateStmt {
	if n == nil {
		return nil
	}
	return &pg.DeallocateStmt{
		Name: n.Name,
	}
}

func convertDeclareCursorStmt(n *nodes.DeclareCursorStmt) *pg.DeclareCursorStmt {
	if n == nil {
		return nil
	}
	return &pg.DeclareCursorStmt{
		Portalname: n.Portalname,
		Options:    n.Options,
		Query:      convertNode(n.Query),
	}
}

func convertDefElem(n *nodes.DefElem) *pg.DefElem {
	if n == nil {
		return nil
	}
	return &pg.DefElem{
		Defnamespace: n.Defnamespace,
		Defname:      n.Defname,
		Arg:          convertNode(n.Arg),
		Defaction:    pg.DefElemAction(n.Defaction),
		Location:     n.Location,
	}
}

func convertDefineStmt(n *nodes.DefineStmt) *pg.DefineStmt {
	if n == nil {
		return nil
	}
	return &pg.DefineStmt{
		Kind:        pg.ObjectType(n.Kind),
		Oldstyle:    n.Oldstyle,
		Defnames:    convertList(n.Defnames),
		Args:        convertList(n.Args),
		Definition:  convertList(n.Definition),
		IfNotExists: n.IfNotExists,
	}
}

func convertDeleteStmt(n *nodes.DeleteStmt) *pg.DeleteStmt {
	if n == nil {
		return nil
	}
	return &pg.DeleteStmt{
		Relation:      convertRangeVar(n.Relation),
		UsingClause:   convertList(n.UsingClause),
		WhereClause:   convertNode(n.WhereClause),
		ReturningList: convertList(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertDiscardStmt(n *nodes.DiscardStmt) *pg.DiscardStmt {
	if n == nil {
		return nil
	}
	return &pg.DiscardStmt{
		Target: pg.DiscardMode(n.Target),
	}
}

func convertDoStmt(n *nodes.DoStmt) *pg.DoStmt {
	if n == nil {
		return nil
	}
	return &pg.DoStmt{
		Args: convertList(n.Args),
	}
}

func convertDropOwnedStmt(n *nodes.DropOwnedStmt) *pg.DropOwnedStmt {
	if n == nil {
		return nil
	}
	return &pg.DropOwnedStmt{
		Roles:    convertList(n.Roles),
		Behavior: pg.DropBehavior(n.Behavior),
	}
}

func convertDropRoleStmt(n *nodes.DropRoleStmt) *pg.DropRoleStmt {
	if n == nil {
		return nil
	}
	return &pg.DropRoleStmt{
		Roles:     convertList(n.Roles),
		MissingOk: n.MissingOk,
	}
}

func convertDropStmt(n *nodes.DropStmt) *pg.DropStmt {
	if n == nil {
		return nil
	}
	return &pg.DropStmt{
		Objects:    convertList(n.Objects),
		RemoveType: pg.ObjectType(n.RemoveType),
		Behavior:   pg.DropBehavior(n.Behavior),
		MissingOk:  n.MissingOk,
		Concurrent: n.Concurrent,
	}
}

func convertDropSubscriptionStmt(n *nodes.DropSubscriptionStmt) *pg.DropSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &pg.DropSubscriptionStmt{
		Subname:   n.Subname,
		MissingOk: n.MissingOk,
		Behavior:  pg.DropBehavior(n.Behavior),
	}
}

func convertDropTableSpaceStmt(n *nodes.DropTableSpaceStmt) *pg.DropTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &pg.DropTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		MissingOk:      n.MissingOk,
	}
}

func convertDropUserMappingStmt(n *nodes.DropUserMappingStmt) *pg.DropUserMappingStmt {
	if n == nil {
		return nil
	}
	return &pg.DropUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: n.Servername,
		MissingOk:  n.MissingOk,
	}
}

func convertDropdbStmt(n *nodes.DropdbStmt) *pg.DropdbStmt {
	if n == nil {
		return nil
	}
	return &pg.DropdbStmt{
		Dbname:    n.Dbname,
		MissingOk: n.MissingOk,
	}
}

func convertExecuteStmt(n *nodes.ExecuteStmt) *pg.ExecuteStmt {
	if n == nil {
		return nil
	}
	return &pg.ExecuteStmt{
		Name:   n.Name,
		Params: convertList(n.Params),
	}
}

func convertExplainStmt(n *nodes.ExplainStmt) *pg.ExplainStmt {
	if n == nil {
		return nil
	}
	return &pg.ExplainStmt{
		Query:   convertNode(n.Query),
		Options: convertList(n.Options),
	}
}

func convertExpr(n *nodes.Expr) *pg.Expr {
	if n == nil {
		return nil
	}
	return &pg.Expr{}
}

func convertFetchStmt(n *nodes.FetchStmt) *pg.FetchStmt {
	if n == nil {
		return nil
	}
	return &pg.FetchStmt{
		Direction:  pg.FetchDirection(n.Direction),
		HowMany:    n.HowMany,
		Portalname: n.Portalname,
		Ismove:     n.Ismove,
	}
}

func convertFieldSelect(n *nodes.FieldSelect) *pg.FieldSelect {
	if n == nil {
		return nil
	}
	return &pg.FieldSelect{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Fieldnum:     pg.AttrNumber(n.Fieldnum),
		Resulttype:   pg.Oid(n.Resulttype),
		Resulttypmod: n.Resulttypmod,
		Resultcollid: pg.Oid(n.Resultcollid),
	}
}

func convertFieldStore(n *nodes.FieldStore) *pg.FieldStore {
	if n == nil {
		return nil
	}
	return &pg.FieldStore{
		Xpr:        convertNode(n.Xpr),
		Arg:        convertNode(n.Arg),
		Newvals:    convertList(n.Newvals),
		Fieldnums:  convertList(n.Fieldnums),
		Resulttype: pg.Oid(n.Resulttype),
	}
}

func convertFloat(n *nodes.Float) *pg.Float {
	if n == nil {
		return nil
	}
	return &pg.Float{
		Str: n.Str,
	}
}

func convertFromExpr(n *nodes.FromExpr) *pg.FromExpr {
	if n == nil {
		return nil
	}
	return &pg.FromExpr{
		Fromlist: convertList(n.Fromlist),
		Quals:    convertNode(n.Quals),
	}
}

func convertFuncCall(n *nodes.FuncCall) *ast.FuncCall {
	if n == nil {
		return nil
	}
	fn, err := parseFuncName(n.Funcname)
	if err != nil {
		// TODO: How should we handle errors?
		panic(err)
	}
	return &ast.FuncCall{
		Func:           fn,
		Funcname:       convertList(n.Funcname),
		Args:           convertList(n.Args),
		AggOrder:       convertList(n.AggOrder),
		AggFilter:      convertNode(n.AggFilter),
		AggWithinGroup: n.AggWithinGroup,
		AggStar:        n.AggStar,
		AggDistinct:    n.AggDistinct,
		FuncVariadic:   n.FuncVariadic,
		Over:           convertWindowDef(n.Over),
		Location:       n.Location,
	}
}

func convertFuncExpr(n *nodes.FuncExpr) *pg.FuncExpr {
	if n == nil {
		return nil
	}
	return &pg.FuncExpr{
		Xpr:            convertNode(n.Xpr),
		Funcid:         pg.Oid(n.Funcid),
		Funcresulttype: pg.Oid(n.Funcresulttype),
		Funcretset:     n.Funcretset,
		Funcvariadic:   n.Funcvariadic,
		Funcformat:     pg.CoercionForm(n.Funcformat),
		Funccollid:     pg.Oid(n.Funccollid),
		Inputcollid:    pg.Oid(n.Inputcollid),
		Args:           convertList(n.Args),
		Location:       n.Location,
	}
}

func convertFunctionParameter(n *nodes.FunctionParameter) *pg.FunctionParameter {
	if n == nil {
		return nil
	}
	return &pg.FunctionParameter{
		Name:    n.Name,
		ArgType: convertTypeName(n.ArgType),
		Mode:    pg.FunctionParameterMode(n.Mode),
		Defexpr: convertNode(n.Defexpr),
	}
}

func convertGrantRoleStmt(n *nodes.GrantRoleStmt) *pg.GrantRoleStmt {
	if n == nil {
		return nil
	}
	return &pg.GrantRoleStmt{
		GrantedRoles: convertList(n.GrantedRoles),
		GranteeRoles: convertList(n.GranteeRoles),
		IsGrant:      n.IsGrant,
		AdminOpt:     n.AdminOpt,
		Grantor:      convertRoleSpec(n.Grantor),
		Behavior:     pg.DropBehavior(n.Behavior),
	}
}

func convertGrantStmt(n *nodes.GrantStmt) *pg.GrantStmt {
	if n == nil {
		return nil
	}
	return &pg.GrantStmt{
		IsGrant:     n.IsGrant,
		Targtype:    pg.GrantTargetType(n.Targtype),
		Objtype:     pg.GrantObjectType(n.Objtype),
		Objects:     convertList(n.Objects),
		Privileges:  convertList(n.Privileges),
		Grantees:    convertList(n.Grantees),
		GrantOption: n.GrantOption,
		Behavior:    pg.DropBehavior(n.Behavior),
	}
}

func convertGroupingFunc(n *nodes.GroupingFunc) *pg.GroupingFunc {
	if n == nil {
		return nil
	}
	return &pg.GroupingFunc{
		Xpr:         convertNode(n.Xpr),
		Args:        convertList(n.Args),
		Refs:        convertList(n.Refs),
		Cols:        convertList(n.Cols),
		Agglevelsup: pg.Index(n.Agglevelsup),
		Location:    n.Location,
	}
}

func convertGroupingSet(n *nodes.GroupingSet) *pg.GroupingSet {
	if n == nil {
		return nil
	}
	return &pg.GroupingSet{
		Kind:     pg.GroupingSetKind(n.Kind),
		Content:  convertList(n.Content),
		Location: n.Location,
	}
}

func convertImportForeignSchemaStmt(n *nodes.ImportForeignSchemaStmt) *pg.ImportForeignSchemaStmt {
	if n == nil {
		return nil
	}
	return &pg.ImportForeignSchemaStmt{
		ServerName:   n.ServerName,
		RemoteSchema: n.RemoteSchema,
		LocalSchema:  n.LocalSchema,
		ListType:     pg.ImportForeignSchemaType(n.ListType),
		TableList:    convertList(n.TableList),
		Options:      convertList(n.Options),
	}
}

func convertIndexElem(n *nodes.IndexElem) *pg.IndexElem {
	if n == nil {
		return nil
	}
	return &pg.IndexElem{
		Name:          n.Name,
		Expr:          convertNode(n.Expr),
		Indexcolname:  n.Indexcolname,
		Collation:     convertList(n.Collation),
		Opclass:       convertList(n.Opclass),
		Ordering:      pg.SortByDir(n.Ordering),
		NullsOrdering: pg.SortByNulls(n.NullsOrdering),
	}
}

func convertIndexStmt(n *nodes.IndexStmt) *pg.IndexStmt {
	if n == nil {
		return nil
	}
	return &pg.IndexStmt{
		Idxname:        n.Idxname,
		Relation:       convertRangeVar(n.Relation),
		AccessMethod:   n.AccessMethod,
		TableSpace:     n.TableSpace,
		IndexParams:    convertList(n.IndexParams),
		Options:        convertList(n.Options),
		WhereClause:    convertNode(n.WhereClause),
		ExcludeOpNames: convertList(n.ExcludeOpNames),
		Idxcomment:     n.Idxcomment,
		IndexOid:       pg.Oid(n.IndexOid),
		OldNode:        pg.Oid(n.OldNode),
		Unique:         n.Unique,
		Primary:        n.Primary,
		Isconstraint:   n.Isconstraint,
		Deferrable:     n.Deferrable,
		Initdeferred:   n.Initdeferred,
		Transformed:    n.Transformed,
		Concurrent:     n.Concurrent,
		IfNotExists:    n.IfNotExists,
	}
}

func convertInferClause(n *nodes.InferClause) *pg.InferClause {
	if n == nil {
		return nil
	}
	return &pg.InferClause{
		IndexElems:  convertList(n.IndexElems),
		WhereClause: convertNode(n.WhereClause),
		Conname:     n.Conname,
		Location:    n.Location,
	}
}

func convertInferenceElem(n *nodes.InferenceElem) *pg.InferenceElem {
	if n == nil {
		return nil
	}
	return &pg.InferenceElem{
		Xpr:          convertNode(n.Xpr),
		Expr:         convertNode(n.Expr),
		Infercollid:  pg.Oid(n.Infercollid),
		Inferopclass: pg.Oid(n.Inferopclass),
	}
}

func convertInlineCodeBlock(n *nodes.InlineCodeBlock) *pg.InlineCodeBlock {
	if n == nil {
		return nil
	}
	return &pg.InlineCodeBlock{
		SourceText:    n.SourceText,
		LangOid:       pg.Oid(n.LangOid),
		LangIsTrusted: n.LangIsTrusted,
	}
}

func convertInsertStmt(n *nodes.InsertStmt) *pg.InsertStmt {
	if n == nil {
		return nil
	}
	return &pg.InsertStmt{
		Relation:         convertRangeVar(n.Relation),
		Cols:             convertList(n.Cols),
		SelectStmt:       convertNode(n.SelectStmt),
		OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		ReturningList:    convertList(n.ReturningList),
		WithClause:       convertWithClause(n.WithClause),
		Override:         pg.OverridingKind(n.Override),
	}
}

func convertInteger(n *nodes.Integer) *pg.Integer {
	if n == nil {
		return nil
	}
	return &pg.Integer{
		Ival: n.Ival,
	}
}

func convertIntoClause(n *nodes.IntoClause) *pg.IntoClause {
	if n == nil {
		return nil
	}
	return &pg.IntoClause{
		Rel:            convertRangeVar(n.Rel),
		ColNames:       convertList(n.ColNames),
		Options:        convertList(n.Options),
		OnCommit:       pg.OnCommitAction(n.OnCommit),
		TableSpaceName: n.TableSpaceName,
		ViewQuery:      convertNode(n.ViewQuery),
		SkipData:       n.SkipData,
	}
}

func convertJoinExpr(n *nodes.JoinExpr) *pg.JoinExpr {
	if n == nil {
		return nil
	}
	return &pg.JoinExpr{
		Jointype:    pg.JoinType(n.Jointype),
		IsNatural:   n.IsNatural,
		Larg:        convertNode(n.Larg),
		Rarg:        convertNode(n.Rarg),
		UsingClause: convertList(n.UsingClause),
		Quals:       convertNode(n.Quals),
		Alias:       convertAlias(n.Alias),
		Rtindex:     n.Rtindex,
	}
}

func convertListenStmt(n *nodes.ListenStmt) *pg.ListenStmt {
	if n == nil {
		return nil
	}
	return &pg.ListenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertLoadStmt(n *nodes.LoadStmt) *pg.LoadStmt {
	if n == nil {
		return nil
	}
	return &pg.LoadStmt{
		Filename: n.Filename,
	}
}

func convertLockStmt(n *nodes.LockStmt) *pg.LockStmt {
	if n == nil {
		return nil
	}
	return &pg.LockStmt{
		Relations: convertList(n.Relations),
		Mode:      n.Mode,
		Nowait:    n.Nowait,
	}
}

func convertLockingClause(n *nodes.LockingClause) *pg.LockingClause {
	if n == nil {
		return nil
	}
	return &pg.LockingClause{
		LockedRels: convertList(n.LockedRels),
		Strength:   pg.LockClauseStrength(n.Strength),
		WaitPolicy: pg.LockWaitPolicy(n.WaitPolicy),
	}
}

func convertMinMaxExpr(n *nodes.MinMaxExpr) *pg.MinMaxExpr {
	if n == nil {
		return nil
	}
	return &pg.MinMaxExpr{
		Xpr:          convertNode(n.Xpr),
		Minmaxtype:   pg.Oid(n.Minmaxtype),
		Minmaxcollid: pg.Oid(n.Minmaxcollid),
		Inputcollid:  pg.Oid(n.Inputcollid),
		Op:           pg.MinMaxOp(n.Op),
		Args:         convertList(n.Args),
		Location:     n.Location,
	}
}

func convertMultiAssignRef(n *nodes.MultiAssignRef) *pg.MultiAssignRef {
	if n == nil {
		return nil
	}
	return &pg.MultiAssignRef{
		Source:   convertNode(n.Source),
		Colno:    n.Colno,
		Ncolumns: n.Ncolumns,
	}
}

func convertNamedArgExpr(n *nodes.NamedArgExpr) *pg.NamedArgExpr {
	if n == nil {
		return nil
	}
	return &pg.NamedArgExpr{
		Xpr:       convertNode(n.Xpr),
		Arg:       convertNode(n.Arg),
		Name:      n.Name,
		Argnumber: n.Argnumber,
		Location:  n.Location,
	}
}

func convertNextValueExpr(n *nodes.NextValueExpr) *pg.NextValueExpr {
	if n == nil {
		return nil
	}
	return &pg.NextValueExpr{
		Xpr:    convertNode(n.Xpr),
		Seqid:  pg.Oid(n.Seqid),
		TypeId: pg.Oid(n.TypeId),
	}
}

func convertNotifyStmt(n *nodes.NotifyStmt) *pg.NotifyStmt {
	if n == nil {
		return nil
	}
	return &pg.NotifyStmt{
		Conditionname: n.Conditionname,
		Payload:       n.Payload,
	}
}

func convertNull(n *nodes.Null) *pg.Null {
	if n == nil {
		return nil
	}
	return &pg.Null{}
}

func convertNullTest(n *nodes.NullTest) *pg.NullTest {
	if n == nil {
		return nil
	}
	return &pg.NullTest{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Nulltesttype: pg.NullTestType(n.Nulltesttype),
		Argisrow:     n.Argisrow,
		Location:     n.Location,
	}
}

func convertObjectWithArgs(n *nodes.ObjectWithArgs) *pg.ObjectWithArgs {
	if n == nil {
		return nil
	}
	return &pg.ObjectWithArgs{
		Objname:         convertList(n.Objname),
		Objargs:         convertList(n.Objargs),
		ArgsUnspecified: n.ArgsUnspecified,
	}
}

func convertOnConflictClause(n *nodes.OnConflictClause) *pg.OnConflictClause {
	if n == nil {
		return nil
	}
	return &pg.OnConflictClause{
		Action:      pg.OnConflictAction(n.Action),
		Infer:       convertInferClause(n.Infer),
		TargetList:  convertList(n.TargetList),
		WhereClause: convertNode(n.WhereClause),
		Location:    n.Location,
	}
}

func convertOnConflictExpr(n *nodes.OnConflictExpr) *pg.OnConflictExpr {
	if n == nil {
		return nil
	}
	return &pg.OnConflictExpr{
		Action:          pg.OnConflictAction(n.Action),
		ArbiterElems:    convertList(n.ArbiterElems),
		ArbiterWhere:    convertNode(n.ArbiterWhere),
		Constraint:      pg.Oid(n.Constraint),
		OnConflictSet:   convertList(n.OnConflictSet),
		OnConflictWhere: convertNode(n.OnConflictWhere),
		ExclRelIndex:    n.ExclRelIndex,
		ExclRelTlist:    convertList(n.ExclRelTlist),
	}
}

func convertOpExpr(n *nodes.OpExpr) *pg.OpExpr {
	if n == nil {
		return nil
	}
	return &pg.OpExpr{
		Xpr:          convertNode(n.Xpr),
		Opno:         pg.Oid(n.Opno),
		Opfuncid:     pg.Oid(n.Opfuncid),
		Opresulttype: pg.Oid(n.Opresulttype),
		Opretset:     n.Opretset,
		Opcollid:     pg.Oid(n.Opcollid),
		Inputcollid:  pg.Oid(n.Inputcollid),
		Args:         convertList(n.Args),
		Location:     n.Location,
	}
}

func convertParam(n *nodes.Param) *pg.Param {
	if n == nil {
		return nil
	}
	return &pg.Param{
		Xpr:         convertNode(n.Xpr),
		Paramkind:   pg.ParamKind(n.Paramkind),
		Paramid:     n.Paramid,
		Paramtype:   pg.Oid(n.Paramtype),
		Paramtypmod: n.Paramtypmod,
		Paramcollid: pg.Oid(n.Paramcollid),
		Location:    n.Location,
	}
}

func convertParamExecData(n *nodes.ParamExecData) *pg.ParamExecData {
	if n == nil {
		return nil
	}
	return &pg.ParamExecData{
		ExecPlan: &ast.TODO{},
		Value:    pg.Datum(n.Value),
		Isnull:   n.Isnull,
	}
}

func convertParamExternData(n *nodes.ParamExternData) *pg.ParamExternData {
	if n == nil {
		return nil
	}
	return &pg.ParamExternData{
		Value:  pg.Datum(n.Value),
		Isnull: n.Isnull,
		Pflags: n.Pflags,
		Ptype:  pg.Oid(n.Ptype),
	}
}

func convertParamListInfoData(n *nodes.ParamListInfoData) *pg.ParamListInfoData {
	if n == nil {
		return nil
	}
	return &pg.ParamListInfoData{
		ParamFetchArg:  &ast.TODO{},
		ParserSetupArg: &ast.TODO{},
		NumParams:      n.NumParams,
		ParamMask:      n.ParamMask,
	}
}

func convertParamRef(n *nodes.ParamRef) *pg.ParamRef {
	if n == nil {
		return nil
	}
	return &pg.ParamRef{
		Number:   n.Number,
		Location: n.Location,
	}
}

func convertPartitionBoundSpec(n *nodes.PartitionBoundSpec) *pg.PartitionBoundSpec {
	if n == nil {
		return nil
	}
	return &pg.PartitionBoundSpec{
		Strategy:    n.Strategy,
		Listdatums:  convertList(n.Listdatums),
		Lowerdatums: convertList(n.Lowerdatums),
		Upperdatums: convertList(n.Upperdatums),
		Location:    n.Location,
	}
}

func convertPartitionCmd(n *nodes.PartitionCmd) *pg.PartitionCmd {
	if n == nil {
		return nil
	}
	return &pg.PartitionCmd{
		Name:  convertRangeVar(n.Name),
		Bound: convertPartitionBoundSpec(n.Bound),
	}
}

func convertPartitionElem(n *nodes.PartitionElem) *pg.PartitionElem {
	if n == nil {
		return nil
	}
	return &pg.PartitionElem{
		Name:      n.Name,
		Expr:      convertNode(n.Expr),
		Collation: convertList(n.Collation),
		Opclass:   convertList(n.Opclass),
		Location:  n.Location,
	}
}

func convertPartitionRangeDatum(n *nodes.PartitionRangeDatum) *pg.PartitionRangeDatum {
	if n == nil {
		return nil
	}
	return &pg.PartitionRangeDatum{
		Kind:     pg.PartitionRangeDatumKind(n.Kind),
		Value:    convertNode(n.Value),
		Location: n.Location,
	}
}

func convertPartitionSpec(n *nodes.PartitionSpec) *pg.PartitionSpec {
	if n == nil {
		return nil
	}
	return &pg.PartitionSpec{
		Strategy:   n.Strategy,
		PartParams: convertList(n.PartParams),
		Location:   n.Location,
	}
}

func convertPrepareStmt(n *nodes.PrepareStmt) *pg.PrepareStmt {
	if n == nil {
		return nil
	}
	return &pg.PrepareStmt{
		Name:     n.Name,
		Argtypes: convertList(n.Argtypes),
		Query:    convertNode(n.Query),
	}
}

func convertQuery(n *nodes.Query) *pg.Query {
	if n == nil {
		return nil
	}
	return &pg.Query{
		CommandType:      pg.CmdType(n.CommandType),
		QuerySource:      pg.QuerySource(n.QuerySource),
		QueryId:          n.QueryId,
		CanSetTag:        n.CanSetTag,
		UtilityStmt:      convertNode(n.UtilityStmt),
		ResultRelation:   n.ResultRelation,
		HasAggs:          n.HasAggs,
		HasWindowFuncs:   n.HasWindowFuncs,
		HasTargetSrfs:    n.HasTargetSrfs,
		HasSubLinks:      n.HasSubLinks,
		HasDistinctOn:    n.HasDistinctOn,
		HasRecursive:     n.HasRecursive,
		HasModifyingCte:  n.HasModifyingCte,
		HasForUpdate:     n.HasForUpdate,
		HasRowSecurity:   n.HasRowSecurity,
		CteList:          convertList(n.CteList),
		Rtable:           convertList(n.Rtable),
		Jointree:         convertFromExpr(n.Jointree),
		TargetList:       convertList(n.TargetList),
		Override:         pg.OverridingKind(n.Override),
		OnConflict:       convertOnConflictExpr(n.OnConflict),
		ReturningList:    convertList(n.ReturningList),
		GroupClause:      convertList(n.GroupClause),
		GroupingSets:     convertList(n.GroupingSets),
		HavingQual:       convertNode(n.HavingQual),
		WindowClause:     convertList(n.WindowClause),
		DistinctClause:   convertList(n.DistinctClause),
		SortClause:       convertList(n.SortClause),
		LimitOffset:      convertNode(n.LimitOffset),
		LimitCount:       convertNode(n.LimitCount),
		RowMarks:         convertList(n.RowMarks),
		SetOperations:    convertNode(n.SetOperations),
		ConstraintDeps:   convertList(n.ConstraintDeps),
		WithCheckOptions: convertList(n.WithCheckOptions),
		StmtLocation:     n.StmtLocation,
		StmtLen:          n.StmtLen,
	}
}

func convertRangeFunction(n *nodes.RangeFunction) *pg.RangeFunction {
	if n == nil {
		return nil
	}
	return &pg.RangeFunction{
		Lateral:    n.Lateral,
		Ordinality: n.Ordinality,
		IsRowsfrom: n.IsRowsfrom,
		Functions:  convertList(n.Functions),
		Alias:      convertAlias(n.Alias),
		Coldeflist: convertList(n.Coldeflist),
	}
}

func convertRangeSubselect(n *nodes.RangeSubselect) *pg.RangeSubselect {
	if n == nil {
		return nil
	}
	return &pg.RangeSubselect{
		Lateral:  n.Lateral,
		Subquery: convertNode(n.Subquery),
		Alias:    convertAlias(n.Alias),
	}
}

func convertRangeTableFunc(n *nodes.RangeTableFunc) *pg.RangeTableFunc {
	if n == nil {
		return nil
	}
	return &pg.RangeTableFunc{
		Lateral:    n.Lateral,
		Docexpr:    convertNode(n.Docexpr),
		Rowexpr:    convertNode(n.Rowexpr),
		Namespaces: convertList(n.Namespaces),
		Columns:    convertList(n.Columns),
		Alias:      convertAlias(n.Alias),
		Location:   n.Location,
	}
}

func convertRangeTableFuncCol(n *nodes.RangeTableFuncCol) *pg.RangeTableFuncCol {
	if n == nil {
		return nil
	}
	return &pg.RangeTableFuncCol{
		Colname:       n.Colname,
		TypeName:      convertTypeName(n.TypeName),
		ForOrdinality: n.ForOrdinality,
		IsNotNull:     n.IsNotNull,
		Colexpr:       convertNode(n.Colexpr),
		Coldefexpr:    convertNode(n.Coldefexpr),
		Location:      n.Location,
	}
}

func convertRangeTableSample(n *nodes.RangeTableSample) *pg.RangeTableSample {
	if n == nil {
		return nil
	}
	return &pg.RangeTableSample{
		Relation:   convertNode(n.Relation),
		Method:     convertList(n.Method),
		Args:       convertList(n.Args),
		Repeatable: convertNode(n.Repeatable),
		Location:   n.Location,
	}
}

func convertRangeTblEntry(n *nodes.RangeTblEntry) *pg.RangeTblEntry {
	if n == nil {
		return nil
	}
	return &pg.RangeTblEntry{
		Rtekind:         pg.RTEKind(n.Rtekind),
		Relid:           pg.Oid(n.Relid),
		Relkind:         n.Relkind,
		Tablesample:     convertTableSampleClause(n.Tablesample),
		Subquery:        convertQuery(n.Subquery),
		SecurityBarrier: n.SecurityBarrier,
		Jointype:        pg.JoinType(n.Jointype),
		Joinaliasvars:   convertList(n.Joinaliasvars),
		Functions:       convertList(n.Functions),
		Funcordinality:  n.Funcordinality,
		Tablefunc:       convertTableFunc(n.Tablefunc),
		ValuesLists:     convertList(n.ValuesLists),
		Ctename:         n.Ctename,
		Ctelevelsup:     pg.Index(n.Ctelevelsup),
		SelfReference:   n.SelfReference,
		Coltypes:        convertList(n.Coltypes),
		Coltypmods:      convertList(n.Coltypmods),
		Colcollations:   convertList(n.Colcollations),
		Enrname:         n.Enrname,
		Enrtuples:       n.Enrtuples,
		Alias:           convertAlias(n.Alias),
		Eref:            convertAlias(n.Eref),
		Lateral:         n.Lateral,
		Inh:             n.Inh,
		InFromCl:        n.InFromCl,
		RequiredPerms:   pg.AclMode(n.RequiredPerms),
		CheckAsUser:     pg.Oid(n.CheckAsUser),
		SelectedCols:    n.SelectedCols,
		InsertedCols:    n.InsertedCols,
		UpdatedCols:     n.UpdatedCols,
		SecurityQuals:   convertList(n.SecurityQuals),
	}
}

func convertRangeTblFunction(n *nodes.RangeTblFunction) *pg.RangeTblFunction {
	if n == nil {
		return nil
	}
	return &pg.RangeTblFunction{
		Funcexpr:          convertNode(n.Funcexpr),
		Funccolcount:      n.Funccolcount,
		Funccolnames:      convertList(n.Funccolnames),
		Funccoltypes:      convertList(n.Funccoltypes),
		Funccoltypmods:    convertList(n.Funccoltypmods),
		Funccolcollations: convertList(n.Funccolcollations),
		Funcparams:        n.Funcparams,
	}
}

func convertRangeTblRef(n *nodes.RangeTblRef) *pg.RangeTblRef {
	if n == nil {
		return nil
	}
	return &pg.RangeTblRef{
		Rtindex: n.Rtindex,
	}
}

func convertRangeVar(n *nodes.RangeVar) *pg.RangeVar {
	if n == nil {
		return nil
	}
	return &pg.RangeVar{
		Catalogname:    n.Catalogname,
		Schemaname:     n.Schemaname,
		Relname:        n.Relname,
		Inh:            n.Inh,
		Relpersistence: n.Relpersistence,
		Alias:          convertAlias(n.Alias),
		Location:       n.Location,
	}
}

func convertRawStmt(n *nodes.RawStmt) *pg.RawStmt {
	if n == nil {
		return nil
	}
	return &pg.RawStmt{
		Stmt:         convertNode(n.Stmt),
		StmtLocation: n.StmtLocation,
		StmtLen:      n.StmtLen,
	}
}

func convertReassignOwnedStmt(n *nodes.ReassignOwnedStmt) *pg.ReassignOwnedStmt {
	if n == nil {
		return nil
	}
	return &pg.ReassignOwnedStmt{
		Roles:   convertList(n.Roles),
		Newrole: convertRoleSpec(n.Newrole),
	}
}

func convertRefreshMatViewStmt(n *nodes.RefreshMatViewStmt) *pg.RefreshMatViewStmt {
	if n == nil {
		return nil
	}
	return &pg.RefreshMatViewStmt{
		Concurrent: n.Concurrent,
		SkipData:   n.SkipData,
		Relation:   convertRangeVar(n.Relation),
	}
}

func convertReindexStmt(n *nodes.ReindexStmt) *pg.ReindexStmt {
	if n == nil {
		return nil
	}
	return &pg.ReindexStmt{
		Kind:     pg.ReindexObjectType(n.Kind),
		Relation: convertRangeVar(n.Relation),
		Name:     n.Name,
		Options:  n.Options,
	}
}

func convertRelabelType(n *nodes.RelabelType) *pg.RelabelType {
	if n == nil {
		return nil
	}
	return &pg.RelabelType{
		Xpr:           convertNode(n.Xpr),
		Arg:           convertNode(n.Arg),
		Resulttype:    pg.Oid(n.Resulttype),
		Resulttypmod:  n.Resulttypmod,
		Resultcollid:  pg.Oid(n.Resultcollid),
		Relabelformat: pg.CoercionForm(n.Relabelformat),
		Location:      n.Location,
	}
}

func convertRenameStmt(n *nodes.RenameStmt) *pg.RenameStmt {
	if n == nil {
		return nil
	}
	return &pg.RenameStmt{
		RenameType:   pg.ObjectType(n.RenameType),
		RelationType: pg.ObjectType(n.RelationType),
		Relation:     convertRangeVar(n.Relation),
		Object:       convertNode(n.Object),
		Subname:      n.Subname,
		Newname:      n.Newname,
		Behavior:     pg.DropBehavior(n.Behavior),
		MissingOk:    n.MissingOk,
	}
}

func convertReplicaIdentityStmt(n *nodes.ReplicaIdentityStmt) *pg.ReplicaIdentityStmt {
	if n == nil {
		return nil
	}
	return &pg.ReplicaIdentityStmt{
		IdentityType: n.IdentityType,
		Name:         n.Name,
	}
}

func convertResTarget(n *nodes.ResTarget) *pg.ResTarget {
	if n == nil {
		return nil
	}
	return &pg.ResTarget{
		Name:        n.Name,
		Indirection: convertList(n.Indirection),
		Val:         convertNode(n.Val),
		Location:    n.Location,
	}
}

func convertRoleSpec(n *nodes.RoleSpec) *pg.RoleSpec {
	if n == nil {
		return nil
	}
	return &pg.RoleSpec{
		Roletype: pg.RoleSpecType(n.Roletype),
		Rolename: n.Rolename,
		Location: n.Location,
	}
}

func convertRowCompareExpr(n *nodes.RowCompareExpr) *pg.RowCompareExpr {
	if n == nil {
		return nil
	}
	return &pg.RowCompareExpr{
		Xpr:          convertNode(n.Xpr),
		Rctype:       pg.RowCompareType(n.Rctype),
		Opnos:        convertList(n.Opnos),
		Opfamilies:   convertList(n.Opfamilies),
		Inputcollids: convertList(n.Inputcollids),
		Largs:        convertList(n.Largs),
		Rargs:        convertList(n.Rargs),
	}
}

func convertRowExpr(n *nodes.RowExpr) *pg.RowExpr {
	if n == nil {
		return nil
	}
	return &pg.RowExpr{
		Xpr:       convertNode(n.Xpr),
		Args:      convertList(n.Args),
		RowTypeid: pg.Oid(n.RowTypeid),
		RowFormat: pg.CoercionForm(n.RowFormat),
		Colnames:  convertList(n.Colnames),
		Location:  n.Location,
	}
}

func convertRowMarkClause(n *nodes.RowMarkClause) *pg.RowMarkClause {
	if n == nil {
		return nil
	}
	return &pg.RowMarkClause{
		Rti:        pg.Index(n.Rti),
		Strength:   pg.LockClauseStrength(n.Strength),
		WaitPolicy: pg.LockWaitPolicy(n.WaitPolicy),
		PushedDown: n.PushedDown,
	}
}

func convertRuleStmt(n *nodes.RuleStmt) *pg.RuleStmt {
	if n == nil {
		return nil
	}
	return &pg.RuleStmt{
		Relation:    convertRangeVar(n.Relation),
		Rulename:    n.Rulename,
		WhereClause: convertNode(n.WhereClause),
		Event:       pg.CmdType(n.Event),
		Instead:     n.Instead,
		Actions:     convertList(n.Actions),
		Replace:     n.Replace,
	}
}

func convertSQLValueFunction(n *nodes.SQLValueFunction) *pg.SQLValueFunction {
	if n == nil {
		return nil
	}
	return &pg.SQLValueFunction{
		Xpr:      convertNode(n.Xpr),
		Op:       pg.SQLValueFunctionOp(n.Op),
		Type:     pg.Oid(n.Type),
		Typmod:   n.Typmod,
		Location: n.Location,
	}
}

func convertScalarArrayOpExpr(n *nodes.ScalarArrayOpExpr) *pg.ScalarArrayOpExpr {
	if n == nil {
		return nil
	}
	return &pg.ScalarArrayOpExpr{
		Xpr:         convertNode(n.Xpr),
		Opno:        pg.Oid(n.Opno),
		Opfuncid:    pg.Oid(n.Opfuncid),
		UseOr:       n.UseOr,
		Inputcollid: pg.Oid(n.Inputcollid),
		Args:        convertList(n.Args),
		Location:    n.Location,
	}
}

func convertSecLabelStmt(n *nodes.SecLabelStmt) *pg.SecLabelStmt {
	if n == nil {
		return nil
	}
	return &pg.SecLabelStmt{
		Objtype:  pg.ObjectType(n.Objtype),
		Object:   convertNode(n.Object),
		Provider: n.Provider,
		Label:    n.Label,
	}
}

func convertSelectStmt(n *nodes.SelectStmt) *pg.SelectStmt {
	if n == nil {
		return nil
	}
	return &pg.SelectStmt{
		DistinctClause: convertList(n.DistinctClause),
		IntoClause:     convertIntoClause(n.IntoClause),
		TargetList:     convertList(n.TargetList),
		FromClause:     convertList(n.FromClause),
		WhereClause:    convertNode(n.WhereClause),
		GroupClause:    convertList(n.GroupClause),
		HavingClause:   convertNode(n.HavingClause),
		WindowClause:   convertList(n.WindowClause),
		ValuesLists:    convertValuesList(n.ValuesLists),
		SortClause:     convertList(n.SortClause),
		LimitOffset:    convertNode(n.LimitOffset),
		LimitCount:     convertNode(n.LimitCount),
		LockingClause:  convertList(n.LockingClause),
		WithClause:     convertWithClause(n.WithClause),
		Op:             pg.SetOperation(n.Op),
		All:            n.All,
		Larg:           convertSelectStmt(n.Larg),
		Rarg:           convertSelectStmt(n.Rarg),
	}
}

func convertSetOperationStmt(n *nodes.SetOperationStmt) *pg.SetOperationStmt {
	if n == nil {
		return nil
	}
	return &pg.SetOperationStmt{
		Op:            pg.SetOperation(n.Op),
		All:           n.All,
		Larg:          convertNode(n.Larg),
		Rarg:          convertNode(n.Rarg),
		ColTypes:      convertList(n.ColTypes),
		ColTypmods:    convertList(n.ColTypmods),
		ColCollations: convertList(n.ColCollations),
		GroupClauses:  convertList(n.GroupClauses),
	}
}

func convertSetToDefault(n *nodes.SetToDefault) *pg.SetToDefault {
	if n == nil {
		return nil
	}
	return &pg.SetToDefault{
		Xpr:       convertNode(n.Xpr),
		TypeId:    pg.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: pg.Oid(n.Collation),
		Location:  n.Location,
	}
}

func convertSortBy(n *nodes.SortBy) *pg.SortBy {
	if n == nil {
		return nil
	}
	return &pg.SortBy{
		Node:        convertNode(n.Node),
		SortbyDir:   pg.SortByDir(n.SortbyDir),
		SortbyNulls: pg.SortByNulls(n.SortbyNulls),
		UseOp:       convertList(n.UseOp),
		Location:    n.Location,
	}
}

func convertSortGroupClause(n *nodes.SortGroupClause) *pg.SortGroupClause {
	if n == nil {
		return nil
	}
	return &pg.SortGroupClause{
		TleSortGroupRef: pg.Index(n.TleSortGroupRef),
		Eqop:            pg.Oid(n.Eqop),
		Sortop:          pg.Oid(n.Sortop),
		NullsFirst:      n.NullsFirst,
		Hashable:        n.Hashable,
	}
}

func convertString(n *nodes.String) *pg.String {
	if n == nil {
		return nil
	}
	return &pg.String{
		Str: n.Str,
	}
}

func convertSubLink(n *nodes.SubLink) *pg.SubLink {
	if n == nil {
		return nil
	}
	return &pg.SubLink{
		Xpr:         convertNode(n.Xpr),
		SubLinkType: pg.SubLinkType(n.SubLinkType),
		SubLinkId:   n.SubLinkId,
		Testexpr:    convertNode(n.Testexpr),
		OperName:    convertList(n.OperName),
		Subselect:   convertNode(n.Subselect),
		Location:    n.Location,
	}
}

func convertSubPlan(n *nodes.SubPlan) *pg.SubPlan {
	if n == nil {
		return nil
	}
	return &pg.SubPlan{
		Xpr:               convertNode(n.Xpr),
		SubLinkType:       pg.SubLinkType(n.SubLinkType),
		Testexpr:          convertNode(n.Testexpr),
		ParamIds:          convertList(n.ParamIds),
		PlanId:            n.PlanId,
		PlanName:          n.PlanName,
		FirstColType:      pg.Oid(n.FirstColType),
		FirstColTypmod:    n.FirstColTypmod,
		FirstColCollation: pg.Oid(n.FirstColCollation),
		UseHashTable:      n.UseHashTable,
		UnknownEqFalse:    n.UnknownEqFalse,
		ParallelSafe:      n.ParallelSafe,
		SetParam:          convertList(n.SetParam),
		ParParam:          convertList(n.ParParam),
		Args:              convertList(n.Args),
		StartupCost:       pg.Cost(n.StartupCost),
		PerCallCost:       pg.Cost(n.PerCallCost),
	}
}

func convertTableFunc(n *nodes.TableFunc) *pg.TableFunc {
	if n == nil {
		return nil
	}
	return &pg.TableFunc{
		NsUris:        convertList(n.NsUris),
		NsNames:       convertList(n.NsNames),
		Docexpr:       convertNode(n.Docexpr),
		Rowexpr:       convertNode(n.Rowexpr),
		Colnames:      convertList(n.Colnames),
		Coltypes:      convertList(n.Coltypes),
		Coltypmods:    convertList(n.Coltypmods),
		Colcollations: convertList(n.Colcollations),
		Colexprs:      convertList(n.Colexprs),
		Coldefexprs:   convertList(n.Coldefexprs),
		Notnulls:      n.Notnulls,
		Ordinalitycol: n.Ordinalitycol,
		Location:      n.Location,
	}
}

func convertTableLikeClause(n *nodes.TableLikeClause) *pg.TableLikeClause {
	if n == nil {
		return nil
	}
	return &pg.TableLikeClause{
		Relation: convertRangeVar(n.Relation),
		Options:  n.Options,
	}
}

func convertTableSampleClause(n *nodes.TableSampleClause) *pg.TableSampleClause {
	if n == nil {
		return nil
	}
	return &pg.TableSampleClause{
		Tsmhandler: pg.Oid(n.Tsmhandler),
		Args:       convertList(n.Args),
		Repeatable: convertNode(n.Repeatable),
	}
}

func convertTargetEntry(n *nodes.TargetEntry) *pg.TargetEntry {
	if n == nil {
		return nil
	}
	return &pg.TargetEntry{
		Xpr:             convertNode(n.Xpr),
		Expr:            convertNode(n.Expr),
		Resno:           pg.AttrNumber(n.Resno),
		Resname:         n.Resname,
		Ressortgroupref: pg.Index(n.Ressortgroupref),
		Resorigtbl:      pg.Oid(n.Resorigtbl),
		Resorigcol:      pg.AttrNumber(n.Resorigcol),
		Resjunk:         n.Resjunk,
	}
}

func convertTransactionStmt(n *nodes.TransactionStmt) *pg.TransactionStmt {
	if n == nil {
		return nil
	}
	return &pg.TransactionStmt{
		Kind:    pg.TransactionStmtKind(n.Kind),
		Options: convertList(n.Options),
		Gid:     n.Gid,
	}
}

func convertTriggerTransition(n *nodes.TriggerTransition) *pg.TriggerTransition {
	if n == nil {
		return nil
	}
	return &pg.TriggerTransition{
		Name:    n.Name,
		IsNew:   n.IsNew,
		IsTable: n.IsTable,
	}
}

func convertTruncateStmt(n *nodes.TruncateStmt) *pg.TruncateStmt {
	if n == nil {
		return nil
	}
	return &pg.TruncateStmt{
		Relations:   convertList(n.Relations),
		RestartSeqs: n.RestartSeqs,
		Behavior:    pg.DropBehavior(n.Behavior),
	}
}

func convertTypeCast(n *nodes.TypeCast) *pg.TypeCast {
	if n == nil {
		return nil
	}
	return &pg.TypeCast{
		Arg:      convertNode(n.Arg),
		TypeName: convertTypeName(n.TypeName),
		Location: n.Location,
	}
}

func convertTypeName(n *nodes.TypeName) *pg.TypeName {
	if n == nil {
		return nil
	}
	return &pg.TypeName{
		Names:       convertList(n.Names),
		TypeOid:     pg.Oid(n.TypeOid),
		Setof:       n.Setof,
		PctType:     n.PctType,
		Typmods:     convertList(n.Typmods),
		Typemod:     n.Typemod,
		ArrayBounds: convertList(n.ArrayBounds),
		Location:    n.Location,
	}
}

func convertUnlistenStmt(n *nodes.UnlistenStmt) *pg.UnlistenStmt {
	if n == nil {
		return nil
	}
	return &pg.UnlistenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertUpdateStmt(n *nodes.UpdateStmt) *pg.UpdateStmt {
	if n == nil {
		return nil
	}
	return &pg.UpdateStmt{
		Relation:      convertRangeVar(n.Relation),
		TargetList:    convertList(n.TargetList),
		WhereClause:   convertNode(n.WhereClause),
		FromClause:    convertList(n.FromClause),
		ReturningList: convertList(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertVacuumStmt(n *nodes.VacuumStmt) *pg.VacuumStmt {
	if n == nil {
		return nil
	}
	return &pg.VacuumStmt{
		Options:  n.Options,
		Relation: convertRangeVar(n.Relation),
		VaCols:   convertList(n.VaCols),
	}
}

func convertVar(n *nodes.Var) *pg.Var {
	if n == nil {
		return nil
	}
	return &pg.Var{
		Xpr:         convertNode(n.Xpr),
		Varno:       pg.Index(n.Varno),
		Varattno:    pg.AttrNumber(n.Varattno),
		Vartype:     pg.Oid(n.Vartype),
		Vartypmod:   n.Vartypmod,
		Varcollid:   pg.Oid(n.Varcollid),
		Varlevelsup: pg.Index(n.Varlevelsup),
		Varnoold:    pg.Index(n.Varnoold),
		Varoattno:   pg.AttrNumber(n.Varoattno),
		Location:    n.Location,
	}
}

func convertVariableSetStmt(n *nodes.VariableSetStmt) *pg.VariableSetStmt {
	if n == nil {
		return nil
	}
	return &pg.VariableSetStmt{
		Kind:    pg.VariableSetKind(n.Kind),
		Name:    n.Name,
		Args:    convertList(n.Args),
		IsLocal: n.IsLocal,
	}
}

func convertVariableShowStmt(n *nodes.VariableShowStmt) *pg.VariableShowStmt {
	if n == nil {
		return nil
	}
	return &pg.VariableShowStmt{
		Name: n.Name,
	}
}

func convertViewStmt(n *nodes.ViewStmt) *pg.ViewStmt {
	if n == nil {
		return nil
	}
	return &pg.ViewStmt{
		View:            convertRangeVar(n.View),
		Aliases:         convertList(n.Aliases),
		Query:           convertNode(n.Query),
		Replace:         n.Replace,
		Options:         convertList(n.Options),
		WithCheckOption: pg.ViewCheckOption(n.WithCheckOption),
	}
}

func convertWindowClause(n *nodes.WindowClause) *pg.WindowClause {
	if n == nil {
		return nil
	}
	return &pg.WindowClause{
		Name:            n.Name,
		Refname:         n.Refname,
		PartitionClause: convertList(n.PartitionClause),
		OrderClause:     convertList(n.OrderClause),
		FrameOptions:    n.FrameOptions,
		StartOffset:     convertNode(n.StartOffset),
		EndOffset:       convertNode(n.EndOffset),
		Winref:          pg.Index(n.Winref),
		CopiedOrder:     n.CopiedOrder,
	}
}

func convertWindowDef(n *nodes.WindowDef) *ast.WindowDef {
	if n == nil {
		return nil
	}
	return &ast.WindowDef{
		Name:            n.Name,
		Refname:         n.Refname,
		PartitionClause: convertList(n.PartitionClause),
		OrderClause:     convertList(n.OrderClause),
		FrameOptions:    n.FrameOptions,
		StartOffset:     convertNode(n.StartOffset),
		EndOffset:       convertNode(n.EndOffset),
		Location:        n.Location,
	}
}

func convertWindowFunc(n *nodes.WindowFunc) *pg.WindowFunc {
	if n == nil {
		return nil
	}
	return &pg.WindowFunc{
		Xpr:         convertNode(n.Xpr),
		Winfnoid:    pg.Oid(n.Winfnoid),
		Wintype:     pg.Oid(n.Wintype),
		Wincollid:   pg.Oid(n.Wincollid),
		Inputcollid: pg.Oid(n.Inputcollid),
		Args:        convertList(n.Args),
		Aggfilter:   convertNode(n.Aggfilter),
		Winref:      pg.Index(n.Winref),
		Winstar:     n.Winstar,
		Winagg:      n.Winagg,
		Location:    n.Location,
	}
}

func convertWithCheckOption(n *nodes.WithCheckOption) *pg.WithCheckOption {
	if n == nil {
		return nil
	}
	return &pg.WithCheckOption{
		Kind:     pg.WCOKind(n.Kind),
		Relname:  n.Relname,
		Polname:  n.Polname,
		Qual:     convertNode(n.Qual),
		Cascaded: n.Cascaded,
	}
}

func convertWithClause(n *nodes.WithClause) *pg.WithClause {
	if n == nil {
		return nil
	}
	return &pg.WithClause{
		Ctes:      convertList(n.Ctes),
		Recursive: n.Recursive,
		Location:  n.Location,
	}
}

func convertXmlExpr(n *nodes.XmlExpr) *pg.XmlExpr {
	if n == nil {
		return nil
	}
	return &pg.XmlExpr{
		Xpr:       convertNode(n.Xpr),
		Op:        pg.XmlExprOp(n.Op),
		Name:      n.Name,
		NamedArgs: convertList(n.NamedArgs),
		ArgNames:  convertList(n.ArgNames),
		Args:      convertList(n.Args),
		Xmloption: pg.XmlOptionType(n.Xmloption),
		Type:      pg.Oid(n.Type),
		Typmod:    n.Typmod,
		Location:  n.Location,
	}
}

func convertXmlSerialize(n *nodes.XmlSerialize) *pg.XmlSerialize {
	if n == nil {
		return nil
	}
	return &pg.XmlSerialize{
		Xmloption: pg.XmlOptionType(n.Xmloption),
		Expr:      convertNode(n.Expr),
		TypeName:  convertTypeName(n.TypeName),
		Location:  n.Location,
	}
}

func convertNode(node nodes.Node) ast.Node {
	switch n := node.(type) {

	case nodes.A_ArrayExpr:
		return convertA_ArrayExpr(&n)

	case nodes.A_Const:
		return convertA_Const(&n)

	case nodes.A_Expr:
		return convertA_Expr(&n)

	case nodes.A_Indices:
		return convertA_Indices(&n)

	case nodes.A_Indirection:
		return convertA_Indirection(&n)

	case nodes.A_Star:
		return convertA_Star(&n)

	case nodes.AccessPriv:
		return convertAccessPriv(&n)

	case nodes.Aggref:
		return convertAggref(&n)

	case nodes.Alias:
		return convertAlias(&n)

	case nodes.AlterCollationStmt:
		return convertAlterCollationStmt(&n)

	case nodes.AlterDatabaseSetStmt:
		return convertAlterDatabaseSetStmt(&n)

	case nodes.AlterDatabaseStmt:
		return convertAlterDatabaseStmt(&n)

	case nodes.AlterDefaultPrivilegesStmt:
		return convertAlterDefaultPrivilegesStmt(&n)

	case nodes.AlterDomainStmt:
		return convertAlterDomainStmt(&n)

	case nodes.AlterEnumStmt:
		return convertAlterEnumStmt(&n)

	case nodes.AlterEventTrigStmt:
		return convertAlterEventTrigStmt(&n)

	case nodes.AlterExtensionContentsStmt:
		return convertAlterExtensionContentsStmt(&n)

	case nodes.AlterExtensionStmt:
		return convertAlterExtensionStmt(&n)

	case nodes.AlterFdwStmt:
		return convertAlterFdwStmt(&n)

	case nodes.AlterForeignServerStmt:
		return convertAlterForeignServerStmt(&n)

	case nodes.AlterFunctionStmt:
		return convertAlterFunctionStmt(&n)

	case nodes.AlterObjectDependsStmt:
		return convertAlterObjectDependsStmt(&n)

	case nodes.AlterObjectSchemaStmt:
		return convertAlterObjectSchemaStmt(&n)

	case nodes.AlterOpFamilyStmt:
		return convertAlterOpFamilyStmt(&n)

	case nodes.AlterOperatorStmt:
		return convertAlterOperatorStmt(&n)

	case nodes.AlterOwnerStmt:
		return convertAlterOwnerStmt(&n)

	case nodes.AlterPolicyStmt:
		return convertAlterPolicyStmt(&n)

	case nodes.AlterPublicationStmt:
		return convertAlterPublicationStmt(&n)

	case nodes.AlterRoleSetStmt:
		return convertAlterRoleSetStmt(&n)

	case nodes.AlterRoleStmt:
		return convertAlterRoleStmt(&n)

	case nodes.AlterSeqStmt:
		return convertAlterSeqStmt(&n)

	case nodes.AlterSubscriptionStmt:
		return convertAlterSubscriptionStmt(&n)

	case nodes.AlterSystemStmt:
		return convertAlterSystemStmt(&n)

	case nodes.AlterTSConfigurationStmt:
		return convertAlterTSConfigurationStmt(&n)

	case nodes.AlterTSDictionaryStmt:
		return convertAlterTSDictionaryStmt(&n)

	case nodes.AlterTableCmd:
		return convertAlterTableCmd(&n)

	case nodes.AlterTableMoveAllStmt:
		return convertAlterTableMoveAllStmt(&n)

	case nodes.AlterTableSpaceOptionsStmt:
		return convertAlterTableSpaceOptionsStmt(&n)

	case nodes.AlterTableStmt:
		return convertAlterTableStmt(&n)

	case nodes.AlterUserMappingStmt:
		return convertAlterUserMappingStmt(&n)

	case nodes.AlternativeSubPlan:
		return convertAlternativeSubPlan(&n)

	case nodes.ArrayCoerceExpr:
		return convertArrayCoerceExpr(&n)

	case nodes.ArrayExpr:
		return convertArrayExpr(&n)

	case nodes.ArrayRef:
		return convertArrayRef(&n)

	case nodes.BitString:
		return convertBitString(&n)

	case nodes.BlockIdData:
		return convertBlockIdData(&n)

	case nodes.BoolExpr:
		return convertBoolExpr(&n)

	case nodes.BooleanTest:
		return convertBooleanTest(&n)

	case nodes.CaseExpr:
		return convertCaseExpr(&n)

	case nodes.CaseTestExpr:
		return convertCaseTestExpr(&n)

	case nodes.CaseWhen:
		return convertCaseWhen(&n)

	case nodes.CheckPointStmt:
		return convertCheckPointStmt(&n)

	case nodes.ClosePortalStmt:
		return convertClosePortalStmt(&n)

	case nodes.ClusterStmt:
		return convertClusterStmt(&n)

	case nodes.CoalesceExpr:
		return convertCoalesceExpr(&n)

	case nodes.CoerceToDomain:
		return convertCoerceToDomain(&n)

	case nodes.CoerceToDomainValue:
		return convertCoerceToDomainValue(&n)

	case nodes.CoerceViaIO:
		return convertCoerceViaIO(&n)

	case nodes.CollateClause:
		return convertCollateClause(&n)

	case nodes.CollateExpr:
		return convertCollateExpr(&n)

	case nodes.ColumnDef:
		return convertColumnDef(&n)

	case nodes.ColumnRef:
		return convertColumnRef(&n)

	case nodes.CommentStmt:
		return convertCommentStmt(&n)

	case nodes.CommonTableExpr:
		return convertCommonTableExpr(&n)

	case nodes.CompositeTypeStmt:
		return convertCompositeTypeStmt(&n)

	case nodes.Const:
		return convertConst(&n)

	case nodes.Constraint:
		return convertConstraint(&n)

	case nodes.ConstraintsSetStmt:
		return convertConstraintsSetStmt(&n)

	case nodes.ConvertRowtypeExpr:
		return convertConvertRowtypeExpr(&n)

	case nodes.CopyStmt:
		return convertCopyStmt(&n)

	case nodes.CreateAmStmt:
		return convertCreateAmStmt(&n)

	case nodes.CreateCastStmt:
		return convertCreateCastStmt(&n)

	case nodes.CreateConversionStmt:
		return convertCreateConversionStmt(&n)

	case nodes.CreateDomainStmt:
		return convertCreateDomainStmt(&n)

	case nodes.CreateEnumStmt:
		return convertCreateEnumStmt(&n)

	case nodes.CreateEventTrigStmt:
		return convertCreateEventTrigStmt(&n)

	case nodes.CreateExtensionStmt:
		return convertCreateExtensionStmt(&n)

	case nodes.CreateFdwStmt:
		return convertCreateFdwStmt(&n)

	case nodes.CreateForeignServerStmt:
		return convertCreateForeignServerStmt(&n)

	case nodes.CreateForeignTableStmt:
		return convertCreateForeignTableStmt(&n)

	case nodes.CreateFunctionStmt:
		return convertCreateFunctionStmt(&n)

	case nodes.CreateOpClassItem:
		return convertCreateOpClassItem(&n)

	case nodes.CreateOpClassStmt:
		return convertCreateOpClassStmt(&n)

	case nodes.CreateOpFamilyStmt:
		return convertCreateOpFamilyStmt(&n)

	case nodes.CreatePLangStmt:
		return convertCreatePLangStmt(&n)

	case nodes.CreatePolicyStmt:
		return convertCreatePolicyStmt(&n)

	case nodes.CreatePublicationStmt:
		return convertCreatePublicationStmt(&n)

	case nodes.CreateRangeStmt:
		return convertCreateRangeStmt(&n)

	case nodes.CreateRoleStmt:
		return convertCreateRoleStmt(&n)

	case nodes.CreateSchemaStmt:
		return convertCreateSchemaStmt(&n)

	case nodes.CreateSeqStmt:
		return convertCreateSeqStmt(&n)

	case nodes.CreateStatsStmt:
		return convertCreateStatsStmt(&n)

	case nodes.CreateStmt:
		return convertCreateStmt(&n)

	case nodes.CreateSubscriptionStmt:
		return convertCreateSubscriptionStmt(&n)

	case nodes.CreateTableAsStmt:
		return convertCreateTableAsStmt(&n)

	case nodes.CreateTableSpaceStmt:
		return convertCreateTableSpaceStmt(&n)

	case nodes.CreateTransformStmt:
		return convertCreateTransformStmt(&n)

	case nodes.CreateTrigStmt:
		return convertCreateTrigStmt(&n)

	case nodes.CreateUserMappingStmt:
		return convertCreateUserMappingStmt(&n)

	case nodes.CreatedbStmt:
		return convertCreatedbStmt(&n)

	case nodes.CurrentOfExpr:
		return convertCurrentOfExpr(&n)

	case nodes.DeallocateStmt:
		return convertDeallocateStmt(&n)

	case nodes.DeclareCursorStmt:
		return convertDeclareCursorStmt(&n)

	case nodes.DefElem:
		return convertDefElem(&n)

	case nodes.DefineStmt:
		return convertDefineStmt(&n)

	case nodes.DeleteStmt:
		return convertDeleteStmt(&n)

	case nodes.DiscardStmt:
		return convertDiscardStmt(&n)

	case nodes.DoStmt:
		return convertDoStmt(&n)

	case nodes.DropOwnedStmt:
		return convertDropOwnedStmt(&n)

	case nodes.DropRoleStmt:
		return convertDropRoleStmt(&n)

	case nodes.DropStmt:
		return convertDropStmt(&n)

	case nodes.DropSubscriptionStmt:
		return convertDropSubscriptionStmt(&n)

	case nodes.DropTableSpaceStmt:
		return convertDropTableSpaceStmt(&n)

	case nodes.DropUserMappingStmt:
		return convertDropUserMappingStmt(&n)

	case nodes.DropdbStmt:
		return convertDropdbStmt(&n)

	case nodes.ExecuteStmt:
		return convertExecuteStmt(&n)

	case nodes.ExplainStmt:
		return convertExplainStmt(&n)

	case nodes.Expr:
		return convertExpr(&n)

	case nodes.FetchStmt:
		return convertFetchStmt(&n)

	case nodes.FieldSelect:
		return convertFieldSelect(&n)

	case nodes.FieldStore:
		return convertFieldStore(&n)

	case nodes.Float:
		return convertFloat(&n)

	case nodes.FromExpr:
		return convertFromExpr(&n)

	case nodes.FuncCall:
		return convertFuncCall(&n)

	case nodes.FuncExpr:
		return convertFuncExpr(&n)

	case nodes.FunctionParameter:
		return convertFunctionParameter(&n)

	case nodes.GrantRoleStmt:
		return convertGrantRoleStmt(&n)

	case nodes.GrantStmt:
		return convertGrantStmt(&n)

	case nodes.GroupingFunc:
		return convertGroupingFunc(&n)

	case nodes.GroupingSet:
		return convertGroupingSet(&n)

	case nodes.ImportForeignSchemaStmt:
		return convertImportForeignSchemaStmt(&n)

	case nodes.IndexElem:
		return convertIndexElem(&n)

	case nodes.IndexStmt:
		return convertIndexStmt(&n)

	case nodes.InferClause:
		return convertInferClause(&n)

	case nodes.InferenceElem:
		return convertInferenceElem(&n)

	case nodes.InlineCodeBlock:
		return convertInlineCodeBlock(&n)

	case nodes.InsertStmt:
		return convertInsertStmt(&n)

	case nodes.Integer:
		return convertInteger(&n)

	case nodes.IntoClause:
		return convertIntoClause(&n)

	case nodes.JoinExpr:
		return convertJoinExpr(&n)

	case nodes.List:
		return convertList(n)

	case nodes.ListenStmt:
		return convertListenStmt(&n)

	case nodes.LoadStmt:
		return convertLoadStmt(&n)

	case nodes.LockStmt:
		return convertLockStmt(&n)

	case nodes.LockingClause:
		return convertLockingClause(&n)

	case nodes.MinMaxExpr:
		return convertMinMaxExpr(&n)

	case nodes.MultiAssignRef:
		return convertMultiAssignRef(&n)

	case nodes.NamedArgExpr:
		return convertNamedArgExpr(&n)

	case nodes.NextValueExpr:
		return convertNextValueExpr(&n)

	case nodes.NotifyStmt:
		return convertNotifyStmt(&n)

	case nodes.Null:
		return convertNull(&n)

	case nodes.NullTest:
		return convertNullTest(&n)

	case nodes.ObjectWithArgs:
		return convertObjectWithArgs(&n)

	case nodes.OnConflictClause:
		return convertOnConflictClause(&n)

	case nodes.OnConflictExpr:
		return convertOnConflictExpr(&n)

	case nodes.OpExpr:
		return convertOpExpr(&n)

	case nodes.Param:
		return convertParam(&n)

	case nodes.ParamExecData:
		return convertParamExecData(&n)

	case nodes.ParamExternData:
		return convertParamExternData(&n)

	case nodes.ParamListInfoData:
		return convertParamListInfoData(&n)

	case nodes.ParamRef:
		return convertParamRef(&n)

	case nodes.PartitionBoundSpec:
		return convertPartitionBoundSpec(&n)

	case nodes.PartitionCmd:
		return convertPartitionCmd(&n)

	case nodes.PartitionElem:
		return convertPartitionElem(&n)

	case nodes.PartitionRangeDatum:
		return convertPartitionRangeDatum(&n)

	case nodes.PartitionSpec:
		return convertPartitionSpec(&n)

	case nodes.PrepareStmt:
		return convertPrepareStmt(&n)

	case nodes.Query:
		return convertQuery(&n)

	case nodes.RangeFunction:
		return convertRangeFunction(&n)

	case nodes.RangeSubselect:
		return convertRangeSubselect(&n)

	case nodes.RangeTableFunc:
		return convertRangeTableFunc(&n)

	case nodes.RangeTableFuncCol:
		return convertRangeTableFuncCol(&n)

	case nodes.RangeTableSample:
		return convertRangeTableSample(&n)

	case nodes.RangeTblEntry:
		return convertRangeTblEntry(&n)

	case nodes.RangeTblFunction:
		return convertRangeTblFunction(&n)

	case nodes.RangeTblRef:
		return convertRangeTblRef(&n)

	case nodes.RangeVar:
		return convertRangeVar(&n)

	case nodes.RawStmt:
		return convertRawStmt(&n)

	case nodes.ReassignOwnedStmt:
		return convertReassignOwnedStmt(&n)

	case nodes.RefreshMatViewStmt:
		return convertRefreshMatViewStmt(&n)

	case nodes.ReindexStmt:
		return convertReindexStmt(&n)

	case nodes.RelabelType:
		return convertRelabelType(&n)

	case nodes.RenameStmt:
		return convertRenameStmt(&n)

	case nodes.ReplicaIdentityStmt:
		return convertReplicaIdentityStmt(&n)

	case nodes.ResTarget:
		return convertResTarget(&n)

	case nodes.RoleSpec:
		return convertRoleSpec(&n)

	case nodes.RowCompareExpr:
		return convertRowCompareExpr(&n)

	case nodes.RowExpr:
		return convertRowExpr(&n)

	case nodes.RowMarkClause:
		return convertRowMarkClause(&n)

	case nodes.RuleStmt:
		return convertRuleStmt(&n)

	case nodes.SQLValueFunction:
		return convertSQLValueFunction(&n)

	case nodes.ScalarArrayOpExpr:
		return convertScalarArrayOpExpr(&n)

	case nodes.SecLabelStmt:
		return convertSecLabelStmt(&n)

	case nodes.SelectStmt:
		return convertSelectStmt(&n)

	case nodes.SetOperationStmt:
		return convertSetOperationStmt(&n)

	case nodes.SetToDefault:
		return convertSetToDefault(&n)

	case nodes.SortBy:
		return convertSortBy(&n)

	case nodes.SortGroupClause:
		return convertSortGroupClause(&n)

	case nodes.String:
		return convertString(&n)

	case nodes.SubLink:
		return convertSubLink(&n)

	case nodes.SubPlan:
		return convertSubPlan(&n)

	case nodes.TableFunc:
		return convertTableFunc(&n)

	case nodes.TableLikeClause:
		return convertTableLikeClause(&n)

	case nodes.TableSampleClause:
		return convertTableSampleClause(&n)

	case nodes.TargetEntry:
		return convertTargetEntry(&n)

	case nodes.TransactionStmt:
		return convertTransactionStmt(&n)

	case nodes.TriggerTransition:
		return convertTriggerTransition(&n)

	case nodes.TruncateStmt:
		return convertTruncateStmt(&n)

	case nodes.TypeCast:
		return convertTypeCast(&n)

	case nodes.TypeName:
		return convertTypeName(&n)

	case nodes.UnlistenStmt:
		return convertUnlistenStmt(&n)

	case nodes.UpdateStmt:
		return convertUpdateStmt(&n)

	case nodes.VacuumStmt:
		return convertVacuumStmt(&n)

	case nodes.Var:
		return convertVar(&n)

	case nodes.VariableSetStmt:
		return convertVariableSetStmt(&n)

	case nodes.VariableShowStmt:
		return convertVariableShowStmt(&n)

	case nodes.ViewStmt:
		return convertViewStmt(&n)

	case nodes.WindowClause:
		return convertWindowClause(&n)

	case nodes.WindowDef:
		return convertWindowDef(&n)

	case nodes.WindowFunc:
		return convertWindowFunc(&n)

	case nodes.WithCheckOption:
		return convertWithCheckOption(&n)

	case nodes.WithClause:
		return convertWithClause(&n)

	case nodes.XmlExpr:
		return convertXmlExpr(&n)

	case nodes.XmlSerialize:
		return convertXmlSerialize(&n)

	default:
		return &ast.TODO{}
	}
}
