// sqlc-engine-sqlite3 demonstrates a custom database engine plugin.
//
// This plugin provides SQLite3 SQL parsing for sqlc. It shows how external
// repositories can implement database support without modifying sqlc core.
//
// Build: go build -o sqlc-engine-sqlite3 .
package main

import (
	"encoding/json"
	"strings"

	"github.com/sqlc-dev/sqlc/pkg/engine"
)

func main() {
	engine.Run(engine.Handler{
		PluginName:        "sqlite3",
		PluginVersion:     "1.0.0",
		Parse:             handleParse,
		GetCatalog:        handleGetCatalog,
		IsReservedKeyword: handleIsReservedKeyword,
		GetCommentSyntax:  handleGetCommentSyntax,
		GetDialect:        handleGetDialect,
	})
}

func handleParse(req *engine.ParseRequest) (*engine.ParseResponse, error) {
	sql := req.GetSql()
	statements := splitStatements(sql)
	var result []*engine.Statement

	for _, stmt := range statements {
		ast := map[string]interface{}{
			"node_type": detectStatementType(stmt),
			"raw":       stmt,
		}
		astJSON, _ := json.Marshal(ast)

		result = append(result, &engine.Statement{
			RawSql:       stmt,
			StmtLocation: int32(strings.Index(sql, stmt)),
			StmtLen:      int32(len(stmt)),
			AstJson:      astJSON,
		})
	}

	return &engine.ParseResponse{Statements: result}, nil
}

func handleGetCatalog(req *engine.GetCatalogRequest) (*engine.GetCatalogResponse, error) {
	return &engine.GetCatalogResponse{
		Catalog: &engine.Catalog{
			DefaultSchema: "main",
			Name:          "sqlite3",
			Schemas: []*engine.Schema{
				{
					Name:   "main",
					Tables: []*engine.Table{},
				},
			},
		},
	}, nil
}

func handleIsReservedKeyword(req *engine.IsReservedKeywordRequest) (*engine.IsReservedKeywordResponse, error) {
	reserved := map[string]bool{
		"abort": true, "action": true, "add": true, "after": true,
		"all": true, "alter": true, "analyze": true, "and": true,
		"as": true, "asc": true, "attach": true, "autoincrement": true,
		"before": true, "begin": true, "between": true, "by": true,
		"cascade": true, "case": true, "cast": true, "check": true,
		"collate": true, "column": true, "commit": true, "conflict": true,
		"constraint": true, "create": true, "cross": true, "current_date": true,
		"current_time": true, "current_timestamp": true, "database": true,
		"default": true, "deferrable": true, "deferred": true, "delete": true,
		"desc": true, "detach": true, "distinct": true, "drop": true,
		"each": true, "else": true, "end": true, "escape": true,
		"except": true, "exclusive": true, "exists": true, "explain": true,
		"fail": true, "for": true, "foreign": true, "from": true,
		"full": true, "glob": true, "group": true, "having": true,
		"if": true, "ignore": true, "immediate": true, "in": true,
		"index": true, "indexed": true, "initially": true, "inner": true,
		"insert": true, "instead": true, "intersect": true, "into": true,
		"is": true, "isnull": true, "join": true, "key": true,
		"left": true, "like": true, "limit": true, "match": true,
		"natural": true, "no": true, "not": true, "notnull": true,
		"null": true, "of": true, "offset": true, "on": true,
		"or": true, "order": true, "outer": true, "plan": true,
		"pragma": true, "primary": true, "query": true, "raise": true,
		"recursive": true, "references": true, "regexp": true, "reindex": true,
		"release": true, "rename": true, "replace": true, "restrict": true,
		"right": true, "rollback": true, "row": true, "savepoint": true,
		"select": true, "set": true, "table": true, "temp": true,
		"temporary": true, "then": true, "to": true, "transaction": true,
		"trigger": true, "union": true, "unique": true, "update": true,
		"using": true, "vacuum": true, "values": true, "view": true,
		"virtual": true, "when": true, "where": true, "with": true,
		"without": true,
	}
	return &engine.IsReservedKeywordResponse{
		IsReserved: reserved[strings.ToLower(req.GetKeyword())],
	}, nil
}

func handleGetCommentSyntax(req *engine.GetCommentSyntaxRequest) (*engine.GetCommentSyntaxResponse, error) {
	return &engine.GetCommentSyntaxResponse{
		Dash:      true,
		SlashStar: true,
		Hash:      false,
	}, nil
}

func handleGetDialect(req *engine.GetDialectRequest) (*engine.GetDialectResponse, error) {
	return &engine.GetDialectResponse{
		QuoteChar:   `"`,
		ParamStyle:  "question",
		ParamPrefix: "?",
		CastSyntax:  "cast_function",
	}, nil
}

func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	for _, line := range strings.Split(sql, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		// Include sqlc metadata comments (-- name: ...) with the statement
		if strings.HasPrefix(trimmedLine, "--") {
			// Check if it's a sqlc query annotation
			if strings.Contains(trimmedLine, "name:") {
				current.WriteString(trimmedLine)
				current.WriteString("\n")
			}
			// Skip other comments
			continue
		}
		current.WriteString(trimmedLine)
		current.WriteString(" ")
		if strings.HasSuffix(trimmedLine, ";") {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" && stmt != ";" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}
	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	return statements
}

func detectStatementType(sql string) string {
	sql = strings.ToUpper(strings.TrimSpace(sql))
	switch {
	case strings.HasPrefix(sql, "SELECT"):
		return "SelectStmt"
	case strings.HasPrefix(sql, "INSERT"):
		return "InsertStmt"
	case strings.HasPrefix(sql, "UPDATE"):
		return "UpdateStmt"
	case strings.HasPrefix(sql, "DELETE"):
		return "DeleteStmt"
	case strings.HasPrefix(sql, "CREATE TABLE"):
		return "CreateTableStmt"
	default:
		return "Unknown"
	}
}
