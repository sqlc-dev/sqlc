package bundler

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/quickdb"
	pb "github.com/sqlc-dev/sqlc/internal/quickdb/v1"
)

var ErrNoProject = errors.New(`project uploads require a cloud project

If you don't have a project, you can create one from the sqlc Cloud
dashboard at https://dashboard.sqlc.dev/. If you have a project, ensure
you've set its id as the value of the "project" field within the "cloud"
section of your sqlc configuration. The id will look similar to
"01HA8TWGMYPHK0V2GGMB3R2TP9".`)
var ErrNoAuthToken = errors.New(`project uploads require an auth token

If you don't have an auth token, you can create one from the sqlc Cloud
dashboard at https://dashboard.sqlc.dev/. If you have an auth token, ensure
you've set it as the value of the SQLC_AUTH_TOKEN environment variable.`)

type Uploader struct {
	configPath string
	config     *config.Config
	dir        string
	client     pb.QuickClient
}

func NewUploader(configPath, dir string, conf *config.Config) *Uploader {
	return &Uploader{
		configPath: configPath,
		config:     conf,
		dir:        dir,
	}
}

func (up *Uploader) Validate() error {
	if up.config.Cloud.Project == "" {
		return ErrNoProject
	}
	if up.config.Cloud.AuthToken == "" {
		return ErrNoAuthToken
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
