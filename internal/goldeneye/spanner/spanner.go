// Package spanner generates the GoogleSQL dialect seed under
// internal/engine/googlesql/dialect from a live Spanner server, and
// verifies the GoogleSQL analyze cases under internal/endtoend/testdata
// against the same server. The server is Spanner Omni, the downloadable
// Spanner, run from its container image; the package writes into the
// googlesql directory, since sqlc's engine is named after the language
// Spanner speaks rather than the database.
//
// Spanner keeps no catalog of its types, functions or operators, so
// types.jsonl, functions.jsonl and operators.jsonl are written by hand and
// are not this package's business. What it does describe is its
// information schema: relations.jsonl is every view of INFORMATION_SCHEMA
// and SPANNER_SYS, read from INFORMATION_SCHEMA itself.
//
// The package also verifies the GoogleSQL analyze cases against the same
// server: each case's schema becomes a database of its own, its fixture is
// written there, and each query is compiled in PLAN mode, which runs
// nothing and reports the type of each result column and of each
// parameter the query leaves undeclared, and the query plan, which says
// which table column each result column and parameter stands for. The
// answer is printed in the JSON shape sqlc analyze prints and compared
// with the case's committed stdout.json byte for byte.
//
// The server is named by SPANNER_SERVER_URI, the gRPC endpoint of a
// Spanner Omni server such as localhost:15000, which is reached without
// TLS or credentials, as Omni is. Databases are created in the instance
// Omni's single server provides, projects/default/instances/default.
package spanner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	spannerapi "cloud.google.com/go/spanner/apiv1"
	"cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
)

// Engine is the name goldeneye knows the database by.
const Engine = "spanner"

// Dir is the name of the engine directory the dialect lives under, and
// Cases the name of the analyze case directories under
// internal/endtoend/testdata: both are named after the language, GoogleSQL,
// which is how sqlc knows the dialect.
const (
	Dir   = "googlesql"
	Cases = "googlesql"
)

// Instance is the instance a Spanner Omni single server provides, which
// every database is created in.
const Instance = "projects/default/instances/default"

// Locate returns the server to generate from, named by SPANNER_SERVER_URI.
func Locate() (string, error) {
	endpoint := os.Getenv("SPANNER_SERVER_URI")
	if endpoint == "" {
		return "", errors.New("SPANNER_SERVER_URI is not set")
	}
	if strings.Contains(endpoint, "://") {
		return "", fmt.Errorf("SPANNER_SERVER_URI: %q should be a host:port, the gRPC endpoint of a Spanner Omni server", endpoint)
	}
	return endpoint, nil
}

// server holds the clients a run talks to the server through.
type server struct {
	endpoint  string
	instances *instance.InstanceAdminClient
	databases *database.DatabaseAdminClient
	data      *spannerapi.Client
}

// open connects to the server. Omni serves plaintext gRPC and asks for no
// credentials.
func open(ctx context.Context, endpoint string) (*server, error) {
	opts := []option.ClientOption{
		option.WithEndpoint(endpoint),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	}
	s := &server{endpoint: endpoint}
	var err error
	if s.instances, err = instance.NewInstanceAdminClient(ctx, opts...); err != nil {
		return nil, err
	}
	if s.databases, err = database.NewDatabaseAdminClient(ctx, opts...); err != nil {
		s.Close()
		return nil, err
	}
	if s.data, err = spannerapi.NewClient(ctx, opts...); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

func (s *server) Close() {
	for _, c := range []interface{ Close() error }{s.instances, s.databases, s.data} {
		if c != nil {
			c.Close()
		}
	}
}

// Version describes the server: Spanner Omni reports no release of its
// own through the API, so the instance the databases are created in is
// described instead.
func Version(ctx context.Context, endpoint string) (string, error) {
	s, err := open(ctx, endpoint)
	if err != nil {
		return "", err
	}
	defer s.Close()
	inst, err := s.instances.GetInstance(ctx, &instancepb.GetInstanceRequest{Name: Instance})
	if err != nil {
		return "", fmt.Errorf("%s: %w", Instance, err)
	}
	return fmt.Sprintf("Spanner Omni at %s, instance %s (%s)", endpoint, inst.Name, inst.State), nil
}

// createDatabase creates a database in the instance, running the DDL
// statements in it, and returns its full name.
func (s *server) createDatabase(ctx context.Context, name string, ddl []string) (string, error) {
	op, err := s.databases.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          Instance,
		CreateStatement: "CREATE DATABASE `" + name + "`",
		ExtraStatements: ddl,
	})
	if err != nil {
		return "", err
	}
	db, err := op.Wait(ctx)
	if err != nil {
		return "", err
	}
	return db.Name, nil
}

// dropDatabase drops a database by its full name. A database that does
// not exist is not an error, so that a run can clear the way for itself.
func (s *server) dropDatabase(ctx context.Context, db string) error {
	err := s.databases.DropDatabase(ctx, &databasepb.DropDatabaseRequest{Database: db})
	if err != nil && strings.Contains(err.Error(), "NotFound") {
		return nil
	}
	return err
}

// session opens a session on a database, multiplexed as Omni requires.
func (s *server) session(ctx context.Context, db string) (string, error) {
	sess, err := s.data.CreateSession(ctx, &spannerpb.CreateSessionRequest{
		Database: db,
		Session:  &spannerpb.Session{Multiplexed: true},
	})
	if err != nil {
		return "", err
	}
	return sess.Name, nil
}

// query runs a statement and returns its rows as lists of values.
func (s *server) query(ctx context.Context, session, sql string) ([][]*structpb.Value, error) {
	rs, err := s.data.ExecuteSql(ctx, &spannerpb.ExecuteSqlRequest{Session: session, Sql: sql})
	if err != nil {
		return nil, err
	}
	var rows [][]*structpb.Value
	for _, row := range rs.Rows {
		rows = append(rows, row.Values)
	}
	return rows, nil
}

// Generate reads the dialect from the server: relations.jsonl, the views of
// INFORMATION_SCHEMA and SPANNER_SYS, read from a database of their own,
// since every database has the same ones.
func Generate(ctx context.Context, endpoint string) (dialect.Files, error) {
	s, err := open(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	db := Instance + "/databases/goldeneye_dialect"
	if err := s.dropDatabase(ctx, db); err != nil {
		return nil, err
	}
	if _, err := s.createDatabase(ctx, "goldeneye_dialect", nil); err != nil {
		return nil, err
	}
	defer s.dropDatabase(context.WithoutCancel(ctx), db)
	session, err := s.session(ctx, db)
	if err != nil {
		return nil, err
	}
	relations, err := readRelations(ctx, s, session)
	if err != nil {
		return nil, err
	}
	blob, err := dialect.JSONL(relations)
	if err != nil {
		return nil, err
	}
	return dialect.Files{dialect.RelationsFile: blob}, nil
}
