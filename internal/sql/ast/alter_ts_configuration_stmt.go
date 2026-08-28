package ast

type AlterTSConfigurationStmt struct {
	Tag NodeTag[AlterTSConfigurationStmt] `json:"tag"`

	Kind      AlterTSConfigType
	Cfgname   *List
	Tokentype *List
	Dicts     *List
	Override  bool
	Replace   bool
	MissingOk bool
}

func (n *AlterTSConfigurationStmt) Pos() int {
	return 0
}
