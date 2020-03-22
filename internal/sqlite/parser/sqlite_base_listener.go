// Code generated from SQLite.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // SQLite

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSQLiteListener is a complete listener for a parse tree produced by SQLiteParser.
type BaseSQLiteListener struct{}

var _ SQLiteListener = &BaseSQLiteListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSQLiteListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSQLiteListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSQLiteListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSQLiteListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterParse is called when production parse is entered.
func (s *BaseSQLiteListener) EnterParse(ctx *ParseContext) {}

// ExitParse is called when production parse is exited.
func (s *BaseSQLiteListener) ExitParse(ctx *ParseContext) {}

// EnterSql_stmt_list is called when production sql_stmt_list is entered.
func (s *BaseSQLiteListener) EnterSql_stmt_list(ctx *Sql_stmt_listContext) {}

// ExitSql_stmt_list is called when production sql_stmt_list is exited.
func (s *BaseSQLiteListener) ExitSql_stmt_list(ctx *Sql_stmt_listContext) {}

// EnterSql_stmt is called when production sql_stmt is entered.
func (s *BaseSQLiteListener) EnterSql_stmt(ctx *Sql_stmtContext) {}

// ExitSql_stmt is called when production sql_stmt is exited.
func (s *BaseSQLiteListener) ExitSql_stmt(ctx *Sql_stmtContext) {}

// EnterAlter_table_stmt is called when production alter_table_stmt is entered.
func (s *BaseSQLiteListener) EnterAlter_table_stmt(ctx *Alter_table_stmtContext) {}

// ExitAlter_table_stmt is called when production alter_table_stmt is exited.
func (s *BaseSQLiteListener) ExitAlter_table_stmt(ctx *Alter_table_stmtContext) {}

// EnterAnalyze_stmt is called when production analyze_stmt is entered.
func (s *BaseSQLiteListener) EnterAnalyze_stmt(ctx *Analyze_stmtContext) {}

// ExitAnalyze_stmt is called when production analyze_stmt is exited.
func (s *BaseSQLiteListener) ExitAnalyze_stmt(ctx *Analyze_stmtContext) {}

// EnterAttach_stmt is called when production attach_stmt is entered.
func (s *BaseSQLiteListener) EnterAttach_stmt(ctx *Attach_stmtContext) {}

// ExitAttach_stmt is called when production attach_stmt is exited.
func (s *BaseSQLiteListener) ExitAttach_stmt(ctx *Attach_stmtContext) {}

// EnterBegin_stmt is called when production begin_stmt is entered.
func (s *BaseSQLiteListener) EnterBegin_stmt(ctx *Begin_stmtContext) {}

// ExitBegin_stmt is called when production begin_stmt is exited.
func (s *BaseSQLiteListener) ExitBegin_stmt(ctx *Begin_stmtContext) {}

// EnterCommit_stmt is called when production commit_stmt is entered.
func (s *BaseSQLiteListener) EnterCommit_stmt(ctx *Commit_stmtContext) {}

// ExitCommit_stmt is called when production commit_stmt is exited.
func (s *BaseSQLiteListener) ExitCommit_stmt(ctx *Commit_stmtContext) {}

// EnterCompound_select_stmt is called when production compound_select_stmt is entered.
func (s *BaseSQLiteListener) EnterCompound_select_stmt(ctx *Compound_select_stmtContext) {}

// ExitCompound_select_stmt is called when production compound_select_stmt is exited.
func (s *BaseSQLiteListener) ExitCompound_select_stmt(ctx *Compound_select_stmtContext) {}

// EnterCreate_index_stmt is called when production create_index_stmt is entered.
func (s *BaseSQLiteListener) EnterCreate_index_stmt(ctx *Create_index_stmtContext) {}

// ExitCreate_index_stmt is called when production create_index_stmt is exited.
func (s *BaseSQLiteListener) ExitCreate_index_stmt(ctx *Create_index_stmtContext) {}

// EnterCreate_table_stmt is called when production create_table_stmt is entered.
func (s *BaseSQLiteListener) EnterCreate_table_stmt(ctx *Create_table_stmtContext) {}

