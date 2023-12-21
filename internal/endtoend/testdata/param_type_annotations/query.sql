-- name: TestSqlcArg :one
-- @param foo text
SELECT * FROM test WHERE id = sqlc.arg('foo');

-- name: TestAt :one
-- @param foo integer
SELECT * FROM test WHERE name = @foo;

-- name: TestForceNotNull :one
-- @param foo! jsonb
SELECT * FROM test WHERE name = @foo;

-- name: TestForceNullable :one
-- @param foo? uuid
SELECT * FROM test WHERE id = @foo;

-- name: TestGibberish :one
-- @param foo? uuid sdfagyi
SELECT * FROM test WHERE id = @foo;
