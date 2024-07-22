-- name: CreateEvent :exec
INSERT INTO outbox (order_id, event_type, status)
VALUES ($1, $2, 'pending');

-- name: FetchNextMsgs :many
SELECT order_id, event_type
FROM outbox
WHERE status = 'pending'
ORDER BY created_at;

-- name: UpdateStatus :exec
UPDATE outbox
SET status = $1,
    updated_at = NOW()
WHERE order_id = $2
  AND event_type = $3;
