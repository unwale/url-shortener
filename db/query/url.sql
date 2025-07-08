-- name: CreateUrl :one
INSERT INTO urls (original_url, short_url)
VALUES ($1, $2)
RETURNING id, original_url, short_url, created_at, updated_at;

-- name: GetUrlByShort :one
SELECT id, original_url, short_url, click_count, created_at, updated_at
FROM urls
WHERE short_url = $1;

-- name: IncrementClickCount :one
UPDATE urls
SET click_count = click_count + 1, updated_at = NOW()
WHERE short_url = $1
RETURNING id, original_url, short_url, click_count, created_at, updated_at;