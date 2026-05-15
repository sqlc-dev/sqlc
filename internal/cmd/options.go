package cmd

import (
	"io"

	"github.com/sqlc-dev/sqlc/internal/config"
)

type Options struct {
	Env    Env
	Stderr io.Writer
	// TODO: Move these to a command-specific struct
	Tags    []string
	Against string
}

func (o *Options) ReadConfig(dir, filename string) (string, *config.Config, error) {
	return readConfig(o.Stderr, dir, filename)
}
