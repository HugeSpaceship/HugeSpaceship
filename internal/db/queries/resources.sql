-- name: CheckResources :many
SELECT l.hash
from UNNEST(sqlc.arg(resources)::text[]) as l(hash)
         LEFT JOIN resources r on l.hash = r.hash
WHERE r.hash is null;

-- name: InsertResource :exec
INSERT INTO resources (uploader,size,resource_type,hash,backend,backend_name,created) VALUES (
    sqlc.arg(uploader)::uuid,
    sqlc.arg(size)::bigint,
    sqlc.arg(type)::resource_type,
    sqlc.arg(hash)::text,
    sqlc.arg(backend)::resource_backends,
    sqlc.arg(backendName)::text,
    NOW()
);