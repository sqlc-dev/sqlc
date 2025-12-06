-- ARRAY JOIN with arrays
-- name: GetProductTagsArray :many
SELECT
    id,
    name,
    arrayJoin(tags) as tag
FROM products
WHERE id = ?
ORDER BY tag;

-- Array functions
-- name: GetProductsWithArrayFunctions :many
SELECT
    id,
    name,
    length(tags) as tag_count,
    length(ratings) as rating_count
FROM products
ORDER BY tag_count DESC;
