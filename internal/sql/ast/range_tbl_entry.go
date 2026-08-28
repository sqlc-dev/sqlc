package ast

type RangeTblEntry struct {
	Tag NodeTag[RangeTblEntry] `json:"tag"`

	Rtekind         RTEKind            `json:"rtekind"`
	Relid           Oid                `json:"relid"`
	Relkind         byte               `json:"relkind"`
	Tablesample     *TableSampleClause `json:"tablesample,omitempty"`
	Subquery        *Query             `json:"subquery,omitempty"`
	SecurityBarrier bool               `json:"security_barrier"`
	Jointype        JoinType           `json:"jointype"`
	Joinaliasvars   *List              `json:"joinaliasvars,omitempty"`
	Functions       *List              `json:"functions,omitempty"`
	Funcordinality  bool               `json:"funcordinality"`
	Tablefunc       *TableFunc         `json:"tablefunc,omitempty"`
	ValuesLists     *List              `json:"values_lists,omitempty"`
	Ctename         *string            `json:"ctename,omitempty"`
	Ctelevelsup     Index              `json:"ctelevelsup"`
	SelfReference   bool               `json:"self_reference"`
	Coltypes        *List              `json:"coltypes,omitempty"`
	Coltypmods      *List              `json:"coltypmods,omitempty"`
	Colcollations   *List              `json:"colcollations,omitempty"`
	Enrname         *string            `json:"enrname,omitempty"`
	Enrtuples       float64            `json:"enrtuples"`
	Alias           *Alias             `json:"alias,omitempty"`
	Eref            *Alias             `json:"eref,omitempty"`
	Lateral         bool               `json:"lateral"`
	Inh             bool               `json:"inh"`
	InFromCl        bool               `json:"in_from_cl"`
	RequiredPerms   AclMode            `json:"required_perms"`
	CheckAsUser     Oid                `json:"check_as_user"`
	SelectedCols    []uint32           `json:"selected_cols,omitempty"`
	InsertedCols    []uint32           `json:"inserted_cols,omitempty"`
	UpdatedCols     []uint32           `json:"updated_cols,omitempty"`
	SecurityQuals   *List              `json:"security_quals,omitempty"`
}

func (n *RangeTblEntry) Pos() int {
	return 0
}
