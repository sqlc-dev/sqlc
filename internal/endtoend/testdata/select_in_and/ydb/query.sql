-- name: DeleteAuthor :exec
DELETE FROM
  books
WHERE
  author NOT IN (
    SELECT
      a.name
    FROM
      authors a
    WHERE
      a.age >= $min_author_age
  )
  AND translator NOT IN (
    SELECT
      t.name
    FROM
      translators t
    WHERE
      t.age >= $min_translator_age
  )
  AND year <= $max_year;
