-- name: SetStatus :exec
UPDATE orders
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: GetOrder :one
SELECT id, user_id, status
FROM orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT item_id, count
FROM order_items
WHERE order_id = $1;

-- name: Create :one
INSERT INTO orders (id, user_id)
VALUES (nextval('order_id_manual_seq') + $1, $2)
RETURNING id;

-- name: AddItemToOrder :exec
INSERT INTO order_items (order_id, item_id, count, status)
VALUES ($1, $2, $3, 'added');
