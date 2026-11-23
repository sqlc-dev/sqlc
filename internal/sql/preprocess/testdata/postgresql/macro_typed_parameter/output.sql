                                                
                       
                          
                             
            
                                                            

-- name: InRange :many
SELECT * FROM prices WHERE (CAST(amount AS numeric(10, 2)) >= CAST($1 AS double precision) AND CAST(amount AS numeric(10, 2)) <= ($2 + 1));
