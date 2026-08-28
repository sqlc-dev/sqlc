package ast

type ParamListInfoData struct {
	Tag NodeTag[ParamListInfoData] `json:"tag"`

	ParamFetchArg  any      `json:"param_fetch_arg,omitempty"`
	ParserSetupArg any      `json:"parser_setup_arg,omitempty"`
	NumParams      int      `json:"num_params"`
	ParamMask      []uint32 `json:"param_mask,omitempty"`
}

func (n *ParamListInfoData) Pos() int {
	return 0
}
