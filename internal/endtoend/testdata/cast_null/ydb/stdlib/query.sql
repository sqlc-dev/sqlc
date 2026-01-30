-- name: ListNullable :many
SELECT
  Cast(NULL AS Text?) AS a,
  Cast(NULL AS Int32?) AS b,
  Cast(NULL AS Int64?) AS c,
  Cast(NULL AS DateTime?) AS d
FROM foo;




