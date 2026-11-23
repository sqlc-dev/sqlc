                     
           
            
            
                                                         

-- name: Clamped :many
SELECT (least(greatest(amount, $1), $2)) FROM prices WHERE id = $3;
