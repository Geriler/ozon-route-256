package handler

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

var ErrOrderCannotCanceled = errors.New("order can't be canceled")

func (h *OrderHandler) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	ctx, span := h.tracer.Start(ctx, "OrderCancel", trace.WithAttributes(
		attribute.Int("order_id", int(req.OrderId)),
	))
	defer span.End()

	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderCancel", status)
	}(time.Now())

	span.AddEvent("Get order by ID")
	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		status = "error"
		return nil, err
	}

	span.AddEvent("Check order status")
	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		status = "error"
		return nil, ErrOrderCannotCanceled
	}

	span.AddEvent("Cancel order")
	err = h.stocksService.ReserveCancel(ctx, order.Items)
	if err != nil {
		status = "error"
		return nil, err
	}

	span.AddEvent("Set order status to canceled")
	err = h.orderService.SetStatus(ctx, model.OrderID(req.OrderId), model.StatusCanceled)
	if err != nil {
		status = "error"
		return nil, err
	}

	span.AddEvent("Return success message")
	return &loms.OrderCancelResponse{}, nil
}
