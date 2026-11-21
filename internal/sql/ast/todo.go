package ast

type TODO struct {
	Tag                NodeTag[TODO] `json:"tag"`
	UnsupportedFeature string        `json:"unsupportedFeature,omitempty"`
}

func (n *TODO) Pos() int {
	return 0
}
