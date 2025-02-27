package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHttpHandlers) Checkout(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.Checkout"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(slog.Int64("UserID", int64(req.UserID)))

	res, err := h.cartHandler.Checkout(r.Context(), req)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(res)
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