// ExitCreate_table_stmt is called when production create_table_stmt is exited.
func (s *BaseSQLiteListener) ExitCreate_table_stmt(ctx *Create_table_stmtContext) {}

// EnterCreate_trigger_stmt is called when production create_trigger_stmt is entered.
func (s *BaseSQLiteListener) EnterCreate_trigger_stmt(ctx *Create_trigger_stmtContext) {}

// ExitCreate_trigger_stmt is called when production create_trigger_stmt is exited.
func (s *BaseSQLiteListener) ExitCreate_trigger_stmt(ctx *Create_trigger_stmtContext) {}

// EnterCreate_view_stmt is called when production create_view_stmt is entered.
func (s *BaseSQLiteListener) EnterCreate_view_stmt(ctx *Create_view_stmtContext) {}

// ExitCreate_view_stmt is called when production create_view_stmt is exited.
func (s *BaseSQLiteListener) ExitCreate_view_stmt(ctx *Create_view_stmtContext) {}

// EnterCreate_virtual_table_stmt is called when production create_virtual_table_stmt is entered.
func (s *BaseSQLiteListener) EnterCreate_virtual_table_stmt(ctx *Create_virtual_table_stmtContext) {}

// ExitCreate_virtual_table_stmt is called when production create_virtual_table_stmt is exited.
func (s *BaseSQLiteListener) ExitCreate_virtual_table_stmt(ctx *Create_virtual_table_stmtContext) {}

// EnterDelete_stmt is called when production delete_stmt is entered.
func (s *BaseSQLiteListener) EnterDelete_stmt(ctx *Delete_stmtContext) {}

// ExitDelete_stmt is called when production delete_stmt is exited.
func (s *BaseSQLiteListener) ExitDelete_stmt(ctx *Delete_stmtContext) {}

// EnterDelete_stmt_limited is called when production delete_stmt_limited is entered.
func (s *BaseSQLiteListener) EnterDelete_stmt_limited(ctx *Delete_stmt_limitedContext) {}

// ExitDelete_stmt_limited is called when production delete_stmt_limited is exited.
func (s *BaseSQLiteListener) ExitDelete_stmt_limited(ctx *Delete_stmt_limitedContext) {}

// EnterDetach_stmt is called when production detach_stmt is entered.
func (s *BaseSQLiteListener) EnterDetach_stmt(ctx *Detach_stmtContext) {}

// ExitDetach_stmt is called when production detach_stmt is exited.
func (s *BaseSQLiteListener) ExitDetach_stmt(ctx *Detach_stmtContext) {}

// EnterDrop_index_stmt is called when production drop_index_stmt is entered.
func (s *BaseSQLiteListener) EnterDrop_index_stmt(ctx *Drop_index_stmtContext) {}

// ExitDrop_index_stmt is called when production drop_index_stmt is exited.
func (s *BaseSQLiteListener) ExitDrop_index_stmt(ctx *Drop_index_stmtContext) {}

// EnterDrop_table_stmt is called when production drop_table_stmt is entered.
func (s *BaseSQLiteListener) EnterDrop_table_stmt(ctx *Drop_table_stmtContext) {}

// ExitDrop_table_stmt is called when production drop_table_stmt is exited.
func (s *BaseSQLiteListener) ExitDrop_table_stmt(ctx *Drop_table_stmtContext) {}

// EnterDrop_trigger_stmt is called when production drop_trigger_stmt is entered.
func (s *BaseSQLiteListener) EnterDrop_trigger_stmt(ctx *Drop_trigger_stmtContext) {}

// ExitDrop_trigger_stmt is called when production drop_trigger_stmt is exited.
func (s *BaseSQLiteListener) ExitDrop_trigger_stmt(ctx *Drop_trigger_stmtContext) {}

// EnterDrop_view_stmt is called when production drop_view_stmt is entered.
func (s *BaseSQLiteListener) EnterDrop_view_stmt(ctx *Drop_view_stmtContext) {}

