-- name: CreateOrder :one
INSERT INTO orders (status, user_id)
VALUES ($1, $2)
RETURNING id;

-- name: InsertOrderItem :exec
INSERT INTO order_items (order_id, sku_id, items_count)
VALUES ($1, $2, $3)
ON CONFLICT (order_id, sku_id) DO UPDATE
    SET items_count = order_items.items_count + EXCLUDED.items_count;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: GetOrderInfo :one
SELECT id, status, user_id
FROM orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT sku_id, items_count
FROM order_items
WHERE order_id = $1;
