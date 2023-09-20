package quickdb

import (
	"crypto/tls"
	"os"

	"github.com/riza-io/grpc-go/credentials/basic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sqlc-dev/sqlc/internal/config"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
)

const defaultHostname = "grpc.sqlc.dev"

func NewClientFromConfig(cloudConfig config.Cloud) (pb.QuickClient, error) {
	projectID := cloudConfig.Project
	authToken := os.Getenv("SQLC_AUTH_TOKEN")
	return NewClient(projectID, authToken, WithHost(cloudConfig.Hostname))
}

type options struct {
	hostname string
}

type Option func(*options)

func WithHost(host string) Option {
	return func(o *options) {
		o.hostname = host
	}
}

func NewClient(project, token string, opts ...Option) (pb.QuickClient, error) {
	var o options
	for _, apply := range opts {
		apply(&o)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(basic.NewPerRPCCredentials(project, token)),
	}

	hostname := o.hostname
	if hostname == "" {
		hostname = defaultHostname
	}

	conn, err := grpc.Dial(hostname+":443", dialOpts...)
	if err != nil {
		return nil, err
	}

	return pb.NewQuickClient(conn), nil
}
