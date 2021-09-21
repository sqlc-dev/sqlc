// package cli exposes the command-line interface for sqlc. It can be used to
// run sqlc from Go without the overhead of creating a child process.
//
// Example usage:
//
//     package main
//
//     import (
//         "os"
//
//         sqlc "github.com/kyleconroy/sqlc/pkg/cli"
//     )
//
//     func main() {
//         os.Exit(sqlc.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
//     }
//
package cli

import (
	"io"

	"github.com/kyleconroy/sqlc/internal/cmd"
)

// Run the sqlc CLI. It takes an array of command-line arguments
// (excluding the executable argument itself) and returns an exit
// code.
func Run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	return cmd.Do(args, stdin, stdout, stderr)
}
