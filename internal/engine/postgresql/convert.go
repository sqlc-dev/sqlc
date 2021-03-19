// +build !windows

package postgresql

import (
	pg "github.com/pganalyze/pg_query_go/v2"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func makeByte(s string) byte {
	var b byte
	if s == "" {
		return b
	}
	return []byte(s)[0]
}

func convertList(l pg.List) *ast.List {
	out := &ast.List{}
	for _, item := range l.Items {
		out.Items = append(out.Items, convertNode(item))
	}
	return out
}

func convertSlice(nodes []*pg.Node) *ast.List {
	out := &ast.List{}
	for _, n := range nodes {
		out.Items = append(out.Items, convertNode(n))
	}
	return out
}

func convertValuesList(l [][]*pg.Node) *ast.List {
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

func convert(node *pg.Node) (ast.Node, error) {
	return convertNode(node), nil
}

func convertA_ArrayExpr(n *pg.A_ArrayExpr) *ast.A_ArrayExpr {
	if n == nil {
		return nil
	}
	return &ast.A_ArrayExpr{
		Elements: convertSlice(n.Elements),
		Location: int(n.Location),
	}
}

func convertA_Const(n *pg.A_Const) *ast.A_Const {
	if n == nil {
		return nil
	}
	return &ast.A_Const{
		Val:      convertNode(n.Val),
		Location: int(n.Location),
	}
}

func convertA_Expr(n *pg.A_Expr) *ast.A_Expr {
	if n == nil {
		return nil
	}
	return &ast.A_Expr{
		Kind:     ast.A_Expr_Kind(n.Kind),
		Name:     convertSlice(n.Name),
		Lexpr:    convertNode(n.Lexpr),
		Rexpr:    convertNode(n.Rexpr),
		Location: int(n.Location),
	}
}

func convertA_Indices(n *pg.A_Indices) *ast.A_Indices {
	if n == nil {
		return nil
	}
	return &ast.A_Indices{
		IsSlice: n.IsSlice,
		Lidx:    convertNode(n.Lidx),
		Uidx:    convertNode(n.Uidx),
	}
}

func convertA_Indirection(n *pg.A_Indirection) *ast.A_Indirection {
	if n == nil {
		return nil
	}
	return &ast.A_Indirection{
		Arg:         convertNode(n.Arg),
		Indirection: convertSlice(n.Indirection),
	}
}

func convertA_Star(n *pg.A_Star) *ast.A_Star {
	if n == nil {
		return nil
	}
	return &ast.A_Star{}
}

func convertAccessPriv(n *pg.AccessPriv) *ast.AccessPriv {
	if n == nil {
		return nil
	}
	return &ast.AccessPriv{
		PrivName: &n.PrivName,
		Cols:     convertSlice(n.Cols),
	}
}

func convertAggref(n *pg.Aggref) *ast.Aggref {
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
		Aggargtypes:   convertSlice(n.Aggargtypes),
		Aggdirectargs: convertSlice(n.Aggdirectargs),
		Args:          convertSlice(n.Args),
		Aggorder:      convertSlice(n.Aggorder),
		Aggdistinct:   convertSlice(n.Aggdistinct),
		Aggfilter:     convertNode(n.Aggfilter),
		Aggstar:       n.Aggstar,
		Aggvariadic:   n.Aggvariadic,
		Aggkind:       makeByte(n.Aggkind),
		Agglevelsup:   ast.Index(n.Agglevelsup),
		Aggsplit:      ast.AggSplit(n.Aggsplit),
		Location:      int(n.Location),
	}
}

func convertAlias(n *pg.Alias) *ast.Alias {
	if n == nil {
		return nil
	}
	return &ast.Alias{
		Aliasname: &n.Aliasname,
		Colnames:  convertSlice(n.Colnames),
	}
}

func convertAlterCollationStmt(n *pg.AlterCollationStmt) *ast.AlterCollationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterCollationStmt{
		Collname: convertSlice(n.Collname),
	}
}

func convertAlterDatabaseSetStmt(n *pg.AlterDatabaseSetStmt) *ast.AlterDatabaseSetStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDatabaseSetStmt{
		Dbname:  &n.Dbname,
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterDatabaseStmt(n *pg.AlterDatabaseStmt) *ast.AlterDatabaseStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDatabaseStmt{
		Dbname:  &n.Dbname,
		Options: convertSlice(n.Options),
	}
}

func convertAlterDefaultPrivilegesStmt(n *pg.AlterDefaultPrivilegesStmt) *ast.AlterDefaultPrivilegesStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDefaultPrivilegesStmt{
		Options: convertSlice(n.Options),
		Action:  convertGrantStmt(n.Action),
	}
}

func convertAlterDomainStmt(n *pg.AlterDomainStmt) *ast.AlterDomainStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterDomainStmt{
		Subtype:   makeByte(n.Subtype),
		TypeName:  convertSlice(n.TypeName),
		Name:      &n.Name,
		Def:       convertNode(n.Def),
		Behavior:  ast.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterEnumStmt(n *pg.AlterEnumStmt) *ast.AlterEnumStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterEnumStmt{
		TypeName:           convertSlice(n.TypeName),
		OldVal:             &n.OldVal,
		NewVal:             &n.NewVal,
		NewValNeighbor:     &n.NewValNeighbor,
		NewValIsAfter:      n.NewValIsAfter,
		SkipIfNewValExists: n.SkipIfNewValExists,
	}
}

func convertAlterEventTrigStmt(n *pg.AlterEventTrigStmt) *ast.AlterEventTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterEventTrigStmt{
		Trigname:  &n.Trigname,
		Tgenabled: makeByte(n.Tgenabled),
	}
}

func convertAlterExtensionContentsStmt(n *pg.AlterExtensionContentsStmt) *ast.AlterExtensionContentsStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterExtensionContentsStmt{
		Extname: &n.Extname,
		Action:  int(n.Action),
		Objtype: ast.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
	}
}

func convertAlterExtensionStmt(n *pg.AlterExtensionStmt) *ast.AlterExtensionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterExtensionStmt{
		Extname: &n.Extname,
		Options: convertSlice(n.Options),
	}
}

func convertAlterFdwStmt(n *pg.AlterFdwStmt) *ast.AlterFdwStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterFdwStmt{
		Fdwname:     &n.Fdwname,
		FuncOptions: convertSlice(n.FuncOptions),
		Options:     convertSlice(n.Options),
	}
}

func convertAlterForeignServerStmt(n *pg.AlterForeignServerStmt) *ast.AlterForeignServerStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterForeignServerStmt{
		Servername: &n.Servername,
		Version:    &n.Version,
		Options:    convertSlice(n.Options),
		HasVersion: n.HasVersion,
	}
}

func convertAlterFunctionStmt(n *pg.AlterFunctionStmt) *ast.AlterFunctionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterFunctionStmt{
		Func:    convertObjectWithArgs(n.Func),
		Actions: convertSlice(n.Actions),
	}
}

func convertAlterObjectDependsStmt(n *pg.AlterObjectDependsStmt) *ast.AlterObjectDependsStmt {
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

func convertAlterObjectSchemaStmt(n *pg.AlterObjectSchemaStmt) *ast.AlterObjectSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterObjectSchemaStmt{
		ObjectType: ast.ObjectType(n.ObjectType),
		Relation:   convertRangeVar(n.Relation),
		Object:     convertNode(n.Object),
		Newschema:  &n.Newschema,
		MissingOk:  n.MissingOk,
	}
}

func convertAlterOpFamilyStmt(n *pg.AlterOpFamilyStmt) *ast.AlterOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterOpFamilyStmt{
		Opfamilyname: convertSlice(n.Opfamilyname),
		Amname:       &n.Amname,
		IsDrop:       n.IsDrop,
		Items:        convertSlice(n.Items),
	}
}

func convertAlterOperatorStmt(n *pg.AlterOperatorStmt) *ast.AlterOperatorStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterOperatorStmt{
		Opername: convertObjectWithArgs(n.Opername),
		Options:  convertSlice(n.Options),
	}
}

func convertAlterOwnerStmt(n *pg.AlterOwnerStmt) *ast.AlterOwnerStmt {
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

func convertAlterPolicyStmt(n *pg.AlterPolicyStmt) *ast.AlterPolicyStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterPolicyStmt{
		PolicyName: &n.PolicyName,
		Table:      convertRangeVar(n.Table),
		Roles:      convertSlice(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertAlterPublicationStmt(n *pg.AlterPublicationStmt) *ast.AlterPublicationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterPublicationStmt{
		Pubname:      &n.Pubname,
		Options:      convertSlice(n.Options),
		Tables:       convertSlice(n.Tables),
		ForAllTables: n.ForAllTables,
		TableAction:  ast.DefElemAction(n.TableAction),
	}
}

func convertAlterRoleSetStmt(n *pg.AlterRoleSetStmt) *ast.AlterRoleSetStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterRoleSetStmt{
		Role:     convertRoleSpec(n.Role),
		Database: &n.Database,
		Setstmt:  convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterRoleStmt(n *pg.AlterRoleStmt) *ast.AlterRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterRoleStmt{
		Role:    convertRoleSpec(n.Role),
		Options: convertSlice(n.Options),
		Action:  int(n.Action),
	}
}

func convertAlterSeqStmt(n *pg.AlterSeqStmt) *ast.AlterSeqStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertSlice(n.Options),
		ForIdentity: n.ForIdentity,
		MissingOk:   n.MissingOk,
	}
}

func convertAlterSubscriptionStmt(n *pg.AlterSubscriptionStmt) *ast.AlterSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSubscriptionStmt{
		Kind:        ast.AlterSubscriptionType(n.Kind),
		Subname:     &n.Subname,
		Conninfo:    &n.Conninfo,
		Publication: convertSlice(n.Publication),
		Options:     convertSlice(n.Options),
	}
}

func convertAlterSystemStmt(n *pg.AlterSystemStmt) *ast.AlterSystemStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterSystemStmt{
		Setstmt: convertVariableSetStmt(n.Setstmt),
	}
}

func convertAlterTSConfigurationStmt(n *pg.AlterTSConfigurationStmt) *ast.AlterTSConfigurationStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTSConfigurationStmt{
		Kind:      ast.AlterTSConfigType(n.Kind),
		Cfgname:   convertSlice(n.Cfgname),
		Tokentype: convertSlice(n.Tokentype),
		Dicts:     convertSlice(n.Dicts),
		Override:  n.Override,
		Replace:   n.Replace,
		MissingOk: n.MissingOk,
	}
}

func convertAlterTSDictionaryStmt(n *pg.AlterTSDictionaryStmt) *ast.AlterTSDictionaryStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTSDictionaryStmt{
		Dictname: convertSlice(n.Dictname),
		Options:  convertSlice(n.Options),
	}
}

func convertAlterTableCmd(n *pg.AlterTableCmd) *ast.AlterTableCmd {
	if n == nil {
		return nil
	}
	def := convertNode(n.Def)
	columnDef := def.(*ast.ColumnDef)
	return &ast.AlterTableCmd{
		Subtype:   ast.AlterTableType(n.Subtype),
		Name:      &n.Name,
		Newowner:  convertRoleSpec(n.Newowner),
		Def:       columnDef,
		Behavior:  ast.DropBehavior(n.Behavior),
		MissingOk: n.MissingOk,
	}
}

