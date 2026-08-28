package ast

type Query struct {
	Tag NodeTag[Query] `json:"tag"`

	CommandType      CmdType
	QuerySource      QuerySource
	QueryId          uint32
	CanSetTag        bool
	UtilityStmt      Node `json:",omitempty"`
	ResultRelation   int
	HasAggs          bool
	HasWindowFuncs   bool
	HasTargetSrfs    bool
	HasSubLinks      bool
	HasDistinctOn    bool
	HasRecursive     bool
	HasModifyingCte  bool
	HasForUpdate     bool
	HasRowSecurity   bool
	CteList          *List     `json:",omitempty"`
	Rtable           *List     `json:",omitempty"`
	Jointree         *FromExpr `json:",omitempty"`
	TargetList       *List     `json:",omitempty"`
	Override         OverridingKind
	OnConflict       *OnConflictExpr `json:",omitempty"`
	ReturningList    *List           `json:",omitempty"`
	GroupClause      *List           `json:",omitempty"`
	GroupingSets     *List           `json:",omitempty"`
	HavingQual       Node            `json:",omitempty"`
	WindowClause     *List           `json:",omitempty"`
	DistinctClause   *List           `json:",omitempty"`
	SortClause       *List           `json:",omitempty"`
	LimitOffset      Node            `json:",omitempty"`
	LimitCount       Node            `json:",omitempty"`
	RowMarks         *List           `json:",omitempty"`
	SetOperations    Node            `json:",omitempty"`
	ConstraintDeps   *List           `json:",omitempty"`
	WithCheckOptions *List           `json:",omitempty"`
	StmtLocation     int
	StmtLen          int
}

func (n *Query) Pos() int {
	return 0
}
