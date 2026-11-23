                      
                                               

-- name: ListBooks :many
WITH a AS (SELECT * FROM authors WHERE deleted_at IS NULL)
SELECT b.*, a.name FROM a JOIN books b ON b.author_id = a.id WHERE b.id > $1;