func convertAlterTableMoveAllStmt(n *pg.AlterTableMoveAllStmt) *ast.AlterTableMoveAllStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableMoveAllStmt{
		OrigTablespacename: &n.OrigTablespacename,
		Objtype:            ast.ObjectType(n.Objtype),
		Roles:              convertSlice(n.Roles),
		NewTablespacename:  &n.NewTablespacename,
		Nowait:             n.Nowait,
	}
}

func convertAlterTableSpaceOptionsStmt(n *pg.AlterTableSpaceOptionsStmt) *ast.AlterTableSpaceOptionsStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableSpaceOptionsStmt{
		Tablespacename: &n.Tablespacename,
		Options:        convertSlice(n.Options),
		IsReset:        n.IsReset,
	}
}

func convertAlterTableStmt(n *pg.AlterTableStmt) *ast.AlterTableStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterTableStmt{
		Relation:  convertRangeVar(n.Relation),
		Cmds:      convertSlice(n.Cmds),
		Relkind:   ast.ObjectType(n.Relkind),
		MissingOk: n.MissingOk,
	}
}

func convertAlterUserMappingStmt(n *pg.AlterUserMappingStmt) *ast.AlterUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.AlterUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: &n.Servername,
		Options:    convertSlice(n.Options),
	}
}

func convertAlternativeSubPlan(n *pg.AlternativeSubPlan) *ast.AlternativeSubPlan {
	if n == nil {
		return nil
	}
	return &ast.AlternativeSubPlan{
		Xpr:      convertNode(n.Xpr),
		Subplans: convertSlice(n.Subplans),
	}
}

func convertArrayCoerceExpr(n *pg.ArrayCoerceExpr) *ast.ArrayCoerceExpr {
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

func convertArrayExpr(n *pg.ArrayExpr) *ast.ArrayExpr {
	if n == nil {
		return nil
	}
	return &ast.ArrayExpr{
		Xpr:           convertNode(n.Xpr),
		ArrayTypeid:   ast.Oid(n.ArrayTypeid),
		ArrayCollid:   ast.Oid(n.ArrayCollid),
		ElementTypeid: ast.Oid(n.ElementTypeid),
		Elements:      convertSlice(n.Elements),
		Multidims:     n.Multidims,
		Location:      n.Location,
	}
}

func convertArrayRef(n *pg.ArrayRef) *ast.ArrayRef {
	if n == nil {
		return nil
	}
	return &ast.ArrayRef{
		Xpr:             convertNode(n.Xpr),
		Refarraytype:    ast.Oid(n.Refarraytype),
		Refelemtype:     ast.Oid(n.Refelemtype),
		Reftypmod:       n.Reftypmod,
		Refcollid:       ast.Oid(n.Refcollid),
		Refupperindexpr: convertSlice(n.Refupperindexpr),
		Reflowerindexpr: convertSlice(n.Reflowerindexpr),
		Refexpr:         convertNode(n.Refexpr),
		Refassgnexpr:    convertNode(n.Refassgnexpr),
	}
}

func convertBitString(n *pg.BitString) *ast.BitString {
	if n == nil {
		return nil
	}
	return &ast.BitString{
		Str: n.Str,
	}
}

func convertBlockIdData(n *pg.BlockIdData) *ast.BlockIdData {
	if n == nil {
		return nil
	}
	return &ast.BlockIdData{
		BiHi: n.BiHi,
		BiLo: n.BiLo,
	}
}

func convertBoolExpr(n *pg.BoolExpr) *ast.BoolExpr {
	if n == nil {
		return nil
	}
	return &ast.BoolExpr{
		Xpr:      convertNode(n.Xpr),
		Boolop:   ast.BoolExprType(n.Boolop),
		Args:     convertSlice(n.Args),
		Location: n.Location,
	}
}

func convertBooleanTest(n *pg.BooleanTest) *ast.BooleanTest {
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

func convertCaseExpr(n *pg.CaseExpr) *ast.CaseExpr {
	if n == nil {
		return nil
	}
	return &ast.CaseExpr{
		Xpr:        convertNode(n.Xpr),
		Casetype:   ast.Oid(n.Casetype),
		Casecollid: ast.Oid(n.Casecollid),
		Arg:        convertNode(n.Arg),
		Args:       convertSlice(n.Args),
		Defresult:  convertNode(n.Defresult),
		Location:   n.Location,
	}
}

func convertCaseTestExpr(n *pg.CaseTestExpr) *ast.CaseTestExpr {
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

func convertCaseWhen(n *pg.CaseWhen) *ast.CaseWhen {
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

func convertCheckPointStmt(n *pg.CheckPointStmt) *ast.CheckPointStmt {
	if n == nil {
		return nil
	}
	return &ast.CheckPointStmt{}
}

func convertClosePortalStmt(n *pg.ClosePortalStmt) *ast.ClosePortalStmt {
	if n == nil {
		return nil
	}
	return &ast.ClosePortalStmt{
		Portalname: n.Portalname,
	}
}

func convertClusterStmt(n *pg.ClusterStmt) *ast.ClusterStmt {
	if n == nil {
		return nil
	}
	return &ast.ClusterStmt{
		Relation:  convertRangeVar(n.Relation),
		Indexname: n.Indexname,
		Verbose:   n.Verbose,
	}
}

func convertCoalesceExpr(n *pg.CoalesceExpr) *ast.CoalesceExpr {
	if n == nil {
		return nil
	}
	return &ast.CoalesceExpr{
		Xpr:            convertNode(n.Xpr),
		Coalescetype:   ast.Oid(n.Coalescetype),
		Coalescecollid: ast.Oid(n.Coalescecollid),
		Args:           convertSlice(n.Args),
		Location:       n.Location,
	}
}

func convertCoerceToDomain(n *pg.CoerceToDomain) *ast.CoerceToDomain {
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

func convertCoerceToDomainValue(n *pg.CoerceToDomainValue) *ast.CoerceToDomainValue {
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

func convertCoerceViaIO(n *pg.CoerceViaIO) *ast.CoerceViaIO {
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

func convertCollateClause(n *pg.CollateClause) *ast.CollateClause {
	if n == nil {
		return nil
	}
	return &ast.CollateClause{
		Arg:      convertNode(n.Arg),
		Collname: convertSlice(n.Collname),
		Location: n.Location,
	}
}

func convertCollateExpr(n *pg.CollateExpr) *ast.CollateExpr {
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

func convertColumnDef(n *pg.ColumnDef) *ast.ColumnDef {
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
		Constraints:   convertSlice(n.Constraints),
		Fdwoptions:    convertSlice(n.Fdwoptions),
		Location:      n.Location,
	}
}

func convertColumnRef(n *pg.ColumnRef) *ast.ColumnRef {
	if n == nil {
		return nil
	}
	return &ast.ColumnRef{
		Fields:   convertSlice(n.Fields),
		Location: n.Location,
	}
}

func convertCommentStmt(n *pg.CommentStmt) *ast.CommentStmt {
	if n == nil {
		return nil
	}
	return &ast.CommentStmt{
		Objtype: ast.ObjectType(n.Objtype),
		Object:  convertNode(n.Object),
		Comment: n.Comment,
	}
}

func convertCommonTableExpr(n *pg.CommonTableExpr) *ast.CommonTableExpr {
	if n == nil {
		return nil
	}
	return &ast.CommonTableExpr{
		Ctename:          n.Ctename,
		Aliascolnames:    convertSlice(n.Aliascolnames),
		Ctequery:         convertNode(n.Ctequery),
		Location:         n.Location,
		Cterecursive:     n.Cterecursive,
		Cterefcount:      n.Cterefcount,
		Ctecolnames:      convertSlice(n.Ctecolnames),
		Ctecoltypes:      convertSlice(n.Ctecoltypes),
		Ctecoltypmods:    convertSlice(n.Ctecoltypmods),
		Ctecolcollations: convertSlice(n.Ctecolcollations),
	}
}

func convertCompositeTypeStmt(n *pg.CompositeTypeStmt) *ast.CompositeTypeStmt {
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

func convertConst(n *pg.Const) *ast.Const {
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

func convertConstraint(n *pg.Constraint) *ast.Constraint {
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
		Keys:           convertSlice(n.Keys),
		Exclusions:     convertSlice(n.Exclusions),
		Options:        convertSlice(n.Options),
		Indexname:      n.Indexname,
		Indexspace:     n.Indexspace,
		AccessMethod:   n.AccessMethod,
		WhereClause:    convertNode(n.WhereClause),
		Pktable:        convertRangeVar(n.Pktable),
		FkAttrs:        convertSlice(n.FkAttrs),
		PkAttrs:        convertSlice(n.PkAttrs),
		FkMatchtype:    n.FkMatchtype,
		FkUpdAction:    n.FkUpdAction,
		FkDelAction:    n.FkDelAction,
		OldConpfeqop:   convertSlice(n.OldConpfeqop),
		OldPktableOid:  ast.Oid(n.OldPktableOid),
		SkipValidation: n.SkipValidation,
		InitiallyValid: n.InitiallyValid,
	}
}

func convertConstraintsSetStmt(n *pg.ConstraintsSetStmt) *ast.ConstraintsSetStmt {
	if n == nil {
		return nil
	}
	return &ast.ConstraintsSetStmt{
		Constraints: convertSlice(n.Constraints),
		Deferred:    n.Deferred,
	}
}

func convertConvertRowtypeExpr(n *pg.ConvertRowtypeExpr) *ast.ConvertRowtypeExpr {
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

func convertCopyStmt(n *pg.CopyStmt) *ast.CopyStmt {
	if n == nil {
		return nil
	}
	return &ast.CopyStmt{
		Relation:  convertRangeVar(n.Relation),
		Query:     convertNode(n.Query),
		Attlist:   convertSlice(n.Attlist),
		IsFrom:    n.IsFrom,
		IsProgram: n.IsProgram,
		Filename:  n.Filename,
		Options:   convertSlice(n.Options),
	}
}

func convertCreateAmStmt(n *pg.CreateAmStmt) *ast.CreateAmStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateAmStmt{
		Amname:      n.Amname,
		HandlerName: convertSlice(n.HandlerName),
		Amtype:      n.Amtype,
	}
}

func convertCreateCastStmt(n *pg.CreateCastStmt) *ast.CreateCastStmt {
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

func convertCreateConversionStmt(n *pg.CreateConversionStmt) *ast.CreateConversionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateConversionStmt{
		ConversionName:  convertSlice(n.ConversionName),
		ForEncodingName: n.ForEncodingName,
		ToEncodingName:  n.ToEncodingName,
		FuncName:        convertSlice(n.FuncName),
		Def:             n.Def,
	}
}

func convertCreateDomainStmt(n *pg.CreateDomainStmt) *ast.CreateDomainStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateDomainStmt{
		Domainname:  convertSlice(n.Domainname),
		TypeName:    convertTypeName(n.TypeName),
		CollClause:  convertCollateClause(n.CollClause),
		Constraints: convertSlice(n.Constraints),
	}
}

func convertCreateEnumStmt(n *pg.CreateEnumStmt) *ast.CreateEnumStmt {
	if n == nil {
		return nil
	}
	tn, err := parseTypeName(n.TypeName)
	if err != nil {
		panic(err)
	}
	return &ast.CreateEnumStmt{
		TypeName: tn,
		Vals:     convertSlice(n.Vals),
	}
}

func convertCreateEventTrigStmt(n *pg.CreateEventTrigStmt) *ast.CreateEventTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateEventTrigStmt{
		Trigname:   n.Trigname,
		Eventname:  n.Eventname,
		Whenclause: convertSlice(n.Whenclause),
		Funcname:   convertSlice(n.Funcname),
	}
}

func convertCreateExtensionStmt(n *pg.CreateExtensionStmt) *ast.CreateExtensionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateExtensionStmt{
		Extname:     n.Extname,
		IfNotExists: n.IfNotExists,
		Options:     convertSlice(n.Options),
	}
}

func convertCreateFdwStmt(n *pg.CreateFdwStmt) *ast.CreateFdwStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateFdwStmt{
		Fdwname:     n.Fdwname,
		FuncOptions: convertSlice(n.FuncOptions),
		Options:     convertSlice(n.Options),
	}
}

func convertCreateForeignServerStmt(n *pg.CreateForeignServerStmt) *ast.CreateForeignServerStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateForeignServerStmt{
		Servername:  n.Servername,
		Servertype:  n.Servertype,
		Version:     n.Version,
		Fdwname:     n.Fdwname,
		IfNotExists: n.IfNotExists,
		Options:     convertSlice(n.Options),
	}
}

func convertCreateForeignTableStmt(n *pg.CreateForeignTableStmt) *ast.CreateForeignTableStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateForeignTableStmt{
		Base:       convertCreateStmt(&n.Base),
		Servername: n.Servername,
		Options:    convertSlice(n.Options),
	}
}

func convertCreateFunctionStmt(n *pg.CreateFunctionStmt) *ast.CreateFunctionStmt {
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
		Params:     convertSlice(n.Parameters),
		ReturnType: convertTypeName(n.ReturnType),
		Options:    convertSlice(n.Options),
		WithClause: convertSlice(n.WithClause),
	}
}

func convertCreateOpClassItem(n *pg.CreateOpClassItem) *ast.CreateOpClassItem {
	if n == nil {
		return nil
	}
	return &ast.CreateOpClassItem{
		Itemtype:    n.Itemtype,
		Name:        convertObjectWithArgs(n.Name),
		Number:      n.Number,
		OrderFamily: convertSlice(n.OrderFamily),
		ClassArgs:   convertSlice(n.ClassArgs),
		Storedtype:  convertTypeName(n.Storedtype),
	}
}

func convertCreateOpClassStmt(n *pg.CreateOpClassStmt) *ast.CreateOpClassStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateOpClassStmt{
		Opclassname:  convertSlice(n.Opclassname),
		Opfamilyname: convertSlice(n.Opfamilyname),
		Amname:       n.Amname,
		Datatype:     convertTypeName(n.Datatype),
		Items:        convertSlice(n.Items),
		IsDefault:    n.IsDefault,
	}
}

