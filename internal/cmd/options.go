package cmd

import (
	"io"

	"google.golang.org/grpc"

	"github.com/sqlc-dev/sqlc/internal/config"
	pb "github.com/sqlc-dev/sqlc/pkg/engine"
)

type Options struct {
	Env    Env
	Stderr io.Writer
	// TODO: Move these to a command-specific struct
	Tags    []string
	Against string

	// Testing only
	MutateConfig func(*config.Config)
	// CodegenHandlerOverride injects a mock codegen handler instead of spawning a process.
	CodegenHandlerOverride grpc.ClientConnInterface
	// PluginParseFunc, when set, is used in the plugin-engine path instead of invoking the engine process (for tests).
	PluginParseFunc func(schemaSQL, querySQL string) (*pb.ParseResponse, error)
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
