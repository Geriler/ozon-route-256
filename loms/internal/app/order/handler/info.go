package handler

import (
	"context"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, err
	}

	return &loms.OrderInfoResponse{
		Status: string(order.Status),
		UserId: order.UserID,
		Items:  model.ItemsToLomsItems(order.Items),
	}, nil
}