func convertCreateOpFamilyStmt(n *pg.CreateOpFamilyStmt) *ast.CreateOpFamilyStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateOpFamilyStmt{
		Opfamilyname: convertSlice(n.Opfamilyname),
		Amname:       n.Amname,
	}
}

func convertCreatePLangStmt(n *pg.CreatePLangStmt) *ast.CreatePLangStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePLangStmt{
		Replace:     n.Replace,
		Plname:      n.Plname,
		Plhandler:   convertSlice(n.Plhandler),
		Plinline:    convertSlice(n.Plinline),
		Plvalidator: convertSlice(n.Plvalidator),
		Pltrusted:   n.Pltrusted,
	}
}

func convertCreatePolicyStmt(n *pg.CreatePolicyStmt) *ast.CreatePolicyStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePolicyStmt{
		PolicyName: n.PolicyName,
		Table:      convertRangeVar(n.Table),
		CmdName:    n.CmdName,
		Permissive: n.Permissive,
		Roles:      convertSlice(n.Roles),
		Qual:       convertNode(n.Qual),
		WithCheck:  convertNode(n.WithCheck),
	}
}

func convertCreatePublicationStmt(n *pg.CreatePublicationStmt) *ast.CreatePublicationStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatePublicationStmt{
		Pubname:      n.Pubname,
		Options:      convertSlice(n.Options),
		Tables:       convertSlice(n.Tables),
		ForAllTables: n.ForAllTables,
	}
}

func convertCreateRangeStmt(n *pg.CreateRangeStmt) *ast.CreateRangeStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateRangeStmt{
		TypeName: convertSlice(n.TypeName),
		Params:   convertSlice(n.Params),
	}
}

func convertCreateRoleStmt(n *pg.CreateRoleStmt) *ast.CreateRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateRoleStmt{
		StmtType: ast.RoleStmtType(n.StmtType),
		Role:     n.Role,
		Options:  convertSlice(n.Options),
	}
}

func convertCreateSchemaStmt(n *pg.CreateSchemaStmt) *ast.CreateSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSchemaStmt{
		Name:        n.Schemaname,
		Authrole:    convertRoleSpec(n.Authrole),
		SchemaElts:  convertSlice(n.SchemaElts),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateSeqStmt(n *pg.CreateSeqStmt) *ast.CreateSeqStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSeqStmt{
		Sequence:    convertRangeVar(n.Sequence),
		Options:     convertSlice(n.Options),
		OwnerId:     ast.Oid(n.OwnerId),
		ForIdentity: n.ForIdentity,
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStatsStmt(n *pg.CreateStatsStmt) *ast.CreateStatsStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateStatsStmt{
		Defnames:    convertSlice(n.Defnames),
		StatTypes:   convertSlice(n.StatTypes),
		Exprs:       convertSlice(n.Exprs),
		Relations:   convertSlice(n.Relations),
		IfNotExists: n.IfNotExists,
	}
}

func convertCreateStmt(n *pg.CreateStmt) *ast.CreateStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateStmt{
		Relation:       convertRangeVar(n.Relation),
		TableElts:      convertSlice(n.TableElts),
		InhRelations:   convertSlice(n.InhRelations),
		Partbound:      convertPartitionBoundSpec(n.Partbound),
		Partspec:       convertPartitionSpec(n.Partspec),
		OfTypename:     convertTypeName(n.OfTypename),
		Constraints:    convertSlice(n.Constraints),
		Options:        convertSlice(n.Options),
		Oncommit:       ast.OnCommitAction(n.Oncommit),
		Tablespacename: n.Tablespacename,
		IfNotExists:    n.IfNotExists,
	}
}

func convertCreateSubscriptionStmt(n *pg.CreateSubscriptionStmt) *ast.CreateSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateSubscriptionStmt{
		Subname:     n.Subname,
		Conninfo:    n.Conninfo,
		Publication: convertSlice(n.Publication),
		Options:     convertSlice(n.Options),
	}
}

func convertCreateTableAsStmt(n *pg.CreateTableAsStmt) *ast.CreateTableAsStmt {
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

func convertCreateTableSpaceStmt(n *pg.CreateTableSpaceStmt) *ast.CreateTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		Owner:          convertRoleSpec(n.Owner),
		Location:       n.Location,
		Options:        convertSlice(n.Options),
	}
}

func convertCreateTransformStmt(n *pg.CreateTransformStmt) *ast.CreateTransformStmt {
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

func convertCreateTrigStmt(n *pg.CreateTrigStmt) *ast.CreateTrigStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateTrigStmt{
		Trigname:       n.Trigname,
		Relation:       convertRangeVar(n.Relation),
		Funcname:       convertSlice(n.Funcname),
		Args:           convertSlice(n.Args),
		Row:            n.Row,
		Timing:         n.Timing,
		Events:         n.Events,
		Columns:        convertSlice(n.Columns),
		WhenClause:     convertNode(n.WhenClause),
		Isconstraint:   n.Isconstraint,
		TransitionRels: convertSlice(n.TransitionRels),
		Deferrable:     n.Deferrable,
		Initdeferred:   n.Initdeferred,
		Constrrel:      convertRangeVar(n.Constrrel),
	}
}

func convertCreateUserMappingStmt(n *pg.CreateUserMappingStmt) *ast.CreateUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.CreateUserMappingStmt{
		User:        convertRoleSpec(n.User),
		Servername:  n.Servername,
		IfNotExists: n.IfNotExists,
		Options:     convertSlice(n.Options),
	}
}

func convertCreatedbStmt(n *pg.CreatedbStmt) *ast.CreatedbStmt {
	if n == nil {
		return nil
	}
	return &ast.CreatedbStmt{
		Dbname:  n.Dbname,
		Options: convertSlice(n.Options),
	}
}

