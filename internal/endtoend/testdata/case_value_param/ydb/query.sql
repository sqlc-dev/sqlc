-- name: Update :exec
UPDATE testing
SET value = CASE CAST($condition AS Bool) WHEN true THEN Utf8('Hello') WHEN false THEN Utf8('Goodbye') ELSE value END;



