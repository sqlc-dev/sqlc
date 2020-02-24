package sqlite

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/kyleconroy/sqlc/internal/sqlite/parser"
)

type visitor struct {
}

// ParseTreeVisitor interace
func (v *visitor) Visit(tree antlr.ParseTree) interface{}            { return v }
func (v *visitor) VisitChildren(node antlr.RuleNode) interface{}     { return v }
func (v *visitor) VisitTerminal(node antlr.TerminalNode) interface{} { return v }
func (v *visitor) VisitErrorNode(node antlr.ErrorNode) interface{}   { return v }

// SQLiteVisitor interface
func (v *visitor) VisitParse(ctx *parser.ParseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSql_stmt_list(ctx *parser.Sql_stmt_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSql_stmt(ctx *parser.Sql_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitAlter_table_stmt(ctx *parser.Alter_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitAnalyze_stmt(ctx *parser.Analyze_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitAttach_stmt(ctx *parser.Attach_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitBegin_stmt(ctx *parser.Begin_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCommit_stmt(ctx *parser.Commit_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCompound_select_stmt(ctx *parser.Compound_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCreate_index_stmt(ctx *parser.Create_index_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCreate_table_stmt(ctx *parser.Create_table_stmtContext) interface{} {
	fmt.Println("CREATE TABLE", ctx)
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCreate_trigger_stmt(ctx *parser.Create_trigger_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCreate_view_stmt(ctx *parser.Create_view_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCreate_virtual_table_stmt(ctx *parser.Create_virtual_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDelete_stmt(ctx *parser.Delete_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDelete_stmt_limited(ctx *parser.Delete_stmt_limitedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDetach_stmt(ctx *parser.Detach_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDrop_index_stmt(ctx *parser.Drop_index_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDrop_table_stmt(ctx *parser.Drop_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDrop_trigger_stmt(ctx *parser.Drop_trigger_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDrop_view_stmt(ctx *parser.Drop_view_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitFactored_select_stmt(ctx *parser.Factored_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitInsert_stmt(ctx *parser.Insert_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitPragma_stmt(ctx *parser.Pragma_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitReindex_stmt(ctx *parser.Reindex_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitRelease_stmt(ctx *parser.Release_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitRollback_stmt(ctx *parser.Rollback_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSavepoint_stmt(ctx *parser.Savepoint_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSimple_select_stmt(ctx *parser.Simple_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSelect_stmt(ctx *parser.Select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSelect_or_values(ctx *parser.Select_or_valuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitUpdate_stmt(ctx *parser.Update_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitUpdate_stmt_limited(ctx *parser.Update_stmt_limitedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitVacuum_stmt(ctx *parser.Vacuum_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitColumn_def(ctx *parser.Column_defContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitType_name(ctx *parser.Type_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitColumn_constraint(ctx *parser.Column_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitConflict_clause(ctx *parser.Conflict_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitExpr(ctx *parser.ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitForeign_key_clause(ctx *parser.Foreign_key_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitRaise_function(ctx *parser.Raise_functionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitIndexed_column(ctx *parser.Indexed_columnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_constraint(ctx *parser.Table_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitWith_clause(ctx *parser.With_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitQualified_table_name(ctx *parser.Qualified_table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitOrdering_term(ctx *parser.Ordering_termContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitPragma_value(ctx *parser.Pragma_valueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCommon_table_expression(ctx *parser.Common_table_expressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitResult_column(ctx *parser.Result_columnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_or_subquery(ctx *parser.Table_or_subqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitJoin_clause(ctx *parser.Join_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitJoin_operator(ctx *parser.Join_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitJoin_constraint(ctx *parser.Join_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSelect_core(ctx *parser.Select_coreContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCompound_operator(ctx *parser.Compound_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSigned_number(ctx *parser.Signed_numberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitLiteral_value(ctx *parser.Literal_valueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitUnary_operator(ctx *parser.Unary_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitError_message(ctx *parser.Error_messageContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitModule_argument(ctx *parser.Module_argumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitColumn_alias(ctx *parser.Column_aliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitKeyword(ctx *parser.KeywordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitName(ctx *parser.NameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitFunction_name(ctx *parser.Function_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitDatabase_name(ctx *parser.Database_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSchema_name(ctx *parser.Schema_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_function_name(ctx *parser.Table_function_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_name(ctx *parser.Table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_or_index_name(ctx *parser.Table_or_index_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitNew_table_name(ctx *parser.New_table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitColumn_name(ctx *parser.Column_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitCollation_name(ctx *parser.Collation_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitForeign_table(ctx *parser.Foreign_tableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitIndex_name(ctx *parser.Index_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTrigger_name(ctx *parser.Trigger_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitView_name(ctx *parser.View_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitModule_name(ctx *parser.Module_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitPragma_name(ctx *parser.Pragma_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitSavepoint_name(ctx *parser.Savepoint_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTable_alias(ctx *parser.Table_aliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitTransaction_name(ctx *parser.Transaction_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *visitor) VisitAny_name(ctx *parser.Any_nameContext) interface{} {
	return v.VisitChildren(ctx)
}
