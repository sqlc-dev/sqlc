-- name: UpdateReturningOldStar :one
UPDATE foo SET bar = $1 WHERE id = $2 RETURNING OLD.*;

-- name: UpdateReturningNewStar :one
UPDATE foo SET bar = $1 WHERE id = $2 RETURNING NEW.*;

-- name: UpdateReturningOldNewCols :one
UPDATE foo SET bar = $1 WHERE id = $2 RETURNING OLD.bar, NEW.bar;

-- name: DeleteReturningOldStar :one
DELETE FROM foo WHERE id = $1 RETURNING OLD.*;
