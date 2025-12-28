package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/sqlc-dev/sqlc/internal/engine"
	"github.com/sqlc-dev/sqlc/internal/info"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	pb "github.com/sqlc-dev/sqlc/pkg/engine"
)

// ProcessRunner runs an engine plugin as an external process.
type ProcessRunner struct {
	Cmd string
	Env []string

	// Cached responses
	commentSyntax *pb.GetCommentSyntaxResponse
	dialect       *pb.GetDialectResponse
}

// NewProcessRunner creates a new ProcessRunner.
func NewProcessRunner(cmd string, env []string) *ProcessRunner {
	return &ProcessRunner{
		Cmd: cmd,
		Env: env,
	}
}

func (r *ProcessRunner) invoke(ctx context.Context, method string, req, resp proto.Message) error {
	stdin, err := proto.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	// Parse command string to support formats like "go run ./path"
	cmdParts := strings.Fields(r.Cmd)
	if len(cmdParts) == 0 {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	path, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return fmt.Errorf("engine plugin not found: %s\n\nMake sure the plugin is installed and available in PATH.\nInstall with: go install <plugin-module>@latest", r.Cmd)
	}

	// Build arguments: rest of cmdParts + method
	args := append(cmdParts[1:], method)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Stdin = bytes.NewReader(stdin)
	// Inherit the current environment and add SQLC_VERSION
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

// Parse implements engine.Parser.
func (r *ProcessRunner) Parse(reader io.Reader) ([]ast.Statement, error) {
	sql, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	req := &pb.ParseRequest{Sql: string(sql)}
	resp := &pb.ParseResponse{}

	if err := r.invoke(context.Background(), "parse", req, resp); err != nil {
		return nil, err
	}

	var stmts []ast.Statement
	for _, s := range resp.Statements {
		// Parse the AST JSON into an ast.Node
		node, err := parseASTJSON(s.AstJson)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AST: %w", err)
		}

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         node,
				StmtLocation: int(s.StmtLocation),
				StmtLen:      int(s.StmtLen),
			},
		})
	}

	return stmts, nil
}

// CommentSyntax implements engine.Parser.
func (r *ProcessRunner) CommentSyntax() source.CommentSyntax {
	if r.commentSyntax == nil {
		req := &pb.GetCommentSyntaxRequest{}
		resp := &pb.GetCommentSyntaxResponse{}
		if err := r.invoke(context.Background(), "get_comment_syntax", req, resp); err != nil {
			// Default to common SQL comment syntax
			return source.CommentSyntax{
				Dash:      true,
				SlashStar: true,
			}
		}
		r.commentSyntax = resp
	}

	return source.CommentSyntax{
		Dash:      r.commentSyntax.Dash,
		SlashStar: r.commentSyntax.SlashStar,
		Hash:      r.commentSyntax.Hash,
	}
}

// IsReservedKeyword implements engine.Parser.
func (r *ProcessRunner) IsReservedKeyword(s string) bool {
	req := &pb.IsReservedKeywordRequest{Keyword: s}
	resp := &pb.IsReservedKeywordResponse{}
	if err := r.invoke(context.Background(), "is_reserved_keyword", req, resp); err != nil {
		return false
	}
	return resp.IsReserved
}

// GetCatalog returns the initial catalog for this engine.
func (r *ProcessRunner) GetCatalog() (*catalog.Catalog, error) {
	req := &pb.GetCatalogRequest{}
	resp := &pb.GetCatalogResponse{}
	if err := r.invoke(context.Background(), "get_catalog", req, resp); err != nil {
		return nil, err
	}

	return convertCatalog(resp.Catalog), nil
}

// QuoteIdent implements engine.Dialect.
func (r *ProcessRunner) QuoteIdent(s string) string {
	r.ensureDialect()
	if r.IsReservedKeyword(s) && r.dialect.QuoteChar != "" {
		return r.dialect.QuoteChar + s + r.dialect.QuoteChar
	}
	return s
}

