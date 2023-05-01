package remote

import (
	"crypto/tls"
	"os"

	"github.com/riza-io/grpc-go/credentials/bearer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/kyleconroy/sqlc/internal/config"
)

const defaultHostname = "remote.sqlc.dev"

func NewClient(cloudConfig config.Cloud) (GenClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(bearer.NewPerRPCCredentials(os.Getenv("SQLC_AUTH_TOKEN"))),
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
