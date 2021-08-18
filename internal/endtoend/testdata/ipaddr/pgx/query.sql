CREATE TABLE foo (
  present_ip    inet not null,
  nullable_ip   inet,
  present_cidr  cidr not null,
  nullable_cidr cidr
);

-- name: Get :many
SELECT * FROM foo LIMIT $1;
