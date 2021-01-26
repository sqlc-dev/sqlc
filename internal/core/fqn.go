package core

import "github.com/gobwas/glob"

// TODO: This is the last struct left over from the old architecture. Figure
// out how to remove it at some point
type FQN struct {
	Catalog string
	Schema  glob.Glob
	Rel     glob.Glob
}
