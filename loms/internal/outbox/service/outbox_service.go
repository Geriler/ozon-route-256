package service

import (
	"context"

	repository "route256/loms/internal/outbox/repository/sqlc"
)

type OutboxRepository interface {
	SendMessage(ctx context.Context, callback func(ctx context.Context, message *repository.FetchNextMsgsRow) error) error
}

type OutboxService struct {
	repository OutboxRepository
}

func NewOutboxService(repository OutboxRepository) *OutboxService {
	return &OutboxService{
		repository: repository,
	}
}

func (s *OutboxService) SendMessage(ctx context.Context, callback func(ctx context.Context, message *repository.FetchNextMsgsRow) error) error {
	return s.repository.SendMessage(ctx, callback)
}
