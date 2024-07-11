package handler

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "OrderCreate", trace.WithAttributes(
		attribute.Int("user_id", int(req.UserId)),
	))

	span.AddEvent("Convert LomsItems to Items")
	var err error
	items := model.LomsItemsToItems(req.Items)

	span.AddEvent("Create order")
	orderID, err := h.orderService.Create(ctx, &model.Order{
		UserID: req.UserId,
		Status: model.StatusNew,
		Items:  items,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		item.OrderID = int64(orderID)
	}

	span.AddEvent("Reserve stocks")
	err = h.stocksService.Reserve(ctx, items)
	if err != nil {
		errSetStatus := h.orderService.SetStatus(ctx, orderID, model.StatusFailed)
		if errSetStatus != nil {
			return nil, errSetStatus
		}

		return nil, err
	}

	span.AddEvent("Set order status to awaiting payment")
	err = h.orderService.SetStatus(ctx, orderID, model.StatusAwaitingPayment)
	if err != nil {
		return nil, err
	}

	span.AddEvent("Return order ID")
	return &loms.OrderCreateResponse{
		OrderId: int64(orderID),
	}, nil
}
