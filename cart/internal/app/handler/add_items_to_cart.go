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
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.productService.GetProduct(req.SKU)
	if err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	item := model.Item{
		SKU:   req.SKU,
		Name:  product.Name,
		Count: req.Count,
		Price: product.Price,
	}

	h.cartService.AddItemsToCart(req.UserID, item)

	h.sendSuccessResponse(w, model.SuccessResponse{Message: "success add items to cart"})
}
