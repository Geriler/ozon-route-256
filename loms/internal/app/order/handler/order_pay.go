package handler

import (
	"context"
	"errors"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

const ErrOrderCannotPaid = "order can't be paid"

func (h *OrderHandler) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	order, err := h.orderService.OrderServiceGetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		return nil, errors.New(ErrOrderCannotPaid)
	}

	for _, item := range order.Items {
		_ = h.stocksService.StocksServiceReserveRemove(ctx, item.SKU, item.Count)
	}

	_ = h.orderService.OrderServiceSetStatus(ctx, model.OrderID(req.OrderId), model.StatusPaid)

	return &loms.OrderPayResponse{}, nil
}
