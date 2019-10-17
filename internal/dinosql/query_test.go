package dinosql

import (
	"fmt"
	"strconv"
	"testing"

	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	pg "github.com/lfittl/pg_query_go"
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

func public(rel string) core.FQN {
	return core.FQN{
		Catalog: "",
		Schema:  "public",
		Rel:     rel,
	}
}

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
					{Table: public("city"), Name: "slug", DataType: "text", NotNull: true},
					{Table: public("city"), Name: "name", DataType: "text", NotNull: true},
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
					{Table: public("city"), Name: "slug", DataType: "text", NotNull: true},
					{Table: public("city"), Name: "name", DataType: "text", NotNull: true},
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
					{Table: public("city"), Name: "slug", DataType: "text", NotNull: true},
					{Table: public("city"), Name: "name", DataType: "text", NotNull: true},
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
					{Table: public("venue"), Name: "id", DataType: "serial", NotNull: true},
					{Table: public("venue"), Name: "create_at", DataType: "pg_catalog.timestamp", NotNull: true},
					{Table: public("venue"), Name: "status", DataType: "status", NotNull: true},
					{Table: public("venue"), Name: "slug", DataType: "text", NotNull: true},
					{Table: public("venue"), Name: "name", DataType: "pg_catalog.varchar", NotNull: true},
					{Table: public("venue"), Name: "city", DataType: "text", NotNull: true},
					{Table: public("venue"), Name: "spotify_playlist", DataType: "pg_catalog.varchar", NotNull: true},
					{Table: public("venue"), Name: "songkick_id", DataType: "text"},
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
					{Table: public("venue"), Name: "id", DataType: "serial", NotNull: true},
					{Table: public("venue"), Name: "create_at", DataType: "pg_catalog.timestamp", NotNull: true},
					{Table: public("venue"), Name: "status", DataType: "status", NotNull: true},
					{Table: public("venue"), Name: "slug", DataType: "text", NotNull: true},
					{Table: public("venue"), Name: "name", DataType: "pg_catalog.varchar", NotNull: true},
					{Table: public("venue"), Name: "city", DataType: "text", NotNull: true},
					{Table: public("venue"), Name: "spotify_playlist", DataType: "pg_catalog.varchar", NotNull: true},
					{Table: public("venue"), Name: "songkick_id", DataType: "text"},
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
					{Table: public("venue"), Name: "id", DataType: "serial", NotNull: true},
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
					{Table: public("venue"), Name: "id", DataType: "serial", NotNull: true},
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
					{Table: public("venue"), Name: "city", DataType: "text", NotNull: true},
					{Name: "count", DataType: "bigint", NotNull: true},
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
			"table-name",
			`
			CREATE TABLE bar (id serial not null);
			CREATE TABLE foo (id serial not null, bar serial references bar(id));

			SELECT foo.id
			FROM foo
			JOIN bar ON foo.bar = bar.id
			WHERE bar.id = $1 AND foo.id = $2;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
					{2, core.Column{Name: "id", DataType: "serial", NotNull: true}},
				},
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
					{Table: public("bar"), Name: "bid", DataType: "serial", NotNull: true},
					{Table: public("foo"), Name: "fid", DataType: "serial", NotNull: true},
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
					{Name: "count", DataType: "bigint", NotNull: true},
					{Name: "count", DataType: "bigint", NotNull: true},
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
					{Name: "count", DataType: "bigint", NotNull: true},
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
			"update_set_multiple",
			`
			CREATE TABLE foo (name text not null, slug text not null);
			UPDATE foo SET (name, slug) = ($2, $1);
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
					{Table: public("foo"), Name: "other_name", DataType: "text", NotNull: true},
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
					{Table: public("bar"), Name: "tags", DataType: "text", IsArray: true, NotNull: true},
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
					{Table: public("foo"), Name: "bar", DataType: "bool", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "limit", DataType: "integer", NotNull: true}},
				},
			},
		},
		{
			"multifrom",
			`
			CREATE TABLE foo (email text not null);
			CREATE TABLE bar (login text not null);
			SELECT email FROM bar, foo WHERE login = $1;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "email", DataType: "text", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "login", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"column-as",
			`
			CREATE TABLE foo (email text not null);
			SELECT email AS id FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "id", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"join where clause",
			`
			CREATE TABLE foo (barid serial not null);
			CREATE TABLE bar (id serial not null, owner text not null);

			SELECT foo.*
			FROM foo
			JOIN bar ON bar.id = barid
			WHERE owner = $1;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "barid", DataType: "serial", NotNull: true, Scope: "foo"},
				},
				Params: []Parameter{
					{1, core.Column{Name: "owner", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"two joins",
			`
			CREATE TABLE foo (bar_id serial not null, baz_id serial not null);
			CREATE TABLE bar (id serial not null);
			CREATE TABLE baz (id serial not null);

			SELECT foo.*
			FROM foo
			JOIN bar ON bar.id = bar_id
			JOIN baz ON baz.id = baz_id;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "bar_id", DataType: "serial", NotNull: true, Scope: "foo"},
					{Table: public("foo"), Name: "baz_id", DataType: "serial", NotNull: true, Scope: "foo"},
				},
			},
		},
		{
			"coalesce",
			`
			CREATE TABLE foo (bar text);

			SELECT coalesce(bar, '') as login
			FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "login", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"cast coalesce",
			`
			CREATE TABLE foo (bar text);

			SELECT coalesce(bar, '')::text as login
			FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Name: "login", DataType: "text", NotNull: true},
				},
			},
		},
		{
			"in",
			`
			CREATE TABLE bar (id serial not null);

			SELECT *
			FROM bar
			WHERE id IN ($1, $2);
			`,
			Query{
				Columns: []core.Column{
					{Table: public("bar"), Name: "id", DataType: "serial", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
					{2, core.Column{Name: "id", DataType: "serial", NotNull: true}},
				},
			},
		},
		{
			"any",
			`
			CREATE TABLE bar (id bigserial not null);

			SELECT id
			FROM bar
			WHERE foo = ANY($1::bigserial[]);
			`,
			Query{
				Columns: []core.Column{
					{Table: public("bar"), Name: "id", DataType: "bigserial", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "", DataType: "bigserial", NotNull: true, IsArray: true}},
				},
			},
		},
		{
			"schema-scoped list",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (id serial not null);
			SELECT * FROM foo.bar;
			`,
			Query{
				Columns: []core.Column{
					{
						Table: core.FQN{Schema: "foo", Rel: "bar"},
						Name:  "id", DataType: "serial", NotNull: true,
					},
				},
			},
		},
		{
			"schema-scoped filter",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (id serial not null);
			SELECT * FROM foo.bar WHERE id = $1;
			`,
			Query{
				Columns: []core.Column{
					{
						Table:    core.FQN{Schema: "foo", Rel: "bar"},
						Name:     "id",
						DataType: "serial",
						NotNull:  true,
					},
				},
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
				},
			},
		},
		{
			"schema-scoped create",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (id serial not null, name text not null);
			INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING id;
			`,
			Query{
				Columns: []core.Column{
					{
						Table: core.FQN{Schema: "foo", Rel: "bar"},
						Name:  "id", DataType: "serial", NotNull: true,
					},
				},
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
					{2, core.Column{Name: "name", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"schema-scoped update",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (id serial not null, name text not null);
			UPDATE foo.bar SET name = $2 WHERE id = $1;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
					{2, core.Column{Name: "name", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"schema-scoped delete",
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.bar (id serial not null);
			DELETE FROM foo.bar WHERE id = $1;
			`,
			Query{
				Params: []Parameter{
					{1, core.Column{Name: "id", DataType: "serial", NotNull: true}},
				},
			},
		},
		{
			"lower",
			`
			CREATE TABLE foo (bar text not null, bat text not null);
			SELECT bar FROM foo WHERE bar = $1 AND LOWER(bat) = $2;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "bar", DataType: "text", NotNull: true},
				},
				Params: []Parameter{
					{1, core.Column{Name: "bar", DataType: "text", NotNull: true}},
					{2, core.Column{Name: "bat", DataType: "text", NotNull: true}},
				},
			},
		},
		{
			"identical-tables",
			`
			CREATE TABLE foo (id text not null);
			CREATE TABLE bar (id text not null);
			SELECT * FROM foo;
			`,
			Query{
				Columns: []core.Column{
					{Table: public("foo"), Name: "id", DataType: "text", NotNull: true},
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
					{Name: "", DataType: "bool", NotNull: true},
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
