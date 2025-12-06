-- Simple queries on employees table
-- name: GetAllEmployees :many
SELECT
    id,
    name,
    salary,
    department
FROM employees
ORDER BY id;

-- name: GetEmployeesByDepartment :many
SELECT
    id,
    name,
    salary
FROM employees
WHERE department = ?
ORDER BY salary DESC;
