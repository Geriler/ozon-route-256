-- name: SetStatus :exec
UPDATE orders
SET status = $1
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
INSERT INTO orders (user_id, status)
VALUES ($1, 'new')
RETURNING id;

-- name: AddItemToOrder :exec
INSERT INTO order_items (order_id, item_id, count, status)
VALUES ($1, $2, $3, 'added');
