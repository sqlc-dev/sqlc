-- name: Update :exec
UPDATE testing
SET value = CASE $1 WHEN true THEN 'Hello' WHEN false THEN 'Goodbye' ELSE value END;
