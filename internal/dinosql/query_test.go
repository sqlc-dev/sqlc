package dinosql

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	core "github.com/kyleconroy/dinosql/internal/pg"
	pg "github.com/lfittl/pg_query_go"

	"github.com/google/go-cmp/cmp"
)

func parseSQL(in string) (Query, error) {
	tree, err := pg.Parse(in)
	if err != nil {
		return Query{}, err
	}
	c := core.NewCatalog()
	if err := updateCatalog(&c, tree); err != nil {
		return Query{}, err
	}

	q, err := parseQuery(c, tree.Statements[len(tree.Statements)-1], in)
	if q == nil {
		return Query{}, err
	}
	q.SQL = ""
	q.NeedsEdit = false
	return *q, err
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
		query Query
	}{
		{
			"list_cities",
			`
			CREATE TABLE city (slug text primary key, name text not null);
			SELECT * FROM city ORDER BY name;
			`,
			Query{
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
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "slug", DataType: "text", NotNull: true}},
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
			INSERT INTO city (
				name,
				slug
			) VALUES (
				$1,
				$2
			) RETURNING *;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "name", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "slug", DataType: "text", NotNull: true}},
				},
				Columns: []core.Column{
					{Name: "slug", DataType: "text", NotNull: true},
					{Name: "name", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"update_city",
			ondeckSchema + `
			UPDATE city SET name = $2 WHERE slug = $1;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "slug", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "name", DataType: "text", NotNull: true}},
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
			Query{
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
					{1, core.Column{Name: "city", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"delete_venue",
			ondeckSchema + `
			DELETE FROM venue
			WHERE slug = $1 AND slug = $1;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "slug", DataType: "text", NotNull: true}},
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
			Query{
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
					{1, core.Column{Name: "slug", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "city", DataType: "text", NotNull: true}},
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
			Query{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{NotNull: true, DataType: "text", Name: "slug"}},
					{2, core.Column{NotNull: true, DataType: "pg_catalog.varchar", Name: "name"}},
					{3, core.Column{NotNull: true, DataType: "text", Name: "city"}},
					{4, core.Column{NotNull: true, DataType: "pg_catalog.varchar", Name: "spotifyPlaylist"}},
					{5, core.Column{NotNull: true, DataType: "status", Name: "status"}},
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
			Query{
				Columns: []core.Column{
					{Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{DataType: "text", Name: "slug", NotNull: true}},
					{2, core.Column{DataType: "pg_catalog.varchar", Name: "name", NotNull: true}},
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
			Query{
				Columns: []core.Column{
					{Name: "city", DataType: "text", NotNull: true},
					{Name: "count", DataType: "bigint"},
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
			Query{
				Params: []Parameter{{1, core.Column{Name: "id", DataType: "serial", NotNull: true}}},
			},
		},
		{
			"star",
			`
			CREATE TABLE bar (bid serial not null);
			CREATE TABLE foo (fid serial not null);
			SELECT * FROM bar, foo;
			`,
			Query{
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
			Query{
				Columns: []core.Column{
					{Name: "count", DataType: "bigint", NotNull: false},
					{Name: "count", DataType: "bigint", NotNull: false},
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
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "ready", DataType: "bool", NotNull: true}},
				},
				Columns: []core.Column{
					{Name: "count", DataType: "bigint"},
				},
			},
		},
		{
			"update_set",
			`
			CREATE TABLE foo (name text not null, slug text not null);
			UPDATE foo SET name = $2 WHERE slug = $1;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "slug", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "name", DataType: "text", NotNull: true}},
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
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "meta", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "ready", DataType: "bool", NotNull: true}},
				},
			},
		},
		{
			"as",
			`
			CREATE TABLE foo (name text not null);
			SELECT name AS "other_name" FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Name: "other_name", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"text_array",
			`
			CREATE TABLE bar (tags text[] not null);
			SELECT * FROM bar;
			`,
			Query{
				Columns: []core.Column{
					{Name: "tags", DataType: "text", IsArray: true, NotNull: true},
				},
			},
		},
		{
			"select text array",
			`
			SELECT $1::TEXT[];
			`,
			Query{
				Columns: []core.Column{
					{Name: "", DataType: "text", IsArray: true, NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "", DataType: "text", NotNull: true, IsArray: true}},
				},
			},
		},
		{
			"select column cast",
			`
			CREATE TABLE foo (bar bool not null);
			SELECT bar::int FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Name: "bar", DataType: "pg_catalog.int4", NotNull: true},
				},
			},
		},
		{
			"limit",
			`
			CREATE TABLE foo (bar bool not null);
			SELECT bar FROM foo LIMIT $1;
			`,
			Query{
				Columns: []core.Column{
					{Name: "bar", DataType: "bool", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "limit", DataType: "integer", NotNull: true}},
				},
			},
		},
	} {
		test := tc
		t.Run(test.name, func(t *testing.T) {
			q, err := parseSQL(test.stmt)
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
			q, err := parseSQL(fmt.Sprintf(testComparisonSQL, o))
			if err != nil {
				t.Fatal(err)
			}
			expected := Query{
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

func TestStarWalker(t *testing.T) {
	for i, tc := range []struct {
		stmt     string
		expected bool
	}{
		{
			`
			SELECT * FROM city ORDER BY name;
			`,
			true,
		},
		{
			`
			INSERT INTO city (
				name,
				slug
			) VALUES (
				$1,
				$2
			) RETURNING *;
			`,
			true,
		},
		{
			`
			UPDATE city SET name = $2 WHERE slug = $1;
			`,
			false,
		},
		{
			`
			UPDATE venue
			SET name = $2
			WHERE slug = $1
			RETURNING *;
			`,
			true,
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tree, err := pg.Parse(test.stmt)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.expected, needsEdit(tree.Statements[0])); diff != "" {
				spew.Dump(tree.Statements[0])
				t.Errorf("query mismatch: \n%s", diff)
			}
		})
	}
}
