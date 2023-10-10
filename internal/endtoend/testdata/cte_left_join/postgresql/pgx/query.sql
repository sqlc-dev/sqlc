-- name: BadQuery :exec
WITH
	q
		AS (
			SELECT
				authors.name, authors.bio
			FROM
				authors
				LEFT JOIN fake ON authors.name = fake.name
		)
SELECT
	*
FROM
	q AS c1
WHERE c1.name = $1;