package main

import (
	"os"

	"route256/cart/internal/app"
	"route256/cart/internal/app/handler"
	"route256/cart/internal/cart/repository"
	"route256/cart/internal/cart/service"
	"route256/cart/internal/config"
	product "route256/cart/internal/product/service"
	"route256/cart/pkg/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	cartRepository := repository.NewInMemoryCartRepository()
	cartService := service.NewCartService(cartRepository)
	productService := product.NewProductService()
	cartHandler := handler.NewCartHandler(cartService, productService, log)

	routes := app.InitRoutes(cartHandler)
	server := app.NewApp(routes, log)
	err := server.ListenAndServe(cfg.Port)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
