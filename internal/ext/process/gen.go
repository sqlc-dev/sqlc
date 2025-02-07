package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

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

	// Check if the output plugin exists
	path, err := exec.LookPath(r.Cmd)
	if err != nil {
		return fmt.Errorf("process: %s not found", r.Cmd)
	}

	cmd := exec.CommandContext(ctx, path, method)
	cmd.Stdin = bytes.NewReader(stdin)
	cmd.Env = []string{
		fmt.Sprintf("SQLC_VERSION=%s", info.Version),
	}
	for _, key := range r.Env {
		if key == "SQLC_AUTH_TOKEN" {
			continue
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, os.Getenv(key)))
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
