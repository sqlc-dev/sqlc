package ast

type RangeTblEntry struct {
	Tag NodeTag[RangeTblEntry] `json:"tag"`

	Rtekind         RTEKind
	Relid           Oid
	Relkind         byte
	Tablesample     *TableSampleClause `json:",omitempty"`
	Subquery        *Query             `json:",omitempty"`
	SecurityBarrier bool
	Jointype        JoinType
	Joinaliasvars   *List `json:",omitempty"`
	Functions       *List `json:",omitempty"`
	Funcordinality  bool
	Tablefunc       *TableFunc `json:",omitempty"`
	ValuesLists     *List      `json:",omitempty"`
	Ctename         *string    `json:",omitempty"`
	Ctelevelsup     Index
	SelfReference   bool
	Coltypes        *List   `json:",omitempty"`
	Coltypmods      *List   `json:",omitempty"`
	Colcollations   *List   `json:",omitempty"`
	Enrname         *string `json:",omitempty"`
	Enrtuples       float64
	Alias           *Alias `json:",omitempty"`
	Eref            *Alias `json:",omitempty"`
	Lateral         bool
	Inh             bool
	InFromCl        bool
	RequiredPerms   AclMode
	CheckAsUser     Oid
	SelectedCols    []uint32 `json:",omitempty"`
	InsertedCols    []uint32 `json:",omitempty"`
	UpdatedCols     []uint32 `json:",omitempty"`
	SecurityQuals   *List    `json:",omitempty"`
}

func (n *RangeTblEntry) Pos() int {
	return 0
}
