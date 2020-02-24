// Code generated from SQLite.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // SQLite

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by SQLiteParser.
type SQLiteVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SQLiteParser#parse.
	VisitParse(ctx *ParseContext) interface{}

	// Visit a parse tree produced by SQLiteParser#sql_stmt_list.
	VisitSql_stmt_list(ctx *Sql_stmt_listContext) interface{}

	// Visit a parse tree produced by SQLiteParser#sql_stmt.
	VisitSql_stmt(ctx *Sql_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#alter_table_stmt.
	VisitAlter_table_stmt(ctx *Alter_table_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#analyze_stmt.
	VisitAnalyze_stmt(ctx *Analyze_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#attach_stmt.
	VisitAttach_stmt(ctx *Attach_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#begin_stmt.
	VisitBegin_stmt(ctx *Begin_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#commit_stmt.
	VisitCommit_stmt(ctx *Commit_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#compound_select_stmt.
	VisitCompound_select_stmt(ctx *Compound_select_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#create_index_stmt.
	VisitCreate_index_stmt(ctx *Create_index_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#create_table_stmt.
	VisitCreate_table_stmt(ctx *Create_table_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#create_trigger_stmt.
	VisitCreate_trigger_stmt(ctx *Create_trigger_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#create_view_stmt.
	VisitCreate_view_stmt(ctx *Create_view_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#create_virtual_table_stmt.
	VisitCreate_virtual_table_stmt(ctx *Create_virtual_table_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#delete_stmt.
	VisitDelete_stmt(ctx *Delete_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#delete_stmt_limited.
	VisitDelete_stmt_limited(ctx *Delete_stmt_limitedContext) interface{}

	// Visit a parse tree produced by SQLiteParser#detach_stmt.
	VisitDetach_stmt(ctx *Detach_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#drop_index_stmt.
	VisitDrop_index_stmt(ctx *Drop_index_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#drop_table_stmt.
	VisitDrop_table_stmt(ctx *Drop_table_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#drop_trigger_stmt.
	VisitDrop_trigger_stmt(ctx *Drop_trigger_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#drop_view_stmt.
	VisitDrop_view_stmt(ctx *Drop_view_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#factored_select_stmt.
	VisitFactored_select_stmt(ctx *Factored_select_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#insert_stmt.
	VisitInsert_stmt(ctx *Insert_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#pragma_stmt.
	VisitPragma_stmt(ctx *Pragma_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#reindex_stmt.
	VisitReindex_stmt(ctx *Reindex_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#release_stmt.
	VisitRelease_stmt(ctx *Release_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#rollback_stmt.
	VisitRollback_stmt(ctx *Rollback_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#savepoint_stmt.
	VisitSavepoint_stmt(ctx *Savepoint_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#simple_select_stmt.
	VisitSimple_select_stmt(ctx *Simple_select_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#select_stmt.
	VisitSelect_stmt(ctx *Select_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#select_or_values.
	VisitSelect_or_values(ctx *Select_or_valuesContext) interface{}

	// Visit a parse tree produced by SQLiteParser#update_stmt.
	VisitUpdate_stmt(ctx *Update_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#update_stmt_limited.
	VisitUpdate_stmt_limited(ctx *Update_stmt_limitedContext) interface{}

	// Visit a parse tree produced by SQLiteParser#vacuum_stmt.
	VisitVacuum_stmt(ctx *Vacuum_stmtContext) interface{}

	// Visit a parse tree produced by SQLiteParser#column_def.
	VisitColumn_def(ctx *Column_defContext) interface{}

	// Visit a parse tree produced by SQLiteParser#type_name.
	VisitType_name(ctx *Type_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#column_constraint.
	VisitColumn_constraint(ctx *Column_constraintContext) interface{}

	// Visit a parse tree produced by SQLiteParser#conflict_clause.
	VisitConflict_clause(ctx *Conflict_clauseContext) interface{}

	// Visit a parse tree produced by SQLiteParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by SQLiteParser#foreign_key_clause.
	VisitForeign_key_clause(ctx *Foreign_key_clauseContext) interface{}

	// Visit a parse tree produced by SQLiteParser#raise_function.
	VisitRaise_function(ctx *Raise_functionContext) interface{}

	// Visit a parse tree produced by SQLiteParser#indexed_column.
	VisitIndexed_column(ctx *Indexed_columnContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_constraint.
	VisitTable_constraint(ctx *Table_constraintContext) interface{}

	// Visit a parse tree produced by SQLiteParser#with_clause.
	VisitWith_clause(ctx *With_clauseContext) interface{}

	// Visit a parse tree produced by SQLiteParser#qualified_table_name.
	VisitQualified_table_name(ctx *Qualified_table_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#ordering_term.
	VisitOrdering_term(ctx *Ordering_termContext) interface{}

	// Visit a parse tree produced by SQLiteParser#pragma_value.
	VisitPragma_value(ctx *Pragma_valueContext) interface{}

	// Visit a parse tree produced by SQLiteParser#common_table_expression.
	VisitCommon_table_expression(ctx *Common_table_expressionContext) interface{}

	// Visit a parse tree produced by SQLiteParser#result_column.
	VisitResult_column(ctx *Result_columnContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_or_subquery.
	VisitTable_or_subquery(ctx *Table_or_subqueryContext) interface{}

	// Visit a parse tree produced by SQLiteParser#join_clause.
	VisitJoin_clause(ctx *Join_clauseContext) interface{}

	// Visit a parse tree produced by SQLiteParser#join_operator.
	VisitJoin_operator(ctx *Join_operatorContext) interface{}

	// Visit a parse tree produced by SQLiteParser#join_constraint.
	VisitJoin_constraint(ctx *Join_constraintContext) interface{}

	// Visit a parse tree produced by SQLiteParser#select_core.
	VisitSelect_core(ctx *Select_coreContext) interface{}

	// Visit a parse tree produced by SQLiteParser#compound_operator.
	VisitCompound_operator(ctx *Compound_operatorContext) interface{}

	// Visit a parse tree produced by SQLiteParser#signed_number.
	VisitSigned_number(ctx *Signed_numberContext) interface{}

	// Visit a parse tree produced by SQLiteParser#literal_value.
	VisitLiteral_value(ctx *Literal_valueContext) interface{}

	// Visit a parse tree produced by SQLiteParser#unary_operator.
	VisitUnary_operator(ctx *Unary_operatorContext) interface{}

	// Visit a parse tree produced by SQLiteParser#error_message.
	VisitError_message(ctx *Error_messageContext) interface{}

	// Visit a parse tree produced by SQLiteParser#module_argument.
	VisitModule_argument(ctx *Module_argumentContext) interface{}

	// Visit a parse tree produced by SQLiteParser#column_alias.
	VisitColumn_alias(ctx *Column_aliasContext) interface{}

	// Visit a parse tree produced by SQLiteParser#keyword.
	VisitKeyword(ctx *KeywordContext) interface{}

	// Visit a parse tree produced by SQLiteParser#name.
	VisitName(ctx *NameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#function_name.
	VisitFunction_name(ctx *Function_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#database_name.
	VisitDatabase_name(ctx *Database_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#schema_name.
	VisitSchema_name(ctx *Schema_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_function_name.
	VisitTable_function_name(ctx *Table_function_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_name.
	VisitTable_name(ctx *Table_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_or_index_name.
	VisitTable_or_index_name(ctx *Table_or_index_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#new_table_name.
	VisitNew_table_name(ctx *New_table_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#column_name.
	VisitColumn_name(ctx *Column_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#collation_name.
	VisitCollation_name(ctx *Collation_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#foreign_table.
	VisitForeign_table(ctx *Foreign_tableContext) interface{}

	// Visit a parse tree produced by SQLiteParser#index_name.
	VisitIndex_name(ctx *Index_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#trigger_name.
	VisitTrigger_name(ctx *Trigger_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#view_name.
	VisitView_name(ctx *View_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#module_name.
	VisitModule_name(ctx *Module_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#pragma_name.
	VisitPragma_name(ctx *Pragma_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#savepoint_name.
	VisitSavepoint_name(ctx *Savepoint_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#table_alias.
	VisitTable_alias(ctx *Table_aliasContext) interface{}

	// Visit a parse tree produced by SQLiteParser#transaction_name.
	VisitTransaction_name(ctx *Transaction_nameContext) interface{}

	// Visit a parse tree produced by SQLiteParser#any_name.
	VisitAny_name(ctx *Any_nameContext) interface{}
}
