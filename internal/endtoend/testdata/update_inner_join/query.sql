-- name: UpdateXWithY :exec
UPDATE x INNER JOIN y ON y.a = x.a SET x.b = y.b;
