package handler

import (
	"context"
	"errors"
	"time"

	"route256/loms/internal"
	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

var ErrOrderCannotCanceled = errors.New("order can't be canceled")

func (h *OrderHandler) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	status := "ok"
	defer func(createdAt time.Time) {
		internal.SaveLomsMetrics(time.Since(createdAt).Seconds(), "/loms.api.OrderCancel", status)
	}(time.Now())

	order, err := h.orderService.GetOrder(ctx, model.OrderID(req.OrderId))
	if err != nil {
		status = "error"
		return nil, err
	}

	if order.Status == model.StatusFailed || order.Status == model.StatusCanceled || order.Status == model.StatusPaid {
		status = "error"
		return nil, ErrOrderCannotCanceled
	}

	err = h.stocksService.ReserveCancel(ctx, order.Items)
	if err != nil {
		status = "error"
		return nil, err
	}

	err = h.orderService.SetStatus(ctx, model.OrderID(req.OrderId), model.StatusCanceled)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &loms.OrderCancelResponse{}, nil
}
