package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"route256/cart/internal/app/errors"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) DeleteItemsFromCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.DeleteItemsFromCart"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	userIdRaw := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdRaw)
	if err != nil || userId < 1 {
		log.Error(errors.ErrUserIdRequired, slog.String("userId", userIdRaw))
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusBadRequest, Message: errors.ErrUserIdRequired}})
		return
	}

	skuIdRaw := r.PathValue("sku_id")
	skuId, err := strconv.Atoi(skuIdRaw)
	if err != nil || skuId < 1 {
		log.Error(errors.ErrSkuIdRequired, slog.String("skuId", skuIdRaw))
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusBadRequest, Message: errors.ErrSkuIdRequired}})
		return
	}

	h.cartService.DeleteItemsFromCart(model.UserID(userId), model.SkuID(skuId))

	h.sendSuccessResponse(w, model.SuccessResponse{Message: "success delete items from cart"})
}
