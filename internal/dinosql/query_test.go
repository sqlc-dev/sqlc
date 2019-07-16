package dinosql

import (
	"fmt"
	"testing"

	core "github.com/kyleconroy/dinosql/internal/pg"
	pg "github.com/lfittl/pg_query_go"

	"github.com/google/go-cmp/cmp"
)

func parseSQLTwo(in string) (QueryTwo, error) {
	tree, err := pg.Parse(in)
	if err != nil {
		return QueryTwo{}, err
	}
	c := core.NewCatalog()
	if err := updateCatalog(&c, tree); err != nil {
		return QueryTwo{}, err
	}

	q, _, err := parseQuery(c, tree.Statements[len(tree.Statements)-1], in)
	q.Stmt = nil
	return q, err
}

const ondeckSchema = `
CREATE TABLE city (
    slug text PRIMARY KEY,
    name text NOT NULL
);

CREATE TYPE status AS ENUM ('open', 'closed');

CREATE TABLE venue (
    id               SERIAL primary key,
	create_at        timestamp    not null,
    status           status       not null,
    slug             text         not null,
    name             varchar(255) not null,
    city             text         not null references city(slug),
    spotify_playlist varchar      not null,
    songkick_id      text
);
`

func TestQueries(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stmt  string
		query QueryTwo
	}{
		{
			"list_cities",
			`
			CREATE TABLE city (slug text primary key, name text not null);
			SELECT * FROM city ORDER BY name;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"get_city",
			ondeckSchema + `
			SELECT * FROM city WHERE slug = $1;
			`,
			QueryTwo{
				Params: []Parameter{{Number: 1, Name: "slug", Type: "text"}},
				Columns: []core.Column{
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"create_city",
			ondeckSchema + `
			INSERT INTO city (
				name,
				slug
			) VALUES (
				$1,
				$2
			) RETURNING *;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "name", Type: "text"},
					{Number: 2, Name: "slug", Type: "text"},
				},
				Columns: []core.Column{
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"create_city",
			ondeckSchema + `
			UPDATE city SET name = $2 WHERE slug = $1;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "slug", Type: "text"},
					{Number: 2, Name: "name", Type: "text"},
				},
			},
		},
		{
			"list_venues",
			ondeckSchema + `
			SELECT *
			FROM venue
			WHERE city = $1
			ORDER BY name;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
					{Name: "create_at", DataType: "pg_catalog.timestamp", NotNull: true},
					{Name: "status", DataType: "status", NotNull: true},
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "pg_catalog.varchar", NotNull: true},
					{Name: "city", DataType: "text", NotNull: true},
					{Name: "spotify_playlist", DataType: "pg_catalog.varchar", NotNull: true},
					{Name: "songkick_id", DataType: "text"},
				},
				Params: []Parameter{
					{Number: 1, Name: "city", Type: "text"},
				},
			},
		},
		{
			"delete_venue",
			ondeckSchema + `
			DELETE FROM venue
			WHERE slug = $1 AND slug = $1;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "slug", Type: "text"},
				},
			},
		},
		{
			"get_venue",
			ondeckSchema + `
			SELECT *
			FROM venue
			WHERE slug = $1 AND city = $2;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
					{Name: "create_at", DataType: "pg_catalog.timestamp", NotNull: true},
					{Name: "status", DataType: "status", NotNull: true},
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "pg_catalog.varchar", NotNull: true},
					{Name: "city", DataType: "text", NotNull: true},
					{Name: "spotify_playlist", DataType: "pg_catalog.varchar", NotNull: true},
					{Name: "songkick_id", DataType: "text"},
				},
				Params: []Parameter{
					{Number: 1, Name: "slug", Type: "text"},
					{Number: 2, Name: "city", Type: "text"},
				},
			},
		},
		{
			"create_venue",
			ondeckSchema + `
			INSERT INTO venue (
				slug,
				name,
				city,
				created_at,
				spotify_playlist,
				status
			) VALUES (
				$1,
				$2,
				$3,
				NOW(),
				$4,
				$5
			) RETURNING id;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{Number: 1, Type: "text", Name: "slug"},
					{Number: 2, Type: "pg_catalog.varchar", Name: "name"},
					{Number: 3, Type: "text", Name: "city"},
					{Number: 4, Type: "pg_catalog.varchar", Name: "spotifyPlaylist"},
					{Number: 5, Type: "status", Name: "status"},
				},
			},
		},
		{
			"update_venue_name",
			ondeckSchema + `
			UPDATE venue
			SET name = $2
			WHERE slug = $1
			RETURNING id;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{Number: 1, Type: "text", Name: "slug"},
					{Number: 2, Type: "pg_catalog.varchar", Name: "name"},
				},
			},
		},
		{
			"venue_count_by_city",
			ondeckSchema + `
			SELECT city, count(*)
			FROM venue
			GROUP BY 1
			ORDER BY 1;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "city", DataType: "text", NotNull: true},
					{Name: "count", DataType: "integer"},
				},
			},
		},
		{
			"alias",
			`
			CREATE TABLE bar (id serial not null);
			CREATE TABLE foo (id serial not null, bar serial references bar(id));
			
			DELETE FROM foo f USING bar b
			WHERE f.bar = b.id AND b.id = $1;
			`,
			QueryTwo{
				Params: []Parameter{{Number: 1, Name: "id", Type: "serial"}},
			},
		},
		{
			"star",
			`
			CREATE TABLE bar (bid serial not null);
			CREATE TABLE foo (fid serial not null);
			SELECT * FROM bar, foo;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "bid", DataType: "serial", NotNull: true},
					{Name: "fid", DataType: "serial", NotNull: true},
				},
			},
		},
		{
			"cte_count",
			`
			CREATE TABLE bar (ready bool not null);
			WITH all_count AS (
				SELECT count(*) FROM bar
			), ready_count AS (
				SELECT count(*) FROM bar WHERE ready = true
			)
			SELECT all_count.count, ready_count.count
			FROM all_count, ready_count;
			`,
			QueryTwo{
				Columns: []core.Column{
					{Name: "count", DataType: "integer", NotNull: false},
					{Name: "count", DataType: "integer", NotNull: false},
				},
			},
		},
		{
			"cte_filter",
			`
			CREATE TABLE bar (ready bool not null);
			WITH filter_count AS (
				SELECT count(*) FROM bar WHERE ready = $1
			)
			SELECT filter_count.count
			FROM filter_count;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "ready", Type: "bool"},
				},
				Columns: []core.Column{
					{Name: "count", DataType: "integer", NotNull: false},
				},
			},
		},
		{
			"update_set",
			`
			CREATE TABLE foo (name text not null, slug text not null);
			UPDATE foo SET name = $2 WHERE slug = $1;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "slug", Type: "text"},
					{Number: 2, Name: "name", Type: "text"},
				},
			},
		},
		{
			"insert_select",
			`
			CREATE TABLE bar (name text not null, ready bool not null);
			CREATE TABLE foo (name text not null, meta text not null);
			INSERT INTO foo (name, meta)
			SELECT name, $1
			FROM bar WHERE ready = $2;
			`,
			QueryTwo{
				Params: []Parameter{
					{Number: 1, Name: "meta", Type: "text"},
					{Number: 2, Name: "ready", Type: "bool"},
				},
			},
		},
	} {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			q, err := parseSQLTwo(test.stmt)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.query, q); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}

const testComparisonSQL = `
CREATE TABLE bar (id serial not null);
SELECT count(*) %s 0 FROM bar;
`

func TestComparisonOperators(t *testing.T) {
	for _, op := range []string{">", "<", ">=", "<=", "<>", "!=", "="} {
		o := op
		t.Run(o, func(t *testing.T) {
			q, err := parseSQLTwo(fmt.Sprintf(testComparisonSQL, o))
			if err != nil {
				t.Fatal(err)
			}
			expected := QueryTwo{
				Columns: []core.Column{
					{Name: "_", DataType: "bool", NotNull: true},
				},
			}
			if diff := cmp.Diff(expected, q); diff != "" {
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}
