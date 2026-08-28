package ast

type TriggerTransition struct {
	Tag NodeTag[TriggerTransition] `json:"tag"`

	Name    *string `json:"name,omitempty"`
	IsNew   bool    `json:"is_new"`
	IsTable bool    `json:"is_table"`
}

func (n *TriggerTransition) Pos() int {
	return 0
}
