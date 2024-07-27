package repository

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"route256/loms/internal/middleware"
	orderModel "route256/loms/internal/order/model"
	"route256/loms/internal/stocks/model"
	repository "route256/loms/internal/stocks/repository/sqlc"
	"route256/loms/pkg/lib/tracing"
)

type PostgresStocksRepository struct {
	conn   *pgxpool.Pool
	cmd    *repository.Queries
	logger *slog.Logger
}

func NewPostgresStocksRepository(conn *pgxpool.Pool, logger *slog.Logger) *PostgresStocksRepository {
	cmd := repository.New(conn)

	return &PostgresStocksRepository{
		conn:   conn,
		cmd:    cmd,
		logger: logger,
	}
}

func (r *PostgresStocksRepository) Reserve(ctx context.Context, items []*orderModel.Item) error {
	var err, commitErr, rollbackErr error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresStocksRepository.Reserve")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	defer func() {
		if err != nil || commitErr != nil || (rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed)) {
			requestStatus = "error"
		}
	}()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr = tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
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
			return err
		}
	}

	commitErr = tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresStocksRepository.Reserve.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) ReserveRemove(ctx context.Context, items []*orderModel.Item) error {
	var err, commitErr, rollbackErr error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresStocksRepository.ReserveRemove")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	defer func() {
		if err != nil || commitErr != nil || (rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed)) {
			requestStatus = "error"
		}
	}()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr = tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
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
			return err
		}
	}

	commitErr = tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresStocksRepository.ReserveRemove.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) ReserveCancel(ctx context.Context, items []*orderModel.Item) error {
	var err, commitErr, rollbackErr error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresStocksRepository.ReserveCancel")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "UPDATE", requestStatus)
	}(time.Now())

	defer func() {
		if err != nil || commitErr != nil || (rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed)) {
			requestStatus = "error"
		}
	}()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr = tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
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
			return err
		}
	}

	commitErr = tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresStocksRepository.ReserveCancel.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresStocksRepository) GetBySKU(ctx context.Context, sku model.SKU) (*model.Stocks, error) {
	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresStocksRepository.GetBySKU")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "SELECT", requestStatus)
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