// TypeName implements engine.Dialect.
func (r *ProcessRunner) TypeName(ns, name string) string {
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param implements engine.Dialect.
func (r *ProcessRunner) Param(n int) string {
	r.ensureDialect()
	switch r.dialect.ParamStyle {
	case "dollar":
		return fmt.Sprintf("$%d", n)
	case "question":
		return "?"
	case "at":
		return fmt.Sprintf("@p%d", n)
	default:
		return fmt.Sprintf("$%d", n)
	}
}

// NamedParam implements engine.Dialect.
func (r *ProcessRunner) NamedParam(name string) string {
	r.ensureDialect()
	if r.dialect.ParamPrefix != "" {
		return r.dialect.ParamPrefix + name
	}
	return "@" + name
}

// Cast implements engine.Dialect.
func (r *ProcessRunner) Cast(arg, typeName string) string {
	r.ensureDialect()
	switch r.dialect.CastSyntax {
	case "double_colon":
		return arg + "::" + typeName
	default:
		return "CAST(" + arg + " AS " + typeName + ")"
	}
}

func (r *ProcessRunner) ensureDialect() {
	if r.dialect == nil {
		req := &pb.GetDialectRequest{}
		resp := &pb.GetDialectResponse{}
		if err := r.invoke(context.Background(), "get_dialect", req, resp); err != nil {
			// Use defaults
			r.dialect = &pb.GetDialectResponse{
				QuoteChar:   `"`,
				ParamStyle:  "dollar",
				ParamPrefix: "@",
				CastSyntax:  "cast_function",
			}
		} else {
			r.dialect = resp
		}
	}
}

// convertCatalog converts a protobuf Catalog to catalog.Catalog.
func convertCatalog(c *pb.Catalog) *catalog.Catalog {
	if c == nil {
		return catalog.New("")
	}

	cat := catalog.New(c.DefaultSchema)
	cat.Name = c.Name
	cat.Comment = c.Comment

	// Clear default schemas and add from plugin
	cat.Schemas = make([]*catalog.Schema, 0, len(c.Schemas))
	for _, s := range c.Schemas {
		schema := &catalog.Schema{
			Name:    s.Name,
			Comment: s.Comment,
		}

		for _, t := range s.Tables {
			table := &catalog.Table{
				Rel: &ast.TableName{
					Catalog: t.Catalog,
					Schema:  t.Schema,
					Name:    t.Name,
				},
				Comment: t.Comment,
			}
			for _, col := range t.Columns {
				table.Columns = append(table.Columns, &catalog.Column{
					Name:       col.Name,
					Type:       ast.TypeName{Name: col.DataType},
					IsNotNull:  col.NotNull,
					IsArray:    col.IsArray,
					ArrayDims:  int(col.ArrayDims),
					Comment:    col.Comment,
					Length:     toPointer(int(col.Length)),
					IsUnsigned: col.IsUnsigned,
				})
			}
			schema.Tables = append(schema.Tables, table)
		}

		for _, e := range s.Enums {
			enum := &catalog.Enum{
				Name:    e.Name,
				Comment: e.Comment,
			}
			enum.Vals = append(enum.Vals, e.Values...)
			schema.Types = append(schema.Types, enum)
		}

		for _, f := range s.Functions {
			fn := &catalog.Function{
				Name:       f.Name,
				Comment:    f.Comment,
				ReturnType: &ast.TypeName{Schema: f.ReturnType.GetSchema(), Name: f.ReturnType.GetName()},
			}
			for _, arg := range f.Args {
				fn.Args = append(fn.Args, &catalog.Argument{
					Name:       arg.Name,
					Type:       &ast.TypeName{Schema: arg.Type.GetSchema(), Name: arg.Type.GetName()},
					HasDefault: arg.HasDefault,
				})
			}
			schema.Funcs = append(schema.Funcs, fn)
		}

		for _, t := range s.Types {
			schema.Types = append(schema.Types, &catalog.CompositeType{
				Name:    t.Name,
				Comment: t.Comment,
			})
		}

		cat.Schemas = append(cat.Schemas, schema)
	}

	return cat
}

func toPointer(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}

// parseASTJSON parses AST JSON into an ast.Node.
// This is a placeholder - full implementation would require a JSON-to-AST converter.
func parseASTJSON(data []byte) (ast.Node, error) {
	if len(data) == 0 {
		return &ast.TODO{}, nil
	}

	// Parse the JSON to determine the node type
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	// Check for node_type field
	if nodeType, ok := raw["node_type"]; ok {
		var typeName string
		if err := json.Unmarshal(nodeType, &typeName); err != nil {
			return nil, err
		}
		return parseNodeByType(typeName, data)
	}

	// Default to TODO for unknown structures
	return &ast.TODO{}, nil
}

// parseNodeByType parses a node based on its type.
func parseNodeByType(nodeType string, data []byte) (ast.Node, error) {
	switch strings.ToLower(nodeType) {
	case "select", "selectstmt":
		return parseSelectStmt(data)
	case "insert", "insertstmt":
		return parseInsertStmt(data)
	case "update", "updatestmt":
		return parseUpdateStmt(data)
	case "delete", "deletestmt":
		return parseDeleteStmt(data)
	case "createtable", "createtablestmt":
		return parseCreateTableStmt(data)
	default:
		return &ast.TODO{}, nil
	}
}

// Placeholder implementations for statement parsing
func parseSelectStmt(data []byte) (ast.Node, error) {
	return &ast.SelectStmt{}, nil
}

func parseInsertStmt(data []byte) (ast.Node, error) {
	return &ast.InsertStmt{}, nil
}

func parseUpdateStmt(data []byte) (ast.Node, error) {
	return &ast.UpdateStmt{}, nil
}

func parseDeleteStmt(data []byte) (ast.Node, error) {
	return &ast.DeleteStmt{}, nil
}

func parseCreateTableStmt(data []byte) (ast.Node, error) {
	// Try to extract table name from JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return &ast.CreateTableStmt{}, nil
	}

	stmt := &ast.CreateTableStmt{}

	// Check for table_name in JSON first
	if tableName, ok := raw["table_name"].(string); ok && tableName != "" {
		schema := ""
		name := tableName
		if parts := strings.SplitN(tableName, ".", 2); len(parts) == 2 {
			schema = parts[0]
			name = parts[1]
		}
		stmt.Name = &ast.TableName{Schema: schema, Name: name}
		return stmt, nil
	}

	// Try to extract from raw SQL
	if rawSQL, ok := raw["raw"].(string); ok && rawSQL != "" {
		if name := extractTableNameFromCreateSQL(rawSQL); name != "" {
			stmt.Name = &ast.TableName{Name: name}
		}
	}

	return stmt, nil
}

