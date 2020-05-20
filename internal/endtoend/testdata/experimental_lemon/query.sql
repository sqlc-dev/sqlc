CREATE TABLE foo (
        bar text NOT NULL 
);

CREATE TABLE bar (
        baz text NOT NULL 
);

CREATE TABLE baz (name text);
ALTER TABLE baz ADD COLUMN email text;

-- name: ListFoo :many
SELECT bar FROM foo;


