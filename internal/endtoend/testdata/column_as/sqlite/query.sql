CREATE TABLE foo (email text not null);

/* name: ColumnAs :many */
SELECT email AS id FROM foo;
