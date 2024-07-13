-- name: GetUserByName :one
SELECT
    users.*,
    users.entitled_slots - COUNT(s) AS free_slots,
    COUNT(s) AS used_slots
FROM users LEFT JOIN slots AS s ON s.uploader = users.id
WHERE username = sqlc.arg(name) GROUP BY users.id LIMIT 1;

