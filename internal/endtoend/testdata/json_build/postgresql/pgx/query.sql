-- name: SelectJSONBuildObject :one
SELECT
  json_build_object(),
  json_build_object('foo'),
  json_build_object('foo', 1),
  json_build_object('foo', 1, 2),
  json_build_object('foo', 1, 2, 'bar');

-- name: SelectJSONBuildArray :one
SELECT 
  json_build_array(),
  json_build_array(1),
  json_build_array(1, 2),
  json_build_array(1, 2, 'foo'),
  json_build_array(1, 2, 'foo', 4);

-- name: SelectJSONBBuildObject :one
SELECT
  jsonb_build_object(),
  jsonb_build_object('foo'),
  jsonb_build_object('foo', 1),
  jsonb_build_object('foo', 1, 2),
  jsonb_build_object('foo', 1, 2, 'bar');

-- name: SelectJSONBBuildArray :one
SELECT 
  jsonb_build_array(),
  jsonb_build_array(1),
  jsonb_build_array(1, 2),
  jsonb_build_array(1, 2, 'foo'),
  jsonb_build_array(1, 2, 'foo', 4);
