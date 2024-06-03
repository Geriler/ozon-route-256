package handler

import (
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) AddItemsToCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.AddItemsToCart"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserSKUCountRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProduct(req.SKU)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item := model.Item{
		SKU:   req.SKU,
		Name:  product.Name,
		Count: req.Count,
		Price: product.Price,
	}

	h.cartService.AddItemsToCart(req.UserID, item)

	w.WriteHeader(http.StatusNoContent)
}
