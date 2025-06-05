-- name: CreatePost :one
INSERT INTO posts(id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
     $1,
     $2,
     $3,
     $4,
     $5,
     $6,
     $7,
     $8
)
RETURNING *;

-- name: GetPostsForUser :many

WITH feed AS(
     SELECT feed_id
     FROM feed_follows
     WHERE feed_follows.user_id = $1
)
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id 
FROM posts
INNER JOIN feed ON posts.feed_id = feed.feed_id
ORDER BY posts.published_at
LIMIT $2;

-- name: GetPosts :many
select * 
from posts;