// extractTableNameFromCreateSQL extracts table name from CREATE TABLE statement
func extractTableNameFromCreateSQL(sql string) string {
	sql = strings.TrimSpace(sql)
	upper := strings.ToUpper(sql)

	// Handle CREATE TABLE [IF NOT EXISTS] name
	idx := strings.Index(upper, "CREATE TABLE")
	if idx == -1 {
		return ""
	}
	sql = strings.TrimSpace(sql[idx+len("CREATE TABLE"):])
	upper = strings.ToUpper(sql)

	// Skip IF NOT EXISTS
	if strings.HasPrefix(upper, "IF NOT EXISTS") {
		sql = strings.TrimSpace(sql[len("IF NOT EXISTS"):])
	}

	// Extract table name (until space or parenthesis)
	var name strings.Builder
	for _, r := range sql {
		if r == ' ' || r == '(' || r == '\t' || r == '\n' || r == '\r' {
			break
		}
		name.WriteRune(r)
	}

	result := name.String()
	// Remove quotes if present
	result = strings.Trim(result, `"'`+"`")
	return result
}

// PluginEngine wraps a ProcessRunner to implement engine.Engine.
type PluginEngine struct {
	name   string
	runner *ProcessRunner
}

// NewPluginEngine creates a new engine from a process plugin.
func NewPluginEngine(name, cmd string, env []string) *PluginEngine {
	return &PluginEngine{
		name:   name,
		runner: NewProcessRunner(cmd, env),
	}
}

// Name implements engine.Engine.
func (e *PluginEngine) Name() string {
	return e.name
}

// Parser implements engine.Engine.
func (e *PluginEngine) Parser() engine.Parser {
	return e.runner
}

// Catalog implements engine.Engine.
func (e *PluginEngine) Catalog() *catalog.Catalog {
	cat, err := e.runner.GetCatalog()
	if err != nil {
		// Return empty catalog on error
		return catalog.New("")
	}
	return cat
}

// Selector implements engine.Engine.
func (e *PluginEngine) Selector() engine.Selector {
	return &engine.DefaultSelector{}
}

// Dialect implements engine.Engine.
func (e *PluginEngine) Dialect() engine.Dialect {
	return e.runner
}
