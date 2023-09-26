package main

import (
	"os"

	"github.com/Dionid/sqlc/internal/cmd"
)

func main() {
	os.Exit(cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
