-- name: ListContacts :many
SELECT id, phone, code, total FROM contacts WHERE phone = @phone;

-- name: ContactByCode :one
SELECT TOP 1 c.id, c.total FROM contacts c LEFT JOIN contacts o ON o.code = c.code WHERE c.code = @code;
