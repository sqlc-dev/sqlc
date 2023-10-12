-- name: UpdateArticles :exec
UPDATE
	public.articles
SET
	is_deleted = TRUE
WHERE
	id = $1;