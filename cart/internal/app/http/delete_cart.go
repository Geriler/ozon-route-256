package http

import (
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) DeleteCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.cartHandler.DeleteCart(req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
