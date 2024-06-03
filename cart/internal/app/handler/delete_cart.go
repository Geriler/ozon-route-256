package handler

import (
	"log/slog"
	"net/http"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.cartService.DeleteCartByUserID(req.UserID)

	w.WriteHeader(http.StatusNoContent)
}
