package expander

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sqlc-dev/sqlc/internal/engine/dolphin"
	"github.com/sqlc-dev/sqlc/internal/engine/postgresql"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite"
	"github.com/sqlc-dev/sqlc/internal/sqltest/docker"
	"github.com/sqlc-dev/sqlc/internal/sqltest/native"
)

// PostgreSQLColumnGetter implements ColumnGetter for PostgreSQL using pgxpool.
type PostgreSQLColumnGetter struct {
	pool *pgxpool.Pool
}

func (g *PostgreSQLColumnGetter) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	conn, err := g.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	desc, err := conn.Conn().Prepare(ctx, "", query)
	if err != nil {
		return nil, err
	}

	columns := make([]string, len(desc.Fields))
	for i, field := range desc.Fields {
		columns[i] = field.Name
	}

	return columns, nil
}

// MySQLColumnGetter implements ColumnGetter for MySQL using the forked driver's StmtMetadata.
type MySQLColumnGetter struct {
	db *sql.DB
}

func (g *MySQLColumnGetter) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	conn, err := g.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var columns []string
	err = conn.Raw(func(driverConn any) error {
		preparer, ok := driverConn.(driver.ConnPrepareContext)
		if !ok {
			return fmt.Errorf("driver connection does not support PrepareContext")
		}

		stmt, err := preparer.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		meta, ok := stmt.(mysql.StmtMetadata)
		if !ok {
			return fmt.Errorf("prepared statement does not implement StmtMetadata")
		}

		for _, col := range meta.ColumnMetadata() {
			columns = append(columns, col.Name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return columns, nil
}

// SQLiteColumnGetter implements ColumnGetter for SQLite using the native ncruces/go-sqlite3 API.
type SQLiteColumnGetter struct {
	conn *sqlite3.Conn
}

func (g *SQLiteColumnGetter) GetColumnNames(ctx context.Context, query string) ([]string, error) {
	// Prepare the statement - this gives us column metadata without executing
	stmt, _, err := g.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Get column names from the prepared statement
	count := stmt.ColumnCount()
	columns := make([]string, count)
	for i := 0; i < count; i++ {
		columns[i] = stmt.ColumnName(i)
	}

	return columns, nil
}

func TestExpandPostgreSQL(t *testing.T) {
	ctx := context.Background()

	uri := os.Getenv("POSTGRESQL_SERVER_URI")
	if uri == "" {
		if err := docker.Installed(); err == nil {
			u, err := docker.StartPostgreSQLServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			uri = u
		} else if err := native.Supported(); err == nil {
			u, err := native.StartPostgreSQLServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			uri = u
		} else {
			t.Skip("POSTGRESQL_SERVER_URI is empty and neither Docker nor native installation is available")
		}
	}

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		t.Skipf("could not connect to database: %v", err)
	}
	defer pool.Close()

	// Create a test table
	_, err = pool.Exec(ctx, `
		DROP TABLE IF EXISTS authors;
		CREATE TABLE authors (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			bio TEXT
		);
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	defer pool.Exec(ctx, "DROP TABLE IF EXISTS authors")

	// Create the parser which also implements format.Dialect
	parser := postgresql.NewParser()

	// Create the expander
	colGetter := &PostgreSQLColumnGetter{pool: pool}
	exp := New(colGetter, parser, parser)

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple select star",
			query:    "SELECT * FROM authors",
			expected: "SELECT id, name, bio FROM authors;",
		},
		{
			name:     "select with no star",
			query:    "SELECT id, name FROM authors",
			expected: "SELECT id, name FROM authors", // No change, returns original
		},
		{
			name:     "select star with where clause",
			query:    "SELECT * FROM authors WHERE id = 1",
			expected: "SELECT id, name, bio FROM authors WHERE id = 1;",
		},
		{
			name:     "double star",
			query:    "SELECT *, * FROM authors",
			expected: "SELECT id, name, bio, id, name, bio FROM authors;",
		},
		{
			name:     "table qualified star",
			query:    "SELECT authors.* FROM authors",
			expected: "SELECT authors.id, authors.name, authors.bio FROM authors;",
		},
		{
			name:     "star in middle of columns",
			query:    "SELECT id, *, name FROM authors",
			expected: "SELECT id, id, name, bio, name FROM authors;",
		},
		{
			name:     "insert returning star",
			query:    "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING *",
			expected: "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, name, bio;",
		},
		{
			name:     "insert returning mixed",
			query:    "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, *",
			expected: "INSERT INTO authors (name, bio) VALUES ('John', 'A writer') RETURNING id, id, name, bio;",
		},
		{
			name:     "update returning star",
			query:    "UPDATE authors SET name = 'Jane' WHERE id = 1 RETURNING *",
			expected: "UPDATE authors SET name = 'Jane' WHERE id = 1 RETURNING id, name, bio;",
		},
		{
			name:     "delete returning star",
			query:    "DELETE FROM authors WHERE id = 1 RETURNING *",
			expected: "DELETE FROM authors WHERE id = 1 RETURNING id, name, bio;",
		},
		{
			name:     "cte with select star",
			query:    "WITH a AS (SELECT * FROM authors) SELECT * FROM a",
			expected: "WITH a AS (SELECT id, name, bio FROM authors) SELECT id, name, bio FROM a;",
		},
		{
			name:     "multiple ctes with dependency",
			query:    "WITH a AS (SELECT * FROM authors), b AS (SELECT * FROM a) SELECT * FROM b",
			expected: "WITH a AS (SELECT id, name, bio FROM authors), b AS (SELECT id, name, bio FROM a) SELECT id, name, bio FROM b;",
		},
		{
			name:     "count star not expanded",
			query:    "SELECT COUNT(*) FROM authors",
			expected: "SELECT COUNT(*) FROM authors", // No change - COUNT(*) should not be expanded
		},
		{
			name:     "count star with other columns",
			query:    "SELECT COUNT(*), name FROM authors GROUP BY name",
			expected: "SELECT COUNT(*), name FROM authors GROUP BY name", // No change
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := exp.Expand(ctx, tc.query)
			if err != nil {
				t.Fatalf("Expand failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestExpandMySQL(t *testing.T) {
	ctx := context.Background()

	source := os.Getenv("MYSQL_SERVER_URI")
	if source == "" {
		if err := docker.Installed(); err == nil {
			u, err := docker.StartMySQLServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			source = u
		} else if err := native.Supported(); err == nil {
			u, err := native.StartMySQLServer(ctx)
			if err != nil {
				t.Fatal(err)
			}
			source = u
		} else {
			t.Skip("MYSQL_SERVER_URI is empty and neither Docker nor native installation is available")
		}
	}

	db, err := sql.Open("mysql", source)
	if err != nil {
		t.Skipf("could not connect to MySQL: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		t.Skipf("could not ping MySQL: %v", err)
	}

	// Create a test table
	_, err = db.ExecContext(ctx, `DROP TABLE IF EXISTS authors`)
	if err != nil {
		t.Fatalf("failed to drop test table: %v", err)
	}
	_, err = db.ExecContext(ctx, `
		CREATE TABLE authors (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			bio TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	defer db.ExecContext(ctx, "DROP TABLE IF EXISTS authors")

	// Create the parser which also implements format.Dialect
	parser := dolphin.NewParser()

	// Create the expander
	colGetter := &MySQLColumnGetter{db: db}
	exp := New(colGetter, parser, parser)

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple select star",
			query:    "SELECT * FROM authors",
			expected: "SELECT id, name, bio FROM authors;",
		},
		{
			name:     "select with no star",
			query:    "SELECT id, name FROM authors",
			expected: "SELECT id, name FROM authors", // No change, returns original
		},
		{
			name:     "select star with where clause",
			query:    "SELECT * FROM authors WHERE id = 1",
			expected: "SELECT id, name, bio FROM authors WHERE id = 1;",
		},
		{
			name:     "table qualified star",
			query:    "SELECT authors.* FROM authors",
			expected: "SELECT authors.id, authors.name, authors.bio FROM authors;",
		},
		{
			name:     "double table qualified star",
			query:    "SELECT authors.*, authors.* FROM authors",
			expected: "SELECT authors.id, authors.name, authors.bio, authors.id, authors.name, authors.bio FROM authors;",
		},
		{
			name:     "star in middle of columns table qualified",
			query:    "SELECT id, authors.*, name FROM authors",
			expected: "SELECT id, authors.id, authors.name, authors.bio, name FROM authors;",
		},
		{
			name:     "count star not expanded",
			query:    "SELECT COUNT(*) FROM authors",
			expected: "SELECT COUNT(*) FROM authors", // No change - COUNT(*) should not be expanded
		},
		{
			name:     "count star with other columns",
			query:    "SELECT COUNT(*), name FROM authors GROUP BY name",
			expected: "SELECT COUNT(*), name FROM authors GROUP BY name", // No change
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := exp.Expand(ctx, tc.query)
			if err != nil {
				t.Fatalf("Expand failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestExpandSQLite(t *testing.T) {
	ctx := context.Background()

	// Create an in-memory SQLite database using native API
	conn, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatalf("could not open SQLite: %v", err)
	}
	defer conn.Close()

	// Create a test table
	err = conn.Exec(`
		CREATE TABLE authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			bio TEXT
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// Create the parser which also implements format.Dialect
	parser := sqlite.NewParser()

	// Create the expander using native SQLite column getter
	colGetter := &SQLiteColumnGetter{conn: conn}
	exp := New(colGetter, parser, parser)

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "simple select star",
			query:    "SELECT * FROM authors",
			expected: "SELECT id, name, bio FROM authors;",
		},
		{
			name:     "select with no star",
			query:    "SELECT id, name FROM authors",
			expected: "SELECT id, name FROM authors", // No change, returns original
		},
		{
			name:     "select star with where clause",
			query:    "SELECT * FROM authors WHERE id = 1",
			expected: "SELECT id, name, bio FROM authors WHERE id = 1;",
		},
		{
			name:     "double star",
			query:    "SELECT *, * FROM authors",
			expected: "SELECT id, name, bio, id, name, bio FROM authors;",
		},
		{
			name:     "table qualified star",
			query:    "SELECT authors.* FROM authors",
			expected: "SELECT authors.id, authors.name, authors.bio FROM authors;",
		},
		{
			name:     "star in middle of columns",
			query:    "SELECT id, *, name FROM authors",
			expected: "SELECT id, id, name, bio, name FROM authors;",
		},
		{
			name:     "count star not expanded",
			query:    "SELECT COUNT(*) FROM authors",
			expected: "SELECT COUNT(*) FROM authors", // No change - COUNT(*) should not be expanded
		},
		{
			name:     "count star with other columns",
			query:    "SELECT COUNT(*), name FROM authors GROUP BY name",
			expected: "SELECT COUNT(*), name FROM authors GROUP BY name", // No change
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := exp.Expand(ctx, tc.query)
			if err != nil {
				t.Fatalf("Expand failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
