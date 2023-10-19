package remote

import (
	"crypto/tls"

	"github.com/riza-io/grpc-go/credentials/basic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/rpc"
)

const defaultHostname = "remote.sqlc.dev"

func NewClient(cloudConfig config.Cloud) (GenClient, error) {
	authID := cloudConfig.Organization + "/" + cloudConfig.Project
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(basic.NewPerRPCCredentials(authID, cloudConfig.AuthToken)),
		grpc.WithUnaryInterceptor(rpc.UnaryInterceptor),
	}

	hostname := cloudConfig.Hostname
	if hostname == "" {
		hostname = defaultHostname
	}

	conn, err := grpc.Dial(hostname+":443", opts...)
	if err != nil {
		return nil, err
	}

	return NewGenClient(conn), nil
}
