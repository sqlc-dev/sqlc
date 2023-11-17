package bundler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/plugin"
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

type QuerySetArchive struct {
	Queries []string
	Schema  []string
	Request *plugin.GenerateRequest
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

var envvars = []string{
	"GITHUB_REPOSITORY",
	"GITHUB_REF",
	"GITHUB_REF_NAME",
	"GITHUB_REF_TYPE",
	"GITHUB_SHA",
}

func annotate() map[string]string {
	labels := map[string]string{}
	for _, ev := range envvars {
		key := strings.ReplaceAll(strings.ToLower(ev), "_", ".")
		labels[key] = os.Getenv(ev)
	}
	return labels
}

func (up *Uploader) buildRequest(ctx context.Context, results []*QuerySetArchive) (*plugin.UploadArchiveRequest, error) {
	conf, err := readFile(up.dir, up.configPath)
	if err != nil {
		return nil, err
	}
	res := &plugin.UploadArchiveRequest{
		SqlcVersion: info.Version,
		Config:      conf,
		Annotations: annotate(),
	}
	for _, result := range results {
		schema, err := readFiles(up.dir, result.Schema)
		if err != nil {
			return nil, err
		}
		queries, err := readFiles(up.dir, result.Queries)
		if err != nil {
			return nil, err
		}
		res.Archives = append(res.Archives, &plugin.QuerySetArchive{
			Schema:  schema,
			Queries: queries,
			Request: result.Request,
		})
	}
	return res, nil
}

func (up *Uploader) DumpRequestOut(ctx context.Context, result []*QuerySetArchive) error {
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(req))
	return nil
}

func (up *Uploader) Upload(ctx context.Context, result []*QuerySetArchive) error {
	if err := up.Validate(); err != nil {
		return err
	}
	req, err := up.buildRequest(ctx, result)
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(req))
	// if _, err := up.client.UploadArchive(ctx, req); err != nil {
	// 	return fmt.Errorf("upload error: %w", err)
	// }
	return nil
}
