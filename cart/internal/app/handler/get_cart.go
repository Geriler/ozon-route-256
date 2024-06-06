package handler

import (
	"sort"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) GetCart(req *model.UserRequest) (model.CartResponse, error) {
	cart, err := h.cartService.GetCartByUserID(req.UserID)
	if err != nil {
		return model.CartResponse{}, err
	}
	totalPrice := h.cartService.GetTotalPrice(cart)
	items := make([]model.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})
	return model.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
