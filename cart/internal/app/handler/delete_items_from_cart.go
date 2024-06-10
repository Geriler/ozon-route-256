package handler

import (
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteItemsFromCart(req *model.UserSKURequest) error {
	h.cartService.DeleteItemsFromCart(req.UserID, req.SKU)

	return nil
}
