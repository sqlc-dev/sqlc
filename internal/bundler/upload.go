package bundler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"

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
	Name    string
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

func BuildRequest(ctx context.Context, dir, configPath string, results []*QuerySetArchive, tags []string) (*pb.UploadArchiveRequest, error) {
	conf, err := readFile(dir, configPath)
	if err != nil {
		return nil, err
	}
	res := &pb.UploadArchiveRequest{
		SqlcVersion: info.Version,
		Config:      conf,
		Tags:        tags,
		Annotations: annotate(),
	}
	for i, result := range results {
		schema, err := readFiles(dir, result.Schema)
		if err != nil {
			return nil, err
		}
		queries, err := readFiles(dir, result.Queries)
		if err != nil {
			return nil, err
		}
		name := result.Name
		if name == "" {
			name = fmt.Sprintf("queryset_%d", i)
		}
		genreq, err := proto.Marshal(result.Request)
		if err != nil {
			return nil, err
		}
		res.QuerySets = append(res.QuerySets, &pb.QuerySet{
			Name:    name,
			Schema:  schema,
			Queries: queries,
			CodegenRequest: &pb.File{
				Name:     "codegen_request.pb",
				Contents: genreq,
			},
		})
	}
	return res, nil
}

func (up *Uploader) buildRequest(ctx context.Context, results []*QuerySetArchive, tags []string) (*pb.UploadArchiveRequest, error) {
	return BuildRequest(ctx, up.dir, up.configPath, results, tags)
}

func (up *Uploader) DumpRequestOut(ctx context.Context, result []*QuerySetArchive) error {
	req, err := up.buildRequest(ctx, result, []string{})
	if err != nil {
		return err
	}
	slog.Info("config", "file", req.Config.Name, "bytes", len(req.Config.Contents))
	for _, qs := range req.QuerySets {
		slog.Info("codegen_request", "queryset", qs.Name, "file", "codegen_request.pb")
		for _, file := range qs.Schema {
			slog.Info("schema", "queryset", qs.Name, "file", file.Name, "bytes", len(file.Contents))
		}
		for _, file := range qs.Queries {
			slog.Info("query", "queryset", qs.Name, "file", file.Name, "bytes", len(file.Contents))
		}
	}
	return nil
}

func (up *Uploader) Upload(ctx context.Context, result []*QuerySetArchive, tags []string) error {
	if err := up.Validate(); err != nil {
		return err
	}
	req, err := up.buildRequest(ctx, result, tags)
	if err != nil {
		return err
	}
	if _, err := up.client.UploadArchive(ctx, req); err != nil {
		return fmt.Errorf("upload error: %w", err)
	}
	return nil
}
