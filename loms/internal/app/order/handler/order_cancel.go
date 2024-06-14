package handler

import (
	"context"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	order, err := h.orderService.OrderServiceGetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, err
	}

	for _, item := range order.Items {
		_ = h.stocksService.StocksServiceReserveCancel(ctx, item.SKU, item.Count)
	}

	_ = h.orderService.OrderServiceSetStatus(ctx, model.OrderID(req.OrderId), model.StatusCanceled)

	return &loms.OrderCancelResponse{}, nil
}
