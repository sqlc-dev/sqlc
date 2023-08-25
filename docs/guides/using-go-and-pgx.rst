================
Using Go and pgx
================

.. note::
   :code:`pgx/v5` is supported starting from v1.18.0.


pgx is a pure Go driver and toolkit for PostgreSQL. It's become the default
PostgreSQL package for many Gophers since lib/pq was put into maintenance mode.

^^^^^^^^^^^^^^^
Getting started
^^^^^^^^^^^^^^^

To start generating code that uses pgx, set the :code:`sql_package` field in
your :code:`sqlc.yaml` configuration file. Valid options are :code:`pgx/v4` or
:code:`pgx/v5`

.. code-block:: yaml

    version: "2"
    sql:
      - engine: "postgresql"
        queries: "query.sql"
        schema: "query.sql"
        gen:
          go:
            package: "db"
            sql_package: "pgx/v5"
            out: "db"

If you don't have an existing sqlc project on hand, create a directory with the
configuration file above and the following :code:`query.sql` file.

.. code-block:: sql

    CREATE TABLE authors (
      id   BIGSERIAL PRIMARY KEY,
      name text      NOT NULL,
      bio  text
    );

    -- name: GetAuthor :one
    SELECT * FROM authors
    WHERE id = $1 LIMIT 1;
    
    -- name: ListAuthors :many
    SELECT * FROM authors
    ORDER BY name;
    
    -- name: CreateAuthor :one
    INSERT INTO authors (
      name, bio
    ) VALUES (
      $1, $2
    )
    RETURNING *;
    
    -- name: DeleteAuthor :exec
    DELETE FROM authors
    WHERE id = $1;


Generating the code will now give you pgx-compatible database access methods.

.. code-block:: bash

   sqlc generate

^^^^^^^^^^^^^^^^^^^^^^^^^^
Generated code walkthrough
^^^^^^^^^^^^^^^^^^^^^^^^^^

The generated code is very similar to the code generated when using
:code:`lib/pq`. However, instead of using :code:`database/sql`, the code uses
pgx types directly.

.. code-block:: go

    package main
    
    import (
    	"context"
    	"fmt"
    	"os"
    
    	"github.com/jackc/pgx/v5"
        
    	"example.com/sqlc-tutorial/db"
    )
    
    func main() {
    	// urlExample := "postgres://username:password@localhost:5432/database_name"
    	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
    		os.Exit(1)
    	}
    	defer conn.Close(context.Background())

    	q := db.New(conn)
    
    	author, err := q.GetAuthor(context.Background(), 1)
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "GetAuthor failed: %v\n", err)
    		os.Exit(1)
    	}
    
    	fmt.Println(author.Name)
    }
