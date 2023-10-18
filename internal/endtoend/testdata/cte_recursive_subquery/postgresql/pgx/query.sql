-- name: GetLatestVersion :one
WITH RECURSIVE search_tree(id, chain_id, chain_counter) AS (
	SELECT base.id, base.id AS chain_id, 0 as chain_counter
	FROM versions AS base
	WHERE base.previous_version_id IS NULL
  	UNION ALL
	SELECT v.id, search_tree.chain_id, search_tree.chain_counter + 1
	FROM versions AS v
	INNER JOIN search_tree ON search_tree.id = v.previous_version_id
)
SELECT DISTINCT ON (search_tree.chain_id) 
	search_tree.id
FROM search_tree   
ORDER BY search_tree.chain_id, chain_counter DESC;

-- name: GetLatestVersionWithSubquery :one
SELECT id
FROM versions
WHERE versions.id IN (
  WITH RECURSIVE search_tree(id, chain_id, chain_counter) AS (
	SELECT base.id, base.id AS chain_id, 0 as chain_counter
	FROM versions AS base
	WHERE versions.previous_version_id IS NULL
	UNION ALL
	SELECT v.id, search_tree.chain_id, search_tree.chain_counter + 1
	FROM versions AS v
	INNER JOIN search_tree ON search_tree.id = v.previous_version_id 
  )
  SELECT DISTINCT ON (search_tree.chain_id) 
	search_tree.id
  FROM search_tree   
  ORDER BY search_tree.chain_id, chain_counter DESC
);
