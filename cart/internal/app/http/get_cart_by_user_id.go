package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"route256/cart/internal"
	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.GetCartByUserID"
	log := h.logger.With(slog.String("op", op))

	statusCode := http.StatusNoContent

	defer func(createdAt time.Time) {
		internal.SaveMetrics(time.Since(createdAt).Seconds(), "GET /user/{user_id}/cart/list", strconv.Itoa(statusCode))
	}(time.Now())

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		statusCode = http.StatusBadRequest
		return
	}

	log = log.With(slog.Int64("UserID", int64(req.UserID)))

	response, err := h.cartHandler.GetCart(r.Context(), req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		statusCode = http.StatusNotFound
		return
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		statusCode = http.StatusInternalServerError
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		statusCode = http.StatusInternalServerError
		return
	}
}
