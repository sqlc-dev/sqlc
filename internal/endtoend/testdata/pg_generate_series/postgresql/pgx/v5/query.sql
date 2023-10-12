-- name: GenerateSeries :many
SELECT generate_series($1::timestamp, $2::timestamp, '10 hours');
