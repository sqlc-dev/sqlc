package postgresql

import (
	nodes "github.com/lfittl/pg_query_go/nodes"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
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

func convertA_ArrayExpr(n *nodes.A_ArrayExpr) *ast.A_ArrayExpr {
	if n == nil {
		return nil
	}
	return &ast.A_ArrayExpr{
		Elements: convertList(n.Elements),
		Location: n.Location,
	}
}

func convertA_Const(n *nodes.A_Const) *ast.A_Const {
	if n == nil {
		return nil
	}
	return &ast.A_Const{
		Val:      convertNode(n.Val),
		Location: n.Location,
	}
}

func convertA_Expr(n *nodes.A_Expr) *ast.A_Expr {
	if n == nil {
		return nil
	}
	return &ast.A_Expr{
		Kind:     ast.A_Expr_Kind(n.Kind),
		Name:     convertList(n.Name),
		Lexpr:    convertNode(n.Lexpr),
		Rexpr:    convertNode(n.Rexpr),
		Location: n.Location,
	}
}

func convertA_Indices(n *nodes.A_Indices) *ast.A_Indices {
	if n == nil {
		return nil
	}
	return &ast.A_Indices{
		IsSlice: n.IsSlice,
		Lidx:    convertNode(n.Lidx),
		Uidx:    convertNode(n.Uidx),
	}
}

func convertA_Indirection(n *nodes.A_Indirection) *ast.A_Indirection {
	if n == nil {
		return nil
	}
	return &ast.A_Indirection{
		Arg:         convertNode(n.Arg),
		Indirection: convertList(n.Indirection),
	}
}

func convertA_Star(n *nodes.A_Star) *ast.A_Star {
	if n == nil {
		return nil
	}
	return &ast.A_Star{}
}

func convertAccessPriv(n *nodes.AccessPriv) *ast.AccessPriv {
	if n == nil {
		return nil
	}
	return &ast.AccessPriv{
		PrivName: n.PrivName,
		Cols:     convertList(n.Cols),
	}
}

func convertAggref(n *nodes.Aggref) *ast.Aggref {
	if n == nil {
		return nil
	}
	return &ast.Aggref{
		Xpr:           convertNode(n.Xpr),
		Aggfnoid:      ast.Oid(n.Aggfnoid),
		Aggtype:       ast.Oid(n.Aggtype),
		Aggcollid:     ast.Oid(n.Aggcollid),
		Inputcollid:   ast.Oid(n.Inputcollid),
		Aggtranstype:  ast.Oid(n.Aggtranstype),
		Aggargtypes:   convertList(n.Aggargtypes),
		Aggdirectargs: convertList(n.Aggdirectargs),
		Args:          convertList(n.Args),
		Aggorder:      convertList(n.Aggorder),
		Aggdistinct:   convertList(n.Aggdistinct),
		Aggfilter:     convertNode(n.Aggfilter),
		Aggstar:       n.Aggstar,
		Aggvariadic:   n.Aggvariadic,
		Aggkind:       n.Aggkind,
		Agglevelsup:   ast.Index(n.Agglevelsup),
		Aggsplit:      ast.AggSplit(n.Aggsplit),
		Location:      n.Location,
	}
}

func convertAlias(n *nodes.Alias) *ast.Alias {
	if n == nil {
		return nil
	}
	return &ast.Alias{
		Aliasname: n.Aliasname,
		Colnames:  convertList(n.Colnames),
	}
}

func convertAlterCollationStmt(n *nodes.AlterCollationStmt) *ast.AlterCollationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterCollationStmt{
		Collname: convertList(n.Collname),
	}
}

func convertAlterDatabaseSetStmt(n *nodes.AlterDatabaseSetStmt) *ast.AlterDatabaseSetStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDatabaseSetStmt{
		Dbname:  n.Dbname,
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterDatabaseStmt(n *nodes.AlterDatabaseStmt) *ast.AlterDatabaseStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDatabaseStmt{
		Dbname:  n.Dbname,
		Options: convertList(n.Options),
	}
}

func convertAlterDefaultPrivilegesStmt(n *nodes.AlterDefaultPrivilegesStmt) *ast.AlterDefaultPrivilegesStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDefaultPrivilegesStmt{
		Options: convertList(n.Options),
		Action:  convertGrantStmt(n.Action),
	}
}

func convertAlterDomainStmt(n *nodes.AlterDomainStmt) *ast.AlterDomainStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDomainStmt{
		Subtype:   n.Subtype,
		TypeName:  convertList(n.TypeName),
		Name:      n.Name,
		Def:       convertNode(n.Def),
		Behavior:  ast.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterEnumStmt(n *nodes.AlterEnumStmt) *ast.AlterEnumStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterEnumStmt{
		TypeName:           convertList(n.TypeName),
		OldVal:             n.OldVal,
		NewVal:             n.NewVal,
		NewValNeighbor:     n.NewValNeighbor,
		NewValIsAfter:      n.NewValIsAfter,
		SkipIfNewValExists: n.SkipIfNewValExists,
	}
}

func convertAlterEventTrigStmt(n *nodes.AlterEventTrigStmt) *ast.AlterEventTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterEventTrigStmt{
		Trigname:  n.Trigname,
		Tgenabled: n.Tgenabled,
	}
}

func convertAlterExtensionContentsStmt(n *nodes.AlterExtensionContentsStmt) *ast.AlterExtensionContentsStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterExtensionContentsStmt{
		Extname: n.Extname,
		Action:  n.Action,
		Objtype: ast.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
	}
}

func convertAlterExtensionStmt(n *nodes.AlterExtensionStmt) *ast.AlterExtensionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterExtensionStmt{
		Extname: n.Extname,
		Options: convertList(n.Options),
	}
}

func convertAlterFdwStmt(n *nodes.AlterFdwStmt) *ast.AlterFdwStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterFdwStmt{
		Fdwname:     n.Fdwname,
		FuncOptions: convertList(n.FuncOptions),
		Options:     convertList(n.Options),
	}
}

func convertAlterForeignServerStmt(n *nodes.AlterForeignServerStmt) *ast.AlterForeignServerStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterForeignServerStmt{
		Servername: n.Servername,
		Version:    n.Version,
		Options:    convertList(n.Options),
		HasVersion: n.HasVersion,
	}
}

func convertAlterFunctionStmt(n *nodes.AlterFunctionStmt) *ast.AlterFunctionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterFunctionStmt{
		Func:    convertObjectWithArgs(n.Func),
		Actions: convertList(n.Actions),
	}
}

func convertAlterObjectDependsStmt(n *nodes.AlterObjectDependsStmt) *ast.AlterObjectDependsStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterObjectDependsStmt{
		ObjectType: ast.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Extname:    convertNode(n.Extname),
	}
}

func convertAlterObjectSchemaStmt(n *nodes.AlterObjectSchemaStmt) *ast.AlterObjectSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterObjectSchemaStmt{
		ObjectType: ast.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Newschema:  n.Newschema,
		MissingOk:  n.MissingOk,
	}
}

func convertAlterOpFamilyStmt(n *nodes.AlterOpFamilyStmt) *ast.AlterOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterOpFamilyStmt{
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
		IsDrop:       n.IsDrop,
		Items:        convertList(n.Items),
	}
}

func convertAlterOperatorStmt(n *nodes.AlterOperatorStmt) *ast.AlterOperatorStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterOperatorStmt{
		Opername: convertObjectWithArgs(n.Opername),
		Options:  convertList(n.Options),
	}
}

func convertAlterOwnerStmt(n *nodes.AlterOwnerStmt) *ast.AlterOwnerStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterOwnerStmt{
		ObjectType: ast.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Newowner:   convertRoleSpec(n.Newowner),
	}
}

func convertAlterPolicyStmt(n *nodes.AlterPolicyStmt) *ast.AlterPolicyStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterPolicyStmt{
		PolicyName: n.PolicyName,
		Table:      convertRangeVar(n.Table),
		Roles:      convertList(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertAlterPublicationStmt(n *nodes.AlterPublicationStmt) *ast.AlterPublicationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterPublicationStmt{
		Pubname:      n.Pubname,
		Options:      convertList(n.Options),
		Tables:       convertList(n.Tables),
		ForAllTables: n.ForAllTables,
		TableAction:  ast.DefElemAction(n.TableAction),
	}
}

func convertAlterRoleSetStmt(n *nodes.AlterRoleSetStmt) *ast.AlterRoleSetStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterRoleSetStmt{
		Role:     convertRoleSpec(n.Role),
		Database: n.Database,
		Setstmt:  convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterRoleStmt(n *nodes.AlterRoleStmt) *ast.AlterRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterRoleStmt{
		Role:    convertRoleSpec(n.Role),
		Options: convertList(n.Options),
		Action:  n.Action,
	}
}

func convertAlterSeqStmt(n *nodes.AlterSeqStmt) *ast.AlterSeqStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertList(n.Options),
		ForIdentity: n.ForIdentity,
		MissingOk:   n.MissingOk,
	}
}

func convertAlterSubscriptionStmt(n *nodes.AlterSubscriptionStmt) *ast.AlterSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSubscriptionStmt{
		Kind:        ast.AlterSubscriptionType(n.Kind),
		Subname:     n.Subname,
		Conninfo:    n.Conninfo,
		Publication: convertList(n.Publication),
		Options:     convertList(n.Options),
	}
}

func convertAlterSystemStmt(n *nodes.AlterSystemStmt) *ast.AlterSystemStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSystemStmt{
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterTSConfigurationStmt(n *nodes.AlterTSConfigurationStmt) *ast.AlterTSConfigurationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTSConfigurationStmt{
		Kind:      ast.AlterTSConfigType(n.Kind),
		Cfgname:   convertList(n.Cfgname),
		Tokentype: convertList(n.Tokentype),
		Dicts:     convertList(n.Dicts),
		Override:  n.Override,
		Replace:   n.Replace,
		MissingOk: n.MissingOk,
	}
}

