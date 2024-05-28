package main

import (
	"net"
	"net/http"
	"os"

	"route256/cart/internal/app/handler"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/cart/service"
	"route256/cart/internal/lib/logger"
	"route256/cart/internal/middleware"
	product "route256/cart/internal/product/service"
)

func main() {
	log := logger.SetupLogger(logger.Local)

	conn, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	cartRepository := repository.NewInMemoryCartRepository()
	cartService := service.NewCartService(cartRepository)
	productService := product.NewProductService()
	cartHandler := handler.NewCartHandler(cartService, productService, log)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", cartHandler.AddItemsToCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", cartHandler.DeleteItemsFromCart)
	mux.HandleFunc("DELETE /user/{user_id}/cart", cartHandler.DeleteCartByUserID)
	mux.HandleFunc("GET /user/{user_id}/cart", cartHandler.GetCartByUserID)

	logWrapperHandler := middleware.NewLogWrapperHandler(mux, log)

	if err = http.Serve(conn, logWrapperHandler); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
