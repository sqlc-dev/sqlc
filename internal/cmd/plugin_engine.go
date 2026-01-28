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
		return fmt.Errorf("engine plugin not found: %s\n\nSee the engine plugins documentation: https://docs.sqlc.dev/en/latest/guides/engine-plugins.html", r.Cmd)
	}

	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("engine plugin not found: %s\n\nSee the engine plugins documentation: https://docs.sqlc.dev/en/latest/guides/engine-plugins.html", r.Cmd)
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

// runPluginQuerySet runs the plugin-engine path: schema and queries are sent to the
// engine plugin via ParseRequest; the responses are turned into compiler.Result and
// passed to ProcessResult. No AST or compiler parsing is used.
// When inputs.FileContents is set, schema/query bytes are taken from it (no disk read).
func runPluginQuerySet(ctx context.Context, rp resultProcessor, name, dir string, sql outputPair, combo config.CombinedSettings, inputs *sourceFiles, o *Options) error {
	enginePlugin, found := config.FindEnginePlugin(&combo.Global, string(combo.Package.Engine))
	if !found || enginePlugin.Process == nil {
		e := string(combo.Package.Engine)
		return fmt.Errorf("unknown engine: %s\n\nAdd the engine to the 'engines' section of sqlc.yaml. See the engine plugins documentation: https://docs.sqlc.dev/en/latest/guides/engine-plugins.html", e)
	}

	readFile := func(path string) ([]byte, error) {
		if inputs != nil && inputs.FileContents != nil {
			if b, ok := inputs.FileContents[path]; ok {
				return b, nil
			}
		}
		return os.ReadFile(path)
	}

	databaseOnly := combo.Package.Analyzer.Database.IsOnly() && combo.Package.Database != nil && combo.Package.Database.URI != ""

	var schemaSQL string
	if !databaseOnly {
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
	}

	type parseFuncType func(querySQL string) (*pb.ParseResponse, error)
	var parseFn parseFuncType
	if o != nil && o.PluginParseFunc != nil {
		schemaStr := schemaSQL
		if databaseOnly {
			schemaStr = ""
		}
		parseFn = func(querySQL string) (*pb.ParseResponse, error) {
			return o.PluginParseFunc(schemaStr, querySQL)
		}
	} else {
		r := newEngineProcessRunner(enginePlugin.Process.Cmd, combo.Dir, enginePlugin.Env)
		parseFn = func(querySQL string) (*pb.ParseResponse, error) {
			req := &pb.ParseRequest{Sql: querySQL}
			if databaseOnly {
				req.SchemaSource = &pb.ParseRequest_ConnectionParams{
					ConnectionParams: &pb.ConnectionParams{Dsn: combo.Package.Database.URI},
				}
			} else {
				req.SchemaSource = &pb.ParseRequest_SchemaSql{SchemaSql: schemaSQL}
			}
			return r.parseRequest(ctx, req)
		}
	}

	var queryPaths []string
	var err error
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
		resp, err := parseFn(queryContent)
		if err != nil {
			merr.Add(filename, queryContent, 0, err)
			continue
		}
		baseName := filepath.Base(filename)
		stmts := resp.GetStatements()
		for _, st := range stmts {
			q := statementToCompilerQuery(st, baseName)
			if q == nil {
				continue
			}
			qName := st.GetName()
			if _, exists := set[qName]; exists {
				merr.Add(filename, queryContent, 0, fmt.Errorf("duplicate query name: %s", qName))
				continue
			}
			set[qName] = struct{}{}
			queries = append(queries, q)
		}
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

// statementToCompilerQuery converts one engine.Statement from the plugin into a compiler.Query.
func statementToCompilerQuery(st *pb.Statement, filename string) *compiler.Query {
	if st == nil {
		return nil
	}
	sqlTrimmed := strings.TrimSpace(st.GetSql())
	if sqlTrimmed == "" {
		return nil
	}
	var params []compiler.Parameter
	for _, p := range st.GetParameters() {
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
	for _, c := range st.GetColumns() {
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
			Name:     st.GetName(),
			Cmd:      pb.CmdToString(st.GetCmd()),
			Filename: filename,
		},
		Params:  params,
		Columns: columns,
	}
}
