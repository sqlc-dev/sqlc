CREATE TABLE foo (a integer, b integer);

/* name: GetOne :one */
SELECT * FROM foo WHERE a = ? AND b = ? LIMIT 1;

/* name: GetAll :many */
SELECT * FROM foo;

/* name: GetAllAByB :many */
SELECT a FROM foo WHERE b = ?;
