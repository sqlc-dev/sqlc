// Code generated from PlSqlParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // PlSqlParser
import "github.com/antlr4-go/antlr/v4"

// PlSqlParserListener is a complete listener for a parse tree produced by PlSqlParser.
type PlSqlParserListener interface {
	antlr.ParseTreeListener

	// EnterSql_script is called when entering the sql_script production.
	EnterSql_script(c *Sql_scriptContext)

	// EnterUnit_statement is called when entering the unit_statement production.
	EnterUnit_statement(c *Unit_statementContext)

	// EnterAlter_diskgroup is called when entering the alter_diskgroup production.
	EnterAlter_diskgroup(c *Alter_diskgroupContext)

	// EnterAdd_disk_clause is called when entering the add_disk_clause production.
	EnterAdd_disk_clause(c *Add_disk_clauseContext)

	// EnterDrop_disk_clause is called when entering the drop_disk_clause production.
	EnterDrop_disk_clause(c *Drop_disk_clauseContext)

	// EnterResize_disk_clause is called when entering the resize_disk_clause production.
	EnterResize_disk_clause(c *Resize_disk_clauseContext)

	// EnterReplace_disk_clause is called when entering the replace_disk_clause production.
	EnterReplace_disk_clause(c *Replace_disk_clauseContext)

	// EnterWait_nowait is called when entering the wait_nowait production.
	EnterWait_nowait(c *Wait_nowaitContext)

	// EnterRename_disk_clause is called when entering the rename_disk_clause production.
	EnterRename_disk_clause(c *Rename_disk_clauseContext)

	// EnterDisk_online_clause is called when entering the disk_online_clause production.
	EnterDisk_online_clause(c *Disk_online_clauseContext)

	// EnterDisk_offline_clause is called when entering the disk_offline_clause production.
	EnterDisk_offline_clause(c *Disk_offline_clauseContext)

	// EnterTimeout_clause is called when entering the timeout_clause production.
	EnterTimeout_clause(c *Timeout_clauseContext)

	// EnterRebalance_diskgroup_clause is called when entering the rebalance_diskgroup_clause production.
	EnterRebalance_diskgroup_clause(c *Rebalance_diskgroup_clauseContext)

	// EnterPhase is called when entering the phase production.
	EnterPhase(c *PhaseContext)

	// EnterCheck_diskgroup_clause is called when entering the check_diskgroup_clause production.
	EnterCheck_diskgroup_clause(c *Check_diskgroup_clauseContext)

	// EnterDiskgroup_template_clauses is called when entering the diskgroup_template_clauses production.
	EnterDiskgroup_template_clauses(c *Diskgroup_template_clausesContext)

	// EnterQualified_template_clause is called when entering the qualified_template_clause production.
	EnterQualified_template_clause(c *Qualified_template_clauseContext)

	// EnterRedundancy_clause is called when entering the redundancy_clause production.
	EnterRedundancy_clause(c *Redundancy_clauseContext)

	// EnterStriping_clause is called when entering the striping_clause production.
	EnterStriping_clause(c *Striping_clauseContext)

	// EnterForce_noforce is called when entering the force_noforce production.
	EnterForce_noforce(c *Force_noforceContext)

	// EnterDiskgroup_directory_clauses is called when entering the diskgroup_directory_clauses production.
	EnterDiskgroup_directory_clauses(c *Diskgroup_directory_clausesContext)

	// EnterDir_name is called when entering the dir_name production.
	EnterDir_name(c *Dir_nameContext)

	// EnterDiskgroup_alias_clauses is called when entering the diskgroup_alias_clauses production.
	EnterDiskgroup_alias_clauses(c *Diskgroup_alias_clausesContext)

	// EnterDiskgroup_volume_clauses is called when entering the diskgroup_volume_clauses production.
	EnterDiskgroup_volume_clauses(c *Diskgroup_volume_clausesContext)

	// EnterAdd_volume_clause is called when entering the add_volume_clause production.
	EnterAdd_volume_clause(c *Add_volume_clauseContext)

	// EnterModify_volume_clause is called when entering the modify_volume_clause production.
	EnterModify_volume_clause(c *Modify_volume_clauseContext)

	// EnterDiskgroup_attributes is called when entering the diskgroup_attributes production.
	EnterDiskgroup_attributes(c *Diskgroup_attributesContext)

	// EnterDrop_diskgroup_file_clause is called when entering the drop_diskgroup_file_clause production.
	EnterDrop_diskgroup_file_clause(c *Drop_diskgroup_file_clauseContext)

	// EnterConvert_redundancy_clause is called when entering the convert_redundancy_clause production.
	EnterConvert_redundancy_clause(c *Convert_redundancy_clauseContext)

	// EnterUsergroup_clauses is called when entering the usergroup_clauses production.
	EnterUsergroup_clauses(c *Usergroup_clausesContext)

	// EnterUser_clauses is called when entering the user_clauses production.
	EnterUser_clauses(c *User_clausesContext)

	// EnterFile_permissions_clause is called when entering the file_permissions_clause production.
	EnterFile_permissions_clause(c *File_permissions_clauseContext)

	// EnterFile_owner_clause is called when entering the file_owner_clause production.
	EnterFile_owner_clause(c *File_owner_clauseContext)

	// EnterScrub_clause is called when entering the scrub_clause production.
	EnterScrub_clause(c *Scrub_clauseContext)

	// EnterQuotagroup_clauses is called when entering the quotagroup_clauses production.
	EnterQuotagroup_clauses(c *Quotagroup_clausesContext)

	// EnterProperty_name is called when entering the property_name production.
	EnterProperty_name(c *Property_nameContext)

	// EnterProperty_value is called when entering the property_value production.
	EnterProperty_value(c *Property_valueContext)

	// EnterFilegroup_clauses is called when entering the filegroup_clauses production.
	EnterFilegroup_clauses(c *Filegroup_clausesContext)

	// EnterAdd_filegroup_clause is called when entering the add_filegroup_clause production.
	EnterAdd_filegroup_clause(c *Add_filegroup_clauseContext)

	// EnterModify_filegroup_clause is called when entering the modify_filegroup_clause production.
	EnterModify_filegroup_clause(c *Modify_filegroup_clauseContext)

	// EnterMove_to_filegroup_clause is called when entering the move_to_filegroup_clause production.
	EnterMove_to_filegroup_clause(c *Move_to_filegroup_clauseContext)

	// EnterDrop_filegroup_clause is called when entering the drop_filegroup_clause production.
	EnterDrop_filegroup_clause(c *Drop_filegroup_clauseContext)

	// EnterQuorum_regular is called when entering the quorum_regular production.
	EnterQuorum_regular(c *Quorum_regularContext)

	// EnterUndrop_disk_clause is called when entering the undrop_disk_clause production.
	EnterUndrop_disk_clause(c *Undrop_disk_clauseContext)

	// EnterDiskgroup_availability is called when entering the diskgroup_availability production.
	EnterDiskgroup_availability(c *Diskgroup_availabilityContext)

	// EnterEnable_disable_volume is called when entering the enable_disable_volume production.
	EnterEnable_disable_volume(c *Enable_disable_volumeContext)

	// EnterDrop_function is called when entering the drop_function production.
	EnterDrop_function(c *Drop_functionContext)

	// EnterAlter_flashback_archive is called when entering the alter_flashback_archive production.
	EnterAlter_flashback_archive(c *Alter_flashback_archiveContext)

	// EnterAlter_hierarchy is called when entering the alter_hierarchy production.
	EnterAlter_hierarchy(c *Alter_hierarchyContext)

	// EnterAlter_function is called when entering the alter_function production.
	EnterAlter_function(c *Alter_functionContext)

	// EnterAlter_java is called when entering the alter_java production.
	EnterAlter_java(c *Alter_javaContext)

	// EnterMatch_string is called when entering the match_string production.
	EnterMatch_string(c *Match_stringContext)

	// EnterCreate_function_body is called when entering the create_function_body production.
	EnterCreate_function_body(c *Create_function_bodyContext)

	// EnterSql_macro_body is called when entering the sql_macro_body production.
	EnterSql_macro_body(c *Sql_macro_bodyContext)

	// EnterParallel_enable_clause is called when entering the parallel_enable_clause production.
	EnterParallel_enable_clause(c *Parallel_enable_clauseContext)

	// EnterPartition_by_clause is called when entering the partition_by_clause production.
	EnterPartition_by_clause(c *Partition_by_clauseContext)

	// EnterResult_cache_clause is called when entering the result_cache_clause production.
	EnterResult_cache_clause(c *Result_cache_clauseContext)

	// EnterAccessible_by_clause is called when entering the accessible_by_clause production.
	EnterAccessible_by_clause(c *Accessible_by_clauseContext)

	// EnterDefault_collation_clause is called when entering the default_collation_clause production.
	EnterDefault_collation_clause(c *Default_collation_clauseContext)

	// EnterAggregate_clause is called when entering the aggregate_clause production.
	EnterAggregate_clause(c *Aggregate_clauseContext)

	// EnterPipelined_using_clause is called when entering the pipelined_using_clause production.
	EnterPipelined_using_clause(c *Pipelined_using_clauseContext)

	// EnterAccessor is called when entering the accessor production.
	EnterAccessor(c *AccessorContext)

	// EnterRelies_on_part is called when entering the relies_on_part production.
	EnterRelies_on_part(c *Relies_on_partContext)

	// EnterStreaming_clause is called when entering the streaming_clause production.
	EnterStreaming_clause(c *Streaming_clauseContext)

	// EnterAlter_outline is called when entering the alter_outline production.
	EnterAlter_outline(c *Alter_outlineContext)

	// EnterOutline_options is called when entering the outline_options production.
	EnterOutline_options(c *Outline_optionsContext)

	// EnterAlter_lockdown_profile is called when entering the alter_lockdown_profile production.
	EnterAlter_lockdown_profile(c *Alter_lockdown_profileContext)

	// EnterLockdown_feature is called when entering the lockdown_feature production.
	EnterLockdown_feature(c *Lockdown_featureContext)

	// EnterLockdown_options is called when entering the lockdown_options production.
	EnterLockdown_options(c *Lockdown_optionsContext)

	// EnterLockdown_statements is called when entering the lockdown_statements production.
	EnterLockdown_statements(c *Lockdown_statementsContext)

	// EnterStatement_clauses is called when entering the statement_clauses production.
	EnterStatement_clauses(c *Statement_clausesContext)

	// EnterClause_options is called when entering the clause_options production.
	EnterClause_options(c *Clause_optionsContext)

	// EnterOption_values is called when entering the option_values production.
	EnterOption_values(c *Option_valuesContext)

	// EnterString_list is called when entering the string_list production.
	EnterString_list(c *String_listContext)

	// EnterDisable_enable is called when entering the disable_enable production.
	EnterDisable_enable(c *Disable_enableContext)

	// EnterDrop_lockdown_profile is called when entering the drop_lockdown_profile production.
	EnterDrop_lockdown_profile(c *Drop_lockdown_profileContext)

	// EnterDrop_package is called when entering the drop_package production.
	EnterDrop_package(c *Drop_packageContext)

	// EnterAlter_package is called when entering the alter_package production.
	EnterAlter_package(c *Alter_packageContext)

	// EnterCreate_package is called when entering the create_package production.
	EnterCreate_package(c *Create_packageContext)

	// EnterCreate_package_body is called when entering the create_package_body production.
	EnterCreate_package_body(c *Create_package_bodyContext)

	// EnterPackage_obj_spec is called when entering the package_obj_spec production.
	EnterPackage_obj_spec(c *Package_obj_specContext)

	// EnterProcedure_spec is called when entering the procedure_spec production.
	EnterProcedure_spec(c *Procedure_specContext)

	// EnterFunction_spec is called when entering the function_spec production.
	EnterFunction_spec(c *Function_specContext)

	// EnterPackage_obj_body is called when entering the package_obj_body production.
	EnterPackage_obj_body(c *Package_obj_bodyContext)

	// EnterAlter_pmem_filestore is called when entering the alter_pmem_filestore production.
	EnterAlter_pmem_filestore(c *Alter_pmem_filestoreContext)

	// EnterDrop_pmem_filestore is called when entering the drop_pmem_filestore production.
	EnterDrop_pmem_filestore(c *Drop_pmem_filestoreContext)

	// EnterDrop_procedure is called when entering the drop_procedure production.
	EnterDrop_procedure(c *Drop_procedureContext)

	// EnterAlter_procedure is called when entering the alter_procedure production.
	EnterAlter_procedure(c *Alter_procedureContext)

	// EnterFunction_body is called when entering the function_body production.
	EnterFunction_body(c *Function_bodyContext)

	// EnterProcedure_body is called when entering the procedure_body production.
	EnterProcedure_body(c *Procedure_bodyContext)

	// EnterCreate_procedure_body is called when entering the create_procedure_body production.
	EnterCreate_procedure_body(c *Create_procedure_bodyContext)

	// EnterAlter_resource_cost is called when entering the alter_resource_cost production.
	EnterAlter_resource_cost(c *Alter_resource_costContext)

	// EnterDrop_outline is called when entering the drop_outline production.
	EnterDrop_outline(c *Drop_outlineContext)

	// EnterAlter_rollback_segment is called when entering the alter_rollback_segment production.
	EnterAlter_rollback_segment(c *Alter_rollback_segmentContext)

	// EnterDrop_restore_point is called when entering the drop_restore_point production.
	EnterDrop_restore_point(c *Drop_restore_pointContext)

	// EnterDrop_rollback_segment is called when entering the drop_rollback_segment production.
	EnterDrop_rollback_segment(c *Drop_rollback_segmentContext)

	// EnterDrop_role is called when entering the drop_role production.
	EnterDrop_role(c *Drop_roleContext)

	// EnterCreate_pmem_filestore is called when entering the create_pmem_filestore production.
	EnterCreate_pmem_filestore(c *Create_pmem_filestoreContext)

	// EnterPmem_filestore_options is called when entering the pmem_filestore_options production.
	EnterPmem_filestore_options(c *Pmem_filestore_optionsContext)

	// EnterFile_path is called when entering the file_path production.
	EnterFile_path(c *File_pathContext)

	// EnterCreate_rollback_segment is called when entering the create_rollback_segment production.
	EnterCreate_rollback_segment(c *Create_rollback_segmentContext)

	// EnterDrop_trigger is called when entering the drop_trigger production.
	EnterDrop_trigger(c *Drop_triggerContext)

	// EnterAlter_trigger is called when entering the alter_trigger production.
	EnterAlter_trigger(c *Alter_triggerContext)

	// EnterCreate_trigger is called when entering the create_trigger production.
	EnterCreate_trigger(c *Create_triggerContext)

	// EnterTrigger_follows_clause is called when entering the trigger_follows_clause production.
	EnterTrigger_follows_clause(c *Trigger_follows_clauseContext)

	// EnterTrigger_when_clause is called when entering the trigger_when_clause production.
	EnterTrigger_when_clause(c *Trigger_when_clauseContext)

	// EnterSimple_dml_trigger is called when entering the simple_dml_trigger production.
	EnterSimple_dml_trigger(c *Simple_dml_triggerContext)

	// EnterFor_each_row is called when entering the for_each_row production.
	EnterFor_each_row(c *For_each_rowContext)

	// EnterCompound_dml_trigger is called when entering the compound_dml_trigger production.
	EnterCompound_dml_trigger(c *Compound_dml_triggerContext)

	// EnterNon_dml_trigger is called when entering the non_dml_trigger production.
	EnterNon_dml_trigger(c *Non_dml_triggerContext)

	// EnterTrigger_body is called when entering the trigger_body production.
	EnterTrigger_body(c *Trigger_bodyContext)

	// EnterCompound_trigger_block is called when entering the compound_trigger_block production.
	EnterCompound_trigger_block(c *Compound_trigger_blockContext)

	// EnterTiming_point_section is called when entering the timing_point_section production.
	EnterTiming_point_section(c *Timing_point_sectionContext)

	// EnterNon_dml_event is called when entering the non_dml_event production.
	EnterNon_dml_event(c *Non_dml_eventContext)

	// EnterDml_event_clause is called when entering the dml_event_clause production.
	EnterDml_event_clause(c *Dml_event_clauseContext)

	// EnterDml_event_element is called when entering the dml_event_element production.
	EnterDml_event_element(c *Dml_event_elementContext)

	// EnterDml_event_nested_clause is called when entering the dml_event_nested_clause production.
	EnterDml_event_nested_clause(c *Dml_event_nested_clauseContext)

	// EnterReferencing_clause is called when entering the referencing_clause production.
	EnterReferencing_clause(c *Referencing_clauseContext)

	// EnterReferencing_element is called when entering the referencing_element production.
	EnterReferencing_element(c *Referencing_elementContext)

	// EnterDrop_type is called when entering the drop_type production.
	EnterDrop_type(c *Drop_typeContext)

	// EnterAlter_type is called when entering the alter_type production.
	EnterAlter_type(c *Alter_typeContext)

	// EnterCompile_type_clause is called when entering the compile_type_clause production.
	EnterCompile_type_clause(c *Compile_type_clauseContext)

	// EnterReplace_type_clause is called when entering the replace_type_clause production.
	EnterReplace_type_clause(c *Replace_type_clauseContext)

	// EnterAlter_method_spec is called when entering the alter_method_spec production.
	EnterAlter_method_spec(c *Alter_method_specContext)

	// EnterAlter_method_element is called when entering the alter_method_element production.
	EnterAlter_method_element(c *Alter_method_elementContext)

	// EnterAlter_collection_clauses is called when entering the alter_collection_clauses production.
	EnterAlter_collection_clauses(c *Alter_collection_clausesContext)

	// EnterDependent_handling_clause is called when entering the dependent_handling_clause production.
	EnterDependent_handling_clause(c *Dependent_handling_clauseContext)

	// EnterDependent_exceptions_part is called when entering the dependent_exceptions_part production.
	EnterDependent_exceptions_part(c *Dependent_exceptions_partContext)

	// EnterCreate_type is called when entering the create_type production.
	EnterCreate_type(c *Create_typeContext)

	// EnterType_definition is called when entering the type_definition production.
	EnterType_definition(c *Type_definitionContext)

	// EnterObject_type_def is called when entering the object_type_def production.
	EnterObject_type_def(c *Object_type_defContext)

	// EnterObject_as_part is called when entering the object_as_part production.
	EnterObject_as_part(c *Object_as_partContext)

	// EnterObject_under_part is called when entering the object_under_part production.
	EnterObject_under_part(c *Object_under_partContext)

	// EnterNested_table_type_def is called when entering the nested_table_type_def production.
	EnterNested_table_type_def(c *Nested_table_type_defContext)

	// EnterSqlj_object_type is called when entering the sqlj_object_type production.
	EnterSqlj_object_type(c *Sqlj_object_typeContext)

	// EnterType_body is called when entering the type_body production.
	EnterType_body(c *Type_bodyContext)

	// EnterType_body_elements is called when entering the type_body_elements production.
	EnterType_body_elements(c *Type_body_elementsContext)

	// EnterMap_order_func_declaration is called when entering the map_order_func_declaration production.
	EnterMap_order_func_declaration(c *Map_order_func_declarationContext)

	// EnterSubprog_decl_in_type is called when entering the subprog_decl_in_type production.
	EnterSubprog_decl_in_type(c *Subprog_decl_in_typeContext)

	// EnterProc_decl_in_type is called when entering the proc_decl_in_type production.
	EnterProc_decl_in_type(c *Proc_decl_in_typeContext)

	// EnterFunc_decl_in_type is called when entering the func_decl_in_type production.
	EnterFunc_decl_in_type(c *Func_decl_in_typeContext)

	// EnterConstructor_declaration is called when entering the constructor_declaration production.
	EnterConstructor_declaration(c *Constructor_declarationContext)

	// EnterModifier_clause is called when entering the modifier_clause production.
	EnterModifier_clause(c *Modifier_clauseContext)

	// EnterObject_member_spec is called when entering the object_member_spec production.
	EnterObject_member_spec(c *Object_member_specContext)

	// EnterSqlj_object_type_attr is called when entering the sqlj_object_type_attr production.
	EnterSqlj_object_type_attr(c *Sqlj_object_type_attrContext)

	// EnterElement_spec is called when entering the element_spec production.
	EnterElement_spec(c *Element_specContext)

	// EnterElement_spec_options is called when entering the element_spec_options production.
	EnterElement_spec_options(c *Element_spec_optionsContext)

	// EnterSubprogram_spec is called when entering the subprogram_spec production.
	EnterSubprogram_spec(c *Subprogram_specContext)

	// EnterOverriding_subprogram_spec is called when entering the overriding_subprogram_spec production.
	EnterOverriding_subprogram_spec(c *Overriding_subprogram_specContext)

	// EnterOverriding_function_spec is called when entering the overriding_function_spec production.
	EnterOverriding_function_spec(c *Overriding_function_specContext)

	// EnterOverriding_procedure_spec is called when entering the overriding_procedure_spec production.
	EnterOverriding_procedure_spec(c *Overriding_procedure_specContext)

	// EnterType_procedure_spec is called when entering the type_procedure_spec production.
	EnterType_procedure_spec(c *Type_procedure_specContext)

	// EnterType_function_spec is called when entering the type_function_spec production.
	EnterType_function_spec(c *Type_function_specContext)

	// EnterConstructor_spec is called when entering the constructor_spec production.
	EnterConstructor_spec(c *Constructor_specContext)

	// EnterMap_order_function_spec is called when entering the map_order_function_spec production.
	EnterMap_order_function_spec(c *Map_order_function_specContext)

	// EnterPragma_clause is called when entering the pragma_clause production.
	EnterPragma_clause(c *Pragma_clauseContext)

	// EnterPragma_elements is called when entering the pragma_elements production.
	EnterPragma_elements(c *Pragma_elementsContext)

	// EnterType_elements_parameter is called when entering the type_elements_parameter production.
	EnterType_elements_parameter(c *Type_elements_parameterContext)

	// EnterDrop_sequence is called when entering the drop_sequence production.
	EnterDrop_sequence(c *Drop_sequenceContext)

	// EnterAlter_sequence is called when entering the alter_sequence production.
	EnterAlter_sequence(c *Alter_sequenceContext)

	// EnterAlter_session is called when entering the alter_session production.
	EnterAlter_session(c *Alter_sessionContext)

	// EnterAlter_session_set_clause is called when entering the alter_session_set_clause production.
	EnterAlter_session_set_clause(c *Alter_session_set_clauseContext)

	// EnterCreate_sequence is called when entering the create_sequence production.
	EnterCreate_sequence(c *Create_sequenceContext)

	// EnterSequence_spec is called when entering the sequence_spec production.
	EnterSequence_spec(c *Sequence_specContext)

	// EnterSequence_start_clause is called when entering the sequence_start_clause production.
	EnterSequence_start_clause(c *Sequence_start_clauseContext)

	// EnterCreate_analytic_view is called when entering the create_analytic_view production.
	EnterCreate_analytic_view(c *Create_analytic_viewContext)

	// EnterClassification_clause is called when entering the classification_clause production.
	EnterClassification_clause(c *Classification_clauseContext)

	// EnterCaption_clause is called when entering the caption_clause production.
	EnterCaption_clause(c *Caption_clauseContext)

	// EnterDescription_clause is called when entering the description_clause production.
	EnterDescription_clause(c *Description_clauseContext)

	// EnterClassification_item is called when entering the classification_item production.
	EnterClassification_item(c *Classification_itemContext)

	// EnterLanguage is called when entering the language production.
	EnterLanguage(c *LanguageContext)

	// EnterCav_using_clause is called when entering the cav_using_clause production.
	EnterCav_using_clause(c *Cav_using_clauseContext)

	// EnterDim_by_clause is called when entering the dim_by_clause production.
	EnterDim_by_clause(c *Dim_by_clauseContext)

	// EnterDim_key is called when entering the dim_key production.
	EnterDim_key(c *Dim_keyContext)

	// EnterDim_ref is called when entering the dim_ref production.
	EnterDim_ref(c *Dim_refContext)

	// EnterHier_ref is called when entering the hier_ref production.
	EnterHier_ref(c *Hier_refContext)

	// EnterMeasures_clause is called when entering the measures_clause production.
	EnterMeasures_clause(c *Measures_clauseContext)

	// EnterAv_measure is called when entering the av_measure production.
	EnterAv_measure(c *Av_measureContext)

	// EnterBase_meas_clause is called when entering the base_meas_clause production.
	EnterBase_meas_clause(c *Base_meas_clauseContext)

	// EnterMeas_aggregate_clause is called when entering the meas_aggregate_clause production.
	EnterMeas_aggregate_clause(c *Meas_aggregate_clauseContext)

	// EnterCalc_meas_clause is called when entering the calc_meas_clause production.
	EnterCalc_meas_clause(c *Calc_meas_clauseContext)

	// EnterDefault_measure_clause is called when entering the default_measure_clause production.
	EnterDefault_measure_clause(c *Default_measure_clauseContext)

	// EnterDefault_aggregate_clause is called when entering the default_aggregate_clause production.
	EnterDefault_aggregate_clause(c *Default_aggregate_clauseContext)

	// EnterCache_clause is called when entering the cache_clause production.
	EnterCache_clause(c *Cache_clauseContext)

	// EnterCache_specification is called when entering the cache_specification production.
	EnterCache_specification(c *Cache_specificationContext)

	// EnterLevels_clause is called when entering the levels_clause production.
	EnterLevels_clause(c *Levels_clauseContext)

	// EnterLevel_specification is called when entering the level_specification production.
	EnterLevel_specification(c *Level_specificationContext)

	// EnterLevel_group_type is called when entering the level_group_type production.
	EnterLevel_group_type(c *Level_group_typeContext)

	// EnterFact_columns_clause is called when entering the fact_columns_clause production.
	EnterFact_columns_clause(c *Fact_columns_clauseContext)

	// EnterQry_transform_clause is called when entering the qry_transform_clause production.
	EnterQry_transform_clause(c *Qry_transform_clauseContext)

	// EnterCreate_attribute_dimension is called when entering the create_attribute_dimension production.
	EnterCreate_attribute_dimension(c *Create_attribute_dimensionContext)

	// EnterAd_using_clause is called when entering the ad_using_clause production.
	EnterAd_using_clause(c *Ad_using_clauseContext)

	// EnterSource_clause is called when entering the source_clause production.
	EnterSource_clause(c *Source_clauseContext)

	// EnterJoin_path_clause is called when entering the join_path_clause production.
	EnterJoin_path_clause(c *Join_path_clauseContext)

	// EnterJoin_condition is called when entering the join_condition production.
	EnterJoin_condition(c *Join_conditionContext)

	// EnterJoin_condition_item is called when entering the join_condition_item production.
	EnterJoin_condition_item(c *Join_condition_itemContext)

	// EnterAttributes_clause is called when entering the attributes_clause production.
	EnterAttributes_clause(c *Attributes_clauseContext)

	// EnterAd_attributes_clause is called when entering the ad_attributes_clause production.
	EnterAd_attributes_clause(c *Ad_attributes_clauseContext)

	// EnterAd_level_clause is called when entering the ad_level_clause production.
	EnterAd_level_clause(c *Ad_level_clauseContext)

	// EnterKey_clause is called when entering the key_clause production.
	EnterKey_clause(c *Key_clauseContext)

	// EnterAlternate_key_clause is called when entering the alternate_key_clause production.
	EnterAlternate_key_clause(c *Alternate_key_clauseContext)

	// EnterDim_order_clause is called when entering the dim_order_clause production.
	EnterDim_order_clause(c *Dim_order_clauseContext)

	// EnterAll_clause is called when entering the all_clause production.
	EnterAll_clause(c *All_clauseContext)

	// EnterCreate_audit_policy is called when entering the create_audit_policy production.
	EnterCreate_audit_policy(c *Create_audit_policyContext)

	// EnterPrivilege_audit_clause is called when entering the privilege_audit_clause production.
	EnterPrivilege_audit_clause(c *Privilege_audit_clauseContext)

	// EnterAction_audit_clause is called when entering the action_audit_clause production.
	EnterAction_audit_clause(c *Action_audit_clauseContext)

	// EnterSystem_actions is called when entering the system_actions production.
	EnterSystem_actions(c *System_actionsContext)

	// EnterStandard_actions is called when entering the standard_actions production.
	EnterStandard_actions(c *Standard_actionsContext)

	// EnterActions_clause is called when entering the actions_clause production.
	EnterActions_clause(c *Actions_clauseContext)

	// EnterObject_action is called when entering the object_action production.
	EnterObject_action(c *Object_actionContext)

	// EnterSystem_action is called when entering the system_action production.
	EnterSystem_action(c *System_actionContext)

	// EnterComponent_actions is called when entering the component_actions production.
	EnterComponent_actions(c *Component_actionsContext)

	// EnterComponent_action is called when entering the component_action production.
	EnterComponent_action(c *Component_actionContext)

	// EnterRole_audit_clause is called when entering the role_audit_clause production.
	EnterRole_audit_clause(c *Role_audit_clauseContext)

	// EnterCreate_controlfile is called when entering the create_controlfile production.
	EnterCreate_controlfile(c *Create_controlfileContext)

	// EnterControlfile_options is called when entering the controlfile_options production.
	EnterControlfile_options(c *Controlfile_optionsContext)

	// EnterLogfile_clause is called when entering the logfile_clause production.
	EnterLogfile_clause(c *Logfile_clauseContext)

	// EnterCharacter_set_clause is called when entering the character_set_clause production.
	EnterCharacter_set_clause(c *Character_set_clauseContext)

	// EnterFile_specification is called when entering the file_specification production.
	EnterFile_specification(c *File_specificationContext)

	// EnterCreate_diskgroup is called when entering the create_diskgroup production.
	EnterCreate_diskgroup(c *Create_diskgroupContext)

	// EnterQualified_disk_clause is called when entering the qualified_disk_clause production.
	EnterQualified_disk_clause(c *Qualified_disk_clauseContext)

	// EnterCreate_edition is called when entering the create_edition production.
	EnterCreate_edition(c *Create_editionContext)

	// EnterCreate_flashback_archive is called when entering the create_flashback_archive production.
	EnterCreate_flashback_archive(c *Create_flashback_archiveContext)

	// EnterFlashback_archive_quota is called when entering the flashback_archive_quota production.
	EnterFlashback_archive_quota(c *Flashback_archive_quotaContext)

	// EnterFlashback_archive_retention is called when entering the flashback_archive_retention production.
	EnterFlashback_archive_retention(c *Flashback_archive_retentionContext)

	// EnterCreate_hierarchy is called when entering the create_hierarchy production.
	EnterCreate_hierarchy(c *Create_hierarchyContext)

	// EnterHier_using_clause is called when entering the hier_using_clause production.
	EnterHier_using_clause(c *Hier_using_clauseContext)

	// EnterLevel_hier_clause is called when entering the level_hier_clause production.
	EnterLevel_hier_clause(c *Level_hier_clauseContext)

	// EnterHier_attrs_clause is called when entering the hier_attrs_clause production.
	EnterHier_attrs_clause(c *Hier_attrs_clauseContext)

	// EnterHier_attr_clause is called when entering the hier_attr_clause production.
	EnterHier_attr_clause(c *Hier_attr_clauseContext)

	// EnterHier_attr_name is called when entering the hier_attr_name production.
	EnterHier_attr_name(c *Hier_attr_nameContext)

	// EnterCreate_index is called when entering the create_index production.
	EnterCreate_index(c *Create_indexContext)

	// EnterCluster_index_clause is called when entering the cluster_index_clause production.
	EnterCluster_index_clause(c *Cluster_index_clauseContext)

	// EnterCluster_name is called when entering the cluster_name production.
	EnterCluster_name(c *Cluster_nameContext)

	// EnterTable_index_clause is called when entering the table_index_clause production.
	EnterTable_index_clause(c *Table_index_clauseContext)

	// EnterBitmap_join_index_clause is called when entering the bitmap_join_index_clause production.
	EnterBitmap_join_index_clause(c *Bitmap_join_index_clauseContext)

	// EnterIndex_expr is called when entering the index_expr production.
	EnterIndex_expr(c *Index_exprContext)

	// EnterIndex_properties is called when entering the index_properties production.
	EnterIndex_properties(c *Index_propertiesContext)

	// EnterDomain_index_clause is called when entering the domain_index_clause production.
	EnterDomain_index_clause(c *Domain_index_clauseContext)

	// EnterLocal_domain_index_clause is called when entering the local_domain_index_clause production.
	EnterLocal_domain_index_clause(c *Local_domain_index_clauseContext)

	// EnterXmlindex_clause is called when entering the xmlindex_clause production.
	EnterXmlindex_clause(c *Xmlindex_clauseContext)

	// EnterLocal_xmlindex_clause is called when entering the local_xmlindex_clause production.
	EnterLocal_xmlindex_clause(c *Local_xmlindex_clauseContext)

	// EnterGlobal_partitioned_index is called when entering the global_partitioned_index production.
	EnterGlobal_partitioned_index(c *Global_partitioned_indexContext)

	// EnterIndex_partitioning_clause is called when entering the index_partitioning_clause production.
	EnterIndex_partitioning_clause(c *Index_partitioning_clauseContext)

	// EnterIndex_partitioning_values_list is called when entering the index_partitioning_values_list production.
	EnterIndex_partitioning_values_list(c *Index_partitioning_values_listContext)

	// EnterLocal_partitioned_index is called when entering the local_partitioned_index production.
	EnterLocal_partitioned_index(c *Local_partitioned_indexContext)

	// EnterOn_range_partitioned_table is called when entering the on_range_partitioned_table production.
	EnterOn_range_partitioned_table(c *On_range_partitioned_tableContext)

	// EnterOn_list_partitioned_table is called when entering the on_list_partitioned_table production.
	EnterOn_list_partitioned_table(c *On_list_partitioned_tableContext)

	// EnterPartitioned_table is called when entering the partitioned_table production.
	EnterPartitioned_table(c *Partitioned_tableContext)

	// EnterOn_hash_partitioned_table is called when entering the on_hash_partitioned_table production.
	EnterOn_hash_partitioned_table(c *On_hash_partitioned_tableContext)

	// EnterOn_hash_partitioned_clause is called when entering the on_hash_partitioned_clause production.
	EnterOn_hash_partitioned_clause(c *On_hash_partitioned_clauseContext)

	// EnterOn_comp_partitioned_table is called when entering the on_comp_partitioned_table production.
	EnterOn_comp_partitioned_table(c *On_comp_partitioned_tableContext)

	// EnterOn_comp_partitioned_clause is called when entering the on_comp_partitioned_clause production.
	EnterOn_comp_partitioned_clause(c *On_comp_partitioned_clauseContext)

	// EnterIndex_subpartition_clause is called when entering the index_subpartition_clause production.
	EnterIndex_subpartition_clause(c *Index_subpartition_clauseContext)

	// EnterIndex_subpartition_subclause is called when entering the index_subpartition_subclause production.
	EnterIndex_subpartition_subclause(c *Index_subpartition_subclauseContext)

	// EnterOdci_parameters is called when entering the odci_parameters production.
	EnterOdci_parameters(c *Odci_parametersContext)

	// EnterIndextype is called when entering the indextype production.
	EnterIndextype(c *IndextypeContext)

	// EnterAlter_index is called when entering the alter_index production.
	EnterAlter_index(c *Alter_indexContext)

	// EnterAlter_index_ops_set1 is called when entering the alter_index_ops_set1 production.
	EnterAlter_index_ops_set1(c *Alter_index_ops_set1Context)

	// EnterAlter_index_ops_set2 is called when entering the alter_index_ops_set2 production.
	EnterAlter_index_ops_set2(c *Alter_index_ops_set2Context)

	// EnterVisible_or_invisible is called when entering the visible_or_invisible production.
	EnterVisible_or_invisible(c *Visible_or_invisibleContext)

	// EnterMonitoring_nomonitoring is called when entering the monitoring_nomonitoring production.
	EnterMonitoring_nomonitoring(c *Monitoring_nomonitoringContext)

	// EnterRebuild_clause is called when entering the rebuild_clause production.
	EnterRebuild_clause(c *Rebuild_clauseContext)

	// EnterAlter_index_partitioning is called when entering the alter_index_partitioning production.
	EnterAlter_index_partitioning(c *Alter_index_partitioningContext)

	// EnterModify_index_default_attrs is called when entering the modify_index_default_attrs production.
	EnterModify_index_default_attrs(c *Modify_index_default_attrsContext)

	// EnterAdd_hash_index_partition is called when entering the add_hash_index_partition production.
	EnterAdd_hash_index_partition(c *Add_hash_index_partitionContext)

	// EnterCoalesce_index_partition is called when entering the coalesce_index_partition production.
	EnterCoalesce_index_partition(c *Coalesce_index_partitionContext)

	// EnterModify_index_partition is called when entering the modify_index_partition production.
	EnterModify_index_partition(c *Modify_index_partitionContext)

	// EnterModify_index_partitions_ops is called when entering the modify_index_partitions_ops production.
	EnterModify_index_partitions_ops(c *Modify_index_partitions_opsContext)

	// EnterRename_index_partition is called when entering the rename_index_partition production.
	EnterRename_index_partition(c *Rename_index_partitionContext)

	// EnterDrop_index_partition is called when entering the drop_index_partition production.
	EnterDrop_index_partition(c *Drop_index_partitionContext)

	// EnterSplit_index_partition is called when entering the split_index_partition production.
	EnterSplit_index_partition(c *Split_index_partitionContext)

	// EnterIndex_partition_description is called when entering the index_partition_description production.
	EnterIndex_partition_description(c *Index_partition_descriptionContext)

	// EnterModify_index_subpartition is called when entering the modify_index_subpartition production.
	EnterModify_index_subpartition(c *Modify_index_subpartitionContext)

	// EnterPartition_name_old is called when entering the partition_name_old production.
	EnterPartition_name_old(c *Partition_name_oldContext)

	// EnterNew_partition_name is called when entering the new_partition_name production.
	EnterNew_partition_name(c *New_partition_nameContext)

	// EnterNew_index_name is called when entering the new_index_name production.
	EnterNew_index_name(c *New_index_nameContext)

	// EnterAlter_inmemory_join_group is called when entering the alter_inmemory_join_group production.
	EnterAlter_inmemory_join_group(c *Alter_inmemory_join_groupContext)

	// EnterCreate_user is called when entering the create_user production.
	EnterCreate_user(c *Create_userContext)

	// EnterAlter_user is called when entering the alter_user production.
	EnterAlter_user(c *Alter_userContext)

	// EnterDrop_user is called when entering the drop_user production.
	EnterDrop_user(c *Drop_userContext)

	// EnterAlter_identified_by is called when entering the alter_identified_by production.
	EnterAlter_identified_by(c *Alter_identified_byContext)

	// EnterIdentified_by is called when entering the identified_by production.
	EnterIdentified_by(c *Identified_byContext)

	// EnterIdentified_other_clause is called when entering the identified_other_clause production.
	EnterIdentified_other_clause(c *Identified_other_clauseContext)

	// EnterUser_tablespace_clause is called when entering the user_tablespace_clause production.
	EnterUser_tablespace_clause(c *User_tablespace_clauseContext)

	// EnterQuota_clause is called when entering the quota_clause production.
	EnterQuota_clause(c *Quota_clauseContext)

	// EnterProfile_clause is called when entering the profile_clause production.
	EnterProfile_clause(c *Profile_clauseContext)

	// EnterRole_clause is called when entering the role_clause production.
	EnterRole_clause(c *Role_clauseContext)

	// EnterUser_default_role_clause is called when entering the user_default_role_clause production.
	EnterUser_default_role_clause(c *User_default_role_clauseContext)

	// EnterPassword_expire_clause is called when entering the password_expire_clause production.
	EnterPassword_expire_clause(c *Password_expire_clauseContext)

	// EnterUser_lock_clause is called when entering the user_lock_clause production.
	EnterUser_lock_clause(c *User_lock_clauseContext)

	// EnterUser_editions_clause is called when entering the user_editions_clause production.
	EnterUser_editions_clause(c *User_editions_clauseContext)

	// EnterAlter_user_editions_clause is called when entering the alter_user_editions_clause production.
	EnterAlter_user_editions_clause(c *Alter_user_editions_clauseContext)

	// EnterProxy_clause is called when entering the proxy_clause production.
	EnterProxy_clause(c *Proxy_clauseContext)

	// EnterContainer_names is called when entering the container_names production.
	EnterContainer_names(c *Container_namesContext)

	// EnterSet_container_data is called when entering the set_container_data production.
	EnterSet_container_data(c *Set_container_dataContext)

	// EnterAdd_rem_container_data is called when entering the add_rem_container_data production.
	EnterAdd_rem_container_data(c *Add_rem_container_dataContext)

	// EnterContainer_data_clause is called when entering the container_data_clause production.
	EnterContainer_data_clause(c *Container_data_clauseContext)

	// EnterAdminister_key_management is called when entering the administer_key_management production.
	EnterAdminister_key_management(c *Administer_key_managementContext)

	// EnterKeystore_management_clauses is called when entering the keystore_management_clauses production.
	EnterKeystore_management_clauses(c *Keystore_management_clausesContext)

	// EnterCreate_keystore is called when entering the create_keystore production.
	EnterCreate_keystore(c *Create_keystoreContext)

	// EnterOpen_keystore is called when entering the open_keystore production.
	EnterOpen_keystore(c *Open_keystoreContext)

	// EnterForce_keystore is called when entering the force_keystore production.
	EnterForce_keystore(c *Force_keystoreContext)

	// EnterClose_keystore is called when entering the close_keystore production.
	EnterClose_keystore(c *Close_keystoreContext)

	// EnterBackup_keystore is called when entering the backup_keystore production.
	EnterBackup_keystore(c *Backup_keystoreContext)

	// EnterAlter_keystore_password is called when entering the alter_keystore_password production.
	EnterAlter_keystore_password(c *Alter_keystore_passwordContext)

	// EnterMerge_into_new_keystore is called when entering the merge_into_new_keystore production.
	EnterMerge_into_new_keystore(c *Merge_into_new_keystoreContext)

	// EnterMerge_into_existing_keystore is called when entering the merge_into_existing_keystore production.
	EnterMerge_into_existing_keystore(c *Merge_into_existing_keystoreContext)

	// EnterIsolate_keystore is called when entering the isolate_keystore production.
	EnterIsolate_keystore(c *Isolate_keystoreContext)

	// EnterUnite_keystore is called when entering the unite_keystore production.
	EnterUnite_keystore(c *Unite_keystoreContext)

	// EnterKey_management_clauses is called when entering the key_management_clauses production.
	EnterKey_management_clauses(c *Key_management_clausesContext)

	// EnterSet_key is called when entering the set_key production.
	EnterSet_key(c *Set_keyContext)

	// EnterCreate_key is called when entering the create_key production.
	EnterCreate_key(c *Create_keyContext)

	// EnterMkid is called when entering the mkid production.
	EnterMkid(c *MkidContext)

	// EnterMk is called when entering the mk production.
	EnterMk(c *MkContext)

	// EnterUse_key is called when entering the use_key production.
	EnterUse_key(c *Use_keyContext)

	// EnterSet_key_tag is called when entering the set_key_tag production.
	EnterSet_key_tag(c *Set_key_tagContext)

	// EnterExport_keys is called when entering the export_keys production.
	EnterExport_keys(c *Export_keysContext)

	// EnterImport_keys is called when entering the import_keys production.
	EnterImport_keys(c *Import_keysContext)

	// EnterMigrate_keys is called when entering the migrate_keys production.
	EnterMigrate_keys(c *Migrate_keysContext)

	// EnterReverse_migrate_keys is called when entering the reverse_migrate_keys production.
	EnterReverse_migrate_keys(c *Reverse_migrate_keysContext)

	// EnterMove_keys is called when entering the move_keys production.
	EnterMove_keys(c *Move_keysContext)

	// EnterIdentified_by_store is called when entering the identified_by_store production.
	EnterIdentified_by_store(c *Identified_by_storeContext)

	// EnterUsing_algorithm_clause is called when entering the using_algorithm_clause production.
	EnterUsing_algorithm_clause(c *Using_algorithm_clauseContext)

	// EnterUsing_tag_clause is called when entering the using_tag_clause production.
	EnterUsing_tag_clause(c *Using_tag_clauseContext)

	// EnterSecret_management_clauses is called when entering the secret_management_clauses production.
	EnterSecret_management_clauses(c *Secret_management_clausesContext)

	// EnterAdd_update_secret is called when entering the add_update_secret production.
	EnterAdd_update_secret(c *Add_update_secretContext)

	// EnterDelete_secret is called when entering the delete_secret production.
	EnterDelete_secret(c *Delete_secretContext)

	// EnterAdd_update_secret_seps is called when entering the add_update_secret_seps production.
	EnterAdd_update_secret_seps(c *Add_update_secret_sepsContext)

	// EnterDelete_secret_seps is called when entering the delete_secret_seps production.
	EnterDelete_secret_seps(c *Delete_secret_sepsContext)

	// EnterZero_downtime_software_patching_clauses is called when entering the zero_downtime_software_patching_clauses production.
	EnterZero_downtime_software_patching_clauses(c *Zero_downtime_software_patching_clausesContext)

	// EnterWith_backup_clause is called when entering the with_backup_clause production.
	EnterWith_backup_clause(c *With_backup_clauseContext)

	// EnterIdentified_by_password_clause is called when entering the identified_by_password_clause production.
	EnterIdentified_by_password_clause(c *Identified_by_password_clauseContext)

	// EnterKeystore_password is called when entering the keystore_password production.
	EnterKeystore_password(c *Keystore_passwordContext)

	// EnterPath is called when entering the path production.
	EnterPath(c *PathContext)

	// EnterSecret is called when entering the secret production.
	EnterSecret(c *SecretContext)

	// EnterAnalyze is called when entering the analyze production.
	EnterAnalyze(c *AnalyzeContext)

	// EnterPartition_extention_clause is called when entering the partition_extention_clause production.
	EnterPartition_extention_clause(c *Partition_extention_clauseContext)

	// EnterValidation_clauses is called when entering the validation_clauses production.
	EnterValidation_clauses(c *Validation_clausesContext)

	// EnterCompute_clauses is called when entering the compute_clauses production.
	EnterCompute_clauses(c *Compute_clausesContext)

	// EnterFor_clause is called when entering the for_clause production.
	EnterFor_clause(c *For_clauseContext)

	// EnterOnline_or_offline is called when entering the online_or_offline production.
	EnterOnline_or_offline(c *Online_or_offlineContext)

	// EnterInto_clause1 is called when entering the into_clause1 production.
	EnterInto_clause1(c *Into_clause1Context)

	// EnterPartition_key_value is called when entering the partition_key_value production.
	EnterPartition_key_value(c *Partition_key_valueContext)

	// EnterSubpartition_key_value is called when entering the subpartition_key_value production.
	EnterSubpartition_key_value(c *Subpartition_key_valueContext)

	// EnterAssociate_statistics is called when entering the associate_statistics production.
	EnterAssociate_statistics(c *Associate_statisticsContext)

	// EnterColumn_association is called when entering the column_association production.
	EnterColumn_association(c *Column_associationContext)

	// EnterFunction_association is called when entering the function_association production.
	EnterFunction_association(c *Function_associationContext)

	// EnterIndextype_name is called when entering the indextype_name production.
	EnterIndextype_name(c *Indextype_nameContext)

	// EnterUsing_statistics_type is called when entering the using_statistics_type production.
	EnterUsing_statistics_type(c *Using_statistics_typeContext)

	// EnterStatistics_type_name is called when entering the statistics_type_name production.
	EnterStatistics_type_name(c *Statistics_type_nameContext)

	// EnterDefault_cost_clause is called when entering the default_cost_clause production.
	EnterDefault_cost_clause(c *Default_cost_clauseContext)

	// EnterCpu_cost is called when entering the cpu_cost production.
	EnterCpu_cost(c *Cpu_costContext)

	// EnterIo_cost is called when entering the io_cost production.
	EnterIo_cost(c *Io_costContext)

	// EnterNetwork_cost is called when entering the network_cost production.
	EnterNetwork_cost(c *Network_costContext)

	// EnterDefault_selectivity_clause is called when entering the default_selectivity_clause production.
	EnterDefault_selectivity_clause(c *Default_selectivity_clauseContext)

	// EnterDefault_selectivity is called when entering the default_selectivity production.
	EnterDefault_selectivity(c *Default_selectivityContext)

	// EnterStorage_table_clause is called when entering the storage_table_clause production.
	EnterStorage_table_clause(c *Storage_table_clauseContext)

	// EnterUnified_auditing is called when entering the unified_auditing production.
	EnterUnified_auditing(c *Unified_auditingContext)

	// EnterPolicy_name is called when entering the policy_name production.
	EnterPolicy_name(c *Policy_nameContext)

	// EnterAudit_traditional is called when entering the audit_traditional production.
	EnterAudit_traditional(c *Audit_traditionalContext)

	// EnterAudit_direct_path is called when entering the audit_direct_path production.
	EnterAudit_direct_path(c *Audit_direct_pathContext)

	// EnterAudit_container_clause is called when entering the audit_container_clause production.
	EnterAudit_container_clause(c *Audit_container_clauseContext)

	// EnterAudit_operation_clause is called when entering the audit_operation_clause production.
	EnterAudit_operation_clause(c *Audit_operation_clauseContext)

	// EnterAuditing_by_clause is called when entering the auditing_by_clause production.
	EnterAuditing_by_clause(c *Auditing_by_clauseContext)

	// EnterAudit_user is called when entering the audit_user production.
	EnterAudit_user(c *Audit_userContext)

	// EnterAudit_schema_object_clause is called when entering the audit_schema_object_clause production.
	EnterAudit_schema_object_clause(c *Audit_schema_object_clauseContext)

	// EnterSql_operation is called when entering the sql_operation production.
	EnterSql_operation(c *Sql_operationContext)

	// EnterAuditing_on_clause is called when entering the auditing_on_clause production.
	EnterAuditing_on_clause(c *Auditing_on_clauseContext)

	// EnterModel_name is called when entering the model_name production.
	EnterModel_name(c *Model_nameContext)

	// EnterObject_name is called when entering the object_name production.
	EnterObject_name(c *Object_nameContext)

	// EnterProfile_name is called when entering the profile_name production.
	EnterProfile_name(c *Profile_nameContext)

	// EnterSql_statement_shortcut is called when entering the sql_statement_shortcut production.
	EnterSql_statement_shortcut(c *Sql_statement_shortcutContext)

	// EnterDrop_index is called when entering the drop_index production.
	EnterDrop_index(c *Drop_indexContext)

	// EnterDisassociate_statistics is called when entering the disassociate_statistics production.
	EnterDisassociate_statistics(c *Disassociate_statisticsContext)

	// EnterDrop_indextype is called when entering the drop_indextype production.
	EnterDrop_indextype(c *Drop_indextypeContext)

	// EnterDrop_inmemory_join_group is called when entering the drop_inmemory_join_group production.
	EnterDrop_inmemory_join_group(c *Drop_inmemory_join_groupContext)

	// EnterFlashback_table is called when entering the flashback_table production.
	EnterFlashback_table(c *Flashback_tableContext)

	// EnterRestore_point is called when entering the restore_point production.
	EnterRestore_point(c *Restore_pointContext)

	// EnterPurge_statement is called when entering the purge_statement production.
	EnterPurge_statement(c *Purge_statementContext)

	// EnterNoaudit_statement is called when entering the noaudit_statement production.
	EnterNoaudit_statement(c *Noaudit_statementContext)

	// EnterRename_object is called when entering the rename_object production.
	EnterRename_object(c *Rename_objectContext)

	// EnterGrant_statement is called when entering the grant_statement production.
	EnterGrant_statement(c *Grant_statementContext)

	// EnterContainer_clause is called when entering the container_clause production.
	EnterContainer_clause(c *Container_clauseContext)

	// EnterRevoke_statement is called when entering the revoke_statement production.
	EnterRevoke_statement(c *Revoke_statementContext)

	// EnterRevoke_system_privilege is called when entering the revoke_system_privilege production.
	EnterRevoke_system_privilege(c *Revoke_system_privilegeContext)

	// EnterRevokee_clause is called when entering the revokee_clause production.
	EnterRevokee_clause(c *Revokee_clauseContext)

	// EnterRevoke_object_privileges is called when entering the revoke_object_privileges production.
	EnterRevoke_object_privileges(c *Revoke_object_privilegesContext)

	// EnterOn_object_clause is called when entering the on_object_clause production.
	EnterOn_object_clause(c *On_object_clauseContext)

	// EnterRevoke_roles_from_programs is called when entering the revoke_roles_from_programs production.
	EnterRevoke_roles_from_programs(c *Revoke_roles_from_programsContext)

	// EnterProgram_unit is called when entering the program_unit production.
	EnterProgram_unit(c *Program_unitContext)

	// EnterCreate_dimension is called when entering the create_dimension production.
	EnterCreate_dimension(c *Create_dimensionContext)

	// EnterCreate_directory is called when entering the create_directory production.
	EnterCreate_directory(c *Create_directoryContext)

	// EnterDirectory_name is called when entering the directory_name production.
	EnterDirectory_name(c *Directory_nameContext)

	// EnterDirectory_path is called when entering the directory_path production.
	EnterDirectory_path(c *Directory_pathContext)

	// EnterCreate_inmemory_join_group is called when entering the create_inmemory_join_group production.
	EnterCreate_inmemory_join_group(c *Create_inmemory_join_groupContext)

	// EnterDrop_hierarchy is called when entering the drop_hierarchy production.
	EnterDrop_hierarchy(c *Drop_hierarchyContext)

	// EnterAlter_library is called when entering the alter_library production.
	EnterAlter_library(c *Alter_libraryContext)

	// EnterDrop_java is called when entering the drop_java production.
	EnterDrop_java(c *Drop_javaContext)

	// EnterDrop_library is called when entering the drop_library production.
	EnterDrop_library(c *Drop_libraryContext)

	// EnterCreate_java is called when entering the create_java production.
	EnterCreate_java(c *Create_javaContext)

	// EnterCreate_library is called when entering the create_library production.
	EnterCreate_library(c *Create_libraryContext)

	// EnterPlsql_library_source is called when entering the plsql_library_source production.
	EnterPlsql_library_source(c *Plsql_library_sourceContext)

	// EnterCredential_name is called when entering the credential_name production.
	EnterCredential_name(c *Credential_nameContext)

	// EnterLibrary_editionable is called when entering the library_editionable production.
	EnterLibrary_editionable(c *Library_editionableContext)

	// EnterLibrary_debug is called when entering the library_debug production.
	EnterLibrary_debug(c *Library_debugContext)

	// EnterCompiler_parameters_clause is called when entering the compiler_parameters_clause production.
	EnterCompiler_parameters_clause(c *Compiler_parameters_clauseContext)

	// EnterParameter_value is called when entering the parameter_value production.
	EnterParameter_value(c *Parameter_valueContext)

	// EnterLibrary_name is called when entering the library_name production.
	EnterLibrary_name(c *Library_nameContext)

	// EnterAlter_dimension is called when entering the alter_dimension production.
	EnterAlter_dimension(c *Alter_dimensionContext)

	// EnterLevel_clause is called when entering the level_clause production.
	EnterLevel_clause(c *Level_clauseContext)

	// EnterHierarchy_clause is called when entering the hierarchy_clause production.
	EnterHierarchy_clause(c *Hierarchy_clauseContext)

	// EnterDimension_join_clause is called when entering the dimension_join_clause production.
	EnterDimension_join_clause(c *Dimension_join_clauseContext)

	// EnterAttribute_clause is called when entering the attribute_clause production.
	EnterAttribute_clause(c *Attribute_clauseContext)

	// EnterExtended_attribute_clause is called when entering the extended_attribute_clause production.
	EnterExtended_attribute_clause(c *Extended_attribute_clauseContext)

	// EnterColumn_one_or_more_sub_clause is called when entering the column_one_or_more_sub_clause production.
	EnterColumn_one_or_more_sub_clause(c *Column_one_or_more_sub_clauseContext)

	// EnterAlter_view is called when entering the alter_view production.
	EnterAlter_view(c *Alter_viewContext)

	// EnterAlter_view_editionable is called when entering the alter_view_editionable production.
	EnterAlter_view_editionable(c *Alter_view_editionableContext)

	// EnterCreate_view is called when entering the create_view production.
	EnterCreate_view(c *Create_viewContext)

	// EnterEditioning_clause is called when entering the editioning_clause production.
	EnterEditioning_clause(c *Editioning_clauseContext)

	// EnterView_options is called when entering the view_options production.
	EnterView_options(c *View_optionsContext)

	// EnterView_alias_constraint is called when entering the view_alias_constraint production.
	EnterView_alias_constraint(c *View_alias_constraintContext)

	// EnterObject_view_clause is called when entering the object_view_clause production.
	EnterObject_view_clause(c *Object_view_clauseContext)

	// EnterInline_constraint is called when entering the inline_constraint production.
	EnterInline_constraint(c *Inline_constraintContext)

	// EnterInline_ref_constraint is called when entering the inline_ref_constraint production.
	EnterInline_ref_constraint(c *Inline_ref_constraintContext)

	// EnterOut_of_line_ref_constraint is called when entering the out_of_line_ref_constraint production.
	EnterOut_of_line_ref_constraint(c *Out_of_line_ref_constraintContext)

	// EnterOut_of_line_constraint is called when entering the out_of_line_constraint production.
	EnterOut_of_line_constraint(c *Out_of_line_constraintContext)

	// EnterConstraint_state is called when entering the constraint_state production.
	EnterConstraint_state(c *Constraint_stateContext)

	// EnterXmltype_view_clause is called when entering the xmltype_view_clause production.
	EnterXmltype_view_clause(c *Xmltype_view_clauseContext)

	// EnterXml_schema_spec is called when entering the xml_schema_spec production.
	EnterXml_schema_spec(c *Xml_schema_specContext)

	// EnterXml_schema_url is called when entering the xml_schema_url production.
	EnterXml_schema_url(c *Xml_schema_urlContext)

	// EnterElement is called when entering the element production.
	EnterElement(c *ElementContext)

	// EnterAlter_tablespace is called when entering the alter_tablespace production.
	EnterAlter_tablespace(c *Alter_tablespaceContext)

	// EnterDatafile_tempfile_clauses is called when entering the datafile_tempfile_clauses production.
	EnterDatafile_tempfile_clauses(c *Datafile_tempfile_clausesContext)

	// EnterTablespace_logging_clauses is called when entering the tablespace_logging_clauses production.
	EnterTablespace_logging_clauses(c *Tablespace_logging_clausesContext)

	// EnterTablespace_group_clause is called when entering the tablespace_group_clause production.
	EnterTablespace_group_clause(c *Tablespace_group_clauseContext)

	// EnterTablespace_group_name is called when entering the tablespace_group_name production.
	EnterTablespace_group_name(c *Tablespace_group_nameContext)

	// EnterTablespace_state_clauses is called when entering the tablespace_state_clauses production.
	EnterTablespace_state_clauses(c *Tablespace_state_clausesContext)

	// EnterFlashback_mode_clause is called when entering the flashback_mode_clause production.
	EnterFlashback_mode_clause(c *Flashback_mode_clauseContext)

	// EnterNew_tablespace_name is called when entering the new_tablespace_name production.
	EnterNew_tablespace_name(c *New_tablespace_nameContext)

	// EnterCreate_tablespace is called when entering the create_tablespace production.
	EnterCreate_tablespace(c *Create_tablespaceContext)

	// EnterPermanent_tablespace_clause is called when entering the permanent_tablespace_clause production.
	EnterPermanent_tablespace_clause(c *Permanent_tablespace_clauseContext)

	// EnterTablespace_encryption_spec is called when entering the tablespace_encryption_spec production.
	EnterTablespace_encryption_spec(c *Tablespace_encryption_specContext)

	// EnterLogging_clause is called when entering the logging_clause production.
	EnterLogging_clause(c *Logging_clauseContext)

	// EnterExtent_management_clause is called when entering the extent_management_clause production.
	EnterExtent_management_clause(c *Extent_management_clauseContext)

	// EnterSegment_management_clause is called when entering the segment_management_clause production.
	EnterSegment_management_clause(c *Segment_management_clauseContext)

	// EnterTemporary_tablespace_clause is called when entering the temporary_tablespace_clause production.
	EnterTemporary_tablespace_clause(c *Temporary_tablespace_clauseContext)

	// EnterUndo_tablespace_clause is called when entering the undo_tablespace_clause production.
	EnterUndo_tablespace_clause(c *Undo_tablespace_clauseContext)

	// EnterTablespace_retention_clause is called when entering the tablespace_retention_clause production.
	EnterTablespace_retention_clause(c *Tablespace_retention_clauseContext)

	// EnterCreate_tablespace_set is called when entering the create_tablespace_set production.
	EnterCreate_tablespace_set(c *Create_tablespace_setContext)

	// EnterPermanent_tablespace_attrs is called when entering the permanent_tablespace_attrs production.
	EnterPermanent_tablespace_attrs(c *Permanent_tablespace_attrsContext)

	// EnterTablespace_encryption_clause is called when entering the tablespace_encryption_clause production.
	EnterTablespace_encryption_clause(c *Tablespace_encryption_clauseContext)

	// EnterDefault_tablespace_params is called when entering the default_tablespace_params production.
	EnterDefault_tablespace_params(c *Default_tablespace_paramsContext)

	// EnterDefault_table_compression is called when entering the default_table_compression production.
	EnterDefault_table_compression(c *Default_table_compressionContext)

	// EnterLow_high is called when entering the low_high production.
	EnterLow_high(c *Low_highContext)

	// EnterDefault_index_compression is called when entering the default_index_compression production.
	EnterDefault_index_compression(c *Default_index_compressionContext)

	// EnterInmmemory_clause is called when entering the inmmemory_clause production.
	EnterInmmemory_clause(c *Inmmemory_clauseContext)

	// EnterDatafile_specification is called when entering the datafile_specification production.
	EnterDatafile_specification(c *Datafile_specificationContext)

	// EnterTempfile_specification is called when entering the tempfile_specification production.
	EnterTempfile_specification(c *Tempfile_specificationContext)

	// EnterDatafile_tempfile_spec is called when entering the datafile_tempfile_spec production.
	EnterDatafile_tempfile_spec(c *Datafile_tempfile_specContext)

	// EnterRedo_log_file_spec is called when entering the redo_log_file_spec production.
	EnterRedo_log_file_spec(c *Redo_log_file_specContext)

	// EnterAutoextend_clause is called when entering the autoextend_clause production.
	EnterAutoextend_clause(c *Autoextend_clauseContext)

	// EnterMaxsize_clause is called when entering the maxsize_clause production.
	EnterMaxsize_clause(c *Maxsize_clauseContext)

	// EnterBuild_clause is called when entering the build_clause production.
	EnterBuild_clause(c *Build_clauseContext)

	// EnterParallel_clause is called when entering the parallel_clause production.
	EnterParallel_clause(c *Parallel_clauseContext)

	// EnterParallel_instances_clause is called when entering the parallel_instances_clause production.
	EnterParallel_instances_clause(c *Parallel_instances_clauseContext)

	// EnterAlter_materialized_view is called when entering the alter_materialized_view production.
	EnterAlter_materialized_view(c *Alter_materialized_viewContext)

	// EnterAlter_mv_option1 is called when entering the alter_mv_option1 production.
	EnterAlter_mv_option1(c *Alter_mv_option1Context)

	// EnterAlter_mv_refresh is called when entering the alter_mv_refresh production.
	EnterAlter_mv_refresh(c *Alter_mv_refreshContext)

	// EnterRollback_segment is called when entering the rollback_segment production.
	EnterRollback_segment(c *Rollback_segmentContext)

	// EnterModify_mv_column_clause is called when entering the modify_mv_column_clause production.
	EnterModify_mv_column_clause(c *Modify_mv_column_clauseContext)

	// EnterAlter_materialized_view_log is called when entering the alter_materialized_view_log production.
	EnterAlter_materialized_view_log(c *Alter_materialized_view_logContext)

	// EnterAdd_mv_log_column_clause is called when entering the add_mv_log_column_clause production.
	EnterAdd_mv_log_column_clause(c *Add_mv_log_column_clauseContext)

	// EnterMove_mv_log_clause is called when entering the move_mv_log_clause production.
	EnterMove_mv_log_clause(c *Move_mv_log_clauseContext)

	// EnterMv_log_augmentation is called when entering the mv_log_augmentation production.
	EnterMv_log_augmentation(c *Mv_log_augmentationContext)

	// EnterCreate_materialized_view_log is called when entering the create_materialized_view_log production.
	EnterCreate_materialized_view_log(c *Create_materialized_view_logContext)

	// EnterNew_values_clause is called when entering the new_values_clause production.
	EnterNew_values_clause(c *New_values_clauseContext)

	// EnterMv_log_purge_clause is called when entering the mv_log_purge_clause production.
	EnterMv_log_purge_clause(c *Mv_log_purge_clauseContext)

	// EnterCreate_materialized_zonemap is called when entering the create_materialized_zonemap production.
	EnterCreate_materialized_zonemap(c *Create_materialized_zonemapContext)

	// EnterAlter_materialized_zonemap is called when entering the alter_materialized_zonemap production.
	EnterAlter_materialized_zonemap(c *Alter_materialized_zonemapContext)

	// EnterDrop_materialized_zonemap is called when entering the drop_materialized_zonemap production.
	EnterDrop_materialized_zonemap(c *Drop_materialized_zonemapContext)

	// EnterZonemap_refresh_clause is called when entering the zonemap_refresh_clause production.
	EnterZonemap_refresh_clause(c *Zonemap_refresh_clauseContext)

	// EnterZonemap_attributes is called when entering the zonemap_attributes production.
	EnterZonemap_attributes(c *Zonemap_attributesContext)

	// EnterZonemap_name is called when entering the zonemap_name production.
	EnterZonemap_name(c *Zonemap_nameContext)

	// EnterOperator_name is called when entering the operator_name production.
	EnterOperator_name(c *Operator_nameContext)

	// EnterOperator_function_name is called when entering the operator_function_name production.
	EnterOperator_function_name(c *Operator_function_nameContext)

	// EnterCreate_zonemap_on_table is called when entering the create_zonemap_on_table production.
	EnterCreate_zonemap_on_table(c *Create_zonemap_on_tableContext)

	// EnterCreate_zonemap_as_subquery is called when entering the create_zonemap_as_subquery production.
	EnterCreate_zonemap_as_subquery(c *Create_zonemap_as_subqueryContext)

	// EnterAlter_operator is called when entering the alter_operator production.
	EnterAlter_operator(c *Alter_operatorContext)

	// EnterDrop_operator is called when entering the drop_operator production.
	EnterDrop_operator(c *Drop_operatorContext)

	// EnterCreate_operator is called when entering the create_operator production.
	EnterCreate_operator(c *Create_operatorContext)

	// EnterBinding_clause is called when entering the binding_clause production.
	EnterBinding_clause(c *Binding_clauseContext)

	// EnterAdd_binding_clause is called when entering the add_binding_clause production.
	EnterAdd_binding_clause(c *Add_binding_clauseContext)

	// EnterImplementation_clause is called when entering the implementation_clause production.
	EnterImplementation_clause(c *Implementation_clauseContext)

	// EnterPrimary_operator_list is called when entering the primary_operator_list production.
	EnterPrimary_operator_list(c *Primary_operator_listContext)

	// EnterPrimary_operator_item is called when entering the primary_operator_item production.
	EnterPrimary_operator_item(c *Primary_operator_itemContext)

	// EnterOperator_context_clause is called when entering the operator_context_clause production.
	EnterOperator_context_clause(c *Operator_context_clauseContext)

	// EnterUsing_function_clause is called when entering the using_function_clause production.
	EnterUsing_function_clause(c *Using_function_clauseContext)

	// EnterDrop_binding_clause is called when entering the drop_binding_clause production.
	EnterDrop_binding_clause(c *Drop_binding_clauseContext)

	// EnterCreate_materialized_view is called when entering the create_materialized_view production.
	EnterCreate_materialized_view(c *Create_materialized_viewContext)

	// EnterScoped_table_ref_constraint is called when entering the scoped_table_ref_constraint production.
	EnterScoped_table_ref_constraint(c *Scoped_table_ref_constraintContext)

	// EnterMv_column_alias is called when entering the mv_column_alias production.
	EnterMv_column_alias(c *Mv_column_aliasContext)

	// EnterCreate_mv_refresh is called when entering the create_mv_refresh production.
	EnterCreate_mv_refresh(c *Create_mv_refreshContext)

	// EnterQuery_rewrite_clause is called when entering the query_rewrite_clause production.
	EnterQuery_rewrite_clause(c *Query_rewrite_clauseContext)

	// EnterUnusable_editions_clause is called when entering the unusable_editions_clause production.
	EnterUnusable_editions_clause(c *Unusable_editions_clauseContext)

	// EnterDrop_materialized_view is called when entering the drop_materialized_view production.
	EnterDrop_materialized_view(c *Drop_materialized_viewContext)

	// EnterDrop_materialized_view_log is called when entering the drop_materialized_view_log production.
	EnterDrop_materialized_view_log(c *Drop_materialized_view_logContext)

	// EnterCreate_context is called when entering the create_context production.
	EnterCreate_context(c *Create_contextContext)

	// EnterOracle_namespace is called when entering the oracle_namespace production.
	EnterOracle_namespace(c *Oracle_namespaceContext)

	// EnterCreate_cluster is called when entering the create_cluster production.
	EnterCreate_cluster(c *Create_clusterContext)

	// EnterCreate_profile is called when entering the create_profile production.
	EnterCreate_profile(c *Create_profileContext)

	// EnterResource_parameters is called when entering the resource_parameters production.
	EnterResource_parameters(c *Resource_parametersContext)

	// EnterPassword_parameters is called when entering the password_parameters production.
	EnterPassword_parameters(c *Password_parametersContext)

	// EnterCreate_lockdown_profile is called when entering the create_lockdown_profile production.
	EnterCreate_lockdown_profile(c *Create_lockdown_profileContext)

	// EnterStatic_base_profile is called when entering the static_base_profile production.
	EnterStatic_base_profile(c *Static_base_profileContext)

	// EnterDynamic_base_profile is called when entering the dynamic_base_profile production.
	EnterDynamic_base_profile(c *Dynamic_base_profileContext)

	// EnterCreate_outline is called when entering the create_outline production.
	EnterCreate_outline(c *Create_outlineContext)

	// EnterCreate_restore_point is called when entering the create_restore_point production.
	EnterCreate_restore_point(c *Create_restore_pointContext)

	// EnterCreate_role is called when entering the create_role production.
	EnterCreate_role(c *Create_roleContext)

	// EnterCreate_table is called when entering the create_table production.
	EnterCreate_table(c *Create_tableContext)

	// EnterXmltype_table is called when entering the xmltype_table production.
	EnterXmltype_table(c *Xmltype_tableContext)

	// EnterXmltype_virtual_columns is called when entering the xmltype_virtual_columns production.
	EnterXmltype_virtual_columns(c *Xmltype_virtual_columnsContext)

	// EnterXmltype_column_properties is called when entering the xmltype_column_properties production.
	EnterXmltype_column_properties(c *Xmltype_column_propertiesContext)

	// EnterXmltype_storage is called when entering the xmltype_storage production.
	EnterXmltype_storage(c *Xmltype_storageContext)

	// EnterXmlschema_spec is called when entering the xmlschema_spec production.
	EnterXmlschema_spec(c *Xmlschema_specContext)

	// EnterObject_table is called when entering the object_table production.
	EnterObject_table(c *Object_tableContext)

	// EnterObject_type is called when entering the object_type production.
	EnterObject_type(c *Object_typeContext)

	// EnterOid_index_clause is called when entering the oid_index_clause production.
	EnterOid_index_clause(c *Oid_index_clauseContext)

	// EnterOid_clause is called when entering the oid_clause production.
	EnterOid_clause(c *Oid_clauseContext)

	// EnterObject_properties is called when entering the object_properties production.
	EnterObject_properties(c *Object_propertiesContext)

	// EnterObject_table_substitution is called when entering the object_table_substitution production.
	EnterObject_table_substitution(c *Object_table_substitutionContext)

	// EnterRelational_table is called when entering the relational_table production.
	EnterRelational_table(c *Relational_tableContext)

	// EnterRelational_table_properties is called when entering the relational_table_properties production.
	EnterRelational_table_properties(c *Relational_table_propertiesContext)

	// EnterRelational_table_property is called when entering the relational_table_property production.
	EnterRelational_table_property(c *Relational_table_propertyContext)

	// EnterImmutable_table_clauses is called when entering the immutable_table_clauses production.
	EnterImmutable_table_clauses(c *Immutable_table_clausesContext)

	// EnterImmutable_table_no_drop_clause is called when entering the immutable_table_no_drop_clause production.
	EnterImmutable_table_no_drop_clause(c *Immutable_table_no_drop_clauseContext)

	// EnterImmutable_table_no_delete_clause is called when entering the immutable_table_no_delete_clause production.
	EnterImmutable_table_no_delete_clause(c *Immutable_table_no_delete_clauseContext)

	// EnterBlockchain_table_clauses is called when entering the blockchain_table_clauses production.
	EnterBlockchain_table_clauses(c *Blockchain_table_clausesContext)

	// EnterBlockchain_drop_table_clause is called when entering the blockchain_drop_table_clause production.
	EnterBlockchain_drop_table_clause(c *Blockchain_drop_table_clauseContext)

	// EnterBlockchain_row_retention_clause is called when entering the blockchain_row_retention_clause production.
	EnterBlockchain_row_retention_clause(c *Blockchain_row_retention_clauseContext)

	// EnterBlockchain_hash_and_data_format_clause is called when entering the blockchain_hash_and_data_format_clause production.
	EnterBlockchain_hash_and_data_format_clause(c *Blockchain_hash_and_data_format_clauseContext)

	// EnterCollation_name is called when entering the collation_name production.
	EnterCollation_name(c *Collation_nameContext)

	// EnterTable_properties is called when entering the table_properties production.
	EnterTable_properties(c *Table_propertiesContext)

	// EnterRead_only_clause is called when entering the read_only_clause production.
	EnterRead_only_clause(c *Read_only_clauseContext)

	// EnterIndexing_clause is called when entering the indexing_clause production.
	EnterIndexing_clause(c *Indexing_clauseContext)

	// EnterAttribute_clustering_clause is called when entering the attribute_clustering_clause production.
	EnterAttribute_clustering_clause(c *Attribute_clustering_clauseContext)

	// EnterClustering_join is called when entering the clustering_join production.
	EnterClustering_join(c *Clustering_joinContext)

	// EnterClustering_join_item is called when entering the clustering_join_item production.
	EnterClustering_join_item(c *Clustering_join_itemContext)

	// EnterEquijoin_condition is called when entering the equijoin_condition production.
	EnterEquijoin_condition(c *Equijoin_conditionContext)

	// EnterCluster_clause is called when entering the cluster_clause production.
	EnterCluster_clause(c *Cluster_clauseContext)

	// EnterClustering_columns is called when entering the clustering_columns production.
	EnterClustering_columns(c *Clustering_columnsContext)

	// EnterClustering_column_group is called when entering the clustering_column_group production.
	EnterClustering_column_group(c *Clustering_column_groupContext)

	// EnterYes_no is called when entering the yes_no production.
	EnterYes_no(c *Yes_noContext)

	// EnterZonemap_clause is called when entering the zonemap_clause production.
	EnterZonemap_clause(c *Zonemap_clauseContext)

	// EnterLogical_replication_clause is called when entering the logical_replication_clause production.
	EnterLogical_replication_clause(c *Logical_replication_clauseContext)

	// EnterTable_name is called when entering the table_name production.
	EnterTable_name(c *Table_nameContext)

	// EnterRelational_property is called when entering the relational_property production.
	EnterRelational_property(c *Relational_propertyContext)

	// EnterTable_partitioning_clauses is called when entering the table_partitioning_clauses production.
	EnterTable_partitioning_clauses(c *Table_partitioning_clausesContext)

	// EnterRange_partitions is called when entering the range_partitions production.
	EnterRange_partitions(c *Range_partitionsContext)

	// EnterList_partitions is called when entering the list_partitions production.
	EnterList_partitions(c *List_partitionsContext)

	// EnterHash_partitions is called when entering the hash_partitions production.
	EnterHash_partitions(c *Hash_partitionsContext)

	// EnterIndividual_hash_partitions is called when entering the individual_hash_partitions production.
	EnterIndividual_hash_partitions(c *Individual_hash_partitionsContext)

	// EnterHash_partitions_by_quantity is called when entering the hash_partitions_by_quantity production.
	EnterHash_partitions_by_quantity(c *Hash_partitions_by_quantityContext)

	// EnterHash_partition_quantity is called when entering the hash_partition_quantity production.
	EnterHash_partition_quantity(c *Hash_partition_quantityContext)

	// EnterComposite_range_partitions is called when entering the composite_range_partitions production.
	EnterComposite_range_partitions(c *Composite_range_partitionsContext)

	// EnterComposite_list_partitions is called when entering the composite_list_partitions production.
	EnterComposite_list_partitions(c *Composite_list_partitionsContext)

	// EnterComposite_hash_partitions is called when entering the composite_hash_partitions production.
	EnterComposite_hash_partitions(c *Composite_hash_partitionsContext)

	// EnterReference_partitioning is called when entering the reference_partitioning production.
	EnterReference_partitioning(c *Reference_partitioningContext)

	// EnterReference_partition_desc is called when entering the reference_partition_desc production.
	EnterReference_partition_desc(c *Reference_partition_descContext)

	// EnterSystem_partitioning is called when entering the system_partitioning production.
	EnterSystem_partitioning(c *System_partitioningContext)

	// EnterRange_partition_desc is called when entering the range_partition_desc production.
	EnterRange_partition_desc(c *Range_partition_descContext)

	// EnterList_partition_desc is called when entering the list_partition_desc production.
	EnterList_partition_desc(c *List_partition_descContext)

	// EnterSubpartition_template is called when entering the subpartition_template production.
	EnterSubpartition_template(c *Subpartition_templateContext)

	// EnterHash_subpartition_quantity is called when entering the hash_subpartition_quantity production.
	EnterHash_subpartition_quantity(c *Hash_subpartition_quantityContext)

	// EnterSubpartition_by_range is called when entering the subpartition_by_range production.
	EnterSubpartition_by_range(c *Subpartition_by_rangeContext)

	// EnterSubpartition_by_list is called when entering the subpartition_by_list production.
	EnterSubpartition_by_list(c *Subpartition_by_listContext)

	// EnterSubpartition_by_hash is called when entering the subpartition_by_hash production.
	EnterSubpartition_by_hash(c *Subpartition_by_hashContext)

	// EnterSubpartition_name is called when entering the subpartition_name production.
	EnterSubpartition_name(c *Subpartition_nameContext)

	// EnterRange_subpartition_desc is called when entering the range_subpartition_desc production.
	EnterRange_subpartition_desc(c *Range_subpartition_descContext)

	// EnterList_subpartition_desc is called when entering the list_subpartition_desc production.
	EnterList_subpartition_desc(c *List_subpartition_descContext)

	// EnterIndividual_hash_subparts is called when entering the individual_hash_subparts production.
	EnterIndividual_hash_subparts(c *Individual_hash_subpartsContext)

	// EnterHash_subparts_by_quantity is called when entering the hash_subparts_by_quantity production.
	EnterHash_subparts_by_quantity(c *Hash_subparts_by_quantityContext)

	// EnterRange_values_clause is called when entering the range_values_clause production.
	EnterRange_values_clause(c *Range_values_clauseContext)

	// EnterRange_values_list is called when entering the range_values_list production.
	EnterRange_values_list(c *Range_values_listContext)

	// EnterList_values_clause is called when entering the list_values_clause production.
	EnterList_values_clause(c *List_values_clauseContext)

	// EnterTable_partition_description is called when entering the table_partition_description production.
	EnterTable_partition_description(c *Table_partition_descriptionContext)

	// EnterPartitioning_storage_clause is called when entering the partitioning_storage_clause production.
	EnterPartitioning_storage_clause(c *Partitioning_storage_clauseContext)

	// EnterLob_partitioning_storage is called when entering the lob_partitioning_storage production.
	EnterLob_partitioning_storage(c *Lob_partitioning_storageContext)

	// EnterSize_clause is called when entering the size_clause production.
	EnterSize_clause(c *Size_clauseContext)

	// EnterTable_compression is called when entering the table_compression production.
	EnterTable_compression(c *Table_compressionContext)

	// EnterInmemory_table_clause is called when entering the inmemory_table_clause production.
	EnterInmemory_table_clause(c *Inmemory_table_clauseContext)

	// EnterInmemory_attributes is called when entering the inmemory_attributes production.
	EnterInmemory_attributes(c *Inmemory_attributesContext)

	// EnterInmemory_memcompress is called when entering the inmemory_memcompress production.
	EnterInmemory_memcompress(c *Inmemory_memcompressContext)

	// EnterInmemory_priority is called when entering the inmemory_priority production.
	EnterInmemory_priority(c *Inmemory_priorityContext)

	// EnterInmemory_distribute is called when entering the inmemory_distribute production.
	EnterInmemory_distribute(c *Inmemory_distributeContext)

	// EnterInmemory_duplicate is called when entering the inmemory_duplicate production.
	EnterInmemory_duplicate(c *Inmemory_duplicateContext)

	// EnterInmemory_column_clause is called when entering the inmemory_column_clause production.
	EnterInmemory_column_clause(c *Inmemory_column_clauseContext)

	// EnterPhysical_attributes_clause is called when entering the physical_attributes_clause production.
	EnterPhysical_attributes_clause(c *Physical_attributes_clauseContext)

	// EnterStorage_clause is called when entering the storage_clause production.
	EnterStorage_clause(c *Storage_clauseContext)

	// EnterDeferred_segment_creation is called when entering the deferred_segment_creation production.
	EnterDeferred_segment_creation(c *Deferred_segment_creationContext)

	// EnterSegment_attributes_clause is called when entering the segment_attributes_clause production.
	EnterSegment_attributes_clause(c *Segment_attributes_clauseContext)

	// EnterPhysical_properties is called when entering the physical_properties production.
	EnterPhysical_properties(c *Physical_propertiesContext)

	// EnterIlm_clause is called when entering the ilm_clause production.
	EnterIlm_clause(c *Ilm_clauseContext)

	// EnterIlm_policy_clause is called when entering the ilm_policy_clause production.
	EnterIlm_policy_clause(c *Ilm_policy_clauseContext)

	// EnterIlm_compression_policy is called when entering the ilm_compression_policy production.
	EnterIlm_compression_policy(c *Ilm_compression_policyContext)

	// EnterIlm_tiering_policy is called when entering the ilm_tiering_policy production.
	EnterIlm_tiering_policy(c *Ilm_tiering_policyContext)

	// EnterIlm_after_on is called when entering the ilm_after_on production.
	EnterIlm_after_on(c *Ilm_after_onContext)

	// EnterSegment_group is called when entering the segment_group production.
	EnterSegment_group(c *Segment_groupContext)

	// EnterIlm_inmemory_policy is called when entering the ilm_inmemory_policy production.
	EnterIlm_inmemory_policy(c *Ilm_inmemory_policyContext)

	// EnterIlm_time_period is called when entering the ilm_time_period production.
	EnterIlm_time_period(c *Ilm_time_periodContext)

	// EnterHeap_org_table_clause is called when entering the heap_org_table_clause production.
	EnterHeap_org_table_clause(c *Heap_org_table_clauseContext)

	// EnterExternal_table_clause is called when entering the external_table_clause production.
	EnterExternal_table_clause(c *External_table_clauseContext)

	// EnterAccess_driver_type is called when entering the access_driver_type production.
	EnterAccess_driver_type(c *Access_driver_typeContext)

	// EnterExternal_table_data_props is called when entering the external_table_data_props production.
	EnterExternal_table_data_props(c *External_table_data_propsContext)

	// EnterExternal_table_data_format is called when entering the external_table_data_format production.
	EnterExternal_table_data_format(c *External_table_data_formatContext)

	// EnterExternal_table_transform is called when entering the external_table_transform production.
	EnterExternal_table_transform(c *External_table_transformContext)

	// EnterExternal_table_field is called when entering the external_table_field production.
	EnterExternal_table_field(c *External_table_fieldContext)

	// EnterExternal_table_field_list is called when entering the external_table_field_list production.
	EnterExternal_table_field_list(c *External_table_field_listContext)

	// EnterExternal_table_fields_clause is called when entering the external_table_fields_clause production.
	EnterExternal_table_fields_clause(c *External_table_fields_clauseContext)

	// EnterExternal_table_position_clause is called when entering the external_table_position_clause production.
	EnterExternal_table_position_clause(c *External_table_position_clauseContext)

	// EnterExternal_table_datatype_clause is called when entering the external_table_datatype_clause production.
	EnterExternal_table_datatype_clause(c *External_table_datatype_clauseContext)

	// EnterExternal_table_delimit_clause is called when entering the external_table_delimit_clause production.
	EnterExternal_table_delimit_clause(c *External_table_delimit_clauseContext)

	// EnterExternal_table_trim_clause is called when entering the external_table_trim_clause production.
	EnterExternal_table_trim_clause(c *External_table_trim_clauseContext)

	// EnterExternal_table_date_format_clause is called when entering the external_table_date_format_clause production.
	EnterExternal_table_date_format_clause(c *External_table_date_format_clauseContext)

	// EnterExternal_table_init_clause is called when entering the external_table_init_clause production.
	EnterExternal_table_init_clause(c *External_table_init_clauseContext)

	// EnterExternal_table_condition_clause is called when entering the external_table_condition_clause production.
	EnterExternal_table_condition_clause(c *External_table_condition_clauseContext)

	// EnterExternal_table_lls_clause is called when entering the external_table_lls_clause production.
	EnterExternal_table_lls_clause(c *External_table_lls_clauseContext)

	// EnterExternal_table_records is called when entering the external_table_records production.
	EnterExternal_table_records(c *External_table_recordsContext)

	// EnterExternal_table_record_options_clause is called when entering the external_table_record_options_clause production.
	EnterExternal_table_record_options_clause(c *External_table_record_options_clauseContext)

	// EnterExternal_table_output_files is called when entering the external_table_output_files production.
	EnterExternal_table_output_files(c *External_table_output_filesContext)

	// EnterExternal_table_fields is called when entering the external_table_fields production.
	EnterExternal_table_fields(c *External_table_fieldsContext)

	// EnterExternal_table_datapump is called when entering the external_table_datapump production.
	EnterExternal_table_datapump(c *External_table_datapumpContext)

	// EnterExternal_table_hive is called when entering the external_table_hive production.
	EnterExternal_table_hive(c *External_table_hiveContext)

	// EnterExternal_table_hive_parameter_map is called when entering the external_table_hive_parameter_map production.
	EnterExternal_table_hive_parameter_map(c *External_table_hive_parameter_mapContext)

	// EnterExternal_table_hive_parameter_map_entry is called when entering the external_table_hive_parameter_map_entry production.
	EnterExternal_table_hive_parameter_map_entry(c *External_table_hive_parameter_map_entryContext)

	// EnterExternal_table_directory is called when entering the external_table_directory production.
	EnterExternal_table_directory(c *External_table_directoryContext)

	// EnterRow_movement_clause is called when entering the row_movement_clause production.
	EnterRow_movement_clause(c *Row_movement_clauseContext)

	// EnterFlashback_archive_clause is called when entering the flashback_archive_clause production.
	EnterFlashback_archive_clause(c *Flashback_archive_clauseContext)

	// EnterLog_grp is called when entering the log_grp production.
	EnterLog_grp(c *Log_grpContext)

	// EnterSupplemental_table_logging is called when entering the supplemental_table_logging production.
	EnterSupplemental_table_logging(c *Supplemental_table_loggingContext)

	// EnterSupplemental_log_grp_clause is called when entering the supplemental_log_grp_clause production.
	EnterSupplemental_log_grp_clause(c *Supplemental_log_grp_clauseContext)

	// EnterSupplemental_id_key_clause is called when entering the supplemental_id_key_clause production.
	EnterSupplemental_id_key_clause(c *Supplemental_id_key_clauseContext)

	// EnterAllocate_extent_clause is called when entering the allocate_extent_clause production.
	EnterAllocate_extent_clause(c *Allocate_extent_clauseContext)

	// EnterDeallocate_unused_clause is called when entering the deallocate_unused_clause production.
	EnterDeallocate_unused_clause(c *Deallocate_unused_clauseContext)

	// EnterShrink_clause is called when entering the shrink_clause production.
	EnterShrink_clause(c *Shrink_clauseContext)

	// EnterRecords_per_block_clause is called when entering the records_per_block_clause production.
	EnterRecords_per_block_clause(c *Records_per_block_clauseContext)

	// EnterUpgrade_table_clause is called when entering the upgrade_table_clause production.
	EnterUpgrade_table_clause(c *Upgrade_table_clauseContext)

	// EnterTruncate_table is called when entering the truncate_table production.
	EnterTruncate_table(c *Truncate_tableContext)

	// EnterDrop_table is called when entering the drop_table production.
	EnterDrop_table(c *Drop_tableContext)

	// EnterDrop_tablespace is called when entering the drop_tablespace production.
	EnterDrop_tablespace(c *Drop_tablespaceContext)

	// EnterDrop_tablespace_set is called when entering the drop_tablespace_set production.
	EnterDrop_tablespace_set(c *Drop_tablespace_setContext)

	// EnterIncluding_contents_clause is called when entering the including_contents_clause production.
	EnterIncluding_contents_clause(c *Including_contents_clauseContext)

	// EnterDrop_view is called when entering the drop_view production.
	EnterDrop_view(c *Drop_viewContext)

	// EnterComment_on_column is called when entering the comment_on_column production.
	EnterComment_on_column(c *Comment_on_columnContext)

	// EnterEnable_or_disable is called when entering the enable_or_disable production.
	EnterEnable_or_disable(c *Enable_or_disableContext)

	// EnterAllow_or_disallow is called when entering the allow_or_disallow production.
	EnterAllow_or_disallow(c *Allow_or_disallowContext)

	// EnterAlter_synonym is called when entering the alter_synonym production.
	EnterAlter_synonym(c *Alter_synonymContext)

	// EnterCreate_synonym is called when entering the create_synonym production.
	EnterCreate_synonym(c *Create_synonymContext)

	// EnterDrop_synonym is called when entering the drop_synonym production.
	EnterDrop_synonym(c *Drop_synonymContext)

	// EnterCreate_spfile is called when entering the create_spfile production.
	EnterCreate_spfile(c *Create_spfileContext)

	// EnterSpfile_name is called when entering the spfile_name production.
	EnterSpfile_name(c *Spfile_nameContext)

	// EnterPfile_name is called when entering the pfile_name production.
	EnterPfile_name(c *Pfile_nameContext)

	// EnterComment_on_table is called when entering the comment_on_table production.
	EnterComment_on_table(c *Comment_on_tableContext)

	// EnterComment_on_materialized is called when entering the comment_on_materialized production.
	EnterComment_on_materialized(c *Comment_on_materializedContext)

	// EnterAlter_analytic_view is called when entering the alter_analytic_view production.
	EnterAlter_analytic_view(c *Alter_analytic_viewContext)

	// EnterAlter_add_cache_clause is called when entering the alter_add_cache_clause production.
	EnterAlter_add_cache_clause(c *Alter_add_cache_clauseContext)

	// EnterLevels_item is called when entering the levels_item production.
	EnterLevels_item(c *Levels_itemContext)

	// EnterMeasure_list is called when entering the measure_list production.
	EnterMeasure_list(c *Measure_listContext)

	// EnterAlter_drop_cache_clause is called when entering the alter_drop_cache_clause production.
	EnterAlter_drop_cache_clause(c *Alter_drop_cache_clauseContext)

	// EnterAlter_attribute_dimension is called when entering the alter_attribute_dimension production.
	EnterAlter_attribute_dimension(c *Alter_attribute_dimensionContext)

	// EnterAlter_audit_policy is called when entering the alter_audit_policy production.
	EnterAlter_audit_policy(c *Alter_audit_policyContext)

	// EnterAlter_cluster is called when entering the alter_cluster production.
	EnterAlter_cluster(c *Alter_clusterContext)

	// EnterDrop_analytic_view is called when entering the drop_analytic_view production.
	EnterDrop_analytic_view(c *Drop_analytic_viewContext)

	// EnterDrop_attribute_dimension is called when entering the drop_attribute_dimension production.
	EnterDrop_attribute_dimension(c *Drop_attribute_dimensionContext)

	// EnterDrop_audit_policy is called when entering the drop_audit_policy production.
	EnterDrop_audit_policy(c *Drop_audit_policyContext)

	// EnterDrop_flashback_archive is called when entering the drop_flashback_archive production.
	EnterDrop_flashback_archive(c *Drop_flashback_archiveContext)

	// EnterDrop_cluster is called when entering the drop_cluster production.
	EnterDrop_cluster(c *Drop_clusterContext)

	// EnterDrop_context is called when entering the drop_context production.
	EnterDrop_context(c *Drop_contextContext)

	// EnterDrop_directory is called when entering the drop_directory production.
	EnterDrop_directory(c *Drop_directoryContext)

	// EnterDrop_diskgroup is called when entering the drop_diskgroup production.
	EnterDrop_diskgroup(c *Drop_diskgroupContext)

	// EnterDrop_edition is called when entering the drop_edition production.
	EnterDrop_edition(c *Drop_editionContext)

	// EnterTruncate_cluster is called when entering the truncate_cluster production.
	EnterTruncate_cluster(c *Truncate_clusterContext)

	// EnterCache_or_nocache is called when entering the cache_or_nocache production.
	EnterCache_or_nocache(c *Cache_or_nocacheContext)

	// EnterDatabase_name is called when entering the database_name production.
	EnterDatabase_name(c *Database_nameContext)

	// EnterAlter_database is called when entering the alter_database production.
	EnterAlter_database(c *Alter_databaseContext)

	// EnterDatabase_clause is called when entering the database_clause production.
	EnterDatabase_clause(c *Database_clauseContext)

	// EnterStartup_clauses is called when entering the startup_clauses production.
	EnterStartup_clauses(c *Startup_clausesContext)

	// EnterResetlogs_or_noresetlogs is called when entering the resetlogs_or_noresetlogs production.
	EnterResetlogs_or_noresetlogs(c *Resetlogs_or_noresetlogsContext)

	// EnterUpgrade_or_downgrade is called when entering the upgrade_or_downgrade production.
	EnterUpgrade_or_downgrade(c *Upgrade_or_downgradeContext)

	// EnterRecovery_clauses is called when entering the recovery_clauses production.
	EnterRecovery_clauses(c *Recovery_clausesContext)

	// EnterBegin_or_end is called when entering the begin_or_end production.
	EnterBegin_or_end(c *Begin_or_endContext)

	// EnterGeneral_recovery is called when entering the general_recovery production.
	EnterGeneral_recovery(c *General_recoveryContext)

	// EnterFull_database_recovery is called when entering the full_database_recovery production.
	EnterFull_database_recovery(c *Full_database_recoveryContext)

	// EnterPartial_database_recovery is called when entering the partial_database_recovery production.
	EnterPartial_database_recovery(c *Partial_database_recoveryContext)

	// EnterPartial_database_recovery_10g is called when entering the partial_database_recovery_10g production.
	EnterPartial_database_recovery_10g(c *Partial_database_recovery_10gContext)

	// EnterManaged_standby_recovery is called when entering the managed_standby_recovery production.
	EnterManaged_standby_recovery(c *Managed_standby_recoveryContext)

	// EnterDb_name is called when entering the db_name production.
	EnterDb_name(c *Db_nameContext)

	// EnterDatabase_file_clauses is called when entering the database_file_clauses production.
	EnterDatabase_file_clauses(c *Database_file_clausesContext)

	// EnterCreate_datafile_clause is called when entering the create_datafile_clause production.
	EnterCreate_datafile_clause(c *Create_datafile_clauseContext)

	// EnterAlter_datafile_clause is called when entering the alter_datafile_clause production.
	EnterAlter_datafile_clause(c *Alter_datafile_clauseContext)

	// EnterAlter_tempfile_clause is called when entering the alter_tempfile_clause production.
	EnterAlter_tempfile_clause(c *Alter_tempfile_clauseContext)

	// EnterMove_datafile_clause is called when entering the move_datafile_clause production.
	EnterMove_datafile_clause(c *Move_datafile_clauseContext)

	// EnterLogfile_clauses is called when entering the logfile_clauses production.
	EnterLogfile_clauses(c *Logfile_clausesContext)

	// EnterAdd_logfile_clauses is called when entering the add_logfile_clauses production.
	EnterAdd_logfile_clauses(c *Add_logfile_clausesContext)

	// EnterGroup_redo_logfile is called when entering the group_redo_logfile production.
	EnterGroup_redo_logfile(c *Group_redo_logfileContext)

	// EnterDrop_logfile_clauses is called when entering the drop_logfile_clauses production.
	EnterDrop_logfile_clauses(c *Drop_logfile_clausesContext)

	// EnterSwitch_logfile_clause is called when entering the switch_logfile_clause production.
	EnterSwitch_logfile_clause(c *Switch_logfile_clauseContext)

	// EnterSupplemental_db_logging is called when entering the supplemental_db_logging production.
	EnterSupplemental_db_logging(c *Supplemental_db_loggingContext)

	// EnterAdd_or_drop is called when entering the add_or_drop production.
	EnterAdd_or_drop(c *Add_or_dropContext)

	// EnterSupplemental_plsql_clause is called when entering the supplemental_plsql_clause production.
	EnterSupplemental_plsql_clause(c *Supplemental_plsql_clauseContext)

	// EnterLogfile_descriptor is called when entering the logfile_descriptor production.
	EnterLogfile_descriptor(c *Logfile_descriptorContext)

	// EnterControlfile_clauses is called when entering the controlfile_clauses production.
	EnterControlfile_clauses(c *Controlfile_clausesContext)

	// EnterTrace_file_clause is called when entering the trace_file_clause production.
	EnterTrace_file_clause(c *Trace_file_clauseContext)

	// EnterStandby_database_clauses is called when entering the standby_database_clauses production.
	EnterStandby_database_clauses(c *Standby_database_clausesContext)

	// EnterActivate_standby_db_clause is called when entering the activate_standby_db_clause production.
	EnterActivate_standby_db_clause(c *Activate_standby_db_clauseContext)

	// EnterMaximize_standby_db_clause is called when entering the maximize_standby_db_clause production.
	EnterMaximize_standby_db_clause(c *Maximize_standby_db_clauseContext)

	// EnterRegister_logfile_clause is called when entering the register_logfile_clause production.
	EnterRegister_logfile_clause(c *Register_logfile_clauseContext)

	// EnterCommit_switchover_clause is called when entering the commit_switchover_clause production.
	EnterCommit_switchover_clause(c *Commit_switchover_clauseContext)

	// EnterStart_standby_clause is called when entering the start_standby_clause production.
	EnterStart_standby_clause(c *Start_standby_clauseContext)

	// EnterStop_standby_clause is called when entering the stop_standby_clause production.
	EnterStop_standby_clause(c *Stop_standby_clauseContext)

	// EnterConvert_database_clause is called when entering the convert_database_clause production.
	EnterConvert_database_clause(c *Convert_database_clauseContext)

	// EnterDefault_settings_clause is called when entering the default_settings_clause production.
	EnterDefault_settings_clause(c *Default_settings_clauseContext)

	// EnterSet_time_zone_clause is called when entering the set_time_zone_clause production.
	EnterSet_time_zone_clause(c *Set_time_zone_clauseContext)

	// EnterInstance_clauses is called when entering the instance_clauses production.
	EnterInstance_clauses(c *Instance_clausesContext)

	// EnterSecurity_clause is called when entering the security_clause production.
	EnterSecurity_clause(c *Security_clauseContext)

	// EnterDomain is called when entering the domain production.
	EnterDomain(c *DomainContext)

	// EnterDatabase is called when entering the database production.
	EnterDatabase(c *DatabaseContext)

	// EnterEdition_name is called when entering the edition_name production.
	EnterEdition_name(c *Edition_nameContext)

	// EnterFilenumber is called when entering the filenumber production.
	EnterFilenumber(c *FilenumberContext)

	// EnterFilename is called when entering the filename production.
	EnterFilename(c *FilenameContext)

	// EnterPrepare_clause is called when entering the prepare_clause production.
	EnterPrepare_clause(c *Prepare_clauseContext)

	// EnterDrop_mirror_clause is called when entering the drop_mirror_clause production.
	EnterDrop_mirror_clause(c *Drop_mirror_clauseContext)

	// EnterLost_write_protection is called when entering the lost_write_protection production.
	EnterLost_write_protection(c *Lost_write_protectionContext)

	// EnterCdb_fleet_clauses is called when entering the cdb_fleet_clauses production.
	EnterCdb_fleet_clauses(c *Cdb_fleet_clausesContext)

	// EnterLead_cdb_clause is called when entering the lead_cdb_clause production.
	EnterLead_cdb_clause(c *Lead_cdb_clauseContext)

	// EnterLead_cdb_uri_clause is called when entering the lead_cdb_uri_clause production.
	EnterLead_cdb_uri_clause(c *Lead_cdb_uri_clauseContext)

	// EnterProperty_clauses is called when entering the property_clauses production.
	EnterProperty_clauses(c *Property_clausesContext)

	// EnterReplay_upgrade_clauses is called when entering the replay_upgrade_clauses production.
	EnterReplay_upgrade_clauses(c *Replay_upgrade_clausesContext)

	// EnterAlter_database_link is called when entering the alter_database_link production.
	EnterAlter_database_link(c *Alter_database_linkContext)

	// EnterPassword_value is called when entering the password_value production.
	EnterPassword_value(c *Password_valueContext)

	// EnterLink_authentication is called when entering the link_authentication production.
	EnterLink_authentication(c *Link_authenticationContext)

	// EnterCreate_schema is called when entering the create_schema production.
	EnterCreate_schema(c *Create_schemaContext)

	// EnterCreate_database is called when entering the create_database production.
	EnterCreate_database(c *Create_databaseContext)

	// EnterDatabase_logging_clauses is called when entering the database_logging_clauses production.
	EnterDatabase_logging_clauses(c *Database_logging_clausesContext)

	// EnterDatabase_logging_sub_clause is called when entering the database_logging_sub_clause production.
	EnterDatabase_logging_sub_clause(c *Database_logging_sub_clauseContext)

	// EnterTablespace_clauses is called when entering the tablespace_clauses production.
	EnterTablespace_clauses(c *Tablespace_clausesContext)

	// EnterEnable_pluggable_database is called when entering the enable_pluggable_database production.
	EnterEnable_pluggable_database(c *Enable_pluggable_databaseContext)

	// EnterFile_name_convert is called when entering the file_name_convert production.
	EnterFile_name_convert(c *File_name_convertContext)

	// EnterFilename_convert_sub_clause is called when entering the filename_convert_sub_clause production.
	EnterFilename_convert_sub_clause(c *Filename_convert_sub_clauseContext)

	// EnterTablespace_datafile_clauses is called when entering the tablespace_datafile_clauses production.
	EnterTablespace_datafile_clauses(c *Tablespace_datafile_clausesContext)

	// EnterUndo_mode_clause is called when entering the undo_mode_clause production.
	EnterUndo_mode_clause(c *Undo_mode_clauseContext)

	// EnterDefault_tablespace is called when entering the default_tablespace production.
	EnterDefault_tablespace(c *Default_tablespaceContext)

	// EnterDefault_temp_tablespace is called when entering the default_temp_tablespace production.
	EnterDefault_temp_tablespace(c *Default_temp_tablespaceContext)

	// EnterUndo_tablespace is called when entering the undo_tablespace production.
	EnterUndo_tablespace(c *Undo_tablespaceContext)

	// EnterDrop_database is called when entering the drop_database production.
	EnterDrop_database(c *Drop_databaseContext)

	// EnterCreate_database_link is called when entering the create_database_link production.
	EnterCreate_database_link(c *Create_database_linkContext)

	// EnterDrop_database_link is called when entering the drop_database_link production.
	EnterDrop_database_link(c *Drop_database_linkContext)

	// EnterAlter_tablespace_set is called when entering the alter_tablespace_set production.
	EnterAlter_tablespace_set(c *Alter_tablespace_setContext)

	// EnterAlter_tablespace_attrs is called when entering the alter_tablespace_attrs production.
	EnterAlter_tablespace_attrs(c *Alter_tablespace_attrsContext)

	// EnterAlter_tablespace_encryption is called when entering the alter_tablespace_encryption production.
	EnterAlter_tablespace_encryption(c *Alter_tablespace_encryptionContext)

	// EnterTs_file_name_convert is called when entering the ts_file_name_convert production.
	EnterTs_file_name_convert(c *Ts_file_name_convertContext)

	// EnterAlter_role is called when entering the alter_role production.
	EnterAlter_role(c *Alter_roleContext)

	// EnterRole_identified_clause is called when entering the role_identified_clause production.
	EnterRole_identified_clause(c *Role_identified_clauseContext)

	// EnterAlter_table is called when entering the alter_table production.
	EnterAlter_table(c *Alter_tableContext)

	// EnterMemoptimize_read_write_clause is called when entering the memoptimize_read_write_clause production.
	EnterMemoptimize_read_write_clause(c *Memoptimize_read_write_clauseContext)

	// EnterAlter_table_properties is called when entering the alter_table_properties production.
	EnterAlter_table_properties(c *Alter_table_propertiesContext)

	// EnterAlter_table_partitioning is called when entering the alter_table_partitioning production.
	EnterAlter_table_partitioning(c *Alter_table_partitioningContext)

	// EnterAdd_table_partition is called when entering the add_table_partition production.
	EnterAdd_table_partition(c *Add_table_partitionContext)

	// EnterDrop_table_partition is called when entering the drop_table_partition production.
	EnterDrop_table_partition(c *Drop_table_partitionContext)

	// EnterMerge_table_partition is called when entering the merge_table_partition production.
	EnterMerge_table_partition(c *Merge_table_partitionContext)

	// EnterModify_table_partition is called when entering the modify_table_partition production.
	EnterModify_table_partition(c *Modify_table_partitionContext)

	// EnterSplit_table_partition is called when entering the split_table_partition production.
	EnterSplit_table_partition(c *Split_table_partitionContext)

	// EnterTruncate_table_partition is called when entering the truncate_table_partition production.
	EnterTruncate_table_partition(c *Truncate_table_partitionContext)

	// EnterExchange_table_partition is called when entering the exchange_table_partition production.
	EnterExchange_table_partition(c *Exchange_table_partitionContext)

	// EnterCoalesce_table_partition is called when entering the coalesce_table_partition production.
	EnterCoalesce_table_partition(c *Coalesce_table_partitionContext)

	// EnterAlter_interval_partition is called when entering the alter_interval_partition production.
	EnterAlter_interval_partition(c *Alter_interval_partitionContext)

	// EnterMove_table_partition is called when entering the move_table_partition production.
	EnterMove_table_partition(c *Move_table_partitionContext)

	// EnterFilter_condition is called when entering the filter_condition production.
	EnterFilter_condition(c *Filter_conditionContext)

	// EnterRename_table_partition is called when entering the rename_table_partition production.
	EnterRename_table_partition(c *Rename_table_partitionContext)

	// EnterPartition_extended_names is called when entering the partition_extended_names production.
	EnterPartition_extended_names(c *Partition_extended_namesContext)

	// EnterSubpartition_extended_names is called when entering the subpartition_extended_names production.
	EnterSubpartition_extended_names(c *Subpartition_extended_namesContext)

	// EnterAlter_table_properties_1 is called when entering the alter_table_properties_1 production.
	EnterAlter_table_properties_1(c *Alter_table_properties_1Context)

	// EnterAlter_iot_clauses is called when entering the alter_iot_clauses production.
	EnterAlter_iot_clauses(c *Alter_iot_clausesContext)

	// EnterAlter_mapping_table_clause is called when entering the alter_mapping_table_clause production.
	EnterAlter_mapping_table_clause(c *Alter_mapping_table_clauseContext)

	// EnterAlter_overflow_clause is called when entering the alter_overflow_clause production.
	EnterAlter_overflow_clause(c *Alter_overflow_clauseContext)

	// EnterAdd_overflow_clause is called when entering the add_overflow_clause production.
	EnterAdd_overflow_clause(c *Add_overflow_clauseContext)

	// EnterUpdate_index_clauses is called when entering the update_index_clauses production.
	EnterUpdate_index_clauses(c *Update_index_clausesContext)

	// EnterUpdate_global_index_clause is called when entering the update_global_index_clause production.
	EnterUpdate_global_index_clause(c *Update_global_index_clauseContext)

	// EnterUpdate_all_indexes_clause is called when entering the update_all_indexes_clause production.
	EnterUpdate_all_indexes_clause(c *Update_all_indexes_clauseContext)

	// EnterUpdate_all_indexes_index_clause is called when entering the update_all_indexes_index_clause production.
	EnterUpdate_all_indexes_index_clause(c *Update_all_indexes_index_clauseContext)

	// EnterUpdate_index_partition is called when entering the update_index_partition production.
	EnterUpdate_index_partition(c *Update_index_partitionContext)

	// EnterUpdate_index_subpartition is called when entering the update_index_subpartition production.
	EnterUpdate_index_subpartition(c *Update_index_subpartitionContext)

	// EnterEnable_disable_clause is called when entering the enable_disable_clause production.
	EnterEnable_disable_clause(c *Enable_disable_clauseContext)

	// EnterUsing_index_clause is called when entering the using_index_clause production.
	EnterUsing_index_clause(c *Using_index_clauseContext)

	// EnterIndex_attributes is called when entering the index_attributes production.
	EnterIndex_attributes(c *Index_attributesContext)

	// EnterSort_or_nosort is called when entering the sort_or_nosort production.
	EnterSort_or_nosort(c *Sort_or_nosortContext)

	// EnterExceptions_clause is called when entering the exceptions_clause production.
	EnterExceptions_clause(c *Exceptions_clauseContext)

	// EnterMove_table_clause is called when entering the move_table_clause production.
	EnterMove_table_clause(c *Move_table_clauseContext)

	// EnterIndex_org_table_clause is called when entering the index_org_table_clause production.
	EnterIndex_org_table_clause(c *Index_org_table_clauseContext)

	// EnterMapping_table_clause is called when entering the mapping_table_clause production.
	EnterMapping_table_clause(c *Mapping_table_clauseContext)

	// EnterKey_compression is called when entering the key_compression production.
	EnterKey_compression(c *Key_compressionContext)

	// EnterIndex_org_overflow_clause is called when entering the index_org_overflow_clause production.
	EnterIndex_org_overflow_clause(c *Index_org_overflow_clauseContext)

	// EnterColumn_clauses is called when entering the column_clauses production.
	EnterColumn_clauses(c *Column_clausesContext)

	// EnterModify_collection_retrieval is called when entering the modify_collection_retrieval production.
	EnterModify_collection_retrieval(c *Modify_collection_retrievalContext)

	// EnterCollection_item is called when entering the collection_item production.
	EnterCollection_item(c *Collection_itemContext)

	// EnterRename_column_clause is called when entering the rename_column_clause production.
	EnterRename_column_clause(c *Rename_column_clauseContext)

	// EnterOld_column_name is called when entering the old_column_name production.
	EnterOld_column_name(c *Old_column_nameContext)

	// EnterNew_column_name is called when entering the new_column_name production.
	EnterNew_column_name(c *New_column_nameContext)

	// EnterAdd_modify_drop_column_clauses is called when entering the add_modify_drop_column_clauses production.
	EnterAdd_modify_drop_column_clauses(c *Add_modify_drop_column_clausesContext)

	// EnterDrop_column_clause is called when entering the drop_column_clause production.
	EnterDrop_column_clause(c *Drop_column_clauseContext)

	// EnterModify_column_clauses is called when entering the modify_column_clauses production.
	EnterModify_column_clauses(c *Modify_column_clausesContext)

	// EnterModify_col_properties is called when entering the modify_col_properties production.
	EnterModify_col_properties(c *Modify_col_propertiesContext)

	// EnterModify_col_visibility is called when entering the modify_col_visibility production.
	EnterModify_col_visibility(c *Modify_col_visibilityContext)

	// EnterModify_col_substitutable is called when entering the modify_col_substitutable production.
	EnterModify_col_substitutable(c *Modify_col_substitutableContext)

	// EnterAdd_column_clause is called when entering the add_column_clause production.
	EnterAdd_column_clause(c *Add_column_clauseContext)

	// EnterVarray_col_properties is called when entering the varray_col_properties production.
	EnterVarray_col_properties(c *Varray_col_propertiesContext)

	// EnterVarray_storage_clause is called when entering the varray_storage_clause production.
	EnterVarray_storage_clause(c *Varray_storage_clauseContext)

	// EnterLob_segname is called when entering the lob_segname production.
	EnterLob_segname(c *Lob_segnameContext)

	// EnterLob_item is called when entering the lob_item production.
	EnterLob_item(c *Lob_itemContext)

	// EnterLob_storage_parameters is called when entering the lob_storage_parameters production.
	EnterLob_storage_parameters(c *Lob_storage_parametersContext)

	// EnterLob_storage_clause is called when entering the lob_storage_clause production.
	EnterLob_storage_clause(c *Lob_storage_clauseContext)

	// EnterModify_lob_storage_clause is called when entering the modify_lob_storage_clause production.
	EnterModify_lob_storage_clause(c *Modify_lob_storage_clauseContext)

	// EnterModify_lob_parameters is called when entering the modify_lob_parameters production.
	EnterModify_lob_parameters(c *Modify_lob_parametersContext)

	// EnterLob_parameters is called when entering the lob_parameters production.
	EnterLob_parameters(c *Lob_parametersContext)

	// EnterLob_deduplicate_clause is called when entering the lob_deduplicate_clause production.
	EnterLob_deduplicate_clause(c *Lob_deduplicate_clauseContext)

	// EnterLob_compression_clause is called when entering the lob_compression_clause production.
	EnterLob_compression_clause(c *Lob_compression_clauseContext)

	// EnterLob_retention_clause is called when entering the lob_retention_clause production.
	EnterLob_retention_clause(c *Lob_retention_clauseContext)

	// EnterEncryption_spec is called when entering the encryption_spec production.
	EnterEncryption_spec(c *Encryption_specContext)

	// EnterTablespace is called when entering the tablespace production.
	EnterTablespace(c *TablespaceContext)

	// EnterVarray_item is called when entering the varray_item production.
	EnterVarray_item(c *Varray_itemContext)

	// EnterColumn_properties is called when entering the column_properties production.
	EnterColumn_properties(c *Column_propertiesContext)

	// EnterLob_partition_storage is called when entering the lob_partition_storage production.
	EnterLob_partition_storage(c *Lob_partition_storageContext)

	// EnterPeriod_definition is called when entering the period_definition production.
	EnterPeriod_definition(c *Period_definitionContext)

	// EnterStart_time_column is called when entering the start_time_column production.
	EnterStart_time_column(c *Start_time_columnContext)

	// EnterEnd_time_column is called when entering the end_time_column production.
	EnterEnd_time_column(c *End_time_columnContext)

	// EnterColumn_definition is called when entering the column_definition production.
	EnterColumn_definition(c *Column_definitionContext)

	// EnterColumn_collation_name is called when entering the column_collation_name production.
	EnterColumn_collation_name(c *Column_collation_nameContext)

	// EnterIdentity_clause is called when entering the identity_clause production.
	EnterIdentity_clause(c *Identity_clauseContext)

	// EnterIdentity_options_parentheses is called when entering the identity_options_parentheses production.
	EnterIdentity_options_parentheses(c *Identity_options_parenthesesContext)

	// EnterIdentity_options is called when entering the identity_options production.
	EnterIdentity_options(c *Identity_optionsContext)

	// EnterVirtual_column_definition is called when entering the virtual_column_definition production.
	EnterVirtual_column_definition(c *Virtual_column_definitionContext)

	// EnterVirtual_column_expression is called when entering the virtual_column_expression production.
	EnterVirtual_column_expression(c *Virtual_column_expressionContext)

	// EnterAutogenerated_sequence_definition is called when entering the autogenerated_sequence_definition production.
	EnterAutogenerated_sequence_definition(c *Autogenerated_sequence_definitionContext)

	// EnterBy_user_for_statistics_clause is called when entering the by_user_for_statistics_clause production.
	EnterBy_user_for_statistics_clause(c *By_user_for_statistics_clauseContext)

	// EnterEvaluation_edition_clause is called when entering the evaluation_edition_clause production.
	EnterEvaluation_edition_clause(c *Evaluation_edition_clauseContext)

	// EnterNested_table_col_properties is called when entering the nested_table_col_properties production.
	EnterNested_table_col_properties(c *Nested_table_col_propertiesContext)

	// EnterNested_item is called when entering the nested_item production.
	EnterNested_item(c *Nested_itemContext)

	// EnterSubstitutable_column_clause is called when entering the substitutable_column_clause production.
	EnterSubstitutable_column_clause(c *Substitutable_column_clauseContext)

	// EnterPartition_name is called when entering the partition_name production.
	EnterPartition_name(c *Partition_nameContext)

	// EnterSupplemental_logging_props is called when entering the supplemental_logging_props production.
	EnterSupplemental_logging_props(c *Supplemental_logging_propsContext)

	// EnterObject_type_col_properties is called when entering the object_type_col_properties production.
	EnterObject_type_col_properties(c *Object_type_col_propertiesContext)

	// EnterConstraint_clauses is called when entering the constraint_clauses production.
	EnterConstraint_clauses(c *Constraint_clausesContext)

	// EnterOld_constraint_name is called when entering the old_constraint_name production.
	EnterOld_constraint_name(c *Old_constraint_nameContext)

	// EnterNew_constraint_name is called when entering the new_constraint_name production.
	EnterNew_constraint_name(c *New_constraint_nameContext)

	// EnterDrop_constraint_clause is called when entering the drop_constraint_clause production.
	EnterDrop_constraint_clause(c *Drop_constraint_clauseContext)

	// EnterCheck_constraint is called when entering the check_constraint production.
	EnterCheck_constraint(c *Check_constraintContext)

	// EnterForeign_key_clause is called when entering the foreign_key_clause production.
	EnterForeign_key_clause(c *Foreign_key_clauseContext)

	// EnterReferences_clause is called when entering the references_clause production.
	EnterReferences_clause(c *References_clauseContext)

	// EnterOn_delete_clause is called when entering the on_delete_clause production.
	EnterOn_delete_clause(c *On_delete_clauseContext)

	// EnterAnonymous_block is called when entering the anonymous_block production.
	EnterAnonymous_block(c *Anonymous_blockContext)

	// EnterInvoker_rights_clause is called when entering the invoker_rights_clause production.
	EnterInvoker_rights_clause(c *Invoker_rights_clauseContext)

	// EnterCall_spec is called when entering the call_spec production.
	EnterCall_spec(c *Call_specContext)

	// EnterJava_spec is called when entering the java_spec production.
	EnterJava_spec(c *Java_specContext)

	// EnterC_spec is called when entering the c_spec production.
	EnterC_spec(c *C_specContext)

	// EnterC_agent_in_clause is called when entering the c_agent_in_clause production.
	EnterC_agent_in_clause(c *C_agent_in_clauseContext)

	// EnterC_parameters_clause is called when entering the c_parameters_clause production.
	EnterC_parameters_clause(c *C_parameters_clauseContext)

	// EnterC_external_parameter is called when entering the c_external_parameter production.
	EnterC_external_parameter(c *C_external_parameterContext)

	// EnterC_property is called when entering the c_property production.
	EnterC_property(c *C_propertyContext)

	// EnterParameter is called when entering the parameter production.
	EnterParameter(c *ParameterContext)

	// EnterDefault_value_part is called when entering the default_value_part production.
	EnterDefault_value_part(c *Default_value_partContext)

	// EnterSeq_of_declare_specs is called when entering the seq_of_declare_specs production.
	EnterSeq_of_declare_specs(c *Seq_of_declare_specsContext)

	// EnterDeclare_spec is called when entering the declare_spec production.
	EnterDeclare_spec(c *Declare_specContext)

	// EnterVariable_declaration is called when entering the variable_declaration production.
	EnterVariable_declaration(c *Variable_declarationContext)

	// EnterSubtype_declaration is called when entering the subtype_declaration production.
	EnterSubtype_declaration(c *Subtype_declarationContext)

	// EnterCursor_declaration is called when entering the cursor_declaration production.
	EnterCursor_declaration(c *Cursor_declarationContext)

	// EnterParameter_spec is called when entering the parameter_spec production.
	EnterParameter_spec(c *Parameter_specContext)

	// EnterException_declaration is called when entering the exception_declaration production.
	EnterException_declaration(c *Exception_declarationContext)

	// EnterPragma_declaration is called when entering the pragma_declaration production.
	EnterPragma_declaration(c *Pragma_declarationContext)

	// EnterRecord_type_def is called when entering the record_type_def production.
	EnterRecord_type_def(c *Record_type_defContext)

	// EnterField_spec is called when entering the field_spec production.
	EnterField_spec(c *Field_specContext)

	// EnterRef_cursor_type_def is called when entering the ref_cursor_type_def production.
	EnterRef_cursor_type_def(c *Ref_cursor_type_defContext)

	// EnterType_declaration is called when entering the type_declaration production.
	EnterType_declaration(c *Type_declarationContext)

	// EnterTable_type_def is called when entering the table_type_def production.
	EnterTable_type_def(c *Table_type_defContext)

	// EnterTable_indexed_by_part is called when entering the table_indexed_by_part production.
	EnterTable_indexed_by_part(c *Table_indexed_by_partContext)

	// EnterVarray_type_def is called when entering the varray_type_def production.
	EnterVarray_type_def(c *Varray_type_defContext)

	// EnterSeq_of_statements is called when entering the seq_of_statements production.
	EnterSeq_of_statements(c *Seq_of_statementsContext)

	// EnterLabel_declaration is called when entering the label_declaration production.
	EnterLabel_declaration(c *Label_declarationContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterAssignment_statement is called when entering the assignment_statement production.
	EnterAssignment_statement(c *Assignment_statementContext)

	// EnterContinue_statement is called when entering the continue_statement production.
	EnterContinue_statement(c *Continue_statementContext)

	// EnterExit_statement is called when entering the exit_statement production.
	EnterExit_statement(c *Exit_statementContext)

	// EnterGoto_statement is called when entering the goto_statement production.
	EnterGoto_statement(c *Goto_statementContext)

	// EnterIf_statement is called when entering the if_statement production.
	EnterIf_statement(c *If_statementContext)

	// EnterElsif_part is called when entering the elsif_part production.
	EnterElsif_part(c *Elsif_partContext)

	// EnterElse_part is called when entering the else_part production.
	EnterElse_part(c *Else_partContext)

	// EnterLoop_statement is called when entering the loop_statement production.
	EnterLoop_statement(c *Loop_statementContext)

	// EnterCursor_loop_param is called when entering the cursor_loop_param production.
	EnterCursor_loop_param(c *Cursor_loop_paramContext)

	// EnterForall_statement is called when entering the forall_statement production.
	EnterForall_statement(c *Forall_statementContext)

	// EnterBounds_clause is called when entering the bounds_clause production.
	EnterBounds_clause(c *Bounds_clauseContext)

	// EnterBetween_bound is called when entering the between_bound production.
	EnterBetween_bound(c *Between_boundContext)

	// EnterLower_bound is called when entering the lower_bound production.
	EnterLower_bound(c *Lower_boundContext)

	// EnterUpper_bound is called when entering the upper_bound production.
	EnterUpper_bound(c *Upper_boundContext)

	// EnterNull_statement is called when entering the null_statement production.
	EnterNull_statement(c *Null_statementContext)

	// EnterRaise_statement is called when entering the raise_statement production.
	EnterRaise_statement(c *Raise_statementContext)

	// EnterReturn_statement is called when entering the return_statement production.
	EnterReturn_statement(c *Return_statementContext)

	// EnterCall_statement is called when entering the call_statement production.
	EnterCall_statement(c *Call_statementContext)

	// EnterPipe_row_statement is called when entering the pipe_row_statement production.
	EnterPipe_row_statement(c *Pipe_row_statementContext)

	// EnterSelection_directive is called when entering the selection_directive production.
	EnterSelection_directive(c *Selection_directiveContext)

	// EnterError_directive is called when entering the error_directive production.
	EnterError_directive(c *Error_directiveContext)

	// EnterSelection_directive_body is called when entering the selection_directive_body production.
	EnterSelection_directive_body(c *Selection_directive_bodyContext)

	// EnterBody is called when entering the body production.
	EnterBody(c *BodyContext)

	// EnterException_handler is called when entering the exception_handler production.
	EnterException_handler(c *Exception_handlerContext)

	// EnterTrigger_block is called when entering the trigger_block production.
	EnterTrigger_block(c *Trigger_blockContext)

	// EnterTps_block is called when entering the tps_block production.
	EnterTps_block(c *Tps_blockContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterSql_statement is called when entering the sql_statement production.
	EnterSql_statement(c *Sql_statementContext)

	// EnterExecute_immediate is called when entering the execute_immediate production.
	EnterExecute_immediate(c *Execute_immediateContext)

	// EnterDynamic_returning_clause is called when entering the dynamic_returning_clause production.
	EnterDynamic_returning_clause(c *Dynamic_returning_clauseContext)

	// EnterData_manipulation_language_statements is called when entering the data_manipulation_language_statements production.
	EnterData_manipulation_language_statements(c *Data_manipulation_language_statementsContext)

	// EnterCursor_manipulation_statements is called when entering the cursor_manipulation_statements production.
	EnterCursor_manipulation_statements(c *Cursor_manipulation_statementsContext)

	// EnterClose_statement is called when entering the close_statement production.
	EnterClose_statement(c *Close_statementContext)

	// EnterOpen_statement is called when entering the open_statement production.
	EnterOpen_statement(c *Open_statementContext)

	// EnterFetch_statement is called when entering the fetch_statement production.
	EnterFetch_statement(c *Fetch_statementContext)

	// EnterVariable_or_collection is called when entering the variable_or_collection production.
	EnterVariable_or_collection(c *Variable_or_collectionContext)

	// EnterOpen_for_statement is called when entering the open_for_statement production.
	EnterOpen_for_statement(c *Open_for_statementContext)

	// EnterTransaction_control_statements is called when entering the transaction_control_statements production.
	EnterTransaction_control_statements(c *Transaction_control_statementsContext)

	// EnterSet_transaction_command is called when entering the set_transaction_command production.
	EnterSet_transaction_command(c *Set_transaction_commandContext)

	// EnterSet_constraint_command is called when entering the set_constraint_command production.
	EnterSet_constraint_command(c *Set_constraint_commandContext)

	// EnterCommit_statement is called when entering the commit_statement production.
	EnterCommit_statement(c *Commit_statementContext)

	// EnterWrite_clause is called when entering the write_clause production.
	EnterWrite_clause(c *Write_clauseContext)

	// EnterRollback_statement is called when entering the rollback_statement production.
	EnterRollback_statement(c *Rollback_statementContext)

	// EnterSavepoint_statement is called when entering the savepoint_statement production.
	EnterSavepoint_statement(c *Savepoint_statementContext)

	// EnterCollection_method_call is called when entering the collection_method_call production.
	EnterCollection_method_call(c *Collection_method_callContext)

	// EnterExplain_statement is called when entering the explain_statement production.
	EnterExplain_statement(c *Explain_statementContext)

	// EnterSelect_only_statement is called when entering the select_only_statement production.
	EnterSelect_only_statement(c *Select_only_statementContext)

	// EnterSelect_statement is called when entering the select_statement production.
	EnterSelect_statement(c *Select_statementContext)

	// EnterWith_clause is called when entering the with_clause production.
	EnterWith_clause(c *With_clauseContext)

	// EnterWith_factoring_clause is called when entering the with_factoring_clause production.
	EnterWith_factoring_clause(c *With_factoring_clauseContext)

	// EnterSubquery_factoring_clause is called when entering the subquery_factoring_clause production.
	EnterSubquery_factoring_clause(c *Subquery_factoring_clauseContext)

	// EnterSearch_clause is called when entering the search_clause production.
	EnterSearch_clause(c *Search_clauseContext)

	// EnterCycle_clause is called when entering the cycle_clause production.
	EnterCycle_clause(c *Cycle_clauseContext)

	// EnterSubav_factoring_clause is called when entering the subav_factoring_clause production.
	EnterSubav_factoring_clause(c *Subav_factoring_clauseContext)

	// EnterSubav_clause is called when entering the subav_clause production.
	EnterSubav_clause(c *Subav_clauseContext)

	// EnterHierarchies_clause is called when entering the hierarchies_clause production.
	EnterHierarchies_clause(c *Hierarchies_clauseContext)

	// EnterFilter_clauses is called when entering the filter_clauses production.
	EnterFilter_clauses(c *Filter_clausesContext)

	// EnterFilter_clause is called when entering the filter_clause production.
	EnterFilter_clause(c *Filter_clauseContext)

	// EnterAdd_calcs_clause is called when entering the add_calcs_clause production.
	EnterAdd_calcs_clause(c *Add_calcs_clauseContext)

	// EnterAdd_calc_meas_clause is called when entering the add_calc_meas_clause production.
	EnterAdd_calc_meas_clause(c *Add_calc_meas_clauseContext)

	// EnterSubquery is called when entering the subquery production.
	EnterSubquery(c *SubqueryContext)

	// EnterSubquery_basic_elements is called when entering the subquery_basic_elements production.
	EnterSubquery_basic_elements(c *Subquery_basic_elementsContext)

	// EnterSubquery_operation_part is called when entering the subquery_operation_part production.
	EnterSubquery_operation_part(c *Subquery_operation_partContext)

	// EnterQuery_block is called when entering the query_block production.
	EnterQuery_block(c *Query_blockContext)

	// EnterSelected_list is called when entering the selected_list production.
	EnterSelected_list(c *Selected_listContext)

	// EnterFrom_clause is called when entering the from_clause production.
	EnterFrom_clause(c *From_clauseContext)

	// EnterSelect_list_elements is called when entering the select_list_elements production.
	EnterSelect_list_elements(c *Select_list_elementsContext)

	// EnterTable_ref_list is called when entering the table_ref_list production.
	EnterTable_ref_list(c *Table_ref_listContext)

	// EnterTable_ref is called when entering the table_ref production.
	EnterTable_ref(c *Table_refContext)

	// EnterTable_ref_aux is called when entering the table_ref_aux production.
	EnterTable_ref_aux(c *Table_ref_auxContext)

	// EnterTable_ref_aux_internal_one is called when entering the table_ref_aux_internal_one production.
	EnterTable_ref_aux_internal_one(c *Table_ref_aux_internal_oneContext)

	// EnterTable_ref_aux_internal_two is called when entering the table_ref_aux_internal_two production.
	EnterTable_ref_aux_internal_two(c *Table_ref_aux_internal_twoContext)

	// EnterTable_ref_aux_internal_thre is called when entering the table_ref_aux_internal_thre production.
	EnterTable_ref_aux_internal_thre(c *Table_ref_aux_internal_threContext)

	// EnterJoin_clause is called when entering the join_clause production.
	EnterJoin_clause(c *Join_clauseContext)

	// EnterJoin_on_part is called when entering the join_on_part production.
	EnterJoin_on_part(c *Join_on_partContext)

	// EnterJoin_using_part is called when entering the join_using_part production.
	EnterJoin_using_part(c *Join_using_partContext)

	// EnterOuter_join_type is called when entering the outer_join_type production.
	EnterOuter_join_type(c *Outer_join_typeContext)

	// EnterQuery_partition_clause is called when entering the query_partition_clause production.
	EnterQuery_partition_clause(c *Query_partition_clauseContext)

	// EnterFlashback_query_clause is called when entering the flashback_query_clause production.
	EnterFlashback_query_clause(c *Flashback_query_clauseContext)

	// EnterPivot_clause is called when entering the pivot_clause production.
	EnterPivot_clause(c *Pivot_clauseContext)

	// EnterPivot_element is called when entering the pivot_element production.
	EnterPivot_element(c *Pivot_elementContext)

	// EnterPivot_for_clause is called when entering the pivot_for_clause production.
	EnterPivot_for_clause(c *Pivot_for_clauseContext)

	// EnterPivot_in_clause is called when entering the pivot_in_clause production.
	EnterPivot_in_clause(c *Pivot_in_clauseContext)

	// EnterPivot_in_clause_element is called when entering the pivot_in_clause_element production.
	EnterPivot_in_clause_element(c *Pivot_in_clause_elementContext)

	// EnterPivot_in_clause_elements is called when entering the pivot_in_clause_elements production.
	EnterPivot_in_clause_elements(c *Pivot_in_clause_elementsContext)

	// EnterUnpivot_clause is called when entering the unpivot_clause production.
	EnterUnpivot_clause(c *Unpivot_clauseContext)

	// EnterUnpivot_in_clause is called when entering the unpivot_in_clause production.
	EnterUnpivot_in_clause(c *Unpivot_in_clauseContext)

	// EnterUnpivot_in_elements is called when entering the unpivot_in_elements production.
	EnterUnpivot_in_elements(c *Unpivot_in_elementsContext)

	// EnterHierarchical_query_clause is called when entering the hierarchical_query_clause production.
	EnterHierarchical_query_clause(c *Hierarchical_query_clauseContext)

	// EnterStart_part is called when entering the start_part production.
	EnterStart_part(c *Start_partContext)

	// EnterGroup_by_clause is called when entering the group_by_clause production.
	EnterGroup_by_clause(c *Group_by_clauseContext)

	// EnterGroup_by_elements is called when entering the group_by_elements production.
	EnterGroup_by_elements(c *Group_by_elementsContext)

	// EnterRollup_cube_clause is called when entering the rollup_cube_clause production.
	EnterRollup_cube_clause(c *Rollup_cube_clauseContext)

	// EnterGrouping_sets_clause is called when entering the grouping_sets_clause production.
	EnterGrouping_sets_clause(c *Grouping_sets_clauseContext)

	// EnterGrouping_sets_elements is called when entering the grouping_sets_elements production.
	EnterGrouping_sets_elements(c *Grouping_sets_elementsContext)

	// EnterHaving_clause is called when entering the having_clause production.
	EnterHaving_clause(c *Having_clauseContext)

	// EnterModel_clause is called when entering the model_clause production.
	EnterModel_clause(c *Model_clauseContext)

	// EnterCell_reference_options is called when entering the cell_reference_options production.
	EnterCell_reference_options(c *Cell_reference_optionsContext)

	// EnterReturn_rows_clause is called when entering the return_rows_clause production.
	EnterReturn_rows_clause(c *Return_rows_clauseContext)

	// EnterReference_model is called when entering the reference_model production.
	EnterReference_model(c *Reference_modelContext)

	// EnterMain_model is called when entering the main_model production.
	EnterMain_model(c *Main_modelContext)

	// EnterModel_column_clauses is called when entering the model_column_clauses production.
	EnterModel_column_clauses(c *Model_column_clausesContext)

	// EnterModel_column_partition_part is called when entering the model_column_partition_part production.
	EnterModel_column_partition_part(c *Model_column_partition_partContext)

	// EnterModel_column_list is called when entering the model_column_list production.
	EnterModel_column_list(c *Model_column_listContext)

	// EnterModel_column is called when entering the model_column production.
	EnterModel_column(c *Model_columnContext)

	// EnterModel_rules_clause is called when entering the model_rules_clause production.
	EnterModel_rules_clause(c *Model_rules_clauseContext)

	// EnterModel_rules_part is called when entering the model_rules_part production.
	EnterModel_rules_part(c *Model_rules_partContext)

	// EnterModel_rules_element is called when entering the model_rules_element production.
	EnterModel_rules_element(c *Model_rules_elementContext)

	// EnterCell_assignment is called when entering the cell_assignment production.
	EnterCell_assignment(c *Cell_assignmentContext)

	// EnterModel_iterate_clause is called when entering the model_iterate_clause production.
	EnterModel_iterate_clause(c *Model_iterate_clauseContext)

	// EnterUntil_part is called when entering the until_part production.
	EnterUntil_part(c *Until_partContext)

	// EnterOrder_by_clause is called when entering the order_by_clause production.
	EnterOrder_by_clause(c *Order_by_clauseContext)

	// EnterOrder_by_elements is called when entering the order_by_elements production.
	EnterOrder_by_elements(c *Order_by_elementsContext)

	// EnterOffset_clause is called when entering the offset_clause production.
	EnterOffset_clause(c *Offset_clauseContext)

	// EnterFetch_clause is called when entering the fetch_clause production.
	EnterFetch_clause(c *Fetch_clauseContext)

	// EnterFor_update_clause is called when entering the for_update_clause production.
	EnterFor_update_clause(c *For_update_clauseContext)

	// EnterFor_update_of_part is called when entering the for_update_of_part production.
	EnterFor_update_of_part(c *For_update_of_partContext)

	// EnterFor_update_options is called when entering the for_update_options production.
	EnterFor_update_options(c *For_update_optionsContext)

	// EnterUpdate_statement is called when entering the update_statement production.
	EnterUpdate_statement(c *Update_statementContext)

	// EnterUpdate_set_clause is called when entering the update_set_clause production.
	EnterUpdate_set_clause(c *Update_set_clauseContext)

	// EnterColumn_based_update_set_clause is called when entering the column_based_update_set_clause production.
	EnterColumn_based_update_set_clause(c *Column_based_update_set_clauseContext)

	// EnterDelete_statement is called when entering the delete_statement production.
	EnterDelete_statement(c *Delete_statementContext)

	// EnterInsert_statement is called when entering the insert_statement production.
	EnterInsert_statement(c *Insert_statementContext)

	// EnterSingle_table_insert is called when entering the single_table_insert production.
	EnterSingle_table_insert(c *Single_table_insertContext)

	// EnterMulti_table_insert is called when entering the multi_table_insert production.
	EnterMulti_table_insert(c *Multi_table_insertContext)

	// EnterMulti_table_element is called when entering the multi_table_element production.
	EnterMulti_table_element(c *Multi_table_elementContext)

	// EnterConditional_insert_clause is called when entering the conditional_insert_clause production.
	EnterConditional_insert_clause(c *Conditional_insert_clauseContext)

	// EnterConditional_insert_when_part is called when entering the conditional_insert_when_part production.
	EnterConditional_insert_when_part(c *Conditional_insert_when_partContext)

	// EnterConditional_insert_else_part is called when entering the conditional_insert_else_part production.
	EnterConditional_insert_else_part(c *Conditional_insert_else_partContext)

	// EnterInsert_into_clause is called when entering the insert_into_clause production.
	EnterInsert_into_clause(c *Insert_into_clauseContext)

	// EnterValues_clause is called when entering the values_clause production.
	EnterValues_clause(c *Values_clauseContext)

	// EnterMerge_statement is called when entering the merge_statement production.
	EnterMerge_statement(c *Merge_statementContext)

	// EnterMerge_update_clause is called when entering the merge_update_clause production.
	EnterMerge_update_clause(c *Merge_update_clauseContext)

	// EnterMerge_element is called when entering the merge_element production.
	EnterMerge_element(c *Merge_elementContext)

	// EnterMerge_update_delete_part is called when entering the merge_update_delete_part production.
	EnterMerge_update_delete_part(c *Merge_update_delete_partContext)

	// EnterMerge_insert_clause is called when entering the merge_insert_clause production.
	EnterMerge_insert_clause(c *Merge_insert_clauseContext)

	// EnterSelected_tableview is called when entering the selected_tableview production.
	EnterSelected_tableview(c *Selected_tableviewContext)

	// EnterLock_table_statement is called when entering the lock_table_statement production.
	EnterLock_table_statement(c *Lock_table_statementContext)

	// EnterWait_nowait_part is called when entering the wait_nowait_part production.
	EnterWait_nowait_part(c *Wait_nowait_partContext)

	// EnterLock_table_element is called when entering the lock_table_element production.
	EnterLock_table_element(c *Lock_table_elementContext)

	// EnterLock_mode is called when entering the lock_mode production.
	EnterLock_mode(c *Lock_modeContext)

	// EnterGeneral_table_ref is called when entering the general_table_ref production.
	EnterGeneral_table_ref(c *General_table_refContext)

	// EnterStatic_returning_clause is called when entering the static_returning_clause production.
	EnterStatic_returning_clause(c *Static_returning_clauseContext)

	// EnterError_logging_clause is called when entering the error_logging_clause production.
	EnterError_logging_clause(c *Error_logging_clauseContext)

	// EnterError_logging_into_part is called when entering the error_logging_into_part production.
	EnterError_logging_into_part(c *Error_logging_into_partContext)

	// EnterError_logging_reject_part is called when entering the error_logging_reject_part production.
	EnterError_logging_reject_part(c *Error_logging_reject_partContext)

	// EnterDml_table_expression_clause is called when entering the dml_table_expression_clause production.
	EnterDml_table_expression_clause(c *Dml_table_expression_clauseContext)

	// EnterTable_collection_expression is called when entering the table_collection_expression production.
	EnterTable_collection_expression(c *Table_collection_expressionContext)

	// EnterSubquery_restriction_clause is called when entering the subquery_restriction_clause production.
	EnterSubquery_restriction_clause(c *Subquery_restriction_clauseContext)

	// EnterSample_clause is called when entering the sample_clause production.
	EnterSample_clause(c *Sample_clauseContext)

	// EnterSeed_part is called when entering the seed_part production.
	EnterSeed_part(c *Seed_partContext)

	// EnterCondition is called when entering the condition production.
	EnterCondition(c *ConditionContext)

	// EnterExpressions_ is called when entering the expressions_ production.
	EnterExpressions_(c *Expressions_Context)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterCursor_expression is called when entering the cursor_expression production.
	EnterCursor_expression(c *Cursor_expressionContext)

	// EnterLogical_expression is called when entering the logical_expression production.
	EnterLogical_expression(c *Logical_expressionContext)

	// EnterUnary_logical_expression is called when entering the unary_logical_expression production.
	EnterUnary_logical_expression(c *Unary_logical_expressionContext)

	// EnterUnary_logical_operation is called when entering the unary_logical_operation production.
	EnterUnary_logical_operation(c *Unary_logical_operationContext)

	// EnterLogical_operation is called when entering the logical_operation production.
	EnterLogical_operation(c *Logical_operationContext)

	// EnterMultiset_expression is called when entering the multiset_expression production.
	EnterMultiset_expression(c *Multiset_expressionContext)

	// EnterRelational_expression is called when entering the relational_expression production.
	EnterRelational_expression(c *Relational_expressionContext)

	// EnterCompound_expression is called when entering the compound_expression production.
	EnterCompound_expression(c *Compound_expressionContext)

	// EnterRelational_operator is called when entering the relational_operator production.
	EnterRelational_operator(c *Relational_operatorContext)

	// EnterIn_elements is called when entering the in_elements production.
	EnterIn_elements(c *In_elementsContext)

	// EnterBetween_elements is called when entering the between_elements production.
	EnterBetween_elements(c *Between_elementsContext)

	// EnterConcatenation is called when entering the concatenation production.
	EnterConcatenation(c *ConcatenationContext)

	// EnterInterval_expression is called when entering the interval_expression production.
	EnterInterval_expression(c *Interval_expressionContext)

	// EnterModel_expression is called when entering the model_expression production.
	EnterModel_expression(c *Model_expressionContext)

	// EnterModel_expression_element is called when entering the model_expression_element production.
	EnterModel_expression_element(c *Model_expression_elementContext)

	// EnterSingle_column_for_loop is called when entering the single_column_for_loop production.
	EnterSingle_column_for_loop(c *Single_column_for_loopContext)

	// EnterMulti_column_for_loop is called when entering the multi_column_for_loop production.
	EnterMulti_column_for_loop(c *Multi_column_for_loopContext)

	// EnterUnary_expression is called when entering the unary_expression production.
	EnterUnary_expression(c *Unary_expressionContext)

	// EnterUnary_expression_core is called when entering the unary_expression_core production.
	EnterUnary_expression_core(c *Unary_expression_coreContext)

	// EnterImplicit_cursor_expression is called when entering the implicit_cursor_expression production.
	EnterImplicit_cursor_expression(c *Implicit_cursor_expressionContext)

	// EnterCollection_expression is called when entering the collection_expression production.
	EnterCollection_expression(c *Collection_expressionContext)

	// EnterCase_statement is called when entering the case_statement production.
	EnterCase_statement(c *Case_statementContext)

	// EnterSimple_case_statement is called when entering the simple_case_statement production.
	EnterSimple_case_statement(c *Simple_case_statementContext)

	// EnterSearched_case_statement is called when entering the searched_case_statement production.
	EnterSearched_case_statement(c *Searched_case_statementContext)

	// EnterCase_when_part_statement is called when entering the case_when_part_statement production.
	EnterCase_when_part_statement(c *Case_when_part_statementContext)

	// EnterCase_else_part_statement is called when entering the case_else_part_statement production.
	EnterCase_else_part_statement(c *Case_else_part_statementContext)

	// EnterCase_expression is called when entering the case_expression production.
	EnterCase_expression(c *Case_expressionContext)

	// EnterSimple_case_expression is called when entering the simple_case_expression production.
	EnterSimple_case_expression(c *Simple_case_expressionContext)

	// EnterSearched_case_expression is called when entering the searched_case_expression production.
	EnterSearched_case_expression(c *Searched_case_expressionContext)

	// EnterCase_when_part_expression is called when entering the case_when_part_expression production.
	EnterCase_when_part_expression(c *Case_when_part_expressionContext)

	// EnterCase_else_part_expression is called when entering the case_else_part_expression production.
	EnterCase_else_part_expression(c *Case_else_part_expressionContext)

	// EnterAtom is called when entering the atom production.
	EnterAtom(c *AtomContext)

	// EnterQuantified_expression is called when entering the quantified_expression production.
	EnterQuantified_expression(c *Quantified_expressionContext)

	// EnterString_function is called when entering the string_function production.
	EnterString_function(c *String_functionContext)

	// EnterStandard_function is called when entering the standard_function production.
	EnterStandard_function(c *Standard_functionContext)

	// EnterJson_function is called when entering the json_function production.
	EnterJson_function(c *Json_functionContext)

	// EnterJson_object_content is called when entering the json_object_content production.
	EnterJson_object_content(c *Json_object_contentContext)

	// EnterJson_object_entry is called when entering the json_object_entry production.
	EnterJson_object_entry(c *Json_object_entryContext)

	// EnterJson_table_clause is called when entering the json_table_clause production.
	EnterJson_table_clause(c *Json_table_clauseContext)

	// EnterJson_array_element is called when entering the json_array_element production.
	EnterJson_array_element(c *Json_array_elementContext)

	// EnterJson_on_null_clause is called when entering the json_on_null_clause production.
	EnterJson_on_null_clause(c *Json_on_null_clauseContext)

	// EnterJson_return_clause is called when entering the json_return_clause production.
	EnterJson_return_clause(c *Json_return_clauseContext)

	// EnterJson_transform_op is called when entering the json_transform_op production.
	EnterJson_transform_op(c *Json_transform_opContext)

	// EnterJson_column_clause is called when entering the json_column_clause production.
	EnterJson_column_clause(c *Json_column_clauseContext)

	// EnterJson_column_definition is called when entering the json_column_definition production.
	EnterJson_column_definition(c *Json_column_definitionContext)

	// EnterJson_query_returning_clause is called when entering the json_query_returning_clause production.
	EnterJson_query_returning_clause(c *Json_query_returning_clauseContext)

	// EnterJson_query_return_type is called when entering the json_query_return_type production.
	EnterJson_query_return_type(c *Json_query_return_typeContext)

	// EnterJson_query_wrapper_clause is called when entering the json_query_wrapper_clause production.
	EnterJson_query_wrapper_clause(c *Json_query_wrapper_clauseContext)

	// EnterJson_query_on_error_clause is called when entering the json_query_on_error_clause production.
	EnterJson_query_on_error_clause(c *Json_query_on_error_clauseContext)

	// EnterJson_query_on_empty_clause is called when entering the json_query_on_empty_clause production.
	EnterJson_query_on_empty_clause(c *Json_query_on_empty_clauseContext)

	// EnterJson_value_return_clause is called when entering the json_value_return_clause production.
	EnterJson_value_return_clause(c *Json_value_return_clauseContext)

	// EnterJson_value_return_type is called when entering the json_value_return_type production.
	EnterJson_value_return_type(c *Json_value_return_typeContext)

	// EnterJson_value_on_mismatch_clause is called when entering the json_value_on_mismatch_clause production.
	EnterJson_value_on_mismatch_clause(c *Json_value_on_mismatch_clauseContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterNumeric_function_wrapper is called when entering the numeric_function_wrapper production.
	EnterNumeric_function_wrapper(c *Numeric_function_wrapperContext)

	// EnterNumeric_function is called when entering the numeric_function production.
	EnterNumeric_function(c *Numeric_functionContext)

	// EnterListagg_overflow_clause is called when entering the listagg_overflow_clause production.
	EnterListagg_overflow_clause(c *Listagg_overflow_clauseContext)

	// EnterOther_function is called when entering the other_function production.
	EnterOther_function(c *Other_functionContext)

	// EnterOver_clause_keyword is called when entering the over_clause_keyword production.
	EnterOver_clause_keyword(c *Over_clause_keywordContext)

	// EnterWithin_or_over_clause_keyword is called when entering the within_or_over_clause_keyword production.
	EnterWithin_or_over_clause_keyword(c *Within_or_over_clause_keywordContext)

	// EnterStandard_prediction_function_keyword is called when entering the standard_prediction_function_keyword production.
	EnterStandard_prediction_function_keyword(c *Standard_prediction_function_keywordContext)

	// EnterOver_clause is called when entering the over_clause production.
	EnterOver_clause(c *Over_clauseContext)

	// EnterWindowing_clause is called when entering the windowing_clause production.
	EnterWindowing_clause(c *Windowing_clauseContext)

	// EnterWindowing_type is called when entering the windowing_type production.
	EnterWindowing_type(c *Windowing_typeContext)

	// EnterWindowing_elements is called when entering the windowing_elements production.
	EnterWindowing_elements(c *Windowing_elementsContext)

	// EnterUsing_clause is called when entering the using_clause production.
	EnterUsing_clause(c *Using_clauseContext)

	// EnterUsing_element is called when entering the using_element production.
	EnterUsing_element(c *Using_elementContext)

	// EnterAssignable_element is called when entering the assignable_element production.
	EnterAssignable_element(c *Assignable_elementContext)

	// EnterCollect_order_by_part is called when entering the collect_order_by_part production.
	EnterCollect_order_by_part(c *Collect_order_by_partContext)

	// EnterWithin_or_over_part is called when entering the within_or_over_part production.
	EnterWithin_or_over_part(c *Within_or_over_partContext)

	// EnterString_delimiter is called when entering the string_delimiter production.
	EnterString_delimiter(c *String_delimiterContext)

	// EnterCost_matrix_clause is called when entering the cost_matrix_clause production.
	EnterCost_matrix_clause(c *Cost_matrix_clauseContext)

	// EnterXml_passing_clause is called when entering the xml_passing_clause production.
	EnterXml_passing_clause(c *Xml_passing_clauseContext)

	// EnterXml_attributes_clause is called when entering the xml_attributes_clause production.
	EnterXml_attributes_clause(c *Xml_attributes_clauseContext)

	// EnterXml_namespaces_clause is called when entering the xml_namespaces_clause production.
	EnterXml_namespaces_clause(c *Xml_namespaces_clauseContext)

	// EnterXml_table_column is called when entering the xml_table_column production.
	EnterXml_table_column(c *Xml_table_columnContext)

	// EnterXml_general_default_part is called when entering the xml_general_default_part production.
	EnterXml_general_default_part(c *Xml_general_default_partContext)

	// EnterXml_multiuse_expression_element is called when entering the xml_multiuse_expression_element production.
	EnterXml_multiuse_expression_element(c *Xml_multiuse_expression_elementContext)

	// EnterXmlroot_param_version_part is called when entering the xmlroot_param_version_part production.
	EnterXmlroot_param_version_part(c *Xmlroot_param_version_partContext)

	// EnterXmlroot_param_standalone_part is called when entering the xmlroot_param_standalone_part production.
	EnterXmlroot_param_standalone_part(c *Xmlroot_param_standalone_partContext)

	// EnterXmlserialize_param_enconding_part is called when entering the xmlserialize_param_enconding_part production.
	EnterXmlserialize_param_enconding_part(c *Xmlserialize_param_enconding_partContext)

	// EnterXmlserialize_param_version_part is called when entering the xmlserialize_param_version_part production.
	EnterXmlserialize_param_version_part(c *Xmlserialize_param_version_partContext)

	// EnterXmlserialize_param_ident_part is called when entering the xmlserialize_param_ident_part production.
	EnterXmlserialize_param_ident_part(c *Xmlserialize_param_ident_partContext)

	// EnterAnnotations_clause is called when entering the annotations_clause production.
	EnterAnnotations_clause(c *Annotations_clauseContext)

	// EnterAnnotations_list is called when entering the annotations_list production.
	EnterAnnotations_list(c *Annotations_listContext)

	// EnterAnnotation is called when entering the annotation production.
	EnterAnnotation(c *AnnotationContext)

	// EnterSql_plus_command is called when entering the sql_plus_command production.
	EnterSql_plus_command(c *Sql_plus_commandContext)

	// EnterStart_command is called when entering the start_command production.
	EnterStart_command(c *Start_commandContext)

	// EnterSql_plus_filepath is called when entering the sql_plus_filepath production.
	EnterSql_plus_filepath(c *Sql_plus_filepathContext)

	// EnterWhenever_command is called when entering the whenever_command production.
	EnterWhenever_command(c *Whenever_commandContext)

	// EnterSet_command is called when entering the set_command production.
	EnterSet_command(c *Set_commandContext)

	// EnterTiming_command is called when entering the timing_command production.
	EnterTiming_command(c *Timing_commandContext)

	// EnterClear_command is called when entering the clear_command production.
	EnterClear_command(c *Clear_commandContext)

	// EnterPartition_extension_clause is called when entering the partition_extension_clause production.
	EnterPartition_extension_clause(c *Partition_extension_clauseContext)

	// EnterColumn_alias is called when entering the column_alias production.
	EnterColumn_alias(c *Column_aliasContext)

	// EnterTable_alias is called when entering the table_alias production.
	EnterTable_alias(c *Table_aliasContext)

	// EnterWhere_clause is called when entering the where_clause production.
	EnterWhere_clause(c *Where_clauseContext)

	// EnterInto_clause is called when entering the into_clause production.
	EnterInto_clause(c *Into_clauseContext)

	// EnterXml_column_name is called when entering the xml_column_name production.
	EnterXml_column_name(c *Xml_column_nameContext)

	// EnterCost_class_name is called when entering the cost_class_name production.
	EnterCost_class_name(c *Cost_class_nameContext)

	// EnterAttribute_name is called when entering the attribute_name production.
	EnterAttribute_name(c *Attribute_nameContext)

	// EnterSavepoint_name is called when entering the savepoint_name production.
	EnterSavepoint_name(c *Savepoint_nameContext)

	// EnterRollback_segment_name is called when entering the rollback_segment_name production.
	EnterRollback_segment_name(c *Rollback_segment_nameContext)

	// EnterSchema_name is called when entering the schema_name production.
	EnterSchema_name(c *Schema_nameContext)

	// EnterRoutine_name is called when entering the routine_name production.
	EnterRoutine_name(c *Routine_nameContext)

	// EnterPackage_name is called when entering the package_name production.
	EnterPackage_name(c *Package_nameContext)

	// EnterImplementation_type_name is called when entering the implementation_type_name production.
	EnterImplementation_type_name(c *Implementation_type_nameContext)

	// EnterParameter_name is called when entering the parameter_name production.
	EnterParameter_name(c *Parameter_nameContext)

	// EnterReference_model_name is called when entering the reference_model_name production.
	EnterReference_model_name(c *Reference_model_nameContext)

	// EnterMain_model_name is called when entering the main_model_name production.
	EnterMain_model_name(c *Main_model_nameContext)

	// EnterContainer_tableview_name is called when entering the container_tableview_name production.
	EnterContainer_tableview_name(c *Container_tableview_nameContext)

	// EnterAggregate_function_name is called when entering the aggregate_function_name production.
	EnterAggregate_function_name(c *Aggregate_function_nameContext)

	// EnterQuery_name is called when entering the query_name production.
	EnterQuery_name(c *Query_nameContext)

	// EnterGrantee_name is called when entering the grantee_name production.
	EnterGrantee_name(c *Grantee_nameContext)

	// EnterRole_name is called when entering the role_name production.
	EnterRole_name(c *Role_nameContext)

	// EnterConstraint_name is called when entering the constraint_name production.
	EnterConstraint_name(c *Constraint_nameContext)

	// EnterLabel_name is called when entering the label_name production.
	EnterLabel_name(c *Label_nameContext)

	// EnterType_name is called when entering the type_name production.
	EnterType_name(c *Type_nameContext)

	// EnterSequence_name is called when entering the sequence_name production.
	EnterSequence_name(c *Sequence_nameContext)

	// EnterException_name is called when entering the exception_name production.
	EnterException_name(c *Exception_nameContext)

	// EnterFunction_name is called when entering the function_name production.
	EnterFunction_name(c *Function_nameContext)

	// EnterProcedure_name is called when entering the procedure_name production.
	EnterProcedure_name(c *Procedure_nameContext)

	// EnterTrigger_name is called when entering the trigger_name production.
	EnterTrigger_name(c *Trigger_nameContext)

	// EnterVariable_name is called when entering the variable_name production.
	EnterVariable_name(c *Variable_nameContext)

	// EnterIndex_name is called when entering the index_name production.
	EnterIndex_name(c *Index_nameContext)

	// EnterCursor_name is called when entering the cursor_name production.
	EnterCursor_name(c *Cursor_nameContext)

	// EnterRecord_name is called when entering the record_name production.
	EnterRecord_name(c *Record_nameContext)

	// EnterLink_name is called when entering the link_name production.
	EnterLink_name(c *Link_nameContext)

	// EnterLocal_link_name is called when entering the local_link_name production.
	EnterLocal_link_name(c *Local_link_nameContext)

	// EnterConnection_qualifier is called when entering the connection_qualifier production.
	EnterConnection_qualifier(c *Connection_qualifierContext)

	// EnterColumn_name is called when entering the column_name production.
	EnterColumn_name(c *Column_nameContext)

	// EnterTableview_name is called when entering the tableview_name production.
	EnterTableview_name(c *Tableview_nameContext)

	// EnterXmltable is called when entering the xmltable production.
	EnterXmltable(c *XmltableContext)

	// EnterChar_set_name is called when entering the char_set_name production.
	EnterChar_set_name(c *Char_set_nameContext)

	// EnterSynonym_name is called when entering the synonym_name production.
	EnterSynonym_name(c *Synonym_nameContext)

	// EnterSchema_object_name is called when entering the schema_object_name production.
	EnterSchema_object_name(c *Schema_object_nameContext)

	// EnterDir_object_name is called when entering the dir_object_name production.
	EnterDir_object_name(c *Dir_object_nameContext)

	// EnterUser_object_name is called when entering the user_object_name production.
	EnterUser_object_name(c *User_object_nameContext)

	// EnterGrant_object_name is called when entering the grant_object_name production.
	EnterGrant_object_name(c *Grant_object_nameContext)

	// EnterColumn_list is called when entering the column_list production.
	EnterColumn_list(c *Column_listContext)

	// EnterParen_column_list is called when entering the paren_column_list production.
	EnterParen_column_list(c *Paren_column_listContext)

	// EnterKeep_clause is called when entering the keep_clause production.
	EnterKeep_clause(c *Keep_clauseContext)

	// EnterFunction_argument is called when entering the function_argument production.
	EnterFunction_argument(c *Function_argumentContext)

	// EnterFunction_argument_analytic is called when entering the function_argument_analytic production.
	EnterFunction_argument_analytic(c *Function_argument_analyticContext)

	// EnterFunction_argument_modeling is called when entering the function_argument_modeling production.
	EnterFunction_argument_modeling(c *Function_argument_modelingContext)

	// EnterRespect_or_ignore_nulls is called when entering the respect_or_ignore_nulls production.
	EnterRespect_or_ignore_nulls(c *Respect_or_ignore_nullsContext)

	// EnterArgument is called when entering the argument production.
	EnterArgument(c *ArgumentContext)

	// EnterType_spec is called when entering the type_spec production.
	EnterType_spec(c *Type_specContext)

	// EnterDatatype is called when entering the datatype production.
	EnterDatatype(c *DatatypeContext)

	// EnterPrecision_part is called when entering the precision_part production.
	EnterPrecision_part(c *Precision_partContext)

	// EnterNative_datatype_element is called when entering the native_datatype_element production.
	EnterNative_datatype_element(c *Native_datatype_elementContext)

	// EnterBind_variable is called when entering the bind_variable production.
	EnterBind_variable(c *Bind_variableContext)

	// EnterGeneral_element is called when entering the general_element production.
	EnterGeneral_element(c *General_elementContext)

	// EnterGeneral_element_part is called when entering the general_element_part production.
	EnterGeneral_element_part(c *General_element_partContext)

	// EnterTable_element is called when entering the table_element production.
	EnterTable_element(c *Table_elementContext)

	// EnterObject_privilege is called when entering the object_privilege production.
	EnterObject_privilege(c *Object_privilegeContext)

	// EnterSystem_privilege is called when entering the system_privilege production.
	EnterSystem_privilege(c *System_privilegeContext)

	// EnterConstant is called when entering the constant production.
	EnterConstant(c *ConstantContext)

	// EnterNumeric is called when entering the numeric production.
	EnterNumeric(c *NumericContext)

	// EnterNumeric_negative is called when entering the numeric_negative production.
	EnterNumeric_negative(c *Numeric_negativeContext)

	// EnterQuoted_string is called when entering the quoted_string production.
	EnterQuoted_string(c *Quoted_stringContext)

	// EnterIdentifier is called when entering the identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterId_expression is called when entering the id_expression production.
	EnterId_expression(c *Id_expressionContext)

	// EnterInquiry_directive is called when entering the inquiry_directive production.
	EnterInquiry_directive(c *Inquiry_directiveContext)

	// EnterOuter_join_sign is called when entering the outer_join_sign production.
	EnterOuter_join_sign(c *Outer_join_signContext)

	// EnterRegular_id is called when entering the regular_id production.
	EnterRegular_id(c *Regular_idContext)

	// EnterNon_reserved_keywords_in_18c is called when entering the non_reserved_keywords_in_18c production.
	EnterNon_reserved_keywords_in_18c(c *Non_reserved_keywords_in_18cContext)

	// EnterNon_reserved_keywords_in_12c is called when entering the non_reserved_keywords_in_12c production.
	EnterNon_reserved_keywords_in_12c(c *Non_reserved_keywords_in_12cContext)

	// EnterNon_reserved_keywords_pre12c is called when entering the non_reserved_keywords_pre12c production.
	EnterNon_reserved_keywords_pre12c(c *Non_reserved_keywords_pre12cContext)

	// ExitSql_script is called when exiting the sql_script production.
	ExitSql_script(c *Sql_scriptContext)

	// ExitUnit_statement is called when exiting the unit_statement production.
	ExitUnit_statement(c *Unit_statementContext)

	// ExitAlter_diskgroup is called when exiting the alter_diskgroup production.
	ExitAlter_diskgroup(c *Alter_diskgroupContext)

	// ExitAdd_disk_clause is called when exiting the add_disk_clause production.
	ExitAdd_disk_clause(c *Add_disk_clauseContext)

	// ExitDrop_disk_clause is called when exiting the drop_disk_clause production.
	ExitDrop_disk_clause(c *Drop_disk_clauseContext)

	// ExitResize_disk_clause is called when exiting the resize_disk_clause production.
	ExitResize_disk_clause(c *Resize_disk_clauseContext)

	// ExitReplace_disk_clause is called when exiting the replace_disk_clause production.
	ExitReplace_disk_clause(c *Replace_disk_clauseContext)

	// ExitWait_nowait is called when exiting the wait_nowait production.
	ExitWait_nowait(c *Wait_nowaitContext)

	// ExitRename_disk_clause is called when exiting the rename_disk_clause production.
	ExitRename_disk_clause(c *Rename_disk_clauseContext)

	// ExitDisk_online_clause is called when exiting the disk_online_clause production.
	ExitDisk_online_clause(c *Disk_online_clauseContext)

	// ExitDisk_offline_clause is called when exiting the disk_offline_clause production.
	ExitDisk_offline_clause(c *Disk_offline_clauseContext)

	// ExitTimeout_clause is called when exiting the timeout_clause production.
	ExitTimeout_clause(c *Timeout_clauseContext)

	// ExitRebalance_diskgroup_clause is called when exiting the rebalance_diskgroup_clause production.
	ExitRebalance_diskgroup_clause(c *Rebalance_diskgroup_clauseContext)

	// ExitPhase is called when exiting the phase production.
	ExitPhase(c *PhaseContext)

	// ExitCheck_diskgroup_clause is called when exiting the check_diskgroup_clause production.
	ExitCheck_diskgroup_clause(c *Check_diskgroup_clauseContext)

	// ExitDiskgroup_template_clauses is called when exiting the diskgroup_template_clauses production.
	ExitDiskgroup_template_clauses(c *Diskgroup_template_clausesContext)

	// ExitQualified_template_clause is called when exiting the qualified_template_clause production.
	ExitQualified_template_clause(c *Qualified_template_clauseContext)

	// ExitRedundancy_clause is called when exiting the redundancy_clause production.
	ExitRedundancy_clause(c *Redundancy_clauseContext)

	// ExitStriping_clause is called when exiting the striping_clause production.
	ExitStriping_clause(c *Striping_clauseContext)

	// ExitForce_noforce is called when exiting the force_noforce production.
	ExitForce_noforce(c *Force_noforceContext)

	// ExitDiskgroup_directory_clauses is called when exiting the diskgroup_directory_clauses production.
	ExitDiskgroup_directory_clauses(c *Diskgroup_directory_clausesContext)

	// ExitDir_name is called when exiting the dir_name production.
	ExitDir_name(c *Dir_nameContext)

	// ExitDiskgroup_alias_clauses is called when exiting the diskgroup_alias_clauses production.
	ExitDiskgroup_alias_clauses(c *Diskgroup_alias_clausesContext)

	// ExitDiskgroup_volume_clauses is called when exiting the diskgroup_volume_clauses production.
	ExitDiskgroup_volume_clauses(c *Diskgroup_volume_clausesContext)

	// ExitAdd_volume_clause is called when exiting the add_volume_clause production.
	ExitAdd_volume_clause(c *Add_volume_clauseContext)

	// ExitModify_volume_clause is called when exiting the modify_volume_clause production.
	ExitModify_volume_clause(c *Modify_volume_clauseContext)

	// ExitDiskgroup_attributes is called when exiting the diskgroup_attributes production.
	ExitDiskgroup_attributes(c *Diskgroup_attributesContext)

	// ExitDrop_diskgroup_file_clause is called when exiting the drop_diskgroup_file_clause production.
	ExitDrop_diskgroup_file_clause(c *Drop_diskgroup_file_clauseContext)

	// ExitConvert_redundancy_clause is called when exiting the convert_redundancy_clause production.
	ExitConvert_redundancy_clause(c *Convert_redundancy_clauseContext)

	// ExitUsergroup_clauses is called when exiting the usergroup_clauses production.
	ExitUsergroup_clauses(c *Usergroup_clausesContext)

	// ExitUser_clauses is called when exiting the user_clauses production.
	ExitUser_clauses(c *User_clausesContext)

	// ExitFile_permissions_clause is called when exiting the file_permissions_clause production.
	ExitFile_permissions_clause(c *File_permissions_clauseContext)

	// ExitFile_owner_clause is called when exiting the file_owner_clause production.
	ExitFile_owner_clause(c *File_owner_clauseContext)

	// ExitScrub_clause is called when exiting the scrub_clause production.
	ExitScrub_clause(c *Scrub_clauseContext)

	// ExitQuotagroup_clauses is called when exiting the quotagroup_clauses production.
	ExitQuotagroup_clauses(c *Quotagroup_clausesContext)

	// ExitProperty_name is called when exiting the property_name production.
	ExitProperty_name(c *Property_nameContext)

	// ExitProperty_value is called when exiting the property_value production.
	ExitProperty_value(c *Property_valueContext)

	// ExitFilegroup_clauses is called when exiting the filegroup_clauses production.
	ExitFilegroup_clauses(c *Filegroup_clausesContext)

	// ExitAdd_filegroup_clause is called when exiting the add_filegroup_clause production.
	ExitAdd_filegroup_clause(c *Add_filegroup_clauseContext)

	// ExitModify_filegroup_clause is called when exiting the modify_filegroup_clause production.
	ExitModify_filegroup_clause(c *Modify_filegroup_clauseContext)

	// ExitMove_to_filegroup_clause is called when exiting the move_to_filegroup_clause production.
	ExitMove_to_filegroup_clause(c *Move_to_filegroup_clauseContext)

	// ExitDrop_filegroup_clause is called when exiting the drop_filegroup_clause production.
	ExitDrop_filegroup_clause(c *Drop_filegroup_clauseContext)

	// ExitQuorum_regular is called when exiting the quorum_regular production.
	ExitQuorum_regular(c *Quorum_regularContext)

	// ExitUndrop_disk_clause is called when exiting the undrop_disk_clause production.
	ExitUndrop_disk_clause(c *Undrop_disk_clauseContext)

	// ExitDiskgroup_availability is called when exiting the diskgroup_availability production.
	ExitDiskgroup_availability(c *Diskgroup_availabilityContext)

	// ExitEnable_disable_volume is called when exiting the enable_disable_volume production.
	ExitEnable_disable_volume(c *Enable_disable_volumeContext)

	// ExitDrop_function is called when exiting the drop_function production.
	ExitDrop_function(c *Drop_functionContext)

	// ExitAlter_flashback_archive is called when exiting the alter_flashback_archive production.
	ExitAlter_flashback_archive(c *Alter_flashback_archiveContext)

	// ExitAlter_hierarchy is called when exiting the alter_hierarchy production.
	ExitAlter_hierarchy(c *Alter_hierarchyContext)

	// ExitAlter_function is called when exiting the alter_function production.
	ExitAlter_function(c *Alter_functionContext)

	// ExitAlter_java is called when exiting the alter_java production.
	ExitAlter_java(c *Alter_javaContext)

	// ExitMatch_string is called when exiting the match_string production.
	ExitMatch_string(c *Match_stringContext)

	// ExitCreate_function_body is called when exiting the create_function_body production.
	ExitCreate_function_body(c *Create_function_bodyContext)

	// ExitSql_macro_body is called when exiting the sql_macro_body production.
	ExitSql_macro_body(c *Sql_macro_bodyContext)

	// ExitParallel_enable_clause is called when exiting the parallel_enable_clause production.
	ExitParallel_enable_clause(c *Parallel_enable_clauseContext)

	// ExitPartition_by_clause is called when exiting the partition_by_clause production.
	ExitPartition_by_clause(c *Partition_by_clauseContext)

	// ExitResult_cache_clause is called when exiting the result_cache_clause production.
	ExitResult_cache_clause(c *Result_cache_clauseContext)

	// ExitAccessible_by_clause is called when exiting the accessible_by_clause production.
	ExitAccessible_by_clause(c *Accessible_by_clauseContext)

	// ExitDefault_collation_clause is called when exiting the default_collation_clause production.
	ExitDefault_collation_clause(c *Default_collation_clauseContext)

	// ExitAggregate_clause is called when exiting the aggregate_clause production.
	ExitAggregate_clause(c *Aggregate_clauseContext)

	// ExitPipelined_using_clause is called when exiting the pipelined_using_clause production.
	ExitPipelined_using_clause(c *Pipelined_using_clauseContext)

	// ExitAccessor is called when exiting the accessor production.
	ExitAccessor(c *AccessorContext)

	// ExitRelies_on_part is called when exiting the relies_on_part production.
	ExitRelies_on_part(c *Relies_on_partContext)

	// ExitStreaming_clause is called when exiting the streaming_clause production.
	ExitStreaming_clause(c *Streaming_clauseContext)

	// ExitAlter_outline is called when exiting the alter_outline production.
	ExitAlter_outline(c *Alter_outlineContext)

	// ExitOutline_options is called when exiting the outline_options production.
	ExitOutline_options(c *Outline_optionsContext)

	// ExitAlter_lockdown_profile is called when exiting the alter_lockdown_profile production.
	ExitAlter_lockdown_profile(c *Alter_lockdown_profileContext)

	// ExitLockdown_feature is called when exiting the lockdown_feature production.
	ExitLockdown_feature(c *Lockdown_featureContext)

	// ExitLockdown_options is called when exiting the lockdown_options production.
	ExitLockdown_options(c *Lockdown_optionsContext)

	// ExitLockdown_statements is called when exiting the lockdown_statements production.
	ExitLockdown_statements(c *Lockdown_statementsContext)

	// ExitStatement_clauses is called when exiting the statement_clauses production.
	ExitStatement_clauses(c *Statement_clausesContext)

	// ExitClause_options is called when exiting the clause_options production.
	ExitClause_options(c *Clause_optionsContext)

	// ExitOption_values is called when exiting the option_values production.
	ExitOption_values(c *Option_valuesContext)

	// ExitString_list is called when exiting the string_list production.
	ExitString_list(c *String_listContext)

	// ExitDisable_enable is called when exiting the disable_enable production.
	ExitDisable_enable(c *Disable_enableContext)

	// ExitDrop_lockdown_profile is called when exiting the drop_lockdown_profile production.
	ExitDrop_lockdown_profile(c *Drop_lockdown_profileContext)

	// ExitDrop_package is called when exiting the drop_package production.
	ExitDrop_package(c *Drop_packageContext)

	// ExitAlter_package is called when exiting the alter_package production.
	ExitAlter_package(c *Alter_packageContext)

	// ExitCreate_package is called when exiting the create_package production.
	ExitCreate_package(c *Create_packageContext)

	// ExitCreate_package_body is called when exiting the create_package_body production.
	ExitCreate_package_body(c *Create_package_bodyContext)

	// ExitPackage_obj_spec is called when exiting the package_obj_spec production.
	ExitPackage_obj_spec(c *Package_obj_specContext)

	// ExitProcedure_spec is called when exiting the procedure_spec production.
	ExitProcedure_spec(c *Procedure_specContext)

	// ExitFunction_spec is called when exiting the function_spec production.
	ExitFunction_spec(c *Function_specContext)

	// ExitPackage_obj_body is called when exiting the package_obj_body production.
	ExitPackage_obj_body(c *Package_obj_bodyContext)

	// ExitAlter_pmem_filestore is called when exiting the alter_pmem_filestore production.
	ExitAlter_pmem_filestore(c *Alter_pmem_filestoreContext)

	// ExitDrop_pmem_filestore is called when exiting the drop_pmem_filestore production.
	ExitDrop_pmem_filestore(c *Drop_pmem_filestoreContext)

	// ExitDrop_procedure is called when exiting the drop_procedure production.
	ExitDrop_procedure(c *Drop_procedureContext)

	// ExitAlter_procedure is called when exiting the alter_procedure production.
	ExitAlter_procedure(c *Alter_procedureContext)

	// ExitFunction_body is called when exiting the function_body production.
	ExitFunction_body(c *Function_bodyContext)

	// ExitProcedure_body is called when exiting the procedure_body production.
	ExitProcedure_body(c *Procedure_bodyContext)

	// ExitCreate_procedure_body is called when exiting the create_procedure_body production.
	ExitCreate_procedure_body(c *Create_procedure_bodyContext)

	// ExitAlter_resource_cost is called when exiting the alter_resource_cost production.
	ExitAlter_resource_cost(c *Alter_resource_costContext)

	// ExitDrop_outline is called when exiting the drop_outline production.
	ExitDrop_outline(c *Drop_outlineContext)

	// ExitAlter_rollback_segment is called when exiting the alter_rollback_segment production.
	ExitAlter_rollback_segment(c *Alter_rollback_segmentContext)

	// ExitDrop_restore_point is called when exiting the drop_restore_point production.
	ExitDrop_restore_point(c *Drop_restore_pointContext)

	// ExitDrop_rollback_segment is called when exiting the drop_rollback_segment production.
	ExitDrop_rollback_segment(c *Drop_rollback_segmentContext)

	// ExitDrop_role is called when exiting the drop_role production.
	ExitDrop_role(c *Drop_roleContext)

	// ExitCreate_pmem_filestore is called when exiting the create_pmem_filestore production.
	ExitCreate_pmem_filestore(c *Create_pmem_filestoreContext)

	// ExitPmem_filestore_options is called when exiting the pmem_filestore_options production.
	ExitPmem_filestore_options(c *Pmem_filestore_optionsContext)

	// ExitFile_path is called when exiting the file_path production.
	ExitFile_path(c *File_pathContext)

	// ExitCreate_rollback_segment is called when exiting the create_rollback_segment production.
	ExitCreate_rollback_segment(c *Create_rollback_segmentContext)

	// ExitDrop_trigger is called when exiting the drop_trigger production.
	ExitDrop_trigger(c *Drop_triggerContext)

	// ExitAlter_trigger is called when exiting the alter_trigger production.
	ExitAlter_trigger(c *Alter_triggerContext)

	// ExitCreate_trigger is called when exiting the create_trigger production.
	ExitCreate_trigger(c *Create_triggerContext)

	// ExitTrigger_follows_clause is called when exiting the trigger_follows_clause production.
	ExitTrigger_follows_clause(c *Trigger_follows_clauseContext)

	// ExitTrigger_when_clause is called when exiting the trigger_when_clause production.
	ExitTrigger_when_clause(c *Trigger_when_clauseContext)

	// ExitSimple_dml_trigger is called when exiting the simple_dml_trigger production.
	ExitSimple_dml_trigger(c *Simple_dml_triggerContext)

	// ExitFor_each_row is called when exiting the for_each_row production.
	ExitFor_each_row(c *For_each_rowContext)

	// ExitCompound_dml_trigger is called when exiting the compound_dml_trigger production.
	ExitCompound_dml_trigger(c *Compound_dml_triggerContext)

	// ExitNon_dml_trigger is called when exiting the non_dml_trigger production.
	ExitNon_dml_trigger(c *Non_dml_triggerContext)

	// ExitTrigger_body is called when exiting the trigger_body production.
	ExitTrigger_body(c *Trigger_bodyContext)

	// ExitCompound_trigger_block is called when exiting the compound_trigger_block production.
	ExitCompound_trigger_block(c *Compound_trigger_blockContext)

	// ExitTiming_point_section is called when exiting the timing_point_section production.
	ExitTiming_point_section(c *Timing_point_sectionContext)

	// ExitNon_dml_event is called when exiting the non_dml_event production.
	ExitNon_dml_event(c *Non_dml_eventContext)

	// ExitDml_event_clause is called when exiting the dml_event_clause production.
	ExitDml_event_clause(c *Dml_event_clauseContext)

	// ExitDml_event_element is called when exiting the dml_event_element production.
	ExitDml_event_element(c *Dml_event_elementContext)

	// ExitDml_event_nested_clause is called when exiting the dml_event_nested_clause production.
	ExitDml_event_nested_clause(c *Dml_event_nested_clauseContext)

	// ExitReferencing_clause is called when exiting the referencing_clause production.
	ExitReferencing_clause(c *Referencing_clauseContext)

	// ExitReferencing_element is called when exiting the referencing_element production.
	ExitReferencing_element(c *Referencing_elementContext)

	// ExitDrop_type is called when exiting the drop_type production.
	ExitDrop_type(c *Drop_typeContext)

	// ExitAlter_type is called when exiting the alter_type production.
	ExitAlter_type(c *Alter_typeContext)

	// ExitCompile_type_clause is called when exiting the compile_type_clause production.
	ExitCompile_type_clause(c *Compile_type_clauseContext)

	// ExitReplace_type_clause is called when exiting the replace_type_clause production.
	ExitReplace_type_clause(c *Replace_type_clauseContext)

	// ExitAlter_method_spec is called when exiting the alter_method_spec production.
	ExitAlter_method_spec(c *Alter_method_specContext)

	// ExitAlter_method_element is called when exiting the alter_method_element production.
	ExitAlter_method_element(c *Alter_method_elementContext)

	// ExitAlter_collection_clauses is called when exiting the alter_collection_clauses production.
	ExitAlter_collection_clauses(c *Alter_collection_clausesContext)

	// ExitDependent_handling_clause is called when exiting the dependent_handling_clause production.
	ExitDependent_handling_clause(c *Dependent_handling_clauseContext)

	// ExitDependent_exceptions_part is called when exiting the dependent_exceptions_part production.
	ExitDependent_exceptions_part(c *Dependent_exceptions_partContext)

	// ExitCreate_type is called when exiting the create_type production.
	ExitCreate_type(c *Create_typeContext)

	// ExitType_definition is called when exiting the type_definition production.
	ExitType_definition(c *Type_definitionContext)

	// ExitObject_type_def is called when exiting the object_type_def production.
	ExitObject_type_def(c *Object_type_defContext)

	// ExitObject_as_part is called when exiting the object_as_part production.
	ExitObject_as_part(c *Object_as_partContext)

	// ExitObject_under_part is called when exiting the object_under_part production.
	ExitObject_under_part(c *Object_under_partContext)

	// ExitNested_table_type_def is called when exiting the nested_table_type_def production.
	ExitNested_table_type_def(c *Nested_table_type_defContext)

	// ExitSqlj_object_type is called when exiting the sqlj_object_type production.
	ExitSqlj_object_type(c *Sqlj_object_typeContext)

	// ExitType_body is called when exiting the type_body production.
	ExitType_body(c *Type_bodyContext)

	// ExitType_body_elements is called when exiting the type_body_elements production.
	ExitType_body_elements(c *Type_body_elementsContext)

	// ExitMap_order_func_declaration is called when exiting the map_order_func_declaration production.
	ExitMap_order_func_declaration(c *Map_order_func_declarationContext)

	// ExitSubprog_decl_in_type is called when exiting the subprog_decl_in_type production.
	ExitSubprog_decl_in_type(c *Subprog_decl_in_typeContext)

	// ExitProc_decl_in_type is called when exiting the proc_decl_in_type production.
	ExitProc_decl_in_type(c *Proc_decl_in_typeContext)

	// ExitFunc_decl_in_type is called when exiting the func_decl_in_type production.
	ExitFunc_decl_in_type(c *Func_decl_in_typeContext)

	// ExitConstructor_declaration is called when exiting the constructor_declaration production.
	ExitConstructor_declaration(c *Constructor_declarationContext)

	// ExitModifier_clause is called when exiting the modifier_clause production.
	ExitModifier_clause(c *Modifier_clauseContext)

	// ExitObject_member_spec is called when exiting the object_member_spec production.
	ExitObject_member_spec(c *Object_member_specContext)

	// ExitSqlj_object_type_attr is called when exiting the sqlj_object_type_attr production.
	ExitSqlj_object_type_attr(c *Sqlj_object_type_attrContext)

	// ExitElement_spec is called when exiting the element_spec production.
	ExitElement_spec(c *Element_specContext)

	// ExitElement_spec_options is called when exiting the element_spec_options production.
	ExitElement_spec_options(c *Element_spec_optionsContext)

	// ExitSubprogram_spec is called when exiting the subprogram_spec production.
	ExitSubprogram_spec(c *Subprogram_specContext)

	// ExitOverriding_subprogram_spec is called when exiting the overriding_subprogram_spec production.
	ExitOverriding_subprogram_spec(c *Overriding_subprogram_specContext)

	// ExitOverriding_function_spec is called when exiting the overriding_function_spec production.
	ExitOverriding_function_spec(c *Overriding_function_specContext)

	// ExitOverriding_procedure_spec is called when exiting the overriding_procedure_spec production.
	ExitOverriding_procedure_spec(c *Overriding_procedure_specContext)

	// ExitType_procedure_spec is called when exiting the type_procedure_spec production.
	ExitType_procedure_spec(c *Type_procedure_specContext)

	// ExitType_function_spec is called when exiting the type_function_spec production.
	ExitType_function_spec(c *Type_function_specContext)

	// ExitConstructor_spec is called when exiting the constructor_spec production.
	ExitConstructor_spec(c *Constructor_specContext)

	// ExitMap_order_function_spec is called when exiting the map_order_function_spec production.
	ExitMap_order_function_spec(c *Map_order_function_specContext)

	// ExitPragma_clause is called when exiting the pragma_clause production.
	ExitPragma_clause(c *Pragma_clauseContext)

	// ExitPragma_elements is called when exiting the pragma_elements production.
	ExitPragma_elements(c *Pragma_elementsContext)

	// ExitType_elements_parameter is called when exiting the type_elements_parameter production.
	ExitType_elements_parameter(c *Type_elements_parameterContext)

	// ExitDrop_sequence is called when exiting the drop_sequence production.
	ExitDrop_sequence(c *Drop_sequenceContext)

	// ExitAlter_sequence is called when exiting the alter_sequence production.
	ExitAlter_sequence(c *Alter_sequenceContext)

	// ExitAlter_session is called when exiting the alter_session production.
	ExitAlter_session(c *Alter_sessionContext)

	// ExitAlter_session_set_clause is called when exiting the alter_session_set_clause production.
	ExitAlter_session_set_clause(c *Alter_session_set_clauseContext)

	// ExitCreate_sequence is called when exiting the create_sequence production.
	ExitCreate_sequence(c *Create_sequenceContext)

	// ExitSequence_spec is called when exiting the sequence_spec production.
	ExitSequence_spec(c *Sequence_specContext)

	// ExitSequence_start_clause is called when exiting the sequence_start_clause production.
	ExitSequence_start_clause(c *Sequence_start_clauseContext)

	// ExitCreate_analytic_view is called when exiting the create_analytic_view production.
	ExitCreate_analytic_view(c *Create_analytic_viewContext)

	// ExitClassification_clause is called when exiting the classification_clause production.
	ExitClassification_clause(c *Classification_clauseContext)

	// ExitCaption_clause is called when exiting the caption_clause production.
	ExitCaption_clause(c *Caption_clauseContext)

	// ExitDescription_clause is called when exiting the description_clause production.
	ExitDescription_clause(c *Description_clauseContext)

	// ExitClassification_item is called when exiting the classification_item production.
	ExitClassification_item(c *Classification_itemContext)

	// ExitLanguage is called when exiting the language production.
	ExitLanguage(c *LanguageContext)

	// ExitCav_using_clause is called when exiting the cav_using_clause production.
	ExitCav_using_clause(c *Cav_using_clauseContext)

	// ExitDim_by_clause is called when exiting the dim_by_clause production.
	ExitDim_by_clause(c *Dim_by_clauseContext)

	// ExitDim_key is called when exiting the dim_key production.
	ExitDim_key(c *Dim_keyContext)

	// ExitDim_ref is called when exiting the dim_ref production.
	ExitDim_ref(c *Dim_refContext)

	// ExitHier_ref is called when exiting the hier_ref production.
	ExitHier_ref(c *Hier_refContext)

	// ExitMeasures_clause is called when exiting the measures_clause production.
	ExitMeasures_clause(c *Measures_clauseContext)

	// ExitAv_measure is called when exiting the av_measure production.
	ExitAv_measure(c *Av_measureContext)

	// ExitBase_meas_clause is called when exiting the base_meas_clause production.
	ExitBase_meas_clause(c *Base_meas_clauseContext)

	// ExitMeas_aggregate_clause is called when exiting the meas_aggregate_clause production.
	ExitMeas_aggregate_clause(c *Meas_aggregate_clauseContext)

	// ExitCalc_meas_clause is called when exiting the calc_meas_clause production.
	ExitCalc_meas_clause(c *Calc_meas_clauseContext)

	// ExitDefault_measure_clause is called when exiting the default_measure_clause production.
	ExitDefault_measure_clause(c *Default_measure_clauseContext)

	// ExitDefault_aggregate_clause is called when exiting the default_aggregate_clause production.
	ExitDefault_aggregate_clause(c *Default_aggregate_clauseContext)

	// ExitCache_clause is called when exiting the cache_clause production.
	ExitCache_clause(c *Cache_clauseContext)

	// ExitCache_specification is called when exiting the cache_specification production.
	ExitCache_specification(c *Cache_specificationContext)

	// ExitLevels_clause is called when exiting the levels_clause production.
	ExitLevels_clause(c *Levels_clauseContext)

	// ExitLevel_specification is called when exiting the level_specification production.
	ExitLevel_specification(c *Level_specificationContext)

	// ExitLevel_group_type is called when exiting the level_group_type production.
	ExitLevel_group_type(c *Level_group_typeContext)

	// ExitFact_columns_clause is called when exiting the fact_columns_clause production.
	ExitFact_columns_clause(c *Fact_columns_clauseContext)

	// ExitQry_transform_clause is called when exiting the qry_transform_clause production.
	ExitQry_transform_clause(c *Qry_transform_clauseContext)

	// ExitCreate_attribute_dimension is called when exiting the create_attribute_dimension production.
	ExitCreate_attribute_dimension(c *Create_attribute_dimensionContext)

	// ExitAd_using_clause is called when exiting the ad_using_clause production.
	ExitAd_using_clause(c *Ad_using_clauseContext)

	// ExitSource_clause is called when exiting the source_clause production.
	ExitSource_clause(c *Source_clauseContext)

	// ExitJoin_path_clause is called when exiting the join_path_clause production.
	ExitJoin_path_clause(c *Join_path_clauseContext)

	// ExitJoin_condition is called when exiting the join_condition production.
	ExitJoin_condition(c *Join_conditionContext)

	// ExitJoin_condition_item is called when exiting the join_condition_item production.
	ExitJoin_condition_item(c *Join_condition_itemContext)

	// ExitAttributes_clause is called when exiting the attributes_clause production.
	ExitAttributes_clause(c *Attributes_clauseContext)

	// ExitAd_attributes_clause is called when exiting the ad_attributes_clause production.
	ExitAd_attributes_clause(c *Ad_attributes_clauseContext)

	// ExitAd_level_clause is called when exiting the ad_level_clause production.
	ExitAd_level_clause(c *Ad_level_clauseContext)

	// ExitKey_clause is called when exiting the key_clause production.
	ExitKey_clause(c *Key_clauseContext)

	// ExitAlternate_key_clause is called when exiting the alternate_key_clause production.
	ExitAlternate_key_clause(c *Alternate_key_clauseContext)

	// ExitDim_order_clause is called when exiting the dim_order_clause production.
	ExitDim_order_clause(c *Dim_order_clauseContext)

	// ExitAll_clause is called when exiting the all_clause production.
	ExitAll_clause(c *All_clauseContext)

	// ExitCreate_audit_policy is called when exiting the create_audit_policy production.
	ExitCreate_audit_policy(c *Create_audit_policyContext)

	// ExitPrivilege_audit_clause is called when exiting the privilege_audit_clause production.
	ExitPrivilege_audit_clause(c *Privilege_audit_clauseContext)

	// ExitAction_audit_clause is called when exiting the action_audit_clause production.
	ExitAction_audit_clause(c *Action_audit_clauseContext)

	// ExitSystem_actions is called when exiting the system_actions production.
	ExitSystem_actions(c *System_actionsContext)

	// ExitStandard_actions is called when exiting the standard_actions production.
	ExitStandard_actions(c *Standard_actionsContext)

	// ExitActions_clause is called when exiting the actions_clause production.
	ExitActions_clause(c *Actions_clauseContext)

	// ExitObject_action is called when exiting the object_action production.
	ExitObject_action(c *Object_actionContext)

	// ExitSystem_action is called when exiting the system_action production.
	ExitSystem_action(c *System_actionContext)

	// ExitComponent_actions is called when exiting the component_actions production.
	ExitComponent_actions(c *Component_actionsContext)

	// ExitComponent_action is called when exiting the component_action production.
	ExitComponent_action(c *Component_actionContext)

	// ExitRole_audit_clause is called when exiting the role_audit_clause production.
	ExitRole_audit_clause(c *Role_audit_clauseContext)

	// ExitCreate_controlfile is called when exiting the create_controlfile production.
	ExitCreate_controlfile(c *Create_controlfileContext)

	// ExitControlfile_options is called when exiting the controlfile_options production.
	ExitControlfile_options(c *Controlfile_optionsContext)

	// ExitLogfile_clause is called when exiting the logfile_clause production.
	ExitLogfile_clause(c *Logfile_clauseContext)

	// ExitCharacter_set_clause is called when exiting the character_set_clause production.
	ExitCharacter_set_clause(c *Character_set_clauseContext)

	// ExitFile_specification is called when exiting the file_specification production.
	ExitFile_specification(c *File_specificationContext)

	// ExitCreate_diskgroup is called when exiting the create_diskgroup production.
	ExitCreate_diskgroup(c *Create_diskgroupContext)

	// ExitQualified_disk_clause is called when exiting the qualified_disk_clause production.
	ExitQualified_disk_clause(c *Qualified_disk_clauseContext)

	// ExitCreate_edition is called when exiting the create_edition production.
	ExitCreate_edition(c *Create_editionContext)

	// ExitCreate_flashback_archive is called when exiting the create_flashback_archive production.
	ExitCreate_flashback_archive(c *Create_flashback_archiveContext)

	// ExitFlashback_archive_quota is called when exiting the flashback_archive_quota production.
	ExitFlashback_archive_quota(c *Flashback_archive_quotaContext)

	// ExitFlashback_archive_retention is called when exiting the flashback_archive_retention production.
	ExitFlashback_archive_retention(c *Flashback_archive_retentionContext)

	// ExitCreate_hierarchy is called when exiting the create_hierarchy production.
	ExitCreate_hierarchy(c *Create_hierarchyContext)

	// ExitHier_using_clause is called when exiting the hier_using_clause production.
	ExitHier_using_clause(c *Hier_using_clauseContext)

	// ExitLevel_hier_clause is called when exiting the level_hier_clause production.
	ExitLevel_hier_clause(c *Level_hier_clauseContext)

	// ExitHier_attrs_clause is called when exiting the hier_attrs_clause production.
	ExitHier_attrs_clause(c *Hier_attrs_clauseContext)

	// ExitHier_attr_clause is called when exiting the hier_attr_clause production.
	ExitHier_attr_clause(c *Hier_attr_clauseContext)

	// ExitHier_attr_name is called when exiting the hier_attr_name production.
	ExitHier_attr_name(c *Hier_attr_nameContext)

	// ExitCreate_index is called when exiting the create_index production.
	ExitCreate_index(c *Create_indexContext)

	// ExitCluster_index_clause is called when exiting the cluster_index_clause production.
	ExitCluster_index_clause(c *Cluster_index_clauseContext)

	// ExitCluster_name is called when exiting the cluster_name production.
	ExitCluster_name(c *Cluster_nameContext)

	// ExitTable_index_clause is called when exiting the table_index_clause production.
	ExitTable_index_clause(c *Table_index_clauseContext)

	// ExitBitmap_join_index_clause is called when exiting the bitmap_join_index_clause production.
	ExitBitmap_join_index_clause(c *Bitmap_join_index_clauseContext)

	// ExitIndex_expr is called when exiting the index_expr production.
	ExitIndex_expr(c *Index_exprContext)

	// ExitIndex_properties is called when exiting the index_properties production.
	ExitIndex_properties(c *Index_propertiesContext)

	// ExitDomain_index_clause is called when exiting the domain_index_clause production.
	ExitDomain_index_clause(c *Domain_index_clauseContext)

	// ExitLocal_domain_index_clause is called when exiting the local_domain_index_clause production.
	ExitLocal_domain_index_clause(c *Local_domain_index_clauseContext)

	// ExitXmlindex_clause is called when exiting the xmlindex_clause production.
	ExitXmlindex_clause(c *Xmlindex_clauseContext)

	// ExitLocal_xmlindex_clause is called when exiting the local_xmlindex_clause production.
	ExitLocal_xmlindex_clause(c *Local_xmlindex_clauseContext)

	// ExitGlobal_partitioned_index is called when exiting the global_partitioned_index production.
	ExitGlobal_partitioned_index(c *Global_partitioned_indexContext)

	// ExitIndex_partitioning_clause is called when exiting the index_partitioning_clause production.
	ExitIndex_partitioning_clause(c *Index_partitioning_clauseContext)

	// ExitIndex_partitioning_values_list is called when exiting the index_partitioning_values_list production.
	ExitIndex_partitioning_values_list(c *Index_partitioning_values_listContext)

	// ExitLocal_partitioned_index is called when exiting the local_partitioned_index production.
	ExitLocal_partitioned_index(c *Local_partitioned_indexContext)

	// ExitOn_range_partitioned_table is called when exiting the on_range_partitioned_table production.
	ExitOn_range_partitioned_table(c *On_range_partitioned_tableContext)

	// ExitOn_list_partitioned_table is called when exiting the on_list_partitioned_table production.
	ExitOn_list_partitioned_table(c *On_list_partitioned_tableContext)

	// ExitPartitioned_table is called when exiting the partitioned_table production.
	ExitPartitioned_table(c *Partitioned_tableContext)

	// ExitOn_hash_partitioned_table is called when exiting the on_hash_partitioned_table production.
	ExitOn_hash_partitioned_table(c *On_hash_partitioned_tableContext)

	// ExitOn_hash_partitioned_clause is called when exiting the on_hash_partitioned_clause production.
	ExitOn_hash_partitioned_clause(c *On_hash_partitioned_clauseContext)

	// ExitOn_comp_partitioned_table is called when exiting the on_comp_partitioned_table production.
	ExitOn_comp_partitioned_table(c *On_comp_partitioned_tableContext)

	// ExitOn_comp_partitioned_clause is called when exiting the on_comp_partitioned_clause production.
	ExitOn_comp_partitioned_clause(c *On_comp_partitioned_clauseContext)

	// ExitIndex_subpartition_clause is called when exiting the index_subpartition_clause production.
	ExitIndex_subpartition_clause(c *Index_subpartition_clauseContext)

	// ExitIndex_subpartition_subclause is called when exiting the index_subpartition_subclause production.
	ExitIndex_subpartition_subclause(c *Index_subpartition_subclauseContext)

	// ExitOdci_parameters is called when exiting the odci_parameters production.
	ExitOdci_parameters(c *Odci_parametersContext)

	// ExitIndextype is called when exiting the indextype production.
	ExitIndextype(c *IndextypeContext)

	// ExitAlter_index is called when exiting the alter_index production.
	ExitAlter_index(c *Alter_indexContext)

	// ExitAlter_index_ops_set1 is called when exiting the alter_index_ops_set1 production.
	ExitAlter_index_ops_set1(c *Alter_index_ops_set1Context)

	// ExitAlter_index_ops_set2 is called when exiting the alter_index_ops_set2 production.
	ExitAlter_index_ops_set2(c *Alter_index_ops_set2Context)

	// ExitVisible_or_invisible is called when exiting the visible_or_invisible production.
	ExitVisible_or_invisible(c *Visible_or_invisibleContext)

	// ExitMonitoring_nomonitoring is called when exiting the monitoring_nomonitoring production.
	ExitMonitoring_nomonitoring(c *Monitoring_nomonitoringContext)

	// ExitRebuild_clause is called when exiting the rebuild_clause production.
	ExitRebuild_clause(c *Rebuild_clauseContext)

	// ExitAlter_index_partitioning is called when exiting the alter_index_partitioning production.
	ExitAlter_index_partitioning(c *Alter_index_partitioningContext)

	// ExitModify_index_default_attrs is called when exiting the modify_index_default_attrs production.
	ExitModify_index_default_attrs(c *Modify_index_default_attrsContext)

	// ExitAdd_hash_index_partition is called when exiting the add_hash_index_partition production.
	ExitAdd_hash_index_partition(c *Add_hash_index_partitionContext)

	// ExitCoalesce_index_partition is called when exiting the coalesce_index_partition production.
	ExitCoalesce_index_partition(c *Coalesce_index_partitionContext)

	// ExitModify_index_partition is called when exiting the modify_index_partition production.
	ExitModify_index_partition(c *Modify_index_partitionContext)

	// ExitModify_index_partitions_ops is called when exiting the modify_index_partitions_ops production.
	ExitModify_index_partitions_ops(c *Modify_index_partitions_opsContext)

	// ExitRename_index_partition is called when exiting the rename_index_partition production.
	ExitRename_index_partition(c *Rename_index_partitionContext)

	// ExitDrop_index_partition is called when exiting the drop_index_partition production.
	ExitDrop_index_partition(c *Drop_index_partitionContext)

	// ExitSplit_index_partition is called when exiting the split_index_partition production.
	ExitSplit_index_partition(c *Split_index_partitionContext)

	// ExitIndex_partition_description is called when exiting the index_partition_description production.
	ExitIndex_partition_description(c *Index_partition_descriptionContext)

	// ExitModify_index_subpartition is called when exiting the modify_index_subpartition production.
	ExitModify_index_subpartition(c *Modify_index_subpartitionContext)

	// ExitPartition_name_old is called when exiting the partition_name_old production.
	ExitPartition_name_old(c *Partition_name_oldContext)

	// ExitNew_partition_name is called when exiting the new_partition_name production.
	ExitNew_partition_name(c *New_partition_nameContext)

	// ExitNew_index_name is called when exiting the new_index_name production.
	ExitNew_index_name(c *New_index_nameContext)

	// ExitAlter_inmemory_join_group is called when exiting the alter_inmemory_join_group production.
	ExitAlter_inmemory_join_group(c *Alter_inmemory_join_groupContext)

	// ExitCreate_user is called when exiting the create_user production.
	ExitCreate_user(c *Create_userContext)

	// ExitAlter_user is called when exiting the alter_user production.
	ExitAlter_user(c *Alter_userContext)

	// ExitDrop_user is called when exiting the drop_user production.
	ExitDrop_user(c *Drop_userContext)

	// ExitAlter_identified_by is called when exiting the alter_identified_by production.
	ExitAlter_identified_by(c *Alter_identified_byContext)

	// ExitIdentified_by is called when exiting the identified_by production.
	ExitIdentified_by(c *Identified_byContext)

	// ExitIdentified_other_clause is called when exiting the identified_other_clause production.
	ExitIdentified_other_clause(c *Identified_other_clauseContext)

	// ExitUser_tablespace_clause is called when exiting the user_tablespace_clause production.
	ExitUser_tablespace_clause(c *User_tablespace_clauseContext)

	// ExitQuota_clause is called when exiting the quota_clause production.
	ExitQuota_clause(c *Quota_clauseContext)

	// ExitProfile_clause is called when exiting the profile_clause production.
	ExitProfile_clause(c *Profile_clauseContext)

	// ExitRole_clause is called when exiting the role_clause production.
	ExitRole_clause(c *Role_clauseContext)

	// ExitUser_default_role_clause is called when exiting the user_default_role_clause production.
	ExitUser_default_role_clause(c *User_default_role_clauseContext)

	// ExitPassword_expire_clause is called when exiting the password_expire_clause production.
	ExitPassword_expire_clause(c *Password_expire_clauseContext)

	// ExitUser_lock_clause is called when exiting the user_lock_clause production.
	ExitUser_lock_clause(c *User_lock_clauseContext)

	// ExitUser_editions_clause is called when exiting the user_editions_clause production.
	ExitUser_editions_clause(c *User_editions_clauseContext)

	// ExitAlter_user_editions_clause is called when exiting the alter_user_editions_clause production.
	ExitAlter_user_editions_clause(c *Alter_user_editions_clauseContext)

	// ExitProxy_clause is called when exiting the proxy_clause production.
	ExitProxy_clause(c *Proxy_clauseContext)

	// ExitContainer_names is called when exiting the container_names production.
	ExitContainer_names(c *Container_namesContext)

	// ExitSet_container_data is called when exiting the set_container_data production.
	ExitSet_container_data(c *Set_container_dataContext)

	// ExitAdd_rem_container_data is called when exiting the add_rem_container_data production.
	ExitAdd_rem_container_data(c *Add_rem_container_dataContext)

	// ExitContainer_data_clause is called when exiting the container_data_clause production.
	ExitContainer_data_clause(c *Container_data_clauseContext)

	// ExitAdminister_key_management is called when exiting the administer_key_management production.
	ExitAdminister_key_management(c *Administer_key_managementContext)

	// ExitKeystore_management_clauses is called when exiting the keystore_management_clauses production.
	ExitKeystore_management_clauses(c *Keystore_management_clausesContext)

	// ExitCreate_keystore is called when exiting the create_keystore production.
	ExitCreate_keystore(c *Create_keystoreContext)

	// ExitOpen_keystore is called when exiting the open_keystore production.
	ExitOpen_keystore(c *Open_keystoreContext)

	// ExitForce_keystore is called when exiting the force_keystore production.
	ExitForce_keystore(c *Force_keystoreContext)

	// ExitClose_keystore is called when exiting the close_keystore production.
	ExitClose_keystore(c *Close_keystoreContext)

	// ExitBackup_keystore is called when exiting the backup_keystore production.
	ExitBackup_keystore(c *Backup_keystoreContext)

	// ExitAlter_keystore_password is called when exiting the alter_keystore_password production.
	ExitAlter_keystore_password(c *Alter_keystore_passwordContext)

	// ExitMerge_into_new_keystore is called when exiting the merge_into_new_keystore production.
	ExitMerge_into_new_keystore(c *Merge_into_new_keystoreContext)

	// ExitMerge_into_existing_keystore is called when exiting the merge_into_existing_keystore production.
	ExitMerge_into_existing_keystore(c *Merge_into_existing_keystoreContext)

	// ExitIsolate_keystore is called when exiting the isolate_keystore production.
	ExitIsolate_keystore(c *Isolate_keystoreContext)

	// ExitUnite_keystore is called when exiting the unite_keystore production.
	ExitUnite_keystore(c *Unite_keystoreContext)

	// ExitKey_management_clauses is called when exiting the key_management_clauses production.
	ExitKey_management_clauses(c *Key_management_clausesContext)

	// ExitSet_key is called when exiting the set_key production.
	ExitSet_key(c *Set_keyContext)

	// ExitCreate_key is called when exiting the create_key production.
	ExitCreate_key(c *Create_keyContext)

	// ExitMkid is called when exiting the mkid production.
	ExitMkid(c *MkidContext)

	// ExitMk is called when exiting the mk production.
	ExitMk(c *MkContext)

	// ExitUse_key is called when exiting the use_key production.
	ExitUse_key(c *Use_keyContext)

	// ExitSet_key_tag is called when exiting the set_key_tag production.
	ExitSet_key_tag(c *Set_key_tagContext)

	// ExitExport_keys is called when exiting the export_keys production.
	ExitExport_keys(c *Export_keysContext)

	// ExitImport_keys is called when exiting the import_keys production.
	ExitImport_keys(c *Import_keysContext)

	// ExitMigrate_keys is called when exiting the migrate_keys production.
	ExitMigrate_keys(c *Migrate_keysContext)

	// ExitReverse_migrate_keys is called when exiting the reverse_migrate_keys production.
	ExitReverse_migrate_keys(c *Reverse_migrate_keysContext)

	// ExitMove_keys is called when exiting the move_keys production.
	ExitMove_keys(c *Move_keysContext)

	// ExitIdentified_by_store is called when exiting the identified_by_store production.
	ExitIdentified_by_store(c *Identified_by_storeContext)

	// ExitUsing_algorithm_clause is called when exiting the using_algorithm_clause production.
	ExitUsing_algorithm_clause(c *Using_algorithm_clauseContext)

	// ExitUsing_tag_clause is called when exiting the using_tag_clause production.
	ExitUsing_tag_clause(c *Using_tag_clauseContext)

	// ExitSecret_management_clauses is called when exiting the secret_management_clauses production.
	ExitSecret_management_clauses(c *Secret_management_clausesContext)

	// ExitAdd_update_secret is called when exiting the add_update_secret production.
	ExitAdd_update_secret(c *Add_update_secretContext)

	// ExitDelete_secret is called when exiting the delete_secret production.
	ExitDelete_secret(c *Delete_secretContext)

	// ExitAdd_update_secret_seps is called when exiting the add_update_secret_seps production.
	ExitAdd_update_secret_seps(c *Add_update_secret_sepsContext)

	// ExitDelete_secret_seps is called when exiting the delete_secret_seps production.
	ExitDelete_secret_seps(c *Delete_secret_sepsContext)

	// ExitZero_downtime_software_patching_clauses is called when exiting the zero_downtime_software_patching_clauses production.
	ExitZero_downtime_software_patching_clauses(c *Zero_downtime_software_patching_clausesContext)

	// ExitWith_backup_clause is called when exiting the with_backup_clause production.
	ExitWith_backup_clause(c *With_backup_clauseContext)

	// ExitIdentified_by_password_clause is called when exiting the identified_by_password_clause production.
	ExitIdentified_by_password_clause(c *Identified_by_password_clauseContext)

	// ExitKeystore_password is called when exiting the keystore_password production.
	ExitKeystore_password(c *Keystore_passwordContext)

	// ExitPath is called when exiting the path production.
	ExitPath(c *PathContext)

	// ExitSecret is called when exiting the secret production.
	ExitSecret(c *SecretContext)

	// ExitAnalyze is called when exiting the analyze production.
	ExitAnalyze(c *AnalyzeContext)

	// ExitPartition_extention_clause is called when exiting the partition_extention_clause production.
	ExitPartition_extention_clause(c *Partition_extention_clauseContext)

	// ExitValidation_clauses is called when exiting the validation_clauses production.
	ExitValidation_clauses(c *Validation_clausesContext)

	// ExitCompute_clauses is called when exiting the compute_clauses production.
	ExitCompute_clauses(c *Compute_clausesContext)

	// ExitFor_clause is called when exiting the for_clause production.
	ExitFor_clause(c *For_clauseContext)

	// ExitOnline_or_offline is called when exiting the online_or_offline production.
	ExitOnline_or_offline(c *Online_or_offlineContext)

	// ExitInto_clause1 is called when exiting the into_clause1 production.
	ExitInto_clause1(c *Into_clause1Context)

	// ExitPartition_key_value is called when exiting the partition_key_value production.
	ExitPartition_key_value(c *Partition_key_valueContext)

	// ExitSubpartition_key_value is called when exiting the subpartition_key_value production.
	ExitSubpartition_key_value(c *Subpartition_key_valueContext)

	// ExitAssociate_statistics is called when exiting the associate_statistics production.
	ExitAssociate_statistics(c *Associate_statisticsContext)

	// ExitColumn_association is called when exiting the column_association production.
	ExitColumn_association(c *Column_associationContext)

	// ExitFunction_association is called when exiting the function_association production.
	ExitFunction_association(c *Function_associationContext)

	// ExitIndextype_name is called when exiting the indextype_name production.
	ExitIndextype_name(c *Indextype_nameContext)

	// ExitUsing_statistics_type is called when exiting the using_statistics_type production.
	ExitUsing_statistics_type(c *Using_statistics_typeContext)

	// ExitStatistics_type_name is called when exiting the statistics_type_name production.
	ExitStatistics_type_name(c *Statistics_type_nameContext)

	// ExitDefault_cost_clause is called when exiting the default_cost_clause production.
	ExitDefault_cost_clause(c *Default_cost_clauseContext)

	// ExitCpu_cost is called when exiting the cpu_cost production.
	ExitCpu_cost(c *Cpu_costContext)

	// ExitIo_cost is called when exiting the io_cost production.
	ExitIo_cost(c *Io_costContext)

	// ExitNetwork_cost is called when exiting the network_cost production.
	ExitNetwork_cost(c *Network_costContext)

	// ExitDefault_selectivity_clause is called when exiting the default_selectivity_clause production.
	ExitDefault_selectivity_clause(c *Default_selectivity_clauseContext)

	// ExitDefault_selectivity is called when exiting the default_selectivity production.
	ExitDefault_selectivity(c *Default_selectivityContext)

	// ExitStorage_table_clause is called when exiting the storage_table_clause production.
	ExitStorage_table_clause(c *Storage_table_clauseContext)

	// ExitUnified_auditing is called when exiting the unified_auditing production.
	ExitUnified_auditing(c *Unified_auditingContext)

	// ExitPolicy_name is called when exiting the policy_name production.
	ExitPolicy_name(c *Policy_nameContext)

	// ExitAudit_traditional is called when exiting the audit_traditional production.
	ExitAudit_traditional(c *Audit_traditionalContext)

	// ExitAudit_direct_path is called when exiting the audit_direct_path production.
	ExitAudit_direct_path(c *Audit_direct_pathContext)

	// ExitAudit_container_clause is called when exiting the audit_container_clause production.
	ExitAudit_container_clause(c *Audit_container_clauseContext)

	// ExitAudit_operation_clause is called when exiting the audit_operation_clause production.
	ExitAudit_operation_clause(c *Audit_operation_clauseContext)

	// ExitAuditing_by_clause is called when exiting the auditing_by_clause production.
	ExitAuditing_by_clause(c *Auditing_by_clauseContext)

	// ExitAudit_user is called when exiting the audit_user production.
	ExitAudit_user(c *Audit_userContext)

	// ExitAudit_schema_object_clause is called when exiting the audit_schema_object_clause production.
	ExitAudit_schema_object_clause(c *Audit_schema_object_clauseContext)

	// ExitSql_operation is called when exiting the sql_operation production.
	ExitSql_operation(c *Sql_operationContext)

	// ExitAuditing_on_clause is called when exiting the auditing_on_clause production.
	ExitAuditing_on_clause(c *Auditing_on_clauseContext)

	// ExitModel_name is called when exiting the model_name production.
	ExitModel_name(c *Model_nameContext)

	// ExitObject_name is called when exiting the object_name production.
	ExitObject_name(c *Object_nameContext)

	// ExitProfile_name is called when exiting the profile_name production.
	ExitProfile_name(c *Profile_nameContext)

	// ExitSql_statement_shortcut is called when exiting the sql_statement_shortcut production.
	ExitSql_statement_shortcut(c *Sql_statement_shortcutContext)

	// ExitDrop_index is called when exiting the drop_index production.
	ExitDrop_index(c *Drop_indexContext)

	// ExitDisassociate_statistics is called when exiting the disassociate_statistics production.
	ExitDisassociate_statistics(c *Disassociate_statisticsContext)

	// ExitDrop_indextype is called when exiting the drop_indextype production.
	ExitDrop_indextype(c *Drop_indextypeContext)

	// ExitDrop_inmemory_join_group is called when exiting the drop_inmemory_join_group production.
	ExitDrop_inmemory_join_group(c *Drop_inmemory_join_groupContext)

	// ExitFlashback_table is called when exiting the flashback_table production.
	ExitFlashback_table(c *Flashback_tableContext)

	// ExitRestore_point is called when exiting the restore_point production.
	ExitRestore_point(c *Restore_pointContext)

	// ExitPurge_statement is called when exiting the purge_statement production.
	ExitPurge_statement(c *Purge_statementContext)

	// ExitNoaudit_statement is called when exiting the noaudit_statement production.
	ExitNoaudit_statement(c *Noaudit_statementContext)

	// ExitRename_object is called when exiting the rename_object production.
	ExitRename_object(c *Rename_objectContext)

	// ExitGrant_statement is called when exiting the grant_statement production.
	ExitGrant_statement(c *Grant_statementContext)

	// ExitContainer_clause is called when exiting the container_clause production.
	ExitContainer_clause(c *Container_clauseContext)

	// ExitRevoke_statement is called when exiting the revoke_statement production.
	ExitRevoke_statement(c *Revoke_statementContext)

	// ExitRevoke_system_privilege is called when exiting the revoke_system_privilege production.
	ExitRevoke_system_privilege(c *Revoke_system_privilegeContext)

	// ExitRevokee_clause is called when exiting the revokee_clause production.
	ExitRevokee_clause(c *Revokee_clauseContext)

	// ExitRevoke_object_privileges is called when exiting the revoke_object_privileges production.
	ExitRevoke_object_privileges(c *Revoke_object_privilegesContext)

	// ExitOn_object_clause is called when exiting the on_object_clause production.
	ExitOn_object_clause(c *On_object_clauseContext)

	// ExitRevoke_roles_from_programs is called when exiting the revoke_roles_from_programs production.
	ExitRevoke_roles_from_programs(c *Revoke_roles_from_programsContext)

	// ExitProgram_unit is called when exiting the program_unit production.
	ExitProgram_unit(c *Program_unitContext)

	// ExitCreate_dimension is called when exiting the create_dimension production.
	ExitCreate_dimension(c *Create_dimensionContext)

	// ExitCreate_directory is called when exiting the create_directory production.
	ExitCreate_directory(c *Create_directoryContext)

	// ExitDirectory_name is called when exiting the directory_name production.
	ExitDirectory_name(c *Directory_nameContext)

	// ExitDirectory_path is called when exiting the directory_path production.
	ExitDirectory_path(c *Directory_pathContext)

	// ExitCreate_inmemory_join_group is called when exiting the create_inmemory_join_group production.
	ExitCreate_inmemory_join_group(c *Create_inmemory_join_groupContext)

	// ExitDrop_hierarchy is called when exiting the drop_hierarchy production.
	ExitDrop_hierarchy(c *Drop_hierarchyContext)

	// ExitAlter_library is called when exiting the alter_library production.
	ExitAlter_library(c *Alter_libraryContext)

	// ExitDrop_java is called when exiting the drop_java production.
	ExitDrop_java(c *Drop_javaContext)

	// ExitDrop_library is called when exiting the drop_library production.
	ExitDrop_library(c *Drop_libraryContext)

	// ExitCreate_java is called when exiting the create_java production.
	ExitCreate_java(c *Create_javaContext)

	// ExitCreate_library is called when exiting the create_library production.
	ExitCreate_library(c *Create_libraryContext)

	// ExitPlsql_library_source is called when exiting the plsql_library_source production.
	ExitPlsql_library_source(c *Plsql_library_sourceContext)

	// ExitCredential_name is called when exiting the credential_name production.
	ExitCredential_name(c *Credential_nameContext)

	// ExitLibrary_editionable is called when exiting the library_editionable production.
	ExitLibrary_editionable(c *Library_editionableContext)

	// ExitLibrary_debug is called when exiting the library_debug production.
	ExitLibrary_debug(c *Library_debugContext)

	// ExitCompiler_parameters_clause is called when exiting the compiler_parameters_clause production.
	ExitCompiler_parameters_clause(c *Compiler_parameters_clauseContext)

	// ExitParameter_value is called when exiting the parameter_value production.
	ExitParameter_value(c *Parameter_valueContext)

	// ExitLibrary_name is called when exiting the library_name production.
	ExitLibrary_name(c *Library_nameContext)

	// ExitAlter_dimension is called when exiting the alter_dimension production.
	ExitAlter_dimension(c *Alter_dimensionContext)

	// ExitLevel_clause is called when exiting the level_clause production.
	ExitLevel_clause(c *Level_clauseContext)

	// ExitHierarchy_clause is called when exiting the hierarchy_clause production.
	ExitHierarchy_clause(c *Hierarchy_clauseContext)

	// ExitDimension_join_clause is called when exiting the dimension_join_clause production.
	ExitDimension_join_clause(c *Dimension_join_clauseContext)

	// ExitAttribute_clause is called when exiting the attribute_clause production.
	ExitAttribute_clause(c *Attribute_clauseContext)

	// ExitExtended_attribute_clause is called when exiting the extended_attribute_clause production.
	ExitExtended_attribute_clause(c *Extended_attribute_clauseContext)

	// ExitColumn_one_or_more_sub_clause is called when exiting the column_one_or_more_sub_clause production.
	ExitColumn_one_or_more_sub_clause(c *Column_one_or_more_sub_clauseContext)

	// ExitAlter_view is called when exiting the alter_view production.
	ExitAlter_view(c *Alter_viewContext)

	// ExitAlter_view_editionable is called when exiting the alter_view_editionable production.
	ExitAlter_view_editionable(c *Alter_view_editionableContext)

	// ExitCreate_view is called when exiting the create_view production.
	ExitCreate_view(c *Create_viewContext)

	// ExitEditioning_clause is called when exiting the editioning_clause production.
	ExitEditioning_clause(c *Editioning_clauseContext)

	// ExitView_options is called when exiting the view_options production.
	ExitView_options(c *View_optionsContext)

	// ExitView_alias_constraint is called when exiting the view_alias_constraint production.
	ExitView_alias_constraint(c *View_alias_constraintContext)

	// ExitObject_view_clause is called when exiting the object_view_clause production.
	ExitObject_view_clause(c *Object_view_clauseContext)

	// ExitInline_constraint is called when exiting the inline_constraint production.
	ExitInline_constraint(c *Inline_constraintContext)

	// ExitInline_ref_constraint is called when exiting the inline_ref_constraint production.
	ExitInline_ref_constraint(c *Inline_ref_constraintContext)

	// ExitOut_of_line_ref_constraint is called when exiting the out_of_line_ref_constraint production.
	ExitOut_of_line_ref_constraint(c *Out_of_line_ref_constraintContext)

	// ExitOut_of_line_constraint is called when exiting the out_of_line_constraint production.
	ExitOut_of_line_constraint(c *Out_of_line_constraintContext)

	// ExitConstraint_state is called when exiting the constraint_state production.
	ExitConstraint_state(c *Constraint_stateContext)

	// ExitXmltype_view_clause is called when exiting the xmltype_view_clause production.
	ExitXmltype_view_clause(c *Xmltype_view_clauseContext)

	// ExitXml_schema_spec is called when exiting the xml_schema_spec production.
	ExitXml_schema_spec(c *Xml_schema_specContext)

	// ExitXml_schema_url is called when exiting the xml_schema_url production.
	ExitXml_schema_url(c *Xml_schema_urlContext)

	// ExitElement is called when exiting the element production.
	ExitElement(c *ElementContext)

	// ExitAlter_tablespace is called when exiting the alter_tablespace production.
	ExitAlter_tablespace(c *Alter_tablespaceContext)

	// ExitDatafile_tempfile_clauses is called when exiting the datafile_tempfile_clauses production.
	ExitDatafile_tempfile_clauses(c *Datafile_tempfile_clausesContext)

	// ExitTablespace_logging_clauses is called when exiting the tablespace_logging_clauses production.
	ExitTablespace_logging_clauses(c *Tablespace_logging_clausesContext)

	// ExitTablespace_group_clause is called when exiting the tablespace_group_clause production.
	ExitTablespace_group_clause(c *Tablespace_group_clauseContext)

	// ExitTablespace_group_name is called when exiting the tablespace_group_name production.
	ExitTablespace_group_name(c *Tablespace_group_nameContext)

	// ExitTablespace_state_clauses is called when exiting the tablespace_state_clauses production.
	ExitTablespace_state_clauses(c *Tablespace_state_clausesContext)

	// ExitFlashback_mode_clause is called when exiting the flashback_mode_clause production.
	ExitFlashback_mode_clause(c *Flashback_mode_clauseContext)

	// ExitNew_tablespace_name is called when exiting the new_tablespace_name production.
	ExitNew_tablespace_name(c *New_tablespace_nameContext)

	// ExitCreate_tablespace is called when exiting the create_tablespace production.
	ExitCreate_tablespace(c *Create_tablespaceContext)

	// ExitPermanent_tablespace_clause is called when exiting the permanent_tablespace_clause production.
	ExitPermanent_tablespace_clause(c *Permanent_tablespace_clauseContext)

	// ExitTablespace_encryption_spec is called when exiting the tablespace_encryption_spec production.
	ExitTablespace_encryption_spec(c *Tablespace_encryption_specContext)

	// ExitLogging_clause is called when exiting the logging_clause production.
	ExitLogging_clause(c *Logging_clauseContext)

	// ExitExtent_management_clause is called when exiting the extent_management_clause production.
	ExitExtent_management_clause(c *Extent_management_clauseContext)

	// ExitSegment_management_clause is called when exiting the segment_management_clause production.
	ExitSegment_management_clause(c *Segment_management_clauseContext)

	// ExitTemporary_tablespace_clause is called when exiting the temporary_tablespace_clause production.
	ExitTemporary_tablespace_clause(c *Temporary_tablespace_clauseContext)

	// ExitUndo_tablespace_clause is called when exiting the undo_tablespace_clause production.
	ExitUndo_tablespace_clause(c *Undo_tablespace_clauseContext)

	// ExitTablespace_retention_clause is called when exiting the tablespace_retention_clause production.
	ExitTablespace_retention_clause(c *Tablespace_retention_clauseContext)

	// ExitCreate_tablespace_set is called when exiting the create_tablespace_set production.
	ExitCreate_tablespace_set(c *Create_tablespace_setContext)

	// ExitPermanent_tablespace_attrs is called when exiting the permanent_tablespace_attrs production.
	ExitPermanent_tablespace_attrs(c *Permanent_tablespace_attrsContext)

	// ExitTablespace_encryption_clause is called when exiting the tablespace_encryption_clause production.
	ExitTablespace_encryption_clause(c *Tablespace_encryption_clauseContext)

	// ExitDefault_tablespace_params is called when exiting the default_tablespace_params production.
	ExitDefault_tablespace_params(c *Default_tablespace_paramsContext)

	// ExitDefault_table_compression is called when exiting the default_table_compression production.
	ExitDefault_table_compression(c *Default_table_compressionContext)

	// ExitLow_high is called when exiting the low_high production.
	ExitLow_high(c *Low_highContext)

	// ExitDefault_index_compression is called when exiting the default_index_compression production.
	ExitDefault_index_compression(c *Default_index_compressionContext)

	// ExitInmmemory_clause is called when exiting the inmmemory_clause production.
	ExitInmmemory_clause(c *Inmmemory_clauseContext)

	// ExitDatafile_specification is called when exiting the datafile_specification production.
	ExitDatafile_specification(c *Datafile_specificationContext)

	// ExitTempfile_specification is called when exiting the tempfile_specification production.
	ExitTempfile_specification(c *Tempfile_specificationContext)

	// ExitDatafile_tempfile_spec is called when exiting the datafile_tempfile_spec production.
	ExitDatafile_tempfile_spec(c *Datafile_tempfile_specContext)

	// ExitRedo_log_file_spec is called when exiting the redo_log_file_spec production.
	ExitRedo_log_file_spec(c *Redo_log_file_specContext)

	// ExitAutoextend_clause is called when exiting the autoextend_clause production.
	ExitAutoextend_clause(c *Autoextend_clauseContext)

	// ExitMaxsize_clause is called when exiting the maxsize_clause production.
	ExitMaxsize_clause(c *Maxsize_clauseContext)

	// ExitBuild_clause is called when exiting the build_clause production.
	ExitBuild_clause(c *Build_clauseContext)

	// ExitParallel_clause is called when exiting the parallel_clause production.
	ExitParallel_clause(c *Parallel_clauseContext)

	// ExitParallel_instances_clause is called when exiting the parallel_instances_clause production.
	ExitParallel_instances_clause(c *Parallel_instances_clauseContext)

	// ExitAlter_materialized_view is called when exiting the alter_materialized_view production.
	ExitAlter_materialized_view(c *Alter_materialized_viewContext)

	// ExitAlter_mv_option1 is called when exiting the alter_mv_option1 production.
	ExitAlter_mv_option1(c *Alter_mv_option1Context)

	// ExitAlter_mv_refresh is called when exiting the alter_mv_refresh production.
	ExitAlter_mv_refresh(c *Alter_mv_refreshContext)

	// ExitRollback_segment is called when exiting the rollback_segment production.
	ExitRollback_segment(c *Rollback_segmentContext)

	// ExitModify_mv_column_clause is called when exiting the modify_mv_column_clause production.
	ExitModify_mv_column_clause(c *Modify_mv_column_clauseContext)

	// ExitAlter_materialized_view_log is called when exiting the alter_materialized_view_log production.
	ExitAlter_materialized_view_log(c *Alter_materialized_view_logContext)

	// ExitAdd_mv_log_column_clause is called when exiting the add_mv_log_column_clause production.
	ExitAdd_mv_log_column_clause(c *Add_mv_log_column_clauseContext)

	// ExitMove_mv_log_clause is called when exiting the move_mv_log_clause production.
	ExitMove_mv_log_clause(c *Move_mv_log_clauseContext)

	// ExitMv_log_augmentation is called when exiting the mv_log_augmentation production.
	ExitMv_log_augmentation(c *Mv_log_augmentationContext)

	// ExitCreate_materialized_view_log is called when exiting the create_materialized_view_log production.
	ExitCreate_materialized_view_log(c *Create_materialized_view_logContext)

	// ExitNew_values_clause is called when exiting the new_values_clause production.
	ExitNew_values_clause(c *New_values_clauseContext)

	// ExitMv_log_purge_clause is called when exiting the mv_log_purge_clause production.
	ExitMv_log_purge_clause(c *Mv_log_purge_clauseContext)

	// ExitCreate_materialized_zonemap is called when exiting the create_materialized_zonemap production.
	ExitCreate_materialized_zonemap(c *Create_materialized_zonemapContext)

	// ExitAlter_materialized_zonemap is called when exiting the alter_materialized_zonemap production.
	ExitAlter_materialized_zonemap(c *Alter_materialized_zonemapContext)

	// ExitDrop_materialized_zonemap is called when exiting the drop_materialized_zonemap production.
	ExitDrop_materialized_zonemap(c *Drop_materialized_zonemapContext)

	// ExitZonemap_refresh_clause is called when exiting the zonemap_refresh_clause production.
	ExitZonemap_refresh_clause(c *Zonemap_refresh_clauseContext)

	// ExitZonemap_attributes is called when exiting the zonemap_attributes production.
	ExitZonemap_attributes(c *Zonemap_attributesContext)

	// ExitZonemap_name is called when exiting the zonemap_name production.
	ExitZonemap_name(c *Zonemap_nameContext)

	// ExitOperator_name is called when exiting the operator_name production.
	ExitOperator_name(c *Operator_nameContext)

	// ExitOperator_function_name is called when exiting the operator_function_name production.
	ExitOperator_function_name(c *Operator_function_nameContext)

	// ExitCreate_zonemap_on_table is called when exiting the create_zonemap_on_table production.
	ExitCreate_zonemap_on_table(c *Create_zonemap_on_tableContext)

	// ExitCreate_zonemap_as_subquery is called when exiting the create_zonemap_as_subquery production.
	ExitCreate_zonemap_as_subquery(c *Create_zonemap_as_subqueryContext)

	// ExitAlter_operator is called when exiting the alter_operator production.
	ExitAlter_operator(c *Alter_operatorContext)

	// ExitDrop_operator is called when exiting the drop_operator production.
	ExitDrop_operator(c *Drop_operatorContext)

	// ExitCreate_operator is called when exiting the create_operator production.
	ExitCreate_operator(c *Create_operatorContext)

	// ExitBinding_clause is called when exiting the binding_clause production.
	ExitBinding_clause(c *Binding_clauseContext)

	// ExitAdd_binding_clause is called when exiting the add_binding_clause production.
	ExitAdd_binding_clause(c *Add_binding_clauseContext)

	// ExitImplementation_clause is called when exiting the implementation_clause production.
	ExitImplementation_clause(c *Implementation_clauseContext)

	// ExitPrimary_operator_list is called when exiting the primary_operator_list production.
	ExitPrimary_operator_list(c *Primary_operator_listContext)

	// ExitPrimary_operator_item is called when exiting the primary_operator_item production.
	ExitPrimary_operator_item(c *Primary_operator_itemContext)

	// ExitOperator_context_clause is called when exiting the operator_context_clause production.
	ExitOperator_context_clause(c *Operator_context_clauseContext)

	// ExitUsing_function_clause is called when exiting the using_function_clause production.
	ExitUsing_function_clause(c *Using_function_clauseContext)

	// ExitDrop_binding_clause is called when exiting the drop_binding_clause production.
	ExitDrop_binding_clause(c *Drop_binding_clauseContext)

	// ExitCreate_materialized_view is called when exiting the create_materialized_view production.
	ExitCreate_materialized_view(c *Create_materialized_viewContext)

	// ExitScoped_table_ref_constraint is called when exiting the scoped_table_ref_constraint production.
	ExitScoped_table_ref_constraint(c *Scoped_table_ref_constraintContext)

	// ExitMv_column_alias is called when exiting the mv_column_alias production.
	ExitMv_column_alias(c *Mv_column_aliasContext)

	// ExitCreate_mv_refresh is called when exiting the create_mv_refresh production.
	ExitCreate_mv_refresh(c *Create_mv_refreshContext)

	// ExitQuery_rewrite_clause is called when exiting the query_rewrite_clause production.
	ExitQuery_rewrite_clause(c *Query_rewrite_clauseContext)

	// ExitUnusable_editions_clause is called when exiting the unusable_editions_clause production.
	ExitUnusable_editions_clause(c *Unusable_editions_clauseContext)

	// ExitDrop_materialized_view is called when exiting the drop_materialized_view production.
	ExitDrop_materialized_view(c *Drop_materialized_viewContext)

	// ExitDrop_materialized_view_log is called when exiting the drop_materialized_view_log production.
	ExitDrop_materialized_view_log(c *Drop_materialized_view_logContext)

	// ExitCreate_context is called when exiting the create_context production.
	ExitCreate_context(c *Create_contextContext)

	// ExitOracle_namespace is called when exiting the oracle_namespace production.
	ExitOracle_namespace(c *Oracle_namespaceContext)

	// ExitCreate_cluster is called when exiting the create_cluster production.
	ExitCreate_cluster(c *Create_clusterContext)

	// ExitCreate_profile is called when exiting the create_profile production.
	ExitCreate_profile(c *Create_profileContext)

	// ExitResource_parameters is called when exiting the resource_parameters production.
	ExitResource_parameters(c *Resource_parametersContext)

	// ExitPassword_parameters is called when exiting the password_parameters production.
	ExitPassword_parameters(c *Password_parametersContext)

	// ExitCreate_lockdown_profile is called when exiting the create_lockdown_profile production.
	ExitCreate_lockdown_profile(c *Create_lockdown_profileContext)

	// ExitStatic_base_profile is called when exiting the static_base_profile production.
	ExitStatic_base_profile(c *Static_base_profileContext)

	// ExitDynamic_base_profile is called when exiting the dynamic_base_profile production.
	ExitDynamic_base_profile(c *Dynamic_base_profileContext)

	// ExitCreate_outline is called when exiting the create_outline production.
	ExitCreate_outline(c *Create_outlineContext)

	// ExitCreate_restore_point is called when exiting the create_restore_point production.
	ExitCreate_restore_point(c *Create_restore_pointContext)

	// ExitCreate_role is called when exiting the create_role production.
	ExitCreate_role(c *Create_roleContext)

	// ExitCreate_table is called when exiting the create_table production.
	ExitCreate_table(c *Create_tableContext)

	// ExitXmltype_table is called when exiting the xmltype_table production.
	ExitXmltype_table(c *Xmltype_tableContext)

	// ExitXmltype_virtual_columns is called when exiting the xmltype_virtual_columns production.
	ExitXmltype_virtual_columns(c *Xmltype_virtual_columnsContext)

	// ExitXmltype_column_properties is called when exiting the xmltype_column_properties production.
	ExitXmltype_column_properties(c *Xmltype_column_propertiesContext)

	// ExitXmltype_storage is called when exiting the xmltype_storage production.
	ExitXmltype_storage(c *Xmltype_storageContext)

	// ExitXmlschema_spec is called when exiting the xmlschema_spec production.
	ExitXmlschema_spec(c *Xmlschema_specContext)

	// ExitObject_table is called when exiting the object_table production.
	ExitObject_table(c *Object_tableContext)

	// ExitObject_type is called when exiting the object_type production.
	ExitObject_type(c *Object_typeContext)

	// ExitOid_index_clause is called when exiting the oid_index_clause production.
	ExitOid_index_clause(c *Oid_index_clauseContext)

	// ExitOid_clause is called when exiting the oid_clause production.
	ExitOid_clause(c *Oid_clauseContext)

	// ExitObject_properties is called when exiting the object_properties production.
	ExitObject_properties(c *Object_propertiesContext)

	// ExitObject_table_substitution is called when exiting the object_table_substitution production.
	ExitObject_table_substitution(c *Object_table_substitutionContext)

	// ExitRelational_table is called when exiting the relational_table production.
	ExitRelational_table(c *Relational_tableContext)

	// ExitRelational_table_properties is called when exiting the relational_table_properties production.
	ExitRelational_table_properties(c *Relational_table_propertiesContext)

	// ExitRelational_table_property is called when exiting the relational_table_property production.
	ExitRelational_table_property(c *Relational_table_propertyContext)

	// ExitImmutable_table_clauses is called when exiting the immutable_table_clauses production.
	ExitImmutable_table_clauses(c *Immutable_table_clausesContext)

	// ExitImmutable_table_no_drop_clause is called when exiting the immutable_table_no_drop_clause production.
	ExitImmutable_table_no_drop_clause(c *Immutable_table_no_drop_clauseContext)

	// ExitImmutable_table_no_delete_clause is called when exiting the immutable_table_no_delete_clause production.
	ExitImmutable_table_no_delete_clause(c *Immutable_table_no_delete_clauseContext)

	// ExitBlockchain_table_clauses is called when exiting the blockchain_table_clauses production.
	ExitBlockchain_table_clauses(c *Blockchain_table_clausesContext)

	// ExitBlockchain_drop_table_clause is called when exiting the blockchain_drop_table_clause production.
	ExitBlockchain_drop_table_clause(c *Blockchain_drop_table_clauseContext)

	// ExitBlockchain_row_retention_clause is called when exiting the blockchain_row_retention_clause production.
	ExitBlockchain_row_retention_clause(c *Blockchain_row_retention_clauseContext)

	// ExitBlockchain_hash_and_data_format_clause is called when exiting the blockchain_hash_and_data_format_clause production.
	ExitBlockchain_hash_and_data_format_clause(c *Blockchain_hash_and_data_format_clauseContext)

	// ExitCollation_name is called when exiting the collation_name production.
	ExitCollation_name(c *Collation_nameContext)

	// ExitTable_properties is called when exiting the table_properties production.
	ExitTable_properties(c *Table_propertiesContext)

	// ExitRead_only_clause is called when exiting the read_only_clause production.
	ExitRead_only_clause(c *Read_only_clauseContext)

	// ExitIndexing_clause is called when exiting the indexing_clause production.
	ExitIndexing_clause(c *Indexing_clauseContext)

	// ExitAttribute_clustering_clause is called when exiting the attribute_clustering_clause production.
	ExitAttribute_clustering_clause(c *Attribute_clustering_clauseContext)

	// ExitClustering_join is called when exiting the clustering_join production.
	ExitClustering_join(c *Clustering_joinContext)

	// ExitClustering_join_item is called when exiting the clustering_join_item production.
	ExitClustering_join_item(c *Clustering_join_itemContext)

	// ExitEquijoin_condition is called when exiting the equijoin_condition production.
	ExitEquijoin_condition(c *Equijoin_conditionContext)

	// ExitCluster_clause is called when exiting the cluster_clause production.
	ExitCluster_clause(c *Cluster_clauseContext)

	// ExitClustering_columns is called when exiting the clustering_columns production.
	ExitClustering_columns(c *Clustering_columnsContext)

	// ExitClustering_column_group is called when exiting the clustering_column_group production.
	ExitClustering_column_group(c *Clustering_column_groupContext)

	// ExitYes_no is called when exiting the yes_no production.
	ExitYes_no(c *Yes_noContext)

	// ExitZonemap_clause is called when exiting the zonemap_clause production.
	ExitZonemap_clause(c *Zonemap_clauseContext)

	// ExitLogical_replication_clause is called when exiting the logical_replication_clause production.
	ExitLogical_replication_clause(c *Logical_replication_clauseContext)

	// ExitTable_name is called when exiting the table_name production.
	ExitTable_name(c *Table_nameContext)

	// ExitRelational_property is called when exiting the relational_property production.
	ExitRelational_property(c *Relational_propertyContext)

	// ExitTable_partitioning_clauses is called when exiting the table_partitioning_clauses production.
	ExitTable_partitioning_clauses(c *Table_partitioning_clausesContext)

	// ExitRange_partitions is called when exiting the range_partitions production.
	ExitRange_partitions(c *Range_partitionsContext)

	// ExitList_partitions is called when exiting the list_partitions production.
	ExitList_partitions(c *List_partitionsContext)

	// ExitHash_partitions is called when exiting the hash_partitions production.
	ExitHash_partitions(c *Hash_partitionsContext)

	// ExitIndividual_hash_partitions is called when exiting the individual_hash_partitions production.
	ExitIndividual_hash_partitions(c *Individual_hash_partitionsContext)

	// ExitHash_partitions_by_quantity is called when exiting the hash_partitions_by_quantity production.
	ExitHash_partitions_by_quantity(c *Hash_partitions_by_quantityContext)

	// ExitHash_partition_quantity is called when exiting the hash_partition_quantity production.
	ExitHash_partition_quantity(c *Hash_partition_quantityContext)

	// ExitComposite_range_partitions is called when exiting the composite_range_partitions production.
	ExitComposite_range_partitions(c *Composite_range_partitionsContext)

	// ExitComposite_list_partitions is called when exiting the composite_list_partitions production.
	ExitComposite_list_partitions(c *Composite_list_partitionsContext)

	// ExitComposite_hash_partitions is called when exiting the composite_hash_partitions production.
	ExitComposite_hash_partitions(c *Composite_hash_partitionsContext)

	// ExitReference_partitioning is called when exiting the reference_partitioning production.
	ExitReference_partitioning(c *Reference_partitioningContext)

	// ExitReference_partition_desc is called when exiting the reference_partition_desc production.
	ExitReference_partition_desc(c *Reference_partition_descContext)

	// ExitSystem_partitioning is called when exiting the system_partitioning production.
	ExitSystem_partitioning(c *System_partitioningContext)

	// ExitRange_partition_desc is called when exiting the range_partition_desc production.
	ExitRange_partition_desc(c *Range_partition_descContext)

	// ExitList_partition_desc is called when exiting the list_partition_desc production.
	ExitList_partition_desc(c *List_partition_descContext)

	// ExitSubpartition_template is called when exiting the subpartition_template production.
	ExitSubpartition_template(c *Subpartition_templateContext)

	// ExitHash_subpartition_quantity is called when exiting the hash_subpartition_quantity production.
	ExitHash_subpartition_quantity(c *Hash_subpartition_quantityContext)

	// ExitSubpartition_by_range is called when exiting the subpartition_by_range production.
	ExitSubpartition_by_range(c *Subpartition_by_rangeContext)

	// ExitSubpartition_by_list is called when exiting the subpartition_by_list production.
	ExitSubpartition_by_list(c *Subpartition_by_listContext)

	// ExitSubpartition_by_hash is called when exiting the subpartition_by_hash production.
	ExitSubpartition_by_hash(c *Subpartition_by_hashContext)

	// ExitSubpartition_name is called when exiting the subpartition_name production.
	ExitSubpartition_name(c *Subpartition_nameContext)

	// ExitRange_subpartition_desc is called when exiting the range_subpartition_desc production.
	ExitRange_subpartition_desc(c *Range_subpartition_descContext)

	// ExitList_subpartition_desc is called when exiting the list_subpartition_desc production.
	ExitList_subpartition_desc(c *List_subpartition_descContext)

	// ExitIndividual_hash_subparts is called when exiting the individual_hash_subparts production.
	ExitIndividual_hash_subparts(c *Individual_hash_subpartsContext)

	// ExitHash_subparts_by_quantity is called when exiting the hash_subparts_by_quantity production.
	ExitHash_subparts_by_quantity(c *Hash_subparts_by_quantityContext)

	// ExitRange_values_clause is called when exiting the range_values_clause production.
	ExitRange_values_clause(c *Range_values_clauseContext)

	// ExitRange_values_list is called when exiting the range_values_list production.
	ExitRange_values_list(c *Range_values_listContext)

	// ExitList_values_clause is called when exiting the list_values_clause production.
	ExitList_values_clause(c *List_values_clauseContext)

	// ExitTable_partition_description is called when exiting the table_partition_description production.
	ExitTable_partition_description(c *Table_partition_descriptionContext)

	// ExitPartitioning_storage_clause is called when exiting the partitioning_storage_clause production.
	ExitPartitioning_storage_clause(c *Partitioning_storage_clauseContext)

	// ExitLob_partitioning_storage is called when exiting the lob_partitioning_storage production.
	ExitLob_partitioning_storage(c *Lob_partitioning_storageContext)

	// ExitSize_clause is called when exiting the size_clause production.
	ExitSize_clause(c *Size_clauseContext)

	// ExitTable_compression is called when exiting the table_compression production.
	ExitTable_compression(c *Table_compressionContext)

	// ExitInmemory_table_clause is called when exiting the inmemory_table_clause production.
	ExitInmemory_table_clause(c *Inmemory_table_clauseContext)

	// ExitInmemory_attributes is called when exiting the inmemory_attributes production.
	ExitInmemory_attributes(c *Inmemory_attributesContext)

	// ExitInmemory_memcompress is called when exiting the inmemory_memcompress production.
	ExitInmemory_memcompress(c *Inmemory_memcompressContext)

	// ExitInmemory_priority is called when exiting the inmemory_priority production.
	ExitInmemory_priority(c *Inmemory_priorityContext)

	// ExitInmemory_distribute is called when exiting the inmemory_distribute production.
	ExitInmemory_distribute(c *Inmemory_distributeContext)

	// ExitInmemory_duplicate is called when exiting the inmemory_duplicate production.
	ExitInmemory_duplicate(c *Inmemory_duplicateContext)

	// ExitInmemory_column_clause is called when exiting the inmemory_column_clause production.
	ExitInmemory_column_clause(c *Inmemory_column_clauseContext)

	// ExitPhysical_attributes_clause is called when exiting the physical_attributes_clause production.
	ExitPhysical_attributes_clause(c *Physical_attributes_clauseContext)

	// ExitStorage_clause is called when exiting the storage_clause production.
	ExitStorage_clause(c *Storage_clauseContext)

	// ExitDeferred_segment_creation is called when exiting the deferred_segment_creation production.
	ExitDeferred_segment_creation(c *Deferred_segment_creationContext)

	// ExitSegment_attributes_clause is called when exiting the segment_attributes_clause production.
	ExitSegment_attributes_clause(c *Segment_attributes_clauseContext)

	// ExitPhysical_properties is called when exiting the physical_properties production.
	ExitPhysical_properties(c *Physical_propertiesContext)

	// ExitIlm_clause is called when exiting the ilm_clause production.
	ExitIlm_clause(c *Ilm_clauseContext)

	// ExitIlm_policy_clause is called when exiting the ilm_policy_clause production.
	ExitIlm_policy_clause(c *Ilm_policy_clauseContext)

	// ExitIlm_compression_policy is called when exiting the ilm_compression_policy production.
	ExitIlm_compression_policy(c *Ilm_compression_policyContext)

	// ExitIlm_tiering_policy is called when exiting the ilm_tiering_policy production.
	ExitIlm_tiering_policy(c *Ilm_tiering_policyContext)

	// ExitIlm_after_on is called when exiting the ilm_after_on production.
	ExitIlm_after_on(c *Ilm_after_onContext)

	// ExitSegment_group is called when exiting the segment_group production.
	ExitSegment_group(c *Segment_groupContext)

	// ExitIlm_inmemory_policy is called when exiting the ilm_inmemory_policy production.
	ExitIlm_inmemory_policy(c *Ilm_inmemory_policyContext)

	// ExitIlm_time_period is called when exiting the ilm_time_period production.
	ExitIlm_time_period(c *Ilm_time_periodContext)

	// ExitHeap_org_table_clause is called when exiting the heap_org_table_clause production.
	ExitHeap_org_table_clause(c *Heap_org_table_clauseContext)

	// ExitExternal_table_clause is called when exiting the external_table_clause production.
	ExitExternal_table_clause(c *External_table_clauseContext)

	// ExitAccess_driver_type is called when exiting the access_driver_type production.
	ExitAccess_driver_type(c *Access_driver_typeContext)

	// ExitExternal_table_data_props is called when exiting the external_table_data_props production.
	ExitExternal_table_data_props(c *External_table_data_propsContext)

	// ExitExternal_table_data_format is called when exiting the external_table_data_format production.
	ExitExternal_table_data_format(c *External_table_data_formatContext)

	// ExitExternal_table_transform is called when exiting the external_table_transform production.
	ExitExternal_table_transform(c *External_table_transformContext)

	// ExitExternal_table_field is called when exiting the external_table_field production.
	ExitExternal_table_field(c *External_table_fieldContext)

	// ExitExternal_table_field_list is called when exiting the external_table_field_list production.
	ExitExternal_table_field_list(c *External_table_field_listContext)

	// ExitExternal_table_fields_clause is called when exiting the external_table_fields_clause production.
	ExitExternal_table_fields_clause(c *External_table_fields_clauseContext)

	// ExitExternal_table_position_clause is called when exiting the external_table_position_clause production.
	ExitExternal_table_position_clause(c *External_table_position_clauseContext)

	// ExitExternal_table_datatype_clause is called when exiting the external_table_datatype_clause production.
	ExitExternal_table_datatype_clause(c *External_table_datatype_clauseContext)

	// ExitExternal_table_delimit_clause is called when exiting the external_table_delimit_clause production.
	ExitExternal_table_delimit_clause(c *External_table_delimit_clauseContext)

	// ExitExternal_table_trim_clause is called when exiting the external_table_trim_clause production.
	ExitExternal_table_trim_clause(c *External_table_trim_clauseContext)

	// ExitExternal_table_date_format_clause is called when exiting the external_table_date_format_clause production.
	ExitExternal_table_date_format_clause(c *External_table_date_format_clauseContext)

	// ExitExternal_table_init_clause is called when exiting the external_table_init_clause production.
	ExitExternal_table_init_clause(c *External_table_init_clauseContext)

	// ExitExternal_table_condition_clause is called when exiting the external_table_condition_clause production.
	ExitExternal_table_condition_clause(c *External_table_condition_clauseContext)

	// ExitExternal_table_lls_clause is called when exiting the external_table_lls_clause production.
	ExitExternal_table_lls_clause(c *External_table_lls_clauseContext)

	// ExitExternal_table_records is called when exiting the external_table_records production.
	ExitExternal_table_records(c *External_table_recordsContext)

	// ExitExternal_table_record_options_clause is called when exiting the external_table_record_options_clause production.
	ExitExternal_table_record_options_clause(c *External_table_record_options_clauseContext)

	// ExitExternal_table_output_files is called when exiting the external_table_output_files production.
	ExitExternal_table_output_files(c *External_table_output_filesContext)

	// ExitExternal_table_fields is called when exiting the external_table_fields production.
	ExitExternal_table_fields(c *External_table_fieldsContext)

	// ExitExternal_table_datapump is called when exiting the external_table_datapump production.
	ExitExternal_table_datapump(c *External_table_datapumpContext)

	// ExitExternal_table_hive is called when exiting the external_table_hive production.
	ExitExternal_table_hive(c *External_table_hiveContext)

	// ExitExternal_table_hive_parameter_map is called when exiting the external_table_hive_parameter_map production.
	ExitExternal_table_hive_parameter_map(c *External_table_hive_parameter_mapContext)

	// ExitExternal_table_hive_parameter_map_entry is called when exiting the external_table_hive_parameter_map_entry production.
	ExitExternal_table_hive_parameter_map_entry(c *External_table_hive_parameter_map_entryContext)

	// ExitExternal_table_directory is called when exiting the external_table_directory production.
	ExitExternal_table_directory(c *External_table_directoryContext)

	// ExitRow_movement_clause is called when exiting the row_movement_clause production.
	ExitRow_movement_clause(c *Row_movement_clauseContext)

	// ExitFlashback_archive_clause is called when exiting the flashback_archive_clause production.
	ExitFlashback_archive_clause(c *Flashback_archive_clauseContext)

	// ExitLog_grp is called when exiting the log_grp production.
	ExitLog_grp(c *Log_grpContext)

	// ExitSupplemental_table_logging is called when exiting the supplemental_table_logging production.
	ExitSupplemental_table_logging(c *Supplemental_table_loggingContext)

	// ExitSupplemental_log_grp_clause is called when exiting the supplemental_log_grp_clause production.
	ExitSupplemental_log_grp_clause(c *Supplemental_log_grp_clauseContext)

	// ExitSupplemental_id_key_clause is called when exiting the supplemental_id_key_clause production.
	ExitSupplemental_id_key_clause(c *Supplemental_id_key_clauseContext)

	// ExitAllocate_extent_clause is called when exiting the allocate_extent_clause production.
	ExitAllocate_extent_clause(c *Allocate_extent_clauseContext)

	// ExitDeallocate_unused_clause is called when exiting the deallocate_unused_clause production.
	ExitDeallocate_unused_clause(c *Deallocate_unused_clauseContext)

	// ExitShrink_clause is called when exiting the shrink_clause production.
	ExitShrink_clause(c *Shrink_clauseContext)

	// ExitRecords_per_block_clause is called when exiting the records_per_block_clause production.
	ExitRecords_per_block_clause(c *Records_per_block_clauseContext)

	// ExitUpgrade_table_clause is called when exiting the upgrade_table_clause production.
	ExitUpgrade_table_clause(c *Upgrade_table_clauseContext)

	// ExitTruncate_table is called when exiting the truncate_table production.
	ExitTruncate_table(c *Truncate_tableContext)

	// ExitDrop_table is called when exiting the drop_table production.
	ExitDrop_table(c *Drop_tableContext)

	// ExitDrop_tablespace is called when exiting the drop_tablespace production.
	ExitDrop_tablespace(c *Drop_tablespaceContext)

	// ExitDrop_tablespace_set is called when exiting the drop_tablespace_set production.
	ExitDrop_tablespace_set(c *Drop_tablespace_setContext)

	// ExitIncluding_contents_clause is called when exiting the including_contents_clause production.
	ExitIncluding_contents_clause(c *Including_contents_clauseContext)

	// ExitDrop_view is called when exiting the drop_view production.
	ExitDrop_view(c *Drop_viewContext)

	// ExitComment_on_column is called when exiting the comment_on_column production.
	ExitComment_on_column(c *Comment_on_columnContext)

	// ExitEnable_or_disable is called when exiting the enable_or_disable production.
	ExitEnable_or_disable(c *Enable_or_disableContext)

	// ExitAllow_or_disallow is called when exiting the allow_or_disallow production.
	ExitAllow_or_disallow(c *Allow_or_disallowContext)

	// ExitAlter_synonym is called when exiting the alter_synonym production.
	ExitAlter_synonym(c *Alter_synonymContext)

	// ExitCreate_synonym is called when exiting the create_synonym production.
	ExitCreate_synonym(c *Create_synonymContext)

	// ExitDrop_synonym is called when exiting the drop_synonym production.
	ExitDrop_synonym(c *Drop_synonymContext)

	// ExitCreate_spfile is called when exiting the create_spfile production.
	ExitCreate_spfile(c *Create_spfileContext)

	// ExitSpfile_name is called when exiting the spfile_name production.
	ExitSpfile_name(c *Spfile_nameContext)

	// ExitPfile_name is called when exiting the pfile_name production.
	ExitPfile_name(c *Pfile_nameContext)

	// ExitComment_on_table is called when exiting the comment_on_table production.
	ExitComment_on_table(c *Comment_on_tableContext)

	// ExitComment_on_materialized is called when exiting the comment_on_materialized production.
	ExitComment_on_materialized(c *Comment_on_materializedContext)

	// ExitAlter_analytic_view is called when exiting the alter_analytic_view production.
	ExitAlter_analytic_view(c *Alter_analytic_viewContext)

	// ExitAlter_add_cache_clause is called when exiting the alter_add_cache_clause production.
	ExitAlter_add_cache_clause(c *Alter_add_cache_clauseContext)

	// ExitLevels_item is called when exiting the levels_item production.
	ExitLevels_item(c *Levels_itemContext)

	// ExitMeasure_list is called when exiting the measure_list production.
	ExitMeasure_list(c *Measure_listContext)

	// ExitAlter_drop_cache_clause is called when exiting the alter_drop_cache_clause production.
	ExitAlter_drop_cache_clause(c *Alter_drop_cache_clauseContext)

	// ExitAlter_attribute_dimension is called when exiting the alter_attribute_dimension production.
	ExitAlter_attribute_dimension(c *Alter_attribute_dimensionContext)

	// ExitAlter_audit_policy is called when exiting the alter_audit_policy production.
	ExitAlter_audit_policy(c *Alter_audit_policyContext)

	// ExitAlter_cluster is called when exiting the alter_cluster production.
	ExitAlter_cluster(c *Alter_clusterContext)

	// ExitDrop_analytic_view is called when exiting the drop_analytic_view production.
	ExitDrop_analytic_view(c *Drop_analytic_viewContext)

	// ExitDrop_attribute_dimension is called when exiting the drop_attribute_dimension production.
	ExitDrop_attribute_dimension(c *Drop_attribute_dimensionContext)

	// ExitDrop_audit_policy is called when exiting the drop_audit_policy production.
	ExitDrop_audit_policy(c *Drop_audit_policyContext)

	// ExitDrop_flashback_archive is called when exiting the drop_flashback_archive production.
	ExitDrop_flashback_archive(c *Drop_flashback_archiveContext)

	// ExitDrop_cluster is called when exiting the drop_cluster production.
	ExitDrop_cluster(c *Drop_clusterContext)

	// ExitDrop_context is called when exiting the drop_context production.
	ExitDrop_context(c *Drop_contextContext)

	// ExitDrop_directory is called when exiting the drop_directory production.
	ExitDrop_directory(c *Drop_directoryContext)

	// ExitDrop_diskgroup is called when exiting the drop_diskgroup production.
	ExitDrop_diskgroup(c *Drop_diskgroupContext)

	// ExitDrop_edition is called when exiting the drop_edition production.
	ExitDrop_edition(c *Drop_editionContext)

	// ExitTruncate_cluster is called when exiting the truncate_cluster production.
	ExitTruncate_cluster(c *Truncate_clusterContext)

	// ExitCache_or_nocache is called when exiting the cache_or_nocache production.
	ExitCache_or_nocache(c *Cache_or_nocacheContext)

	// ExitDatabase_name is called when exiting the database_name production.
	ExitDatabase_name(c *Database_nameContext)

	// ExitAlter_database is called when exiting the alter_database production.
	ExitAlter_database(c *Alter_databaseContext)

	// ExitDatabase_clause is called when exiting the database_clause production.
	ExitDatabase_clause(c *Database_clauseContext)

	// ExitStartup_clauses is called when exiting the startup_clauses production.
	ExitStartup_clauses(c *Startup_clausesContext)

	// ExitResetlogs_or_noresetlogs is called when exiting the resetlogs_or_noresetlogs production.
	ExitResetlogs_or_noresetlogs(c *Resetlogs_or_noresetlogsContext)

	// ExitUpgrade_or_downgrade is called when exiting the upgrade_or_downgrade production.
	ExitUpgrade_or_downgrade(c *Upgrade_or_downgradeContext)

	// ExitRecovery_clauses is called when exiting the recovery_clauses production.
	ExitRecovery_clauses(c *Recovery_clausesContext)

	// ExitBegin_or_end is called when exiting the begin_or_end production.
	ExitBegin_or_end(c *Begin_or_endContext)

	// ExitGeneral_recovery is called when exiting the general_recovery production.
	ExitGeneral_recovery(c *General_recoveryContext)

	// ExitFull_database_recovery is called when exiting the full_database_recovery production.
	ExitFull_database_recovery(c *Full_database_recoveryContext)

	// ExitPartial_database_recovery is called when exiting the partial_database_recovery production.
	ExitPartial_database_recovery(c *Partial_database_recoveryContext)

	// ExitPartial_database_recovery_10g is called when exiting the partial_database_recovery_10g production.
	ExitPartial_database_recovery_10g(c *Partial_database_recovery_10gContext)

	// ExitManaged_standby_recovery is called when exiting the managed_standby_recovery production.
	ExitManaged_standby_recovery(c *Managed_standby_recoveryContext)

	// ExitDb_name is called when exiting the db_name production.
	ExitDb_name(c *Db_nameContext)

	// ExitDatabase_file_clauses is called when exiting the database_file_clauses production.
	ExitDatabase_file_clauses(c *Database_file_clausesContext)

	// ExitCreate_datafile_clause is called when exiting the create_datafile_clause production.
	ExitCreate_datafile_clause(c *Create_datafile_clauseContext)

	// ExitAlter_datafile_clause is called when exiting the alter_datafile_clause production.
	ExitAlter_datafile_clause(c *Alter_datafile_clauseContext)

	// ExitAlter_tempfile_clause is called when exiting the alter_tempfile_clause production.
	ExitAlter_tempfile_clause(c *Alter_tempfile_clauseContext)

	// ExitMove_datafile_clause is called when exiting the move_datafile_clause production.
	ExitMove_datafile_clause(c *Move_datafile_clauseContext)

	// ExitLogfile_clauses is called when exiting the logfile_clauses production.
	ExitLogfile_clauses(c *Logfile_clausesContext)

	// ExitAdd_logfile_clauses is called when exiting the add_logfile_clauses production.
	ExitAdd_logfile_clauses(c *Add_logfile_clausesContext)

	// ExitGroup_redo_logfile is called when exiting the group_redo_logfile production.
	ExitGroup_redo_logfile(c *Group_redo_logfileContext)

	// ExitDrop_logfile_clauses is called when exiting the drop_logfile_clauses production.
	ExitDrop_logfile_clauses(c *Drop_logfile_clausesContext)

	// ExitSwitch_logfile_clause is called when exiting the switch_logfile_clause production.
	ExitSwitch_logfile_clause(c *Switch_logfile_clauseContext)

	// ExitSupplemental_db_logging is called when exiting the supplemental_db_logging production.
	ExitSupplemental_db_logging(c *Supplemental_db_loggingContext)

	// ExitAdd_or_drop is called when exiting the add_or_drop production.
	ExitAdd_or_drop(c *Add_or_dropContext)

	// ExitSupplemental_plsql_clause is called when exiting the supplemental_plsql_clause production.
	ExitSupplemental_plsql_clause(c *Supplemental_plsql_clauseContext)

	// ExitLogfile_descriptor is called when exiting the logfile_descriptor production.
	ExitLogfile_descriptor(c *Logfile_descriptorContext)

	// ExitControlfile_clauses is called when exiting the controlfile_clauses production.
	ExitControlfile_clauses(c *Controlfile_clausesContext)

	// ExitTrace_file_clause is called when exiting the trace_file_clause production.
	ExitTrace_file_clause(c *Trace_file_clauseContext)

	// ExitStandby_database_clauses is called when exiting the standby_database_clauses production.
	ExitStandby_database_clauses(c *Standby_database_clausesContext)

	// ExitActivate_standby_db_clause is called when exiting the activate_standby_db_clause production.
	ExitActivate_standby_db_clause(c *Activate_standby_db_clauseContext)

	// ExitMaximize_standby_db_clause is called when exiting the maximize_standby_db_clause production.
	ExitMaximize_standby_db_clause(c *Maximize_standby_db_clauseContext)

	// ExitRegister_logfile_clause is called when exiting the register_logfile_clause production.
	ExitRegister_logfile_clause(c *Register_logfile_clauseContext)

	// ExitCommit_switchover_clause is called when exiting the commit_switchover_clause production.
	ExitCommit_switchover_clause(c *Commit_switchover_clauseContext)

	// ExitStart_standby_clause is called when exiting the start_standby_clause production.
	ExitStart_standby_clause(c *Start_standby_clauseContext)

	// ExitStop_standby_clause is called when exiting the stop_standby_clause production.
	ExitStop_standby_clause(c *Stop_standby_clauseContext)

	// ExitConvert_database_clause is called when exiting the convert_database_clause production.
	ExitConvert_database_clause(c *Convert_database_clauseContext)

	// ExitDefault_settings_clause is called when exiting the default_settings_clause production.
	ExitDefault_settings_clause(c *Default_settings_clauseContext)

	// ExitSet_time_zone_clause is called when exiting the set_time_zone_clause production.
	ExitSet_time_zone_clause(c *Set_time_zone_clauseContext)

	// ExitInstance_clauses is called when exiting the instance_clauses production.
	ExitInstance_clauses(c *Instance_clausesContext)

	// ExitSecurity_clause is called when exiting the security_clause production.
	ExitSecurity_clause(c *Security_clauseContext)

	// ExitDomain is called when exiting the domain production.
	ExitDomain(c *DomainContext)

	// ExitDatabase is called when exiting the database production.
	ExitDatabase(c *DatabaseContext)

	// ExitEdition_name is called when exiting the edition_name production.
	ExitEdition_name(c *Edition_nameContext)

	// ExitFilenumber is called when exiting the filenumber production.
	ExitFilenumber(c *FilenumberContext)

	// ExitFilename is called when exiting the filename production.
	ExitFilename(c *FilenameContext)

	// ExitPrepare_clause is called when exiting the prepare_clause production.
	ExitPrepare_clause(c *Prepare_clauseContext)

	// ExitDrop_mirror_clause is called when exiting the drop_mirror_clause production.
	ExitDrop_mirror_clause(c *Drop_mirror_clauseContext)

	// ExitLost_write_protection is called when exiting the lost_write_protection production.
	ExitLost_write_protection(c *Lost_write_protectionContext)

	// ExitCdb_fleet_clauses is called when exiting the cdb_fleet_clauses production.
	ExitCdb_fleet_clauses(c *Cdb_fleet_clausesContext)

	// ExitLead_cdb_clause is called when exiting the lead_cdb_clause production.
	ExitLead_cdb_clause(c *Lead_cdb_clauseContext)

	// ExitLead_cdb_uri_clause is called when exiting the lead_cdb_uri_clause production.
	ExitLead_cdb_uri_clause(c *Lead_cdb_uri_clauseContext)

	// ExitProperty_clauses is called when exiting the property_clauses production.
	ExitProperty_clauses(c *Property_clausesContext)

	// ExitReplay_upgrade_clauses is called when exiting the replay_upgrade_clauses production.
	ExitReplay_upgrade_clauses(c *Replay_upgrade_clausesContext)

	// ExitAlter_database_link is called when exiting the alter_database_link production.
	ExitAlter_database_link(c *Alter_database_linkContext)

	// ExitPassword_value is called when exiting the password_value production.
	ExitPassword_value(c *Password_valueContext)

	// ExitLink_authentication is called when exiting the link_authentication production.
	ExitLink_authentication(c *Link_authenticationContext)

	// ExitCreate_schema is called when exiting the create_schema production.
	ExitCreate_schema(c *Create_schemaContext)

	// ExitCreate_database is called when exiting the create_database production.
	ExitCreate_database(c *Create_databaseContext)

	// ExitDatabase_logging_clauses is called when exiting the database_logging_clauses production.
	ExitDatabase_logging_clauses(c *Database_logging_clausesContext)

	// ExitDatabase_logging_sub_clause is called when exiting the database_logging_sub_clause production.
	ExitDatabase_logging_sub_clause(c *Database_logging_sub_clauseContext)

	// ExitTablespace_clauses is called when exiting the tablespace_clauses production.
	ExitTablespace_clauses(c *Tablespace_clausesContext)

	// ExitEnable_pluggable_database is called when exiting the enable_pluggable_database production.
	ExitEnable_pluggable_database(c *Enable_pluggable_databaseContext)

	// ExitFile_name_convert is called when exiting the file_name_convert production.
	ExitFile_name_convert(c *File_name_convertContext)

	// ExitFilename_convert_sub_clause is called when exiting the filename_convert_sub_clause production.
	ExitFilename_convert_sub_clause(c *Filename_convert_sub_clauseContext)

	// ExitTablespace_datafile_clauses is called when exiting the tablespace_datafile_clauses production.
	ExitTablespace_datafile_clauses(c *Tablespace_datafile_clausesContext)

	// ExitUndo_mode_clause is called when exiting the undo_mode_clause production.
	ExitUndo_mode_clause(c *Undo_mode_clauseContext)

	// ExitDefault_tablespace is called when exiting the default_tablespace production.
	ExitDefault_tablespace(c *Default_tablespaceContext)

	// ExitDefault_temp_tablespace is called when exiting the default_temp_tablespace production.
	ExitDefault_temp_tablespace(c *Default_temp_tablespaceContext)

	// ExitUndo_tablespace is called when exiting the undo_tablespace production.
	ExitUndo_tablespace(c *Undo_tablespaceContext)

	// ExitDrop_database is called when exiting the drop_database production.
	ExitDrop_database(c *Drop_databaseContext)

	// ExitCreate_database_link is called when exiting the create_database_link production.
	ExitCreate_database_link(c *Create_database_linkContext)

	// ExitDrop_database_link is called when exiting the drop_database_link production.
	ExitDrop_database_link(c *Drop_database_linkContext)

	// ExitAlter_tablespace_set is called when exiting the alter_tablespace_set production.
	ExitAlter_tablespace_set(c *Alter_tablespace_setContext)

	// ExitAlter_tablespace_attrs is called when exiting the alter_tablespace_attrs production.
	ExitAlter_tablespace_attrs(c *Alter_tablespace_attrsContext)

	// ExitAlter_tablespace_encryption is called when exiting the alter_tablespace_encryption production.
	ExitAlter_tablespace_encryption(c *Alter_tablespace_encryptionContext)

	// ExitTs_file_name_convert is called when exiting the ts_file_name_convert production.
	ExitTs_file_name_convert(c *Ts_file_name_convertContext)

	// ExitAlter_role is called when exiting the alter_role production.
	ExitAlter_role(c *Alter_roleContext)

	// ExitRole_identified_clause is called when exiting the role_identified_clause production.
	ExitRole_identified_clause(c *Role_identified_clauseContext)

	// ExitAlter_table is called when exiting the alter_table production.
	ExitAlter_table(c *Alter_tableContext)

	// ExitMemoptimize_read_write_clause is called when exiting the memoptimize_read_write_clause production.
	ExitMemoptimize_read_write_clause(c *Memoptimize_read_write_clauseContext)

	// ExitAlter_table_properties is called when exiting the alter_table_properties production.
	ExitAlter_table_properties(c *Alter_table_propertiesContext)

	// ExitAlter_table_partitioning is called when exiting the alter_table_partitioning production.
	ExitAlter_table_partitioning(c *Alter_table_partitioningContext)

	// ExitAdd_table_partition is called when exiting the add_table_partition production.
	ExitAdd_table_partition(c *Add_table_partitionContext)

	// ExitDrop_table_partition is called when exiting the drop_table_partition production.
	ExitDrop_table_partition(c *Drop_table_partitionContext)

	// ExitMerge_table_partition is called when exiting the merge_table_partition production.
	ExitMerge_table_partition(c *Merge_table_partitionContext)

	// ExitModify_table_partition is called when exiting the modify_table_partition production.
	ExitModify_table_partition(c *Modify_table_partitionContext)

	// ExitSplit_table_partition is called when exiting the split_table_partition production.
	ExitSplit_table_partition(c *Split_table_partitionContext)

	// ExitTruncate_table_partition is called when exiting the truncate_table_partition production.
	ExitTruncate_table_partition(c *Truncate_table_partitionContext)

	// ExitExchange_table_partition is called when exiting the exchange_table_partition production.
	ExitExchange_table_partition(c *Exchange_table_partitionContext)

	// ExitCoalesce_table_partition is called when exiting the coalesce_table_partition production.
	ExitCoalesce_table_partition(c *Coalesce_table_partitionContext)

	// ExitAlter_interval_partition is called when exiting the alter_interval_partition production.
	ExitAlter_interval_partition(c *Alter_interval_partitionContext)

	// ExitMove_table_partition is called when exiting the move_table_partition production.
	ExitMove_table_partition(c *Move_table_partitionContext)

	// ExitFilter_condition is called when exiting the filter_condition production.
	ExitFilter_condition(c *Filter_conditionContext)

	// ExitRename_table_partition is called when exiting the rename_table_partition production.
	ExitRename_table_partition(c *Rename_table_partitionContext)

	// ExitPartition_extended_names is called when exiting the partition_extended_names production.
	ExitPartition_extended_names(c *Partition_extended_namesContext)

	// ExitSubpartition_extended_names is called when exiting the subpartition_extended_names production.
	ExitSubpartition_extended_names(c *Subpartition_extended_namesContext)

	// ExitAlter_table_properties_1 is called when exiting the alter_table_properties_1 production.
	ExitAlter_table_properties_1(c *Alter_table_properties_1Context)

	// ExitAlter_iot_clauses is called when exiting the alter_iot_clauses production.
	ExitAlter_iot_clauses(c *Alter_iot_clausesContext)

	// ExitAlter_mapping_table_clause is called when exiting the alter_mapping_table_clause production.
	ExitAlter_mapping_table_clause(c *Alter_mapping_table_clauseContext)

	// ExitAlter_overflow_clause is called when exiting the alter_overflow_clause production.
	ExitAlter_overflow_clause(c *Alter_overflow_clauseContext)

	// ExitAdd_overflow_clause is called when exiting the add_overflow_clause production.
	ExitAdd_overflow_clause(c *Add_overflow_clauseContext)

	// ExitUpdate_index_clauses is called when exiting the update_index_clauses production.
	ExitUpdate_index_clauses(c *Update_index_clausesContext)

	// ExitUpdate_global_index_clause is called when exiting the update_global_index_clause production.
	ExitUpdate_global_index_clause(c *Update_global_index_clauseContext)

	// ExitUpdate_all_indexes_clause is called when exiting the update_all_indexes_clause production.
	ExitUpdate_all_indexes_clause(c *Update_all_indexes_clauseContext)

	// ExitUpdate_all_indexes_index_clause is called when exiting the update_all_indexes_index_clause production.
	ExitUpdate_all_indexes_index_clause(c *Update_all_indexes_index_clauseContext)

	// ExitUpdate_index_partition is called when exiting the update_index_partition production.
	ExitUpdate_index_partition(c *Update_index_partitionContext)

	// ExitUpdate_index_subpartition is called when exiting the update_index_subpartition production.
	ExitUpdate_index_subpartition(c *Update_index_subpartitionContext)

	// ExitEnable_disable_clause is called when exiting the enable_disable_clause production.
	ExitEnable_disable_clause(c *Enable_disable_clauseContext)

	// ExitUsing_index_clause is called when exiting the using_index_clause production.
	ExitUsing_index_clause(c *Using_index_clauseContext)

	// ExitIndex_attributes is called when exiting the index_attributes production.
	ExitIndex_attributes(c *Index_attributesContext)

	// ExitSort_or_nosort is called when exiting the sort_or_nosort production.
	ExitSort_or_nosort(c *Sort_or_nosortContext)

	// ExitExceptions_clause is called when exiting the exceptions_clause production.
	ExitExceptions_clause(c *Exceptions_clauseContext)

	// ExitMove_table_clause is called when exiting the move_table_clause production.
	ExitMove_table_clause(c *Move_table_clauseContext)

	// ExitIndex_org_table_clause is called when exiting the index_org_table_clause production.
	ExitIndex_org_table_clause(c *Index_org_table_clauseContext)

	// ExitMapping_table_clause is called when exiting the mapping_table_clause production.
	ExitMapping_table_clause(c *Mapping_table_clauseContext)

	// ExitKey_compression is called when exiting the key_compression production.
	ExitKey_compression(c *Key_compressionContext)

	// ExitIndex_org_overflow_clause is called when exiting the index_org_overflow_clause production.
	ExitIndex_org_overflow_clause(c *Index_org_overflow_clauseContext)

	// ExitColumn_clauses is called when exiting the column_clauses production.
	ExitColumn_clauses(c *Column_clausesContext)

	// ExitModify_collection_retrieval is called when exiting the modify_collection_retrieval production.
	ExitModify_collection_retrieval(c *Modify_collection_retrievalContext)

	// ExitCollection_item is called when exiting the collection_item production.
	ExitCollection_item(c *Collection_itemContext)

	// ExitRename_column_clause is called when exiting the rename_column_clause production.
	ExitRename_column_clause(c *Rename_column_clauseContext)

	// ExitOld_column_name is called when exiting the old_column_name production.
	ExitOld_column_name(c *Old_column_nameContext)

	// ExitNew_column_name is called when exiting the new_column_name production.
	ExitNew_column_name(c *New_column_nameContext)

	// ExitAdd_modify_drop_column_clauses is called when exiting the add_modify_drop_column_clauses production.
	ExitAdd_modify_drop_column_clauses(c *Add_modify_drop_column_clausesContext)

	// ExitDrop_column_clause is called when exiting the drop_column_clause production.
	ExitDrop_column_clause(c *Drop_column_clauseContext)

	// ExitModify_column_clauses is called when exiting the modify_column_clauses production.
	ExitModify_column_clauses(c *Modify_column_clausesContext)

	// ExitModify_col_properties is called when exiting the modify_col_properties production.
	ExitModify_col_properties(c *Modify_col_propertiesContext)

	// ExitModify_col_visibility is called when exiting the modify_col_visibility production.
	ExitModify_col_visibility(c *Modify_col_visibilityContext)

	// ExitModify_col_substitutable is called when exiting the modify_col_substitutable production.
	ExitModify_col_substitutable(c *Modify_col_substitutableContext)

	// ExitAdd_column_clause is called when exiting the add_column_clause production.
	ExitAdd_column_clause(c *Add_column_clauseContext)

	// ExitVarray_col_properties is called when exiting the varray_col_properties production.
	ExitVarray_col_properties(c *Varray_col_propertiesContext)

	// ExitVarray_storage_clause is called when exiting the varray_storage_clause production.
	ExitVarray_storage_clause(c *Varray_storage_clauseContext)

	// ExitLob_segname is called when exiting the lob_segname production.
	ExitLob_segname(c *Lob_segnameContext)

	// ExitLob_item is called when exiting the lob_item production.
	ExitLob_item(c *Lob_itemContext)

	// ExitLob_storage_parameters is called when exiting the lob_storage_parameters production.
	ExitLob_storage_parameters(c *Lob_storage_parametersContext)

	// ExitLob_storage_clause is called when exiting the lob_storage_clause production.
	ExitLob_storage_clause(c *Lob_storage_clauseContext)

	// ExitModify_lob_storage_clause is called when exiting the modify_lob_storage_clause production.
	ExitModify_lob_storage_clause(c *Modify_lob_storage_clauseContext)

	// ExitModify_lob_parameters is called when exiting the modify_lob_parameters production.
	ExitModify_lob_parameters(c *Modify_lob_parametersContext)

	// ExitLob_parameters is called when exiting the lob_parameters production.
	ExitLob_parameters(c *Lob_parametersContext)

	// ExitLob_deduplicate_clause is called when exiting the lob_deduplicate_clause production.
	ExitLob_deduplicate_clause(c *Lob_deduplicate_clauseContext)

	// ExitLob_compression_clause is called when exiting the lob_compression_clause production.
	ExitLob_compression_clause(c *Lob_compression_clauseContext)

	// ExitLob_retention_clause is called when exiting the lob_retention_clause production.
	ExitLob_retention_clause(c *Lob_retention_clauseContext)

	// ExitEncryption_spec is called when exiting the encryption_spec production.
	ExitEncryption_spec(c *Encryption_specContext)

	// ExitTablespace is called when exiting the tablespace production.
	ExitTablespace(c *TablespaceContext)

	// ExitVarray_item is called when exiting the varray_item production.
	ExitVarray_item(c *Varray_itemContext)

	// ExitColumn_properties is called when exiting the column_properties production.
	ExitColumn_properties(c *Column_propertiesContext)

	// ExitLob_partition_storage is called when exiting the lob_partition_storage production.
	ExitLob_partition_storage(c *Lob_partition_storageContext)

	// ExitPeriod_definition is called when exiting the period_definition production.
	ExitPeriod_definition(c *Period_definitionContext)

	// ExitStart_time_column is called when exiting the start_time_column production.
	ExitStart_time_column(c *Start_time_columnContext)

	// ExitEnd_time_column is called when exiting the end_time_column production.
	ExitEnd_time_column(c *End_time_columnContext)

	// ExitColumn_definition is called when exiting the column_definition production.
	ExitColumn_definition(c *Column_definitionContext)

	// ExitColumn_collation_name is called when exiting the column_collation_name production.
	ExitColumn_collation_name(c *Column_collation_nameContext)

	// ExitIdentity_clause is called when exiting the identity_clause production.
	ExitIdentity_clause(c *Identity_clauseContext)

	// ExitIdentity_options_parentheses is called when exiting the identity_options_parentheses production.
	ExitIdentity_options_parentheses(c *Identity_options_parenthesesContext)

	// ExitIdentity_options is called when exiting the identity_options production.
	ExitIdentity_options(c *Identity_optionsContext)

	// ExitVirtual_column_definition is called when exiting the virtual_column_definition production.
	ExitVirtual_column_definition(c *Virtual_column_definitionContext)

	// ExitVirtual_column_expression is called when exiting the virtual_column_expression production.
	ExitVirtual_column_expression(c *Virtual_column_expressionContext)

	// ExitAutogenerated_sequence_definition is called when exiting the autogenerated_sequence_definition production.
	ExitAutogenerated_sequence_definition(c *Autogenerated_sequence_definitionContext)

	// ExitBy_user_for_statistics_clause is called when exiting the by_user_for_statistics_clause production.
	ExitBy_user_for_statistics_clause(c *By_user_for_statistics_clauseContext)

	// ExitEvaluation_edition_clause is called when exiting the evaluation_edition_clause production.
	ExitEvaluation_edition_clause(c *Evaluation_edition_clauseContext)

	// ExitNested_table_col_properties is called when exiting the nested_table_col_properties production.
	ExitNested_table_col_properties(c *Nested_table_col_propertiesContext)

	// ExitNested_item is called when exiting the nested_item production.
	ExitNested_item(c *Nested_itemContext)

	// ExitSubstitutable_column_clause is called when exiting the substitutable_column_clause production.
	ExitSubstitutable_column_clause(c *Substitutable_column_clauseContext)

	// ExitPartition_name is called when exiting the partition_name production.
	ExitPartition_name(c *Partition_nameContext)

	// ExitSupplemental_logging_props is called when exiting the supplemental_logging_props production.
	ExitSupplemental_logging_props(c *Supplemental_logging_propsContext)

	// ExitObject_type_col_properties is called when exiting the object_type_col_properties production.
	ExitObject_type_col_properties(c *Object_type_col_propertiesContext)

	// ExitConstraint_clauses is called when exiting the constraint_clauses production.
	ExitConstraint_clauses(c *Constraint_clausesContext)

	// ExitOld_constraint_name is called when exiting the old_constraint_name production.
	ExitOld_constraint_name(c *Old_constraint_nameContext)

	// ExitNew_constraint_name is called when exiting the new_constraint_name production.
	ExitNew_constraint_name(c *New_constraint_nameContext)

	// ExitDrop_constraint_clause is called when exiting the drop_constraint_clause production.
	ExitDrop_constraint_clause(c *Drop_constraint_clauseContext)

	// ExitCheck_constraint is called when exiting the check_constraint production.
	ExitCheck_constraint(c *Check_constraintContext)

	// ExitForeign_key_clause is called when exiting the foreign_key_clause production.
	ExitForeign_key_clause(c *Foreign_key_clauseContext)

	// ExitReferences_clause is called when exiting the references_clause production.
	ExitReferences_clause(c *References_clauseContext)

	// ExitOn_delete_clause is called when exiting the on_delete_clause production.
	ExitOn_delete_clause(c *On_delete_clauseContext)

	// ExitAnonymous_block is called when exiting the anonymous_block production.
	ExitAnonymous_block(c *Anonymous_blockContext)

	// ExitInvoker_rights_clause is called when exiting the invoker_rights_clause production.
	ExitInvoker_rights_clause(c *Invoker_rights_clauseContext)

	// ExitCall_spec is called when exiting the call_spec production.
	ExitCall_spec(c *Call_specContext)

	// ExitJava_spec is called when exiting the java_spec production.
	ExitJava_spec(c *Java_specContext)

	// ExitC_spec is called when exiting the c_spec production.
	ExitC_spec(c *C_specContext)

	// ExitC_agent_in_clause is called when exiting the c_agent_in_clause production.
	ExitC_agent_in_clause(c *C_agent_in_clauseContext)

	// ExitC_parameters_clause is called when exiting the c_parameters_clause production.
	ExitC_parameters_clause(c *C_parameters_clauseContext)

	// ExitC_external_parameter is called when exiting the c_external_parameter production.
	ExitC_external_parameter(c *C_external_parameterContext)

	// ExitC_property is called when exiting the c_property production.
	ExitC_property(c *C_propertyContext)

	// ExitParameter is called when exiting the parameter production.
	ExitParameter(c *ParameterContext)

	// ExitDefault_value_part is called when exiting the default_value_part production.
	ExitDefault_value_part(c *Default_value_partContext)

	// ExitSeq_of_declare_specs is called when exiting the seq_of_declare_specs production.
	ExitSeq_of_declare_specs(c *Seq_of_declare_specsContext)

	// ExitDeclare_spec is called when exiting the declare_spec production.
	ExitDeclare_spec(c *Declare_specContext)

	// ExitVariable_declaration is called when exiting the variable_declaration production.
	ExitVariable_declaration(c *Variable_declarationContext)

	// ExitSubtype_declaration is called when exiting the subtype_declaration production.
	ExitSubtype_declaration(c *Subtype_declarationContext)

	// ExitCursor_declaration is called when exiting the cursor_declaration production.
	ExitCursor_declaration(c *Cursor_declarationContext)

	// ExitParameter_spec is called when exiting the parameter_spec production.
	ExitParameter_spec(c *Parameter_specContext)

	// ExitException_declaration is called when exiting the exception_declaration production.
	ExitException_declaration(c *Exception_declarationContext)

	// ExitPragma_declaration is called when exiting the pragma_declaration production.
	ExitPragma_declaration(c *Pragma_declarationContext)

	// ExitRecord_type_def is called when exiting the record_type_def production.
	ExitRecord_type_def(c *Record_type_defContext)

	// ExitField_spec is called when exiting the field_spec production.
	ExitField_spec(c *Field_specContext)

	// ExitRef_cursor_type_def is called when exiting the ref_cursor_type_def production.
	ExitRef_cursor_type_def(c *Ref_cursor_type_defContext)

	// ExitType_declaration is called when exiting the type_declaration production.
	ExitType_declaration(c *Type_declarationContext)

	// ExitTable_type_def is called when exiting the table_type_def production.
	ExitTable_type_def(c *Table_type_defContext)

	// ExitTable_indexed_by_part is called when exiting the table_indexed_by_part production.
	ExitTable_indexed_by_part(c *Table_indexed_by_partContext)

	// ExitVarray_type_def is called when exiting the varray_type_def production.
	ExitVarray_type_def(c *Varray_type_defContext)

	// ExitSeq_of_statements is called when exiting the seq_of_statements production.
	ExitSeq_of_statements(c *Seq_of_statementsContext)

	// ExitLabel_declaration is called when exiting the label_declaration production.
	ExitLabel_declaration(c *Label_declarationContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitAssignment_statement is called when exiting the assignment_statement production.
	ExitAssignment_statement(c *Assignment_statementContext)

	// ExitContinue_statement is called when exiting the continue_statement production.
	ExitContinue_statement(c *Continue_statementContext)

	// ExitExit_statement is called when exiting the exit_statement production.
	ExitExit_statement(c *Exit_statementContext)

	// ExitGoto_statement is called when exiting the goto_statement production.
	ExitGoto_statement(c *Goto_statementContext)

	// ExitIf_statement is called when exiting the if_statement production.
	ExitIf_statement(c *If_statementContext)

	// ExitElsif_part is called when exiting the elsif_part production.
	ExitElsif_part(c *Elsif_partContext)

	// ExitElse_part is called when exiting the else_part production.
	ExitElse_part(c *Else_partContext)

	// ExitLoop_statement is called when exiting the loop_statement production.
	ExitLoop_statement(c *Loop_statementContext)

	// ExitCursor_loop_param is called when exiting the cursor_loop_param production.
	ExitCursor_loop_param(c *Cursor_loop_paramContext)

	// ExitForall_statement is called when exiting the forall_statement production.
	ExitForall_statement(c *Forall_statementContext)

	// ExitBounds_clause is called when exiting the bounds_clause production.
	ExitBounds_clause(c *Bounds_clauseContext)

	// ExitBetween_bound is called when exiting the between_bound production.
	ExitBetween_bound(c *Between_boundContext)

	// ExitLower_bound is called when exiting the lower_bound production.
	ExitLower_bound(c *Lower_boundContext)

	// ExitUpper_bound is called when exiting the upper_bound production.
	ExitUpper_bound(c *Upper_boundContext)

	// ExitNull_statement is called when exiting the null_statement production.
	ExitNull_statement(c *Null_statementContext)

	// ExitRaise_statement is called when exiting the raise_statement production.
	ExitRaise_statement(c *Raise_statementContext)

	// ExitReturn_statement is called when exiting the return_statement production.
	ExitReturn_statement(c *Return_statementContext)

	// ExitCall_statement is called when exiting the call_statement production.
	ExitCall_statement(c *Call_statementContext)

	// ExitPipe_row_statement is called when exiting the pipe_row_statement production.
	ExitPipe_row_statement(c *Pipe_row_statementContext)

	// ExitSelection_directive is called when exiting the selection_directive production.
	ExitSelection_directive(c *Selection_directiveContext)

	// ExitError_directive is called when exiting the error_directive production.
	ExitError_directive(c *Error_directiveContext)

	// ExitSelection_directive_body is called when exiting the selection_directive_body production.
	ExitSelection_directive_body(c *Selection_directive_bodyContext)

	// ExitBody is called when exiting the body production.
	ExitBody(c *BodyContext)

	// ExitException_handler is called when exiting the exception_handler production.
	ExitException_handler(c *Exception_handlerContext)

	// ExitTrigger_block is called when exiting the trigger_block production.
	ExitTrigger_block(c *Trigger_blockContext)

	// ExitTps_block is called when exiting the tps_block production.
	ExitTps_block(c *Tps_blockContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitSql_statement is called when exiting the sql_statement production.
	ExitSql_statement(c *Sql_statementContext)

	// ExitExecute_immediate is called when exiting the execute_immediate production.
	ExitExecute_immediate(c *Execute_immediateContext)

	// ExitDynamic_returning_clause is called when exiting the dynamic_returning_clause production.
	ExitDynamic_returning_clause(c *Dynamic_returning_clauseContext)

	// ExitData_manipulation_language_statements is called when exiting the data_manipulation_language_statements production.
	ExitData_manipulation_language_statements(c *Data_manipulation_language_statementsContext)

	// ExitCursor_manipulation_statements is called when exiting the cursor_manipulation_statements production.
	ExitCursor_manipulation_statements(c *Cursor_manipulation_statementsContext)

	// ExitClose_statement is called when exiting the close_statement production.
	ExitClose_statement(c *Close_statementContext)

	// ExitOpen_statement is called when exiting the open_statement production.
	ExitOpen_statement(c *Open_statementContext)

	// ExitFetch_statement is called when exiting the fetch_statement production.
	ExitFetch_statement(c *Fetch_statementContext)

	// ExitVariable_or_collection is called when exiting the variable_or_collection production.
	ExitVariable_or_collection(c *Variable_or_collectionContext)

	// ExitOpen_for_statement is called when exiting the open_for_statement production.
	ExitOpen_for_statement(c *Open_for_statementContext)

	// ExitTransaction_control_statements is called when exiting the transaction_control_statements production.
	ExitTransaction_control_statements(c *Transaction_control_statementsContext)

	// ExitSet_transaction_command is called when exiting the set_transaction_command production.
	ExitSet_transaction_command(c *Set_transaction_commandContext)

	// ExitSet_constraint_command is called when exiting the set_constraint_command production.
	ExitSet_constraint_command(c *Set_constraint_commandContext)

	// ExitCommit_statement is called when exiting the commit_statement production.
	ExitCommit_statement(c *Commit_statementContext)

	// ExitWrite_clause is called when exiting the write_clause production.
	ExitWrite_clause(c *Write_clauseContext)

	// ExitRollback_statement is called when exiting the rollback_statement production.
	ExitRollback_statement(c *Rollback_statementContext)

	// ExitSavepoint_statement is called when exiting the savepoint_statement production.
	ExitSavepoint_statement(c *Savepoint_statementContext)

	// ExitCollection_method_call is called when exiting the collection_method_call production.
	ExitCollection_method_call(c *Collection_method_callContext)

	// ExitExplain_statement is called when exiting the explain_statement production.
	ExitExplain_statement(c *Explain_statementContext)

	// ExitSelect_only_statement is called when exiting the select_only_statement production.
	ExitSelect_only_statement(c *Select_only_statementContext)

	// ExitSelect_statement is called when exiting the select_statement production.
	ExitSelect_statement(c *Select_statementContext)

	// ExitWith_clause is called when exiting the with_clause production.
	ExitWith_clause(c *With_clauseContext)

	// ExitWith_factoring_clause is called when exiting the with_factoring_clause production.
	ExitWith_factoring_clause(c *With_factoring_clauseContext)

	// ExitSubquery_factoring_clause is called when exiting the subquery_factoring_clause production.
	ExitSubquery_factoring_clause(c *Subquery_factoring_clauseContext)

	// ExitSearch_clause is called when exiting the search_clause production.
	ExitSearch_clause(c *Search_clauseContext)

	// ExitCycle_clause is called when exiting the cycle_clause production.
	ExitCycle_clause(c *Cycle_clauseContext)

	// ExitSubav_factoring_clause is called when exiting the subav_factoring_clause production.
	ExitSubav_factoring_clause(c *Subav_factoring_clauseContext)

	// ExitSubav_clause is called when exiting the subav_clause production.
	ExitSubav_clause(c *Subav_clauseContext)

	// ExitHierarchies_clause is called when exiting the hierarchies_clause production.
	ExitHierarchies_clause(c *Hierarchies_clauseContext)

	// ExitFilter_clauses is called when exiting the filter_clauses production.
	ExitFilter_clauses(c *Filter_clausesContext)

	// ExitFilter_clause is called when exiting the filter_clause production.
	ExitFilter_clause(c *Filter_clauseContext)

	// ExitAdd_calcs_clause is called when exiting the add_calcs_clause production.
	ExitAdd_calcs_clause(c *Add_calcs_clauseContext)

	// ExitAdd_calc_meas_clause is called when exiting the add_calc_meas_clause production.
	ExitAdd_calc_meas_clause(c *Add_calc_meas_clauseContext)

	// ExitSubquery is called when exiting the subquery production.
	ExitSubquery(c *SubqueryContext)

	// ExitSubquery_basic_elements is called when exiting the subquery_basic_elements production.
	ExitSubquery_basic_elements(c *Subquery_basic_elementsContext)

	// ExitSubquery_operation_part is called when exiting the subquery_operation_part production.
	ExitSubquery_operation_part(c *Subquery_operation_partContext)

	// ExitQuery_block is called when exiting the query_block production.
	ExitQuery_block(c *Query_blockContext)

	// ExitSelected_list is called when exiting the selected_list production.
	ExitSelected_list(c *Selected_listContext)

	// ExitFrom_clause is called when exiting the from_clause production.
	ExitFrom_clause(c *From_clauseContext)

	// ExitSelect_list_elements is called when exiting the select_list_elements production.
	ExitSelect_list_elements(c *Select_list_elementsContext)

	// ExitTable_ref_list is called when exiting the table_ref_list production.
	ExitTable_ref_list(c *Table_ref_listContext)

	// ExitTable_ref is called when exiting the table_ref production.
	ExitTable_ref(c *Table_refContext)

	// ExitTable_ref_aux is called when exiting the table_ref_aux production.
	ExitTable_ref_aux(c *Table_ref_auxContext)

	// ExitTable_ref_aux_internal_one is called when exiting the table_ref_aux_internal_one production.
	ExitTable_ref_aux_internal_one(c *Table_ref_aux_internal_oneContext)

	// ExitTable_ref_aux_internal_two is called when exiting the table_ref_aux_internal_two production.
	ExitTable_ref_aux_internal_two(c *Table_ref_aux_internal_twoContext)

	// ExitTable_ref_aux_internal_thre is called when exiting the table_ref_aux_internal_thre production.
	ExitTable_ref_aux_internal_thre(c *Table_ref_aux_internal_threContext)

	// ExitJoin_clause is called when exiting the join_clause production.
	ExitJoin_clause(c *Join_clauseContext)

	// ExitJoin_on_part is called when exiting the join_on_part production.
	ExitJoin_on_part(c *Join_on_partContext)

	// ExitJoin_using_part is called when exiting the join_using_part production.
	ExitJoin_using_part(c *Join_using_partContext)

	// ExitOuter_join_type is called when exiting the outer_join_type production.
	ExitOuter_join_type(c *Outer_join_typeContext)

	// ExitQuery_partition_clause is called when exiting the query_partition_clause production.
	ExitQuery_partition_clause(c *Query_partition_clauseContext)

	// ExitFlashback_query_clause is called when exiting the flashback_query_clause production.
	ExitFlashback_query_clause(c *Flashback_query_clauseContext)

	// ExitPivot_clause is called when exiting the pivot_clause production.
	ExitPivot_clause(c *Pivot_clauseContext)

	// ExitPivot_element is called when exiting the pivot_element production.
	ExitPivot_element(c *Pivot_elementContext)

	// ExitPivot_for_clause is called when exiting the pivot_for_clause production.
	ExitPivot_for_clause(c *Pivot_for_clauseContext)

	// ExitPivot_in_clause is called when exiting the pivot_in_clause production.
	ExitPivot_in_clause(c *Pivot_in_clauseContext)

	// ExitPivot_in_clause_element is called when exiting the pivot_in_clause_element production.
	ExitPivot_in_clause_element(c *Pivot_in_clause_elementContext)

	// ExitPivot_in_clause_elements is called when exiting the pivot_in_clause_elements production.
	ExitPivot_in_clause_elements(c *Pivot_in_clause_elementsContext)

	// ExitUnpivot_clause is called when exiting the unpivot_clause production.
	ExitUnpivot_clause(c *Unpivot_clauseContext)

	// ExitUnpivot_in_clause is called when exiting the unpivot_in_clause production.
	ExitUnpivot_in_clause(c *Unpivot_in_clauseContext)

	// ExitUnpivot_in_elements is called when exiting the unpivot_in_elements production.
	ExitUnpivot_in_elements(c *Unpivot_in_elementsContext)

	// ExitHierarchical_query_clause is called when exiting the hierarchical_query_clause production.
	ExitHierarchical_query_clause(c *Hierarchical_query_clauseContext)

	// ExitStart_part is called when exiting the start_part production.
	ExitStart_part(c *Start_partContext)

	// ExitGroup_by_clause is called when exiting the group_by_clause production.
	ExitGroup_by_clause(c *Group_by_clauseContext)

	// ExitGroup_by_elements is called when exiting the group_by_elements production.
	ExitGroup_by_elements(c *Group_by_elementsContext)

	// ExitRollup_cube_clause is called when exiting the rollup_cube_clause production.
	ExitRollup_cube_clause(c *Rollup_cube_clauseContext)

	// ExitGrouping_sets_clause is called when exiting the grouping_sets_clause production.
	ExitGrouping_sets_clause(c *Grouping_sets_clauseContext)

	// ExitGrouping_sets_elements is called when exiting the grouping_sets_elements production.
	ExitGrouping_sets_elements(c *Grouping_sets_elementsContext)

	// ExitHaving_clause is called when exiting the having_clause production.
	ExitHaving_clause(c *Having_clauseContext)

	// ExitModel_clause is called when exiting the model_clause production.
	ExitModel_clause(c *Model_clauseContext)

	// ExitCell_reference_options is called when exiting the cell_reference_options production.
	ExitCell_reference_options(c *Cell_reference_optionsContext)

	// ExitReturn_rows_clause is called when exiting the return_rows_clause production.
	ExitReturn_rows_clause(c *Return_rows_clauseContext)

	// ExitReference_model is called when exiting the reference_model production.
	ExitReference_model(c *Reference_modelContext)

	// ExitMain_model is called when exiting the main_model production.
	ExitMain_model(c *Main_modelContext)

	// ExitModel_column_clauses is called when exiting the model_column_clauses production.
	ExitModel_column_clauses(c *Model_column_clausesContext)

	// ExitModel_column_partition_part is called when exiting the model_column_partition_part production.
	ExitModel_column_partition_part(c *Model_column_partition_partContext)

	// ExitModel_column_list is called when exiting the model_column_list production.
	ExitModel_column_list(c *Model_column_listContext)

	// ExitModel_column is called when exiting the model_column production.
	ExitModel_column(c *Model_columnContext)

	// ExitModel_rules_clause is called when exiting the model_rules_clause production.
	ExitModel_rules_clause(c *Model_rules_clauseContext)

	// ExitModel_rules_part is called when exiting the model_rules_part production.
	ExitModel_rules_part(c *Model_rules_partContext)

	// ExitModel_rules_element is called when exiting the model_rules_element production.
	ExitModel_rules_element(c *Model_rules_elementContext)

	// ExitCell_assignment is called when exiting the cell_assignment production.
	ExitCell_assignment(c *Cell_assignmentContext)

	// ExitModel_iterate_clause is called when exiting the model_iterate_clause production.
	ExitModel_iterate_clause(c *Model_iterate_clauseContext)

	// ExitUntil_part is called when exiting the until_part production.
	ExitUntil_part(c *Until_partContext)

	// ExitOrder_by_clause is called when exiting the order_by_clause production.
	ExitOrder_by_clause(c *Order_by_clauseContext)

	// ExitOrder_by_elements is called when exiting the order_by_elements production.
	ExitOrder_by_elements(c *Order_by_elementsContext)

	// ExitOffset_clause is called when exiting the offset_clause production.
	ExitOffset_clause(c *Offset_clauseContext)

	// ExitFetch_clause is called when exiting the fetch_clause production.
	ExitFetch_clause(c *Fetch_clauseContext)

	// ExitFor_update_clause is called when exiting the for_update_clause production.
	ExitFor_update_clause(c *For_update_clauseContext)

	// ExitFor_update_of_part is called when exiting the for_update_of_part production.
	ExitFor_update_of_part(c *For_update_of_partContext)

	// ExitFor_update_options is called when exiting the for_update_options production.
	ExitFor_update_options(c *For_update_optionsContext)

	// ExitUpdate_statement is called when exiting the update_statement production.
	ExitUpdate_statement(c *Update_statementContext)

	// ExitUpdate_set_clause is called when exiting the update_set_clause production.
	ExitUpdate_set_clause(c *Update_set_clauseContext)

	// ExitColumn_based_update_set_clause is called when exiting the column_based_update_set_clause production.
	ExitColumn_based_update_set_clause(c *Column_based_update_set_clauseContext)

	// ExitDelete_statement is called when exiting the delete_statement production.
	ExitDelete_statement(c *Delete_statementContext)

	// ExitInsert_statement is called when exiting the insert_statement production.
	ExitInsert_statement(c *Insert_statementContext)

	// ExitSingle_table_insert is called when exiting the single_table_insert production.
	ExitSingle_table_insert(c *Single_table_insertContext)

	// ExitMulti_table_insert is called when exiting the multi_table_insert production.
	ExitMulti_table_insert(c *Multi_table_insertContext)

	// ExitMulti_table_element is called when exiting the multi_table_element production.
	ExitMulti_table_element(c *Multi_table_elementContext)

	// ExitConditional_insert_clause is called when exiting the conditional_insert_clause production.
	ExitConditional_insert_clause(c *Conditional_insert_clauseContext)

	// ExitConditional_insert_when_part is called when exiting the conditional_insert_when_part production.
	ExitConditional_insert_when_part(c *Conditional_insert_when_partContext)

	// ExitConditional_insert_else_part is called when exiting the conditional_insert_else_part production.
	ExitConditional_insert_else_part(c *Conditional_insert_else_partContext)

	// ExitInsert_into_clause is called when exiting the insert_into_clause production.
	ExitInsert_into_clause(c *Insert_into_clauseContext)

	// ExitValues_clause is called when exiting the values_clause production.
	ExitValues_clause(c *Values_clauseContext)

	// ExitMerge_statement is called when exiting the merge_statement production.
	ExitMerge_statement(c *Merge_statementContext)

	// ExitMerge_update_clause is called when exiting the merge_update_clause production.
	ExitMerge_update_clause(c *Merge_update_clauseContext)

	// ExitMerge_element is called when exiting the merge_element production.
	ExitMerge_element(c *Merge_elementContext)

	// ExitMerge_update_delete_part is called when exiting the merge_update_delete_part production.
	ExitMerge_update_delete_part(c *Merge_update_delete_partContext)

	// ExitMerge_insert_clause is called when exiting the merge_insert_clause production.
	ExitMerge_insert_clause(c *Merge_insert_clauseContext)

	// ExitSelected_tableview is called when exiting the selected_tableview production.
	ExitSelected_tableview(c *Selected_tableviewContext)

	// ExitLock_table_statement is called when exiting the lock_table_statement production.
	ExitLock_table_statement(c *Lock_table_statementContext)

	// ExitWait_nowait_part is called when exiting the wait_nowait_part production.
	ExitWait_nowait_part(c *Wait_nowait_partContext)

	// ExitLock_table_element is called when exiting the lock_table_element production.
	ExitLock_table_element(c *Lock_table_elementContext)

	// ExitLock_mode is called when exiting the lock_mode production.
	ExitLock_mode(c *Lock_modeContext)

	// ExitGeneral_table_ref is called when exiting the general_table_ref production.
	ExitGeneral_table_ref(c *General_table_refContext)

	// ExitStatic_returning_clause is called when exiting the static_returning_clause production.
	ExitStatic_returning_clause(c *Static_returning_clauseContext)

	// ExitError_logging_clause is called when exiting the error_logging_clause production.
	ExitError_logging_clause(c *Error_logging_clauseContext)

	// ExitError_logging_into_part is called when exiting the error_logging_into_part production.
	ExitError_logging_into_part(c *Error_logging_into_partContext)

	// ExitError_logging_reject_part is called when exiting the error_logging_reject_part production.
	ExitError_logging_reject_part(c *Error_logging_reject_partContext)

	// ExitDml_table_expression_clause is called when exiting the dml_table_expression_clause production.
	ExitDml_table_expression_clause(c *Dml_table_expression_clauseContext)

	// ExitTable_collection_expression is called when exiting the table_collection_expression production.
	ExitTable_collection_expression(c *Table_collection_expressionContext)

	// ExitSubquery_restriction_clause is called when exiting the subquery_restriction_clause production.
	ExitSubquery_restriction_clause(c *Subquery_restriction_clauseContext)

	// ExitSample_clause is called when exiting the sample_clause production.
	ExitSample_clause(c *Sample_clauseContext)

	// ExitSeed_part is called when exiting the seed_part production.
	ExitSeed_part(c *Seed_partContext)

	// ExitCondition is called when exiting the condition production.
	ExitCondition(c *ConditionContext)

	// ExitExpressions_ is called when exiting the expressions_ production.
	ExitExpressions_(c *Expressions_Context)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitCursor_expression is called when exiting the cursor_expression production.
	ExitCursor_expression(c *Cursor_expressionContext)

	// ExitLogical_expression is called when exiting the logical_expression production.
	ExitLogical_expression(c *Logical_expressionContext)

	// ExitUnary_logical_expression is called when exiting the unary_logical_expression production.
	ExitUnary_logical_expression(c *Unary_logical_expressionContext)

	// ExitUnary_logical_operation is called when exiting the unary_logical_operation production.
	ExitUnary_logical_operation(c *Unary_logical_operationContext)

	// ExitLogical_operation is called when exiting the logical_operation production.
	ExitLogical_operation(c *Logical_operationContext)

	// ExitMultiset_expression is called when exiting the multiset_expression production.
	ExitMultiset_expression(c *Multiset_expressionContext)

	// ExitRelational_expression is called when exiting the relational_expression production.
	ExitRelational_expression(c *Relational_expressionContext)

	// ExitCompound_expression is called when exiting the compound_expression production.
	ExitCompound_expression(c *Compound_expressionContext)

	// ExitRelational_operator is called when exiting the relational_operator production.
	ExitRelational_operator(c *Relational_operatorContext)

	// ExitIn_elements is called when exiting the in_elements production.
	ExitIn_elements(c *In_elementsContext)

	// ExitBetween_elements is called when exiting the between_elements production.
	ExitBetween_elements(c *Between_elementsContext)

	// ExitConcatenation is called when exiting the concatenation production.
	ExitConcatenation(c *ConcatenationContext)

	// ExitInterval_expression is called when exiting the interval_expression production.
	ExitInterval_expression(c *Interval_expressionContext)

	// ExitModel_expression is called when exiting the model_expression production.
	ExitModel_expression(c *Model_expressionContext)

	// ExitModel_expression_element is called when exiting the model_expression_element production.
	ExitModel_expression_element(c *Model_expression_elementContext)

	// ExitSingle_column_for_loop is called when exiting the single_column_for_loop production.
	ExitSingle_column_for_loop(c *Single_column_for_loopContext)

	// ExitMulti_column_for_loop is called when exiting the multi_column_for_loop production.
	ExitMulti_column_for_loop(c *Multi_column_for_loopContext)

	// ExitUnary_expression is called when exiting the unary_expression production.
	ExitUnary_expression(c *Unary_expressionContext)

	// ExitUnary_expression_core is called when exiting the unary_expression_core production.
	ExitUnary_expression_core(c *Unary_expression_coreContext)

	// ExitImplicit_cursor_expression is called when exiting the implicit_cursor_expression production.
	ExitImplicit_cursor_expression(c *Implicit_cursor_expressionContext)

	// ExitCollection_expression is called when exiting the collection_expression production.
	ExitCollection_expression(c *Collection_expressionContext)

	// ExitCase_statement is called when exiting the case_statement production.
	ExitCase_statement(c *Case_statementContext)

	// ExitSimple_case_statement is called when exiting the simple_case_statement production.
	ExitSimple_case_statement(c *Simple_case_statementContext)

	// ExitSearched_case_statement is called when exiting the searched_case_statement production.
	ExitSearched_case_statement(c *Searched_case_statementContext)

	// ExitCase_when_part_statement is called when exiting the case_when_part_statement production.
	ExitCase_when_part_statement(c *Case_when_part_statementContext)

	// ExitCase_else_part_statement is called when exiting the case_else_part_statement production.
	ExitCase_else_part_statement(c *Case_else_part_statementContext)

	// ExitCase_expression is called when exiting the case_expression production.
	ExitCase_expression(c *Case_expressionContext)

	// ExitSimple_case_expression is called when exiting the simple_case_expression production.
	ExitSimple_case_expression(c *Simple_case_expressionContext)

	// ExitSearched_case_expression is called when exiting the searched_case_expression production.
	ExitSearched_case_expression(c *Searched_case_expressionContext)

	// ExitCase_when_part_expression is called when exiting the case_when_part_expression production.
	ExitCase_when_part_expression(c *Case_when_part_expressionContext)

	// ExitCase_else_part_expression is called when exiting the case_else_part_expression production.
	ExitCase_else_part_expression(c *Case_else_part_expressionContext)

	// ExitAtom is called when exiting the atom production.
	ExitAtom(c *AtomContext)

	// ExitQuantified_expression is called when exiting the quantified_expression production.
	ExitQuantified_expression(c *Quantified_expressionContext)

	// ExitString_function is called when exiting the string_function production.
	ExitString_function(c *String_functionContext)

	// ExitStandard_function is called when exiting the standard_function production.
	ExitStandard_function(c *Standard_functionContext)

	// ExitJson_function is called when exiting the json_function production.
	ExitJson_function(c *Json_functionContext)

	// ExitJson_object_content is called when exiting the json_object_content production.
	ExitJson_object_content(c *Json_object_contentContext)

	// ExitJson_object_entry is called when exiting the json_object_entry production.
	ExitJson_object_entry(c *Json_object_entryContext)

	// ExitJson_table_clause is called when exiting the json_table_clause production.
	ExitJson_table_clause(c *Json_table_clauseContext)

	// ExitJson_array_element is called when exiting the json_array_element production.
	ExitJson_array_element(c *Json_array_elementContext)

	// ExitJson_on_null_clause is called when exiting the json_on_null_clause production.
	ExitJson_on_null_clause(c *Json_on_null_clauseContext)

	// ExitJson_return_clause is called when exiting the json_return_clause production.
	ExitJson_return_clause(c *Json_return_clauseContext)

	// ExitJson_transform_op is called when exiting the json_transform_op production.
	ExitJson_transform_op(c *Json_transform_opContext)

	// ExitJson_column_clause is called when exiting the json_column_clause production.
	ExitJson_column_clause(c *Json_column_clauseContext)

	// ExitJson_column_definition is called when exiting the json_column_definition production.
	ExitJson_column_definition(c *Json_column_definitionContext)

	// ExitJson_query_returning_clause is called when exiting the json_query_returning_clause production.
	ExitJson_query_returning_clause(c *Json_query_returning_clauseContext)

	// ExitJson_query_return_type is called when exiting the json_query_return_type production.
	ExitJson_query_return_type(c *Json_query_return_typeContext)

	// ExitJson_query_wrapper_clause is called when exiting the json_query_wrapper_clause production.
	ExitJson_query_wrapper_clause(c *Json_query_wrapper_clauseContext)

	// ExitJson_query_on_error_clause is called when exiting the json_query_on_error_clause production.
	ExitJson_query_on_error_clause(c *Json_query_on_error_clauseContext)

	// ExitJson_query_on_empty_clause is called when exiting the json_query_on_empty_clause production.
	ExitJson_query_on_empty_clause(c *Json_query_on_empty_clauseContext)

	// ExitJson_value_return_clause is called when exiting the json_value_return_clause production.
	ExitJson_value_return_clause(c *Json_value_return_clauseContext)

	// ExitJson_value_return_type is called when exiting the json_value_return_type production.
	ExitJson_value_return_type(c *Json_value_return_typeContext)

	// ExitJson_value_on_mismatch_clause is called when exiting the json_value_on_mismatch_clause production.
	ExitJson_value_on_mismatch_clause(c *Json_value_on_mismatch_clauseContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitNumeric_function_wrapper is called when exiting the numeric_function_wrapper production.
	ExitNumeric_function_wrapper(c *Numeric_function_wrapperContext)

	// ExitNumeric_function is called when exiting the numeric_function production.
	ExitNumeric_function(c *Numeric_functionContext)

	// ExitListagg_overflow_clause is called when exiting the listagg_overflow_clause production.
	ExitListagg_overflow_clause(c *Listagg_overflow_clauseContext)

	// ExitOther_function is called when exiting the other_function production.
	ExitOther_function(c *Other_functionContext)

	// ExitOver_clause_keyword is called when exiting the over_clause_keyword production.
	ExitOver_clause_keyword(c *Over_clause_keywordContext)

	// ExitWithin_or_over_clause_keyword is called when exiting the within_or_over_clause_keyword production.
	ExitWithin_or_over_clause_keyword(c *Within_or_over_clause_keywordContext)

	// ExitStandard_prediction_function_keyword is called when exiting the standard_prediction_function_keyword production.
	ExitStandard_prediction_function_keyword(c *Standard_prediction_function_keywordContext)

	// ExitOver_clause is called when exiting the over_clause production.
	ExitOver_clause(c *Over_clauseContext)

	// ExitWindowing_clause is called when exiting the windowing_clause production.
	ExitWindowing_clause(c *Windowing_clauseContext)

	// ExitWindowing_type is called when exiting the windowing_type production.
	ExitWindowing_type(c *Windowing_typeContext)

	// ExitWindowing_elements is called when exiting the windowing_elements production.
	ExitWindowing_elements(c *Windowing_elementsContext)

	// ExitUsing_clause is called when exiting the using_clause production.
	ExitUsing_clause(c *Using_clauseContext)

	// ExitUsing_element is called when exiting the using_element production.
	ExitUsing_element(c *Using_elementContext)

	// ExitAssignable_element is called when exiting the assignable_element production.
	ExitAssignable_element(c *Assignable_elementContext)

	// ExitCollect_order_by_part is called when exiting the collect_order_by_part production.
	ExitCollect_order_by_part(c *Collect_order_by_partContext)

	// ExitWithin_or_over_part is called when exiting the within_or_over_part production.
	ExitWithin_or_over_part(c *Within_or_over_partContext)

	// ExitString_delimiter is called when exiting the string_delimiter production.
	ExitString_delimiter(c *String_delimiterContext)

	// ExitCost_matrix_clause is called when exiting the cost_matrix_clause production.
	ExitCost_matrix_clause(c *Cost_matrix_clauseContext)

	// ExitXml_passing_clause is called when exiting the xml_passing_clause production.
	ExitXml_passing_clause(c *Xml_passing_clauseContext)

	// ExitXml_attributes_clause is called when exiting the xml_attributes_clause production.
	ExitXml_attributes_clause(c *Xml_attributes_clauseContext)

	// ExitXml_namespaces_clause is called when exiting the xml_namespaces_clause production.
	ExitXml_namespaces_clause(c *Xml_namespaces_clauseContext)

	// ExitXml_table_column is called when exiting the xml_table_column production.
	ExitXml_table_column(c *Xml_table_columnContext)

	// ExitXml_general_default_part is called when exiting the xml_general_default_part production.
	ExitXml_general_default_part(c *Xml_general_default_partContext)

	// ExitXml_multiuse_expression_element is called when exiting the xml_multiuse_expression_element production.
	ExitXml_multiuse_expression_element(c *Xml_multiuse_expression_elementContext)

	// ExitXmlroot_param_version_part is called when exiting the xmlroot_param_version_part production.
	ExitXmlroot_param_version_part(c *Xmlroot_param_version_partContext)

	// ExitXmlroot_param_standalone_part is called when exiting the xmlroot_param_standalone_part production.
	ExitXmlroot_param_standalone_part(c *Xmlroot_param_standalone_partContext)

	// ExitXmlserialize_param_enconding_part is called when exiting the xmlserialize_param_enconding_part production.
	ExitXmlserialize_param_enconding_part(c *Xmlserialize_param_enconding_partContext)

	// ExitXmlserialize_param_version_part is called when exiting the xmlserialize_param_version_part production.
	ExitXmlserialize_param_version_part(c *Xmlserialize_param_version_partContext)

	// ExitXmlserialize_param_ident_part is called when exiting the xmlserialize_param_ident_part production.
	ExitXmlserialize_param_ident_part(c *Xmlserialize_param_ident_partContext)

	// ExitAnnotations_clause is called when exiting the annotations_clause production.
	ExitAnnotations_clause(c *Annotations_clauseContext)

	// ExitAnnotations_list is called when exiting the annotations_list production.
	ExitAnnotations_list(c *Annotations_listContext)

	// ExitAnnotation is called when exiting the annotation production.
	ExitAnnotation(c *AnnotationContext)

	// ExitSql_plus_command is called when exiting the sql_plus_command production.
	ExitSql_plus_command(c *Sql_plus_commandContext)

	// ExitStart_command is called when exiting the start_command production.
	ExitStart_command(c *Start_commandContext)

	// ExitSql_plus_filepath is called when exiting the sql_plus_filepath production.
	ExitSql_plus_filepath(c *Sql_plus_filepathContext)

	// ExitWhenever_command is called when exiting the whenever_command production.
	ExitWhenever_command(c *Whenever_commandContext)

	// ExitSet_command is called when exiting the set_command production.
	ExitSet_command(c *Set_commandContext)

	// ExitTiming_command is called when exiting the timing_command production.
	ExitTiming_command(c *Timing_commandContext)

	// ExitClear_command is called when exiting the clear_command production.
	ExitClear_command(c *Clear_commandContext)

	// ExitPartition_extension_clause is called when exiting the partition_extension_clause production.
	ExitPartition_extension_clause(c *Partition_extension_clauseContext)

	// ExitColumn_alias is called when exiting the column_alias production.
	ExitColumn_alias(c *Column_aliasContext)

	// ExitTable_alias is called when exiting the table_alias production.
	ExitTable_alias(c *Table_aliasContext)

	// ExitWhere_clause is called when exiting the where_clause production.
	ExitWhere_clause(c *Where_clauseContext)

	// ExitInto_clause is called when exiting the into_clause production.
	ExitInto_clause(c *Into_clauseContext)

	// ExitXml_column_name is called when exiting the xml_column_name production.
	ExitXml_column_name(c *Xml_column_nameContext)

	// ExitCost_class_name is called when exiting the cost_class_name production.
	ExitCost_class_name(c *Cost_class_nameContext)

	// ExitAttribute_name is called when exiting the attribute_name production.
	ExitAttribute_name(c *Attribute_nameContext)

	// ExitSavepoint_name is called when exiting the savepoint_name production.
	ExitSavepoint_name(c *Savepoint_nameContext)

	// ExitRollback_segment_name is called when exiting the rollback_segment_name production.
	ExitRollback_segment_name(c *Rollback_segment_nameContext)

	// ExitSchema_name is called when exiting the schema_name production.
	ExitSchema_name(c *Schema_nameContext)

	// ExitRoutine_name is called when exiting the routine_name production.
	ExitRoutine_name(c *Routine_nameContext)

	// ExitPackage_name is called when exiting the package_name production.
	ExitPackage_name(c *Package_nameContext)

	// ExitImplementation_type_name is called when exiting the implementation_type_name production.
	ExitImplementation_type_name(c *Implementation_type_nameContext)

	// ExitParameter_name is called when exiting the parameter_name production.
	ExitParameter_name(c *Parameter_nameContext)

	// ExitReference_model_name is called when exiting the reference_model_name production.
	ExitReference_model_name(c *Reference_model_nameContext)

	// ExitMain_model_name is called when exiting the main_model_name production.
	ExitMain_model_name(c *Main_model_nameContext)

	// ExitContainer_tableview_name is called when exiting the container_tableview_name production.
	ExitContainer_tableview_name(c *Container_tableview_nameContext)

	// ExitAggregate_function_name is called when exiting the aggregate_function_name production.
	ExitAggregate_function_name(c *Aggregate_function_nameContext)

	// ExitQuery_name is called when exiting the query_name production.
	ExitQuery_name(c *Query_nameContext)

	// ExitGrantee_name is called when exiting the grantee_name production.
	ExitGrantee_name(c *Grantee_nameContext)

	// ExitRole_name is called when exiting the role_name production.
	ExitRole_name(c *Role_nameContext)

	// ExitConstraint_name is called when exiting the constraint_name production.
	ExitConstraint_name(c *Constraint_nameContext)

	// ExitLabel_name is called when exiting the label_name production.
	ExitLabel_name(c *Label_nameContext)

	// ExitType_name is called when exiting the type_name production.
	ExitType_name(c *Type_nameContext)

	// ExitSequence_name is called when exiting the sequence_name production.
	ExitSequence_name(c *Sequence_nameContext)

	// ExitException_name is called when exiting the exception_name production.
	ExitException_name(c *Exception_nameContext)

	// ExitFunction_name is called when exiting the function_name production.
	ExitFunction_name(c *Function_nameContext)

	// ExitProcedure_name is called when exiting the procedure_name production.
	ExitProcedure_name(c *Procedure_nameContext)

	// ExitTrigger_name is called when exiting the trigger_name production.
	ExitTrigger_name(c *Trigger_nameContext)

	// ExitVariable_name is called when exiting the variable_name production.
	ExitVariable_name(c *Variable_nameContext)

	// ExitIndex_name is called when exiting the index_name production.
	ExitIndex_name(c *Index_nameContext)

	// ExitCursor_name is called when exiting the cursor_name production.
	ExitCursor_name(c *Cursor_nameContext)

	// ExitRecord_name is called when exiting the record_name production.
	ExitRecord_name(c *Record_nameContext)

	// ExitLink_name is called when exiting the link_name production.
	ExitLink_name(c *Link_nameContext)

	// ExitLocal_link_name is called when exiting the local_link_name production.
	ExitLocal_link_name(c *Local_link_nameContext)

	// ExitConnection_qualifier is called when exiting the connection_qualifier production.
	ExitConnection_qualifier(c *Connection_qualifierContext)

	// ExitColumn_name is called when exiting the column_name production.
	ExitColumn_name(c *Column_nameContext)

	// ExitTableview_name is called when exiting the tableview_name production.
	ExitTableview_name(c *Tableview_nameContext)

	// ExitXmltable is called when exiting the xmltable production.
	ExitXmltable(c *XmltableContext)

	// ExitChar_set_name is called when exiting the char_set_name production.
	ExitChar_set_name(c *Char_set_nameContext)

	// ExitSynonym_name is called when exiting the synonym_name production.
	ExitSynonym_name(c *Synonym_nameContext)

	// ExitSchema_object_name is called when exiting the schema_object_name production.
	ExitSchema_object_name(c *Schema_object_nameContext)

	// ExitDir_object_name is called when exiting the dir_object_name production.
	ExitDir_object_name(c *Dir_object_nameContext)

	// ExitUser_object_name is called when exiting the user_object_name production.
	ExitUser_object_name(c *User_object_nameContext)

	// ExitGrant_object_name is called when exiting the grant_object_name production.
	ExitGrant_object_name(c *Grant_object_nameContext)

	// ExitColumn_list is called when exiting the column_list production.
	ExitColumn_list(c *Column_listContext)

	// ExitParen_column_list is called when exiting the paren_column_list production.
	ExitParen_column_list(c *Paren_column_listContext)

	// ExitKeep_clause is called when exiting the keep_clause production.
	ExitKeep_clause(c *Keep_clauseContext)

	// ExitFunction_argument is called when exiting the function_argument production.
	ExitFunction_argument(c *Function_argumentContext)

	// ExitFunction_argument_analytic is called when exiting the function_argument_analytic production.
	ExitFunction_argument_analytic(c *Function_argument_analyticContext)

	// ExitFunction_argument_modeling is called when exiting the function_argument_modeling production.
	ExitFunction_argument_modeling(c *Function_argument_modelingContext)

	// ExitRespect_or_ignore_nulls is called when exiting the respect_or_ignore_nulls production.
	ExitRespect_or_ignore_nulls(c *Respect_or_ignore_nullsContext)

	// ExitArgument is called when exiting the argument production.
	ExitArgument(c *ArgumentContext)

	// ExitType_spec is called when exiting the type_spec production.
	ExitType_spec(c *Type_specContext)

	// ExitDatatype is called when exiting the datatype production.
	ExitDatatype(c *DatatypeContext)

	// ExitPrecision_part is called when exiting the precision_part production.
	ExitPrecision_part(c *Precision_partContext)

	// ExitNative_datatype_element is called when exiting the native_datatype_element production.
	ExitNative_datatype_element(c *Native_datatype_elementContext)

	// ExitBind_variable is called when exiting the bind_variable production.
	ExitBind_variable(c *Bind_variableContext)

	// ExitGeneral_element is called when exiting the general_element production.
	ExitGeneral_element(c *General_elementContext)

	// ExitGeneral_element_part is called when exiting the general_element_part production.
	ExitGeneral_element_part(c *General_element_partContext)

	// ExitTable_element is called when exiting the table_element production.
	ExitTable_element(c *Table_elementContext)

	// ExitObject_privilege is called when exiting the object_privilege production.
	ExitObject_privilege(c *Object_privilegeContext)

	// ExitSystem_privilege is called when exiting the system_privilege production.
	ExitSystem_privilege(c *System_privilegeContext)

	// ExitConstant is called when exiting the constant production.
	ExitConstant(c *ConstantContext)

	// ExitNumeric is called when exiting the numeric production.
	ExitNumeric(c *NumericContext)

	// ExitNumeric_negative is called when exiting the numeric_negative production.
	ExitNumeric_negative(c *Numeric_negativeContext)

	// ExitQuoted_string is called when exiting the quoted_string production.
	ExitQuoted_string(c *Quoted_stringContext)

	// ExitIdentifier is called when exiting the identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitId_expression is called when exiting the id_expression production.
	ExitId_expression(c *Id_expressionContext)

	// ExitInquiry_directive is called when exiting the inquiry_directive production.
	ExitInquiry_directive(c *Inquiry_directiveContext)

	// ExitOuter_join_sign is called when exiting the outer_join_sign production.
	ExitOuter_join_sign(c *Outer_join_signContext)

	// ExitRegular_id is called when exiting the regular_id production.
	ExitRegular_id(c *Regular_idContext)

	// ExitNon_reserved_keywords_in_18c is called when exiting the non_reserved_keywords_in_18c production.
	ExitNon_reserved_keywords_in_18c(c *Non_reserved_keywords_in_18cContext)

	// ExitNon_reserved_keywords_in_12c is called when exiting the non_reserved_keywords_in_12c production.
	ExitNon_reserved_keywords_in_12c(c *Non_reserved_keywords_in_12cContext)

	// ExitNon_reserved_keywords_pre12c is called when exiting the non_reserved_keywords_pre12c production.
	ExitNon_reserved_keywords_pre12c(c *Non_reserved_keywords_pre12cContext)
}
