// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package repository

import (
	"context"
)

const addItemToOrder = `-- name: AddItemToOrder :exec
INSERT INTO order_items (order_id, item_id, count, status)
VALUES ($1, $2, $3, 'added')
`

type AddItemToOrderParams struct {
	OrderID int32
	ItemID  int32
	Count   int32
}

func (q *Queries) AddItemToOrder(ctx context.Context, arg AddItemToOrderParams) error {
	_, err := q.db.Exec(ctx, addItemToOrder, arg.OrderID, arg.ItemID, arg.Count)
	return err
}

const create = `-- name: Create :one
INSERT INTO orders (id, user_id)
VALUES (nextval('order_id_manual_seq') + $1, $2)
RETURNING id
`

type CreateParams struct {
	Column1 interface{}
	UserID  int32
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) (int32, error) {
	row := q.db.QueryRow(ctx, create, arg.Column1, arg.UserID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const getOrder = `-- name: GetOrder :one
SELECT id, user_id, status
FROM orders
WHERE id = $1
`

type GetOrderRow struct {
	ID     int32
	UserID int32
	Status string
}

func (q *Queries) GetOrder(ctx context.Context, id int32) (GetOrderRow, error) {
	row := q.db.QueryRow(ctx, getOrder, id)
	var i GetOrderRow
	err := row.Scan(&i.ID, &i.UserID, &i.Status)
	return i, err
}

const getOrderItems = `-- name: GetOrderItems :many
SELECT item_id, count
FROM order_items
WHERE order_id = $1
`

type GetOrderItemsRow struct {
	ItemID int32
	Count  int32
}

func (q *Queries) GetOrderItems(ctx context.Context, orderID int32) ([]GetOrderItemsRow, error) {
	rows, err := q.db.Query(ctx, getOrderItems, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderItemsRow
	for rows.Next() {
		var i GetOrderItemsRow
		if err := rows.Scan(&i.ItemID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setStatus = `-- name: SetStatus :exec
UPDATE orders
SET status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
`

type SetStatusParams struct {
	Status string
	ID     int32
}

func (q *Queries) SetStatus(ctx context.Context, arg SetStatusParams) error {
	_, err := q.db.Exec(ctx, setStatus, arg.Status, arg.ID)
	return err
}
