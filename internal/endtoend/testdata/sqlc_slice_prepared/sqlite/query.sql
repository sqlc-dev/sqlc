CREATE TABLE foo (id int not null, name text not null, bar text);

/* name: FuncParamIdent :many */
SELECT name FROM foo
WHERE name = sqlc.arg(slug)
  AND id IN (sqlc.slice(favourites));
