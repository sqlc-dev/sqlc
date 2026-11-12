-- name: SearchBooks :many
SELECT id, title FROM books WHERE CONTAINS(title, @query);

-- name: SearchBooksFreeText :many
SELECT id, title FROM books WHERE FREETEXT(title, @phrase, LANGUAGE 1033) AND id > @min_id;
