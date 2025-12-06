package clickhouse

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	chparser "github.com/AfterShip/clickhouse-sql-parser/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
	// catalog is set by the Compiler after schema parsing is complete
	// It allows the parser to register context-dependent functions
	Catalog interface{}
}

// preprocessNamedParameters converts sqlc named parameter syntax to valid ClickHouse syntax
// that the parser can handle. This allows sqlc.arg(), sqlc.narg(), and sqlc.slice() syntax
// to be used in queries by converting the dot notation to underscore notation, which the
// ClickHouse parser can recognize as a function call.
//
// Conversions:
// - sqlc.arg('name') → sqlc_arg('name')
// - sqlc.narg('name') → sqlc_narg('name')
// - sqlc.slice('name') → sqlc_slice('name')
//
// The original SQL is preserved in the compiler, so rewrite.NamedParameters can still
// find and process the original named parameter syntax. The converter normalizes
// sqlc_* function names back to sqlc.* schema.function format in the AST.
func preprocessNamedParameters(sql string) string {
	// Convert sqlc.arg/narg/slice to sqlc_arg/narg/slice
	// This makes them valid function names in ClickHouse parser
	// Using same-length replacement (sqlc. = 5 chars, sqlc_ = 5 chars) preserves positions
	funcPattern := regexp.MustCompile(`sqlc\.(arg|narg|slice)`)
	sql = funcPattern.ReplaceAllString(sql, "sqlc_$1")

	return sql
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	originalSQL := string(blob)

	// Preprocess to replace named parameter syntax with valid ClickHouse syntax
	processedSQL := preprocessNamedParameters(originalSQL)

	chp := chparser.NewParser(processedSQL)
	stmtNodes, err := chp.ParseStmts()
	if err != nil {
		return nil, normalizeErr(err)
	}

	// Find all -- name: comments in the original SQL
	nameCommentPositions := findNameComments(originalSQL)

	// Clone the catalog for this parse operation to isolate function registrations
	// (like arrayJoin, argMin, etc.) from other queries
	var clonedCatalog *catalog.Catalog
	if p.Catalog != nil {
		cat, ok := p.Catalog.(*catalog.Catalog)
		if !ok {
			return nil, fmt.Errorf("invalid catalog type: expected *catalog.Catalog, got %T", p.Catalog)
		}
		clonedCatalog = cat.Clone()
	}

	var stmts []ast.Statement
	for i := range stmtNodes {
		converter := &cc{catalog: clonedCatalog}
		out := converter.convert(stmtNodes[i])

		var statementStart, statementEnd int

		// Check if we're processing a file with -- name: comments (queries.sql)
		// or without them (schema.sql)
		if len(nameCommentPositions) == len(stmtNodes) {
			// We have a -- name: comment for each statement (queries file)
			statementStart = nameCommentPositions[i]

			if i+1 < len(nameCommentPositions) {
				statementEnd = nameCommentPositions[i+1]
			} else {
				statementEnd = len(originalSQL)
			}
		} else {
			// No name comments, or mismatch (schema file or mixed)
			// Use the parser's positions, but try to find better boundaries
			processedStmtPos := int(stmtNodes[i].Pos())
			if processedStmtPos > 0 {
				processedStmtPos -= 1
			}

			originalStmtPos := findOriginalPosition(originalSQL, processedSQL, processedStmtPos)
			statementStart = findStatementStart(originalSQL, originalStmtPos)

			if i+1 < len(stmtNodes) {
				nextProcessedStmtPos := int(stmtNodes[i+1].Pos())
				if nextProcessedStmtPos > 0 {
					nextProcessedStmtPos -= 1
				}
				nextOriginalStmtPos := findOriginalPosition(originalSQL, processedSQL, nextProcessedStmtPos)
				statementEnd = findStatementStart(originalSQL, nextOriginalStmtPos)
			} else {
				statementEnd = len(originalSQL)
			}
		}

		// Bounds check
		if statementStart < 0 || statementStart >= len(originalSQL) || statementEnd > len(originalSQL) || statementStart >= statementEnd {
			continue
		}

		segment := originalSQL[statementStart:statementEnd]
		// Trim trailing whitespace but preserve the content
		segment = strings.TrimRight(segment, " \t\r\n")

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: statementStart,
				StmtLen:      len(segment),
			},
		})
	}

	return stmts, nil
}

// findNameComments finds all positions of -- name: comments in the SQL
// Returns a slice of positions where each -- name: comment starts
func findNameComments(sql string) []int {
	var positions []int
	lines := strings.Split(sql, "\n")
	currentPos := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Check if this line contains a -- name: comment
		if strings.HasPrefix(trimmed, "--") && strings.Contains(trimmed, "name:") {
			// Find the actual position of the start of this line (not trimmed)
			lineStart := currentPos
			// Walk backwards to find any leading whitespace
			for lineStart < currentPos+len(line) && (sql[lineStart] == ' ' || sql[lineStart] == '\t') {
				lineStart++
			}
			// Actually, we want the start of the line including whitespace
			positions = append(positions, currentPos)
		}
		// Move to next line (including the \n character)
		currentPos += len(line) + 1
	}

	return positions
}

// findOriginalPosition maps a position in processedSQL back to the original SQL
// Preprocessing only changes sqlc.arg to sqlc_arg (same length), so positions are mostly 1:1
func findOriginalPosition(originalSQL, processedSQL string, processedPos int) int {
	// Since sqlc. → sqlc_ is same length (5 chars), positions are virtually identical
	// Just ensure we don't go out of bounds
	if processedPos >= len(originalSQL) {
		return len(originalSQL)
	}
	return processedPos
}

