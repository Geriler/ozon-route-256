package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"route256/cart/internal/app/errors"
	"route256/cart/internal/cart/model"
)

func (h *CartHandler) AddItemsToCart(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.AddItemsToCart"
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

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}})
		return
	}
	defer r.Body.Close()

	var req model.AddItemsToCartRequest
	if err = json.Unmarshal(buf, &req); err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}})
		return
	}

	product, err := h.productService.GetProduct(model.SkuID(skuId))
	if err != nil {
		log.Error(err.Error())
		h.sendErrorResponse(w, model.ErrorResponse{Error: model.Error{Code: http.StatusBadRequest, Message: err.Error()}})
		return
	}

	item := model.Item{
		SKU:   model.SkuID(skuId),
		Name:  product.Name,
		Count: uint16(req.Count),
		Price: product.Price,
	}

	h.cartService.AddItemsToCart(model.UserID(userId), item)

	h.sendSuccessResponse(w, model.SuccessResponse{Message: "success add items to cart"})
}
