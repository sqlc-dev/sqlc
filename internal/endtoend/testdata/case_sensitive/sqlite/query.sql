CREATE TABLE contacts (
	pid	TEXT,
	CustomerName	TEXT
);

-- name: InsertContact :exec
INSERT INTO contacts (
    pid,
    CustomerName
)
VALUES (?,?)
;
