-- name: CreateFeedFollow :one
WITH insert_feed_follow AS (
    INSERT INTO feed_follows(
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    ) values (
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
    RETURNING *
)
SELECT ff.*, u.name as Username, f.name as Feedname
FROM insert_feed_follow ff
INNER JOIN users u ON u.id = ff.user_id
INNER JOIN feeds f ON f.id = ff.feed_id;