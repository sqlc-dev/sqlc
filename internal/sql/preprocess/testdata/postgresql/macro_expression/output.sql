                        
                        
              
                                                

-- name: ListRecent :many
SELECT * FROM authors WHERE NOT (CAST(created_at AS timestamptz) > now() - 7 * interval '1 day') OR id = $1;
