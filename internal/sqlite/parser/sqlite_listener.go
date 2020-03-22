// Code generated from SQLite.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // SQLite

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SQLiteListener is a complete listener for a parse tree produced by SQLiteParser.
type SQLiteListener interface {
	antlr.ParseTreeListener

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterSql_stmt_list is called when entering the sql_stmt_list production.
	EnterSql_stmt_list(c *Sql_stmt_listContext)

	// EnterSql_stmt is called when entering the sql_stmt production.
	EnterSql_stmt(c *Sql_stmtContext)

	// EnterAlter_table_stmt is called when entering the alter_table_stmt production.
	EnterAlter_table_stmt(c *Alter_table_stmtContext)

	// EnterAnalyze_stmt is called when entering the analyze_stmt production.
	EnterAnalyze_stmt(c *Analyze_stmtContext)

	// EnterAttach_stmt is called when entering the attach_stmt production.
	EnterAttach_stmt(c *Attach_stmtContext)

	// EnterBegin_stmt is called when entering the begin_stmt production.
	EnterBegin_stmt(c *Begin_stmtContext)

	// EnterCommit_stmt is called when entering the commit_stmt production.
	EnterCommit_stmt(c *Commit_stmtContext)

	// EnterCompound_select_stmt is called when entering the compound_select_stmt production.
	EnterCompound_select_stmt(c *Compound_select_stmtContext)

	// EnterCreate_index_stmt is called when entering the create_index_stmt production.
	EnterCreate_index_stmt(c *Create_index_stmtContext)

	// EnterCreate_table_stmt is called when entering the create_table_stmt production.
	EnterCreate_table_stmt(c *Create_table_stmtContext)

	// EnterCreate_trigger_stmt is called when entering the create_trigger_stmt production.
	EnterCreate_trigger_stmt(c *Create_trigger_stmtContext)

	// EnterCreate_view_stmt is called when entering the create_view_stmt production.
	EnterCreate_view_stmt(c *Create_view_stmtContext)

	// EnterCreate_virtual_table_stmt is called when entering the create_virtual_table_stmt production.
	EnterCreate_virtual_table_stmt(c *Create_virtual_table_stmtContext)

	// EnterDelete_stmt is called when entering the delete_stmt production.
	EnterDelete_stmt(c *Delete_stmtContext)

	// EnterDelete_stmt_limited is called when entering the delete_stmt_limited production.
	EnterDelete_stmt_limited(c *Delete_stmt_limitedContext)

	// EnterDetach_stmt is called when entering the detach_stmt production.
	EnterDetach_stmt(c *Detach_stmtContext)

	// EnterDrop_index_stmt is called when entering the drop_index_stmt production.
	EnterDrop_index_stmt(c *Drop_index_stmtContext)

	// EnterDrop_table_stmt is called when entering the drop_table_stmt production.
	EnterDrop_table_stmt(c *Drop_table_stmtContext)

	// EnterDrop_trigger_stmt is called when entering the drop_trigger_stmt production.
	EnterDrop_trigger_stmt(c *Drop_trigger_stmtContext)

	// EnterDrop_view_stmt is called when entering the drop_view_stmt production.
	EnterDrop_view_stmt(c *Drop_view_stmtContext)

	// EnterFactored_select_stmt is called when entering the factored_select_stmt production.
	EnterFactored_select_stmt(c *Factored_select_stmtContext)

	// EnterInsert_stmt is called when entering the insert_stmt production.
	EnterInsert_stmt(c *Insert_stmtContext)

	// EnterPragma_stmt is called when entering the pragma_stmt production.
	EnterPragma_stmt(c *Pragma_stmtContext)

	// EnterReindex_stmt is called when entering the reindex_stmt production.
	EnterReindex_stmt(c *Reindex_stmtContext)

	// EnterRelease_stmt is called when entering the release_stmt production.
	EnterRelease_stmt(c *Release_stmtContext)

	// EnterRollback_stmt is called when entering the rollback_stmt production.
	EnterRollback_stmt(c *Rollback_stmtContext)

	// EnterSavepoint_stmt is called when entering the savepoint_stmt production.
	EnterSavepoint_stmt(c *Savepoint_stmtContext)

	// EnterSimple_select_stmt is called when entering the simple_select_stmt production.
	EnterSimple_select_stmt(c *Simple_select_stmtContext)

	// EnterSelect_stmt is called when entering the select_stmt production.
	EnterSelect_stmt(c *Select_stmtContext)

	// EnterSelect_or_values is called when entering the select_or_values production.
	EnterSelect_or_values(c *Select_or_valuesContext)

	// EnterUpdate_stmt is called when entering the update_stmt production.
	EnterUpdate_stmt(c *Update_stmtContext)

	// EnterUpdate_stmt_limited is called when entering the update_stmt_limited production.
	EnterUpdate_stmt_limited(c *Update_stmt_limitedContext)

	// EnterVacuum_stmt is called when entering the vacuum_stmt production.
	EnterVacuum_stmt(c *Vacuum_stmtContext)

	// EnterColumn_def is called when entering the column_def production.
	EnterColumn_def(c *Column_defContext)

	// EnterType_name is called when entering the type_name production.
	EnterType_name(c *Type_nameContext)

	// EnterColumn_constraint is called when entering the column_constraint production.
	EnterColumn_constraint(c *Column_constraintContext)

	// EnterConflict_clause is called when entering the conflict_clause production.
	EnterConflict_clause(c *Conflict_clauseContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterForeign_key_clause is called when entering the foreign_key_clause production.
	EnterForeign_key_clause(c *Foreign_key_clauseContext)

	// EnterRaise_function is called when entering the raise_function production.
	EnterRaise_function(c *Raise_functionContext)

	// EnterIndexed_column is called when entering the indexed_column production.
	EnterIndexed_column(c *Indexed_columnContext)

	// EnterTable_constraint is called when entering the table_constraint production.
	EnterTable_constraint(c *Table_constraintContext)

	// EnterWith_clause is called when entering the with_clause production.
	EnterWith_clause(c *With_clauseContext)

	// EnterQualified_table_name is called when entering the qualified_table_name production.
	EnterQualified_table_name(c *Qualified_table_nameContext)

	// EnterOrdering_term is called when entering the ordering_term production.
	EnterOrdering_term(c *Ordering_termContext)

	// EnterPragma_value is called when entering the pragma_value production.
	EnterPragma_value(c *Pragma_valueContext)

	// EnterCommon_table_expression is called when entering the common_table_expression production.
	EnterCommon_table_expression(c *Common_table_expressionContext)

	// EnterResult_column is called when entering the result_column production.
	EnterResult_column(c *Result_columnContext)

	// EnterTable_or_subquery is called when entering the table_or_subquery production.
	EnterTable_or_subquery(c *Table_or_subqueryContext)

	// EnterJoin_clause is called when entering the join_clause production.
	EnterJoin_clause(c *Join_clauseContext)

	// EnterJoin_operator is called when entering the join_operator production.
	EnterJoin_operator(c *Join_operatorContext)

	// EnterJoin_constraint is called when entering the join_constraint production.
	EnterJoin_constraint(c *Join_constraintContext)

	// EnterSelect_core is called when entering the select_core production.
	EnterSelect_core(c *Select_coreContext)

	// EnterCompound_operator is called when entering the compound_operator production.
	EnterCompound_operator(c *Compound_operatorContext)

	// EnterSigned_number is called when entering the signed_number production.
	EnterSigned_number(c *Signed_numberContext)

	// EnterLiteral_value is called when entering the literal_value production.
	EnterLiteral_value(c *Literal_valueContext)

	// EnterUnary_operator is called when entering the unary_operator production.
	EnterUnary_operator(c *Unary_operatorContext)

	// EnterError_message is called when entering the error_message production.
	EnterError_message(c *Error_messageContext)

	// EnterModule_argument is called when entering the module_argument production.
	EnterModule_argument(c *Module_argumentContext)

	// EnterColumn_alias is called when entering the column_alias production.
	EnterColumn_alias(c *Column_aliasContext)

	// EnterKeyword is called when entering the keyword production.
	EnterKeyword(c *KeywordContext)

	// EnterName is called when entering the name production.
	EnterName(c *NameContext)

	// EnterFunction_name is called when entering the function_name production.
	EnterFunction_name(c *Function_nameContext)

	// EnterDatabase_name is called when entering the database_name production.
	EnterDatabase_name(c *Database_nameContext)

	// EnterSchema_name is called when entering the schema_name production.
	EnterSchema_name(c *Schema_nameContext)

	// EnterTable_function_name is called when entering the table_function_name production.
	EnterTable_function_name(c *Table_function_nameContext)

	// EnterTable_name is called when entering the table_name production.
	EnterTable_name(c *Table_nameContext)

	// EnterTable_or_index_name is called when entering the table_or_index_name production.
	EnterTable_or_index_name(c *Table_or_index_nameContext)

	// EnterNew_table_name is called when entering the new_table_name production.
	EnterNew_table_name(c *New_table_nameContext)

	// EnterColumn_name is called when entering the column_name production.
	EnterColumn_name(c *Column_nameContext)

	// EnterNew_column_name is called when entering the new_column_name production.
	EnterNew_column_name(c *New_column_nameContext)

	// EnterCollation_name is called when entering the collation_name production.
	EnterCollation_name(c *Collation_nameContext)

	// EnterForeign_table is called when entering the foreign_table production.
	EnterForeign_table(c *Foreign_tableContext)

	// EnterIndex_name is called when entering the index_name production.
	EnterIndex_name(c *Index_nameContext)

	// EnterTrigger_name is called when entering the trigger_name production.
	EnterTrigger_name(c *Trigger_nameContext)

	// EnterView_name is called when entering the view_name production.
	EnterView_name(c *View_nameContext)

	// EnterModule_name is called when entering the module_name production.
	EnterModule_name(c *Module_nameContext)

	// EnterPragma_name is called when entering the pragma_name production.
	EnterPragma_name(c *Pragma_nameContext)

	// EnterSavepoint_name is called when entering the savepoint_name production.
	EnterSavepoint_name(c *Savepoint_nameContext)

	// EnterTable_alias is called when entering the table_alias production.
	EnterTable_alias(c *Table_aliasContext)

	// EnterTransaction_name is called when entering the transaction_name production.
	EnterTransaction_name(c *Transaction_nameContext)

	// EnterAny_name is called when entering the any_name production.
	EnterAny_name(c *Any_nameContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitSql_stmt_list is called when exiting the sql_stmt_list production.
	ExitSql_stmt_list(c *Sql_stmt_listContext)

	// ExitSql_stmt is called when exiting the sql_stmt production.
	ExitSql_stmt(c *Sql_stmtContext)

	// ExitAlter_table_stmt is called when exiting the alter_table_stmt production.
	ExitAlter_table_stmt(c *Alter_table_stmtContext)

	// ExitAnalyze_stmt is called when exiting the analyze_stmt production.
	ExitAnalyze_stmt(c *Analyze_stmtContext)

	// ExitAttach_stmt is called when exiting the attach_stmt production.
	ExitAttach_stmt(c *Attach_stmtContext)

	// ExitBegin_stmt is called when exiting the begin_stmt production.
	ExitBegin_stmt(c *Begin_stmtContext)

	// ExitCommit_stmt is called when exiting the commit_stmt production.
	ExitCommit_stmt(c *Commit_stmtContext)

	// ExitCompound_select_stmt is called when exiting the compound_select_stmt production.
	ExitCompound_select_stmt(c *Compound_select_stmtContext)

	// ExitCreate_index_stmt is called when exiting the create_index_stmt production.
	ExitCreate_index_stmt(c *Create_index_stmtContext)

	// ExitCreate_table_stmt is called when exiting the create_table_stmt production.
	ExitCreate_table_stmt(c *Create_table_stmtContext)

	// ExitCreate_trigger_stmt is called when exiting the create_trigger_stmt production.
	ExitCreate_trigger_stmt(c *Create_trigger_stmtContext)

	// ExitCreate_view_stmt is called when exiting the create_view_stmt production.
	ExitCreate_view_stmt(c *Create_view_stmtContext)

	// ExitCreate_virtual_table_stmt is called when exiting the create_virtual_table_stmt production.
	ExitCreate_virtual_table_stmt(c *Create_virtual_table_stmtContext)

	// ExitDelete_stmt is called when exiting the delete_stmt production.
	ExitDelete_stmt(c *Delete_stmtContext)

	// ExitDelete_stmt_limited is called when exiting the delete_stmt_limited production.
	ExitDelete_stmt_limited(c *Delete_stmt_limitedContext)

	// ExitDetach_stmt is called when exiting the detach_stmt production.
	ExitDetach_stmt(c *Detach_stmtContext)

	// ExitDrop_index_stmt is called when exiting the drop_index_stmt production.
	ExitDrop_index_stmt(c *Drop_index_stmtContext)

	// ExitDrop_table_stmt is called when exiting the drop_table_stmt production.
	ExitDrop_table_stmt(c *Drop_table_stmtContext)

	// ExitDrop_trigger_stmt is called when exiting the drop_trigger_stmt production.
	ExitDrop_trigger_stmt(c *Drop_trigger_stmtContext)

	// ExitDrop_view_stmt is called when exiting the drop_view_stmt production.
	ExitDrop_view_stmt(c *Drop_view_stmtContext)

	// ExitFactored_select_stmt is called when exiting the factored_select_stmt production.
	ExitFactored_select_stmt(c *Factored_select_stmtContext)

	// ExitInsert_stmt is called when exiting the insert_stmt production.
	ExitInsert_stmt(c *Insert_stmtContext)

	// ExitPragma_stmt is called when exiting the pragma_stmt production.
	ExitPragma_stmt(c *Pragma_stmtContext)

	// ExitReindex_stmt is called when exiting the reindex_stmt production.
	ExitReindex_stmt(c *Reindex_stmtContext)

	// ExitRelease_stmt is called when exiting the release_stmt production.
	ExitRelease_stmt(c *Release_stmtContext)

	// ExitRollback_stmt is called when exiting the rollback_stmt production.
	ExitRollback_stmt(c *Rollback_stmtContext)

	// ExitSavepoint_stmt is called when exiting the savepoint_stmt production.
	ExitSavepoint_stmt(c *Savepoint_stmtContext)

	// ExitSimple_select_stmt is called when exiting the simple_select_stmt production.
	ExitSimple_select_stmt(c *Simple_select_stmtContext)

	// ExitSelect_stmt is called when exiting the select_stmt production.
	ExitSelect_stmt(c *Select_stmtContext)

	// ExitSelect_or_values is called when exiting the select_or_values production.
	ExitSelect_or_values(c *Select_or_valuesContext)

	// ExitUpdate_stmt is called when exiting the update_stmt production.
	ExitUpdate_stmt(c *Update_stmtContext)

	// ExitUpdate_stmt_limited is called when exiting the update_stmt_limited production.
	ExitUpdate_stmt_limited(c *Update_stmt_limitedContext)

	// ExitVacuum_stmt is called when exiting the vacuum_stmt production.
	ExitVacuum_stmt(c *Vacuum_stmtContext)

	// ExitColumn_def is called when exiting the column_def production.
	ExitColumn_def(c *Column_defContext)

	// ExitType_name is called when exiting the type_name production.
	ExitType_name(c *Type_nameContext)

	// ExitColumn_constraint is called when exiting the column_constraint production.
	ExitColumn_constraint(c *Column_constraintContext)

	// ExitConflict_clause is called when exiting the conflict_clause production.
	ExitConflict_clause(c *Conflict_clauseContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitForeign_key_clause is called when exiting the foreign_key_clause production.
	ExitForeign_key_clause(c *Foreign_key_clauseContext)

	// ExitRaise_function is called when exiting the raise_function production.
	ExitRaise_function(c *Raise_functionContext)

	// ExitIndexed_column is called when exiting the indexed_column production.
	ExitIndexed_column(c *Indexed_columnContext)

	// ExitTable_constraint is called when exiting the table_constraint production.
	ExitTable_constraint(c *Table_constraintContext)

	// ExitWith_clause is called when exiting the with_clause production.
	ExitWith_clause(c *With_clauseContext)

	// ExitQualified_table_name is called when exiting the qualified_table_name production.
	ExitQualified_table_name(c *Qualified_table_nameContext)

	// ExitOrdering_term is called when exiting the ordering_term production.
	ExitOrdering_term(c *Ordering_termContext)

	// ExitPragma_value is called when exiting the pragma_value production.
	ExitPragma_value(c *Pragma_valueContext)

	// ExitCommon_table_expression is called when exiting the common_table_expression production.
	ExitCommon_table_expression(c *Common_table_expressionContext)

	// ExitResult_column is called when exiting the result_column production.
	ExitResult_column(c *Result_columnContext)

	// ExitTable_or_subquery is called when exiting the table_or_subquery production.
	ExitTable_or_subquery(c *Table_or_subqueryContext)

	// ExitJoin_clause is called when exiting the join_clause production.
	ExitJoin_clause(c *Join_clauseContext)

	// ExitJoin_operator is called when exiting the join_operator production.
	ExitJoin_operator(c *Join_operatorContext)

	// ExitJoin_constraint is called when exiting the join_constraint production.
	ExitJoin_constraint(c *Join_constraintContext)

	// ExitSelect_core is called when exiting the select_core production.
	ExitSelect_core(c *Select_coreContext)

	// ExitCompound_operator is called when exiting the compound_operator production.
	ExitCompound_operator(c *Compound_operatorContext)

	// ExitSigned_number is called when exiting the signed_number production.
	ExitSigned_number(c *Signed_numberContext)

	// ExitLiteral_value is called when exiting the literal_value production.
	ExitLiteral_value(c *Literal_valueContext)

	// ExitUnary_operator is called when exiting the unary_operator production.
	ExitUnary_operator(c *Unary_operatorContext)

	// ExitError_message is called when exiting the error_message production.
	ExitError_message(c *Error_messageContext)

	// ExitModule_argument is called when exiting the module_argument production.
	ExitModule_argument(c *Module_argumentContext)

	// ExitColumn_alias is called when exiting the column_alias production.
	ExitColumn_alias(c *Column_aliasContext)

	// ExitKeyword is called when exiting the keyword production.
	ExitKeyword(c *KeywordContext)

	// ExitName is called when exiting the name production.
	ExitName(c *NameContext)

	// ExitFunction_name is called when exiting the function_name production.
	ExitFunction_name(c *Function_nameContext)

	// ExitDatabase_name is called when exiting the database_name production.
	ExitDatabase_name(c *Database_nameContext)

	// ExitSchema_name is called when exiting the schema_name production.
	ExitSchema_name(c *Schema_nameContext)

	// ExitTable_function_name is called when exiting the table_function_name production.
	ExitTable_function_name(c *Table_function_nameContext)

	// ExitTable_name is called when exiting the table_name production.
	ExitTable_name(c *Table_nameContext)

	// ExitTable_or_index_name is called when exiting the table_or_index_name production.
	ExitTable_or_index_name(c *Table_or_index_nameContext)

	// ExitNew_table_name is called when exiting the new_table_name production.
	ExitNew_table_name(c *New_table_nameContext)

	// ExitColumn_name is called when exiting the column_name production.
	ExitColumn_name(c *Column_nameContext)

	// ExitNew_column_name is called when exiting the new_column_name production.
	ExitNew_column_name(c *New_column_nameContext)

	// ExitCollation_name is called when exiting the collation_name production.
	ExitCollation_name(c *Collation_nameContext)

	// ExitForeign_table is called when exiting the foreign_table production.
	ExitForeign_table(c *Foreign_tableContext)

	// ExitIndex_name is called when exiting the index_name production.
	ExitIndex_name(c *Index_nameContext)

	// ExitTrigger_name is called when exiting the trigger_name production.
	ExitTrigger_name(c *Trigger_nameContext)

	// ExitView_name is called when exiting the view_name production.
	ExitView_name(c *View_nameContext)

	// ExitModule_name is called when exiting the module_name production.
	ExitModule_name(c *Module_nameContext)

	// ExitPragma_name is called when exiting the pragma_name production.
	ExitPragma_name(c *Pragma_nameContext)

	// ExitSavepoint_name is called when exiting the savepoint_name production.
	ExitSavepoint_name(c *Savepoint_nameContext)

	// ExitTable_alias is called when exiting the table_alias production.
	ExitTable_alias(c *Table_aliasContext)

	// ExitTransaction_name is called when exiting the transaction_name production.
	ExitTransaction_name(c *Transaction_nameContext)

	// ExitAny_name is called when exiting the any_name production.
	ExitAny_name(c *Any_nameContext)
}
