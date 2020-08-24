package pg

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
	CteList          *ast.List
	Rtable           *ast.List
	Jointree         *FromExpr
	TargetList       *ast.List
	Override         OverridingKind
	OnConflict       *OnConflictExpr
	ReturningList    *ast.List
	GroupClause      *ast.List
	GroupingSets     *ast.List
	HavingQual       ast.Node
	WindowClause     *ast.List
	DistinctClause   *ast.List
	SortClause       *ast.List
	LimitOffset      ast.Node
	LimitCount       ast.Node
	RowMarks         *ast.List
	SetOperations    ast.Node
	ConstraintDeps   *ast.List
	WithCheckOptions *ast.List
	StmtLocation     int
	StmtLen          int
}

func (n *Query) Pos() int {
	return 0
}
