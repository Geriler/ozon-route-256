package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"route256/loms/internal/order/model"
	repository "route256/loms/internal/order/repository/sqlc"
	modelStocks "route256/loms/internal/stocks/model"
)

type PostgresOrderRepository struct {
	conn *pgx.Conn
	cmd  *repository.Queries
}

func NewPostgresOrderRepository(conn *pgx.Conn) *PostgresOrderRepository {
	cmd := repository.New(conn)

	return &PostgresOrderRepository{
		conn: conn,
		cmd:  cmd,
	}
}

func (r *PostgresOrderRepository) SetStatus(ctx context.Context, orderID model.OrderID, status model.Status) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = r.cmd.WithTx(tx).SetStatus(ctx, repository.SetStatusParams{
		Status: string(status),
		ID:     int32(orderID),
	})
	if err != nil {
		return err
	}

	tx.Commit(ctx)
	return nil
}

func (r *PostgresOrderRepository) GetOrder(ctx context.Context, orderID model.OrderID) (*model.Order, error) {
	row, err := r.cmd.GetOrder(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, model.ErrOrderNotFound
	}

	orderItems, err := r.cmd.GetOrderItems(ctx, int32(orderID))
	if err != nil {
		return nil, err
	}

	var items []*model.Item
	for _, orderItem := range orderItems {
		items = append(items, &model.Item{
			OrderID: int64(orderID),
			SKU:     modelStocks.SKU(orderItem.ItemID),
			Count:   int64(orderItem.Count),
		})
	}

	return &model.Order{
		UserID: int64(row.UserID),
		Status: model.Status(row.Status),
		Items:  items,
	}, nil
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *model.Order) (model.OrderID, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	create, err := r.cmd.WithTx(tx).Create(ctx, int32(order.UserID))
	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		err = r.cmd.WithTx(tx).AddItemToOrder(ctx, repository.AddItemToOrderParams{
			OrderID: create,
			ItemID:  int32(item.SKU),
			Count:   int32(item.Count),
		})
		if err != nil {
			return 0, err
		}
	}

	tx.Commit(ctx)
	return model.OrderID(create), nil
}