// ExitDrop_view_stmt is called when production drop_view_stmt is exited.
func (s *BaseSQLiteListener) ExitDrop_view_stmt(ctx *Drop_view_stmtContext) {}

// EnterFactored_select_stmt is called when production factored_select_stmt is entered.
func (s *BaseSQLiteListener) EnterFactored_select_stmt(ctx *Factored_select_stmtContext) {}

// ExitFactored_select_stmt is called when production factored_select_stmt is exited.
func (s *BaseSQLiteListener) ExitFactored_select_stmt(ctx *Factored_select_stmtContext) {}

// EnterInsert_stmt is called when production insert_stmt is entered.
func (s *BaseSQLiteListener) EnterInsert_stmt(ctx *Insert_stmtContext) {}

// ExitInsert_stmt is called when production insert_stmt is exited.
func (s *BaseSQLiteListener) ExitInsert_stmt(ctx *Insert_stmtContext) {}

// EnterPragma_stmt is called when production pragma_stmt is entered.
func (s *BaseSQLiteListener) EnterPragma_stmt(ctx *Pragma_stmtContext) {}

// ExitPragma_stmt is called when production pragma_stmt is exited.
func (s *BaseSQLiteListener) ExitPragma_stmt(ctx *Pragma_stmtContext) {}

// EnterReindex_stmt is called when production reindex_stmt is entered.
func (s *BaseSQLiteListener) EnterReindex_stmt(ctx *Reindex_stmtContext) {}

// ExitReindex_stmt is called when production reindex_stmt is exited.
func (s *BaseSQLiteListener) ExitReindex_stmt(ctx *Reindex_stmtContext) {}

// EnterRelease_stmt is called when production release_stmt is entered.
func (s *BaseSQLiteListener) EnterRelease_stmt(ctx *Release_stmtContext) {}

// ExitRelease_stmt is called when production release_stmt is exited.
func (s *BaseSQLiteListener) ExitRelease_stmt(ctx *Release_stmtContext) {}

// EnterRollback_stmt is called when production rollback_stmt is entered.
func (s *BaseSQLiteListener) EnterRollback_stmt(ctx *Rollback_stmtContext) {}

// ExitRollback_stmt is called when production rollback_stmt is exited.
func (s *BaseSQLiteListener) ExitRollback_stmt(ctx *Rollback_stmtContext) {}

// EnterSavepoint_stmt is called when production savepoint_stmt is entered.
func (s *BaseSQLiteListener) EnterSavepoint_stmt(ctx *Savepoint_stmtContext) {}

// ExitSavepoint_stmt is called when production savepoint_stmt is exited.
func (s *BaseSQLiteListener) ExitSavepoint_stmt(ctx *Savepoint_stmtContext) {}

// EnterSimple_select_stmt is called when production simple_select_stmt is entered.
func (s *BaseSQLiteListener) EnterSimple_select_stmt(ctx *Simple_select_stmtContext) {}

// ExitSimple_select_stmt is called when production simple_select_stmt is exited.
func (s *BaseSQLiteListener) ExitSimple_select_stmt(ctx *Simple_select_stmtContext) {}

// EnterSelect_stmt is called when production select_stmt is entered.
func (s *BaseSQLiteListener) EnterSelect_stmt(ctx *Select_stmtContext) {}

// ExitSelect_stmt is called when production select_stmt is exited.
func (s *BaseSQLiteListener) ExitSelect_stmt(ctx *Select_stmtContext) {}

// EnterSelect_or_values is called when production select_or_values is entered.
func (s *BaseSQLiteListener) EnterSelect_or_values(ctx *Select_or_valuesContext) {}

// ExitSelect_or_values is called when production select_or_values is exited.
func (s *BaseSQLiteListener) ExitSelect_or_values(ctx *Select_or_valuesContext) {}

// EnterUpdate_stmt is called when production update_stmt is entered.
func (s *BaseSQLiteListener) EnterUpdate_stmt(ctx *Update_stmtContext) {}

// ExitUpdate_stmt is called when production update_stmt is exited.
func (s *BaseSQLiteListener) ExitUpdate_stmt(ctx *Update_stmtContext) {}

