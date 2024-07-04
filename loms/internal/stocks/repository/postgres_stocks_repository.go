package repository

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"route256/loms/internal"
	orderModel "route256/loms/internal/order/model"
	"route256/loms/internal/stocks/model"
	repository "route256/loms/internal/stocks/repository/sqlc"
)

type PostgresStocksRepository struct {
	conn   *pgx.Conn
	cmd    *repository.Queries
	logger *slog.Logger
}

func NewPostgresStocksRepository(conn *pgx.Conn, logger *slog.Logger) *PostgresStocksRepository {
	cmd := repository.New(conn)

	return &PostgresStocksRepository{
		conn:   conn,
		cmd:    cmd,
		logger: logger,
	}
}

func (r *PostgresStocksRepository) Reserve(ctx context.Context, items []*orderModel.Item) error {
	requestStatus := "ok"
	defer func(createdAt time.Time) {
		internal.SaveDatabaseMetrics(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		requestStatus = "error"
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			requestStatus = "error"
			r.logger.Error("Error in PostgresStocksRepository.Reserve.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).Reserve(ctx, repository.ReserveParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			requestStatus = "error"
			return err
		}
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		requestStatus = "error"
		r.logger.Error("Error in PostgresStocksRepository.Reserve.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) ReserveRemove(ctx context.Context, items []*orderModel.Item) error {
	requestStatus := "ok"
	defer func(createdAt time.Time) {
		internal.SaveDatabaseMetrics(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		requestStatus = "error"
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			requestStatus = "error"
			r.logger.Error("Error in PostgresStocksRepository.ReserveRemove.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).ReserveRemove(ctx, repository.ReserveRemoveParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			requestStatus = "error"
			return err
		}
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		requestStatus = "error"
		r.logger.Error("Error in PostgresStocksRepository.ReserveRemove.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) ReserveCancel(ctx context.Context, items []*orderModel.Item) error {
	requestStatus := "ok"
	defer func(createdAt time.Time) {
		internal.SaveDatabaseMetrics(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		requestStatus = "error"
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			requestStatus = "error"
			r.logger.Error("Error in PostgresStocksRepository.ReserveCancel.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	for _, item := range items {
		err = r.cmd.WithTx(tx).ReserveCancel(ctx, repository.ReserveCancelParams{
			OrderID: int32(item.OrderID),
			ItemID:  int32(item.SKU),
		})
		if err != nil {
			requestStatus = "error"
			return err
		}
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		requestStatus = "error"
		r.logger.Error("Error in PostgresStocksRepository.ReserveCancel.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) GetBySKU(ctx context.Context, sku model.SKU) (*model.Stocks, error) {
	requestStatus := "ok"
	defer func(createdAt time.Time) {
		internal.SaveDatabaseMetrics(time.Since(createdAt).Seconds(), "SELECT", requestStatus)
	}(time.Now())

	row, err := r.cmd.GetBySKU(ctx, int32(sku))
	if err != nil {
		requestStatus = "error"
		return nil, err
	}

	return &model.Stocks{
		SKU:           model.SKU(row.Sku),
		TotalCount:    int64(row.TotalCount),
		ReservedCount: row.ReservedCount,
	}, nil
}
