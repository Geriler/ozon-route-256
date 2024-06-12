package http

import (
	"log/slog"
	"net/http"

	"route256/cart/internal/app/handler"
	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) AddItemsToCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.AddItemsToCart"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserSKUCountRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(slog.Int64("UserID", int64(req.UserID))).
		With(slog.Int64("SKU", int64(req.SKU)))

	err = h.cartHandler.AddItemsToCart(r.Context(), req)
	if err != nil {
		log.Error(err.Error())

		if err.Error() == handler.ErrNotEnoughStock {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
