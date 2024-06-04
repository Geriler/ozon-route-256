package app

import (
	"net/http"

	"route256/cart/internal/app/handler"
)

func InitRoutes(handlers *handler.CartHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handlers.AddItemsToCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", handlers.DeleteItemsFromCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart", handlers.DeleteCartByUserID)
	mux.HandleFunc("GET /user/{user_id}/cart/list", handlers.GetCartByUserID)

	return mux
}
