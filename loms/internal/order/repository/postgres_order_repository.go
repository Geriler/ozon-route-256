package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"route256/loms/internal/middleware"
	"route256/loms/internal/order/model"
	repository "route256/loms/internal/order/repository/sqlc"
	ourboxRepo "route256/loms/internal/outbox/repository/sqlc"
	modelStocks "route256/loms/internal/stocks/model"
	"route256/loms/pkg/lib/tracing"
)

type PostgresOrderRepository struct {
	conn   *pgx.Conn
	cmd    *repository.Queries
	logger *slog.Logger
	outbox *ourboxRepo.Queries
}

func NewPostgresOrderRepository(conn *pgx.Conn, logger *slog.Logger) *PostgresOrderRepository {
	cmd := repository.New(conn)

	return &PostgresOrderRepository{
		conn:   conn,
		cmd:    cmd,
		logger: logger,
	}
}

func (r *PostgresOrderRepository) SetStatus(ctx context.Context, orderID model.OrderID, status model.Status) error {
	var err, commitErr, rollbackErr error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresOrderRepository.SetStatus")
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
			r.logger.Error("Error in PostgresOrderRepository.SetStatus.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	err = r.cmd.WithTx(tx).SetStatus(ctx, repository.SetStatusParams{
		Status: string(status),
		ID:     int32(orderID),
	})
	if err != nil {
		return err
	}

	err = r.outbox.WithTx(tx).CreateEvent(ctx, ourboxRepo.CreateEventParams{
		OrderID:   int32(orderID),
		EventType: string(status),
	})
	if err != nil {
		return err
	}

	commitErr = tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresOrderRepository.SetStatus.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}

func (r *PostgresOrderRepository) GetOrder(ctx context.Context, orderID model.OrderID) (*model.Order, error) {
	var err error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresOrderRepository.GetOrder")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "SELECT", requestStatus)
	}(time.Now())

	defer func() {
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			requestStatus = "error"
		}
	}()

	row, err := r.cmd.GetOrder(ctx, int32(orderID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
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
	var err, commitErr, rollbackErr error

	ctx, span := tracing.StartSpanFromContext(ctx, "PostgresOrderRepository.Create")
	defer span.End()

	requestStatus := "ok"
	defer func(createdAt time.Time) {
		middleware.ObserveRequestDatabaseDurationSeconds(time.Since(createdAt).Seconds(), "INSERT", requestStatus)
	}(time.Now())

	defer func() {
		if err != nil || commitErr != nil || (rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed)) {
			requestStatus = "error"
		}
	}()

	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr = tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			r.logger.Error("Error in PostgresOrderRepository.Create.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	orderID, err := r.cmd.WithTx(tx).Create(ctx, int32(order.UserID))
	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		err = r.cmd.WithTx(tx).AddItemToOrder(ctx, repository.AddItemToOrderParams{
			OrderID: orderID,
			ItemID:  int32(item.SKU),
			Count:   int32(item.Count),
		})
		if err != nil {
			return 0, err
		}
	}

	err = r.outbox.WithTx(tx).CreateEvent(ctx, ourboxRepo.CreateEventParams{
		OrderID:   orderID,
		EventType: string(model.StatusNew),
	})
	if err != nil {
		return 0, err
	}

	commitErr = tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresOrderRepository.Create.Commit",
			slog.String("error", commitErr.Error()),
		)
		return 0, commitErr
	}
	return model.OrderID(orderID), nil
}
