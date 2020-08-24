package ast

type ParamListInfoData struct {
	ParamFetchArg  interface{}
	ParserSetupArg interface{}
	NumParams      int
	ParamMask      []uint32
}

func (n *ParamListInfoData) Pos() int {
	return 0
}
