package handler

import (
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteCart(req *model.UserRequest) error {
	h.cartService.DeleteCartByUserID(req.UserID)

	return nil
}
