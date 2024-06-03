package main

import (
	"os"

	"route256/cart/internal/app"
	"route256/cart/internal/app/handler"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/cart/service"
	product "route256/cart/internal/product/service"
	"route256/cart/pkg/lib/logger"
)

func main() {
	log := logger.SetupLogger(logger.Local)

	cartRepository := repository.NewInMemoryCartRepository()
	cartService := service.NewCartService(cartRepository)
	productService := product.NewProductService()
	cartHandler := handler.NewCartHandler(cartService, productService, log)

	routes := app.InitRoutes(cartHandler)
	server := app.NewApp(routes, log)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
