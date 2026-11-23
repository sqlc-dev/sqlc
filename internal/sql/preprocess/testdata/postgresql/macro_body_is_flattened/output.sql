                                                                              
                                                 
                      
               
                     
                                          
                                           
             
                 
                   

-- name: ListRecent :many
SELECT *
FROM (SELECT id, 'not -- a comment' AS note, $tag$ /* nor this */ $tag$ AS tagged FROM events ORDER BY id DESC LIMIT CAST(10 AS int)) r
WHERE r.id > $1;