// EnterUpdate_stmt_limited is called when production update_stmt_limited is entered.
func (s *BaseSQLiteListener) EnterUpdate_stmt_limited(ctx *Update_stmt_limitedContext) {}

// ExitUpdate_stmt_limited is called when production update_stmt_limited is exited.
func (s *BaseSQLiteListener) ExitUpdate_stmt_limited(ctx *Update_stmt_limitedContext) {}

// EnterVacuum_stmt is called when production vacuum_stmt is entered.
func (s *BaseSQLiteListener) EnterVacuum_stmt(ctx *Vacuum_stmtContext) {}

// ExitVacuum_stmt is called when production vacuum_stmt is exited.
func (s *BaseSQLiteListener) ExitVacuum_stmt(ctx *Vacuum_stmtContext) {}

// EnterColumn_def is called when production column_def is entered.
func (s *BaseSQLiteListener) EnterColumn_def(ctx *Column_defContext) {}

// ExitColumn_def is called when production column_def is exited.
func (s *BaseSQLiteListener) ExitColumn_def(ctx *Column_defContext) {}

// EnterType_name is called when production type_name is entered.
func (s *BaseSQLiteListener) EnterType_name(ctx *Type_nameContext) {}

// ExitType_name is called when production type_name is exited.
func (s *BaseSQLiteListener) ExitType_name(ctx *Type_nameContext) {}

// EnterColumn_constraint is called when production column_constraint is entered.
func (s *BaseSQLiteListener) EnterColumn_constraint(ctx *Column_constraintContext) {}

// ExitColumn_constraint is called when production column_constraint is exited.
func (s *BaseSQLiteListener) ExitColumn_constraint(ctx *Column_constraintContext) {}

// EnterConflict_clause is called when production conflict_clause is entered.
func (s *BaseSQLiteListener) EnterConflict_clause(ctx *Conflict_clauseContext) {}

// ExitConflict_clause is called when production conflict_clause is exited.
func (s *BaseSQLiteListener) ExitConflict_clause(ctx *Conflict_clauseContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSQLiteListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSQLiteListener) ExitExpr(ctx *ExprContext) {}

// EnterForeign_key_clause is called when production foreign_key_clause is entered.
func (s *BaseSQLiteListener) EnterForeign_key_clause(ctx *Foreign_key_clauseContext) {}

// ExitForeign_key_clause is called when production foreign_key_clause is exited.
func (s *BaseSQLiteListener) ExitForeign_key_clause(ctx *Foreign_key_clauseContext) {}

// EnterRaise_function is called when production raise_function is entered.
func (s *BaseSQLiteListener) EnterRaise_function(ctx *Raise_functionContext) {}

// ExitRaise_function is called when production raise_function is exited.
func (s *BaseSQLiteListener) ExitRaise_function(ctx *Raise_functionContext) {}

// EnterIndexed_column is called when production indexed_column is entered.
func (s *BaseSQLiteListener) EnterIndexed_column(ctx *Indexed_columnContext) {}

// ExitIndexed_column is called when production indexed_column is exited.
func (s *BaseSQLiteListener) ExitIndexed_column(ctx *Indexed_columnContext) {}

// EnterTable_constraint is called when production table_constraint is entered.
func (s *BaseSQLiteListener) EnterTable_constraint(ctx *Table_constraintContext) {}

// ExitTable_constraint is called when production table_constraint is exited.
func (s *BaseSQLiteListener) ExitTable_constraint(ctx *Table_constraintContext) {}

// EnterWith_clause is called when production with_clause is entered.
func (s *BaseSQLiteListener) EnterWith_clause(ctx *With_clauseContext) {}

// ExitWith_clause is called when production with_clause is exited.
func (s *BaseSQLiteListener) ExitWith_clause(ctx *With_clauseContext) {}

// EnterQualified_table_name is called when production qualified_table_name is entered.
func (s *BaseSQLiteListener) EnterQualified_table_name(ctx *Qualified_table_nameContext) {}

