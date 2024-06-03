package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.GetCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart := h.cartService.GetCartByUserID(req.UserID)
	totalPrice := h.cartService.GetTotalPrice(cart)
	items := make([]model.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, *item)
	}
	response := model.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
