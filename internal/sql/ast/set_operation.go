package ast

import "strconv"

const (
	None SetOperation = iota
	Union
	Intersect
	Except
)

type SetOperation uint

func (n *SetOperation) Pos() int {
	return 0
}

func (n SetOperation) String() string {
	switch n {
	case None:
		return "None"
	case Union:
		return "Union"
	case Intersect:
		return "Intersect"
	case Except:
		return "Except"
	default:
		return "Unknown(" + strconv.FormatUint(uint64(n), 10) + ")"
	}
}
