                                                                           
                                
                       
           
            
            
                                                            

-- name: InRange :many
SELECT * FROM prices WHERE (1 <= ? AND ? <= `max`);
