-- name: GetFeeds :many
SELECT f.name as FeedName, f.url as FeedUrl, u.name as UserName
FROM feeds f
JOIN users u ON f.user_id = u.id;