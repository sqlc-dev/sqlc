package ast

type TriggerTransition struct {
	Tag NodeTag[TriggerTransition] `json:"tag"`

	Name    *string
	IsNew   bool
	IsTable bool
}

func (n *TriggerTransition) Pos() int {
	return 0
}
