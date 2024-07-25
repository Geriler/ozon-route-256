package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"route256/loms/internal/outbox/model"
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
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	for _, message := range messages {
		err = callback(ctx, &message)

		var status model.Status

		if err != nil {
			status = model.StatusError
		} else {
			status = model.StatusSuccess
		}

		markErr := r.cmd.WithTx(tx).UpdateStatus(ctx, repository.UpdateStatusParams{
			Status: pgtype.Text{
				String: string(status),
				Valid:  true,
			},
			OrderID:   message.OrderID,
			EventType: message.EventType,
		})
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

func (r *PostgresOutboxRepository) ClearOutbox(ctx context.Context, oldDataDuration time.Duration) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			r.logger.Error("Error in PostgresOutboxRepository.ClearOutbox.Rollback",
				slog.String("error", rollbackErr.Error()),
			)
		}
	}(tx, ctx)

	err = r.cmd.WithTx(tx).ClearOutbox(ctx, pgtype.Timestamp{
		Time:  time.Now().Add(-oldDataDuration),
		Valid: true,
	})
	if err != nil {
		return err
	}

	commitErr := tx.Commit(ctx)
	if commitErr != nil {
		r.logger.Error("Error in PostgresOutboxRepository.ClearOutbox.Commit",
			slog.String("error", commitErr.Error()),
		)
		return commitErr
	}
	return nil
}
