-- name: CreateUser :execresult
INSERT INTO user (username) VALUES (?);

-- name: CreateUserWithPassword :execresult
INSERT INTO user (username,password_hash) VALUES(?,?); 

-- name: GetUser :one
SELECT * FROM user WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM user;

-- name: UpdateUser :exec
UPDATE user SET username = ? WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM user WHERE id = ?;
