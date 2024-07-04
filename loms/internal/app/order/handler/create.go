package handler

import (
	"context"
	"time"

	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderCreate", status)
	}(time.Now())

	var err error
	items := model.LomsItemsToItems(req.Items)

	orderID, err := h.orderService.Create(ctx, &model.Order{
		UserID: req.UserId,
		Status: model.StatusNew,
		Items:  items,
	})
	if err != nil {
		status = "error"
		return nil, err
	}

	for _, item := range items {
		item.OrderID = int64(orderID)
	}

	err = h.stocksService.Reserve(ctx, items)
	if err != nil {
		status = "error"
		errSetStatus := h.orderService.SetStatus(ctx, orderID, model.StatusFailed)
		if errSetStatus != nil {
			return nil, errSetStatus
		}

		return nil, err
	}

	err = h.orderService.SetStatus(ctx, orderID, model.StatusAwaitingPayment)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &loms.OrderCreateResponse{
		OrderId: int64(orderID),
	}, nil
}
