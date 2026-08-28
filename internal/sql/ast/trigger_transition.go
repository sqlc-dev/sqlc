package ast

type TriggerTransition struct {
	Tag NodeTag[TriggerTransition] `json:"tag"`

	Name    *string `json:",omitempty"`
	IsNew   bool
	IsTable bool
}

func (n *TriggerTransition) Pos() int {
	return 0
}