func convertAlterTSDictionaryStmt(n *nodes.AlterTSDictionaryStmt) *ast.AlterTSDictionaryStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTSDictionaryStmt{
		Dictname: convertList(n.Dictname),
		Options:  convertList(n.Options),
	}
}

func convertAlterTableCmd(n *nodes.AlterTableCmd) *ast.AlterTableCmd {
	if n == nil {
		return nil
	}
	def := convertNode(n.Def)
	columnDef := def.(*ast.ColumnDef)
	return &ast.AlterTableCmd{
		Subtype:   ast.AlterTableType(n.Subtype),
		Name:      n.Name,
		Newowner:  convertRoleSpec(n.Newowner),
		Def:       columnDef,
		Behavior:  ast.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterTableMoveAllStmt(n *nodes.AlterTableMoveAllStmt) *ast.AlterTableMoveAllStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableMoveAllStmt{
		OrigTablespacename: n.OrigTablespacename,
		Objtype:            ast.ObjectType(n.Objtype),
		Roles:              convertList(n.Roles),
		NewTablespacename:  n.NewTablespacename,
		Nowait:             n.Nowait,
	}
}

func convertAlterTableSpaceOptionsStmt(n *nodes.AlterTableSpaceOptionsStmt) *ast.AlterTableSpaceOptionsStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableSpaceOptionsStmt{
		Tablespacename: n.Tablespacename,
		Options:        convertList(n.Options),
		IsReset:        n.IsReset,
	}
}

func convertAlterTableStmt(n *nodes.AlterTableStmt) *ast.AlterTableStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableStmt{
		Relation:  convertRangeVar(n.Relation),
		Cmds:      convertList(n.Cmds),
		Relkind:   ast.ObjectType(n.Relkind),
		MissingOk: n.MissingOk,
	}
}

func convertAlterUserMappingStmt(n *nodes.AlterUserMappingStmt) *ast.AlterUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: n.Servername,
		Options:    convertList(n.Options),
	}
}

func convertAlternativeSubPlan(n *nodes.AlternativeSubPlan) *ast.AlternativeSubPlan {
	if n == nil {
		return nil
	}
	return &ast.AlternativeSubPlan{
		Xpr:      convertNode(n.Xpr),
		Subplans: convertList(n.Subplans),
	}
}

func convertArrayCoerceExpr(n *nodes.ArrayCoerceExpr) *ast.ArrayCoerceExpr {
	if n == nil {
		return nil
	}
	return &ast.ArrayCoerceExpr{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Elemfuncid:   ast.Oid(n.Elemfuncid),
		Resulttype:   ast.Oid(n.Resulttype),
		Resulttypmod: n.Resulttypmod,
		Resultcollid: ast.Oid(n.Resultcollid),
		IsExplicit:   n.IsExplicit,
		Coerceformat: ast.CoercionForm(n.Coerceformat),
		Location:     n.Location,
	}
}

func convertArrayExpr(n *nodes.ArrayExpr) *ast.ArrayExpr {
	if n == nil {
		return nil
	}
	return &ast.ArrayExpr{
		Xpr:           convertNode(n.Xpr),
		ArrayTypeid:   ast.Oid(n.ArrayTypeid),
		ArrayCollid:   ast.Oid(n.ArrayCollid),
		ElementTypeid: ast.Oid(n.ElementTypeid),
		Elements:      convertList(n.Elements),
		Multidims:     n.Multidims,
		Location:      n.Location,
	}
}

func convertArrayRef(n *nodes.ArrayRef) *ast.ArrayRef {
	if n == nil {
		return nil
	}
	return &ast.ArrayRef{
		Xpr:             convertNode(n.Xpr),
		Refarraytype:    ast.Oid(n.Refarraytype),
		Refelemtype:     ast.Oid(n.Refelemtype),
		Reftypmod:       n.Reftypmod,
		Refcollid:       ast.Oid(n.Refcollid),
		Refupperindexpr: convertList(n.Refupperindexpr),
		Reflowerindexpr: convertList(n.Reflowerindexpr),
		Refexpr:         convertNode(n.Refexpr),
		Refassgnexpr:    convertNode(n.Refassgnexpr),
	}
}

func convertBitString(n *nodes.BitString) *ast.BitString {
	if n == nil {
		return nil
	}
	return &ast.BitString{
		Str: n.Str,
	}
}

func convertBlockIdData(n *nodes.BlockIdData) *ast.BlockIdData {
	if n == nil {
		return nil
	}
	return &ast.BlockIdData{
		BiHi: n.BiHi,
		BiLo: n.BiLo,
	}
}

func convertBoolExpr(n *nodes.BoolExpr) *ast.BoolExpr {
	if n == nil {
		return nil
	}
	return &ast.BoolExpr{
		Xpr:      convertNode(n.Xpr),
		Boolop:   ast.BoolExprType(n.Boolop),
		Args:     convertList(n.Args),
		Location: n.Location,
	}
}

func convertBooleanTest(n *nodes.BooleanTest) *ast.BooleanTest {
	if n == nil {
		return nil
	}
	return &ast.BooleanTest{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Booltesttype: ast.BoolTestType(n.Booltesttype),
		Location:     n.Location,
	}
}

func convertCaseExpr(n *nodes.CaseExpr) *ast.CaseExpr {
	if n == nil {
		return nil
	}
	return &ast.CaseExpr{
		Xpr:        convertNode(n.Xpr),
		Casetype:   ast.Oid(n.Casetype),
		Casecollid: ast.Oid(n.Casecollid),
		Arg:        convertNode(n.Arg),
		Args:       convertList(n.Args),
		Defresult:  convertNode(n.Defresult),
		Location:   n.Location,
	}
}

func convertCaseTestExpr(n *nodes.CaseTestExpr) *ast.CaseTestExpr {
	if n == nil {
		return nil
	}
	return &ast.CaseTestExpr{
		Xpr:       convertNode(n.Xpr),
		TypeId:    ast.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: ast.Oid(n.Collation),
	}
}

func convertCaseWhen(n *nodes.CaseWhen) *ast.CaseWhen {
	if n == nil {
		return nil
	}
	return &ast.CaseWhen{
		Xpr:      convertNode(n.Xpr),
		Expr:     convertNode(n.Expr),
		Result:   convertNode(n.Result),
		Location: n.Location,
	}
}

func convertCheckPointStmt(n *nodes.CheckPointStmt) *ast.CheckPointStmt {
	if n == nil {
		return nil
	}
	return &ast.CheckPointStmt{}
}

func convertClosePortalStmt(n *nodes.ClosePortalStmt) *ast.ClosePortalStmt {
	if n == nil {
		return nil
	}
	return &ast.ClosePortalStmt{
		Portalname: n.Portalname,
	}
}

func convertClusterStmt(n *nodes.ClusterStmt) *ast.ClusterStmt {
	if n == nil {
		return nil
	}
	return &ast.ClusterStmt{
		Relation:  convertRangeVar(n.Relation),
		Indexname: n.Indexname,
		Verbose:   n.Verbose,
	}
}

func convertCoalesceExpr(n *nodes.CoalesceExpr) *ast.CoalesceExpr {
	if n == nil {
		return nil
	}
	return &ast.CoalesceExpr{
		Xpr:            convertNode(n.Xpr),
		Coalescetype:   ast.Oid(n.Coalescetype),
		Coalescecollid: ast.Oid(n.Coalescecollid),
		Args:           convertList(n.Args),
		Location:       n.Location,
	}
}

func convertCoerceToDomain(n *nodes.CoerceToDomain) *ast.CoerceToDomain {
	if n == nil {
		return nil
	}
	return &ast.CoerceToDomain{
		Xpr:            convertNode(n.Xpr),
		Arg:            convertNode(n.Arg),
		Resulttype:     ast.Oid(n.Resulttype),
		Resulttypmod:   n.Resulttypmod,
		Resultcollid:   ast.Oid(n.Resultcollid),
		Coercionformat: ast.CoercionForm(n.Coercionformat),
		Location:       n.Location,
	}
}

func convertCoerceToDomainValue(n *nodes.CoerceToDomainValue) *ast.CoerceToDomainValue {
	if n == nil {
		return nil
	}
	return &ast.CoerceToDomainValue{
		Xpr:       convertNode(n.Xpr),
		TypeId:    ast.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: ast.Oid(n.Collation),
		Location:  n.Location,
	}
}

func convertCoerceViaIO(n *nodes.CoerceViaIO) *ast.CoerceViaIO {
	if n == nil {
		return nil
	}
	return &ast.CoerceViaIO{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Resulttype:   ast.Oid(n.Resulttype),
		Resultcollid: ast.Oid(n.Resultcollid),
		Coerceformat: ast.CoercionForm(n.Coerceformat),
		Location:     n.Location,
	}
}

func convertCollateClause(n *nodes.CollateClause) *ast.CollateClause {
	if n == nil {
		return nil
	}
	return &ast.CollateClause{
		Arg:      convertNode(n.Arg),
		Collname: convertList(n.Collname),
		Location: n.Location,
	}
}

func convertCollateExpr(n *nodes.CollateExpr) *ast.CollateExpr {
	if n == nil {
		return nil
	}
	return &ast.CollateExpr{
		Xpr:      convertNode(n.Xpr),
		Arg:      convertNode(n.Arg),
		CollOid:  ast.Oid(n.CollOid),
		Location: n.Location,
	}
}

func convertColumnDef(n *nodes.ColumnDef) *ast.ColumnDef {
	if n == nil {
		return nil
	}
	colname := ""
	if n.Colname != nil {
		colname = *n.Colname
	}
	return &ast.ColumnDef{
		Colname:       colname,
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
		CollOid:       ast.Oid(n.CollOid),
		Constraints:   convertList(n.Constraints),
		Fdwoptions:    convertList(n.Fdwoptions),
		Location:      n.Location,
	}
}

