// This file runs a database-engine plugin as an external process (parse RPC over stdin/stdout).
// It is used only by the plugin-engine generate path (runPluginQuerySet). Vet does not support plugin engines.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/info"
	pb "github.com/sqlc-dev/sqlc/pkg/engine"
)

// engineProcessRunner runs an engine plugin as an external process.
type engineProcessRunner struct {
	Cmd string
	Dir string // Working directory for the plugin (config file directory)
	Env []string
}

func newEngineProcessRunner(cmd, dir string, env []string) *engineProcessRunner {
	return &engineProcessRunner{Cmd: cmd, Dir: dir, Env: env}
}

func (r *engineProcessRunner) invoke(ctx context.Context, method string, req, resp proto.Message) error {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	cmdParts := strings.Fields(r.Cmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	args := append(cmdParts[1:], method)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = bytes.NewReader(stdin)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("SQLC_VERSION=%s", info.Version))

	out, err := cmd.Output()
	if err != nil {
		stderr := err.Error()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			stderr = string(exit.Stderr)
		}
		return fmt.Errorf("engine plugin error: %s", stderr)
	}

	if err := proto.Unmarshal(out, resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

// parseRequest invokes the plugin's Parse RPC. Used by runPluginQuerySet.
func (r *engineProcessRunner) parseRequest(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	resp := &pb.ParseResponse{}
	if err := r.invoke(ctx, "parse", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
