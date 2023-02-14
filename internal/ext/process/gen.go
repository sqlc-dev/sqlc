package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/kyleconroy/sqlc/internal/plugin"
)

type Runner struct {
	Cmd string
}

// TODO: Update the gen func signature to take a ctx
func (r Runner) Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode codegen request: %s", err)
	}

	cmdAndArgs := strings.Split(r.Cmd, " ")

	// Check if the output plugin exists
	path, err := exec.LookPath(cmdAndArgs[0])
	if err != nil {
		return nil, fmt.Errorf("process: %s not found", r.Cmd)
	}

	cmd := exec.CommandContext(ctx, path, cmdAndArgs[1:]...)
	cmd.Stdin = bytes.NewReader(stdin)
	// Make sure to propagate current env before modifying it.
	env := make([]string, len(os.Environ()))
	copy(env, os.Environ())
	env = append(env, fmt.Sprintf("SQLC_VERSION=%s", req.SqlcVersion))
	cmd.Env = env

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
