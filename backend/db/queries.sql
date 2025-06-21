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

-- name: GetUserIDByUsername :one
SELECT id FROM user WHERE username = ?;

-- name: CreateScore :execresult
INSERT INTO score (user_id,value) VALUES(?,?);

-- name: GetTop10Scores :many
SELECT user.username,score.value FROM score JOIN user ON score.user_id = user.id ORDER BY value DESC limit 10;