// ExitQualified_table_name is called when production qualified_table_name is exited.
func (s *BaseSQLiteListener) ExitQualified_table_name(ctx *Qualified_table_nameContext) {}

// EnterOrdering_term is called when production ordering_term is entered.
func (s *BaseSQLiteListener) EnterOrdering_term(ctx *Ordering_termContext) {}

// ExitOrdering_term is called when production ordering_term is exited.
func (s *BaseSQLiteListener) ExitOrdering_term(ctx *Ordering_termContext) {}

// EnterPragma_value is called when production pragma_value is entered.
func (s *BaseSQLiteListener) EnterPragma_value(ctx *Pragma_valueContext) {}

// ExitPragma_value is called when production pragma_value is exited.
func (s *BaseSQLiteListener) ExitPragma_value(ctx *Pragma_valueContext) {}

// EnterCommon_table_expression is called when production common_table_expression is entered.
func (s *BaseSQLiteListener) EnterCommon_table_expression(ctx *Common_table_expressionContext) {}

// ExitCommon_table_expression is called when production common_table_expression is exited.
func (s *BaseSQLiteListener) ExitCommon_table_expression(ctx *Common_table_expressionContext) {}

// EnterResult_column is called when production result_column is entered.
func (s *BaseSQLiteListener) EnterResult_column(ctx *Result_columnContext) {}

// ExitResult_column is called when production result_column is exited.
func (s *BaseSQLiteListener) ExitResult_column(ctx *Result_columnContext) {}

// EnterTable_or_subquery is called when production table_or_subquery is entered.
func (s *BaseSQLiteListener) EnterTable_or_subquery(ctx *Table_or_subqueryContext) {}

// ExitTable_or_subquery is called when production table_or_subquery is exited.
func (s *BaseSQLiteListener) ExitTable_or_subquery(ctx *Table_or_subqueryContext) {}

// EnterJoin_clause is called when production join_clause is entered.
func (s *BaseSQLiteListener) EnterJoin_clause(ctx *Join_clauseContext) {}

// ExitJoin_clause is called when production join_clause is exited.
func (s *BaseSQLiteListener) ExitJoin_clause(ctx *Join_clauseContext) {}

// EnterJoin_operator is called when production join_operator is entered.
func (s *BaseSQLiteListener) EnterJoin_operator(ctx *Join_operatorContext) {}

// ExitJoin_operator is called when production join_operator is exited.
func (s *BaseSQLiteListener) ExitJoin_operator(ctx *Join_operatorContext) {}

// EnterJoin_constraint is called when production join_constraint is entered.
func (s *BaseSQLiteListener) EnterJoin_constraint(ctx *Join_constraintContext) {}

// ExitJoin_constraint is called when production join_constraint is exited.
func (s *BaseSQLiteListener) ExitJoin_constraint(ctx *Join_constraintContext) {}

// EnterSelect_core is called when production select_core is entered.
func (s *BaseSQLiteListener) EnterSelect_core(ctx *Select_coreContext) {}

// ExitSelect_core is called when production select_core is exited.
func (s *BaseSQLiteListener) ExitSelect_core(ctx *Select_coreContext) {}

// EnterCompound_operator is called when production compound_operator is entered.
func (s *BaseSQLiteListener) EnterCompound_operator(ctx *Compound_operatorContext) {}

// ExitCompound_operator is called when production compound_operator is exited.
func (s *BaseSQLiteListener) ExitCompound_operator(ctx *Compound_operatorContext) {}

// EnterSigned_number is called when production signed_number is entered.
func (s *BaseSQLiteListener) EnterSigned_number(ctx *Signed_numberContext) {}

// ExitSigned_number is called when production signed_number is exited.
func (s *BaseSQLiteListener) ExitSigned_number(ctx *Signed_numberContext) {}

// EnterLiteral_value is called when production literal_value is entered.
func (s *BaseSQLiteListener) EnterLiteral_value(ctx *Literal_valueContext) {}

// ExitLiteral_value is called when production literal_value is exited.
func (s *BaseSQLiteListener) ExitLiteral_value(ctx *Literal_valueContext) {}

