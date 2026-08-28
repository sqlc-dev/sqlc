package ast

type AlterTSConfigurationStmt struct {
	Tag NodeTag[AlterTSConfigurationStmt] `json:"tag"`

	Kind      AlterTSConfigType `json:"kind"`
	Cfgname   *List             `json:"cfgname,omitempty"`
	Tokentype *List             `json:"tokentype,omitempty"`
	Dicts     *List             `json:"dicts,omitempty"`
	Override  bool              `json:"override"`
	Replace   bool              `json:"replace"`
	MissingOk bool              `json:"missing_ok"`
}

func (n *AlterTSConfigurationStmt) Pos() int {
	return 0
}
