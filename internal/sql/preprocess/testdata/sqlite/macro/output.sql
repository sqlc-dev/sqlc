                             
                     
                                                                 

-- A numbered placeholder keeps its number, and the macro's arguments are
-- numbered around it.
-- name: ListActiveAuthors :many
SELECT a.id
FROM (SELECT * FROM authors WHERE org_id = CAST(?2 AS integer) AND deleted_at IS NULL) a
JOIN (SELECT * FROM authors WHERE org_id = CAST(?1 AS integer) AND deleted_at IS NULL) b ON b.id = a.id
WHERE a.name = ?3 AND a.id IN (/*SLICE:ids*/?);
