-- name: GetAuthor :one
WITH person AS(
    WITH summary AS(
        WITH mb AS(
            select
                id, name, bio
            from authors
        )
        SELECT
            bio
        FROM mb
        )
    SELECT
       count(*) as total
    FROM summary
)
select total from person;
