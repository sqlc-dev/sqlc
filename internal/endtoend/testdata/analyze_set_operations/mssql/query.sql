-- name: Everyone :many
SELECT id, name, code, score, joined FROM people WHERE name <> @skip
UNION ALL
SELECT id, title, code, budget, founded FROM orgs;

-- name: Names :many
SELECT name FROM people
UNION
SELECT nickname FROM people
UNION
SELECT NULL;

-- name: SharedNicknames :many
SELECT nickname, id FROM people
INTERSECT
SELECT title, id FROM orgs;

-- name: OnlyPeople :many
SELECT id, score FROM people
EXCEPT
SELECT id, budget FROM orgs;

