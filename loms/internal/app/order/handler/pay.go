package handler

import (
	"context"
	"errors"
	"time"

	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

var ErrOrderCannotPaid = errors.New("order can't be paid")

func (h *OrderHandler) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderPay", status)
	}(time.Now())

	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		status = "error"
		return nil, err
	}

	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		status = "error"
		return nil, ErrOrderCannotPaid
	}

	err = h.stocksService.ReserveRemove(ctx, order.Items)
	if err != nil {
		status = "error"
		return nil, err
	}

	err = h.orderService.SetStatus(ctx, model.OrderID(req.OrderId), model.StatusPaid)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &loms.OrderPayResponse{}, nil
}
