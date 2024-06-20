-- name: Reserve :exec
UPDATE order_items
SET status = 'reserved'
WHERE order_id = $1 AND item_id = $2;

-- name: ReserveRemove :exec
UPDATE order_items
SET status = 'sold'
WHERE order_id = $1 AND item_id = $2;

-- name: ReserveCancel :exec
UPDATE order_items
SET status = 'canceled'
WHERE order_id = $1 AND item_id = $2;

-- name: GetBySKU :one
WITH sold_order_items AS (
  SELECT item_id, count
  FROM order_items oi
  WHERE oi.item_id = $1 AND oi.status = 'sold'
), reserved_order_items AS (
  SELECT item_id, count
  FROM order_items oi
  WHERE oi.item_id = $1 AND oi.status = 'reserved'
)
SELECT si.sku, (si.total_count - COALESCE(soi.count, 0)) AS total_count, SUM(COALESCE(roi.count, 0)) AS reserved_count
FROM stocks_items si
LEFT JOIN sold_order_items soi ON soi.item_id = si.sku
LEFT JOIN reserved_order_items roi ON roi.item_id = si.sku
WHERE si.sku = $1
GROUP BY si.sku, si.total_count, soi.count;
