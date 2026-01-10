/* name: GetOne :one */
SELECT * FROM foo WHERE a = $a AND b = $b LIMIT 1;

/* name: GetAll :many */
SELECT * FROM foo;

/* name: GetAllAByB :many */
SELECT a FROM foo WHERE b = $b;
