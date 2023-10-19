package rpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errMessageUnauthenticated = `rpc authentication failed

You may be using a sqlc auth token that was created for a different project,
or your auth token may have expired.`

var ErrUnauthenticated = status.New(codes.Unauthenticated, errMessageUnauthenticated).Err()