// findStatementStart finds the start of a statement in SQL, including preceding -- name: annotation
func findStatementStart(sql string, stmtPos int) int {
	if stmtPos <= 0 {
		return 0
	}

	// Walk backwards through lines to find the -- name: annotation
	// The stmtPos usually points to the SELECT/INSERT/etc line, but we need to include
	// the -- name: comment that precedes it.

	currentPos := stmtPos

	// Keep walking backwards through lines
	for currentPos > 0 {
		// Find the start of the current line
		lineStart := currentPos - 1
		for lineStart >= 0 && sql[lineStart] != '\n' {
			lineStart--
		}
		lineStart++ // Move past the newline (or stay at 0)

		// Extract the line content, skipping leading whitespace
		checkPos := lineStart
		for checkPos < len(sql) && (sql[checkPos] == ' ' || sql[checkPos] == '\t') {
			checkPos++
		}

		// Find the end of this line
		lineEnd := checkPos
		for lineEnd < len(sql) && sql[lineEnd] != '\n' {
			lineEnd++
		}

		if checkPos < lineEnd {
			lineText := sql[checkPos:lineEnd]

			// Check if this line is a -- name: annotation
			if strings.HasPrefix(lineText, "--") && strings.Contains(lineText, "name:") {
				// Found it! Return the start of this line
				return lineStart
			}

			// Check if this is a non-empty, non-comment line
			// If we find actual SQL before finding the -- name: comment, stop
			if len(strings.TrimSpace(lineText)) > 0 && !strings.HasPrefix(lineText, "--") {
				// This is SQL content, stop searching backwards
				break
			}
		}

		// Move to the previous line
		if lineStart == 0 {
			break
		}
		currentPos = lineStart - 1
	}

	// Didn't find a -- name: annotation, return the original position
	return stmtPos
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// normalizeErr converts ClickHouse parser errors to sqlc error format
func normalizeErr(err error) error {
	if err == nil {
		return err
	}

	// For now, wrap the error as a generic syntax error
	// ClickHouse parser may provide better error information in future versions
	return &sqlerr.Error{
		Message: "syntax error",
		Err:     errors.New(err.Error()),
		Line:    1,
		Column:  1,
	}
}

// CommentSyntax returns the comment syntax for ClickHouse
// ClickHouse supports:
// - Line comments: -- (with optional space after)
// - Block comments: /* */
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		SlashStar: true,
	}
}

// IsReservedKeyword checks if a string is a ClickHouse reserved keyword
func (p *Parser) IsReservedKeyword(s string) bool {
	return isReservedKeyword(strings.ToUpper(s))
}

// ClickHouse reserved keywords
// Based on https://clickhouse.com/docs/sql-reference/syntax#keywords
var reservedKeywords = map[string]bool{
	"ALTER":       true,
	"AND":         true,
	"ARRAY":       true,
	"AS":          true,
	"ASCENDING":   true,
	"ASOF":        true,
	"BETWEEN":     true,
	"BY":          true,
	"CASE":        true,
	"CAST":        true,
	"CHECK":       true,
	"CLUSTER":     true,
	"CODEC":       true,
	"COLLATE":     true,
	"COLUMN":      true,
	"CONSTRAINT":  true,
	"CREATE":      true,
	"CROSS":       true,
	"CUBE":        true,
	"CURRENT":     true,
	"DATABASE":    true,
	"DAY":         true,
	"DEDUPLICATE": true,
	"DEFAULT":     true,
	"DEFINER":     true,
	"DELETE":      true,
	"DESC":        true,
	"DESCENDING":  true,
	"DESCRIBE":    true,
	"DISTINCT":    true,
	"DROP":        true,
	"ELSE":        true,
	"END":         true,
	"ESCAPING":    true,
	"EXCEPT":      true,
	"EXCHANGE":    true,
	"EXPLAIN":     true,
	"FETCH":       true,
	"FILL":        true,
	"FINAL":       true,
	"FIRST":       true,
	"FOR":         true,
	"FOREGROUND":  true,
	"FROM":        true,
	"FULL":        true,
	"FUNCTION":    true,
	"GLOBAL":      true,
	"GRANT":       true,
	"GROUP":       true,
	"HAVING":      true,
	"HOUR":        true,
	"IF":          true,
	"ILIKE":       true,
	"IN":          true,
	"INNER":       true,
	"INTERSECT":   true,
	"INTO":        true,
	"IS":          true,
	"ISNULL":      true,
	"JOIN":        true,
	"KEY":         true,
	"KILL":        true,
	"LAST":        true,
	"LATERAL":     true,
	"LEFT":        true,
	"LIKE":        true,
	"LIMIT":       true,
	"LOCAL":       true,
	"NATURAL":     true,
	"NOT":         true,
	"NOTNULL":     true,
	"NULL":        true,
	"OFFSET":      true,
	"ON":          true,
	"OR":          true,
	"ORDER":       true,
	"OUTER":       true,
	"PARTITION":   true,
	"PREWHERE":    true,
	"PRIMARY":     true,
	"REVOKE":      true,
	"RIGHT":       true,
	"ROLLUP":      true,
	"ROW":         true,
	"ROWS":        true,
	"SAMPLE":      true,
	"SELECT":      true,
	"SEMI":        true,
	"SET":         true,
	"SETTINGS":    true,
	"SHOW":        true,
	"SOME":        true,
	"SUBJECT":     true,
	"TABLE":       true,
	"THEN":        true,
	"TIES":        true,
	"TRUNCATE":    true,
	"UNION":       true,
	"UPDATE":      true,
	"USING":       true,
	"VIEW":        true,
	"WHEN":        true,
	"WHERE":       true,
	"WINDOW":      true,
	"WITH":        true,
	"YEAR":        true,
}

func isReservedKeyword(s string) bool {
	return reservedKeywords[s]
}
