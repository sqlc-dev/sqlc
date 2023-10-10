-- name: Percentile :one
select percentile_disc(0.5) within group (order by authors.name)
from authors;
