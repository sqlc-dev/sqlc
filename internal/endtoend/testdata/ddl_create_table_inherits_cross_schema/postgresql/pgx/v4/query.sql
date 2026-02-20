-- name: GetAllParties :many
SELECT * FROM parent.party;

-- name: GetAllPeople :many
SELECT * FROM child.person;

