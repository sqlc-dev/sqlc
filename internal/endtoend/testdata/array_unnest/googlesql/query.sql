-- name: SingersByIDs :many
SELECT singer_id, name FROM singers WHERE singer_id IN UNNEST(@ids);

-- name: SingersNotIn :many
SELECT singer_id FROM singers WHERE singer_id NOT IN UNNEST(@ids);

-- name: SingerTags :many
SELECT singer_id, tag FROM singers, UNNEST(tags) AS tag;

-- name: AlbumTitles :one
SELECT singer_id, ARRAY(SELECT title FROM albums WHERE albums.singer_id = singers.singer_id) AS titles
FROM singers WHERE singer_id = @singer_id;

-- name: FirstTag :one
SELECT tags[OFFSET(0)] AS first_tag, tags[SAFE_ORDINAL(@n)] AS nth_tag, scores[SAFE_OFFSET(0)] AS first_score
FROM singers WHERE singer_id = @singer_id;
