// This file runs a database-engine plugin as an external process (parse RPC over stdin/stdout).
// It is used only by the plugin-engine generate path (runPluginQuerySet). Vet does not support plugin engines.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/compiler"
	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/multierr"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlpath"
	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/info"
	pb "github.com/sqlc-dev/sqlc/pkg/engine"
)

// engineProcessRunner runs an engine plugin as an external process.
type engineProcessRunner struct {
	Cmd string
	Dir string // Working directory for the plugin (config file directory)
	Env []string
}

func newEngineProcessRunner(cmd, dir string, env []string) *engineProcessRunner {
	return &engineProcessRunner{Cmd: cmd, Dir: dir, Env: env}
}

func (r *engineProcessRunner) invoke(ctx context.Context, method string, req, resp proto.Message) error {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	cmdParts := strings.Fields(r.Cmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	args := append(cmdParts[1:], method)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = bytes.NewReader(stdin)
	if r.Dir != "" {
		cmd.Dir = r.Dir
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("SQLC_VERSION=%s", info.Version))

	out, err := cmd.Output()
	if err != nil {
		stderr := err.Error()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			stderr = string(exit.Stderr)
		}
		return fmt.Errorf("engine plugin error: %s", stderr)
	}

	if err := proto.Unmarshal(out, resp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	return nil
}

// parseRequest invokes the plugin's Parse RPC. Used by runPluginQuerySet.
func (r *engineProcessRunner) parseRequest(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	resp := &pb.ParseResponse{}
	if err := r.invoke(ctx, "parse", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// defaultCommentSyntax is used when parsing query names from plugin-engine query files.
var defaultCommentSyntax = metadata.CommentSyntax(source.CommentSyntax{Dash: true, SlashStar: true, Hash: false})

// runPluginQuerySet runs the plugin-engine path: schema and queries are sent to the
// engine plugin via ParseRequest; the responses are turned into compiler.Result and
// passed to ProcessResult. No AST or compiler parsing is used.
// When inputs.FileContents is set, schema/query bytes are taken from it (no disk read).
func runPluginQuerySet(ctx context.Context, rp ResultProcessor, name, dir string, sql OutputPair, combo config.CombinedSettings, inputs *sourceFiles, o *Options) error {
	enginePlugin, found := config.FindEnginePlugin(&combo.Global, string(combo.Package.Engine))
	if !found || enginePlugin.Process == nil {
		e := string(combo.Package.Engine)
		return fmt.Errorf("unknown engine: %s\n\nTo use a custom database engine, add it to the 'engines' section of sqlc.yaml:\n\n  engines:\n    - name: %s\n      process:\n        cmd: sqlc-engine-%s\n\nThen install the plugin: go install github.com/example/sqlc-engine-%s@latest",
			e, e, e, e)
	}

	var parseFn func(schemaSQL, querySQL string) (*pb.ParseResponse, error)
	if o != nil && o.PluginParseFunc != nil {
		parseFn = o.PluginParseFunc
	} else {
		r := newEngineProcessRunner(enginePlugin.Process.Cmd, combo.Dir, enginePlugin.Env)
		parseFn = func(schemaSQL, querySQL string) (*pb.ParseResponse, error) {
			req := &pb.ParseRequest{Sql: querySQL}
			if schemaSQL != "" {
				req.SchemaSource = &pb.ParseRequest_SchemaSql{SchemaSql: schemaSQL}
			}
			return r.parseRequest(ctx, req)
		}
	}

	readFile := func(path string) ([]byte, error) {
		if inputs != nil && inputs.FileContents != nil {
			if b, ok := inputs.FileContents[path]; ok {
				return b, nil
			}
		}
		return os.ReadFile(path)
	}

	var schemaSQL string
	var err error
	if inputs != nil && inputs.FileContents != nil {
		var parts []string
		for _, p := range sql.Schema {
			if b, ok := inputs.FileContents[p]; ok {
				parts = append(parts, string(b))
			}
		}
		schemaSQL = strings.Join(parts, "\n")
	} else {
		schemaSQL, err = loadSchemaSQL(sql.Schema, readFile)
		if err != nil {
			return err
		}
	}

	var queryPaths []string
	if inputs != nil && inputs.FileContents != nil {
		queryPaths = sql.Queries
	} else {
		queryPaths, err = sqlpath.Glob(sql.Queries)
		if err != nil {
			return err
		}
	}

	var queries []*compiler.Query
	merr := multierr.New()
	set := map[string]struct{}{}

	for _, filename := range queryPaths {
		blob, err := readFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		queryContent := string(blob)
		nameStr, cmd, err := metadata.ParseQueryNameAndType(queryContent, defaultCommentSyntax)
		if err != nil {
			merr.Add(filename, queryContent, 0, err)
			continue
		}
		if nameStr == "" {
			continue
		}

		resp, err := parseFn(schemaSQL, queryContent)
		if err != nil {
			merr.Add(filename, queryContent, 0, err)
			continue
		}

		q := pluginResponseToCompilerQuery(nameStr, cmd, filepath.Base(filename), resp)
		if q == nil {
			continue
		}

		if _, exists := set[nameStr]; exists {
			merr.Add(filename, queryContent, 0, fmt.Errorf("duplicate query name: %s", nameStr))
			continue
		}
		set[nameStr] = struct{}{}
		queries = append(queries, q)
	}

	if len(merr.Errs()) > 0 {
		return merr
	}
	if len(queries) == 0 {
		return fmt.Errorf("no queries in paths %s", strings.Join(sql.Queries, ","))
	}

	result := &compiler.Result{
		Catalog: catalog.New(""),
		Queries: queries,
	}
	return rp.ProcessResult(ctx, combo, sql, result)
}

func loadSchemaSQL(schemaPaths []string, readFile func(string) ([]byte, error)) (string, error) {
	var parts []string
	for _, p := range schemaPaths {
		files, err := sqlpath.Glob([]string{p})
		if err != nil {
			return "", err
		}
		if len(files) == 0 {
			files = []string{p}
		}
		for _, f := range files {
			b, err := readFile(f)
			if err != nil {
				return "", err
			}
			parts = append(parts, string(b))
		}
	}
	return strings.Join(parts, "\n"), nil
}

func pluginResponseToCompilerQuery(name, cmd, filename string, resp *pb.ParseResponse) *compiler.Query {
	sqlTrimmed := strings.TrimSpace(resp.GetSql())
	if sqlTrimmed == "" {
		return nil
	}

	var params []compiler.Parameter
	for _, p := range resp.GetParameters() {
		col := &compiler.Column{
			DataType:  p.GetDataType(),
			NotNull:   !p.GetNullable(),
			IsArray:   p.GetIsArray(),
			ArrayDims: int(p.GetArrayDims()),
		}
		pos := int(p.GetPosition())
		if pos < 1 {
			pos = len(params) + 1
		}
		params = append(params, compiler.Parameter{Number: pos, Column: col})
	}

	var columns []*compiler.Column
	for _, c := range resp.GetColumns() {
		columns = append(columns, &compiler.Column{
			Name:      c.GetName(),
			DataType:  c.GetDataType(),
			NotNull:   !c.GetNullable(),
			IsArray:   c.GetIsArray(),
			ArrayDims: int(c.GetArrayDims()),
		})
	}

	return &compiler.Query{
		SQL: sqlTrimmed,
		Metadata: metadata.Metadata{
			Name:     name,
			Cmd:      cmd,
			Filename: filename,
		},
		Params:  params,
		Columns: columns,
	}
}
