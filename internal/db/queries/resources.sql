-- name: CheckResources :many
SELECT l.hash
from UNNEST(sqlc.arg(resources)::text[]) as l(hash)
         LEFT JOIN resources r on l.hash = r.hash
WHERE r.hash is null;