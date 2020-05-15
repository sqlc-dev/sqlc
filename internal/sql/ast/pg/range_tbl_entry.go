package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RangeTblEntry struct {
	Rtekind         RTEKind
	Relid           Oid
	Relkind         byte
	Tablesample     *TableSampleClause
	Subquery        *Query
	SecurityBarrier bool
	Jointype        JoinType
	Joinaliasvars   *ast.List
	Functions       *ast.List
	Funcordinality  bool
	Tablefunc       *TableFunc
	ValuesLists     *ast.List
	Ctename         *string
	Ctelevelsup     Index
	SelfReference   bool
	Coltypes        *ast.List
	Coltypmods      *ast.List
	Colcollations   *ast.List
	Enrname         *string
	Enrtuples       float64
	Alias           *Alias
	Eref            *Alias
	Lateral         bool
	Inh             bool
	InFromCl        bool
	RequiredPerms   AclMode
	CheckAsUser     Oid
	SelectedCols    []uint32
	InsertedCols    []uint32
	UpdatedCols     []uint32
	SecurityQuals   *ast.List
}

func (n *RangeTblEntry) Pos() int {
	return 0
}
