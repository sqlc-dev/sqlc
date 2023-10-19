-- name: GetLatestVersionWithSubquery :one
SELECT * 
FROM versions
WHERE versions.id IN (
  WITH RECURSIVE search_tree(id) AS (
	SELECT id, 0 as chain_id, 0 as chain_counter
    FROM versions
  )
  SELECT DISTINCT ON (search_tree.chain_id) 
	search_tree.id
  FROM search_tree   
  ORDER BY search_tree.chain_id, chain_counter DESC
);
