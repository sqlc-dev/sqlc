package pg

type OverridingKind uint

func (n *OverridingKind) Pos() int {
	return 0
}
