-- name: CreateUser :exec
INSERT INTO users (id, name, bio) VALUES (@id, @name, @bio);

-- name: CreateUserReturning :one
INSERT INTO users (id, name) VALUES (@id, @name) THEN RETURN id, name;

-- name: UpdateBio :exec
UPDATE users SET bio = @bio WHERE id = @id;

-- name: UpdateNameReturning :one
UPDATE users SET name = @name WHERE id = @id THEN RETURN id, name;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = @id;

-- name: DeleteUserReturning :one
DELETE FROM users WHERE id = @id THEN RETURN id;
