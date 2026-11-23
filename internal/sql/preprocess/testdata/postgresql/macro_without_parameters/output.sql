                       
          

-- name: Names :many
SELECT coalesce(name, ('unnamed')) FROM authors;
