                                                                       
                       
               
                                                                     

-- name: ListDocs :many
SELECT * FROM docs d
WHERE (d.owner_id = $1 OR $2::bool IS TRUE)
  AND d.kind = ANY($3)
  AND d.tenant_id = $1;
