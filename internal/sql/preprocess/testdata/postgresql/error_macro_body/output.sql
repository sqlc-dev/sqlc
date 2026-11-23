                                                                            
                                            
                          
                                

                        
            
                                           

                      
             
                                                  

                      
                                     

-- name: UsePositional :many
SELECT * FROM sqlc.macro.Positional() p;

-- name: UseNullable :many
SELECT * FROM sqlc.macro.Nullable(1) n;

-- name: UseSliced :many
SELECT * FROM sqlc.macro.Sliced(sqlc.arg(ids)) s;

-- An error in an argument is reported at the argument.
-- name: BadArgument :many
SELECT * FROM sqlc.macro.Nullable(1 + sqlc.macro.Nowhere()) n;
