package core

// TODO: This is the last struct left over from the old architecture. Figure
// out how to remove it at some point
type FQN struct {
	Catalog string
	Schema  string
	Rel     string
}
