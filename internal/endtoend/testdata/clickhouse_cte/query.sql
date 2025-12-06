-- name: GetHighEarners :many
WITH high_earners AS (
  SELECT id, name, salary FROM employees WHERE salary > ?
)
SELECT id, name, salary FROM high_earners ORDER BY salary DESC;

-- name: GetMultipleCTEs :many
WITH emp_data AS (
  SELECT id, name, salary FROM employees
),
filtered AS (
  SELECT id, name, salary FROM emp_data WHERE salary > ?
)
SELECT id, name FROM filtered;
