package golang

import (
	core "github.com/kyleconroy/sqlc/internal/pg"
)

type Struct struct {
	Table   core.FQN
	Name    string
	Fields  []Field
	Comment string
}
