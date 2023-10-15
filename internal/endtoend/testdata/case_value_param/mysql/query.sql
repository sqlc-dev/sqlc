-- name: Update :exec
UPDATE testing
SET value = CASE ? WHEN true THEN 'Hello' WHEN false THEN 'Goodbye' ELSE value END;
