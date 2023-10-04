/* name: FuncParamIdent :many */
SELECT name FROM foo WHERE name = sqlc.arg(slug);

/* name: FuncParamString :many */
SELECT name FROM foo WHERE name = sqlc.arg('slug');
