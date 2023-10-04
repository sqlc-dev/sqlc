/* name: FuncParamIdent :many */
SELECT name FROM foo
WHERE name = sqlc.arg(slug)
  AND id IN (sqlc.slice(favourites));



/* name: FuncParamString :many */
SELECT name FROM foo
WHERE name = sqlc.arg('slug')
  AND id IN (sqlc.slice('favourites'));
