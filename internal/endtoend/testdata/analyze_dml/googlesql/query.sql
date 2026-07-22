-- name: CreateUser :exec
INSERT INTO users (id, name, bio) VALUES (@id, @name, @bio);

-- name: CreateUserReturning :one
INSERT INTO users (id, name) VALUES (@id, @name) THEN RETURN id, name;

-- name: UpdateBio :exec
UPDATE users SET bio = @bio WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;
