CREATE TABLE foo (name text not null);

/* name: FuncParamIdent :many */
SELECT name FROM foo WHERE name = sqlc_arg(slug);

/* name: FuncParamString :many */
SELECT name FROM foo WHERE name = sqlc_arg('slug');
