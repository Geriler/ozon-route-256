package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	repository "route256/loms/internal/outbox/repository/sqlc"
)

type PostgresOutboxRepository struct {
	conn   *pgx.Conn
	cmd    *repository.Queries
	logger *slog.Logger
}

func NewPostgresOutboxRepository(conn *pgx.Conn, logger *slog.Logger) *PostgresOutboxRepository {
	cmd := repository.New(conn)

	return &PostgresOutboxRepository{
		conn:   conn,
		cmd:    cmd,
		logger: logger,
	}
}

func (r *PostgresOutboxRepository) SendMessage(ctx context.Context, callback func(ctx context.Context, message *repository.FetchNextMsgsRow) error) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			r.logger.Error("Error in PostgresOutboxRepository.SendMessages.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	messages, err := r.cmd.WithTx(tx).FetchNextMsgs(ctx)
	if err != nil && err.Error() == "no rows in result set" {
		return nil
	}
	if err != nil {
		return err
	}

	for _, message := range messages {
		err = callback(ctx, &message)

		var markErr error
		if err != nil {
			markErr = r.cmd.WithTx(tx).MarkAsError(ctx, repository.MarkAsErrorParams{
				OrderID:   message.OrderID,
				EventType: message.EventType,
			})
		} else {
			markErr = r.cmd.WithTx(tx).MarkAsSuccess(ctx, repository.MarkAsSuccessParams{
				OrderID:   message.OrderID,
				EventType: message.EventType,
			})
		}
		if markErr != nil {
			return err
		}
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresOutboxRepository.SendMessages.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}
