package ast

type AlterTSConfigurationStmt struct {
	Tag NodeTag[AlterTSConfigurationStmt] `json:"tag"`

	Kind      AlterTSConfigType
	Cfgname   *List `json:",omitempty"`
	Tokentype *List `json:",omitempty"`
	Dicts     *List `json:",omitempty"`
	Override  bool
	Replace   bool
	MissingOk bool
}

func (n *AlterTSConfigurationStmt) Pos() int {
	return 0
}
