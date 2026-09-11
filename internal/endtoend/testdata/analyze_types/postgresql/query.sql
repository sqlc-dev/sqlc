-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  $1::numeric(10,2) AS a,
  $2::int[] AS b,
  price::text AS c,
  $3::mood AS d,
  $4::varchar(20) AS e,
  $5::posint AS f,
  1::int8 AS g,
  $6::int4[][] AS h,
  $7::numeric(5,1) AS i,
  ARRAY[1, 2] AS j,
  $8::myschema.mood AS k,
  $9::varchar(10)[] AS l,
  $10::interval day to second AS m,
  $11::myschema.mood[] AS n,
  $12::mood[] AS o
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE price = $1 AND ints = $2 AND m = $3 AND p = $4 AND title = $5 AND grid = $6 AND sn = $7 AND fr = $8 AND ivd = $9;