// EnterUnary_operator is called when production unary_operator is entered.
func (s *BaseSQLiteListener) EnterUnary_operator(ctx *Unary_operatorContext) {}

// ExitUnary_operator is called when production unary_operator is exited.
func (s *BaseSQLiteListener) ExitUnary_operator(ctx *Unary_operatorContext) {}

// EnterError_message is called when production error_message is entered.
func (s *BaseSQLiteListener) EnterError_message(ctx *Error_messageContext) {}

// ExitError_message is called when production error_message is exited.
func (s *BaseSQLiteListener) ExitError_message(ctx *Error_messageContext) {}

// EnterModule_argument is called when production module_argument is entered.
func (s *BaseSQLiteListener) EnterModule_argument(ctx *Module_argumentContext) {}

// ExitModule_argument is called when production module_argument is exited.
func (s *BaseSQLiteListener) ExitModule_argument(ctx *Module_argumentContext) {}

// EnterColumn_alias is called when production column_alias is entered.
func (s *BaseSQLiteListener) EnterColumn_alias(ctx *Column_aliasContext) {}

// ExitColumn_alias is called when production column_alias is exited.
func (s *BaseSQLiteListener) ExitColumn_alias(ctx *Column_aliasContext) {}

// EnterKeyword is called when production keyword is entered.
func (s *BaseSQLiteListener) EnterKeyword(ctx *KeywordContext) {}

// ExitKeyword is called when production keyword is exited.
func (s *BaseSQLiteListener) ExitKeyword(ctx *KeywordContext) {}

// EnterName is called when production name is entered.
func (s *BaseSQLiteListener) EnterName(ctx *NameContext) {}

// ExitName is called when production name is exited.
func (s *BaseSQLiteListener) ExitName(ctx *NameContext) {}

// EnterFunction_name is called when production function_name is entered.
func (s *BaseSQLiteListener) EnterFunction_name(ctx *Function_nameContext) {}

// ExitFunction_name is called when production function_name is exited.
func (s *BaseSQLiteListener) ExitFunction_name(ctx *Function_nameContext) {}

// EnterDatabase_name is called when production database_name is entered.
func (s *BaseSQLiteListener) EnterDatabase_name(ctx *Database_nameContext) {}

// ExitDatabase_name is called when production database_name is exited.
func (s *BaseSQLiteListener) ExitDatabase_name(ctx *Database_nameContext) {}

// EnterSchema_name is called when production schema_name is entered.
func (s *BaseSQLiteListener) EnterSchema_name(ctx *Schema_nameContext) {}

// ExitSchema_name is called when production schema_name is exited.
func (s *BaseSQLiteListener) ExitSchema_name(ctx *Schema_nameContext) {}

// EnterTable_function_name is called when production table_function_name is entered.
func (s *BaseSQLiteListener) EnterTable_function_name(ctx *Table_function_nameContext) {}

// ExitTable_function_name is called when production table_function_name is exited.
func (s *BaseSQLiteListener) ExitTable_function_name(ctx *Table_function_nameContext) {}

// EnterTable_name is called when production table_name is entered.
func (s *BaseSQLiteListener) EnterTable_name(ctx *Table_nameContext) {}

// ExitTable_name is called when production table_name is exited.
func (s *BaseSQLiteListener) ExitTable_name(ctx *Table_nameContext) {}

// EnterTable_or_index_name is called when production table_or_index_name is entered.
func (s *BaseSQLiteListener) EnterTable_or_index_name(ctx *Table_or_index_nameContext) {}

// ExitTable_or_index_name is called when production table_or_index_name is exited.
func (s *BaseSQLiteListener) ExitTable_or_index_name(ctx *Table_or_index_nameContext) {}

// EnterNew_table_name is called when production new_table_name is entered.
func (s *BaseSQLiteListener) EnterNew_table_name(ctx *New_table_nameContext) {}

// ExitNew_table_name is called when production new_table_name is exited.
func (s *BaseSQLiteListener) ExitNew_table_name(ctx *New_table_nameContext) {}

