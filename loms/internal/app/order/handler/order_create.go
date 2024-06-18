package handler

import (
	"context"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	var err error
	reserved := make([]*model.Item, 0)
	items := model.LomsItemsToItems(req.Items)

	orderID := h.orderService.OrderServiceCreate(ctx, &model.Order{
		UserID: req.UserId,
		Status: model.StatusNew,
		Items:  items,
	})

	for _, item := range items {
		err = h.stocksService.StocksServiceReserve(ctx, item.SKU, item.Count)
		if err != nil {
			errSetStatus := h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusFailed)
			if errSetStatus != nil {
				return nil, errSetStatus
			}

			for _, item := range reserved {
				err = h.stocksService.StocksServiceReserveCancel(ctx, item.SKU, item.Count)
				if err != nil {
					return nil, err
				}
			}

			return nil, err
		}
		reserved = append(reserved, &model.Item{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	err = h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusAwaitingPayment)
	if err != nil {
		return nil, err
	}

	return &loms.OrderCreateResponse{
		OrderId: int64(orderID),
	}, nil
}
