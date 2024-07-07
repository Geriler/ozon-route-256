package handler

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	ctx, span := h.tracer.Start(ctx, "OrderInfo", trace.WithAttributes(
		attribute.Int("order_id", int(req.OrderId)),
	))
	defer span.End()

	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderInfo", status)
	}(time.Now())

	span.AddEvent("Get order by ID")
	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		status = "error"
		return nil, err
	}

	span.AddEvent("Return order information")
	return &loms.OrderInfoResponse{
		Status: string(order.Status),
		UserId: order.UserID,
		Items:  model.ItemsToLomsItems(order.Items),
	}, nil
}
