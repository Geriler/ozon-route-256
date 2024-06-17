package handler

import (
	"context"
	"fmt"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

var ErrOrderCannotPaid = fmt.Errorf("%s", "order can't be paid")

func (h *OrderHandler) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	order, err := h.orderService.OrderServiceGetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, err
	}

	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		return nil, ErrOrderCannotPaid
	}

	for _, item := range order.Items {
		err = h.stocksService.StocksServiceReserveRemove(ctx, item.SKU, item.Count)
		if err != nil {
			return nil, err
		}
	}

	err = h.orderService.OrderServiceSetStatus(ctx, model.OrderID(req.OrderId), model.StatusPaid)
	if err != nil {
		return nil, err
	}

	return &loms.OrderPayResponse{}, nil
}
