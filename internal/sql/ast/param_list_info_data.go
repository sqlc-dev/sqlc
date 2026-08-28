package ast

type ParamListInfoData struct {
	Tag NodeTag[ParamListInfoData] `json:"tag"`

	ParamFetchArg  any `json:",omitempty"`
	ParserSetupArg any `json:",omitempty"`
	NumParams      int
	ParamMask      []uint32 `json:",omitempty"`
}

func (n *ParamListInfoData) Pos() int {
	return 0
}
