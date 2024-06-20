package handler

import (
	"context"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	var err error
	items := model.LomsItemsToItems(req.Items)

	orderID := h.orderService.OrderServiceCreate(ctx, &model.Order{
		UserID: req.UserId,
		Status: model.StatusNew,
		Items:  items,
	})

	err = h.stocksService.StocksServiceReserve(ctx, items)
	if err != nil {
		errSetStatus := h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusFailed)
		if errSetStatus != nil {
			return nil, errSetStatus
		}

		return nil, err
	}

	err = h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusAwaitingPayment)
	if err != nil {
		return nil, err
	}

	return &loms.OrderCreateResponse{
		OrderId: int64(orderID),
	}, nil
}
