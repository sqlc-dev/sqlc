package ast

type ParamListInfoData struct {
	Tag NodeTag[ParamListInfoData] `json:"tag"`

	ParamFetchArg  any
	ParserSetupArg any
	NumParams      int
	ParamMask      []uint32
}

func (n *ParamListInfoData) Pos() int {
	return 0
}
