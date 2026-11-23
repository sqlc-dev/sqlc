                                                     
                             
                    
                                        
              
                                         
                         

-- Every ? is an argument of its own, numbered in the order the expanded query
-- sends them.
-- name: ListActiveAuthors :many
SELECT a.id
FROM authors x
JOIN (SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM `authors` WHERE org_id = CAST(? AS signed) AND deleted_at IS NULL) a ON a.id = x.id
JOIN (SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM `authors` WHERE org_id = CAST(? AS signed) AND deleted_at IS NULL) b ON b.id = x.id
WHERE x.name = ? AND x.bio <> ? AND @uservar IS NULL;
