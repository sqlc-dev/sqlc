package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sqlc-dev/sqlc/internal/info"
)

type Runner struct {
	Cmd    string
	Dir    string // Working directory for the plugin (config file directory)
	Format string
	Env    []string
}

func (r *Runner) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	req, ok := args.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("args isn't a protoreflect.ProtoMessage")
	}

	var stdin []byte
	var err error
	switch r.Format {
	case "json":
		m := &protojson.MarshalOptions{
			EmitUnpopulated: true,
			Indent:          "",
			UseProtoNames:   true,
		}
		stdin, err = m.Marshal(req)

		if err != nil {
			return fmt.Errorf("failed to encode codegen request: %w", err)
		}
	case "", "protobuf":
		stdin, err = proto.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to encode codegen request: %w", err)
		}
	default:
		return fmt.Errorf("unknown plugin format: %s", r.Format)
	}

	// Parse command string to support formats like "go run ./path"
	cmdParts := strings.Fields(r.Cmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("process: %s not found", r.Cmd)
	}

	// Check if the output plugin exists
	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("process: %s not found", r.Cmd)
	}

	// Build arguments: rest of cmdParts + method
	cmdArgs := append(cmdParts[1:], method)
	cmd := exec.CommandContext(ctx, path, cmdArgs...)
	cmd.Stdin = bytes.NewReader(stdin)
	// Set working directory to config file directory for relative paths
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	// Pass only SQLC_VERSION and explicitly configured environment variables
	cmd.Env = []string{
		fmt.Sprintf("SQLC_VERSION=%s", info.Version),
	}
	for _, key := range r.Env {
		if key == "SQLC_AUTH_TOKEN" {
			continue
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, os.Getenv(key)))
	}
	// For "go run" commands, inherit PATH and Go-related environment
	if len(cmdParts) > 1 && cmdParts[0] == "go" {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "PATH=") ||
				strings.HasPrefix(env, "GOPATH=") ||
				strings.HasPrefix(env, "GOROOT=") ||
				strings.HasPrefix(env, "GOWORK=") ||
				strings.HasPrefix(env, "HOME=") ||
				strings.HasPrefix(env, "GOCACHE=") ||
				strings.HasPrefix(env, "GOMODCACHE=") {
				cmd.Env = append(cmd.Env, env)
			}
		}
	}

	out, err := cmd.Output()
	if err != nil {
		stderr := err.Error()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			stderr = string(exit.Stderr)
		}
		return fmt.Errorf("process: error running command %s", stderr)
	}

	resp, ok := reply.(protoreflect.ProtoMessage)
	if !ok {
		return fmt.Errorf("reply isn't a protoreflect.ProtoMessage")
	}

	switch r.Format {
	case "json":
		if err := protojson.Unmarshal(out, resp); err != nil {
			return fmt.Errorf("process: failed to read codegen resp: %w", err)
		}
	default:
		if err := proto.Unmarshal(out, resp); err != nil {
			return fmt.Errorf("process: failed to read codegen resp: %w", err)
		}
	}

	return nil
}

func (r *Runner) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
