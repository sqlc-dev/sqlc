CREATE TABLE foo (
    bar citext,
    bat citext not null
);

-- name: GetCitexts :many
SELECT bar, bat
FROM foo;

