package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"route256/cart/internal/app/errors"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.GetCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdRaw)
	if err != nil || userId == 0 {
		log.Error(errors.ErrUserIdRequired, slog.String("userId", userIdRaw))
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusBadRequest, Message: errors.ErrUserIdRequired}})
		return
	}

	cart := h.cartService.GetCartByUserID(model.UserID(userId))
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
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}})
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}})
		return
	}
}
