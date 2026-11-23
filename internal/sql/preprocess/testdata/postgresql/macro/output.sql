                                                      
                             
                    
        
            
                                          
                         

-- name: ListActiveAuthors :many
SELECT a.id FROM (SELECT * FROM authors WHERE org_id = CAST($1 AS bigint) AND deleted_at IS NULL) a WHERE a.name = $2;

-- name: CountAuthors :one
SELECT count(*) FROM authors WHERE org_id = $1;
