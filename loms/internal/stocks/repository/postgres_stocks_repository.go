package repository

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	orderModel "route256/loms/internal/order/model"
	"route256/loms/internal/stocks/model"
	repository "route256/loms/internal/stocks/repository/sqlc"
)

type PostgresStocksRepository struct {
	conn *pgx.Conn
	cmd  *repository.Queries
}

func NewPostgresStocksRepository(conn *pgx.Conn) *PostgresStocksRepository {
	cmd := repository.New(conn)

	return &PostgresStocksRepository{
		conn: conn,
		cmd:  cmd,
	}
}

func (r *PostgresStocksRepository) Reserve(ctx context.Context, items []*orderModel.Item) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).Reserve(ctx, repository.ReserveParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (r *PostgresStocksRepository) ReserveRemove(ctx context.Context, items []*orderModel.Item) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).ReserveRemove(ctx, repository.ReserveRemoveParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (r *PostgresStocksRepository) ReserveCancel(ctx context.Context, items []*orderModel.Item) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).ReserveCancel(ctx, repository.ReserveCancelParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			return err
		}
	}

	tx.Commit(ctx)
	return nil
}

func (r *PostgresStocksRepository) GetBySKU(ctx context.Context, sku model.SKU) (*model.Stocks, error) {
	row, err := r.cmd.GetBySKU(ctx, int32(sku))
	if err != nil {
		return nil, err
	}

	return &model.Stocks{
		SKU:           model.SKU(row.Sku),
		TotalCount:    int64(row.TotalCount),
		ReservedCount: row.ReservedCount,
	}, nil
}
