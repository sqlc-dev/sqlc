// Code generated from SQLite.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // SQLite

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type BaseSQLiteVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSQLiteVisitor) VisitParse(ctx *ParseContext) interface{} {
	fmt.Println(ctx)
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSql_stmt_list(ctx *Sql_stmt_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSql_stmt(ctx *Sql_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitAlter_table_stmt(ctx *Alter_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitAnalyze_stmt(ctx *Analyze_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitAttach_stmt(ctx *Attach_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitBegin_stmt(ctx *Begin_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCommit_stmt(ctx *Commit_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCompound_select_stmt(ctx *Compound_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCreate_index_stmt(ctx *Create_index_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCreate_table_stmt(ctx *Create_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCreate_trigger_stmt(ctx *Create_trigger_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCreate_view_stmt(ctx *Create_view_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCreate_virtual_table_stmt(ctx *Create_virtual_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDelete_stmt(ctx *Delete_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDelete_stmt_limited(ctx *Delete_stmt_limitedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDetach_stmt(ctx *Detach_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDrop_index_stmt(ctx *Drop_index_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDrop_table_stmt(ctx *Drop_table_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDrop_trigger_stmt(ctx *Drop_trigger_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDrop_view_stmt(ctx *Drop_view_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitFactored_select_stmt(ctx *Factored_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitInsert_stmt(ctx *Insert_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitPragma_stmt(ctx *Pragma_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitReindex_stmt(ctx *Reindex_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitRelease_stmt(ctx *Release_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitRollback_stmt(ctx *Rollback_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSavepoint_stmt(ctx *Savepoint_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSimple_select_stmt(ctx *Simple_select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSelect_stmt(ctx *Select_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSelect_or_values(ctx *Select_or_valuesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitUpdate_stmt(ctx *Update_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitUpdate_stmt_limited(ctx *Update_stmt_limitedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitVacuum_stmt(ctx *Vacuum_stmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitColumn_def(ctx *Column_defContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitType_name(ctx *Type_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitColumn_constraint(ctx *Column_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitConflict_clause(ctx *Conflict_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitForeign_key_clause(ctx *Foreign_key_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitRaise_function(ctx *Raise_functionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitIndexed_column(ctx *Indexed_columnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_constraint(ctx *Table_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitWith_clause(ctx *With_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitQualified_table_name(ctx *Qualified_table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitOrdering_term(ctx *Ordering_termContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitPragma_value(ctx *Pragma_valueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCommon_table_expression(ctx *Common_table_expressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitResult_column(ctx *Result_columnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_or_subquery(ctx *Table_or_subqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitJoin_clause(ctx *Join_clauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitJoin_operator(ctx *Join_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitJoin_constraint(ctx *Join_constraintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSelect_core(ctx *Select_coreContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCompound_operator(ctx *Compound_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSigned_number(ctx *Signed_numberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitLiteral_value(ctx *Literal_valueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitUnary_operator(ctx *Unary_operatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitError_message(ctx *Error_messageContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitModule_argument(ctx *Module_argumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitColumn_alias(ctx *Column_aliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitKeyword(ctx *KeywordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitName(ctx *NameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitFunction_name(ctx *Function_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitDatabase_name(ctx *Database_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSchema_name(ctx *Schema_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_function_name(ctx *Table_function_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_name(ctx *Table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_or_index_name(ctx *Table_or_index_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitNew_table_name(ctx *New_table_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitColumn_name(ctx *Column_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitCollation_name(ctx *Collation_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitForeign_table(ctx *Foreign_tableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitIndex_name(ctx *Index_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTrigger_name(ctx *Trigger_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitView_name(ctx *View_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitModule_name(ctx *Module_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitPragma_name(ctx *Pragma_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitSavepoint_name(ctx *Savepoint_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTable_alias(ctx *Table_aliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitTransaction_name(ctx *Transaction_nameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSQLiteVisitor) VisitAny_name(ctx *Any_nameContext) interface{} {
	return v.VisitChildren(ctx)
}
