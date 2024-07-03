package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"route256/cart/internal"
	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) DeleteCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteCartByUserID"
	log := h.logger.With(slog.String("op", op))

	statusCode := http.StatusNoContent

	defer func(createdAt time.Time) {
		internal.SaveMetrics(time.Since(createdAt).Seconds(), "DELETE /user/{user_id}/cart", strconv.Itoa(statusCode))
	}(time.Now())

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		statusCode = http.StatusBadRequest
		return
	}

	err = h.cartHandler.DeleteCart(r.Context(), req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		statusCode = http.StatusInternalServerError
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