func convertColumnRef(n *nodes.ColumnRef) *ast.ColumnRef {
	if n == nil {
		return nil
	}
	return &ast.ColumnRef{
		Fields:   convertList(n.Fields),
		Location: n.Location,
	}
}

func convertCommentStmt(n *nodes.CommentStmt) *ast.CommentStmt {
	if n == nil {
		return nil
	}
	return &ast.CommentStmt{
		Objtype: ast.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
		Comment: n.Comment,
	}
}

func convertCommonTableExpr(n *nodes.CommonTableExpr) *ast.CommonTableExpr {
	if n == nil {
		return nil
	}
	return &ast.CommonTableExpr{
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

func convertCompositeTypeStmt(n *nodes.CompositeTypeStmt) *ast.CompositeTypeStmt {
	if n == nil {
		return nil
	}
	tn, err := parseTypeName(n.Typevar)
	if err != nil {
		panic(err)
	}
	return &ast.CompositeTypeStmt{
		TypeName: tn,
	}
}

func convertConst(n *nodes.Const) *ast.Const {
	if n == nil {
		return nil
	}
	return &ast.Const{
		Xpr:         convertNode(n.Xpr),
		Consttype:   ast.Oid(n.Consttype),
		Consttypmod: n.Consttypmod,
		Constcollid: ast.Oid(n.Constcollid),
		Constlen:    n.Constlen,
		Constvalue:  ast.Datum(n.Constvalue),
		Constisnull: n.Constisnull,
		Constbyval:  n.Constbyval,
		Location:    n.Location,
	}
}

func convertConstraint(n *nodes.Constraint) *ast.Constraint {
	if n == nil {
		return nil
	}
	return &ast.Constraint{
		Contype:        ast.ConstrType(n.Contype),
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
		OldPktableOid:  ast.Oid(n.OldPktableOid),
		SkipValidation: n.SkipValidation,
		InitiallyValid: n.InitiallyValid,
	}
}

func convertConstraintsSetStmt(n *nodes.ConstraintsSetStmt) *ast.ConstraintsSetStmt {
	if n == nil {
		return nil
	}
	return &ast.ConstraintsSetStmt{
		Constraints: convertList(n.Constraints),
		Deferred:    n.Deferred,
	}
}

func convertConvertRowtypeExpr(n *nodes.ConvertRowtypeExpr) *ast.ConvertRowtypeExpr {
	if n == nil {
		return nil
	}
	return &ast.ConvertRowtypeExpr{
		Xpr:           convertNode(n.Xpr),
		Arg:           convertNode(n.Arg),
		Resulttype:    ast.Oid(n.Resulttype),
		Convertformat: ast.CoercionForm(n.Convertformat),
		Location:      n.Location,
	}
}

func convertCopyStmt(n *nodes.CopyStmt) *ast.CopyStmt {
	if n == nil {
		return nil
	}
	return &ast.CopyStmt{
		Relation:  convertRangeVar(n.Relation),
		Query:     convertNode(n.Query),
		Attlist:   convertList(n.Attlist),
		IsFrom:    n.IsFrom,
		IsProgram: n.IsProgram,
		Filename:  n.Filename,
		Options:   convertList(n.Options),
	}
}

func convertCreateAmStmt(n *nodes.CreateAmStmt) *ast.CreateAmStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateAmStmt{
		Amname:      n.Amname,
		HandlerName: convertList(n.HandlerName),
		Amtype:      n.Amtype,
	}
}

func convertCreateCastStmt(n *nodes.CreateCastStmt) *ast.CreateCastStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateCastStmt{
		Sourcetype: convertTypeName(n.Sourcetype),
		Targettype: convertTypeName(n.Targettype),
		Func:       convertObjectWithArgs(n.Func),
		Context:    ast.CoercionContext(n.Context),
		Inout:      n.Inout,
	}
}

func convertCreateConversionStmt(n *nodes.CreateConversionStmt) *ast.CreateConversionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateConversionStmt{
		ConversionName:  convertList(n.ConversionName),
		ForEncodingName: n.ForEncodingName,
		ToEncodingName:  n.ToEncodingName,
		FuncName:        convertList(n.FuncName),
		Def:             n.Def,
	}
}

func convertCreateDomainStmt(n *nodes.CreateDomainStmt) *ast.CreateDomainStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateDomainStmt{
		Domainname:  convertList(n.Domainname),
		TypeName:    convertTypeName(n.TypeName),
		CollClause:  convertCollateClause(n.CollClause),
		Constraints: convertList(n.Constraints),
	}
}

func convertCreateEnumStmt(n *nodes.CreateEnumStmt) *ast.CreateEnumStmt {
	if n == nil {
		return nil
	}
	tn, err := parseTypeName(n.TypeName)
	if err != nil {
		panic(err)
	}
	return &ast.CreateEnumStmt{
		TypeName: tn,
		Vals:     convertList(n.Vals),
	}
}

func convertCreateEventTrigStmt(n *nodes.CreateEventTrigStmt) *ast.CreateEventTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateEventTrigStmt{
		Trigname:   n.Trigname,
		Eventname:  n.Eventname,
		Whenclause: convertList(n.Whenclause),
		Funcname:   convertList(n.Funcname),
	}
}

func convertCreateExtensionStmt(n *nodes.CreateExtensionStmt) *ast.CreateExtensionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateExtensionStmt{
		Extname:     n.Extname,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreateFdwStmt(n *nodes.CreateFdwStmt) *ast.CreateFdwStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateFdwStmt{
		Fdwname:     n.Fdwname,
		FuncOptions: convertList(n.FuncOptions),
		Options:     convertList(n.Options),
	}
}

func convertCreateForeignServerStmt(n *nodes.CreateForeignServerStmt) *ast.CreateForeignServerStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateForeignServerStmt{
		Servername:  n.Servername,
		Servertype:  n.Servertype,
		Version:     n.Version,
		Fdwname:     n.Fdwname,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreateForeignTableStmt(n *nodes.CreateForeignTableStmt) *ast.CreateForeignTableStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateForeignTableStmt{
		Base:       convertCreateStmt(&n.Base),
		Servername: n.Servername,
		Options:    convertList(n.Options),
	}
}

func convertCreateFunctionStmt(n *nodes.CreateFunctionStmt) *ast.CreateFunctionStmt {
	if n == nil {
		return nil
	}
	fn, err := parseFuncName(n.Funcname)
	if err != nil {
		panic(err)
	}
	return &ast.CreateFunctionStmt{
		Replace:    n.Replace,
		Func:       fn,
		Params:     convertList(n.Parameters),
		ReturnType: convertTypeName(n.ReturnType),
		Options:    convertList(n.Options),
		WithClause: convertList(n.WithClause),
	}
}

func convertCreateOpClassItem(n *nodes.CreateOpClassItem) *ast.CreateOpClassItem {
	if n == nil {
		return nil
	}
	return &ast.CreateOpClassItem{
		Itemtype:    n.Itemtype,
		Name:        convertObjectWithArgs(n.Name),
		Number:      n.Number,
		OrderFamily: convertList(n.OrderFamily),
		ClassArgs:   convertList(n.ClassArgs),
		Storedtype:  convertTypeName(n.Storedtype),
	}
}

func convertCreateOpClassStmt(n *nodes.CreateOpClassStmt) *ast.CreateOpClassStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateOpClassStmt{
		Opclassname:  convertList(n.Opclassname),
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
		Datatype:     convertTypeName(n.Datatype),
		Items:        convertList(n.Items),
		IsDefault:    n.IsDefault,
	}
}

func convertCreateOpFamilyStmt(n *nodes.CreateOpFamilyStmt) *ast.CreateOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateOpFamilyStmt{
		Opfamilyname: convertList(n.Opfamilyname),
		Amname:       n.Amname,
	}
}

func convertCreatePLangStmt(n *nodes.CreatePLangStmt) *ast.CreatePLangStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePLangStmt{
		Replace:     n.Replace,
		Plname:      n.Plname,
		Plhandler:   convertList(n.Plhandler),
		Plinline:    convertList(n.Plinline),
		Plvalidator: convertList(n.Plvalidator),
		Pltrusted:   n.Pltrusted,
	}
}

func convertCreatePolicyStmt(n *nodes.CreatePolicyStmt) *ast.CreatePolicyStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePolicyStmt{
		PolicyName: n.PolicyName,
		Table:      convertRangeVar(n.Table),
		CmdName:    n.CmdName,
		Permissive: n.Permissive,
		Roles:      convertList(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertCreatePublicationStmt(n *nodes.CreatePublicationStmt) *ast.CreatePublicationStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePublicationStmt{
		Pubname:      n.Pubname,
		Options:      convertList(n.Options),
		Tables:       convertList(n.Tables),
		ForAllTables: n.ForAllTables,
	}
}

func convertCreateRangeStmt(n *nodes.CreateRangeStmt) *ast.CreateRangeStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateRangeStmt{
		TypeName: convertList(n.TypeName),
		Params:   convertList(n.Params),
	}
}

func convertCreateRoleStmt(n *nodes.CreateRoleStmt) *ast.CreateRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateRoleStmt{
		StmtType: ast.RoleStmtType(n.StmtType),
		Role:     n.Role,
		Options:  convertList(n.Options),
	}
}

