package handler

import (
	"context"

	"route256/cart/internal/cart/model"
	loms "route256/cart/pb/api"
)

func (h *CartHandler) Checkout(ctx context.Context, req *model.UserRequest) (model.CartCheckoutResponse, error) {
	cart, err := h.cartService.GetCartByUserID(ctx, req.UserID)
	if err != nil {
		return model.CartCheckoutResponse{}, err
	}

	items := make([]*loms.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, &loms.Item{
			SkuId: int64(item.SKU),
			Count: item.Count,
		})
	}

	orderCreateResponse, err := h.grpcClient.Order.OrderCreate(ctx, &loms.OrderCreateRequest{
		UserId: int64(req.UserID),
		Items:  items,
	})
	if err != nil {
		return model.CartCheckoutResponse{}, err
	}

	return model.CartCheckoutResponse{
		OrderID: orderCreateResponse.OrderId,
	}, nil
}
