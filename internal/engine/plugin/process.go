// Package plugin implements running database-engine plugins as external processes.
//
// It is used only by the generate path (cmd runPluginQuerySet): schema and queries
// are sent via ParseRequest to the plugin; the compiler is not used for plugin engines.
// Vet does not support plugin engines.
package plugin

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

// ProcessRunner runs an engine plugin as an external process.
type ProcessRunner struct {
	Cmd string
	Dir string // Working directory for the plugin (config file directory)
	Env []string
}

// NewProcessRunner creates a new ProcessRunner.
func NewProcessRunner(cmd, dir string, env []string) *ProcessRunner {
	return &ProcessRunner{
		Cmd: cmd,
		Dir: dir,
		Env: env,
	}
}

func (r *ProcessRunner) invoke(ctx context.Context, method string, req, resp proto.Message) error {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	// Parse command string to support formats like "go run ./path"
	cmdParts := strings.Fields(r.Cmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	// Build arguments: rest of cmdParts + method
	args := append(cmdParts[1:], method)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = bytes.NewReader(stdin)
	// Set working directory to config file directory for relative paths
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	// Inherit the current environment and add SQLC_VERSION
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

// ParseRequest invokes the plugin's Parse RPC with the given request (sql and optional schema_sql).
// The cmd layer uses this for the plugin-engine generate path instead of the compiler.
func (r *ProcessRunner) ParseRequest(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	resp := &pb.ParseResponse{}
	if err := r.invoke(ctx, "parse", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
