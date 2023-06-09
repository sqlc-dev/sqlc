-- name: LoadFoo :many
SELECT * FROM foo WHERE id = $1;

-- name: LoadFooWithAliases :many
SELECT
  id AS aliased_id,
  other_id AS aliased_other_id,
  age AS aliased_age,
  balance AS aliased_balance,
  bio AS aliased_bio,
  about AS aliased_about
FROM foo
WHERE id = @named_parameter_id;
