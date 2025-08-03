-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;