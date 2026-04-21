-- name: NearestTo :many
SELECT x, y
FROM nearest_to(@p::point_input);
