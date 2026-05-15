-- name: LeftJoinSubquery :many
SELECT * FROM a AS table_a
LEFT JOIN (SELECT * FROM b WHERE b.name IS NOT NULL) si ON si.a_id = table_a.id;

-- name: LeftJoinSubqueryExplicitColumns :many
SELECT
    table_a.id,
    table_a.name,
    si.id,
    si.a_id,
    si.name
FROM a AS table_a
LEFT JOIN (SELECT id, a_id, name FROM b WHERE b.name IS NOT NULL) si ON si.a_id = table_a.id;

-- name: LeftJoinSubqueryNoAlias :many
SELECT * FROM a
LEFT JOIN (SELECT * FROM b) subquery ON subquery.a_id = a.id;

-- name: RightJoinSubquery :many
SELECT * FROM a AS table_a
RIGHT JOIN (SELECT * FROM b WHERE b.name IS NOT NULL) si ON si.a_id = table_a.id;

-- name: FullOuterJoinSubquery :many
SELECT * FROM a AS table_a
FULL OUTER JOIN (SELECT * FROM b WHERE b.name IS NOT NULL) si ON si.a_id = table_a.id;
