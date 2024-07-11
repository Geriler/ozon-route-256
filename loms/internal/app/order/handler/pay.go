package handler

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

var ErrOrderCannotPaid = errors.New("order can't be paid")

func (h *OrderHandler) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	ctx, span := h.tracer.Start(ctx, "OrderPay", trace.WithAttributes(
		attribute.Int("order_id", int(req.OrderId)),
	))
	defer span.End()

	span.AddEvent("Get order by ID")
	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, err
	}

	span.AddEvent("Check order status")
	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		return nil, ErrOrderCannotPaid
	}

	span.AddEvent("Remove stocks")
	err = h.stocksService.ReserveRemove(ctx, order.Items)
	if err != nil {
		return nil, err
	}

	span.AddEvent("Set order status to paid")
	err = h.orderService.SetStatus(ctx, model.OrderID(req.OrderId), model.StatusPaid)
	if err != nil {
		return nil, err
	}

	span.AddEvent("Return success message")
	return &loms.OrderPayResponse{}, nil
}
