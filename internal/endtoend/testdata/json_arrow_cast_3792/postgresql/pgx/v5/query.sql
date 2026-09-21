-- name: ListJobs :many
-- A `data ->> 'key'` lookup can return SQL NULL when the key is missing.
-- A surrounding `::text` cast must preserve that nullability instead of
-- emitting a non-nullable Go string. See issue #3792.
SELECT id,
       (data ->> 'PhoneNumber')::text AS phone_number,
       (data ->> 'ContactName')::text AS contact_name,
       (data ->> 'State')::text       AS state
FROM jobs;
