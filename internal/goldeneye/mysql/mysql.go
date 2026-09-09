// Package mysql generates the MySQL dialect seed under
// internal/engine/dolphin/dialect from a live server, and verifies the
// MySQL analyze cases under internal/endtoend/testdata against the same
// server.
//
// MySQL keeps no catalog of its types, functions or operators — the server
// describes a built-in function no further than its name, in the help
// tables — so types.jsonl, functions.jsonl and operators.jsonl are written
// by hand and are not this package's business. What it does describe is
// its data dictionary: relations.jsonl is every view of information_schema,
// read from information_schema itself.
//
// The package also verifies the MySQL analyze cases under
// internal/endtoend/testdata against the same server: each case's schema
// and fixture are loaded into a database of their own and its queries run
// there. Result columns come from the result set's metadata as the driver
// reports it, provenance and parameters from the optimizer trace's
// expanded_query and the note EXPLAIN leaves, since MySQL reports nothing
// about a parameter but its position. The answer is printed in the JSON
// shape sqlc analyze prints and compared with the case's committed
// stdout.json byte for byte.
//
// The server is named by MYSQL_SERVER_URI, in the form the go-sql-driver
// DSN takes, and has to be the major release in Major, since each release
// adds to information_schema.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	driver "github.com/go-sql-driver/mysql"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name sqlc knows the dialect by, and the name of the analyze
// case directories under internal/endtoend/testdata.
const Engine = "mysql"

// Dir is the name of the engine directory the dialect lives under, which
// is named after the parser sqlc reads MySQL with rather than the engine.
const Dir = "dolphin"

// Major is the MySQL major release the dialect is generated from. Bumping
// it is a deliberate change: every release adds tables and columns to the
// system schemas, so regenerate and review the dialect after changing it.
const Major = 26

// Locate returns the server to generate from, named by MYSQL_SERVER_URI.
func Locate() (string, error) {
	dsn := os.Getenv("MYSQL_SERVER_URI")
	if dsn == "" {
		return "", errors.New("MYSQL_SERVER_URI is not set")
	}
	if _, err := driver.ParseDSN(dsn); err != nil {
		return "", fmt.Errorf("MYSQL_SERVER_URI: %w", err)
	}
	return dsn, nil
}

// open connects to the server. The connection runs several statements at
// once, so that a case's schema loads as one script, and holds one
// session, since the optimizer trace and the warnings a statement leaves
// are both per session.
func open(ctx context.Context, dsn string) (*sql.Conn, error) {
	cfg, err := driver.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MultiStatements = true
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	conn, err := db.Conn(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return conn, nil
}

// Version reports the release a server is.
func Version(ctx context.Context, dsn string) (string, error) {
	conn, err := open(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return version(ctx, conn)
}

func version(ctx context.Context, conn *sql.Conn) (string, error) {
	var v string
	if err := conn.QueryRowContext(ctx, "SELECT version()").Scan(&v); err != nil {
		return "", err
	}
	return "MySQL " + v, nil
}

// checkVersion refuses a server of another major release than the dialect
// is generated from, whose system schemas would differ from the committed
// ones without anything being wrong.
func checkVersion(ctx context.Context, conn *sql.Conn) error {
	var v string
	if err := conn.QueryRowContext(ctx, "SELECT version()").Scan(&v); err != nil {
		return err
	}
	head, _, _ := strings.Cut(v, ".")
	major, err := strconv.Atoi(head)
	if err != nil {
		return fmt.Errorf("version %q: %w", v, err)
	}
	if major != Major {
		return fmt.Errorf("the dialect is generated from MySQL %d, but the server is MySQL %d", Major, major)
	}
	return nil
}

// Generate reads the dialect from the server: relations.jsonl, the views of
// information_schema.
func Generate(ctx context.Context, dsn string) (dialect.Files, error) {
	conn, err := open(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := checkVersion(ctx, conn); err != nil {
		return nil, err
	}
	var relations []dialect.Relation
	for _, schema := range systemSchemas {
		rels, err := readRelations(ctx, conn, schema)
		if err != nil {
			return nil, err
		}
		relations = append(relations, rels...)
	}
	blob, err := dialect.JSONL(relations)
	if err != nil {
		return nil, err
	}
	return dialect.Files{dialect.RelationsFile: blob}, nil
}
