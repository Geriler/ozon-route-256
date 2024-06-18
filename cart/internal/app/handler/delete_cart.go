package handler

import (
	"context"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteCart(ctx context.Context, req *model.UserRequest) error {
	h.cartService.DeleteCartByUserID(ctx, req.UserID)

	return nil
}
