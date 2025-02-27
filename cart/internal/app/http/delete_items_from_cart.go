package http

import (
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) DeleteItemsFromCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteItemsFromCart"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserSKURequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartHandler.DeleteItemsFromCart(r.Context(), req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
