package ext

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Handler interface {
	Generate(context.Context, *plugin.GenerateRequest) (*plugin.GenerateResponse, error)

	Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

type wrapper struct {
	fn func(context.Context, *plugin.GenerateRequest) (*plugin.GenerateResponse, error)
}

func (w *wrapper) Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	return w.fn(ctx, req)
}

func (w *wrapper) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	req, ok := args.(*plugin.GenerateRequest)
	if !ok {
		return fmt.Errorf("args isn't a GenerateRequest")
	}
	resp, ok := reply.(*plugin.GenerateResponse)
	if !ok {
		return fmt.Errorf("reply isn't a GenerateResponse")
	}
	res, err := w.Generate(ctx, req)
	if err != nil {
		return err
	}
	resp.Files = res.Files
	return nil
}

func (w *wrapper) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func HandleFunc(fn func(context.Context, *plugin.GenerateRequest) (*plugin.GenerateResponse, error)) Handler {
	return &wrapper{fn}
}
