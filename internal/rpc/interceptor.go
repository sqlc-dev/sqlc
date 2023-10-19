package rpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errMessageUnauthenticated = `rpc authentication failed

You may be using a sqlc auth token that was created for a different project,
or your auth token may have expired.`

func UnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)

	switch status.Convert(err).Code() {
	case codes.OK:
		return nil
	case codes.Unauthenticated:
		return status.New(codes.Unauthenticated, errMessageUnauthenticated).Err()
	default:
		return err
	}
}
