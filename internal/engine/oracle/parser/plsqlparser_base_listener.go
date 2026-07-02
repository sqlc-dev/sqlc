// Code generated from PlSqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // PlSqlParser
import "github.com/antlr4-go/antlr/v4"

// BasePlSqlParserListener is a complete listener for a parse tree produced by PlSqlParser.
type BasePlSqlParserListener struct{}

var _ PlSqlParserListener = &BasePlSqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasePlSqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasePlSqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasePlSqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasePlSqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSql_script is called when production sql_script is entered.
func (s *BasePlSqlParserListener) EnterSql_script(ctx *Sql_scriptContext) {}

// ExitSql_script is called when production sql_script is exited.
func (s *BasePlSqlParserListener) ExitSql_script(ctx *Sql_scriptContext) {}

// EnterUnit_statement is called when production unit_statement is entered.
func (s *BasePlSqlParserListener) EnterUnit_statement(ctx *Unit_statementContext) {}

// ExitUnit_statement is called when production unit_statement is exited.
func (s *BasePlSqlParserListener) ExitUnit_statement(ctx *Unit_statementContext) {}

// EnterAlter_diskgroup is called when production alter_diskgroup is entered.
func (s *BasePlSqlParserListener) EnterAlter_diskgroup(ctx *Alter_diskgroupContext) {}

// ExitAlter_diskgroup is called when production alter_diskgroup is exited.
func (s *BasePlSqlParserListener) ExitAlter_diskgroup(ctx *Alter_diskgroupContext) {}

// EnterAdd_disk_clause is called when production add_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_disk_clause(ctx *Add_disk_clauseContext) {}

// ExitAdd_disk_clause is called when production add_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_disk_clause(ctx *Add_disk_clauseContext) {}

// EnterDrop_disk_clause is called when production drop_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_disk_clause(ctx *Drop_disk_clauseContext) {}

// ExitDrop_disk_clause is called when production drop_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_disk_clause(ctx *Drop_disk_clauseContext) {}

// EnterResize_disk_clause is called when production resize_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterResize_disk_clause(ctx *Resize_disk_clauseContext) {}

// ExitResize_disk_clause is called when production resize_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitResize_disk_clause(ctx *Resize_disk_clauseContext) {}

// EnterReplace_disk_clause is called when production replace_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterReplace_disk_clause(ctx *Replace_disk_clauseContext) {}

// ExitReplace_disk_clause is called when production replace_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitReplace_disk_clause(ctx *Replace_disk_clauseContext) {}

// EnterWait_nowait is called when production wait_nowait is entered.
func (s *BasePlSqlParserListener) EnterWait_nowait(ctx *Wait_nowaitContext) {}

// ExitWait_nowait is called when production wait_nowait is exited.
func (s *BasePlSqlParserListener) ExitWait_nowait(ctx *Wait_nowaitContext) {}

// EnterRename_disk_clause is called when production rename_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterRename_disk_clause(ctx *Rename_disk_clauseContext) {}

// ExitRename_disk_clause is called when production rename_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitRename_disk_clause(ctx *Rename_disk_clauseContext) {}

// EnterDisk_online_clause is called when production disk_online_clause is entered.
func (s *BasePlSqlParserListener) EnterDisk_online_clause(ctx *Disk_online_clauseContext) {}

// ExitDisk_online_clause is called when production disk_online_clause is exited.
func (s *BasePlSqlParserListener) ExitDisk_online_clause(ctx *Disk_online_clauseContext) {}

// EnterDisk_offline_clause is called when production disk_offline_clause is entered.
func (s *BasePlSqlParserListener) EnterDisk_offline_clause(ctx *Disk_offline_clauseContext) {}

// ExitDisk_offline_clause is called when production disk_offline_clause is exited.
func (s *BasePlSqlParserListener) ExitDisk_offline_clause(ctx *Disk_offline_clauseContext) {}

// EnterTimeout_clause is called when production timeout_clause is entered.
func (s *BasePlSqlParserListener) EnterTimeout_clause(ctx *Timeout_clauseContext) {}

// ExitTimeout_clause is called when production timeout_clause is exited.
func (s *BasePlSqlParserListener) ExitTimeout_clause(ctx *Timeout_clauseContext) {}

// EnterRebalance_diskgroup_clause is called when production rebalance_diskgroup_clause is entered.
func (s *BasePlSqlParserListener) EnterRebalance_diskgroup_clause(ctx *Rebalance_diskgroup_clauseContext) {
}

// ExitRebalance_diskgroup_clause is called when production rebalance_diskgroup_clause is exited.
func (s *BasePlSqlParserListener) ExitRebalance_diskgroup_clause(ctx *Rebalance_diskgroup_clauseContext) {
}

// EnterPhase is called when production phase is entered.
func (s *BasePlSqlParserListener) EnterPhase(ctx *PhaseContext) {}

// ExitPhase is called when production phase is exited.
func (s *BasePlSqlParserListener) ExitPhase(ctx *PhaseContext) {}

// EnterCheck_diskgroup_clause is called when production check_diskgroup_clause is entered.
func (s *BasePlSqlParserListener) EnterCheck_diskgroup_clause(ctx *Check_diskgroup_clauseContext) {}

// ExitCheck_diskgroup_clause is called when production check_diskgroup_clause is exited.
func (s *BasePlSqlParserListener) ExitCheck_diskgroup_clause(ctx *Check_diskgroup_clauseContext) {}

// EnterDiskgroup_template_clauses is called when production diskgroup_template_clauses is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_template_clauses(ctx *Diskgroup_template_clausesContext) {
}

// ExitDiskgroup_template_clauses is called when production diskgroup_template_clauses is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_template_clauses(ctx *Diskgroup_template_clausesContext) {
}

// EnterQualified_template_clause is called when production qualified_template_clause is entered.
func (s *BasePlSqlParserListener) EnterQualified_template_clause(ctx *Qualified_template_clauseContext) {
}

// ExitQualified_template_clause is called when production qualified_template_clause is exited.
func (s *BasePlSqlParserListener) ExitQualified_template_clause(ctx *Qualified_template_clauseContext) {
}

// EnterRedundancy_clause is called when production redundancy_clause is entered.
func (s *BasePlSqlParserListener) EnterRedundancy_clause(ctx *Redundancy_clauseContext) {}

// ExitRedundancy_clause is called when production redundancy_clause is exited.
func (s *BasePlSqlParserListener) ExitRedundancy_clause(ctx *Redundancy_clauseContext) {}

// EnterStriping_clause is called when production striping_clause is entered.
func (s *BasePlSqlParserListener) EnterStriping_clause(ctx *Striping_clauseContext) {}

// ExitStriping_clause is called when production striping_clause is exited.
func (s *BasePlSqlParserListener) ExitStriping_clause(ctx *Striping_clauseContext) {}

// EnterForce_noforce is called when production force_noforce is entered.
func (s *BasePlSqlParserListener) EnterForce_noforce(ctx *Force_noforceContext) {}

// ExitForce_noforce is called when production force_noforce is exited.
func (s *BasePlSqlParserListener) ExitForce_noforce(ctx *Force_noforceContext) {}

// EnterDiskgroup_directory_clauses is called when production diskgroup_directory_clauses is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_directory_clauses(ctx *Diskgroup_directory_clausesContext) {
}

// ExitDiskgroup_directory_clauses is called when production diskgroup_directory_clauses is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_directory_clauses(ctx *Diskgroup_directory_clausesContext) {
}

// EnterDir_name is called when production dir_name is entered.
func (s *BasePlSqlParserListener) EnterDir_name(ctx *Dir_nameContext) {}

// ExitDir_name is called when production dir_name is exited.
func (s *BasePlSqlParserListener) ExitDir_name(ctx *Dir_nameContext) {}

// EnterDiskgroup_alias_clauses is called when production diskgroup_alias_clauses is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_alias_clauses(ctx *Diskgroup_alias_clausesContext) {}

// ExitDiskgroup_alias_clauses is called when production diskgroup_alias_clauses is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_alias_clauses(ctx *Diskgroup_alias_clausesContext) {}

// EnterDiskgroup_volume_clauses is called when production diskgroup_volume_clauses is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_volume_clauses(ctx *Diskgroup_volume_clausesContext) {
}

// ExitDiskgroup_volume_clauses is called when production diskgroup_volume_clauses is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_volume_clauses(ctx *Diskgroup_volume_clausesContext) {
}

// EnterAdd_volume_clause is called when production add_volume_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_volume_clause(ctx *Add_volume_clauseContext) {}

// ExitAdd_volume_clause is called when production add_volume_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_volume_clause(ctx *Add_volume_clauseContext) {}

// EnterModify_volume_clause is called when production modify_volume_clause is entered.
func (s *BasePlSqlParserListener) EnterModify_volume_clause(ctx *Modify_volume_clauseContext) {}

// ExitModify_volume_clause is called when production modify_volume_clause is exited.
func (s *BasePlSqlParserListener) ExitModify_volume_clause(ctx *Modify_volume_clauseContext) {}

// EnterDiskgroup_attributes is called when production diskgroup_attributes is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_attributes(ctx *Diskgroup_attributesContext) {}

// ExitDiskgroup_attributes is called when production diskgroup_attributes is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_attributes(ctx *Diskgroup_attributesContext) {}

// EnterDrop_diskgroup_file_clause is called when production drop_diskgroup_file_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_diskgroup_file_clause(ctx *Drop_diskgroup_file_clauseContext) {
}

// ExitDrop_diskgroup_file_clause is called when production drop_diskgroup_file_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_diskgroup_file_clause(ctx *Drop_diskgroup_file_clauseContext) {
}

// EnterConvert_redundancy_clause is called when production convert_redundancy_clause is entered.
func (s *BasePlSqlParserListener) EnterConvert_redundancy_clause(ctx *Convert_redundancy_clauseContext) {
}

// ExitConvert_redundancy_clause is called when production convert_redundancy_clause is exited.
func (s *BasePlSqlParserListener) ExitConvert_redundancy_clause(ctx *Convert_redundancy_clauseContext) {
}

// EnterUsergroup_clauses is called when production usergroup_clauses is entered.
func (s *BasePlSqlParserListener) EnterUsergroup_clauses(ctx *Usergroup_clausesContext) {}

// ExitUsergroup_clauses is called when production usergroup_clauses is exited.
func (s *BasePlSqlParserListener) ExitUsergroup_clauses(ctx *Usergroup_clausesContext) {}

// EnterUser_clauses is called when production user_clauses is entered.
func (s *BasePlSqlParserListener) EnterUser_clauses(ctx *User_clausesContext) {}

// ExitUser_clauses is called when production user_clauses is exited.
func (s *BasePlSqlParserListener) ExitUser_clauses(ctx *User_clausesContext) {}

// EnterFile_permissions_clause is called when production file_permissions_clause is entered.
func (s *BasePlSqlParserListener) EnterFile_permissions_clause(ctx *File_permissions_clauseContext) {}

// ExitFile_permissions_clause is called when production file_permissions_clause is exited.
func (s *BasePlSqlParserListener) ExitFile_permissions_clause(ctx *File_permissions_clauseContext) {}

// EnterFile_owner_clause is called when production file_owner_clause is entered.
func (s *BasePlSqlParserListener) EnterFile_owner_clause(ctx *File_owner_clauseContext) {}

// ExitFile_owner_clause is called when production file_owner_clause is exited.
func (s *BasePlSqlParserListener) ExitFile_owner_clause(ctx *File_owner_clauseContext) {}

// EnterScrub_clause is called when production scrub_clause is entered.
func (s *BasePlSqlParserListener) EnterScrub_clause(ctx *Scrub_clauseContext) {}

// ExitScrub_clause is called when production scrub_clause is exited.
func (s *BasePlSqlParserListener) ExitScrub_clause(ctx *Scrub_clauseContext) {}

// EnterQuotagroup_clauses is called when production quotagroup_clauses is entered.
func (s *BasePlSqlParserListener) EnterQuotagroup_clauses(ctx *Quotagroup_clausesContext) {}

// ExitQuotagroup_clauses is called when production quotagroup_clauses is exited.
func (s *BasePlSqlParserListener) ExitQuotagroup_clauses(ctx *Quotagroup_clausesContext) {}

// EnterProperty_name is called when production property_name is entered.
func (s *BasePlSqlParserListener) EnterProperty_name(ctx *Property_nameContext) {}

// ExitProperty_name is called when production property_name is exited.
func (s *BasePlSqlParserListener) ExitProperty_name(ctx *Property_nameContext) {}

// EnterProperty_value is called when production property_value is entered.
func (s *BasePlSqlParserListener) EnterProperty_value(ctx *Property_valueContext) {}

// ExitProperty_value is called when production property_value is exited.
func (s *BasePlSqlParserListener) ExitProperty_value(ctx *Property_valueContext) {}

// EnterFilegroup_clauses is called when production filegroup_clauses is entered.
func (s *BasePlSqlParserListener) EnterFilegroup_clauses(ctx *Filegroup_clausesContext) {}

// ExitFilegroup_clauses is called when production filegroup_clauses is exited.
func (s *BasePlSqlParserListener) ExitFilegroup_clauses(ctx *Filegroup_clausesContext) {}

// EnterAdd_filegroup_clause is called when production add_filegroup_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_filegroup_clause(ctx *Add_filegroup_clauseContext) {}

// ExitAdd_filegroup_clause is called when production add_filegroup_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_filegroup_clause(ctx *Add_filegroup_clauseContext) {}

// EnterModify_filegroup_clause is called when production modify_filegroup_clause is entered.
func (s *BasePlSqlParserListener) EnterModify_filegroup_clause(ctx *Modify_filegroup_clauseContext) {}

// ExitModify_filegroup_clause is called when production modify_filegroup_clause is exited.
func (s *BasePlSqlParserListener) ExitModify_filegroup_clause(ctx *Modify_filegroup_clauseContext) {}

// EnterMove_to_filegroup_clause is called when production move_to_filegroup_clause is entered.
func (s *BasePlSqlParserListener) EnterMove_to_filegroup_clause(ctx *Move_to_filegroup_clauseContext) {
}

// ExitMove_to_filegroup_clause is called when production move_to_filegroup_clause is exited.
func (s *BasePlSqlParserListener) ExitMove_to_filegroup_clause(ctx *Move_to_filegroup_clauseContext) {
}

// EnterDrop_filegroup_clause is called when production drop_filegroup_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_filegroup_clause(ctx *Drop_filegroup_clauseContext) {}

// ExitDrop_filegroup_clause is called when production drop_filegroup_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_filegroup_clause(ctx *Drop_filegroup_clauseContext) {}

// EnterQuorum_regular is called when production quorum_regular is entered.
func (s *BasePlSqlParserListener) EnterQuorum_regular(ctx *Quorum_regularContext) {}

// ExitQuorum_regular is called when production quorum_regular is exited.
func (s *BasePlSqlParserListener) ExitQuorum_regular(ctx *Quorum_regularContext) {}

// EnterUndrop_disk_clause is called when production undrop_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterUndrop_disk_clause(ctx *Undrop_disk_clauseContext) {}

// ExitUndrop_disk_clause is called when production undrop_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitUndrop_disk_clause(ctx *Undrop_disk_clauseContext) {}

// EnterDiskgroup_availability is called when production diskgroup_availability is entered.
func (s *BasePlSqlParserListener) EnterDiskgroup_availability(ctx *Diskgroup_availabilityContext) {}

// ExitDiskgroup_availability is called when production diskgroup_availability is exited.
func (s *BasePlSqlParserListener) ExitDiskgroup_availability(ctx *Diskgroup_availabilityContext) {}

// EnterEnable_disable_volume is called when production enable_disable_volume is entered.
func (s *BasePlSqlParserListener) EnterEnable_disable_volume(ctx *Enable_disable_volumeContext) {}

// ExitEnable_disable_volume is called when production enable_disable_volume is exited.
func (s *BasePlSqlParserListener) ExitEnable_disable_volume(ctx *Enable_disable_volumeContext) {}

// EnterDrop_function is called when production drop_function is entered.
func (s *BasePlSqlParserListener) EnterDrop_function(ctx *Drop_functionContext) {}

// ExitDrop_function is called when production drop_function is exited.
func (s *BasePlSqlParserListener) ExitDrop_function(ctx *Drop_functionContext) {}

// EnterAlter_flashback_archive is called when production alter_flashback_archive is entered.
func (s *BasePlSqlParserListener) EnterAlter_flashback_archive(ctx *Alter_flashback_archiveContext) {}

// ExitAlter_flashback_archive is called when production alter_flashback_archive is exited.
func (s *BasePlSqlParserListener) ExitAlter_flashback_archive(ctx *Alter_flashback_archiveContext) {}

// EnterAlter_hierarchy is called when production alter_hierarchy is entered.
func (s *BasePlSqlParserListener) EnterAlter_hierarchy(ctx *Alter_hierarchyContext) {}

// ExitAlter_hierarchy is called when production alter_hierarchy is exited.
func (s *BasePlSqlParserListener) ExitAlter_hierarchy(ctx *Alter_hierarchyContext) {}

// EnterAlter_function is called when production alter_function is entered.
func (s *BasePlSqlParserListener) EnterAlter_function(ctx *Alter_functionContext) {}

// ExitAlter_function is called when production alter_function is exited.
func (s *BasePlSqlParserListener) ExitAlter_function(ctx *Alter_functionContext) {}

// EnterAlter_java is called when production alter_java is entered.
func (s *BasePlSqlParserListener) EnterAlter_java(ctx *Alter_javaContext) {}

// ExitAlter_java is called when production alter_java is exited.
func (s *BasePlSqlParserListener) ExitAlter_java(ctx *Alter_javaContext) {}

// EnterMatch_string is called when production match_string is entered.
func (s *BasePlSqlParserListener) EnterMatch_string(ctx *Match_stringContext) {}

// ExitMatch_string is called when production match_string is exited.
func (s *BasePlSqlParserListener) ExitMatch_string(ctx *Match_stringContext) {}

// EnterCreate_function_body is called when production create_function_body is entered.
func (s *BasePlSqlParserListener) EnterCreate_function_body(ctx *Create_function_bodyContext) {}

// ExitCreate_function_body is called when production create_function_body is exited.
func (s *BasePlSqlParserListener) ExitCreate_function_body(ctx *Create_function_bodyContext) {}

// EnterSql_macro_body is called when production sql_macro_body is entered.
func (s *BasePlSqlParserListener) EnterSql_macro_body(ctx *Sql_macro_bodyContext) {}

// ExitSql_macro_body is called when production sql_macro_body is exited.
func (s *BasePlSqlParserListener) ExitSql_macro_body(ctx *Sql_macro_bodyContext) {}

// EnterParallel_enable_clause is called when production parallel_enable_clause is entered.
func (s *BasePlSqlParserListener) EnterParallel_enable_clause(ctx *Parallel_enable_clauseContext) {}

// ExitParallel_enable_clause is called when production parallel_enable_clause is exited.
func (s *BasePlSqlParserListener) ExitParallel_enable_clause(ctx *Parallel_enable_clauseContext) {}

// EnterPartition_by_clause is called when production partition_by_clause is entered.
func (s *BasePlSqlParserListener) EnterPartition_by_clause(ctx *Partition_by_clauseContext) {}

// ExitPartition_by_clause is called when production partition_by_clause is exited.
func (s *BasePlSqlParserListener) ExitPartition_by_clause(ctx *Partition_by_clauseContext) {}

// EnterResult_cache_clause is called when production result_cache_clause is entered.
func (s *BasePlSqlParserListener) EnterResult_cache_clause(ctx *Result_cache_clauseContext) {}

// ExitResult_cache_clause is called when production result_cache_clause is exited.
func (s *BasePlSqlParserListener) ExitResult_cache_clause(ctx *Result_cache_clauseContext) {}

// EnterAccessible_by_clause is called when production accessible_by_clause is entered.
func (s *BasePlSqlParserListener) EnterAccessible_by_clause(ctx *Accessible_by_clauseContext) {}

// ExitAccessible_by_clause is called when production accessible_by_clause is exited.
func (s *BasePlSqlParserListener) ExitAccessible_by_clause(ctx *Accessible_by_clauseContext) {}

// EnterDefault_collation_clause is called when production default_collation_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_collation_clause(ctx *Default_collation_clauseContext) {
}

// ExitDefault_collation_clause is called when production default_collation_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_collation_clause(ctx *Default_collation_clauseContext) {
}

// EnterAggregate_clause is called when production aggregate_clause is entered.
func (s *BasePlSqlParserListener) EnterAggregate_clause(ctx *Aggregate_clauseContext) {}

// ExitAggregate_clause is called when production aggregate_clause is exited.
func (s *BasePlSqlParserListener) ExitAggregate_clause(ctx *Aggregate_clauseContext) {}

// EnterPipelined_using_clause is called when production pipelined_using_clause is entered.
func (s *BasePlSqlParserListener) EnterPipelined_using_clause(ctx *Pipelined_using_clauseContext) {}

// ExitPipelined_using_clause is called when production pipelined_using_clause is exited.
func (s *BasePlSqlParserListener) ExitPipelined_using_clause(ctx *Pipelined_using_clauseContext) {}

// EnterAccessor is called when production accessor is entered.
func (s *BasePlSqlParserListener) EnterAccessor(ctx *AccessorContext) {}

// ExitAccessor is called when production accessor is exited.
func (s *BasePlSqlParserListener) ExitAccessor(ctx *AccessorContext) {}

// EnterRelies_on_part is called when production relies_on_part is entered.
func (s *BasePlSqlParserListener) EnterRelies_on_part(ctx *Relies_on_partContext) {}

// ExitRelies_on_part is called when production relies_on_part is exited.
func (s *BasePlSqlParserListener) ExitRelies_on_part(ctx *Relies_on_partContext) {}

// EnterStreaming_clause is called when production streaming_clause is entered.
func (s *BasePlSqlParserListener) EnterStreaming_clause(ctx *Streaming_clauseContext) {}

// ExitStreaming_clause is called when production streaming_clause is exited.
func (s *BasePlSqlParserListener) ExitStreaming_clause(ctx *Streaming_clauseContext) {}

// EnterAlter_outline is called when production alter_outline is entered.
func (s *BasePlSqlParserListener) EnterAlter_outline(ctx *Alter_outlineContext) {}

// ExitAlter_outline is called when production alter_outline is exited.
func (s *BasePlSqlParserListener) ExitAlter_outline(ctx *Alter_outlineContext) {}

// EnterOutline_options is called when production outline_options is entered.
func (s *BasePlSqlParserListener) EnterOutline_options(ctx *Outline_optionsContext) {}

// ExitOutline_options is called when production outline_options is exited.
func (s *BasePlSqlParserListener) ExitOutline_options(ctx *Outline_optionsContext) {}

// EnterAlter_lockdown_profile is called when production alter_lockdown_profile is entered.
func (s *BasePlSqlParserListener) EnterAlter_lockdown_profile(ctx *Alter_lockdown_profileContext) {}

// ExitAlter_lockdown_profile is called when production alter_lockdown_profile is exited.
func (s *BasePlSqlParserListener) ExitAlter_lockdown_profile(ctx *Alter_lockdown_profileContext) {}

// EnterLockdown_feature is called when production lockdown_feature is entered.
func (s *BasePlSqlParserListener) EnterLockdown_feature(ctx *Lockdown_featureContext) {}

// ExitLockdown_feature is called when production lockdown_feature is exited.
func (s *BasePlSqlParserListener) ExitLockdown_feature(ctx *Lockdown_featureContext) {}

// EnterLockdown_options is called when production lockdown_options is entered.
func (s *BasePlSqlParserListener) EnterLockdown_options(ctx *Lockdown_optionsContext) {}

// ExitLockdown_options is called when production lockdown_options is exited.
func (s *BasePlSqlParserListener) ExitLockdown_options(ctx *Lockdown_optionsContext) {}

// EnterLockdown_statements is called when production lockdown_statements is entered.
func (s *BasePlSqlParserListener) EnterLockdown_statements(ctx *Lockdown_statementsContext) {}

// ExitLockdown_statements is called when production lockdown_statements is exited.
func (s *BasePlSqlParserListener) ExitLockdown_statements(ctx *Lockdown_statementsContext) {}

// EnterStatement_clauses is called when production statement_clauses is entered.
func (s *BasePlSqlParserListener) EnterStatement_clauses(ctx *Statement_clausesContext) {}

// ExitStatement_clauses is called when production statement_clauses is exited.
func (s *BasePlSqlParserListener) ExitStatement_clauses(ctx *Statement_clausesContext) {}

// EnterClause_options is called when production clause_options is entered.
func (s *BasePlSqlParserListener) EnterClause_options(ctx *Clause_optionsContext) {}

// ExitClause_options is called when production clause_options is exited.
func (s *BasePlSqlParserListener) ExitClause_options(ctx *Clause_optionsContext) {}

// EnterOption_values is called when production option_values is entered.
func (s *BasePlSqlParserListener) EnterOption_values(ctx *Option_valuesContext) {}

// ExitOption_values is called when production option_values is exited.
func (s *BasePlSqlParserListener) ExitOption_values(ctx *Option_valuesContext) {}

// EnterString_list is called when production string_list is entered.
func (s *BasePlSqlParserListener) EnterString_list(ctx *String_listContext) {}

// ExitString_list is called when production string_list is exited.
func (s *BasePlSqlParserListener) ExitString_list(ctx *String_listContext) {}

// EnterDisable_enable is called when production disable_enable is entered.
func (s *BasePlSqlParserListener) EnterDisable_enable(ctx *Disable_enableContext) {}

// ExitDisable_enable is called when production disable_enable is exited.
func (s *BasePlSqlParserListener) ExitDisable_enable(ctx *Disable_enableContext) {}

// EnterDrop_lockdown_profile is called when production drop_lockdown_profile is entered.
func (s *BasePlSqlParserListener) EnterDrop_lockdown_profile(ctx *Drop_lockdown_profileContext) {}

// ExitDrop_lockdown_profile is called when production drop_lockdown_profile is exited.
func (s *BasePlSqlParserListener) ExitDrop_lockdown_profile(ctx *Drop_lockdown_profileContext) {}

// EnterDrop_package is called when production drop_package is entered.
func (s *BasePlSqlParserListener) EnterDrop_package(ctx *Drop_packageContext) {}

// ExitDrop_package is called when production drop_package is exited.
func (s *BasePlSqlParserListener) ExitDrop_package(ctx *Drop_packageContext) {}

// EnterAlter_package is called when production alter_package is entered.
func (s *BasePlSqlParserListener) EnterAlter_package(ctx *Alter_packageContext) {}

// ExitAlter_package is called when production alter_package is exited.
func (s *BasePlSqlParserListener) ExitAlter_package(ctx *Alter_packageContext) {}

// EnterCreate_package is called when production create_package is entered.
func (s *BasePlSqlParserListener) EnterCreate_package(ctx *Create_packageContext) {}

// ExitCreate_package is called when production create_package is exited.
func (s *BasePlSqlParserListener) ExitCreate_package(ctx *Create_packageContext) {}

// EnterCreate_package_body is called when production create_package_body is entered.
func (s *BasePlSqlParserListener) EnterCreate_package_body(ctx *Create_package_bodyContext) {}

// ExitCreate_package_body is called when production create_package_body is exited.
func (s *BasePlSqlParserListener) ExitCreate_package_body(ctx *Create_package_bodyContext) {}

// EnterPackage_obj_spec is called when production package_obj_spec is entered.
func (s *BasePlSqlParserListener) EnterPackage_obj_spec(ctx *Package_obj_specContext) {}

// ExitPackage_obj_spec is called when production package_obj_spec is exited.
func (s *BasePlSqlParserListener) ExitPackage_obj_spec(ctx *Package_obj_specContext) {}

// EnterProcedure_spec is called when production procedure_spec is entered.
func (s *BasePlSqlParserListener) EnterProcedure_spec(ctx *Procedure_specContext) {}

// ExitProcedure_spec is called when production procedure_spec is exited.
func (s *BasePlSqlParserListener) ExitProcedure_spec(ctx *Procedure_specContext) {}

// EnterFunction_spec is called when production function_spec is entered.
func (s *BasePlSqlParserListener) EnterFunction_spec(ctx *Function_specContext) {}

// ExitFunction_spec is called when production function_spec is exited.
func (s *BasePlSqlParserListener) ExitFunction_spec(ctx *Function_specContext) {}

// EnterPackage_obj_body is called when production package_obj_body is entered.
func (s *BasePlSqlParserListener) EnterPackage_obj_body(ctx *Package_obj_bodyContext) {}

// ExitPackage_obj_body is called when production package_obj_body is exited.
func (s *BasePlSqlParserListener) ExitPackage_obj_body(ctx *Package_obj_bodyContext) {}

// EnterAlter_pmem_filestore is called when production alter_pmem_filestore is entered.
func (s *BasePlSqlParserListener) EnterAlter_pmem_filestore(ctx *Alter_pmem_filestoreContext) {}

// ExitAlter_pmem_filestore is called when production alter_pmem_filestore is exited.
func (s *BasePlSqlParserListener) ExitAlter_pmem_filestore(ctx *Alter_pmem_filestoreContext) {}

// EnterDrop_pmem_filestore is called when production drop_pmem_filestore is entered.
func (s *BasePlSqlParserListener) EnterDrop_pmem_filestore(ctx *Drop_pmem_filestoreContext) {}

// ExitDrop_pmem_filestore is called when production drop_pmem_filestore is exited.
func (s *BasePlSqlParserListener) ExitDrop_pmem_filestore(ctx *Drop_pmem_filestoreContext) {}

// EnterDrop_procedure is called when production drop_procedure is entered.
func (s *BasePlSqlParserListener) EnterDrop_procedure(ctx *Drop_procedureContext) {}

// ExitDrop_procedure is called when production drop_procedure is exited.
func (s *BasePlSqlParserListener) ExitDrop_procedure(ctx *Drop_procedureContext) {}

// EnterAlter_procedure is called when production alter_procedure is entered.
func (s *BasePlSqlParserListener) EnterAlter_procedure(ctx *Alter_procedureContext) {}

// ExitAlter_procedure is called when production alter_procedure is exited.
func (s *BasePlSqlParserListener) ExitAlter_procedure(ctx *Alter_procedureContext) {}

// EnterFunction_body is called when production function_body is entered.
func (s *BasePlSqlParserListener) EnterFunction_body(ctx *Function_bodyContext) {}

// ExitFunction_body is called when production function_body is exited.
func (s *BasePlSqlParserListener) ExitFunction_body(ctx *Function_bodyContext) {}

// EnterProcedure_body is called when production procedure_body is entered.
func (s *BasePlSqlParserListener) EnterProcedure_body(ctx *Procedure_bodyContext) {}

// ExitProcedure_body is called when production procedure_body is exited.
func (s *BasePlSqlParserListener) ExitProcedure_body(ctx *Procedure_bodyContext) {}

// EnterCreate_procedure_body is called when production create_procedure_body is entered.
func (s *BasePlSqlParserListener) EnterCreate_procedure_body(ctx *Create_procedure_bodyContext) {}

// ExitCreate_procedure_body is called when production create_procedure_body is exited.
func (s *BasePlSqlParserListener) ExitCreate_procedure_body(ctx *Create_procedure_bodyContext) {}

// EnterAlter_resource_cost is called when production alter_resource_cost is entered.
func (s *BasePlSqlParserListener) EnterAlter_resource_cost(ctx *Alter_resource_costContext) {}

// ExitAlter_resource_cost is called when production alter_resource_cost is exited.
func (s *BasePlSqlParserListener) ExitAlter_resource_cost(ctx *Alter_resource_costContext) {}

// EnterDrop_outline is called when production drop_outline is entered.
func (s *BasePlSqlParserListener) EnterDrop_outline(ctx *Drop_outlineContext) {}

// ExitDrop_outline is called when production drop_outline is exited.
func (s *BasePlSqlParserListener) ExitDrop_outline(ctx *Drop_outlineContext) {}

// EnterAlter_rollback_segment is called when production alter_rollback_segment is entered.
func (s *BasePlSqlParserListener) EnterAlter_rollback_segment(ctx *Alter_rollback_segmentContext) {}

// ExitAlter_rollback_segment is called when production alter_rollback_segment is exited.
func (s *BasePlSqlParserListener) ExitAlter_rollback_segment(ctx *Alter_rollback_segmentContext) {}

// EnterDrop_restore_point is called when production drop_restore_point is entered.
func (s *BasePlSqlParserListener) EnterDrop_restore_point(ctx *Drop_restore_pointContext) {}

// ExitDrop_restore_point is called when production drop_restore_point is exited.
func (s *BasePlSqlParserListener) ExitDrop_restore_point(ctx *Drop_restore_pointContext) {}

// EnterDrop_rollback_segment is called when production drop_rollback_segment is entered.
func (s *BasePlSqlParserListener) EnterDrop_rollback_segment(ctx *Drop_rollback_segmentContext) {}

// ExitDrop_rollback_segment is called when production drop_rollback_segment is exited.
func (s *BasePlSqlParserListener) ExitDrop_rollback_segment(ctx *Drop_rollback_segmentContext) {}

// EnterDrop_role is called when production drop_role is entered.
func (s *BasePlSqlParserListener) EnterDrop_role(ctx *Drop_roleContext) {}

// ExitDrop_role is called when production drop_role is exited.
func (s *BasePlSqlParserListener) ExitDrop_role(ctx *Drop_roleContext) {}

// EnterCreate_pmem_filestore is called when production create_pmem_filestore is entered.
func (s *BasePlSqlParserListener) EnterCreate_pmem_filestore(ctx *Create_pmem_filestoreContext) {}

// ExitCreate_pmem_filestore is called when production create_pmem_filestore is exited.
func (s *BasePlSqlParserListener) ExitCreate_pmem_filestore(ctx *Create_pmem_filestoreContext) {}

// EnterPmem_filestore_options is called when production pmem_filestore_options is entered.
func (s *BasePlSqlParserListener) EnterPmem_filestore_options(ctx *Pmem_filestore_optionsContext) {}

// ExitPmem_filestore_options is called when production pmem_filestore_options is exited.
func (s *BasePlSqlParserListener) ExitPmem_filestore_options(ctx *Pmem_filestore_optionsContext) {}

// EnterFile_path is called when production file_path is entered.
func (s *BasePlSqlParserListener) EnterFile_path(ctx *File_pathContext) {}

// ExitFile_path is called when production file_path is exited.
func (s *BasePlSqlParserListener) ExitFile_path(ctx *File_pathContext) {}

// EnterCreate_rollback_segment is called when production create_rollback_segment is entered.
func (s *BasePlSqlParserListener) EnterCreate_rollback_segment(ctx *Create_rollback_segmentContext) {}

// ExitCreate_rollback_segment is called when production create_rollback_segment is exited.
func (s *BasePlSqlParserListener) ExitCreate_rollback_segment(ctx *Create_rollback_segmentContext) {}

// EnterDrop_trigger is called when production drop_trigger is entered.
func (s *BasePlSqlParserListener) EnterDrop_trigger(ctx *Drop_triggerContext) {}

// ExitDrop_trigger is called when production drop_trigger is exited.
func (s *BasePlSqlParserListener) ExitDrop_trigger(ctx *Drop_triggerContext) {}

// EnterAlter_trigger is called when production alter_trigger is entered.
func (s *BasePlSqlParserListener) EnterAlter_trigger(ctx *Alter_triggerContext) {}

// ExitAlter_trigger is called when production alter_trigger is exited.
func (s *BasePlSqlParserListener) ExitAlter_trigger(ctx *Alter_triggerContext) {}

// EnterCreate_trigger is called when production create_trigger is entered.
func (s *BasePlSqlParserListener) EnterCreate_trigger(ctx *Create_triggerContext) {}

// ExitCreate_trigger is called when production create_trigger is exited.
func (s *BasePlSqlParserListener) ExitCreate_trigger(ctx *Create_triggerContext) {}

// EnterTrigger_follows_clause is called when production trigger_follows_clause is entered.
func (s *BasePlSqlParserListener) EnterTrigger_follows_clause(ctx *Trigger_follows_clauseContext) {}

// ExitTrigger_follows_clause is called when production trigger_follows_clause is exited.
func (s *BasePlSqlParserListener) ExitTrigger_follows_clause(ctx *Trigger_follows_clauseContext) {}

// EnterTrigger_when_clause is called when production trigger_when_clause is entered.
func (s *BasePlSqlParserListener) EnterTrigger_when_clause(ctx *Trigger_when_clauseContext) {}

// ExitTrigger_when_clause is called when production trigger_when_clause is exited.
func (s *BasePlSqlParserListener) ExitTrigger_when_clause(ctx *Trigger_when_clauseContext) {}

// EnterSimple_dml_trigger is called when production simple_dml_trigger is entered.
func (s *BasePlSqlParserListener) EnterSimple_dml_trigger(ctx *Simple_dml_triggerContext) {}

// ExitSimple_dml_trigger is called when production simple_dml_trigger is exited.
func (s *BasePlSqlParserListener) ExitSimple_dml_trigger(ctx *Simple_dml_triggerContext) {}

// EnterFor_each_row is called when production for_each_row is entered.
func (s *BasePlSqlParserListener) EnterFor_each_row(ctx *For_each_rowContext) {}

// ExitFor_each_row is called when production for_each_row is exited.
func (s *BasePlSqlParserListener) ExitFor_each_row(ctx *For_each_rowContext) {}

// EnterCompound_dml_trigger is called when production compound_dml_trigger is entered.
func (s *BasePlSqlParserListener) EnterCompound_dml_trigger(ctx *Compound_dml_triggerContext) {}

// ExitCompound_dml_trigger is called when production compound_dml_trigger is exited.
func (s *BasePlSqlParserListener) ExitCompound_dml_trigger(ctx *Compound_dml_triggerContext) {}

// EnterNon_dml_trigger is called when production non_dml_trigger is entered.
func (s *BasePlSqlParserListener) EnterNon_dml_trigger(ctx *Non_dml_triggerContext) {}

// ExitNon_dml_trigger is called when production non_dml_trigger is exited.
func (s *BasePlSqlParserListener) ExitNon_dml_trigger(ctx *Non_dml_triggerContext) {}

// EnterTrigger_body is called when production trigger_body is entered.
func (s *BasePlSqlParserListener) EnterTrigger_body(ctx *Trigger_bodyContext) {}

// ExitTrigger_body is called when production trigger_body is exited.
func (s *BasePlSqlParserListener) ExitTrigger_body(ctx *Trigger_bodyContext) {}

// EnterCompound_trigger_block is called when production compound_trigger_block is entered.
func (s *BasePlSqlParserListener) EnterCompound_trigger_block(ctx *Compound_trigger_blockContext) {}

// ExitCompound_trigger_block is called when production compound_trigger_block is exited.
func (s *BasePlSqlParserListener) ExitCompound_trigger_block(ctx *Compound_trigger_blockContext) {}

// EnterTiming_point_section is called when production timing_point_section is entered.
func (s *BasePlSqlParserListener) EnterTiming_point_section(ctx *Timing_point_sectionContext) {}

// ExitTiming_point_section is called when production timing_point_section is exited.
func (s *BasePlSqlParserListener) ExitTiming_point_section(ctx *Timing_point_sectionContext) {}

// EnterNon_dml_event is called when production non_dml_event is entered.
func (s *BasePlSqlParserListener) EnterNon_dml_event(ctx *Non_dml_eventContext) {}

// ExitNon_dml_event is called when production non_dml_event is exited.
func (s *BasePlSqlParserListener) ExitNon_dml_event(ctx *Non_dml_eventContext) {}

// EnterDml_event_clause is called when production dml_event_clause is entered.
func (s *BasePlSqlParserListener) EnterDml_event_clause(ctx *Dml_event_clauseContext) {}

// ExitDml_event_clause is called when production dml_event_clause is exited.
func (s *BasePlSqlParserListener) ExitDml_event_clause(ctx *Dml_event_clauseContext) {}

// EnterDml_event_element is called when production dml_event_element is entered.
func (s *BasePlSqlParserListener) EnterDml_event_element(ctx *Dml_event_elementContext) {}

// ExitDml_event_element is called when production dml_event_element is exited.
func (s *BasePlSqlParserListener) ExitDml_event_element(ctx *Dml_event_elementContext) {}

// EnterDml_event_nested_clause is called when production dml_event_nested_clause is entered.
func (s *BasePlSqlParserListener) EnterDml_event_nested_clause(ctx *Dml_event_nested_clauseContext) {}

// ExitDml_event_nested_clause is called when production dml_event_nested_clause is exited.
func (s *BasePlSqlParserListener) ExitDml_event_nested_clause(ctx *Dml_event_nested_clauseContext) {}

// EnterReferencing_clause is called when production referencing_clause is entered.
func (s *BasePlSqlParserListener) EnterReferencing_clause(ctx *Referencing_clauseContext) {}

// ExitReferencing_clause is called when production referencing_clause is exited.
func (s *BasePlSqlParserListener) ExitReferencing_clause(ctx *Referencing_clauseContext) {}

// EnterReferencing_element is called when production referencing_element is entered.
func (s *BasePlSqlParserListener) EnterReferencing_element(ctx *Referencing_elementContext) {}

// ExitReferencing_element is called when production referencing_element is exited.
func (s *BasePlSqlParserListener) ExitReferencing_element(ctx *Referencing_elementContext) {}

// EnterDrop_type is called when production drop_type is entered.
func (s *BasePlSqlParserListener) EnterDrop_type(ctx *Drop_typeContext) {}

// ExitDrop_type is called when production drop_type is exited.
func (s *BasePlSqlParserListener) ExitDrop_type(ctx *Drop_typeContext) {}

// EnterAlter_type is called when production alter_type is entered.
func (s *BasePlSqlParserListener) EnterAlter_type(ctx *Alter_typeContext) {}

// ExitAlter_type is called when production alter_type is exited.
func (s *BasePlSqlParserListener) ExitAlter_type(ctx *Alter_typeContext) {}

// EnterCompile_type_clause is called when production compile_type_clause is entered.
func (s *BasePlSqlParserListener) EnterCompile_type_clause(ctx *Compile_type_clauseContext) {}

// ExitCompile_type_clause is called when production compile_type_clause is exited.
func (s *BasePlSqlParserListener) ExitCompile_type_clause(ctx *Compile_type_clauseContext) {}

// EnterReplace_type_clause is called when production replace_type_clause is entered.
func (s *BasePlSqlParserListener) EnterReplace_type_clause(ctx *Replace_type_clauseContext) {}

// ExitReplace_type_clause is called when production replace_type_clause is exited.
func (s *BasePlSqlParserListener) ExitReplace_type_clause(ctx *Replace_type_clauseContext) {}

// EnterAlter_method_spec is called when production alter_method_spec is entered.
func (s *BasePlSqlParserListener) EnterAlter_method_spec(ctx *Alter_method_specContext) {}

// ExitAlter_method_spec is called when production alter_method_spec is exited.
func (s *BasePlSqlParserListener) ExitAlter_method_spec(ctx *Alter_method_specContext) {}

// EnterAlter_method_element is called when production alter_method_element is entered.
func (s *BasePlSqlParserListener) EnterAlter_method_element(ctx *Alter_method_elementContext) {}

// ExitAlter_method_element is called when production alter_method_element is exited.
func (s *BasePlSqlParserListener) ExitAlter_method_element(ctx *Alter_method_elementContext) {}

// EnterAlter_collection_clauses is called when production alter_collection_clauses is entered.
func (s *BasePlSqlParserListener) EnterAlter_collection_clauses(ctx *Alter_collection_clausesContext) {
}

// ExitAlter_collection_clauses is called when production alter_collection_clauses is exited.
func (s *BasePlSqlParserListener) ExitAlter_collection_clauses(ctx *Alter_collection_clausesContext) {
}

// EnterDependent_handling_clause is called when production dependent_handling_clause is entered.
func (s *BasePlSqlParserListener) EnterDependent_handling_clause(ctx *Dependent_handling_clauseContext) {
}

// ExitDependent_handling_clause is called when production dependent_handling_clause is exited.
func (s *BasePlSqlParserListener) ExitDependent_handling_clause(ctx *Dependent_handling_clauseContext) {
}

// EnterDependent_exceptions_part is called when production dependent_exceptions_part is entered.
func (s *BasePlSqlParserListener) EnterDependent_exceptions_part(ctx *Dependent_exceptions_partContext) {
}

// ExitDependent_exceptions_part is called when production dependent_exceptions_part is exited.
func (s *BasePlSqlParserListener) ExitDependent_exceptions_part(ctx *Dependent_exceptions_partContext) {
}

// EnterCreate_type is called when production create_type is entered.
func (s *BasePlSqlParserListener) EnterCreate_type(ctx *Create_typeContext) {}

// ExitCreate_type is called when production create_type is exited.
func (s *BasePlSqlParserListener) ExitCreate_type(ctx *Create_typeContext) {}

// EnterType_definition is called when production type_definition is entered.
func (s *BasePlSqlParserListener) EnterType_definition(ctx *Type_definitionContext) {}

// ExitType_definition is called when production type_definition is exited.
func (s *BasePlSqlParserListener) ExitType_definition(ctx *Type_definitionContext) {}

// EnterObject_type_def is called when production object_type_def is entered.
func (s *BasePlSqlParserListener) EnterObject_type_def(ctx *Object_type_defContext) {}

// ExitObject_type_def is called when production object_type_def is exited.
func (s *BasePlSqlParserListener) ExitObject_type_def(ctx *Object_type_defContext) {}

// EnterObject_as_part is called when production object_as_part is entered.
func (s *BasePlSqlParserListener) EnterObject_as_part(ctx *Object_as_partContext) {}

// ExitObject_as_part is called when production object_as_part is exited.
func (s *BasePlSqlParserListener) ExitObject_as_part(ctx *Object_as_partContext) {}

// EnterObject_under_part is called when production object_under_part is entered.
func (s *BasePlSqlParserListener) EnterObject_under_part(ctx *Object_under_partContext) {}

// ExitObject_under_part is called when production object_under_part is exited.
func (s *BasePlSqlParserListener) ExitObject_under_part(ctx *Object_under_partContext) {}

// EnterNested_table_type_def is called when production nested_table_type_def is entered.
func (s *BasePlSqlParserListener) EnterNested_table_type_def(ctx *Nested_table_type_defContext) {}

// ExitNested_table_type_def is called when production nested_table_type_def is exited.
func (s *BasePlSqlParserListener) ExitNested_table_type_def(ctx *Nested_table_type_defContext) {}

// EnterSqlj_object_type is called when production sqlj_object_type is entered.
func (s *BasePlSqlParserListener) EnterSqlj_object_type(ctx *Sqlj_object_typeContext) {}

// ExitSqlj_object_type is called when production sqlj_object_type is exited.
func (s *BasePlSqlParserListener) ExitSqlj_object_type(ctx *Sqlj_object_typeContext) {}

// EnterType_body is called when production type_body is entered.
func (s *BasePlSqlParserListener) EnterType_body(ctx *Type_bodyContext) {}

// ExitType_body is called when production type_body is exited.
func (s *BasePlSqlParserListener) ExitType_body(ctx *Type_bodyContext) {}

// EnterType_body_elements is called when production type_body_elements is entered.
func (s *BasePlSqlParserListener) EnterType_body_elements(ctx *Type_body_elementsContext) {}

// ExitType_body_elements is called when production type_body_elements is exited.
func (s *BasePlSqlParserListener) ExitType_body_elements(ctx *Type_body_elementsContext) {}

// EnterMap_order_func_declaration is called when production map_order_func_declaration is entered.
func (s *BasePlSqlParserListener) EnterMap_order_func_declaration(ctx *Map_order_func_declarationContext) {
}

// ExitMap_order_func_declaration is called when production map_order_func_declaration is exited.
func (s *BasePlSqlParserListener) ExitMap_order_func_declaration(ctx *Map_order_func_declarationContext) {
}

// EnterSubprog_decl_in_type is called when production subprog_decl_in_type is entered.
func (s *BasePlSqlParserListener) EnterSubprog_decl_in_type(ctx *Subprog_decl_in_typeContext) {}

// ExitSubprog_decl_in_type is called when production subprog_decl_in_type is exited.
func (s *BasePlSqlParserListener) ExitSubprog_decl_in_type(ctx *Subprog_decl_in_typeContext) {}

// EnterProc_decl_in_type is called when production proc_decl_in_type is entered.
func (s *BasePlSqlParserListener) EnterProc_decl_in_type(ctx *Proc_decl_in_typeContext) {}

// ExitProc_decl_in_type is called when production proc_decl_in_type is exited.
func (s *BasePlSqlParserListener) ExitProc_decl_in_type(ctx *Proc_decl_in_typeContext) {}

// EnterFunc_decl_in_type is called when production func_decl_in_type is entered.
func (s *BasePlSqlParserListener) EnterFunc_decl_in_type(ctx *Func_decl_in_typeContext) {}

// ExitFunc_decl_in_type is called when production func_decl_in_type is exited.
func (s *BasePlSqlParserListener) ExitFunc_decl_in_type(ctx *Func_decl_in_typeContext) {}

// EnterConstructor_declaration is called when production constructor_declaration is entered.
func (s *BasePlSqlParserListener) EnterConstructor_declaration(ctx *Constructor_declarationContext) {}

// ExitConstructor_declaration is called when production constructor_declaration is exited.
func (s *BasePlSqlParserListener) ExitConstructor_declaration(ctx *Constructor_declarationContext) {}

// EnterModifier_clause is called when production modifier_clause is entered.
func (s *BasePlSqlParserListener) EnterModifier_clause(ctx *Modifier_clauseContext) {}

// ExitModifier_clause is called when production modifier_clause is exited.
func (s *BasePlSqlParserListener) ExitModifier_clause(ctx *Modifier_clauseContext) {}

// EnterObject_member_spec is called when production object_member_spec is entered.
func (s *BasePlSqlParserListener) EnterObject_member_spec(ctx *Object_member_specContext) {}

// ExitObject_member_spec is called when production object_member_spec is exited.
func (s *BasePlSqlParserListener) ExitObject_member_spec(ctx *Object_member_specContext) {}

// EnterSqlj_object_type_attr is called when production sqlj_object_type_attr is entered.
func (s *BasePlSqlParserListener) EnterSqlj_object_type_attr(ctx *Sqlj_object_type_attrContext) {}

// ExitSqlj_object_type_attr is called when production sqlj_object_type_attr is exited.
func (s *BasePlSqlParserListener) ExitSqlj_object_type_attr(ctx *Sqlj_object_type_attrContext) {}

// EnterElement_spec is called when production element_spec is entered.
func (s *BasePlSqlParserListener) EnterElement_spec(ctx *Element_specContext) {}

// ExitElement_spec is called when production element_spec is exited.
func (s *BasePlSqlParserListener) ExitElement_spec(ctx *Element_specContext) {}

// EnterElement_spec_options is called when production element_spec_options is entered.
func (s *BasePlSqlParserListener) EnterElement_spec_options(ctx *Element_spec_optionsContext) {}

// ExitElement_spec_options is called when production element_spec_options is exited.
func (s *BasePlSqlParserListener) ExitElement_spec_options(ctx *Element_spec_optionsContext) {}

// EnterSubprogram_spec is called when production subprogram_spec is entered.
func (s *BasePlSqlParserListener) EnterSubprogram_spec(ctx *Subprogram_specContext) {}

// ExitSubprogram_spec is called when production subprogram_spec is exited.
func (s *BasePlSqlParserListener) ExitSubprogram_spec(ctx *Subprogram_specContext) {}

// EnterOverriding_subprogram_spec is called when production overriding_subprogram_spec is entered.
func (s *BasePlSqlParserListener) EnterOverriding_subprogram_spec(ctx *Overriding_subprogram_specContext) {
}

// ExitOverriding_subprogram_spec is called when production overriding_subprogram_spec is exited.
func (s *BasePlSqlParserListener) ExitOverriding_subprogram_spec(ctx *Overriding_subprogram_specContext) {
}

// EnterOverriding_function_spec is called when production overriding_function_spec is entered.
func (s *BasePlSqlParserListener) EnterOverriding_function_spec(ctx *Overriding_function_specContext) {
}

// ExitOverriding_function_spec is called when production overriding_function_spec is exited.
func (s *BasePlSqlParserListener) ExitOverriding_function_spec(ctx *Overriding_function_specContext) {
}

// EnterOverriding_procedure_spec is called when production overriding_procedure_spec is entered.
func (s *BasePlSqlParserListener) EnterOverriding_procedure_spec(ctx *Overriding_procedure_specContext) {
}

// ExitOverriding_procedure_spec is called when production overriding_procedure_spec is exited.
func (s *BasePlSqlParserListener) ExitOverriding_procedure_spec(ctx *Overriding_procedure_specContext) {
}

// EnterType_procedure_spec is called when production type_procedure_spec is entered.
func (s *BasePlSqlParserListener) EnterType_procedure_spec(ctx *Type_procedure_specContext) {}

// ExitType_procedure_spec is called when production type_procedure_spec is exited.
func (s *BasePlSqlParserListener) ExitType_procedure_spec(ctx *Type_procedure_specContext) {}

// EnterType_function_spec is called when production type_function_spec is entered.
func (s *BasePlSqlParserListener) EnterType_function_spec(ctx *Type_function_specContext) {}

// ExitType_function_spec is called when production type_function_spec is exited.
func (s *BasePlSqlParserListener) ExitType_function_spec(ctx *Type_function_specContext) {}

// EnterConstructor_spec is called when production constructor_spec is entered.
func (s *BasePlSqlParserListener) EnterConstructor_spec(ctx *Constructor_specContext) {}

// ExitConstructor_spec is called when production constructor_spec is exited.
func (s *BasePlSqlParserListener) ExitConstructor_spec(ctx *Constructor_specContext) {}

// EnterMap_order_function_spec is called when production map_order_function_spec is entered.
func (s *BasePlSqlParserListener) EnterMap_order_function_spec(ctx *Map_order_function_specContext) {}

// ExitMap_order_function_spec is called when production map_order_function_spec is exited.
func (s *BasePlSqlParserListener) ExitMap_order_function_spec(ctx *Map_order_function_specContext) {}

// EnterPragma_clause is called when production pragma_clause is entered.
func (s *BasePlSqlParserListener) EnterPragma_clause(ctx *Pragma_clauseContext) {}

// ExitPragma_clause is called when production pragma_clause is exited.
func (s *BasePlSqlParserListener) ExitPragma_clause(ctx *Pragma_clauseContext) {}

// EnterPragma_elements is called when production pragma_elements is entered.
func (s *BasePlSqlParserListener) EnterPragma_elements(ctx *Pragma_elementsContext) {}

// ExitPragma_elements is called when production pragma_elements is exited.
func (s *BasePlSqlParserListener) ExitPragma_elements(ctx *Pragma_elementsContext) {}

// EnterType_elements_parameter is called when production type_elements_parameter is entered.
func (s *BasePlSqlParserListener) EnterType_elements_parameter(ctx *Type_elements_parameterContext) {}

// ExitType_elements_parameter is called when production type_elements_parameter is exited.
func (s *BasePlSqlParserListener) ExitType_elements_parameter(ctx *Type_elements_parameterContext) {}

// EnterDrop_sequence is called when production drop_sequence is entered.
func (s *BasePlSqlParserListener) EnterDrop_sequence(ctx *Drop_sequenceContext) {}

// ExitDrop_sequence is called when production drop_sequence is exited.
func (s *BasePlSqlParserListener) ExitDrop_sequence(ctx *Drop_sequenceContext) {}

// EnterAlter_sequence is called when production alter_sequence is entered.
func (s *BasePlSqlParserListener) EnterAlter_sequence(ctx *Alter_sequenceContext) {}

// ExitAlter_sequence is called when production alter_sequence is exited.
func (s *BasePlSqlParserListener) ExitAlter_sequence(ctx *Alter_sequenceContext) {}

// EnterAlter_session is called when production alter_session is entered.
func (s *BasePlSqlParserListener) EnterAlter_session(ctx *Alter_sessionContext) {}

// ExitAlter_session is called when production alter_session is exited.
func (s *BasePlSqlParserListener) ExitAlter_session(ctx *Alter_sessionContext) {}

// EnterAlter_session_set_clause is called when production alter_session_set_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_session_set_clause(ctx *Alter_session_set_clauseContext) {
}

// ExitAlter_session_set_clause is called when production alter_session_set_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_session_set_clause(ctx *Alter_session_set_clauseContext) {
}

// EnterCreate_sequence is called when production create_sequence is entered.
func (s *BasePlSqlParserListener) EnterCreate_sequence(ctx *Create_sequenceContext) {}

// ExitCreate_sequence is called when production create_sequence is exited.
func (s *BasePlSqlParserListener) ExitCreate_sequence(ctx *Create_sequenceContext) {}

// EnterSequence_spec is called when production sequence_spec is entered.
func (s *BasePlSqlParserListener) EnterSequence_spec(ctx *Sequence_specContext) {}

// ExitSequence_spec is called when production sequence_spec is exited.
func (s *BasePlSqlParserListener) ExitSequence_spec(ctx *Sequence_specContext) {}

// EnterSequence_start_clause is called when production sequence_start_clause is entered.
func (s *BasePlSqlParserListener) EnterSequence_start_clause(ctx *Sequence_start_clauseContext) {}

// ExitSequence_start_clause is called when production sequence_start_clause is exited.
func (s *BasePlSqlParserListener) ExitSequence_start_clause(ctx *Sequence_start_clauseContext) {}

// EnterCreate_analytic_view is called when production create_analytic_view is entered.
func (s *BasePlSqlParserListener) EnterCreate_analytic_view(ctx *Create_analytic_viewContext) {}

// ExitCreate_analytic_view is called when production create_analytic_view is exited.
func (s *BasePlSqlParserListener) ExitCreate_analytic_view(ctx *Create_analytic_viewContext) {}

// EnterClassification_clause is called when production classification_clause is entered.
func (s *BasePlSqlParserListener) EnterClassification_clause(ctx *Classification_clauseContext) {}

// ExitClassification_clause is called when production classification_clause is exited.
func (s *BasePlSqlParserListener) ExitClassification_clause(ctx *Classification_clauseContext) {}

// EnterCaption_clause is called when production caption_clause is entered.
func (s *BasePlSqlParserListener) EnterCaption_clause(ctx *Caption_clauseContext) {}

// ExitCaption_clause is called when production caption_clause is exited.
func (s *BasePlSqlParserListener) ExitCaption_clause(ctx *Caption_clauseContext) {}

// EnterDescription_clause is called when production description_clause is entered.
func (s *BasePlSqlParserListener) EnterDescription_clause(ctx *Description_clauseContext) {}

// ExitDescription_clause is called when production description_clause is exited.
func (s *BasePlSqlParserListener) ExitDescription_clause(ctx *Description_clauseContext) {}

// EnterClassification_item is called when production classification_item is entered.
func (s *BasePlSqlParserListener) EnterClassification_item(ctx *Classification_itemContext) {}

// ExitClassification_item is called when production classification_item is exited.
func (s *BasePlSqlParserListener) ExitClassification_item(ctx *Classification_itemContext) {}

// EnterLanguage is called when production language is entered.
func (s *BasePlSqlParserListener) EnterLanguage(ctx *LanguageContext) {}

// ExitLanguage is called when production language is exited.
func (s *BasePlSqlParserListener) ExitLanguage(ctx *LanguageContext) {}

// EnterCav_using_clause is called when production cav_using_clause is entered.
func (s *BasePlSqlParserListener) EnterCav_using_clause(ctx *Cav_using_clauseContext) {}

// ExitCav_using_clause is called when production cav_using_clause is exited.
func (s *BasePlSqlParserListener) ExitCav_using_clause(ctx *Cav_using_clauseContext) {}

// EnterDim_by_clause is called when production dim_by_clause is entered.
func (s *BasePlSqlParserListener) EnterDim_by_clause(ctx *Dim_by_clauseContext) {}

// ExitDim_by_clause is called when production dim_by_clause is exited.
func (s *BasePlSqlParserListener) ExitDim_by_clause(ctx *Dim_by_clauseContext) {}

// EnterDim_key is called when production dim_key is entered.
func (s *BasePlSqlParserListener) EnterDim_key(ctx *Dim_keyContext) {}

// ExitDim_key is called when production dim_key is exited.
func (s *BasePlSqlParserListener) ExitDim_key(ctx *Dim_keyContext) {}

// EnterDim_ref is called when production dim_ref is entered.
func (s *BasePlSqlParserListener) EnterDim_ref(ctx *Dim_refContext) {}

// ExitDim_ref is called when production dim_ref is exited.
func (s *BasePlSqlParserListener) ExitDim_ref(ctx *Dim_refContext) {}

// EnterHier_ref is called when production hier_ref is entered.
func (s *BasePlSqlParserListener) EnterHier_ref(ctx *Hier_refContext) {}

// ExitHier_ref is called when production hier_ref is exited.
func (s *BasePlSqlParserListener) ExitHier_ref(ctx *Hier_refContext) {}

// EnterMeasures_clause is called when production measures_clause is entered.
func (s *BasePlSqlParserListener) EnterMeasures_clause(ctx *Measures_clauseContext) {}

// ExitMeasures_clause is called when production measures_clause is exited.
func (s *BasePlSqlParserListener) ExitMeasures_clause(ctx *Measures_clauseContext) {}

// EnterAv_measure is called when production av_measure is entered.
func (s *BasePlSqlParserListener) EnterAv_measure(ctx *Av_measureContext) {}

// ExitAv_measure is called when production av_measure is exited.
func (s *BasePlSqlParserListener) ExitAv_measure(ctx *Av_measureContext) {}

// EnterBase_meas_clause is called when production base_meas_clause is entered.
func (s *BasePlSqlParserListener) EnterBase_meas_clause(ctx *Base_meas_clauseContext) {}

// ExitBase_meas_clause is called when production base_meas_clause is exited.
func (s *BasePlSqlParserListener) ExitBase_meas_clause(ctx *Base_meas_clauseContext) {}

// EnterMeas_aggregate_clause is called when production meas_aggregate_clause is entered.
func (s *BasePlSqlParserListener) EnterMeas_aggregate_clause(ctx *Meas_aggregate_clauseContext) {}

// ExitMeas_aggregate_clause is called when production meas_aggregate_clause is exited.
func (s *BasePlSqlParserListener) ExitMeas_aggregate_clause(ctx *Meas_aggregate_clauseContext) {}

// EnterCalc_meas_clause is called when production calc_meas_clause is entered.
func (s *BasePlSqlParserListener) EnterCalc_meas_clause(ctx *Calc_meas_clauseContext) {}

// ExitCalc_meas_clause is called when production calc_meas_clause is exited.
func (s *BasePlSqlParserListener) ExitCalc_meas_clause(ctx *Calc_meas_clauseContext) {}

// EnterDefault_measure_clause is called when production default_measure_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_measure_clause(ctx *Default_measure_clauseContext) {}

// ExitDefault_measure_clause is called when production default_measure_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_measure_clause(ctx *Default_measure_clauseContext) {}

// EnterDefault_aggregate_clause is called when production default_aggregate_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_aggregate_clause(ctx *Default_aggregate_clauseContext) {
}

// ExitDefault_aggregate_clause is called when production default_aggregate_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_aggregate_clause(ctx *Default_aggregate_clauseContext) {
}

// EnterCache_clause is called when production cache_clause is entered.
func (s *BasePlSqlParserListener) EnterCache_clause(ctx *Cache_clauseContext) {}

// ExitCache_clause is called when production cache_clause is exited.
func (s *BasePlSqlParserListener) ExitCache_clause(ctx *Cache_clauseContext) {}

// EnterCache_specification is called when production cache_specification is entered.
func (s *BasePlSqlParserListener) EnterCache_specification(ctx *Cache_specificationContext) {}

// ExitCache_specification is called when production cache_specification is exited.
func (s *BasePlSqlParserListener) ExitCache_specification(ctx *Cache_specificationContext) {}

// EnterLevels_clause is called when production levels_clause is entered.
func (s *BasePlSqlParserListener) EnterLevels_clause(ctx *Levels_clauseContext) {}

// ExitLevels_clause is called when production levels_clause is exited.
func (s *BasePlSqlParserListener) ExitLevels_clause(ctx *Levels_clauseContext) {}

// EnterLevel_specification is called when production level_specification is entered.
func (s *BasePlSqlParserListener) EnterLevel_specification(ctx *Level_specificationContext) {}

// ExitLevel_specification is called when production level_specification is exited.
func (s *BasePlSqlParserListener) ExitLevel_specification(ctx *Level_specificationContext) {}

// EnterLevel_group_type is called when production level_group_type is entered.
func (s *BasePlSqlParserListener) EnterLevel_group_type(ctx *Level_group_typeContext) {}

// ExitLevel_group_type is called when production level_group_type is exited.
func (s *BasePlSqlParserListener) ExitLevel_group_type(ctx *Level_group_typeContext) {}

// EnterFact_columns_clause is called when production fact_columns_clause is entered.
func (s *BasePlSqlParserListener) EnterFact_columns_clause(ctx *Fact_columns_clauseContext) {}

// ExitFact_columns_clause is called when production fact_columns_clause is exited.
func (s *BasePlSqlParserListener) ExitFact_columns_clause(ctx *Fact_columns_clauseContext) {}

// EnterQry_transform_clause is called when production qry_transform_clause is entered.
func (s *BasePlSqlParserListener) EnterQry_transform_clause(ctx *Qry_transform_clauseContext) {}

// ExitQry_transform_clause is called when production qry_transform_clause is exited.
func (s *BasePlSqlParserListener) ExitQry_transform_clause(ctx *Qry_transform_clauseContext) {}

// EnterCreate_attribute_dimension is called when production create_attribute_dimension is entered.
func (s *BasePlSqlParserListener) EnterCreate_attribute_dimension(ctx *Create_attribute_dimensionContext) {
}

// ExitCreate_attribute_dimension is called when production create_attribute_dimension is exited.
func (s *BasePlSqlParserListener) ExitCreate_attribute_dimension(ctx *Create_attribute_dimensionContext) {
}

// EnterAd_using_clause is called when production ad_using_clause is entered.
func (s *BasePlSqlParserListener) EnterAd_using_clause(ctx *Ad_using_clauseContext) {}

// ExitAd_using_clause is called when production ad_using_clause is exited.
func (s *BasePlSqlParserListener) ExitAd_using_clause(ctx *Ad_using_clauseContext) {}

// EnterSource_clause is called when production source_clause is entered.
func (s *BasePlSqlParserListener) EnterSource_clause(ctx *Source_clauseContext) {}

// ExitSource_clause is called when production source_clause is exited.
func (s *BasePlSqlParserListener) ExitSource_clause(ctx *Source_clauseContext) {}

// EnterJoin_path_clause is called when production join_path_clause is entered.
func (s *BasePlSqlParserListener) EnterJoin_path_clause(ctx *Join_path_clauseContext) {}

// ExitJoin_path_clause is called when production join_path_clause is exited.
func (s *BasePlSqlParserListener) ExitJoin_path_clause(ctx *Join_path_clauseContext) {}

// EnterJoin_condition is called when production join_condition is entered.
func (s *BasePlSqlParserListener) EnterJoin_condition(ctx *Join_conditionContext) {}

// ExitJoin_condition is called when production join_condition is exited.
func (s *BasePlSqlParserListener) ExitJoin_condition(ctx *Join_conditionContext) {}

// EnterJoin_condition_item is called when production join_condition_item is entered.
func (s *BasePlSqlParserListener) EnterJoin_condition_item(ctx *Join_condition_itemContext) {}

// ExitJoin_condition_item is called when production join_condition_item is exited.
func (s *BasePlSqlParserListener) ExitJoin_condition_item(ctx *Join_condition_itemContext) {}

// EnterAttributes_clause is called when production attributes_clause is entered.
func (s *BasePlSqlParserListener) EnterAttributes_clause(ctx *Attributes_clauseContext) {}

// ExitAttributes_clause is called when production attributes_clause is exited.
func (s *BasePlSqlParserListener) ExitAttributes_clause(ctx *Attributes_clauseContext) {}

// EnterAd_attributes_clause is called when production ad_attributes_clause is entered.
func (s *BasePlSqlParserListener) EnterAd_attributes_clause(ctx *Ad_attributes_clauseContext) {}

// ExitAd_attributes_clause is called when production ad_attributes_clause is exited.
func (s *BasePlSqlParserListener) ExitAd_attributes_clause(ctx *Ad_attributes_clauseContext) {}

// EnterAd_level_clause is called when production ad_level_clause is entered.
func (s *BasePlSqlParserListener) EnterAd_level_clause(ctx *Ad_level_clauseContext) {}

// ExitAd_level_clause is called when production ad_level_clause is exited.
func (s *BasePlSqlParserListener) ExitAd_level_clause(ctx *Ad_level_clauseContext) {}

// EnterKey_clause is called when production key_clause is entered.
func (s *BasePlSqlParserListener) EnterKey_clause(ctx *Key_clauseContext) {}

// ExitKey_clause is called when production key_clause is exited.
func (s *BasePlSqlParserListener) ExitKey_clause(ctx *Key_clauseContext) {}

// EnterAlternate_key_clause is called when production alternate_key_clause is entered.
func (s *BasePlSqlParserListener) EnterAlternate_key_clause(ctx *Alternate_key_clauseContext) {}

// ExitAlternate_key_clause is called when production alternate_key_clause is exited.
func (s *BasePlSqlParserListener) ExitAlternate_key_clause(ctx *Alternate_key_clauseContext) {}

// EnterDim_order_clause is called when production dim_order_clause is entered.
func (s *BasePlSqlParserListener) EnterDim_order_clause(ctx *Dim_order_clauseContext) {}

// ExitDim_order_clause is called when production dim_order_clause is exited.
func (s *BasePlSqlParserListener) ExitDim_order_clause(ctx *Dim_order_clauseContext) {}

// EnterAll_clause is called when production all_clause is entered.
func (s *BasePlSqlParserListener) EnterAll_clause(ctx *All_clauseContext) {}

// ExitAll_clause is called when production all_clause is exited.
func (s *BasePlSqlParserListener) ExitAll_clause(ctx *All_clauseContext) {}

// EnterCreate_audit_policy is called when production create_audit_policy is entered.
func (s *BasePlSqlParserListener) EnterCreate_audit_policy(ctx *Create_audit_policyContext) {}

// ExitCreate_audit_policy is called when production create_audit_policy is exited.
func (s *BasePlSqlParserListener) ExitCreate_audit_policy(ctx *Create_audit_policyContext) {}

// EnterPrivilege_audit_clause is called when production privilege_audit_clause is entered.
func (s *BasePlSqlParserListener) EnterPrivilege_audit_clause(ctx *Privilege_audit_clauseContext) {}

// ExitPrivilege_audit_clause is called when production privilege_audit_clause is exited.
func (s *BasePlSqlParserListener) ExitPrivilege_audit_clause(ctx *Privilege_audit_clauseContext) {}

// EnterAction_audit_clause is called when production action_audit_clause is entered.
func (s *BasePlSqlParserListener) EnterAction_audit_clause(ctx *Action_audit_clauseContext) {}

// ExitAction_audit_clause is called when production action_audit_clause is exited.
func (s *BasePlSqlParserListener) ExitAction_audit_clause(ctx *Action_audit_clauseContext) {}

// EnterSystem_actions is called when production system_actions is entered.
func (s *BasePlSqlParserListener) EnterSystem_actions(ctx *System_actionsContext) {}

// ExitSystem_actions is called when production system_actions is exited.
func (s *BasePlSqlParserListener) ExitSystem_actions(ctx *System_actionsContext) {}

// EnterStandard_actions is called when production standard_actions is entered.
func (s *BasePlSqlParserListener) EnterStandard_actions(ctx *Standard_actionsContext) {}

// ExitStandard_actions is called when production standard_actions is exited.
func (s *BasePlSqlParserListener) ExitStandard_actions(ctx *Standard_actionsContext) {}

// EnterActions_clause is called when production actions_clause is entered.
func (s *BasePlSqlParserListener) EnterActions_clause(ctx *Actions_clauseContext) {}

// ExitActions_clause is called when production actions_clause is exited.
func (s *BasePlSqlParserListener) ExitActions_clause(ctx *Actions_clauseContext) {}

// EnterObject_action is called when production object_action is entered.
func (s *BasePlSqlParserListener) EnterObject_action(ctx *Object_actionContext) {}

// ExitObject_action is called when production object_action is exited.
func (s *BasePlSqlParserListener) ExitObject_action(ctx *Object_actionContext) {}

// EnterSystem_action is called when production system_action is entered.
func (s *BasePlSqlParserListener) EnterSystem_action(ctx *System_actionContext) {}

// ExitSystem_action is called when production system_action is exited.
func (s *BasePlSqlParserListener) ExitSystem_action(ctx *System_actionContext) {}

// EnterComponent_actions is called when production component_actions is entered.
func (s *BasePlSqlParserListener) EnterComponent_actions(ctx *Component_actionsContext) {}

// ExitComponent_actions is called when production component_actions is exited.
func (s *BasePlSqlParserListener) ExitComponent_actions(ctx *Component_actionsContext) {}

// EnterComponent_action is called when production component_action is entered.
func (s *BasePlSqlParserListener) EnterComponent_action(ctx *Component_actionContext) {}

// ExitComponent_action is called when production component_action is exited.
func (s *BasePlSqlParserListener) ExitComponent_action(ctx *Component_actionContext) {}

// EnterRole_audit_clause is called when production role_audit_clause is entered.
func (s *BasePlSqlParserListener) EnterRole_audit_clause(ctx *Role_audit_clauseContext) {}

// ExitRole_audit_clause is called when production role_audit_clause is exited.
func (s *BasePlSqlParserListener) ExitRole_audit_clause(ctx *Role_audit_clauseContext) {}

// EnterCreate_controlfile is called when production create_controlfile is entered.
func (s *BasePlSqlParserListener) EnterCreate_controlfile(ctx *Create_controlfileContext) {}

// ExitCreate_controlfile is called when production create_controlfile is exited.
func (s *BasePlSqlParserListener) ExitCreate_controlfile(ctx *Create_controlfileContext) {}

// EnterControlfile_options is called when production controlfile_options is entered.
func (s *BasePlSqlParserListener) EnterControlfile_options(ctx *Controlfile_optionsContext) {}

// ExitControlfile_options is called when production controlfile_options is exited.
func (s *BasePlSqlParserListener) ExitControlfile_options(ctx *Controlfile_optionsContext) {}

// EnterLogfile_clause is called when production logfile_clause is entered.
func (s *BasePlSqlParserListener) EnterLogfile_clause(ctx *Logfile_clauseContext) {}

// ExitLogfile_clause is called when production logfile_clause is exited.
func (s *BasePlSqlParserListener) ExitLogfile_clause(ctx *Logfile_clauseContext) {}

// EnterCharacter_set_clause is called when production character_set_clause is entered.
func (s *BasePlSqlParserListener) EnterCharacter_set_clause(ctx *Character_set_clauseContext) {}

// ExitCharacter_set_clause is called when production character_set_clause is exited.
func (s *BasePlSqlParserListener) ExitCharacter_set_clause(ctx *Character_set_clauseContext) {}

// EnterFile_specification is called when production file_specification is entered.
func (s *BasePlSqlParserListener) EnterFile_specification(ctx *File_specificationContext) {}

// ExitFile_specification is called when production file_specification is exited.
func (s *BasePlSqlParserListener) ExitFile_specification(ctx *File_specificationContext) {}

// EnterCreate_diskgroup is called when production create_diskgroup is entered.
func (s *BasePlSqlParserListener) EnterCreate_diskgroup(ctx *Create_diskgroupContext) {}

// ExitCreate_diskgroup is called when production create_diskgroup is exited.
func (s *BasePlSqlParserListener) ExitCreate_diskgroup(ctx *Create_diskgroupContext) {}

// EnterQualified_disk_clause is called when production qualified_disk_clause is entered.
func (s *BasePlSqlParserListener) EnterQualified_disk_clause(ctx *Qualified_disk_clauseContext) {}

// ExitQualified_disk_clause is called when production qualified_disk_clause is exited.
func (s *BasePlSqlParserListener) ExitQualified_disk_clause(ctx *Qualified_disk_clauseContext) {}

// EnterCreate_edition is called when production create_edition is entered.
func (s *BasePlSqlParserListener) EnterCreate_edition(ctx *Create_editionContext) {}

// ExitCreate_edition is called when production create_edition is exited.
func (s *BasePlSqlParserListener) ExitCreate_edition(ctx *Create_editionContext) {}

// EnterCreate_flashback_archive is called when production create_flashback_archive is entered.
func (s *BasePlSqlParserListener) EnterCreate_flashback_archive(ctx *Create_flashback_archiveContext) {
}

// ExitCreate_flashback_archive is called when production create_flashback_archive is exited.
func (s *BasePlSqlParserListener) ExitCreate_flashback_archive(ctx *Create_flashback_archiveContext) {
}

// EnterFlashback_archive_quota is called when production flashback_archive_quota is entered.
func (s *BasePlSqlParserListener) EnterFlashback_archive_quota(ctx *Flashback_archive_quotaContext) {}

// ExitFlashback_archive_quota is called when production flashback_archive_quota is exited.
func (s *BasePlSqlParserListener) ExitFlashback_archive_quota(ctx *Flashback_archive_quotaContext) {}

// EnterFlashback_archive_retention is called when production flashback_archive_retention is entered.
func (s *BasePlSqlParserListener) EnterFlashback_archive_retention(ctx *Flashback_archive_retentionContext) {
}

// ExitFlashback_archive_retention is called when production flashback_archive_retention is exited.
func (s *BasePlSqlParserListener) ExitFlashback_archive_retention(ctx *Flashback_archive_retentionContext) {
}

// EnterCreate_hierarchy is called when production create_hierarchy is entered.
func (s *BasePlSqlParserListener) EnterCreate_hierarchy(ctx *Create_hierarchyContext) {}

// ExitCreate_hierarchy is called when production create_hierarchy is exited.
func (s *BasePlSqlParserListener) ExitCreate_hierarchy(ctx *Create_hierarchyContext) {}

// EnterHier_using_clause is called when production hier_using_clause is entered.
func (s *BasePlSqlParserListener) EnterHier_using_clause(ctx *Hier_using_clauseContext) {}

// ExitHier_using_clause is called when production hier_using_clause is exited.
func (s *BasePlSqlParserListener) ExitHier_using_clause(ctx *Hier_using_clauseContext) {}

// EnterLevel_hier_clause is called when production level_hier_clause is entered.
func (s *BasePlSqlParserListener) EnterLevel_hier_clause(ctx *Level_hier_clauseContext) {}

// ExitLevel_hier_clause is called when production level_hier_clause is exited.
func (s *BasePlSqlParserListener) ExitLevel_hier_clause(ctx *Level_hier_clauseContext) {}

// EnterHier_attrs_clause is called when production hier_attrs_clause is entered.
func (s *BasePlSqlParserListener) EnterHier_attrs_clause(ctx *Hier_attrs_clauseContext) {}

// ExitHier_attrs_clause is called when production hier_attrs_clause is exited.
func (s *BasePlSqlParserListener) ExitHier_attrs_clause(ctx *Hier_attrs_clauseContext) {}

// EnterHier_attr_clause is called when production hier_attr_clause is entered.
func (s *BasePlSqlParserListener) EnterHier_attr_clause(ctx *Hier_attr_clauseContext) {}

// ExitHier_attr_clause is called when production hier_attr_clause is exited.
func (s *BasePlSqlParserListener) ExitHier_attr_clause(ctx *Hier_attr_clauseContext) {}

// EnterHier_attr_name is called when production hier_attr_name is entered.
func (s *BasePlSqlParserListener) EnterHier_attr_name(ctx *Hier_attr_nameContext) {}

// ExitHier_attr_name is called when production hier_attr_name is exited.
func (s *BasePlSqlParserListener) ExitHier_attr_name(ctx *Hier_attr_nameContext) {}

// EnterCreate_index is called when production create_index is entered.
func (s *BasePlSqlParserListener) EnterCreate_index(ctx *Create_indexContext) {}

// ExitCreate_index is called when production create_index is exited.
func (s *BasePlSqlParserListener) ExitCreate_index(ctx *Create_indexContext) {}

// EnterCluster_index_clause is called when production cluster_index_clause is entered.
func (s *BasePlSqlParserListener) EnterCluster_index_clause(ctx *Cluster_index_clauseContext) {}

// ExitCluster_index_clause is called when production cluster_index_clause is exited.
func (s *BasePlSqlParserListener) ExitCluster_index_clause(ctx *Cluster_index_clauseContext) {}

// EnterCluster_name is called when production cluster_name is entered.
func (s *BasePlSqlParserListener) EnterCluster_name(ctx *Cluster_nameContext) {}

// ExitCluster_name is called when production cluster_name is exited.
func (s *BasePlSqlParserListener) ExitCluster_name(ctx *Cluster_nameContext) {}

// EnterTable_index_clause is called when production table_index_clause is entered.
func (s *BasePlSqlParserListener) EnterTable_index_clause(ctx *Table_index_clauseContext) {}

// ExitTable_index_clause is called when production table_index_clause is exited.
func (s *BasePlSqlParserListener) ExitTable_index_clause(ctx *Table_index_clauseContext) {}

// EnterBitmap_join_index_clause is called when production bitmap_join_index_clause is entered.
func (s *BasePlSqlParserListener) EnterBitmap_join_index_clause(ctx *Bitmap_join_index_clauseContext) {
}

// ExitBitmap_join_index_clause is called when production bitmap_join_index_clause is exited.
func (s *BasePlSqlParserListener) ExitBitmap_join_index_clause(ctx *Bitmap_join_index_clauseContext) {
}

// EnterIndex_expr is called when production index_expr is entered.
func (s *BasePlSqlParserListener) EnterIndex_expr(ctx *Index_exprContext) {}

// ExitIndex_expr is called when production index_expr is exited.
func (s *BasePlSqlParserListener) ExitIndex_expr(ctx *Index_exprContext) {}

// EnterIndex_properties is called when production index_properties is entered.
func (s *BasePlSqlParserListener) EnterIndex_properties(ctx *Index_propertiesContext) {}

// ExitIndex_properties is called when production index_properties is exited.
func (s *BasePlSqlParserListener) ExitIndex_properties(ctx *Index_propertiesContext) {}

// EnterDomain_index_clause is called when production domain_index_clause is entered.
func (s *BasePlSqlParserListener) EnterDomain_index_clause(ctx *Domain_index_clauseContext) {}

// ExitDomain_index_clause is called when production domain_index_clause is exited.
func (s *BasePlSqlParserListener) ExitDomain_index_clause(ctx *Domain_index_clauseContext) {}

// EnterLocal_domain_index_clause is called when production local_domain_index_clause is entered.
func (s *BasePlSqlParserListener) EnterLocal_domain_index_clause(ctx *Local_domain_index_clauseContext) {
}

// ExitLocal_domain_index_clause is called when production local_domain_index_clause is exited.
func (s *BasePlSqlParserListener) ExitLocal_domain_index_clause(ctx *Local_domain_index_clauseContext) {
}

// EnterXmlindex_clause is called when production xmlindex_clause is entered.
func (s *BasePlSqlParserListener) EnterXmlindex_clause(ctx *Xmlindex_clauseContext) {}

// ExitXmlindex_clause is called when production xmlindex_clause is exited.
func (s *BasePlSqlParserListener) ExitXmlindex_clause(ctx *Xmlindex_clauseContext) {}

// EnterLocal_xmlindex_clause is called when production local_xmlindex_clause is entered.
func (s *BasePlSqlParserListener) EnterLocal_xmlindex_clause(ctx *Local_xmlindex_clauseContext) {}

// ExitLocal_xmlindex_clause is called when production local_xmlindex_clause is exited.
func (s *BasePlSqlParserListener) ExitLocal_xmlindex_clause(ctx *Local_xmlindex_clauseContext) {}

// EnterGlobal_partitioned_index is called when production global_partitioned_index is entered.
func (s *BasePlSqlParserListener) EnterGlobal_partitioned_index(ctx *Global_partitioned_indexContext) {
}

// ExitGlobal_partitioned_index is called when production global_partitioned_index is exited.
func (s *BasePlSqlParserListener) ExitGlobal_partitioned_index(ctx *Global_partitioned_indexContext) {
}

// EnterIndex_partitioning_clause is called when production index_partitioning_clause is entered.
func (s *BasePlSqlParserListener) EnterIndex_partitioning_clause(ctx *Index_partitioning_clauseContext) {
}

// ExitIndex_partitioning_clause is called when production index_partitioning_clause is exited.
func (s *BasePlSqlParserListener) ExitIndex_partitioning_clause(ctx *Index_partitioning_clauseContext) {
}

// EnterIndex_partitioning_values_list is called when production index_partitioning_values_list is entered.
func (s *BasePlSqlParserListener) EnterIndex_partitioning_values_list(ctx *Index_partitioning_values_listContext) {
}

// ExitIndex_partitioning_values_list is called when production index_partitioning_values_list is exited.
func (s *BasePlSqlParserListener) ExitIndex_partitioning_values_list(ctx *Index_partitioning_values_listContext) {
}

// EnterLocal_partitioned_index is called when production local_partitioned_index is entered.
func (s *BasePlSqlParserListener) EnterLocal_partitioned_index(ctx *Local_partitioned_indexContext) {}

// ExitLocal_partitioned_index is called when production local_partitioned_index is exited.
func (s *BasePlSqlParserListener) ExitLocal_partitioned_index(ctx *Local_partitioned_indexContext) {}

// EnterOn_range_partitioned_table is called when production on_range_partitioned_table is entered.
func (s *BasePlSqlParserListener) EnterOn_range_partitioned_table(ctx *On_range_partitioned_tableContext) {
}

// ExitOn_range_partitioned_table is called when production on_range_partitioned_table is exited.
func (s *BasePlSqlParserListener) ExitOn_range_partitioned_table(ctx *On_range_partitioned_tableContext) {
}

// EnterOn_list_partitioned_table is called when production on_list_partitioned_table is entered.
func (s *BasePlSqlParserListener) EnterOn_list_partitioned_table(ctx *On_list_partitioned_tableContext) {
}

// ExitOn_list_partitioned_table is called when production on_list_partitioned_table is exited.
func (s *BasePlSqlParserListener) ExitOn_list_partitioned_table(ctx *On_list_partitioned_tableContext) {
}

// EnterPartitioned_table is called when production partitioned_table is entered.
func (s *BasePlSqlParserListener) EnterPartitioned_table(ctx *Partitioned_tableContext) {}

// ExitPartitioned_table is called when production partitioned_table is exited.
func (s *BasePlSqlParserListener) ExitPartitioned_table(ctx *Partitioned_tableContext) {}

// EnterOn_hash_partitioned_table is called when production on_hash_partitioned_table is entered.
func (s *BasePlSqlParserListener) EnterOn_hash_partitioned_table(ctx *On_hash_partitioned_tableContext) {
}

// ExitOn_hash_partitioned_table is called when production on_hash_partitioned_table is exited.
func (s *BasePlSqlParserListener) ExitOn_hash_partitioned_table(ctx *On_hash_partitioned_tableContext) {
}

// EnterOn_hash_partitioned_clause is called when production on_hash_partitioned_clause is entered.
func (s *BasePlSqlParserListener) EnterOn_hash_partitioned_clause(ctx *On_hash_partitioned_clauseContext) {
}

// ExitOn_hash_partitioned_clause is called when production on_hash_partitioned_clause is exited.
func (s *BasePlSqlParserListener) ExitOn_hash_partitioned_clause(ctx *On_hash_partitioned_clauseContext) {
}

// EnterOn_comp_partitioned_table is called when production on_comp_partitioned_table is entered.
func (s *BasePlSqlParserListener) EnterOn_comp_partitioned_table(ctx *On_comp_partitioned_tableContext) {
}

// ExitOn_comp_partitioned_table is called when production on_comp_partitioned_table is exited.
func (s *BasePlSqlParserListener) ExitOn_comp_partitioned_table(ctx *On_comp_partitioned_tableContext) {
}

// EnterOn_comp_partitioned_clause is called when production on_comp_partitioned_clause is entered.
func (s *BasePlSqlParserListener) EnterOn_comp_partitioned_clause(ctx *On_comp_partitioned_clauseContext) {
}

// ExitOn_comp_partitioned_clause is called when production on_comp_partitioned_clause is exited.
func (s *BasePlSqlParserListener) ExitOn_comp_partitioned_clause(ctx *On_comp_partitioned_clauseContext) {
}

// EnterIndex_subpartition_clause is called when production index_subpartition_clause is entered.
func (s *BasePlSqlParserListener) EnterIndex_subpartition_clause(ctx *Index_subpartition_clauseContext) {
}

// ExitIndex_subpartition_clause is called when production index_subpartition_clause is exited.
func (s *BasePlSqlParserListener) ExitIndex_subpartition_clause(ctx *Index_subpartition_clauseContext) {
}

// EnterIndex_subpartition_subclause is called when production index_subpartition_subclause is entered.
func (s *BasePlSqlParserListener) EnterIndex_subpartition_subclause(ctx *Index_subpartition_subclauseContext) {
}

// ExitIndex_subpartition_subclause is called when production index_subpartition_subclause is exited.
func (s *BasePlSqlParserListener) ExitIndex_subpartition_subclause(ctx *Index_subpartition_subclauseContext) {
}

// EnterOdci_parameters is called when production odci_parameters is entered.
func (s *BasePlSqlParserListener) EnterOdci_parameters(ctx *Odci_parametersContext) {}

// ExitOdci_parameters is called when production odci_parameters is exited.
func (s *BasePlSqlParserListener) ExitOdci_parameters(ctx *Odci_parametersContext) {}

// EnterIndextype is called when production indextype is entered.
func (s *BasePlSqlParserListener) EnterIndextype(ctx *IndextypeContext) {}

// ExitIndextype is called when production indextype is exited.
func (s *BasePlSqlParserListener) ExitIndextype(ctx *IndextypeContext) {}

// EnterAlter_index is called when production alter_index is entered.
func (s *BasePlSqlParserListener) EnterAlter_index(ctx *Alter_indexContext) {}

// ExitAlter_index is called when production alter_index is exited.
func (s *BasePlSqlParserListener) ExitAlter_index(ctx *Alter_indexContext) {}

// EnterAlter_index_ops_set1 is called when production alter_index_ops_set1 is entered.
func (s *BasePlSqlParserListener) EnterAlter_index_ops_set1(ctx *Alter_index_ops_set1Context) {}

// ExitAlter_index_ops_set1 is called when production alter_index_ops_set1 is exited.
func (s *BasePlSqlParserListener) ExitAlter_index_ops_set1(ctx *Alter_index_ops_set1Context) {}

// EnterAlter_index_ops_set2 is called when production alter_index_ops_set2 is entered.
func (s *BasePlSqlParserListener) EnterAlter_index_ops_set2(ctx *Alter_index_ops_set2Context) {}

// ExitAlter_index_ops_set2 is called when production alter_index_ops_set2 is exited.
func (s *BasePlSqlParserListener) ExitAlter_index_ops_set2(ctx *Alter_index_ops_set2Context) {}

// EnterVisible_or_invisible is called when production visible_or_invisible is entered.
func (s *BasePlSqlParserListener) EnterVisible_or_invisible(ctx *Visible_or_invisibleContext) {}

// ExitVisible_or_invisible is called when production visible_or_invisible is exited.
func (s *BasePlSqlParserListener) ExitVisible_or_invisible(ctx *Visible_or_invisibleContext) {}

// EnterMonitoring_nomonitoring is called when production monitoring_nomonitoring is entered.
func (s *BasePlSqlParserListener) EnterMonitoring_nomonitoring(ctx *Monitoring_nomonitoringContext) {}

// ExitMonitoring_nomonitoring is called when production monitoring_nomonitoring is exited.
func (s *BasePlSqlParserListener) ExitMonitoring_nomonitoring(ctx *Monitoring_nomonitoringContext) {}

// EnterRebuild_clause is called when production rebuild_clause is entered.
func (s *BasePlSqlParserListener) EnterRebuild_clause(ctx *Rebuild_clauseContext) {}

// ExitRebuild_clause is called when production rebuild_clause is exited.
func (s *BasePlSqlParserListener) ExitRebuild_clause(ctx *Rebuild_clauseContext) {}

// EnterAlter_index_partitioning is called when production alter_index_partitioning is entered.
func (s *BasePlSqlParserListener) EnterAlter_index_partitioning(ctx *Alter_index_partitioningContext) {
}

// ExitAlter_index_partitioning is called when production alter_index_partitioning is exited.
func (s *BasePlSqlParserListener) ExitAlter_index_partitioning(ctx *Alter_index_partitioningContext) {
}

// EnterModify_index_default_attrs is called when production modify_index_default_attrs is entered.
func (s *BasePlSqlParserListener) EnterModify_index_default_attrs(ctx *Modify_index_default_attrsContext) {
}

// ExitModify_index_default_attrs is called when production modify_index_default_attrs is exited.
func (s *BasePlSqlParserListener) ExitModify_index_default_attrs(ctx *Modify_index_default_attrsContext) {
}

// EnterAdd_hash_index_partition is called when production add_hash_index_partition is entered.
func (s *BasePlSqlParserListener) EnterAdd_hash_index_partition(ctx *Add_hash_index_partitionContext) {
}

// ExitAdd_hash_index_partition is called when production add_hash_index_partition is exited.
func (s *BasePlSqlParserListener) ExitAdd_hash_index_partition(ctx *Add_hash_index_partitionContext) {
}

// EnterCoalesce_index_partition is called when production coalesce_index_partition is entered.
func (s *BasePlSqlParserListener) EnterCoalesce_index_partition(ctx *Coalesce_index_partitionContext) {
}

// ExitCoalesce_index_partition is called when production coalesce_index_partition is exited.
func (s *BasePlSqlParserListener) ExitCoalesce_index_partition(ctx *Coalesce_index_partitionContext) {
}

// EnterModify_index_partition is called when production modify_index_partition is entered.
func (s *BasePlSqlParserListener) EnterModify_index_partition(ctx *Modify_index_partitionContext) {}

// ExitModify_index_partition is called when production modify_index_partition is exited.
func (s *BasePlSqlParserListener) ExitModify_index_partition(ctx *Modify_index_partitionContext) {}

// EnterModify_index_partitions_ops is called when production modify_index_partitions_ops is entered.
func (s *BasePlSqlParserListener) EnterModify_index_partitions_ops(ctx *Modify_index_partitions_opsContext) {
}

// ExitModify_index_partitions_ops is called when production modify_index_partitions_ops is exited.
func (s *BasePlSqlParserListener) ExitModify_index_partitions_ops(ctx *Modify_index_partitions_opsContext) {
}

// EnterRename_index_partition is called when production rename_index_partition is entered.
func (s *BasePlSqlParserListener) EnterRename_index_partition(ctx *Rename_index_partitionContext) {}

// ExitRename_index_partition is called when production rename_index_partition is exited.
func (s *BasePlSqlParserListener) ExitRename_index_partition(ctx *Rename_index_partitionContext) {}

// EnterDrop_index_partition is called when production drop_index_partition is entered.
func (s *BasePlSqlParserListener) EnterDrop_index_partition(ctx *Drop_index_partitionContext) {}

// ExitDrop_index_partition is called when production drop_index_partition is exited.
func (s *BasePlSqlParserListener) ExitDrop_index_partition(ctx *Drop_index_partitionContext) {}

// EnterSplit_index_partition is called when production split_index_partition is entered.
func (s *BasePlSqlParserListener) EnterSplit_index_partition(ctx *Split_index_partitionContext) {}

// ExitSplit_index_partition is called when production split_index_partition is exited.
func (s *BasePlSqlParserListener) ExitSplit_index_partition(ctx *Split_index_partitionContext) {}

// EnterIndex_partition_description is called when production index_partition_description is entered.
func (s *BasePlSqlParserListener) EnterIndex_partition_description(ctx *Index_partition_descriptionContext) {
}

// ExitIndex_partition_description is called when production index_partition_description is exited.
func (s *BasePlSqlParserListener) ExitIndex_partition_description(ctx *Index_partition_descriptionContext) {
}

// EnterModify_index_subpartition is called when production modify_index_subpartition is entered.
func (s *BasePlSqlParserListener) EnterModify_index_subpartition(ctx *Modify_index_subpartitionContext) {
}

// ExitModify_index_subpartition is called when production modify_index_subpartition is exited.
func (s *BasePlSqlParserListener) ExitModify_index_subpartition(ctx *Modify_index_subpartitionContext) {
}

// EnterPartition_name_old is called when production partition_name_old is entered.
func (s *BasePlSqlParserListener) EnterPartition_name_old(ctx *Partition_name_oldContext) {}

// ExitPartition_name_old is called when production partition_name_old is exited.
func (s *BasePlSqlParserListener) ExitPartition_name_old(ctx *Partition_name_oldContext) {}

// EnterNew_partition_name is called when production new_partition_name is entered.
func (s *BasePlSqlParserListener) EnterNew_partition_name(ctx *New_partition_nameContext) {}

// ExitNew_partition_name is called when production new_partition_name is exited.
func (s *BasePlSqlParserListener) ExitNew_partition_name(ctx *New_partition_nameContext) {}

// EnterNew_index_name is called when production new_index_name is entered.
func (s *BasePlSqlParserListener) EnterNew_index_name(ctx *New_index_nameContext) {}

// ExitNew_index_name is called when production new_index_name is exited.
func (s *BasePlSqlParserListener) ExitNew_index_name(ctx *New_index_nameContext) {}

// EnterAlter_inmemory_join_group is called when production alter_inmemory_join_group is entered.
func (s *BasePlSqlParserListener) EnterAlter_inmemory_join_group(ctx *Alter_inmemory_join_groupContext) {
}

// ExitAlter_inmemory_join_group is called when production alter_inmemory_join_group is exited.
func (s *BasePlSqlParserListener) ExitAlter_inmemory_join_group(ctx *Alter_inmemory_join_groupContext) {
}

// EnterCreate_user is called when production create_user is entered.
func (s *BasePlSqlParserListener) EnterCreate_user(ctx *Create_userContext) {}

// ExitCreate_user is called when production create_user is exited.
func (s *BasePlSqlParserListener) ExitCreate_user(ctx *Create_userContext) {}

// EnterAlter_user is called when production alter_user is entered.
func (s *BasePlSqlParserListener) EnterAlter_user(ctx *Alter_userContext) {}

// ExitAlter_user is called when production alter_user is exited.
func (s *BasePlSqlParserListener) ExitAlter_user(ctx *Alter_userContext) {}

// EnterDrop_user is called when production drop_user is entered.
func (s *BasePlSqlParserListener) EnterDrop_user(ctx *Drop_userContext) {}

// ExitDrop_user is called when production drop_user is exited.
func (s *BasePlSqlParserListener) ExitDrop_user(ctx *Drop_userContext) {}

// EnterAlter_identified_by is called when production alter_identified_by is entered.
func (s *BasePlSqlParserListener) EnterAlter_identified_by(ctx *Alter_identified_byContext) {}

// ExitAlter_identified_by is called when production alter_identified_by is exited.
func (s *BasePlSqlParserListener) ExitAlter_identified_by(ctx *Alter_identified_byContext) {}

// EnterIdentified_by is called when production identified_by is entered.
func (s *BasePlSqlParserListener) EnterIdentified_by(ctx *Identified_byContext) {}

// ExitIdentified_by is called when production identified_by is exited.
func (s *BasePlSqlParserListener) ExitIdentified_by(ctx *Identified_byContext) {}

// EnterIdentified_other_clause is called when production identified_other_clause is entered.
func (s *BasePlSqlParserListener) EnterIdentified_other_clause(ctx *Identified_other_clauseContext) {}

// ExitIdentified_other_clause is called when production identified_other_clause is exited.
func (s *BasePlSqlParserListener) ExitIdentified_other_clause(ctx *Identified_other_clauseContext) {}

// EnterUser_tablespace_clause is called when production user_tablespace_clause is entered.
func (s *BasePlSqlParserListener) EnterUser_tablespace_clause(ctx *User_tablespace_clauseContext) {}

// ExitUser_tablespace_clause is called when production user_tablespace_clause is exited.
func (s *BasePlSqlParserListener) ExitUser_tablespace_clause(ctx *User_tablespace_clauseContext) {}

// EnterQuota_clause is called when production quota_clause is entered.
func (s *BasePlSqlParserListener) EnterQuota_clause(ctx *Quota_clauseContext) {}

// ExitQuota_clause is called when production quota_clause is exited.
func (s *BasePlSqlParserListener) ExitQuota_clause(ctx *Quota_clauseContext) {}

// EnterProfile_clause is called when production profile_clause is entered.
func (s *BasePlSqlParserListener) EnterProfile_clause(ctx *Profile_clauseContext) {}

// ExitProfile_clause is called when production profile_clause is exited.
func (s *BasePlSqlParserListener) ExitProfile_clause(ctx *Profile_clauseContext) {}

// EnterRole_clause is called when production role_clause is entered.
func (s *BasePlSqlParserListener) EnterRole_clause(ctx *Role_clauseContext) {}

// ExitRole_clause is called when production role_clause is exited.
func (s *BasePlSqlParserListener) ExitRole_clause(ctx *Role_clauseContext) {}

// EnterUser_default_role_clause is called when production user_default_role_clause is entered.
func (s *BasePlSqlParserListener) EnterUser_default_role_clause(ctx *User_default_role_clauseContext) {
}

// ExitUser_default_role_clause is called when production user_default_role_clause is exited.
func (s *BasePlSqlParserListener) ExitUser_default_role_clause(ctx *User_default_role_clauseContext) {
}

// EnterPassword_expire_clause is called when production password_expire_clause is entered.
func (s *BasePlSqlParserListener) EnterPassword_expire_clause(ctx *Password_expire_clauseContext) {}

// ExitPassword_expire_clause is called when production password_expire_clause is exited.
func (s *BasePlSqlParserListener) ExitPassword_expire_clause(ctx *Password_expire_clauseContext) {}

// EnterUser_lock_clause is called when production user_lock_clause is entered.
func (s *BasePlSqlParserListener) EnterUser_lock_clause(ctx *User_lock_clauseContext) {}

// ExitUser_lock_clause is called when production user_lock_clause is exited.
func (s *BasePlSqlParserListener) ExitUser_lock_clause(ctx *User_lock_clauseContext) {}

// EnterUser_editions_clause is called when production user_editions_clause is entered.
func (s *BasePlSqlParserListener) EnterUser_editions_clause(ctx *User_editions_clauseContext) {}

// ExitUser_editions_clause is called when production user_editions_clause is exited.
func (s *BasePlSqlParserListener) ExitUser_editions_clause(ctx *User_editions_clauseContext) {}

// EnterAlter_user_editions_clause is called when production alter_user_editions_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_user_editions_clause(ctx *Alter_user_editions_clauseContext) {
}

// ExitAlter_user_editions_clause is called when production alter_user_editions_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_user_editions_clause(ctx *Alter_user_editions_clauseContext) {
}

// EnterProxy_clause is called when production proxy_clause is entered.
func (s *BasePlSqlParserListener) EnterProxy_clause(ctx *Proxy_clauseContext) {}

// ExitProxy_clause is called when production proxy_clause is exited.
func (s *BasePlSqlParserListener) ExitProxy_clause(ctx *Proxy_clauseContext) {}

// EnterContainer_names is called when production container_names is entered.
func (s *BasePlSqlParserListener) EnterContainer_names(ctx *Container_namesContext) {}

// ExitContainer_names is called when production container_names is exited.
func (s *BasePlSqlParserListener) ExitContainer_names(ctx *Container_namesContext) {}

// EnterSet_container_data is called when production set_container_data is entered.
func (s *BasePlSqlParserListener) EnterSet_container_data(ctx *Set_container_dataContext) {}

// ExitSet_container_data is called when production set_container_data is exited.
func (s *BasePlSqlParserListener) ExitSet_container_data(ctx *Set_container_dataContext) {}

// EnterAdd_rem_container_data is called when production add_rem_container_data is entered.
func (s *BasePlSqlParserListener) EnterAdd_rem_container_data(ctx *Add_rem_container_dataContext) {}

// ExitAdd_rem_container_data is called when production add_rem_container_data is exited.
func (s *BasePlSqlParserListener) ExitAdd_rem_container_data(ctx *Add_rem_container_dataContext) {}

// EnterContainer_data_clause is called when production container_data_clause is entered.
func (s *BasePlSqlParserListener) EnterContainer_data_clause(ctx *Container_data_clauseContext) {}

// ExitContainer_data_clause is called when production container_data_clause is exited.
func (s *BasePlSqlParserListener) ExitContainer_data_clause(ctx *Container_data_clauseContext) {}

// EnterAdminister_key_management is called when production administer_key_management is entered.
func (s *BasePlSqlParserListener) EnterAdminister_key_management(ctx *Administer_key_managementContext) {
}

// ExitAdminister_key_management is called when production administer_key_management is exited.
func (s *BasePlSqlParserListener) ExitAdminister_key_management(ctx *Administer_key_managementContext) {
}

// EnterKeystore_management_clauses is called when production keystore_management_clauses is entered.
func (s *BasePlSqlParserListener) EnterKeystore_management_clauses(ctx *Keystore_management_clausesContext) {
}

// ExitKeystore_management_clauses is called when production keystore_management_clauses is exited.
func (s *BasePlSqlParserListener) ExitKeystore_management_clauses(ctx *Keystore_management_clausesContext) {
}

// EnterCreate_keystore is called when production create_keystore is entered.
func (s *BasePlSqlParserListener) EnterCreate_keystore(ctx *Create_keystoreContext) {}

// ExitCreate_keystore is called when production create_keystore is exited.
func (s *BasePlSqlParserListener) ExitCreate_keystore(ctx *Create_keystoreContext) {}

// EnterOpen_keystore is called when production open_keystore is entered.
func (s *BasePlSqlParserListener) EnterOpen_keystore(ctx *Open_keystoreContext) {}

// ExitOpen_keystore is called when production open_keystore is exited.
func (s *BasePlSqlParserListener) ExitOpen_keystore(ctx *Open_keystoreContext) {}

// EnterForce_keystore is called when production force_keystore is entered.
func (s *BasePlSqlParserListener) EnterForce_keystore(ctx *Force_keystoreContext) {}

// ExitForce_keystore is called when production force_keystore is exited.
func (s *BasePlSqlParserListener) ExitForce_keystore(ctx *Force_keystoreContext) {}

// EnterClose_keystore is called when production close_keystore is entered.
func (s *BasePlSqlParserListener) EnterClose_keystore(ctx *Close_keystoreContext) {}

// ExitClose_keystore is called when production close_keystore is exited.
func (s *BasePlSqlParserListener) ExitClose_keystore(ctx *Close_keystoreContext) {}

// EnterBackup_keystore is called when production backup_keystore is entered.
func (s *BasePlSqlParserListener) EnterBackup_keystore(ctx *Backup_keystoreContext) {}

// ExitBackup_keystore is called when production backup_keystore is exited.
func (s *BasePlSqlParserListener) ExitBackup_keystore(ctx *Backup_keystoreContext) {}

// EnterAlter_keystore_password is called when production alter_keystore_password is entered.
func (s *BasePlSqlParserListener) EnterAlter_keystore_password(ctx *Alter_keystore_passwordContext) {}

// ExitAlter_keystore_password is called when production alter_keystore_password is exited.
func (s *BasePlSqlParserListener) ExitAlter_keystore_password(ctx *Alter_keystore_passwordContext) {}

// EnterMerge_into_new_keystore is called when production merge_into_new_keystore is entered.
func (s *BasePlSqlParserListener) EnterMerge_into_new_keystore(ctx *Merge_into_new_keystoreContext) {}

// ExitMerge_into_new_keystore is called when production merge_into_new_keystore is exited.
func (s *BasePlSqlParserListener) ExitMerge_into_new_keystore(ctx *Merge_into_new_keystoreContext) {}

// EnterMerge_into_existing_keystore is called when production merge_into_existing_keystore is entered.
func (s *BasePlSqlParserListener) EnterMerge_into_existing_keystore(ctx *Merge_into_existing_keystoreContext) {
}

// ExitMerge_into_existing_keystore is called when production merge_into_existing_keystore is exited.
func (s *BasePlSqlParserListener) ExitMerge_into_existing_keystore(ctx *Merge_into_existing_keystoreContext) {
}

// EnterIsolate_keystore is called when production isolate_keystore is entered.
func (s *BasePlSqlParserListener) EnterIsolate_keystore(ctx *Isolate_keystoreContext) {}

// ExitIsolate_keystore is called when production isolate_keystore is exited.
func (s *BasePlSqlParserListener) ExitIsolate_keystore(ctx *Isolate_keystoreContext) {}

// EnterUnite_keystore is called when production unite_keystore is entered.
func (s *BasePlSqlParserListener) EnterUnite_keystore(ctx *Unite_keystoreContext) {}

// ExitUnite_keystore is called when production unite_keystore is exited.
func (s *BasePlSqlParserListener) ExitUnite_keystore(ctx *Unite_keystoreContext) {}

// EnterKey_management_clauses is called when production key_management_clauses is entered.
func (s *BasePlSqlParserListener) EnterKey_management_clauses(ctx *Key_management_clausesContext) {}

// ExitKey_management_clauses is called when production key_management_clauses is exited.
func (s *BasePlSqlParserListener) ExitKey_management_clauses(ctx *Key_management_clausesContext) {}

// EnterSet_key is called when production set_key is entered.
func (s *BasePlSqlParserListener) EnterSet_key(ctx *Set_keyContext) {}

// ExitSet_key is called when production set_key is exited.
func (s *BasePlSqlParserListener) ExitSet_key(ctx *Set_keyContext) {}

// EnterCreate_key is called when production create_key is entered.
func (s *BasePlSqlParserListener) EnterCreate_key(ctx *Create_keyContext) {}

// ExitCreate_key is called when production create_key is exited.
func (s *BasePlSqlParserListener) ExitCreate_key(ctx *Create_keyContext) {}

// EnterMkid is called when production mkid is entered.
func (s *BasePlSqlParserListener) EnterMkid(ctx *MkidContext) {}

// ExitMkid is called when production mkid is exited.
func (s *BasePlSqlParserListener) ExitMkid(ctx *MkidContext) {}

// EnterMk is called when production mk is entered.
func (s *BasePlSqlParserListener) EnterMk(ctx *MkContext) {}

// ExitMk is called when production mk is exited.
func (s *BasePlSqlParserListener) ExitMk(ctx *MkContext) {}

// EnterUse_key is called when production use_key is entered.
func (s *BasePlSqlParserListener) EnterUse_key(ctx *Use_keyContext) {}

// ExitUse_key is called when production use_key is exited.
func (s *BasePlSqlParserListener) ExitUse_key(ctx *Use_keyContext) {}

// EnterSet_key_tag is called when production set_key_tag is entered.
func (s *BasePlSqlParserListener) EnterSet_key_tag(ctx *Set_key_tagContext) {}

// ExitSet_key_tag is called when production set_key_tag is exited.
func (s *BasePlSqlParserListener) ExitSet_key_tag(ctx *Set_key_tagContext) {}

// EnterExport_keys is called when production export_keys is entered.
func (s *BasePlSqlParserListener) EnterExport_keys(ctx *Export_keysContext) {}

// ExitExport_keys is called when production export_keys is exited.
func (s *BasePlSqlParserListener) ExitExport_keys(ctx *Export_keysContext) {}

// EnterImport_keys is called when production import_keys is entered.
func (s *BasePlSqlParserListener) EnterImport_keys(ctx *Import_keysContext) {}

// ExitImport_keys is called when production import_keys is exited.
func (s *BasePlSqlParserListener) ExitImport_keys(ctx *Import_keysContext) {}

// EnterMigrate_keys is called when production migrate_keys is entered.
func (s *BasePlSqlParserListener) EnterMigrate_keys(ctx *Migrate_keysContext) {}

// ExitMigrate_keys is called when production migrate_keys is exited.
func (s *BasePlSqlParserListener) ExitMigrate_keys(ctx *Migrate_keysContext) {}

// EnterReverse_migrate_keys is called when production reverse_migrate_keys is entered.
func (s *BasePlSqlParserListener) EnterReverse_migrate_keys(ctx *Reverse_migrate_keysContext) {}

// ExitReverse_migrate_keys is called when production reverse_migrate_keys is exited.
func (s *BasePlSqlParserListener) ExitReverse_migrate_keys(ctx *Reverse_migrate_keysContext) {}

// EnterMove_keys is called when production move_keys is entered.
func (s *BasePlSqlParserListener) EnterMove_keys(ctx *Move_keysContext) {}

// ExitMove_keys is called when production move_keys is exited.
func (s *BasePlSqlParserListener) ExitMove_keys(ctx *Move_keysContext) {}

// EnterIdentified_by_store is called when production identified_by_store is entered.
func (s *BasePlSqlParserListener) EnterIdentified_by_store(ctx *Identified_by_storeContext) {}

// ExitIdentified_by_store is called when production identified_by_store is exited.
func (s *BasePlSqlParserListener) ExitIdentified_by_store(ctx *Identified_by_storeContext) {}

// EnterUsing_algorithm_clause is called when production using_algorithm_clause is entered.
func (s *BasePlSqlParserListener) EnterUsing_algorithm_clause(ctx *Using_algorithm_clauseContext) {}

// ExitUsing_algorithm_clause is called when production using_algorithm_clause is exited.
func (s *BasePlSqlParserListener) ExitUsing_algorithm_clause(ctx *Using_algorithm_clauseContext) {}

// EnterUsing_tag_clause is called when production using_tag_clause is entered.
func (s *BasePlSqlParserListener) EnterUsing_tag_clause(ctx *Using_tag_clauseContext) {}

// ExitUsing_tag_clause is called when production using_tag_clause is exited.
func (s *BasePlSqlParserListener) ExitUsing_tag_clause(ctx *Using_tag_clauseContext) {}

// EnterSecret_management_clauses is called when production secret_management_clauses is entered.
func (s *BasePlSqlParserListener) EnterSecret_management_clauses(ctx *Secret_management_clausesContext) {
}

// ExitSecret_management_clauses is called when production secret_management_clauses is exited.
func (s *BasePlSqlParserListener) ExitSecret_management_clauses(ctx *Secret_management_clausesContext) {
}

// EnterAdd_update_secret is called when production add_update_secret is entered.
func (s *BasePlSqlParserListener) EnterAdd_update_secret(ctx *Add_update_secretContext) {}

// ExitAdd_update_secret is called when production add_update_secret is exited.
func (s *BasePlSqlParserListener) ExitAdd_update_secret(ctx *Add_update_secretContext) {}

// EnterDelete_secret is called when production delete_secret is entered.
func (s *BasePlSqlParserListener) EnterDelete_secret(ctx *Delete_secretContext) {}

// ExitDelete_secret is called when production delete_secret is exited.
func (s *BasePlSqlParserListener) ExitDelete_secret(ctx *Delete_secretContext) {}

// EnterAdd_update_secret_seps is called when production add_update_secret_seps is entered.
func (s *BasePlSqlParserListener) EnterAdd_update_secret_seps(ctx *Add_update_secret_sepsContext) {}

// ExitAdd_update_secret_seps is called when production add_update_secret_seps is exited.
func (s *BasePlSqlParserListener) ExitAdd_update_secret_seps(ctx *Add_update_secret_sepsContext) {}

// EnterDelete_secret_seps is called when production delete_secret_seps is entered.
func (s *BasePlSqlParserListener) EnterDelete_secret_seps(ctx *Delete_secret_sepsContext) {}

// ExitDelete_secret_seps is called when production delete_secret_seps is exited.
func (s *BasePlSqlParserListener) ExitDelete_secret_seps(ctx *Delete_secret_sepsContext) {}

// EnterZero_downtime_software_patching_clauses is called when production zero_downtime_software_patching_clauses is entered.
func (s *BasePlSqlParserListener) EnterZero_downtime_software_patching_clauses(ctx *Zero_downtime_software_patching_clausesContext) {
}

// ExitZero_downtime_software_patching_clauses is called when production zero_downtime_software_patching_clauses is exited.
func (s *BasePlSqlParserListener) ExitZero_downtime_software_patching_clauses(ctx *Zero_downtime_software_patching_clausesContext) {
}

// EnterWith_backup_clause is called when production with_backup_clause is entered.
func (s *BasePlSqlParserListener) EnterWith_backup_clause(ctx *With_backup_clauseContext) {}

// ExitWith_backup_clause is called when production with_backup_clause is exited.
func (s *BasePlSqlParserListener) ExitWith_backup_clause(ctx *With_backup_clauseContext) {}

// EnterIdentified_by_password_clause is called when production identified_by_password_clause is entered.
func (s *BasePlSqlParserListener) EnterIdentified_by_password_clause(ctx *Identified_by_password_clauseContext) {
}

// ExitIdentified_by_password_clause is called when production identified_by_password_clause is exited.
func (s *BasePlSqlParserListener) ExitIdentified_by_password_clause(ctx *Identified_by_password_clauseContext) {
}

// EnterKeystore_password is called when production keystore_password is entered.
func (s *BasePlSqlParserListener) EnterKeystore_password(ctx *Keystore_passwordContext) {}

// ExitKeystore_password is called when production keystore_password is exited.
func (s *BasePlSqlParserListener) ExitKeystore_password(ctx *Keystore_passwordContext) {}

// EnterPath is called when production path is entered.
func (s *BasePlSqlParserListener) EnterPath(ctx *PathContext) {}

// ExitPath is called when production path is exited.
func (s *BasePlSqlParserListener) ExitPath(ctx *PathContext) {}

// EnterSecret is called when production secret is entered.
func (s *BasePlSqlParserListener) EnterSecret(ctx *SecretContext) {}

// ExitSecret is called when production secret is exited.
func (s *BasePlSqlParserListener) ExitSecret(ctx *SecretContext) {}

// EnterAnalyze is called when production analyze is entered.
func (s *BasePlSqlParserListener) EnterAnalyze(ctx *AnalyzeContext) {}

// ExitAnalyze is called when production analyze is exited.
func (s *BasePlSqlParserListener) ExitAnalyze(ctx *AnalyzeContext) {}

// EnterPartition_extention_clause is called when production partition_extention_clause is entered.
func (s *BasePlSqlParserListener) EnterPartition_extention_clause(ctx *Partition_extention_clauseContext) {
}

// ExitPartition_extention_clause is called when production partition_extention_clause is exited.
func (s *BasePlSqlParserListener) ExitPartition_extention_clause(ctx *Partition_extention_clauseContext) {
}

// EnterValidation_clauses is called when production validation_clauses is entered.
func (s *BasePlSqlParserListener) EnterValidation_clauses(ctx *Validation_clausesContext) {}

// ExitValidation_clauses is called when production validation_clauses is exited.
func (s *BasePlSqlParserListener) ExitValidation_clauses(ctx *Validation_clausesContext) {}

// EnterCompute_clauses is called when production compute_clauses is entered.
func (s *BasePlSqlParserListener) EnterCompute_clauses(ctx *Compute_clausesContext) {}

// ExitCompute_clauses is called when production compute_clauses is exited.
func (s *BasePlSqlParserListener) ExitCompute_clauses(ctx *Compute_clausesContext) {}

// EnterFor_clause is called when production for_clause is entered.
func (s *BasePlSqlParserListener) EnterFor_clause(ctx *For_clauseContext) {}

// ExitFor_clause is called when production for_clause is exited.
func (s *BasePlSqlParserListener) ExitFor_clause(ctx *For_clauseContext) {}

// EnterOnline_or_offline is called when production online_or_offline is entered.
func (s *BasePlSqlParserListener) EnterOnline_or_offline(ctx *Online_or_offlineContext) {}

// ExitOnline_or_offline is called when production online_or_offline is exited.
func (s *BasePlSqlParserListener) ExitOnline_or_offline(ctx *Online_or_offlineContext) {}

// EnterInto_clause1 is called when production into_clause1 is entered.
func (s *BasePlSqlParserListener) EnterInto_clause1(ctx *Into_clause1Context) {}

// ExitInto_clause1 is called when production into_clause1 is exited.
func (s *BasePlSqlParserListener) ExitInto_clause1(ctx *Into_clause1Context) {}

// EnterPartition_key_value is called when production partition_key_value is entered.
func (s *BasePlSqlParserListener) EnterPartition_key_value(ctx *Partition_key_valueContext) {}

// ExitPartition_key_value is called when production partition_key_value is exited.
func (s *BasePlSqlParserListener) ExitPartition_key_value(ctx *Partition_key_valueContext) {}

// EnterSubpartition_key_value is called when production subpartition_key_value is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_key_value(ctx *Subpartition_key_valueContext) {}

// ExitSubpartition_key_value is called when production subpartition_key_value is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_key_value(ctx *Subpartition_key_valueContext) {}

// EnterAssociate_statistics is called when production associate_statistics is entered.
func (s *BasePlSqlParserListener) EnterAssociate_statistics(ctx *Associate_statisticsContext) {}

// ExitAssociate_statistics is called when production associate_statistics is exited.
func (s *BasePlSqlParserListener) ExitAssociate_statistics(ctx *Associate_statisticsContext) {}

// EnterColumn_association is called when production column_association is entered.
func (s *BasePlSqlParserListener) EnterColumn_association(ctx *Column_associationContext) {}

// ExitColumn_association is called when production column_association is exited.
func (s *BasePlSqlParserListener) ExitColumn_association(ctx *Column_associationContext) {}

// EnterFunction_association is called when production function_association is entered.
func (s *BasePlSqlParserListener) EnterFunction_association(ctx *Function_associationContext) {}

// ExitFunction_association is called when production function_association is exited.
func (s *BasePlSqlParserListener) ExitFunction_association(ctx *Function_associationContext) {}

// EnterIndextype_name is called when production indextype_name is entered.
func (s *BasePlSqlParserListener) EnterIndextype_name(ctx *Indextype_nameContext) {}

// ExitIndextype_name is called when production indextype_name is exited.
func (s *BasePlSqlParserListener) ExitIndextype_name(ctx *Indextype_nameContext) {}

// EnterUsing_statistics_type is called when production using_statistics_type is entered.
func (s *BasePlSqlParserListener) EnterUsing_statistics_type(ctx *Using_statistics_typeContext) {}

// ExitUsing_statistics_type is called when production using_statistics_type is exited.
func (s *BasePlSqlParserListener) ExitUsing_statistics_type(ctx *Using_statistics_typeContext) {}

// EnterStatistics_type_name is called when production statistics_type_name is entered.
func (s *BasePlSqlParserListener) EnterStatistics_type_name(ctx *Statistics_type_nameContext) {}

// ExitStatistics_type_name is called when production statistics_type_name is exited.
func (s *BasePlSqlParserListener) ExitStatistics_type_name(ctx *Statistics_type_nameContext) {}

// EnterDefault_cost_clause is called when production default_cost_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_cost_clause(ctx *Default_cost_clauseContext) {}

// ExitDefault_cost_clause is called when production default_cost_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_cost_clause(ctx *Default_cost_clauseContext) {}

// EnterCpu_cost is called when production cpu_cost is entered.
func (s *BasePlSqlParserListener) EnterCpu_cost(ctx *Cpu_costContext) {}

// ExitCpu_cost is called when production cpu_cost is exited.
func (s *BasePlSqlParserListener) ExitCpu_cost(ctx *Cpu_costContext) {}

// EnterIo_cost is called when production io_cost is entered.
func (s *BasePlSqlParserListener) EnterIo_cost(ctx *Io_costContext) {}

// ExitIo_cost is called when production io_cost is exited.
func (s *BasePlSqlParserListener) ExitIo_cost(ctx *Io_costContext) {}

// EnterNetwork_cost is called when production network_cost is entered.
func (s *BasePlSqlParserListener) EnterNetwork_cost(ctx *Network_costContext) {}

// ExitNetwork_cost is called when production network_cost is exited.
func (s *BasePlSqlParserListener) ExitNetwork_cost(ctx *Network_costContext) {}

// EnterDefault_selectivity_clause is called when production default_selectivity_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_selectivity_clause(ctx *Default_selectivity_clauseContext) {
}

// ExitDefault_selectivity_clause is called when production default_selectivity_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_selectivity_clause(ctx *Default_selectivity_clauseContext) {
}

// EnterDefault_selectivity is called when production default_selectivity is entered.
func (s *BasePlSqlParserListener) EnterDefault_selectivity(ctx *Default_selectivityContext) {}

// ExitDefault_selectivity is called when production default_selectivity is exited.
func (s *BasePlSqlParserListener) ExitDefault_selectivity(ctx *Default_selectivityContext) {}

// EnterStorage_table_clause is called when production storage_table_clause is entered.
func (s *BasePlSqlParserListener) EnterStorage_table_clause(ctx *Storage_table_clauseContext) {}

// ExitStorage_table_clause is called when production storage_table_clause is exited.
func (s *BasePlSqlParserListener) ExitStorage_table_clause(ctx *Storage_table_clauseContext) {}

// EnterUnified_auditing is called when production unified_auditing is entered.
func (s *BasePlSqlParserListener) EnterUnified_auditing(ctx *Unified_auditingContext) {}

// ExitUnified_auditing is called when production unified_auditing is exited.
func (s *BasePlSqlParserListener) ExitUnified_auditing(ctx *Unified_auditingContext) {}

// EnterPolicy_name is called when production policy_name is entered.
func (s *BasePlSqlParserListener) EnterPolicy_name(ctx *Policy_nameContext) {}

// ExitPolicy_name is called when production policy_name is exited.
func (s *BasePlSqlParserListener) ExitPolicy_name(ctx *Policy_nameContext) {}

// EnterAudit_traditional is called when production audit_traditional is entered.
func (s *BasePlSqlParserListener) EnterAudit_traditional(ctx *Audit_traditionalContext) {}

// ExitAudit_traditional is called when production audit_traditional is exited.
func (s *BasePlSqlParserListener) ExitAudit_traditional(ctx *Audit_traditionalContext) {}

// EnterAudit_direct_path is called when production audit_direct_path is entered.
func (s *BasePlSqlParserListener) EnterAudit_direct_path(ctx *Audit_direct_pathContext) {}

// ExitAudit_direct_path is called when production audit_direct_path is exited.
func (s *BasePlSqlParserListener) ExitAudit_direct_path(ctx *Audit_direct_pathContext) {}

// EnterAudit_container_clause is called when production audit_container_clause is entered.
func (s *BasePlSqlParserListener) EnterAudit_container_clause(ctx *Audit_container_clauseContext) {}

// ExitAudit_container_clause is called when production audit_container_clause is exited.
func (s *BasePlSqlParserListener) ExitAudit_container_clause(ctx *Audit_container_clauseContext) {}

// EnterAudit_operation_clause is called when production audit_operation_clause is entered.
func (s *BasePlSqlParserListener) EnterAudit_operation_clause(ctx *Audit_operation_clauseContext) {}

// ExitAudit_operation_clause is called when production audit_operation_clause is exited.
func (s *BasePlSqlParserListener) ExitAudit_operation_clause(ctx *Audit_operation_clauseContext) {}

// EnterAuditing_by_clause is called when production auditing_by_clause is entered.
func (s *BasePlSqlParserListener) EnterAuditing_by_clause(ctx *Auditing_by_clauseContext) {}

// ExitAuditing_by_clause is called when production auditing_by_clause is exited.
func (s *BasePlSqlParserListener) ExitAuditing_by_clause(ctx *Auditing_by_clauseContext) {}

// EnterAudit_user is called when production audit_user is entered.
func (s *BasePlSqlParserListener) EnterAudit_user(ctx *Audit_userContext) {}

// ExitAudit_user is called when production audit_user is exited.
func (s *BasePlSqlParserListener) ExitAudit_user(ctx *Audit_userContext) {}

// EnterAudit_schema_object_clause is called when production audit_schema_object_clause is entered.
func (s *BasePlSqlParserListener) EnterAudit_schema_object_clause(ctx *Audit_schema_object_clauseContext) {
}

// ExitAudit_schema_object_clause is called when production audit_schema_object_clause is exited.
func (s *BasePlSqlParserListener) ExitAudit_schema_object_clause(ctx *Audit_schema_object_clauseContext) {
}

// EnterSql_operation is called when production sql_operation is entered.
func (s *BasePlSqlParserListener) EnterSql_operation(ctx *Sql_operationContext) {}

// ExitSql_operation is called when production sql_operation is exited.
func (s *BasePlSqlParserListener) ExitSql_operation(ctx *Sql_operationContext) {}

// EnterAuditing_on_clause is called when production auditing_on_clause is entered.
func (s *BasePlSqlParserListener) EnterAuditing_on_clause(ctx *Auditing_on_clauseContext) {}

// ExitAuditing_on_clause is called when production auditing_on_clause is exited.
func (s *BasePlSqlParserListener) ExitAuditing_on_clause(ctx *Auditing_on_clauseContext) {}

// EnterModel_name is called when production model_name is entered.
func (s *BasePlSqlParserListener) EnterModel_name(ctx *Model_nameContext) {}

// ExitModel_name is called when production model_name is exited.
func (s *BasePlSqlParserListener) ExitModel_name(ctx *Model_nameContext) {}

// EnterObject_name is called when production object_name is entered.
func (s *BasePlSqlParserListener) EnterObject_name(ctx *Object_nameContext) {}

// ExitObject_name is called when production object_name is exited.
func (s *BasePlSqlParserListener) ExitObject_name(ctx *Object_nameContext) {}

// EnterProfile_name is called when production profile_name is entered.
func (s *BasePlSqlParserListener) EnterProfile_name(ctx *Profile_nameContext) {}

// ExitProfile_name is called when production profile_name is exited.
func (s *BasePlSqlParserListener) ExitProfile_name(ctx *Profile_nameContext) {}

// EnterSql_statement_shortcut is called when production sql_statement_shortcut is entered.
func (s *BasePlSqlParserListener) EnterSql_statement_shortcut(ctx *Sql_statement_shortcutContext) {}

// ExitSql_statement_shortcut is called when production sql_statement_shortcut is exited.
func (s *BasePlSqlParserListener) ExitSql_statement_shortcut(ctx *Sql_statement_shortcutContext) {}

// EnterDrop_index is called when production drop_index is entered.
func (s *BasePlSqlParserListener) EnterDrop_index(ctx *Drop_indexContext) {}

// ExitDrop_index is called when production drop_index is exited.
func (s *BasePlSqlParserListener) ExitDrop_index(ctx *Drop_indexContext) {}

// EnterDisassociate_statistics is called when production disassociate_statistics is entered.
func (s *BasePlSqlParserListener) EnterDisassociate_statistics(ctx *Disassociate_statisticsContext) {}

// ExitDisassociate_statistics is called when production disassociate_statistics is exited.
func (s *BasePlSqlParserListener) ExitDisassociate_statistics(ctx *Disassociate_statisticsContext) {}

// EnterDrop_indextype is called when production drop_indextype is entered.
func (s *BasePlSqlParserListener) EnterDrop_indextype(ctx *Drop_indextypeContext) {}

// ExitDrop_indextype is called when production drop_indextype is exited.
func (s *BasePlSqlParserListener) ExitDrop_indextype(ctx *Drop_indextypeContext) {}

// EnterDrop_inmemory_join_group is called when production drop_inmemory_join_group is entered.
func (s *BasePlSqlParserListener) EnterDrop_inmemory_join_group(ctx *Drop_inmemory_join_groupContext) {
}

// ExitDrop_inmemory_join_group is called when production drop_inmemory_join_group is exited.
func (s *BasePlSqlParserListener) ExitDrop_inmemory_join_group(ctx *Drop_inmemory_join_groupContext) {
}

// EnterFlashback_table is called when production flashback_table is entered.
func (s *BasePlSqlParserListener) EnterFlashback_table(ctx *Flashback_tableContext) {}

// ExitFlashback_table is called when production flashback_table is exited.
func (s *BasePlSqlParserListener) ExitFlashback_table(ctx *Flashback_tableContext) {}

// EnterRestore_point is called when production restore_point is entered.
func (s *BasePlSqlParserListener) EnterRestore_point(ctx *Restore_pointContext) {}

// ExitRestore_point is called when production restore_point is exited.
func (s *BasePlSqlParserListener) ExitRestore_point(ctx *Restore_pointContext) {}

// EnterPurge_statement is called when production purge_statement is entered.
func (s *BasePlSqlParserListener) EnterPurge_statement(ctx *Purge_statementContext) {}

// ExitPurge_statement is called when production purge_statement is exited.
func (s *BasePlSqlParserListener) ExitPurge_statement(ctx *Purge_statementContext) {}

// EnterNoaudit_statement is called when production noaudit_statement is entered.
func (s *BasePlSqlParserListener) EnterNoaudit_statement(ctx *Noaudit_statementContext) {}

// ExitNoaudit_statement is called when production noaudit_statement is exited.
func (s *BasePlSqlParserListener) ExitNoaudit_statement(ctx *Noaudit_statementContext) {}

// EnterRename_object is called when production rename_object is entered.
func (s *BasePlSqlParserListener) EnterRename_object(ctx *Rename_objectContext) {}

// ExitRename_object is called when production rename_object is exited.
func (s *BasePlSqlParserListener) ExitRename_object(ctx *Rename_objectContext) {}

// EnterGrant_statement is called when production grant_statement is entered.
func (s *BasePlSqlParserListener) EnterGrant_statement(ctx *Grant_statementContext) {}

// ExitGrant_statement is called when production grant_statement is exited.
func (s *BasePlSqlParserListener) ExitGrant_statement(ctx *Grant_statementContext) {}

// EnterContainer_clause is called when production container_clause is entered.
func (s *BasePlSqlParserListener) EnterContainer_clause(ctx *Container_clauseContext) {}

// ExitContainer_clause is called when production container_clause is exited.
func (s *BasePlSqlParserListener) ExitContainer_clause(ctx *Container_clauseContext) {}

// EnterRevoke_statement is called when production revoke_statement is entered.
func (s *BasePlSqlParserListener) EnterRevoke_statement(ctx *Revoke_statementContext) {}

// ExitRevoke_statement is called when production revoke_statement is exited.
func (s *BasePlSqlParserListener) ExitRevoke_statement(ctx *Revoke_statementContext) {}

// EnterRevoke_system_privilege is called when production revoke_system_privilege is entered.
func (s *BasePlSqlParserListener) EnterRevoke_system_privilege(ctx *Revoke_system_privilegeContext) {}

// ExitRevoke_system_privilege is called when production revoke_system_privilege is exited.
func (s *BasePlSqlParserListener) ExitRevoke_system_privilege(ctx *Revoke_system_privilegeContext) {}

// EnterRevokee_clause is called when production revokee_clause is entered.
func (s *BasePlSqlParserListener) EnterRevokee_clause(ctx *Revokee_clauseContext) {}

// ExitRevokee_clause is called when production revokee_clause is exited.
func (s *BasePlSqlParserListener) ExitRevokee_clause(ctx *Revokee_clauseContext) {}

// EnterRevoke_object_privileges is called when production revoke_object_privileges is entered.
func (s *BasePlSqlParserListener) EnterRevoke_object_privileges(ctx *Revoke_object_privilegesContext) {
}

// ExitRevoke_object_privileges is called when production revoke_object_privileges is exited.
func (s *BasePlSqlParserListener) ExitRevoke_object_privileges(ctx *Revoke_object_privilegesContext) {
}

// EnterOn_object_clause is called when production on_object_clause is entered.
func (s *BasePlSqlParserListener) EnterOn_object_clause(ctx *On_object_clauseContext) {}

// ExitOn_object_clause is called when production on_object_clause is exited.
func (s *BasePlSqlParserListener) ExitOn_object_clause(ctx *On_object_clauseContext) {}

// EnterRevoke_roles_from_programs is called when production revoke_roles_from_programs is entered.
func (s *BasePlSqlParserListener) EnterRevoke_roles_from_programs(ctx *Revoke_roles_from_programsContext) {
}

// ExitRevoke_roles_from_programs is called when production revoke_roles_from_programs is exited.
func (s *BasePlSqlParserListener) ExitRevoke_roles_from_programs(ctx *Revoke_roles_from_programsContext) {
}

// EnterProgram_unit is called when production program_unit is entered.
func (s *BasePlSqlParserListener) EnterProgram_unit(ctx *Program_unitContext) {}

// ExitProgram_unit is called when production program_unit is exited.
func (s *BasePlSqlParserListener) ExitProgram_unit(ctx *Program_unitContext) {}

// EnterCreate_dimension is called when production create_dimension is entered.
func (s *BasePlSqlParserListener) EnterCreate_dimension(ctx *Create_dimensionContext) {}

// ExitCreate_dimension is called when production create_dimension is exited.
func (s *BasePlSqlParserListener) ExitCreate_dimension(ctx *Create_dimensionContext) {}

// EnterCreate_directory is called when production create_directory is entered.
func (s *BasePlSqlParserListener) EnterCreate_directory(ctx *Create_directoryContext) {}

// ExitCreate_directory is called when production create_directory is exited.
func (s *BasePlSqlParserListener) ExitCreate_directory(ctx *Create_directoryContext) {}

// EnterDirectory_name is called when production directory_name is entered.
func (s *BasePlSqlParserListener) EnterDirectory_name(ctx *Directory_nameContext) {}

// ExitDirectory_name is called when production directory_name is exited.
func (s *BasePlSqlParserListener) ExitDirectory_name(ctx *Directory_nameContext) {}

// EnterDirectory_path is called when production directory_path is entered.
func (s *BasePlSqlParserListener) EnterDirectory_path(ctx *Directory_pathContext) {}

// ExitDirectory_path is called when production directory_path is exited.
func (s *BasePlSqlParserListener) ExitDirectory_path(ctx *Directory_pathContext) {}

// EnterCreate_inmemory_join_group is called when production create_inmemory_join_group is entered.
func (s *BasePlSqlParserListener) EnterCreate_inmemory_join_group(ctx *Create_inmemory_join_groupContext) {
}

// ExitCreate_inmemory_join_group is called when production create_inmemory_join_group is exited.
func (s *BasePlSqlParserListener) ExitCreate_inmemory_join_group(ctx *Create_inmemory_join_groupContext) {
}

// EnterDrop_hierarchy is called when production drop_hierarchy is entered.
func (s *BasePlSqlParserListener) EnterDrop_hierarchy(ctx *Drop_hierarchyContext) {}

// ExitDrop_hierarchy is called when production drop_hierarchy is exited.
func (s *BasePlSqlParserListener) ExitDrop_hierarchy(ctx *Drop_hierarchyContext) {}

// EnterAlter_library is called when production alter_library is entered.
func (s *BasePlSqlParserListener) EnterAlter_library(ctx *Alter_libraryContext) {}

// ExitAlter_library is called when production alter_library is exited.
func (s *BasePlSqlParserListener) ExitAlter_library(ctx *Alter_libraryContext) {}

// EnterDrop_java is called when production drop_java is entered.
func (s *BasePlSqlParserListener) EnterDrop_java(ctx *Drop_javaContext) {}

// ExitDrop_java is called when production drop_java is exited.
func (s *BasePlSqlParserListener) ExitDrop_java(ctx *Drop_javaContext) {}

// EnterDrop_library is called when production drop_library is entered.
func (s *BasePlSqlParserListener) EnterDrop_library(ctx *Drop_libraryContext) {}

// ExitDrop_library is called when production drop_library is exited.
func (s *BasePlSqlParserListener) ExitDrop_library(ctx *Drop_libraryContext) {}

// EnterCreate_java is called when production create_java is entered.
func (s *BasePlSqlParserListener) EnterCreate_java(ctx *Create_javaContext) {}

// ExitCreate_java is called when production create_java is exited.
func (s *BasePlSqlParserListener) ExitCreate_java(ctx *Create_javaContext) {}

// EnterCreate_library is called when production create_library is entered.
func (s *BasePlSqlParserListener) EnterCreate_library(ctx *Create_libraryContext) {}

// ExitCreate_library is called when production create_library is exited.
func (s *BasePlSqlParserListener) ExitCreate_library(ctx *Create_libraryContext) {}

// EnterPlsql_library_source is called when production plsql_library_source is entered.
func (s *BasePlSqlParserListener) EnterPlsql_library_source(ctx *Plsql_library_sourceContext) {}

// ExitPlsql_library_source is called when production plsql_library_source is exited.
func (s *BasePlSqlParserListener) ExitPlsql_library_source(ctx *Plsql_library_sourceContext) {}

// EnterCredential_name is called when production credential_name is entered.
func (s *BasePlSqlParserListener) EnterCredential_name(ctx *Credential_nameContext) {}

// ExitCredential_name is called when production credential_name is exited.
func (s *BasePlSqlParserListener) ExitCredential_name(ctx *Credential_nameContext) {}

// EnterLibrary_editionable is called when production library_editionable is entered.
func (s *BasePlSqlParserListener) EnterLibrary_editionable(ctx *Library_editionableContext) {}

// ExitLibrary_editionable is called when production library_editionable is exited.
func (s *BasePlSqlParserListener) ExitLibrary_editionable(ctx *Library_editionableContext) {}

// EnterLibrary_debug is called when production library_debug is entered.
func (s *BasePlSqlParserListener) EnterLibrary_debug(ctx *Library_debugContext) {}

// ExitLibrary_debug is called when production library_debug is exited.
func (s *BasePlSqlParserListener) ExitLibrary_debug(ctx *Library_debugContext) {}

// EnterCompiler_parameters_clause is called when production compiler_parameters_clause is entered.
func (s *BasePlSqlParserListener) EnterCompiler_parameters_clause(ctx *Compiler_parameters_clauseContext) {
}

// ExitCompiler_parameters_clause is called when production compiler_parameters_clause is exited.
func (s *BasePlSqlParserListener) ExitCompiler_parameters_clause(ctx *Compiler_parameters_clauseContext) {
}

// EnterParameter_value is called when production parameter_value is entered.
func (s *BasePlSqlParserListener) EnterParameter_value(ctx *Parameter_valueContext) {}

// ExitParameter_value is called when production parameter_value is exited.
func (s *BasePlSqlParserListener) ExitParameter_value(ctx *Parameter_valueContext) {}

// EnterLibrary_name is called when production library_name is entered.
func (s *BasePlSqlParserListener) EnterLibrary_name(ctx *Library_nameContext) {}

// ExitLibrary_name is called when production library_name is exited.
func (s *BasePlSqlParserListener) ExitLibrary_name(ctx *Library_nameContext) {}

// EnterAlter_dimension is called when production alter_dimension is entered.
func (s *BasePlSqlParserListener) EnterAlter_dimension(ctx *Alter_dimensionContext) {}

// ExitAlter_dimension is called when production alter_dimension is exited.
func (s *BasePlSqlParserListener) ExitAlter_dimension(ctx *Alter_dimensionContext) {}

// EnterLevel_clause is called when production level_clause is entered.
func (s *BasePlSqlParserListener) EnterLevel_clause(ctx *Level_clauseContext) {}

// ExitLevel_clause is called when production level_clause is exited.
func (s *BasePlSqlParserListener) ExitLevel_clause(ctx *Level_clauseContext) {}

// EnterHierarchy_clause is called when production hierarchy_clause is entered.
func (s *BasePlSqlParserListener) EnterHierarchy_clause(ctx *Hierarchy_clauseContext) {}

// ExitHierarchy_clause is called when production hierarchy_clause is exited.
func (s *BasePlSqlParserListener) ExitHierarchy_clause(ctx *Hierarchy_clauseContext) {}

// EnterDimension_join_clause is called when production dimension_join_clause is entered.
func (s *BasePlSqlParserListener) EnterDimension_join_clause(ctx *Dimension_join_clauseContext) {}

// ExitDimension_join_clause is called when production dimension_join_clause is exited.
func (s *BasePlSqlParserListener) ExitDimension_join_clause(ctx *Dimension_join_clauseContext) {}

// EnterAttribute_clause is called when production attribute_clause is entered.
func (s *BasePlSqlParserListener) EnterAttribute_clause(ctx *Attribute_clauseContext) {}

// ExitAttribute_clause is called when production attribute_clause is exited.
func (s *BasePlSqlParserListener) ExitAttribute_clause(ctx *Attribute_clauseContext) {}

// EnterExtended_attribute_clause is called when production extended_attribute_clause is entered.
func (s *BasePlSqlParserListener) EnterExtended_attribute_clause(ctx *Extended_attribute_clauseContext) {
}

// ExitExtended_attribute_clause is called when production extended_attribute_clause is exited.
func (s *BasePlSqlParserListener) ExitExtended_attribute_clause(ctx *Extended_attribute_clauseContext) {
}

// EnterColumn_one_or_more_sub_clause is called when production column_one_or_more_sub_clause is entered.
func (s *BasePlSqlParserListener) EnterColumn_one_or_more_sub_clause(ctx *Column_one_or_more_sub_clauseContext) {
}

// ExitColumn_one_or_more_sub_clause is called when production column_one_or_more_sub_clause is exited.
func (s *BasePlSqlParserListener) ExitColumn_one_or_more_sub_clause(ctx *Column_one_or_more_sub_clauseContext) {
}

// EnterAlter_view is called when production alter_view is entered.
func (s *BasePlSqlParserListener) EnterAlter_view(ctx *Alter_viewContext) {}

// ExitAlter_view is called when production alter_view is exited.
func (s *BasePlSqlParserListener) ExitAlter_view(ctx *Alter_viewContext) {}

// EnterAlter_view_editionable is called when production alter_view_editionable is entered.
func (s *BasePlSqlParserListener) EnterAlter_view_editionable(ctx *Alter_view_editionableContext) {}

// ExitAlter_view_editionable is called when production alter_view_editionable is exited.
func (s *BasePlSqlParserListener) ExitAlter_view_editionable(ctx *Alter_view_editionableContext) {}

// EnterCreate_view is called when production create_view is entered.
func (s *BasePlSqlParserListener) EnterCreate_view(ctx *Create_viewContext) {}

// ExitCreate_view is called when production create_view is exited.
func (s *BasePlSqlParserListener) ExitCreate_view(ctx *Create_viewContext) {}

// EnterEditioning_clause is called when production editioning_clause is entered.
func (s *BasePlSqlParserListener) EnterEditioning_clause(ctx *Editioning_clauseContext) {}

// ExitEditioning_clause is called when production editioning_clause is exited.
func (s *BasePlSqlParserListener) ExitEditioning_clause(ctx *Editioning_clauseContext) {}

// EnterView_options is called when production view_options is entered.
func (s *BasePlSqlParserListener) EnterView_options(ctx *View_optionsContext) {}

// ExitView_options is called when production view_options is exited.
func (s *BasePlSqlParserListener) ExitView_options(ctx *View_optionsContext) {}

// EnterView_alias_constraint is called when production view_alias_constraint is entered.
func (s *BasePlSqlParserListener) EnterView_alias_constraint(ctx *View_alias_constraintContext) {}

// ExitView_alias_constraint is called when production view_alias_constraint is exited.
func (s *BasePlSqlParserListener) ExitView_alias_constraint(ctx *View_alias_constraintContext) {}

// EnterObject_view_clause is called when production object_view_clause is entered.
func (s *BasePlSqlParserListener) EnterObject_view_clause(ctx *Object_view_clauseContext) {}

// ExitObject_view_clause is called when production object_view_clause is exited.
func (s *BasePlSqlParserListener) ExitObject_view_clause(ctx *Object_view_clauseContext) {}

// EnterInline_constraint is called when production inline_constraint is entered.
func (s *BasePlSqlParserListener) EnterInline_constraint(ctx *Inline_constraintContext) {}

// ExitInline_constraint is called when production inline_constraint is exited.
func (s *BasePlSqlParserListener) ExitInline_constraint(ctx *Inline_constraintContext) {}

// EnterInline_ref_constraint is called when production inline_ref_constraint is entered.
func (s *BasePlSqlParserListener) EnterInline_ref_constraint(ctx *Inline_ref_constraintContext) {}

// ExitInline_ref_constraint is called when production inline_ref_constraint is exited.
func (s *BasePlSqlParserListener) ExitInline_ref_constraint(ctx *Inline_ref_constraintContext) {}

// EnterOut_of_line_ref_constraint is called when production out_of_line_ref_constraint is entered.
func (s *BasePlSqlParserListener) EnterOut_of_line_ref_constraint(ctx *Out_of_line_ref_constraintContext) {
}

// ExitOut_of_line_ref_constraint is called when production out_of_line_ref_constraint is exited.
func (s *BasePlSqlParserListener) ExitOut_of_line_ref_constraint(ctx *Out_of_line_ref_constraintContext) {
}

// EnterOut_of_line_constraint is called when production out_of_line_constraint is entered.
func (s *BasePlSqlParserListener) EnterOut_of_line_constraint(ctx *Out_of_line_constraintContext) {}

// ExitOut_of_line_constraint is called when production out_of_line_constraint is exited.
func (s *BasePlSqlParserListener) ExitOut_of_line_constraint(ctx *Out_of_line_constraintContext) {}

// EnterConstraint_state is called when production constraint_state is entered.
func (s *BasePlSqlParserListener) EnterConstraint_state(ctx *Constraint_stateContext) {}

// ExitConstraint_state is called when production constraint_state is exited.
func (s *BasePlSqlParserListener) ExitConstraint_state(ctx *Constraint_stateContext) {}

// EnterXmltype_view_clause is called when production xmltype_view_clause is entered.
func (s *BasePlSqlParserListener) EnterXmltype_view_clause(ctx *Xmltype_view_clauseContext) {}

// ExitXmltype_view_clause is called when production xmltype_view_clause is exited.
func (s *BasePlSqlParserListener) ExitXmltype_view_clause(ctx *Xmltype_view_clauseContext) {}

// EnterXml_schema_spec is called when production xml_schema_spec is entered.
func (s *BasePlSqlParserListener) EnterXml_schema_spec(ctx *Xml_schema_specContext) {}

// ExitXml_schema_spec is called when production xml_schema_spec is exited.
func (s *BasePlSqlParserListener) ExitXml_schema_spec(ctx *Xml_schema_specContext) {}

// EnterXml_schema_url is called when production xml_schema_url is entered.
func (s *BasePlSqlParserListener) EnterXml_schema_url(ctx *Xml_schema_urlContext) {}

// ExitXml_schema_url is called when production xml_schema_url is exited.
func (s *BasePlSqlParserListener) ExitXml_schema_url(ctx *Xml_schema_urlContext) {}

// EnterElement is called when production element is entered.
func (s *BasePlSqlParserListener) EnterElement(ctx *ElementContext) {}

// ExitElement is called when production element is exited.
func (s *BasePlSqlParserListener) ExitElement(ctx *ElementContext) {}

// EnterAlter_tablespace is called when production alter_tablespace is entered.
func (s *BasePlSqlParserListener) EnterAlter_tablespace(ctx *Alter_tablespaceContext) {}

// ExitAlter_tablespace is called when production alter_tablespace is exited.
func (s *BasePlSqlParserListener) ExitAlter_tablespace(ctx *Alter_tablespaceContext) {}

// EnterDatafile_tempfile_clauses is called when production datafile_tempfile_clauses is entered.
func (s *BasePlSqlParserListener) EnterDatafile_tempfile_clauses(ctx *Datafile_tempfile_clausesContext) {
}

// ExitDatafile_tempfile_clauses is called when production datafile_tempfile_clauses is exited.
func (s *BasePlSqlParserListener) ExitDatafile_tempfile_clauses(ctx *Datafile_tempfile_clausesContext) {
}

// EnterTablespace_logging_clauses is called when production tablespace_logging_clauses is entered.
func (s *BasePlSqlParserListener) EnterTablespace_logging_clauses(ctx *Tablespace_logging_clausesContext) {
}

// ExitTablespace_logging_clauses is called when production tablespace_logging_clauses is exited.
func (s *BasePlSqlParserListener) ExitTablespace_logging_clauses(ctx *Tablespace_logging_clausesContext) {
}

// EnterTablespace_group_clause is called when production tablespace_group_clause is entered.
func (s *BasePlSqlParserListener) EnterTablespace_group_clause(ctx *Tablespace_group_clauseContext) {}

// ExitTablespace_group_clause is called when production tablespace_group_clause is exited.
func (s *BasePlSqlParserListener) ExitTablespace_group_clause(ctx *Tablespace_group_clauseContext) {}

// EnterTablespace_group_name is called when production tablespace_group_name is entered.
func (s *BasePlSqlParserListener) EnterTablespace_group_name(ctx *Tablespace_group_nameContext) {}

// ExitTablespace_group_name is called when production tablespace_group_name is exited.
func (s *BasePlSqlParserListener) ExitTablespace_group_name(ctx *Tablespace_group_nameContext) {}

// EnterTablespace_state_clauses is called when production tablespace_state_clauses is entered.
func (s *BasePlSqlParserListener) EnterTablespace_state_clauses(ctx *Tablespace_state_clausesContext) {
}

// ExitTablespace_state_clauses is called when production tablespace_state_clauses is exited.
func (s *BasePlSqlParserListener) ExitTablespace_state_clauses(ctx *Tablespace_state_clausesContext) {
}

// EnterFlashback_mode_clause is called when production flashback_mode_clause is entered.
func (s *BasePlSqlParserListener) EnterFlashback_mode_clause(ctx *Flashback_mode_clauseContext) {}

// ExitFlashback_mode_clause is called when production flashback_mode_clause is exited.
func (s *BasePlSqlParserListener) ExitFlashback_mode_clause(ctx *Flashback_mode_clauseContext) {}

// EnterNew_tablespace_name is called when production new_tablespace_name is entered.
func (s *BasePlSqlParserListener) EnterNew_tablespace_name(ctx *New_tablespace_nameContext) {}

// ExitNew_tablespace_name is called when production new_tablespace_name is exited.
func (s *BasePlSqlParserListener) ExitNew_tablespace_name(ctx *New_tablespace_nameContext) {}

// EnterCreate_tablespace is called when production create_tablespace is entered.
func (s *BasePlSqlParserListener) EnterCreate_tablespace(ctx *Create_tablespaceContext) {}

// ExitCreate_tablespace is called when production create_tablespace is exited.
func (s *BasePlSqlParserListener) ExitCreate_tablespace(ctx *Create_tablespaceContext) {}

// EnterPermanent_tablespace_clause is called when production permanent_tablespace_clause is entered.
func (s *BasePlSqlParserListener) EnterPermanent_tablespace_clause(ctx *Permanent_tablespace_clauseContext) {
}

// ExitPermanent_tablespace_clause is called when production permanent_tablespace_clause is exited.
func (s *BasePlSqlParserListener) ExitPermanent_tablespace_clause(ctx *Permanent_tablespace_clauseContext) {
}

// EnterTablespace_encryption_spec is called when production tablespace_encryption_spec is entered.
func (s *BasePlSqlParserListener) EnterTablespace_encryption_spec(ctx *Tablespace_encryption_specContext) {
}

// ExitTablespace_encryption_spec is called when production tablespace_encryption_spec is exited.
func (s *BasePlSqlParserListener) ExitTablespace_encryption_spec(ctx *Tablespace_encryption_specContext) {
}

// EnterLogging_clause is called when production logging_clause is entered.
func (s *BasePlSqlParserListener) EnterLogging_clause(ctx *Logging_clauseContext) {}

// ExitLogging_clause is called when production logging_clause is exited.
func (s *BasePlSqlParserListener) ExitLogging_clause(ctx *Logging_clauseContext) {}

// EnterExtent_management_clause is called when production extent_management_clause is entered.
func (s *BasePlSqlParserListener) EnterExtent_management_clause(ctx *Extent_management_clauseContext) {
}

// ExitExtent_management_clause is called when production extent_management_clause is exited.
func (s *BasePlSqlParserListener) ExitExtent_management_clause(ctx *Extent_management_clauseContext) {
}

// EnterSegment_management_clause is called when production segment_management_clause is entered.
func (s *BasePlSqlParserListener) EnterSegment_management_clause(ctx *Segment_management_clauseContext) {
}

// ExitSegment_management_clause is called when production segment_management_clause is exited.
func (s *BasePlSqlParserListener) ExitSegment_management_clause(ctx *Segment_management_clauseContext) {
}

// EnterTemporary_tablespace_clause is called when production temporary_tablespace_clause is entered.
func (s *BasePlSqlParserListener) EnterTemporary_tablespace_clause(ctx *Temporary_tablespace_clauseContext) {
}

// ExitTemporary_tablespace_clause is called when production temporary_tablespace_clause is exited.
func (s *BasePlSqlParserListener) ExitTemporary_tablespace_clause(ctx *Temporary_tablespace_clauseContext) {
}

// EnterUndo_tablespace_clause is called when production undo_tablespace_clause is entered.
func (s *BasePlSqlParserListener) EnterUndo_tablespace_clause(ctx *Undo_tablespace_clauseContext) {}

// ExitUndo_tablespace_clause is called when production undo_tablespace_clause is exited.
func (s *BasePlSqlParserListener) ExitUndo_tablespace_clause(ctx *Undo_tablespace_clauseContext) {}

// EnterTablespace_retention_clause is called when production tablespace_retention_clause is entered.
func (s *BasePlSqlParserListener) EnterTablespace_retention_clause(ctx *Tablespace_retention_clauseContext) {
}

// ExitTablespace_retention_clause is called when production tablespace_retention_clause is exited.
func (s *BasePlSqlParserListener) ExitTablespace_retention_clause(ctx *Tablespace_retention_clauseContext) {
}

// EnterCreate_tablespace_set is called when production create_tablespace_set is entered.
func (s *BasePlSqlParserListener) EnterCreate_tablespace_set(ctx *Create_tablespace_setContext) {}

// ExitCreate_tablespace_set is called when production create_tablespace_set is exited.
func (s *BasePlSqlParserListener) ExitCreate_tablespace_set(ctx *Create_tablespace_setContext) {}

// EnterPermanent_tablespace_attrs is called when production permanent_tablespace_attrs is entered.
func (s *BasePlSqlParserListener) EnterPermanent_tablespace_attrs(ctx *Permanent_tablespace_attrsContext) {
}

// ExitPermanent_tablespace_attrs is called when production permanent_tablespace_attrs is exited.
func (s *BasePlSqlParserListener) ExitPermanent_tablespace_attrs(ctx *Permanent_tablespace_attrsContext) {
}

// EnterTablespace_encryption_clause is called when production tablespace_encryption_clause is entered.
func (s *BasePlSqlParserListener) EnterTablespace_encryption_clause(ctx *Tablespace_encryption_clauseContext) {
}

// ExitTablespace_encryption_clause is called when production tablespace_encryption_clause is exited.
func (s *BasePlSqlParserListener) ExitTablespace_encryption_clause(ctx *Tablespace_encryption_clauseContext) {
}

// EnterDefault_tablespace_params is called when production default_tablespace_params is entered.
func (s *BasePlSqlParserListener) EnterDefault_tablespace_params(ctx *Default_tablespace_paramsContext) {
}

// ExitDefault_tablespace_params is called when production default_tablespace_params is exited.
func (s *BasePlSqlParserListener) ExitDefault_tablespace_params(ctx *Default_tablespace_paramsContext) {
}

// EnterDefault_table_compression is called when production default_table_compression is entered.
func (s *BasePlSqlParserListener) EnterDefault_table_compression(ctx *Default_table_compressionContext) {
}

// ExitDefault_table_compression is called when production default_table_compression is exited.
func (s *BasePlSqlParserListener) ExitDefault_table_compression(ctx *Default_table_compressionContext) {
}

// EnterLow_high is called when production low_high is entered.
func (s *BasePlSqlParserListener) EnterLow_high(ctx *Low_highContext) {}

// ExitLow_high is called when production low_high is exited.
func (s *BasePlSqlParserListener) ExitLow_high(ctx *Low_highContext) {}

// EnterDefault_index_compression is called when production default_index_compression is entered.
func (s *BasePlSqlParserListener) EnterDefault_index_compression(ctx *Default_index_compressionContext) {
}

// ExitDefault_index_compression is called when production default_index_compression is exited.
func (s *BasePlSqlParserListener) ExitDefault_index_compression(ctx *Default_index_compressionContext) {
}

// EnterInmmemory_clause is called when production inmmemory_clause is entered.
func (s *BasePlSqlParserListener) EnterInmmemory_clause(ctx *Inmmemory_clauseContext) {}

// ExitInmmemory_clause is called when production inmmemory_clause is exited.
func (s *BasePlSqlParserListener) ExitInmmemory_clause(ctx *Inmmemory_clauseContext) {}

// EnterDatafile_specification is called when production datafile_specification is entered.
func (s *BasePlSqlParserListener) EnterDatafile_specification(ctx *Datafile_specificationContext) {}

// ExitDatafile_specification is called when production datafile_specification is exited.
func (s *BasePlSqlParserListener) ExitDatafile_specification(ctx *Datafile_specificationContext) {}

// EnterTempfile_specification is called when production tempfile_specification is entered.
func (s *BasePlSqlParserListener) EnterTempfile_specification(ctx *Tempfile_specificationContext) {}

// ExitTempfile_specification is called when production tempfile_specification is exited.
func (s *BasePlSqlParserListener) ExitTempfile_specification(ctx *Tempfile_specificationContext) {}

// EnterDatafile_tempfile_spec is called when production datafile_tempfile_spec is entered.
func (s *BasePlSqlParserListener) EnterDatafile_tempfile_spec(ctx *Datafile_tempfile_specContext) {}

// ExitDatafile_tempfile_spec is called when production datafile_tempfile_spec is exited.
func (s *BasePlSqlParserListener) ExitDatafile_tempfile_spec(ctx *Datafile_tempfile_specContext) {}

// EnterRedo_log_file_spec is called when production redo_log_file_spec is entered.
func (s *BasePlSqlParserListener) EnterRedo_log_file_spec(ctx *Redo_log_file_specContext) {}

// ExitRedo_log_file_spec is called when production redo_log_file_spec is exited.
func (s *BasePlSqlParserListener) ExitRedo_log_file_spec(ctx *Redo_log_file_specContext) {}

// EnterAutoextend_clause is called when production autoextend_clause is entered.
func (s *BasePlSqlParserListener) EnterAutoextend_clause(ctx *Autoextend_clauseContext) {}

// ExitAutoextend_clause is called when production autoextend_clause is exited.
func (s *BasePlSqlParserListener) ExitAutoextend_clause(ctx *Autoextend_clauseContext) {}

// EnterMaxsize_clause is called when production maxsize_clause is entered.
func (s *BasePlSqlParserListener) EnterMaxsize_clause(ctx *Maxsize_clauseContext) {}

// ExitMaxsize_clause is called when production maxsize_clause is exited.
func (s *BasePlSqlParserListener) ExitMaxsize_clause(ctx *Maxsize_clauseContext) {}

// EnterBuild_clause is called when production build_clause is entered.
func (s *BasePlSqlParserListener) EnterBuild_clause(ctx *Build_clauseContext) {}

// ExitBuild_clause is called when production build_clause is exited.
func (s *BasePlSqlParserListener) ExitBuild_clause(ctx *Build_clauseContext) {}

// EnterParallel_clause is called when production parallel_clause is entered.
func (s *BasePlSqlParserListener) EnterParallel_clause(ctx *Parallel_clauseContext) {}

// ExitParallel_clause is called when production parallel_clause is exited.
func (s *BasePlSqlParserListener) ExitParallel_clause(ctx *Parallel_clauseContext) {}

// EnterParallel_instances_clause is called when production parallel_instances_clause is entered.
func (s *BasePlSqlParserListener) EnterParallel_instances_clause(ctx *Parallel_instances_clauseContext) {
}

// ExitParallel_instances_clause is called when production parallel_instances_clause is exited.
func (s *BasePlSqlParserListener) ExitParallel_instances_clause(ctx *Parallel_instances_clauseContext) {
}

// EnterAlter_materialized_view is called when production alter_materialized_view is entered.
func (s *BasePlSqlParserListener) EnterAlter_materialized_view(ctx *Alter_materialized_viewContext) {}

// ExitAlter_materialized_view is called when production alter_materialized_view is exited.
func (s *BasePlSqlParserListener) ExitAlter_materialized_view(ctx *Alter_materialized_viewContext) {}

// EnterAlter_mv_option1 is called when production alter_mv_option1 is entered.
func (s *BasePlSqlParserListener) EnterAlter_mv_option1(ctx *Alter_mv_option1Context) {}

// ExitAlter_mv_option1 is called when production alter_mv_option1 is exited.
func (s *BasePlSqlParserListener) ExitAlter_mv_option1(ctx *Alter_mv_option1Context) {}

// EnterAlter_mv_refresh is called when production alter_mv_refresh is entered.
func (s *BasePlSqlParserListener) EnterAlter_mv_refresh(ctx *Alter_mv_refreshContext) {}

// ExitAlter_mv_refresh is called when production alter_mv_refresh is exited.
func (s *BasePlSqlParserListener) ExitAlter_mv_refresh(ctx *Alter_mv_refreshContext) {}

// EnterRollback_segment is called when production rollback_segment is entered.
func (s *BasePlSqlParserListener) EnterRollback_segment(ctx *Rollback_segmentContext) {}

// ExitRollback_segment is called when production rollback_segment is exited.
func (s *BasePlSqlParserListener) ExitRollback_segment(ctx *Rollback_segmentContext) {}

// EnterModify_mv_column_clause is called when production modify_mv_column_clause is entered.
func (s *BasePlSqlParserListener) EnterModify_mv_column_clause(ctx *Modify_mv_column_clauseContext) {}

// ExitModify_mv_column_clause is called when production modify_mv_column_clause is exited.
func (s *BasePlSqlParserListener) ExitModify_mv_column_clause(ctx *Modify_mv_column_clauseContext) {}

// EnterAlter_materialized_view_log is called when production alter_materialized_view_log is entered.
func (s *BasePlSqlParserListener) EnterAlter_materialized_view_log(ctx *Alter_materialized_view_logContext) {
}

// ExitAlter_materialized_view_log is called when production alter_materialized_view_log is exited.
func (s *BasePlSqlParserListener) ExitAlter_materialized_view_log(ctx *Alter_materialized_view_logContext) {
}

// EnterAdd_mv_log_column_clause is called when production add_mv_log_column_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_mv_log_column_clause(ctx *Add_mv_log_column_clauseContext) {
}

// ExitAdd_mv_log_column_clause is called when production add_mv_log_column_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_mv_log_column_clause(ctx *Add_mv_log_column_clauseContext) {
}

// EnterMove_mv_log_clause is called when production move_mv_log_clause is entered.
func (s *BasePlSqlParserListener) EnterMove_mv_log_clause(ctx *Move_mv_log_clauseContext) {}

// ExitMove_mv_log_clause is called when production move_mv_log_clause is exited.
func (s *BasePlSqlParserListener) ExitMove_mv_log_clause(ctx *Move_mv_log_clauseContext) {}

// EnterMv_log_augmentation is called when production mv_log_augmentation is entered.
func (s *BasePlSqlParserListener) EnterMv_log_augmentation(ctx *Mv_log_augmentationContext) {}

// ExitMv_log_augmentation is called when production mv_log_augmentation is exited.
func (s *BasePlSqlParserListener) ExitMv_log_augmentation(ctx *Mv_log_augmentationContext) {}

// EnterCreate_materialized_view_log is called when production create_materialized_view_log is entered.
func (s *BasePlSqlParserListener) EnterCreate_materialized_view_log(ctx *Create_materialized_view_logContext) {
}

// ExitCreate_materialized_view_log is called when production create_materialized_view_log is exited.
func (s *BasePlSqlParserListener) ExitCreate_materialized_view_log(ctx *Create_materialized_view_logContext) {
}

// EnterNew_values_clause is called when production new_values_clause is entered.
func (s *BasePlSqlParserListener) EnterNew_values_clause(ctx *New_values_clauseContext) {}

// ExitNew_values_clause is called when production new_values_clause is exited.
func (s *BasePlSqlParserListener) ExitNew_values_clause(ctx *New_values_clauseContext) {}

// EnterMv_log_purge_clause is called when production mv_log_purge_clause is entered.
func (s *BasePlSqlParserListener) EnterMv_log_purge_clause(ctx *Mv_log_purge_clauseContext) {}

// ExitMv_log_purge_clause is called when production mv_log_purge_clause is exited.
func (s *BasePlSqlParserListener) ExitMv_log_purge_clause(ctx *Mv_log_purge_clauseContext) {}

// EnterCreate_materialized_zonemap is called when production create_materialized_zonemap is entered.
func (s *BasePlSqlParserListener) EnterCreate_materialized_zonemap(ctx *Create_materialized_zonemapContext) {
}

// ExitCreate_materialized_zonemap is called when production create_materialized_zonemap is exited.
func (s *BasePlSqlParserListener) ExitCreate_materialized_zonemap(ctx *Create_materialized_zonemapContext) {
}

// EnterAlter_materialized_zonemap is called when production alter_materialized_zonemap is entered.
func (s *BasePlSqlParserListener) EnterAlter_materialized_zonemap(ctx *Alter_materialized_zonemapContext) {
}

// ExitAlter_materialized_zonemap is called when production alter_materialized_zonemap is exited.
func (s *BasePlSqlParserListener) ExitAlter_materialized_zonemap(ctx *Alter_materialized_zonemapContext) {
}

// EnterDrop_materialized_zonemap is called when production drop_materialized_zonemap is entered.
func (s *BasePlSqlParserListener) EnterDrop_materialized_zonemap(ctx *Drop_materialized_zonemapContext) {
}

// ExitDrop_materialized_zonemap is called when production drop_materialized_zonemap is exited.
func (s *BasePlSqlParserListener) ExitDrop_materialized_zonemap(ctx *Drop_materialized_zonemapContext) {
}

// EnterZonemap_refresh_clause is called when production zonemap_refresh_clause is entered.
func (s *BasePlSqlParserListener) EnterZonemap_refresh_clause(ctx *Zonemap_refresh_clauseContext) {}

// ExitZonemap_refresh_clause is called when production zonemap_refresh_clause is exited.
func (s *BasePlSqlParserListener) ExitZonemap_refresh_clause(ctx *Zonemap_refresh_clauseContext) {}

// EnterZonemap_attributes is called when production zonemap_attributes is entered.
func (s *BasePlSqlParserListener) EnterZonemap_attributes(ctx *Zonemap_attributesContext) {}

// ExitZonemap_attributes is called when production zonemap_attributes is exited.
func (s *BasePlSqlParserListener) ExitZonemap_attributes(ctx *Zonemap_attributesContext) {}

// EnterZonemap_name is called when production zonemap_name is entered.
func (s *BasePlSqlParserListener) EnterZonemap_name(ctx *Zonemap_nameContext) {}

// ExitZonemap_name is called when production zonemap_name is exited.
func (s *BasePlSqlParserListener) ExitZonemap_name(ctx *Zonemap_nameContext) {}

// EnterOperator_name is called when production operator_name is entered.
func (s *BasePlSqlParserListener) EnterOperator_name(ctx *Operator_nameContext) {}

// ExitOperator_name is called when production operator_name is exited.
func (s *BasePlSqlParserListener) ExitOperator_name(ctx *Operator_nameContext) {}

// EnterOperator_function_name is called when production operator_function_name is entered.
func (s *BasePlSqlParserListener) EnterOperator_function_name(ctx *Operator_function_nameContext) {}

// ExitOperator_function_name is called when production operator_function_name is exited.
func (s *BasePlSqlParserListener) ExitOperator_function_name(ctx *Operator_function_nameContext) {}

// EnterCreate_zonemap_on_table is called when production create_zonemap_on_table is entered.
func (s *BasePlSqlParserListener) EnterCreate_zonemap_on_table(ctx *Create_zonemap_on_tableContext) {}

// ExitCreate_zonemap_on_table is called when production create_zonemap_on_table is exited.
func (s *BasePlSqlParserListener) ExitCreate_zonemap_on_table(ctx *Create_zonemap_on_tableContext) {}

// EnterCreate_zonemap_as_subquery is called when production create_zonemap_as_subquery is entered.
func (s *BasePlSqlParserListener) EnterCreate_zonemap_as_subquery(ctx *Create_zonemap_as_subqueryContext) {
}

// ExitCreate_zonemap_as_subquery is called when production create_zonemap_as_subquery is exited.
func (s *BasePlSqlParserListener) ExitCreate_zonemap_as_subquery(ctx *Create_zonemap_as_subqueryContext) {
}

// EnterAlter_operator is called when production alter_operator is entered.
func (s *BasePlSqlParserListener) EnterAlter_operator(ctx *Alter_operatorContext) {}

// ExitAlter_operator is called when production alter_operator is exited.
func (s *BasePlSqlParserListener) ExitAlter_operator(ctx *Alter_operatorContext) {}

// EnterDrop_operator is called when production drop_operator is entered.
func (s *BasePlSqlParserListener) EnterDrop_operator(ctx *Drop_operatorContext) {}

// ExitDrop_operator is called when production drop_operator is exited.
func (s *BasePlSqlParserListener) ExitDrop_operator(ctx *Drop_operatorContext) {}

// EnterCreate_operator is called when production create_operator is entered.
func (s *BasePlSqlParserListener) EnterCreate_operator(ctx *Create_operatorContext) {}

// ExitCreate_operator is called when production create_operator is exited.
func (s *BasePlSqlParserListener) ExitCreate_operator(ctx *Create_operatorContext) {}

// EnterBinding_clause is called when production binding_clause is entered.
func (s *BasePlSqlParserListener) EnterBinding_clause(ctx *Binding_clauseContext) {}

// ExitBinding_clause is called when production binding_clause is exited.
func (s *BasePlSqlParserListener) ExitBinding_clause(ctx *Binding_clauseContext) {}

// EnterAdd_binding_clause is called when production add_binding_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_binding_clause(ctx *Add_binding_clauseContext) {}

// ExitAdd_binding_clause is called when production add_binding_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_binding_clause(ctx *Add_binding_clauseContext) {}

// EnterImplementation_clause is called when production implementation_clause is entered.
func (s *BasePlSqlParserListener) EnterImplementation_clause(ctx *Implementation_clauseContext) {}

// ExitImplementation_clause is called when production implementation_clause is exited.
func (s *BasePlSqlParserListener) ExitImplementation_clause(ctx *Implementation_clauseContext) {}

// EnterPrimary_operator_list is called when production primary_operator_list is entered.
func (s *BasePlSqlParserListener) EnterPrimary_operator_list(ctx *Primary_operator_listContext) {}

// ExitPrimary_operator_list is called when production primary_operator_list is exited.
func (s *BasePlSqlParserListener) ExitPrimary_operator_list(ctx *Primary_operator_listContext) {}

// EnterPrimary_operator_item is called when production primary_operator_item is entered.
func (s *BasePlSqlParserListener) EnterPrimary_operator_item(ctx *Primary_operator_itemContext) {}

// ExitPrimary_operator_item is called when production primary_operator_item is exited.
func (s *BasePlSqlParserListener) ExitPrimary_operator_item(ctx *Primary_operator_itemContext) {}

// EnterOperator_context_clause is called when production operator_context_clause is entered.
func (s *BasePlSqlParserListener) EnterOperator_context_clause(ctx *Operator_context_clauseContext) {}

// ExitOperator_context_clause is called when production operator_context_clause is exited.
func (s *BasePlSqlParserListener) ExitOperator_context_clause(ctx *Operator_context_clauseContext) {}

// EnterUsing_function_clause is called when production using_function_clause is entered.
func (s *BasePlSqlParserListener) EnterUsing_function_clause(ctx *Using_function_clauseContext) {}

// ExitUsing_function_clause is called when production using_function_clause is exited.
func (s *BasePlSqlParserListener) ExitUsing_function_clause(ctx *Using_function_clauseContext) {}

// EnterDrop_binding_clause is called when production drop_binding_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_binding_clause(ctx *Drop_binding_clauseContext) {}

// ExitDrop_binding_clause is called when production drop_binding_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_binding_clause(ctx *Drop_binding_clauseContext) {}

// EnterCreate_materialized_view is called when production create_materialized_view is entered.
func (s *BasePlSqlParserListener) EnterCreate_materialized_view(ctx *Create_materialized_viewContext) {
}

// ExitCreate_materialized_view is called when production create_materialized_view is exited.
func (s *BasePlSqlParserListener) ExitCreate_materialized_view(ctx *Create_materialized_viewContext) {
}

// EnterScoped_table_ref_constraint is called when production scoped_table_ref_constraint is entered.
func (s *BasePlSqlParserListener) EnterScoped_table_ref_constraint(ctx *Scoped_table_ref_constraintContext) {
}

// ExitScoped_table_ref_constraint is called when production scoped_table_ref_constraint is exited.
func (s *BasePlSqlParserListener) ExitScoped_table_ref_constraint(ctx *Scoped_table_ref_constraintContext) {
}

// EnterMv_column_alias is called when production mv_column_alias is entered.
func (s *BasePlSqlParserListener) EnterMv_column_alias(ctx *Mv_column_aliasContext) {}

// ExitMv_column_alias is called when production mv_column_alias is exited.
func (s *BasePlSqlParserListener) ExitMv_column_alias(ctx *Mv_column_aliasContext) {}

// EnterCreate_mv_refresh is called when production create_mv_refresh is entered.
func (s *BasePlSqlParserListener) EnterCreate_mv_refresh(ctx *Create_mv_refreshContext) {}

// ExitCreate_mv_refresh is called when production create_mv_refresh is exited.
func (s *BasePlSqlParserListener) ExitCreate_mv_refresh(ctx *Create_mv_refreshContext) {}

// EnterQuery_rewrite_clause is called when production query_rewrite_clause is entered.
func (s *BasePlSqlParserListener) EnterQuery_rewrite_clause(ctx *Query_rewrite_clauseContext) {}

// ExitQuery_rewrite_clause is called when production query_rewrite_clause is exited.
func (s *BasePlSqlParserListener) ExitQuery_rewrite_clause(ctx *Query_rewrite_clauseContext) {}

// EnterUnusable_editions_clause is called when production unusable_editions_clause is entered.
func (s *BasePlSqlParserListener) EnterUnusable_editions_clause(ctx *Unusable_editions_clauseContext) {
}

// ExitUnusable_editions_clause is called when production unusable_editions_clause is exited.
func (s *BasePlSqlParserListener) ExitUnusable_editions_clause(ctx *Unusable_editions_clauseContext) {
}

// EnterDrop_materialized_view is called when production drop_materialized_view is entered.
func (s *BasePlSqlParserListener) EnterDrop_materialized_view(ctx *Drop_materialized_viewContext) {}

// ExitDrop_materialized_view is called when production drop_materialized_view is exited.
func (s *BasePlSqlParserListener) ExitDrop_materialized_view(ctx *Drop_materialized_viewContext) {}

// EnterDrop_materialized_view_log is called when production drop_materialized_view_log is entered.
func (s *BasePlSqlParserListener) EnterDrop_materialized_view_log(ctx *Drop_materialized_view_logContext) {
}

// ExitDrop_materialized_view_log is called when production drop_materialized_view_log is exited.
func (s *BasePlSqlParserListener) ExitDrop_materialized_view_log(ctx *Drop_materialized_view_logContext) {
}

// EnterCreate_context is called when production create_context is entered.
func (s *BasePlSqlParserListener) EnterCreate_context(ctx *Create_contextContext) {}

// ExitCreate_context is called when production create_context is exited.
func (s *BasePlSqlParserListener) ExitCreate_context(ctx *Create_contextContext) {}

// EnterOracle_namespace is called when production oracle_namespace is entered.
func (s *BasePlSqlParserListener) EnterOracle_namespace(ctx *Oracle_namespaceContext) {}

// ExitOracle_namespace is called when production oracle_namespace is exited.
func (s *BasePlSqlParserListener) ExitOracle_namespace(ctx *Oracle_namespaceContext) {}

// EnterCreate_cluster is called when production create_cluster is entered.
func (s *BasePlSqlParserListener) EnterCreate_cluster(ctx *Create_clusterContext) {}

// ExitCreate_cluster is called when production create_cluster is exited.
func (s *BasePlSqlParserListener) ExitCreate_cluster(ctx *Create_clusterContext) {}

// EnterCreate_profile is called when production create_profile is entered.
func (s *BasePlSqlParserListener) EnterCreate_profile(ctx *Create_profileContext) {}

// ExitCreate_profile is called when production create_profile is exited.
func (s *BasePlSqlParserListener) ExitCreate_profile(ctx *Create_profileContext) {}

// EnterResource_parameters is called when production resource_parameters is entered.
func (s *BasePlSqlParserListener) EnterResource_parameters(ctx *Resource_parametersContext) {}

// ExitResource_parameters is called when production resource_parameters is exited.
func (s *BasePlSqlParserListener) ExitResource_parameters(ctx *Resource_parametersContext) {}

// EnterPassword_parameters is called when production password_parameters is entered.
func (s *BasePlSqlParserListener) EnterPassword_parameters(ctx *Password_parametersContext) {}

// ExitPassword_parameters is called when production password_parameters is exited.
func (s *BasePlSqlParserListener) ExitPassword_parameters(ctx *Password_parametersContext) {}

// EnterCreate_lockdown_profile is called when production create_lockdown_profile is entered.
func (s *BasePlSqlParserListener) EnterCreate_lockdown_profile(ctx *Create_lockdown_profileContext) {}

// ExitCreate_lockdown_profile is called when production create_lockdown_profile is exited.
func (s *BasePlSqlParserListener) ExitCreate_lockdown_profile(ctx *Create_lockdown_profileContext) {}

// EnterStatic_base_profile is called when production static_base_profile is entered.
func (s *BasePlSqlParserListener) EnterStatic_base_profile(ctx *Static_base_profileContext) {}

// ExitStatic_base_profile is called when production static_base_profile is exited.
func (s *BasePlSqlParserListener) ExitStatic_base_profile(ctx *Static_base_profileContext) {}

// EnterDynamic_base_profile is called when production dynamic_base_profile is entered.
func (s *BasePlSqlParserListener) EnterDynamic_base_profile(ctx *Dynamic_base_profileContext) {}

// ExitDynamic_base_profile is called when production dynamic_base_profile is exited.
func (s *BasePlSqlParserListener) ExitDynamic_base_profile(ctx *Dynamic_base_profileContext) {}

// EnterCreate_outline is called when production create_outline is entered.
func (s *BasePlSqlParserListener) EnterCreate_outline(ctx *Create_outlineContext) {}

// ExitCreate_outline is called when production create_outline is exited.
func (s *BasePlSqlParserListener) ExitCreate_outline(ctx *Create_outlineContext) {}

// EnterCreate_restore_point is called when production create_restore_point is entered.
func (s *BasePlSqlParserListener) EnterCreate_restore_point(ctx *Create_restore_pointContext) {}

// ExitCreate_restore_point is called when production create_restore_point is exited.
func (s *BasePlSqlParserListener) ExitCreate_restore_point(ctx *Create_restore_pointContext) {}

// EnterCreate_role is called when production create_role is entered.
func (s *BasePlSqlParserListener) EnterCreate_role(ctx *Create_roleContext) {}

// ExitCreate_role is called when production create_role is exited.
func (s *BasePlSqlParserListener) ExitCreate_role(ctx *Create_roleContext) {}

// EnterCreate_table is called when production create_table is entered.
func (s *BasePlSqlParserListener) EnterCreate_table(ctx *Create_tableContext) {}

// ExitCreate_table is called when production create_table is exited.
func (s *BasePlSqlParserListener) ExitCreate_table(ctx *Create_tableContext) {}

// EnterXmltype_table is called when production xmltype_table is entered.
func (s *BasePlSqlParserListener) EnterXmltype_table(ctx *Xmltype_tableContext) {}

// ExitXmltype_table is called when production xmltype_table is exited.
func (s *BasePlSqlParserListener) ExitXmltype_table(ctx *Xmltype_tableContext) {}

// EnterXmltype_virtual_columns is called when production xmltype_virtual_columns is entered.
func (s *BasePlSqlParserListener) EnterXmltype_virtual_columns(ctx *Xmltype_virtual_columnsContext) {}

// ExitXmltype_virtual_columns is called when production xmltype_virtual_columns is exited.
func (s *BasePlSqlParserListener) ExitXmltype_virtual_columns(ctx *Xmltype_virtual_columnsContext) {}

// EnterXmltype_column_properties is called when production xmltype_column_properties is entered.
func (s *BasePlSqlParserListener) EnterXmltype_column_properties(ctx *Xmltype_column_propertiesContext) {
}

// ExitXmltype_column_properties is called when production xmltype_column_properties is exited.
func (s *BasePlSqlParserListener) ExitXmltype_column_properties(ctx *Xmltype_column_propertiesContext) {
}

// EnterXmltype_storage is called when production xmltype_storage is entered.
func (s *BasePlSqlParserListener) EnterXmltype_storage(ctx *Xmltype_storageContext) {}

// ExitXmltype_storage is called when production xmltype_storage is exited.
func (s *BasePlSqlParserListener) ExitXmltype_storage(ctx *Xmltype_storageContext) {}

// EnterXmlschema_spec is called when production xmlschema_spec is entered.
func (s *BasePlSqlParserListener) EnterXmlschema_spec(ctx *Xmlschema_specContext) {}

// ExitXmlschema_spec is called when production xmlschema_spec is exited.
func (s *BasePlSqlParserListener) ExitXmlschema_spec(ctx *Xmlschema_specContext) {}

// EnterObject_table is called when production object_table is entered.
func (s *BasePlSqlParserListener) EnterObject_table(ctx *Object_tableContext) {}

// ExitObject_table is called when production object_table is exited.
func (s *BasePlSqlParserListener) ExitObject_table(ctx *Object_tableContext) {}

// EnterObject_type is called when production object_type is entered.
func (s *BasePlSqlParserListener) EnterObject_type(ctx *Object_typeContext) {}

// ExitObject_type is called when production object_type is exited.
func (s *BasePlSqlParserListener) ExitObject_type(ctx *Object_typeContext) {}

// EnterOid_index_clause is called when production oid_index_clause is entered.
func (s *BasePlSqlParserListener) EnterOid_index_clause(ctx *Oid_index_clauseContext) {}

// ExitOid_index_clause is called when production oid_index_clause is exited.
func (s *BasePlSqlParserListener) ExitOid_index_clause(ctx *Oid_index_clauseContext) {}

// EnterOid_clause is called when production oid_clause is entered.
func (s *BasePlSqlParserListener) EnterOid_clause(ctx *Oid_clauseContext) {}

// ExitOid_clause is called when production oid_clause is exited.
func (s *BasePlSqlParserListener) ExitOid_clause(ctx *Oid_clauseContext) {}

// EnterObject_properties is called when production object_properties is entered.
func (s *BasePlSqlParserListener) EnterObject_properties(ctx *Object_propertiesContext) {}

// ExitObject_properties is called when production object_properties is exited.
func (s *BasePlSqlParserListener) ExitObject_properties(ctx *Object_propertiesContext) {}

// EnterObject_table_substitution is called when production object_table_substitution is entered.
func (s *BasePlSqlParserListener) EnterObject_table_substitution(ctx *Object_table_substitutionContext) {
}

// ExitObject_table_substitution is called when production object_table_substitution is exited.
func (s *BasePlSqlParserListener) ExitObject_table_substitution(ctx *Object_table_substitutionContext) {
}

// EnterRelational_table is called when production relational_table is entered.
func (s *BasePlSqlParserListener) EnterRelational_table(ctx *Relational_tableContext) {}

// ExitRelational_table is called when production relational_table is exited.
func (s *BasePlSqlParserListener) ExitRelational_table(ctx *Relational_tableContext) {}

// EnterRelational_table_properties is called when production relational_table_properties is entered.
func (s *BasePlSqlParserListener) EnterRelational_table_properties(ctx *Relational_table_propertiesContext) {
}

// ExitRelational_table_properties is called when production relational_table_properties is exited.
func (s *BasePlSqlParserListener) ExitRelational_table_properties(ctx *Relational_table_propertiesContext) {
}

// EnterRelational_table_property is called when production relational_table_property is entered.
func (s *BasePlSqlParserListener) EnterRelational_table_property(ctx *Relational_table_propertyContext) {
}

// ExitRelational_table_property is called when production relational_table_property is exited.
func (s *BasePlSqlParserListener) ExitRelational_table_property(ctx *Relational_table_propertyContext) {
}

// EnterImmutable_table_clauses is called when production immutable_table_clauses is entered.
func (s *BasePlSqlParserListener) EnterImmutable_table_clauses(ctx *Immutable_table_clausesContext) {}

// ExitImmutable_table_clauses is called when production immutable_table_clauses is exited.
func (s *BasePlSqlParserListener) ExitImmutable_table_clauses(ctx *Immutable_table_clausesContext) {}

// EnterImmutable_table_no_drop_clause is called when production immutable_table_no_drop_clause is entered.
func (s *BasePlSqlParserListener) EnterImmutable_table_no_drop_clause(ctx *Immutable_table_no_drop_clauseContext) {
}

// ExitImmutable_table_no_drop_clause is called when production immutable_table_no_drop_clause is exited.
func (s *BasePlSqlParserListener) ExitImmutable_table_no_drop_clause(ctx *Immutable_table_no_drop_clauseContext) {
}

// EnterImmutable_table_no_delete_clause is called when production immutable_table_no_delete_clause is entered.
func (s *BasePlSqlParserListener) EnterImmutable_table_no_delete_clause(ctx *Immutable_table_no_delete_clauseContext) {
}

// ExitImmutable_table_no_delete_clause is called when production immutable_table_no_delete_clause is exited.
func (s *BasePlSqlParserListener) ExitImmutable_table_no_delete_clause(ctx *Immutable_table_no_delete_clauseContext) {
}

// EnterBlockchain_table_clauses is called when production blockchain_table_clauses is entered.
func (s *BasePlSqlParserListener) EnterBlockchain_table_clauses(ctx *Blockchain_table_clausesContext) {
}

// ExitBlockchain_table_clauses is called when production blockchain_table_clauses is exited.
func (s *BasePlSqlParserListener) ExitBlockchain_table_clauses(ctx *Blockchain_table_clausesContext) {
}

// EnterBlockchain_drop_table_clause is called when production blockchain_drop_table_clause is entered.
func (s *BasePlSqlParserListener) EnterBlockchain_drop_table_clause(ctx *Blockchain_drop_table_clauseContext) {
}

// ExitBlockchain_drop_table_clause is called when production blockchain_drop_table_clause is exited.
func (s *BasePlSqlParserListener) ExitBlockchain_drop_table_clause(ctx *Blockchain_drop_table_clauseContext) {
}

// EnterBlockchain_row_retention_clause is called when production blockchain_row_retention_clause is entered.
func (s *BasePlSqlParserListener) EnterBlockchain_row_retention_clause(ctx *Blockchain_row_retention_clauseContext) {
}

// ExitBlockchain_row_retention_clause is called when production blockchain_row_retention_clause is exited.
func (s *BasePlSqlParserListener) ExitBlockchain_row_retention_clause(ctx *Blockchain_row_retention_clauseContext) {
}

// EnterBlockchain_hash_and_data_format_clause is called when production blockchain_hash_and_data_format_clause is entered.
func (s *BasePlSqlParserListener) EnterBlockchain_hash_and_data_format_clause(ctx *Blockchain_hash_and_data_format_clauseContext) {
}

// ExitBlockchain_hash_and_data_format_clause is called when production blockchain_hash_and_data_format_clause is exited.
func (s *BasePlSqlParserListener) ExitBlockchain_hash_and_data_format_clause(ctx *Blockchain_hash_and_data_format_clauseContext) {
}

// EnterCollation_name is called when production collation_name is entered.
func (s *BasePlSqlParserListener) EnterCollation_name(ctx *Collation_nameContext) {}

// ExitCollation_name is called when production collation_name is exited.
func (s *BasePlSqlParserListener) ExitCollation_name(ctx *Collation_nameContext) {}

// EnterTable_properties is called when production table_properties is entered.
func (s *BasePlSqlParserListener) EnterTable_properties(ctx *Table_propertiesContext) {}

// ExitTable_properties is called when production table_properties is exited.
func (s *BasePlSqlParserListener) ExitTable_properties(ctx *Table_propertiesContext) {}

// EnterRead_only_clause is called when production read_only_clause is entered.
func (s *BasePlSqlParserListener) EnterRead_only_clause(ctx *Read_only_clauseContext) {}

// ExitRead_only_clause is called when production read_only_clause is exited.
func (s *BasePlSqlParserListener) ExitRead_only_clause(ctx *Read_only_clauseContext) {}

// EnterIndexing_clause is called when production indexing_clause is entered.
func (s *BasePlSqlParserListener) EnterIndexing_clause(ctx *Indexing_clauseContext) {}

// ExitIndexing_clause is called when production indexing_clause is exited.
func (s *BasePlSqlParserListener) ExitIndexing_clause(ctx *Indexing_clauseContext) {}

// EnterAttribute_clustering_clause is called when production attribute_clustering_clause is entered.
func (s *BasePlSqlParserListener) EnterAttribute_clustering_clause(ctx *Attribute_clustering_clauseContext) {
}

// ExitAttribute_clustering_clause is called when production attribute_clustering_clause is exited.
func (s *BasePlSqlParserListener) ExitAttribute_clustering_clause(ctx *Attribute_clustering_clauseContext) {
}

// EnterClustering_join is called when production clustering_join is entered.
func (s *BasePlSqlParserListener) EnterClustering_join(ctx *Clustering_joinContext) {}

// ExitClustering_join is called when production clustering_join is exited.
func (s *BasePlSqlParserListener) ExitClustering_join(ctx *Clustering_joinContext) {}

// EnterClustering_join_item is called when production clustering_join_item is entered.
func (s *BasePlSqlParserListener) EnterClustering_join_item(ctx *Clustering_join_itemContext) {}

// ExitClustering_join_item is called when production clustering_join_item is exited.
func (s *BasePlSqlParserListener) ExitClustering_join_item(ctx *Clustering_join_itemContext) {}

// EnterEquijoin_condition is called when production equijoin_condition is entered.
func (s *BasePlSqlParserListener) EnterEquijoin_condition(ctx *Equijoin_conditionContext) {}

// ExitEquijoin_condition is called when production equijoin_condition is exited.
func (s *BasePlSqlParserListener) ExitEquijoin_condition(ctx *Equijoin_conditionContext) {}

// EnterCluster_clause is called when production cluster_clause is entered.
func (s *BasePlSqlParserListener) EnterCluster_clause(ctx *Cluster_clauseContext) {}

// ExitCluster_clause is called when production cluster_clause is exited.
func (s *BasePlSqlParserListener) ExitCluster_clause(ctx *Cluster_clauseContext) {}

// EnterClustering_columns is called when production clustering_columns is entered.
func (s *BasePlSqlParserListener) EnterClustering_columns(ctx *Clustering_columnsContext) {}

// ExitClustering_columns is called when production clustering_columns is exited.
func (s *BasePlSqlParserListener) ExitClustering_columns(ctx *Clustering_columnsContext) {}

// EnterClustering_column_group is called when production clustering_column_group is entered.
func (s *BasePlSqlParserListener) EnterClustering_column_group(ctx *Clustering_column_groupContext) {}

// ExitClustering_column_group is called when production clustering_column_group is exited.
func (s *BasePlSqlParserListener) ExitClustering_column_group(ctx *Clustering_column_groupContext) {}

// EnterYes_no is called when production yes_no is entered.
func (s *BasePlSqlParserListener) EnterYes_no(ctx *Yes_noContext) {}

// ExitYes_no is called when production yes_no is exited.
func (s *BasePlSqlParserListener) ExitYes_no(ctx *Yes_noContext) {}

// EnterZonemap_clause is called when production zonemap_clause is entered.
func (s *BasePlSqlParserListener) EnterZonemap_clause(ctx *Zonemap_clauseContext) {}

// ExitZonemap_clause is called when production zonemap_clause is exited.
func (s *BasePlSqlParserListener) ExitZonemap_clause(ctx *Zonemap_clauseContext) {}

// EnterLogical_replication_clause is called when production logical_replication_clause is entered.
func (s *BasePlSqlParserListener) EnterLogical_replication_clause(ctx *Logical_replication_clauseContext) {
}

// ExitLogical_replication_clause is called when production logical_replication_clause is exited.
func (s *BasePlSqlParserListener) ExitLogical_replication_clause(ctx *Logical_replication_clauseContext) {
}

// EnterTable_name is called when production table_name is entered.
func (s *BasePlSqlParserListener) EnterTable_name(ctx *Table_nameContext) {}

// ExitTable_name is called when production table_name is exited.
func (s *BasePlSqlParserListener) ExitTable_name(ctx *Table_nameContext) {}

// EnterRelational_property is called when production relational_property is entered.
func (s *BasePlSqlParserListener) EnterRelational_property(ctx *Relational_propertyContext) {}

// ExitRelational_property is called when production relational_property is exited.
func (s *BasePlSqlParserListener) ExitRelational_property(ctx *Relational_propertyContext) {}

// EnterTable_partitioning_clauses is called when production table_partitioning_clauses is entered.
func (s *BasePlSqlParserListener) EnterTable_partitioning_clauses(ctx *Table_partitioning_clausesContext) {
}

// ExitTable_partitioning_clauses is called when production table_partitioning_clauses is exited.
func (s *BasePlSqlParserListener) ExitTable_partitioning_clauses(ctx *Table_partitioning_clausesContext) {
}

// EnterRange_partitions is called when production range_partitions is entered.
func (s *BasePlSqlParserListener) EnterRange_partitions(ctx *Range_partitionsContext) {}

// ExitRange_partitions is called when production range_partitions is exited.
func (s *BasePlSqlParserListener) ExitRange_partitions(ctx *Range_partitionsContext) {}

// EnterList_partitions is called when production list_partitions is entered.
func (s *BasePlSqlParserListener) EnterList_partitions(ctx *List_partitionsContext) {}

// ExitList_partitions is called when production list_partitions is exited.
func (s *BasePlSqlParserListener) ExitList_partitions(ctx *List_partitionsContext) {}

// EnterHash_partitions is called when production hash_partitions is entered.
func (s *BasePlSqlParserListener) EnterHash_partitions(ctx *Hash_partitionsContext) {}

// ExitHash_partitions is called when production hash_partitions is exited.
func (s *BasePlSqlParserListener) ExitHash_partitions(ctx *Hash_partitionsContext) {}

// EnterIndividual_hash_partitions is called when production individual_hash_partitions is entered.
func (s *BasePlSqlParserListener) EnterIndividual_hash_partitions(ctx *Individual_hash_partitionsContext) {
}

// ExitIndividual_hash_partitions is called when production individual_hash_partitions is exited.
func (s *BasePlSqlParserListener) ExitIndividual_hash_partitions(ctx *Individual_hash_partitionsContext) {
}

// EnterHash_partitions_by_quantity is called when production hash_partitions_by_quantity is entered.
func (s *BasePlSqlParserListener) EnterHash_partitions_by_quantity(ctx *Hash_partitions_by_quantityContext) {
}

// ExitHash_partitions_by_quantity is called when production hash_partitions_by_quantity is exited.
func (s *BasePlSqlParserListener) ExitHash_partitions_by_quantity(ctx *Hash_partitions_by_quantityContext) {
}

// EnterHash_partition_quantity is called when production hash_partition_quantity is entered.
func (s *BasePlSqlParserListener) EnterHash_partition_quantity(ctx *Hash_partition_quantityContext) {}

// ExitHash_partition_quantity is called when production hash_partition_quantity is exited.
func (s *BasePlSqlParserListener) ExitHash_partition_quantity(ctx *Hash_partition_quantityContext) {}

// EnterComposite_range_partitions is called when production composite_range_partitions is entered.
func (s *BasePlSqlParserListener) EnterComposite_range_partitions(ctx *Composite_range_partitionsContext) {
}

// ExitComposite_range_partitions is called when production composite_range_partitions is exited.
func (s *BasePlSqlParserListener) ExitComposite_range_partitions(ctx *Composite_range_partitionsContext) {
}

// EnterComposite_list_partitions is called when production composite_list_partitions is entered.
func (s *BasePlSqlParserListener) EnterComposite_list_partitions(ctx *Composite_list_partitionsContext) {
}

// ExitComposite_list_partitions is called when production composite_list_partitions is exited.
func (s *BasePlSqlParserListener) ExitComposite_list_partitions(ctx *Composite_list_partitionsContext) {
}

// EnterComposite_hash_partitions is called when production composite_hash_partitions is entered.
func (s *BasePlSqlParserListener) EnterComposite_hash_partitions(ctx *Composite_hash_partitionsContext) {
}

// ExitComposite_hash_partitions is called when production composite_hash_partitions is exited.
func (s *BasePlSqlParserListener) ExitComposite_hash_partitions(ctx *Composite_hash_partitionsContext) {
}

// EnterReference_partitioning is called when production reference_partitioning is entered.
func (s *BasePlSqlParserListener) EnterReference_partitioning(ctx *Reference_partitioningContext) {}

// ExitReference_partitioning is called when production reference_partitioning is exited.
func (s *BasePlSqlParserListener) ExitReference_partitioning(ctx *Reference_partitioningContext) {}

// EnterReference_partition_desc is called when production reference_partition_desc is entered.
func (s *BasePlSqlParserListener) EnterReference_partition_desc(ctx *Reference_partition_descContext) {
}

// ExitReference_partition_desc is called when production reference_partition_desc is exited.
func (s *BasePlSqlParserListener) ExitReference_partition_desc(ctx *Reference_partition_descContext) {
}

// EnterSystem_partitioning is called when production system_partitioning is entered.
func (s *BasePlSqlParserListener) EnterSystem_partitioning(ctx *System_partitioningContext) {}

// ExitSystem_partitioning is called when production system_partitioning is exited.
func (s *BasePlSqlParserListener) ExitSystem_partitioning(ctx *System_partitioningContext) {}

// EnterRange_partition_desc is called when production range_partition_desc is entered.
func (s *BasePlSqlParserListener) EnterRange_partition_desc(ctx *Range_partition_descContext) {}

// ExitRange_partition_desc is called when production range_partition_desc is exited.
func (s *BasePlSqlParserListener) ExitRange_partition_desc(ctx *Range_partition_descContext) {}

// EnterList_partition_desc is called when production list_partition_desc is entered.
func (s *BasePlSqlParserListener) EnterList_partition_desc(ctx *List_partition_descContext) {}

// ExitList_partition_desc is called when production list_partition_desc is exited.
func (s *BasePlSqlParserListener) ExitList_partition_desc(ctx *List_partition_descContext) {}

// EnterSubpartition_template is called when production subpartition_template is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_template(ctx *Subpartition_templateContext) {}

// ExitSubpartition_template is called when production subpartition_template is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_template(ctx *Subpartition_templateContext) {}

// EnterHash_subpartition_quantity is called when production hash_subpartition_quantity is entered.
func (s *BasePlSqlParserListener) EnterHash_subpartition_quantity(ctx *Hash_subpartition_quantityContext) {
}

// ExitHash_subpartition_quantity is called when production hash_subpartition_quantity is exited.
func (s *BasePlSqlParserListener) ExitHash_subpartition_quantity(ctx *Hash_subpartition_quantityContext) {
}

// EnterSubpartition_by_range is called when production subpartition_by_range is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_by_range(ctx *Subpartition_by_rangeContext) {}

// ExitSubpartition_by_range is called when production subpartition_by_range is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_by_range(ctx *Subpartition_by_rangeContext) {}

// EnterSubpartition_by_list is called when production subpartition_by_list is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_by_list(ctx *Subpartition_by_listContext) {}

// ExitSubpartition_by_list is called when production subpartition_by_list is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_by_list(ctx *Subpartition_by_listContext) {}

// EnterSubpartition_by_hash is called when production subpartition_by_hash is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_by_hash(ctx *Subpartition_by_hashContext) {}

// ExitSubpartition_by_hash is called when production subpartition_by_hash is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_by_hash(ctx *Subpartition_by_hashContext) {}

// EnterSubpartition_name is called when production subpartition_name is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_name(ctx *Subpartition_nameContext) {}

// ExitSubpartition_name is called when production subpartition_name is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_name(ctx *Subpartition_nameContext) {}

// EnterRange_subpartition_desc is called when production range_subpartition_desc is entered.
func (s *BasePlSqlParserListener) EnterRange_subpartition_desc(ctx *Range_subpartition_descContext) {}

// ExitRange_subpartition_desc is called when production range_subpartition_desc is exited.
func (s *BasePlSqlParserListener) ExitRange_subpartition_desc(ctx *Range_subpartition_descContext) {}

// EnterList_subpartition_desc is called when production list_subpartition_desc is entered.
func (s *BasePlSqlParserListener) EnterList_subpartition_desc(ctx *List_subpartition_descContext) {}

// ExitList_subpartition_desc is called when production list_subpartition_desc is exited.
func (s *BasePlSqlParserListener) ExitList_subpartition_desc(ctx *List_subpartition_descContext) {}

// EnterIndividual_hash_subparts is called when production individual_hash_subparts is entered.
func (s *BasePlSqlParserListener) EnterIndividual_hash_subparts(ctx *Individual_hash_subpartsContext) {
}

// ExitIndividual_hash_subparts is called when production individual_hash_subparts is exited.
func (s *BasePlSqlParserListener) ExitIndividual_hash_subparts(ctx *Individual_hash_subpartsContext) {
}

// EnterHash_subparts_by_quantity is called when production hash_subparts_by_quantity is entered.
func (s *BasePlSqlParserListener) EnterHash_subparts_by_quantity(ctx *Hash_subparts_by_quantityContext) {
}

// ExitHash_subparts_by_quantity is called when production hash_subparts_by_quantity is exited.
func (s *BasePlSqlParserListener) ExitHash_subparts_by_quantity(ctx *Hash_subparts_by_quantityContext) {
}

// EnterRange_values_clause is called when production range_values_clause is entered.
func (s *BasePlSqlParserListener) EnterRange_values_clause(ctx *Range_values_clauseContext) {}

// ExitRange_values_clause is called when production range_values_clause is exited.
func (s *BasePlSqlParserListener) ExitRange_values_clause(ctx *Range_values_clauseContext) {}

// EnterRange_values_list is called when production range_values_list is entered.
func (s *BasePlSqlParserListener) EnterRange_values_list(ctx *Range_values_listContext) {}

// ExitRange_values_list is called when production range_values_list is exited.
func (s *BasePlSqlParserListener) ExitRange_values_list(ctx *Range_values_listContext) {}

// EnterList_values_clause is called when production list_values_clause is entered.
func (s *BasePlSqlParserListener) EnterList_values_clause(ctx *List_values_clauseContext) {}

// ExitList_values_clause is called when production list_values_clause is exited.
func (s *BasePlSqlParserListener) ExitList_values_clause(ctx *List_values_clauseContext) {}

// EnterTable_partition_description is called when production table_partition_description is entered.
func (s *BasePlSqlParserListener) EnterTable_partition_description(ctx *Table_partition_descriptionContext) {
}

// ExitTable_partition_description is called when production table_partition_description is exited.
func (s *BasePlSqlParserListener) ExitTable_partition_description(ctx *Table_partition_descriptionContext) {
}

// EnterPartitioning_storage_clause is called when production partitioning_storage_clause is entered.
func (s *BasePlSqlParserListener) EnterPartitioning_storage_clause(ctx *Partitioning_storage_clauseContext) {
}

// ExitPartitioning_storage_clause is called when production partitioning_storage_clause is exited.
func (s *BasePlSqlParserListener) ExitPartitioning_storage_clause(ctx *Partitioning_storage_clauseContext) {
}

// EnterLob_partitioning_storage is called when production lob_partitioning_storage is entered.
func (s *BasePlSqlParserListener) EnterLob_partitioning_storage(ctx *Lob_partitioning_storageContext) {
}

// ExitLob_partitioning_storage is called when production lob_partitioning_storage is exited.
func (s *BasePlSqlParserListener) ExitLob_partitioning_storage(ctx *Lob_partitioning_storageContext) {
}

// EnterSize_clause is called when production size_clause is entered.
func (s *BasePlSqlParserListener) EnterSize_clause(ctx *Size_clauseContext) {}

// ExitSize_clause is called when production size_clause is exited.
func (s *BasePlSqlParserListener) ExitSize_clause(ctx *Size_clauseContext) {}

// EnterTable_compression is called when production table_compression is entered.
func (s *BasePlSqlParserListener) EnterTable_compression(ctx *Table_compressionContext) {}

// ExitTable_compression is called when production table_compression is exited.
func (s *BasePlSqlParserListener) ExitTable_compression(ctx *Table_compressionContext) {}

// EnterInmemory_table_clause is called when production inmemory_table_clause is entered.
func (s *BasePlSqlParserListener) EnterInmemory_table_clause(ctx *Inmemory_table_clauseContext) {}

// ExitInmemory_table_clause is called when production inmemory_table_clause is exited.
func (s *BasePlSqlParserListener) ExitInmemory_table_clause(ctx *Inmemory_table_clauseContext) {}

// EnterInmemory_attributes is called when production inmemory_attributes is entered.
func (s *BasePlSqlParserListener) EnterInmemory_attributes(ctx *Inmemory_attributesContext) {}

// ExitInmemory_attributes is called when production inmemory_attributes is exited.
func (s *BasePlSqlParserListener) ExitInmemory_attributes(ctx *Inmemory_attributesContext) {}

// EnterInmemory_memcompress is called when production inmemory_memcompress is entered.
func (s *BasePlSqlParserListener) EnterInmemory_memcompress(ctx *Inmemory_memcompressContext) {}

// ExitInmemory_memcompress is called when production inmemory_memcompress is exited.
func (s *BasePlSqlParserListener) ExitInmemory_memcompress(ctx *Inmemory_memcompressContext) {}

// EnterInmemory_priority is called when production inmemory_priority is entered.
func (s *BasePlSqlParserListener) EnterInmemory_priority(ctx *Inmemory_priorityContext) {}

// ExitInmemory_priority is called when production inmemory_priority is exited.
func (s *BasePlSqlParserListener) ExitInmemory_priority(ctx *Inmemory_priorityContext) {}

// EnterInmemory_distribute is called when production inmemory_distribute is entered.
func (s *BasePlSqlParserListener) EnterInmemory_distribute(ctx *Inmemory_distributeContext) {}

// ExitInmemory_distribute is called when production inmemory_distribute is exited.
func (s *BasePlSqlParserListener) ExitInmemory_distribute(ctx *Inmemory_distributeContext) {}

// EnterInmemory_duplicate is called when production inmemory_duplicate is entered.
func (s *BasePlSqlParserListener) EnterInmemory_duplicate(ctx *Inmemory_duplicateContext) {}

// ExitInmemory_duplicate is called when production inmemory_duplicate is exited.
func (s *BasePlSqlParserListener) ExitInmemory_duplicate(ctx *Inmemory_duplicateContext) {}

// EnterInmemory_column_clause is called when production inmemory_column_clause is entered.
func (s *BasePlSqlParserListener) EnterInmemory_column_clause(ctx *Inmemory_column_clauseContext) {}

// ExitInmemory_column_clause is called when production inmemory_column_clause is exited.
func (s *BasePlSqlParserListener) ExitInmemory_column_clause(ctx *Inmemory_column_clauseContext) {}

// EnterPhysical_attributes_clause is called when production physical_attributes_clause is entered.
func (s *BasePlSqlParserListener) EnterPhysical_attributes_clause(ctx *Physical_attributes_clauseContext) {
}

// ExitPhysical_attributes_clause is called when production physical_attributes_clause is exited.
func (s *BasePlSqlParserListener) ExitPhysical_attributes_clause(ctx *Physical_attributes_clauseContext) {
}

// EnterStorage_clause is called when production storage_clause is entered.
func (s *BasePlSqlParserListener) EnterStorage_clause(ctx *Storage_clauseContext) {}

// ExitStorage_clause is called when production storage_clause is exited.
func (s *BasePlSqlParserListener) ExitStorage_clause(ctx *Storage_clauseContext) {}

// EnterDeferred_segment_creation is called when production deferred_segment_creation is entered.
func (s *BasePlSqlParserListener) EnterDeferred_segment_creation(ctx *Deferred_segment_creationContext) {
}

// ExitDeferred_segment_creation is called when production deferred_segment_creation is exited.
func (s *BasePlSqlParserListener) ExitDeferred_segment_creation(ctx *Deferred_segment_creationContext) {
}

// EnterSegment_attributes_clause is called when production segment_attributes_clause is entered.
func (s *BasePlSqlParserListener) EnterSegment_attributes_clause(ctx *Segment_attributes_clauseContext) {
}

// ExitSegment_attributes_clause is called when production segment_attributes_clause is exited.
func (s *BasePlSqlParserListener) ExitSegment_attributes_clause(ctx *Segment_attributes_clauseContext) {
}

// EnterPhysical_properties is called when production physical_properties is entered.
func (s *BasePlSqlParserListener) EnterPhysical_properties(ctx *Physical_propertiesContext) {}

// ExitPhysical_properties is called when production physical_properties is exited.
func (s *BasePlSqlParserListener) ExitPhysical_properties(ctx *Physical_propertiesContext) {}

// EnterIlm_clause is called when production ilm_clause is entered.
func (s *BasePlSqlParserListener) EnterIlm_clause(ctx *Ilm_clauseContext) {}

// ExitIlm_clause is called when production ilm_clause is exited.
func (s *BasePlSqlParserListener) ExitIlm_clause(ctx *Ilm_clauseContext) {}

// EnterIlm_policy_clause is called when production ilm_policy_clause is entered.
func (s *BasePlSqlParserListener) EnterIlm_policy_clause(ctx *Ilm_policy_clauseContext) {}

// ExitIlm_policy_clause is called when production ilm_policy_clause is exited.
func (s *BasePlSqlParserListener) ExitIlm_policy_clause(ctx *Ilm_policy_clauseContext) {}

// EnterIlm_compression_policy is called when production ilm_compression_policy is entered.
func (s *BasePlSqlParserListener) EnterIlm_compression_policy(ctx *Ilm_compression_policyContext) {}

// ExitIlm_compression_policy is called when production ilm_compression_policy is exited.
func (s *BasePlSqlParserListener) ExitIlm_compression_policy(ctx *Ilm_compression_policyContext) {}

// EnterIlm_tiering_policy is called when production ilm_tiering_policy is entered.
func (s *BasePlSqlParserListener) EnterIlm_tiering_policy(ctx *Ilm_tiering_policyContext) {}

// ExitIlm_tiering_policy is called when production ilm_tiering_policy is exited.
func (s *BasePlSqlParserListener) ExitIlm_tiering_policy(ctx *Ilm_tiering_policyContext) {}

// EnterIlm_after_on is called when production ilm_after_on is entered.
func (s *BasePlSqlParserListener) EnterIlm_after_on(ctx *Ilm_after_onContext) {}

// ExitIlm_after_on is called when production ilm_after_on is exited.
func (s *BasePlSqlParserListener) ExitIlm_after_on(ctx *Ilm_after_onContext) {}

// EnterSegment_group is called when production segment_group is entered.
func (s *BasePlSqlParserListener) EnterSegment_group(ctx *Segment_groupContext) {}

// ExitSegment_group is called when production segment_group is exited.
func (s *BasePlSqlParserListener) ExitSegment_group(ctx *Segment_groupContext) {}

// EnterIlm_inmemory_policy is called when production ilm_inmemory_policy is entered.
func (s *BasePlSqlParserListener) EnterIlm_inmemory_policy(ctx *Ilm_inmemory_policyContext) {}

// ExitIlm_inmemory_policy is called when production ilm_inmemory_policy is exited.
func (s *BasePlSqlParserListener) ExitIlm_inmemory_policy(ctx *Ilm_inmemory_policyContext) {}

// EnterIlm_time_period is called when production ilm_time_period is entered.
func (s *BasePlSqlParserListener) EnterIlm_time_period(ctx *Ilm_time_periodContext) {}

// ExitIlm_time_period is called when production ilm_time_period is exited.
func (s *BasePlSqlParserListener) ExitIlm_time_period(ctx *Ilm_time_periodContext) {}

// EnterHeap_org_table_clause is called when production heap_org_table_clause is entered.
func (s *BasePlSqlParserListener) EnterHeap_org_table_clause(ctx *Heap_org_table_clauseContext) {}

// ExitHeap_org_table_clause is called when production heap_org_table_clause is exited.
func (s *BasePlSqlParserListener) ExitHeap_org_table_clause(ctx *Heap_org_table_clauseContext) {}

// EnterExternal_table_clause is called when production external_table_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_clause(ctx *External_table_clauseContext) {}

// ExitExternal_table_clause is called when production external_table_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_clause(ctx *External_table_clauseContext) {}

// EnterAccess_driver_type is called when production access_driver_type is entered.
func (s *BasePlSqlParserListener) EnterAccess_driver_type(ctx *Access_driver_typeContext) {}

// ExitAccess_driver_type is called when production access_driver_type is exited.
func (s *BasePlSqlParserListener) ExitAccess_driver_type(ctx *Access_driver_typeContext) {}

// EnterExternal_table_data_props is called when production external_table_data_props is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_data_props(ctx *External_table_data_propsContext) {
}

// ExitExternal_table_data_props is called when production external_table_data_props is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_data_props(ctx *External_table_data_propsContext) {
}

// EnterExternal_table_data_format is called when production external_table_data_format is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_data_format(ctx *External_table_data_formatContext) {
}

// ExitExternal_table_data_format is called when production external_table_data_format is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_data_format(ctx *External_table_data_formatContext) {
}

// EnterExternal_table_transform is called when production external_table_transform is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_transform(ctx *External_table_transformContext) {
}

// ExitExternal_table_transform is called when production external_table_transform is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_transform(ctx *External_table_transformContext) {
}

// EnterExternal_table_field is called when production external_table_field is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_field(ctx *External_table_fieldContext) {}

// ExitExternal_table_field is called when production external_table_field is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_field(ctx *External_table_fieldContext) {}

// EnterExternal_table_field_list is called when production external_table_field_list is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_field_list(ctx *External_table_field_listContext) {
}

// ExitExternal_table_field_list is called when production external_table_field_list is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_field_list(ctx *External_table_field_listContext) {
}

// EnterExternal_table_fields_clause is called when production external_table_fields_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_fields_clause(ctx *External_table_fields_clauseContext) {
}

// ExitExternal_table_fields_clause is called when production external_table_fields_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_fields_clause(ctx *External_table_fields_clauseContext) {
}

// EnterExternal_table_position_clause is called when production external_table_position_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_position_clause(ctx *External_table_position_clauseContext) {
}

// ExitExternal_table_position_clause is called when production external_table_position_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_position_clause(ctx *External_table_position_clauseContext) {
}

// EnterExternal_table_datatype_clause is called when production external_table_datatype_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_datatype_clause(ctx *External_table_datatype_clauseContext) {
}

// ExitExternal_table_datatype_clause is called when production external_table_datatype_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_datatype_clause(ctx *External_table_datatype_clauseContext) {
}

// EnterExternal_table_delimit_clause is called when production external_table_delimit_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_delimit_clause(ctx *External_table_delimit_clauseContext) {
}

// ExitExternal_table_delimit_clause is called when production external_table_delimit_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_delimit_clause(ctx *External_table_delimit_clauseContext) {
}

// EnterExternal_table_trim_clause is called when production external_table_trim_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_trim_clause(ctx *External_table_trim_clauseContext) {
}

// ExitExternal_table_trim_clause is called when production external_table_trim_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_trim_clause(ctx *External_table_trim_clauseContext) {
}

// EnterExternal_table_date_format_clause is called when production external_table_date_format_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_date_format_clause(ctx *External_table_date_format_clauseContext) {
}

// ExitExternal_table_date_format_clause is called when production external_table_date_format_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_date_format_clause(ctx *External_table_date_format_clauseContext) {
}

// EnterExternal_table_init_clause is called when production external_table_init_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_init_clause(ctx *External_table_init_clauseContext) {
}

// ExitExternal_table_init_clause is called when production external_table_init_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_init_clause(ctx *External_table_init_clauseContext) {
}

// EnterExternal_table_condition_clause is called when production external_table_condition_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_condition_clause(ctx *External_table_condition_clauseContext) {
}

// ExitExternal_table_condition_clause is called when production external_table_condition_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_condition_clause(ctx *External_table_condition_clauseContext) {
}

// EnterExternal_table_lls_clause is called when production external_table_lls_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_lls_clause(ctx *External_table_lls_clauseContext) {
}

// ExitExternal_table_lls_clause is called when production external_table_lls_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_lls_clause(ctx *External_table_lls_clauseContext) {
}

// EnterExternal_table_records is called when production external_table_records is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_records(ctx *External_table_recordsContext) {}

// ExitExternal_table_records is called when production external_table_records is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_records(ctx *External_table_recordsContext) {}

// EnterExternal_table_record_options_clause is called when production external_table_record_options_clause is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_record_options_clause(ctx *External_table_record_options_clauseContext) {
}

// ExitExternal_table_record_options_clause is called when production external_table_record_options_clause is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_record_options_clause(ctx *External_table_record_options_clauseContext) {
}

// EnterExternal_table_output_files is called when production external_table_output_files is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_output_files(ctx *External_table_output_filesContext) {
}

// ExitExternal_table_output_files is called when production external_table_output_files is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_output_files(ctx *External_table_output_filesContext) {
}

// EnterExternal_table_fields is called when production external_table_fields is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_fields(ctx *External_table_fieldsContext) {}

// ExitExternal_table_fields is called when production external_table_fields is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_fields(ctx *External_table_fieldsContext) {}

// EnterExternal_table_datapump is called when production external_table_datapump is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_datapump(ctx *External_table_datapumpContext) {}

// ExitExternal_table_datapump is called when production external_table_datapump is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_datapump(ctx *External_table_datapumpContext) {}

// EnterExternal_table_hive is called when production external_table_hive is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_hive(ctx *External_table_hiveContext) {}

// ExitExternal_table_hive is called when production external_table_hive is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_hive(ctx *External_table_hiveContext) {}

// EnterExternal_table_hive_parameter_map is called when production external_table_hive_parameter_map is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_hive_parameter_map(ctx *External_table_hive_parameter_mapContext) {
}

// ExitExternal_table_hive_parameter_map is called when production external_table_hive_parameter_map is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_hive_parameter_map(ctx *External_table_hive_parameter_mapContext) {
}

// EnterExternal_table_hive_parameter_map_entry is called when production external_table_hive_parameter_map_entry is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_hive_parameter_map_entry(ctx *External_table_hive_parameter_map_entryContext) {
}

// ExitExternal_table_hive_parameter_map_entry is called when production external_table_hive_parameter_map_entry is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_hive_parameter_map_entry(ctx *External_table_hive_parameter_map_entryContext) {
}

// EnterExternal_table_directory is called when production external_table_directory is entered.
func (s *BasePlSqlParserListener) EnterExternal_table_directory(ctx *External_table_directoryContext) {
}

// ExitExternal_table_directory is called when production external_table_directory is exited.
func (s *BasePlSqlParserListener) ExitExternal_table_directory(ctx *External_table_directoryContext) {
}

// EnterRow_movement_clause is called when production row_movement_clause is entered.
func (s *BasePlSqlParserListener) EnterRow_movement_clause(ctx *Row_movement_clauseContext) {}

// ExitRow_movement_clause is called when production row_movement_clause is exited.
func (s *BasePlSqlParserListener) ExitRow_movement_clause(ctx *Row_movement_clauseContext) {}

// EnterFlashback_archive_clause is called when production flashback_archive_clause is entered.
func (s *BasePlSqlParserListener) EnterFlashback_archive_clause(ctx *Flashback_archive_clauseContext) {
}

// ExitFlashback_archive_clause is called when production flashback_archive_clause is exited.
func (s *BasePlSqlParserListener) ExitFlashback_archive_clause(ctx *Flashback_archive_clauseContext) {
}

// EnterLog_grp is called when production log_grp is entered.
func (s *BasePlSqlParserListener) EnterLog_grp(ctx *Log_grpContext) {}

// ExitLog_grp is called when production log_grp is exited.
func (s *BasePlSqlParserListener) ExitLog_grp(ctx *Log_grpContext) {}

// EnterSupplemental_table_logging is called when production supplemental_table_logging is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_table_logging(ctx *Supplemental_table_loggingContext) {
}

// ExitSupplemental_table_logging is called when production supplemental_table_logging is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_table_logging(ctx *Supplemental_table_loggingContext) {
}

// EnterSupplemental_log_grp_clause is called when production supplemental_log_grp_clause is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_log_grp_clause(ctx *Supplemental_log_grp_clauseContext) {
}

// ExitSupplemental_log_grp_clause is called when production supplemental_log_grp_clause is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_log_grp_clause(ctx *Supplemental_log_grp_clauseContext) {
}

// EnterSupplemental_id_key_clause is called when production supplemental_id_key_clause is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_id_key_clause(ctx *Supplemental_id_key_clauseContext) {
}

// ExitSupplemental_id_key_clause is called when production supplemental_id_key_clause is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_id_key_clause(ctx *Supplemental_id_key_clauseContext) {
}

// EnterAllocate_extent_clause is called when production allocate_extent_clause is entered.
func (s *BasePlSqlParserListener) EnterAllocate_extent_clause(ctx *Allocate_extent_clauseContext) {}

// ExitAllocate_extent_clause is called when production allocate_extent_clause is exited.
func (s *BasePlSqlParserListener) ExitAllocate_extent_clause(ctx *Allocate_extent_clauseContext) {}

// EnterDeallocate_unused_clause is called when production deallocate_unused_clause is entered.
func (s *BasePlSqlParserListener) EnterDeallocate_unused_clause(ctx *Deallocate_unused_clauseContext) {
}

// ExitDeallocate_unused_clause is called when production deallocate_unused_clause is exited.
func (s *BasePlSqlParserListener) ExitDeallocate_unused_clause(ctx *Deallocate_unused_clauseContext) {
}

// EnterShrink_clause is called when production shrink_clause is entered.
func (s *BasePlSqlParserListener) EnterShrink_clause(ctx *Shrink_clauseContext) {}

// ExitShrink_clause is called when production shrink_clause is exited.
func (s *BasePlSqlParserListener) ExitShrink_clause(ctx *Shrink_clauseContext) {}

// EnterRecords_per_block_clause is called when production records_per_block_clause is entered.
func (s *BasePlSqlParserListener) EnterRecords_per_block_clause(ctx *Records_per_block_clauseContext) {
}

// ExitRecords_per_block_clause is called when production records_per_block_clause is exited.
func (s *BasePlSqlParserListener) ExitRecords_per_block_clause(ctx *Records_per_block_clauseContext) {
}

// EnterUpgrade_table_clause is called when production upgrade_table_clause is entered.
func (s *BasePlSqlParserListener) EnterUpgrade_table_clause(ctx *Upgrade_table_clauseContext) {}

// ExitUpgrade_table_clause is called when production upgrade_table_clause is exited.
func (s *BasePlSqlParserListener) ExitUpgrade_table_clause(ctx *Upgrade_table_clauseContext) {}

// EnterTruncate_table is called when production truncate_table is entered.
func (s *BasePlSqlParserListener) EnterTruncate_table(ctx *Truncate_tableContext) {}

// ExitTruncate_table is called when production truncate_table is exited.
func (s *BasePlSqlParserListener) ExitTruncate_table(ctx *Truncate_tableContext) {}

// EnterDrop_table is called when production drop_table is entered.
func (s *BasePlSqlParserListener) EnterDrop_table(ctx *Drop_tableContext) {}

// ExitDrop_table is called when production drop_table is exited.
func (s *BasePlSqlParserListener) ExitDrop_table(ctx *Drop_tableContext) {}

// EnterDrop_tablespace is called when production drop_tablespace is entered.
func (s *BasePlSqlParserListener) EnterDrop_tablespace(ctx *Drop_tablespaceContext) {}

// ExitDrop_tablespace is called when production drop_tablespace is exited.
func (s *BasePlSqlParserListener) ExitDrop_tablespace(ctx *Drop_tablespaceContext) {}

// EnterDrop_tablespace_set is called when production drop_tablespace_set is entered.
func (s *BasePlSqlParserListener) EnterDrop_tablespace_set(ctx *Drop_tablespace_setContext) {}

// ExitDrop_tablespace_set is called when production drop_tablespace_set is exited.
func (s *BasePlSqlParserListener) ExitDrop_tablespace_set(ctx *Drop_tablespace_setContext) {}

// EnterIncluding_contents_clause is called when production including_contents_clause is entered.
func (s *BasePlSqlParserListener) EnterIncluding_contents_clause(ctx *Including_contents_clauseContext) {
}

// ExitIncluding_contents_clause is called when production including_contents_clause is exited.
func (s *BasePlSqlParserListener) ExitIncluding_contents_clause(ctx *Including_contents_clauseContext) {
}

// EnterDrop_view is called when production drop_view is entered.
func (s *BasePlSqlParserListener) EnterDrop_view(ctx *Drop_viewContext) {}

// ExitDrop_view is called when production drop_view is exited.
func (s *BasePlSqlParserListener) ExitDrop_view(ctx *Drop_viewContext) {}

// EnterComment_on_column is called when production comment_on_column is entered.
func (s *BasePlSqlParserListener) EnterComment_on_column(ctx *Comment_on_columnContext) {}

// ExitComment_on_column is called when production comment_on_column is exited.
func (s *BasePlSqlParserListener) ExitComment_on_column(ctx *Comment_on_columnContext) {}

// EnterEnable_or_disable is called when production enable_or_disable is entered.
func (s *BasePlSqlParserListener) EnterEnable_or_disable(ctx *Enable_or_disableContext) {}

// ExitEnable_or_disable is called when production enable_or_disable is exited.
func (s *BasePlSqlParserListener) ExitEnable_or_disable(ctx *Enable_or_disableContext) {}

// EnterAllow_or_disallow is called when production allow_or_disallow is entered.
func (s *BasePlSqlParserListener) EnterAllow_or_disallow(ctx *Allow_or_disallowContext) {}

// ExitAllow_or_disallow is called when production allow_or_disallow is exited.
func (s *BasePlSqlParserListener) ExitAllow_or_disallow(ctx *Allow_or_disallowContext) {}

// EnterAlter_synonym is called when production alter_synonym is entered.
func (s *BasePlSqlParserListener) EnterAlter_synonym(ctx *Alter_synonymContext) {}

// ExitAlter_synonym is called when production alter_synonym is exited.
func (s *BasePlSqlParserListener) ExitAlter_synonym(ctx *Alter_synonymContext) {}

// EnterCreate_synonym is called when production create_synonym is entered.
func (s *BasePlSqlParserListener) EnterCreate_synonym(ctx *Create_synonymContext) {}

// ExitCreate_synonym is called when production create_synonym is exited.
func (s *BasePlSqlParserListener) ExitCreate_synonym(ctx *Create_synonymContext) {}

// EnterDrop_synonym is called when production drop_synonym is entered.
func (s *BasePlSqlParserListener) EnterDrop_synonym(ctx *Drop_synonymContext) {}

// ExitDrop_synonym is called when production drop_synonym is exited.
func (s *BasePlSqlParserListener) ExitDrop_synonym(ctx *Drop_synonymContext) {}

// EnterCreate_spfile is called when production create_spfile is entered.
func (s *BasePlSqlParserListener) EnterCreate_spfile(ctx *Create_spfileContext) {}

// ExitCreate_spfile is called when production create_spfile is exited.
func (s *BasePlSqlParserListener) ExitCreate_spfile(ctx *Create_spfileContext) {}

// EnterSpfile_name is called when production spfile_name is entered.
func (s *BasePlSqlParserListener) EnterSpfile_name(ctx *Spfile_nameContext) {}

// ExitSpfile_name is called when production spfile_name is exited.
func (s *BasePlSqlParserListener) ExitSpfile_name(ctx *Spfile_nameContext) {}

// EnterPfile_name is called when production pfile_name is entered.
func (s *BasePlSqlParserListener) EnterPfile_name(ctx *Pfile_nameContext) {}

// ExitPfile_name is called when production pfile_name is exited.
func (s *BasePlSqlParserListener) ExitPfile_name(ctx *Pfile_nameContext) {}

// EnterComment_on_table is called when production comment_on_table is entered.
func (s *BasePlSqlParserListener) EnterComment_on_table(ctx *Comment_on_tableContext) {}

// ExitComment_on_table is called when production comment_on_table is exited.
func (s *BasePlSqlParserListener) ExitComment_on_table(ctx *Comment_on_tableContext) {}

// EnterComment_on_materialized is called when production comment_on_materialized is entered.
func (s *BasePlSqlParserListener) EnterComment_on_materialized(ctx *Comment_on_materializedContext) {}

// ExitComment_on_materialized is called when production comment_on_materialized is exited.
func (s *BasePlSqlParserListener) ExitComment_on_materialized(ctx *Comment_on_materializedContext) {}

// EnterAlter_analytic_view is called when production alter_analytic_view is entered.
func (s *BasePlSqlParserListener) EnterAlter_analytic_view(ctx *Alter_analytic_viewContext) {}

// ExitAlter_analytic_view is called when production alter_analytic_view is exited.
func (s *BasePlSqlParserListener) ExitAlter_analytic_view(ctx *Alter_analytic_viewContext) {}

// EnterAlter_add_cache_clause is called when production alter_add_cache_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_add_cache_clause(ctx *Alter_add_cache_clauseContext) {}

// ExitAlter_add_cache_clause is called when production alter_add_cache_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_add_cache_clause(ctx *Alter_add_cache_clauseContext) {}

// EnterLevels_item is called when production levels_item is entered.
func (s *BasePlSqlParserListener) EnterLevels_item(ctx *Levels_itemContext) {}

// ExitLevels_item is called when production levels_item is exited.
func (s *BasePlSqlParserListener) ExitLevels_item(ctx *Levels_itemContext) {}

// EnterMeasure_list is called when production measure_list is entered.
func (s *BasePlSqlParserListener) EnterMeasure_list(ctx *Measure_listContext) {}

// ExitMeasure_list is called when production measure_list is exited.
func (s *BasePlSqlParserListener) ExitMeasure_list(ctx *Measure_listContext) {}

// EnterAlter_drop_cache_clause is called when production alter_drop_cache_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_drop_cache_clause(ctx *Alter_drop_cache_clauseContext) {}

// ExitAlter_drop_cache_clause is called when production alter_drop_cache_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_drop_cache_clause(ctx *Alter_drop_cache_clauseContext) {}

// EnterAlter_attribute_dimension is called when production alter_attribute_dimension is entered.
func (s *BasePlSqlParserListener) EnterAlter_attribute_dimension(ctx *Alter_attribute_dimensionContext) {
}

// ExitAlter_attribute_dimension is called when production alter_attribute_dimension is exited.
func (s *BasePlSqlParserListener) ExitAlter_attribute_dimension(ctx *Alter_attribute_dimensionContext) {
}

// EnterAlter_audit_policy is called when production alter_audit_policy is entered.
func (s *BasePlSqlParserListener) EnterAlter_audit_policy(ctx *Alter_audit_policyContext) {}

// ExitAlter_audit_policy is called when production alter_audit_policy is exited.
func (s *BasePlSqlParserListener) ExitAlter_audit_policy(ctx *Alter_audit_policyContext) {}

// EnterAlter_cluster is called when production alter_cluster is entered.
func (s *BasePlSqlParserListener) EnterAlter_cluster(ctx *Alter_clusterContext) {}

// ExitAlter_cluster is called when production alter_cluster is exited.
func (s *BasePlSqlParserListener) ExitAlter_cluster(ctx *Alter_clusterContext) {}

// EnterDrop_analytic_view is called when production drop_analytic_view is entered.
func (s *BasePlSqlParserListener) EnterDrop_analytic_view(ctx *Drop_analytic_viewContext) {}

// ExitDrop_analytic_view is called when production drop_analytic_view is exited.
func (s *BasePlSqlParserListener) ExitDrop_analytic_view(ctx *Drop_analytic_viewContext) {}

// EnterDrop_attribute_dimension is called when production drop_attribute_dimension is entered.
func (s *BasePlSqlParserListener) EnterDrop_attribute_dimension(ctx *Drop_attribute_dimensionContext) {
}

// ExitDrop_attribute_dimension is called when production drop_attribute_dimension is exited.
func (s *BasePlSqlParserListener) ExitDrop_attribute_dimension(ctx *Drop_attribute_dimensionContext) {
}

// EnterDrop_audit_policy is called when production drop_audit_policy is entered.
func (s *BasePlSqlParserListener) EnterDrop_audit_policy(ctx *Drop_audit_policyContext) {}

// ExitDrop_audit_policy is called when production drop_audit_policy is exited.
func (s *BasePlSqlParserListener) ExitDrop_audit_policy(ctx *Drop_audit_policyContext) {}

// EnterDrop_flashback_archive is called when production drop_flashback_archive is entered.
func (s *BasePlSqlParserListener) EnterDrop_flashback_archive(ctx *Drop_flashback_archiveContext) {}

// ExitDrop_flashback_archive is called when production drop_flashback_archive is exited.
func (s *BasePlSqlParserListener) ExitDrop_flashback_archive(ctx *Drop_flashback_archiveContext) {}

// EnterDrop_cluster is called when production drop_cluster is entered.
func (s *BasePlSqlParserListener) EnterDrop_cluster(ctx *Drop_clusterContext) {}

// ExitDrop_cluster is called when production drop_cluster is exited.
func (s *BasePlSqlParserListener) ExitDrop_cluster(ctx *Drop_clusterContext) {}

// EnterDrop_context is called when production drop_context is entered.
func (s *BasePlSqlParserListener) EnterDrop_context(ctx *Drop_contextContext) {}

// ExitDrop_context is called when production drop_context is exited.
func (s *BasePlSqlParserListener) ExitDrop_context(ctx *Drop_contextContext) {}

// EnterDrop_directory is called when production drop_directory is entered.
func (s *BasePlSqlParserListener) EnterDrop_directory(ctx *Drop_directoryContext) {}

// ExitDrop_directory is called when production drop_directory is exited.
func (s *BasePlSqlParserListener) ExitDrop_directory(ctx *Drop_directoryContext) {}

// EnterDrop_diskgroup is called when production drop_diskgroup is entered.
func (s *BasePlSqlParserListener) EnterDrop_diskgroup(ctx *Drop_diskgroupContext) {}

// ExitDrop_diskgroup is called when production drop_diskgroup is exited.
func (s *BasePlSqlParserListener) ExitDrop_diskgroup(ctx *Drop_diskgroupContext) {}

// EnterDrop_edition is called when production drop_edition is entered.
func (s *BasePlSqlParserListener) EnterDrop_edition(ctx *Drop_editionContext) {}

// ExitDrop_edition is called when production drop_edition is exited.
func (s *BasePlSqlParserListener) ExitDrop_edition(ctx *Drop_editionContext) {}

// EnterTruncate_cluster is called when production truncate_cluster is entered.
func (s *BasePlSqlParserListener) EnterTruncate_cluster(ctx *Truncate_clusterContext) {}

// ExitTruncate_cluster is called when production truncate_cluster is exited.
func (s *BasePlSqlParserListener) ExitTruncate_cluster(ctx *Truncate_clusterContext) {}

// EnterCache_or_nocache is called when production cache_or_nocache is entered.
func (s *BasePlSqlParserListener) EnterCache_or_nocache(ctx *Cache_or_nocacheContext) {}

// ExitCache_or_nocache is called when production cache_or_nocache is exited.
func (s *BasePlSqlParserListener) ExitCache_or_nocache(ctx *Cache_or_nocacheContext) {}

// EnterDatabase_name is called when production database_name is entered.
func (s *BasePlSqlParserListener) EnterDatabase_name(ctx *Database_nameContext) {}

// ExitDatabase_name is called when production database_name is exited.
func (s *BasePlSqlParserListener) ExitDatabase_name(ctx *Database_nameContext) {}

// EnterAlter_database is called when production alter_database is entered.
func (s *BasePlSqlParserListener) EnterAlter_database(ctx *Alter_databaseContext) {}

// ExitAlter_database is called when production alter_database is exited.
func (s *BasePlSqlParserListener) ExitAlter_database(ctx *Alter_databaseContext) {}

// EnterDatabase_clause is called when production database_clause is entered.
func (s *BasePlSqlParserListener) EnterDatabase_clause(ctx *Database_clauseContext) {}

// ExitDatabase_clause is called when production database_clause is exited.
func (s *BasePlSqlParserListener) ExitDatabase_clause(ctx *Database_clauseContext) {}

// EnterStartup_clauses is called when production startup_clauses is entered.
func (s *BasePlSqlParserListener) EnterStartup_clauses(ctx *Startup_clausesContext) {}

// ExitStartup_clauses is called when production startup_clauses is exited.
func (s *BasePlSqlParserListener) ExitStartup_clauses(ctx *Startup_clausesContext) {}

// EnterResetlogs_or_noresetlogs is called when production resetlogs_or_noresetlogs is entered.
func (s *BasePlSqlParserListener) EnterResetlogs_or_noresetlogs(ctx *Resetlogs_or_noresetlogsContext) {
}

// ExitResetlogs_or_noresetlogs is called when production resetlogs_or_noresetlogs is exited.
func (s *BasePlSqlParserListener) ExitResetlogs_or_noresetlogs(ctx *Resetlogs_or_noresetlogsContext) {
}

// EnterUpgrade_or_downgrade is called when production upgrade_or_downgrade is entered.
func (s *BasePlSqlParserListener) EnterUpgrade_or_downgrade(ctx *Upgrade_or_downgradeContext) {}

// ExitUpgrade_or_downgrade is called when production upgrade_or_downgrade is exited.
func (s *BasePlSqlParserListener) ExitUpgrade_or_downgrade(ctx *Upgrade_or_downgradeContext) {}

// EnterRecovery_clauses is called when production recovery_clauses is entered.
func (s *BasePlSqlParserListener) EnterRecovery_clauses(ctx *Recovery_clausesContext) {}

// ExitRecovery_clauses is called when production recovery_clauses is exited.
func (s *BasePlSqlParserListener) ExitRecovery_clauses(ctx *Recovery_clausesContext) {}

// EnterBegin_or_end is called when production begin_or_end is entered.
func (s *BasePlSqlParserListener) EnterBegin_or_end(ctx *Begin_or_endContext) {}

// ExitBegin_or_end is called when production begin_or_end is exited.
func (s *BasePlSqlParserListener) ExitBegin_or_end(ctx *Begin_or_endContext) {}

// EnterGeneral_recovery is called when production general_recovery is entered.
func (s *BasePlSqlParserListener) EnterGeneral_recovery(ctx *General_recoveryContext) {}

// ExitGeneral_recovery is called when production general_recovery is exited.
func (s *BasePlSqlParserListener) ExitGeneral_recovery(ctx *General_recoveryContext) {}

// EnterFull_database_recovery is called when production full_database_recovery is entered.
func (s *BasePlSqlParserListener) EnterFull_database_recovery(ctx *Full_database_recoveryContext) {}

// ExitFull_database_recovery is called when production full_database_recovery is exited.
func (s *BasePlSqlParserListener) ExitFull_database_recovery(ctx *Full_database_recoveryContext) {}

// EnterPartial_database_recovery is called when production partial_database_recovery is entered.
func (s *BasePlSqlParserListener) EnterPartial_database_recovery(ctx *Partial_database_recoveryContext) {
}

// ExitPartial_database_recovery is called when production partial_database_recovery is exited.
func (s *BasePlSqlParserListener) ExitPartial_database_recovery(ctx *Partial_database_recoveryContext) {
}

// EnterPartial_database_recovery_10g is called when production partial_database_recovery_10g is entered.
func (s *BasePlSqlParserListener) EnterPartial_database_recovery_10g(ctx *Partial_database_recovery_10gContext) {
}

// ExitPartial_database_recovery_10g is called when production partial_database_recovery_10g is exited.
func (s *BasePlSqlParserListener) ExitPartial_database_recovery_10g(ctx *Partial_database_recovery_10gContext) {
}

// EnterManaged_standby_recovery is called when production managed_standby_recovery is entered.
func (s *BasePlSqlParserListener) EnterManaged_standby_recovery(ctx *Managed_standby_recoveryContext) {
}

// ExitManaged_standby_recovery is called when production managed_standby_recovery is exited.
func (s *BasePlSqlParserListener) ExitManaged_standby_recovery(ctx *Managed_standby_recoveryContext) {
}

// EnterDb_name is called when production db_name is entered.
func (s *BasePlSqlParserListener) EnterDb_name(ctx *Db_nameContext) {}

// ExitDb_name is called when production db_name is exited.
func (s *BasePlSqlParserListener) ExitDb_name(ctx *Db_nameContext) {}

// EnterDatabase_file_clauses is called when production database_file_clauses is entered.
func (s *BasePlSqlParserListener) EnterDatabase_file_clauses(ctx *Database_file_clausesContext) {}

// ExitDatabase_file_clauses is called when production database_file_clauses is exited.
func (s *BasePlSqlParserListener) ExitDatabase_file_clauses(ctx *Database_file_clausesContext) {}

// EnterCreate_datafile_clause is called when production create_datafile_clause is entered.
func (s *BasePlSqlParserListener) EnterCreate_datafile_clause(ctx *Create_datafile_clauseContext) {}

// ExitCreate_datafile_clause is called when production create_datafile_clause is exited.
func (s *BasePlSqlParserListener) ExitCreate_datafile_clause(ctx *Create_datafile_clauseContext) {}

// EnterAlter_datafile_clause is called when production alter_datafile_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_datafile_clause(ctx *Alter_datafile_clauseContext) {}

// ExitAlter_datafile_clause is called when production alter_datafile_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_datafile_clause(ctx *Alter_datafile_clauseContext) {}

// EnterAlter_tempfile_clause is called when production alter_tempfile_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_tempfile_clause(ctx *Alter_tempfile_clauseContext) {}

// ExitAlter_tempfile_clause is called when production alter_tempfile_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_tempfile_clause(ctx *Alter_tempfile_clauseContext) {}

// EnterMove_datafile_clause is called when production move_datafile_clause is entered.
func (s *BasePlSqlParserListener) EnterMove_datafile_clause(ctx *Move_datafile_clauseContext) {}

// ExitMove_datafile_clause is called when production move_datafile_clause is exited.
func (s *BasePlSqlParserListener) ExitMove_datafile_clause(ctx *Move_datafile_clauseContext) {}

// EnterLogfile_clauses is called when production logfile_clauses is entered.
func (s *BasePlSqlParserListener) EnterLogfile_clauses(ctx *Logfile_clausesContext) {}

// ExitLogfile_clauses is called when production logfile_clauses is exited.
func (s *BasePlSqlParserListener) ExitLogfile_clauses(ctx *Logfile_clausesContext) {}

// EnterAdd_logfile_clauses is called when production add_logfile_clauses is entered.
func (s *BasePlSqlParserListener) EnterAdd_logfile_clauses(ctx *Add_logfile_clausesContext) {}

// ExitAdd_logfile_clauses is called when production add_logfile_clauses is exited.
func (s *BasePlSqlParserListener) ExitAdd_logfile_clauses(ctx *Add_logfile_clausesContext) {}

// EnterGroup_redo_logfile is called when production group_redo_logfile is entered.
func (s *BasePlSqlParserListener) EnterGroup_redo_logfile(ctx *Group_redo_logfileContext) {}

// ExitGroup_redo_logfile is called when production group_redo_logfile is exited.
func (s *BasePlSqlParserListener) ExitGroup_redo_logfile(ctx *Group_redo_logfileContext) {}

// EnterDrop_logfile_clauses is called when production drop_logfile_clauses is entered.
func (s *BasePlSqlParserListener) EnterDrop_logfile_clauses(ctx *Drop_logfile_clausesContext) {}

// ExitDrop_logfile_clauses is called when production drop_logfile_clauses is exited.
func (s *BasePlSqlParserListener) ExitDrop_logfile_clauses(ctx *Drop_logfile_clausesContext) {}

// EnterSwitch_logfile_clause is called when production switch_logfile_clause is entered.
func (s *BasePlSqlParserListener) EnterSwitch_logfile_clause(ctx *Switch_logfile_clauseContext) {}

// ExitSwitch_logfile_clause is called when production switch_logfile_clause is exited.
func (s *BasePlSqlParserListener) ExitSwitch_logfile_clause(ctx *Switch_logfile_clauseContext) {}

// EnterSupplemental_db_logging is called when production supplemental_db_logging is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_db_logging(ctx *Supplemental_db_loggingContext) {}

// ExitSupplemental_db_logging is called when production supplemental_db_logging is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_db_logging(ctx *Supplemental_db_loggingContext) {}

// EnterAdd_or_drop is called when production add_or_drop is entered.
func (s *BasePlSqlParserListener) EnterAdd_or_drop(ctx *Add_or_dropContext) {}

// ExitAdd_or_drop is called when production add_or_drop is exited.
func (s *BasePlSqlParserListener) ExitAdd_or_drop(ctx *Add_or_dropContext) {}

// EnterSupplemental_plsql_clause is called when production supplemental_plsql_clause is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_plsql_clause(ctx *Supplemental_plsql_clauseContext) {
}

// ExitSupplemental_plsql_clause is called when production supplemental_plsql_clause is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_plsql_clause(ctx *Supplemental_plsql_clauseContext) {
}

// EnterLogfile_descriptor is called when production logfile_descriptor is entered.
func (s *BasePlSqlParserListener) EnterLogfile_descriptor(ctx *Logfile_descriptorContext) {}

// ExitLogfile_descriptor is called when production logfile_descriptor is exited.
func (s *BasePlSqlParserListener) ExitLogfile_descriptor(ctx *Logfile_descriptorContext) {}

// EnterControlfile_clauses is called when production controlfile_clauses is entered.
func (s *BasePlSqlParserListener) EnterControlfile_clauses(ctx *Controlfile_clausesContext) {}

// ExitControlfile_clauses is called when production controlfile_clauses is exited.
func (s *BasePlSqlParserListener) ExitControlfile_clauses(ctx *Controlfile_clausesContext) {}

// EnterTrace_file_clause is called when production trace_file_clause is entered.
func (s *BasePlSqlParserListener) EnterTrace_file_clause(ctx *Trace_file_clauseContext) {}

// ExitTrace_file_clause is called when production trace_file_clause is exited.
func (s *BasePlSqlParserListener) ExitTrace_file_clause(ctx *Trace_file_clauseContext) {}

// EnterStandby_database_clauses is called when production standby_database_clauses is entered.
func (s *BasePlSqlParserListener) EnterStandby_database_clauses(ctx *Standby_database_clausesContext) {
}

// ExitStandby_database_clauses is called when production standby_database_clauses is exited.
func (s *BasePlSqlParserListener) ExitStandby_database_clauses(ctx *Standby_database_clausesContext) {
}

// EnterActivate_standby_db_clause is called when production activate_standby_db_clause is entered.
func (s *BasePlSqlParserListener) EnterActivate_standby_db_clause(ctx *Activate_standby_db_clauseContext) {
}

// ExitActivate_standby_db_clause is called when production activate_standby_db_clause is exited.
func (s *BasePlSqlParserListener) ExitActivate_standby_db_clause(ctx *Activate_standby_db_clauseContext) {
}

// EnterMaximize_standby_db_clause is called when production maximize_standby_db_clause is entered.
func (s *BasePlSqlParserListener) EnterMaximize_standby_db_clause(ctx *Maximize_standby_db_clauseContext) {
}

// ExitMaximize_standby_db_clause is called when production maximize_standby_db_clause is exited.
func (s *BasePlSqlParserListener) ExitMaximize_standby_db_clause(ctx *Maximize_standby_db_clauseContext) {
}

// EnterRegister_logfile_clause is called when production register_logfile_clause is entered.
func (s *BasePlSqlParserListener) EnterRegister_logfile_clause(ctx *Register_logfile_clauseContext) {}

// ExitRegister_logfile_clause is called when production register_logfile_clause is exited.
func (s *BasePlSqlParserListener) ExitRegister_logfile_clause(ctx *Register_logfile_clauseContext) {}

// EnterCommit_switchover_clause is called when production commit_switchover_clause is entered.
func (s *BasePlSqlParserListener) EnterCommit_switchover_clause(ctx *Commit_switchover_clauseContext) {
}

// ExitCommit_switchover_clause is called when production commit_switchover_clause is exited.
func (s *BasePlSqlParserListener) ExitCommit_switchover_clause(ctx *Commit_switchover_clauseContext) {
}

// EnterStart_standby_clause is called when production start_standby_clause is entered.
func (s *BasePlSqlParserListener) EnterStart_standby_clause(ctx *Start_standby_clauseContext) {}

// ExitStart_standby_clause is called when production start_standby_clause is exited.
func (s *BasePlSqlParserListener) ExitStart_standby_clause(ctx *Start_standby_clauseContext) {}

// EnterStop_standby_clause is called when production stop_standby_clause is entered.
func (s *BasePlSqlParserListener) EnterStop_standby_clause(ctx *Stop_standby_clauseContext) {}

// ExitStop_standby_clause is called when production stop_standby_clause is exited.
func (s *BasePlSqlParserListener) ExitStop_standby_clause(ctx *Stop_standby_clauseContext) {}

// EnterConvert_database_clause is called when production convert_database_clause is entered.
func (s *BasePlSqlParserListener) EnterConvert_database_clause(ctx *Convert_database_clauseContext) {}

// ExitConvert_database_clause is called when production convert_database_clause is exited.
func (s *BasePlSqlParserListener) ExitConvert_database_clause(ctx *Convert_database_clauseContext) {}

// EnterDefault_settings_clause is called when production default_settings_clause is entered.
func (s *BasePlSqlParserListener) EnterDefault_settings_clause(ctx *Default_settings_clauseContext) {}

// ExitDefault_settings_clause is called when production default_settings_clause is exited.
func (s *BasePlSqlParserListener) ExitDefault_settings_clause(ctx *Default_settings_clauseContext) {}

// EnterSet_time_zone_clause is called when production set_time_zone_clause is entered.
func (s *BasePlSqlParserListener) EnterSet_time_zone_clause(ctx *Set_time_zone_clauseContext) {}

// ExitSet_time_zone_clause is called when production set_time_zone_clause is exited.
func (s *BasePlSqlParserListener) ExitSet_time_zone_clause(ctx *Set_time_zone_clauseContext) {}

// EnterInstance_clauses is called when production instance_clauses is entered.
func (s *BasePlSqlParserListener) EnterInstance_clauses(ctx *Instance_clausesContext) {}

// ExitInstance_clauses is called when production instance_clauses is exited.
func (s *BasePlSqlParserListener) ExitInstance_clauses(ctx *Instance_clausesContext) {}

// EnterSecurity_clause is called when production security_clause is entered.
func (s *BasePlSqlParserListener) EnterSecurity_clause(ctx *Security_clauseContext) {}

// ExitSecurity_clause is called when production security_clause is exited.
func (s *BasePlSqlParserListener) ExitSecurity_clause(ctx *Security_clauseContext) {}

// EnterDomain is called when production domain is entered.
func (s *BasePlSqlParserListener) EnterDomain(ctx *DomainContext) {}

// ExitDomain is called when production domain is exited.
func (s *BasePlSqlParserListener) ExitDomain(ctx *DomainContext) {}

// EnterDatabase is called when production database is entered.
func (s *BasePlSqlParserListener) EnterDatabase(ctx *DatabaseContext) {}

// ExitDatabase is called when production database is exited.
func (s *BasePlSqlParserListener) ExitDatabase(ctx *DatabaseContext) {}

// EnterEdition_name is called when production edition_name is entered.
func (s *BasePlSqlParserListener) EnterEdition_name(ctx *Edition_nameContext) {}

// ExitEdition_name is called when production edition_name is exited.
func (s *BasePlSqlParserListener) ExitEdition_name(ctx *Edition_nameContext) {}

// EnterFilenumber is called when production filenumber is entered.
func (s *BasePlSqlParserListener) EnterFilenumber(ctx *FilenumberContext) {}

// ExitFilenumber is called when production filenumber is exited.
func (s *BasePlSqlParserListener) ExitFilenumber(ctx *FilenumberContext) {}

// EnterFilename is called when production filename is entered.
func (s *BasePlSqlParserListener) EnterFilename(ctx *FilenameContext) {}

// ExitFilename is called when production filename is exited.
func (s *BasePlSqlParserListener) ExitFilename(ctx *FilenameContext) {}

// EnterPrepare_clause is called when production prepare_clause is entered.
func (s *BasePlSqlParserListener) EnterPrepare_clause(ctx *Prepare_clauseContext) {}

// ExitPrepare_clause is called when production prepare_clause is exited.
func (s *BasePlSqlParserListener) ExitPrepare_clause(ctx *Prepare_clauseContext) {}

// EnterDrop_mirror_clause is called when production drop_mirror_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_mirror_clause(ctx *Drop_mirror_clauseContext) {}

// ExitDrop_mirror_clause is called when production drop_mirror_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_mirror_clause(ctx *Drop_mirror_clauseContext) {}

// EnterLost_write_protection is called when production lost_write_protection is entered.
func (s *BasePlSqlParserListener) EnterLost_write_protection(ctx *Lost_write_protectionContext) {}

// ExitLost_write_protection is called when production lost_write_protection is exited.
func (s *BasePlSqlParserListener) ExitLost_write_protection(ctx *Lost_write_protectionContext) {}

// EnterCdb_fleet_clauses is called when production cdb_fleet_clauses is entered.
func (s *BasePlSqlParserListener) EnterCdb_fleet_clauses(ctx *Cdb_fleet_clausesContext) {}

// ExitCdb_fleet_clauses is called when production cdb_fleet_clauses is exited.
func (s *BasePlSqlParserListener) ExitCdb_fleet_clauses(ctx *Cdb_fleet_clausesContext) {}

// EnterLead_cdb_clause is called when production lead_cdb_clause is entered.
func (s *BasePlSqlParserListener) EnterLead_cdb_clause(ctx *Lead_cdb_clauseContext) {}

// ExitLead_cdb_clause is called when production lead_cdb_clause is exited.
func (s *BasePlSqlParserListener) ExitLead_cdb_clause(ctx *Lead_cdb_clauseContext) {}

// EnterLead_cdb_uri_clause is called when production lead_cdb_uri_clause is entered.
func (s *BasePlSqlParserListener) EnterLead_cdb_uri_clause(ctx *Lead_cdb_uri_clauseContext) {}

// ExitLead_cdb_uri_clause is called when production lead_cdb_uri_clause is exited.
func (s *BasePlSqlParserListener) ExitLead_cdb_uri_clause(ctx *Lead_cdb_uri_clauseContext) {}

// EnterProperty_clauses is called when production property_clauses is entered.
func (s *BasePlSqlParserListener) EnterProperty_clauses(ctx *Property_clausesContext) {}

// ExitProperty_clauses is called when production property_clauses is exited.
func (s *BasePlSqlParserListener) ExitProperty_clauses(ctx *Property_clausesContext) {}

// EnterReplay_upgrade_clauses is called when production replay_upgrade_clauses is entered.
func (s *BasePlSqlParserListener) EnterReplay_upgrade_clauses(ctx *Replay_upgrade_clausesContext) {}

// ExitReplay_upgrade_clauses is called when production replay_upgrade_clauses is exited.
func (s *BasePlSqlParserListener) ExitReplay_upgrade_clauses(ctx *Replay_upgrade_clausesContext) {}

// EnterAlter_database_link is called when production alter_database_link is entered.
func (s *BasePlSqlParserListener) EnterAlter_database_link(ctx *Alter_database_linkContext) {}

// ExitAlter_database_link is called when production alter_database_link is exited.
func (s *BasePlSqlParserListener) ExitAlter_database_link(ctx *Alter_database_linkContext) {}

// EnterPassword_value is called when production password_value is entered.
func (s *BasePlSqlParserListener) EnterPassword_value(ctx *Password_valueContext) {}

// ExitPassword_value is called when production password_value is exited.
func (s *BasePlSqlParserListener) ExitPassword_value(ctx *Password_valueContext) {}

// EnterLink_authentication is called when production link_authentication is entered.
func (s *BasePlSqlParserListener) EnterLink_authentication(ctx *Link_authenticationContext) {}

// ExitLink_authentication is called when production link_authentication is exited.
func (s *BasePlSqlParserListener) ExitLink_authentication(ctx *Link_authenticationContext) {}

// EnterCreate_schema is called when production create_schema is entered.
func (s *BasePlSqlParserListener) EnterCreate_schema(ctx *Create_schemaContext) {}

// ExitCreate_schema is called when production create_schema is exited.
func (s *BasePlSqlParserListener) ExitCreate_schema(ctx *Create_schemaContext) {}

// EnterCreate_database is called when production create_database is entered.
func (s *BasePlSqlParserListener) EnterCreate_database(ctx *Create_databaseContext) {}

// ExitCreate_database is called when production create_database is exited.
func (s *BasePlSqlParserListener) ExitCreate_database(ctx *Create_databaseContext) {}

// EnterDatabase_logging_clauses is called when production database_logging_clauses is entered.
func (s *BasePlSqlParserListener) EnterDatabase_logging_clauses(ctx *Database_logging_clausesContext) {
}

// ExitDatabase_logging_clauses is called when production database_logging_clauses is exited.
func (s *BasePlSqlParserListener) ExitDatabase_logging_clauses(ctx *Database_logging_clausesContext) {
}

// EnterDatabase_logging_sub_clause is called when production database_logging_sub_clause is entered.
func (s *BasePlSqlParserListener) EnterDatabase_logging_sub_clause(ctx *Database_logging_sub_clauseContext) {
}

// ExitDatabase_logging_sub_clause is called when production database_logging_sub_clause is exited.
func (s *BasePlSqlParserListener) ExitDatabase_logging_sub_clause(ctx *Database_logging_sub_clauseContext) {
}

// EnterTablespace_clauses is called when production tablespace_clauses is entered.
func (s *BasePlSqlParserListener) EnterTablespace_clauses(ctx *Tablespace_clausesContext) {}

// ExitTablespace_clauses is called when production tablespace_clauses is exited.
func (s *BasePlSqlParserListener) ExitTablespace_clauses(ctx *Tablespace_clausesContext) {}

// EnterEnable_pluggable_database is called when production enable_pluggable_database is entered.
func (s *BasePlSqlParserListener) EnterEnable_pluggable_database(ctx *Enable_pluggable_databaseContext) {
}

// ExitEnable_pluggable_database is called when production enable_pluggable_database is exited.
func (s *BasePlSqlParserListener) ExitEnable_pluggable_database(ctx *Enable_pluggable_databaseContext) {
}

// EnterFile_name_convert is called when production file_name_convert is entered.
func (s *BasePlSqlParserListener) EnterFile_name_convert(ctx *File_name_convertContext) {}

// ExitFile_name_convert is called when production file_name_convert is exited.
func (s *BasePlSqlParserListener) ExitFile_name_convert(ctx *File_name_convertContext) {}

// EnterFilename_convert_sub_clause is called when production filename_convert_sub_clause is entered.
func (s *BasePlSqlParserListener) EnterFilename_convert_sub_clause(ctx *Filename_convert_sub_clauseContext) {
}

// ExitFilename_convert_sub_clause is called when production filename_convert_sub_clause is exited.
func (s *BasePlSqlParserListener) ExitFilename_convert_sub_clause(ctx *Filename_convert_sub_clauseContext) {
}

// EnterTablespace_datafile_clauses is called when production tablespace_datafile_clauses is entered.
func (s *BasePlSqlParserListener) EnterTablespace_datafile_clauses(ctx *Tablespace_datafile_clausesContext) {
}

// ExitTablespace_datafile_clauses is called when production tablespace_datafile_clauses is exited.
func (s *BasePlSqlParserListener) ExitTablespace_datafile_clauses(ctx *Tablespace_datafile_clausesContext) {
}

// EnterUndo_mode_clause is called when production undo_mode_clause is entered.
func (s *BasePlSqlParserListener) EnterUndo_mode_clause(ctx *Undo_mode_clauseContext) {}

// ExitUndo_mode_clause is called when production undo_mode_clause is exited.
func (s *BasePlSqlParserListener) ExitUndo_mode_clause(ctx *Undo_mode_clauseContext) {}

// EnterDefault_tablespace is called when production default_tablespace is entered.
func (s *BasePlSqlParserListener) EnterDefault_tablespace(ctx *Default_tablespaceContext) {}

// ExitDefault_tablespace is called when production default_tablespace is exited.
func (s *BasePlSqlParserListener) ExitDefault_tablespace(ctx *Default_tablespaceContext) {}

// EnterDefault_temp_tablespace is called when production default_temp_tablespace is entered.
func (s *BasePlSqlParserListener) EnterDefault_temp_tablespace(ctx *Default_temp_tablespaceContext) {}

// ExitDefault_temp_tablespace is called when production default_temp_tablespace is exited.
func (s *BasePlSqlParserListener) ExitDefault_temp_tablespace(ctx *Default_temp_tablespaceContext) {}

// EnterUndo_tablespace is called when production undo_tablespace is entered.
func (s *BasePlSqlParserListener) EnterUndo_tablespace(ctx *Undo_tablespaceContext) {}

// ExitUndo_tablespace is called when production undo_tablespace is exited.
func (s *BasePlSqlParserListener) ExitUndo_tablespace(ctx *Undo_tablespaceContext) {}

// EnterDrop_database is called when production drop_database is entered.
func (s *BasePlSqlParserListener) EnterDrop_database(ctx *Drop_databaseContext) {}

// ExitDrop_database is called when production drop_database is exited.
func (s *BasePlSqlParserListener) ExitDrop_database(ctx *Drop_databaseContext) {}

// EnterCreate_database_link is called when production create_database_link is entered.
func (s *BasePlSqlParserListener) EnterCreate_database_link(ctx *Create_database_linkContext) {}

// ExitCreate_database_link is called when production create_database_link is exited.
func (s *BasePlSqlParserListener) ExitCreate_database_link(ctx *Create_database_linkContext) {}

// EnterDrop_database_link is called when production drop_database_link is entered.
func (s *BasePlSqlParserListener) EnterDrop_database_link(ctx *Drop_database_linkContext) {}

// ExitDrop_database_link is called when production drop_database_link is exited.
func (s *BasePlSqlParserListener) ExitDrop_database_link(ctx *Drop_database_linkContext) {}

// EnterAlter_tablespace_set is called when production alter_tablespace_set is entered.
func (s *BasePlSqlParserListener) EnterAlter_tablespace_set(ctx *Alter_tablespace_setContext) {}

// ExitAlter_tablespace_set is called when production alter_tablespace_set is exited.
func (s *BasePlSqlParserListener) ExitAlter_tablespace_set(ctx *Alter_tablespace_setContext) {}

// EnterAlter_tablespace_attrs is called when production alter_tablespace_attrs is entered.
func (s *BasePlSqlParserListener) EnterAlter_tablespace_attrs(ctx *Alter_tablespace_attrsContext) {}

// ExitAlter_tablespace_attrs is called when production alter_tablespace_attrs is exited.
func (s *BasePlSqlParserListener) ExitAlter_tablespace_attrs(ctx *Alter_tablespace_attrsContext) {}

// EnterAlter_tablespace_encryption is called when production alter_tablespace_encryption is entered.
func (s *BasePlSqlParserListener) EnterAlter_tablespace_encryption(ctx *Alter_tablespace_encryptionContext) {
}

// ExitAlter_tablespace_encryption is called when production alter_tablespace_encryption is exited.
func (s *BasePlSqlParserListener) ExitAlter_tablespace_encryption(ctx *Alter_tablespace_encryptionContext) {
}

// EnterTs_file_name_convert is called when production ts_file_name_convert is entered.
func (s *BasePlSqlParserListener) EnterTs_file_name_convert(ctx *Ts_file_name_convertContext) {}

// ExitTs_file_name_convert is called when production ts_file_name_convert is exited.
func (s *BasePlSqlParserListener) ExitTs_file_name_convert(ctx *Ts_file_name_convertContext) {}

// EnterAlter_role is called when production alter_role is entered.
func (s *BasePlSqlParserListener) EnterAlter_role(ctx *Alter_roleContext) {}

// ExitAlter_role is called when production alter_role is exited.
func (s *BasePlSqlParserListener) ExitAlter_role(ctx *Alter_roleContext) {}

// EnterRole_identified_clause is called when production role_identified_clause is entered.
func (s *BasePlSqlParserListener) EnterRole_identified_clause(ctx *Role_identified_clauseContext) {}

// ExitRole_identified_clause is called when production role_identified_clause is exited.
func (s *BasePlSqlParserListener) ExitRole_identified_clause(ctx *Role_identified_clauseContext) {}

// EnterAlter_table is called when production alter_table is entered.
func (s *BasePlSqlParserListener) EnterAlter_table(ctx *Alter_tableContext) {}

// ExitAlter_table is called when production alter_table is exited.
func (s *BasePlSqlParserListener) ExitAlter_table(ctx *Alter_tableContext) {}

// EnterMemoptimize_read_write_clause is called when production memoptimize_read_write_clause is entered.
func (s *BasePlSqlParserListener) EnterMemoptimize_read_write_clause(ctx *Memoptimize_read_write_clauseContext) {
}

// ExitMemoptimize_read_write_clause is called when production memoptimize_read_write_clause is exited.
func (s *BasePlSqlParserListener) ExitMemoptimize_read_write_clause(ctx *Memoptimize_read_write_clauseContext) {
}

// EnterAlter_table_properties is called when production alter_table_properties is entered.
func (s *BasePlSqlParserListener) EnterAlter_table_properties(ctx *Alter_table_propertiesContext) {}

// ExitAlter_table_properties is called when production alter_table_properties is exited.
func (s *BasePlSqlParserListener) ExitAlter_table_properties(ctx *Alter_table_propertiesContext) {}

// EnterAlter_table_partitioning is called when production alter_table_partitioning is entered.
func (s *BasePlSqlParserListener) EnterAlter_table_partitioning(ctx *Alter_table_partitioningContext) {
}

// ExitAlter_table_partitioning is called when production alter_table_partitioning is exited.
func (s *BasePlSqlParserListener) ExitAlter_table_partitioning(ctx *Alter_table_partitioningContext) {
}

// EnterAdd_table_partition is called when production add_table_partition is entered.
func (s *BasePlSqlParserListener) EnterAdd_table_partition(ctx *Add_table_partitionContext) {}

// ExitAdd_table_partition is called when production add_table_partition is exited.
func (s *BasePlSqlParserListener) ExitAdd_table_partition(ctx *Add_table_partitionContext) {}

// EnterDrop_table_partition is called when production drop_table_partition is entered.
func (s *BasePlSqlParserListener) EnterDrop_table_partition(ctx *Drop_table_partitionContext) {}

// ExitDrop_table_partition is called when production drop_table_partition is exited.
func (s *BasePlSqlParserListener) ExitDrop_table_partition(ctx *Drop_table_partitionContext) {}

// EnterMerge_table_partition is called when production merge_table_partition is entered.
func (s *BasePlSqlParserListener) EnterMerge_table_partition(ctx *Merge_table_partitionContext) {}

// ExitMerge_table_partition is called when production merge_table_partition is exited.
func (s *BasePlSqlParserListener) ExitMerge_table_partition(ctx *Merge_table_partitionContext) {}

// EnterModify_table_partition is called when production modify_table_partition is entered.
func (s *BasePlSqlParserListener) EnterModify_table_partition(ctx *Modify_table_partitionContext) {}

// ExitModify_table_partition is called when production modify_table_partition is exited.
func (s *BasePlSqlParserListener) ExitModify_table_partition(ctx *Modify_table_partitionContext) {}

// EnterSplit_table_partition is called when production split_table_partition is entered.
func (s *BasePlSqlParserListener) EnterSplit_table_partition(ctx *Split_table_partitionContext) {}

// ExitSplit_table_partition is called when production split_table_partition is exited.
func (s *BasePlSqlParserListener) ExitSplit_table_partition(ctx *Split_table_partitionContext) {}

// EnterTruncate_table_partition is called when production truncate_table_partition is entered.
func (s *BasePlSqlParserListener) EnterTruncate_table_partition(ctx *Truncate_table_partitionContext) {
}

// ExitTruncate_table_partition is called when production truncate_table_partition is exited.
func (s *BasePlSqlParserListener) ExitTruncate_table_partition(ctx *Truncate_table_partitionContext) {
}

// EnterExchange_table_partition is called when production exchange_table_partition is entered.
func (s *BasePlSqlParserListener) EnterExchange_table_partition(ctx *Exchange_table_partitionContext) {
}

// ExitExchange_table_partition is called when production exchange_table_partition is exited.
func (s *BasePlSqlParserListener) ExitExchange_table_partition(ctx *Exchange_table_partitionContext) {
}

// EnterCoalesce_table_partition is called when production coalesce_table_partition is entered.
func (s *BasePlSqlParserListener) EnterCoalesce_table_partition(ctx *Coalesce_table_partitionContext) {
}

// ExitCoalesce_table_partition is called when production coalesce_table_partition is exited.
func (s *BasePlSqlParserListener) ExitCoalesce_table_partition(ctx *Coalesce_table_partitionContext) {
}

// EnterAlter_interval_partition is called when production alter_interval_partition is entered.
func (s *BasePlSqlParserListener) EnterAlter_interval_partition(ctx *Alter_interval_partitionContext) {
}

// ExitAlter_interval_partition is called when production alter_interval_partition is exited.
func (s *BasePlSqlParserListener) ExitAlter_interval_partition(ctx *Alter_interval_partitionContext) {
}

// EnterMove_table_partition is called when production move_table_partition is entered.
func (s *BasePlSqlParserListener) EnterMove_table_partition(ctx *Move_table_partitionContext) {}

// ExitMove_table_partition is called when production move_table_partition is exited.
func (s *BasePlSqlParserListener) ExitMove_table_partition(ctx *Move_table_partitionContext) {}

// EnterFilter_condition is called when production filter_condition is entered.
func (s *BasePlSqlParserListener) EnterFilter_condition(ctx *Filter_conditionContext) {}

// ExitFilter_condition is called when production filter_condition is exited.
func (s *BasePlSqlParserListener) ExitFilter_condition(ctx *Filter_conditionContext) {}

// EnterRename_table_partition is called when production rename_table_partition is entered.
func (s *BasePlSqlParserListener) EnterRename_table_partition(ctx *Rename_table_partitionContext) {}

// ExitRename_table_partition is called when production rename_table_partition is exited.
func (s *BasePlSqlParserListener) ExitRename_table_partition(ctx *Rename_table_partitionContext) {}

// EnterPartition_extended_names is called when production partition_extended_names is entered.
func (s *BasePlSqlParserListener) EnterPartition_extended_names(ctx *Partition_extended_namesContext) {
}

// ExitPartition_extended_names is called when production partition_extended_names is exited.
func (s *BasePlSqlParserListener) ExitPartition_extended_names(ctx *Partition_extended_namesContext) {
}

// EnterSubpartition_extended_names is called when production subpartition_extended_names is entered.
func (s *BasePlSqlParserListener) EnterSubpartition_extended_names(ctx *Subpartition_extended_namesContext) {
}

// ExitSubpartition_extended_names is called when production subpartition_extended_names is exited.
func (s *BasePlSqlParserListener) ExitSubpartition_extended_names(ctx *Subpartition_extended_namesContext) {
}

// EnterAlter_table_properties_1 is called when production alter_table_properties_1 is entered.
func (s *BasePlSqlParserListener) EnterAlter_table_properties_1(ctx *Alter_table_properties_1Context) {
}

// ExitAlter_table_properties_1 is called when production alter_table_properties_1 is exited.
func (s *BasePlSqlParserListener) ExitAlter_table_properties_1(ctx *Alter_table_properties_1Context) {
}

// EnterAlter_iot_clauses is called when production alter_iot_clauses is entered.
func (s *BasePlSqlParserListener) EnterAlter_iot_clauses(ctx *Alter_iot_clausesContext) {}

// ExitAlter_iot_clauses is called when production alter_iot_clauses is exited.
func (s *BasePlSqlParserListener) ExitAlter_iot_clauses(ctx *Alter_iot_clausesContext) {}

// EnterAlter_mapping_table_clause is called when production alter_mapping_table_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_mapping_table_clause(ctx *Alter_mapping_table_clauseContext) {
}

// ExitAlter_mapping_table_clause is called when production alter_mapping_table_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_mapping_table_clause(ctx *Alter_mapping_table_clauseContext) {
}

// EnterAlter_overflow_clause is called when production alter_overflow_clause is entered.
func (s *BasePlSqlParserListener) EnterAlter_overflow_clause(ctx *Alter_overflow_clauseContext) {}

// ExitAlter_overflow_clause is called when production alter_overflow_clause is exited.
func (s *BasePlSqlParserListener) ExitAlter_overflow_clause(ctx *Alter_overflow_clauseContext) {}

// EnterAdd_overflow_clause is called when production add_overflow_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_overflow_clause(ctx *Add_overflow_clauseContext) {}

// ExitAdd_overflow_clause is called when production add_overflow_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_overflow_clause(ctx *Add_overflow_clauseContext) {}

// EnterUpdate_index_clauses is called when production update_index_clauses is entered.
func (s *BasePlSqlParserListener) EnterUpdate_index_clauses(ctx *Update_index_clausesContext) {}

// ExitUpdate_index_clauses is called when production update_index_clauses is exited.
func (s *BasePlSqlParserListener) ExitUpdate_index_clauses(ctx *Update_index_clausesContext) {}

// EnterUpdate_global_index_clause is called when production update_global_index_clause is entered.
func (s *BasePlSqlParserListener) EnterUpdate_global_index_clause(ctx *Update_global_index_clauseContext) {
}

// ExitUpdate_global_index_clause is called when production update_global_index_clause is exited.
func (s *BasePlSqlParserListener) ExitUpdate_global_index_clause(ctx *Update_global_index_clauseContext) {
}

// EnterUpdate_all_indexes_clause is called when production update_all_indexes_clause is entered.
func (s *BasePlSqlParserListener) EnterUpdate_all_indexes_clause(ctx *Update_all_indexes_clauseContext) {
}

// ExitUpdate_all_indexes_clause is called when production update_all_indexes_clause is exited.
func (s *BasePlSqlParserListener) ExitUpdate_all_indexes_clause(ctx *Update_all_indexes_clauseContext) {
}

// EnterUpdate_all_indexes_index_clause is called when production update_all_indexes_index_clause is entered.
func (s *BasePlSqlParserListener) EnterUpdate_all_indexes_index_clause(ctx *Update_all_indexes_index_clauseContext) {
}

// ExitUpdate_all_indexes_index_clause is called when production update_all_indexes_index_clause is exited.
func (s *BasePlSqlParserListener) ExitUpdate_all_indexes_index_clause(ctx *Update_all_indexes_index_clauseContext) {
}

// EnterUpdate_index_partition is called when production update_index_partition is entered.
func (s *BasePlSqlParserListener) EnterUpdate_index_partition(ctx *Update_index_partitionContext) {}

// ExitUpdate_index_partition is called when production update_index_partition is exited.
func (s *BasePlSqlParserListener) ExitUpdate_index_partition(ctx *Update_index_partitionContext) {}

// EnterUpdate_index_subpartition is called when production update_index_subpartition is entered.
func (s *BasePlSqlParserListener) EnterUpdate_index_subpartition(ctx *Update_index_subpartitionContext) {
}

// ExitUpdate_index_subpartition is called when production update_index_subpartition is exited.
func (s *BasePlSqlParserListener) ExitUpdate_index_subpartition(ctx *Update_index_subpartitionContext) {
}

// EnterEnable_disable_clause is called when production enable_disable_clause is entered.
func (s *BasePlSqlParserListener) EnterEnable_disable_clause(ctx *Enable_disable_clauseContext) {}

// ExitEnable_disable_clause is called when production enable_disable_clause is exited.
func (s *BasePlSqlParserListener) ExitEnable_disable_clause(ctx *Enable_disable_clauseContext) {}

// EnterUsing_index_clause is called when production using_index_clause is entered.
func (s *BasePlSqlParserListener) EnterUsing_index_clause(ctx *Using_index_clauseContext) {}

// ExitUsing_index_clause is called when production using_index_clause is exited.
func (s *BasePlSqlParserListener) ExitUsing_index_clause(ctx *Using_index_clauseContext) {}

// EnterIndex_attributes is called when production index_attributes is entered.
func (s *BasePlSqlParserListener) EnterIndex_attributes(ctx *Index_attributesContext) {}

// ExitIndex_attributes is called when production index_attributes is exited.
func (s *BasePlSqlParserListener) ExitIndex_attributes(ctx *Index_attributesContext) {}

// EnterSort_or_nosort is called when production sort_or_nosort is entered.
func (s *BasePlSqlParserListener) EnterSort_or_nosort(ctx *Sort_or_nosortContext) {}

// ExitSort_or_nosort is called when production sort_or_nosort is exited.
func (s *BasePlSqlParserListener) ExitSort_or_nosort(ctx *Sort_or_nosortContext) {}

// EnterExceptions_clause is called when production exceptions_clause is entered.
func (s *BasePlSqlParserListener) EnterExceptions_clause(ctx *Exceptions_clauseContext) {}

// ExitExceptions_clause is called when production exceptions_clause is exited.
func (s *BasePlSqlParserListener) ExitExceptions_clause(ctx *Exceptions_clauseContext) {}

// EnterMove_table_clause is called when production move_table_clause is entered.
func (s *BasePlSqlParserListener) EnterMove_table_clause(ctx *Move_table_clauseContext) {}

// ExitMove_table_clause is called when production move_table_clause is exited.
func (s *BasePlSqlParserListener) ExitMove_table_clause(ctx *Move_table_clauseContext) {}

// EnterIndex_org_table_clause is called when production index_org_table_clause is entered.
func (s *BasePlSqlParserListener) EnterIndex_org_table_clause(ctx *Index_org_table_clauseContext) {}

// ExitIndex_org_table_clause is called when production index_org_table_clause is exited.
func (s *BasePlSqlParserListener) ExitIndex_org_table_clause(ctx *Index_org_table_clauseContext) {}

// EnterMapping_table_clause is called when production mapping_table_clause is entered.
func (s *BasePlSqlParserListener) EnterMapping_table_clause(ctx *Mapping_table_clauseContext) {}

// ExitMapping_table_clause is called when production mapping_table_clause is exited.
func (s *BasePlSqlParserListener) ExitMapping_table_clause(ctx *Mapping_table_clauseContext) {}

// EnterKey_compression is called when production key_compression is entered.
func (s *BasePlSqlParserListener) EnterKey_compression(ctx *Key_compressionContext) {}

// ExitKey_compression is called when production key_compression is exited.
func (s *BasePlSqlParserListener) ExitKey_compression(ctx *Key_compressionContext) {}

// EnterIndex_org_overflow_clause is called when production index_org_overflow_clause is entered.
func (s *BasePlSqlParserListener) EnterIndex_org_overflow_clause(ctx *Index_org_overflow_clauseContext) {
}

// ExitIndex_org_overflow_clause is called when production index_org_overflow_clause is exited.
func (s *BasePlSqlParserListener) ExitIndex_org_overflow_clause(ctx *Index_org_overflow_clauseContext) {
}

// EnterColumn_clauses is called when production column_clauses is entered.
func (s *BasePlSqlParserListener) EnterColumn_clauses(ctx *Column_clausesContext) {}

// ExitColumn_clauses is called when production column_clauses is exited.
func (s *BasePlSqlParserListener) ExitColumn_clauses(ctx *Column_clausesContext) {}

// EnterModify_collection_retrieval is called when production modify_collection_retrieval is entered.
func (s *BasePlSqlParserListener) EnterModify_collection_retrieval(ctx *Modify_collection_retrievalContext) {
}

// ExitModify_collection_retrieval is called when production modify_collection_retrieval is exited.
func (s *BasePlSqlParserListener) ExitModify_collection_retrieval(ctx *Modify_collection_retrievalContext) {
}

// EnterCollection_item is called when production collection_item is entered.
func (s *BasePlSqlParserListener) EnterCollection_item(ctx *Collection_itemContext) {}

// ExitCollection_item is called when production collection_item is exited.
func (s *BasePlSqlParserListener) ExitCollection_item(ctx *Collection_itemContext) {}

// EnterRename_column_clause is called when production rename_column_clause is entered.
func (s *BasePlSqlParserListener) EnterRename_column_clause(ctx *Rename_column_clauseContext) {}

// ExitRename_column_clause is called when production rename_column_clause is exited.
func (s *BasePlSqlParserListener) ExitRename_column_clause(ctx *Rename_column_clauseContext) {}

// EnterOld_column_name is called when production old_column_name is entered.
func (s *BasePlSqlParserListener) EnterOld_column_name(ctx *Old_column_nameContext) {}

// ExitOld_column_name is called when production old_column_name is exited.
func (s *BasePlSqlParserListener) ExitOld_column_name(ctx *Old_column_nameContext) {}

// EnterNew_column_name is called when production new_column_name is entered.
func (s *BasePlSqlParserListener) EnterNew_column_name(ctx *New_column_nameContext) {}

// ExitNew_column_name is called when production new_column_name is exited.
func (s *BasePlSqlParserListener) ExitNew_column_name(ctx *New_column_nameContext) {}

// EnterAdd_modify_drop_column_clauses is called when production add_modify_drop_column_clauses is entered.
func (s *BasePlSqlParserListener) EnterAdd_modify_drop_column_clauses(ctx *Add_modify_drop_column_clausesContext) {
}

// ExitAdd_modify_drop_column_clauses is called when production add_modify_drop_column_clauses is exited.
func (s *BasePlSqlParserListener) ExitAdd_modify_drop_column_clauses(ctx *Add_modify_drop_column_clausesContext) {
}

// EnterDrop_column_clause is called when production drop_column_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_column_clause(ctx *Drop_column_clauseContext) {}

// ExitDrop_column_clause is called when production drop_column_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_column_clause(ctx *Drop_column_clauseContext) {}

// EnterModify_column_clauses is called when production modify_column_clauses is entered.
func (s *BasePlSqlParserListener) EnterModify_column_clauses(ctx *Modify_column_clausesContext) {}

// ExitModify_column_clauses is called when production modify_column_clauses is exited.
func (s *BasePlSqlParserListener) ExitModify_column_clauses(ctx *Modify_column_clausesContext) {}

// EnterModify_col_properties is called when production modify_col_properties is entered.
func (s *BasePlSqlParserListener) EnterModify_col_properties(ctx *Modify_col_propertiesContext) {}

// ExitModify_col_properties is called when production modify_col_properties is exited.
func (s *BasePlSqlParserListener) ExitModify_col_properties(ctx *Modify_col_propertiesContext) {}

// EnterModify_col_visibility is called when production modify_col_visibility is entered.
func (s *BasePlSqlParserListener) EnterModify_col_visibility(ctx *Modify_col_visibilityContext) {}

// ExitModify_col_visibility is called when production modify_col_visibility is exited.
func (s *BasePlSqlParserListener) ExitModify_col_visibility(ctx *Modify_col_visibilityContext) {}

// EnterModify_col_substitutable is called when production modify_col_substitutable is entered.
func (s *BasePlSqlParserListener) EnterModify_col_substitutable(ctx *Modify_col_substitutableContext) {
}

// ExitModify_col_substitutable is called when production modify_col_substitutable is exited.
func (s *BasePlSqlParserListener) ExitModify_col_substitutable(ctx *Modify_col_substitutableContext) {
}

// EnterAdd_column_clause is called when production add_column_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_column_clause(ctx *Add_column_clauseContext) {}

// ExitAdd_column_clause is called when production add_column_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_column_clause(ctx *Add_column_clauseContext) {}

// EnterVarray_col_properties is called when production varray_col_properties is entered.
func (s *BasePlSqlParserListener) EnterVarray_col_properties(ctx *Varray_col_propertiesContext) {}

// ExitVarray_col_properties is called when production varray_col_properties is exited.
func (s *BasePlSqlParserListener) ExitVarray_col_properties(ctx *Varray_col_propertiesContext) {}

// EnterVarray_storage_clause is called when production varray_storage_clause is entered.
func (s *BasePlSqlParserListener) EnterVarray_storage_clause(ctx *Varray_storage_clauseContext) {}

// ExitVarray_storage_clause is called when production varray_storage_clause is exited.
func (s *BasePlSqlParserListener) ExitVarray_storage_clause(ctx *Varray_storage_clauseContext) {}

// EnterLob_segname is called when production lob_segname is entered.
func (s *BasePlSqlParserListener) EnterLob_segname(ctx *Lob_segnameContext) {}

// ExitLob_segname is called when production lob_segname is exited.
func (s *BasePlSqlParserListener) ExitLob_segname(ctx *Lob_segnameContext) {}

// EnterLob_item is called when production lob_item is entered.
func (s *BasePlSqlParserListener) EnterLob_item(ctx *Lob_itemContext) {}

// ExitLob_item is called when production lob_item is exited.
func (s *BasePlSqlParserListener) ExitLob_item(ctx *Lob_itemContext) {}

// EnterLob_storage_parameters is called when production lob_storage_parameters is entered.
func (s *BasePlSqlParserListener) EnterLob_storage_parameters(ctx *Lob_storage_parametersContext) {}

// ExitLob_storage_parameters is called when production lob_storage_parameters is exited.
func (s *BasePlSqlParserListener) ExitLob_storage_parameters(ctx *Lob_storage_parametersContext) {}

// EnterLob_storage_clause is called when production lob_storage_clause is entered.
func (s *BasePlSqlParserListener) EnterLob_storage_clause(ctx *Lob_storage_clauseContext) {}

// ExitLob_storage_clause is called when production lob_storage_clause is exited.
func (s *BasePlSqlParserListener) ExitLob_storage_clause(ctx *Lob_storage_clauseContext) {}

// EnterModify_lob_storage_clause is called when production modify_lob_storage_clause is entered.
func (s *BasePlSqlParserListener) EnterModify_lob_storage_clause(ctx *Modify_lob_storage_clauseContext) {
}

// ExitModify_lob_storage_clause is called when production modify_lob_storage_clause is exited.
func (s *BasePlSqlParserListener) ExitModify_lob_storage_clause(ctx *Modify_lob_storage_clauseContext) {
}

// EnterModify_lob_parameters is called when production modify_lob_parameters is entered.
func (s *BasePlSqlParserListener) EnterModify_lob_parameters(ctx *Modify_lob_parametersContext) {}

// ExitModify_lob_parameters is called when production modify_lob_parameters is exited.
func (s *BasePlSqlParserListener) ExitModify_lob_parameters(ctx *Modify_lob_parametersContext) {}

// EnterLob_parameters is called when production lob_parameters is entered.
func (s *BasePlSqlParserListener) EnterLob_parameters(ctx *Lob_parametersContext) {}

// ExitLob_parameters is called when production lob_parameters is exited.
func (s *BasePlSqlParserListener) ExitLob_parameters(ctx *Lob_parametersContext) {}

// EnterLob_deduplicate_clause is called when production lob_deduplicate_clause is entered.
func (s *BasePlSqlParserListener) EnterLob_deduplicate_clause(ctx *Lob_deduplicate_clauseContext) {}

// ExitLob_deduplicate_clause is called when production lob_deduplicate_clause is exited.
func (s *BasePlSqlParserListener) ExitLob_deduplicate_clause(ctx *Lob_deduplicate_clauseContext) {}

// EnterLob_compression_clause is called when production lob_compression_clause is entered.
func (s *BasePlSqlParserListener) EnterLob_compression_clause(ctx *Lob_compression_clauseContext) {}

// ExitLob_compression_clause is called when production lob_compression_clause is exited.
func (s *BasePlSqlParserListener) ExitLob_compression_clause(ctx *Lob_compression_clauseContext) {}

// EnterLob_retention_clause is called when production lob_retention_clause is entered.
func (s *BasePlSqlParserListener) EnterLob_retention_clause(ctx *Lob_retention_clauseContext) {}

// ExitLob_retention_clause is called when production lob_retention_clause is exited.
func (s *BasePlSqlParserListener) ExitLob_retention_clause(ctx *Lob_retention_clauseContext) {}

// EnterEncryption_spec is called when production encryption_spec is entered.
func (s *BasePlSqlParserListener) EnterEncryption_spec(ctx *Encryption_specContext) {}

// ExitEncryption_spec is called when production encryption_spec is exited.
func (s *BasePlSqlParserListener) ExitEncryption_spec(ctx *Encryption_specContext) {}

// EnterTablespace is called when production tablespace is entered.
func (s *BasePlSqlParserListener) EnterTablespace(ctx *TablespaceContext) {}

// ExitTablespace is called when production tablespace is exited.
func (s *BasePlSqlParserListener) ExitTablespace(ctx *TablespaceContext) {}

// EnterVarray_item is called when production varray_item is entered.
func (s *BasePlSqlParserListener) EnterVarray_item(ctx *Varray_itemContext) {}

// ExitVarray_item is called when production varray_item is exited.
func (s *BasePlSqlParserListener) ExitVarray_item(ctx *Varray_itemContext) {}

// EnterColumn_properties is called when production column_properties is entered.
func (s *BasePlSqlParserListener) EnterColumn_properties(ctx *Column_propertiesContext) {}

// ExitColumn_properties is called when production column_properties is exited.
func (s *BasePlSqlParserListener) ExitColumn_properties(ctx *Column_propertiesContext) {}

// EnterLob_partition_storage is called when production lob_partition_storage is entered.
func (s *BasePlSqlParserListener) EnterLob_partition_storage(ctx *Lob_partition_storageContext) {}

// ExitLob_partition_storage is called when production lob_partition_storage is exited.
func (s *BasePlSqlParserListener) ExitLob_partition_storage(ctx *Lob_partition_storageContext) {}

// EnterPeriod_definition is called when production period_definition is entered.
func (s *BasePlSqlParserListener) EnterPeriod_definition(ctx *Period_definitionContext) {}

// ExitPeriod_definition is called when production period_definition is exited.
func (s *BasePlSqlParserListener) ExitPeriod_definition(ctx *Period_definitionContext) {}

// EnterStart_time_column is called when production start_time_column is entered.
func (s *BasePlSqlParserListener) EnterStart_time_column(ctx *Start_time_columnContext) {}

// ExitStart_time_column is called when production start_time_column is exited.
func (s *BasePlSqlParserListener) ExitStart_time_column(ctx *Start_time_columnContext) {}

// EnterEnd_time_column is called when production end_time_column is entered.
func (s *BasePlSqlParserListener) EnterEnd_time_column(ctx *End_time_columnContext) {}

// ExitEnd_time_column is called when production end_time_column is exited.
func (s *BasePlSqlParserListener) ExitEnd_time_column(ctx *End_time_columnContext) {}

// EnterColumn_definition is called when production column_definition is entered.
func (s *BasePlSqlParserListener) EnterColumn_definition(ctx *Column_definitionContext) {}

// ExitColumn_definition is called when production column_definition is exited.
func (s *BasePlSqlParserListener) ExitColumn_definition(ctx *Column_definitionContext) {}

// EnterColumn_collation_name is called when production column_collation_name is entered.
func (s *BasePlSqlParserListener) EnterColumn_collation_name(ctx *Column_collation_nameContext) {}

// ExitColumn_collation_name is called when production column_collation_name is exited.
func (s *BasePlSqlParserListener) ExitColumn_collation_name(ctx *Column_collation_nameContext) {}

// EnterIdentity_clause is called when production identity_clause is entered.
func (s *BasePlSqlParserListener) EnterIdentity_clause(ctx *Identity_clauseContext) {}

// ExitIdentity_clause is called when production identity_clause is exited.
func (s *BasePlSqlParserListener) ExitIdentity_clause(ctx *Identity_clauseContext) {}

// EnterIdentity_options_parentheses is called when production identity_options_parentheses is entered.
func (s *BasePlSqlParserListener) EnterIdentity_options_parentheses(ctx *Identity_options_parenthesesContext) {
}

// ExitIdentity_options_parentheses is called when production identity_options_parentheses is exited.
func (s *BasePlSqlParserListener) ExitIdentity_options_parentheses(ctx *Identity_options_parenthesesContext) {
}

// EnterIdentity_options is called when production identity_options is entered.
func (s *BasePlSqlParserListener) EnterIdentity_options(ctx *Identity_optionsContext) {}

// ExitIdentity_options is called when production identity_options is exited.
func (s *BasePlSqlParserListener) ExitIdentity_options(ctx *Identity_optionsContext) {}

// EnterVirtual_column_definition is called when production virtual_column_definition is entered.
func (s *BasePlSqlParserListener) EnterVirtual_column_definition(ctx *Virtual_column_definitionContext) {
}

// ExitVirtual_column_definition is called when production virtual_column_definition is exited.
func (s *BasePlSqlParserListener) ExitVirtual_column_definition(ctx *Virtual_column_definitionContext) {
}

// EnterVirtual_column_expression is called when production virtual_column_expression is entered.
func (s *BasePlSqlParserListener) EnterVirtual_column_expression(ctx *Virtual_column_expressionContext) {
}

// ExitVirtual_column_expression is called when production virtual_column_expression is exited.
func (s *BasePlSqlParserListener) ExitVirtual_column_expression(ctx *Virtual_column_expressionContext) {
}

// EnterAutogenerated_sequence_definition is called when production autogenerated_sequence_definition is entered.
func (s *BasePlSqlParserListener) EnterAutogenerated_sequence_definition(ctx *Autogenerated_sequence_definitionContext) {
}

// ExitAutogenerated_sequence_definition is called when production autogenerated_sequence_definition is exited.
func (s *BasePlSqlParserListener) ExitAutogenerated_sequence_definition(ctx *Autogenerated_sequence_definitionContext) {
}

// EnterBy_user_for_statistics_clause is called when production by_user_for_statistics_clause is entered.
func (s *BasePlSqlParserListener) EnterBy_user_for_statistics_clause(ctx *By_user_for_statistics_clauseContext) {
}

// ExitBy_user_for_statistics_clause is called when production by_user_for_statistics_clause is exited.
func (s *BasePlSqlParserListener) ExitBy_user_for_statistics_clause(ctx *By_user_for_statistics_clauseContext) {
}

// EnterEvaluation_edition_clause is called when production evaluation_edition_clause is entered.
func (s *BasePlSqlParserListener) EnterEvaluation_edition_clause(ctx *Evaluation_edition_clauseContext) {
}

// ExitEvaluation_edition_clause is called when production evaluation_edition_clause is exited.
func (s *BasePlSqlParserListener) ExitEvaluation_edition_clause(ctx *Evaluation_edition_clauseContext) {
}

// EnterNested_table_col_properties is called when production nested_table_col_properties is entered.
func (s *BasePlSqlParserListener) EnterNested_table_col_properties(ctx *Nested_table_col_propertiesContext) {
}

// ExitNested_table_col_properties is called when production nested_table_col_properties is exited.
func (s *BasePlSqlParserListener) ExitNested_table_col_properties(ctx *Nested_table_col_propertiesContext) {
}

// EnterNested_item is called when production nested_item is entered.
func (s *BasePlSqlParserListener) EnterNested_item(ctx *Nested_itemContext) {}

// ExitNested_item is called when production nested_item is exited.
func (s *BasePlSqlParserListener) ExitNested_item(ctx *Nested_itemContext) {}

// EnterSubstitutable_column_clause is called when production substitutable_column_clause is entered.
func (s *BasePlSqlParserListener) EnterSubstitutable_column_clause(ctx *Substitutable_column_clauseContext) {
}

// ExitSubstitutable_column_clause is called when production substitutable_column_clause is exited.
func (s *BasePlSqlParserListener) ExitSubstitutable_column_clause(ctx *Substitutable_column_clauseContext) {
}

// EnterPartition_name is called when production partition_name is entered.
func (s *BasePlSqlParserListener) EnterPartition_name(ctx *Partition_nameContext) {}

// ExitPartition_name is called when production partition_name is exited.
func (s *BasePlSqlParserListener) ExitPartition_name(ctx *Partition_nameContext) {}

// EnterSupplemental_logging_props is called when production supplemental_logging_props is entered.
func (s *BasePlSqlParserListener) EnterSupplemental_logging_props(ctx *Supplemental_logging_propsContext) {
}

// ExitSupplemental_logging_props is called when production supplemental_logging_props is exited.
func (s *BasePlSqlParserListener) ExitSupplemental_logging_props(ctx *Supplemental_logging_propsContext) {
}

// EnterObject_type_col_properties is called when production object_type_col_properties is entered.
func (s *BasePlSqlParserListener) EnterObject_type_col_properties(ctx *Object_type_col_propertiesContext) {
}

// ExitObject_type_col_properties is called when production object_type_col_properties is exited.
func (s *BasePlSqlParserListener) ExitObject_type_col_properties(ctx *Object_type_col_propertiesContext) {
}

// EnterConstraint_clauses is called when production constraint_clauses is entered.
func (s *BasePlSqlParserListener) EnterConstraint_clauses(ctx *Constraint_clausesContext) {}

// ExitConstraint_clauses is called when production constraint_clauses is exited.
func (s *BasePlSqlParserListener) ExitConstraint_clauses(ctx *Constraint_clausesContext) {}

// EnterOld_constraint_name is called when production old_constraint_name is entered.
func (s *BasePlSqlParserListener) EnterOld_constraint_name(ctx *Old_constraint_nameContext) {}

// ExitOld_constraint_name is called when production old_constraint_name is exited.
func (s *BasePlSqlParserListener) ExitOld_constraint_name(ctx *Old_constraint_nameContext) {}

// EnterNew_constraint_name is called when production new_constraint_name is entered.
func (s *BasePlSqlParserListener) EnterNew_constraint_name(ctx *New_constraint_nameContext) {}

// ExitNew_constraint_name is called when production new_constraint_name is exited.
func (s *BasePlSqlParserListener) ExitNew_constraint_name(ctx *New_constraint_nameContext) {}

// EnterDrop_constraint_clause is called when production drop_constraint_clause is entered.
func (s *BasePlSqlParserListener) EnterDrop_constraint_clause(ctx *Drop_constraint_clauseContext) {}

// ExitDrop_constraint_clause is called when production drop_constraint_clause is exited.
func (s *BasePlSqlParserListener) ExitDrop_constraint_clause(ctx *Drop_constraint_clauseContext) {}

// EnterCheck_constraint is called when production check_constraint is entered.
func (s *BasePlSqlParserListener) EnterCheck_constraint(ctx *Check_constraintContext) {}

// ExitCheck_constraint is called when production check_constraint is exited.
func (s *BasePlSqlParserListener) ExitCheck_constraint(ctx *Check_constraintContext) {}

// EnterForeign_key_clause is called when production foreign_key_clause is entered.
func (s *BasePlSqlParserListener) EnterForeign_key_clause(ctx *Foreign_key_clauseContext) {}

// ExitForeign_key_clause is called when production foreign_key_clause is exited.
func (s *BasePlSqlParserListener) ExitForeign_key_clause(ctx *Foreign_key_clauseContext) {}

// EnterReferences_clause is called when production references_clause is entered.
func (s *BasePlSqlParserListener) EnterReferences_clause(ctx *References_clauseContext) {}

// ExitReferences_clause is called when production references_clause is exited.
func (s *BasePlSqlParserListener) ExitReferences_clause(ctx *References_clauseContext) {}

// EnterOn_delete_clause is called when production on_delete_clause is entered.
func (s *BasePlSqlParserListener) EnterOn_delete_clause(ctx *On_delete_clauseContext) {}

// ExitOn_delete_clause is called when production on_delete_clause is exited.
func (s *BasePlSqlParserListener) ExitOn_delete_clause(ctx *On_delete_clauseContext) {}

// EnterAnonymous_block is called when production anonymous_block is entered.
func (s *BasePlSqlParserListener) EnterAnonymous_block(ctx *Anonymous_blockContext) {}

// ExitAnonymous_block is called when production anonymous_block is exited.
func (s *BasePlSqlParserListener) ExitAnonymous_block(ctx *Anonymous_blockContext) {}

// EnterInvoker_rights_clause is called when production invoker_rights_clause is entered.
func (s *BasePlSqlParserListener) EnterInvoker_rights_clause(ctx *Invoker_rights_clauseContext) {}

// ExitInvoker_rights_clause is called when production invoker_rights_clause is exited.
func (s *BasePlSqlParserListener) ExitInvoker_rights_clause(ctx *Invoker_rights_clauseContext) {}

// EnterCall_spec is called when production call_spec is entered.
func (s *BasePlSqlParserListener) EnterCall_spec(ctx *Call_specContext) {}

// ExitCall_spec is called when production call_spec is exited.
func (s *BasePlSqlParserListener) ExitCall_spec(ctx *Call_specContext) {}

// EnterJava_spec is called when production java_spec is entered.
func (s *BasePlSqlParserListener) EnterJava_spec(ctx *Java_specContext) {}

// ExitJava_spec is called when production java_spec is exited.
func (s *BasePlSqlParserListener) ExitJava_spec(ctx *Java_specContext) {}

// EnterC_spec is called when production c_spec is entered.
func (s *BasePlSqlParserListener) EnterC_spec(ctx *C_specContext) {}

// ExitC_spec is called when production c_spec is exited.
func (s *BasePlSqlParserListener) ExitC_spec(ctx *C_specContext) {}

// EnterC_agent_in_clause is called when production c_agent_in_clause is entered.
func (s *BasePlSqlParserListener) EnterC_agent_in_clause(ctx *C_agent_in_clauseContext) {}

// ExitC_agent_in_clause is called when production c_agent_in_clause is exited.
func (s *BasePlSqlParserListener) ExitC_agent_in_clause(ctx *C_agent_in_clauseContext) {}

// EnterC_parameters_clause is called when production c_parameters_clause is entered.
func (s *BasePlSqlParserListener) EnterC_parameters_clause(ctx *C_parameters_clauseContext) {}

// ExitC_parameters_clause is called when production c_parameters_clause is exited.
func (s *BasePlSqlParserListener) ExitC_parameters_clause(ctx *C_parameters_clauseContext) {}

// EnterC_external_parameter is called when production c_external_parameter is entered.
func (s *BasePlSqlParserListener) EnterC_external_parameter(ctx *C_external_parameterContext) {}

// ExitC_external_parameter is called when production c_external_parameter is exited.
func (s *BasePlSqlParserListener) ExitC_external_parameter(ctx *C_external_parameterContext) {}

// EnterC_property is called when production c_property is entered.
func (s *BasePlSqlParserListener) EnterC_property(ctx *C_propertyContext) {}

// ExitC_property is called when production c_property is exited.
func (s *BasePlSqlParserListener) ExitC_property(ctx *C_propertyContext) {}

// EnterParameter is called when production parameter is entered.
func (s *BasePlSqlParserListener) EnterParameter(ctx *ParameterContext) {}

// ExitParameter is called when production parameter is exited.
func (s *BasePlSqlParserListener) ExitParameter(ctx *ParameterContext) {}

// EnterDefault_value_part is called when production default_value_part is entered.
func (s *BasePlSqlParserListener) EnterDefault_value_part(ctx *Default_value_partContext) {}

// ExitDefault_value_part is called when production default_value_part is exited.
func (s *BasePlSqlParserListener) ExitDefault_value_part(ctx *Default_value_partContext) {}

// EnterSeq_of_declare_specs is called when production seq_of_declare_specs is entered.
func (s *BasePlSqlParserListener) EnterSeq_of_declare_specs(ctx *Seq_of_declare_specsContext) {}

// ExitSeq_of_declare_specs is called when production seq_of_declare_specs is exited.
func (s *BasePlSqlParserListener) ExitSeq_of_declare_specs(ctx *Seq_of_declare_specsContext) {}

// EnterDeclare_spec is called when production declare_spec is entered.
func (s *BasePlSqlParserListener) EnterDeclare_spec(ctx *Declare_specContext) {}

// ExitDeclare_spec is called when production declare_spec is exited.
func (s *BasePlSqlParserListener) ExitDeclare_spec(ctx *Declare_specContext) {}

// EnterVariable_declaration is called when production variable_declaration is entered.
func (s *BasePlSqlParserListener) EnterVariable_declaration(ctx *Variable_declarationContext) {}

// ExitVariable_declaration is called when production variable_declaration is exited.
func (s *BasePlSqlParserListener) ExitVariable_declaration(ctx *Variable_declarationContext) {}

// EnterSubtype_declaration is called when production subtype_declaration is entered.
func (s *BasePlSqlParserListener) EnterSubtype_declaration(ctx *Subtype_declarationContext) {}

// ExitSubtype_declaration is called when production subtype_declaration is exited.
func (s *BasePlSqlParserListener) ExitSubtype_declaration(ctx *Subtype_declarationContext) {}

// EnterCursor_declaration is called when production cursor_declaration is entered.
func (s *BasePlSqlParserListener) EnterCursor_declaration(ctx *Cursor_declarationContext) {}

// ExitCursor_declaration is called when production cursor_declaration is exited.
func (s *BasePlSqlParserListener) ExitCursor_declaration(ctx *Cursor_declarationContext) {}

// EnterParameter_spec is called when production parameter_spec is entered.
func (s *BasePlSqlParserListener) EnterParameter_spec(ctx *Parameter_specContext) {}

// ExitParameter_spec is called when production parameter_spec is exited.
func (s *BasePlSqlParserListener) ExitParameter_spec(ctx *Parameter_specContext) {}

// EnterException_declaration is called when production exception_declaration is entered.
func (s *BasePlSqlParserListener) EnterException_declaration(ctx *Exception_declarationContext) {}

// ExitException_declaration is called when production exception_declaration is exited.
func (s *BasePlSqlParserListener) ExitException_declaration(ctx *Exception_declarationContext) {}

// EnterPragma_declaration is called when production pragma_declaration is entered.
func (s *BasePlSqlParserListener) EnterPragma_declaration(ctx *Pragma_declarationContext) {}

// ExitPragma_declaration is called when production pragma_declaration is exited.
func (s *BasePlSqlParserListener) ExitPragma_declaration(ctx *Pragma_declarationContext) {}

// EnterRecord_type_def is called when production record_type_def is entered.
func (s *BasePlSqlParserListener) EnterRecord_type_def(ctx *Record_type_defContext) {}

// ExitRecord_type_def is called when production record_type_def is exited.
func (s *BasePlSqlParserListener) ExitRecord_type_def(ctx *Record_type_defContext) {}

// EnterField_spec is called when production field_spec is entered.
func (s *BasePlSqlParserListener) EnterField_spec(ctx *Field_specContext) {}

// ExitField_spec is called when production field_spec is exited.
func (s *BasePlSqlParserListener) ExitField_spec(ctx *Field_specContext) {}

// EnterRef_cursor_type_def is called when production ref_cursor_type_def is entered.
func (s *BasePlSqlParserListener) EnterRef_cursor_type_def(ctx *Ref_cursor_type_defContext) {}

// ExitRef_cursor_type_def is called when production ref_cursor_type_def is exited.
func (s *BasePlSqlParserListener) ExitRef_cursor_type_def(ctx *Ref_cursor_type_defContext) {}

// EnterType_declaration is called when production type_declaration is entered.
func (s *BasePlSqlParserListener) EnterType_declaration(ctx *Type_declarationContext) {}

// ExitType_declaration is called when production type_declaration is exited.
func (s *BasePlSqlParserListener) ExitType_declaration(ctx *Type_declarationContext) {}

// EnterTable_type_def is called when production table_type_def is entered.
func (s *BasePlSqlParserListener) EnterTable_type_def(ctx *Table_type_defContext) {}

// ExitTable_type_def is called when production table_type_def is exited.
func (s *BasePlSqlParserListener) ExitTable_type_def(ctx *Table_type_defContext) {}

// EnterTable_indexed_by_part is called when production table_indexed_by_part is entered.
func (s *BasePlSqlParserListener) EnterTable_indexed_by_part(ctx *Table_indexed_by_partContext) {}

// ExitTable_indexed_by_part is called when production table_indexed_by_part is exited.
func (s *BasePlSqlParserListener) ExitTable_indexed_by_part(ctx *Table_indexed_by_partContext) {}

// EnterVarray_type_def is called when production varray_type_def is entered.
func (s *BasePlSqlParserListener) EnterVarray_type_def(ctx *Varray_type_defContext) {}

// ExitVarray_type_def is called when production varray_type_def is exited.
func (s *BasePlSqlParserListener) ExitVarray_type_def(ctx *Varray_type_defContext) {}

// EnterSeq_of_statements is called when production seq_of_statements is entered.
func (s *BasePlSqlParserListener) EnterSeq_of_statements(ctx *Seq_of_statementsContext) {}

// ExitSeq_of_statements is called when production seq_of_statements is exited.
func (s *BasePlSqlParserListener) ExitSeq_of_statements(ctx *Seq_of_statementsContext) {}

// EnterLabel_declaration is called when production label_declaration is entered.
func (s *BasePlSqlParserListener) EnterLabel_declaration(ctx *Label_declarationContext) {}

// ExitLabel_declaration is called when production label_declaration is exited.
func (s *BasePlSqlParserListener) ExitLabel_declaration(ctx *Label_declarationContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasePlSqlParserListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasePlSqlParserListener) ExitStatement(ctx *StatementContext) {}

// EnterAssignment_statement is called when production assignment_statement is entered.
func (s *BasePlSqlParserListener) EnterAssignment_statement(ctx *Assignment_statementContext) {}

// ExitAssignment_statement is called when production assignment_statement is exited.
func (s *BasePlSqlParserListener) ExitAssignment_statement(ctx *Assignment_statementContext) {}

// EnterContinue_statement is called when production continue_statement is entered.
func (s *BasePlSqlParserListener) EnterContinue_statement(ctx *Continue_statementContext) {}

// ExitContinue_statement is called when production continue_statement is exited.
func (s *BasePlSqlParserListener) ExitContinue_statement(ctx *Continue_statementContext) {}

// EnterExit_statement is called when production exit_statement is entered.
func (s *BasePlSqlParserListener) EnterExit_statement(ctx *Exit_statementContext) {}

// ExitExit_statement is called when production exit_statement is exited.
func (s *BasePlSqlParserListener) ExitExit_statement(ctx *Exit_statementContext) {}

// EnterGoto_statement is called when production goto_statement is entered.
func (s *BasePlSqlParserListener) EnterGoto_statement(ctx *Goto_statementContext) {}

// ExitGoto_statement is called when production goto_statement is exited.
func (s *BasePlSqlParserListener) ExitGoto_statement(ctx *Goto_statementContext) {}

// EnterIf_statement is called when production if_statement is entered.
func (s *BasePlSqlParserListener) EnterIf_statement(ctx *If_statementContext) {}

// ExitIf_statement is called when production if_statement is exited.
func (s *BasePlSqlParserListener) ExitIf_statement(ctx *If_statementContext) {}

// EnterElsif_part is called when production elsif_part is entered.
func (s *BasePlSqlParserListener) EnterElsif_part(ctx *Elsif_partContext) {}

// ExitElsif_part is called when production elsif_part is exited.
func (s *BasePlSqlParserListener) ExitElsif_part(ctx *Elsif_partContext) {}

// EnterElse_part is called when production else_part is entered.
func (s *BasePlSqlParserListener) EnterElse_part(ctx *Else_partContext) {}

// ExitElse_part is called when production else_part is exited.
func (s *BasePlSqlParserListener) ExitElse_part(ctx *Else_partContext) {}

// EnterLoop_statement is called when production loop_statement is entered.
func (s *BasePlSqlParserListener) EnterLoop_statement(ctx *Loop_statementContext) {}

// ExitLoop_statement is called when production loop_statement is exited.
func (s *BasePlSqlParserListener) ExitLoop_statement(ctx *Loop_statementContext) {}

// EnterCursor_loop_param is called when production cursor_loop_param is entered.
func (s *BasePlSqlParserListener) EnterCursor_loop_param(ctx *Cursor_loop_paramContext) {}

// ExitCursor_loop_param is called when production cursor_loop_param is exited.
func (s *BasePlSqlParserListener) ExitCursor_loop_param(ctx *Cursor_loop_paramContext) {}

// EnterForall_statement is called when production forall_statement is entered.
func (s *BasePlSqlParserListener) EnterForall_statement(ctx *Forall_statementContext) {}

// ExitForall_statement is called when production forall_statement is exited.
func (s *BasePlSqlParserListener) ExitForall_statement(ctx *Forall_statementContext) {}

// EnterBounds_clause is called when production bounds_clause is entered.
func (s *BasePlSqlParserListener) EnterBounds_clause(ctx *Bounds_clauseContext) {}

// ExitBounds_clause is called when production bounds_clause is exited.
func (s *BasePlSqlParserListener) ExitBounds_clause(ctx *Bounds_clauseContext) {}

// EnterBetween_bound is called when production between_bound is entered.
func (s *BasePlSqlParserListener) EnterBetween_bound(ctx *Between_boundContext) {}

// ExitBetween_bound is called when production between_bound is exited.
func (s *BasePlSqlParserListener) ExitBetween_bound(ctx *Between_boundContext) {}

// EnterLower_bound is called when production lower_bound is entered.
func (s *BasePlSqlParserListener) EnterLower_bound(ctx *Lower_boundContext) {}

// ExitLower_bound is called when production lower_bound is exited.
func (s *BasePlSqlParserListener) ExitLower_bound(ctx *Lower_boundContext) {}

// EnterUpper_bound is called when production upper_bound is entered.
func (s *BasePlSqlParserListener) EnterUpper_bound(ctx *Upper_boundContext) {}

// ExitUpper_bound is called when production upper_bound is exited.
func (s *BasePlSqlParserListener) ExitUpper_bound(ctx *Upper_boundContext) {}

// EnterNull_statement is called when production null_statement is entered.
func (s *BasePlSqlParserListener) EnterNull_statement(ctx *Null_statementContext) {}

// ExitNull_statement is called when production null_statement is exited.
func (s *BasePlSqlParserListener) ExitNull_statement(ctx *Null_statementContext) {}

// EnterRaise_statement is called when production raise_statement is entered.
func (s *BasePlSqlParserListener) EnterRaise_statement(ctx *Raise_statementContext) {}

// ExitRaise_statement is called when production raise_statement is exited.
func (s *BasePlSqlParserListener) ExitRaise_statement(ctx *Raise_statementContext) {}

// EnterReturn_statement is called when production return_statement is entered.
func (s *BasePlSqlParserListener) EnterReturn_statement(ctx *Return_statementContext) {}

// ExitReturn_statement is called when production return_statement is exited.
func (s *BasePlSqlParserListener) ExitReturn_statement(ctx *Return_statementContext) {}

// EnterCall_statement is called when production call_statement is entered.
func (s *BasePlSqlParserListener) EnterCall_statement(ctx *Call_statementContext) {}

// ExitCall_statement is called when production call_statement is exited.
func (s *BasePlSqlParserListener) ExitCall_statement(ctx *Call_statementContext) {}

// EnterPipe_row_statement is called when production pipe_row_statement is entered.
func (s *BasePlSqlParserListener) EnterPipe_row_statement(ctx *Pipe_row_statementContext) {}

// ExitPipe_row_statement is called when production pipe_row_statement is exited.
func (s *BasePlSqlParserListener) ExitPipe_row_statement(ctx *Pipe_row_statementContext) {}

// EnterSelection_directive is called when production selection_directive is entered.
func (s *BasePlSqlParserListener) EnterSelection_directive(ctx *Selection_directiveContext) {}

// ExitSelection_directive is called when production selection_directive is exited.
func (s *BasePlSqlParserListener) ExitSelection_directive(ctx *Selection_directiveContext) {}

// EnterError_directive is called when production error_directive is entered.
func (s *BasePlSqlParserListener) EnterError_directive(ctx *Error_directiveContext) {}

// ExitError_directive is called when production error_directive is exited.
func (s *BasePlSqlParserListener) ExitError_directive(ctx *Error_directiveContext) {}

// EnterSelection_directive_body is called when production selection_directive_body is entered.
func (s *BasePlSqlParserListener) EnterSelection_directive_body(ctx *Selection_directive_bodyContext) {
}

// ExitSelection_directive_body is called when production selection_directive_body is exited.
func (s *BasePlSqlParserListener) ExitSelection_directive_body(ctx *Selection_directive_bodyContext) {
}

// EnterBody is called when production body is entered.
func (s *BasePlSqlParserListener) EnterBody(ctx *BodyContext) {}

// ExitBody is called when production body is exited.
func (s *BasePlSqlParserListener) ExitBody(ctx *BodyContext) {}

// EnterException_handler is called when production exception_handler is entered.
func (s *BasePlSqlParserListener) EnterException_handler(ctx *Exception_handlerContext) {}

// ExitException_handler is called when production exception_handler is exited.
func (s *BasePlSqlParserListener) ExitException_handler(ctx *Exception_handlerContext) {}

// EnterTrigger_block is called when production trigger_block is entered.
func (s *BasePlSqlParserListener) EnterTrigger_block(ctx *Trigger_blockContext) {}

// ExitTrigger_block is called when production trigger_block is exited.
func (s *BasePlSqlParserListener) ExitTrigger_block(ctx *Trigger_blockContext) {}

// EnterTps_block is called when production tps_block is entered.
func (s *BasePlSqlParserListener) EnterTps_block(ctx *Tps_blockContext) {}

// ExitTps_block is called when production tps_block is exited.
func (s *BasePlSqlParserListener) ExitTps_block(ctx *Tps_blockContext) {}

// EnterBlock is called when production block is entered.
func (s *BasePlSqlParserListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BasePlSqlParserListener) ExitBlock(ctx *BlockContext) {}

// EnterSql_statement is called when production sql_statement is entered.
func (s *BasePlSqlParserListener) EnterSql_statement(ctx *Sql_statementContext) {}

// ExitSql_statement is called when production sql_statement is exited.
func (s *BasePlSqlParserListener) ExitSql_statement(ctx *Sql_statementContext) {}

// EnterExecute_immediate is called when production execute_immediate is entered.
func (s *BasePlSqlParserListener) EnterExecute_immediate(ctx *Execute_immediateContext) {}

// ExitExecute_immediate is called when production execute_immediate is exited.
func (s *BasePlSqlParserListener) ExitExecute_immediate(ctx *Execute_immediateContext) {}

// EnterDynamic_returning_clause is called when production dynamic_returning_clause is entered.
func (s *BasePlSqlParserListener) EnterDynamic_returning_clause(ctx *Dynamic_returning_clauseContext) {
}

// ExitDynamic_returning_clause is called when production dynamic_returning_clause is exited.
func (s *BasePlSqlParserListener) ExitDynamic_returning_clause(ctx *Dynamic_returning_clauseContext) {
}

// EnterData_manipulation_language_statements is called when production data_manipulation_language_statements is entered.
func (s *BasePlSqlParserListener) EnterData_manipulation_language_statements(ctx *Data_manipulation_language_statementsContext) {
}

// ExitData_manipulation_language_statements is called when production data_manipulation_language_statements is exited.
func (s *BasePlSqlParserListener) ExitData_manipulation_language_statements(ctx *Data_manipulation_language_statementsContext) {
}

// EnterCursor_manipulation_statements is called when production cursor_manipulation_statements is entered.
func (s *BasePlSqlParserListener) EnterCursor_manipulation_statements(ctx *Cursor_manipulation_statementsContext) {
}

// ExitCursor_manipulation_statements is called when production cursor_manipulation_statements is exited.
func (s *BasePlSqlParserListener) ExitCursor_manipulation_statements(ctx *Cursor_manipulation_statementsContext) {
}

// EnterClose_statement is called when production close_statement is entered.
func (s *BasePlSqlParserListener) EnterClose_statement(ctx *Close_statementContext) {}

// ExitClose_statement is called when production close_statement is exited.
func (s *BasePlSqlParserListener) ExitClose_statement(ctx *Close_statementContext) {}

// EnterOpen_statement is called when production open_statement is entered.
func (s *BasePlSqlParserListener) EnterOpen_statement(ctx *Open_statementContext) {}

// ExitOpen_statement is called when production open_statement is exited.
func (s *BasePlSqlParserListener) ExitOpen_statement(ctx *Open_statementContext) {}

// EnterFetch_statement is called when production fetch_statement is entered.
func (s *BasePlSqlParserListener) EnterFetch_statement(ctx *Fetch_statementContext) {}

// ExitFetch_statement is called when production fetch_statement is exited.
func (s *BasePlSqlParserListener) ExitFetch_statement(ctx *Fetch_statementContext) {}

// EnterVariable_or_collection is called when production variable_or_collection is entered.
func (s *BasePlSqlParserListener) EnterVariable_or_collection(ctx *Variable_or_collectionContext) {}

// ExitVariable_or_collection is called when production variable_or_collection is exited.
func (s *BasePlSqlParserListener) ExitVariable_or_collection(ctx *Variable_or_collectionContext) {}

// EnterOpen_for_statement is called when production open_for_statement is entered.
func (s *BasePlSqlParserListener) EnterOpen_for_statement(ctx *Open_for_statementContext) {}

// ExitOpen_for_statement is called when production open_for_statement is exited.
func (s *BasePlSqlParserListener) ExitOpen_for_statement(ctx *Open_for_statementContext) {}

// EnterTransaction_control_statements is called when production transaction_control_statements is entered.
func (s *BasePlSqlParserListener) EnterTransaction_control_statements(ctx *Transaction_control_statementsContext) {
}

// ExitTransaction_control_statements is called when production transaction_control_statements is exited.
func (s *BasePlSqlParserListener) ExitTransaction_control_statements(ctx *Transaction_control_statementsContext) {
}

// EnterSet_transaction_command is called when production set_transaction_command is entered.
func (s *BasePlSqlParserListener) EnterSet_transaction_command(ctx *Set_transaction_commandContext) {}

// ExitSet_transaction_command is called when production set_transaction_command is exited.
func (s *BasePlSqlParserListener) ExitSet_transaction_command(ctx *Set_transaction_commandContext) {}

// EnterSet_constraint_command is called when production set_constraint_command is entered.
func (s *BasePlSqlParserListener) EnterSet_constraint_command(ctx *Set_constraint_commandContext) {}

// ExitSet_constraint_command is called when production set_constraint_command is exited.
func (s *BasePlSqlParserListener) ExitSet_constraint_command(ctx *Set_constraint_commandContext) {}

// EnterCommit_statement is called when production commit_statement is entered.
func (s *BasePlSqlParserListener) EnterCommit_statement(ctx *Commit_statementContext) {}

// ExitCommit_statement is called when production commit_statement is exited.
func (s *BasePlSqlParserListener) ExitCommit_statement(ctx *Commit_statementContext) {}

// EnterWrite_clause is called when production write_clause is entered.
func (s *BasePlSqlParserListener) EnterWrite_clause(ctx *Write_clauseContext) {}

// ExitWrite_clause is called when production write_clause is exited.
func (s *BasePlSqlParserListener) ExitWrite_clause(ctx *Write_clauseContext) {}

// EnterRollback_statement is called when production rollback_statement is entered.
func (s *BasePlSqlParserListener) EnterRollback_statement(ctx *Rollback_statementContext) {}

// ExitRollback_statement is called when production rollback_statement is exited.
func (s *BasePlSqlParserListener) ExitRollback_statement(ctx *Rollback_statementContext) {}

// EnterSavepoint_statement is called when production savepoint_statement is entered.
func (s *BasePlSqlParserListener) EnterSavepoint_statement(ctx *Savepoint_statementContext) {}

// ExitSavepoint_statement is called when production savepoint_statement is exited.
func (s *BasePlSqlParserListener) ExitSavepoint_statement(ctx *Savepoint_statementContext) {}

// EnterCollection_method_call is called when production collection_method_call is entered.
func (s *BasePlSqlParserListener) EnterCollection_method_call(ctx *Collection_method_callContext) {}

// ExitCollection_method_call is called when production collection_method_call is exited.
func (s *BasePlSqlParserListener) ExitCollection_method_call(ctx *Collection_method_callContext) {}

// EnterExplain_statement is called when production explain_statement is entered.
func (s *BasePlSqlParserListener) EnterExplain_statement(ctx *Explain_statementContext) {}

// ExitExplain_statement is called when production explain_statement is exited.
func (s *BasePlSqlParserListener) ExitExplain_statement(ctx *Explain_statementContext) {}

// EnterSelect_only_statement is called when production select_only_statement is entered.
func (s *BasePlSqlParserListener) EnterSelect_only_statement(ctx *Select_only_statementContext) {}

// ExitSelect_only_statement is called when production select_only_statement is exited.
func (s *BasePlSqlParserListener) ExitSelect_only_statement(ctx *Select_only_statementContext) {}

// EnterSelect_statement is called when production select_statement is entered.
func (s *BasePlSqlParserListener) EnterSelect_statement(ctx *Select_statementContext) {}

// ExitSelect_statement is called when production select_statement is exited.
func (s *BasePlSqlParserListener) ExitSelect_statement(ctx *Select_statementContext) {}

// EnterWith_clause is called when production with_clause is entered.
func (s *BasePlSqlParserListener) EnterWith_clause(ctx *With_clauseContext) {}

// ExitWith_clause is called when production with_clause is exited.
func (s *BasePlSqlParserListener) ExitWith_clause(ctx *With_clauseContext) {}

// EnterWith_factoring_clause is called when production with_factoring_clause is entered.
func (s *BasePlSqlParserListener) EnterWith_factoring_clause(ctx *With_factoring_clauseContext) {}

// ExitWith_factoring_clause is called when production with_factoring_clause is exited.
func (s *BasePlSqlParserListener) ExitWith_factoring_clause(ctx *With_factoring_clauseContext) {}

// EnterSubquery_factoring_clause is called when production subquery_factoring_clause is entered.
func (s *BasePlSqlParserListener) EnterSubquery_factoring_clause(ctx *Subquery_factoring_clauseContext) {
}

// ExitSubquery_factoring_clause is called when production subquery_factoring_clause is exited.
func (s *BasePlSqlParserListener) ExitSubquery_factoring_clause(ctx *Subquery_factoring_clauseContext) {
}

// EnterSearch_clause is called when production search_clause is entered.
func (s *BasePlSqlParserListener) EnterSearch_clause(ctx *Search_clauseContext) {}

// ExitSearch_clause is called when production search_clause is exited.
func (s *BasePlSqlParserListener) ExitSearch_clause(ctx *Search_clauseContext) {}

// EnterCycle_clause is called when production cycle_clause is entered.
func (s *BasePlSqlParserListener) EnterCycle_clause(ctx *Cycle_clauseContext) {}

// ExitCycle_clause is called when production cycle_clause is exited.
func (s *BasePlSqlParserListener) ExitCycle_clause(ctx *Cycle_clauseContext) {}

// EnterSubav_factoring_clause is called when production subav_factoring_clause is entered.
func (s *BasePlSqlParserListener) EnterSubav_factoring_clause(ctx *Subav_factoring_clauseContext) {}

// ExitSubav_factoring_clause is called when production subav_factoring_clause is exited.
func (s *BasePlSqlParserListener) ExitSubav_factoring_clause(ctx *Subav_factoring_clauseContext) {}

// EnterSubav_clause is called when production subav_clause is entered.
func (s *BasePlSqlParserListener) EnterSubav_clause(ctx *Subav_clauseContext) {}

// ExitSubav_clause is called when production subav_clause is exited.
func (s *BasePlSqlParserListener) ExitSubav_clause(ctx *Subav_clauseContext) {}

// EnterHierarchies_clause is called when production hierarchies_clause is entered.
func (s *BasePlSqlParserListener) EnterHierarchies_clause(ctx *Hierarchies_clauseContext) {}

// ExitHierarchies_clause is called when production hierarchies_clause is exited.
func (s *BasePlSqlParserListener) ExitHierarchies_clause(ctx *Hierarchies_clauseContext) {}

// EnterFilter_clauses is called when production filter_clauses is entered.
func (s *BasePlSqlParserListener) EnterFilter_clauses(ctx *Filter_clausesContext) {}

// ExitFilter_clauses is called when production filter_clauses is exited.
func (s *BasePlSqlParserListener) ExitFilter_clauses(ctx *Filter_clausesContext) {}

// EnterFilter_clause is called when production filter_clause is entered.
func (s *BasePlSqlParserListener) EnterFilter_clause(ctx *Filter_clauseContext) {}

// ExitFilter_clause is called when production filter_clause is exited.
func (s *BasePlSqlParserListener) ExitFilter_clause(ctx *Filter_clauseContext) {}

// EnterAdd_calcs_clause is called when production add_calcs_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_calcs_clause(ctx *Add_calcs_clauseContext) {}

// ExitAdd_calcs_clause is called when production add_calcs_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_calcs_clause(ctx *Add_calcs_clauseContext) {}

// EnterAdd_calc_meas_clause is called when production add_calc_meas_clause is entered.
func (s *BasePlSqlParserListener) EnterAdd_calc_meas_clause(ctx *Add_calc_meas_clauseContext) {}

// ExitAdd_calc_meas_clause is called when production add_calc_meas_clause is exited.
func (s *BasePlSqlParserListener) ExitAdd_calc_meas_clause(ctx *Add_calc_meas_clauseContext) {}

// EnterSubquery is called when production subquery is entered.
func (s *BasePlSqlParserListener) EnterSubquery(ctx *SubqueryContext) {}

// ExitSubquery is called when production subquery is exited.
func (s *BasePlSqlParserListener) ExitSubquery(ctx *SubqueryContext) {}

// EnterSubquery_basic_elements is called when production subquery_basic_elements is entered.
func (s *BasePlSqlParserListener) EnterSubquery_basic_elements(ctx *Subquery_basic_elementsContext) {}

// ExitSubquery_basic_elements is called when production subquery_basic_elements is exited.
func (s *BasePlSqlParserListener) ExitSubquery_basic_elements(ctx *Subquery_basic_elementsContext) {}

// EnterSubquery_operation_part is called when production subquery_operation_part is entered.
func (s *BasePlSqlParserListener) EnterSubquery_operation_part(ctx *Subquery_operation_partContext) {}

// ExitSubquery_operation_part is called when production subquery_operation_part is exited.
func (s *BasePlSqlParserListener) ExitSubquery_operation_part(ctx *Subquery_operation_partContext) {}

// EnterQuery_block is called when production query_block is entered.
func (s *BasePlSqlParserListener) EnterQuery_block(ctx *Query_blockContext) {}

// ExitQuery_block is called when production query_block is exited.
func (s *BasePlSqlParserListener) ExitQuery_block(ctx *Query_blockContext) {}

// EnterSelected_list is called when production selected_list is entered.
func (s *BasePlSqlParserListener) EnterSelected_list(ctx *Selected_listContext) {}

// ExitSelected_list is called when production selected_list is exited.
func (s *BasePlSqlParserListener) ExitSelected_list(ctx *Selected_listContext) {}

// EnterFrom_clause is called when production from_clause is entered.
func (s *BasePlSqlParserListener) EnterFrom_clause(ctx *From_clauseContext) {}

// ExitFrom_clause is called when production from_clause is exited.
func (s *BasePlSqlParserListener) ExitFrom_clause(ctx *From_clauseContext) {}

// EnterSelect_list_elements is called when production select_list_elements is entered.
func (s *BasePlSqlParserListener) EnterSelect_list_elements(ctx *Select_list_elementsContext) {}

// ExitSelect_list_elements is called when production select_list_elements is exited.
func (s *BasePlSqlParserListener) ExitSelect_list_elements(ctx *Select_list_elementsContext) {}

// EnterTable_ref_list is called when production table_ref_list is entered.
func (s *BasePlSqlParserListener) EnterTable_ref_list(ctx *Table_ref_listContext) {}

// ExitTable_ref_list is called when production table_ref_list is exited.
func (s *BasePlSqlParserListener) ExitTable_ref_list(ctx *Table_ref_listContext) {}

// EnterTable_ref is called when production table_ref is entered.
func (s *BasePlSqlParserListener) EnterTable_ref(ctx *Table_refContext) {}

// ExitTable_ref is called when production table_ref is exited.
func (s *BasePlSqlParserListener) ExitTable_ref(ctx *Table_refContext) {}

// EnterTable_ref_aux is called when production table_ref_aux is entered.
func (s *BasePlSqlParserListener) EnterTable_ref_aux(ctx *Table_ref_auxContext) {}

// ExitTable_ref_aux is called when production table_ref_aux is exited.
func (s *BasePlSqlParserListener) ExitTable_ref_aux(ctx *Table_ref_auxContext) {}

// EnterTable_ref_aux_internal_one is called when production table_ref_aux_internal_one is entered.
func (s *BasePlSqlParserListener) EnterTable_ref_aux_internal_one(ctx *Table_ref_aux_internal_oneContext) {
}

// ExitTable_ref_aux_internal_one is called when production table_ref_aux_internal_one is exited.
func (s *BasePlSqlParserListener) ExitTable_ref_aux_internal_one(ctx *Table_ref_aux_internal_oneContext) {
}

// EnterTable_ref_aux_internal_two is called when production table_ref_aux_internal_two is entered.
func (s *BasePlSqlParserListener) EnterTable_ref_aux_internal_two(ctx *Table_ref_aux_internal_twoContext) {
}

// ExitTable_ref_aux_internal_two is called when production table_ref_aux_internal_two is exited.
func (s *BasePlSqlParserListener) ExitTable_ref_aux_internal_two(ctx *Table_ref_aux_internal_twoContext) {
}

// EnterTable_ref_aux_internal_thre is called when production table_ref_aux_internal_thre is entered.
func (s *BasePlSqlParserListener) EnterTable_ref_aux_internal_thre(ctx *Table_ref_aux_internal_threContext) {
}

// ExitTable_ref_aux_internal_thre is called when production table_ref_aux_internal_thre is exited.
func (s *BasePlSqlParserListener) ExitTable_ref_aux_internal_thre(ctx *Table_ref_aux_internal_threContext) {
}

// EnterJoin_clause is called when production join_clause is entered.
func (s *BasePlSqlParserListener) EnterJoin_clause(ctx *Join_clauseContext) {}

// ExitJoin_clause is called when production join_clause is exited.
func (s *BasePlSqlParserListener) ExitJoin_clause(ctx *Join_clauseContext) {}

// EnterJoin_on_part is called when production join_on_part is entered.
func (s *BasePlSqlParserListener) EnterJoin_on_part(ctx *Join_on_partContext) {}

// ExitJoin_on_part is called when production join_on_part is exited.
func (s *BasePlSqlParserListener) ExitJoin_on_part(ctx *Join_on_partContext) {}

// EnterJoin_using_part is called when production join_using_part is entered.
func (s *BasePlSqlParserListener) EnterJoin_using_part(ctx *Join_using_partContext) {}

// ExitJoin_using_part is called when production join_using_part is exited.
func (s *BasePlSqlParserListener) ExitJoin_using_part(ctx *Join_using_partContext) {}

// EnterOuter_join_type is called when production outer_join_type is entered.
func (s *BasePlSqlParserListener) EnterOuter_join_type(ctx *Outer_join_typeContext) {}

// ExitOuter_join_type is called when production outer_join_type is exited.
func (s *BasePlSqlParserListener) ExitOuter_join_type(ctx *Outer_join_typeContext) {}

// EnterQuery_partition_clause is called when production query_partition_clause is entered.
func (s *BasePlSqlParserListener) EnterQuery_partition_clause(ctx *Query_partition_clauseContext) {}

// ExitQuery_partition_clause is called when production query_partition_clause is exited.
func (s *BasePlSqlParserListener) ExitQuery_partition_clause(ctx *Query_partition_clauseContext) {}

// EnterFlashback_query_clause is called when production flashback_query_clause is entered.
func (s *BasePlSqlParserListener) EnterFlashback_query_clause(ctx *Flashback_query_clauseContext) {}

// ExitFlashback_query_clause is called when production flashback_query_clause is exited.
func (s *BasePlSqlParserListener) ExitFlashback_query_clause(ctx *Flashback_query_clauseContext) {}

// EnterPivot_clause is called when production pivot_clause is entered.
func (s *BasePlSqlParserListener) EnterPivot_clause(ctx *Pivot_clauseContext) {}

// ExitPivot_clause is called when production pivot_clause is exited.
func (s *BasePlSqlParserListener) ExitPivot_clause(ctx *Pivot_clauseContext) {}

// EnterPivot_element is called when production pivot_element is entered.
func (s *BasePlSqlParserListener) EnterPivot_element(ctx *Pivot_elementContext) {}

// ExitPivot_element is called when production pivot_element is exited.
func (s *BasePlSqlParserListener) ExitPivot_element(ctx *Pivot_elementContext) {}

// EnterPivot_for_clause is called when production pivot_for_clause is entered.
func (s *BasePlSqlParserListener) EnterPivot_for_clause(ctx *Pivot_for_clauseContext) {}

// ExitPivot_for_clause is called when production pivot_for_clause is exited.
func (s *BasePlSqlParserListener) ExitPivot_for_clause(ctx *Pivot_for_clauseContext) {}

// EnterPivot_in_clause is called when production pivot_in_clause is entered.
func (s *BasePlSqlParserListener) EnterPivot_in_clause(ctx *Pivot_in_clauseContext) {}

// ExitPivot_in_clause is called when production pivot_in_clause is exited.
func (s *BasePlSqlParserListener) ExitPivot_in_clause(ctx *Pivot_in_clauseContext) {}

// EnterPivot_in_clause_element is called when production pivot_in_clause_element is entered.
func (s *BasePlSqlParserListener) EnterPivot_in_clause_element(ctx *Pivot_in_clause_elementContext) {}

// ExitPivot_in_clause_element is called when production pivot_in_clause_element is exited.
func (s *BasePlSqlParserListener) ExitPivot_in_clause_element(ctx *Pivot_in_clause_elementContext) {}

// EnterPivot_in_clause_elements is called when production pivot_in_clause_elements is entered.
func (s *BasePlSqlParserListener) EnterPivot_in_clause_elements(ctx *Pivot_in_clause_elementsContext) {
}

// ExitPivot_in_clause_elements is called when production pivot_in_clause_elements is exited.
func (s *BasePlSqlParserListener) ExitPivot_in_clause_elements(ctx *Pivot_in_clause_elementsContext) {
}

// EnterUnpivot_clause is called when production unpivot_clause is entered.
func (s *BasePlSqlParserListener) EnterUnpivot_clause(ctx *Unpivot_clauseContext) {}

// ExitUnpivot_clause is called when production unpivot_clause is exited.
func (s *BasePlSqlParserListener) ExitUnpivot_clause(ctx *Unpivot_clauseContext) {}

// EnterUnpivot_in_clause is called when production unpivot_in_clause is entered.
func (s *BasePlSqlParserListener) EnterUnpivot_in_clause(ctx *Unpivot_in_clauseContext) {}

// ExitUnpivot_in_clause is called when production unpivot_in_clause is exited.
func (s *BasePlSqlParserListener) ExitUnpivot_in_clause(ctx *Unpivot_in_clauseContext) {}

// EnterUnpivot_in_elements is called when production unpivot_in_elements is entered.
func (s *BasePlSqlParserListener) EnterUnpivot_in_elements(ctx *Unpivot_in_elementsContext) {}

// ExitUnpivot_in_elements is called when production unpivot_in_elements is exited.
func (s *BasePlSqlParserListener) ExitUnpivot_in_elements(ctx *Unpivot_in_elementsContext) {}

// EnterHierarchical_query_clause is called when production hierarchical_query_clause is entered.
func (s *BasePlSqlParserListener) EnterHierarchical_query_clause(ctx *Hierarchical_query_clauseContext) {
}

// ExitHierarchical_query_clause is called when production hierarchical_query_clause is exited.
func (s *BasePlSqlParserListener) ExitHierarchical_query_clause(ctx *Hierarchical_query_clauseContext) {
}

// EnterStart_part is called when production start_part is entered.
func (s *BasePlSqlParserListener) EnterStart_part(ctx *Start_partContext) {}

// ExitStart_part is called when production start_part is exited.
func (s *BasePlSqlParserListener) ExitStart_part(ctx *Start_partContext) {}

// EnterGroup_by_clause is called when production group_by_clause is entered.
func (s *BasePlSqlParserListener) EnterGroup_by_clause(ctx *Group_by_clauseContext) {}

// ExitGroup_by_clause is called when production group_by_clause is exited.
func (s *BasePlSqlParserListener) ExitGroup_by_clause(ctx *Group_by_clauseContext) {}

// EnterGroup_by_elements is called when production group_by_elements is entered.
func (s *BasePlSqlParserListener) EnterGroup_by_elements(ctx *Group_by_elementsContext) {}

// ExitGroup_by_elements is called when production group_by_elements is exited.
func (s *BasePlSqlParserListener) ExitGroup_by_elements(ctx *Group_by_elementsContext) {}

// EnterRollup_cube_clause is called when production rollup_cube_clause is entered.
func (s *BasePlSqlParserListener) EnterRollup_cube_clause(ctx *Rollup_cube_clauseContext) {}

// ExitRollup_cube_clause is called when production rollup_cube_clause is exited.
func (s *BasePlSqlParserListener) ExitRollup_cube_clause(ctx *Rollup_cube_clauseContext) {}

// EnterGrouping_sets_clause is called when production grouping_sets_clause is entered.
func (s *BasePlSqlParserListener) EnterGrouping_sets_clause(ctx *Grouping_sets_clauseContext) {}

// ExitGrouping_sets_clause is called when production grouping_sets_clause is exited.
func (s *BasePlSqlParserListener) ExitGrouping_sets_clause(ctx *Grouping_sets_clauseContext) {}

// EnterGrouping_sets_elements is called when production grouping_sets_elements is entered.
func (s *BasePlSqlParserListener) EnterGrouping_sets_elements(ctx *Grouping_sets_elementsContext) {}

// ExitGrouping_sets_elements is called when production grouping_sets_elements is exited.
func (s *BasePlSqlParserListener) ExitGrouping_sets_elements(ctx *Grouping_sets_elementsContext) {}

// EnterHaving_clause is called when production having_clause is entered.
func (s *BasePlSqlParserListener) EnterHaving_clause(ctx *Having_clauseContext) {}

// ExitHaving_clause is called when production having_clause is exited.
func (s *BasePlSqlParserListener) ExitHaving_clause(ctx *Having_clauseContext) {}

// EnterModel_clause is called when production model_clause is entered.
func (s *BasePlSqlParserListener) EnterModel_clause(ctx *Model_clauseContext) {}

// ExitModel_clause is called when production model_clause is exited.
func (s *BasePlSqlParserListener) ExitModel_clause(ctx *Model_clauseContext) {}

// EnterCell_reference_options is called when production cell_reference_options is entered.
func (s *BasePlSqlParserListener) EnterCell_reference_options(ctx *Cell_reference_optionsContext) {}

// ExitCell_reference_options is called when production cell_reference_options is exited.
func (s *BasePlSqlParserListener) ExitCell_reference_options(ctx *Cell_reference_optionsContext) {}

// EnterReturn_rows_clause is called when production return_rows_clause is entered.
func (s *BasePlSqlParserListener) EnterReturn_rows_clause(ctx *Return_rows_clauseContext) {}

// ExitReturn_rows_clause is called when production return_rows_clause is exited.
func (s *BasePlSqlParserListener) ExitReturn_rows_clause(ctx *Return_rows_clauseContext) {}

// EnterReference_model is called when production reference_model is entered.
func (s *BasePlSqlParserListener) EnterReference_model(ctx *Reference_modelContext) {}

// ExitReference_model is called when production reference_model is exited.
func (s *BasePlSqlParserListener) ExitReference_model(ctx *Reference_modelContext) {}

// EnterMain_model is called when production main_model is entered.
func (s *BasePlSqlParserListener) EnterMain_model(ctx *Main_modelContext) {}

// ExitMain_model is called when production main_model is exited.
func (s *BasePlSqlParserListener) ExitMain_model(ctx *Main_modelContext) {}

// EnterModel_column_clauses is called when production model_column_clauses is entered.
func (s *BasePlSqlParserListener) EnterModel_column_clauses(ctx *Model_column_clausesContext) {}

// ExitModel_column_clauses is called when production model_column_clauses is exited.
func (s *BasePlSqlParserListener) ExitModel_column_clauses(ctx *Model_column_clausesContext) {}

// EnterModel_column_partition_part is called when production model_column_partition_part is entered.
func (s *BasePlSqlParserListener) EnterModel_column_partition_part(ctx *Model_column_partition_partContext) {
}

// ExitModel_column_partition_part is called when production model_column_partition_part is exited.
func (s *BasePlSqlParserListener) ExitModel_column_partition_part(ctx *Model_column_partition_partContext) {
}

// EnterModel_column_list is called when production model_column_list is entered.
func (s *BasePlSqlParserListener) EnterModel_column_list(ctx *Model_column_listContext) {}

// ExitModel_column_list is called when production model_column_list is exited.
func (s *BasePlSqlParserListener) ExitModel_column_list(ctx *Model_column_listContext) {}

// EnterModel_column is called when production model_column is entered.
func (s *BasePlSqlParserListener) EnterModel_column(ctx *Model_columnContext) {}

// ExitModel_column is called when production model_column is exited.
func (s *BasePlSqlParserListener) ExitModel_column(ctx *Model_columnContext) {}

// EnterModel_rules_clause is called when production model_rules_clause is entered.
func (s *BasePlSqlParserListener) EnterModel_rules_clause(ctx *Model_rules_clauseContext) {}

// ExitModel_rules_clause is called when production model_rules_clause is exited.
func (s *BasePlSqlParserListener) ExitModel_rules_clause(ctx *Model_rules_clauseContext) {}

// EnterModel_rules_part is called when production model_rules_part is entered.
func (s *BasePlSqlParserListener) EnterModel_rules_part(ctx *Model_rules_partContext) {}

// ExitModel_rules_part is called when production model_rules_part is exited.
func (s *BasePlSqlParserListener) ExitModel_rules_part(ctx *Model_rules_partContext) {}

// EnterModel_rules_element is called when production model_rules_element is entered.
func (s *BasePlSqlParserListener) EnterModel_rules_element(ctx *Model_rules_elementContext) {}

// ExitModel_rules_element is called when production model_rules_element is exited.
func (s *BasePlSqlParserListener) ExitModel_rules_element(ctx *Model_rules_elementContext) {}

// EnterCell_assignment is called when production cell_assignment is entered.
func (s *BasePlSqlParserListener) EnterCell_assignment(ctx *Cell_assignmentContext) {}

// ExitCell_assignment is called when production cell_assignment is exited.
func (s *BasePlSqlParserListener) ExitCell_assignment(ctx *Cell_assignmentContext) {}

// EnterModel_iterate_clause is called when production model_iterate_clause is entered.
func (s *BasePlSqlParserListener) EnterModel_iterate_clause(ctx *Model_iterate_clauseContext) {}

// ExitModel_iterate_clause is called when production model_iterate_clause is exited.
func (s *BasePlSqlParserListener) ExitModel_iterate_clause(ctx *Model_iterate_clauseContext) {}

// EnterUntil_part is called when production until_part is entered.
func (s *BasePlSqlParserListener) EnterUntil_part(ctx *Until_partContext) {}

// ExitUntil_part is called when production until_part is exited.
func (s *BasePlSqlParserListener) ExitUntil_part(ctx *Until_partContext) {}

// EnterOrder_by_clause is called when production order_by_clause is entered.
func (s *BasePlSqlParserListener) EnterOrder_by_clause(ctx *Order_by_clauseContext) {}

// ExitOrder_by_clause is called when production order_by_clause is exited.
func (s *BasePlSqlParserListener) ExitOrder_by_clause(ctx *Order_by_clauseContext) {}

// EnterOrder_by_elements is called when production order_by_elements is entered.
func (s *BasePlSqlParserListener) EnterOrder_by_elements(ctx *Order_by_elementsContext) {}

// ExitOrder_by_elements is called when production order_by_elements is exited.
func (s *BasePlSqlParserListener) ExitOrder_by_elements(ctx *Order_by_elementsContext) {}

// EnterOffset_clause is called when production offset_clause is entered.
func (s *BasePlSqlParserListener) EnterOffset_clause(ctx *Offset_clauseContext) {}

// ExitOffset_clause is called when production offset_clause is exited.
func (s *BasePlSqlParserListener) ExitOffset_clause(ctx *Offset_clauseContext) {}

// EnterFetch_clause is called when production fetch_clause is entered.
func (s *BasePlSqlParserListener) EnterFetch_clause(ctx *Fetch_clauseContext) {}

// ExitFetch_clause is called when production fetch_clause is exited.
func (s *BasePlSqlParserListener) ExitFetch_clause(ctx *Fetch_clauseContext) {}

// EnterFor_update_clause is called when production for_update_clause is entered.
func (s *BasePlSqlParserListener) EnterFor_update_clause(ctx *For_update_clauseContext) {}

// ExitFor_update_clause is called when production for_update_clause is exited.
func (s *BasePlSqlParserListener) ExitFor_update_clause(ctx *For_update_clauseContext) {}

// EnterFor_update_of_part is called when production for_update_of_part is entered.
func (s *BasePlSqlParserListener) EnterFor_update_of_part(ctx *For_update_of_partContext) {}

// ExitFor_update_of_part is called when production for_update_of_part is exited.
func (s *BasePlSqlParserListener) ExitFor_update_of_part(ctx *For_update_of_partContext) {}

// EnterFor_update_options is called when production for_update_options is entered.
func (s *BasePlSqlParserListener) EnterFor_update_options(ctx *For_update_optionsContext) {}

// ExitFor_update_options is called when production for_update_options is exited.
func (s *BasePlSqlParserListener) ExitFor_update_options(ctx *For_update_optionsContext) {}

// EnterUpdate_statement is called when production update_statement is entered.
func (s *BasePlSqlParserListener) EnterUpdate_statement(ctx *Update_statementContext) {}

// ExitUpdate_statement is called when production update_statement is exited.
func (s *BasePlSqlParserListener) ExitUpdate_statement(ctx *Update_statementContext) {}

// EnterUpdate_set_clause is called when production update_set_clause is entered.
func (s *BasePlSqlParserListener) EnterUpdate_set_clause(ctx *Update_set_clauseContext) {}

// ExitUpdate_set_clause is called when production update_set_clause is exited.
func (s *BasePlSqlParserListener) ExitUpdate_set_clause(ctx *Update_set_clauseContext) {}

// EnterColumn_based_update_set_clause is called when production column_based_update_set_clause is entered.
func (s *BasePlSqlParserListener) EnterColumn_based_update_set_clause(ctx *Column_based_update_set_clauseContext) {
}

// ExitColumn_based_update_set_clause is called when production column_based_update_set_clause is exited.
func (s *BasePlSqlParserListener) ExitColumn_based_update_set_clause(ctx *Column_based_update_set_clauseContext) {
}

// EnterDelete_statement is called when production delete_statement is entered.
func (s *BasePlSqlParserListener) EnterDelete_statement(ctx *Delete_statementContext) {}

// ExitDelete_statement is called when production delete_statement is exited.
func (s *BasePlSqlParserListener) ExitDelete_statement(ctx *Delete_statementContext) {}

// EnterInsert_statement is called when production insert_statement is entered.
func (s *BasePlSqlParserListener) EnterInsert_statement(ctx *Insert_statementContext) {}

// ExitInsert_statement is called when production insert_statement is exited.
func (s *BasePlSqlParserListener) ExitInsert_statement(ctx *Insert_statementContext) {}

// EnterSingle_table_insert is called when production single_table_insert is entered.
func (s *BasePlSqlParserListener) EnterSingle_table_insert(ctx *Single_table_insertContext) {}

// ExitSingle_table_insert is called when production single_table_insert is exited.
func (s *BasePlSqlParserListener) ExitSingle_table_insert(ctx *Single_table_insertContext) {}

// EnterMulti_table_insert is called when production multi_table_insert is entered.
func (s *BasePlSqlParserListener) EnterMulti_table_insert(ctx *Multi_table_insertContext) {}

// ExitMulti_table_insert is called when production multi_table_insert is exited.
func (s *BasePlSqlParserListener) ExitMulti_table_insert(ctx *Multi_table_insertContext) {}

// EnterMulti_table_element is called when production multi_table_element is entered.
func (s *BasePlSqlParserListener) EnterMulti_table_element(ctx *Multi_table_elementContext) {}

// ExitMulti_table_element is called when production multi_table_element is exited.
func (s *BasePlSqlParserListener) ExitMulti_table_element(ctx *Multi_table_elementContext) {}

// EnterConditional_insert_clause is called when production conditional_insert_clause is entered.
func (s *BasePlSqlParserListener) EnterConditional_insert_clause(ctx *Conditional_insert_clauseContext) {
}

// ExitConditional_insert_clause is called when production conditional_insert_clause is exited.
func (s *BasePlSqlParserListener) ExitConditional_insert_clause(ctx *Conditional_insert_clauseContext) {
}

// EnterConditional_insert_when_part is called when production conditional_insert_when_part is entered.
func (s *BasePlSqlParserListener) EnterConditional_insert_when_part(ctx *Conditional_insert_when_partContext) {
}

// ExitConditional_insert_when_part is called when production conditional_insert_when_part is exited.
func (s *BasePlSqlParserListener) ExitConditional_insert_when_part(ctx *Conditional_insert_when_partContext) {
}

// EnterConditional_insert_else_part is called when production conditional_insert_else_part is entered.
func (s *BasePlSqlParserListener) EnterConditional_insert_else_part(ctx *Conditional_insert_else_partContext) {
}

// ExitConditional_insert_else_part is called when production conditional_insert_else_part is exited.
func (s *BasePlSqlParserListener) ExitConditional_insert_else_part(ctx *Conditional_insert_else_partContext) {
}

// EnterInsert_into_clause is called when production insert_into_clause is entered.
func (s *BasePlSqlParserListener) EnterInsert_into_clause(ctx *Insert_into_clauseContext) {}

// ExitInsert_into_clause is called when production insert_into_clause is exited.
func (s *BasePlSqlParserListener) ExitInsert_into_clause(ctx *Insert_into_clauseContext) {}

// EnterValues_clause is called when production values_clause is entered.
func (s *BasePlSqlParserListener) EnterValues_clause(ctx *Values_clauseContext) {}

// ExitValues_clause is called when production values_clause is exited.
func (s *BasePlSqlParserListener) ExitValues_clause(ctx *Values_clauseContext) {}

// EnterMerge_statement is called when production merge_statement is entered.
func (s *BasePlSqlParserListener) EnterMerge_statement(ctx *Merge_statementContext) {}

// ExitMerge_statement is called when production merge_statement is exited.
func (s *BasePlSqlParserListener) ExitMerge_statement(ctx *Merge_statementContext) {}

// EnterMerge_update_clause is called when production merge_update_clause is entered.
func (s *BasePlSqlParserListener) EnterMerge_update_clause(ctx *Merge_update_clauseContext) {}

// ExitMerge_update_clause is called when production merge_update_clause is exited.
func (s *BasePlSqlParserListener) ExitMerge_update_clause(ctx *Merge_update_clauseContext) {}

// EnterMerge_element is called when production merge_element is entered.
func (s *BasePlSqlParserListener) EnterMerge_element(ctx *Merge_elementContext) {}

// ExitMerge_element is called when production merge_element is exited.
func (s *BasePlSqlParserListener) ExitMerge_element(ctx *Merge_elementContext) {}

// EnterMerge_update_delete_part is called when production merge_update_delete_part is entered.
func (s *BasePlSqlParserListener) EnterMerge_update_delete_part(ctx *Merge_update_delete_partContext) {
}

// ExitMerge_update_delete_part is called when production merge_update_delete_part is exited.
func (s *BasePlSqlParserListener) ExitMerge_update_delete_part(ctx *Merge_update_delete_partContext) {
}

// EnterMerge_insert_clause is called when production merge_insert_clause is entered.
func (s *BasePlSqlParserListener) EnterMerge_insert_clause(ctx *Merge_insert_clauseContext) {}

// ExitMerge_insert_clause is called when production merge_insert_clause is exited.
func (s *BasePlSqlParserListener) ExitMerge_insert_clause(ctx *Merge_insert_clauseContext) {}

// EnterSelected_tableview is called when production selected_tableview is entered.
func (s *BasePlSqlParserListener) EnterSelected_tableview(ctx *Selected_tableviewContext) {}

// ExitSelected_tableview is called when production selected_tableview is exited.
func (s *BasePlSqlParserListener) ExitSelected_tableview(ctx *Selected_tableviewContext) {}

// EnterLock_table_statement is called when production lock_table_statement is entered.
func (s *BasePlSqlParserListener) EnterLock_table_statement(ctx *Lock_table_statementContext) {}

// ExitLock_table_statement is called when production lock_table_statement is exited.
func (s *BasePlSqlParserListener) ExitLock_table_statement(ctx *Lock_table_statementContext) {}

// EnterWait_nowait_part is called when production wait_nowait_part is entered.
func (s *BasePlSqlParserListener) EnterWait_nowait_part(ctx *Wait_nowait_partContext) {}

// ExitWait_nowait_part is called when production wait_nowait_part is exited.
func (s *BasePlSqlParserListener) ExitWait_nowait_part(ctx *Wait_nowait_partContext) {}

// EnterLock_table_element is called when production lock_table_element is entered.
func (s *BasePlSqlParserListener) EnterLock_table_element(ctx *Lock_table_elementContext) {}

// ExitLock_table_element is called when production lock_table_element is exited.
func (s *BasePlSqlParserListener) ExitLock_table_element(ctx *Lock_table_elementContext) {}

// EnterLock_mode is called when production lock_mode is entered.
func (s *BasePlSqlParserListener) EnterLock_mode(ctx *Lock_modeContext) {}

// ExitLock_mode is called when production lock_mode is exited.
func (s *BasePlSqlParserListener) ExitLock_mode(ctx *Lock_modeContext) {}

// EnterGeneral_table_ref is called when production general_table_ref is entered.
func (s *BasePlSqlParserListener) EnterGeneral_table_ref(ctx *General_table_refContext) {}

// ExitGeneral_table_ref is called when production general_table_ref is exited.
func (s *BasePlSqlParserListener) ExitGeneral_table_ref(ctx *General_table_refContext) {}

// EnterStatic_returning_clause is called when production static_returning_clause is entered.
func (s *BasePlSqlParserListener) EnterStatic_returning_clause(ctx *Static_returning_clauseContext) {}

// ExitStatic_returning_clause is called when production static_returning_clause is exited.
func (s *BasePlSqlParserListener) ExitStatic_returning_clause(ctx *Static_returning_clauseContext) {}

// EnterError_logging_clause is called when production error_logging_clause is entered.
func (s *BasePlSqlParserListener) EnterError_logging_clause(ctx *Error_logging_clauseContext) {}

// ExitError_logging_clause is called when production error_logging_clause is exited.
func (s *BasePlSqlParserListener) ExitError_logging_clause(ctx *Error_logging_clauseContext) {}

// EnterError_logging_into_part is called when production error_logging_into_part is entered.
func (s *BasePlSqlParserListener) EnterError_logging_into_part(ctx *Error_logging_into_partContext) {}

// ExitError_logging_into_part is called when production error_logging_into_part is exited.
func (s *BasePlSqlParserListener) ExitError_logging_into_part(ctx *Error_logging_into_partContext) {}

// EnterError_logging_reject_part is called when production error_logging_reject_part is entered.
func (s *BasePlSqlParserListener) EnterError_logging_reject_part(ctx *Error_logging_reject_partContext) {
}

// ExitError_logging_reject_part is called when production error_logging_reject_part is exited.
func (s *BasePlSqlParserListener) ExitError_logging_reject_part(ctx *Error_logging_reject_partContext) {
}

// EnterDml_table_expression_clause is called when production dml_table_expression_clause is entered.
func (s *BasePlSqlParserListener) EnterDml_table_expression_clause(ctx *Dml_table_expression_clauseContext) {
}

// ExitDml_table_expression_clause is called when production dml_table_expression_clause is exited.
func (s *BasePlSqlParserListener) ExitDml_table_expression_clause(ctx *Dml_table_expression_clauseContext) {
}

// EnterTable_collection_expression is called when production table_collection_expression is entered.
func (s *BasePlSqlParserListener) EnterTable_collection_expression(ctx *Table_collection_expressionContext) {
}

// ExitTable_collection_expression is called when production table_collection_expression is exited.
func (s *BasePlSqlParserListener) ExitTable_collection_expression(ctx *Table_collection_expressionContext) {
}

// EnterSubquery_restriction_clause is called when production subquery_restriction_clause is entered.
func (s *BasePlSqlParserListener) EnterSubquery_restriction_clause(ctx *Subquery_restriction_clauseContext) {
}

// ExitSubquery_restriction_clause is called when production subquery_restriction_clause is exited.
func (s *BasePlSqlParserListener) ExitSubquery_restriction_clause(ctx *Subquery_restriction_clauseContext) {
}

// EnterSample_clause is called when production sample_clause is entered.
func (s *BasePlSqlParserListener) EnterSample_clause(ctx *Sample_clauseContext) {}

// ExitSample_clause is called when production sample_clause is exited.
func (s *BasePlSqlParserListener) ExitSample_clause(ctx *Sample_clauseContext) {}

// EnterSeed_part is called when production seed_part is entered.
func (s *BasePlSqlParserListener) EnterSeed_part(ctx *Seed_partContext) {}

// ExitSeed_part is called when production seed_part is exited.
func (s *BasePlSqlParserListener) ExitSeed_part(ctx *Seed_partContext) {}

// EnterCondition is called when production condition is entered.
func (s *BasePlSqlParserListener) EnterCondition(ctx *ConditionContext) {}

// ExitCondition is called when production condition is exited.
func (s *BasePlSqlParserListener) ExitCondition(ctx *ConditionContext) {}

// EnterExpressions_ is called when production expressions_ is entered.
func (s *BasePlSqlParserListener) EnterExpressions_(ctx *Expressions_Context) {}

// ExitExpressions_ is called when production expressions_ is exited.
func (s *BasePlSqlParserListener) ExitExpressions_(ctx *Expressions_Context) {}

// EnterExpression is called when production expression is entered.
func (s *BasePlSqlParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasePlSqlParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterCursor_expression is called when production cursor_expression is entered.
func (s *BasePlSqlParserListener) EnterCursor_expression(ctx *Cursor_expressionContext) {}

// ExitCursor_expression is called when production cursor_expression is exited.
func (s *BasePlSqlParserListener) ExitCursor_expression(ctx *Cursor_expressionContext) {}

// EnterLogical_expression is called when production logical_expression is entered.
func (s *BasePlSqlParserListener) EnterLogical_expression(ctx *Logical_expressionContext) {}

// ExitLogical_expression is called when production logical_expression is exited.
func (s *BasePlSqlParserListener) ExitLogical_expression(ctx *Logical_expressionContext) {}

// EnterUnary_logical_expression is called when production unary_logical_expression is entered.
func (s *BasePlSqlParserListener) EnterUnary_logical_expression(ctx *Unary_logical_expressionContext) {
}

// ExitUnary_logical_expression is called when production unary_logical_expression is exited.
func (s *BasePlSqlParserListener) ExitUnary_logical_expression(ctx *Unary_logical_expressionContext) {
}

// EnterUnary_logical_operation is called when production unary_logical_operation is entered.
func (s *BasePlSqlParserListener) EnterUnary_logical_operation(ctx *Unary_logical_operationContext) {}

// ExitUnary_logical_operation is called when production unary_logical_operation is exited.
func (s *BasePlSqlParserListener) ExitUnary_logical_operation(ctx *Unary_logical_operationContext) {}

// EnterLogical_operation is called when production logical_operation is entered.
func (s *BasePlSqlParserListener) EnterLogical_operation(ctx *Logical_operationContext) {}

// ExitLogical_operation is called when production logical_operation is exited.
func (s *BasePlSqlParserListener) ExitLogical_operation(ctx *Logical_operationContext) {}

// EnterMultiset_expression is called when production multiset_expression is entered.
func (s *BasePlSqlParserListener) EnterMultiset_expression(ctx *Multiset_expressionContext) {}

// ExitMultiset_expression is called when production multiset_expression is exited.
func (s *BasePlSqlParserListener) ExitMultiset_expression(ctx *Multiset_expressionContext) {}

// EnterRelational_expression is called when production relational_expression is entered.
func (s *BasePlSqlParserListener) EnterRelational_expression(ctx *Relational_expressionContext) {}

// ExitRelational_expression is called when production relational_expression is exited.
func (s *BasePlSqlParserListener) ExitRelational_expression(ctx *Relational_expressionContext) {}

// EnterCompound_expression is called when production compound_expression is entered.
func (s *BasePlSqlParserListener) EnterCompound_expression(ctx *Compound_expressionContext) {}

// ExitCompound_expression is called when production compound_expression is exited.
func (s *BasePlSqlParserListener) ExitCompound_expression(ctx *Compound_expressionContext) {}

// EnterRelational_operator is called when production relational_operator is entered.
func (s *BasePlSqlParserListener) EnterRelational_operator(ctx *Relational_operatorContext) {}

// ExitRelational_operator is called when production relational_operator is exited.
func (s *BasePlSqlParserListener) ExitRelational_operator(ctx *Relational_operatorContext) {}

// EnterIn_elements is called when production in_elements is entered.
func (s *BasePlSqlParserListener) EnterIn_elements(ctx *In_elementsContext) {}

// ExitIn_elements is called when production in_elements is exited.
func (s *BasePlSqlParserListener) ExitIn_elements(ctx *In_elementsContext) {}

// EnterBetween_elements is called when production between_elements is entered.
func (s *BasePlSqlParserListener) EnterBetween_elements(ctx *Between_elementsContext) {}

// ExitBetween_elements is called when production between_elements is exited.
func (s *BasePlSqlParserListener) ExitBetween_elements(ctx *Between_elementsContext) {}

// EnterConcatenation is called when production concatenation is entered.
func (s *BasePlSqlParserListener) EnterConcatenation(ctx *ConcatenationContext) {}

// ExitConcatenation is called when production concatenation is exited.
func (s *BasePlSqlParserListener) ExitConcatenation(ctx *ConcatenationContext) {}

// EnterInterval_expression is called when production interval_expression is entered.
func (s *BasePlSqlParserListener) EnterInterval_expression(ctx *Interval_expressionContext) {}

// ExitInterval_expression is called when production interval_expression is exited.
func (s *BasePlSqlParserListener) ExitInterval_expression(ctx *Interval_expressionContext) {}

// EnterModel_expression is called when production model_expression is entered.
func (s *BasePlSqlParserListener) EnterModel_expression(ctx *Model_expressionContext) {}

// ExitModel_expression is called when production model_expression is exited.
func (s *BasePlSqlParserListener) ExitModel_expression(ctx *Model_expressionContext) {}

// EnterModel_expression_element is called when production model_expression_element is entered.
func (s *BasePlSqlParserListener) EnterModel_expression_element(ctx *Model_expression_elementContext) {
}

// ExitModel_expression_element is called when production model_expression_element is exited.
func (s *BasePlSqlParserListener) ExitModel_expression_element(ctx *Model_expression_elementContext) {
}

// EnterSingle_column_for_loop is called when production single_column_for_loop is entered.
func (s *BasePlSqlParserListener) EnterSingle_column_for_loop(ctx *Single_column_for_loopContext) {}

// ExitSingle_column_for_loop is called when production single_column_for_loop is exited.
func (s *BasePlSqlParserListener) ExitSingle_column_for_loop(ctx *Single_column_for_loopContext) {}

// EnterMulti_column_for_loop is called when production multi_column_for_loop is entered.
func (s *BasePlSqlParserListener) EnterMulti_column_for_loop(ctx *Multi_column_for_loopContext) {}

// ExitMulti_column_for_loop is called when production multi_column_for_loop is exited.
func (s *BasePlSqlParserListener) ExitMulti_column_for_loop(ctx *Multi_column_for_loopContext) {}

// EnterUnary_expression is called when production unary_expression is entered.
func (s *BasePlSqlParserListener) EnterUnary_expression(ctx *Unary_expressionContext) {}

// ExitUnary_expression is called when production unary_expression is exited.
func (s *BasePlSqlParserListener) ExitUnary_expression(ctx *Unary_expressionContext) {}

// EnterUnary_expression_core is called when production unary_expression_core is entered.
func (s *BasePlSqlParserListener) EnterUnary_expression_core(ctx *Unary_expression_coreContext) {}

// ExitUnary_expression_core is called when production unary_expression_core is exited.
func (s *BasePlSqlParserListener) ExitUnary_expression_core(ctx *Unary_expression_coreContext) {}

// EnterImplicit_cursor_expression is called when production implicit_cursor_expression is entered.
func (s *BasePlSqlParserListener) EnterImplicit_cursor_expression(ctx *Implicit_cursor_expressionContext) {
}

// ExitImplicit_cursor_expression is called when production implicit_cursor_expression is exited.
func (s *BasePlSqlParserListener) ExitImplicit_cursor_expression(ctx *Implicit_cursor_expressionContext) {
}

// EnterCollection_expression is called when production collection_expression is entered.
func (s *BasePlSqlParserListener) EnterCollection_expression(ctx *Collection_expressionContext) {}

// ExitCollection_expression is called when production collection_expression is exited.
func (s *BasePlSqlParserListener) ExitCollection_expression(ctx *Collection_expressionContext) {}

// EnterCase_statement is called when production case_statement is entered.
func (s *BasePlSqlParserListener) EnterCase_statement(ctx *Case_statementContext) {}

// ExitCase_statement is called when production case_statement is exited.
func (s *BasePlSqlParserListener) ExitCase_statement(ctx *Case_statementContext) {}

// EnterSimple_case_statement is called when production simple_case_statement is entered.
func (s *BasePlSqlParserListener) EnterSimple_case_statement(ctx *Simple_case_statementContext) {}

// ExitSimple_case_statement is called when production simple_case_statement is exited.
func (s *BasePlSqlParserListener) ExitSimple_case_statement(ctx *Simple_case_statementContext) {}

// EnterSearched_case_statement is called when production searched_case_statement is entered.
func (s *BasePlSqlParserListener) EnterSearched_case_statement(ctx *Searched_case_statementContext) {}

// ExitSearched_case_statement is called when production searched_case_statement is exited.
func (s *BasePlSqlParserListener) ExitSearched_case_statement(ctx *Searched_case_statementContext) {}

// EnterCase_when_part_statement is called when production case_when_part_statement is entered.
func (s *BasePlSqlParserListener) EnterCase_when_part_statement(ctx *Case_when_part_statementContext) {
}

// ExitCase_when_part_statement is called when production case_when_part_statement is exited.
func (s *BasePlSqlParserListener) ExitCase_when_part_statement(ctx *Case_when_part_statementContext) {
}

// EnterCase_else_part_statement is called when production case_else_part_statement is entered.
func (s *BasePlSqlParserListener) EnterCase_else_part_statement(ctx *Case_else_part_statementContext) {
}

// ExitCase_else_part_statement is called when production case_else_part_statement is exited.
func (s *BasePlSqlParserListener) ExitCase_else_part_statement(ctx *Case_else_part_statementContext) {
}

// EnterCase_expression is called when production case_expression is entered.
func (s *BasePlSqlParserListener) EnterCase_expression(ctx *Case_expressionContext) {}

// ExitCase_expression is called when production case_expression is exited.
func (s *BasePlSqlParserListener) ExitCase_expression(ctx *Case_expressionContext) {}

// EnterSimple_case_expression is called when production simple_case_expression is entered.
func (s *BasePlSqlParserListener) EnterSimple_case_expression(ctx *Simple_case_expressionContext) {}

// ExitSimple_case_expression is called when production simple_case_expression is exited.
func (s *BasePlSqlParserListener) ExitSimple_case_expression(ctx *Simple_case_expressionContext) {}

// EnterSearched_case_expression is called when production searched_case_expression is entered.
func (s *BasePlSqlParserListener) EnterSearched_case_expression(ctx *Searched_case_expressionContext) {
}

// ExitSearched_case_expression is called when production searched_case_expression is exited.
func (s *BasePlSqlParserListener) ExitSearched_case_expression(ctx *Searched_case_expressionContext) {
}

// EnterCase_when_part_expression is called when production case_when_part_expression is entered.
func (s *BasePlSqlParserListener) EnterCase_when_part_expression(ctx *Case_when_part_expressionContext) {
}

// ExitCase_when_part_expression is called when production case_when_part_expression is exited.
func (s *BasePlSqlParserListener) ExitCase_when_part_expression(ctx *Case_when_part_expressionContext) {
}

// EnterCase_else_part_expression is called when production case_else_part_expression is entered.
func (s *BasePlSqlParserListener) EnterCase_else_part_expression(ctx *Case_else_part_expressionContext) {
}

// ExitCase_else_part_expression is called when production case_else_part_expression is exited.
func (s *BasePlSqlParserListener) ExitCase_else_part_expression(ctx *Case_else_part_expressionContext) {
}

// EnterAtom is called when production atom is entered.
func (s *BasePlSqlParserListener) EnterAtom(ctx *AtomContext) {}

// ExitAtom is called when production atom is exited.
func (s *BasePlSqlParserListener) ExitAtom(ctx *AtomContext) {}

// EnterQuantified_expression is called when production quantified_expression is entered.
func (s *BasePlSqlParserListener) EnterQuantified_expression(ctx *Quantified_expressionContext) {}

// ExitQuantified_expression is called when production quantified_expression is exited.
func (s *BasePlSqlParserListener) ExitQuantified_expression(ctx *Quantified_expressionContext) {}

// EnterString_function is called when production string_function is entered.
func (s *BasePlSqlParserListener) EnterString_function(ctx *String_functionContext) {}

// ExitString_function is called when production string_function is exited.
func (s *BasePlSqlParserListener) ExitString_function(ctx *String_functionContext) {}

// EnterStandard_function is called when production standard_function is entered.
func (s *BasePlSqlParserListener) EnterStandard_function(ctx *Standard_functionContext) {}

// ExitStandard_function is called when production standard_function is exited.
func (s *BasePlSqlParserListener) ExitStandard_function(ctx *Standard_functionContext) {}

// EnterJson_function is called when production json_function is entered.
func (s *BasePlSqlParserListener) EnterJson_function(ctx *Json_functionContext) {}

// ExitJson_function is called when production json_function is exited.
func (s *BasePlSqlParserListener) ExitJson_function(ctx *Json_functionContext) {}

// EnterJson_object_content is called when production json_object_content is entered.
func (s *BasePlSqlParserListener) EnterJson_object_content(ctx *Json_object_contentContext) {}

// ExitJson_object_content is called when production json_object_content is exited.
func (s *BasePlSqlParserListener) ExitJson_object_content(ctx *Json_object_contentContext) {}

// EnterJson_object_entry is called when production json_object_entry is entered.
func (s *BasePlSqlParserListener) EnterJson_object_entry(ctx *Json_object_entryContext) {}

// ExitJson_object_entry is called when production json_object_entry is exited.
func (s *BasePlSqlParserListener) ExitJson_object_entry(ctx *Json_object_entryContext) {}

// EnterJson_table_clause is called when production json_table_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_table_clause(ctx *Json_table_clauseContext) {}

// ExitJson_table_clause is called when production json_table_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_table_clause(ctx *Json_table_clauseContext) {}

// EnterJson_array_element is called when production json_array_element is entered.
func (s *BasePlSqlParserListener) EnterJson_array_element(ctx *Json_array_elementContext) {}

// ExitJson_array_element is called when production json_array_element is exited.
func (s *BasePlSqlParserListener) ExitJson_array_element(ctx *Json_array_elementContext) {}

// EnterJson_on_null_clause is called when production json_on_null_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_on_null_clause(ctx *Json_on_null_clauseContext) {}

// ExitJson_on_null_clause is called when production json_on_null_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_on_null_clause(ctx *Json_on_null_clauseContext) {}

// EnterJson_return_clause is called when production json_return_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_return_clause(ctx *Json_return_clauseContext) {}

// ExitJson_return_clause is called when production json_return_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_return_clause(ctx *Json_return_clauseContext) {}

// EnterJson_transform_op is called when production json_transform_op is entered.
func (s *BasePlSqlParserListener) EnterJson_transform_op(ctx *Json_transform_opContext) {}

// ExitJson_transform_op is called when production json_transform_op is exited.
func (s *BasePlSqlParserListener) ExitJson_transform_op(ctx *Json_transform_opContext) {}

// EnterJson_column_clause is called when production json_column_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_column_clause(ctx *Json_column_clauseContext) {}

// ExitJson_column_clause is called when production json_column_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_column_clause(ctx *Json_column_clauseContext) {}

// EnterJson_column_definition is called when production json_column_definition is entered.
func (s *BasePlSqlParserListener) EnterJson_column_definition(ctx *Json_column_definitionContext) {}

// ExitJson_column_definition is called when production json_column_definition is exited.
func (s *BasePlSqlParserListener) ExitJson_column_definition(ctx *Json_column_definitionContext) {}

// EnterJson_query_returning_clause is called when production json_query_returning_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_query_returning_clause(ctx *Json_query_returning_clauseContext) {
}

// ExitJson_query_returning_clause is called when production json_query_returning_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_query_returning_clause(ctx *Json_query_returning_clauseContext) {
}

// EnterJson_query_return_type is called when production json_query_return_type is entered.
func (s *BasePlSqlParserListener) EnterJson_query_return_type(ctx *Json_query_return_typeContext) {}

// ExitJson_query_return_type is called when production json_query_return_type is exited.
func (s *BasePlSqlParserListener) ExitJson_query_return_type(ctx *Json_query_return_typeContext) {}

// EnterJson_query_wrapper_clause is called when production json_query_wrapper_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_query_wrapper_clause(ctx *Json_query_wrapper_clauseContext) {
}

// ExitJson_query_wrapper_clause is called when production json_query_wrapper_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_query_wrapper_clause(ctx *Json_query_wrapper_clauseContext) {
}

// EnterJson_query_on_error_clause is called when production json_query_on_error_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_query_on_error_clause(ctx *Json_query_on_error_clauseContext) {
}

// ExitJson_query_on_error_clause is called when production json_query_on_error_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_query_on_error_clause(ctx *Json_query_on_error_clauseContext) {
}

// EnterJson_query_on_empty_clause is called when production json_query_on_empty_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_query_on_empty_clause(ctx *Json_query_on_empty_clauseContext) {
}

// ExitJson_query_on_empty_clause is called when production json_query_on_empty_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_query_on_empty_clause(ctx *Json_query_on_empty_clauseContext) {
}

// EnterJson_value_return_clause is called when production json_value_return_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_value_return_clause(ctx *Json_value_return_clauseContext) {
}

// ExitJson_value_return_clause is called when production json_value_return_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_value_return_clause(ctx *Json_value_return_clauseContext) {
}

// EnterJson_value_return_type is called when production json_value_return_type is entered.
func (s *BasePlSqlParserListener) EnterJson_value_return_type(ctx *Json_value_return_typeContext) {}

// ExitJson_value_return_type is called when production json_value_return_type is exited.
func (s *BasePlSqlParserListener) ExitJson_value_return_type(ctx *Json_value_return_typeContext) {}

// EnterJson_value_on_mismatch_clause is called when production json_value_on_mismatch_clause is entered.
func (s *BasePlSqlParserListener) EnterJson_value_on_mismatch_clause(ctx *Json_value_on_mismatch_clauseContext) {
}

// ExitJson_value_on_mismatch_clause is called when production json_value_on_mismatch_clause is exited.
func (s *BasePlSqlParserListener) ExitJson_value_on_mismatch_clause(ctx *Json_value_on_mismatch_clauseContext) {
}

// EnterLiteral is called when production literal is entered.
func (s *BasePlSqlParserListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BasePlSqlParserListener) ExitLiteral(ctx *LiteralContext) {}

// EnterNumeric_function_wrapper is called when production numeric_function_wrapper is entered.
func (s *BasePlSqlParserListener) EnterNumeric_function_wrapper(ctx *Numeric_function_wrapperContext) {
}

// ExitNumeric_function_wrapper is called when production numeric_function_wrapper is exited.
func (s *BasePlSqlParserListener) ExitNumeric_function_wrapper(ctx *Numeric_function_wrapperContext) {
}

// EnterNumeric_function is called when production numeric_function is entered.
func (s *BasePlSqlParserListener) EnterNumeric_function(ctx *Numeric_functionContext) {}

// ExitNumeric_function is called when production numeric_function is exited.
func (s *BasePlSqlParserListener) ExitNumeric_function(ctx *Numeric_functionContext) {}

// EnterListagg_overflow_clause is called when production listagg_overflow_clause is entered.
func (s *BasePlSqlParserListener) EnterListagg_overflow_clause(ctx *Listagg_overflow_clauseContext) {}

// ExitListagg_overflow_clause is called when production listagg_overflow_clause is exited.
func (s *BasePlSqlParserListener) ExitListagg_overflow_clause(ctx *Listagg_overflow_clauseContext) {}

// EnterOther_function is called when production other_function is entered.
func (s *BasePlSqlParserListener) EnterOther_function(ctx *Other_functionContext) {}

// ExitOther_function is called when production other_function is exited.
func (s *BasePlSqlParserListener) ExitOther_function(ctx *Other_functionContext) {}

// EnterOver_clause_keyword is called when production over_clause_keyword is entered.
func (s *BasePlSqlParserListener) EnterOver_clause_keyword(ctx *Over_clause_keywordContext) {}

// ExitOver_clause_keyword is called when production over_clause_keyword is exited.
func (s *BasePlSqlParserListener) ExitOver_clause_keyword(ctx *Over_clause_keywordContext) {}

// EnterWithin_or_over_clause_keyword is called when production within_or_over_clause_keyword is entered.
func (s *BasePlSqlParserListener) EnterWithin_or_over_clause_keyword(ctx *Within_or_over_clause_keywordContext) {
}

// ExitWithin_or_over_clause_keyword is called when production within_or_over_clause_keyword is exited.
func (s *BasePlSqlParserListener) ExitWithin_or_over_clause_keyword(ctx *Within_or_over_clause_keywordContext) {
}

// EnterStandard_prediction_function_keyword is called when production standard_prediction_function_keyword is entered.
func (s *BasePlSqlParserListener) EnterStandard_prediction_function_keyword(ctx *Standard_prediction_function_keywordContext) {
}

// ExitStandard_prediction_function_keyword is called when production standard_prediction_function_keyword is exited.
func (s *BasePlSqlParserListener) ExitStandard_prediction_function_keyword(ctx *Standard_prediction_function_keywordContext) {
}

// EnterOver_clause is called when production over_clause is entered.
func (s *BasePlSqlParserListener) EnterOver_clause(ctx *Over_clauseContext) {}

// ExitOver_clause is called when production over_clause is exited.
func (s *BasePlSqlParserListener) ExitOver_clause(ctx *Over_clauseContext) {}

// EnterWindowing_clause is called when production windowing_clause is entered.
func (s *BasePlSqlParserListener) EnterWindowing_clause(ctx *Windowing_clauseContext) {}

// ExitWindowing_clause is called when production windowing_clause is exited.
func (s *BasePlSqlParserListener) ExitWindowing_clause(ctx *Windowing_clauseContext) {}

// EnterWindowing_type is called when production windowing_type is entered.
func (s *BasePlSqlParserListener) EnterWindowing_type(ctx *Windowing_typeContext) {}

// ExitWindowing_type is called when production windowing_type is exited.
func (s *BasePlSqlParserListener) ExitWindowing_type(ctx *Windowing_typeContext) {}

// EnterWindowing_elements is called when production windowing_elements is entered.
func (s *BasePlSqlParserListener) EnterWindowing_elements(ctx *Windowing_elementsContext) {}

// ExitWindowing_elements is called when production windowing_elements is exited.
func (s *BasePlSqlParserListener) ExitWindowing_elements(ctx *Windowing_elementsContext) {}

// EnterUsing_clause is called when production using_clause is entered.
func (s *BasePlSqlParserListener) EnterUsing_clause(ctx *Using_clauseContext) {}

// ExitUsing_clause is called when production using_clause is exited.
func (s *BasePlSqlParserListener) ExitUsing_clause(ctx *Using_clauseContext) {}

// EnterUsing_element is called when production using_element is entered.
func (s *BasePlSqlParserListener) EnterUsing_element(ctx *Using_elementContext) {}

// ExitUsing_element is called when production using_element is exited.
func (s *BasePlSqlParserListener) ExitUsing_element(ctx *Using_elementContext) {}

// EnterAssignable_element is called when production assignable_element is entered.
func (s *BasePlSqlParserListener) EnterAssignable_element(ctx *Assignable_elementContext) {}

// ExitAssignable_element is called when production assignable_element is exited.
func (s *BasePlSqlParserListener) ExitAssignable_element(ctx *Assignable_elementContext) {}

// EnterCollect_order_by_part is called when production collect_order_by_part is entered.
func (s *BasePlSqlParserListener) EnterCollect_order_by_part(ctx *Collect_order_by_partContext) {}

// ExitCollect_order_by_part is called when production collect_order_by_part is exited.
func (s *BasePlSqlParserListener) ExitCollect_order_by_part(ctx *Collect_order_by_partContext) {}

// EnterWithin_or_over_part is called when production within_or_over_part is entered.
func (s *BasePlSqlParserListener) EnterWithin_or_over_part(ctx *Within_or_over_partContext) {}

// ExitWithin_or_over_part is called when production within_or_over_part is exited.
func (s *BasePlSqlParserListener) ExitWithin_or_over_part(ctx *Within_or_over_partContext) {}

// EnterString_delimiter is called when production string_delimiter is entered.
func (s *BasePlSqlParserListener) EnterString_delimiter(ctx *String_delimiterContext) {}

// ExitString_delimiter is called when production string_delimiter is exited.
func (s *BasePlSqlParserListener) ExitString_delimiter(ctx *String_delimiterContext) {}

// EnterCost_matrix_clause is called when production cost_matrix_clause is entered.
func (s *BasePlSqlParserListener) EnterCost_matrix_clause(ctx *Cost_matrix_clauseContext) {}

// ExitCost_matrix_clause is called when production cost_matrix_clause is exited.
func (s *BasePlSqlParserListener) ExitCost_matrix_clause(ctx *Cost_matrix_clauseContext) {}

// EnterXml_passing_clause is called when production xml_passing_clause is entered.
func (s *BasePlSqlParserListener) EnterXml_passing_clause(ctx *Xml_passing_clauseContext) {}

// ExitXml_passing_clause is called when production xml_passing_clause is exited.
func (s *BasePlSqlParserListener) ExitXml_passing_clause(ctx *Xml_passing_clauseContext) {}

// EnterXml_attributes_clause is called when production xml_attributes_clause is entered.
func (s *BasePlSqlParserListener) EnterXml_attributes_clause(ctx *Xml_attributes_clauseContext) {}

// ExitXml_attributes_clause is called when production xml_attributes_clause is exited.
func (s *BasePlSqlParserListener) ExitXml_attributes_clause(ctx *Xml_attributes_clauseContext) {}

// EnterXml_namespaces_clause is called when production xml_namespaces_clause is entered.
func (s *BasePlSqlParserListener) EnterXml_namespaces_clause(ctx *Xml_namespaces_clauseContext) {}

// ExitXml_namespaces_clause is called when production xml_namespaces_clause is exited.
func (s *BasePlSqlParserListener) ExitXml_namespaces_clause(ctx *Xml_namespaces_clauseContext) {}

// EnterXml_table_column is called when production xml_table_column is entered.
func (s *BasePlSqlParserListener) EnterXml_table_column(ctx *Xml_table_columnContext) {}

// ExitXml_table_column is called when production xml_table_column is exited.
func (s *BasePlSqlParserListener) ExitXml_table_column(ctx *Xml_table_columnContext) {}

// EnterXml_general_default_part is called when production xml_general_default_part is entered.
func (s *BasePlSqlParserListener) EnterXml_general_default_part(ctx *Xml_general_default_partContext) {
}

// ExitXml_general_default_part is called when production xml_general_default_part is exited.
func (s *BasePlSqlParserListener) ExitXml_general_default_part(ctx *Xml_general_default_partContext) {
}

// EnterXml_multiuse_expression_element is called when production xml_multiuse_expression_element is entered.
func (s *BasePlSqlParserListener) EnterXml_multiuse_expression_element(ctx *Xml_multiuse_expression_elementContext) {
}

// ExitXml_multiuse_expression_element is called when production xml_multiuse_expression_element is exited.
func (s *BasePlSqlParserListener) ExitXml_multiuse_expression_element(ctx *Xml_multiuse_expression_elementContext) {
}

// EnterXmlroot_param_version_part is called when production xmlroot_param_version_part is entered.
func (s *BasePlSqlParserListener) EnterXmlroot_param_version_part(ctx *Xmlroot_param_version_partContext) {
}

// ExitXmlroot_param_version_part is called when production xmlroot_param_version_part is exited.
func (s *BasePlSqlParserListener) ExitXmlroot_param_version_part(ctx *Xmlroot_param_version_partContext) {
}

// EnterXmlroot_param_standalone_part is called when production xmlroot_param_standalone_part is entered.
func (s *BasePlSqlParserListener) EnterXmlroot_param_standalone_part(ctx *Xmlroot_param_standalone_partContext) {
}

// ExitXmlroot_param_standalone_part is called when production xmlroot_param_standalone_part is exited.
func (s *BasePlSqlParserListener) ExitXmlroot_param_standalone_part(ctx *Xmlroot_param_standalone_partContext) {
}

// EnterXmlserialize_param_enconding_part is called when production xmlserialize_param_enconding_part is entered.
func (s *BasePlSqlParserListener) EnterXmlserialize_param_enconding_part(ctx *Xmlserialize_param_enconding_partContext) {
}

// ExitXmlserialize_param_enconding_part is called when production xmlserialize_param_enconding_part is exited.
func (s *BasePlSqlParserListener) ExitXmlserialize_param_enconding_part(ctx *Xmlserialize_param_enconding_partContext) {
}

// EnterXmlserialize_param_version_part is called when production xmlserialize_param_version_part is entered.
func (s *BasePlSqlParserListener) EnterXmlserialize_param_version_part(ctx *Xmlserialize_param_version_partContext) {
}

// ExitXmlserialize_param_version_part is called when production xmlserialize_param_version_part is exited.
func (s *BasePlSqlParserListener) ExitXmlserialize_param_version_part(ctx *Xmlserialize_param_version_partContext) {
}

// EnterXmlserialize_param_ident_part is called when production xmlserialize_param_ident_part is entered.
func (s *BasePlSqlParserListener) EnterXmlserialize_param_ident_part(ctx *Xmlserialize_param_ident_partContext) {
}

// ExitXmlserialize_param_ident_part is called when production xmlserialize_param_ident_part is exited.
func (s *BasePlSqlParserListener) ExitXmlserialize_param_ident_part(ctx *Xmlserialize_param_ident_partContext) {
}

// EnterAnnotations_clause is called when production annotations_clause is entered.
func (s *BasePlSqlParserListener) EnterAnnotations_clause(ctx *Annotations_clauseContext) {}

// ExitAnnotations_clause is called when production annotations_clause is exited.
func (s *BasePlSqlParserListener) ExitAnnotations_clause(ctx *Annotations_clauseContext) {}

// EnterAnnotations_list is called when production annotations_list is entered.
func (s *BasePlSqlParserListener) EnterAnnotations_list(ctx *Annotations_listContext) {}

// ExitAnnotations_list is called when production annotations_list is exited.
func (s *BasePlSqlParserListener) ExitAnnotations_list(ctx *Annotations_listContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *BasePlSqlParserListener) EnterAnnotation(ctx *AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *BasePlSqlParserListener) ExitAnnotation(ctx *AnnotationContext) {}

// EnterSql_plus_command is called when production sql_plus_command is entered.
func (s *BasePlSqlParserListener) EnterSql_plus_command(ctx *Sql_plus_commandContext) {}

// ExitSql_plus_command is called when production sql_plus_command is exited.
func (s *BasePlSqlParserListener) ExitSql_plus_command(ctx *Sql_plus_commandContext) {}

// EnterStart_command is called when production start_command is entered.
func (s *BasePlSqlParserListener) EnterStart_command(ctx *Start_commandContext) {}

// ExitStart_command is called when production start_command is exited.
func (s *BasePlSqlParserListener) ExitStart_command(ctx *Start_commandContext) {}

// EnterSql_plus_filepath is called when production sql_plus_filepath is entered.
func (s *BasePlSqlParserListener) EnterSql_plus_filepath(ctx *Sql_plus_filepathContext) {}

// ExitSql_plus_filepath is called when production sql_plus_filepath is exited.
func (s *BasePlSqlParserListener) ExitSql_plus_filepath(ctx *Sql_plus_filepathContext) {}

// EnterWhenever_command is called when production whenever_command is entered.
func (s *BasePlSqlParserListener) EnterWhenever_command(ctx *Whenever_commandContext) {}

// ExitWhenever_command is called when production whenever_command is exited.
func (s *BasePlSqlParserListener) ExitWhenever_command(ctx *Whenever_commandContext) {}

// EnterSet_command is called when production set_command is entered.
func (s *BasePlSqlParserListener) EnterSet_command(ctx *Set_commandContext) {}

// ExitSet_command is called when production set_command is exited.
func (s *BasePlSqlParserListener) ExitSet_command(ctx *Set_commandContext) {}

// EnterTiming_command is called when production timing_command is entered.
func (s *BasePlSqlParserListener) EnterTiming_command(ctx *Timing_commandContext) {}

// ExitTiming_command is called when production timing_command is exited.
func (s *BasePlSqlParserListener) ExitTiming_command(ctx *Timing_commandContext) {}

// EnterClear_command is called when production clear_command is entered.
func (s *BasePlSqlParserListener) EnterClear_command(ctx *Clear_commandContext) {}

// ExitClear_command is called when production clear_command is exited.
func (s *BasePlSqlParserListener) ExitClear_command(ctx *Clear_commandContext) {}

// EnterPartition_extension_clause is called when production partition_extension_clause is entered.
func (s *BasePlSqlParserListener) EnterPartition_extension_clause(ctx *Partition_extension_clauseContext) {
}

// ExitPartition_extension_clause is called when production partition_extension_clause is exited.
func (s *BasePlSqlParserListener) ExitPartition_extension_clause(ctx *Partition_extension_clauseContext) {
}

// EnterColumn_alias is called when production column_alias is entered.
func (s *BasePlSqlParserListener) EnterColumn_alias(ctx *Column_aliasContext) {}

// ExitColumn_alias is called when production column_alias is exited.
func (s *BasePlSqlParserListener) ExitColumn_alias(ctx *Column_aliasContext) {}

// EnterTable_alias is called when production table_alias is entered.
func (s *BasePlSqlParserListener) EnterTable_alias(ctx *Table_aliasContext) {}

// ExitTable_alias is called when production table_alias is exited.
func (s *BasePlSqlParserListener) ExitTable_alias(ctx *Table_aliasContext) {}

// EnterWhere_clause is called when production where_clause is entered.
func (s *BasePlSqlParserListener) EnterWhere_clause(ctx *Where_clauseContext) {}

// ExitWhere_clause is called when production where_clause is exited.
func (s *BasePlSqlParserListener) ExitWhere_clause(ctx *Where_clauseContext) {}

// EnterInto_clause is called when production into_clause is entered.
func (s *BasePlSqlParserListener) EnterInto_clause(ctx *Into_clauseContext) {}

// ExitInto_clause is called when production into_clause is exited.
func (s *BasePlSqlParserListener) ExitInto_clause(ctx *Into_clauseContext) {}

// EnterXml_column_name is called when production xml_column_name is entered.
func (s *BasePlSqlParserListener) EnterXml_column_name(ctx *Xml_column_nameContext) {}

// ExitXml_column_name is called when production xml_column_name is exited.
func (s *BasePlSqlParserListener) ExitXml_column_name(ctx *Xml_column_nameContext) {}

// EnterCost_class_name is called when production cost_class_name is entered.
func (s *BasePlSqlParserListener) EnterCost_class_name(ctx *Cost_class_nameContext) {}

// ExitCost_class_name is called when production cost_class_name is exited.
func (s *BasePlSqlParserListener) ExitCost_class_name(ctx *Cost_class_nameContext) {}

// EnterAttribute_name is called when production attribute_name is entered.
func (s *BasePlSqlParserListener) EnterAttribute_name(ctx *Attribute_nameContext) {}

// ExitAttribute_name is called when production attribute_name is exited.
func (s *BasePlSqlParserListener) ExitAttribute_name(ctx *Attribute_nameContext) {}

// EnterSavepoint_name is called when production savepoint_name is entered.
func (s *BasePlSqlParserListener) EnterSavepoint_name(ctx *Savepoint_nameContext) {}

// ExitSavepoint_name is called when production savepoint_name is exited.
func (s *BasePlSqlParserListener) ExitSavepoint_name(ctx *Savepoint_nameContext) {}

// EnterRollback_segment_name is called when production rollback_segment_name is entered.
func (s *BasePlSqlParserListener) EnterRollback_segment_name(ctx *Rollback_segment_nameContext) {}

// ExitRollback_segment_name is called when production rollback_segment_name is exited.
func (s *BasePlSqlParserListener) ExitRollback_segment_name(ctx *Rollback_segment_nameContext) {}

// EnterSchema_name is called when production schema_name is entered.
func (s *BasePlSqlParserListener) EnterSchema_name(ctx *Schema_nameContext) {}

// ExitSchema_name is called when production schema_name is exited.
func (s *BasePlSqlParserListener) ExitSchema_name(ctx *Schema_nameContext) {}

// EnterRoutine_name is called when production routine_name is entered.
func (s *BasePlSqlParserListener) EnterRoutine_name(ctx *Routine_nameContext) {}

// ExitRoutine_name is called when production routine_name is exited.
func (s *BasePlSqlParserListener) ExitRoutine_name(ctx *Routine_nameContext) {}

// EnterPackage_name is called when production package_name is entered.
func (s *BasePlSqlParserListener) EnterPackage_name(ctx *Package_nameContext) {}

// ExitPackage_name is called when production package_name is exited.
func (s *BasePlSqlParserListener) ExitPackage_name(ctx *Package_nameContext) {}

// EnterImplementation_type_name is called when production implementation_type_name is entered.
func (s *BasePlSqlParserListener) EnterImplementation_type_name(ctx *Implementation_type_nameContext) {
}

// ExitImplementation_type_name is called when production implementation_type_name is exited.
func (s *BasePlSqlParserListener) ExitImplementation_type_name(ctx *Implementation_type_nameContext) {
}

// EnterParameter_name is called when production parameter_name is entered.
func (s *BasePlSqlParserListener) EnterParameter_name(ctx *Parameter_nameContext) {}

// ExitParameter_name is called when production parameter_name is exited.
func (s *BasePlSqlParserListener) ExitParameter_name(ctx *Parameter_nameContext) {}

// EnterReference_model_name is called when production reference_model_name is entered.
func (s *BasePlSqlParserListener) EnterReference_model_name(ctx *Reference_model_nameContext) {}

// ExitReference_model_name is called when production reference_model_name is exited.
func (s *BasePlSqlParserListener) ExitReference_model_name(ctx *Reference_model_nameContext) {}

// EnterMain_model_name is called when production main_model_name is entered.
func (s *BasePlSqlParserListener) EnterMain_model_name(ctx *Main_model_nameContext) {}

// ExitMain_model_name is called when production main_model_name is exited.
func (s *BasePlSqlParserListener) ExitMain_model_name(ctx *Main_model_nameContext) {}

// EnterContainer_tableview_name is called when production container_tableview_name is entered.
func (s *BasePlSqlParserListener) EnterContainer_tableview_name(ctx *Container_tableview_nameContext) {
}

// ExitContainer_tableview_name is called when production container_tableview_name is exited.
func (s *BasePlSqlParserListener) ExitContainer_tableview_name(ctx *Container_tableview_nameContext) {
}

// EnterAggregate_function_name is called when production aggregate_function_name is entered.
func (s *BasePlSqlParserListener) EnterAggregate_function_name(ctx *Aggregate_function_nameContext) {}

// ExitAggregate_function_name is called when production aggregate_function_name is exited.
func (s *BasePlSqlParserListener) ExitAggregate_function_name(ctx *Aggregate_function_nameContext) {}

// EnterQuery_name is called when production query_name is entered.
func (s *BasePlSqlParserListener) EnterQuery_name(ctx *Query_nameContext) {}

// ExitQuery_name is called when production query_name is exited.
func (s *BasePlSqlParserListener) ExitQuery_name(ctx *Query_nameContext) {}

// EnterGrantee_name is called when production grantee_name is entered.
func (s *BasePlSqlParserListener) EnterGrantee_name(ctx *Grantee_nameContext) {}

// ExitGrantee_name is called when production grantee_name is exited.
func (s *BasePlSqlParserListener) ExitGrantee_name(ctx *Grantee_nameContext) {}

// EnterRole_name is called when production role_name is entered.
func (s *BasePlSqlParserListener) EnterRole_name(ctx *Role_nameContext) {}

// ExitRole_name is called when production role_name is exited.
func (s *BasePlSqlParserListener) ExitRole_name(ctx *Role_nameContext) {}

// EnterConstraint_name is called when production constraint_name is entered.
func (s *BasePlSqlParserListener) EnterConstraint_name(ctx *Constraint_nameContext) {}

// ExitConstraint_name is called when production constraint_name is exited.
func (s *BasePlSqlParserListener) ExitConstraint_name(ctx *Constraint_nameContext) {}

// EnterLabel_name is called when production label_name is entered.
func (s *BasePlSqlParserListener) EnterLabel_name(ctx *Label_nameContext) {}

// ExitLabel_name is called when production label_name is exited.
func (s *BasePlSqlParserListener) ExitLabel_name(ctx *Label_nameContext) {}

// EnterType_name is called when production type_name is entered.
func (s *BasePlSqlParserListener) EnterType_name(ctx *Type_nameContext) {}

// ExitType_name is called when production type_name is exited.
func (s *BasePlSqlParserListener) ExitType_name(ctx *Type_nameContext) {}

// EnterSequence_name is called when production sequence_name is entered.
func (s *BasePlSqlParserListener) EnterSequence_name(ctx *Sequence_nameContext) {}

// ExitSequence_name is called when production sequence_name is exited.
func (s *BasePlSqlParserListener) ExitSequence_name(ctx *Sequence_nameContext) {}

// EnterException_name is called when production exception_name is entered.
func (s *BasePlSqlParserListener) EnterException_name(ctx *Exception_nameContext) {}

// ExitException_name is called when production exception_name is exited.
func (s *BasePlSqlParserListener) ExitException_name(ctx *Exception_nameContext) {}

// EnterFunction_name is called when production function_name is entered.
func (s *BasePlSqlParserListener) EnterFunction_name(ctx *Function_nameContext) {}

// ExitFunction_name is called when production function_name is exited.
func (s *BasePlSqlParserListener) ExitFunction_name(ctx *Function_nameContext) {}

// EnterProcedure_name is called when production procedure_name is entered.
func (s *BasePlSqlParserListener) EnterProcedure_name(ctx *Procedure_nameContext) {}

// ExitProcedure_name is called when production procedure_name is exited.
func (s *BasePlSqlParserListener) ExitProcedure_name(ctx *Procedure_nameContext) {}

// EnterTrigger_name is called when production trigger_name is entered.
func (s *BasePlSqlParserListener) EnterTrigger_name(ctx *Trigger_nameContext) {}

// ExitTrigger_name is called when production trigger_name is exited.
func (s *BasePlSqlParserListener) ExitTrigger_name(ctx *Trigger_nameContext) {}

// EnterVariable_name is called when production variable_name is entered.
func (s *BasePlSqlParserListener) EnterVariable_name(ctx *Variable_nameContext) {}

// ExitVariable_name is called when production variable_name is exited.
func (s *BasePlSqlParserListener) ExitVariable_name(ctx *Variable_nameContext) {}

// EnterIndex_name is called when production index_name is entered.
func (s *BasePlSqlParserListener) EnterIndex_name(ctx *Index_nameContext) {}

// ExitIndex_name is called when production index_name is exited.
func (s *BasePlSqlParserListener) ExitIndex_name(ctx *Index_nameContext) {}

// EnterCursor_name is called when production cursor_name is entered.
func (s *BasePlSqlParserListener) EnterCursor_name(ctx *Cursor_nameContext) {}

// ExitCursor_name is called when production cursor_name is exited.
func (s *BasePlSqlParserListener) ExitCursor_name(ctx *Cursor_nameContext) {}

// EnterRecord_name is called when production record_name is entered.
func (s *BasePlSqlParserListener) EnterRecord_name(ctx *Record_nameContext) {}

// ExitRecord_name is called when production record_name is exited.
func (s *BasePlSqlParserListener) ExitRecord_name(ctx *Record_nameContext) {}

// EnterLink_name is called when production link_name is entered.
func (s *BasePlSqlParserListener) EnterLink_name(ctx *Link_nameContext) {}

// ExitLink_name is called when production link_name is exited.
func (s *BasePlSqlParserListener) ExitLink_name(ctx *Link_nameContext) {}

// EnterLocal_link_name is called when production local_link_name is entered.
func (s *BasePlSqlParserListener) EnterLocal_link_name(ctx *Local_link_nameContext) {}

// ExitLocal_link_name is called when production local_link_name is exited.
func (s *BasePlSqlParserListener) ExitLocal_link_name(ctx *Local_link_nameContext) {}

// EnterConnection_qualifier is called when production connection_qualifier is entered.
func (s *BasePlSqlParserListener) EnterConnection_qualifier(ctx *Connection_qualifierContext) {}

// ExitConnection_qualifier is called when production connection_qualifier is exited.
func (s *BasePlSqlParserListener) ExitConnection_qualifier(ctx *Connection_qualifierContext) {}

// EnterColumn_name is called when production column_name is entered.
func (s *BasePlSqlParserListener) EnterColumn_name(ctx *Column_nameContext) {}

// ExitColumn_name is called when production column_name is exited.
func (s *BasePlSqlParserListener) ExitColumn_name(ctx *Column_nameContext) {}

// EnterTableview_name is called when production tableview_name is entered.
func (s *BasePlSqlParserListener) EnterTableview_name(ctx *Tableview_nameContext) {}

// ExitTableview_name is called when production tableview_name is exited.
func (s *BasePlSqlParserListener) ExitTableview_name(ctx *Tableview_nameContext) {}

// EnterXmltable is called when production xmltable is entered.
func (s *BasePlSqlParserListener) EnterXmltable(ctx *XmltableContext) {}

// ExitXmltable is called when production xmltable is exited.
func (s *BasePlSqlParserListener) ExitXmltable(ctx *XmltableContext) {}

// EnterChar_set_name is called when production char_set_name is entered.
func (s *BasePlSqlParserListener) EnterChar_set_name(ctx *Char_set_nameContext) {}

// ExitChar_set_name is called when production char_set_name is exited.
func (s *BasePlSqlParserListener) ExitChar_set_name(ctx *Char_set_nameContext) {}

// EnterSynonym_name is called when production synonym_name is entered.
func (s *BasePlSqlParserListener) EnterSynonym_name(ctx *Synonym_nameContext) {}

// ExitSynonym_name is called when production synonym_name is exited.
func (s *BasePlSqlParserListener) ExitSynonym_name(ctx *Synonym_nameContext) {}

// EnterSchema_object_name is called when production schema_object_name is entered.
func (s *BasePlSqlParserListener) EnterSchema_object_name(ctx *Schema_object_nameContext) {}

// ExitSchema_object_name is called when production schema_object_name is exited.
func (s *BasePlSqlParserListener) ExitSchema_object_name(ctx *Schema_object_nameContext) {}

// EnterDir_object_name is called when production dir_object_name is entered.
func (s *BasePlSqlParserListener) EnterDir_object_name(ctx *Dir_object_nameContext) {}

// ExitDir_object_name is called when production dir_object_name is exited.
func (s *BasePlSqlParserListener) ExitDir_object_name(ctx *Dir_object_nameContext) {}

// EnterUser_object_name is called when production user_object_name is entered.
func (s *BasePlSqlParserListener) EnterUser_object_name(ctx *User_object_nameContext) {}

// ExitUser_object_name is called when production user_object_name is exited.
func (s *BasePlSqlParserListener) ExitUser_object_name(ctx *User_object_nameContext) {}

// EnterGrant_object_name is called when production grant_object_name is entered.
func (s *BasePlSqlParserListener) EnterGrant_object_name(ctx *Grant_object_nameContext) {}

// ExitGrant_object_name is called when production grant_object_name is exited.
func (s *BasePlSqlParserListener) ExitGrant_object_name(ctx *Grant_object_nameContext) {}

// EnterColumn_list is called when production column_list is entered.
func (s *BasePlSqlParserListener) EnterColumn_list(ctx *Column_listContext) {}

// ExitColumn_list is called when production column_list is exited.
func (s *BasePlSqlParserListener) ExitColumn_list(ctx *Column_listContext) {}

// EnterParen_column_list is called when production paren_column_list is entered.
func (s *BasePlSqlParserListener) EnterParen_column_list(ctx *Paren_column_listContext) {}

// ExitParen_column_list is called when production paren_column_list is exited.
func (s *BasePlSqlParserListener) ExitParen_column_list(ctx *Paren_column_listContext) {}

// EnterKeep_clause is called when production keep_clause is entered.
func (s *BasePlSqlParserListener) EnterKeep_clause(ctx *Keep_clauseContext) {}

// ExitKeep_clause is called when production keep_clause is exited.
func (s *BasePlSqlParserListener) ExitKeep_clause(ctx *Keep_clauseContext) {}

// EnterFunction_argument is called when production function_argument is entered.
func (s *BasePlSqlParserListener) EnterFunction_argument(ctx *Function_argumentContext) {}

// ExitFunction_argument is called when production function_argument is exited.
func (s *BasePlSqlParserListener) ExitFunction_argument(ctx *Function_argumentContext) {}

// EnterFunction_argument_analytic is called when production function_argument_analytic is entered.
func (s *BasePlSqlParserListener) EnterFunction_argument_analytic(ctx *Function_argument_analyticContext) {
}

// ExitFunction_argument_analytic is called when production function_argument_analytic is exited.
func (s *BasePlSqlParserListener) ExitFunction_argument_analytic(ctx *Function_argument_analyticContext) {
}

// EnterFunction_argument_modeling is called when production function_argument_modeling is entered.
func (s *BasePlSqlParserListener) EnterFunction_argument_modeling(ctx *Function_argument_modelingContext) {
}

// ExitFunction_argument_modeling is called when production function_argument_modeling is exited.
func (s *BasePlSqlParserListener) ExitFunction_argument_modeling(ctx *Function_argument_modelingContext) {
}

// EnterRespect_or_ignore_nulls is called when production respect_or_ignore_nulls is entered.
func (s *BasePlSqlParserListener) EnterRespect_or_ignore_nulls(ctx *Respect_or_ignore_nullsContext) {}

// ExitRespect_or_ignore_nulls is called when production respect_or_ignore_nulls is exited.
func (s *BasePlSqlParserListener) ExitRespect_or_ignore_nulls(ctx *Respect_or_ignore_nullsContext) {}

// EnterArgument is called when production argument is entered.
func (s *BasePlSqlParserListener) EnterArgument(ctx *ArgumentContext) {}

// ExitArgument is called when production argument is exited.
func (s *BasePlSqlParserListener) ExitArgument(ctx *ArgumentContext) {}

// EnterType_spec is called when production type_spec is entered.
func (s *BasePlSqlParserListener) EnterType_spec(ctx *Type_specContext) {}

// ExitType_spec is called when production type_spec is exited.
func (s *BasePlSqlParserListener) ExitType_spec(ctx *Type_specContext) {}

// EnterDatatype is called when production datatype is entered.
func (s *BasePlSqlParserListener) EnterDatatype(ctx *DatatypeContext) {}

// ExitDatatype is called when production datatype is exited.
func (s *BasePlSqlParserListener) ExitDatatype(ctx *DatatypeContext) {}

// EnterPrecision_part is called when production precision_part is entered.
func (s *BasePlSqlParserListener) EnterPrecision_part(ctx *Precision_partContext) {}

// ExitPrecision_part is called when production precision_part is exited.
func (s *BasePlSqlParserListener) ExitPrecision_part(ctx *Precision_partContext) {}

// EnterNative_datatype_element is called when production native_datatype_element is entered.
func (s *BasePlSqlParserListener) EnterNative_datatype_element(ctx *Native_datatype_elementContext) {}

// ExitNative_datatype_element is called when production native_datatype_element is exited.
func (s *BasePlSqlParserListener) ExitNative_datatype_element(ctx *Native_datatype_elementContext) {}

// EnterBind_variable is called when production bind_variable is entered.
func (s *BasePlSqlParserListener) EnterBind_variable(ctx *Bind_variableContext) {}

// ExitBind_variable is called when production bind_variable is exited.
func (s *BasePlSqlParserListener) ExitBind_variable(ctx *Bind_variableContext) {}

// EnterGeneral_element is called when production general_element is entered.
func (s *BasePlSqlParserListener) EnterGeneral_element(ctx *General_elementContext) {}

// ExitGeneral_element is called when production general_element is exited.
func (s *BasePlSqlParserListener) ExitGeneral_element(ctx *General_elementContext) {}

// EnterGeneral_element_part is called when production general_element_part is entered.
func (s *BasePlSqlParserListener) EnterGeneral_element_part(ctx *General_element_partContext) {}

// ExitGeneral_element_part is called when production general_element_part is exited.
func (s *BasePlSqlParserListener) ExitGeneral_element_part(ctx *General_element_partContext) {}

// EnterTable_element is called when production table_element is entered.
func (s *BasePlSqlParserListener) EnterTable_element(ctx *Table_elementContext) {}

// ExitTable_element is called when production table_element is exited.
func (s *BasePlSqlParserListener) ExitTable_element(ctx *Table_elementContext) {}

// EnterObject_privilege is called when production object_privilege is entered.
func (s *BasePlSqlParserListener) EnterObject_privilege(ctx *Object_privilegeContext) {}

// ExitObject_privilege is called when production object_privilege is exited.
func (s *BasePlSqlParserListener) ExitObject_privilege(ctx *Object_privilegeContext) {}

// EnterSystem_privilege is called when production system_privilege is entered.
func (s *BasePlSqlParserListener) EnterSystem_privilege(ctx *System_privilegeContext) {}

// ExitSystem_privilege is called when production system_privilege is exited.
func (s *BasePlSqlParserListener) ExitSystem_privilege(ctx *System_privilegeContext) {}

// EnterConstant is called when production constant is entered.
func (s *BasePlSqlParserListener) EnterConstant(ctx *ConstantContext) {}

// ExitConstant is called when production constant is exited.
func (s *BasePlSqlParserListener) ExitConstant(ctx *ConstantContext) {}

// EnterNumeric is called when production numeric is entered.
func (s *BasePlSqlParserListener) EnterNumeric(ctx *NumericContext) {}

// ExitNumeric is called when production numeric is exited.
func (s *BasePlSqlParserListener) ExitNumeric(ctx *NumericContext) {}

// EnterNumeric_negative is called when production numeric_negative is entered.
func (s *BasePlSqlParserListener) EnterNumeric_negative(ctx *Numeric_negativeContext) {}

// ExitNumeric_negative is called when production numeric_negative is exited.
func (s *BasePlSqlParserListener) ExitNumeric_negative(ctx *Numeric_negativeContext) {}

// EnterQuoted_string is called when production quoted_string is entered.
func (s *BasePlSqlParserListener) EnterQuoted_string(ctx *Quoted_stringContext) {}

// ExitQuoted_string is called when production quoted_string is exited.
func (s *BasePlSqlParserListener) ExitQuoted_string(ctx *Quoted_stringContext) {}

// EnterIdentifier is called when production identifier is entered.
func (s *BasePlSqlParserListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production identifier is exited.
func (s *BasePlSqlParserListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterId_expression is called when production id_expression is entered.
func (s *BasePlSqlParserListener) EnterId_expression(ctx *Id_expressionContext) {}

// ExitId_expression is called when production id_expression is exited.
func (s *BasePlSqlParserListener) ExitId_expression(ctx *Id_expressionContext) {}

// EnterInquiry_directive is called when production inquiry_directive is entered.
func (s *BasePlSqlParserListener) EnterInquiry_directive(ctx *Inquiry_directiveContext) {}

// ExitInquiry_directive is called when production inquiry_directive is exited.
func (s *BasePlSqlParserListener) ExitInquiry_directive(ctx *Inquiry_directiveContext) {}

// EnterOuter_join_sign is called when production outer_join_sign is entered.
func (s *BasePlSqlParserListener) EnterOuter_join_sign(ctx *Outer_join_signContext) {}

// ExitOuter_join_sign is called when production outer_join_sign is exited.
func (s *BasePlSqlParserListener) ExitOuter_join_sign(ctx *Outer_join_signContext) {}

// EnterRegular_id is called when production regular_id is entered.
func (s *BasePlSqlParserListener) EnterRegular_id(ctx *Regular_idContext) {}

// ExitRegular_id is called when production regular_id is exited.
func (s *BasePlSqlParserListener) ExitRegular_id(ctx *Regular_idContext) {}

// EnterNon_reserved_keywords_in_18c is called when production non_reserved_keywords_in_18c is entered.
func (s *BasePlSqlParserListener) EnterNon_reserved_keywords_in_18c(ctx *Non_reserved_keywords_in_18cContext) {
}

// ExitNon_reserved_keywords_in_18c is called when production non_reserved_keywords_in_18c is exited.
func (s *BasePlSqlParserListener) ExitNon_reserved_keywords_in_18c(ctx *Non_reserved_keywords_in_18cContext) {
}

// EnterNon_reserved_keywords_in_12c is called when production non_reserved_keywords_in_12c is entered.
func (s *BasePlSqlParserListener) EnterNon_reserved_keywords_in_12c(ctx *Non_reserved_keywords_in_12cContext) {
}

// ExitNon_reserved_keywords_in_12c is called when production non_reserved_keywords_in_12c is exited.
func (s *BasePlSqlParserListener) ExitNon_reserved_keywords_in_12c(ctx *Non_reserved_keywords_in_12cContext) {
}

// EnterNon_reserved_keywords_pre12c is called when production non_reserved_keywords_pre12c is entered.
func (s *BasePlSqlParserListener) EnterNon_reserved_keywords_pre12c(ctx *Non_reserved_keywords_pre12cContext) {
}

// ExitNon_reserved_keywords_pre12c is called when production non_reserved_keywords_pre12c is exited.
func (s *BasePlSqlParserListener) ExitNon_reserved_keywords_pre12c(ctx *Non_reserved_keywords_pre12cContext) {
}
