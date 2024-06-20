package service

import (
	"context"

	"route256/loms/internal/order/model"
)

type OrderRepository interface {
	SetStatus(ctx context.Context, orderID model.OrderID, status model.Status) error
	GetOrder(ctx context.Context, orderID model.OrderID) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) (model.OrderID, error)
}

type OrderService struct {
	orderRepository OrderRepository
}

func NewOrderService(orderRepository OrderRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
	}
}

func (os *OrderService) OrderServiceSetStatus(ctx context.Context, orderID model.OrderID, status model.Status) error {
	return os.orderRepository.SetStatus(ctx, orderID, status)
}

func (os *OrderService) OrderServiceGetOrder(ctx context.Context, orderID model.OrderID) (*model.Order, error) {
	return os.orderRepository.GetOrder(ctx, orderID)
}

func (os *OrderService) OrderServiceCreate(ctx context.Context, order *model.Order) (model.OrderID, error) {
	return os.orderRepository.Create(ctx, order)
}
