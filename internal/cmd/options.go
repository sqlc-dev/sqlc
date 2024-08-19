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

	// Testing only
	MutateConfig func(*config.Config)
}

func (o *Options) ReadConfig(dir, filename string) (string, *config.Config, error) {
	path, conf, err := readConfig(o.Stderr, dir, filename)
	if err != nil {
		return path, conf, err
	}
	if o.MutateConfig != nil {
		o.MutateConfig(conf)
	}
	return path, conf, nil
}
