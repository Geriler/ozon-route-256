-- name: CreateEvent :exec
INSERT INTO outbox (order_id, event_type)
VALUES ($1, $2);

-- name: FetchNextMsgs :many
SELECT order_id, event_type
FROM outbox
WHERE status = 'pending'
ORDER BY created_at;

-- name: MarkAsSuccess :exec
UPDATE outbox
SET status = 'success',
    updated_at = NOW()
WHERE order_id = $1
  AND event_type = $2;

-- name: MarkAsError :exec
UPDATE outbox
SET status = 'error',
    updated_at = NOW()
WHERE order_id = $1
  AND event_type = $2;
