-- name: GetUsersWithProfiles :many
SELECT id, name FROM users WHERE EXISTS (SELECT 1 FROM profiles WHERE profiles.user_id = users.id);

-- name: GetUsersWithoutProfiles :many
SELECT id, name FROM users WHERE NOT EXISTS (SELECT 1 FROM profiles WHERE profiles.user_id = users.id);
