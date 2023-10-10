-- name: UpsertServer :exec
INSERT INTO servers(code, name) VALUES ($1, $2) 
ON CONFLICT (code) 
DO UPDATE SET name_typo = 1111;