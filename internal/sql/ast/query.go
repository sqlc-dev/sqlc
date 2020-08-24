package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Query struct {
	CommandType      CmdType
	QuerySource      QuerySource
	QueryId          uint32
	CanSetTag        bool
	UtilityStmt      ast.Node
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
	CteList          *List
	Rtable           *List
	Jointree         *FromExpr
	TargetList       *List
	Override         OverridingKind
	OnConflict       *OnConflictExpr
	ReturningList    *List
	GroupClause      *List
	GroupingSets     *List
	HavingQual       ast.Node
	WindowClause     *List
	DistinctClause   *List
	SortClause       *List
	LimitOffset      ast.Node
	LimitCount       ast.Node
	RowMarks         *List
	SetOperations    ast.Node
	ConstraintDeps   *List
	WithCheckOptions *List
	StmtLocation     int
	StmtLen          int
}

func (n *Query) Pos() int {
	return 0
}
