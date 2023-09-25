package bundler

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
)

type Uploader struct {
	token      string
	configPath string
	config     *config.Config
	dir        string
	client     pb.QuickClient
}

func NewUploader(configPath, dir string, conf *config.Config) *Uploader {
	return &Uploader{
		token:      os.Getenv("SQLC_AUTH_TOKEN"),
		configPath: configPath,
		config:     conf,
		dir:        dir,
	}
}

func (up *Uploader) Validate() error {
	if up.config.Cloud.Project == "" {
		return fmt.Errorf("cloud.project is not set")
	}
	if up.token == "" {
		return fmt.Errorf("SQLC_AUTH_TOKEN environment variable is not set")
	}
	if up.client == nil {
		client, err := quickdb.NewClientFromConfig(up.config.Cloud)
		if err != nil {
			return fmt.Errorf("client init failed: %w", err)
		}
		up.client = client
	}
	return nil
}

func (up *Uploader) buildRequest(ctx context.Context, result map[string]string) (*pb.UploadArchiveRequest, error) {
	ins, err := readInputs(up.configPath, up.config)
	if err != nil {
		return nil, err
	}
	outs, err := readOutputs(up.dir, result)
	if err != nil {
		return nil, err
	}
	return &pb.UploadArchiveRequest{
		SqlcVersion: info.Version,
		Inputs:      ins,
		Outputs:     outs,
	}, nil
}

func (up *Uploader) DumpRequestOut(ctx context.Context, result map[string]string) error {
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(req))
	return nil
}

func (up *Uploader) Upload(ctx context.Context, result map[string]string) error {
	if err := up.Validate(); err != nil {
		return err
	}
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	if _, err := up.client.UploadArchive(ctx, req); err != nil {
		return fmt.Errorf("upload error: %w", err)
	}
	return nil
}
