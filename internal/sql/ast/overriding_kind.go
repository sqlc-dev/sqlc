package ast

type OverridingKind uint

// Values match the PostgreSQL OverridingKind enum
const (
	OverridingNotSet      OverridingKind = 1
	OverridingUserValue   OverridingKind = 2
	OverridingSystemValue OverridingKind = 3
)

func (n *OverridingKind) Pos() int {
	return 0
}
