                    
                                  

                    
                                  

                    
                                  

-- name: Direct :many
SELECT * FROM sqlc.macro.Self() s;

-- name: Mutual :many
SELECT * FROM sqlc.macro.Ping() p;
