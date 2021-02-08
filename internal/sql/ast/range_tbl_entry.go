package ast

type RangeTblEntry struct {
	Rtekind         RTEKind
	Relid           Oid
	Relkind         byte
	Tablesample     *TableSampleClause
	Subquery        *Query
	SecurityBarrier bool
	Jointype        JoinType
	Joinaliasvars   *List
	Functions       *List
	Funcordinality  bool
	Tablefunc       *TableFunc
	ValuesLists     *List
	Ctename         *string
	Ctelevelsup     Index
	SelfReference   bool
	Coltypes        *List
	Coltypmods      *List
	Colcollations   *List
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
	SecurityQuals   *List
}

func (n *RangeTblEntry) Pos() int {
	return 0
}
