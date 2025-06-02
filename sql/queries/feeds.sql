-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
-- name: GetAllFeeds :many
select * from feeds;
-- name: GetFeedByUrl :one
select * 
from feeds 
where url = $1;
-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;
-- name: GetNextFeedToFetch :one
select * 
from feeds
order by last_fetched_at ASC NULLS FIRST
LIMIT 1;