func convertCreateSchemaStmt(n *nodes.CreateSchemaStmt) *ast.CreateSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSchemaStmt{
		Name:        n.Schemaname,
		Authrole:    convertRoleSpec(n.Authrole),
		SchemaElts:  convertList(n.SchemaElts),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateSeqStmt(n *nodes.CreateSeqStmt) *ast.CreateSeqStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertList(n.Options),
		OwnerId:     ast.Oid(n.OwnerId),
		ForIdentity: n.ForIdentity,
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStatsStmt(n *nodes.CreateStatsStmt) *ast.CreateStatsStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateStatsStmt{
		Defnames:    convertList(n.Defnames),
		StatTypes:   convertList(n.StatTypes),
		Exprs:       convertList(n.Exprs),
		Relations:   convertList(n.Relations),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStmt(n *nodes.CreateStmt) *ast.CreateStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateStmt{
		Relation:       convertRangeVar(n.Relation),
		TableElts:      convertList(n.TableElts),
		InhRelations:   convertList(n.InhRelations),
		Partbound:      convertPartitionBoundSpec(n.Partbound),
		Partspec:       convertPartitionSpec(n.Partspec),
		OfTypename:     convertTypeName(n.OfTypename),
		Constraints:    convertList(n.Constraints),
		Options:        convertList(n.Options),
		Oncommit:       ast.OnCommitAction(n.Oncommit),
		Tablespacename: n.Tablespacename,
		IfNotExists:    n.IfNotExists,
	}
}

func convertCreateSubscriptionStmt(n *nodes.CreateSubscriptionStmt) *ast.CreateSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSubscriptionStmt{
		Subname:     n.Subname,
		Conninfo:    n.Conninfo,
		Publication: convertList(n.Publication),
		Options:     convertList(n.Options),
	}
}

func convertCreateTableAsStmt(n *nodes.CreateTableAsStmt) *ast.CreateTableAsStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTableAsStmt{
		Query:        convertNode(n.Query),
		Into:         convertIntoClause(n.Into),
		Relkind:      ast.ObjectType(n.Relkind),
		IsSelectInto: n.IsSelectInto,
		IfNotExists:  n.IfNotExists,
	}
}

func convertCreateTableSpaceStmt(n *nodes.CreateTableSpaceStmt) *ast.CreateTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		Owner:          convertRoleSpec(n.Owner),
		Location:       n.Location,
		Options:        convertList(n.Options),
	}
}

func convertCreateTransformStmt(n *nodes.CreateTransformStmt) *ast.CreateTransformStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTransformStmt{
		Replace:  n.Replace,
		TypeName: convertTypeName(n.TypeName),
		Lang:     n.Lang,
		Fromsql:  convertObjectWithArgs(n.Fromsql),
		Tosql:    convertObjectWithArgs(n.Tosql),
	}
}

func convertCreateTrigStmt(n *nodes.CreateTrigStmt) *ast.CreateTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTrigStmt{
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

func convertCreateUserMappingStmt(n *nodes.CreateUserMappingStmt) *ast.CreateUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateUserMappingStmt{
		User:        convertRoleSpec(n.User),
		Servername:  n.Servername,
		IfNotExists: n.IfNotExists,
		Options:     convertList(n.Options),
	}
}

func convertCreatedbStmt(n *nodes.CreatedbStmt) *ast.CreatedbStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatedbStmt{
		Dbname:  n.Dbname,
		Options: convertList(n.Options),
	}
}

func convertCurrentOfExpr(n *nodes.CurrentOfExpr) *ast.CurrentOfExpr {
	if n == nil {
		return nil
	}
	return &ast.CurrentOfExpr{
		Xpr:         convertNode(n.Xpr),
		Cvarno:      ast.Index(n.Cvarno),
		CursorName:  n.CursorName,
		CursorParam: n.CursorParam,
	}
}

func convertDeallocateStmt(n *nodes.DeallocateStmt) *ast.DeallocateStmt {
	if n == nil {
		return nil
	}
	return &ast.DeallocateStmt{
		Name: n.Name,
	}
}

func convertDeclareCursorStmt(n *nodes.DeclareCursorStmt) *ast.DeclareCursorStmt {
	if n == nil {
		return nil
	}
	return &ast.DeclareCursorStmt{
		Portalname: n.Portalname,
		Options:    n.Options,
		Query:      convertNode(n.Query),
	}
}

func convertDefElem(n *nodes.DefElem) *ast.DefElem {
	if n == nil {
		return nil
	}
	return &ast.DefElem{
		Defnamespace: n.Defnamespace,
		Defname:      n.Defname,
		Arg:          convertNode(n.Arg),
		Defaction:    ast.DefElemAction(n.Defaction),
		Location:     n.Location,
	}
}

func convertDefineStmt(n *nodes.DefineStmt) *ast.DefineStmt {
	if n == nil {
		return nil
	}
	return &ast.DefineStmt{
		Kind:        ast.ObjectType(n.Kind),
		Oldstyle:    n.Oldstyle,
		Defnames:    convertList(n.Defnames),
		Args:        convertList(n.Args),
		Definition:  convertList(n.Definition),
		IfNotExists: n.IfNotExists,
	}
}

func convertDeleteStmt(n *nodes.DeleteStmt) *ast.DeleteStmt {
	if n == nil {
		return nil
	}
	return &ast.DeleteStmt{
		Relation:      convertRangeVar(n.Relation),
		UsingClause:   convertList(n.UsingClause),
		WhereClause:   convertNode(n.WhereClause),
		ReturningList: convertList(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertDiscardStmt(n *nodes.DiscardStmt) *ast.DiscardStmt {
	if n == nil {
		return nil
	}
	return &ast.DiscardStmt{
		Target: ast.DiscardMode(n.Target),
	}
}

func convertDoStmt(n *nodes.DoStmt) *ast.DoStmt {
	if n == nil {
		return nil
	}
	return &ast.DoStmt{
		Args: convertList(n.Args),
	}
}

func convertDropOwnedStmt(n *nodes.DropOwnedStmt) *ast.DropOwnedStmt {
	if n == nil {
		return nil
	}
	return &ast.DropOwnedStmt{
		Roles:    convertList(n.Roles),
		Behavior: ast.DropBehavior(n.Behavior),
	}
}

func convertDropRoleStmt(n *nodes.DropRoleStmt) *ast.DropRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.DropRoleStmt{
		Roles:     convertList(n.Roles),
		MissingOk: n.MissingOk,
	}
}

func convertDropStmt(n *nodes.DropStmt) *ast.DropStmt {
	if n == nil {
		return nil
	}
	return &ast.DropStmt{
		Objects:    convertList(n.Objects),
		RemoveType: ast.ObjectType(n.RemoveType),
		Behavior:   ast.DropBehavior(n.Behavior),
		MissingOk:  n.MissingOk,
		Concurrent: n.Concurrent,
	}
}

func convertDropSubscriptionStmt(n *nodes.DropSubscriptionStmt) *ast.DropSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.DropSubscriptionStmt{
		Subname:   n.Subname,
		MissingOk: n.MissingOk,
		Behavior:  ast.DropBehavior(n.Behavior),
	}
}

func convertDropTableSpaceStmt(n *nodes.DropTableSpaceStmt) *ast.DropTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &ast.DropTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		MissingOk:      n.MissingOk,
	}
}

func convertDropUserMappingStmt(n *nodes.DropUserMappingStmt) *ast.DropUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.DropUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: n.Servername,
		MissingOk:  n.MissingOk,
	}
}

func convertDropdbStmt(n *nodes.DropdbStmt) *ast.DropdbStmt {
	if n == nil {
		return nil
	}
	return &ast.DropdbStmt{
		Dbname:    n.Dbname,
		MissingOk: n.MissingOk,
	}
}

func convertExecuteStmt(n *nodes.ExecuteStmt) *ast.ExecuteStmt {
	if n == nil {
		return nil
	}
	return &ast.ExecuteStmt{
		Name:   n.Name,
		Params: convertList(n.Params),
	}
}

func convertExplainStmt(n *nodes.ExplainStmt) *ast.ExplainStmt {
	if n == nil {
		return nil
	}
	return &ast.ExplainStmt{
		Query:   convertNode(n.Query),
		Options: convertList(n.Options),
	}
}

func convertExpr(n *nodes.Expr) *ast.Expr {
	if n == nil {
		return nil
	}
	return &ast.Expr{}
}

func convertFetchStmt(n *nodes.FetchStmt) *ast.FetchStmt {
	if n == nil {
		return nil
	}
	return &ast.FetchStmt{
		Direction:  ast.FetchDirection(n.Direction),
		HowMany:    n.HowMany,
		Portalname: n.Portalname,
		Ismove:     n.Ismove,
	}
}

func convertFieldSelect(n *nodes.FieldSelect) *ast.FieldSelect {
	if n == nil {
		return nil
	}
	return &ast.FieldSelect{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Fieldnum:     ast.AttrNumber(n.Fieldnum),
		Resulttype:   ast.Oid(n.Resulttype),
		Resulttypmod: n.Resulttypmod,
		Resultcollid: ast.Oid(n.Resultcollid),
	}
}

func convertFieldStore(n *nodes.FieldStore) *ast.FieldStore {
	if n == nil {
		return nil
	}
	return &ast.FieldStore{
		Xpr:        convertNode(n.Xpr),
		Arg:        convertNode(n.Arg),
		Newvals:    convertList(n.Newvals),
		Fieldnums:  convertList(n.Fieldnums),
		Resulttype: ast.Oid(n.Resulttype),
	}
}

func convertFloat(n *nodes.Float) *ast.Float {
	if n == nil {
		return nil
	}
	return &ast.Float{
		Str: n.Str,
	}
}

