CREATE TABLE IF NOT EXISTS employees
(
    id UInt32,
    name String,
    salary Float64,
    department String,
    hire_date DateTime
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE IF NOT EXISTS salaries_history
(
    id UInt32,
    employee_id UInt32,
    salary Float64,
    effective_date DateTime
)
ENGINE = MergeTree()
ORDER BY (id, employee_id, effective_date);
