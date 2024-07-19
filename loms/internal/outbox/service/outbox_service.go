package service

import (
	"context"

	repository "route256/loms/internal/outbox/repository/sqlc"
)

type OrderRepository interface {
	SendMessage(ctx context.Context, callback func(ctx context.Context, message *repository.FetchNextMsgsRow) error) error
}

type OutboxService struct {
	repository OrderRepository
}

func NewOutboxService(repository OrderRepository) *OutboxService {
	return &OutboxService{
		repository: repository,
	}
}

func (s *OutboxService) SendMessage(ctx context.Context, callback func(ctx context.Context, message *repository.FetchNextMsgsRow) error) error {
	return s.repository.SendMessage(ctx, callback)
}
