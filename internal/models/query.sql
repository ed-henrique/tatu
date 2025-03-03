-- name: AddSecret :one
INSERT INTO secrets (value) VALUES (?) RETURNING id;
