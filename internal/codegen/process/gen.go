package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"

	"google.golang.org/protobuf/proto"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Runner struct {
	Config config.Config
	Plugin string
}

func (r Runner) pluginCmd() (string, error) {
	for _, plug := range r.Config.Plugins {
		if plug.Name != r.Plugin {
			continue
		}
		if plug.Process == nil {
			continue
		}
		return plug.Process.Cmd, nil
	}
	return "", fmt.Errorf("plugin not found")
}

// TODO: Update the gen func signature to take a ctx
func (r Runner) Generate(req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode codegen request: %s", err)
	}

	name, err := r.pluginCmd()
	if err != nil {
		return nil, fmt.Errorf("process: unknown plugin %s", r.Plugin)
	}

	// Check if the output plugin exists
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("process: %s not found", name)
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = bytes.NewReader(stdin)
	cmd.Env = []string{
		fmt.Sprintf("SQLC_VERSION=%s", req.SqlcVersion),
	}

	out, err := cmd.Output()
	if err != nil {
		stderr := err.Error()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			stderr = string(exit.Stderr)
		}
		return nil, fmt.Errorf("process: error running command %s", stderr)
	}

	var resp plugin.CodeGenResponse
	if err := proto.Unmarshal(out, &resp); err != nil {
		return nil, fmt.Errorf("process: failed to read codegen resp: %s", err)
	}

	return &resp, nil
}
