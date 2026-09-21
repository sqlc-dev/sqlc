-- name: CreateAuthor :exec
insert into authors (id, name, bio)
with potential_authors as (
  select id, name, bio
  from dummy
)
select id, name, bio
from potential_authors;
