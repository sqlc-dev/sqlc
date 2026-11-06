-- name: LeftSideOfLeftJoin :many
SELECT sqlc.embed(p), u.name FROM posts AS p
LEFT JOIN users AS u ON p.user_id = u.id;

-- name: RightSideOfRightJoin :many
SELECT sqlc.embed(u), p.id FROM posts AS p
RIGHT JOIN users AS u ON p.user_id = u.id;

-- name: InnerJoin :many
SELECT sqlc.embed(p), sqlc.embed(u) FROM posts AS p
INNER JOIN users AS u ON p.user_id = u.id;

-- name: SelfJoinRequiredAlias :many
SELECT sqlc.embed(a), b.name FROM users AS a
LEFT JOIN users AS b ON a.id = b.id;
