package handler

import (
	"context"

	"route256/loms/internal/app/stocks/handler"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

type OrderService interface {
	SetStatus(ctx context.Context, orderID model.OrderID, status model.Status) error
	GetOrder(ctx context.Context, orderID model.OrderID) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) (model.OrderID, error)
}

type OrderHandler struct {
	loms.UnimplementedOrderServer
	orderService  OrderService
	stocksService handler.StocksService
}

func NewOrderHandler(orderService OrderService, stocksService handler.StocksService) *OrderHandler {
	return &OrderHandler{
		orderService:  orderService,
		stocksService: stocksService,
	}
}
