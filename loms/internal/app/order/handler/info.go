package handler

import (
	"context"
	"time"

	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderInfo", status)
	}(time.Now())

	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		status = "error"
		return nil, err
	}

	return &loms.OrderInfoResponse{
		Status: string(order.Status),
		UserId: order.UserID,
		Items:  model.ItemsToLomsItems(order.Items),
	}, nil
}
