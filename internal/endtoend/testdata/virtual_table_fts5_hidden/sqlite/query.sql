-- name: SearchRecipes :many
SELECT rowid, name
FROM recipes_fts
WHERE recipes_fts MATCH ?;

-- name: SearchRecipesRanked :many
SELECT rowid, name, rank
FROM recipes_fts
WHERE recipes_fts MATCH ?
ORDER BY rank;
