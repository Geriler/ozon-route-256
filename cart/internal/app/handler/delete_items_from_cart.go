package handler

import (
	"context"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteItemsFromCart(ctx context.Context, req *model.UserSKURequest) error {
	h.cartService.DeleteItemsFromCart(ctx, req.UserID, req.SKU)

	return nil
}
