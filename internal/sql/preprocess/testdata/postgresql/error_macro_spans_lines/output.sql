                                                                              
                                             
                    
    
       

                    
           
            

-- name: Body :one
SELECT sqlc.macro.Note();

-- name: Argument :one
SELECT sqlc.macro.Echo(1), sqlc.macro.Echo($$ two
lines $$);
