-- name: GetUsers :many
SELECT 
    users.id,
    users.fname,
    users.lname,
    users.email,
    users.created_at,
    rank_email,
    rank_fname,
    rank_lname,
    similarity
FROM 
    users, 
    to_tsvector(users.email || users.fname || users.lname) document,
    to_tsquery(@search_term::TEXT) query,
    NULLIF(ts_rank(to_tsvector(users.email), query), 0) rank_email,
    NULLIF(ts_rank(to_tsvector(users.fname), query), 0) rank_fname,
    NULLIF(ts_rank(to_tsvector(users.lname), query), 0) rank_lname,
    SIMILARITY(@search_term::TEXT, users.email || users.fname || users.lname) similarity
WHERE query @@ document OR similarity > 0
ORDER BY rank_email, rank_lname, rank_fname, similarity DESC NULLS LAST;
