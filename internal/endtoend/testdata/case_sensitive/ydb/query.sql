-- name: InsertContact :exec
INSERT INTO contacts (
    pid,
    CustomerName
)
VALUES ($pid, $customer_name);




