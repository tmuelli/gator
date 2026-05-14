-- name: GetFeedFollowsByUser :many
SELECT ff.*, f.name as Feedname, u.name as Username
FROM feed_follows ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;