func convertFromExpr(n *nodes.FromExpr) *ast.FromExpr {
	if n == nil {
		return nil
	}
	return &ast.FromExpr{
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

func convertFuncExpr(n *nodes.FuncExpr) *ast.FuncExpr {
	if n == nil {
		return nil
	}
	return &ast.FuncExpr{
		Xpr:            convertNode(n.Xpr),
		Funcid:         ast.Oid(n.Funcid),
		Funcresulttype: ast.Oid(n.Funcresulttype),
		Funcretset:     n.Funcretset,
		Funcvariadic:   n.Funcvariadic,
		Funcformat:     ast.CoercionForm(n.Funcformat),
		Funccollid:     ast.Oid(n.Funccollid),
		Inputcollid:    ast.Oid(n.Inputcollid),
		Args:           convertList(n.Args),
		Location:       n.Location,
	}
}

func convertFunctionParameter(n *nodes.FunctionParameter) *ast.FunctionParameter {
	if n == nil {
		return nil
	}
	return &ast.FunctionParameter{
		Name:    n.Name,
		ArgType: convertTypeName(n.ArgType),
		Mode:    ast.FunctionParameterMode(n.Mode),
		Defexpr: convertNode(n.Defexpr),
	}
}

func convertGrantRoleStmt(n *nodes.GrantRoleStmt) *ast.GrantRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.GrantRoleStmt{
		GrantedRoles: convertList(n.GrantedRoles),
		GranteeRoles: convertList(n.GranteeRoles),
		IsGrant:      n.IsGrant,
		AdminOpt:     n.AdminOpt,
		Grantor:      convertRoleSpec(n.Grantor),
		Behavior:     ast.DropBehavior(n.Behavior),
	}
}

func convertGrantStmt(n *nodes.GrantStmt) *ast.GrantStmt {
	if n == nil {
		return nil
	}
	return &ast.GrantStmt{
		IsGrant:     n.IsGrant,
		Targtype:    ast.GrantTargetType(n.Targtype),
		Objtype:     ast.GrantObjectType(n.Objtype),
		Objects:     convertList(n.Objects),
		Privileges:  convertList(n.Privileges),
		Grantees:    convertList(n.Grantees),
		GrantOption: n.GrantOption,
		Behavior:    ast.DropBehavior(n.Behavior),
	}
}

func convertGroupingFunc(n *nodes.GroupingFunc) *ast.GroupingFunc {
	if n == nil {
		return nil
	}
	return &ast.GroupingFunc{
		Xpr:         convertNode(n.Xpr),
		Args:        convertList(n.Args),
		Refs:        convertList(n.Refs),
		Cols:        convertList(n.Cols),
		Agglevelsup: ast.Index(n.Agglevelsup),
		Location:    n.Location,
	}
}

func convertGroupingSet(n *nodes.GroupingSet) *ast.GroupingSet {
	if n == nil {
		return nil
	}
	return &ast.GroupingSet{
		Kind:     ast.GroupingSetKind(n.Kind),
		Content:  convertList(n.Content),
		Location: n.Location,
	}
}

func convertImportForeignSchemaStmt(n *nodes.ImportForeignSchemaStmt) *ast.ImportForeignSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.ImportForeignSchemaStmt{
		ServerName:   n.ServerName,
		RemoteSchema: n.RemoteSchema,
		LocalSchema:  n.LocalSchema,
		ListType:     ast.ImportForeignSchemaType(n.ListType),
		TableList:    convertList(n.TableList),
		Options:      convertList(n.Options),
	}
}

func convertIndexElem(n *nodes.IndexElem) *ast.IndexElem {
	if n == nil {
		return nil
	}
	return &ast.IndexElem{
		Name:          n.Name,
		Expr:          convertNode(n.Expr),
		Indexcolname:  n.Indexcolname,
		Collation:     convertList(n.Collation),
		Opclass:       convertList(n.Opclass),
		Ordering:      ast.SortByDir(n.Ordering),
		NullsOrdering: ast.SortByNulls(n.NullsOrdering),
	}
}

