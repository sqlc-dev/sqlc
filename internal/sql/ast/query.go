package ast

type Query struct {
	Tag NodeTag[Query] `json:"tag"`

	CommandType      CmdType         `json:"command_type"`
	QuerySource      QuerySource     `json:"query_source"`
	QueryId          uint32          `json:"query_id"`
	CanSetTag        bool            `json:"can_set_tag"`
	UtilityStmt      Node            `json:"utility_stmt,omitempty"`
	ResultRelation   int             `json:"result_relation"`
	HasAggs          bool            `json:"has_aggs"`
	HasWindowFuncs   bool            `json:"has_window_funcs"`
	HasTargetSrfs    bool            `json:"has_target_srfs"`
	HasSubLinks      bool            `json:"has_sub_links"`
	HasDistinctOn    bool            `json:"has_distinct_on"`
	HasRecursive     bool            `json:"has_recursive"`
	HasModifyingCte  bool            `json:"has_modifying_cte"`
	HasForUpdate     bool            `json:"has_for_update"`
	HasRowSecurity   bool            `json:"has_row_security"`
	CteList          *List           `json:"cte_list,omitempty"`
	Rtable           *List           `json:"rtable,omitempty"`
	Jointree         *FromExpr       `json:"jointree,omitempty"`
	TargetList       *List           `json:"target_list,omitempty"`
	Override         OverridingKind  `json:"override"`
	OnConflict       *OnConflictExpr `json:"on_conflict,omitempty"`
	ReturningList    *List           `json:"returning_list,omitempty"`
	GroupClause      *List           `json:"group_clause,omitempty"`
	GroupingSets     *List           `json:"grouping_sets,omitempty"`
	HavingQual       Node            `json:"having_qual,omitempty"`
	WindowClause     *List           `json:"window_clause,omitempty"`
	DistinctClause   *List           `json:"distinct_clause,omitempty"`
	SortClause       *List           `json:"sort_clause,omitempty"`
	LimitOffset      Node            `json:"limit_offset,omitempty"`
	LimitCount       Node            `json:"limit_count,omitempty"`
	RowMarks         *List           `json:"row_marks,omitempty"`
	SetOperations    Node            `json:"set_operations,omitempty"`
	ConstraintDeps   *List           `json:"constraint_deps,omitempty"`
	WithCheckOptions *List           `json:"with_check_options,omitempty"`
	StmtLocation     int             `json:"stmt_location"`
	StmtLen          int             `json:"stmt_len"`
}

func (n *Query) Pos() int {
	return 0
}
