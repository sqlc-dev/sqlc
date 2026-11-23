                     
           
                

                     
           
                                                    

-- A macro may be an argument to itself: an argument belongs to the caller.
-- name: Nested :one
SELECT (((($1 + 1) + 1) + 1) * (2 + 1));
