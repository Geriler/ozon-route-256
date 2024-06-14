package handler

import (
	"context"

	"route256/loms/internal/order/model"
	loms "route256/loms/pb/api"
)

func (h *OrderHandler) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
	reserved := make([]*model.Item, len(req.Items))
	items := model.LomsItemsToItems(req.Items)

	orderID := h.orderService.OrderServiceCreate(ctx, &model.Order{
		UserID: req.UserId,
		Status: model.StatusNew,
		Items:  items,
	})

	for _, item := range items {
		err := h.stocksService.StocksServiceReserve(ctx, item.SKU, item.Count)
		if err != nil {
			_ = h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusFailed)

			for _, item := range reserved {
				_ = h.stocksService.StocksServiceReserveCancel(ctx, item.SKU, item.Count)
			}

			return nil, err
		}
		reserved = append(reserved, &model.Item{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	_ = h.orderService.OrderServiceSetStatus(ctx, orderID, model.StatusAwaitingPayment)

	return &loms.OrderCreateResponse{
		OrderId: int64(orderID),
	}, nil
}
