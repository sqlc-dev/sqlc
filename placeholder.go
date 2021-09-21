package sqlc

import (
	"io"

	"github.com/kyleconroy/sqlc/internal/cmd"
)

// This is a dummy file that allows SQLC to be "installed" as a module and locked using
// go.mod and then run using "go run github.com/kyleconroy/sqlc"

type Placeholder struct{}

// Run the sqlc command
func Run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	return cmd.Do(args, stdin, stdout, stderr)
}
