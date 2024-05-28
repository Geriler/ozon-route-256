package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"route256/cart/internal/app/errors"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdRaw)
	if err != nil || userId < 1 {
		log.Error(errors.ErrUserIdRequired, slog.String("userId", userIdRaw))
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusBadRequest, Message: errors.ErrUserIdRequired}})
		return
	}

	h.cartService.DeleteCartByUserID(model.UserID(userId))

	h.sendSuccessResponse(w, model.SuccessResponse{Message: "success delete cart"})
}
