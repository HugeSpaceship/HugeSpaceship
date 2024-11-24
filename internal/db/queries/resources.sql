-- name: CheckResources :many
SELECT (l.hash)::text
from UNNEST(sqlc.arg(resources)::text[]) as l(hash)
         LEFT JOIN resources r on l.hash = r.hash
WHERE r.hash is null;

-- name: InsertResource :exec
INSERT INTO resources (uploader,size,resource_type,hash,created) VALUES (
    sqlc.arg(uploader)::uuid,
    sqlc.arg(size)::bigint,
    sqlc.arg(type)::resource_type,
    sqlc.arg(hash)::text,
    NOW()
);

-- name: CheckResource :one
SELECT EXISTS (
    SELECT hash FROM resources WHERE hash = sqlc.arg(hash)::text);

-- name: DeleteResource :exec
DELETE FROM resources WHERE hash = sqlc.arg(hash)::text;