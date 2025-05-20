-- name: CreateFeedFollow :one
WITH feed_follow as (INSERT INTO feed_follows(id, created_at, updated_at, feed_id, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *)
SELECT 
feed_follow.*,
feeds.name as feed_name,
users.name as user_name
FROM feed_follow
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
INNER JOIN users ON feed_follow.user_id = users.id;
-- name: GetFeedFollowsByUserId :many
WITH feed_follow as(select * 
from feed_follows
where feed_follows.user_id = $1)
SELECT 
feed_follow.*,
feeds.name as feed_name,
users.name as user_name
FROM feed_follow
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
INNER JOIN users ON feed_follow.user_id = users.id;
;