// EnterColumn_name is called when production column_name is entered.
func (s *BaseSQLiteListener) EnterColumn_name(ctx *Column_nameContext) {}

// ExitColumn_name is called when production column_name is exited.
func (s *BaseSQLiteListener) ExitColumn_name(ctx *Column_nameContext) {}

// EnterNew_column_name is called when production new_column_name is entered.
func (s *BaseSQLiteListener) EnterNew_column_name(ctx *New_column_nameContext) {}

// ExitNew_column_name is called when production new_column_name is exited.
func (s *BaseSQLiteListener) ExitNew_column_name(ctx *New_column_nameContext) {}

// EnterCollation_name is called when production collation_name is entered.
func (s *BaseSQLiteListener) EnterCollation_name(ctx *Collation_nameContext) {}

// ExitCollation_name is called when production collation_name is exited.
func (s *BaseSQLiteListener) ExitCollation_name(ctx *Collation_nameContext) {}

// EnterForeign_table is called when production foreign_table is entered.
func (s *BaseSQLiteListener) EnterForeign_table(ctx *Foreign_tableContext) {}

// ExitForeign_table is called when production foreign_table is exited.
func (s *BaseSQLiteListener) ExitForeign_table(ctx *Foreign_tableContext) {}

// EnterIndex_name is called when production index_name is entered.
func (s *BaseSQLiteListener) EnterIndex_name(ctx *Index_nameContext) {}

// ExitIndex_name is called when production index_name is exited.
func (s *BaseSQLiteListener) ExitIndex_name(ctx *Index_nameContext) {}

// EnterTrigger_name is called when production trigger_name is entered.
func (s *BaseSQLiteListener) EnterTrigger_name(ctx *Trigger_nameContext) {}

// ExitTrigger_name is called when production trigger_name is exited.
func (s *BaseSQLiteListener) ExitTrigger_name(ctx *Trigger_nameContext) {}

// EnterView_name is called when production view_name is entered.
func (s *BaseSQLiteListener) EnterView_name(ctx *View_nameContext) {}

// ExitView_name is called when production view_name is exited.
func (s *BaseSQLiteListener) ExitView_name(ctx *View_nameContext) {}

// EnterModule_name is called when production module_name is entered.
func (s *BaseSQLiteListener) EnterModule_name(ctx *Module_nameContext) {}

// ExitModule_name is called when production module_name is exited.
func (s *BaseSQLiteListener) ExitModule_name(ctx *Module_nameContext) {}

// EnterPragma_name is called when production pragma_name is entered.
func (s *BaseSQLiteListener) EnterPragma_name(ctx *Pragma_nameContext) {}

// ExitPragma_name is called when production pragma_name is exited.
func (s *BaseSQLiteListener) ExitPragma_name(ctx *Pragma_nameContext) {}

// EnterSavepoint_name is called when production savepoint_name is entered.
func (s *BaseSQLiteListener) EnterSavepoint_name(ctx *Savepoint_nameContext) {}

// ExitSavepoint_name is called when production savepoint_name is exited.
func (s *BaseSQLiteListener) ExitSavepoint_name(ctx *Savepoint_nameContext) {}

// EnterTable_alias is called when production table_alias is entered.
func (s *BaseSQLiteListener) EnterTable_alias(ctx *Table_aliasContext) {}

// ExitTable_alias is called when production table_alias is exited.
func (s *BaseSQLiteListener) ExitTable_alias(ctx *Table_aliasContext) {}

// EnterTransaction_name is called when production transaction_name is entered.
func (s *BaseSQLiteListener) EnterTransaction_name(ctx *Transaction_nameContext) {}

// ExitTransaction_name is called when production transaction_name is exited.
func (s *BaseSQLiteListener) ExitTransaction_name(ctx *Transaction_nameContext) {}

// EnterAny_name is called when production any_name is entered.
func (s *BaseSQLiteListener) EnterAny_name(ctx *Any_nameContext) {}

// ExitAny_name is called when production any_name is exited.
func (s *BaseSQLiteListener) ExitAny_name(ctx *Any_nameContext) {}
