package handler

import (
	"context"

	"go.opentelemetry.io/otel/trace"
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
	tracer        trace.Tracer
}

func NewOrderHandler(orderService OrderService, stocksService handler.StocksService, tracer trace.Tracer) *OrderHandler {
	return &OrderHandler{
		orderService:  orderService,
		stocksService: stocksService,
		tracer:        tracer,
	}
}
