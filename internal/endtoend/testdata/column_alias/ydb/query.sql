-- name: GetUsers :many
SELECT 
    users.id,
    users.fname,
    users.lname,
    users.email,
    users.created_at,
    CASE WHEN users.email LIKE '%' || $search_term || '%' THEN 1 ELSE 0 END AS rank_email,
    CASE WHEN users.fname LIKE '%' || $search_term || '%' THEN 1 ELSE 0 END AS rank_fname,
    CASE WHEN users.lname LIKE '%' || $search_term || '%' THEN 1 ELSE 0 END AS rank_lname,
    CASE WHEN (users.email || users.fname || users.lname) LIKE '%' || $search_term || '%' THEN 1 ELSE 0 END AS similarity
FROM users
WHERE users.email LIKE '%' || $search_term || '%' 
   OR users.fname LIKE '%' || $search_term || '%' 
   OR users.lname LIKE '%' || $search_term || '%'
ORDER BY rank_email DESC, rank_lname DESC, rank_fname DESC, similarity DESC;
