package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"route256/cart/internal"
	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) DeleteItemsFromCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteItemsFromCart"
	log := h.logger.With(slog.String("op", op))

	statusCode := http.StatusNoContent

	defer func(createdAt time.Time) {
		internal.SaveMetrics(time.Since(createdAt).Seconds(), "DELETE /user/{user_id}/cart/{sku_id}", strconv.Itoa(statusCode))
	}(time.Now())

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserSKURequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		statusCode = http.StatusBadRequest
		return
	}

	err = h.cartHandler.DeleteItemsFromCart(r.Context(), req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		statusCode = http.StatusBadRequest
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
