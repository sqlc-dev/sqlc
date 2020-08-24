package ast

type TriggerTransition struct {
	Name    *string
	IsNew   bool
	IsTable bool
}

func (n *TriggerTransition) Pos() int {
	return 0
}
