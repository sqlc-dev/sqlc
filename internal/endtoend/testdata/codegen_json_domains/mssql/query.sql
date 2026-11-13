-- name: ListContacts :many
SELECT id, phone, code, total FROM contacts WHERE phone = @phone;
