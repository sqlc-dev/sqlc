-- name: TestSelect :many
SELECT
  ID, USERNAME
FROM
  USERS;

-- name: TestInsert :exec
INSERT INTO USERS (
  ID, USERNAME
) VALUES (
  ?, ?
);

-- name: TestUpdate :exec
UPDATE USERS 
SET USERNAME = ?
WHERE ID = ?;

-- name: TestDelete :exec
DELETE FROM USERS WHERE ID = ?;
