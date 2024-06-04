package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"

	"route256/cart/internal/cart/model"
)

func (h *CartHandler) GetCartByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CartHandler.GetCartByUserID"
	log := h.logger.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	req, err := model.GetValidateUserRequest(r)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(slog.Int64("UserID", int64(req.UserID)))

	cart, err := h.cartService.GetCartByUserID(req.UserID)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	totalPrice := h.cartService.GetTotalPrice(cart)
	items := make([]model.Item, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].SKU < items[j].SKU
	})
	response := model.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}

	bytes, err := json.Marshal(response)
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
