package handler

import "route256/cart/internal/cart/model"

func (h *CartHandler) AddItemsToCart(req *model.UserSKUCountRequest) error {
	product, err := h.productService.GetProduct(req.SKU)
	if err != nil {
		return err
	}

	item := model.Item{
		SKU:   req.SKU,
		Name:  product.Name,
		Count: req.Count,
		Price: product.Price,
	}

	h.cartService.AddItemsToCart(req.UserID, item)

	return nil
}