func convertCurrentOfExpr(n *pg.CurrentOfExpr) *ast.CurrentOfExpr {
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

func convertDeallocateStmt(n *pg.DeallocateStmt) *ast.DeallocateStmt {
	if n == nil {
		return nil
	}
	return &ast.DeallocateStmt{
		Name: n.Name,
	}
}

func convertDeclareCursorStmt(n *pg.DeclareCursorStmt) *ast.DeclareCursorStmt {
	if n == nil {
		return nil
	}
	return &ast.DeclareCursorStmt{
		Portalname: n.Portalname,
		Options:    n.Options,
		Query:      convertNode(n.Query),
	}
}

func convertDefElem(n *pg.DefElem) *ast.DefElem {
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

func convertDefineStmt(n *pg.DefineStmt) *ast.DefineStmt {
	if n == nil {
		return nil
	}
	return &ast.DefineStmt{
		Kind:        ast.ObjectType(n.Kind),
		Oldstyle:    n.Oldstyle,
		Defnames:    convertSlice(n.Defnames),
		Args:        convertSlice(n.Args),
		Definition:  convertSlice(n.Definition),
		IfNotExists: n.IfNotExists,
	}
}

func convertDeleteStmt(n *pg.DeleteStmt) *ast.DeleteStmt {
	if n == nil {
		return nil
	}
	return &ast.DeleteStmt{
		Relation:      convertRangeVar(n.Relation),
		UsingClause:   convertSlice(n.UsingClause),
		WhereClause:   convertNode(n.WhereClause),
		ReturningList: convertSlice(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertDiscardStmt(n *pg.DiscardStmt) *ast.DiscardStmt {
	if n == nil {
		return nil
	}
	return &ast.DiscardStmt{
		Target: ast.DiscardMode(n.Target),
	}
}

func convertDoStmt(n *pg.DoStmt) *ast.DoStmt {
	if n == nil {
		return nil
	}
	return &ast.DoStmt{
		Args: convertSlice(n.Args),
	}
}

func convertDropOwnedStmt(n *pg.DropOwnedStmt) *ast.DropOwnedStmt {
	if n == nil {
		return nil
	}
	return &ast.DropOwnedStmt{
		Roles:    convertSlice(n.Roles),
		Behavior: ast.DropBehavior(n.Behavior),
	}
}

func convertDropRoleStmt(n *pg.DropRoleStmt) *ast.DropRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.DropRoleStmt{
		Roles:     convertSlice(n.Roles),
		MissingOk: n.MissingOk,
	}
}

func convertDropStmt(n *pg.DropStmt) *ast.DropStmt {
	if n == nil {
		return nil
	}
	return &ast.DropStmt{
		Objects:    convertSlice(n.Objects),
		RemoveType: ast.ObjectType(n.RemoveType),
		Behavior:   ast.DropBehavior(n.Behavior),
		MissingOk:  n.MissingOk,
		Concurrent: n.Concurrent,
	}
}

func convertDropSubscriptionStmt(n *pg.DropSubscriptionStmt) *ast.DropSubscriptionStmt {
	if n == nil {
		return nil
	}
	return &ast.DropSubscriptionStmt{
		Subname:   n.Subname,
		MissingOk: n.MissingOk,
		Behavior:  ast.DropBehavior(n.Behavior),
	}
}

func convertDropTableSpaceStmt(n *pg.DropTableSpaceStmt) *ast.DropTableSpaceStmt {
	if n == nil {
		return nil
	}
	return &ast.DropTableSpaceStmt{
		Tablespacename: n.Tablespacename,
		MissingOk:      n.MissingOk,
	}
}

func convertDropUserMappingStmt(n *pg.DropUserMappingStmt) *ast.DropUserMappingStmt {
	if n == nil {
		return nil
	}
	return &ast.DropUserMappingStmt{
		User:       convertRoleSpec(n.User),
		Servername: n.Servername,
		MissingOk:  n.MissingOk,
	}
}

func convertDropdbStmt(n *pg.DropdbStmt) *ast.DropdbStmt {
	if n == nil {
		return nil
	}
	return &ast.DropdbStmt{
		Dbname:    n.Dbname,
		MissingOk: n.MissingOk,
	}
}

func convertExecuteStmt(n *pg.ExecuteStmt) *ast.ExecuteStmt {
	if n == nil {
		return nil
	}
	return &ast.ExecuteStmt{
		Name:   n.Name,
		Params: convertSlice(n.Params),
	}
}

func convertExplainStmt(n *pg.ExplainStmt) *ast.ExplainStmt {
	if n == nil {
		return nil
	}
	return &ast.ExplainStmt{
		Query:   convertNode(n.Query),
		Options: convertSlice(n.Options),
	}
}

func convertExpr(n *pg.Expr) *ast.Expr {
	if n == nil {
		return nil
	}
	return &ast.Expr{}
}

func convertFetchStmt(n *pg.FetchStmt) *ast.FetchStmt {
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

func convertFieldSelect(n *pg.FieldSelect) *ast.FieldSelect {
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

func convertFieldStore(n *pg.FieldStore) *ast.FieldStore {
	if n == nil {
		return nil
	}
	return &ast.FieldStore{
		Xpr:        convertNode(n.Xpr),
		Arg:        convertNode(n.Arg),
		Newvals:    convertSlice(n.Newvals),
		Fieldnums:  convertSlice(n.Fieldnums),
		Resulttype: ast.Oid(n.Resulttype),
	}
}

func convertFloat(n *pg.Float) *ast.Float {
	if n == nil {
		return nil
	}
	return &ast.Float{
		Str: n.Str,
	}
}

func convertFromExpr(n *pg.FromExpr) *ast.FromExpr {
	if n == nil {
		return nil
	}
	return &ast.FromExpr{
		Fromlist: convertSlice(n.Fromlist),
		Quals:    convertNode(n.Quals),
	}
}

func convertFuncCall(n *pg.FuncCall) *ast.FuncCall {
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
		Funcname:       convertSlice(n.Funcname),
		Args:           convertSlice(n.Args),
		AggOrder:       convertSlice(n.AggOrder),
		AggFilter:      convertNode(n.AggFilter),
		AggWithinGroup: n.AggWithinGroup,
		AggStar:        n.AggStar,
		AggDistinct:    n.AggDistinct,
		FuncVariadic:   n.FuncVariadic,
		Over:           convertWindowDef(n.Over),
		Location:       n.Location,
	}
}

func convertFuncExpr(n *pg.FuncExpr) *ast.FuncExpr {
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
		Args:           convertSlice(n.Args),
		Location:       n.Location,
	}
}

func convertFunctionParameter(n *pg.FunctionParameter) *ast.FunctionParameter {
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

func convertGrantRoleStmt(n *pg.GrantRoleStmt) *ast.GrantRoleStmt {
	if n == nil {
		return nil
	}
	return &ast.GrantRoleStmt{
		GrantedRoles: convertSlice(n.GrantedRoles),
		GranteeRoles: convertSlice(n.GranteeRoles),
		IsGrant:      n.IsGrant,
		AdminOpt:     n.AdminOpt,
		Grantor:      convertRoleSpec(n.Grantor),
		Behavior:     ast.DropBehavior(n.Behavior),
	}
}

func convertGrantStmt(n *pg.GrantStmt) *ast.GrantStmt {
	if n == nil {
		return nil
	}
	return &ast.GrantStmt{
		IsGrant:     n.IsGrant,
		Targtype:    ast.GrantTargetType(n.Targtype),
		Objtype:     ast.GrantObjectType(n.Objtype),
		Objects:     convertSlice(n.Objects),
		Privileges:  convertSlice(n.Privileges),
		Grantees:    convertSlice(n.Grantees),
		GrantOption: n.GrantOption,
		Behavior:    ast.DropBehavior(n.Behavior),
	}
}

func convertGroupingFunc(n *pg.GroupingFunc) *ast.GroupingFunc {
	if n == nil {
		return nil
	}
	return &ast.GroupingFunc{
		Xpr:         convertNode(n.Xpr),
		Args:        convertSlice(n.Args),
		Refs:        convertSlice(n.Refs),
		Cols:        convertSlice(n.Cols),
		Agglevelsup: ast.Index(n.Agglevelsup),
		Location:    n.Location,
	}
}

func convertGroupingSet(n *pg.GroupingSet) *ast.GroupingSet {
	if n == nil {
		return nil
	}
	return &ast.GroupingSet{
		Kind:     ast.GroupingSetKind(n.Kind),
		Content:  convertSlice(n.Content),
		Location: n.Location,
	}
}

func convertImportForeignSchemaStmt(n *pg.ImportForeignSchemaStmt) *ast.ImportForeignSchemaStmt {
	if n == nil {
		return nil
	}
	return &ast.ImportForeignSchemaStmt{
		ServerName:   n.ServerName,
		RemoteSchema: n.RemoteSchema,
		LocalSchema:  n.LocalSchema,
		ListType:     ast.ImportForeignSchemaType(n.ListType),
		TableList:    convertSlice(n.TableList),
		Options:      convertSlice(n.Options),
	}
}

func convertIndexElem(n *pg.IndexElem) *ast.IndexElem {
	if n == nil {
		return nil
	}
	return &ast.IndexElem{
		Name:          n.Name,
		Expr:          convertNode(n.Expr),
		Indexcolname:  n.Indexcolname,
		Collation:     convertSlice(n.Collation),
		Opclass:       convertSlice(n.Opclass),
		Ordering:      ast.SortByDir(n.Ordering),
		NullsOrdering: ast.SortByNulls(n.NullsOrdering),
	}
}

func convertIndexStmt(n *pg.IndexStmt) *ast.IndexStmt {
	if n == nil {
		return nil
	}
	return &ast.IndexStmt{
		Idxname:        n.Idxname,
		Relation:       convertRangeVar(n.Relation),
		AccessMethod:   n.AccessMethod,
		TableSpace:     n.TableSpace,
		IndexParams:    convertSlice(n.IndexParams),
		Options:        convertSlice(n.Options),
		WhereClause:    convertNode(n.WhereClause),
		ExcludeOpNames: convertSlice(n.ExcludeOpNames),
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

func convertInferClause(n *pg.InferClause) *ast.InferClause {
	if n == nil {
		return nil
	}
	return &ast.InferClause{
		IndexElems:  convertSlice(n.IndexElems),
		WhereClause: convertNode(n.WhereClause),
		Conname:     n.Conname,
		Location:    n.Location,
	}
}

func convertInferenceElem(n *pg.InferenceElem) *ast.InferenceElem {
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

func convertInlineCodeBlock(n *pg.InlineCodeBlock) *ast.InlineCodeBlock {
	if n == nil {
		return nil
	}
	return &ast.InlineCodeBlock{
		SourceText:    n.SourceText,
		LangOid:       ast.Oid(n.LangOid),
		LangIsTrusted: n.LangIsTrusted,
	}
}

func convertInsertStmt(n *pg.InsertStmt) *ast.InsertStmt {
	if n == nil {
		return nil
	}
	return &ast.InsertStmt{
		Relation:         convertRangeVar(n.Relation),
		Cols:             convertSlice(n.Cols),
		SelectStmt:       convertNode(n.SelectStmt),
		OnConflictClause: convertOnConflictClause(n.OnConflictClause),
		ReturningList:    convertSlice(n.ReturningList),
		WithClause:       convertWithClause(n.WithClause),
		Override:         ast.OverridingKind(n.Override),
	}
}

func convertInteger(n *pg.Integer) *ast.Integer {
	if n == nil {
		return nil
	}
	return &ast.Integer{
		Ival: n.Ival,
	}
}

func convertIntoClause(n *pg.IntoClause) *ast.IntoClause {
	if n == nil {
		return nil
	}
	return &ast.IntoClause{
		Rel:            convertRangeVar(n.Rel),
		ColNames:       convertSlice(n.ColNames),
		Options:        convertSlice(n.Options),
		OnCommit:       ast.OnCommitAction(n.OnCommit),
		TableSpaceName: n.TableSpaceName,
		ViewQuery:      convertNode(n.ViewQuery),
		SkipData:       n.SkipData,
	}
}

func convertJoinExpr(n *pg.JoinExpr) *ast.JoinExpr {
	if n == nil {
		return nil
	}
	return &ast.JoinExpr{
		Jointype:    ast.JoinType(n.Jointype),
		IsNatural:   n.IsNatural,
		Larg:        convertNode(n.Larg),
		Rarg:        convertNode(n.Rarg),
		UsingClause: convertSlice(n.UsingClause),
		Quals:       convertNode(n.Quals),
		Alias:       convertAlias(n.Alias),
		Rtindex:     n.Rtindex,
	}
}

func convertSliceenStmt(n *pg.ListenStmt) *ast.ListenStmt {
	if n == nil {
		return nil
	}
	return &ast.ListenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertLoadStmt(n *pg.LoadStmt) *ast.LoadStmt {
	if n == nil {
		return nil
	}
	return &ast.LoadStmt{
		Filename: n.Filename,
	}
}

func convertLockStmt(n *pg.LockStmt) *ast.LockStmt {
	if n == nil {
		return nil
	}
	return &ast.LockStmt{
		Relations: convertSlice(n.Relations),
		Mode:      n.Mode,
		Nowait:    n.Nowait,
	}
}

func convertLockingClause(n *pg.LockingClause) *ast.LockingClause {
	if n == nil {
		return nil
	}
	return &ast.LockingClause{
		LockedRels: convertSlice(n.LockedRels),
		Strength:   ast.LockClauseStrength(n.Strength),
		WaitPolicy: ast.LockWaitPolicy(n.WaitPolicy),
	}
}

func convertMinMaxExpr(n *pg.MinMaxExpr) *ast.MinMaxExpr {
	if n == nil {
		return nil
	}
	return &ast.MinMaxExpr{
		Xpr:          convertNode(n.Xpr),
		Minmaxtype:   ast.Oid(n.Minmaxtype),
		Minmaxcollid: ast.Oid(n.Minmaxcollid),
		Inputcollid:  ast.Oid(n.Inputcollid),
		Op:           ast.MinMaxOp(n.Op),
		Args:         convertSlice(n.Args),
		Location:     n.Location,
	}
}

func convertMultiAssignRef(n *pg.MultiAssignRef) *ast.MultiAssignRef {
	if n == nil {
		return nil
	}
	return &ast.MultiAssignRef{
		Source:   convertNode(n.Source),
		Colno:    n.Colno,
		Ncolumns: n.Ncolumns,
	}
}

func convertNamedArgExpr(n *pg.NamedArgExpr) *ast.NamedArgExpr {
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

func convertNextValueExpr(n *pg.NextValueExpr) *ast.NextValueExpr {
	if n == nil {
		return nil
	}
	return &ast.NextValueExpr{
		Xpr:    convertNode(n.Xpr),
		Seqid:  ast.Oid(n.Seqid),
		TypeId: ast.Oid(n.TypeId),
	}
}

func convertNotifyStmt(n *pg.NotifyStmt) *ast.NotifyStmt {
	if n == nil {
		return nil
	}
	return &ast.NotifyStmt{
		Conditionname: n.Conditionname,
		Payload:       n.Payload,
	}
}

func convertNull(n *pg.Null) *ast.Null {
	if n == nil {
		return nil
	}
	return &ast.Null{}
}

func convertNullTest(n *pg.NullTest) *ast.NullTest {
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

func convertObjectWithArgs(n *pg.ObjectWithArgs) *ast.ObjectWithArgs {
	if n == nil {
		return nil
	}
	return &ast.ObjectWithArgs{
		Objname:         convertSlice(n.Objname),
		Objargs:         convertSlice(n.Objargs),
		ArgsUnspecified: n.ArgsUnspecified,
	}
}

func convertOnConflictClause(n *pg.OnConflictClause) *ast.OnConflictClause {
	if n == nil {
		return nil
	}
	return &ast.OnConflictClause{
		Action:      ast.OnConflictAction(n.Action),
		Infer:       convertInferClause(n.Infer),
		TargetList:  convertSlice(n.TargetList),
		WhereClause: convertNode(n.WhereClause),
		Location:    n.Location,
	}
}

func convertOnConflictExpr(n *pg.OnConflictExpr) *ast.OnConflictExpr {
	if n == nil {
		return nil
	}
	return &ast.OnConflictExpr{
		Action:          ast.OnConflictAction(n.Action),
		ArbiterElems:    convertSlice(n.ArbiterElems),
		ArbiterWhere:    convertNode(n.ArbiterWhere),
		Constraint:      ast.Oid(n.Constraint),
		OnConflictSet:   convertSlice(n.OnConflictSet),
		OnConflictWhere: convertNode(n.OnConflictWhere),
		ExclRelIndex:    n.ExclRelIndex,
		ExclRelTlist:    convertSlice(n.ExclRelTlist),
	}
}

func convertOpExpr(n *pg.OpExpr) *ast.OpExpr {
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
		Args:         convertSlice(n.Args),
		Location:     n.Location,
	}
}

func convertParam(n *pg.Param) *ast.Param {
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

func convertParamExecData(n *pg.ParamExecData) *ast.ParamExecData {
	if n == nil {
		return nil
	}
	return &ast.ParamExecData{
		ExecPlan: &ast.TODO{},
		Value:    ast.Datum(n.Value),
		Isnull:   n.Isnull,
	}
}

func convertParamExternData(n *pg.ParamExternData) *ast.ParamExternData {
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

func convertParamListInfoData(n *pg.ParamListInfoData) *ast.ParamListInfoData {
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

func convertParamRef(n *pg.ParamRef) *ast.ParamRef {
	if n == nil {
		return nil
	}
	return &ast.ParamRef{
		Number:   n.Number,
		Location: n.Location,
	}
}

func convertPartitionBoundSpec(n *pg.PartitionBoundSpec) *ast.PartitionBoundSpec {
	if n == nil {
		return nil
	}
	return &ast.PartitionBoundSpec{
		Strategy:    n.Strategy,
		Listdatums:  convertSlice(n.Listdatums),
		Lowerdatums: convertSlice(n.Lowerdatums),
		Upperdatums: convertSlice(n.Upperdatums),
		Location:    n.Location,
	}
}

func convertPartitionCmd(n *pg.PartitionCmd) *ast.PartitionCmd {
	if n == nil {
		return nil
	}
	return &ast.PartitionCmd{
		Name:  convertRangeVar(n.Name),
		Bound: convertPartitionBoundSpec(n.Bound),
	}
}

func convertPartitionElem(n *pg.PartitionElem) *ast.PartitionElem {
	if n == nil {
		return nil
	}
	return &ast.PartitionElem{
		Name:      n.Name,
		Expr:      convertNode(n.Expr),
		Collation: convertSlice(n.Collation),
		Opclass:   convertSlice(n.Opclass),
		Location:  n.Location,
	}
}

func convertPartitionRangeDatum(n *pg.PartitionRangeDatum) *ast.PartitionRangeDatum {
	if n == nil {
		return nil
	}
	return &ast.PartitionRangeDatum{
		Kind:     ast.PartitionRangeDatumKind(n.Kind),
		Value:    convertNode(n.Value),
		Location: n.Location,
	}
}

func convertPartitionSpec(n *pg.PartitionSpec) *ast.PartitionSpec {
	if n == nil {
		return nil
	}
	return &ast.PartitionSpec{
		Strategy:   n.Strategy,
		PartParams: convertSlice(n.PartParams),
		Location:   n.Location,
	}
}

func convertPrepareStmt(n *pg.PrepareStmt) *ast.PrepareStmt {
	if n == nil {
		return nil
	}
	return &ast.PrepareStmt{
		Name:     n.Name,
		Argtypes: convertSlice(n.Argtypes),
		Query:    convertNode(n.Query),
	}
}

func convertQuery(n *pg.Query) *ast.Query {
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
		CteList:          convertSlice(n.CteList),
		Rtable:           convertSlice(n.Rtable),
		Jointree:         convertFromExpr(n.Jointree),
		TargetList:       convertSlice(n.TargetList),
		Override:         ast.OverridingKind(n.Override),
		OnConflict:       convertOnConflictExpr(n.OnConflict),
		ReturningList:    convertSlice(n.ReturningList),
		GroupClause:      convertSlice(n.GroupClause),
		GroupingSets:     convertSlice(n.GroupingSets),
		HavingQual:       convertNode(n.HavingQual),
		WindowClause:     convertSlice(n.WindowClause),
		DistinctClause:   convertSlice(n.DistinctClause),
		SortClause:       convertSlice(n.SortClause),
		LimitOffset:      convertNode(n.LimitOffset),
		LimitCount:       convertNode(n.LimitCount),
		RowMarks:         convertSlice(n.RowMarks),
		SetOperations:    convertNode(n.SetOperations),
		ConstraintDeps:   convertSlice(n.ConstraintDeps),
		WithCheckOptions: convertSlice(n.WithCheckOptions),
		StmtLocation:     n.StmtLocation,
		StmtLen:          n.StmtLen,
	}
}

func convertRangeFunction(n *pg.RangeFunction) *ast.RangeFunction {
	if n == nil {
		return nil
	}
	return &ast.RangeFunction{
		Lateral:    n.Lateral,
		Ordinality: n.Ordinality,
		IsRowsfrom: n.IsRowsfrom,
		Functions:  convertSlice(n.Functions),
		Alias:      convertAlias(n.Alias),
		Coldeflist: convertSlice(n.Coldeflist),
	}
}

func convertRangeSubselect(n *pg.RangeSubselect) *ast.RangeSubselect {
	if n == nil {
		return nil
	}
	return &ast.RangeSubselect{
		Lateral:  n.Lateral,
		Subquery: convertNode(n.Subquery),
		Alias:    convertAlias(n.Alias),
	}
}

func convertRangeTableFunc(n *pg.RangeTableFunc) *ast.RangeTableFunc {
	if n == nil {
		return nil
	}
	return &ast.RangeTableFunc{
		Lateral:    n.Lateral,
		Docexpr:    convertNode(n.Docexpr),
		Rowexpr:    convertNode(n.Rowexpr),
		Namespaces: convertSlice(n.Namespaces),
		Columns:    convertSlice(n.Columns),
		Alias:      convertAlias(n.Alias),
		Location:   n.Location,
	}
}

func convertRangeTableFuncCol(n *pg.RangeTableFuncCol) *ast.RangeTableFuncCol {
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

func convertRangeTableSample(n *pg.RangeTableSample) *ast.RangeTableSample {
	if n == nil {
		return nil
	}
	return &ast.RangeTableSample{
		Relation:   convertNode(n.Relation),
		Method:     convertSlice(n.Method),
		Args:       convertSlice(n.Args),
		Repeatable: convertNode(n.Repeatable),
		Location:   n.Location,
	}
}

func convertRangeTblEntry(n *pg.RangeTblEntry) *ast.RangeTblEntry {
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
		Joinaliasvars:   convertSlice(n.Joinaliasvars),
		Functions:       convertSlice(n.Functions),
		Funcordinality:  n.Funcordinality,
		Tablefunc:       convertTableFunc(n.Tablefunc),
		ValuesLists:     convertSlice(n.ValuesLists),
		Ctename:         n.Ctename,
		Ctelevelsup:     ast.Index(n.Ctelevelsup),
		SelfReference:   n.SelfReference,
		Coltypes:        convertSlice(n.Coltypes),
		Coltypmods:      convertSlice(n.Coltypmods),
		Colcollations:   convertSlice(n.Colcollations),
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
		SecurityQuals:   convertSlice(n.SecurityQuals),
	}
}

func convertRangeTblFunction(n *pg.RangeTblFunction) *ast.RangeTblFunction {
	if n == nil {
		return nil
	}
	return &ast.RangeTblFunction{
		Funcexpr:          convertNode(n.Funcexpr),
		Funccolcount:      n.Funccolcount,
		Funccolnames:      convertSlice(n.Funccolnames),
		Funccoltypes:      convertSlice(n.Funccoltypes),
		Funccoltypmods:    convertSlice(n.Funccoltypmods),
		Funccolcollations: convertSlice(n.Funccolcollations),
		Funcparams:        n.Funcparams,
	}
}

func convertRangeTblRef(n *pg.RangeTblRef) *ast.RangeTblRef {
	if n == nil {
		return nil
	}
	return &ast.RangeTblRef{
		Rtindex: n.Rtindex,
	}
}

func convertRangeVar(n *pg.RangeVar) *ast.RangeVar {
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

func convertRawStmt(n *pg.RawStmt) *ast.RawStmt {
	if n == nil {
		return nil
	}
	return &ast.RawStmt{
		Stmt:         convertNode(n.Stmt),
		StmtLocation: n.StmtLocation,
		StmtLen:      n.StmtLen,
	}
}

func convertReassignOwnedStmt(n *pg.ReassignOwnedStmt) *ast.ReassignOwnedStmt {
	if n == nil {
		return nil
	}
	return &ast.ReassignOwnedStmt{
		Roles:   convertSlice(n.Roles),
		Newrole: convertRoleSpec(n.Newrole),
	}
}

func convertRefreshMatViewStmt(n *pg.RefreshMatViewStmt) *ast.RefreshMatViewStmt {
	if n == nil {
		return nil
	}
	return &ast.RefreshMatViewStmt{
		Concurrent: n.Concurrent,
		SkipData:   n.SkipData,
		Relation:   convertRangeVar(n.Relation),
	}
}

func convertReindexStmt(n *pg.ReindexStmt) *ast.ReindexStmt {
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

func convertRelabelType(n *pg.RelabelType) *ast.RelabelType {
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

func convertRenameStmt(n *pg.RenameStmt) *ast.RenameStmt {
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

func convertReplicaIdentityStmt(n *pg.ReplicaIdentityStmt) *ast.ReplicaIdentityStmt {
	if n == nil {
		return nil
	}
	return &ast.ReplicaIdentityStmt{
		IdentityType: n.IdentityType,
		Name:         n.Name,
	}
}

func convertResTarget(n *pg.ResTarget) *ast.ResTarget {
	if n == nil {
		return nil
	}
	return &ast.ResTarget{
		Name:        n.Name,
		Indirection: convertSlice(n.Indirection),
		Val:         convertNode(n.Val),
		Location:    n.Location,
	}
}

func convertRoleSpec(n *pg.RoleSpec) *ast.RoleSpec {
	if n == nil {
		return nil
	}
	return &ast.RoleSpec{
		Roletype: ast.RoleSpecType(n.Roletype),
		Rolename: n.Rolename,
		Location: n.Location,
	}
}

func convertRowCompareExpr(n *pg.RowCompareExpr) *ast.RowCompareExpr {
	if n == nil {
		return nil
	}
	return &ast.RowCompareExpr{
		Xpr:          convertNode(n.Xpr),
		Rctype:       ast.RowCompareType(n.Rctype),
		Opnos:        convertSlice(n.Opnos),
		Opfamilies:   convertSlice(n.Opfamilies),
		Inputcollids: convertSlice(n.Inputcollids),
		Largs:        convertSlice(n.Largs),
		Rargs:        convertSlice(n.Rargs),
	}
}

func convertRowExpr(n *pg.RowExpr) *ast.RowExpr {
	if n == nil {
		return nil
	}
	return &ast.RowExpr{
		Xpr:       convertNode(n.Xpr),
		Args:      convertSlice(n.Args),
		RowTypeid: ast.Oid(n.RowTypeid),
		RowFormat: ast.CoercionForm(n.RowFormat),
		Colnames:  convertSlice(n.Colnames),
		Location:  n.Location,
	}
}

func convertRowMarkClause(n *pg.RowMarkClause) *ast.RowMarkClause {
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

func convertRuleStmt(n *pg.RuleStmt) *ast.RuleStmt {
	if n == nil {
		return nil
	}
	return &ast.RuleStmt{
		Relation:    convertRangeVar(n.Relation),
		Rulename:    n.Rulename,
		WhereClause: convertNode(n.WhereClause),
		Event:       ast.CmdType(n.Event),
		Instead:     n.Instead,
		Actions:     convertSlice(n.Actions),
		Replace:     n.Replace,
	}
}

func convertSQLValueFunction(n *pg.SQLValueFunction) *ast.SQLValueFunction {
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

func convertScalarArrayOpExpr(n *pg.ScalarArrayOpExpr) *ast.ScalarArrayOpExpr {
	if n == nil {
		return nil
	}
	return &ast.ScalarArrayOpExpr{
		Xpr:         convertNode(n.Xpr),
		Opno:        ast.Oid(n.Opno),
		Opfuncid:    ast.Oid(n.Opfuncid),
		UseOr:       n.UseOr,
		Inputcollid: ast.Oid(n.Inputcollid),
		Args:        convertSlice(n.Args),
		Location:    n.Location,
	}
}

func convertSecLabelStmt(n *pg.SecLabelStmt) *ast.SecLabelStmt {
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

func convertSelectStmt(n *pg.SelectStmt) *ast.SelectStmt {
	if n == nil {
		return nil
	}
	return &ast.SelectStmt{
		DistinctClause: convertSlice(n.DistinctClause),
		IntoClause:     convertIntoClause(n.IntoClause),
		TargetList:     convertSlice(n.TargetList),
		FromClause:     convertSlice(n.FromClause),
		WhereClause:    convertNode(n.WhereClause),
		GroupClause:    convertSlice(n.GroupClause),
		HavingClause:   convertNode(n.HavingClause),
		WindowClause:   convertSlice(n.WindowClause),
		ValuesLists:    convertValuesList(n.ValuesLists),
		SortClause:     convertSlice(n.SortClause),
		LimitOffset:    convertNode(n.LimitOffset),
		LimitCount:     convertNode(n.LimitCount),
		LockingClause:  convertSlice(n.LockingClause),
		WithClause:     convertWithClause(n.WithClause),
		Op:             ast.SetOperation(n.Op),
		All:            n.All,
		Larg:           convertSelectStmt(n.Larg),
		Rarg:           convertSelectStmt(n.Rarg),
	}
}

func convertSetOperationStmt(n *pg.SetOperationStmt) *ast.SetOperationStmt {
	if n == nil {
		return nil
	}
	return &ast.SetOperationStmt{
		Op:            ast.SetOperation(n.Op),
		All:           n.All,
		Larg:          convertNode(n.Larg),
		Rarg:          convertNode(n.Rarg),
		ColTypes:      convertSlice(n.ColTypes),
		ColTypmods:    convertSlice(n.ColTypmods),
		ColCollations: convertSlice(n.ColCollations),
		GroupClauses:  convertSlice(n.GroupClauses),
	}
}

func convertSetToDefault(n *pg.SetToDefault) *ast.SetToDefault {
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

func convertSortBy(n *pg.SortBy) *ast.SortBy {
	if n == nil {
		return nil
	}
	return &ast.SortBy{
		Node:        convertNode(n.Node),
		SortbyDir:   ast.SortByDir(n.SortbyDir),
		SortbyNulls: ast.SortByNulls(n.SortbyNulls),
		UseOp:       convertSlice(n.UseOp),
		Location:    n.Location,
	}
}

func convertSortGroupClause(n *pg.SortGroupClause) *ast.SortGroupClause {
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

func convertString(n *pg.String) *ast.String {
	if n == nil {
		return nil
	}
	return &ast.String{
		Str: n.Str,
	}
}

func convertSubLink(n *pg.SubLink) *ast.SubLink {
	if n == nil {
		return nil
	}
	return &ast.SubLink{
		Xpr:         convertNode(n.Xpr),
		SubLinkType: ast.SubLinkType(n.SubLinkType),
		SubLinkId:   n.SubLinkId,
		Testexpr:    convertNode(n.Testexpr),
		OperName:    convertSlice(n.OperName),
		Subselect:   convertNode(n.Subselect),
		Location:    n.Location,
	}
}

func convertSubPlan(n *pg.SubPlan) *ast.SubPlan {
	if n == nil {
		return nil
	}
	return &ast.SubPlan{
		Xpr:               convertNode(n.Xpr),
		SubLinkType:       ast.SubLinkType(n.SubLinkType),
		Testexpr:          convertNode(n.Testexpr),
		ParamIds:          convertSlice(n.ParamIds),
		PlanId:            n.PlanId,
		PlanName:          n.PlanName,
		FirstColType:      ast.Oid(n.FirstColType),
		FirstColTypmod:    n.FirstColTypmod,
		FirstColCollation: ast.Oid(n.FirstColCollation),
		UseHashTable:      n.UseHashTable,
		UnknownEqFalse:    n.UnknownEqFalse,
		ParallelSafe:      n.ParallelSafe,
		SetParam:          convertSlice(n.SetParam),
		ParParam:          convertSlice(n.ParParam),
		Args:              convertSlice(n.Args),
		StartupCost:       ast.Cost(n.StartupCost),
		PerCallCost:       ast.Cost(n.PerCallCost),
	}
}

func convertTableFunc(n *pg.TableFunc) *ast.TableFunc {
	if n == nil {
		return nil
	}
	return &ast.TableFunc{
		NsUris:        convertSlice(n.NsUris),
		NsNames:       convertSlice(n.NsNames),
		Docexpr:       convertNode(n.Docexpr),
		Rowexpr:       convertNode(n.Rowexpr),
		Colnames:      convertSlice(n.Colnames),
		Coltypes:      convertSlice(n.Coltypes),
		Coltypmods:    convertSlice(n.Coltypmods),
		Colcollations: convertSlice(n.Colcollations),
		Colexprs:      convertSlice(n.Colexprs),
		Coldefexprs:   convertSlice(n.Coldefexprs),
		Notnulls:      n.Notnulls,
		Ordinalitycol: n.Ordinalitycol,
		Location:      n.Location,
	}
}

func convertTableLikeClause(n *pg.TableLikeClause) *ast.TableLikeClause {
	if n == nil {
		return nil
	}
	return &ast.TableLikeClause{
		Relation: convertRangeVar(n.Relation),
		Options:  n.Options,
	}
}

func convertTableSampleClause(n *pg.TableSampleClause) *ast.TableSampleClause {
	if n == nil {
		return nil
	}
	return &ast.TableSampleClause{
		Tsmhandler: ast.Oid(n.Tsmhandler),
		Args:       convertSlice(n.Args),
		Repeatable: convertNode(n.Repeatable),
	}
}

func convertTargetEntry(n *pg.TargetEntry) *ast.TargetEntry {
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

func convertTransactionStmt(n *pg.TransactionStmt) *ast.TransactionStmt {
	if n == nil {
		return nil
	}
	return &ast.TransactionStmt{
		Kind:    ast.TransactionStmtKind(n.Kind),
		Options: convertSlice(n.Options),
		Gid:     n.Gid,
	}
}

func convertTriggerTransition(n *pg.TriggerTransition) *ast.TriggerTransition {
	if n == nil {
		return nil
	}
	return &ast.TriggerTransition{
		Name:    n.Name,
		IsNew:   n.IsNew,
		IsTable: n.IsTable,
	}
}

func convertTruncateStmt(n *pg.TruncateStmt) *ast.TruncateStmt {
	if n == nil {
		return nil
	}
	return &ast.TruncateStmt{
		Relations:   convertSlice(n.Relations),
		RestartSeqs: n.RestartSeqs,
		Behavior:    ast.DropBehavior(n.Behavior),
	}
}

func convertTypeCast(n *pg.TypeCast) *ast.TypeCast {
	if n == nil {
		return nil
	}
	return &ast.TypeCast{
		Arg:      convertNode(n.Arg),
		TypeName: convertTypeName(n.TypeName),
		Location: n.Location,
	}
}

func convertTypeName(n *pg.TypeName) *ast.TypeName {
	if n == nil {
		return nil
	}
	return &ast.TypeName{
		Names:       convertSlice(n.Names),
		TypeOid:     ast.Oid(n.TypeOid),
		Setof:       n.Setof,
		PctType:     n.PctType,
		Typmods:     convertSlice(n.Typmods),
		Typemod:     n.Typemod,
		ArrayBounds: convertSlice(n.ArrayBounds),
		Location:    n.Location,
	}
}

func convertUnlistenStmt(n *pg.UnlistenStmt) *ast.UnlistenStmt {
	if n == nil {
		return nil
	}
	return &ast.UnlistenStmt{
		Conditionname: n.Conditionname,
	}
}

func convertUpdateStmt(n *pg.UpdateStmt) *ast.UpdateStmt {
	if n == nil {
		return nil
	}
	return &ast.UpdateStmt{
		Relation:      convertRangeVar(n.Relation),
		TargetList:    convertSlice(n.TargetList),
		WhereClause:   convertNode(n.WhereClause),
		FromClause:    convertSlice(n.FromClause),
		ReturningList: convertSlice(n.ReturningList),
		WithClause:    convertWithClause(n.WithClause),
	}
}

func convertVacuumStmt(n *pg.VacuumStmt) *ast.VacuumStmt {
	if n == nil {
		return nil
	}
	return &ast.VacuumStmt{
		Options:  n.Options,
		Relation: convertRangeVar(n.Relation),
		VaCols:   convertSlice(n.VaCols),
	}
}

func convertVar(n *pg.Var) *ast.Var {
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

func convertVariableSetStmt(n *pg.VariableSetStmt) *ast.VariableSetStmt {
	if n == nil {
		return nil
	}
	return &ast.VariableSetStmt{
		Kind:    ast.VariableSetKind(n.Kind),
		Name:    n.Name,
		Args:    convertSlice(n.Args),
		IsLocal: n.IsLocal,
	}
}

func convertVariableShowStmt(n *pg.VariableShowStmt) *ast.VariableShowStmt {
	if n == nil {
		return nil
	}
	return &ast.VariableShowStmt{
		Name: n.Name,
	}
}

func convertViewStmt(n *pg.ViewStmt) *ast.ViewStmt {
	if n == nil {
		return nil
	}
	return &ast.ViewStmt{
		View:            convertRangeVar(n.View),
		Aliases:         convertSlice(n.Aliases),
		Query:           convertNode(n.Query),
		Replace:         n.Replace,
		Options:         convertSlice(n.Options),
		WithCheckOption: ast.ViewCheckOption(n.WithCheckOption),
	}
}

func convertWindowClause(n *pg.WindowClause) *ast.WindowClause {
	if n == nil {
		return nil
	}
	return &ast.WindowClause{
		Name:            n.Name,
		Refname:         n.Refname,
		PartitionClause: convertSlice(n.PartitionClause),
		OrderClause:     convertSlice(n.OrderClause),
		FrameOptions:    n.FrameOptions,
		StartOffset:     convertNode(n.StartOffset),
		EndOffset:       convertNode(n.EndOffset),
		Winref:          ast.Index(n.Winref),
		CopiedOrder:     n.CopiedOrder,
	}
}

func convertWindowDef(n *pg.WindowDef) *ast.WindowDef {
	if n == nil {
		return nil
	}
	return &ast.WindowDef{
		Name:            n.Name,
		Refname:         n.Refname,
		PartitionClause: convertSlice(n.PartitionClause),
		OrderClause:     convertSlice(n.OrderClause),
		FrameOptions:    n.FrameOptions,
		StartOffset:     convertNode(n.StartOffset),
		EndOffset:       convertNode(n.EndOffset),
		Location:        n.Location,
	}
}

func convertWindowFunc(n *pg.WindowFunc) *ast.WindowFunc {
	if n == nil {
		return nil
	}
	return &ast.WindowFunc{
		Xpr:         convertNode(n.Xpr),
		Winfnoid:    ast.Oid(n.Winfnoid),
		Wintype:     ast.Oid(n.Wintype),
		Wincollid:   ast.Oid(n.Wincollid),
		Inputcollid: ast.Oid(n.Inputcollid),
		Args:        convertSlice(n.Args),
		Aggfilter:   convertNode(n.Aggfilter),
		Winref:      ast.Index(n.Winref),
		Winstar:     n.Winstar,
		Winagg:      n.Winagg,
		Location:    n.Location,
	}
}

func convertWithCheckOption(n *pg.WithCheckOption) *ast.WithCheckOption {
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

func convertWithClause(n *pg.WithClause) *ast.WithClause {
	if n == nil {
		return nil
	}
	return &ast.WithClause{
		Ctes:      convertSlice(n.Ctes),
		Recursive: n.Recursive,
		Location:  n.Location,
	}
}

func convertXmlExpr(n *pg.XmlExpr) *ast.XmlExpr {
	if n == nil {
		return nil
	}
	return &ast.XmlExpr{
		Xpr:       convertNode(n.Xpr),
		Op:        ast.XmlExprOp(n.Op),
		Name:      n.Name,
		NamedArgs: convertSlice(n.NamedArgs),
		ArgNames:  convertSlice(n.ArgNames),
		Args:      convertSlice(n.Args),
		Xmloption: ast.XmlOptionType(n.Xmloption),
		Type:      ast.Oid(n.Type),
		Typmod:    n.Typmod,
		Location:  n.Location,
	}
}

func convertXmlSerialize(n *pg.XmlSerialize) *ast.XmlSerialize {
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

func convertNode(node *pg.Node) ast.Node {
	switch n := node.Node.(type) {

	case *pg.Node_AArrayExpr:
		return convertA_ArrayExpr(n.AArrayExpr)

	case pg.A_Const:
		return convertA_Const(&n)

	case pg.A_Expr:
		return convertA_Expr(&n)

	case pg.A_Indices:
		return convertA_Indices(&n)

	case pg.A_Indirection:
		return convertA_Indirection(&n)

	case pg.A_Star:
		return convertA_Star(&n)

	case pg.AccessPriv:
		return convertAccessPriv(&n)

	case pg.Aggref:
		return convertAggref(&n)

	case pg.Alias:
		return convertAlias(&n)

	case pg.AlterCollationStmt:
		return convertAlterCollationStmt(&n)

	case pg.AlterDatabaseSetStmt:
		return convertAlterDatabaseSetStmt(&n)

	case pg.AlterDatabaseStmt:
		return convertAlterDatabaseStmt(&n)

	case pg.AlterDefaultPrivilegesStmt:
		return convertAlterDefaultPrivilegesStmt(&n)

	case pg.AlterDomainStmt:
		return convertAlterDomainStmt(&n)

	case pg.AlterEnumStmt:
		return convertAlterEnumStmt(&n)

	case pg.AlterEventTrigStmt:
		return convertAlterEventTrigStmt(&n)

	case pg.AlterExtensionContentsStmt:
		return convertAlterExtensionContentsStmt(&n)

	case pg.AlterExtensionStmt:
		return convertAlterExtensionStmt(&n)

	case pg.AlterFdwStmt:
		return convertAlterFdwStmt(&n)

	case pg.AlterForeignServerStmt:
		return convertAlterForeignServerStmt(&n)

	case pg.AlterFunctionStmt:
		return convertAlterFunctionStmt(&n)

	case pg.AlterObjectDependsStmt:
		return convertAlterObjectDependsStmt(&n)

	case pg.AlterObjectSchemaStmt:
		return convertAlterObjectSchemaStmt(&n)

	case pg.AlterOpFamilyStmt:
		return convertAlterOpFamilyStmt(&n)

	case pg.AlterOperatorStmt:
		return convertAlterOperatorStmt(&n)

	case pg.AlterOwnerStmt:
		return convertAlterOwnerStmt(&n)

	case pg.AlterPolicyStmt:
		return convertAlterPolicyStmt(&n)

	case pg.AlterPublicationStmt:
		return convertAlterPublicationStmt(&n)

	case pg.AlterRoleSetStmt:
		return convertAlterRoleSetStmt(&n)

	case pg.AlterRoleStmt:
		return convertAlterRoleStmt(&n)

	case pg.AlterSeqStmt:
		return convertAlterSeqStmt(&n)

	case pg.AlterSubscriptionStmt:
		return convertAlterSubscriptionStmt(&n)

	case pg.AlterSystemStmt:
		return convertAlterSystemStmt(&n)

	case pg.AlterTSConfigurationStmt:
		return convertAlterTSConfigurationStmt(&n)

	case pg.AlterTSDictionaryStmt:
		return convertAlterTSDictionaryStmt(&n)

	case pg.AlterTableCmd:
		return convertAlterTableCmd(&n)

	case pg.AlterTableMoveAllStmt:
		return convertAlterTableMoveAllStmt(&n)

	case pg.AlterTableSpaceOptionsStmt:
		return convertAlterTableSpaceOptionsStmt(&n)

	case pg.AlterTableStmt:
		return convertAlterTableStmt(&n)

	case pg.AlterUserMappingStmt:
		return convertAlterUserMappingStmt(&n)

	case pg.AlternativeSubPlan:
		return convertAlternativeSubPlan(&n)

	case pg.ArrayCoerceExpr:
		return convertArrayCoerceExpr(&n)

	case pg.ArrayExpr:
		return convertArrayExpr(&n)

	case pg.ArrayRef:
		return convertArrayRef(&n)

	case pg.BitString:
		return convertBitString(&n)

	case pg.BlockIdData:
		return convertBlockIdData(&n)

	case pg.BoolExpr:
		return convertBoolExpr(&n)

	case pg.BooleanTest:
		return convertBooleanTest(&n)

	case pg.CaseExpr:
		return convertCaseExpr(&n)

	case pg.CaseTestExpr:
		return convertCaseTestExpr(&n)

	case pg.CaseWhen:
		return convertCaseWhen(&n)

	case pg.CheckPointStmt:
		return convertCheckPointStmt(&n)

	case pg.ClosePortalStmt:
		return convertClosePortalStmt(&n)

	case pg.ClusterStmt:
		return convertClusterStmt(&n)

	case pg.CoalesceExpr:
		return convertCoalesceExpr(&n)

	case pg.CoerceToDomain:
		return convertCoerceToDomain(&n)

	case pg.CoerceToDomainValue:
		return convertCoerceToDomainValue(&n)

	case pg.CoerceViaIO:
		return convertCoerceViaIO(&n)

	case pg.CollateClause:
		return convertCollateClause(&n)

	case pg.CollateExpr:
		return convertCollateExpr(&n)

	case pg.ColumnDef:
		return convertColumnDef(&n)

	case pg.ColumnRef:
		return convertColumnRef(&n)

	case pg.CommentStmt:
		return convertCommentStmt(&n)

	case pg.CommonTableExpr:
		return convertCommonTableExpr(&n)

	case pg.CompositeTypeStmt:
		return convertCompositeTypeStmt(&n)

	case pg.Const:
		return convertConst(&n)

	case pg.Constraint:
		return convertConstraint(&n)

	case pg.ConstraintsSetStmt:
		return convertConstraintsSetStmt(&n)

	case pg.ConvertRowtypeExpr:
		return convertConvertRowtypeExpr(&n)

	case pg.CopyStmt:
		return convertCopyStmt(&n)

	case pg.CreateAmStmt:
		return convertCreateAmStmt(&n)

	case pg.CreateCastStmt:
		return convertCreateCastStmt(&n)

	case pg.CreateConversionStmt:
		return convertCreateConversionStmt(&n)

	case pg.CreateDomainStmt:
		return convertCreateDomainStmt(&n)

	case pg.CreateEnumStmt:
		return convertCreateEnumStmt(&n)

	case pg.CreateEventTrigStmt:
		return convertCreateEventTrigStmt(&n)

	case pg.CreateExtensionStmt:
		return convertCreateExtensionStmt(&n)

	case pg.CreateFdwStmt:
		return convertCreateFdwStmt(&n)

	case pg.CreateForeignServerStmt:
		return convertCreateForeignServerStmt(&n)

	case pg.CreateForeignTableStmt:
		return convertCreateForeignTableStmt(&n)

	case pg.CreateFunctionStmt:
		return convertCreateFunctionStmt(&n)

	case pg.CreateOpClassItem:
		return convertCreateOpClassItem(&n)

	case pg.CreateOpClassStmt:
		return convertCreateOpClassStmt(&n)

	case pg.CreateOpFamilyStmt:
		return convertCreateOpFamilyStmt(&n)

	case pg.CreatePLangStmt:
		return convertCreatePLangStmt(&n)

	case pg.CreatePolicyStmt:
		return convertCreatePolicyStmt(&n)

	case pg.CreatePublicationStmt:
		return convertCreatePublicationStmt(&n)

	case pg.CreateRangeStmt:
		return convertCreateRangeStmt(&n)

	case pg.CreateRoleStmt:
		return convertCreateRoleStmt(&n)

	case pg.CreateSchemaStmt:
		return convertCreateSchemaStmt(&n)

	case pg.CreateSeqStmt:
		return convertCreateSeqStmt(&n)

	case pg.CreateStatsStmt:
		return convertCreateStatsStmt(&n)

	case pg.CreateStmt:
		return convertCreateStmt(&n)

	case pg.CreateSubscriptionStmt:
		return convertCreateSubscriptionStmt(&n)

	case pg.CreateTableAsStmt:
		return convertCreateTableAsStmt(&n)

	case pg.CreateTableSpaceStmt:
		return convertCreateTableSpaceStmt(&n)

	case pg.CreateTransformStmt:
		return convertCreateTransformStmt(&n)

	case pg.CreateTrigStmt:
		return convertCreateTrigStmt(&n)

	case pg.CreateUserMappingStmt:
		return convertCreateUserMappingStmt(&n)

	case pg.CreatedbStmt:
		return convertCreatedbStmt(&n)

	case pg.CurrentOfExpr:
		return convertCurrentOfExpr(&n)

	case pg.DeallocateStmt:
		return convertDeallocateStmt(&n)

	case pg.DeclareCursorStmt:
		return convertDeclareCursorStmt(&n)

	case pg.DefElem:
		return convertDefElem(&n)

	case pg.DefineStmt:
		return convertDefineStmt(&n)

	case pg.DeleteStmt:
		return convertDeleteStmt(&n)

	case pg.DiscardStmt:
		return convertDiscardStmt(&n)

	case pg.DoStmt:
		return convertDoStmt(&n)

	case pg.DropOwnedStmt:
		return convertDropOwnedStmt(&n)

	case pg.DropRoleStmt:
		return convertDropRoleStmt(&n)

	case pg.DropStmt:
		return convertDropStmt(&n)

	case pg.DropSubscriptionStmt:
		return convertDropSubscriptionStmt(&n)

	case pg.DropTableSpaceStmt:
		return convertDropTableSpaceStmt(&n)

	case pg.DropUserMappingStmt:
		return convertDropUserMappingStmt(&n)

	case pg.DropdbStmt:
		return convertDropdbStmt(&n)

	case pg.ExecuteStmt:
		return convertExecuteStmt(&n)

	case pg.ExplainStmt:
		return convertExplainStmt(&n)

	case pg.Expr:
		return convertExpr(&n)

	case pg.FetchStmt:
		return convertFetchStmt(&n)

	case pg.FieldSelect:
		return convertFieldSelect(&n)

	case pg.FieldStore:
		return convertFieldStore(&n)

	case pg.Float:
		return convertFloat(&n)

	case pg.FromExpr:
		return convertFromExpr(&n)

	case pg.FuncCall:
		return convertFuncCall(&n)

	case pg.FuncExpr:
		return convertFuncExpr(&n)

	case pg.FunctionParameter:
		return convertFunctionParameter(&n)

	case pg.GrantRoleStmt:
		return convertGrantRoleStmt(&n)

	case pg.GrantStmt:
		return convertGrantStmt(&n)

	case pg.GroupingFunc:
		return convertGroupingFunc(&n)

	case pg.GroupingSet:
		return convertGroupingSet(&n)

	case pg.ImportForeignSchemaStmt:
		return convertImportForeignSchemaStmt(&n)

	case pg.IndexElem:
		return convertIndexElem(&n)

	case pg.IndexStmt:
		return convertIndexStmt(&n)

	case pg.InferClause:
		return convertInferClause(&n)

	case pg.InferenceElem:
		return convertInferenceElem(&n)

	case pg.InlineCodeBlock:
		return convertInlineCodeBlock(&n)

	case pg.InsertStmt:
		return convertInsertStmt(&n)

	case pg.Integer:
		return convertInteger(&n)

	case pg.IntoClause:
		return convertIntoClause(&n)

	case pg.JoinExpr:
		return convertJoinExpr(&n)

	case pg.List:
		return convertSlice(n)

	case pg.ListenStmt:
		return convertSliceenStmt(&n)

	case pg.LoadStmt:
		return convertLoadStmt(&n)

	case pg.LockStmt:
		return convertLockStmt(&n)

	case pg.LockingClause:
		return convertLockingClause(&n)

	case pg.MinMaxExpr:
		return convertMinMaxExpr(&n)

	case pg.MultiAssignRef:
		return convertMultiAssignRef(&n)

	case pg.NamedArgExpr:
		return convertNamedArgExpr(&n)

	case pg.NextValueExpr:
		return convertNextValueExpr(&n)

	case pg.NotifyStmt:
		return convertNotifyStmt(&n)

	case pg.Null:
		return convertNull(&n)

	case pg.NullTest:
		return convertNullTest(&n)

	case pg.ObjectWithArgs:
		return convertObjectWithArgs(&n)

	case pg.OnConflictClause:
		return convertOnConflictClause(&n)

	case pg.OnConflictExpr:
		return convertOnConflictExpr(&n)

	case pg.OpExpr:
		return convertOpExpr(&n)

	case pg.Param:
		return convertParam(&n)

	case pg.ParamExecData:
		return convertParamExecData(&n)

	case pg.ParamExternData:
		return convertParamExternData(&n)

	case pg.ParamListInfoData:
		return convertParamListInfoData(&n)

	case pg.ParamRef:
		return convertParamRef(&n)

	case pg.PartitionBoundSpec:
		return convertPartitionBoundSpec(&n)

	case pg.PartitionCmd:
		return convertPartitionCmd(&n)

	case pg.PartitionElem:
		return convertPartitionElem(&n)

	case pg.PartitionRangeDatum:
		return convertPartitionRangeDatum(&n)

	case pg.PartitionSpec:
		return convertPartitionSpec(&n)

	case pg.PrepareStmt:
		return convertPrepareStmt(&n)

	case pg.Query:
		return convertQuery(&n)

	case pg.RangeFunction:
		return convertRangeFunction(&n)

	case pg.RangeSubselect:
		return convertRangeSubselect(&n)

	case pg.RangeTableFunc:
		return convertRangeTableFunc(&n)

	case pg.RangeTableFuncCol:
		return convertRangeTableFuncCol(&n)

	case pg.RangeTableSample:
		return convertRangeTableSample(&n)

	case pg.RangeTblEntry:
		return convertRangeTblEntry(&n)

	case pg.RangeTblFunction:
		return convertRangeTblFunction(&n)

	case pg.RangeTblRef:
		return convertRangeTblRef(&n)

	case pg.RangeVar:
		return convertRangeVar(&n)

	case pg.RawStmt:
		return convertRawStmt(&n)

	case pg.ReassignOwnedStmt:
		return convertReassignOwnedStmt(&n)

	case pg.RefreshMatViewStmt:
		return convertRefreshMatViewStmt(&n)

	case pg.ReindexStmt:
		return convertReindexStmt(&n)

	case pg.RelabelType:
		return convertRelabelType(&n)

	case pg.RenameStmt:
		return convertRenameStmt(&n)

	case pg.ReplicaIdentityStmt:
		return convertReplicaIdentityStmt(&n)

	case pg.ResTarget:
		return convertResTarget(&n)

	case pg.RoleSpec:
		return convertRoleSpec(&n)

	case pg.RowCompareExpr:
		return convertRowCompareExpr(&n)

	case pg.RowExpr:
		return convertRowExpr(&n)

	case pg.RowMarkClause:
		return convertRowMarkClause(&n)

	case pg.RuleStmt:
		return convertRuleStmt(&n)

	case pg.SQLValueFunction:
		return convertSQLValueFunction(&n)

	case pg.ScalarArrayOpExpr:
		return convertScalarArrayOpExpr(&n)

	case pg.SecLabelStmt:
		return convertSecLabelStmt(&n)

	case pg.SelectStmt:
		return convertSelectStmt(&n)

	case pg.SetOperationStmt:
		return convertSetOperationStmt(&n)

	case pg.SetToDefault:
		return convertSetToDefault(&n)

	case pg.SortBy:
		return convertSortBy(&n)

	case pg.SortGroupClause:
		return convertSortGroupClause(&n)

	case pg.String:
		return convertString(&n)

	case pg.SubLink:
		return convertSubLink(&n)

	case pg.SubPlan:
		return convertSubPlan(&n)

	case pg.TableFunc:
		return convertTableFunc(&n)

	case pg.TableLikeClause:
		return convertTableLikeClause(&n)

	case pg.TableSampleClause:
		return convertTableSampleClause(&n)

	case pg.TargetEntry:
		return convertTargetEntry(&n)

	case pg.TransactionStmt:
		return convertTransactionStmt(&n)

	case pg.TriggerTransition:
		return convertTriggerTransition(&n)

	case pg.TruncateStmt:
		return convertTruncateStmt(&n)

	case pg.TypeCast:
		return convertTypeCast(&n)

	case pg.TypeName:
		return convertTypeName(&n)

	case pg.UnlistenStmt:
		return convertUnlistenStmt(&n)

	case pg.UpdateStmt:
		return convertUpdateStmt(&n)

	case pg.VacuumStmt:
		return convertVacuumStmt(&n)

	case pg.Var:
		return convertVar(&n)

	case pg.VariableSetStmt:
		return convertVariableSetStmt(&n)

	case pg.VariableShowStmt:
		return convertVariableShowStmt(&n)

	case pg.ViewStmt:
		return convertViewStmt(&n)

	case pg.WindowClause:
		return convertWindowClause(&n)

	case pg.WindowDef:
		return convertWindowDef(&n)

	case pg.WindowFunc:
		return convertWindowFunc(&n)

	case pg.WithCheckOption:
		return convertWithCheckOption(&n)

	case pg.WithClause:
		return convertWithClause(&n)

	case pg.XmlExpr:
		return convertXmlExpr(&n)

	case pg.XmlSerialize:
		return convertXmlSerialize(&n)

	default:
		return &ast.TODO{}
	}
}
