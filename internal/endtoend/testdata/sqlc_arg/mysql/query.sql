/* name: FuncParamIdent :many */
SELECT name FROM foo WHERE name = sqlc.arg(slug);

/* name: FuncParamString :many */
SELECT name FROM foo WHERE name = sqlc.arg('slug');

/* name: Complicated :many */
WITH names AS (SELECT name from foo WHERE foo.name = sqlc.arg('slug'))
SELECT name FROM names WHERE name IN (SELECT name FROM foo WHERE foo.name = sqlc.arg('slug'));
