-- name: BadQuery :one
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
	q AS c1,
	q as c2;