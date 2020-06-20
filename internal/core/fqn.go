package core

// TODO: This is the last struct left over from the old architecture. Figure
// out how to remove it at some point
type FQN struct {
	Catalog string
	Schema  string
	Rel     string
}

func (f FQN) String() string {
	s := f.Rel
	if f.Schema != "" {
		s = f.Schema + "." + s
	}
	if f.Catalog != "" {
		s = f.Catalog + "." + s
	}
	return s
}
