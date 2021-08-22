CREATE TABLE foo (
  present_ip    inet not null,
  nullable_ip   inet,
  present_cidr  cidr not null,
  nullable_cidr cidr
);

CREATE TABLE bar (
  addr          macaddr not null,
  nullable_addr macaddr
);

-- name: ListFoo :many
SELECT * FROM foo;

-- name: FindFooByIP :one
SELECT * FROM foo
WHERE present_ip = $1;

-- name: FindFooByCIDR :one
SELECT * FROM foo
WHERE present_cidr = $1;

-- name: ListBar :many
SELECT * FROM bar;

-- name: FindBarByAddr :one
SELECT * FROM bar
WHERE addr = $1;
