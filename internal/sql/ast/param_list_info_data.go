package ast

type ParamListInfoData struct {
	ParamFetchArg  any
	ParserSetupArg any
	NumParams      int
	ParamMask      []uint32
}

func (n *ParamListInfoData) Pos() int {
	return 0
}