func convertIndexStmt(n *nodes.IndexStmt) *ast.IndexStmt {
	if n == nil {
		return nil
	}
	return &ast.IndexStmt{
		Idxname:        n.Idxname,
		Relation:       convertRangeVar(n.Relation),
		AccessMethod:   n.AccessMethod,
		TableSpace:     n.TableSpace,
		IndexParams:    convertList(n.IndexParams),
		Options:        convertList(n.Options),
		WhereClause:    convertNode(n.WhereClause),
		ExcludeOpNames: convertList(n.ExcludeOpNames),
		Idxcomment:     n.Idxcomment,
		IndexOid:       ast.Oid(n.IndexOid),
		OldNode:        ast.Oid(n.OldNode),
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

func convertInferClause(n *nodes.InferClause) *ast.InferClause {
	if n == nil {
		return nil
	}
	return &ast.InferClause{
		IndexElems:  convertList(n.IndexElems),
		WhereClause: convertNode(n.WhereClause),
		Conname:     n.Conname,
		Location:    n.Location,
	}
}

func convertInferenceElem(n *nodes.InferenceElem) *ast.InferenceElem {
	if n == nil {
		return nil
	}
	return &ast.InferenceElem{
		Xpr:          convertNode(n.Xpr),
		Expr:         convertNode(n.Expr),
		Infercollid:  ast.Oid(n.Infercollid),
		Inferopclass: ast.Oid(n.Inferopclass),
	}
}

func convertInlineCodeBlock(n *nodes.InlineCodeBlock) *ast.InlineCodeBlock {
	if n == nil {
		return nil
	}
	return &ast.InlineCodeBlock{
		SourceText:    n.SourceText,
		LangOid:       ast.Oid(n.LangOid),
		LangIsTrusted: n.LangIsTrusted,
	}
}

func convertInsertStmt(n *nodes.InsertStmt) *ast.InsertStmt {
	if n == nil {
		return nil
	}
	return &ast.InsertStmt{
		Relation:         convertRangeVar(n.Relation),
		Cols:             convertList(n.Cols),
		SelectStmt:       convertNode(n.SelectStmt),
		OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		ReturningList:    convertList(n.ReturningList),
		WithClause:       convertWithClause(n.WithClause),
		Override:         ast.OverridingKind(n.Override),
	}
}

func convertInteger(n *nodes.Integer) *ast.Integer {
	if n == nil {
		return nil
	}
	return &ast.Integer{
		Ival: n.Ival,
	}
}

func convertIntoClause(n *nodes.IntoClause) *ast.IntoClause {
	if n == nil {
		return nil
	}
	return &ast.IntoClause{
		Rel:            convertRangeVar(n.Rel),
		ColNames:       convertList(n.ColNames),
		Options:        convertList(n.Options),
		OnCommit:       ast.OnCommitAction(n.OnCommit),
		TableSpaceName: n.TableSpaceName,
		ViewQuery:      convertNode(n.ViewQuery),
		SkipData:       n.SkipData,
	}
}

func convertJoinExpr(n *nodes.JoinExpr) *ast.JoinExpr {
	if n == nil {
		return nil
	}
	return &ast.JoinExpr{
		Jointype:    ast.JoinType(n.Jointype),
		IsNatural:   n.IsNatural,
		Larg:        convertNode(n.Larg),
		Rarg:        convertNode(n.Rarg),
		UsingClause: convertList(n.UsingClause),
		Quals:       convertNode(n.Quals),
		Alias:       convertAlias(n.Alias),
		Rtindex:     n.Rtindex,
	}
}

func convertListenStmt(n *nodes.ListenStmt) *ast.ListenStmt {
	if n == nil {
		return nil
	}
	return &ast.ListenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertLoadStmt(n *nodes.LoadStmt) *ast.LoadStmt {
	if n == nil {
		return nil
	}
	return &ast.LoadStmt{
		Filename: n.Filename,
	}
}

func convertLockStmt(n *nodes.LockStmt) *ast.LockStmt {
	if n == nil {
		return nil
	}
	return &ast.LockStmt{
		Relations: convertList(n.Relations),
		Mode:      n.Mode,
		Nowait:    n.Nowait,
	}
}

func convertLockingClause(n *nodes.LockingClause) *ast.LockingClause {
	if n == nil {
		return nil
	}
	return &ast.LockingClause{
		LockedRels: convertList(n.LockedRels),
		Strength:   ast.LockClauseStrength(n.Strength),
		WaitPolicy: ast.LockWaitPolicy(n.WaitPolicy),
	}
}

func convertMinMaxExpr(n *nodes.MinMaxExpr) *ast.MinMaxExpr {
	if n == nil {
		return nil
	}
	return &ast.MinMaxExpr{
		Xpr:          convertNode(n.Xpr),
		Minmaxtype:   ast.Oid(n.Minmaxtype),
		Minmaxcollid: ast.Oid(n.Minmaxcollid),
		Inputcollid:  ast.Oid(n.Inputcollid),
		Op:           ast.MinMaxOp(n.Op),
		Args:         convertList(n.Args),
		Location:     n.Location,
	}
}

func convertMultiAssignRef(n *nodes.MultiAssignRef) *ast.MultiAssignRef {
	if n == nil {
		return nil
	}
	return &ast.MultiAssignRef{
		Source:   convertNode(n.Source),
		Colno:    n.Colno,
		Ncolumns: n.Ncolumns,
	}
}

func convertNamedArgExpr(n *nodes.NamedArgExpr) *ast.NamedArgExpr {
	if n == nil {
		return nil
	}
	return &ast.NamedArgExpr{
		Xpr:       convertNode(n.Xpr),
		Arg:       convertNode(n.Arg),
		Name:      n.Name,
		Argnumber: n.Argnumber,
		Location:  n.Location,
	}
}

func convertNextValueExpr(n *nodes.NextValueExpr) *ast.NextValueExpr {
	if n == nil {
		return nil
	}
	return &ast.NextValueExpr{
		Xpr:    convertNode(n.Xpr),
		Seqid:  ast.Oid(n.Seqid),
		TypeId: ast.Oid(n.TypeId),
	}
}

func convertNotifyStmt(n *nodes.NotifyStmt) *ast.NotifyStmt {
	if n == nil {
		return nil
	}
	return &ast.NotifyStmt{
		Conditionname: n.Conditionname,
		Payload:       n.Payload,
	}
}

func convertNull(n *nodes.Null) *ast.Null {
	if n == nil {
		return nil
	}
	return &ast.Null{}
}

func convertNullTest(n *nodes.NullTest) *ast.NullTest {
	if n == nil {
		return nil
	}
	return &ast.NullTest{
		Xpr:          convertNode(n.Xpr),
		Arg:          convertNode(n.Arg),
		Nulltesttype: ast.NullTestType(n.Nulltesttype),
		Argisrow:     n.Argisrow,
		Location:     n.Location,
	}
}

func convertObjectWithArgs(n *nodes.ObjectWithArgs) *ast.ObjectWithArgs {
	if n == nil {
		return nil
	}
	return &ast.ObjectWithArgs{
		Objname:         convertList(n.Objname),
		Objargs:         convertList(n.Objargs),
		ArgsUnspecified: n.ArgsUnspecified,
	}
}

func convertOnConflictClause(n *nodes.OnConflictClause) *ast.OnConflictClause {
	if n == nil {
		return nil
	}
	return &ast.OnConflictClause{
		Action:      ast.OnConflictAction(n.Action),
		Infer:       convertInferClause(n.Infer),
		TargetList:  convertList(n.TargetList),
		WhereClause: convertNode(n.WhereClause),
		Location:    n.Location,
	}
}

func convertOnConflictExpr(n *nodes.OnConflictExpr) *ast.OnConflictExpr {
	if n == nil {
		return nil
	}
	return &ast.OnConflictExpr{
		Action:          ast.OnConflictAction(n.Action),
		ArbiterElems:    convertList(n.ArbiterElems),
		ArbiterWhere:    convertNode(n.ArbiterWhere),
		Constraint:      ast.Oid(n.Constraint),
		OnConflictSet:   convertList(n.OnConflictSet),
		OnConflictWhere: convertNode(n.OnConflictWhere),
		ExclRelIndex:    n.ExclRelIndex,
		ExclRelTlist:    convertList(n.ExclRelTlist),
	}
}

func convertOpExpr(n *nodes.OpExpr) *ast.OpExpr {
	if n == nil {
		return nil
	}
	return &ast.OpExpr{
		Xpr:          convertNode(n.Xpr),
		Opno:         ast.Oid(n.Opno),
		Opfuncid:     ast.Oid(n.Opfuncid),
		Opresulttype: ast.Oid(n.Opresulttype),
		Opretset:     n.Opretset,
		Opcollid:     ast.Oid(n.Opcollid),
		Inputcollid:  ast.Oid(n.Inputcollid),
		Args:         convertList(n.Args),
		Location:     n.Location,
	}
}

func convertParam(n *nodes.Param) *ast.Param {
	if n == nil {
		return nil
	}
	return &ast.Param{
		Xpr:         convertNode(n.Xpr),
		Paramkind:   ast.ParamKind(n.Paramkind),
		Paramid:     n.Paramid,
		Paramtype:   ast.Oid(n.Paramtype),
		Paramtypmod: n.Paramtypmod,
		Paramcollid: ast.Oid(n.Paramcollid),
		Location:    n.Location,
	}
}

func convertParamExecData(n *nodes.ParamExecData) *ast.ParamExecData {
	if n == nil {
		return nil
	}
	return &ast.ParamExecData{
		ExecPlan: &ast.TODO{},
		Value:    ast.Datum(n.Value),
		Isnull:   n.Isnull,
	}
}

func convertParamExternData(n *nodes.ParamExternData) *ast.ParamExternData {
	if n == nil {
		return nil
	}
	return &ast.ParamExternData{
		Value:  ast.Datum(n.Value),
		Isnull: n.Isnull,
		Pflags: n.Pflags,
		Ptype:  ast.Oid(n.Ptype),
	}
}

func convertParamListInfoData(n *nodes.ParamListInfoData) *ast.ParamListInfoData {
	if n == nil {
		return nil
	}
	return &ast.ParamListInfoData{
		ParamFetchArg:  &ast.TODO{},
		ParserSetupArg: &ast.TODO{},
		NumParams:      n.NumParams,
		ParamMask:      n.ParamMask,
	}
}

func convertParamRef(n *nodes.ParamRef) *ast.ParamRef {
	if n == nil {
		return nil
	}
	return &ast.ParamRef{
		Number:   n.Number,
		Location: n.Location,
	}
}

func convertPartitionBoundSpec(n *nodes.PartitionBoundSpec) *ast.PartitionBoundSpec {
	if n == nil {
		return nil
	}
	return &ast.PartitionBoundSpec{
		Strategy:    n.Strategy,
		Listdatums:  convertList(n.Listdatums),
		Lowerdatums: convertList(n.Lowerdatums),
		Upperdatums: convertList(n.Upperdatums),
		Location:    n.Location,
	}
}

func convertPartitionCmd(n *nodes.PartitionCmd) *ast.PartitionCmd {
	if n == nil {
		return nil
	}
	return &ast.PartitionCmd{
		Name:  convertRangeVar(n.Name),
		Bound: convertPartitionBoundSpec(n.Bound),
	}
}

func convertPartitionElem(n *nodes.PartitionElem) *ast.PartitionElem {
	if n == nil {
		return nil
	}
	return &ast.PartitionElem{
		Name:      n.Name,
		Expr:      convertNode(n.Expr),
		Collation: convertList(n.Collation),
		Opclass:   convertList(n.Opclass),
		Location:  n.Location,
	}
}

func convertPartitionRangeDatum(n *nodes.PartitionRangeDatum) *ast.PartitionRangeDatum {
	if n == nil {
		return nil
	}
	return &ast.PartitionRangeDatum{
		Kind:     ast.PartitionRangeDatumKind(n.Kind),
		Value:    convertNode(n.Value),
		Location: n.Location,
	}
}

func convertPartitionSpec(n *nodes.PartitionSpec) *ast.PartitionSpec {
	if n == nil {
		return nil
	}
	return &ast.PartitionSpec{
		Strategy:   n.Strategy,
		PartParams: convertList(n.PartParams),
		Location:   n.Location,
	}
}

func convertPrepareStmt(n *nodes.PrepareStmt) *ast.PrepareStmt {
	if n == nil {
		return nil
	}
	return &ast.PrepareStmt{
		Name:     n.Name,
		Argtypes: convertList(n.Argtypes),
		Query:    convertNode(n.Query),
	}
}

func convertQuery(n *nodes.Query) *ast.Query {
	if n == nil {
		return nil
	}
	return &ast.Query{
		CommandType:      ast.CmdType(n.CommandType),
		QuerySource:      ast.QuerySource(n.QuerySource),
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
		Override:         ast.OverridingKind(n.Override),
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

func convertRangeFunction(n *nodes.RangeFunction) *ast.RangeFunction {
	if n == nil {
		return nil
	}
	return &ast.RangeFunction{
		Lateral:    n.Lateral,
		Ordinality: n.Ordinality,
		IsRowsfrom: n.IsRowsfrom,
		Functions:  convertList(n.Functions),
		Alias:      convertAlias(n.Alias),
		Coldeflist: convertList(n.Coldeflist),
	}
}

func convertRangeSubselect(n *nodes.RangeSubselect) *ast.RangeSubselect {
	if n == nil {
		return nil
	}
	return &ast.RangeSubselect{
		Lateral:  n.Lateral,
		Subquery: convertNode(n.Subquery),
		Alias:    convertAlias(n.Alias),
	}
}

func convertRangeTableFunc(n *nodes.RangeTableFunc) *ast.RangeTableFunc {
	if n == nil {
		return nil
	}
	return &ast.RangeTableFunc{
		Lateral:    n.Lateral,
		Docexpr:    convertNode(n.Docexpr),
		Rowexpr:    convertNode(n.Rowexpr),
		Namespaces: convertList(n.Namespaces),
		Columns:    convertList(n.Columns),
		Alias:      convertAlias(n.Alias),
		Location:   n.Location,
	}
}

func convertRangeTableFuncCol(n *nodes.RangeTableFuncCol) *ast.RangeTableFuncCol {
	if n == nil {
		return nil
	}
	return &ast.RangeTableFuncCol{
		Colname:       n.Colname,
		TypeName:      convertTypeName(n.TypeName),
		ForOrdinality: n.ForOrdinality,
		IsNotNull:     n.IsNotNull,
		Colexpr:       convertNode(n.Colexpr),
		Coldefexpr:    convertNode(n.Coldefexpr),
		Location:      n.Location,
	}
}

func convertRangeTableSample(n *nodes.RangeTableSample) *ast.RangeTableSample {
	if n == nil {
		return nil
	}
	return &ast.RangeTableSample{
		Relation:   convertNode(n.Relation),
		Method:     convertList(n.Method),
		Args:       convertList(n.Args),
		Repeatable: convertNode(n.Repeatable),
		Location:   n.Location,
	}
}

func convertRangeTblEntry(n *nodes.RangeTblEntry) *ast.RangeTblEntry {
	if n == nil {
		return nil
	}
	return &ast.RangeTblEntry{
		Rtekind:         ast.RTEKind(n.Rtekind),
		Relid:           ast.Oid(n.Relid),
		Relkind:         n.Relkind,
		Tablesample:     convertTableSampleClause(n.Tablesample),
		Subquery:        convertQuery(n.Subquery),
		SecurityBarrier: n.SecurityBarrier,
		Jointype:        ast.JoinType(n.Jointype),
		Joinaliasvars:   convertList(n.Joinaliasvars),
		Functions:       convertList(n.Functions),
		Funcordinality:  n.Funcordinality,
		Tablefunc:       convertTableFunc(n.Tablefunc),
		ValuesLists:     convertList(n.ValuesLists),
		Ctename:         n.Ctename,
		Ctelevelsup:     ast.Index(n.Ctelevelsup),
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
		RequiredPerms:   ast.AclMode(n.RequiredPerms),
		CheckAsUser:     ast.Oid(n.CheckAsUser),
		SelectedCols:    n.SelectedCols,
		InsertedCols:    n.InsertedCols,
		UpdatedCols:     n.UpdatedCols,
		SecurityQuals:   convertList(n.SecurityQuals),
	}
}

func convertRangeTblFunction(n *nodes.RangeTblFunction) *ast.RangeTblFunction {
	if n == nil {
		return nil
	}
	return &ast.RangeTblFunction{
		Funcexpr:          convertNode(n.Funcexpr),
		Funccolcount:      n.Funccolcount,
		Funccolnames:      convertList(n.Funccolnames),
		Funccoltypes:      convertList(n.Funccoltypes),
		Funccoltypmods:    convertList(n.Funccoltypmods),
		Funccolcollations: convertList(n.Funccolcollations),
		Funcparams:        n.Funcparams,
	}
}

func convertRangeTblRef(n *nodes.RangeTblRef) *ast.RangeTblRef {
	if n == nil {
		return nil
	}
	return &ast.RangeTblRef{
		Rtindex: n.Rtindex,
	}
}

func convertRangeVar(n *nodes.RangeVar) *ast.RangeVar {
	if n == nil {
		return nil
	}
	return &ast.RangeVar{
		Catalogname:    n.Catalogname,
		Schemaname:     n.Schemaname,
		Relname:        n.Relname,
		Inh:            n.Inh,
		Relpersistence: n.Relpersistence,
		Alias:          convertAlias(n.Alias),
		Location:       n.Location,
	}
}

func convertRawStmt(n *nodes.RawStmt) *ast.RawStmt {
	if n == nil {
		return nil
	}
	return &ast.RawStmt{
		Stmt:         convertNode(n.Stmt),
		StmtLocation: n.StmtLocation,
		StmtLen:      n.StmtLen,
	}
}

func convertReassignOwnedStmt(n *nodes.ReassignOwnedStmt) *ast.ReassignOwnedStmt {
	if n == nil {
		return nil
	}
	return &ast.ReassignOwnedStmt{
		Roles:   convertList(n.Roles),
		Newrole: convertRoleSpec(n.Newrole),
	}
}

func convertRefreshMatViewStmt(n *nodes.RefreshMatViewStmt) *ast.RefreshMatViewStmt {
	if n == nil {
		return nil
	}
	return &ast.RefreshMatViewStmt{
		Concurrent: n.Concurrent,
		SkipData:   n.SkipData,
		Relation:   convertRangeVar(n.Relation),
	}
}

func convertReindexStmt(n *nodes.ReindexStmt) *ast.ReindexStmt {
	if n == nil {
		return nil
	}
	return &ast.ReindexStmt{
		Kind:     ast.ReindexObjectType(n.Kind),
		Relation: convertRangeVar(n.Relation),
		Name:     n.Name,
		Options:  n.Options,
	}
}

func convertRelabelType(n *nodes.RelabelType) *ast.RelabelType {
	if n == nil {
		return nil
	}
	return &ast.RelabelType{
		Xpr:           convertNode(n.Xpr),
		Arg:           convertNode(n.Arg),
		Resulttype:    ast.Oid(n.Resulttype),
		Resulttypmod:  n.Resulttypmod,
		Resultcollid:  ast.Oid(n.Resultcollid),
		Relabelformat: ast.CoercionForm(n.Relabelformat),
		Location:      n.Location,
	}
}

func convertRenameStmt(n *nodes.RenameStmt) *ast.RenameStmt {
	if n == nil {
		return nil
	}
	return &ast.RenameStmt{
		RenameType:   ast.ObjectType(n.RenameType),
		RelationType: ast.ObjectType(n.RelationType),
		Relation:     convertRangeVar(n.Relation),
		Object:       convertNode(n.Object),
		Subname:      n.Subname,
		Newname:      n.Newname,
		Behavior:     ast.DropBehavior(n.Behavior),
		MissingOk:    n.MissingOk,
	}
}

func convertReplicaIdentityStmt(n *nodes.ReplicaIdentityStmt) *ast.ReplicaIdentityStmt {
	if n == nil {
		return nil
	}
	return &ast.ReplicaIdentityStmt{
		IdentityType: n.IdentityType,
		Name:         n.Name,
	}
}

func convertResTarget(n *nodes.ResTarget) *ast.ResTarget {
	if n == nil {
		return nil
	}
	return &ast.ResTarget{
		Name:        n.Name,
		Indirection: convertList(n.Indirection),
		Val:         convertNode(n.Val),
		Location:    n.Location,
	}
}

func convertRoleSpec(n *nodes.RoleSpec) *ast.RoleSpec {
	if n == nil {
		return nil
	}
	return &ast.RoleSpec{
		Roletype: ast.RoleSpecType(n.Roletype),
		Rolename: n.Rolename,
		Location: n.Location,
	}
}

func convertRowCompareExpr(n *nodes.RowCompareExpr) *ast.RowCompareExpr {
	if n == nil {
		return nil
	}
	return &ast.RowCompareExpr{
		Xpr:          convertNode(n.Xpr),
		Rctype:       ast.RowCompareType(n.Rctype),
		Opnos:        convertList(n.Opnos),
		Opfamilies:   convertList(n.Opfamilies),
		Inputcollids: convertList(n.Inputcollids),
		Largs:        convertList(n.Largs),
		Rargs:        convertList(n.Rargs),
	}
}

func convertRowExpr(n *nodes.RowExpr) *ast.RowExpr {
	if n == nil {
		return nil
	}
	return &ast.RowExpr{
		Xpr:       convertNode(n.Xpr),
		Args:      convertList(n.Args),
		RowTypeid: ast.Oid(n.RowTypeid),
		RowFormat: ast.CoercionForm(n.RowFormat),
		Colnames:  convertList(n.Colnames),
		Location:  n.Location,
	}
}

func convertRowMarkClause(n *nodes.RowMarkClause) *ast.RowMarkClause {
	if n == nil {
		return nil
	}
	return &ast.RowMarkClause{
		Rti:        ast.Index(n.Rti),
		Strength:   ast.LockClauseStrength(n.Strength),
		WaitPolicy: ast.LockWaitPolicy(n.WaitPolicy),
		PushedDown: n.PushedDown,
	}
}

func convertRuleStmt(n *nodes.RuleStmt) *ast.RuleStmt {
	if n == nil {
		return nil
	}
	return &ast.RuleStmt{
		Relation:    convertRangeVar(n.Relation),
		Rulename:    n.Rulename,
		WhereClause: convertNode(n.WhereClause),
		Event:       ast.CmdType(n.Event),
		Instead:     n.Instead,
		Actions:     convertList(n.Actions),
		Replace:     n.Replace,
	}
}

func convertSQLValueFunction(n *nodes.SQLValueFunction) *ast.SQLValueFunction {
	if n == nil {
		return nil
	}
	return &ast.SQLValueFunction{
		Xpr:      convertNode(n.Xpr),
		Op:       ast.SQLValueFunctionOp(n.Op),
		Type:     ast.Oid(n.Type),
		Typmod:   n.Typmod,
		Location: n.Location,
	}
}

func convertScalarArrayOpExpr(n *nodes.ScalarArrayOpExpr) *ast.ScalarArrayOpExpr {
	if n == nil {
		return nil
	}
	return &ast.ScalarArrayOpExpr{
		Xpr:         convertNode(n.Xpr),
		Opno:        ast.Oid(n.Opno),
		Opfuncid:    ast.Oid(n.Opfuncid),
		UseOr:       n.UseOr,
		Inputcollid: ast.Oid(n.Inputcollid),
		Args:        convertList(n.Args),
		Location:    n.Location,
	}
}

func convertSecLabelStmt(n *nodes.SecLabelStmt) *ast.SecLabelStmt {
	if n == nil {
		return nil
	}
	return &ast.SecLabelStmt{
		Objtype:  ast.ObjectType(n.Objtype),
		Object:   convertNode(n.Object),
		Provider: n.Provider,
		Label:    n.Label,
	}
}

func convertSelectStmt(n *nodes.SelectStmt) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	return &ast.SelectStmt{
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
		Op:             ast.SetOperation(n.Op),
		All:            n.All,
		Larg:           convertSelectStmt(n.Larg),
		Rarg:           convertSelectStmt(n.Rarg),
	}
}

func convertSetOperationStmt(n *nodes.SetOperationStmt) *ast.SetOperationStmt {
	if n == nil {
		return nil
	}
	return &ast.SetOperationStmt{
		Op:            ast.SetOperation(n.Op),
		All:           n.All,
		Larg:          convertNode(n.Larg),
		Rarg:          convertNode(n.Rarg),
		ColTypes:      convertList(n.ColTypes),
		ColTypmods:    convertList(n.ColTypmods),
		ColCollations: convertList(n.ColCollations),
		GroupClauses:  convertList(n.GroupClauses),
	}
}

func convertSetToDefault(n *nodes.SetToDefault) *ast.SetToDefault {
	if n == nil {
		return nil
	}
	return &ast.SetToDefault{
		Xpr:       convertNode(n.Xpr),
		TypeId:    ast.Oid(n.TypeId),
		TypeMod:   n.TypeMod,
		Collation: ast.Oid(n.Collation),
		Location:  n.Location,
	}
}

func convertSortBy(n *nodes.SortBy) *ast.SortBy {
	if n == nil {
		return nil
	}
	return &ast.SortBy{
		Node:        convertNode(n.Node),
		SortbyDir:   ast.SortByDir(n.SortbyDir),
		SortbyNulls: ast.SortByNulls(n.SortbyNulls),
		UseOp:       convertList(n.UseOp),
		Location:    n.Location,
	}
}

func convertSortGroupClause(n *nodes.SortGroupClause) *ast.SortGroupClause {
	if n == nil {
		return nil
	}
	return &ast.SortGroupClause{
		TleSortGroupRef: ast.Index(n.TleSortGroupRef),
		Eqop:            ast.Oid(n.Eqop),
		Sortop:          ast.Oid(n.Sortop),
		NullsFirst:      n.NullsFirst,
		Hashable:        n.Hashable,
	}
}

func convertString(n *nodes.String) *ast.String {
	if n == nil {
		return nil
	}
	return &ast.String{
		Str: n.Str,
	}
}

func convertSubLink(n *nodes.SubLink) *ast.SubLink {
	if n == nil {
		return nil
	}
	return &ast.SubLink{
		Xpr:         convertNode(n.Xpr),
		SubLinkType: ast.SubLinkType(n.SubLinkType),
		SubLinkId:   n.SubLinkId,
		Testexpr:    convertNode(n.Testexpr),
		OperName:    convertList(n.OperName),
		Subselect:   convertNode(n.Subselect),
		Location:    n.Location,
	}
}

func convertSubPlan(n *nodes.SubPlan) *ast.SubPlan {
	if n == nil {
		return nil
	}
	return &ast.SubPlan{
		Xpr:               convertNode(n.Xpr),
		SubLinkType:       ast.SubLinkType(n.SubLinkType),
		Testexpr:          convertNode(n.Testexpr),
		ParamIds:          convertList(n.ParamIds),
		PlanId:            n.PlanId,
		PlanName:          n.PlanName,
		FirstColType:      ast.Oid(n.FirstColType),
		FirstColTypmod:    n.FirstColTypmod,
		FirstColCollation: ast.Oid(n.FirstColCollation),
		UseHashTable:      n.UseHashTable,
		UnknownEqFalse:    n.UnknownEqFalse,
		ParallelSafe:      n.ParallelSafe,
		SetParam:          convertList(n.SetParam),
		ParParam:          convertList(n.ParParam),
		Args:              convertList(n.Args),
		StartupCost:       ast.Cost(n.StartupCost),
		PerCallCost:       ast.Cost(n.PerCallCost),
	}
}

func convertTableFunc(n *nodes.TableFunc) *ast.TableFunc {
	if n == nil {
		return nil
	}
	return &ast.TableFunc{
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

func convertTableLikeClause(n *nodes.TableLikeClause) *ast.TableLikeClause {
	if n == nil {
		return nil
	}
	return &ast.TableLikeClause{
		Relation: convertRangeVar(n.Relation),
		Options:  n.Options,
	}
}

func convertTableSampleClause(n *nodes.TableSampleClause) *ast.TableSampleClause {
	if n == nil {
		return nil
	}
	return &ast.TableSampleClause{
		Tsmhandler: ast.Oid(n.Tsmhandler),
		Args:       convertList(n.Args),
		Repeatable: convertNode(n.Repeatable),
	}
}

func convertTargetEntry(n *nodes.TargetEntry) *ast.TargetEntry {
	if n == nil {
		return nil
	}
	return &ast.TargetEntry{
		Xpr:             convertNode(n.Xpr),
		Expr:            convertNode(n.Expr),
		Resno:           ast.AttrNumber(n.Resno),
		Resname:         n.Resname,
		Ressortgroupref: ast.Index(n.Ressortgroupref),
		Resorigtbl:      ast.Oid(n.Resorigtbl),
		Resorigcol:      ast.AttrNumber(n.Resorigcol),
		Resjunk:         n.Resjunk,
	}
}

func convertTransactionStmt(n *nodes.TransactionStmt) *ast.TransactionStmt {
	if n == nil {
		return nil
	}
	return &ast.TransactionStmt{
		Kind:    ast.TransactionStmtKind(n.Kind),
		Options: convertList(n.Options),
		Gid:     n.Gid,
	}
}

func convertTriggerTransition(n *nodes.TriggerTransition) *ast.TriggerTransition {
	if n == nil {
		return nil
	}
	return &ast.TriggerTransition{
		Name:    n.Name,
		IsNew:   n.IsNew,
		IsTable: n.IsTable,
	}
}

func convertTruncateStmt(n *nodes.TruncateStmt) *ast.TruncateStmt {
	if n == nil {
		return nil
	}
	return &ast.TruncateStmt{
		Relations:   convertList(n.Relations),
		RestartSeqs: n.RestartSeqs,
		Behavior:    ast.DropBehavior(n.Behavior),
	}
}

func convertTypeCast(n *nodes.TypeCast) *ast.TypeCast {
	if n == nil {
		return nil
	}
	return &ast.TypeCast{
		Arg:      convertNode(n.Arg),
		TypeName: convertTypeName(n.TypeName),
		Location: n.Location,
	}
}

func convertTypeName(n *nodes.TypeName) *ast.TypeName {
	if n == nil {
		return nil
	}
	return &ast.TypeName{
		Names:       convertList(n.Names),
		TypeOid:     ast.Oid(n.TypeOid),
		Setof:       n.Setof,
		PctType:     n.PctType,
		Typmods:     convertList(n.Typmods),
		Typemod:     n.Typemod,
		ArrayBounds: convertList(n.ArrayBounds),
		Location:    n.Location,
	}
}

func convertUnlistenStmt(n *nodes.UnlistenStmt) *ast.UnlistenStmt {
	if n == nil {
		return nil
	}
	return &ast.UnlistenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertUpdateStmt(n *nodes.UpdateStmt) *ast.UpdateStmt {
	if n == nil {
		return nil
	}
	return &ast.UpdateStmt{
		Relation:      convertRangeVar(n.Relation),
		TargetList:    convertList(n.TargetList),
		WhereClause:   convertNode(n.WhereClause),
		FromClause:    convertList(n.FromClause),
		ReturningList: convertList(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertVacuumStmt(n *nodes.VacuumStmt) *ast.VacuumStmt {
	if n == nil {
		return nil
	}
	return &ast.VacuumStmt{
		Options:  n.Options,
		Relation: convertRangeVar(n.Relation),
		VaCols:   convertList(n.VaCols),
	}
}

func convertVar(n *nodes.Var) *ast.Var {
	if n == nil {
		return nil
	}
	return &ast.Var{
		Xpr:         convertNode(n.Xpr),
		Varno:       ast.Index(n.Varno),
		Varattno:    ast.AttrNumber(n.Varattno),
		Vartype:     ast.Oid(n.Vartype),
		Vartypmod:   n.Vartypmod,
		Varcollid:   ast.Oid(n.Varcollid),
		Varlevelsup: ast.Index(n.Varlevelsup),
		Varnoold:    ast.Index(n.Varnoold),
		Varoattno:   ast.AttrNumber(n.Varoattno),
		Location:    n.Location,
	}
}

func convertVariableSetStmt(n *nodes.VariableSetStmt) *ast.VariableSetStmt {
	if n == nil {
		return nil
	}
	return &ast.VariableSetStmt{
		Kind:    ast.VariableSetKind(n.Kind),
		Name:    n.Name,
		Args:    convertList(n.Args),
		IsLocal: n.IsLocal,
	}
}

func convertVariableShowStmt(n *nodes.VariableShowStmt) *ast.VariableShowStmt {
	if n == nil {
		return nil
	}
	return &ast.VariableShowStmt{
		Name: n.Name,
	}
}

func convertViewStmt(n *nodes.ViewStmt) *ast.ViewStmt {
	if n == nil {
		return nil
	}
	return &ast.ViewStmt{
		View:            convertRangeVar(n.View),
		Aliases:         convertList(n.Aliases),
		Query:           convertNode(n.Query),
		Replace:         n.Replace,
		Options:         convertList(n.Options),
		WithCheckOption: ast.ViewCheckOption(n.WithCheckOption),
	}
}

func convertWindowClause(n *nodes.WindowClause) *ast.WindowClause {
	if n == nil {
		return nil
	}
	return &ast.WindowClause{
		Name:            n.Name,
		Refname:         n.Refname,
		PartitionClause: convertList(n.PartitionClause),
		OrderClause:     convertList(n.OrderClause),
		FrameOptions:    n.FrameOptions,
		StartOffset:     convertNode(n.StartOffset),
		EndOffset:       convertNode(n.EndOffset),
		Winref:          ast.Index(n.Winref),
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

func convertWindowFunc(n *nodes.WindowFunc) *ast.WindowFunc {
	if n == nil {
		return nil
	}
	return &ast.WindowFunc{
		Xpr:         convertNode(n.Xpr),
		Winfnoid:    ast.Oid(n.Winfnoid),
		Wintype:     ast.Oid(n.Wintype),
		Wincollid:   ast.Oid(n.Wincollid),
		Inputcollid: ast.Oid(n.Inputcollid),
		Args:        convertList(n.Args),
		Aggfilter:   convertNode(n.Aggfilter),
		Winref:      ast.Index(n.Winref),
		Winstar:     n.Winstar,
		Winagg:      n.Winagg,
		Location:    n.Location,
	}
}

func convertWithCheckOption(n *nodes.WithCheckOption) *ast.WithCheckOption {
	if n == nil {
		return nil
	}
	return &ast.WithCheckOption{
		Kind:     ast.WCOKind(n.Kind),
		Relname:  n.Relname,
		Polname:  n.Polname,
		Qual:     convertNode(n.Qual),
		Cascaded: n.Cascaded,
	}
}

func convertWithClause(n *nodes.WithClause) *ast.WithClause {
	if n == nil {
		return nil
	}
	return &ast.WithClause{
		Ctes:      convertList(n.Ctes),
		Recursive: n.Recursive,
		Location:  n.Location,
	}
}

func convertXmlExpr(n *nodes.XmlExpr) *ast.XmlExpr {
	if n == nil {
		return nil
	}
	return &ast.XmlExpr{
		Xpr:       convertNode(n.Xpr),
		Op:        ast.XmlExprOp(n.Op),
		Name:      n.Name,
		NamedArgs: convertList(n.NamedArgs),
		ArgNames:  convertList(n.ArgNames),
		Args:      convertList(n.Args),
		Xmloption: ast.XmlOptionType(n.Xmloption),
		Type:      ast.Oid(n.Type),
		Typmod:    n.Typmod,
		Location:  n.Location,
	}
}

func convertXmlSerialize(n *nodes.XmlSerialize) *ast.XmlSerialize {
	if n == nil {
		return nil
	}
	return &ast.XmlSerialize{
		Xmloption: ast.XmlOptionType(n.Xmloption),
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
