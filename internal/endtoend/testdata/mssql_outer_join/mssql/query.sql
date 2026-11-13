-- name: LeftJoin :many
SELECT a.id, a.name, b.title FROM authors a LEFT JOIN books b ON b.author_id = a.id;

-- name: RightJoin :many
SELECT a.name, b.title FROM authors a RIGHT JOIN books b ON b.author_id = a.id;

-- name: FullJoin :many
SELECT a.name, b.title FROM authors a FULL OUTER JOIN books b ON b.author_id = a.id;

-- name: LeftJoinStar :many
SELECT a.name, b.* FROM authors a LEFT JOIN books b ON b.author_id = a.id;
