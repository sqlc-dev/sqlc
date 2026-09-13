// Package mssql generates the SQL Server dialect seed under
// internal/engine/mssql/dialect from a live server, and verifies the SQL
// Server analyze cases under internal/endtoend/testdata against the same
// server.
//
// SQL Server keeps no catalog of its built-in functions or operators — the
// intrinsic functions such as GETDATE and LEN are not objects — so
// types.jsonl and functions.jsonl are written by hand and are not this
// package's business. What it does describe is its catalog: relations.jsonl
// is every view of the sys and INFORMATION_SCHEMA schemas, with each view's
// columns as the server itself describes a SELECT * from it.
//
// The package also verifies the SQL Server analyze cases under
// internal/endtoend/testdata against the same server: each case's schema
// and fixture are loaded into a database of their own and its queries
// described there, without being run. Result columns come from
// sys.dm_exec_describe_first_result_set, which says what each column is
// called, what type it has, whether it can be NULL and which table column
// it is read from; parameters from sp_describe_undeclared_parameters, which
// says what type the server would give each, and from the estimated
// showplan, which says which column each is compared with or assigned to.
// The answer is printed in the JSON shape sqlc analyze prints and compared
// with the case's committed stdout.json byte for byte.
//
// The server is named by MSSQL_SERVER_URI, in any form the go-mssqldb
// driver accepts, such as
// sqlserver://sa:password@127.0.0.1:1433?encrypt=disable, and has to be
// the major release in Major, since each release adds to the catalog views.
package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/microsoft/go-mssqldb/msdsn"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name sqlc knows the dialect by, and the name of the analyze
// case directories under internal/endtoend/testdata.
const Engine = "mssql"

// Major is the SQL Server major release the dialect is generated from, as
// SERVERPROPERTY('ProductMajorVersion') reports it: 17 is SQL Server 2025.
// Bumping it is a deliberate change: every release adds views and columns
// to the catalog, so regenerate and review the dialect after changing it.
const Major = 17

// Locate returns the server to generate from, named by MSSQL_SERVER_URI.
func Locate() (string, error) {
	dsn := os.Getenv("MSSQL_SERVER_URI")
	if dsn == "" {
		return "", errors.New("MSSQL_SERVER_URI is not set")
	}
	if _, err := msdsn.Parse(dsn); err != nil {
		return "", fmt.Errorf("MSSQL_SERVER_URI: %w", err)
	}
	return dsn, nil
}

// open connects to the server and holds one session, since a case's
// database is the session's current database and the showplan is a
// session setting.
func open(ctx context.Context, dsn string) (*sql.Conn, error) {
	db, err := sql.Open("sqlserver", dsn)
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
	var v string
	if err := conn.QueryRowContext(ctx, "SELECT @@VERSION").Scan(&v); err != nil {
		return "", err
	}
	// The first line names the product and build; the rest is the copyright
	// and the platform.
	v, _, _ = strings.Cut(v, "\n")
	return strings.TrimSpace(v), nil
}

// checkVersion refuses a server of another major release than the dialect
// is generated from, whose catalog would differ from the committed one
// without anything being wrong.
func checkVersion(ctx context.Context, conn *sql.Conn) error {
	var v string
	if err := conn.QueryRowContext(ctx, "SELECT CAST(SERVERPROPERTY('ProductMajorVersion') AS nvarchar(10))").Scan(&v); err != nil {
		return err
	}
	major, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("ProductMajorVersion %q: %w", v, err)
	}
	if major != Major {
		return fmt.Errorf("the dialect is generated from SQL Server major release %d, but the server is release %d", Major, major)
	}
	return nil
}

// Generate reads the dialect from the server: relations.jsonl, the views of
// sys and INFORMATION_SCHEMA. They are read from a database of their own,
// since the views a query sees are the ones a user database has, and
// master lists internal views of its own that no query can name.
func Generate(ctx context.Context, dsn string) (dialect.Files, error) {
	conn, err := open(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := checkVersion(ctx, conn); err != nil {
		return nil, err
	}
	db := quote("goldeneye_dialect")
	for _, stmt := range []string{
		"USE master",
		"DROP DATABASE IF EXISTS " + db,
		"CREATE DATABASE " + db,
		"USE " + db,
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return nil, err
		}
	}
	defer conn.ExecContext(context.WithoutCancel(ctx), "USE master; DROP DATABASE IF EXISTS "+db)
	relations, err := readRelations(ctx, conn)
	if err != nil {
		return nil, err
	}
	blob, err := dialect.JSONL(relations)
	if err != nil {
		return nil, err
	}
	return dialect.Files{dialect.RelationsFile: blob}, nil
}
