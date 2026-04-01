-- name: UpsertServerSetColumnTypo :exec
INSERT INTO servers(code, name) VALUES ($1, $2)
ON CONFLICT (code)
DO UPDATE SET name_typo = 1111;

-- name: UpsertServerConflictTargetTypo :exec
INSERT INTO servers(code, name) VALUES ($1, $2)
ON CONFLICT (code_typo)
DO UPDATE SET name = 1111;

-- name: UpsertServerExcludedColumnTypo :exec
INSERT INTO servers(code, name) VALUES ($1, $2)
ON CONFLICT (code)
DO UPDATE SET name = EXCLUDED.name_typo;
